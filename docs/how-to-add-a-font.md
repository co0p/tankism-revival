# How to add a font

See [ADR 0007](adr/0007-scene-system.md) for the scene system design and [ADR 0004](adr/0004-resource-loading-at-composition-root.md) for the asset-loading rationale.

## Steps

1. Copy the TTF into `resources/fonts/`. The existing embed glob `fonts/*.ttf` in `resources/assets.go` picks it up automatically.
2. Font filenames may contain spaces (e.g. `"Kenney Blocks.ttf"`). Preserve the exact name. Loader paths are relative inside the embedded FS.
3. Add a named convenience method to `lib/fonts/loader.go` (follow `LoadKenneyBlocksFace` as the pattern).
4. At the composition root (`app/tankism/main.go`), the single `fonts.NewLoader(resources.FS)` instance is already created. Inject it into any new scene factory via `ctx[scene.FontLoaderContextKey] = fontLoader`.
5. In the scene constructor, call `loader.LoadYourFace(size)` and pass the result to `fonts.NewRenderer(face)`. Both functions fail hard — there is no silent fallback. A missing loader or a failed face load returns an error from the constructor, which surfaces at startup.

## No silent fallbacks

`fonts.NewRenderer(nil)` panics. `DrawText` returns an error if the face is nil. Scene constructors return an error if the loader is absent from context. If text is missing at runtime, the loader was not injected correctly into the scene factory.

## Current state

Loading, Menu, and GameOver scenes receive the font loader. PlayScene does not — add it to the Play factory in `app/tankism/main.go` before adding any text rendering to that scene.
