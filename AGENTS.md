# Tankism agent guide

Build a compact, playable top-down tank game in small learning slices.

## Read first

1. `README.md`
2. `docs/DESIGN.md`
3. `docs/ROADMAP.md` (pick first unblocked slice)
4. `docs/adr/` (architecture source of truth)

Code + tests are authoritative for current behavior.

## Core rules (with ADR references)

- Deliver one small playable slice at a time; no speculative engine work.
- Keep gameplay deterministic and framework-independent in core logic ([ADR 0002](docs/adr/0002-deterministic-core-framework-edges.md)).
- Keep ECS minimal and in-house: entities are IDs, components are data, systems are ordered behavior ([ADR 0001](docs/adr/0001-minimal-in-house-ecs.md), [ADR 0005](docs/adr/0005-explicit-ecs-boundaries-and-ordered-pipeline.md)).
- Keep one explicit simulation pipeline; order/dependency changes require tests or rationale ([ADR 0005](docs/adr/0005-explicit-ecs-boundaries-and-ordered-pipeline.md)).
- Load assets at composition root using `resources` + `lib` loaders; gameplay stores IDs, not framework handles ([ADR 0004](docs/adr/0004-resource-loading-at-composition-root.md)).
- Preserve package dependency direction; avoid upward imports and cycles ([ADR 0006](docs/adr/0006-go-package-dependency-direction.md)).
- Do not expand MVP scope (no multiplayer/progression/procgen/editor-engine).
- Add attribution + license files for every external asset before use.
- Preserve user changes; do not revert unrelated work.

## Package intent (short)

- `app/tankism/`: process entrypoint only.
- `app/mapism/`: map editor entrypoint only.
- `game/`: runtime composition and loop orchestration.
- `game/ecs`, `game/sim`, `game/worldseed`: pure gameplay/core.
- `game/input`, `game/render`: Ebitengine adapters.
- `lib/`: reusable loaders/helpers.
- `resources/`: embedded assets + attribution/licenses.

## Verification + handoff

```sh
gofmt -w <changed-go-files>
go test ./...
go build ./...
go run ./app/tankism
```

If a check cannot run, report exact blocker. Handoff should include outcome, files changed, checks run, risks, and next smallest slice.
