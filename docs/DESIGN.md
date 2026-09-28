# Game and technical design

## MVP: single-player tank hunt

One player tank fights several AI hunters in a small arena with walls and cover. Tanks move, turn, fire finite-lived projectiles, take damage, and explode. The match ends in victory when every hunter is destroyed or defeat when the player is destroyed; the player can restart immediately.

**Controls (initial):** `W/S` or arrows move forward/back, `A/D` or left/right turn, `Space` fires, `Escape` pauses/back, and `R` restarts a finished hunt. Start with hull-direction aim. Mouse/turret aiming is an experiment for after a playable hunt exists.

## Feel and readability goals

- A direct hit matters but normally does not destroy a full-health tank.
- Cooldown and obstacles reward deliberate shots, movement, and flanking.
- Hunters differ through capability data—perception, reaction rate, turn rate, engagement distance, and weapon—not opaque difficulty classes.
- Effects make motion, shots, hits, damage, danger, and victory/defeat obvious.
- The first arena favors clear cover and flow over visual complexity.

Keep tuneable gameplay values grouped and named by intent. Validate these hypotheses by playtesting rather than treating them as fixed rules.

## ECS model

Entities are opaque IDs. Components are plain data. Systems query relevant component sets and run in the order specified in `AGENTS.md`.

Initial components will likely include `Transform`, `Motion`, `TankControlIntent`, `Collider`, `Health`, `Weapon`, `Projectile`, `Lifetime`, `Sprite`, `Light`, `ParticleEmitter`, and AI capability components. Rendering resources must remain outside simulation components: a sprite or effect identifier may be data, but Ebitengine image, shader, and audio handles belong in adapters.

Enemy AI produces the same movement and fire intents used by the player. It does not get a separate movement, collision, or weapon implementation.

## Physics and collision progression

Use a fixed game-update delta supplied by Ebitengine's normal update loop. Begin with circle colliders for tanks/projectiles and axis-aligned rectangles for arena walls/obstacles. First identify candidate pairs with simple loops, then use narrow overlap tests and separate resolution from detection. Resolve tank-vs-solid penetration along the minimum translation direction; projectiles apply one hit then expire.

Only introduce a spatial hash or more advanced shapes when profiling or arena scale demonstrates a real need.

## Visual effects and lighting progression

Particles are short-lived simulation data: spawn position, velocity, age/lifetime, color/size interpolation, and optional drag. Start with muzzle flashes, impact sparks, tread dust, and explosions. Prefer fixed-capacity storage or a simple cap only after measuring allocations or frame instability.

Lighting begins as a rendering pass: draw the world, then combine a uniform ambient level with soft local radial light masks for muzzle flashes, explosions, and damaged tanks. Light components express intent; Ebitengine images, blend modes, and shaders remain in `game/render`. Keep the first version simple and visually legible before adding normal maps, shadow casting, or post-processing.

## Test seams

Unit test vector/math helpers, movement, cooldowns, collision queries and resolution, damage/lifetime, AI decisions, and match transitions. Integration test a small in-memory world through the ordered update pipeline. Manually smoke test window creation, input, rendering, sound, lighting, particles, and complete rounds.

## Non-goals for the MVP

No multiplayer, accounts, matchmaking, progression, shops, procedural generation, terrain destruction, advanced pathfinding, generic ECS framework, or general-purpose physics engine.
