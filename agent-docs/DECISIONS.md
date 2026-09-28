# Decision log

Use this lightweight log for choices that constrain future work. Add entries; do not rewrite accepted decisions to hide history. If a decision changes, mark the old entry superseded and link the replacement.

## D-001 — Ship a single-player duel before multiplayer

- **Status:** superseded by D-006
- **Date:** 2026-09-28
- **Context:** Tankism could grow into networked and larger game modes, but doing so before proving the local combat loop would multiply complexity and uncertainty.
- **Decision:** Finish and validate one player-versus-AI duel before designing or implementing networking.
- **Consequences:** Gameplay rules should remain deterministic and avoid needless coupling to local input, but agents must not build speculative network infrastructure during the MVP.

## D-002 — Start from a fresh repository

- **Status:** accepted
- **Date:** 2026-09-28
- **Context:** The owner wants to restart agentically without carrying the previous implementation into the new codebase.
- **Decision:** Begin with this documentation-only starter kit. Build new implementation code incrementally; do not copy legacy source by default.
- **Consequences:** Old behavior is inspiration, not compatibility scope. Every capability must be explicitly implemented and verified in the new repository.

## D-003 — Keep deterministic rules independent of graphics initialization

- **Status:** accepted
- **Date:** 2026-09-28
- **Context:** Gameplay development and CI should not require a display or graphics device for every test.
- **Decision:** Deterministic gameplay rules must be testable without creating a window or initializing graphics. Framework-bound rendering/input remains integration or manual-test territory.
- **Consequences:** Input snapshots, plain math helpers, and injected randomness are preferred where they create a useful test seam. This is not permission to abstract every framework API.

## D-004 — Use Kenney Topdown Tanks Redux as the visual style

- **Status:** accepted
- **Date:** 2026-09-28
- **Context:** The previous Tankism project used Kenney's Topdown Tanks Redux pack, and the owner wants that visual identity retained without importing assets before they are needed.
- **Decision:** Use the pinned [historical Tankism visuals](REFERENCES.md#visual-direction) and [Kenney's Topdown Tanks Redux](https://kenney.nl/assets/topdown-tanks-redux) as Tankism's required visual style and future asset source. Import only assets needed by an active slice and store all runtime visuals under `resources/` with source and CC0 license records.
- **Consequences:** Placeholder and custom visuals must remain compatible with this style. Documentation embeds historical examples but does not duplicate runtime assets. The legacy repository is a visual reference, not a source-code migration requirement.

## D-005 — Use a minimal explicit ECS

- **Status:** accepted
- **Date:** 2026-09-28
- **Context:** The owner requires ECS to remain a foundational part of Tankism, while the fresh implementation still needs to avoid inheriting the legacy project's complexity.
- **Decision:** Follow the pinned [historical Tankism ECS model](REFERENCES.md#ecs-direction): represent game objects as entity identifiers plus plain data components, and implement behavior through explicitly ordered systems. Build the smallest Tankism-specific ECS needed by each playable slice rather than copying the legacy ECS or adopting a speculative engine.
- **Consequences:** The first tank and every later gameplay object must participate in the ECS. World queries and deterministic systems must be graphics-independent and testable. Rendering consumes ECS state and resources but does not own gameplay state.

## D-006 — Make multiple composable AI tanks hunt the player

- **Status:** accepted
- **Date:** 2026-09-28
- **Context:** A one-on-one duel does not match the intended game mode. The player should face several pursuing tanks, with ECS composition controlling how each enemy hunts and what it can do.
- **Decision:** The first complete mode is a single-player hunt against multiple enemy tanks. Enemy behavior and capabilities are selected by component presence and tuned by component data, including perception, reaction, targeting, aiming, movement, and weapon capabilities. Do not model variants as subclasses or only as a global difficulty multiplier.
- **Consequences:** Victory requires defeating all active hunters; defeat occurs when the player is destroyed. AI systems must query capability components, support multiple enemies at once, and remain deterministic in tests. Initial slices may add one hunter first to validate a system, but the MVP is not complete until multiple hunters operate together.

## D-007 — Optimize for playable learning

- **Status:** accepted
- **Date:** 2026-09-28
- **Context:** The purpose of Tankism is not only to produce a game, but to learn ECS and game-programming patterns by implementing them in a real, enjoyable project.
- **Decision:** Every substantive slice must deliver a concrete game improvement and answer a focused learning question. Use existing licensed visual resources, while implementing understandable Tankism-specific ECS, collision, lighting, ambient audio, combat, AI, and lifecycle systems instead of outsourcing those learning goals to a general-purpose engine or opaque framework.
- **Consequences:** Tasks and handoffs record learning evidence and tradeoffs alongside tests and player-visible results. Implementations should favor clarity and inspectability over premature performance sophistication. Ebitengine and standard libraries still provide platform, graphics, input, image/audio decoding, and device integration; rebuilding those is not a learning goal.

## Proposed decisions for M0.1

- Go module path.
- Go and Ebitengine versions.
- primary development/release platform.
- reference logical resolution and window policy.
- initial aiming/control model.

## Entry template

```md
## D-NNN — Short title

- **Status:** proposed | accepted | superseded by D-NNN
- **Date:** YYYY-MM-DD
- **Context:** What forced a choice?
- **Decision:** What are we doing?
- **Consequences:** What becomes easier, harder, required, or forbidden?
```
