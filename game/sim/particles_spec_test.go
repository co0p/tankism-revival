package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestSmokeParticlesSpecDependencies(t *testing.T) {
	assertSystemSpec(t, smokeParticlesSpec(), SystemSpec{Name: "tank_smoke_particles", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindHealth, ecscomponents.KindTankDrive, ecscomponents.KindSmokeParticles}, Writes: []ecscomponents.Kind{ecscomponents.KindSmokeParticles}})
}
