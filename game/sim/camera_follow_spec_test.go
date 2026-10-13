package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestCameraFollowSpecDependencies(t *testing.T) {
	assertSystemSpec(t, cameraFollowSpec(), SystemSpec{Name: "camera_follow_player", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds}, Writes: []ecscomponents.Kind{ecscomponents.KindCamera}})
}
