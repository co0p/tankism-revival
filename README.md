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

## Runtime package architecture

- `app/mapism`: Ebitengine map editor app for authoring map descriptions
- `game/ecs`: minimal entity-component storage and queries
- `game/sim`: deterministic game rules/systems (movement, AI, shooting, collisions)
- `game/worldseed`: deterministic map/entity seeding for startup
- `game/input`: Ebitengine input adapter (`keyboard -> intent snapshot`)
- `game/render`: Ebitengine rendering adapter (sprites, lighting, HUD)

## Assets

- Current player placeholder sprite: `resources/images/kenney_topdown_tanks_redux/tank_red_small.png`
- Attribution and license details: [`resources/ATTRIBUTION.md`](resources/ATTRIBUTION.md)

## Local commands

Run these from the repository root:

```sh
go test ./...
go build ./...
go run ./app/tankism
go run ./app/mapism
```

Optional gameplay map path:

```sh
go run ./app/tankism --map resources/maps/mapism_default.tmj
```

See the roadmap before starting work; implement the earliest unblocked slice and keep the game runnable after it.

## Map editor (Mapism)

Run:

```sh
go run ./app/mapism
```

Optional map path:

```sh
go run ./app/mapism --map resources/maps/my_map.tmj
```

Optional map size (used when creating a new map file):

```sh
go run ./app/mapism --map resources/maps/my_map.tmj --width 60 --height 36
```

Controls:

- Scroll wheel (over right palette): select sprite
- Left drag: paint ground tiles
- Left click: place selected props sprite as a free object (non-tile)
- Right drag: erase ground tiles
- Right click: remove nearest selected props object
- `R` / `Shift+R`: rotate selected placement by 1° clockwise/counter-clockwise (hold to repeat; ground tiles snap to 90°)
- `P`: set player spawn at cursor
- `E`: add enemy spawn at cursor
- `Shift+E`: remove nearest enemy spawn
- `Ctrl+Shift+Arrow keys`: resize map in tiles
- `Ctrl+S`: save map
- `Ctrl+L`: reload map from disk
