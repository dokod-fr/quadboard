package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDC struct {
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	secretKey    []byte
	secure       bool
}

type Session struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
	Email    string   `json:"email"`
}

func NewOIDC(ctx context.Context, issuer, clientID, clientSecret, baseURL, secretKey string, secure bool) (*OIDC, error) {
	provider, err := createProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	redirectURL, err := GetOIDCRedirectURL(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC redirectURL: %w", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"}, // rip off "groups"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	return &OIDC{
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		secretKey:    []byte(secretKey),
		secure:       secure,
	}, nil
}

func GetOIDCRedirectURL(baseURL string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base_url: %w", err)
	}

	u = u.JoinPath("/auth/callback")

	return u.String(), nil
}

func createProvider(ctx context.Context, issuer string) (*oidc.Provider, error) {
	const (
		maxRetries = 10              // TODO: Create variable for this constant
		delay      = 2 * time.Second // TODO: Create variable for this constant
	)

	var lastErr error

	for i := 1; i <= maxRetries; i++ {
		provider, err := oidc.NewProvider(ctx, issuer)
		if err == nil {
			if i > 1 {
				slog.Info("OIDC provider available",
					"issuer", issuer,
					"attempt", i,
				)
			}
			return provider, nil
		}

		lastErr = err

		slog.Warn("OIDC provider not yet available",
			"issuer", issuer,
			"attempt", i,
			"max_attempts", maxRetries,
			"error", err,
		)

		time.Sleep(delay)
	}

	return nil, fmt.Errorf("Attempts %w failed after %d", lastErr, maxRetries)
}

// generateState random state string for OIDC login to prevent CSRF attacks. It returns a base64 encoded string.
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (o *OIDC) LoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// Generate verifier PKCE
	verifier := oauth2.GenerateVerifier()

	// Keep it temporarily in a cookie ...
	http.SetCookie(w, &http.Cookie{
		Name:     "quadboard_pkce_verifier",
		Value:    verifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   o.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})

	// ... and state for CSRF protection
	http.SetCookie(w, &http.Cookie{
		Name:     "quadboard_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   o.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})

	// Create the auth URL with the state and PKCE challenge
	url := o.oauth2Config.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(verifier),
	)

	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OIDC) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Validate the state parameter to prevent CSRF attacks
	stateCookie, err := r.Cookie("quadboard_state")
	if err != nil {
		http.Error(w, "State cookie not found", http.StatusBadRequest)
		return
	}
	stateParam := r.URL.Query().Get("state")
	if stateParam != stateCookie.Value {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "quadboard_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Get the PKCE verifier from the cookie
	pkceCookie, err := r.Cookie("quadboard_pkce_verifier")
	if err != nil {
		http.Error(w, "PKCE verifier cookie not found", http.StatusBadRequest)
		return
	}
	verifier := pkceCookie.Value

	http.SetCookie(w, &http.Cookie{
		Name:   "quadboard_pkce_verifier",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Exchange the code for a token
	code := r.URL.Query().Get("code")
	oauth2Token, err := o.oauth2Config.Exchange(
		r.Context(),
		code,
		oauth2.VerifierOption(verifier),
	)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	if errCode := r.URL.Query().Get("error"); errCode != "" {
		desc := r.URL.Query().Get("error_description")
		slog.Error("OIDC authentication failed",
			"error", errCode,
			"description", desc,
		)
		http.Error(w, desc, http.StatusUnauthorized)
		return
	}

	// Extract the ID Token from the OAuth2 token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token", http.StatusInternalServerError)
		return
	}

	idToken, err := o.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token", http.StatusInternalServerError)
		return
	}

	var claims struct {
		PreferredUsername string   `json:"preferred_username"`
		Groups            []string `json:"groups"`
		Email             string   `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
		return
	}
	slog.Debug("claims struct", slog.Any("claims", claims))

	session := Session{
		Username: claims.PreferredUsername,
		Groups:   claims.Groups,
		Email:    claims.Email,
	}

	if err := o.setSessionCookie(w, session); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (o *OIDC) GetSessionFromRequest(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("quadboard_session")
	if err != nil {
		return nil, err
	}
	return o.decodeSession(cookie.Value)
}

func (o *OIDC) setSessionCookie(w http.ResponseWriter, session Session) error {
	encoded, err := o.encodeSession(session)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "quadboard_session",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   o.secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	return nil
}

func (o *OIDC) encodeSession(session Session) (string, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	encodedData := base64.StdEncoding.EncodeToString(data)

	mac := hmac.New(sha256.New, o.secretKey)
	mac.Write([]byte(encodedData))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s.%s", encodedData, signature), nil
}

func (o *OIDC) decodeSession(s string) (*Session, error) {
	// Look for the last '.' to separate the encoded data from the signature
	idx := strings.LastIndex(s, ".")
	if idx == -1 {
		return nil, errors.New("invalid session format")
	}

	encodedData := s[:idx]
	signature := s[idx+1:]

	mac := hmac.New(sha256.New, o.secretKey)
	mac.Write([]byte(encodedData))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, errors.New("invalid session signature")
	}

	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}
