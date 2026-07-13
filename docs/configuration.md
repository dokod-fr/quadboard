# QuadBoard Configuration

QuadBoard offers a flexible configuration system. You can configure it using a YAML file, environment variables, or a combination of both.

## Priority Order

Configuration is loaded using the following priority (from highest to lowest):

  - `Environment Variables`: Override all other settings.
  - `YAML File`: Overrides default constants.
  - `Default Constants`: Built-in sensible defaults.

## Configuration File

By default, QuadBoard will look for a config.yaml file located next to the QuadBoard executable. 

You can specify a custom path for the configuration file by setting the QUADBOARD_CONFIG_FILE environment variable.
Example config.yaml

```yml
server:
  address: "0.0.0.0:8080"
  read_timeout: 5
  write_timeout: 10
logging:
  level: "info"
  format: "text"
providers:
  quadlet:
    paths:
      - /etc/containers/systemd/
      - /opt/quadboard/quadlets
# Required if OIDC is enabled  
auth:
  secret_key: "your-strong-hmac-secret-key"
  secure: true # Set to false for local HTTP testing
  oidc:
    issuer: "https://auth.example.com"
    client_id: "quadboard"
    client_secret: "super-secret"
    redirect_url: "https://quadboard.example.com/auth/callback"
```

## Environment Variables

All environment variables start with the QUADBOARD_ prefix.

### Server Configuration

| Variable                      | Description               | Default                |
|-------------------------------|---------------------------|------------------------|	
| QUADBOARD_CONFIG_FILE	        | Path to a custom YAML configuration file.        |	(Checks for config.yaml next to the binary) |
| QUADBOARD_SERVER_ADDRESS	    | The address and port the HTTP server listens on. |	0.0.0.0:8080 |
| QUADBOARD_SERVER_READ_TIMEOUT	| HTTP read timeout in seconds.                    |	5  |
| QUADBOARD_SERVER_WRITE_TIMEOUT|	HTTP write timeout in seconds.                   | 10  |
  
### Logging Configuration

| Variable                      | Description               | Default                |
|-------------------------------|---------------------------|------------------------|
| QUADBOARD_LOGGING_LEVEL	      | Log level (e.g., debug, info, warn, error).	| info |
| QUADBOARD_LOGGING_FORMAT	    | Log format (e.g., text, json).	            | text |
  
### Providers Configuration

| Variable                      | Description               | Default                |
|-------------------------------|---------------------------|------------------------|	
| QUADBOARD_QUADLET_PATHS	      | Comma-separated list of directories to scan for Quadlet files.	| /etc/containers/systemd/, ~/.config/containers/systemd/ |
  

> Note: When using QUADBOARD_QUADLET_PATHS, it completely replaces the default paths and the YAML configuration paths.

### Authentication Configuration

| Variable                           | Description               | Default                |
|------------------------------------|---------------------------|------------------------|	
| QUADBOARD_AUTH_SECRET_KEY	         | Secret key used to sign the session cookie (HMAC).	| (Empty) |
| QUADBOARD_AUTH_OIDC_ISSUER	       | The issuer URL of your OIDC provider (e.g., Authelia, Authentik). | (Empty - Auth disabled) |
| QUADBOARD_AUTH_OIDC_CLIENT_ID	     | The Client ID configured in your OIDC provider.	| (Empty) |
| QUADBOARD_AUTH_OIDC_CLIENT_SECRET	 | The Client Secret configured in your OIDC provider. |	(Empty) |
| QUADBOARD_AUTH_OIDC_REDIRECT_URL	 | The callback URL QuadBoard will use (must match the OIDC provider config). |	(Empty) |
| QUADBOARD_AUTH_SECURE | Set to false to allow session cookies over plain HTTP (for local dev only). | true |


> Note: If QUADBOARD_AUTH_OIDC_ISSUER is not set, authentication is completely disabled and all resources are visible.

## Docker Configuration

When running QuadBoard inside a container, configuration should generally be provided from outside the image.

The recommended approach is to mount a configuration file:

```bash
docker run \
  -v ./config.yaml:/etc/quadboard/config.yaml:ro \
  -e QUADBOARD_CONFIG_FILE=/etc/quadboard/config.yaml \
  ghcr.io/dokod-fr/quadboard:latest
```

The container image does not contain environment-specific configuration.

This allows the same image to be reused with different:

* OIDC providers
* Quadlet locations
* logging settings
* deployment environments

For production deployments, prefer mounting the configuration file as read-only.
