package sim

import (
	"math"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
)

func TestIntegrateTankMovementMovesForwardFromZeroRotation(t *testing.T) {
	w := ecs.NewWorld()
	e := w.CreateEntity()

	must(t, w.SetTransform(e, ecs.Transform{X: 480, Y: 270, RotationRadians: 0}))
	must(t, w.SetTankDrive(e, ecs.TankDrive{MoveSpeed: 120, TurnSpeedRadiansPS: math.Pi}))
	must(t, w.SetTankControlIntent(e, ecs.TankControlIntent{Throttle: 1, Turn: 0}))

	IntegrateTankMovement(w, 1)

	got, ok := w.Transform(e)
	if !ok {
		t.Fatalf("expected transform")
	}

	if got.X != 480 {
		t.Fatalf("expected X unchanged, got %f", got.X)
	}
	if got.Y != 150 {
		t.Fatalf("expected Y moved up by 120, got %f", got.Y)
	}
}

func TestIntegrateTankMovementRotatesClockwiseAndWraps(t *testing.T) {
	w := ecs.NewWorld()
	e := w.CreateEntity()

	must(t, w.SetTransform(e, ecs.Transform{X: 0, Y: 0, RotationRadians: 2*math.Pi - 0.1}))
	must(t, w.SetTankDrive(e, ecs.TankDrive{MoveSpeed: 0, TurnSpeedRadiansPS: 1}))
	must(t, w.SetTankControlIntent(e, ecs.TankControlIntent{Throttle: 0, Turn: 1}))

	IntegrateTankMovement(w, 1)

	got, ok := w.Transform(e)
	if !ok {
		t.Fatalf("expected transform")
	}

	if got.RotationRadians < 0 || got.RotationRadians >= 2*math.Pi {
		t.Fatalf("rotation should be wrapped to [0, 2pi), got %f", got.RotationRadians)
	}

	approxEqual(t, got.RotationRadians, 0.9, 1e-9)
}

func TestIntegrateTankMovementClampsIntentToUnitRange(t *testing.T) {
	w := ecs.NewWorld()
	e := w.CreateEntity()

	must(t, w.SetTransform(e, ecs.Transform{X: 10, Y: 10, RotationRadians: 0}))
	must(t, w.SetTankDrive(e, ecs.TankDrive{MoveSpeed: 100, TurnSpeedRadiansPS: math.Pi}))
	must(t, w.SetTankControlIntent(e, ecs.TankControlIntent{Throttle: 3, Turn: -2}))

	IntegrateTankMovement(w, 0.5)

	got, ok := w.Transform(e)
	if !ok {
		t.Fatalf("expected transform")
	}

	// Throttle 3 is clamped to 1 -> 50 units forward in the post-turn heading.
	approxEqual(t, got.X, -40, 1e-9)
	approxEqual(t, got.Y, 10, 1e-9)
	// Turn -2 is clamped to -1 -> rotate by -pi/2, then wrapped to 3pi/2.
	approxEqual(t, got.RotationRadians, 1.5*math.Pi, 1e-9)
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func approxEqual(t *testing.T, got, want, tolerance float64) {
	t.Helper()
	if math.Abs(got-want) > tolerance {
		t.Fatalf("got %f want %f (tol=%f)", got, want, tolerance)
	}
}
