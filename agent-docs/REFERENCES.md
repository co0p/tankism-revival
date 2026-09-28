# Visual and ECS references

This document records references for later implementation. It does **not** import legacy code or runtime assets into this documentation-only repository.

References are starting points for experiments, not specifications to reproduce blindly. Each implementation should be explained and evaluated against Tankism's current gameplay needs.

## Upstream snapshot

Use [`brotherhood-of-recursive-descent/tankism` at commit `8e07b3f`](https://github.com/brotherhood-of-recursive-descent/tankism/tree/8e07b3fc71fbc94f1a66d697e7a33203929039a7) as the stable historical reference. Pinning the commit prevents future upstream changes from silently changing the intended reference.

## Visual direction

The visual target is the readable, colorful top-down style demonstrated by the historical Tankism project and built primarily from [Kenney's Top-down Tanks Redux](https://kenney.nl/assets/topdown-tanks-redux):

- separate top-down tank hull and barrel sprites;
- grass and sand terrain tiles with simple roads;
- crates, barrels, barricades, trees, sandbags, and wire as readable cover;
- small, high-contrast projectiles, muzzle effects, explosions, and track marks;
- restrained ambient/point lighting that preserves gameplay readability.

Historical visual examples:

| Tank movement and scale | Lighting treatment | Effects treatment |
| --- | --- | --- |
| ![Historical Tankism tank movement](https://raw.githubusercontent.com/brotherhood-of-recursive-descent/tankism/8e07b3fc71fbc94f1a66d697e7a33203929039a7/docs/tank_move.gif) | ![Historical Tankism lighting](https://raw.githubusercontent.com/brotherhood-of-recursive-descent/tankism/8e07b3fc71fbc94f1a66d697e7a33203929039a7/docs/tank_lighting.gif) | ![Historical Tankism explosion animation](https://raw.githubusercontent.com/brotherhood-of-recursive-descent/tankism/8e07b3fc71fbc94f1a66d697e7a33203929039a7/docs/tank_boom.gif) |

The [historical `resources/` tree](https://github.com/brotherhood-of-recursive-descent/tankism/tree/8e07b3fc71fbc94f1a66d697e7a33203929039a7/resources) is an asset-selection reference, not permission to restore the entire old resource directory. When a playable slice needs art:

1. import only the required files into this repository's `resources/` directory;
2. load runtime visuals from `resources/`, not from remote URLs or documentation paths;
3. record each asset's source, author, and license in the local asset-attribution file;
4. retain a local copy of the applicable license alongside the resources;
5. verify that each selected file is covered by the upstream [Kenney CC0 notice](https://github.com/brotherhood-of-recursive-descent/tankism/blob/8e07b3fc71fbc94f1a66d697e7a33203929039a7/resources/License.txt), rather than assuming the repository-wide LGPL license applies to it.

Documentation images above are externally linked references only and are not runtime resources.

## ECS direction

Use the historical Tankism ECS as the conceptual reference:

- [ECS overview](https://github.com/brotherhood-of-recursive-descent/tankism/blob/8e07b3fc71fbc94f1a66d697e7a33203929039a7/docs/ecs.md)
- [generic entity/component definitions](https://github.com/brotherhood-of-recursive-descent/tankism/tree/8e07b3fc71fbc94f1a66d697e7a33203929039a7/lib/ecs)
- [Tankism components](https://github.com/brotherhood-of-recursive-descent/tankism/tree/8e07b3fc71fbc94f1a66d697e7a33203929039a7/game/ecs/components)
- [Tankism systems](https://github.com/brotherhood-of-recursive-descent/tankism/tree/8e07b3fc71fbc94f1a66d697e7a33203929039a7/game/ecs/systems)

The required model is:

- entities identify and group components;
- components are plain gameplay data;
- systems query the components they need and apply behavior;
- system execution order is explicit in the composition root;
- rendering reads ECS transform/sprite data and resolves sprite identifiers through the resource layer;
- deterministic gameplay systems remain testable without Ebitengine or a graphics device.

Do not copy the historical implementation wholesale. In particular, its generic `System` interface directly imports Ebitengine; the revival must preserve the graphics-independent gameplay boundary in D-003. Implement only the ECS capabilities demanded by the active playable slice, while retaining the historical data-oriented model.
