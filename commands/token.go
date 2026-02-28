package commands

import (
	"fmt"

	"github.com/owenHochwald/egg-carton/cli/auth"
	"github.com/owenHochwald/egg-carton/cli/config"
)

// ensureValidToken loads tokens from disk, refreshes them if expired, saves the
// refreshed tokens, and returns a valid TokenData. If the user is not logged in
// or refresh fails, it returns a user-friendly error.
func ensureValidToken(cfg *config.Config) (*config.TokenData, error) {
	tokens, err := cfg.LoadTokens()
	if err != nil {
		return nil, fmt.Errorf("you are not logged in — run 'egg login' to authenticate")
	}
	if tokens == nil {
		return nil, fmt.Errorf("you are not logged in — run 'egg login' to authenticate")
	}
	if tokens.IsTokenValid() {
		return tokens, nil
	}
	fmt.Println("Session expired, refreshing...")
	newTokens, err := auth.RefreshAccessToken(cfg.GetTokenURL(), cfg.CognitoConfig.ClientID, tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("session refresh failed — run 'egg login' to re-authenticate")
	}
	if err := cfg.SaveTokens(newTokens); err != nil {
		return nil, fmt.Errorf("failed to save refreshed session: %w", err)
	}
	return newTokens, nil
}
