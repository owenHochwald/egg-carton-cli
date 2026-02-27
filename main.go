package main

import (
	"fmt"
	"os"

	"github.com/owenHochwald/egg-carton/cli/commands"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/spf13/cobra"
)

// Build-time version metadata — set via ldflags
var (
	Version   string
	BuildTime string
	GitCommit string
)

// Compiled-in config values — set via ldflags at build time.
// Zero value is intentional: LoadConfig falls back to env vars when these are empty (local dev).
var (
	compiledAPIEndpoint   string
	compiledUserPoolID    string
	compiledClientID      string
	compiledCognitoDomain string
	compiledRegion        string
)

var rootCmd = &cobra.Command{
	Use:   "egg",
	Short: "🥚 EggCarton - Secure secret management CLI",
	Long: `EggCarton is a secure CLI tool for managing secrets with an egg theme!

Commands:
  🔐 login           - Authenticate with OAuth
  🐔 lay (add)       - Store a secret (lay an egg)
  🥚 get             - Retrieve secrets from your vault
  🐣 hatch (run)     - Inject secrets and run a command (hatch your eggs)
  💥 break           - Delete a secret from your vault

It uses AWS Lambda, DynamoDB, and KMS for encryption,
with Cognito authentication via OAuth PKCE flow.`,
}

func main() {
	// Wire compiled-in config values into the config package before any command runs
	config.CompiledAPIEndpoint = compiledAPIEndpoint
	config.CompiledUserPoolID = compiledUserPoolID
	config.CompiledClientID = compiledClientID
	config.CompiledCognitoDomain = compiledCognitoDomain
	config.CompiledRegion = compiledRegion

	// Add all subcommands
	rootCmd.AddCommand(commands.LoginCmd)
	rootCmd.AddCommand(commands.AddCmd)
	rootCmd.AddCommand(commands.GetCmd)
	rootCmd.AddCommand(commands.BreakCmd)
	rootCmd.AddCommand(commands.RunCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
