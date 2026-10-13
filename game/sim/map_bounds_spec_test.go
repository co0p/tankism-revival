package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestPlayerMapBoundsSpecDependencies(t *testing.T) {
	assertSystemSpec(t, playerMapBoundsSpec(), SystemSpec{Name: "player_clamp_to_map_bounds", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}})
}

func TestEnemyMapBoundsSpecDependencies(t *testing.T) {
	assertSystemSpec(t, enemyMapBoundsSpec(), SystemSpec{Name: "enemy_clamp_to_map_bounds", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds, ecscomponents.KindEnemyAI}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}})
}
