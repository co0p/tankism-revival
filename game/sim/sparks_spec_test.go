package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestSparkIntegrationSpecDependencies(t *testing.T) {
	assertSystemSpec(t, sparkIntegrationSpec(), SystemSpec{Name: "spark_particle_integration", Reads: []ecscomponents.Kind{ecscomponents.KindSparkParticles}, Writes: []ecscomponents.Kind{ecscomponents.KindSparkParticles}})
}
