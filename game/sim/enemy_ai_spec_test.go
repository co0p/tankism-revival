package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestEnemyAISpecDependencies(t *testing.T) {
	assertSystemSpec(t, enemyAISpec(), SystemSpec{
		Name:   "ai_enemy_markov_astar_intent",
		Reads:  []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindHealth, ecscomponents.KindWeapon, ecscomponents.KindCircleCollider, ecscomponents.KindSolidObstacle, ecscomponents.KindMapBounds, ecscomponents.KindEnemyAI},
		Writes: []ecscomponents.Kind{ecscomponents.KindTankControlIntent, ecscomponents.KindEnemyAI},
	})
}
