# QuadBoard Authentication

QuadBoard supports optional OpenID Connect (OIDC) authentication. This is particularly useful when hosting QuadBoard behind a reverse proxy like Traefik or Caddy, alongside an identity provider like Authelia or Authentik.

## Behavior

- **Optional:** If the `oidc` block is not configured in `config.yaml` or via environment variables, QuadBoard runs without authentication. All resources are visible to everyone.
- **Group Extraction:** When a user logs in, QuadBoard extracts the `groups` claim from the ID Token. These groups will be used in future versions to filter resources based on Quadlet labels.
- **Stateless Sessions:** QuadBoard does not use a database or an in-memory session store. Upon successful login, a session cookie is created. This cookie contains the username and groups, and is signed with an HMAC (using the `secret_key` from the configuration) to prevent tampering.

## Configuration

To enable authentication, you must provide an OIDC configuration and a secret key to sign the cookies.

### Provider Configuration (Authelia / Authentik / Keycloak)

You need to create a new Client/Application in your Identity Provider with the following settings:
- **Client ID:** Your choice (e.g., `quadboard`)
- **Client Secret:** Your choice (use a strong secret)
- **Redirect URI:** `https://<your-quadboard-domain>/auth/callback`
- **Scopes:** `openid`, `profile`, `groups` (ensure your provider issues the `groups` claim in the ID token).

### Local Development

**Important:** The session cookie is set with the `Secure: true` flag, meaning it will only be transmitted over HTTPS. 

If you are testing locally over plain HTTP (`http://localhost:8080`), your browser will reject the cookie, resulting in an infinite redirect loop between `/` and `/login`.

To test authentication locally, you have two options:
1. Use a reverse proxy with self-signed certificates (e.g., `mkcert`) in front of QuadBoard.
2. Temporarily set `Secure: false` in the `setSessionCookie` function inside `internal/auth/auth.go` during your local testing. **Do not commit this change.**