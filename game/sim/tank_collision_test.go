package sim

import (
	"math"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestResolveTankTankCollisionsPushesBothTanksApart(t *testing.T) {
	w := ecs.NewWorld()

	a := w.CreateEntity()
	b := w.CreateEntity()

	must(t, w.SetTransform(a, ecscomponents.Transform{X: 0, Y: 0}))
	must(t, w.SetTankDrive(a, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(a, ecscomponents.CircleCollider{Radius: 10}))

	must(t, w.SetTransform(b, ecscomponents.Transform{X: 15, Y: 0}))
	must(t, w.SetTankDrive(b, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(b, ecscomponents.CircleCollider{Radius: 10}))

	if err := ResolveTankTankCollisions(w); err != nil {
		t.Fatalf("resolve tank collisions: %v", err)
	}

	aTransform, _ := w.Transform(a)
	bTransform, _ := w.Transform(b)

	approxEqual(t, aTransform.X, -2.5, 1e-9)
	approxEqual(t, bTransform.X, 17.5, 1e-9)
	approxEqual(t, aTransform.Y, 0, 1e-9)
	approxEqual(t, bTransform.Y, 0, 1e-9)
}

func TestResolveTankTankCollisionsNoOverlapNoMovement(t *testing.T) {
	w := ecs.NewWorld()

	a := w.CreateEntity()
	b := w.CreateEntity()

	must(t, w.SetTransform(a, ecscomponents.Transform{X: 0, Y: 0}))
	must(t, w.SetTankDrive(a, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(a, ecscomponents.CircleCollider{Radius: 10}))

	must(t, w.SetTransform(b, ecscomponents.Transform{X: 50, Y: 0}))
	must(t, w.SetTankDrive(b, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(b, ecscomponents.CircleCollider{Radius: 10}))

	if err := ResolveTankTankCollisions(w); err != nil {
		t.Fatalf("resolve tank collisions: %v", err)
	}

	aTransform, _ := w.Transform(a)
	bTransform, _ := w.Transform(b)

	approxEqual(t, aTransform.X, 0, 1e-9)
	approxEqual(t, bTransform.X, 50, 1e-9)
}

func TestResolveTankTankCollisionsHandlesCoincidentCenters(t *testing.T) {
	w := ecs.NewWorld()

	a := w.CreateEntity()
	b := w.CreateEntity()

	must(t, w.SetTransform(a, ecscomponents.Transform{X: 10, Y: 10}))
	must(t, w.SetTankDrive(a, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(a, ecscomponents.CircleCollider{Radius: 10}))

	must(t, w.SetTransform(b, ecscomponents.Transform{X: 10, Y: 10}))
	must(t, w.SetTankDrive(b, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetCircleCollider(b, ecscomponents.CircleCollider{Radius: 10}))

	if err := ResolveTankTankCollisions(w); err != nil {
		t.Fatalf("resolve tank collisions: %v", err)
	}

	aTransform, _ := w.Transform(a)
	bTransform, _ := w.Transform(b)
	if math.IsNaN(aTransform.X) || math.IsNaN(aTransform.Y) || math.IsNaN(bTransform.X) || math.IsNaN(bTransform.Y) {
		t.Fatalf("resolved tank transforms should not be NaN: a=%+v b=%+v", aTransform, bTransform)
	}
}
