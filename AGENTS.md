# Tankism agent guide

This repository is an agentically developed learning project: build a fun, compact top-down tank game in Go with Ebitengine while keeping the ECS, collision, effects, and rendering code understandable.

## Read first

1. `README.md` — project intent and local commands.
2. `docs/DESIGN.md` — playable target, scope, and technical model.
3. `docs/ROADMAP.md` — ordered thin slices; choose the first unblocked item.
4. `docs/adr/` — accepted decisions and their reasons.

Code and tests are the authority for implemented behavior. Documentation is the authority for intended behavior only where code does not exist yet.

## Operating rules

- Deliver one small, playable vertical slice at a time. Do not build a generic engine first.
- Before coding, state the slice's acceptance criteria and learning question in the pull request, task note, or handoff.
- Keep game rules testable without Ebitengine, a GPU, a window, live input, or audio hardware.
- ECS entities are IDs; components are plain data; systems are ordered behavior. Do not place gameplay rules in render or input code.
- Keep system order explicit in one update pipeline. Changes to system order are gameplay changes and need tests or documented reasoning.
- Prefer a minimal in-house ECS for the MVP. Do not add an ECS, physics, or particle dependency without an ADR.
- Use straightforward collision and effects implementations first; profile before adding pooling, spatial indexes, shaders, or parallelism.
- Do not expand MVP scope to multiplayer, progression, procedural levels, or a general-purpose editor/engine.
- Add source and license information before adding any external asset.
- Preserve existing user changes. Do not revert, reset, or overwrite unrelated work.

## Expected layout

```text
app/tankism/       executable and composition root
game/              pure Tankism world, ECS, rules, and systems
game/input/        Ebitengine input adapter and input snapshots
game/render/       Ebitengine-only rendering, lighting, and particles
resources/         assets and ATTRIBUTION.md
docs/              design, roadmap, and architecture decisions
```

Create packages only when a current slice needs them. `game/` must not import `app/tankism/`; pure simulation packages must not import Ebitengine.

## Default update model

1. Capture live input into an input snapshot.
2. Produce player and AI intent.
3. Advance cooldowns and create spawn requests.
4. Integrate movement.
5. Detect and resolve collisions; apply damage.
6. Expire projectiles, particles, and destroyed entities.
7. Evaluate win/lose/restart state.
8. Render world, lights/effects, and UI from current state.

Start with center-based positions, radians, and positive clockwise screen rotation; confirm these conventions with tests before collision or rendering expands.

## Verification and handoff

For every completed code slice:

```sh
gofmt -w <changed-go-files>
go test ./...
go build ./...
go run ./app/tankism
```

Run the commands that exist; record exact blockers rather than claiming unrun checks passed. Manually verify visual/audio changes. A handoff must list the outcome, learning result, files changed, checks run, known risks, and the next smallest roadmap item.
