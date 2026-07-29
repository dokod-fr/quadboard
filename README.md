<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="internal/http/view/assets/img/quadboard-bw.svg">
    <source media="(prefers-color-scheme: light)" srcset="internal/http/view/assets/img/quadboard-color.svg">
    <img alt="QuadBoard Logo" src="internal/http/view/assets/img/quadboard-color.svg" width="200">
  </picture>
</p>

<h1 align="center">QuadBoard</h1>

<p align="center">
  Zero-config application portal for Podman Quadlets.
</p>

# QuadBoard

> Zero-config application portal for Podman Quadlets.

QuadBoard is a lightweight application portal that automatically discovers services deployed with Podman Quadlets and presents them in a clean, responsive web interface.

The goal is simple:

> **Deploy an application, and it appears automatically.**

No database. No manual dashboard configuration. No duplicated metadata.

## Features

* Automatic discovery of Podman Quadlets
* Automatic service metadata detection (Traefik routing, descritions, icons)
* Responsive web interface (dark/light mode, search, grouped cards)
* OIDC authentication support (Authelia, Keycloak, ...)
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

Environment=QUADBOARD_BASE_URL=https://auth.example.com

# OIDC config
Environment=QUADBOARD_AUTH_OIDC_ISSUER=https://auth.example.com/realms/homelab
Environment=QUADBOARD_AUTH_OIDC_CLIENT_ID=quadboard

Secret=oidc_client_secret,type=env,target=QUADBOARD_AUTH_OIDC_CLIENT_SECRET
Secret=session_secret,type=env,target=QUADBOARD_AUTH_SECRET_KEY

[Install]
WantedBy=default.target
```

## Configuration

QuadBoard is designed to work out of the box with zero configuration. However, you can customize its behavior using environment variables or a YAML configuration file.

### Server Configuration

By default, QuadBoard listens on 0.0.0.0:8080 and automatically discovers Podman Quadlets in the standard directories (/etc/containers/systemd/ and ~/.config/containers/systemd/).

You can override server settings (like address, timeouts, logging, or authentication) using:

* YAML configuration file
* Environment variables
* Built-in defaults

See: [Configuration](docs/configuration.md)

### Service Discovery & Labels

For the Quadlet provider, QuadBoard automatically detects services and extracts metadata (like URLs from Traefik routing rules or descriptions from standard container labels) to create dashboard cards.

You can help QuadBoard accurately represent your applications by using specific quadboard.* labels in your .container or .pod files to define custom names, icons, logos, descriptions, and groups.

Note: By default, QuadBoard hides its own container in the dashboard. You can disable this behavior by setting show_itself: true in your configuration.

See: [Labels Documentation](docs/labels.md)

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
