package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func playerMapBoundsSystem() runtimeSystem {
	return runtimeSystem{spec: playerMapBoundsSpec(), run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
		return clampEntityToBounds(world, cfg.PlayerEntity, cfg.MapEntity, true)
	}}
}

func enemyMapBoundsSystem() runtimeSystem {
	return runtimeSystem{spec: enemyMapBoundsSpec(), run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
		return clampEnemyEntitiesToBounds(world, cfg.MapEntity)
	}}
}

func playerMapBoundsSpec() SystemSpec {
	return SystemSpec{Name: "player_clamp_to_map_bounds", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}}
}

func enemyMapBoundsSpec() SystemSpec {
	return SystemSpec{Name: "enemy_clamp_to_map_bounds", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds, ecscomponents.KindEnemyAI}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}}
}
