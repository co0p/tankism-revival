package sim

import (
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"github.com/co0p/tankism-revival/lib/mathutil"
)

func IntegrateTankMovement(world *ecs.World, deltaSeconds float64) {
	for _, mover := range world.TankMovers() {
		throttle := mathutil.ClampUnit(mover.Intent.Throttle)
		turn := mathutil.ClampUnit(mover.Intent.Turn)

		transform := mover.Transform

		transform.RotationRadians = mathutil.NormalizeRadians(
			transform.RotationRadians + turn*mover.Drive.TurnSpeedRadiansPS*deltaSeconds,
		)

		currentForwardSpeed := mover.Velocity.ForwardUnitsPerSecond
		targetForwardSpeed := throttle * mover.Drive.MoveSpeed
		maxDelta := movementDeltaForIntent(mover, throttle, currentForwardSpeed, deltaSeconds)

		forwardSpeed := targetForwardSpeed
		if maxDelta > 0 {
			forwardSpeed = mathutil.MoveTowards(
				currentForwardSpeed,
				targetForwardSpeed,
				maxDelta,
			)
		}

		forwardX, forwardY := forwardVector(transform.RotationRadians)
		distance := forwardSpeed * deltaSeconds
		transform.X += forwardX * distance
		transform.Y += forwardY * distance

		_ = world.SetTransform(mover.Entity, transform)
		_ = world.SetTankVelocity(mover.Entity, ecscomponents.TankVelocity{ForwardUnitsPerSecond: forwardSpeed})
	}
}

func movementDeltaForIntent(mover ecs.TankMover, throttle, currentForwardSpeed, deltaSeconds float64) float64 {
	if throttle == 0 {
		if mover.Deceleration.CoastUnitsPerSecondSquared > 0 {
			return mover.Deceleration.CoastUnitsPerSecondSquared * deltaSeconds
		}
		if mover.Acceleration.ForwardUnitsPerSecondSquared > 0 {
			return mover.Acceleration.ForwardUnitsPerSecondSquared * deltaSeconds
		}

		return 0
	}

	throttleDirection := directionOf(throttle)
	currentDirection := directionOf(currentForwardSpeed)

	if currentDirection != 0 && throttleDirection != currentDirection {
		if mover.Deceleration.BrakeUnitsPerSecondSquared > 0 {
			return mover.Deceleration.BrakeUnitsPerSecondSquared * deltaSeconds
		}
		if mover.Deceleration.CoastUnitsPerSecondSquared > 0 {
			return mover.Deceleration.CoastUnitsPerSecondSquared * deltaSeconds
		}
	}

	if mover.Acceleration.ForwardUnitsPerSecondSquared > 0 {
		return mover.Acceleration.ForwardUnitsPerSecondSquared * deltaSeconds
	}

	return 0
}

func directionOf(value float64) float64 {
	if value > 0 {
		return 1
	}
	if value < 0 {
		return -1
	}

	return 0
}

func forwardVector(rotationRadians float64) (x float64, y float64) {
	return math.Sin(rotationRadians), -math.Cos(rotationRadians)
}
