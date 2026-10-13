package sim

import (
	"math"
	"reflect"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestTrackSpacingReverseAndStationaryTurn(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	tracks := NewTracks()
	tracks.Update(world, 0)
	move := func(y, rotation float64) {
		must(t, world.SetTransform(tank, ecscomponents.Transform{X: 100, Y: y, RotationRadians: rotation}))
		tracks.Update(world, 1.0/60)
	}
	move(98, 0)
	if len(tracks.Marks()) != 0 {
		t.Fatal("marks should be distance spaced")
	}
	move(90, 0)
	marks := tracks.Marks()
	if len(marks) != 2 || marks[0].Y != 96 || marks[1].Y != 92 {
		t.Fatalf("unexpected track spacing: %v", marks)
	}
	move(94, 0)
	if len(tracks.Marks()) != 3 {
		t.Fatal("reverse travel should stamp tracks")
	}
	move(94, math.Pi/2)
	if len(tracks.Marks()) != 3 {
		t.Fatal("stationary turn should not stamp")
	}
	must(t, world.SetTransform(tank, ecscomponents.Transform{X: 108, Y: 94, RotationRadians: math.Pi / 2}))
	tracks.Update(world, 1.0/60)
	marks = tracks.Marks()
	if marks[len(marks)-1].RotationRadians != math.Pi/2 {
		t.Fatal("tracks should follow hull turns")
	}
}

func TestTracksCleanupFadeAndCap(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	tracks := NewTracks()
	tracks.Update(world, 0)
	must(t, world.SetTransform(tank, ecscomponents.Transform{X: 100, Y: 80}))
	tracks.Update(world, 0)
	world.DestroyEntity(tank)
	tracks.Update(world, 0)
	if len(tracks.emitters) != 0 || len(tracks.Marks()) != 5 {
		t.Fatal("destroyed tank should leave its marks but no emitter")
	}
	tracks.Advance(TrackLifetime)
	if len(tracks.Marks()) != 0 {
		t.Fatal("marks did not fade out")
	}
	for i := 0; i < MaxTrackMarks+10; i++ {
		tracks.add(TrackMark{X: float64(i)})
	}
	marks := tracks.Marks()
	if len(marks) != MaxTrackMarks || marks[0].X != 10 {
		t.Fatal("cap should replace oldest marks")
	}
	marks[0].X = -99
	if tracks.marks[0].X == -99 {
		t.Fatal("query exposed mutable storage")
	}
}

func TestTracksIgnoreTeleportAndDeadTank(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	tracks := NewTracks()
	tracks.Update(world, 0)
	must(t, world.SetTransform(tank, ecscomponents.Transform{X: 100, Y: 400}))
	tracks.Update(world, 0)
	if len(tracks.Marks()) != 0 {
		t.Fatal("teleport should not leave a long trail")
	}
	must(t, world.SetHealth(tank, ecscomponents.Health{Max: 100}))
	must(t, world.SetTransform(tank, ecscomponents.Transform{X: 100, Y: 380}))
	tracks.Update(world, 0)
	if len(tracks.Marks()) != 0 || len(tracks.emitters) != 0 {
		t.Fatal("dead tank should not stamp")
	}
}

func TestTracksDeterministic(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	a, b := NewTracks(), NewTracks()
	for tick := 0; tick < 120; tick++ {
		must(t, world.SetTransform(tank, ecscomponents.Transform{X: 100, Y: 400 - float64(tick)*2}))
		a.Update(world, 1.0/60)
		b.Update(world, 1.0/60)
		if !reflect.DeepEqual(a, b) {
			t.Fatal("identical travel produced different marks")
		}
	}
}

func TestPipelineTracksUseResolvedPosition(t *testing.T) {
	world := ecs.NewWorld()
	tank := hitFeedbackTank(t, world)
	mapEntity := world.CreateEntity()
	must(t, world.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: 500, Height: 500}))
	cover := world.CreateEntity()
	must(t, world.SetTransform(cover, ecscomponents.Transform{X: 100, Y: 68}))
	must(t, world.SetCircleCollider(cover, ecscomponents.CircleCollider{Radius: 16}))
	must(t, world.SetSolidObstacle(cover))
	must(t, world.SetTankDrive(tank, ecscomponents.TankDrive{MoveSpeed: 60}))
	must(t, world.SetTankControlIntent(tank, ecscomponents.TankControlIntent{Throttle: 1}))
	tracks := NewTracks()
	tracks.Update(world, 0)
	for tick := 0; tick < 10; tick++ {
		must(t, RunDefaultPipeline(world, 0.1, PipelineConfig{PlayerEntity: tank, MapEntity: mapEntity, Tracks: tracks}))
	}
	if len(tracks.Marks()) != 0 {
		t.Fatal("tank pushing against cover should not leave marks from unresolved movement")
	}
}
