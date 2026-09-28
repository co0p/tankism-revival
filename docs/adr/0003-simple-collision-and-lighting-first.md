# ADR 0003: Begin with simple collision and lighting implementations

- **Status:** Accepted
- **Date:** 2026-09-28

## Context

Collision and lighting are learning goals, but sophisticated physics, spatial partitioning, dynamic shadows, and shader pipelines would delay a playable combat loop.

## Decision

Use circles for mobile tank/projectile colliders and axis-aligned rectangles for static arena solids. Use simple candidate loops, explicit overlap tests, and separate resolution. Render lighting as ambient illumination plus soft local light masks; particles are short-lived data rendered in a separate effects pass.

## Consequences

The initial game prioritizes clarity and inspectability over physical realism. Add spatial hashes, pooling, shader-based lighting, shadowing, or complex shapes only after a profile or playtest demonstrates a specific limitation; record the change in a new ADR.
