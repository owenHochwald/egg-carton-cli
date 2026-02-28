package commands

import (
	"fmt"

	"github.com/owenHochwald/egg-carton/cli/api"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/spf13/cobra"
)

// BreakCmd represents the break command
var BreakCmd = &cobra.Command{
	Use:   "break [key]",
	Short: "Delete a secret",
	Long:  `Permanently delete a secret from your EggCarton vault.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runBreak,
}

func runBreak(cmd *cobra.Command, args []string) error {
	key := args[0]

	fmt.Printf("Breaking egg: %s\n", key)

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

	if err := client.BreakEgg(owner, key); err != nil {
		return fmt.Errorf("failed to delete secret %q: %w", key, err)
	}

	fmt.Printf("Successfully deleted secret: %s\n", key)

	return nil
}
