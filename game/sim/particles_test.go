package sim

import (
	"reflect"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func smokeTestWorld(t *testing.T, health float64) (*ecs.World, ecs.Entity) {
	t.Helper()
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	must(t, world.SetHealth(tank, ecscomponents.Health{Current: health, Max: 100}))
	return world, tank
}

func TestSmokeHealthScaling(t *testing.T) {
	counts := make(map[float64]int)
	for _, health := range []float64{100, 75, 60, 10, 0} {
		world, _ := smokeTestWorld(t, health)
		smoke := NewSmokeParticles(10)
		for tick := 0; tick < 60; tick++ {
			smoke.Update(world, 1.0/60)
		}
		counts[health] = len(smoke.Particles())
	}
	if counts[100] != 0 || counts[75] != 0 || counts[0] != 0 {
		t.Fatal("healthy and dead tanks should not emit smoke")
	}
	if counts[60] == 0 || counts[10] <= counts[60] {
		t.Fatalf("lower health should emit more smoke: %v", counts)
	}
}

func TestSmokeDriftGrowthAndExpiration(t *testing.T) {
	world, tank := smokeTestWorld(t, 10)
	smoke := NewSmokeParticles(10)
	smoke.Update(world, 0.1)
	before := smoke.Particles()[0]
	if before.EndSize <= before.StartSize || before.Variant < 0 || before.Variant >= 10 {
		t.Fatal("expected growing puff with a valid variant")
	}
	world.DestroyEntity(tank)
	smoke.Update(world, 0.1)
	after := smoke.Particles()[0]
	if after.Y >= before.Y || after.Age <= before.Age || after.X == before.X {
		t.Fatal("puffs should drift and age independently of their tank")
	}
	if len(smoke.emission) != 0 {
		t.Fatal("destroyed tank retained emitter bookkeeping")
	}
	smoke.Advance(3)
	if len(smoke.Particles()) != 0 {
		t.Fatal("puffs did not expire")
	}
}

func TestSmokeStopsOnHealing(t *testing.T) {
	world, tank := smokeTestWorld(t, 10)
	smoke := NewSmokeParticles(10)
	smoke.Update(world, 0.1)
	count := len(smoke.Particles())
	must(t, world.SetHealth(tank, ecscomponents.Health{Current: 100, Max: 100}))
	smoke.Update(world, 0.1)
	if len(smoke.Particles()) != count || len(smoke.emission) != 0 {
		t.Fatal("healed tank should stop emitting without clearing existing puffs")
	}
}

func TestSmokeCapAndDeterminism(t *testing.T) {
	world := ecs.NewWorld()
	for i := 0; i < 40; i++ {
		tank := hitFeedbackTank(t, world)
		must(t, world.SetHealth(tank, ecscomponents.Health{Current: 1, Max: 100}))
	}
	a, b := NewSmokeParticles(10), NewSmokeParticles(10)
	reachedCap := false
	for tick := 0; tick < 180; tick++ {
		a.Update(world, 1.0/60)
		b.Update(world, 1.0/60)
		if len(a.particles) > MaxSmokeParticles {
			t.Fatal("particle cap exceeded")
		}
		reachedCap = reachedCap || len(a.particles) == MaxSmokeParticles
		if !reflect.DeepEqual(a, b) {
			t.Fatal("identical updates produced different smoke")
		}
	}
	if !reachedCap {
		t.Fatal("stress world should reach cap")
	}
	snapshot := a.Particles()
	snapshot[0].X = -999
	if a.particles[0].X == -999 {
		t.Fatal("particle query must not expose mutable storage")
	}
}

func TestPipelineEmitsSmokeAfterDamage(t *testing.T) {
	world, tank := smokeTestWorld(t, 76)
	mapEntity := world.CreateEntity()
	must(t, world.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: 500, Height: 500}))
	hitFeedbackProjectile(t, world, 0, 50)
	smoke := NewSmokeParticles(10)
	must(t, RunDefaultPipeline(world, 0.1, PipelineConfig{PlayerEntity: tank, MapEntity: mapEntity, Smoke: smoke}))
	if len(smoke.Particles()) == 0 {
		t.Fatal("pipeline must emit smoke using health after projectile damage")
	}
}
