# QuadBoard

> Zero-config application portal for Podman Quadlets.

QuadBoard is a lightweight application portal that automatically discovers services deployed with Podman Quadlets and presents them in a clean, responsive web interface.

The goal is simple:

> **Deploy an application, and it appears automatically.**

No database. No manual dashboard configuration. No duplicated metadata.

## Features (planned)

- Automatic discovery of Podman Quadlets
- Automatic Traefik URL detection
- Theme support
- Reverse proxy authentication (Authelia, Authentik, OAuth2 Proxy…)
- Responsive interface
- Application search
- Metadata through Quadlet labels
- Single static Go binary
- First-class Nix support

## Project Status

🚧 Early development.

The architecture is being designed before the first implementation.

## Philosophy

QuadBoard follows a few simple principles:

- Convention over configuration
- Zero-config by default
- Single binary
- No database
- Extensible architecture
- Fast startup and low resource usage

## Documentation

- [Architecture](ARCHITECTURE.md)

## Roadmap

v0.1.0

✓ CLI
✓ Config
✓ HTTP
✓ Dashboard
✓ Resource model
✓ Quadlet provider
✓ Traefik provider
□ Groups
□ Authelia Openid Connect
□ Theme engine

v0.2.0

□ Static provider
□ Search
□ Tags
□ User preferences

v0.3.0

□ Authentication hooks
□ Plugins
□ API

## License

[GPLv3](LICENSE.md)
