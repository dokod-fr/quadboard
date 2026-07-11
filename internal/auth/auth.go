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
	"net/http"
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
}

func NewOIDC(ctx context.Context, issuer, clientID, clientSecret, redirectURL, secretKey string, secure bool) (*OIDC, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "groups"},
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

	// 1. On stocke le state dans un cookie temporaire pour pouvoir le vérifier au callback
	http.SetCookie(w, &http.Cookie{
		Name:     "quadboard_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Passe à false pour les tests locaux en HTTP
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // Expire après 5 minutes (suffisant pour le flux de login)
	})

	// 2. On envoie l'utilisateur vers le fournisseur avec ce state
	url := o.oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OIDC) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get state from the cookie
	stateCookie, err := r.Cookie("quadboard_state")
	if err != nil {
		http.Error(w, "State cookie not found", http.StatusBadRequest)
		return
	}

	// 2. Get state from provider callback
	stateParam := r.URL.Query().Get("state")

	// 3. Check state concordance (CSRF protection)
	if stateParam != stateCookie.Value {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	// 4. Clean up the state cookie after validation
	http.SetCookie(w, &http.Cookie{
		Name:   "quadboard_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Supprime le cookie
	})

	code := r.URL.Query().Get("code")
	oauth2Token, err := o.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

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
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
		return
	}

	session := Session{
		Username: claims.PreferredUsername,
		Groups:   claims.Groups,
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
		Secure:   true, // S'assurer que l'app est servie en HTTPS
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
