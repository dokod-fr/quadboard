# QuadBoard Labels

QuadBoard uses a convention-over-configuration approach to gather metadata about your services. It reads standard container labels and specific `quadboard.*` labels from your Podman Quadlet files (`.container` and `.pod`).

## Priority

Labels are evaluated in the following order of priority:

1.  **`quadboard.*` labels**: These are first-class citizens. If present, they always override any other detected value.
2.  **Standard container labels**: If a `quadboard.*` label is not present, QuadBoard will look for standard OCI/Podman labels (e.g., `io.containers.description`).
3.  **Infrastructure detection**: If no label is found, QuadBoard attempts to infer the value (e.g., extracting the URL from Traefik routing rules).

## Supported Labels

You can add these labels to your `.container` or `.pod` files using the `Label=` directive.

### `quadboard.url`
Defines the URL that the dashboard card will link to.
*   **If omitted:** QuadBoard will automatically parse `traefik.http.routers.*.rule` labels. If a `Host(\`example.com\`)` rule is found, the URL will be set to `https://example.com`.
*   **Example:** `Label="quadboard.url=https://my-app.local"`

### `quadboard.description`
A short description of the service displayed on the card.
*   **If omitted:** QuadBoard will look for `io.containers.description` or `org.opencontainers.image.description`.
*   **Example:** `Label="quadboard.description=My favorite home automation tool"`

### `quadboard.icon`
Specifies the name of a [Lucide icon](https://lucide.dev/icons/) to display on the card.
*   **If omitted:** QuadBoard checks for `quadboard.logo`. If neither is found, a default `box` icon is used.
*   **Example:** `Label="quadboard.icon=cloud"` (This will render the cloud icon).

### `quadboard.logo`
Specifies a direct URL to an image (PNG, SVG, etc.) to use as the card's logo.
*   **If omitted:** QuadBoard checks for `quadboard.icon`.
*   **Note:** `quadboard.logo` and `quadboard.icon` are mutually exclusive in the UI, but `quadboard.logo` takes precedence if both are present.
*   **Example:** `Label="quadboard.logo=https://raw.githubusercontent.com/amir20/dozzle/master/assets/logo.svg"`

### `quadboard.group`
Defines the security group and category under which the service will be displayed.

When OIDC authentication is enabled, this value must match the group name declared in your identity provider (e.g., Authelia, Authentik). QuadBoard uses this to grant or restrict access to the service based on the user's groups.

If you want to display a cleaner, human-readable name on the dashboard instead of the raw technical group name, you can map it using the group_labels setting in your UI configuration.

* If omitted: The service is considered public (if auth is enabled) and placed in the "Default" visual group.
* Example: `Label="quadboard.group=security"` (If mapped in config, it will display as "Security & Access" on the dashboard).

## Example Quadlet

Here is a complete example of a `.container` Quadlet file utilizing QuadBoard labels:

```ini
[Unit]
Description=My Custom Application

[Container]
Image=nginx:latest
ContainerName=my-app

# Traefik routing (QuadBoard will automatically detect https://app.local)
Label="traefik.http.routers.my-app.rule=Host(`app.local`)"

# QuadBoard specific overrides
Label="quadboard.description=Production Nginx web server"
Label="quadboard.icon=server"
Label="quadboard.group=Web"

[Install]
WantedBy=multi-user.target
```