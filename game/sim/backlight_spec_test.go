package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestBacklightSpecDependencies(t *testing.T) {
	assertSystemSpec(t, backlightSpec(), SystemSpec{Name: "player_backlight_follow", Reads: []ecscomponents.Kind{ecscomponents.KindTransform}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}})
}
