# Tankism learning project

This folder gives humans and coding agents a shared plan for building a fun top-down tank game while learning ECS and game-programming patterns through working implementations.

## Document map

| Document | Purpose | Update when |
| --- | --- | --- |
| [`VISION.md`](VISION.md) | Product goal, core loop, MVP boundary, open product questions | Product direction changes |
| [`VISUAL_STYLE.md`](VISUAL_STYLE.md) | Required art reference, resource policy, source, and license | Visual direction or asset usage changes |
| [`CURRENT_STATE.md`](CURRENT_STATE.md) | Evidence-based inventory and known problems | Repository behavior changes |
| [`ROADMAP.md`](ROADMAP.md) | Ordered milestones and completion criteria | A slice starts, finishes, or is reprioritized |
| [`ARCHITECTURE.md`](ARCHITECTURE.md) | Package map and design constraints | A structural decision changes |
| [`REFERENCES.md`](REFERENCES.md) | Pinned visual and ECS references for later implementation | The selected visual source or ECS direction changes |
| [`LEARNINGS.md`](LEARNINGS.md) | Evidence-based lessons from implemented systems | A learning slice is verified |
| [`WORKFLOW.md`](WORKFLOW.md) | How to turn an idea into a verified change | Team process changes |
| [`DECISIONS.md`](DECISIONS.md) | Lightweight architecture/product decision log | A consequential choice is accepted |
| [`ACTIVE_TASK.md`](ACTIVE_TASK.md) | The single task currently recommended | Work starts or hands off |
| [`templates/TASK.md`](templates/TASK.md) | Template for a bounded implementation task | Task format needs improvement |
| [`templates/HANDOFF.md`](templates/HANDOFF.md) | Template for cross-session context | Handoff format needs improvement |

## Source-of-truth order

1. Current code and executable tests describe what exists once implementation begins.
2. `AGENTS.md` describes how agents must work.
3. This folder describes the product direction and current plan.

When these disagree, do not silently choose. Correct stale factual documentation or record a product decision before implementing an ambiguous feature.

## Recommended build-and-learn sequence

1. Complete the task in `ACTIVE_TASK.md`.
2. Scaffold a green, documented development baseline.
3. Import a minimal set of licensed tiles and sprites as resources.
4. Learn ECS fundamentals by rendering and moving a component-composed player tank.
5. Implement collision, combat, and multiple component-configured hunters in small playable slices.
6. Implement lighting, ambient sound, and feedback to learn presentation systems and make the game enjoyable.
7. Record what each experiment taught before optimizing or generalizing it.
8. Reassess multiplayer and other modes only after the single-player game is fun.

The roadmap is ordered to keep the game playable and the learning observable at every stage. Shipping a feature without understanding or recording its design tradeoff is incomplete; studying a pattern without integrating it into the game is also incomplete.
