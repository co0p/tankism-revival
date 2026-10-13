package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func proximitySoundSystem() runtimeSystem {
	return runtimeSystem{spec: proximitySoundSpec(), run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
		return ApplyProximitySoundVolume(world, cfg.AudioListenerEntity)
	}}
}

func proximitySoundSpec() SystemSpec {
	return SystemSpec{Name: "sound_volume_by_proximity", Reads: []ecscomponents.Kind{ecscomponents.KindPlaySoundIntent, ecscomponents.KindTransform, ecscomponents.KindSoundVolumeByDistance}, Writes: []ecscomponents.Kind{ecscomponents.KindPlaySoundIntent}}
}
