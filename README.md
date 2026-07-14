# QuadBoard

> Zero-config application portal for Podman Quadlets.

QuadBoard is a lightweight application portal that automatically discovers services deployed with Podman Quadlets and presents them in a clean, responsive web interface.

The goal is simple:

> **Deploy an application, and it appears automatically.**

No database. No manual dashboard configuration. No duplicated metadata.

## Features

* Automatic discovery of Podman Quadlets
* Automatic service metadata detection
* Responsive web interface
* OIDC authentication support
* Configurable providers architecture
* YAML configuration support
* Environment variable overrides
* Single static Go binary
* Docker image support
* First-class Nix support (planned)

## Project Status

🚧 Early development / beta release.

QuadBoard is functional and actively evolving. The current release includes:

* Podman Quadlet discovery
* Web dashboard
* Configuration management
* OIDC authentication
* Container image distribution

Breaking changes may still occur while the project matures.

## Quick Start

### Binary

Download the latest release and run:

```bash
./quadboard serve
```

By default QuadBoard will:

* listen on `0.0.0.0:8080`
* discover Quadlet files from the classic Quadlet paths `/etc/containers/systemd/` or `~/.config/containers/systemd/`

### Docker

Run QuadBoard with an external configuration file:

```bash
docker run \
  -p 8080:8080 \
  -v ./config.yaml:/etc/quadboard/config.yaml:ro \
  -v ./providers/quadlet/testdata:/etc/containers/systemd:ro \
  -e QUADBOARD_CONFIG_FILE=/etc/quadboard/config.yaml \
  ghcr.io/dokod-fr/quadboard:latest
```

> `/etc/quadboard` is created into Dockerfile to receive `config.yaml` configuration file.

### Quadlet

```ini
[Container]
Image=ghcr.io/dokod-fr/quadboard:latest

PublishPort=8080:8080

# Or the place 
Volume=/etc/containers/systemd:/etc/containers/systemd:ro

# OIDC config
Environment=QUADBOARD_AUTH_OIDC_ISSUER=https://auth.example.com/realms/homelab
Environment=QUADBOARD_AUTH_OIDC_CLIENT_ID=quadboard
Environment=QUADBOARD_AUTH_OIDC_REDIRECT_URL=https://auth.example.com/callback

Secret=oidc_client_secret,type=env,target=QUADBOARD_AUTH_OIDC_CLIENT_SECRET
Secret=session_secret,type=env,target=QUADBOARD_AUTH_SECRET_KEY

[Install]
WantedBy=default.target
```

## Configuration

QuadBoard supports:

* YAML configuration files
* Environment variables
* Built-in defaults

Configuration priority:

1. Environment variables
2. YAML configuration file
3. Default values

See: [Configuration](docs/configuration.md)

## Authentication

QuadBoard supports OIDC authentication providers such as:

* Keycloak
* Authentik
* Authelia
* OAuth2-compatible identity providers

When authentication is enabled, users are authenticated through the configured OIDC provider.

See:

[Authentication documentation](docs/authentication.md)

## Philosophy

QuadBoard follows a few simple principles:

* Convention over configuration
* Zero-config by default
* Single binary
* No database
* Extensible architecture
* Fast startup and low resource usage

## Documentation

* [Architecture](ARCHITECTURE.md)
* [Configuration](docs/configuration.md)
* [Authentication](docs/authentication.md)

## Roadmap

[Roadmap](docs/roadmap.md)

## License

[GPLv3](LICENSE.md)
