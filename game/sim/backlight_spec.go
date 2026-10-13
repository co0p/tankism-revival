package sim

import (
	"fmt"
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func backlightSystem() runtimeSystem {
	return runtimeSystem{spec: backlightSpec(), run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
		if cfg.PlayerBackLightEntity == 0 {
			return nil
		}
		transform, ok := world.Transform(cfg.PlayerEntity)
		if !ok {
			return fmt.Errorf("player transform missing")
		}
		x, y := TankBackLightPosition(transform.X, transform.Y, transform.RotationRadians, cfg.PlayerBackLightDistance)
		if err := world.SetTransform(cfg.PlayerBackLightEntity, ecscomponents.Transform{X: x, Y: y, RotationRadians: 0}); err != nil {
			return fmt.Errorf("update player back light transform: %w", err)
		}
		return nil
	}}
}

func backlightSpec() SystemSpec {
	return SystemSpec{Name: "player_backlight_follow", Reads: []ecscomponents.Kind{ecscomponents.KindTransform}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform}}
}
