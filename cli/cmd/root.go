package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	serverURL string
	verbose   bool
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
  envreceive <link> > .env                       # Receive and decrypt`,
	Version: "1.0.0",
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
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "EnvSend server URL")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
