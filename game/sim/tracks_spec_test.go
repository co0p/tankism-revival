package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestTracksSpecDependencies(t *testing.T) {
	assertSystemSpec(t, tracksSpec(), SystemSpec{Name: "tank_tracks", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindHealth, ecscomponents.KindTracks}, Writes: []ecscomponents.Kind{ecscomponents.KindTracks}})
}
