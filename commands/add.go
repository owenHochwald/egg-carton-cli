package commands

import (
	"fmt"

	"github.com/owenHochwald/egg-carton/cli/api"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/spf13/cobra"
)

// AddCmd represents the lay command (alias: add)
var AddCmd = &cobra.Command{
	Use:     "lay [key] [value]",
	Aliases: []string{"add"},
	Short:   "Store a secret (lay an egg)",
	Long:    `Encrypt and store a secret in your EggCarton vault.`,
	Args:    cobra.ExactArgs(2),
	RunE:    runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	fmt.Printf("Laying egg: %s\n", key)

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

	if err := client.PutEgg(owner, key, value); err != nil {
		return fmt.Errorf("failed to store secret %q: %w", key, err)
	}

	fmt.Printf("Successfully stored secret: %s\n", key)

	return nil
}
