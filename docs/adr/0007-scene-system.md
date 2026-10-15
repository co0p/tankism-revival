# ADR 0007: Scene system

**Status:** Accepted  
**Date:** 2026-10-15

## Context

The game needs to transition between distinct application states:
1. A loading screen (assets, fake progress bar)
2. A main menu (navigation, play/exit)
3. Gameplay (the current `Game` struct and simulation)
4. A game-over screen (outcome display, restart/menu options)

Previously, the entire app was a single monolithic `Game` struct. Managing transitions, input routing, and rendering across these states requires a clear contract for what a "scene" is and how they communicate.

## Decision

We introduce a **Scene system** with the following characteristics:

1. **Scene Interface**: All scenes implement a minimal interface:
   ```go
   type Scene interface {
       Update(deltaSeconds float64) SceneResult
       Draw(screen *ebiten.Image)
   }
   ```

2. **SceneResult**: Scenes signal their intent via result values, not direct references:
   ```go
   type SceneResult interface { /* marker */ }
   
   type StayInScene struct{}
   type TransitionToScene struct {
       NextSceneID SceneID
       Context     map[string]interface{}
   }
   type QuitGame struct{}
   ```

3. **SceneManager**: A single coordinator (owned by the app layer, implements `ebiten.Game`) that:
   - Holds the currently active scene
   - Receives scene results and constructs the next scene via a factory function
   - Manages fade-in and fade-out transitions (swappable, currently 0.3s black fade)
   - Drives the top-level game loop

4. **No Scene-to-Scene Coupling**: Scenes do not know or reference each other. Navigation intent is expressed as a `SceneID`. Scene-specific context passes through `TransitionToScene.Context` (a map).

5. **Scenes Are Stateless at the Top Level**: Each transition results in a fresh scene instance. Persistent state (e.g., which map to play) is passed via context.

6. **Factory Pattern**: Each scene has a factory function that receives optional context and returns a constructed scene. Factories are registered in the app layer and passed to `SceneManager` on startup.

## Rationale

- **Testability**: Scenes can be unit-tested independently with mock inputs and outputs.
- **Decoupling**: No circular imports or scene-to-scene dependencies. The app layer owns the routing.
- **Replaceable Transitions**: Fade logic is isolated in `SceneManager`. Swapping for a wipe, slide, or instant cut requires only a `SceneManager` change.
- **Framework-Independent Core**: Game simulation (`game/sim`) remains untouched. Scenes are composition-only (ADR 0002, ADR 0005 still apply).
- **Clear Intent**: `SceneResult` makes it explicit what a scene wants to do next; there are no side effects like calling `os.Exit()` from within a scene.

## Consequences

1. **Easier to Add Scenes**: A new scene (pause menu, settings, etc.) requires only a new scene implementation and factory registration.

2. **Consistent Input Handling**: All scenes use `inpututil` to poll input directly. Input intent flows from input → scene logic → `SceneResult`.

3. **Context Passing**: Complex state (e.g., game outcome, map path) is passed via a map. This is less type-safe than function parameters but avoids circular dependencies.

4. **No Scene Retention**: Old scenes are discarded after a transition. If a scene needs to remember something (e.g., "player won"), it passes that via context to the next scene.

5. **Fade Transition Cost**: Rendering a fade-out overlay each frame has a small cost. Disable fades by setting `FadeDurationSeconds = 0` if needed.

## Alternatives Considered

1. **Direct Scene References**: `Scene` could hold a reference to the next scene:
   - **Rejected**: Creates cycles and circular imports. A menu scene would need to import the play scene.

2. **Centralized State Machine**: A separate state machine package that scenes query and command:
   - **Rejected**: Adds complexity. A simple result type is clearer.

3. **Singleton Game State**: A global game state object that scenes read/write:
   - **Rejected**: Breaks testability and makes dependencies implicit.

4. **Coroutine-Like Async Transitions**: Scenes could `yield` a transition:
   - **Rejected**: Overcomplicates the update loop for this use case.

## Future Work

- Add a pause menu scene (transition from play, return on resume).
- Implement settings scene (persist gameplay tuning values).
- Add a credits/attribution scene.
- Integrate replay/replay-save via context (pass replay data between scenes).
- Add transition sound effects (short audio played by `SceneManager` during fades).
