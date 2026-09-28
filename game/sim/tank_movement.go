package sim

import (
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
)

func IntegrateTankMovement(world *ecs.World, deltaSeconds float64) {
	for _, mover := range world.TankMovers() {
		throttle := clampUnit(mover.Intent.Throttle)
		turn := clampUnit(mover.Intent.Turn)

		transform := mover.Transform

		transform.RotationRadians = normalizeRadians(
			transform.RotationRadians + turn*mover.Drive.TurnSpeedRadiansPS*deltaSeconds,
		)

		forwardX, forwardY := forwardVector(transform.RotationRadians)
		distance := throttle * mover.Drive.MoveSpeed * deltaSeconds
		transform.X += forwardX * distance
		transform.Y += forwardY * distance

		_ = world.SetTransform(mover.Entity, transform)
	}
}

func clampUnit(value float64) float64 {
	if value < -1 {
		return -1
	}
	if value > 1 {
		return 1
	}

	return value
}

func normalizeRadians(radians float64) float64 {
	wrapped := math.Mod(radians, 2*math.Pi)
	if wrapped < 0 {
		wrapped += 2 * math.Pi
	}

	return wrapped
}

func forwardVector(rotationRadians float64) (x float64, y float64) {
	return math.Sin(rotationRadians), -math.Cos(rotationRadians)
}
