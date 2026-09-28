# Roadmap

Work from top to bottom. Each item is a thin slice: write its acceptance criteria and learning question before implementation, then keep the executable runnable.

## 0. Bootstrap

- Create the Go module, Ebitengine window, `app/tankism` entry point, and package boundaries.
- Show a fixed logical-resolution background and close cleanly.
- Add one pure game-package test and document canonical local commands.

**Learn:** isolate framework composition from deterministic game code.

## 1. Minimal ECS and one player tank

- Implement ID allocation, plain component storage, creation/removal, and only the queries required now.
- Render a player entity with center-based transform and radians.
- Add tests for entity/component lifecycle and coordinate conventions.

**Learn:** entity/component separation and system-owned behavior.

## 2. Movement, arena, and solid collisions

- Add tank forward/back motion, turning, walls, and obstacle colliders.
- Detect tank/solid overlap and resolve it predictably.
- Manually verify the tank cannot leave the arena or pass through obstacles.

**Learn:** integration, broad versus narrow collision work, and resolution.

## 3. Combat loop

- Add firing cooldown, projectile spawning, lifetime, projectile collision, damage, and cleanup.
- Add one stationary target and hit feedback.
- Unit/integration test cooldown, one-hit behavior, and projectile expiration.

**Learn:** system ordering, spawn requests, ownership, and lifecycle.

## 4. Hunt and match flow

- Add multiple component-configured hunters that perceive, pursue, aim, and fire.
- Add health HUD, victory, defeat, and restart.
- Playtest toward 30–120 second rounds.

**Learn:** AI as intent-producing systems and composition-based variations.

## 5. Effects, lighting, and sound

- Add capped particles for firing, impacts, movement, and explosions.
- Add ambient plus local effect lights and essential audio cues.
- Manually verify effects improve readability without harming stable play.

**Learn:** simulation-driven effects and rendering/audio adapters at the boundary.

## 6. Polish and release candidate

- Tune combat, arenas, hunter variants, controls, and feedback through repeatable playtests.
- Profile before optimizing; add only measured improvements.
- Add asset attribution, release notes, and platform smoke-test results.

**Learn:** use observations and measurements to guide technical and design changes.
