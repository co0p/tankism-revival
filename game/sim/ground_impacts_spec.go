package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func groundImpactsSpec() SystemSpec {
	return SystemSpec{Name: "ground_impact_effects", Reads: []ecscomponents.Kind{ecscomponents.KindGroundImpacts}, Writes: []ecscomponents.Kind{ecscomponents.KindGroundImpacts}}
}
func groundImpactsSystem() runtimeSystem {
	return runtimeSystem{spec: groundImpactsSpec(), run: func(_ *ecs.World, dt float64, cfg PipelineConfig) error {
		if cfg.GroundImpacts != nil {
			cfg.GroundImpacts.Advance(dt)
		}
		return nil
	}}
}
