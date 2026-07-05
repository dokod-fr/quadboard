# Resource Model

A Resource represents any service, application or endpoint exposed through QuadBoard.

Resources are the primary building blocks of the dashboard. Every provider (Quadlet, Traefik, static configuration, etc.) exposes its data by producing one or more `Resource` instances.

A Resource is independent from its presentation and from its origin.

---

## Core fields

Every Resource provides the following information:

- ID
- Name
- Description
- Group
- Tags
- Health
- Actions

Future versions may introduce additional metadata without changing the core model.

Examples include:

- Icon
- URL
- Labels
- Owner
- Visibility
- Metrics

---

## Actions

Actions represent entry points associated with the Resource.

Examples:

- Open
- Admin
- Documentation
- Repository
- Metrics
- Logs

A Resource may expose zero, one or many actions.

---

## Rendering

A Resource may be rendered using different presentation modes.

- Compact
- Standard
- Detailed

Rendering never changes the Resource itself.

The same Resource instance should be displayable using any supported layout.

---

## Providers

Resources are produced by providers.

Examples include:

- Quadlet
- Traefik
- Static configuration

Providers should never generate HTML.
They only expose Resources.

---

## Design Principles

The Resource model should remain:

- Generic
- Stable
- Provider-agnostic
- Presentation-agnostic

The dashboard is responsible for rendering Resources.
Providers are responsible for discovering them.