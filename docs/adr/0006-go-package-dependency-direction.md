# ADR 0006: Enforce one-way dependency direction for Go packages

- **Status:** Accepted
- **Date:** 2026-09-29

## Context

As features expanded, package boundaries became easier to blur (core rules, adapters, runtime composition, and asset loaders). Without explicit direction, import cycles and framework leakage into core gameplay become likely.

## Decision

Use one-way dependency flow:

- `app/tankism` is entrypoint-only and depends on `game`.
- `game` (runtime composition) may depend on core packages (`game/ecs`, `game/sim`, `game/worldseed`), adapters (`game/input`, `game/render`), `lib/*`, and `resources`.
- `game/ecs`, `game/sim`, and `game/worldseed` must stay framework-independent (no Ebitengine imports).
- Adapters (`game/input`, `game/render`) may depend on Ebitengine and core data, but core packages must not depend on adapters.
- `lib/*` stays reusable and must not depend on `game/*`.
- `resources` exposes embedded asset paths/files only.

In short: dependencies may point downward to core/shared packages, never upward into composition or entrypoint.

## Consequences

This keeps the core deterministic and test-friendly, prevents import cycles, and makes architecture intent obvious during review. Some boundary types/interfaces are needed to avoid back-edges, but that cost is preferred over hidden coupling.
