package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestProximitySoundSpecDependencies(t *testing.T) {
	assertSystemSpec(t, proximitySoundSpec(), SystemSpec{Name: "sound_volume_by_proximity", Reads: []ecscomponents.Kind{ecscomponents.KindPlaySoundIntent, ecscomponents.KindTransform, ecscomponents.KindSoundVolumeByDistance}, Writes: []ecscomponents.Kind{ecscomponents.KindPlaySoundIntent}})
}
