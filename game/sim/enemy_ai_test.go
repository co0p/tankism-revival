package sim

import (
	"math"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestUpdateChaserIntentTargetAhead(t *testing.T) {
	w := ecs.NewWorld()
	chaser := w.CreateEntity()
	target := w.CreateEntity()

	must(t, w.SetTransform(chaser, ecscomponents.Transform{X: 0, Y: 0, RotationRadians: 0}))
	must(t, w.SetTransform(target, ecscomponents.Transform{X: 0, Y: -100, RotationRadians: 0}))
	must(t, w.SetTankControlIntent(chaser, ecscomponents.TankControlIntent{}))

	if err := UpdateChaserIntent(w, chaser, target, 10); err != nil {
		t.Fatalf("update chaser intent: %v", err)
	}

	intent, ok := w.TankControlIntent(chaser)
	if !ok {
		t.Fatalf("expected chaser intent")
	}

	approxEqual(t, intent.Turn, 0, 1e-9)
	approxEqual(t, intent.Throttle, 1, 1e-9)
	if !intent.Fire {
		t.Fatalf("expected chaser to fire when target is ahead and in range")
	}
}

func TestUpdateChaserIntentTurnsTowardRightTarget(t *testing.T) {
	w := ecs.NewWorld()
	chaser := w.CreateEntity()
	target := w.CreateEntity()

	must(t, w.SetTransform(chaser, ecscomponents.Transform{X: 0, Y: 0, RotationRadians: 0}))
	must(t, w.SetTransform(target, ecscomponents.Transform{X: 100, Y: 0, RotationRadians: 0}))

	if err := UpdateChaserIntent(w, chaser, target, 10); err != nil {
		t.Fatalf("update chaser intent: %v", err)
	}

	intent, ok := w.TankControlIntent(chaser)
	if !ok {
		t.Fatalf("expected chaser intent")
	}

	if intent.Turn <= 0 {
		t.Fatalf("expected positive turn toward right-side target, got %f", intent.Turn)
	}

	if intent.Fire {
		t.Fatalf("chaser should not fire while not aligned")
	}
}

func TestUpdateChaserIntentStopsWithinStopDistance(t *testing.T) {
	w := ecs.NewWorld()
	chaser := w.CreateEntity()
	target := w.CreateEntity()

	must(t, w.SetTransform(chaser, ecscomponents.Transform{X: 0, Y: 0, RotationRadians: math.Pi / 3}))
	must(t, w.SetTransform(target, ecscomponents.Transform{X: 20, Y: 0, RotationRadians: 0}))

	if err := UpdateChaserIntent(w, chaser, target, 40); err != nil {
		t.Fatalf("update chaser intent: %v", err)
	}

	intent, ok := w.TankControlIntent(chaser)
	if !ok {
		t.Fatalf("expected chaser intent")
	}

	approxEqual(t, intent.Throttle, 0, 1e-9)
	if intent.Fire {
		t.Fatalf("chaser should not fire while heavily misaligned")
	}
}
