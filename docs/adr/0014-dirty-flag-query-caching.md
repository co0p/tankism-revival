# ADR 0014 — Dirty Flag: Query Result Caching

**Status:** Proposed  
**Pattern source:** [gameprogrammingpatterns.com/dirty-flag.html](https://gameprogrammingpatterns.com/dirty-flag.html)

---

## Context

Every method in `world_queries.go` (`Renderables()`, `TankMovers()`,
`ArmedTanks()`, etc.) rebuilds and re-sorts its result slice on **every
call**. `Renderables()` sorts by ZIndex then Entity every frame even when
nothing has changed.

In the current codebase many query methods are called multiple times per
frame (once from the simulation pipeline, once from the render system).
Each call iterates all relevant maps and allocates a new slice.

## Decision

Add a structural dirty flag to `World`:

```go
// in World
structurallyDirty bool

func (w *World) markDirty()      { w.structurallyDirty = true }
func (w *World) ClearDirty()     { w.structurallyDirty = false }
func (w *World) IsDirty() bool   { return w.structurallyDirty }
```

`markDirty()` is called inside `CreateEntity`, `DestroyEntity`, and
`AttachComponent`.

Each query method caches its last result:

```go
// in World
cachedRenderables    []Renderable
cachedTankMovers     []TankMover
// … one per query

func (w *World) Renderables() []Renderable {
    if w.structurallyDirty || w.cachedRenderables == nil {
        w.cachedRenderables = w.buildRenderables()
    }
    return w.cachedRenderables
}
```

`RunDefaultPipeline` calls `world.ClearDirty()` at the start of each tick
(after draining queues, before running systems). Any system that creates or
destroys entities will re-mark dirty, causing one rebuild for that tick.

## Consequences

**Good:**
- Query rebuild + sort cost is paid **at most once per frame**, and zero
  times in steady state (no entity creation/destruction).
- Multiple callers of the same query method within one frame share the cached
  result at no extra cost.

**Bad / risks:**
- Stale caches if a system reads a query after another system modifies
  structure in the same tick. `ClearDirty` at tick start + `markDirty` on
  mutation means the second system correctly sees dirty and rebuilds.
- Cache slices must not be mutated by callers. Return copies or document the
  read-only contract clearly.

## Alternatives considered

**Per-query dirty flags.** More granular but more state to manage; a single
flag is sufficient given the pipeline's linear order.

**Rebuild once at the start of each tick unconditionally.** Simpler but
wastes work in steady state; the dirty flag is trivially cheap.
