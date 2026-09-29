package sim

import (
	"math"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestUpdateWeaponsAndSpawnProjectiles(t *testing.T) {
	w := ecs.NewWorld()
	tank := w.CreateEntity()

	must(t, w.SetTransform(tank, ecscomponents.Transform{X: 10, Y: 20, RotationRadians: 0}))
	must(t, w.SetTankDrive(tank, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetTankControlIntent(tank, ecscomponents.TankControlIntent{Fire: true}))
	must(t, w.SetWeapon(tank, ecscomponents.Weapon{
		CooldownSeconds:           0.5,
		ProjectileSpeed:           100,
		ProjectileLifetime:        1,
		ProjectileDamage:          5,
		MuzzleOffset:              12,
		ProjectileSpriteID:        "p",
		ProjectileRadius:          2,
		ProjectileLightRadius:     20,
		ProjectileLightIntensity:  0.4,
		ProjectileLightColor:      ecscomponents.LightColor{R: 1, G: 0.5, B: 0.3},
		MuzzleFlashSpriteID:       "flash",
		MuzzleFlashLifetime:       0.08,
		MuzzleFlashOffset:         6,
		MuzzleFlashLightRadius:    24,
		MuzzleFlashLightIntensity: 0.6,
		MuzzleFlashLightColor:     ecscomponents.LightColor{R: 1, G: 0.7, B: 0.4},
	}))
	must(t, w.SetAudioPlayer(tank, ecscomponents.AudioPlayer{ShootSoundID: "cannon_fire", Volume: 0.6}))

	if err := UpdateWeaponsAndSpawnProjectiles(w, 0.1, 99, 100); err != nil {
		t.Fatalf("update weapons: %v", err)
	}

	projectiles := w.ProjectileStates()
	if len(projectiles) != 1 {
		t.Fatalf("expected 1 projectile, got %d", len(projectiles))
	}

	p := projectiles[0]
	if p.Transform.Y >= 20 {
		t.Fatalf("expected projectile spawned forward (upward), got y=%f", p.Transform.Y)
	}

	renderables := w.Renderables()
	if len(renderables) != 2 {
		t.Fatalf("expected projectile + muzzle flash renderables, got %d", len(renderables))
	}

	pointLights := w.PointLightSources()
	if len(pointLights) != 2 {
		t.Fatalf("expected projectile + muzzle flash lights, got %d", len(pointLights))
	}

	projectileHasLight := false
	for _, source := range pointLights {
		if source.Entity == p.Entity {
			projectileHasLight = true
			break
		}
	}
	if !projectileHasLight {
		t.Fatalf("expected projectile to have point light")
	}

	lifetimes := w.LifetimeStates()
	if len(lifetimes) != 1 {
		t.Fatalf("expected exactly 1 lifetime entity for muzzle flash, got %d", len(lifetimes))
	}

	for entity := range lifetimes {
		muzzleFlashTransform, ok := w.Transform(entity)
		if !ok {
			t.Fatalf("expected muzzle flash transform")
		}

		if muzzleFlashTransform.Y >= p.Transform.Y {
			t.Fatalf("expected muzzle flash in front of projectile spawn, got flashY=%f projectileY=%f", muzzleFlashTransform.Y, p.Transform.Y)
		}

		muzzleFlashHasLight := false
		for _, source := range pointLights {
			if source.Entity == entity {
				muzzleFlashHasLight = true
				break
			}
		}
		if !muzzleFlashHasLight {
			t.Fatalf("expected muzzle flash to carry the transient point light")
		}
	}

	weapon, ok := w.Weapon(tank)
	if !ok {
		t.Fatalf("expected weapon")
	}
	if weapon.CooldownRemaining <= 0 {
		t.Fatalf("expected weapon cooldown to be set after firing")
	}

	intents := w.SoundIntents()
	if len(intents) != 1 {
		t.Fatalf("expected one play-sound intent, got %d", len(intents))
	}
	if intents[0].Entity != tank {
		t.Fatalf("unexpected sound intent entity: got %d want %d", intents[0].Entity, tank)
	}
	if intents[0].Intent.SoundID != "cannon_fire" {
		t.Fatalf("unexpected sound ID: %q", intents[0].Intent.SoundID)
	}
}

func TestApplyProximitySoundVolume(t *testing.T) {
	w := ecs.NewWorld()

	listener := w.CreateEntity()
	must(t, w.SetTransform(listener, ecscomponents.Transform{X: 0, Y: 0}))

	near := w.CreateEntity()
	must(t, w.SetTransform(near, ecscomponents.Transform{X: 20, Y: 0}))
	must(t, w.SetPlaySoundIntent(near, ecscomponents.PlaySoundIntent{SoundID: "cannon_fire", Volume: 1}))
	must(t, w.SetSoundVolumeByDistance(near, ecscomponents.SoundVolumeByDistance{MaxDistance: 100, MinFactor: 0.2}))

	far := w.CreateEntity()
	must(t, w.SetTransform(far, ecscomponents.Transform{X: 400, Y: 0}))
	must(t, w.SetPlaySoundIntent(far, ecscomponents.PlaySoundIntent{SoundID: "cannon_fire", Volume: 1}))
	must(t, w.SetSoundVolumeByDistance(far, ecscomponents.SoundVolumeByDistance{MaxDistance: 100, MinFactor: 0.2}))

	if err := ApplyProximitySoundVolume(w, listener); err != nil {
		t.Fatalf("apply proximity sound volume: %v", err)
	}

	nearIntent, ok := w.PlaySoundIntent(near)
	if !ok {
		t.Fatalf("near sound intent missing")
	}
	if math.Abs(nearIntent.Volume-0.8) > 0.0001 {
		t.Fatalf("near volume mismatch: got %f want %f", nearIntent.Volume, 0.8)
	}

	farIntent, ok := w.PlaySoundIntent(far)
	if !ok {
		t.Fatalf("far sound intent missing")
	}
	if math.Abs(farIntent.Volume-0.2) > 0.0001 {
		t.Fatalf("far volume mismatch: got %f want %f", farIntent.Volume, 0.2)
	}
}

func TestUpdateLifetimesExpiresTemporaryEntity(t *testing.T) {
	w := ecs.NewWorld()
	e := w.CreateEntity()

	must(t, w.SetLifetime(e, ecscomponents.Lifetime{RemainingSeconds: 0.05}))

	if err := UpdateLifetimes(w, 0.1); err != nil {
		t.Fatalf("update lifetimes: %v", err)
	}

	if w.IsAlive(e) {
		t.Fatalf("temporary lifetime entity should be destroyed")
	}
}

func TestIntegrateProjectilesExpiresByLifetime(t *testing.T) {
	w := ecs.NewWorld()
	p := w.CreateEntity()

	must(t, w.SetTransform(p, ecscomponents.Transform{X: 0, Y: 0}))
	must(t, w.SetCircleCollider(p, ecscomponents.CircleCollider{Radius: 2}))
	must(t, w.SetProjectile(p, ecscomponents.Projectile{VelocityX: 10, VelocityY: 0, RemainingLifetime: 0.25, Damage: 1}))

	if err := IntegrateProjectiles(w, 0.5); err != nil {
		t.Fatalf("integrate projectiles: %v", err)
	}

	if w.IsAlive(p) {
		t.Fatalf("projectile should expire and be destroyed")
	}
}

func TestResolveProjectileImpactsDamagesTank(t *testing.T) {
	w := ecs.NewWorld()

	owner := w.CreateEntity()
	must(t, w.SetTransform(owner, ecscomponents.Transform{X: 0, Y: 0}))
	must(t, w.SetTankDrive(owner, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(owner, ecscomponents.CircleCollider{Radius: 8}))
	must(t, w.SetHealth(owner, ecscomponents.Health{Current: 100, Max: 100}))

	target := w.CreateEntity()
	must(t, w.SetTransform(target, ecscomponents.Transform{X: 20, Y: 0}))
	must(t, w.SetTankDrive(target, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(target, ecscomponents.CircleCollider{Radius: 8}))
	must(t, w.SetHealth(target, ecscomponents.Health{Current: 30, Max: 30}))

	projectile := w.CreateEntity()
	must(t, w.SetTransform(projectile, ecscomponents.Transform{X: 20, Y: 0}))
	must(t, w.SetCircleCollider(projectile, ecscomponents.CircleCollider{Radius: 2}))
	must(t, w.SetProjectile(projectile, ecscomponents.Projectile{Owner: owner, Damage: 12, RemainingLifetime: 1}))

	if err := ResolveProjectileImpacts(w); err != nil {
		t.Fatalf("resolve projectile impacts: %v", err)
	}

	if w.IsAlive(projectile) {
		t.Fatalf("projectile should be destroyed on hit")
	}

	health, ok := w.Health(target)
	if !ok {
		t.Fatalf("target health missing")
	}
	if health.Current != 18 {
		t.Fatalf("expected target health 18, got %f", health.Current)
	}
}

func TestResolveProjectileImpactsStopsAtObstacle(t *testing.T) {
	w := ecs.NewWorld()

	owner := w.CreateEntity()
	must(t, w.SetTransform(owner, ecscomponents.Transform{X: 0, Y: 0}))
	must(t, w.SetTankDrive(owner, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(owner, ecscomponents.CircleCollider{Radius: 8}))

	target := w.CreateEntity()
	must(t, w.SetTransform(target, ecscomponents.Transform{X: 40, Y: 0}))
	must(t, w.SetTankDrive(target, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(target, ecscomponents.CircleCollider{Radius: 8}))
	must(t, w.SetHealth(target, ecscomponents.Health{Current: 30, Max: 30}))

	obstacle := w.CreateEntity()
	must(t, w.SetTransform(obstacle, ecscomponents.Transform{X: 20, Y: 0}))
	must(t, w.SetCircleCollider(obstacle, ecscomponents.CircleCollider{Radius: 6}))
	must(t, w.SetSolidObstacle(obstacle))

	projectile := w.CreateEntity()
	must(t, w.SetTransform(projectile, ecscomponents.Transform{X: 20, Y: 0}))
	must(t, w.SetCircleCollider(projectile, ecscomponents.CircleCollider{Radius: 2}))
	must(t, w.SetProjectile(projectile, ecscomponents.Projectile{Owner: owner, Damage: 12, RemainingLifetime: 1}))

	if err := ResolveProjectileImpacts(w); err != nil {
		t.Fatalf("resolve projectile impacts: %v", err)
	}

	health, _ := w.Health(target)
	if health.Current != 30 {
		t.Fatalf("target should not take damage through obstacle, got %f", health.Current)
	}
}

func TestEffectiveProjectileSpreadImprovesWhenStoppedAndRelaxed(t *testing.T) {
	w := ecs.NewWorld()
	tank := w.CreateEntity()

	must(t, w.SetTankDrive(tank, ecscomponents.TankDrive{MoveSpeed: 120, TurnSpeedRadiansPS: 1}))
	must(t, w.SetTankVelocity(tank, ecscomponents.TankVelocity{ForwardUnitsPerSecond: 0}))
	must(t, w.SetEnemyAI(tank, ecscomponents.EnemyAI{Behavior: ecscomponents.EnemyBehaviorCautious}))

	weapon := ecscomponents.Weapon{
		AccuracySpreadRadians:     0.10,
		AccuracyMinSpreadRadians:  0.01,
		AccuracyMovingSpread:      0.12,
		AccuracyStoppedMultiplier: 0.5,
		AccuracyRelaxedMultiplier: 0.7,
	}

	stoppedRelaxedSpread := effectiveProjectileSpread(w, tank, weapon)

	must(t, w.SetEnemyAI(tank, ecscomponents.EnemyAI{Behavior: ecscomponents.EnemyBehaviorAggressive}))
	must(t, w.SetTankVelocity(tank, ecscomponents.TankVelocity{ForwardUnitsPerSecond: 120}))
	movingAggressiveSpread := effectiveProjectileSpread(w, tank, weapon)

	if stoppedRelaxedSpread >= movingAggressiveSpread {
		t.Fatalf("expected stopped+relaxed spread to be tighter: stoppedRelaxed=%f movingAggressive=%f", stoppedRelaxedSpread, movingAggressiveSpread)
	}
}

func TestUpdateWeaponsAndSpawnProjectilesAppliesDeterministicSpreadOffset(t *testing.T) {
	w := ecs.NewWorld()
	tank := w.CreateEntity()

	must(t, w.SetTransform(tank, ecscomponents.Transform{X: 10, Y: 20, RotationRadians: 0}))
	must(t, w.SetTankDrive(tank, ecscomponents.TankDrive{MoveSpeed: 100, TurnSpeedRadiansPS: 1}))
	must(t, w.SetTankVelocity(tank, ecscomponents.TankVelocity{ForwardUnitsPerSecond: 0}))
	must(t, w.SetTankControlIntent(tank, ecscomponents.TankControlIntent{Fire: true}))

	weapon := ecscomponents.Weapon{
		CooldownSeconds:           0.5,
		ProjectileSpeed:           100,
		ProjectileLifetime:        1,
		ProjectileDamage:          5,
		MuzzleOffset:              12,
		AccuracySpreadRadians:     0.12,
		AccuracyMinSpreadRadians:  0.12,
		AccuracyStoppedMultiplier: 1,
		ProjectileSpriteID:        "p",
		ProjectileRadius:          2,
	}
	must(t, w.SetWeapon(tank, weapon))

	if err := UpdateWeaponsAndSpawnProjectiles(w, 0.1, 99, 100); err != nil {
		t.Fatalf("update weapons: %v", err)
	}

	projectiles := w.ProjectileStates()
	if len(projectiles) != 1 {
		t.Fatalf("expected 1 projectile, got %d", len(projectiles))
	}

	spread := effectiveProjectileSpread(w, tank, weapon)
	expectedRotation := deterministicSpreadOffset(tank, 0, spread)
	approxEqual(t, projectiles[0].Transform.RotationRadians, expectedRotation, 1e-9)
}
