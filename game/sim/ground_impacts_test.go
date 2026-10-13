package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"reflect"
	"testing"
)

func TestProjectileLandsAtLifetimeEndpoint(t *testing.T) {
	world := ecs.NewWorld()
	projectile := hitFeedbackProjectile(t, world, 0, 6)
	must(t, world.SetProjectile(projectile, ecscomponents.Projectile{VelocityX: 100, RemainingLifetime: 0.1, Damage: 6}))
	ground := NewGroundImpacts(8)
	must(t, integrateProjectiles(world, 0.5, ground))
	if world.IsAlive(projectile) || len(ground.Splashes()) != 1 || len(ground.Dust()) != GroundDustCount {
		t.Fatal("expired projectile should land once with dust")
	}
	splash := ground.Splashes()[0]
	if splash.X != 110 || splash.Y != 100 || splash.Age != 0 {
		t.Fatalf("projectile overshot endpoint: %v", splash)
	}
	must(t, integrateProjectiles(world, 0.5, ground))
	if len(ground.Splashes()) != 1 {
		t.Fatal("consumed projectile landed twice")
	}
}

func TestFinalTickCollisionPrecedesGroundImpact(t *testing.T) {
	for _, cover := range []bool{true, false} {
		world := ecs.NewWorld()
		tank := hitFeedbackTank(t, world)
		if cover {
			must(t, world.SetSolidObstacle(tank))
		}
		mapEntity := world.CreateEntity()
		must(t, world.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: 500, Height: 500}))
		projectile := hitFeedbackProjectile(t, world, 0, 6)
		must(t, world.SetProjectile(projectile, ecscomponents.Projectile{RemainingLifetime: 0.01, Damage: 6}))
		ground, sparks := NewGroundImpacts(8), NewSparkParticles()
		must(t, RunDefaultPipeline(world, 0.1, PipelineConfig{PlayerEntity: tank, MapEntity: mapEntity, GroundImpacts: ground, Sparks: sparks}))
		if world.IsAlive(projectile) || len(ground.Splashes()) != 0 || len(sparks.Particles()) != ImpactSparkCount {
			t.Fatal("final tick collision should create sparks, not ground splash")
		}
		health, _ := world.Health(tank)
		if !cover && health.Current != 94 {
			t.Fatal("final tick tank collision should deal damage")
		}
	}
}

func TestGroundImpactCapsDeterminismAndExpiration(t *testing.T) {
	a, b := NewGroundImpacts(8), NewGroundImpacts(8)
	for i := 0; i < 200; i++ {
		projectile := ecs.ProjectileState{Entity: ecs.Entity(i + 1), Transform: ecscomponents.Transform{X: float64(i)}}
		a.Burst(projectile)
		b.Burst(projectile)
	}
	if len(a.Splashes()) != MaxGroundSplashes || len(a.Dust()) != MaxDustParticles || !reflect.DeepEqual(a, b) {
		t.Fatal("ground effects should be capped and deterministic")
	}
	if a.Splashes()[0].X != 72 {
		t.Fatal("mark cap should replace oldest")
	}
	before := a.Dust()[0]
	a.Advance(0.1)
	b.Advance(0.1)
	after := a.Dust()[0]
	if after.X == before.X && after.Y == before.Y {
		t.Fatal("dust should scatter")
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("dust motion should be deterministic")
	}
	a.Advance(1)
	if len(a.Dust()) != 0 || len(a.Splashes()) == 0 {
		t.Fatal("dust should expire before ground marks")
	}
	a.Advance(GroundSplashLifetime)
	if len(a.Splashes()) != 0 {
		t.Fatal("marks should expire offscreen too")
	}
}

func TestGroundImpactSpawnsAfterEffectAdvance(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	mapEntity := world.CreateEntity()
	must(t, world.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: 500, Height: 500}))
	projectile := hitFeedbackProjectile(t, world, tank, 6)
	ground := NewGroundImpacts(8)
	must(t, RunDefaultPipeline(world, 1, PipelineConfig{PlayerEntity: tank, MapEntity: mapEntity, GroundImpacts: ground}))
	if world.IsAlive(projectile) || len(ground.Dust()) != GroundDustCount || ground.Dust()[0].Age != 0 {
		t.Fatal("new dust should have full lifetime; owner overlap must not block landing")
	}
}
