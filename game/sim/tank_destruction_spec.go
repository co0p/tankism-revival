package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func tankDestructionSystem() runtimeSystem {
	return runtimeSystem{spec: tankDestructionSpec(), run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
		return ResolveTankDestructions(world, cfg.PlayerEntity, cfg.TankDestruction)
	}}
}

func tankDestructionSpec() SystemSpec {
	return SystemSpec{Name: "tank_destruction_effects", Reads: []ecscomponents.Kind{ecscomponents.KindHealth, ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindEnemyAI}, Writes: []ecscomponents.Kind{ecscomponents.KindHealth, ecscomponents.KindTransform, ecscomponents.KindSpriteRenderer, ecscomponents.KindCircleCollider, ecscomponents.KindSolidObstacle, ecscomponents.KindLifetime, ecscomponents.KindPointLight, ecscomponents.KindTankDrive, ecscomponents.KindEnemyAI}}
}
