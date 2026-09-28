package input

import (
	"github.com/co0p/tankism-revival/game/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// ReadPlayerTankIntent captures keyboard state and produces a tank control snapshot.
func ReadPlayerTankIntent() ecs.TankControlIntent {
	var intent ecs.TankControlIntent

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		intent.Throttle += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		intent.Throttle -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		intent.Turn -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		intent.Turn += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		intent.Fire = true
	}

	return intent
}
