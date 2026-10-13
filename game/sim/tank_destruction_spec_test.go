package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestTankDestructionSpecDependencies(t *testing.T) {
	assertSystemSpec(t, tankDestructionSpec(), SystemSpec{Name: "tank_destruction_effects", Reads: []ecscomponents.Kind{ecscomponents.KindHealth, ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindEnemyAI}, Writes: []ecscomponents.Kind{ecscomponents.KindHealth, ecscomponents.KindTransform, ecscomponents.KindSpriteRenderer, ecscomponents.KindCircleCollider, ecscomponents.KindSolidObstacle, ecscomponents.KindLifetime, ecscomponents.KindPointLight, ecscomponents.KindTankDrive, ecscomponents.KindEnemyAI}})
}
