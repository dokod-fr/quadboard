# QuadBoard Philosophy

QuadBoard is intentionally focused.

It is not an orchestrator.

It is not a deployment platform.

It is not a replacement for existing administration tools.

Its purpose is simple:

> Discover resources, organise them and provide the best possible entry point for users.

## Principles

### Resources are provider-agnostic

The UI never knows where a Resource comes from.

### ProviderRegistry is separate from presentation

Providers discover.

The dashboard presents.

### Simplicity first

Prefer small, understandable components over complex abstractions.

### Explicit configuration

Configuration should remain predictable, reproducible and versionable.

### Lightweight by default

QuadBoard should remain fast, responsive and easy to deploy.

### Themes change appearance, not behaviour

Changing the visual identity must never change the functionality.

### Build on standards

Whenever possible, rely on standard Go libraries and well-established ecosystem projects.