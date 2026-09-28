package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
)

func ResolveTankObstacleCollisions(world *ecs.World) error {
	obstacles := world.SolidObstacleColliderBodies()

	for _, tank := range world.TankColliderBodies() {
		resolved := tank.Transform

		for _, obstacle := range obstacles {
			dx := resolved.X - obstacle.Transform.X
			dy := resolved.Y - obstacle.Transform.Y

			totalRadius := tank.Collider.Radius + obstacle.Collider.Radius
			if totalRadius <= 0 {
				continue
			}

			distanceSquared := dx*dx + dy*dy
			if distanceSquared >= totalRadius*totalRadius {
				continue
			}

			if distanceSquared == 0 {
				dx = 1
				dy = 0
				distanceSquared = 1
			}

			distance := math.Sqrt(distanceSquared)
			penetration := totalRadius - distance
			normalX := dx / distance
			normalY := dy / distance

			resolved.X += normalX * penetration
			resolved.Y += normalY * penetration
		}

		if err := world.SetTransform(tank.Entity, resolved); err != nil {
			return fmt.Errorf("set resolved tank transform: %w", err)
		}
	}

	return nil
}
