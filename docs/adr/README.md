# Architecture Decision Records

This directory contains the Architecture Decision Records (ADRs) for QuadBoard.

An ADR documents a significant architectural decision made during the development of the project, along with its context and consequences.

The goal is to capture **why** a decision was made, not only **what** was implemented.

## Format

Each ADR follows the same structure:

* **Status**
* **Date**
* **Context**
* **Decision**
* **Consequences**

## Naming

Files are numbered in chronological order.

Example:

```text
0001-project-layout.md
0002-configuration.md
0003-http-router.md
```

The number is permanent and should never be reused.

## Updating a decision

Architecture evolves over time.

Existing ADRs should not be rewritten to reflect new decisions.

Instead, create a new ADR that supersedes or amends the previous one.

## Status values

Typical statuses are:

* **Proposed**
* **Accepted**
* **Superseded**
* **Deprecated**

## Principles

An ADR should:

* explain the context;
* describe the decision;
* document the consequences and trade-offs;
* remain concise and focused on a single topic.

The collection of ADRs provides a historical record of the project's architectural evolution.
