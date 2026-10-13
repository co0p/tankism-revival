package sim

import (
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"testing"
)

func TestWeaponSpawnSpecDependencies(t *testing.T) {
	assertSystemSpec(t, weaponSpawnSpec(), SystemSpec{Name: "weapon_cooldowns_and_projectile_spawn", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindTankDrive, ecscomponents.KindTankVelocity, ecscomponents.KindTankControlIntent, ecscomponents.KindWeapon, ecscomponents.KindEnemyAI, ecscomponents.KindAudioPlayer}, Writes: []ecscomponents.Kind{ecscomponents.KindWeapon, ecscomponents.KindTransform, ecscomponents.KindSpriteRenderer, ecscomponents.KindCircleCollider, ecscomponents.KindProjectile, ecscomponents.KindLifetime, ecscomponents.KindPointLight, ecscomponents.KindPlaySoundIntent}})
}

func TestProjectileIntegrationSpecDependencies(t *testing.T) {
	assertSystemSpec(t, projectileIntegrationSpec(), SystemSpec{Name: "projectile_integration", Reads: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindProjectile, ecscomponents.KindCircleCollider, ecscomponents.KindSolidObstacle, ecscomponents.KindTankDrive}, Writes: []ecscomponents.Kind{ecscomponents.KindTransform, ecscomponents.KindProjectile, ecscomponents.KindGroundImpacts}})
}

func TestProjectileImpactsSpecDependencies(t *testing.T) {
	assertSystemSpec(t, projectileImpactsSpec(), SystemSpec{Name: "projectile_impacts", Reads: []ecscomponents.Kind{ecscomponents.KindProjectile, ecscomponents.KindTransform, ecscomponents.KindCircleCollider, ecscomponents.KindSolidObstacle, ecscomponents.KindTankDrive, ecscomponents.KindHealth}, Writes: []ecscomponents.Kind{ecscomponents.KindHealth, ecscomponents.KindProjectile, ecscomponents.KindHitFlash, ecscomponents.KindSparkParticles}})
}

func TestLifetimeSpecDependencies(t *testing.T) {
	assertSystemSpec(t, lifetimeSpec(), SystemSpec{Name: "lifetime_expiration", Reads: []ecscomponents.Kind{ecscomponents.KindLifetime}, Writes: []ecscomponents.Kind{ecscomponents.KindLifetime}})
}
