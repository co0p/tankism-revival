# Active task: <short outcome>

**Roadmap:** <milestone/slice>  
**Status:** ready | active | blocked | verification  
**Owner:** <human or agent/session>

## Outcome

Describe the observable player or maintainer result in one or two sentences.

## Current evidence

- Relevant files and behavior observed before editing.
- Existing tests or exact failures.
- Constraints from prior decisions.

## Learning question

- Name one ECS or game-programming concept this slice will make concrete.
- State what evidence will show that the concept works and is understood.

## Assumptions

- Reversible choices being made for this slice.

## Acceptance criteria

- [ ] One externally observable or testable result per line.
- [ ] Include error/edge behavior where relevant.
- [ ] Include documentation or migration expectations if needed.

## Non-goals

- Explicitly list nearby work that this task will not include.

## Likely files

- `path/to/file.go` — reason it may change.

## Verification plan

```sh
go test ./path/to/package
go test ./...
go build ./...
```

Manual check: describe exact steps and expected visible result.

## Notes / blockers

- Record discoveries that the next session would otherwise have to repeat.

## Learning notes

- Pattern implemented and its role in Tankism.
- Observed tradeoff or limitation.
- Follow-up question, if any.
- On completion, promote durable findings to `agent-docs/LEARNINGS.md`.
