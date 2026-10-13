package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
)

const HitFlashDurationSeconds = 0.10

// UpdateHitFlashes runs before impacts so new hits get a full flash duration.
func UpdateHitFlashes(world *ecs.World, deltaSeconds float64) error {
	for entity, flash := range world.HitFlashStates() {
		flash.RemainingSeconds -= math.Max(0, deltaSeconds)
		if flash.RemainingSeconds <= 0 {
			world.RemoveHitFlash(entity)
			continue
		}
		if err := world.SetHitFlash(entity, flash); err != nil {
			return fmt.Errorf("update hit flash: %w", err)
		}
	}
	return nil
}
