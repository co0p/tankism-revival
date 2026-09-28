# Step-by-step roadmap

This is an ordered build-and-learn plan, not a calendar. Complete the smallest unchecked slice before starting later milestones unless it is explicitly blocked or reprioritized by the owner. Every slice should produce a playable or observable result, focused verification, and a short learning note about the pattern and tradeoffs encountered.

Status legend: `[ ]` not started, `[~]` active/partial, `[x]` verified.

## Milestone 0 — bootstrap a trustworthy baseline

Goal: every later agent can build, test, and launch the same project with known commands.

- [~] **M0.1 Confirm foundations.** Choose the module path, supported Go/Ebitengine versions, primary platform, logical resolution, and first aiming model; record decisions. ECS and visual direction are confirmed; remaining choices are open.
- [ ] **M0.2 Scaffold the application.** Add the smallest Go module and Ebitengine entry point that opens a titled window, draws a background consistent with the documented visual palette, and exits cleanly.
- [ ] **M0.3 Establish reliable checks.** Add a graphics-independent unit test, formatting/vetting/test commands, and CI that runs on a clean checkout.
- [ ] **M0.4 Document development.** Add a concise root README and stable run/build/test commands.
- [ ] **M0.5 Capture a reference screenshot.** Store the first application screenshot in `agent-docs/evidence/` with platform and commit details.

Exit criteria:

- a clean checkout has documented setup commands;
- build and automated tests pass on the supported environment;
- the application opens a visible window;
- known failures are recorded rather than left as tribal knowledge.

## Milestone 1 — learn ECS through a controllable arena

Goal: the player can drive one tank around a bounded arena with predictable collisions.

- [ ] **M1.1 Define coordinate conventions.** Decide transform origin, sprite anchor, rotation units/direction, and world-vs-screen coordinates; document with focused tests.
- [ ] **M1.2 Establish the playable ECS/resource path.** Using [`REFERENCES.md`](REFERENCES.md), add only the entity, transform, sprite-reference, query, and ordered-system capabilities needed for one tank. Import the selected art into `resources/` with attribution/license records, and render by resolving ECS sprite references through the resource layer. Do not copy the historical implementation wholesale.
- [ ] **M1.3 Integrate player movement.** Wire ECS controller and movement systems into the single-player scene with one explicit update order.
- [ ] **M1.4 Learn collision detection with arena boundaries.** Detect boundary overlap independently from movement and prevent the player from leaving the playable space.
- [ ] **M1.5 Learn collision resolution with one obstacle.** Separate collision detection from response and resolve tank/obstacle overlap without jitter or permanent overlap.
- [ ] **M1.6 Add a minimal controls overlay.** Show movement and fire controls in game.

Exit criteria: from the menu, the player can enter an arena, drive reliably, collide with boundaries/cover, and return to the menu.

## Milestone 2 — learn lifecycle and combat systems

Goal: shooting any tank produces deterministic damage and death with multiple targets present.

- [ ] **M2.1 Define combat data.** Add named configuration for tank health, projectile damage/speed/lifetime, and fire cooldown.
- [ ] **M2.2 Integrate firing.** Spawn projectiles at the correct muzzle position and heading.
- [ ] **M2.3 Integrate projectile movement/lifetime.** Projectiles move consistently and are cleaned up off-range or after timeout.
- [ ] **M2.4 Resolve projectile collisions.** A projectile damages the intended target once, ignores its owner, and is removed on impact.
- [ ] **M2.5 Resolve death.** Zero health disables/removes the tank and emits a clear event/state transition.
- [ ] **M2.6 Add health HUD.** Display player health and readable status for multiple enemies with obvious hit feedback.

Exit criteria: the player and multiple stationary enemy tanks can shoot, take damage, and die independently; the rules have deterministic tests independent of rendering.

## Milestone 3 — learn component-composed AI in a complete hunt

Goal: the game has a beginning, contest, ending, and replay.

- [ ] **M3.1 Add composable hunting behavior.** Add hunter, target, perception, reaction, and aiming components with explicitly ordered acquire/turn/approach/fire systems. Start with one hunter to validate the behavior; avoid pathfinding until obstacles prove it necessary.
- [ ] **M3.2 Add multiple distinct hunters.** Run several enemies concurrently and vary their intelligence and abilities through component presence and values, not subclasses or a monolithic difficulty flag.
- [ ] **M3.3 Keep hunters fair and legible.** Give each hunter reaction delay, aim tolerance, and the same movement/combat constraints used by the player; avoid unavoidable synchronized fire.
- [ ] **M3.4 Add win/lose states.** Win when all hunters are destroyed, lose when the player is destroyed, then stop combat and clearly announce the result.
- [ ] **M3.5 Add restart/menu flow.** Restart without relaunching the process and return cleanly to menu.
- [ ] **M3.6 Add a representative multi-hunter simulation.** Verify independent targeting, component-driven capability differences, and win/lose transitions without requiring pixel-perfect rendering.

Exit criteria: a new player can launch, survive and defeat multiple component-configured hunters (or be destroyed), and replay without developer assistance.

## Milestone 4 — learn presentation systems and make the game fun

Goal: the hunt is readable and satisfying rather than merely functional.

- [ ] **M4.1 Learn event-driven effects.** Add muzzle flash, impact, destruction, and restrained screen shake/particles without making effects authoritative gameplay state.
- [ ] **M4.2 Learn ambient lighting.** Implement an ambient-light pass, then add component-driven local lights where they improve readability and atmosphere.
- [ ] **M4.3 Learn ambient audio.** Add a looping arena soundscape, event cues, clean transitions, and independent volume control through an audio adapter.
- [ ] **M4.4 Improve game feel.** Tune spawn positions, hunter composition, arena layout, cover, controls, and combat pacing from playtest evidence.
- [ ] **M4.5 Add pause and concise settings/help.** Make lighting and audio settings understandable and preserve clean state transitions.
- [ ] **M4.6 Learn measurement-driven optimization.** Profile update/draw hot spots, establish a performance target, and optimize only a demonstrated bottleneck.
- [ ] **M4.7 Conduct at least three short playtests.** Record fun/readability findings separately from technical findings.

Exit criteria: players understand what happened, controls feel intentional, and common sessions do not expose severe bugs.

## Milestone 5 — release-quality single-player build

Goal: produce a distributable, documented MVP.

- [ ] Choose and verify target OS/architectures.
- [ ] Add version/build metadata and reproducible release commands.
- [ ] Verify licenses and credits for every shipped asset.
- [ ] Add smoke-test checklist and known-issues document.
- [ ] Package a release candidate and test it outside the repository.

## After MVP — decide, do not assume

Candidate directions:

- recruit/major/general progression;
- repair, armor, and ammunition pickups;
- more arenas and enemy behaviors;
- local multiplayer/controller support;
- fog of war and richer lighting;
- networking.

Multiplayer is a later learning and product phase, not an implied extension of the MVP. Before starting it, confirm that the single-player hunt is fun and write a decision covering local versus network play, protocol, authority, simulation model, latency strategy, cheating assumptions, and deployment cost.
