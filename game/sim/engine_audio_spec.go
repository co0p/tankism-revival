package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func engineAudioSystem() runtimeSystem {
	return runtimeSystem{spec: engineAudioSpec(), run: func(world *ecs.World, dt float64, _ PipelineConfig) error {
		return UpdateEngineAudioVolumes(world, dt)
	}}
}

func engineAudioSpec() SystemSpec {
	return SystemSpec{Name: "engine_audio_volume", Reads: []ecscomponents.Kind{ecscomponents.KindTankDrive, ecscomponents.KindTankVelocity, ecscomponents.KindEngineAudio}, Writes: []ecscomponents.Kind{ecscomponents.KindEngineAudio}}
}
