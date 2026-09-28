# Product vision

## One sentence

Tankism is a fun top-down tank game built agentically as a hands-on way to learn ECS architecture and game-programming patterns, where multiple AI tanks hunt the player in a compact arena.

## Why build it

The project has two equally important outcomes:

1. **A fun game:** movement, combat, enemies, atmosphere, and feedback should combine into a replayable experience.
2. **Practical understanding:** each playable slice should teach how a game system works, why it is structured that way, and what tradeoffs appeared in use.

Existing licensed tiles and sprites provide the visual vocabulary. The learning work is in composing those resources through systems we understand: ECS storage and queries, transforms, movement, collision detection and resolution, rendering, lighting, ambient audio, combat, AI, state transitions, and testing.

## Experience goals

- **Immediate:** the player understands movement, aiming, shooting, health, and the objective without a tutorial wall.
- **Tactical:** facing, momentum, obstacles, and shot timing create meaningful choices.
- **Readable:** hits, damage, danger, victory, and defeat are unmistakable.
- **Cohesive:** visuals follow the [historical Tankism and Kenney reference](REFERENCES.md#visual-direction) rather than mixing unrelated styles.
- **Pressured:** several hunters force the player to move, prioritize threats, and use cover.
- **Replayable:** a hunt is short enough to restart immediately and varied enough to retry.
- **Focused:** finish one strong mode before adding breadth.
- **Educational:** important systems are small enough to read, test, explain, and change.

## MVP: single-player hunt

The first releasable slice contains:

- one local player tank;
- multiple AI enemy tanks actively hunting the player;
- one compact arena with solid boundaries and a few obstacles;
- tank movement and aiming;
- projectiles with cooldown, collision, damage, and finite lifetime;
- health for the player and every enemy tank;
- clear win and lose states;
- restart and return-to-menu controls;
- a minimal HUD and essential audiovisual feedback;
- ambient sound and a simple lighting model that strengthen the arena atmosphere;
- stable performance on the primary desktop target.

Enemy behavior is composed through ECS components. Components determine whether an enemy can hunt and tune capabilities such as perception, reaction time, aiming, movement, and firing. This allows several enemy variants without hard-coded tank subclasses or one monolithic difficulty flag.

The MVP is complete when a new player can launch the game, enter a hunt, understand the controls, defeat all hunters or be destroyed, and replay without developer intervention.

## Proposed default controls

These are defaults, not immutable design decisions:

| Action | Control |
| --- | --- |
| Move forward/backward | Up/Down or W/S |
| Rotate hull | Left/Right or A/D |
| Aim | Hull direction initially; independent mouse aim is a later evaluated slice |
| Fire | Space |
| Pause/back | Escape |
| Restart finished hunt | Enter or R |

Use a single aiming model for the first playable version. Do not build both hull-aim and independent turret-aim at once.

## Explicitly out of MVP

- Network multiplayer and authoritative servers.
- Sixteen-player deathmatch.
- Accounts, matchmaking, anti-cheat, economy, or paid cosmetics.
- Procedural level generation.
- Destructible terrain.
- Full progression through recruit, major, and general classes.
- A large pickup catalog.
- Advanced fog of war or field-of-view simulation.
- A general-purpose game engine or ECS rewrite.
- Production-scale tooling, editor support, or technically impressive systems that do not improve the current game.

Progression, pickups, and additional modes remain candidate ideas for later phases.

## Learning goals

By completing the single-player game, the repository should demonstrate:

- why entities, components, queries, and ordered systems are useful for this game;
- how fixed gameplay data remains separate from Ebitengine rendering, input, and audio adapters;
- how transforms connect world position, sprite placement, rotation, and parent-child relationships;
- how broad/narrow collision checks and collision resolution differ;
- how projectiles, damage, lifetime, and cleanup interact through system order;
- how component composition creates hunters with different perception, reactions, aim, movement, and weapons;
- how ambient light, local light sources, particles, and sound communicate game state;
- how deterministic tests and manual playtesting complement each other.

Learning should come from the smallest understandable implementation that supports current gameplay. It does not require recreating operating-system, graphics-driver, image-decoding, or audio-decoding facilities already supplied by Go and Ebitengine.

## Gameplay rules to validate early

Use these as starting hypotheses and tune through playtesting:

- A hunt should usually last 30–120 seconds.
- A direct hit should be significant but not usually an instant kill.
- Fire cooldown should make deliberate shots more effective than holding fire.
- Arena obstacles should create cover without trapping tanks.
- Hunters should coordinate only through simple, legible rules and remain beatable before their intelligence becomes sophisticated.
- Enemy composition should create recognizable strengths and weaknesses rather than only increasing health or damage.

Keep tunable values together and name them by intent. Avoid scattering unexplained literals through systems.

## Open product questions

The owner should answer these through small prototypes rather than prolonged up-front design:

1. Should the first version use hull-direction aiming or independent turret/mouse aiming?
2. Should movement retain the current momentum-based “Asteroids style,” or feel heavier and more immediate?
3. Which desktop OS is the primary release target?
4. Is local controller support important before or after MVP?

Until answered, agents should choose the smallest reversible implementation and record the assumption in `ACTIVE_TASK.md`.
