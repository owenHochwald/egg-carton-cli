package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/owenHochwald/egg-carton/cli/api"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/spf13/cobra"
)

// RunCmd represents the hatch command (alias: run)
var RunCmd = &cobra.Command{
	Use:     "hatch -- [command]",
	Aliases: []string{"run"},
	Short:   "Inject secrets and run a command (hatch your eggs)",
	Long: `Fetch all secrets, set them as environment variables, and execute a command.

Example:
  egg hatch -- go run main.go
  egg hatch -- npm start
  egg hatch -- ./my-script.sh`,
	RunE: runRun,
	// DisableFlagParsing allows passing flags to the subprocess
	DisableFlagParsing: true,
}

func runRun(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("failed to load secrets: %w", err)
	}

	secretEnvVars := make(map[string]string)
	for _, egg := range eggs {
		// Convert secret_id to uppercase env var format (e.g., api_key -> API_KEY)
		envVarName := strings.ToUpper(egg.SecretID)
		secretEnvVars[envVarName] = egg.Plaintext
	}

	// Find the "--" separator in args
	dashIndex := -1
	for i, arg := range args {
		if arg == "--" {
			dashIndex = i
			break
		}
	}

	if dashIndex == -1 || dashIndex == len(args)-1 {
		return fmt.Errorf("usage: egg hatch -- <command> [args...]")
	}

	// Extract command and arguments after "--"
	commandArgs := args[dashIndex+1:]
	if len(commandArgs) == 0 {
		return fmt.Errorf("no command specified after '--'")
	}

	commandName := commandArgs[0]
	commandArguments := commandArgs[1:]

	currentEnv := os.Environ()

	mergedEnv := append([]string{}, currentEnv...)
	for key, value := range secretEnvVars {
		mergedEnv = append(mergedEnv, fmt.Sprintf("%s=%s", key, value))
	}

	fmt.Printf("Hatching %d egg(s) into your environment...\n", len(secretEnvVars))
	for key := range secretEnvVars {
		fmt.Printf("   + %s\n", key)
	}
	fmt.Println()

	command := exec.Command(commandName, commandArguments...)
	command.Env = mergedEnv
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}
