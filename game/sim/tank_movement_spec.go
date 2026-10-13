package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func tankMovementSystem() runtimeSystem {
	return runtimeSystem{spec: tankMovementSpec(), run: func(world *ecs.World, dt float64, _ PipelineConfig) error {
		IntegrateTankMovement(world, dt)
		return nil
	}}
}

func tankMovementSpec() SystemSpec {
	return SystemSpec{Name: "tank_movement_integration", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindTankAcceleration, ecscomponents.KindTankDeceleration, ecscomponents.KindTankVelocity, ecscomponents.KindTankControlIntent}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankVelocity}}
}
