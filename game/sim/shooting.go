package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func UpdateWeaponsAndSpawnProjectiles(world *ecs.World, deltaSeconds float64, projectileZIndex, muzzleFlashZIndex int) error {
	for _, tank := range world.ArmedTanks() {
		weapon := tank.Weapon

		if weapon.CooldownRemaining > 0 {
			weapon.CooldownRemaining -= deltaSeconds
			if weapon.CooldownRemaining < 0 {
				weapon.CooldownRemaining = 0
			}
		}

		shouldFire := tank.Intent.Fire && weapon.CooldownRemaining <= 0
		if shouldFire && canSpawnProjectile(weapon) {
			forwardX, forwardY := forwardVector(tank.Transform.RotationRadians)

			spawnX := tank.Transform.X + forwardX*weapon.MuzzleOffset
			spawnY := tank.Transform.Y + forwardY*weapon.MuzzleOffset

			projectile := world.CreateEntity()

			if err := world.SetTransform(projectile, ecscomponents.Transform{X: spawnX, Y: spawnY, RotationRadians: tank.Transform.RotationRadians}); err != nil {
				return fmt.Errorf("set projectile transform: %w", err)
			}

			if err := world.SetSpriteRenderer(projectile, ecscomponents.SpriteRenderer{
				SpriteID:              weapon.ProjectileSpriteID,
				ZIndex:                projectileZIndex,
				RotationOffsetRadians: 0,
			}); err != nil {
				return fmt.Errorf("set projectile sprite: %w", err)
			}

			if err := world.SetCircleCollider(projectile, ecscomponents.CircleCollider{Radius: weapon.ProjectileRadius}); err != nil {
				return fmt.Errorf("set projectile collider: %w", err)
			}

			if err := world.SetProjectile(projectile, ecscomponents.Projectile{
				Owner:             tank.Entity,
				VelocityX:         forwardX * weapon.ProjectileSpeed,
				VelocityY:         forwardY * weapon.ProjectileSpeed,
				RemainingLifetime: weapon.ProjectileLifetime,
				Damage:            weapon.ProjectileDamage,
			}); err != nil {
				return fmt.Errorf("set projectile state: %w", err)
			}

			if weapon.ProjectileLightRadius > 0 && weapon.ProjectileLightIntensity > 0 {
				if err := world.SetPointLight(projectile, ecscomponents.PointLight{
					Radius:    weapon.ProjectileLightRadius,
					Intensity: weapon.ProjectileLightIntensity,
					Color:     weapon.ProjectileLightColor,
				}); err != nil {
					return fmt.Errorf("set projectile light: %w", err)
				}
			}

			if shouldSpawnMuzzleFlash(weapon) {
				muzzleFlashX := spawnX + forwardX*weapon.MuzzleFlashOffset
				muzzleFlashY := spawnY + forwardY*weapon.MuzzleFlashOffset

				muzzleFlash := world.CreateEntity()
				if err := world.SetTransform(muzzleFlash, ecscomponents.Transform{X: muzzleFlashX, Y: muzzleFlashY, RotationRadians: tank.Transform.RotationRadians}); err != nil {
					return fmt.Errorf("set muzzle flash transform: %w", err)
				}

				if weapon.MuzzleFlashSpriteID != "" {
					if err := world.SetSpriteRenderer(muzzleFlash, ecscomponents.SpriteRenderer{
						SpriteID:              weapon.MuzzleFlashSpriteID,
						ZIndex:                muzzleFlashZIndex,
						RotationOffsetRadians: math.Pi,
					}); err != nil {
						return fmt.Errorf("set muzzle flash sprite: %w", err)
					}
				}

				if err := world.SetLifetime(muzzleFlash, ecscomponents.Lifetime{RemainingSeconds: weapon.MuzzleFlashLifetime}); err != nil {
					return fmt.Errorf("set muzzle flash lifetime: %w", err)
				}

				if weapon.MuzzleFlashLightRadius > 0 && weapon.MuzzleFlashLightIntensity > 0 {
					if err := world.SetPointLight(muzzleFlash, ecscomponents.PointLight{
						Radius:    weapon.MuzzleFlashLightRadius,
						Intensity: weapon.MuzzleFlashLightIntensity,
						Color:     weapon.MuzzleFlashLightColor,
					}); err != nil {
						return fmt.Errorf("set muzzle flash light: %w", err)
					}
				}
			}

			weapon.CooldownRemaining = weapon.CooldownSeconds

			if audioPlayer, ok := world.AudioPlayer(tank.Entity); ok && audioPlayer.ShootSoundID != "" {
				soundVolume := audioPlayer.Volume
				if soundVolume <= 0 {
					soundVolume = 1
				}
				if err := world.SetPlaySoundIntent(tank.Entity, ecscomponents.PlaySoundIntent{
					SoundID: audioPlayer.ShootSoundID,
					Volume:  soundVolume,
				}); err != nil {
					return fmt.Errorf("set play sound intent: %w", err)
				}
			}
		}

		if err := world.SetWeapon(tank.Entity, weapon); err != nil {
			return fmt.Errorf("update weapon cooldown: %w", err)
		}
	}

	return nil
}

