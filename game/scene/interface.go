package scene

import "github.com/hajimehoshi/ebiten/v2"

// SceneID uniquely identifies a scene type.
type SceneID string

const (
	SceneIDLoading  SceneID = "loading"
	SceneIDMenu     SceneID = "menu"
	SceneIDPlay     SceneID = "play"
	SceneIDGameOver SceneID = "gameover"
)

// Scene is the interface all scene implementations must satisfy.
type Scene interface {
	// Update processes one frame of logic, returning a result that may signal a scene transition.
	Update(deltaSeconds float64) SceneResult

	// Draw renders the scene to the screen.
	Draw(screen *ebiten.Image)
}

// SceneResult describes the outcome of a scene's Update call.
type SceneResult interface {
	isSceneResult() // marker method to prevent external implementations
}

// StayInScene signals that the scene should remain active.
type StayInScene struct{}

func (StayInScene) isSceneResult() {}

// TransitionToScene signals a request to transition to a different scene.
type TransitionToScene struct {
	NextSceneID SceneID
	// Context can hold arbitrary scene-transition data (e.g., game outcome).
	// The SceneManager will make this available to the next scene's factory.
	Context map[string]interface{}
}

func (TransitionToScene) isSceneResult() {}

// QuitGame signals that the application should exit.
type QuitGame struct{}

func (QuitGame) isSceneResult() {}

// NewTransitionToScene creates a transition result with an empty context.
func NewTransitionToScene(nextID SceneID) TransitionToScene {
	return TransitionToScene{
		NextSceneID: nextID,
		Context:     make(map[string]interface{}),
	}
}

// WithContext adds a key-value pair to the transition context.
func (t TransitionToScene) WithContext(key string, value interface{}) TransitionToScene {
	t.Context[key] = value
	return t
}
