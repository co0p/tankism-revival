package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

type TankDestructionConfig struct {
	ExplosionSpriteID        string
	ExplosionBurstSpriteID   string
	ExplosionSmokeSpriteID   string
	ExplosionPlumeSpriteID   string
	EnemyWreckSpriteID       string
	ExplosionZIndex          int
	EnemyWreckZIndex         int
	ExplosionLifetime        float64
	ExplosionBurstLifetime   float64
	ExplosionSmokeLifetime   float64
	ExplosionPlumeLifetime   float64
	ExplosionLightRadius     float64
	ExplosionLightIntensity  float64
	BurstLightRadius         float64
	BurstLightIntensity      float64
	ExplosionLightColor      ecscomponents.LightColor
	ExplosionSmokeRotationPI bool
	EnemyWreckColliderRadius float64
}

func ResolveTankDestructions(world *ecs.World, playerEntity ecs.Entity, cfg TankDestructionConfig) error {
	for _, tank := range world.TankHealthStates() {
		health := tank.Health
		if health.Current > 0 || health.Max <= 0 {
			continue
		}

		if err := spawnTankExplosion(world, tank.Transform, cfg); err != nil {
			return err
		}

		if tank.Entity == playerEntity {
			health.Max = 0
			health.Current = 0
			if err := world.SetHealth(tank.Entity, health); err != nil {
				return fmt.Errorf("mark player destroyed: %w", err)
			}
			continue
		}

		if _, enemy := world.EnemyAI(tank.Entity); enemy {
			world.DestroyEntity(tank.Entity)
			if err := spawnEnemyWreck(world, tank.Transform, cfg); err != nil {
				return err
			}
			continue
		}

		world.DestroyEntity(tank.Entity)
	}

	return nil
}

func spawnTankExplosion(world *ecs.World, transform ecscomponents.Transform, cfg TankDestructionConfig) error {
	if cfg.ExplosionSpriteID == "" || cfg.ExplosionLifetime <= 0 {
		return nil
	}

	explosion := world.CreateEntity()
	if err := world.SetTransform(explosion, ecscomponents.Transform{X: transform.X, Y: transform.Y, RotationRadians: transform.RotationRadians}); err != nil {
		return fmt.Errorf("set explosion transform: %w", err)
	}
	if err := world.SetSpriteRenderer(explosion, ecscomponents.SpriteRenderer{SpriteID: cfg.ExplosionSpriteID, ZIndex: cfg.ExplosionZIndex}); err != nil {
		return fmt.Errorf("set explosion sprite: %w", err)
	}
	if err := world.SetLifetime(explosion, ecscomponents.Lifetime{RemainingSeconds: cfg.ExplosionLifetime}); err != nil {
		return fmt.Errorf("set explosion lifetime: %w", err)
	}

	if cfg.ExplosionLightRadius > 0 && cfg.ExplosionLightIntensity > 0 {
		if err := world.SetPointLight(explosion, ecscomponents.PointLight{
			Radius:    cfg.ExplosionLightRadius,
			Intensity: cfg.ExplosionLightIntensity,
			Color:     cfg.ExplosionLightColor,
		}); err != nil {
			return fmt.Errorf("set explosion light: %w", err)
		}
	}

	if cfg.ExplosionBurstSpriteID != "" && cfg.ExplosionBurstLifetime > 0 {
		burst := world.CreateEntity()
		if err := world.SetTransform(burst, ecscomponents.Transform{X: transform.X, Y: transform.Y, RotationRadians: transform.RotationRadians}); err != nil {
			return fmt.Errorf("set explosion burst transform: %w", err)
		}
		if err := world.SetSpriteRenderer(burst, ecscomponents.SpriteRenderer{SpriteID: cfg.ExplosionBurstSpriteID, ZIndex: cfg.ExplosionZIndex + 1}); err != nil {
			return fmt.Errorf("set explosion burst sprite: %w", err)
		}
		if err := world.SetLifetime(burst, ecscomponents.Lifetime{RemainingSeconds: cfg.ExplosionBurstLifetime}); err != nil {
			return fmt.Errorf("set explosion burst lifetime: %w", err)
		}

		if cfg.BurstLightRadius > 0 && cfg.BurstLightIntensity > 0 {
			if err := world.SetPointLight(burst, ecscomponents.PointLight{
				Radius:    cfg.BurstLightRadius,
				Intensity: cfg.BurstLightIntensity,
				Color:     cfg.ExplosionLightColor,
			}); err != nil {
				return fmt.Errorf("set explosion burst light: %w", err)
			}
		}
	}

	if cfg.ExplosionSmokeSpriteID == "" || cfg.ExplosionSmokeLifetime <= 0 {
		if cfg.ExplosionPlumeSpriteID == "" || cfg.ExplosionPlumeLifetime <= 0 {
			return nil
		}
		return spawnExplosionPlume(world, transform, cfg)
	}

	rotationOffset := 0.0
	if cfg.ExplosionSmokeRotationPI {
		rotationOffset = math.Pi
	}

	smoke := world.CreateEntity()
	if err := world.SetTransform(smoke, ecscomponents.Transform{X: transform.X, Y: transform.Y, RotationRadians: transform.RotationRadians}); err != nil {
		return fmt.Errorf("set explosion smoke transform: %w", err)
	}
	if err := world.SetSpriteRenderer(smoke, ecscomponents.SpriteRenderer{SpriteID: cfg.ExplosionSmokeSpriteID, ZIndex: cfg.ExplosionZIndex, RotationOffsetRadians: rotationOffset}); err != nil {
		return fmt.Errorf("set explosion smoke sprite: %w", err)
	}
	if err := world.SetLifetime(smoke, ecscomponents.Lifetime{RemainingSeconds: cfg.ExplosionSmokeLifetime}); err != nil {
		return fmt.Errorf("set explosion smoke lifetime: %w", err)
	}

	if cfg.ExplosionPlumeSpriteID != "" && cfg.ExplosionPlumeLifetime > 0 {
		if err := spawnExplosionPlume(world, transform, cfg); err != nil {
			return err
		}
	}

	return nil
}

