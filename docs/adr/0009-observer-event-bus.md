# ADR 0009 — Observer: Typed Event Bus for Game Outcomes

**Status:** Proposed  
**Pattern source:** [gameprogrammingpatterns.com/observer.html](https://gameprogrammingpatterns.com/observer.html)

---

## Context

`PlayScene.Update` currently polls `game.Game` for outcome predicates every
frame:

```go
if ps.gameInstance.IsPlayerDestroyed() { … }
if ps.gameInstance.AllEnemiesDestroyed() { … }
```

This has two problems:

1. **Coupling:** `game.Game` must expose predicate methods that let scenes
   reach into simulation internals.
2. **Polling cost:** the checks run every tick even in steady state, and every
   new outcome type (time limit, score threshold) requires a new predicate
   method on `game.Game`.

## Decision

Add a small typed event bus in `game/events/`:

```go
// game/events/events.go
type GameEvent int

const (
    PlayerDestroyed GameEvent = iota
    AllEnemiesDestroyed
)

type Bus struct { … }
func (b *Bus) Publish(e GameEvent)
func (b *Bus) Subscribe(handler func(GameEvent)) (cancel func())
```

- Simulation systems post events via `Bus.Publish` (e.g.
  `tankDestructionSystem` publishes `PlayerDestroyed`).
- `PlayScene` subscribes once in `NewPlayScene`; the callback sets an
  internal flag that `Update` checks.
- `game.Game` no longer needs predicate methods.

The bus is synchronous (handlers called inside `Publish`) and single-threaded
— no goroutines, matching the deterministic-core constraint (ADR 0002).

## Consequences

**Good:**
- Scenes are decoupled from simulation predicates.
- New outcomes add a new `GameEvent` constant and a publish call; no changes
  to `game.Game`'s public API.
- `game/events` has no Ebitengine import and no `game/scene` import, keeping
  the dependency direction clean (ADR 0006).

**Bad / risks:**
- One more package to understand.
- Subscribe/cancel pattern must be used correctly to avoid stale handlers
  across scene restarts.

## Alternatives considered

**Channels:** Would require goroutines or careful select logic. Overkill for
a single-threaded game loop.

**Return value from `game.Game.Update`:** Simpler but still couples the
scene to the game's return type; harder to extend.
