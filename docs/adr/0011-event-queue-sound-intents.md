# ADR 0011 — Event Queue: Sound Intents

**Status:** Proposed  
**Pattern source:** [gameprogrammingpatterns.com/event-queue.html](https://gameprogrammingpatterns.com/event-queue.html)

---

## Context

`World` stores sound intents as `map[Entity]PlaySoundIntent`. This means
**at most one sound per entity per frame** can be queued. If a tank fires
and is hit in the same tick, the second `SetPlaySoundIntent` call silently
overwrites the first. Players lose audio feedback.

`HitFlash` has the same flaw — only one flash per entity per frame.

## Decision

Replace the `map[Entity]PlaySoundIntent` with a `[]PlaySoundIntent` slice
acting as a queue:

```go
// in World
soundIntents []PlaySoundIntent

func (w *World) EnqueueSoundIntent(entity Entity, intent PlaySoundIntent)
func (w *World) DrainSoundIntents() []PlaySoundIntent  // returns and clears
```

- **Producers** (weapon system, hit system, etc.) call `EnqueueSoundIntent`.
- **Consumer** (audio system) calls `DrainSoundIntents` once per frame,
  processes all intents, and the queue is empty for the next tick.

The same pattern applies to `HitFlash` if multiple flashes per entity per
frame are desired; at minimum, change it to a slice keyed off a separate
flash ID.

## Consequences

**Good:**
- Multiple sounds per entity per frame are preserved — no silent drops.
- Producer and consumer are decoupled in time (within the same frame),
  matching the Event Queue pattern intent.
- Queue ordering is deterministic (append order = system pipeline order),
  which satisfies ADR 0002.

**Bad / risks:**
- The audio system must process a variable-length slice rather than a fixed
  map. This is strictly simpler.
- Unbounded queue: if an entity enqueues sounds in a loop, the slice grows.
  Accept this for now; add a per-entity cap later if needed.

## Alternatives considered

**`map[Entity][]PlaySoundIntent`:** Groups by entity but adds a map lookup.
A flat slice is simpler and sufficient since the audio system processes all
intents anyway.

**Keep the map, document the limitation.** Cheap but leaves a real bug in
place (sounds drop silently).
