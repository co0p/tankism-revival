# Refactoring Plan — Game Programming Patterns

This document maps every pattern from [gameprogrammingpatterns.com](https://gameprogrammingpatterns.com/contents.html)
to the Tankism Revival codebase. For each pattern the status is recorded, the
specific problem in our code is described, and a concrete change is proposed.

---

## Legend

| Status | Meaning |
|--------|---------|
| ✅ Already applied | Pattern is correctly in place; no change needed. |
| 🔶 Partial / improvable | Concept exists but the implementation has a known weakness. |
| 🔴 Should adopt | Clear benefit; not yet present. |
| ⬜ Not applicable | Pattern does not fit this project's scope or size. |

---

## Design Patterns Revisited

### Command ✅
**Where:** `components.TankControlIntent` — player input is captured as a
plain data struct (`Throttle`, `Turn`, `Fire`) before any simulation system
reads it. This decouples input from movement and weapon logic.  
**ADR:** no new ADR needed.

---

### Flyweight 🔶
**Problem:** `components.Weapon` has ≈ 20 fields, most of which are
**read-only configuration** (`CooldownSeconds`, `ProjectileSpeed`,
`ProjectileLifetime`, `ProjectileDamage`, accuracy spreads, sprite IDs,
light fields, muzzle-flash fields). Every tank entity carries a full copy of
these values even though they never change per-entity after spawn.

The same applies to `EnemyAI`: `FavoredBehavior` is set at spawn time and
never mutates, yet it lives alongside mutable AI state.

**Change:**
- Split `Weapon` into `WeaponDef` (shared, immutable config — stored once per
  weapon type, keyed by a `WeaponDefID` string) and `WeaponState` (per-entity
  mutable: `CooldownRemaining`, `ShotsFired`).
- Add a `game/defs/` package that holds all `WeaponDef` presets. Systems
  look up the def by ID; they never copy it.
- Mirror the split for `EnemyAI` → `EnemyAIDef` + `EnemyAIState`.

**Benefit:** Entities become lighter. Config changes (tuning damage, speed)
happen in one place rather than at every spawn site.

**ADR:** `docs/adr/0008-flyweight-shared-weapon-and-ai-config.md`

---

### Observer 🔶
**Problem:** `PlayScene.Update` polls `gameInstance` for outcome every frame:
```go
if ps.gameInstance.IsPlayerDestroyed() { … }
if ps.gameInstance.AllEnemiesDestroyed() { … }
```
This is polling, not observation. It also means `game.Game` must expose
predicate methods that couple the scene to internal simulation state.

**Change:**
- Add a small `EventBus` (or typed channel) in `game/events/`. Systems post
  `PlayerDestroyedEvent`, `AllEnemiesDestroyedEvent`, etc.
- `PlayScene` subscribes; removes the polling predicates from `game.Game`.

**Benefit:** Scenes stop reaching into simulation internals. New outcomes
(time limit, draw condition) add a new event type, not a new predicate.

**ADR:** `docs/adr/0009-observer-event-bus.md`

---

### Prototype ⬜
Prototype (cloneable object factories) does not apply — Go does not use
class hierarchies or virtual constructors, and entities are assembled
from component data, not cloned from prototypes.

---

### Singleton ⬜
No global state in the hot path. `resources.FS` is a package-level
`embed.FS`, which is acceptable (immutable, no side effects). No changes
needed.

---

### State ✅
**Where (scene level):** `SceneManager` transitions between
`SceneIDLoading → SceneIDMenu → SceneIDPlay → SceneIDGameOver` using
`SceneResult` values. Each scene is a self-contained state object.

**Where (entity level):** `EnemyBehavior` enum (`Aggressive`, `Cautious`,
`Fleeing`) drives the AI decision system.

**ADR:** ADR 0007 documents the scene-level state machine.

---

## Sequencing Patterns

### Double Buffer 🔶
**Problem:** `SceneManager.Draw` creates a temporary `*ebiten.Image` on
**every frame** during fade transitions to composite the fade overlay. This
allocates and GC-s a GPU-backed image 60 times per second during every
transition.

**Change:** Pre-allocate the fade buffer once in `SceneManager` at
construction time (same logical dimensions). Reuse it each frame.

**Benefit:** Zero per-frame allocations during transitions; no GC pressure
on GPU images.

**ADR:** small enough for an inline code comment; no separate ADR needed.

---

### Game Loop ✅
**Where:** `SceneManager` implements `ebiten.Game` (`Update`, `Draw`,
`Layout`). Ebitengine drives a fixed 60 Hz update tick; `deltaSeconds` is
passed down to scenes and the simulation pipeline.

---

### Update Method ✅
**Where:** `game/sim/pipeline.go` — `defaultPipeline` is an ordered slice of
`runtimeSystem` values, each with a `run(world, deltaSeconds, cfg)` function.
`RunDefaultPipeline` iterates them in declaration order.

**ADR:** ADR 0005 documents this pipeline.

---

## Behavioral Patterns

### Bytecode ⬜
Tankism's AI is data-driven but simple enough that a bytecode VM would add
complexity without benefit. Not applicable at MVP scope.

---

### Subclass Sandbox ⬜
No inheritance hierarchy exists (Go favors composition). Not applicable.

---

### Type Object 🔶
**Problem:** Enemy tank variants (Aggressive / Cautious / Fleeing) are
created by inline construction calls scattered across `worldseed`. Config
values (move speed, weapon cooldown, AI behavior, sprite) are embedded at
each call site. Adding a new variant means finding and copying call sites.

**Change:**
- Add `game/defs/tankdefs.go` with named `TankDef` presets:
  `TankDefAggressive`, `TankDefCautious`, `TankDefFleeing`.
- `TankDef` holds all spawn-time config (speeds, weapon def ID, sprite ID,
  AI behavior).
- `worldseed` calls `defs.SpawnTank(world, def, x, y)` instead of
  assembling components manually.

**Benefit:** All tank tuning lives in one file. New variant = new `TankDef`
constant.

**ADR:** `docs/adr/0010-type-object-tank-definitions.md`

---

## Decoupling Patterns

### Component ✅
**Where:** `game/ecs/` — entities are `uint64` IDs; components are plain
structs in `components.go`; systems in `game/sim/pipeline.go` operate on
queries over the component maps in `World`.

**ADR:** ADR 0001, ADR 0005.

---

### Event Queue 🔶
**Problem:** Sound intents are stored as `map[Entity]PlaySoundIntent` in
`World`. This means **only one sound per entity per frame** can be enqueued.
If a tank is hit and fires in the same tick, one sound is silently dropped.
`HitFlash` is stored the same way and has the same limitation.

**Change:**
- Replace `map[Entity]PlaySoundIntent` with `[]PlaySoundIntent` in `World`.
- Add `EnqueueSoundIntent(entity Entity, intent PlaySoundIntent)` and
  `DrainSoundIntents() []SoundIntent`.
- The audio system drains the queue each frame; the queue is empty at the
  start of the next frame.

**Benefit:** Multiple sounds per entity per frame are preserved. The pattern
matches the book's Event Queue chapter exactly — producers and consumers are
decoupled in time.

**ADR:** `docs/adr/0011-event-queue-sound-intents.md`

---

### Service Locator 🔶
**Problem:** Scene context is `map[string]interface{}` — an untyped bag of
services. Type assertions (`context["fontLoader"].(*fonts.Loader)`) are
scattered across every `NewXxxScene` constructor and fail silently if a
key is missing or mistyped.

**Change:**
- Replace `map[string]interface{}` with a typed `SceneContext` struct:
  ```go
  type SceneContext struct {
      FontLoader  *fonts.Loader
      BlurShader  *ebiten.Shader
      Background  *ebiten.Image  // optional
      MapPath     string
  }
  ```
- `SceneFactory` becomes `func(SceneContext) (Scene, error)`.
- `TransitionToScene.WithContext` is replaced by direct struct field
  assignment.

**Benefit:** Compiler-checked dependencies; no runtime key typos; easier
to read what a scene requires.

**ADR:** `docs/adr/0012-typed-scene-context.md`

---

## Optimization Patterns

### Data Locality 🔶
**Problem:** `World` stores one `map[Entity]T` per component type (20 maps).
Iterating e.g. `TankMovers()` touches 6 separate maps, each with random
heap addresses. Cache misses accumulate as entity count grows.

**Change (deferred — profile first):** For the highest-churn component
groups (`Transform`, `TankVelocity`, `TankDrive`), replace the `map` with a
dense slice indexed by a stable slot, reserving the `map` for sparse
components (`HitFlash`, `PlaySoundIntent`).

**Constraint:** Do not change until a profiling run shows this is a
meaningful bottleneck. The current entity count (< 50) makes this
premature.

**ADR:** `docs/adr/0013-data-locality-dense-component-storage.md` *(deferred)*

---

### Dirty Flag 🔶
**Problem:** `world_queries.go` rebuilds and re-sorts query slices
(`Renderables()`, `TankMovers()`, `ArmedTanks()`, etc.) **on every call**.
`Renderables()` sorts by ZIndex then Entity every frame even when no
entity has been created or destroyed.

**Change:**
- Add a `dirty bool` field to `World`.
- `CreateEntity`, `DestroyEntity`, and `AttachComponent` set `dirty = true`.
- Each query caches its last result. If `!dirty`, return the cached slice.
- Clear `dirty` at the end of each pipeline tick via `world.ClearDirty()`.

**Benefit:** Query rebuild cost is paid once per frame at most, and zero
times when no structural change occurred (common in mid-game steady state).

**ADR:** `docs/adr/0014-dirty-flag-query-caching.md`

---

### Object Pool 🔶
**Problem:** Projectiles are spawned and destroyed frequently
(`weaponSpawnSystem` creates entities; `lifetimeSystem` + `projectileImpactsSystem`
destroy them). Each `DestroyEntity` call deletes from 20 separate maps;
each spawn inserts into several maps. At high fire rates this produces
steady GC pressure.

**Change:**
- Add a free-list (`[]Entity`) to `World`.
- `CreateEntity` pops from the free-list before incrementing `nextEntityID`.
- `DestroyEntity` pushes to the free-list instead of deleting (components
  are overwritten on reuse, so stale data is never visible).
- Add `World.MaxLiveEntities() int` for observability.

**Benefit:** Entity churn becomes allocation-free. Matches the book's
Object Pool chapter — the pool is inside `World` itself.

**ADR:** `docs/adr/0015-object-pool-entity-recycling.md`

---

### Spatial Partition ⬜ (deferred)
**Problem potential:** `obstacleCollisionSystem` and `projectileImpactsSystem`
loop over all tanks × all obstacles / projectiles — O(n²) pairs. With
< 20 moving entities this is irrelevant.

**Change (deferred):** Add a uniform-grid spatial hash once profiling
confirms collision is a bottleneck. Cell size = 2 × largest collider
radius.

**ADR:** `docs/adr/0016-spatial-partition-collision-broadphase.md` *(deferred)*

---

## Prioritised Action List

### P1 — Do now (correctness / low risk)
1. **Event Queue** — fix dropped sounds (map → slice).  
   Files: `game/ecs/world.go`, `game/ecs/world_components.go`, audio system.  
   ADR: 0011
2. **Double Buffer** — pre-allocate fade buffer in `SceneManager`.  
   File: `game/scene/manager.go`. No ADR needed.

### P2 — Before Slice 6 (quality / testability)
3. **Service Locator → Typed SceneContext** — remove `map[string]interface{}`.  
   Files: `game/scene/interface.go`, all `NewXxxScene` constructors,
   `app/tankism/main.go`.  
   ADR: 0012
4. **Flyweight** — split `Weapon` into `WeaponDef` + `WeaponState`.  
   Files: `game/ecs/components/components.go`, `game/sim/` weapon system,
   new `game/defs/weapondefs.go`.  
   ADR: 0008
5. **Type Object** — `TankDef` presets in `game/defs/`.  
   Files: `game/worldseed/tiled_seed.go`, new `game/defs/tankdefs.go`.  
   ADR: 0010
6. **Observer** — typed event bus for game outcomes.  
   Files: new `game/events/`, `game/scene/play.go`, `game/*.go`.  
   ADR: 0009

### P3 — After profiling (optimisation)
7. **Dirty Flag** — cache query slices.  
   Files: `game/ecs/world.go`, `game/ecs/world_queries.go`.  
   ADR: 0014
8. **Object Pool** — entity free-list.  
   Files: `game/ecs/world.go`.  
   ADR: 0015
9. **Data Locality** — dense slice storage for hot components. *(deferred)*  
   ADR: 0013
10. **Spatial Partition** — grid broadphase for collision. *(deferred)*  
    ADR: 0016

---

## ADR Summary

| # | Title | Priority |
|---|-------|----------|
| 0008 | Flyweight — Shared Weapon and AI Config | P2 |
| 0009 | Observer — Typed Event Bus for Game Outcomes | P2 |
| 0010 | Type Object — Tank Definitions | P2 |
| 0011 | Event Queue — Sound Intents | P1 |
| 0012 | Typed Scene Context | P2 |
| 0013 | Data Locality — Dense Component Storage *(deferred)* | P3 |
| 0014 | Dirty Flag — Query Result Caching | P3 |
| 0015 | Object Pool — Entity Recycling | P3 |
| 0016 | Spatial Partition — Collision Broadphase *(deferred)* | P3 |
