# Contributing to QuadBoard

First of all, thank you for considering contributing to QuadBoard!

The goal of this project is to provide a lightweight, zero-configuration application portal for Podman Quadlets.

## Development philosophy

When contributing, please keep the following principles in mind:

* Keep things simple.
* Prefer convention over configuration.
* Avoid unnecessary dependencies.
* Document architectural decisions.
* Write readable code before clever code.
* Keep commits focused and atomic.

## Getting started

Clone the repository and enter the development environment:

```bash
git clone <repository>
cd quadboard
nix develop
```

Available development tasks:

```bash
task build
task run
task fmt
task lint
task test
```

## Pull Requests

Before opening a pull request:

* Make sure the project builds.
* Run formatting and linting.
* Add or update tests when appropriate.
* Update documentation if needed.

## Commit messages

This project follows the Conventional Commits specification.

Examples:

* `feat: add quadlet provider`
* `fix: handle missing labels`
* `docs: update architecture`
* `build: add Task task`
* `refactor: simplify provider interface`

## Discussions

For larger features or architectural changes, please open an issue or start a discussion before implementing them.

This helps keep the project consistent and avoids duplicated work.

Thank you for contributing!
