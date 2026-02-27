package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the CLI configuration
type Config struct {
	APIEndpoint   string        `json:"api_endpoint"`
	CognitoConfig CognitoConfig `json:"cognito"`
	TokenPath     string        `json:"-"` // Not serialized
}

// CognitoConfig holds Cognito-specific configuration
type CognitoConfig struct {
	UserPoolID string `json:"user_pool_id"`
	ClientID   string `json:"client_id"`
	Domain     string `json:"domain"`
	Region     string `json:"region"`
}

// TokenData holds the OAuth tokens
type TokenData struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	IssuedAt     int64  `json:"issued_at"` // Unix timestamp when token was received
}

// Set by main.go from ldflags values before LoadConfig is called.
// When all five are non-empty the binary is self-contained; users need no .env file.
var (
	CompiledAPIEndpoint   string
	CompiledUserPoolID    string
	CompiledClientID      string
	CompiledCognitoDomain string
	CompiledRegion        string
)

// LoadConfig returns configuration for the CLI.
// Priority: compiled-in ldflags values (release builds) → env vars / .env file (local dev).
func LoadConfig() (*Config, error) {
	var (
		apiEndpoint string
		userPoolID  string
		clientID    string
		domain      string
		region      string
	)

	// Use compiled-in values when all five are present (standard release binary).
	if CompiledAPIEndpoint != "" && CompiledUserPoolID != "" &&
		CompiledClientID != "" && CompiledCognitoDomain != "" && CompiledRegion != "" {
		apiEndpoint = CompiledAPIEndpoint
		userPoolID = CompiledUserPoolID
		clientID = CompiledClientID
		domain = CompiledCognitoDomain
		region = CompiledRegion
	} else {
		// Fall back to .env file + environment variables (local dev override).
		_ = godotenv.Load() // ignore error — .env is optional
		apiEndpoint = os.Getenv("API_ENDPOINT")
		userPoolID = os.Getenv("COGNITO_USER_POOL_ID")
		clientID = os.Getenv("COGNITO_CLIENT_ID")
		domain = os.Getenv("COGNITO_DOMAIN")
		region = os.Getenv("COGNITO_REGION")
	}

	if apiEndpoint == "" || userPoolID == "" || clientID == "" || domain == "" || region == "" {
		return nil, fmt.Errorf("configuration not available — this binary may not have been built with required config values")
	}

	cfg := &Config{
		APIEndpoint: apiEndpoint,
		CognitoConfig: CognitoConfig{
			UserPoolID: userPoolID,
			ClientID:   clientID,
			Domain:     domain,
			Region:     region,
		},
	}

	// Set token path
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	cfg.TokenPath = filepath.Join(home, ".eggcarton", "credentials.json")

	return cfg, nil
}

// Should save tokens to ~/.eggcarton/credentials.json with 0600 permissions
func (c *Config) SaveTokens(tokens *TokenData) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(c.TokenPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal tokens to JSON
	b, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tokens: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(c.TokenPath, b, 0600); err != nil {
		return fmt.Errorf("failed to write tokens to file: %w", err)
	}

	return nil
}

// Should load tokens from ~/.eggcarton/credentials.json
func (c *Config) LoadTokens() (*TokenData, error) {
	data, err := os.ReadFile(c.TokenPath)
	if err != nil {
		return nil, err
	}

	var tokens TokenData
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tokens: %w", err)
	}

	return &tokens, nil
}

// Should check if access token is still valid (not expired)
func (t *TokenData) IsTokenValid() bool {
	now := time.Now().Unix()
	if now > t.IssuedAt+int64(t.ExpiresIn)-300 { // 5 minute buffer
		return false
	}
	return true
}

// Returns the OAuth redirect URI for the callback server
func (c *Config) GetRedirectURI() string {
	return "http://localhost:8080/callback"
}

// Returns the full authorization URL for Cognito
func (c *Config) GetAuthorizationURL() string {

	return fmt.Sprintf("https://%s/oauth2/authorize", c.CognitoConfig.Domain)
}

// Returns the token exchange endpoint
func (c *Config) GetTokenURL() string {
	return fmt.Sprintf("https://%s/oauth2/token", c.CognitoConfig.Domain)
}

// Returns the API base URL
func (c *Config) GetAPIBaseURL() string {
	return c.APIEndpoint
}

// GetOwner extracts the owner (user ID) from the access token
func (c *Config) GetOwner() (string, error) {
	tokens, err := c.LoadTokens()
	if err != nil {
		return "", fmt.Errorf("failed to load tokens: %w", err)
	}

	// Extract owner (sub claim) from the access token JWT
	return extractOwnerFromToken(tokens.AccessToken)
}

// extractOwnerFromToken decodes JWT and extracts the 'sub' claim (user ID)
func extractOwnerFromToken(accessToken string) (string, error) {
	parts := strings.Split(accessToken, ".")

	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	var claims struct {
		Sub string `json:"sub"`
	}

	if err = json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	if claims.Sub == "" {
		return "", fmt.Errorf("sub claim not found in token")
	}

	return claims.Sub, nil
}
