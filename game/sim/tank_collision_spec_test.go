package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestTankCollisionSpecDependencies(t *testing.T) {
	assertSystemSpec(t, tankCollisionSpec(), SystemSpec{Name: "tank_tank_collision_resolution", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindCircleCollider, ecscomponents.KindTankDrive}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}})
}
