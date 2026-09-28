# Tankism agent guide

This file is the operating contract for coding agents working in this repository.

## Mission

Build Tankism from a documentation-only repository into a fun, polished top-down tank game while deliberately learning ECS design and practical game-programming patterns through short, playable increments.

The first target is a **single-player hunt in which multiple AI tanks pursue the player**. Enemy hunting intelligence and other capabilities are composed from ECS components. Networking, large deathmatches, shops, and procedural generation are later possibilities, not current requirements.

Use existing licensed tiles and sprites so effort can focus on implementing and understanding gameplay systems such as transforms, rendering, collision detection and resolution, lighting, ambient sound, combat, and AI. Do not replace those learning goals with a ready-made game engine subsystem, but also do not reinvent platform APIs, image/audio codecs, or other infrastructure that Ebitengine already provides.

## Start here

Before changing code, read:

1. [`agent-docs/README.md`](agent-docs/README.md) — documentation map and working loop.
2. [`agent-docs/VISION.md`](agent-docs/VISION.md) — product direction and MVP boundaries.
3. [`agent-docs/CURRENT_STATE.md`](agent-docs/CURRENT_STATE.md) — what has actually been built so far.
4. [`agent-docs/ROADMAP.md`](agent-docs/ROADMAP.md) — ordered milestones and the next recommended work.
5. [`agent-docs/ARCHITECTURE.md`](agent-docs/ARCHITECTURE.md) — package boundaries and technical constraints.
6. [`agent-docs/LEARNINGS.md`](agent-docs/LEARNINGS.md) — evidence-based lessons from completed implementation slices.

The repository intentionally starts with documentation only. The documents under `agent-docs/` are the source of truth until code and executable tests exist; afterward, code and tests are the final authority for current behavior.

## Working principles

- Work on one thin, end-to-end slice at a time. Prefer a visible playable improvement over a broad framework rewrite.
- Treat every slice as both a game increment and a learning experiment: state the pattern or question being explored, then record what the implementation demonstrated.
- Keep each task small enough to review and verify independently.
- Inspect relevant code before changing it. Do not replace working systems without evidence that the active slice requires it.
- Separate facts from proposals. If a gameplay or architecture choice is not documented, ask the owner or record the assumption in the active task.
- Keep the game runnable after each completed task whenever the local environment permits it.
- Add or update tests for deterministic logic. Do not require a graphics device for logic tests where that dependency can reasonably be avoided.
- Prefer small in-house implementations of ECS and game-specific systems when they expose the concepts being studied. Avoid abstraction or optimization without a current gameplay need.
- Do not silently expand the MVP. Put good ideas in the future section of the roadmap instead.
- Do not edit generated or third-party assets unless the task explicitly requires it.
- Record the source and license of every added asset in the repository's asset-attribution file.

## Standard task loop

1. Read `agent-docs/ROADMAP.md` and select the earliest unblocked slice.
2. Create or update `agent-docs/ACTIVE_TASK.md` using `agent-docs/templates/TASK.md`.
3. Confirm current behavior and identify the smallest relevant set of files.
4. State acceptance criteria before implementation.
5. State the learning question and expected design lesson.
6. Implement the smallest complete change.
7. Run formatting and the narrowest useful tests, then the repository-wide checks when possible.
8. Manually exercise player-visible changes and record what was verified and learned.
9. Update `CURRENT_STATE.md`, `ROADMAP.md`, and `DECISIONS.md` only when the facts, plan, architecture, or durable lessons changed.
10. Replace `ACTIVE_TASK.md` with the next task or remove it when no task is active.

Do not mark a slice complete unless its acceptance criteria are met. If verification is blocked by the environment, describe the exact blocker and leave the item unverified.

## Commands

Run from the repository root:

```sh
# Format changed Go files.
gofmt -w path/to/file.go

# Run all tests after the Go module exists.
go test ./...

# Compile all packages without launching a window.
go build ./...

# Run the game after the initial scaffold defines this entry point.
go run ./app/tankism
```

The bootstrap task must add canonical build, test, and run commands to the root README and CI configuration. Until then, do not assume platform-specific graphics setup.

## Definition of done

A code task is done when all applicable items are true:

- Acceptance criteria are demonstrably satisfied.
- Changed Go files are formatted.
- New deterministic behavior has focused tests.
- `go test ./...` and `go build ./...` pass once a Go module exists, or an environmental failure is recorded precisely.
- Player-visible behavior was manually checked, or the reason it could not be checked is recorded.
- The implemented pattern, tradeoff, and observed lesson are briefly documented where future work will find them.
- No unrelated files or user changes were overwritten.
- Relevant agent documentation reflects the resulting state.

## Repository conventions

- The initial scaffold must confirm the module path; `github.com/co0p/tankism` is only the current proposal.
- The proposed main game entry point is `app/tankism/main.go`.
- Reusable, game-agnostic code may belong in `lib/` only when a concrete need exists.
- Tankism-specific rules and runtime code belong in `game/` or `app/tankism/` as described in `agent-docs/ARCHITECTURE.md`.
- Static assets and attribution belong under `resources/`.
- Use idiomatic Go for new APIs and filenames.
- Avoid hidden global state and nondeterministic logic in gameplay rules. Inject clocks, random sources, or input snapshots when tests need control.
- System order is gameplay behavior. Document and test changes to it.

## Scope and safety

- Never discard changes that were not created by the current task.
- Avoid dependency upgrades mixed with gameplay features; make upgrades their own reviewed slice.
- Avoid speculative abstractions. Introduce an interface or subsystem when a current feature needs it.
- Do not begin multiplayer work before the single-player MVP exit criteria in `ROADMAP.md` are met and the owner approves the scope change.
- When requirements conflict, use this priority: explicit owner request, accepted entries in `agent-docs/DECISIONS.md`, `agent-docs/VISION.md`, active task acceptance criteria, then the roadmap.

## Handoff

At the end of a work session, report:

- outcome and player-visible change;
- learning outcome and relevant tradeoffs;
- files changed;
- checks run and their results;
- unresolved risks or blockers;
- the smallest recommended next task.

Use [`agent-docs/templates/HANDOFF.md`](agent-docs/templates/HANDOFF.md) for longer handoffs.
