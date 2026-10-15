# Tankism agent guide

Build a compact, playable top-down tank game in small learning slices.

## Read first

1. `README.md`
2. `docs/DESIGN.md`
3. `docs/ROADMAP.md` (pick the first unblocked slice)
4. `docs/adr/` (architecture source of truth)

Code + tests are authoritative for current behavior.

## Core rules (with ADR references)

- Deliver one small playable slice at a time; no speculative engine work.
- Keep gameplay deterministic and framework-independent in core logic ([ADR 0002](docs/adr/0002-deterministic-core-framework-edges.md)).
- Keep ECS minimal and in-house: entities are IDs, components are data, systems are ordered behavior ([ADR 0001](docs/adr/0001-minimal-in-house-ecs.md), [ADR 0005](docs/adr/0005-explicit-ecs-boundaries-and-ordered-pipeline.md)).
- Keep one explicit simulation pipeline in `game/sim/pipeline.go`; order/dependency changes require tests or rationale ([ADR 0005](docs/adr/0005-explicit-ecs-boundaries-and-ordered-pipeline.md)).
- Load assets at composition root using `resources` + `lib` loaders; gameplay stores IDs, not framework handles ([ADR 0004](docs/adr/0004-resource-loading-at-composition-root.md), partly superseded by [ADR 0007](docs/adr/0007-scene-system.md) for scene-scoped fonts).
- Preserve package dependency direction; avoid upward imports and cycles ([ADR 0006](docs/adr/0006-go-package-dependency-direction.md)).
- `game/ecs`, `game/sim`, `game/worldseed` must have zero Ebitengine imports.
- `lib/*` must not import `game/*`.
- `ebiten.Image`, shaders, and audio handles must not appear in simulation components — only IDs.
- Do not expand MVP scope (no multiplayer/progression/procgen/editor-engine).
- Add attribution + license files for every external asset before use.
- Preserve user changes; do not revert unrelated work.

## Package intent

- `app/tankism/`: process entrypoint, SceneManager wiring, game loop. No gameplay logic here.
- `app/mapism/`: map editor entrypoint. Accepts `--map`, `--width`, `--height`.
- `game/`: runtime composition and loop orchestration.
- `game/ecs`, `game/sim`, `game/worldseed`: pure gameplay/core (no Ebitengine).
- `game/input`, `game/render`: Ebitengine adapters.
- `game/scene/`: scene system (loading, menu, play, gameover) + SceneManager.
- `lib/`: reusable loaders/helpers (fonts, audio, sprites, etc.).
- `resources/`: embedded assets + attribution/licenses.

## Entrypoints

- **Game:** `go run ./app/tankism [--map path/to/map.tmj]`
- **Map editor:** `go run ./app/mapism [--map path/to/map.tmj] [--width N] [--height N]`
- Default map path: `game.DefaultMapPath` (defined in `game/` package).

## Verification + handoff

```sh
gofmt -w <changed-go-files>
go test ./...
go build ./...
go run ./app/tankism
```

If a check cannot run, report exact blocker. Handoff should include outcome, files changed, checks run, risks, and next smallest slice.

## Gotchas

**go.sum / text/v2:** When you import `ebiten/v2/text/v2` (or any new package) for the first time, run `go mod tidy` before `go build`. Missing go.sum entries cause build failures.

**Embedded assets:** `resources/assets.go` is the single go:embed source. It covers three globs:
```
images/kenney_topdown_tanks_redux/*.png
audio/opengameart_top_down_shooter/*.ogg
fonts/*.ttf
```
New assets must go into one of these directories, or you must update the embed directive. Failing to do so causes silent runtime load errors.

**Fonts:** See [docs/how-to-add-a-font.md](docs/how-to-add-a-font.md). There are no silent fallbacks — missing font loader or failed face load errors at startup.

**Shaders:** `lib/shaders/loader.go` compiles `.kage` files from the embedded FS. Shaders must be under `resources/shaders/` (covered by `shaders/*.kage` in the embed directive). Pass compiled `*ebiten.Shader` to scenes via context the same way the font loader is passed.

**Scene wiring:** See [ADR 0007](docs/adr/0007-scene-system.md). All scene factories live in `app/tankism/main.go`. `SceneManager` is the `ebiten.Game` implementor; individual scenes are not. `SceneResult` is sealed — no new result types can be added outside `game/scene/interface.go`.

**`ObstacleBarrelRedTopPath`** is a compatibility alias for the gray barrel sprite — code using the "red" name gets the gray image silently.

**Magic exit string:** `main.go` treats an error whose message contains `"quit game"` as a normal exit, not a real error. This is how the Exit menu option terminates cleanly.

**Window scale:** `windowScale = 2` in main.go. The OS window is 2× the logical resolution. Logical dimensions come from `game.LogicalWidth` / `game.LogicalHeight`.

