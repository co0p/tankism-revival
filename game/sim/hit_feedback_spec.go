package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func hitFlashSystem() runtimeSystem {
	return runtimeSystem{spec: hitFlashSpec(), run: func(world *ecs.World, dt float64, _ PipelineConfig) error {
		return UpdateHitFlashes(world, dt)
	}}
}

func hitFlashSpec() SystemSpec {
	return SystemSpec{Name: "hit_flash_expiration", Reads: []ecscomponents.Kind{ecscomponents.KindHitFlash}, Writes: []ecscomponents.Kind{ecscomponents.KindHitFlash}}
}
