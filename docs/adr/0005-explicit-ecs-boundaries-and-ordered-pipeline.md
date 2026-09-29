# ADR 0005: Use explicit ECS component boundaries and an ordered simulation pipeline

- **Status:** Accepted
- **Date:** 2026-09-29

## Context

As the simulation grew (movement, shooting, collision, camera, map bounds, lighting, audio intents, and proximity/engine volume), implicit dependencies and ad-hoc world state made execution order and data ownership harder to reason about.

## Decision

Define canonical component types in `game/ecs/components` and reference them directly across systems. Keep world attachment generic through `ecs.AttachComponent[T]` and provide explicit world query/setter methods for each component slice. Represent map bounds and audio behavior as ECS data where possible. Execute gameplay through one default ordered pipeline in `game/sim/pipeline.go`, with each system declaring read/write dependency metadata.

## Consequences

System order is visible, testable, and reviewable as gameplay behavior. Dependency declarations act as lightweight contracts for future refactors. The tradeoff is more explicit boilerplate (component APIs and query structs), but it preserves determinism, keeps boundaries clear, and avoids hidden coupling between render/input/audio adapters and core simulation rules.
