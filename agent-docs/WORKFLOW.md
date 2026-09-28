# Agentic development workflow

## 1. Choose a thin slice

Select the earliest unblocked roadmap item. A good slice changes one observable behavior and normally touches one subsystem. Examples:

- “The player tank cannot leave the arena,” not “finish physics.”
- “A bullet damages an enemy exactly once,” not “finish combat.”
- “Pressing R after defeat starts a fresh hunt,” not “finish scenes.”

If the task needs more than a short list of acceptance criteria, split it.

## 2. Write the contract first

Copy `templates/TASK.md` to `ACTIVE_TASK.md` and fill in:

- player or maintainer outcome;
- facts observed in the current code;
- assumptions and explicit non-goals;
- testable acceptance criteria;
- planned verification.
- one focused learning question and the evidence that will answer it.

Ask the owner when an irreversible product choice is required. For a reversible detail, state the assumption and proceed with the simplest option.

## 3. Establish the baseline

Before editing:

```sh
git status --short --branch
go test ./...          # after a Go module exists
go build ./...         # after a Go module exists
```

Do not take ownership of unrelated working-tree changes. Record pre-existing failures in the task instead of masking them.

## 4. Implement vertically

Prefer this order:

1. pure gameplay rule/data;
2. focused automated test;
3. system integration;
4. scene wiring;
5. rendering/audio feedback;
6. manual verification.

This makes the rule testable without requiring a window and minimizes debugging through visuals alone.

## 5. Close the learning loop

Before calling a slice complete, record a short note in [`LEARNINGS.md`](LEARNINGS.md):

- what pattern or subsystem was explored;
- what data entered and left the system;
- where it sits in update/draw order;
- what the tests and playtest demonstrated;
- one tradeoff or limitation intentionally left for later.

Prefer a concrete explanation tied to Tankism over a generic tutorial. Agents may research established patterns, but generated code must be inspected, tested, and understood in the context of this repository.

## 6. Verify proportionally

Minimum checks for Go changes:

```sh
gofmt -w <changed-go-files>
go test ./path/to/changed/package
go test ./...
go build ./...
```

For player-visible work, also use the run command established by M0.2 and check the exact acceptance criteria. Capture a screenshot or short recording for milestone exits when useful. Never claim manual verification if the environment could not display or control the game.

For documentation-only work, inspect links/paths and run `git diff --check`.

## 7. Review the diff

Before handing off:

- Is every changed file required for the task?
- Did the change preserve package direction?
- Are system-order changes intentional?
- Are gameplay constants named and centralized appropriately?
- Can deterministic logic run without graphics initialization?
- Did tests prove behavior rather than implementation details?
- Did documentation become stale because of the change?
- Can a future reader explain why the chosen pattern fits this slice?
- Did the work improve the game as well as exercise a concept?

## 8. Update shared memory

- Update `CURRENT_STATE.md` for facts that changed.
- Check roadmap items only after verification.
- Add a decision entry for meaningful tradeoffs or constraints.
- Keep transient investigation details in `ACTIVE_TASK.md`; promote demonstrated, durable lessons to `LEARNINGS.md`.
- End with a concise handoff using `templates/HANDOFF.md` when another session must continue.

## Bug workflow

For a bug:

1. Write exact reproduction steps and expected/actual results.
2. Add the smallest failing regression test when feasible.
3. Locate the responsible rule/system rather than patching the visual symptom.
4. Fix narrowly.
5. Run adjacent and repository-wide checks.
6. Record environmental or nondeterministic factors.

## Dependency workflow

Dependency changes are isolated tasks:

1. Record current versions and failures.
2. Read migration/release notes for the target versions.
3. Choose one target set rather than repeatedly chasing latest versions.
4. Update code and CI without gameplay changes.
5. Run tests/build and manually launch representative demos plus the main game.
6. Record the version decision and remaining compatibility issues.
