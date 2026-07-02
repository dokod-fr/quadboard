# ADR-0001: Project Layout

* **Status:** Accepted
* **Date:** 2026-07-02

## Context

QuadBoard is intended to be a long-lived open source project.

The project should remain easy to navigate, test and extend while keeping the codebase small and idiomatic.

Several architectural choices had to be made before implementing the first features:

* project layout;
* CLI framework;
* logging;
* development environment;
* build tooling.

## Decision

### Project layout

The project follows the standard Go project layout.

```
cmd/
    quadboard/

internal/
    cli/
    config/
    http/
    logging/
    provider/
    theme/
    version/

web/
    assets/
    templates/
```

* `cmd/` contains application entry points.
* `internal/` contains all application-specific code.
* `web/` contains static assets and HTML templates.

The application logic lives in `internal/`.

The executable entry point (`main.go`) should remain as small as possible and is responsible only for bootstrapping the application.

### CLI

The command-line interface is implemented using Cobra.

Reasons:

* mature ecosystem;
* excellent documentation;
* built-in help and shell completion;
* extensible command hierarchy.

### Logging

The project uses the Go standard library `log/slog`.

Reasons:

* no external dependency;
* structured logging;
* configurable handlers;
* modern Go API.

### Build information

Version information is exposed through the `internal/version` package.

Build metadata is injected at compile time using Go linker flags (`-ldflags`).

This information can be reused by the CLI, HTTP endpoints, logs and diagnostics.

### Development environment

The project uses:

* Nix Flakes for reproducible development environments;
* Task as the single task runner.

All common development operations should be available through `task`.

## Consequences

### Advantages

* Clear separation of responsibilities.
* Small and maintainable packages.
* Minimal `main.go`.
* Reproducible development environment.
* Consistent developer experience.
* Easy to extend with additional providers and interfaces.

### Trade-offs

* Slightly more upfront structure.
* Architectural decisions are made early to avoid future refactoring.

The project intentionally favors clarity and long-term maintainability over rapid prototyping.
