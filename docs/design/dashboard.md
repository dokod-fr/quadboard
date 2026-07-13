# Dashboard Design

This document describes the user interface concepts of QuadBoard.

It is intentionally descriptive rather than normative.

## Resource

A resource is the primary element displayed on the dashboard.

Resources may represent:

- Applications
- Services
- Administrative interfaces
- External links
- Future provider-specific objects

## Groups

Resources may be organised into groups.

Examples:

- Monitoring
- Infrastructure
- Development
- Storage
- Security

Groups improve readability when many resources are available.

## Display Modes

QuadBoard supports multiple presentation modes for the same Resource.

### Compact

Designed to maximise information density.

Typical use cases:

- Experienced users
- Large dashboards
- Small screens

### Standard

The default presentation.

Displays the most useful information while remaining compact.

### Detailed

Provides additional metadata and actions.

This view may later be implemented as a dedicated page or side panel.

## Resource Card

A Resource is displayed using a card.

The card is responsible only for presentation.

Typical information includes:

- Icon
- Name
- Description
- Health status
- Primary action
- Optional actions
- Tags

Cards should remain lightweight and easy to scan.

## Actions

Resources may expose one or more actions.

Examples:

- Open
- Admin
- Documentation
- Repository

Actions are provided by the Resource itself.

The dashboard does not invent actions.

## Health

Resources may expose a health status.

Initial states include:

- Healthy
- Starting
- Unreachable
- Unknown

Future versions may include additional metrics such as latency or last successful check.