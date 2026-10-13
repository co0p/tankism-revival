package sim

import (
	"fmt"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

type CameraFollower interface {
	FollowTarget(mapWidth, mapHeight, targetX, targetY float64)
}

type PipelineConfig struct {
	PlayerEntity        ecs.Entity
	MapEntity           ecs.Entity
	AudioListenerEntity ecs.Entity

	EnemyStopDistance float64

	ProjectileZIndex  int
	MuzzleFlashZIndex int

	TankDestruction TankDestructionConfig

	Camera        CameraFollower
	Smoke         *SmokeParticles
	Sparks        *SparkParticles
	Tracks        *Tracks
	GroundImpacts *GroundImpacts

	PlayerBackLightEntity   ecs.Entity
	PlayerBackLightDistance float64
}

type runtimeSystem struct {
	spec SystemSpec
	run  func(world *ecs.World, deltaSeconds float64, cfg PipelineConfig) error
}

var defaultPipeline = []runtimeSystem{
	enemyAISystem(),
	tankMovementSystem(),
	engineAudioSystem(),
	weaponSpawnSystem(),
	proximitySoundSystem(),
	groundImpactsSystem(),
	projectileIntegrationSystem(),
	hitFlashSystem(),
	sparkIntegrationSystem(),
	projectileImpactsSystem(),
	tankDestructionSystem(),
	tankCollisionSystem(),
	obstacleCollisionSystem(),
	lifetimeSystem(),
	playerMapBoundsSystem(),
	enemyMapBoundsSystem(),
	tracksSystem(),
	smokeParticlesSystem(),
	cameraFollowSystem(),
	backlightSystem(),
}

func DefaultPipelineSpecs() []SystemSpec {
	specs := make([]SystemSpec, len(defaultPipeline))
	for i, step := range defaultPipeline {
		specs[i] = SystemSpec{
			Name:   step.spec.Name,
			Reads:  append([]ecscomponents.Kind(nil), step.spec.Reads...),
			Writes: append([]ecscomponents.Kind(nil), step.spec.Writes...),
		}
	}
	return specs
}

func RunDefaultPipeline(world *ecs.World, deltaSeconds float64, cfg PipelineConfig) error {
	for _, step := range defaultPipeline {
		if err := step.run(world, deltaSeconds, cfg); err != nil {
			return fmt.Errorf("run system %q: %w", step.spec.Name, err)
		}
	}
	return nil
}
