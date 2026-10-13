package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func smokeParticlesSystem() runtimeSystem {
	return runtimeSystem{spec: smokeParticlesSpec(), run: func(world *ecs.World, dt float64, cfg PipelineConfig) error {
		if cfg.Smoke != nil {
			cfg.Smoke.Update(world, dt)
		}
		return nil
	}}
}

func smokeParticlesSpec() SystemSpec {
	return SystemSpec{Name: "tank_smoke_particles", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindHealth, ecscomponents.KindTankDrive, ecscomponents.KindSmokeParticles}, Writes: []ecscomponents.Kind{ecscomponents.KindSmokeParticles}}
}
