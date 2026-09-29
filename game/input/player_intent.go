package input

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"github.com/hajimehoshi/ebiten/v2"
)

// ReadPlayerTankIntent captures keyboard state and produces a tank control snapshot.
func ReadPlayerTankIntent() ecscomponents.TankControlIntent {
	var intent ecscomponents.TankControlIntent

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
