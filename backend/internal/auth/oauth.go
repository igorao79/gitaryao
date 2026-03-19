package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitserv/internal/models"
)

// OAuthHandler handles OAuth login flows.
type OAuthHandler struct {
	Providers   *ProviderConfig
	JWT         *JWTManager
	Users       *models.UserStore
	FrontendURL string
	states      map[string]bool // simple in-memory state store (use Redis in production)
}

// NewOAuthHandler creates a new OAuth handler.
func NewOAuthHandler(providers *ProviderConfig, jwtMgr *JWTManager, users *models.UserStore, frontendURL string) *OAuthHandler {
	return &OAuthHandler{
		Providers:   providers,
		JWT:         jwtMgr,
		Users:       users,
		FrontendURL: frontendURL,
		states:      make(map[string]bool),
	}
}

func (h *OAuthHandler) generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := hex.EncodeToString(b)
	h.states[state] = true
	return state
}

func (h *OAuthHandler) validateState(state string) bool {
	if h.states[state] {
		delete(h.states, state)
		return true
	}
	return false
}

// redirectWithToken redirects to the frontend with the JWT token.
func (h *OAuthHandler) redirectWithToken(w http.ResponseWriter, r *http.Request, token string) {
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", h.FrontendURL, url.QueryEscape(token))
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// redirectWithError redirects to the frontend with an error message.
func (h *OAuthHandler) redirectWithError(w http.ResponseWriter, r *http.Request, errMsg string) {
	redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.FrontendURL, url.QueryEscape(errMsg))
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// --- GitHub OAuth ---

// GithubLogin redirects to GitHub OAuth authorize page.
func (h *OAuthHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	if h.Providers.GitHub == nil {
		http.Error(w, "GitHub OAuth not configured", http.StatusNotImplemented)
		return
	}
	state := h.generateState()
	url := h.Providers.GitHub.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GithubCallback handles the OAuth callback from GitHub.
func (h *OAuthHandler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	if h.Providers.GitHub == nil {
		http.Error(w, "GitHub OAuth not configured", http.StatusNotImplemented)
		return
	}

	state := r.URL.Query().Get("state")
	if !h.validateState(state) {
		h.redirectWithError(w, r, "Invalid state parameter")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		h.redirectWithError(w, r, "No code provided")
		return
	}

	// Exchange code for token
	token, err := h.Providers.GitHub.Exchange(r.Context(), code)
	if err != nil {
		h.redirectWithError(w, r, "Failed to exchange code")
		return
	}

	// Get user info from GitHub API
	ghUser, err := h.getGithubUser(token.AccessToken)
	if err != nil {
		h.redirectWithError(w, r, "Failed to get user info")
		return
	}

	// Get user email if not public
	email := ghUser.Email
	if email == "" {
		email, _ = h.getGithubEmail(token.AccessToken)
	}
	if email == "" {
		email = fmt.Sprintf("%s@users.gitserv.local", ghUser.Login)
	}

	// Upsert user in database
	user, err := h.Users.UpsertGithubUser(ghUser.ID, ghUser.Login, email, ghUser.AvatarURL)
	if err != nil {
		h.redirectWithError(w, r, "Failed to save user")
		return
	}

	// Generate JWT
	jwtToken, err := h.JWT.Generate(user.ID, user.Username, user.Email)
	if err != nil {
		h.redirectWithError(w, r, "Failed to generate token")
		return
	}

	h.redirectWithToken(w, r, jwtToken)
}

type githubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (h *OAuthHandler) getGithubUser(accessToken string) (*githubUser, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user githubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (h *OAuthHandler) getGithubEmail(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []githubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}
	return "", fmt.Errorf("no verified email found")
}

// --- Google OAuth ---

// GoogleLogin redirects to Google OAuth authorize page.
func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.Providers.Google == nil {
		http.Error(w, "Google OAuth not configured", http.StatusNotImplemented)
		return
	}
	state := h.generateState()
	url := h.Providers.Google.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the OAuth callback from Google.
func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if h.Providers.Google == nil {
		http.Error(w, "Google OAuth not configured", http.StatusNotImplemented)
		return
	}

	state := r.URL.Query().Get("state")
	if !h.validateState(state) {
		h.redirectWithError(w, r, "Invalid state parameter")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		h.redirectWithError(w, r, "No code provided")
		return
	}

	token, err := h.Providers.Google.Exchange(r.Context(), code)
	if err != nil {
		h.redirectWithError(w, r, "Failed to exchange code")
		return
	}

	// Get user info from Google
	goUser, err := h.getGoogleUser(token.AccessToken)
	if err != nil {
		h.redirectWithError(w, r, "Failed to get user info")
		return
	}

	// Generate a username from email (part before @)
	username := goUser.Email
	if at := len(goUser.Email); at > 0 {
		for i, c := range goUser.Email {
			if c == '@' {
				username = goUser.Email[:i]
				break
			}
		}
	}

	user, err := h.Users.UpsertGoogleUser(goUser.Sub, username, goUser.Email, goUser.Picture)
	if err != nil {
		h.redirectWithError(w, r, "Failed to save user")
		return
	}

	jwtToken, err := h.JWT.Generate(user.ID, user.Username, user.Email)
	if err != nil {
		h.redirectWithError(w, r, "Failed to generate token")
		return
	}

	h.redirectWithToken(w, r, jwtToken)
}

type googleUser struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (h *OAuthHandler) getGoogleUser(accessToken string) (*googleUser, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var user googleUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// --- User info endpoint ---

// CurrentUser returns the currently authenticated user.
func (h *OAuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	claims := GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.Users.FindByID(claims.UserID)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
