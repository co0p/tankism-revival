# ADR 0016 — Spatial Partition: Collision Broadphase (Deferred)

**Status:** Deferred — do not implement until profiling confirms a bottleneck  
**Pattern source:** [gameprogrammingpatterns.com/spatial-partition.html](https://gameprogrammingpatterns.com/spatial-partition.html)

---

## Context

`obstacleCollisionSystem` and `projectileImpactsSystem` loop over all
projectiles × all tanks and all projectiles × all obstacles to detect
overlaps. This is O(n × m) per frame. At current MVP scale (< 10 projectiles,
< 20 obstacles, < 10 tanks) this is negligible.

A future slice might increase entity counts substantially (e.g. area-of-effect
weapons, many simultaneous projectiles, large procedural maps).

## Decision

**Defer this change.** When profiling shows collision checking in the top
CPU consumers, add a uniform spatial hash:

- Divide the map into a grid of fixed-size cells.
- Cell size = 2 × the largest collider radius in the game.
- Each frame, insert all `ColliderBody` entities into the grid by their
  cell(s).
- Collision pairs are tested only between entities in the same or adjacent
  cells.

The grid is rebuilt from scratch each frame (no incremental update) because
entity movement invalidates positions anyway.

## Consequences (when implemented)

**Good:**
- Collision check count drops from O(n × m) to O(n) amortised for sparse
  distributions.

**Bad / risks:**
- Grid cell size must be tuned to the actual collider radius distribution.
  Too large → many entities per cell, little benefit. Too small → entities
  span many cells, insertion cost rises.
- Grid must live in `game/sim/` or `game/ecs/`, not in `game/render/`.
- Adds meaningful code complexity. Only worthwhile above ~100 collidable
  entities.

## Trigger condition

Implement when `go tool pprof` on a real play session shows collision loop
iteration is a top-5 CPU consumer **and** collidable entity count exceeds 100.
