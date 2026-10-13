package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func enemyAISystem() runtimeSystem {
	return runtimeSystem{spec: enemyAISpec(), run: func(world *ecs.World, dt float64, cfg PipelineConfig) error {
		return UpdateEnemyTankIntents(world, cfg.PlayerEntity, cfg.MapEntity, cfg.EnemyStopDistance, dt)
	}}
}

func enemyAISpec() SystemSpec {
	return SystemSpec{
		Name:   "ai_enemy_markov_astar_intent",
		Reads:  []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindHealth, ecscomponents.KindWeapon, ecscomponents.KindCircleCollider, ecscomponents.KindSolidObstacle, ecscomponents.KindMapBounds, ecscomponents.KindEnemyAI},
		Writes: []ecscomponents.Kind{ecscomponents.KindTankControlIntent, ecscomponents.KindEnemyAI},
	}
}
