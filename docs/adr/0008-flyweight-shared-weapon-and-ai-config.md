# ADR 0008 — Flyweight: Shared Weapon and AI Config

**Status:** Proposed  
**Pattern source:** [gameprogrammingpatterns.com/flyweight.html](https://gameprogrammingpatterns.com/flyweight.html)

---

## Context

`components.Weapon` has ≈ 20 fields. Most of them are **read-only
configuration** that is identical across every tank of the same type:
`CooldownSeconds`, `ProjectileSpeed`, `ProjectileLifetime`,
`ProjectileDamage`, accuracy spreads, `ProjectileSpriteID`, all light and
muzzle-flash fields.

Right now every entity stores a full copy of this config in its own map slot.
If 10 Aggressive hunters exist, 10 identical copies of the same weapon config
exist in memory. Tuning a value means changing every spawn call.

The same applies to `EnemyAI`: `FavoredBehavior` is set at spawn and never
changes, yet it lives alongside the mutable fields (`RNGState`, `Path`,
`PathIndex`, `DecisionCooldownSeconds`).

## Decision

Split each component into a **def** (immutable, shared) and **state**
(mutable, per-entity):

```
components.WeaponDef    — all read-only config (CooldownSeconds, speeds, sprites, …)
components.WeaponState  — CooldownRemaining, ShotsFired

components.EnemyAIDef   — FavoredBehavior (and future tuning knobs)
components.EnemyAIState — Behavior, DecisionCooldownSeconds, RNGState, Path, PathIndex
```

A new `game/defs/` package holds named presets (e.g. `WeaponDefStandard`,
`EnemyAIDefAggressive`). Presets are package-level `var` values — plain Go
structs, zero overhead.

Systems that need config look up the def by value (passed directly); they
never mutate it. Systems that need mutable state use the state component.

## Consequences

**Good:**
- Tuning config lives in one place (`game/defs/`), not scattered across
  spawn sites.
- Entity memory footprint shrinks (mutable fields only in the component map).
- New weapon or AI variant = new def constant; no spawn-site changes.

**Bad / risks:**
- `World` gains two new component kinds per split (`WeaponDef`, `WeaponState`,
  `EnemyAIDef`, `EnemyAIState`), adding four more map entries.
- All callers of `Weapon()` and `EnemyAI()` must be updated.
- `world_test.go` tests that set `Weapon` directly must be updated.

## Alternatives considered

**Keep the monolith and just document which fields are config.** Cheaper
short-term but does not enforce immutability; a system could accidentally
mutate config.

**Use a pointer to a shared def struct inside the component.** Works but
makes the component non-comparable and complicates serialisation.
