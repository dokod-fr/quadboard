package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// startMockOIDC crée un faux serveur OIDC qui répond à la découverte
func startMockOIDC(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	// Endpoint de découverte
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"issuer":                 "http://" + r.Host,
			"authorization_endpoint": "http://" + r.Host + "/auth",
			"token_endpoint":         "http://" + r.Host + "/token",
			"jwks_uri":               "http://" + r.Host + "/jwks",
		})
	})

	// Réponse fictive pour les clés publiques (JWKS)
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"keys":[]}`)) // Vide, car on ne va pas vraiment vérifier de vrai JWT dans ce test d'init
	})

	return httptest.NewServer(mux)
}

func TestNewOIDC(t *testing.T) {
	mockServer := startMockOIDC(t)
	defer mockServer.Close()

	// On teste juste que l'initialisation (découverte) fonctionne
	o, err := NewOIDC(
		context.Background(),
		mockServer.URL, // On passe l'URL de notre faux serveur
		"test-client",
		"test-secret",
		"http://localhost/callback",
		"test-secret-key",
	)

	if err != nil {
		t.Fatalf("NewOIDC() failed to initialize with mock server: %v", err)
	}

	if o == nil {
		t.Fatal("Expected OIDC instance, got nil")
	}
}

func TestEncodeDecodeSession(t *testing.T) {
	o := &OIDC{
		secretKey: []byte("test-secret-key"),
	}

	originalSession := Session{
		Username: "alice",
		Groups:   []string{"admins", "users"},
	}

	encoded, err := o.encodeSession(originalSession)
	if err != nil {
		t.Fatalf("encodeSession failed: %v", err)
	}

	decoded, err := o.decodeSession(encoded)
	if err != nil {
		t.Fatalf("decodeSession failed: %v", err)
	}

	if decoded.Username != originalSession.Username {
		t.Errorf("Expected username %s, got %s", originalSession.Username, decoded.Username)
	}

	if len(decoded.Groups) != len(originalSession.Groups) {
		t.Errorf("Expected %d groups, got %d", len(originalSession.Groups), len(decoded.Groups))
	}
}

func TestDecodeInvalidSignature(t *testing.T) {
	o := &OIDC{
		secretKey: []byte("test-secret-key"),
	}

	// On essaie de décoder un cookie signé avec une autre clé
	tamperedCookie := "eyJ1c2VybmFtZSI6ImFsaWNlIn0=.invalidSignature"
	_, err := o.decodeSession(tamperedCookie)
	if err == nil {
		t.Fatal("Expected error when decoding invalid signature, got nil")
	}
}
