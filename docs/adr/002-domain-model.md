# ADR 0002 — Resource Domain Model

- Status: Accepted
- Date: 2026-07-02

## Context

QuadBoard aims to provide a single entry point to the services running on a server.

Initially, the project focused on "applications". During the design phase, it became clear that not every item displayed on the dashboard is necessarily an application.

Examples include:

- A web application (Grafana)
- An identity provider (Authelia)
- A reverse proxy (Traefik)
- A documentation website
- An external link
- A database administration interface
- Future integrations

Using the term "Application" would unnecessarily restrict the domain model.

## Decision

The central concept of QuadBoard is the **Resource**.

A Resource represents anything that can be presented to the user on the dashboard.

Providers are responsible for discovering resources.

The dashboard is responsible for displaying resources.

The dashboard must not know how a resource was discovered.

## Consequences

Providers expose Resources.

Examples of future providers include:

- Quadlet
- Traefik
- Static configuration
- Kubernetes
- Docker

Each provider produces the same domain objects.

The UI only consumes Resources.

This separation allows new providers to be added without changing the presentation layer.

## Future evolution

Resources may later support:

- Groups
- Tags
- Health information
- Multiple actions
- Metadata
- Provider-specific information

without changing the overall architecture.