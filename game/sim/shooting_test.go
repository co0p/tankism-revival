package sim

import (
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
)

func TestUpdateWeaponsAndSpawnProjectiles(t *testing.T) {
	w := ecs.NewWorld()
	tank := w.CreateEntity()

	must(t, w.SetTransform(tank, ecs.Transform{X: 10, Y: 20, RotationRadians: 0}))
	must(t, w.SetTankDrive(tank, ecs.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetTankControlIntent(tank, ecs.TankControlIntent{Fire: true}))
	must(t, w.SetWeapon(tank, ecs.Weapon{
		CooldownSeconds:           0.5,
		ProjectileSpeed:           100,
		ProjectileLifetime:        1,
		ProjectileDamage:          5,
		MuzzleOffset:              12,
		ProjectileSpriteID:        "p",
		ProjectileRadius:          2,
		ProjectileLightRadius:     20,
		ProjectileLightIntensity:  0.4,
		ProjectileLightColor:      ecs.LightColor{R: 1, G: 0.5, B: 0.3},
		MuzzleFlashSpriteID:       "flash",
		MuzzleFlashLifetime:       0.08,
		MuzzleFlashLightRadius:    24,
		MuzzleFlashLightIntensity: 0.6,
		MuzzleFlashLightColor:     ecs.LightColor{R: 1, G: 0.7, B: 0.4},
	}))

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

	weapon, ok := w.Weapon(tank)
	if !ok {
		t.Fatalf("expected weapon")
	}
	if weapon.CooldownRemaining <= 0 {
		t.Fatalf("expected weapon cooldown to be set after firing")
	}
}

func TestUpdateLifetimesExpiresTemporaryEntity(t *testing.T) {
	w := ecs.NewWorld()
	e := w.CreateEntity()

	must(t, w.SetLifetime(e, ecs.Lifetime{RemainingSeconds: 0.05}))

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

	must(t, w.SetTransform(p, ecs.Transform{X: 0, Y: 0}))
	must(t, w.SetCircleCollider(p, ecs.CircleCollider{Radius: 2}))
	must(t, w.SetProjectile(p, ecs.Projectile{VelocityX: 10, VelocityY: 0, RemainingLifetime: 0.25, Damage: 1}))

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
	must(t, w.SetTransform(owner, ecs.Transform{X: 0, Y: 0}))
	must(t, w.SetTankDrive(owner, ecs.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(owner, ecs.CircleCollider{Radius: 8}))
	must(t, w.SetHealth(owner, ecs.Health{Current: 100, Max: 100}))

	target := w.CreateEntity()
	must(t, w.SetTransform(target, ecs.Transform{X: 20, Y: 0}))
	must(t, w.SetTankDrive(target, ecs.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(target, ecs.CircleCollider{Radius: 8}))
	must(t, w.SetHealth(target, ecs.Health{Current: 30, Max: 30}))

	projectile := w.CreateEntity()
	must(t, w.SetTransform(projectile, ecs.Transform{X: 20, Y: 0}))
	must(t, w.SetCircleCollider(projectile, ecs.CircleCollider{Radius: 2}))
	must(t, w.SetProjectile(projectile, ecs.Projectile{Owner: owner, Damage: 12, RemainingLifetime: 1}))

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
	must(t, w.SetTransform(owner, ecs.Transform{X: 0, Y: 0}))
	must(t, w.SetTankDrive(owner, ecs.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(owner, ecs.CircleCollider{Radius: 8}))

	target := w.CreateEntity()
	must(t, w.SetTransform(target, ecs.Transform{X: 40, Y: 0}))
	must(t, w.SetTankDrive(target, ecs.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(target, ecs.CircleCollider{Radius: 8}))
	must(t, w.SetHealth(target, ecs.Health{Current: 30, Max: 30}))

	obstacle := w.CreateEntity()
	must(t, w.SetTransform(obstacle, ecs.Transform{X: 20, Y: 0}))
	must(t, w.SetCircleCollider(obstacle, ecs.CircleCollider{Radius: 6}))
	must(t, w.SetSolidObstacle(obstacle))

	projectile := w.CreateEntity()
	must(t, w.SetTransform(projectile, ecs.Transform{X: 20, Y: 0}))
	must(t, w.SetCircleCollider(projectile, ecs.CircleCollider{Radius: 2}))
	must(t, w.SetProjectile(projectile, ecs.Projectile{Owner: owner, Damage: 12, RemainingLifetime: 1}))

	if err := ResolveProjectileImpacts(w); err != nil {
		t.Fatalf("resolve projectile impacts: %v", err)
	}

	health, _ := w.Health(target)
	if health.Current != 30 {
		t.Fatalf("target should not take damage through obstacle, got %f", health.Current)
	}
}