func spawnExplosionPlume(world *ecs.World, transform ecscomponents.Transform, cfg TankDestructionConfig) error {
	plume := world.CreateEntity()
	if err := world.SetTransform(plume, ecscomponents.Transform{X: transform.X, Y: transform.Y, RotationRadians: transform.RotationRadians}); err != nil {
		return fmt.Errorf("set explosion plume transform: %w", err)
	}
	if err := world.SetSpriteRenderer(plume, ecscomponents.SpriteRenderer{SpriteID: cfg.ExplosionPlumeSpriteID, ZIndex: cfg.ExplosionZIndex}); err != nil {
		return fmt.Errorf("set explosion plume sprite: %w", err)
	}
	if err := world.SetLifetime(plume, ecscomponents.Lifetime{RemainingSeconds: cfg.ExplosionPlumeLifetime}); err != nil {
		return fmt.Errorf("set explosion plume lifetime: %w", err)
	}

	return nil
}

func spawnEnemyWreck(world *ecs.World, transform ecscomponents.Transform, cfg TankDestructionConfig) error {
	if cfg.EnemyWreckSpriteID == "" {
		return nil
	}

	wreck := world.CreateEntity()
	if err := world.SetTransform(wreck, ecscomponents.Transform{X: transform.X, Y: transform.Y, RotationRadians: transform.RotationRadians}); err != nil {
		return fmt.Errorf("set enemy wreck transform: %w", err)
	}
	if err := world.SetSpriteRenderer(wreck, ecscomponents.SpriteRenderer{SpriteID: cfg.EnemyWreckSpriteID, ZIndex: cfg.EnemyWreckZIndex}); err != nil {
		return fmt.Errorf("set enemy wreck sprite: %w", err)
	}

	if cfg.EnemyWreckColliderRadius > 0 {
		if err := world.SetCircleCollider(wreck, ecscomponents.CircleCollider{Radius: cfg.EnemyWreckColliderRadius}); err != nil {
			return fmt.Errorf("set enemy wreck collider: %w", err)
		}
		if err := world.SetSolidObstacle(wreck); err != nil {
			return fmt.Errorf("set enemy wreck solid obstacle: %w", err)
		}
	}

	return nil
}
