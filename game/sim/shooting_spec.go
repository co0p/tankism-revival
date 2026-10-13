package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func weaponSpawnSystem() runtimeSystem {
	return runtimeSystem{spec: weaponSpawnSpec(), run: func(world *ecs.World, dt float64, cfg PipelineConfig) error {
		return UpdateWeaponsAndSpawnProjectiles(world, dt, cfg.ProjectileZIndex, cfg.MuzzleFlashZIndex)
	}}
}

func projectileIntegrationSystem() runtimeSystem {
	return runtimeSystem{spec: projectileIntegrationSpec(), run: func(world *ecs.World, dt float64, cfg PipelineConfig) error {
		return integrateProjectiles(world, dt, cfg.GroundImpacts)
	}}
}

func projectileImpactsSystem() runtimeSystem {
	return runtimeSystem{spec: projectileImpactsSpec(), run: func(world *ecs.World, _ float64, cfg PipelineConfig) error {
		return resolveProjectileImpacts(world, cfg.Sparks)
	}}
}

func lifetimeSystem() runtimeSystem {
	return runtimeSystem{spec: lifetimeSpec(), run: func(world *ecs.World, dt float64, _ PipelineConfig) error {
		return UpdateLifetimes(world, dt)
	}}
}

func weaponSpawnSpec() SystemSpec {
	return SystemSpec{Name: "weapon_cooldowns_and_projectile_spawn", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindTankVelocity, ecscomponents.KindTankControlIntent, ecscomponents.KindWeapon, ecscomponents.KindEnemyAI, ecscomponents.KindAudioPlayer}, Writes: []ecscomponents.Kind{ecscomponents.KindWeapon, ecscomponents.KindTransform, ecscomponents.KindSpriteRenderer, ecscomponents.KindCircleCollider, ecscomponents.KindProjectile, ecscomponents.KindLifetime, ecscomponents.KindPointLight, ecscomponents.KindPlaySoundIntent}}
}

func projectileIntegrationSpec() SystemSpec {
	return SystemSpec{Name: "projectile_integration", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindProjectile, ecscomponents.KindCircleCollider, ecscomponents.KindSolidObstacle, ecscomponents.KindTankDrive}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindProjectile, ecscomponents.KindGroundImpacts}}
}

func projectileImpactsSpec() SystemSpec {
	return SystemSpec{Name: "projectile_impacts", Reads: []ecscomponents.Kind{ecscomponents.KindProjectile, ecscomponents.KindTransform, ecscomponents.KindCircleCollider, ecscomponents.KindSolidObstacle, ecscomponents.KindTankDrive, ecscomponents.KindHealth}, Writes: []ecscomponents.Kind{ecscomponents.KindHealth, ecscomponents.KindProjectile, ecscomponents.KindHitFlash, ecscomponents.KindSparkParticles}}
}

func lifetimeSpec() SystemSpec {
	return SystemSpec{Name: "lifetime_expiration", Reads: []ecscomponents.Kind{ecscomponents.KindLifetime}, Writes: []ecscomponents.Kind{ecscomponents.KindLifetime}}
}
