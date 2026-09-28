# Visual style

## Reference

Tankism's required visual direction is [Kenney's Topdown Tanks Redux](https://kenney.nl/assets/topdown-tanks-redux): clean, colorful, readable top-down arcade military art with clearly differentiated tanks, projectiles, terrain, obstacles, and effects.

The pinned [historical Tankism examples and resource tree](REFERENCES.md#visual-direction) demonstrate this asset pack in context. They are references for later asset selection; the old source code is not a migration requirement.

## Usage rules

- Use Topdown Tanks Redux assets for the first playable arena instead of creating a competing placeholder art style.
- Keep gameplay readable: player and enemy silhouettes/colors must remain distinguishable, projectiles must stand out from terrain, and obstacles must clearly communicate collision.
- Prefer the smallest asset subset needed by the active slice. Do not import the complete pack preemptively.
- Every imported runtime asset must live under `resources/`, grouped by purpose (for example `resources/tanks/`, `resources/tiles/`, and `resources/objects/`). Documentation may link to this reference page without duplicating the assets.
- Rendering systems consume resource identifiers or loaded images; asset paths and image dimensions must not become authoritative gameplay state.
- Do not copy legacy Go resource-loading code by default. Add a small loader when the first rendered sprite requires it.

## Source and license

- **Creator:** Kenney (Kenney Vleugels)
- **Pack:** [Topdown Tanks Redux](https://kenney.nl/assets/topdown-tanks-redux)
- **License:** [Creative Commons CC0 1.0](https://creativecommons.org/publicdomain/zero/1.0/)
- **Legacy reference:** [brotherhood-of-recursive-descent/tankism](https://github.com/brotherhood-of-recursive-descent/tankism)

No visual files are stored in this repository yet. When assets are first imported, add the pack's license file and a precise inventory/source note under `resources/` in the same change.
