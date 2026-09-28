# Architecture guide

This file describes the intended starting architecture. Nothing here should be mistaken for already-implemented code.

## Guiding constraints

- Optimize for a small, complete single-player game rather than a reusable engine.
- Optimize for understandable implementations: each major system should make one game-programming pattern concrete and document its data flow and tradeoffs.
- Keep deterministic gameplay rules separable from rendering, audio, and live input.
- Add structure only when the next playable slice needs it.
- Prefer explicit data flow and update order over hidden global state.
- Keep platform/framework code at the edges so core rules are easy to test.
- Use existing licensed art as data, but implement the game-specific behavior that transforms it into a game.

## Proposed initial layout

Start smaller than the final layout. The bootstrap may create only the files needed for a window and one pure test.

```text
app/tankism/       application entry point and composition root
game/              Tankism-specific world state, rules, and systems
game/render/       Ebitengine rendering adapters when needed
game/input/        live input adapter and testable input snapshots
resources/         assets plus attribution/license records
```

Create `lib/` only after genuinely game-agnostic code has at least two uses. The ECS required below should remain Tankism-specific under `game/ecs/` until a concrete second use justifies extracting generic machinery.

## Dependency direction

```text
app/tankism -> game
app/tankism -> framework adapters
framework adapters -> game contracts/state
game -> Go standard library where practical
```

Rules:

- `game/` must not import `app/tankism/`.
- Pure game rules should not require a window, graphics device, or live keyboard.
- Scene and dependency composition belongs at the application boundary.
- Assets are inputs to rendering/audio, not the source of gameplay state.

## Runtime model

Begin with the simplest explicit model that supports the current milestone:

1. Read live controls into an input snapshot.
2. Convert player and AI input into intent.
3. Process cooldowns and spawn requests.
4. Integrate movement.
5. Detect and resolve collisions/damage.
6. Remove expired or destroyed objects.
7. Evaluate match state.
8. Render world, effects, and UI from the resulting state.

System order is gameplay behavior. Keep it visible in one composition location and cover important ordering assumptions with tests.

This pipeline is also the learning map. New systems should identify their inputs, outputs, order dependency, deterministic test seam, and visible or audible effect before implementation.

## ECS model

Tankism uses a small, explicit entity-component-system architecture from the first rendered game object onward:

- **Entities** are opaque identifiers and contain no behavior.
- **Components** are plain data, such as transform, sprite, movement intent, collider, health, or projectile ownership.
- **Systems** query the exact component sets they need and implement behavior in the documented runtime order.
- The application composition root owns the world and ordered system list; systems must not register themselves through hidden global state.

The bootstrap may draw only a background before an entity exists, but it must establish the minimal ECS before adding the first tank. Use the pinned [historical Tankism ECS](REFERENCES.md#ecs-direction) as the conceptual reference. Do not copy it wholesale, introduce a third-party ECS without a separate decision, or build speculative features such as archetype storage, parallel scheduling, serialization, or a general-purpose engine. Deterministic systems and world queries must remain testable without Ebitengine or a graphics device.

Rendering is an ECS system at the framework boundary. A sprite component refers to a loaded visual resource, while transform and gameplay components remain independent of file paths and image types. See [`VISUAL_STYLE.md`](VISUAL_STYLE.md) for the required asset source and resource layout and [`REFERENCES.md`](REFERENCES.md) for the pinned historical examples.

### Enemy hunting and abilities

Enemy tanks use the same base tank components and combat rules as the player. Their behavior is enabled and varied through component composition rather than enemy subclasses or behavior embedded in rendering code.

The first hunting slice should introduce only the data required by its behavior, likely including:

- a hunter/AI-control marker and current target;
- perception data such as detection range and target-memory duration;
- decision data such as reaction interval and preferred engagement distance;
- aiming data such as turn rate and accuracy tolerance;
- ordinary movement, weapon, cooldown, health, and collision components shared with other tanks.

Presence of a capability component determines whether the matching system processes an enemy; values on that component tune its competence. For example, an enemy without a weapon component can pursue but cannot fire, while reaction and aim values produce different hunting intelligence without a single opaque difficulty number. Keep perception, target selection, movement intent, aiming, and firing as deterministic, explicitly ordered systems. Add coordination or pathfinding components only when a concrete playable slice requires them.

## Systems to learn by implementing

- **Transforms and rendering:** convert center-based world transforms into ordered sprite draw operations while keeping image resources outside gameplay state.
- **Collision:** first implement understandable broad candidate checks and narrow overlap tests, then separate detection from resolution. Optimize only after measurement.
- **Combat and lifecycle:** model spawn requests, cooldowns, ownership, damage, death, and cleanup so update-order effects are visible and testable.
- **Lighting:** build an ambient-light baseline and composable local light sources as a rendering pass. Light components describe game intent; Ebitengine images/shaders stay in the adapter.
- **Ambient sound:** represent desired ambience and event cues independently from playback handles. An audio system maps that state to loaded resources and manages looping, transitions, and volume.
- **AI:** turn perception and capability components into intent; reuse ordinary movement and combat systems rather than giving AI separate physics or weapons.
- **Game state:** model start, active hunt, victory, defeat, and restart explicitly rather than scattering scene flags.

For each system, prefer a small complete version, observe it in play, record what was learned, and extend it only when the game exposes a limitation.

## Coordinate conventions to decide early

Before implementing movement/collision, define and test:

- whether object position represents its center or top-left corner;
- sprite anchor and rotation pivot;
- radians versus degrees and positive rotation direction;
- local, world, camera, and screen spaces;
- collision-shape dimensions and units;
- fixed logical resolution and window scaling behavior.

The recommended default is center-based world positions and radians, but M0.1 must confirm it.

## Test strategy

Use three layers:

1. **Pure unit tests:** movement math, cooldowns, health, collision calculations, AI decisions, and match transitions. No graphics initialization.
2. **Integration tests:** multiple rules/systems operating on a small in-memory world.
3. **Manual smoke tests:** window creation, rendering, live input, audio, scene flow, and visual feedback.

Avoid pixel-perfect screenshot tests at the beginning. Capture milestone screenshots as human-review evidence instead.

Tests are part of the learning mechanism: pure tests expose assumptions in rules and system order, integration tests show component interaction, and manual playtests reveal feel and communication that unit tests cannot measure.

## Error handling

- Return contextual errors from initialization and resource loading.
- Panic only for programmer invariants that cannot reasonably be recovered from.
- Do not ignore asset-loading errors.
- Avoid `os.Exit` below `main`, because it prevents cleanup and complicates testing.

## Performance guidance

Target a stable 60 updates/frames per second on the primary platform, subject to confirmation in M0.1. Use straightforward algorithms at MVP scale. Profile before adding spatial indexes, object pools, or other optimization infrastructure.

## Decisions to defer until required

- ECS storage/query optimization beyond the minimal required model;
- fixed timestep details beyond Ebitengine's normal update model;
- transient events versus explicit command/event queues;
- data-driven entity definitions;
- save compatibility;
- advanced collision shapes and spatial partitioning;
- networking architecture.
