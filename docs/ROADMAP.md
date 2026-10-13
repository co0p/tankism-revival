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

### 5a. Tank hit feedback

**Acceptance:** A damaging projectile hit creates a small, expanding ring flash at the impact point without tinting the tank sprite. The ring fades in 0.10 seconds and stays at the impact position as the tank moves. Player damage also triggers a small, fading camera shake without changing world positions or controls. Repeated hits refresh rather than stack the effects. Player destruction keeps its stronger shake. Tests cover damage-only triggering, effect expiration, and pipeline order.

**Learn:** Keep hit feedback driven by deterministic simulation data, with sprite drawing and camera offsets at the boundary.

This slice does not complete item #5; smoke particles are covered by slice 5b, while firing/impact/movement particles and additional audio cues remain separate work.

### 5b. Health-driven tank smoke

**Acceptance:** Player and enemy tanks below 75% health emit gray smoke circles extracted from Kenney's `explosionSmoke5.png`. Lower health produces denser, larger smoke. Puffs drift from their spawn point, grow, and fade, with at most 256 particles alive. Healthy or destroyed tanks stop emitting; existing puffs expire, including after game over. Ring impacts and camera shake remain unchanged. Tests cover extraction, health scaling, cleanup, the cap, and deterministic updates.

**Learn:** Drive a bounded particle simulation from health data while keeping asset extraction and rendering at the boundary.

### 5c. Impact sparks

**Acceptance:** Each projectile collision with a tank or solid cover creates one burst of eight orange-yellow sparks at the impact point. Sparks scatter, slow, and fade over 0.18–0.32 seconds, with at most 192 alive. Owner overlaps and projectile expiration do not spawn sparks. Effects continue fading after game over. Existing ring flashes, health-driven smoke, and shake stay unchanged. Tests cover one-burst behavior, collision types, lifetime, drag, the cap, and determinism.

**Learn:** Spawn bounded particle bursts from collision resolution, then draw short trails from pure motion data without adding assets or changing combat rules.

### 5d. Tank tracks (instead of tread dust)

**Acceptance:** Moving player and enemy tanks leave paired tread marks using Kenney's `tracksSmall.png`. Marks are spaced by actual travel, follow hull turns and reverse movement, stay below tanks/props/lighting, and fade over 14 seconds. Stationary turning does not pile up marks. Storage is capped at 512 marks, replacing the oldest when full. Destroyed tanks stop stamping; existing marks fade, including after game over. Tests cover spacing, turns, reverse movement, cleanup, the cap, and determinism.

**Learn:** Keep distance-driven ground effects independent from tank entities and draw them in an explicit ground pass.

### 5e. Projectile ground impacts

**Acceptance:** Missed projectiles land at their finite lifetime endpoint, leaving a small mark using `oilSpill_small.png` and six short-lived dust puffs. Travel cannot overshoot the lifetime. Tank/cover impacts on the final tick take priority over ground splashes. Marks fade in eight seconds; dust fades in under one second. Storage is capped at 128 marks and 192 puffs. Only effects overlapping the current viewport are drawn; simulation and expiration continue offscreen and after game over.

**Learn:** Separate finite projectile travel, collision priority, and capped effects from camera-dependent rendering.

## 6. Polish and release candidate

- Tune combat, arenas, hunter variants, controls, and feedback through repeatable playtests.
- Profile before optimizing; add only measured improvements.
- Add asset attribution, release notes, and platform smoke-test results.

**Learn:** use observations and measurements to guide technical and design changes.
