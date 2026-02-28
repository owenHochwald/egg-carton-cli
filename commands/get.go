package commands

import (
	"fmt"

	"github.com/owenHochwald/egg-carton/cli/api"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/spf13/cobra"
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Retrieve a secret",
	Long:  `Decrypt and retrieve a secret from your EggCarton vault.`,
	Args:  cobra.MaximumNArgs(1), // 0 or 1 args - if no key, list all
	RunE:  runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("configuration error — the binary may not be configured correctly")
	}

	tokens, err := ensureValidToken(cfg)
	if err != nil {
		return err
	}

	owner, err := cfg.GetOwner()
	if err != nil {
		return fmt.Errorf("could not identify your account — try running 'egg login' again")
	}

	client := api.NewClient(cfg.GetAPIBaseURL(), tokens.AccessToken)

	eggs, err := client.GetEgg(owner)
	if err != nil {
		return fmt.Errorf("failed to retrieve secrets: %w", err)
	}

	if len(args) == 1 {
		key := args[0]
		found := false
		for _, egg := range eggs {
			if egg.SecretID == key {
				fmt.Printf("Secret: %s\n", key)
				fmt.Printf("Value: %s\n", egg.Plaintext)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("secret %q not found in your vault", key)
		}
	} else {
		if len(eggs) == 0 {
			fmt.Println("No secrets found in your vault.")
			return nil
		}

		fmt.Printf("Found %d secret(s):\n\n", len(eggs))
		for _, egg := range eggs {
			fmt.Printf("Key: %s\n", egg.SecretID)
			fmt.Printf("Value: %s\n", egg.Plaintext)
			fmt.Printf("Created: %s\n", egg.CreatedAt)
			fmt.Println("---")
		}
	}

	return nil
}
