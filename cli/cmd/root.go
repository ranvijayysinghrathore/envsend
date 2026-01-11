package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	serverURL string
	verbose   bool

	// DefaultServerURL is the default server URL.
	// It can be overridden at build time using -ldflags.
	DefaultServerURL = "http://localhost:8080"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "envsend",
	Short: "Secure, zero-knowledge secret transfer CLI",
	Long: `EnvSend is a production-grade CLI tool for securely transferring .env files
and secrets between developers without ever exposing plaintext to the server.

All encryption happens client-side using AES-256-GCM with Argon2id key derivation.
Secrets are ephemeral by default and auto-expire or self-destruct after access.

Examples:
  envsend .env                                    # Send with defaults (10m, 1 view)
  envsend .env --expires 1h --max-views 3        # Custom expiry and views
  envsend .env --require-passphrase              # Passphrase-protected
  envsend .env --ssh github:username             # Encrypt for specific recipient
  cat .env | envsend                             # Pipe support
  envsend "https://..."                          # Receive and decrypt (Smart Mode)
  envsend receive <link> > .env                  # Legacy receive command`,
	Version: "1.0.0",
	// If arguments are provided and it's not a subcommand, treat it as a file to send
	// OR a URL to receive (Smart Mode)
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for piped input
		stat, _ := os.Stdin.Stat()
		hasPipe := (stat.Mode() & os.ModeCharDevice) == 0

		if len(args) > 0 {
			arg := args[0]
			// Check if argument looks like a URL or relative path
			if strings.HasPrefix(arg, "http:") || strings.HasPrefix(arg, "https:") || strings.HasPrefix(arg, "/s/") {
				return runReceive(cmd, args)
			}
			// Otherwise treat as file to send
			return runSend(cmd, args)
		} else if hasPipe {
			// If no args but piped input, send it
			return runSend(cmd, args)
		}
		return cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", DefaultServerURL, "EnvSend server URL")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Register send flags on root command for default action
	rootCmd.Flags().StringVar(&expiresIn, "expires", "10m", "Expiration time (e.g., 10m, 1h, 24h)")
	rootCmd.Flags().IntVar(&maxViews, "max-views", 1, "Maximum number of views before destruction")
	rootCmd.Flags().BoolVar(&requirePassphrase, "require-passphrase", false, "Require passphrase for decryption")
	rootCmd.Flags().BoolVar(&ipLock, "ip-lock", false, "Lock secret to sender's IP address")
	rootCmd.Flags().StringVar(&sshRecipient, "ssh", "", "Encrypt for specific recipient (github:username or gitlab:username)")
	rootCmd.Flags().IntVar(&shamirThreshold, "shamir-threshold", 0, "Shamir secret sharing threshold (advanced)")
	rootCmd.Flags().IntVar(&shamirShares, "shamir-shares", 0, "Shamir secret sharing total shares (advanced)")
}
