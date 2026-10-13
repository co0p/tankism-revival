package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestTankMovementSpecDependencies(t *testing.T) {
	assertSystemSpec(t, tankMovementSpec(), SystemSpec{Name: "tank_movement_integration", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindTankAcceleration, ecscomponents.KindTankDeceleration, ecscomponents.KindTankVelocity, ecscomponents.KindTankControlIntent}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankVelocity}})
}
