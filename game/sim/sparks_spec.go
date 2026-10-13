package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func sparkIntegrationSystem() runtimeSystem {
	return runtimeSystem{spec: sparkIntegrationSpec(), run: func(_ *ecs.World, dt float64, cfg PipelineConfig) error {
		if cfg.Sparks != nil {
			cfg.Sparks.Advance(dt)
		}
		return nil
	}}
}

func sparkIntegrationSpec() SystemSpec {
	return SystemSpec{Name: "spark_particle_integration", Reads: []ecscomponents.Kind{ecscomponents.KindSparkParticles}, Writes: []ecscomponents.Kind{ecscomponents.KindSparkParticles}}
}
