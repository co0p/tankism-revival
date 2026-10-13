package sim

import (
	"math"
	"reflect"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestSparksOnProjectileImpacts(t *testing.T) {
	for _, target := range []string{"tank", "cover", "owner", "miss", "expired"} {
		t.Run(target, func(t *testing.T) {
			world := ecs.NewWorld()
			tank := hitFeedbackTank(t, world)
			owner := ecs.Entity(0)
			if target == "owner" {
				owner = tank
			}
			if target == "cover" {
				cover := world.CreateEntity()
				must(t, world.SetTransform(cover, ecscomponents.Transform{X: 100, Y: 100}))
				must(t, world.SetCircleCollider(cover, ecscomponents.CircleCollider{Radius: 12}))
				must(t, world.SetSolidObstacle(cover))
			}
			projectile := hitFeedbackProjectile(t, world, owner, 6)
			if target == "miss" {
				must(t, world.SetTransform(projectile, ecscomponents.Transform{X: 300, Y: 300}))
			}
			if target == "expired" {
				must(t, IntegrateProjectiles(world, 1))
			}
			sparks := NewSparkParticles()
			must(t, resolveProjectileImpacts(world, sparks))
			want := 0
			if target == "tank" || target == "cover" {
				want = ImpactSparkCount
			}
			if len(sparks.Particles()) != want {
				t.Fatalf("got %d sparks, want %d", len(sparks.Particles()), want)
			}
			for _, p := range sparks.Particles() {
				if p.X != 100 || p.Y != 100 || p.Age != 0 {
					t.Fatal("sparks must spawn at the impact point with full lifetime")
				}
			}
			must(t, resolveProjectileImpacts(world, sparks))
			if len(sparks.Particles()) != want {
				t.Fatal("consumed projectile spawned a second burst")
			}
		})
	}
}

func TestSparkMotionCapAndDeterminism(t *testing.T) {
	a, b := NewSparkParticles(), NewSparkParticles()
	for i := 0; i < 100; i++ {
		projectile := ecs.ProjectileState{Entity: ecs.Entity(i + 1)}
		a.Burst(projectile)
		b.Burst(projectile)
	}
	if len(a.Particles()) != MaxSparkParticles || !reflect.DeepEqual(a, b) {
		t.Fatal("bursts must be capped and deterministic")
	}
	before := a.Particles()[0]
	a.Advance(0.05)
	b.Advance(0.05)
	after := a.Particles()[0]
	if after.X == before.X || after.Age != 0.05 || math.Hypot(after.VelocityX, after.VelocityY) >= math.Hypot(before.VelocityX, before.VelocityY) {
		t.Fatal("spark should move, age, and slow down")
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("motion must be deterministic")
	}
	snapshot := a.Particles()
	snapshot[0].X = -999
	if a.particles[0].X == -999 {
		t.Fatal("snapshot exposed mutable particle storage")
	}
	a.Advance(1)
	if len(a.Particles()) != 0 {
		t.Fatal("sparks did not expire")
	}
	a.Burst(ecs.ProjectileState{Entity: 101})
	if len(a.Particles()) != ImpactSparkCount {
		t.Fatal("expired particles did not free capacity")
	}
}

func TestPipelinePreservesNewSparkBurst(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	mapEntity := world.CreateEntity()
	must(t, world.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: 500, Height: 500}))
	hitFeedbackProjectile(t, world, 0, 6)
	sparks := NewSparkParticles()
	cfg := PipelineConfig{PlayerEntity: tank, MapEntity: mapEntity, Sparks: sparks}
	must(t, RunDefaultPipeline(world, 1.0/60, cfg))
	if len(sparks.Particles()) != ImpactSparkCount || sparks.Particles()[0].Age != 0 {
		t.Fatal("pipeline must advance old sparks before spawning new bursts")
	}
	must(t, RunDefaultPipeline(world, 1, cfg))
	if len(sparks.Particles()) != 0 {
		t.Fatal("pipeline did not expire sparks")
	}
}
