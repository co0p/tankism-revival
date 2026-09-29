# ADR 0004: Load runtime resources at the composition root via embedded assets and lib loaders

- **Status:** Accepted
- **Date:** 2026-09-29

## Context

The project now uses multiple runtime asset types (sprites and audio effects/loops). Resource handling spread across gameplay systems would couple pure simulation code to framework-specific loading concerns and make licensing/attribution drift likely.

## Decision

Keep embedded asset paths in `resources/assets.go` and load runtime resources in `game.NewGame()` (composition root). Use focused helpers under `lib/` (for example `lib/sprites` and `lib/audio`) to decode and prepare assets. Gameplay components store stable IDs (like sprite or sound IDs), not framework handles. Every newly added external asset must include attribution and a license file under `resources/`.

## Consequences

Gameplay systems remain testable without filesystem, decoder, or Ebitengine object dependencies. Startup performs explicit fail-fast validation of required resources. This adds some upfront wiring in composition and may duplicate decoded data in memory for convenience, but keeps resource boundaries clear and auditable.
