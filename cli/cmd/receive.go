package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ranvijayysinghrathore/envsend/cli/client"
	"github.com/ranvijayysinghrathore/envsend/cli/crypto"
	"github.com/ranvijayysinghrathore/envsend/cli/utils"
)

var (
	// Receive command flags
	outputFile string
	passphrase string
)

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive <url>",
	Short: "Receive and decrypt a secret",
	Long: `Download and decrypt a secret from the EnvSend server.
Decryption happens locally - the server never sees your decryption key.

Examples:
  envsend receive https://envsend.io/s/abc123#key
  envsend receive https://envsend.io/s/abc123 --passphrase
  envreceive https://envsend.io/s/abc123#key > .env`,
	Aliases: []string{"get", "download"},
	Args:    cobra.ExactArgs(1),
	RunE:    runReceive,
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	receiveCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	receiveCmd.Flags().StringVar(&passphrase, "passphrase", "", "Passphrase for decryption (will prompt if not provided)")
}

func runReceive(cmd *cobra.Command, args []string) error {
	secretURL := args[0]

	// Parse URL to extract secret ID and key
	parsedURL, err := url.Parse(secretURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Extract secret ID from path
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) == 0 {
		return fmt.Errorf("invalid URL: missing secret ID")
	}
	secretID := pathParts[len(pathParts)-1]

	// Extract encryption key from fragment (after #)
	var encryptionKey []byte
	fragment := parsedURL.Fragment
	if fragment != "" {
		encryptionKey, err = crypto.DecodeKey(fragment)
		if err != nil {
			return fmt.Errorf("invalid encryption key in URL: %w", err)
		}
		defer crypto.ZeroBytes(encryptionKey)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Downloading secret: %s\n", secretID)
	}

	// Download encrypted secret
	apiClient := client.NewAPIClient(serverURL)
	resp, err := apiClient.DownloadSecret(secretID)
	if err != nil {
		return fmt.Errorf("failed to download secret: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Downloaded encrypted data (%d bytes)\n", len(resp.EncryptedBlob))
		fmt.Fprintf(os.Stderr, "Views remaining: %d\n", resp.ViewsRemaining)
	}

	// Parse encryption metadata
	metadataJSON, err := json.Marshal(resp.EncryptionMetadata)
	if err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	metadata, err := crypto.MetadataFromJSON(string(metadataJSON))
	if err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Decrypt based on key derivation method
	var plaintext []byte

	switch metadata.KeyDerivation {
	case "argon2id":
		// Passphrase-based decryption
		if passphrase == "" {
			passphrase, err = utils.PromptPassword("Enter passphrase: ")
			if err != nil {
				return fmt.Errorf("failed to read passphrase: %w", err)
			}
		}
		defer crypto.ZeroString(&passphrase)

		plaintext, err = crypto.DecryptWithPassphrase(resp.EncryptedBlob, passphrase, metadata)
		if err != nil {
			return fmt.Errorf("decryption failed (wrong passphrase?): %w", err)
		}
		return handlePlaintextOutput(plaintext, outputFile, resp.ViewsRemaining)

	case "x25519":
		// SSH-based decryption
		return fmt.Errorf("SSH-based decryption not yet implemented")

	case "shamir":
		// Shamir Secret Sharing decryption
		fmt.Printf("☢️  SHAMIR SECRET SHARING DETECTED\n")

		var shares []crypto.ShamirShare

		// Process initial share from URL if present
		if fragment != "" && strings.HasPrefix(fragment, "s") {
			parts := strings.SplitN(fragment[1:], "-", 2)
			if len(parts) == 2 {
				var index int
				fmt.Sscanf(parts[0], "%d", &index)
				shares = append(shares, crypto.ShamirShare{
					Index: index,
					Data:  parts[1],
				})
				fmt.Printf("✓ Share %d collected from URL\n", index)
			}
		}

		// Prompt for additional shares until we can decrypt
		for {
			if len(shares) >= 2 { // Min threshold for Shamir is usually 2
				// Try combining and decrypting
				plaintext, err = crypto.DecryptWithShamir(resp.EncryptedBlob, metadata, shares)
				if err == nil {
					// Success!
					return handlePlaintextOutput(plaintext, outputFile, resp.ViewsRemaining)
				}
				if verbose {
					fmt.Fprintf(os.Stderr, "Not enough shares yet (or invalid shares): %v\n", err)
				}
			}
			fmt.Printf("Need more shares. Collected: %d\n", len(shares))
			input, err := utils.PromptUser("Enter next share key (part after #): ")
			if err != nil {
				return fmt.Errorf("failed to read share: %w", err)
			}
			input = strings.TrimSpace(input)
			if input == "" {
				return fmt.Errorf("decryption aborted: not enough shares")
			}

			// Clean input (in case they pasted the whole URL)
			if idx := strings.LastIndex(input, "#"); idx != -1 {
				input = input[idx+1:]
			}

			if !strings.HasPrefix(input, "s") {
				fmt.Println("❌ Invalid share format. Expected 's1-base64...'")
				continue
			}

			parts := strings.SplitN(input[1:], "-", 2)
			if len(parts) == 2 {
				var index int
				fmt.Sscanf(parts[0], "%d", &index)
				
				// Check for duplicates
				duplicate := false
				for _, s := range shares {
					if s.Index == index {
						duplicate = true
						break
					}
				}
				if duplicate {
					fmt.Printf("⚠️  Share %d already collected. Please provide a different share.\n", index)
					continue
				}

				shares = append(shares, crypto.ShamirShare{
					Index: index,
					Data:  parts[1],
				})
				fmt.Printf("✓ Share %d collected\n", index)
			} else {
				fmt.Println("❌ Invalid share format.")
			}
		}

	case "none":
		// Key-based decryption
		if encryptionKey == nil {
			return fmt.Errorf("encryption key required but not found in URL")
		}

		plaintext, err = crypto.DecryptAES256GCM(resp.EncryptedBlob, encryptionKey, metadata)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
		return handlePlaintextOutput(plaintext, outputFile, resp.ViewsRemaining)

	default:
		return fmt.Errorf("unsupported key derivation method: %s", metadata.KeyDerivation)
	}
}

func handlePlaintextOutput(plaintext []byte, outputFile string, viewsRemaining int) error {
	defer crypto.ZeroBytes(plaintext)

	// Output plaintext
	if outputFile != "" {
		if err := utils.WriteFile(outputFile, plaintext, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "✅ Secret saved to: %s\n", outputFile)
	} else {
		if err := utils.WriteToStdout(plaintext); err != nil {
			return fmt.Errorf("failed to write to stdout: %w", err)
		}
		fmt.Println() // Add newline for console output
	}

	// Warn about remaining views
	if viewsRemaining > 0 {
		fmt.Fprintf(os.Stderr, "⚠️  Warning: %d view(s) remaining before destruction\n", viewsRemaining)
	} else {
		fmt.Fprintf(os.Stderr, "🔥 Secret has been destroyed (max views reached)\n")
	}

	return nil
}
