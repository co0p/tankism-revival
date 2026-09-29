package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"github.com/co0p/tankism-revival/lib/mathutil"
)

type CameraFollower interface {
	FollowTarget(mapWidth, mapHeight, targetX, targetY float64)
}

type SystemSpec struct {
	Name   string
	Reads  []ecscomponents.Kind
	Writes []ecscomponents.Kind
}

type PipelineConfig struct {
	PlayerEntity        ecs.Entity
	MapEntity           ecs.Entity
	AudioListenerEntity ecs.Entity

	EnemyStopDistance float64

	ProjectileZIndex  int
	MuzzleFlashZIndex int

	TankDestruction TankDestructionConfig

	Camera CameraFollower

	PlayerBackLightEntity   ecs.Entity
	PlayerBackLightDistance float64
}

type runtimeSystem struct {
	spec SystemSpec
	run  func(world *ecs.World, deltaSeconds float64, cfg PipelineConfig) error
}

var defaultPipeline = []runtimeSystem{
	{
		spec: SystemSpec{
			Name: "ai_enemy_markov_astar_intent",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindTransform,
				ecscomponents.KindTankDrive,
				ecscomponents.KindHealth,
				ecscomponents.KindWeapon,
				ecscomponents.KindCircleCollider,
				ecscomponents.KindSolidObstacle,
				ecscomponents.KindMapBounds,
				ecscomponents.KindEnemyAI,
			},
			Writes: []ecscomponents.Kind{ecscomponents.KindTankControlIntent, ecscomponents.KindEnemyAI},
		},
		run: func(world *ecs.World, deltaSeconds float64, cfg PipelineConfig) error {
			return UpdateEnemyTankIntents(world, cfg.PlayerEntity, cfg.MapEntity, cfg.EnemyStopDistance, deltaSeconds)
		},
	},
	{
		spec: SystemSpec{
			Name: "tank_movement_integration",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindTransform,
				ecscomponents.KindTankDrive,
				ecscomponents.KindTankAcceleration,
				ecscomponents.KindTankDeceleration,
				ecscomponents.KindTankVelocity,
				ecscomponents.KindTankControlIntent,
			},
			Writes: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankVelocity},
		},
		run: func(world *ecs.World, deltaSeconds float64, _ PipelineConfig) error {
			IntegrateTankMovement(world, deltaSeconds)
			return nil
		},
	},
	{
		spec: SystemSpec{
			Name: "engine_audio_volume",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindTankDrive,
				ecscomponents.KindTankVelocity,
				ecscomponents.KindEngineAudio,
			},
			Writes: []ecscomponents.Kind{ecscomponents.KindEngineAudio},
		},
		run: func(world *ecs.World, deltaSeconds float64, _ PipelineConfig) error {
			return UpdateEngineAudioVolumes(world, deltaSeconds)
		},
	},
	{
		spec: SystemSpec{
			Name: "weapon_cooldowns_and_projectile_spawn",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindTransform,
				ecscomponents.KindTankDrive,
				ecscomponents.KindTankVelocity,
				ecscomponents.KindTankControlIntent,
				ecscomponents.KindWeapon,
				ecscomponents.KindEnemyAI,
				ecscomponents.KindAudioPlayer,
			},
			Writes: []ecscomponents.Kind{
				ecscomponents.KindWeapon,
				ecscomponents.KindTransform,
				ecscomponents.KindSpriteRenderer,
				ecscomponents.KindCircleCollider,
				ecscomponents.KindProjectile,
				ecscomponents.KindLifetime,
				ecscomponents.KindPointLight,
				ecscomponents.KindPlaySoundIntent,
			},
		},
		run: func(world *ecs.World, deltaSeconds float64, cfg PipelineConfig) error {
			return UpdateWeaponsAndSpawnProjectiles(world, deltaSeconds, cfg.ProjectileZIndex, cfg.MuzzleFlashZIndex)
		},
	},
	{
		spec: SystemSpec{
			Name: "sound_volume_by_proximity",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindPlaySoundIntent,
				ecscomponents.KindTransform,
				ecscomponents.KindSoundVolumeByDistance,
			},
			Writes: []ecscomponents.Kind{ecscomponents.KindPlaySoundIntent},
		},
		run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
			return ApplyProximitySoundVolume(world, cfg.AudioListenerEntity)
		},
	},
	{
		spec: SystemSpec{
			Name:   "projectile_integration",
			Reads:  []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindProjectile},
			Writes: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindProjectile},
		},
		run: func(world *ecs.World, deltaSeconds float64, _ PipelineConfig) error {
			return IntegrateProjectiles(world, deltaSeconds)
		},
	},
	{
		spec: SystemSpec{
			Name: "projectile_impacts",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindProjectile,
				ecscomponents.KindTransform,
				ecscomponents.KindCircleCollider,
				ecscomponents.KindSolidObstacle,
				ecscomponents.KindTankDrive,
				ecscomponents.KindHealth,
			},
			Writes: []ecscomponents.Kind{ecscomponents.KindHealth, ecscomponents.KindProjectile},
		},
		run: func(world *ecs.World, _ float64, _ PipelineConfig) error {
			return ResolveProjectileImpacts(world)
		},
	},
	{
		spec: SystemSpec{
			Name: "tank_destruction_effects",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindHealth,
				ecscomponents.KindTransform,
				ecscomponents.KindTankDrive,
				ecscomponents.KindEnemyAI,
			},
			Writes: []ecscomponents.Kind{
				ecscomponents.KindHealth,
				ecscomponents.KindTransform,
				ecscomponents.KindSpriteRenderer,
				ecscomponents.KindCircleCollider,
				ecscomponents.KindSolidObstacle,
				ecscomponents.KindLifetime,
				ecscomponents.KindPointLight,
				ecscomponents.KindTankDrive,
				ecscomponents.KindEnemyAI,
			},
		},
		run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
			return ResolveTankDestructions(world, cfg.PlayerEntity, cfg.TankDestruction)
		},
	},
	{
		spec: SystemSpec{
			Name: "tank_tank_collision_resolution",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindTransform,
				ecscomponents.KindCircleCollider,
				ecscomponents.KindTankDrive,
			},
			Writes: []ecscomponents.Kind{ecscomponents.KindTransform},
		},
		run: func(world *ecs.World, _ float64, _ PipelineConfig) error {
			return ResolveTankTankCollisions(world)
		},
	},
	{
		spec: SystemSpec{
			Name: "tank_obstacle_collision_resolution",
			Reads: []ecscomponents.Kind{
				ecscomponents.KindTransform,
				ecscomponents.KindCircleCollider,
				ecscomponents.KindTankDrive,
				ecscomponents.KindSolidObstacle,
			},
			Writes: []ecscomponents.Kind{ecscomponents.KindTransform},
		},
		run: func(world *ecs.World, _ float64, _ PipelineConfig) error {
			return ResolveTankObstacleCollisions(world)
		},
	},
	{
		spec: SystemSpec{
			Name:   "lifetime_expiration",
			Reads:  []ecscomponents.Kind{ecscomponents.KindLifetime},
			Writes: []ecscomponents.Kind{ecscomponents.KindLifetime},
		},
		run: func(world *ecs.World, deltaSeconds float64, _ PipelineConfig) error {
			return UpdateLifetimes(world, deltaSeconds)
		},
	},
	{
		spec: SystemSpec{
			Name:   "player_clamp_to_map_bounds",
			Reads:  []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds},
			Writes: []ecscomponents.Kind{ecscomponents.KindTransform},
		},
		run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
			return clampEntityToBounds(world, cfg.PlayerEntity, cfg.MapEntity, true)
		},
	},
	{
		spec: SystemSpec{
			Name:   "enemy_clamp_to_map_bounds",
			Reads:  []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds, ecscomponents.KindEnemyAI},
			Writes: []ecscomponents.Kind{ecscomponents.KindTransform},
		},
		run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
			return clampEnemyEntitiesToBounds(world, cfg.MapEntity)
		},
	},
	{
		spec: SystemSpec{
			Name:   "camera_follow_player",
			Reads:  []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindMapBounds},
			Writes: []ecscomponents.Kind{ecscomponents.KindCamera},
		},
		run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
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
		},
	},
	{
		spec: SystemSpec{
			Name:   "player_backlight_follow",
			Reads:  []ecscomponents.Kind{ecscomponents.KindTransform},
			Writes: []ecscomponents.Kind{ecscomponents.KindTransform},
		},
		run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
			if cfg.PlayerBackLightEntity == 0 {
				return nil
			}

			transform, ok := world.Transform(cfg.PlayerEntity)
			if !ok {
				return fmt.Errorf("player transform missing")
			}

			backLightX, backLightY := TankBackLightPosition(
				transform.X,
				transform.Y,
				transform.RotationRadians,
				cfg.PlayerBackLightDistance,
			)

			if err := world.SetTransform(cfg.PlayerBackLightEntity, ecscomponents.Transform{
				X:               backLightX,
				Y:               backLightY,
				RotationRadians: 0,
			}); err != nil {
				return fmt.Errorf("update player back light transform: %w", err)
			}

			return nil
		},
	},
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

