package auth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// ProviderConfig holds OAuth configs for all providers.
type ProviderConfig struct {
	GitHub *oauth2.Config
	Google *oauth2.Config
}

// LoadProviders creates OAuth configs from environment variables.
func LoadProviders(backendURL string) *ProviderConfig {
	p := &ProviderConfig{}

	// GitHub OAuth
	ghID := os.Getenv("GITHUB_CLIENT_ID")
	ghSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if ghID != "" && ghSecret != "" {
		p.GitHub = &oauth2.Config{
			ClientID:     ghID,
			ClientSecret: ghSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  backendURL + "/auth/github/callback",
			Scopes:       []string{"user:email", "read:user"},
		}
	}

	// Google OAuth
	goID := os.Getenv("GOOGLE_CLIENT_ID")
	goSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if goID != "" && goSecret != "" {
		p.Google = &oauth2.Config{
			ClientID:     goID,
			ClientSecret: goSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  backendURL + "/auth/google/callback",
			Scopes:       []string{"openid", "profile", "email"},
		}
	}

	return p
}
