package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func tankCollisionSystem() runtimeSystem {
	return runtimeSystem{spec: tankCollisionSpec(), run: func(world *ecs.World, _ float64, _ PipelineConfig) error {
		return ResolveTankTankCollisions(world)
	}}
}

func tankCollisionSpec() SystemSpec {
	return SystemSpec{Name: "tank_tank_collision_resolution", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindCircleCollider, ecscomponents.KindTankDrive}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}}
}
