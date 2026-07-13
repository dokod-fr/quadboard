# QuadBoard Architecture

## Vision

QuadBoard is a lightweight dashboard focused on self-hosted services.

The primary goal of v0.1.0 is to discover services from Podman Quadlets and expose them as dashboard resources.

Future versions may support additional providers, we'll see.

---

# Layers

```
cmd/
        │
        ▼
internal/cli
        │
        ▼
internal/http
        │
        ▼
internal/app
        │
        ▼
internal/domain
        ▲
        │
internal/providers
```

Each layer only depends on the layer below.

---

# Domain

The domain contains business concepts only.

```
internal/domain
```

It defines:

- Resource
- Action
- Health
- Provider interface

The domain never depends on infrastructure.

Example:

```go
type Provider interface {
    Resources() ([]Resource, error)
}
```

---

# Application

```
internal/app
```

The application orchestrates providers.

Current service:

```
Discovery
```

Responsibilities:

- call one or more providers
- aggregate resources
- expose them to the HTTP layer

It does **not** know anything about Quadlets.

---

# Providers

Providers convert external systems into domain resources.

Current provider:

```
quadlet
```

Future providers:

```
docker
kubernetes
nomad
...
```

Each provider implements:

```go
type Provider interface {
    Resources() ([]domain.Resource, error)
}
```

---

# Quadlet Provider

Pipeline:

```
Filesystem
      │
      ▼
Load()
      │
      ▼
Model
      │
      ▼
Parse()
      │
      ▼
Build()
      │
      ▼
[]domain.Resource
```

## Load

Responsible for filesystem discovery.

Produces a Model containing:

- Pods
- Containers
- Volumes

Only names and paths are known.

---

## Parse

Responsible for enriching the model.

It parses Quadlet files and extracts metadata.

Examples:

- Description
- PodName
- ContainerName
- Pod
- ...

---

## Build

Responsible for applying QuadBoard business rules.

Current rules:

- one Pod ⇒ one Resource
- one standalone Container ⇒ one Resource
- containers belonging to a Pod enrich the Pod Resource

The Builder knows nothing about the filesystem.

---

# UnitFile

```
internal/providers/quadlet/unitfile
```

This package parses Quadlet/systemd files.

It knows nothing about QuadBoard.

API:

```go
Parse(io.Reader) (*File, error)

ParseFile(path string) (*File, error)
```

Responsibilities:

- parse sections
- decode keys
- expose typed structures

Example:

```go
file.Container.Image
file.Container.Label
file.Pod.Network
```

---

# HTTP

HTTP only renders Resources.

Templates should never assume optional fields exist.

Example:

```go
len(resource.Actions) > 0
```

must be checked before using:

```go
resource.Actions[0]
```

---

# Testing Strategy

## unitfile

Uses:

```
strings.NewReader(...)
```

Tests only parsing logic.

---

## quadlet

Uses:

```
testdata/
```

Tests:

- loader
- parser
- builder

using real Quadlet files.

---

## HTTP

Tests handlers independently from providers.

---

# Current Status (v0.1.0)

Implemented:

- HTTP server
- templ views
- assets
- discovery service
- Quadlet loader
- Quadlet parser
- UnitFile parser
- Builder
- tests

Remaining work:

- configuration
- real provider wiring
- OpenID (Authelia)
- dashboard improvements
- contribution documentation

# Development Guidelines

- Keep packages focused on one responsibility.
- Prefer composition over inheritance.
- Test each layer independently.
- Providers convert external systems into domain resources.
- The domain never depends on infrastructure.
- Mock as little as possible.
- Prefer using real testdata whenever possible.