func UpdateLifetimes(world *ecs.World, deltaSeconds float64) error {
	for entity, lifetime := range world.LifetimeStates() {
		lifetime.RemainingSeconds -= deltaSeconds
		if lifetime.RemainingSeconds <= 0 {
			world.DestroyEntity(entity)
			continue
		}

		if err := world.SetLifetime(entity, lifetime); err != nil {
			return fmt.Errorf("set lifetime: %w", err)
		}
	}

	return nil
}

func IntegrateProjectiles(world *ecs.World, deltaSeconds float64) error {
	for _, state := range world.ProjectileStates() {
		transform := state.Transform
		projectile := state.Projectile

		transform.X += projectile.VelocityX * deltaSeconds
		transform.Y += projectile.VelocityY * deltaSeconds
		projectile.RemainingLifetime -= deltaSeconds

		if err := world.SetTransform(state.Entity, transform); err != nil {
			return fmt.Errorf("set projectile transform: %w", err)
		}

		if projectile.RemainingLifetime <= 0 {
			world.DestroyEntity(state.Entity)
			continue
		}

		if err := world.SetProjectile(state.Entity, projectile); err != nil {
			return fmt.Errorf("set projectile state: %w", err)
		}
	}

	return nil
}

func ResolveProjectileImpacts(world *ecs.World) error {
	projectiles := world.ProjectileStates()
	obstacles := world.SolidObstacleColliderBodies()
	tanks := world.TankColliderBodies()

	for _, projectile := range projectiles {
		if hitsObstacle(projectile, obstacles) {
			world.DestroyEntity(projectile.Entity)
			continue
		}

		for _, tank := range tanks {
			if tank.Entity == projectile.Projectile.Owner {
				continue
			}

			if !circlesOverlap(
				projectile.Transform.X,
				projectile.Transform.Y,
				projectile.Collider.Radius,
				tank.Transform.X,
				tank.Transform.Y,
				tank.Collider.Radius,
			) {
				continue
			}

			health, ok := world.Health(tank.Entity)
			if ok {
				health.Current -= projectile.Projectile.Damage
				if health.Current < 0 {
					health.Current = 0
				}
				if err := world.SetHealth(tank.Entity, health); err != nil {
					return fmt.Errorf("set tank health: %w", err)
				}
			}

			world.DestroyEntity(projectile.Entity)
			break
		}
	}

	return nil
}

func canSpawnProjectile(weapon ecscomponents.Weapon) bool {
	return weapon.ProjectileSpeed > 0 &&
		weapon.ProjectileLifetime > 0 &&
		weapon.ProjectileDamage > 0 &&
		weapon.ProjectileRadius > 0 &&
		weapon.ProjectileSpriteID != ""
}

func shouldSpawnMuzzleFlash(weapon ecscomponents.Weapon) bool {
	if weapon.MuzzleFlashLifetime <= 0 {
		return false
	}

	hasVisual := weapon.MuzzleFlashSpriteID != ""
	hasLight := weapon.MuzzleFlashLightRadius > 0 && weapon.MuzzleFlashLightIntensity > 0

	return hasVisual || hasLight
}

func hitsObstacle(projectile ecs.ProjectileState, obstacles []ecs.ColliderBody) bool {
	for _, obstacle := range obstacles {
		if circlesOverlap(
			projectile.Transform.X,
			projectile.Transform.Y,
			projectile.Collider.Radius,
			obstacle.Transform.X,
			obstacle.Transform.Y,
			obstacle.Collider.Radius,
		) {
			return true
		}
	}

	return false
}

func circlesOverlap(ax, ay, ar, bx, by, br float64) bool {
	dx := ax - bx
	dy := ay - by
	r := ar + br

	return dx*dx+dy*dy <= r*r
}
