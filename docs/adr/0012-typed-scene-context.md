# ADR 0012 — Typed Scene Context (replacing Service Locator map)

**Status:** Proposed  
**Pattern source:** [gameprogrammingpatterns.com/service-locator.html](https://gameprogrammingpatterns.com/service-locator.html)

---

## Context

`TransitionToScene.Context` is `map[string]interface{}`. Every `NewXxxScene`
constructor does untyped assertions such as:

```go
fontLoader, ok := context[scene.FontLoaderContextKey].(*fonts.Loader)
```

Problems:
- A missing or mistyped key panics at runtime, not compile time.
- Callers must know magic string constants.
- Adding a new service means touching `constants.go`, `main.go`, and every
  scene that needs the service.
- The compiler cannot tell you when a required field is missing.

## Decision

Replace `map[string]interface{}` with a typed `SceneContext` struct:

```go
// game/scene/interface.go
type SceneContext struct {
    FontLoader  *fonts.Loader   // required by all scenes
    BlurShader  *ebiten.Shader  // required by Menu, GameOver
    Background  *ebiten.Image   // optional; nil = no background
    MapPath     string          // required by Play
}
```

`SceneFactory` becomes:

```go
type SceneFactory func(ctx SceneContext) (Scene, error)
```

`TransitionToScene` carries a `SceneContext` value instead of a map:

```go
type TransitionToScene struct {
    NextSceneID SceneID
    Context     SceneContext
}
```

`app/tankism/main.go` builds a base `SceneContext` once and passes it (with
overrides per transition) to each factory.

## Consequences

**Good:**
- Missing fields are compiler errors, not runtime panics.
- No magic string constants; `constants.go` loses the context-key section.
- Easier to read what a scene actually needs (struct fields vs. map keys).
- Autocomplete works.

**Bad / risks:**
- `SceneContext` grows as new services are added. Fields that only some scenes
  need (e.g. `Background`) stay optional (zero value = not set).
- If a scene is added by an external contributor who cannot modify
  `SceneContext`, they must add a field. This is a feature, not a bug — it
  makes dependencies explicit.

## Alternatives considered

**Keep the map but add typed accessor helpers.** Reduces runtime panics but
does not eliminate them; string constants remain.

**Dependency injection framework.** Excessive for a small fixed set of
services.
