package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func ResolveTankTankCollisions(world *ecs.World) error {
	bodies := world.TankColliderBodies()
	if len(bodies) < 2 {
		return nil
	}

	resolved := make(map[ecs.Entity]ecscomponents.Transform, len(bodies))
	for _, body := range bodies {
		resolved[body.Entity] = body.Transform
	}

	for i := 0; i < len(bodies)-1; i++ {
		a := bodies[i]
		aTransform := resolved[a.Entity]
		for j := i + 1; j < len(bodies); j++ {
			b := bodies[j]
			bTransform := resolved[b.Entity]

			dx := aTransform.X - bTransform.X
			dy := aTransform.Y - bTransform.Y

			totalRadius := a.Collider.Radius + b.Collider.Radius
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

			halfPush := penetration * 0.5
			aTransform.X += normalX * halfPush
			aTransform.Y += normalY * halfPush
			bTransform.X -= normalX * halfPush
			bTransform.Y -= normalY * halfPush

			resolved[a.Entity] = aTransform
			resolved[b.Entity] = bTransform
		}
	}

	for _, body := range bodies {
		if err := world.SetTransform(body.Entity, resolved[body.Entity]); err != nil {
			return fmt.Errorf("set resolved tank transform: %w", err)
		}
	}

	return nil
}
