package sim

import (
	"math"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func hitFeedbackTank(t *testing.T, world *ecs.World) ecs.Entity {
	t.Helper()
	tank := world.CreateEntity()
	must(t, world.SetTransform(tank, ecscomponents.Transform{X: 100, Y: 100}))
	must(t, world.SetTankDrive(tank, ecscomponents.TankDrive{}))
	must(t, world.SetCircleCollider(tank, ecscomponents.CircleCollider{Radius: 16}))
	must(t, world.SetHealth(tank, ecscomponents.Health{Current: 100, Max: 100}))
	return tank
}

func hitFeedbackProjectile(t *testing.T, world *ecs.World, owner ecs.Entity, damage float64) ecs.Entity {
	t.Helper()
	projectile := world.CreateEntity()
	must(t, world.SetTransform(projectile, ecscomponents.Transform{X: 100, Y: 100}))
	must(t, world.SetCircleCollider(projectile, ecscomponents.CircleCollider{Radius: 4}))
	must(t, world.SetProjectile(projectile, ecscomponents.Projectile{Owner: owner, Damage: damage, RemainingLifetime: 1}))
	return projectile
}

func TestHitFlashOnlyOnDamagingTankImpact(t *testing.T) {
	for _, test := range []struct {
		name     string
		damage   float64
		selfShot bool
		obstacle bool
		dead     bool
		wantHit  bool
	}{
		{name: "damage", damage: 6, wantHit: true},
		{name: "lethal damage", damage: 100, wantHit: true},
		{name: "zero damage"},
		{name: "negative damage", damage: -6},
		{name: "owner ignored", damage: 6, selfShot: true},
		{name: "obstacle blocks", damage: 6, obstacle: true},
		{name: "dead tank", damage: 6, dead: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			world := ecs.NewWorld()
			tank := hitFeedbackTank(t, world)
			if test.dead {
				must(t, world.SetHealth(tank, ecscomponents.Health{Max: 100}))
			}
			owner := ecs.Entity(0)
			if test.selfShot {
				owner = tank
			}
			if test.obstacle {
				obstacle := world.CreateEntity()
				must(t, world.SetTransform(obstacle, ecscomponents.Transform{X: 100, Y: 100}))
				must(t, world.SetCircleCollider(obstacle, ecscomponents.CircleCollider{Radius: 10}))
				must(t, world.SetSolidObstacle(obstacle))
			}
			projectile := hitFeedbackProjectile(t, world, owner, test.damage)
			must(t, ResolveProjectileImpacts(world))
			flash, hit := world.HitFlash(tank)
			if hit != test.wantHit {
				t.Fatalf("hit flash present = %v, want %v", hit, test.wantHit)
			}
			if test.wantHit {
				if flash.RemainingSeconds != HitFlashDurationSeconds || world.IsAlive(projectile) {
					t.Fatal("expected full-duration flash and consumed projectile")
				}
				health, _ := world.Health(tank)
				must(t, ResolveProjectileImpacts(world))
				after, _ := world.Health(tank)
				if after != health {
					t.Fatal("consumed projectile dealt damage twice")
				}
			}
		})
	}
}

func TestHitFlashRefreshAndExpiration(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	hitFeedbackProjectile(t, world, 0, 6)
	must(t, ResolveProjectileImpacts(world))
	must(t, UpdateHitFlashes(world, HitFlashDurationSeconds/2))
	flash, _ := world.HitFlash(tank)
	if flash.X != 100 || flash.Y != 100 || flash.DurationSeconds != HitFlashDurationSeconds {
		t.Fatal("ring must store the impact position and full duration")
	}
	must(t, world.SetTransform(tank, ecscomponents.Transform{X: 105, Y: 100}))
	if movedFlash, _ := world.HitFlash(tank); movedFlash != flash {
		t.Fatal("ring should stay at the impact point when the tank moves")
	}
	if flash.RemainingSeconds != HitFlashDurationSeconds/2 {
		t.Fatal("flash did not age")
	}
	hitFeedbackProjectile(t, world, 0, 6)
	must(t, ResolveProjectileImpacts(world))
	flash, _ = world.HitFlash(tank)
	if flash.RemainingSeconds != HitFlashDurationSeconds {
		t.Fatal("repeated hit did not refresh flash")
	}
	must(t, UpdateHitFlashes(world, HitFlashDurationSeconds))
	if _, ok := world.HitFlash(tank); ok || !world.IsAlive(tank) {
		t.Fatal("expiration should remove only the flash, not the tank")
	}
}

func TestPipelinePreservesNewHitFlashDuration(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	mapEntity := world.CreateEntity()
	must(t, world.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: 500, Height: 500}))
	hitFeedbackProjectile(t, world, 0, 6)
	cfg := PipelineConfig{PlayerEntity: tank, MapEntity: mapEntity}
	must(t, RunDefaultPipeline(world, 1.0/60, cfg))
	flash, ok := world.HitFlash(tank)
	if !ok || flash.RemainingSeconds != HitFlashDurationSeconds {
		t.Fatal("new hit flash must not age on its spawning tick")
	}
	must(t, RunDefaultPipeline(world, HitFlashDurationSeconds, cfg))
	if _, ok := world.HitFlash(tank); ok {
		t.Fatal("pipeline did not expire flash")
	}
}

func TestPlayerHitShake(t *testing.T) {
	var shake PlayerHitShake
	shake.Update(1.0/60, false)
	if shake.X != 0 || shake.Y != 0 {
		t.Fatal("shake without damage")
	}
	shake.Update(1.0/60, true)
	if shake.X == 0 && shake.Y == 0 {
		t.Fatal("damage should immediately shake the camera")
	}
	shake.Update(PlayerHitShakeDurationSeconds/2, false)
	if math.Abs(shake.X) > PlayerHitShakeAmplitude/2 || math.Abs(shake.Y) > PlayerHitShakeAmplitude/2 {
		t.Fatal("shake did not fade within its amplitude envelope")
	}
	shake.Update(1.0/60, true)
	if shake.RemainingSeconds != PlayerHitShakeDurationSeconds || math.Abs(shake.X) > PlayerHitShakeAmplitude {
		t.Fatal("repeated damage should refresh, not stack shake")
	}
	shake.Update(PlayerHitShakeDurationSeconds, false)
	if shake.RemainingSeconds != 0 || shake.X != 0 || shake.Y != 0 {
		t.Fatal("shake must return to exactly zero")
	}
}

func TestPlayerHitShakeDeterministic(t *testing.T) {
	var a, b PlayerHitShake
	for tick := 0; tick < 60; tick++ {
		a.Update(1.0/60, tick == 0 || tick == 4)
		b.Update(1.0/60, tick == 0 || tick == 4)
		if a != b {
			t.Fatal("same updates produced different shakes")
		}
	}
}
