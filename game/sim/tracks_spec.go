package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func tracksSystem() runtimeSystem {
	return runtimeSystem{spec: tracksSpec(), run: func(world *ecs.World, dt float64, cfg PipelineConfig) error {
		if cfg.Tracks != nil {
			cfg.Tracks.Update(world, dt)
		}
		return nil
	}}
}

func tracksSpec() SystemSpec {
	return SystemSpec{Name: "tank_tracks", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindHealth, ecscomponents.KindTracks}, Writes: []ecscomponents.Kind{ecscomponents.KindTracks}}
}
