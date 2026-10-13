package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestEngineAudioSpecDependencies(t *testing.T) {
	assertSystemSpec(t, engineAudioSpec(), SystemSpec{Name: "engine_audio_volume", Reads: []ecscomponents.Kind{ecscomponents.KindTankDrive, ecscomponents.KindTankVelocity, ecscomponents.KindEngineAudio}, Writes: []ecscomponents.Kind{ecscomponents.KindEngineAudio}})
}
