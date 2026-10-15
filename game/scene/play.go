package scene

import (
	"github.com/co0p/tankism-revival/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MapPathContextKey = "mapPath"
)

// PlayScene wraps the game logic and manages game-over transitions.
type PlayScene struct {
	gameInstance  *game.Game
	lastFrame     *ebiten.Image // snapshot of the last rendered frame
	logicalWidth  int
	logicalHeight int
}

// NewPlayScene creates a play scene, initializing a fresh game instance.
// Context can provide "mapPath" to use a custom map, otherwise uses the default.
func NewPlayScene(context map[string]interface{}) (Scene, error) {
	mapPath := game.DefaultMapPath
	if context != nil {
		if path, ok := context[MapPathContextKey].(string); ok {
			mapPath = path
		}
	}

	gameInstance, err := game.NewGameWithMap(mapPath)
	if err != nil {
		return nil, err
	}

	return &PlayScene{
		gameInstance:  gameInstance,
		logicalWidth:  game.LogicalWidth,
		logicalHeight: game.LogicalHeight,
	}, nil
}

// Update processes one frame of game logic and checks for game-over conditions.
func (ps *PlayScene) Update(deltaSeconds float64) SceneResult {
	if err := ps.gameInstance.Update(); err != nil {
		_ = err // Update errors are non-fatal; game signals state via accessors.
	}

	if ps.gameInstance.IsPlayerDestroyed() {
		return NewTransitionToScene(SceneIDGameOver).
			WithContext(outcomeContextKey, outcomeDefeat).
			WithContext(BackgroundImageContextKey, ps.lastFrame)
	}

	if ps.gameInstance.AreAllEnemiesDestroyed() {
		return NewTransitionToScene(SceneIDGameOver).
			WithContext(outcomeContextKey, outcomeVictory).
			WithContext(BackgroundImageContextKey, ps.lastFrame)
	}

	return StayInScene{}
}

// Draw renders the game and captures the frame for use as a menu background.
func (ps *PlayScene) Draw(screen *ebiten.Image) {
	ps.gameInstance.Draw(screen)

	// Keep a snapshot of this frame so the menu can use it as a blurred background.
	if ps.lastFrame == nil {
		ps.lastFrame = ebiten.NewImage(ps.logicalWidth, ps.logicalHeight)
	}
	ps.lastFrame.DrawImage(screen, nil)
}
