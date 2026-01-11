package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/envsend/cli/client"
	"github.com/yourusername/envsend/cli/crypto"
	"github.com/yourusername/envsend/cli/utils"
)

var (
	// Send command flags
	expiresIn         string
	maxViews          int
	requirePassphrase bool
	ipLock            bool
	sshRecipient      string
	shamirThreshold   int
	shamirShares      int
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send [file]",
	Short: "Send an encrypted secret",
	Long: `Encrypt and send a secret file to the EnvSend server.
The file is encrypted locally before upload. The server never sees plaintext.

Examples:
  envsend send .env
  envsend send .env --expires 30m --max-views 2
  envsend send .env --require-passphrase
  envsend send .env --ssh github:username
  cat .env | envsend send`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSend,
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVar(&expiresIn, "expires", "10m", "Expiration time (e.g., 10m, 1h, 24h)")
	sendCmd.Flags().IntVar(&maxViews, "max-views", 1, "Maximum number of views before destruction")
	sendCmd.Flags().BoolVar(&requirePassphrase, "require-passphrase", false, "Require passphrase for decryption")
	sendCmd.Flags().BoolVar(&ipLock, "ip-lock", false, "Lock secret to sender's IP address")
	sendCmd.Flags().StringVar(&sshRecipient, "ssh", "", "Encrypt for specific recipient (github:username or gitlab:username)")
	sendCmd.Flags().IntVar(&shamirThreshold, "shamir-threshold", 0, "Shamir secret sharing threshold (advanced)")
	sendCmd.Flags().IntVar(&shamirShares, "shamir-shares", 0, "Shamir secret sharing total shares (advanced)")
}

func runSend(cmd *cobra.Command, args []string) error {
	// Read file or stdin
	var filePath string
	if len(args) > 0 {
		filePath = args[0]
	}

	plaintext, err := utils.ReadFileOrStdin(filePath)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	defer crypto.ZeroBytes(plaintext)

	if verbose {
		fmt.Fprintf(os.Stderr, "Read %d bytes from input\n", len(plaintext))
	}

	// Determine encryption method
	var encryptedBlob string
	var metadata crypto.EncryptionMetadata
	var encryptionKey []byte

	if shamirThreshold > 0 && shamirShares > 0 {
		// Shamir Secret Sharing mode
		return fmt.Errorf("Shamir mode not yet implemented in this example")
	} else if sshRecipient != "" {
		// SSH-based encryption
		return runSSHEncryption(plaintext)
	} else if requirePassphrase {
		// Passphrase-based encryption
		passphrase, err := utils.PromptPassword("Enter passphrase: ")
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}
		defer crypto.ZeroString(&passphrase)

		confirmPassphrase, err := utils.PromptPassword("Confirm passphrase: ")
		if err != nil {
			return fmt.Errorf("failed to read passphrase confirmation: %w", err)
		}
		defer crypto.ZeroString(&confirmPassphrase)

		if passphrase != confirmPassphrase {
			return fmt.Errorf("passphrases do not match")
		}

		encryptedBlob, metadata, err = crypto.EncryptWithPassphrase(plaintext, passphrase)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
	} else {
		// Default: random key encryption
		encryptionKey, err = crypto.GenerateKey()
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}
		defer crypto.ZeroBytes(encryptionKey)

		encryptedBlob, metadata, err = crypto.EncryptAES256GCM(plaintext, encryptionKey)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Encrypted data successfully\n")
	}

	// Convert metadata to map for JSON
	metadataMap := map[string]interface{}{
		"algorithm":      metadata.Algorithm,
		"iv":             metadata.IV,
		"keyDerivation":  metadata.KeyDerivation,
		"version":        metadata.Version,
	}
	if metadata.Salt != "" {
		metadataMap["salt"] = metadata.Salt
	}

	// Upload to server
	apiClient := client.NewAPIClient(serverURL)

	uploadReq := client.UploadSecretRequest{
		EncryptedBlob:      encryptedBlob,
		EncryptionMetadata: metadataMap,
		ExpiresIn:          expiresIn,
		MaxViews:           maxViews,
	}

	if ipLock {
		// TODO: Get current IP address
		uploadReq.IPLock = ""
	}

	resp, err := apiClient.UploadSecret(uploadReq)
	if err != nil {
		return fmt.Errorf("failed to upload secret: %w", err)
	}

	// Display results
	fmt.Println("\n✅ Secret uploaded successfully!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if requirePassphrase {
		fmt.Printf("🔗 Share this link: %s\n", resp.URL)
		fmt.Println("🔑 Recipient will need the passphrase you set")
	} else {
		// Encode key in URL
		encodedKey := crypto.EncodeKey(encryptionKey)
		fullURL := fmt.Sprintf("%s#%s", resp.URL, encodedKey)
		fmt.Printf("🔗 Share this link: %s\n", fullURL)
		fmt.Println("⚠️  The key is in the URL fragment (after #) - it's never sent to the server")
	}

	fmt.Printf("⏰ Expires: %s (%s)\n", resp.ExpiresAt.Format(time.RFC3339), expiresIn)
	fmt.Printf("👁️  Max views: %d\n", resp.MaxViews)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\n⚠️  WARNING: This link will self-destruct after use!")

	return nil
}

func runSSHEncryption(plaintext []byte) error {
	fmt.Fprintf(os.Stderr, "Fetching SSH public key for %s...\n", sshRecipient)

	// Fetch recipient's SSH public key
	sshKey, err := client.FetchSSHKey(sshRecipient)
	if err != nil {
		return fmt.Errorf("failed to fetch SSH key: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Retrieved SSH key: %s...\n", sshKey[:50])
	}

	// Parse SSH key to X25519 format
	recipientPublicKey, err := crypto.ParseSSHPublicKey(sshKey)
	if err != nil {
		return fmt.Errorf("failed to parse SSH key: %w", err)
	}

	// Generate symmetric key
	symmetricKey, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}
	defer crypto.ZeroBytes(symmetricKey)

	// Encrypt data with symmetric key
	encryptedBlob, metadata, err := crypto.EncryptAES256GCM(plaintext, symmetricKey)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Encrypt symmetric key for recipient
	encryptedKey, ephemeralPubKey, err := crypto.EncryptSymmetricKeyForRecipient(symmetricKey, recipientPublicKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt key for recipient: %w", err)
	}

	// Add encrypted key info to metadata
	metadataMap := map[string]interface{}{
		"algorithm":       metadata.Algorithm,
		"iv":              metadata.IV,
		"keyDerivation":   "x25519",
		"version":         metadata.Version,
		"encryptedKey":    crypto.EncodeKey(encryptedKey),
		"ephemeralPubKey": crypto.EncodePublicKey(ephemeralPubKey),
	}

	// Upload to server
	apiClient := client.NewAPIClient(serverURL)

	uploadReq := client.UploadSecretRequest{
		EncryptedBlob:      encryptedBlob,
		EncryptionMetadata: metadataMap,
		ExpiresIn:          expiresIn,
		MaxViews:           maxViews,
		RecipientID:        sshRecipient,
	}

	resp, err := apiClient.UploadSecret(uploadReq)
	if err != nil {
		return fmt.Errorf("failed to upload secret: %w", err)
	}

	// Display results
	fmt.Println("\n✅ Secret uploaded successfully!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🔗 Share this link: %s\n", resp.URL)
	fmt.Printf("🔑 Encrypted for: %s\n", sshRecipient)
	fmt.Printf("⏰ Expires: %s (%s)\n", resp.ExpiresAt.Format(time.RFC3339), expiresIn)
	fmt.Printf("👁️  Max views: %d\n", resp.MaxViews)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}
