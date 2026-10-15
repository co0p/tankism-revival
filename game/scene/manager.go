package scene

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// SceneFactory is a function that constructs a scene, optionally using context from a prior transition.
type SceneFactory func(context map[string]interface{}) (Scene, error)

// ManagerConfig holds configuration for the SceneManager.
type ManagerConfig struct {
	// FadeDurationSeconds is the duration of each fade in/out. Set to 0 to disable fades.
	FadeDurationSeconds float64

	// FadeColor is the color to fade to/from.
	FadeColor color.RGBA

	// SceneFactories maps SceneID to its factory function.
	SceneFactories map[SceneID]SceneFactory

	// LogicalWidth and LogicalHeight define the screen dimensions.
	LogicalWidth  int
	LogicalHeight int
}

// FadeState tracks the state of a fade transition.
type FadeState struct {
	Active      bool
	Direction   FadeDirection // In or Out
	Elapsed     float64
	Duration    float64
	FadeColor   color.RGBA
	NextSceneID SceneID
	NextContext map[string]interface{}
}

type FadeDirection int

const (
	FadeDirectionOut FadeDirection = iota
	FadeDirectionIn
)

// SceneManager is the top-level game orchestrator that owns the active scene and drives transitions.
type SceneManager struct {
	config        ManagerConfig
	activeSceneID SceneID
	activeScene   Scene
	fadeState     FadeState
	nextScene     Scene // constructed during fade-out, swapped in during fade-in
	lastContext   map[string]interface{}
}

// NewSceneManager creates a scene manager and initializes with the loading scene.
func NewSceneManager(cfg ManagerConfig) (*SceneManager, error) {
	if cfg.FadeDurationSeconds < 0 {
		cfg.FadeDurationSeconds = 0
	}

	// Initialize the loading scene
	loadingFactory, ok := cfg.SceneFactories[SceneIDLoading]
	if !ok {
		return nil, fmt.Errorf("no factory provided for scene %q", SceneIDLoading)
	}

	loadingScene, err := loadingFactory(nil)
	if err != nil {
		return nil, fmt.Errorf("construct loading scene: %w", err)
	}

	return &SceneManager{
		config:        cfg,
		activeSceneID: SceneIDLoading,
		activeScene:   loadingScene,
		fadeState: FadeState{
			FadeColor: cfg.FadeColor,
		},
		lastContext: make(map[string]interface{}),
	}, nil
}

// Update processes one frame of scene logic and manages transitions.
func (sm *SceneManager) Update() error {
	deltaSeconds := 1.0 / 60.0

	// Update fade state
	if sm.fadeState.Active {
		sm.fadeState.Elapsed += deltaSeconds
		if sm.fadeState.Elapsed >= sm.fadeState.Duration {
			sm.fadeState.Elapsed = sm.fadeState.Duration
			if sm.fadeState.Direction == FadeDirectionOut {
				// Swap to the next scene
				sm.activeSceneID = sm.fadeState.NextSceneID
				sm.activeScene = sm.nextScene
				sm.nextScene = nil
				sm.lastContext = sm.fadeState.NextContext
				// Start fade-in
				sm.fadeState.Direction = FadeDirectionIn
				sm.fadeState.Elapsed = 0
			} else {
				// Fade-in complete
				sm.fadeState.Active = false
				sm.fadeState.Elapsed = 0
			}
		}
	}

	// Update the active scene only if not fading out
	if !sm.fadeState.Active || sm.fadeState.Direction == FadeDirectionIn {
		result := sm.activeScene.Update(deltaSeconds)
		if err := sm.handleSceneResult(result); err != nil {
			return err
		}
	}

	return nil
}

// handleSceneResult processes the result of a scene's Update call.
func (sm *SceneManager) handleSceneResult(result SceneResult) error {
	switch r := result.(type) {
	case StayInScene:
		// Nothing to do
		return nil

	case TransitionToScene:
		if sm.fadeState.Active {
			// Ignore transition requests while fading
			return nil
		}

		// Construct the next scene
		factory, ok := sm.config.SceneFactories[r.NextSceneID]
		if !ok {
			return fmt.Errorf("no factory provided for scene %q", r.NextSceneID)
		}

		nextScene, err := factory(r.Context)
		if err != nil {
			return fmt.Errorf("construct scene %q: %w", r.NextSceneID, err)
		}

		sm.nextScene = nextScene
		sm.fadeState.NextSceneID = r.NextSceneID
		sm.fadeState.NextContext = r.Context

		// Start fade transition if enabled; otherwise swap immediately
		if sm.config.FadeDurationSeconds > 0 {
			sm.fadeState.Active = true
			sm.fadeState.Direction = FadeDirectionOut
			sm.fadeState.Duration = sm.config.FadeDurationSeconds
			sm.fadeState.Elapsed = 0
		} else {
			// No fade: swap immediately
			sm.activeSceneID = r.NextSceneID
			sm.activeScene = nextScene
			sm.lastContext = r.Context
			sm.nextScene = nil
		}

		return nil

	case QuitGame:
		return fmt.Errorf("quit game")

	default:
		return fmt.Errorf("unknown scene result type: %T", result)
	}
}

// Draw renders the current scene and any active fade overlay.
func (sm *SceneManager) Draw(screen *ebiten.Image) {
	sm.activeScene.Draw(screen)

	// Draw fade overlay if active
	if sm.fadeState.Active {
		fadedScreen := ebiten.NewImage(sm.config.LogicalWidth, sm.config.LogicalHeight)
		alpha := sm.fadeState.Elapsed / sm.fadeState.Duration
		if sm.fadeState.Direction == FadeDirectionOut {
			// Fade out: alpha increases from 0 to 1
		} else {
			// Fade in: alpha decreases from 1 to 0
			alpha = 1.0 - alpha
		}

		fadeColor := color.RGBA{
			R: sm.fadeState.FadeColor.R,
			G: sm.fadeState.FadeColor.G,
			B: sm.fadeState.FadeColor.B,
			A: uint8(255 * alpha),
		}
		fadedScreen.Fill(fadeColor)
		screen.DrawImage(fadedScreen, nil)
		fadedScreen.Dispose()
	}
}

// Layout returns the logical screen dimensions.
func (sm *SceneManager) Layout(_, _ int) (screenWidth, screenHeight int) {
	return sm.config.LogicalWidth, sm.config.LogicalHeight
}
