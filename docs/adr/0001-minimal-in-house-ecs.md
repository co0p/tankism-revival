# ADR 0001: Use a minimal in-house ECS for the MVP

- **Status:** Accepted
- **Date:** 2026-09-28

## Context

The project is intended to teach ECS concepts while building one small game. Third-party ECS libraries and highly optimized archetype storage would obscure mechanics and add APIs before performance evidence exists.

## Decision

Implement the smallest Tankism-specific ECS required by active slices. Entities are opaque IDs, components are plain data held in explicit stores, and systems query the component combinations they need. The world and ordered systems are composed at the application boundary.

## Consequences

The first implementation may use simple maps or slices and linear queries. It must support deterministic tests and safe entity cleanup. Future profiling may justify an internal storage redesign, but a dependency or generic extraction requires a new ADR and a demonstrated need.
