package sim

import (
	"fmt"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"github.com/co0p/tankism-revival/lib/mathutil"
)

func clampEntityToBounds(world *ecs.World, entity, mapEntity ecs.Entity, requireTransform bool) error {
	mapBounds, err := mapBoundsFromEntity(world, mapEntity)
	if err != nil {
		return err
	}
	transform, ok := world.Transform(entity)
	if !ok {
		if requireTransform {
			return fmt.Errorf("entity %d transform missing", entity)
		}
		return nil
	}
	clampedX := mathutil.Clamp(transform.X, 0, mapBounds.Width)
	clampedY := mathutil.Clamp(transform.Y, 0, mapBounds.Height)
	if clampedX == transform.X && clampedY == transform.Y {
		return nil
	}
	transform.X = clampedX
	transform.Y = clampedY
	if err := world.SetTransform(entity, transform); err != nil {
		return fmt.Errorf("set clamped transform for entity %d: %w", entity, err)
	}
	return nil
}

func clampEnemyEntitiesToBounds(world *ecs.World, mapEntity ecs.Entity) error {
	for _, enemy := range world.EnemyAITanks() {
		if err := clampEntityToBounds(world, enemy.Entity, mapEntity, false); err != nil {
			return err
		}
	}
	return nil
}

func mapBoundsFromEntity(world *ecs.World, entity ecs.Entity) (ecscomponents.MapBounds, error) {
	if entity == 0 {
		return ecscomponents.MapBounds{}, fmt.Errorf("map entity is required")
	}
	bounds, ok := world.MapBounds(entity)
	if !ok {
		return ecscomponents.MapBounds{}, fmt.Errorf("map bounds missing for entity %d", entity)
	}
	return bounds, nil
}
