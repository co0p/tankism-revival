package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestObstacleCollisionSpecDependencies(t *testing.T) {
	assertSystemSpec(t, obstacleCollisionSpec(), SystemSpec{Name: "tank_obstacle_collision_resolution", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindCircleCollider, ecscomponents.KindTankDrive, ecscomponents.KindSolidObstacle}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}})
}
