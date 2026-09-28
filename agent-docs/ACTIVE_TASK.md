# Active task: confirm foundations and scaffold the first window

**Roadmap:** M0.1–M0.2  
**Status:** ready  
**Owner:** unassigned

## Outcome

A clean repository has an agreed technical baseline and the smallest runnable Tankism application: it opens a window with a background compatible with the documented Kenney visual style and title, then exits cleanly. The bootstrap documents the Ebitengine lifecycle learned, and the next object-rendering slice has an explicit minimal-ECS constraint.

## Current evidence

- The repository contains only the documentation in this starter kit.
- Go and Ebitengine are proposed but not installed as project dependencies.
- The desired first complete game mode is a desktop single-player hunt against multiple AI tanks.
- No legacy source or assets should be copied by default.
- The owner selected Kenney's Topdown Tanks Redux as the visual style and required ECS for gameplay objects.
- The historical visual examples, resource tree, and ECS model are pinned in [`REFERENCES.md`](REFERENCES.md) for later implementation.
- Enemy hunting intelligence and abilities must be selected and tuned through ECS components; a one-enemy slice may validate behavior but does not complete the MVP.

## Assumptions requiring confirmation

- Proposed module path: `github.com/co0p/tankism`.
- Use the latest stable Go version supported by the selected stable Ebitengine v2 release.
- Primary development target: macOS; CI target: Linux. Confirm with the owner.
- Proposed logical resolution: 1280×720 in a resizable window.
- Do not import visual assets during bootstrap; use a simple background compatible with the documented Topdown Tanks Redux palette.
- Use hull-direction aiming for the first combat slice.

## Learning question

- Which responsibilities belong to Ebitengine's `Game` lifecycle and which belong to Tankism's future ECS world and systems?
- Evidence: the bootstrap can explain its update, draw, layout, and shutdown boundaries without placing gameplay objects or rules in framework-owned state.

## Acceptance criteria

- [ ] The owner confirms or replaces each assumption above.
- [ ] Accepted choices are recorded in `DECISIONS.md`.
- [x] Visual style, future resource location, and license expectations are documented without importing assets.
- [x] The ECS requirement and scope constraints are documented as an accepted decision.
- [ ] A minimal Go module uses pinned, supported dependency versions.
- [ ] `go run ./app/tankism` opens a titled 1280×720 development window with a visible placeholder background.
- [ ] Closing the window exits without an application error.
- [ ] At least one pure unit test runs without initializing graphics.
- [ ] `go test ./...` and `go build ./...` pass on the primary development machine.
- [ ] A root README documents prerequisites and exact run/build/test commands.
- [ ] `CURRENT_STATE.md` is updated with actual versions, files, and verification results.
- [ ] A short entry in [`LEARNINGS.md`](LEARNINGS.md) records the application lifecycle and the boundary reserved for ECS composition.

## Non-goals

- Player movement, tanks, ECS implementation, scenes, collision, combat, AI, audio, or asset imports.
- Copying code from the previous Tankism repository.
- Building abstractions for future gameplay.
- Release packaging.

## Likely files

- `go.mod`, `go.sum` — module and pinned dependency graph.
- `app/tankism/main.go` — application boundary and Ebitengine startup.
- `game/game.go` — minimal game implementation only if needed.
- `game/game_test.go` — a graphics-independent test of pure initial state/configuration.
- `README.md` — setup and commands.
- `.github/workflows/ci.yml` — preferably M0.3, unless included as a separate reviewable commit.

## Verification plan

```sh
gofmt -w <changed-go-files>
go test ./...
go build ./...
go run ./app/tankism
```

Manual check: confirm window title, dimensions, visible background, resize behavior, and clean close.

The first later slice that renders a tank must first add the minimal ECS world/entity/component/system path; direct ad hoc tank ownership by the Ebitengine game object is not acceptable.

## Handoff requirement

Report exact tool/dependency versions, commands and results, assumptions confirmed by the owner, and the next smallest roadmap slice. Do not mark M0.3 complete unless CI runs successfully on a clean checkout.
