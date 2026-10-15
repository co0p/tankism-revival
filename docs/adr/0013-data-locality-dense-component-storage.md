# ADR 0013 — Data Locality: Dense Component Storage (Deferred)

**Status:** Deferred — do not implement until profiling confirms a bottleneck  
**Pattern source:** [gameprogrammingpatterns.com/data-locality.html](https://gameprogrammingpatterns.com/data-locality.html)

---

## Context

`World` uses one `map[Entity]T` per component type (20 maps). Go maps store
values in hash buckets on the heap. Iterating a query like `TankMovers()`
touches 6 separate maps, each with randomly addressed heap nodes. As entity
count grows, CPU cache misses on hot-path components (`Transform`,
`TankVelocity`, `TankDrive`) will become measurable.

At the current MVP entity count (< 50 entities) this is **not** a
bottleneck.

## Decision

**Defer this change.** When a profiling run (`go tool pprof`) shows that
component map iteration is in the top 5% of CPU time, switch the highest-
churn component groups to dense slice storage:

- Replace `map[Entity]Transform` with a `[]Transform` indexed by a stable
  slot allocator (free-list of indices, matching ADR 0015 entity recycling).
- Keep `map` storage for sparse components (`HitFlash`, `PlaySoundIntent`,
  `PointLight`) that few entities hold at any given time.

The query methods in `world_queries.go` become tight loops over contiguous
memory, which the CPU prefetcher can stride through efficiently.

## Consequences (when implemented)

**Good:**
- Hot-path iteration is cache-friendly.
- No per-element map overhead (hash computation, pointer indirection).

**Bad / risks:**
- Slot allocator adds complexity; must coordinate with the object pool
  (ADR 0015).
- `AttachComponent` can no longer use a simple type switch; it needs to know
  whether a kind uses dense or sparse storage.
- World internals change significantly — all tests must pass.

## Trigger condition

Implement when `go tool pprof` on a real play session shows component
iteration in the hot path **and** entity count exceeds 200.
