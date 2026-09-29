package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	"github.com/co0p/tankism-revival/lib/mathutil"
)

func UpdateEngineAudioVolumes(world *ecs.World, deltaSeconds float64) error {
	for _, emitter := range world.EngineAudioEmitters() {
		engineAudio := emitter.EngineAudio

		speedFactor := 0.0
		if emitter.Drive.MoveSpeed > 0 {
			speedFactor = math.Abs(emitter.Velocity.ForwardUnitsPerSecond) / emitter.Drive.MoveSpeed
		}
		speedFactor = mathutil.Clamp(speedFactor, 0, 1)

		idleVolume := mathutil.Clamp(engineAudio.IdleVolume, 0, 1)
		movingVolume := mathutil.Clamp(engineAudio.MovingVolume, 0, 1)
		targetVolume := idleVolume + (movingVolume-idleVolume)*speedFactor
		currentVolume := mathutil.Clamp(engineAudio.CurrentVolume, 0, 1)

		attackPerSec := engineAudio.AttackPerSec
		releasePerSec := engineAudio.ReleasePerSec
		if attackPerSec <= 0 {
			attackPerSec = math.MaxFloat64
		}
		if releasePerSec <= 0 {
			releasePerSec = math.MaxFloat64
		}

		volumeDelta := targetVolume - currentVolume
		ratePerSec := releasePerSec
		if volumeDelta > 0 {
			ratePerSec = attackPerSec
		}

		maxStep := ratePerSec * deltaSeconds
		if volumeDelta > maxStep {
			volumeDelta = maxStep
		} else if volumeDelta < -maxStep {
			volumeDelta = -maxStep
		}

		engineAudio.CurrentVolume = mathutil.Clamp(currentVolume+volumeDelta, 0, 1)

		if err := world.SetEngineAudio(emitter.Entity, engineAudio); err != nil {
			return fmt.Errorf("set engine audio volume for entity %d: %w", emitter.Entity, err)
		}
	}

	return nil
}
