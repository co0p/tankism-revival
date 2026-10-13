package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func obstacleCollisionSystem() runtimeSystem {
	return runtimeSystem{spec: obstacleCollisionSpec(), run: func(world *ecs.World, _ float64, _ PipelineConfig) error {
		return ResolveTankObstacleCollisions(world)
	}}
}

func obstacleCollisionSpec() SystemSpec {
	return SystemSpec{Name: "tank_obstacle_collision_resolution", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindCircleCollider, ecscomponents.KindTankDrive, ecscomponents.KindSolidObstacle}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}}
}
