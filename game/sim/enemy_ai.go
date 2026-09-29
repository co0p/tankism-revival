package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func UpdateChaserIntent(world *ecs.World, chaser, target ecs.Entity, stopDistance float64) error {
	chaserTransform, ok := world.Transform(chaser)
	if !ok {
		return fmt.Errorf("chaser transform missing")
	}

	targetTransform, ok := world.Transform(target)
	if !ok {
		return fmt.Errorf("target transform missing")
	}

	dx := targetTransform.X - chaserTransform.X
	dy := targetTransform.Y - chaserTransform.Y
	distanceSquared := dx*dx + dy*dy

	intent := ecscomponents.TankControlIntent{}

	if stopDistance <= 0 || distanceSquared > stopDistance*stopDistance {
		intent.Throttle = 1
	}

	if distanceSquared > 0.0001 {
		desiredRotation := math.Atan2(dx, -dy)
		turnDelta := shortestAngleDelta(chaserTransform.RotationRadians, desiredRotation)
		intent.Turn = clampUnit(turnDelta / (math.Pi / 4))

		if intent.Throttle > 0 && math.Abs(turnDelta) > (math.Pi*0.75) {
			intent.Throttle = 0.25
		}

		if distanceSquared <= 420*420 && math.Abs(turnDelta) <= 0.22 {
			intent.Fire = true
		}
	}

	if err := world.SetTankControlIntent(chaser, intent); err != nil {
		return fmt.Errorf("set chaser intent: %w", err)
	}

	return nil
}

func shortestAngleDelta(current, target float64) float64 {
	delta := math.Mod((target-current)+math.Pi, 2*math.Pi)
	if delta < 0 {
		delta += 2 * math.Pi
	}

	return delta - math.Pi
}
