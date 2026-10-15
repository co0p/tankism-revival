# ADR 0010 — Type Object: Tank Definitions

**Status:** Proposed  
**Pattern source:** [gameprogrammingpatterns.com/type-object.html](https://gameprogrammingpatterns.com/type-object.html)

---

## Context

Enemy tank variants (Aggressive, Cautious, Fleeing) are constructed in
`worldseed/tiled_seed.go` by inline component assembly. Config values
(move speed, weapon cooldown, AI behavior, sprite ID, collider radius) are
embedded at each call site. To add a new variant you must find all spawn
sites and copy-paste a block of component assignments.

The same pattern afflicts the player tank, which has its own inline
construction block.

## Decision

Introduce `TankDef` in a new `game/defs/` package:

```go
// game/defs/tankdefs.go
type TankDef struct {
    SpriteID        string
    MoveSpeed       float64
    TurnSpeedRad    float64
    Acceleration    float64
    CoastDecel      float64
    BrakeDecel      float64
    ColliderRadius  float64
    WeaponDefID     string
    EnemyAIDefID    string   // empty for player
    AudioDefID      string
}

var (
    TankDefPlayer     = TankDef{ … }
    TankDefAggressive = TankDef{ … }
    TankDefCautious   = TankDef{ … }
    TankDefFleeing    = TankDef{ … }
)

func SpawnTank(world *ecs.World, def TankDef, x, y float64) ecs.Entity
```

`worldseed` calls `defs.SpawnTank(world, defs.TankDefAggressive, x, y)`.
Component assembly lives in one place.

## Consequences

**Good:**
- All tank tuning is in `game/defs/` — one file, one place.
- `worldseed` becomes shorter and purely about reading Tiled data.
- New variant = new `TankDef` var + optional `SpawnTank` override; no spawn
  loop changes.
- `TankDef` is a plain struct value — serialisable, testable in isolation.

**Bad / risks:**
- `game/defs` must not import `game/scene` or `game/render` (ADR 0006).
  It may import `game/ecs` for `SpawnTank`.
- If def IDs are strings, a typo produces a silent nil lookup. Prefer typed
  constants or direct struct references over ID strings where possible.

## Alternatives considered

**Keep inline construction.** Low effort, but the problem grows with each
new variant or tuning pass.

**YAML/JSON data files.** More flexible, but adds a data-loading step and
runtime error surface. Out of scope for MVP.