func TankBackLightPosition(x, y, rotationRadians, distance float64) (float64, float64) {
	forwardX := math.Sin(rotationRadians)
	forwardY := -math.Cos(rotationRadians)

	return x - forwardX*distance, y - forwardY*distance
}

func clampEntityToBounds(world *ecs.World, entity, mapEntity ecs.Entity, requireTransform bool) error {
	mapBounds, err := mapBoundsFromEntity(world, mapEntity)
	if err != nil {
		return err
	}

	transform, ok := world.Transform(entity)
	if !ok {
		if requireTransform {
			return fmt.Errorf("entity %d transform missing", entity)
		}
		return nil
	}

	clampedX := mathutil.Clamp(transform.X, 0, mapBounds.Width)
	clampedY := mathutil.Clamp(transform.Y, 0, mapBounds.Height)
	if clampedX == transform.X && clampedY == transform.Y {
		return nil
	}

	transform.X = clampedX
	transform.Y = clampedY

	if err := world.SetTransform(entity, transform); err != nil {
		return fmt.Errorf("set clamped transform for entity %d: %w", entity, err)
	}

	return nil
}

func clampEnemyEntitiesToBounds(world *ecs.World, mapEntity ecs.Entity) error {
	for _, enemy := range world.EnemyAITanks() {
		if err := clampEntityToBounds(world, enemy.Entity, mapEntity, false); err != nil {
			return err
		}
	}

	return nil
}

func mapBoundsFromEntity(world *ecs.World, entity ecs.Entity) (ecscomponents.MapBounds, error) {
	if entity == 0 {
		return ecscomponents.MapBounds{}, fmt.Errorf("map entity is required")
	}

	bounds, ok := world.MapBounds(entity)
	if !ok {
		return ecscomponents.MapBounds{}, fmt.Errorf("map bounds missing for entity %d", entity)
	}

	return bounds, nil
}
