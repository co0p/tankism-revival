# Tankism Revival

Tankism is a small single-player top-down tank hunt built in **Go** with **Ebitengine**. It is intentionally developed in short, agent-friendly slices to learn practical ECS design, collision detection, lighting, particles, and game feel—not to create a reusable engine.

## Target experience

Drive a tank through a compact arena, use cover to survive several hunters, and destroy every enemy before they destroy you. A round should be readable, tense, and restartable in roughly 30–120 seconds.

The first release includes player movement and firing, enemy hunters, obstacles, projectiles, health, win/lose/restart flow, a basic HUD, particles, ambient/local lighting, and essential sound feedback. Multiplayer, progression, procedural generation, and destructible terrain are deliberately out of scope.

## Agent and design documentation

- [`AGENTS.md`](AGENTS.md): repository contract for coding agents.
- [`docs/DESIGN.md`](docs/DESIGN.md): gameplay and technical design.
- [`docs/ROADMAP.md`](docs/ROADMAP.md): implementation slices and completion criteria.
- [`docs/adr/`](docs/adr): durable architecture decisions.

## Assets

- Current player placeholder sprite: `resources/images/kenney_topdown_tanks_redux/tank_red_small.png`
- Attribution and license details: [`resources/ATTRIBUTION.md`](resources/ATTRIBUTION.md)

## Local commands

Run these from the repository root:

```sh
go test ./...
go build ./...
go run ./app/tankism
```

See the roadmap before starting work; implement the earliest unblocked slice and keep the game runnable after it.
