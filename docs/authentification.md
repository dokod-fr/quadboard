# QuadBoard Authentication

QuadBoard supports robust authentication through a hybrid approach: it can either transparently trust identity headers forwarded by a secure reverse proxy (like Authelia/Authentik via Traefik), or perform a standard OpenID Connect (OIDC) authorization flow.

## Authentication Modes

1. Reverse Proxy Header Auth (Recommended for Self-Hosting)

When deployed behind a reverse proxy that handles central authentication (e.g., Traefik with an Authelia Forward Auth middleware), QuadBoard can automatically sign in users using secure HTTP headers.

If the proxy validates the user, it injects identity headers that QuadBoard trusts immediately. This bypasses the need for OIDC client secrets and eliminates double-redirection loops.
2. Autonomous OIDC (Fallback)

If no trusted proxy headers are detected, and the oidc configuration is present, QuadBoard acts as an autonomous OIDC client. It will initiate the authorization code flow, exchange the token, and manage its own secure session.

- **Optional**: If neither proxy headers nor OIDC configuration are present, QuadBoard runs without authentication. All resources are visible to everyone.
- **Group Extraction**: Whether using headers or OIDC, QuadBoard extracts the user's `groups` (from the headers or the OIDC ID Token). These groups are used to filter resources based on Quadlet labels.

- **Stateless Sessions** (for OIDC): For autonomous OIDC, QuadBoard does not use a database. Upon successful login, a session cookie is created containing the username and groups, signed with an HMAC (using the `secret_key` from the configuration) to prevent tampering.

## Configuration

### A. Setup for Reverse Proxy Header Auth (e.g., Authelia + Traefik)

To use this zero-friction mode, configure your reverse proxy middleware to pass the authenticated user's identity via headers.

For Traefik + Authelia, ensure your Forward Auth middleware configuration forwards the required response headers:

```yaml
# Traefik Dynamic Configuration
http:
  middlewares:
    authelia:
      forwardAuth:
        address: "http://authelia:9091/api/verify?rd=https://auth.dokod.fr/"
        trustForwardHeader: true
        authResponseHeaders:
          - "Remote-User"   # Maps to Username
          - "Remote-Groups" # Maps to Groups (comma-separated)
```

In QuadBoard, no further OIDC setup is needed. It will automatically detect these headers and authenticate the session.

### B. Setup for Autonomous OIDC (Authelia / Authentik / Keycloak)

If you are not using forward auth headers, you must register QuadBoard as a client in your Identity Provider (IdP):

- **Client ID**: Your choice (e.g., quadboard)

- **Client Secret**: Your choice (use a strong secret)

- **Redirect URI**: https://<your-quadboard-domain>/auth/callback

- **Scopes**: openid, profile, groups (ensure your provider issues the groups claim in the ID token).

Provide these details to QuadBoard via config.yaml or environment variables:
```yaml
# config.yaml
auth:
  secret_key: "a-very-strong-random-key-to-sign-cookies"
  oidc:
    issuer: "https://auth.dokod.fr"
    client_id: "quadboard"
    client_secret: "your-client-secret"
```

### Local Development & Testing

Testing Header Auth (No OIDC required)

You can easily mock the Authelia/Traefik header integration in development using curl:
Bash

curl -H "Remote-User: liam" -H "Remote-Groups: admins,users" http://localhost:8080/

Testing Autonomous OIDC

**Important**: The OIDC session cookie is set with the `Secure: true` flag, meaning it will only be transmitted over HTTPS.

If you are testing the OIDC flow locally over plain HTTP (`http://localhost:8080`), your browser will reject the cookie, resulting in an infinite redirect loop between `/` and `/login`.

To test the OIDC flow locally:

1. Use a reverse proxy with self-signed certificates (e.g., mkcert or Caddy) in front of QuadBoard.

2. Or temporarily set Secure: false in the setSessionCookie function inside internal/auth/auth.go during local testing. Do not commit this change.