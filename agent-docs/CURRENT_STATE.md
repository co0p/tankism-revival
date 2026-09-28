# Current project state

Last updated: 2026-09-28.

## Status

Tankism is a greenfield, documentation-only project. No application code, Go module, assets, tests, build scripts, or CI configuration have been created yet.

The files in this starter kit define:

- the dual goal of building a fun game and learning ECS/game-programming patterns;
- the product vision and first single-player MVP;
- the proposed technical direction;
- an incremental implementation roadmap;
- working rules for coding agents;
- task, decision, and handoff templates.

The owner has additionally confirmed two implementation constraints:

- [Kenney's Topdown Tanks Redux](VISUAL_STYLE.md) is the required visual style and future asset source, with pinned historical examples in [`REFERENCES.md`](REFERENCES.md), but no assets have been imported yet;
- gameplay objects must use the minimal, explicit ECS model linked in [`REFERENCES.md`](REFERENCES.md), built incrementally for Tankism rather than copied from the legacy repository.
- the first complete game mode has multiple AI tanks hunting the player; component composition determines each enemy's hunting intelligence and other capabilities.
- existing licensed tiles and sprites will be used as resources so implementation effort can focus on ECS, transforms, collision, combat, lighting, ambient sound, AI, and game-state patterns.

Do not claim that a gameplay feature exists until its code and verification are present in the new repository.

## Proposed technology

These choices are proposals to confirm in the first task, not existing dependencies:

- Go as the implementation language;
- Ebitengine v2 for rendering, input, audio, and the game loop;
- Go's standard testing tools for deterministic logic;
- GitHub Actions for continuous integration if the repository is hosted on GitHub.

Use currently supported stable versions when bootstrapping. Record exact versions and rationale in `DECISIONS.md`; do not recreate the legacy project's old dependency pins.

## Initial repository contents

When this starter kit is copied into a fresh repository, the expected initial tree is:

```text
AGENTS.md
agent-docs/
  README.md
  VISION.md
  CURRENT_STATE.md
  ROADMAP.md
  ARCHITECTURE.md
  REFERENCES.md
  VISUAL_STYLE.md
  LEARNINGS.md
  WORKFLOW.md
  DECISIONS.md
  ACTIVE_TASK.md
  templates/
    TASK.md
    HANDOFF.md
```

The bootstrap task may then add source, tests, CI, and a root project README.

## Immediate unknowns

Resolve these before or during the first implementation slice:

1. Final Go module path and repository hosting location.
2. Supported Go and Ebitengine versions.
3. Primary development and release operating systems.
4. Reference logical resolution and development window size.
5. Whether hull-direction aiming or independent turret aiming is the first control model.

Prefer reversible defaults when the owner has not decided. Record each assumption in `ACTIVE_TASK.md`.

## First evidence to collect

The first implementation task should leave behind:

- a module manifest and locked dependencies;
- a minimal application that opens a window and exits cleanly;
- one test that runs without initializing a graphics device;
- passing documented build/test commands;
- a CI workflow running the same checks;
- an updated snapshot in this file;
- a concise learning note explaining the first implemented pattern and its boundary.
