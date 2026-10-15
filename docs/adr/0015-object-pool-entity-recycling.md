# ADR 0015 — Object Pool: Entity Recycling

**Status:** Proposed  
**Pattern source:** [gameprogrammingpatterns.com/object-pool.html](https://gameprogrammingpatterns.com/object-pool.html)

---

## Context

Projectiles are the highest-churn entities in Tankism. `weaponSpawnSystem`
creates a new entity (inserts into several component maps) every time a tank
fires. `lifetimeSystem` and `projectileImpactsSystem` destroy entities
(delete from all 20 component maps) on expiry or impact.

At a fire rate of 2 shots/second with 10 tanks firing, this is ~20
create/destroy cycles per second. Each `DestroyEntity` call performs 20 map
deletions; each spawn performs 5–7 map insertions. This generates steady GC
pressure from map bucket churn.

## Decision

Add a free-list inside `World`:

```go
// in World
freeList []Entity

func (w *World) CreateEntity() Entity {
    if len(w.freeList) > 0 {
        id := w.freeList[len(w.freeList)-1]
        w.freeList = w.freeList[:len(w.freeList)-1]
        w.alive[id] = struct{}{}
        return id
    }
    // existing path: increment nextEntityID
    …
}

func (w *World) DestroyEntity(entity Entity) bool {
    if !w.IsAlive(entity) { return false }
    delete(w.alive, entity)
    // remove from component maps as before
    w.freeList = append(w.freeList, entity)
    return true
}
```

Recycled entity IDs are reused for new entities. Components are set by
the spawn call, so stale data from the previous occupant is always
overwritten before the entity is observable.

Add an observability helper:

```go
func (w *World) LiveEntityCount() int  { return len(w.alive) }
func (w *World) PooledEntityCount() int { return len(w.freeList) }
```

## Consequences

**Good:**
- Entity create/destroy becomes allocation-free after the pool warms up.
- GC pressure from map churn is eliminated for recycled IDs.
- Matches the Object Pool pattern intent: pre-allocated slots, reused
  without heap allocation.

**Bad / risks:**
- Entity IDs are reused. Any system that holds a stale `Entity` value from a
  previous occupant would see a live (but different) entity. Existing code
  uses `IsAlive` checks; this must remain the contract.
- The free-list itself is a slice that grows. In the worst case (all entities
  destroyed at once) it holds O(n) entries. Acceptable.
- Coordination with ADR 0013 (dense storage): if dense slices are adopted,
  the free-list must also track slot indices, not just entity IDs.

## Alternatives considered

**Pre-allocated ring buffer for projectiles only.** More targeted but
requires special-casing in `World`. A general free-list is simpler and
benefits all short-lived entities (sparks, hit flashes, tread marks).
