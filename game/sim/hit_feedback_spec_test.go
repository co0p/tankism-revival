package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestHitFlashSpecDependencies(t *testing.T) {
	assertSystemSpec(t, hitFlashSpec(), SystemSpec{Name: "hit_flash_expiration", Reads: []ecscomponents.Kind{ecscomponents.KindHitFlash}, Writes: []ecscomponents.Kind{ecscomponents.KindHitFlash}})
}
