package sim

import (
	"fmt"
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func cameraFollowSystem() runtimeSystem {
	return runtimeSystem{spec: cameraFollowSpec(), run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
		if cfg.Camera == nil {
			return nil
		}
		mapBounds, err := mapBoundsFromEntity(world, cfg.MapEntity)
		if err != nil {
			return err
		}
		transform, ok := world.Transform(cfg.PlayerEntity)
		if !ok {
			return fmt.Errorf("player transform missing")
		}
		cfg.Camera.FollowTarget(mapBounds.Width, mapBounds.Height, transform.X, transform.Y)
		return nil
	}}
}

func cameraFollowSpec() SystemSpec {
	return SystemSpec{Name: "camera_follow_player", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds}, Writes: []ecscomponents.Kind{ecscomponents.KindCamera}}
}
