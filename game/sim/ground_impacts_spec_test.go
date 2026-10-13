package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestGroundImpactsSpecDependencies(t *testing.T) {
	assertSystemSpec(t, groundImpactsSpec(), SystemSpec{Name: "ground_impact_effects", Reads: []ecscomponents.Kind{ecscomponents.KindGroundImpacts}, Writes: []ecscomponents.Kind{ecscomponents.KindGroundImpacts}})
}
