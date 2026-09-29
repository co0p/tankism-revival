package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
)

func ApplyProximitySoundVolume(world *ecs.World, listenerEntity ecs.Entity) error {
	if listenerEntity == 0 {
		return nil
	}

	listenerTransform, ok := world.Transform(listenerEntity)
	if !ok {
		return fmt.Errorf("audio listener transform missing")
	}

	for _, queued := range world.SoundIntents() {
		settings, ok := world.SoundVolumeByDistance(queued.Entity)
		if !ok {
			continue
		}

		emitterTransform, ok := world.Transform(queued.Entity)
		if !ok {
			continue
		}

		factor := proximityFactor(listenerTransform.X, listenerTransform.Y, emitterTransform.X, emitterTransform.Y, settings.MaxDistance)
		minFactor := clamp(settings.MinFactor, 0, 1)
		if factor < minFactor {
			factor = minFactor
		}

		intent := queued.Intent
		intent.Volume *= factor
		if intent.Volume < 0 {
			intent.Volume = 0
		}

		if err := world.SetPlaySoundIntent(queued.Entity, intent); err != nil {
			return fmt.Errorf("set proximity sound intent volume: %w", err)
		}
	}

	return nil
}

func proximityFactor(listenerX, listenerY, emitterX, emitterY, maxDistance float64) float64 {
	if maxDistance <= 0 {
		return 1
	}

	distance := math.Hypot(emitterX-listenerX, emitterY-listenerY)
	factor := 1 - (distance / maxDistance)
	return clamp(factor, 0, 1)
}
