package sim

import (
	"math"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
)

func TestResolveTankObstacleCollisionsPushesTankOut(t *testing.T) {
	w := ecs.NewWorld()

	tank := w.CreateEntity()
	must(t, w.SetTransform(tank, ecs.Transform{X: 0, Y: 0}))
	must(t, w.SetTankDrive(tank, ecs.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(tank, ecs.CircleCollider{Radius: 10}))

	obstacle := w.CreateEntity()
	must(t, w.SetTransform(obstacle, ecs.Transform{X: 15, Y: 0}))
	must(t, w.SetCircleCollider(obstacle, ecs.CircleCollider{Radius: 10}))
	must(t, w.SetSolidObstacle(obstacle))

	if err := ResolveTankObstacleCollisions(w); err != nil {
		t.Fatalf("resolve collisions: %v", err)
	}

	tankTransform, ok := w.Transform(tank)
	if !ok {
		t.Fatalf("tank transform missing")
	}

	approxEqual(t, tankTransform.X, -5, 1e-9)
	approxEqual(t, tankTransform.Y, 0, 1e-9)
}

func TestResolveTankObstacleCollisionsNoOverlapDoesNotMoveTank(t *testing.T) {
	w := ecs.NewWorld()

	tank := w.CreateEntity()
	must(t, w.SetTransform(tank, ecs.Transform{X: 0, Y: 0}))
	must(t, w.SetTankDrive(tank, ecs.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(tank, ecs.CircleCollider{Radius: 10}))

	obstacle := w.CreateEntity()
	must(t, w.SetTransform(obstacle, ecs.Transform{X: 40, Y: 0}))
	must(t, w.SetCircleCollider(obstacle, ecs.CircleCollider{Radius: 10}))
	must(t, w.SetSolidObstacle(obstacle))

	if err := ResolveTankObstacleCollisions(w); err != nil {
		t.Fatalf("resolve collisions: %v", err)
	}

	tankTransform, ok := w.Transform(tank)
	if !ok {
		t.Fatalf("tank transform missing")
	}

	approxEqual(t, tankTransform.X, 0, 1e-9)
	approxEqual(t, tankTransform.Y, 0, 1e-9)
}

func TestResolveTankObstacleCollisionsHandlesCoincidentCenters(t *testing.T) {
	w := ecs.NewWorld()

	tank := w.CreateEntity()
	must(t, w.SetTransform(tank, ecs.Transform{X: 10, Y: 10}))
	must(t, w.SetTankDrive(tank, ecs.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(tank, ecs.CircleCollider{Radius: 10}))

	obstacle := w.CreateEntity()
	must(t, w.SetTransform(obstacle, ecs.Transform{X: 10, Y: 10}))
	must(t, w.SetCircleCollider(obstacle, ecs.CircleCollider{Radius: 10}))
	must(t, w.SetSolidObstacle(obstacle))

	if err := ResolveTankObstacleCollisions(w); err != nil {
		t.Fatalf("resolve collisions: %v", err)
	}

	tankTransform, ok := w.Transform(tank)
	if !ok {
		t.Fatalf("tank transform missing")
	}

	if math.IsNaN(tankTransform.X) || math.IsNaN(tankTransform.Y) {
		t.Fatalf("resolved tank transform should not be NaN: %+v", tankTransform)
	}
}
