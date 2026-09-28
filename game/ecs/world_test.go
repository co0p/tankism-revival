package ecs

import (
	"errors"
	"math"
	"testing"
)

func TestCreateEntityAllocatesUniqueAliveIDs(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()

	if e1 == 0 {
		t.Fatalf("first entity id should be non-zero")
	}

	if e2 <= e1 {
		t.Fatalf("entity ids should be monotonic: e1=%d e2=%d", e1, e2)
	}

	if !w.IsAlive(e1) || !w.IsAlive(e2) {
		t.Fatalf("new entities must be alive")
	}
}

func TestDestroyEntityRemovesAllComponents(t *testing.T) {
	w := NewWorld()
	e := w.CreateEntity()

	if err := w.SetTransform(e, Transform{X: 10, Y: 20, RotationRadians: math.Pi / 2}); err != nil {
		t.Fatalf("set transform: %v", err)
	}
	if err := w.SetSpriteRenderer(e, SpriteRenderer{SpriteID: "tank_red_small", ZIndex: 1}); err != nil {
		t.Fatalf("set sprite renderer: %v", err)
	}
	if err := w.SetTankControlIntent(e, TankControlIntent{Throttle: 1, Turn: -0.5}); err != nil {
		t.Fatalf("set tank control intent: %v", err)
	}
	if err := w.SetTankDrive(e, TankDrive{MoveSpeed: 100, TurnSpeedRadiansPS: math.Pi}); err != nil {
		t.Fatalf("set tank drive: %v", err)
	}
	if err := w.SetTankAcceleration(e, TankAcceleration{ForwardUnitsPerSecondSquared: 200}); err != nil {
		t.Fatalf("set tank acceleration: %v", err)
	}
	if err := w.SetTankDeceleration(e, TankDeceleration{CoastUnitsPerSecondSquared: 120, BrakeUnitsPerSecondSquared: 320}); err != nil {
		t.Fatalf("set tank deceleration: %v", err)
	}
	if err := w.SetTankVelocity(e, TankVelocity{ForwardUnitsPerSecond: 45}); err != nil {
		t.Fatalf("set tank velocity: %v", err)
	}
	if err := w.SetCircleCollider(e, CircleCollider{Radius: 10}); err != nil {
		t.Fatalf("set circle collider: %v", err)
	}
	if err := w.SetSolidObstacle(e); err != nil {
		t.Fatalf("set solid obstacle: %v", err)
	}

	if ok := w.DestroyEntity(e); !ok {
		t.Fatalf("destroy entity should succeed")
	}

	if w.IsAlive(e) {
		t.Fatalf("entity should not be alive after destroy")
	}

	if _, ok := w.Transform(e); ok {
		t.Fatalf("transform should be removed on destroy")
	}

	if _, ok := w.SpriteRenderer(e); ok {
		t.Fatalf("sprite renderer should be removed on destroy")
	}

	if _, ok := w.TankControlIntent(e); ok {
		t.Fatalf("tank control intent should be removed on destroy")
	}

	if _, ok := w.TankDrive(e); ok {
		t.Fatalf("tank drive should be removed on destroy")
	}

	if _, ok := w.TankAcceleration(e); ok {
		t.Fatalf("tank acceleration should be removed on destroy")
	}

	if _, ok := w.TankDeceleration(e); ok {
		t.Fatalf("tank deceleration should be removed on destroy")
	}

	if _, ok := w.TankVelocity(e); ok {
		t.Fatalf("tank velocity should be removed on destroy")
	}

	if _, ok := w.CircleCollider(e); ok {
		t.Fatalf("circle collider should be removed on destroy")
	}

	if w.IsSolidObstacle(e) {
		t.Fatalf("solid obstacle marker should be removed on destroy")
	}

	if ok := w.DestroyEntity(e); ok {
		t.Fatalf("destroying an already-destroyed entity should return false")
	}
}

func TestSetComponentOnDeadEntityFails(t *testing.T) {
	w := NewWorld()

	err := w.SetTransform(999, Transform{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set transform on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetSpriteRenderer(999, SpriteRenderer{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set sprite renderer on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankControlIntent(999, TankControlIntent{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank control intent on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankDrive(999, TankDrive{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank drive on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankAcceleration(999, TankAcceleration{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank acceleration on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankDeceleration(999, TankDeceleration{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank deceleration on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankVelocity(999, TankVelocity{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank velocity on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetCircleCollider(999, CircleCollider{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set circle collider on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetSolidObstacle(999)
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set solid obstacle on dead entity should return ErrEntityNotAlive, got: %v", err)
	}
}

func TestRenderablesReturnsEntitiesWithTransformAndSpriteSortedByZIndexThenEntity(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()
	e3 := w.CreateEntity()

	if err := w.SetTransform(e1, Transform{X: 100, Y: 200, RotationRadians: math.Pi}); err != nil {
		t.Fatalf("set transform e1: %v", err)
	}
	if err := w.SetSpriteRenderer(e1, SpriteRenderer{SpriteID: "red", ZIndex: 5}); err != nil {
		t.Fatalf("set sprite e1: %v", err)
	}

	if err := w.SetTransform(e2, Transform{X: 1, Y: 2, RotationRadians: 0}); err != nil {
		t.Fatalf("set transform e2: %v", err)
	}
	if err := w.SetSpriteRenderer(e2, SpriteRenderer{SpriteID: "blue", ZIndex: 1}); err != nil {
		t.Fatalf("set sprite e2: %v", err)
	}

	if err := w.SetTransform(e3, Transform{X: 5, Y: 6, RotationRadians: math.Pi / 4}); err != nil {
		t.Fatalf("set transform e3: %v", err)
	}

	renderables := w.Renderables()

	if len(renderables) != 2 {
		t.Fatalf("expected 2 renderables, got %d", len(renderables))
	}

	if renderables[0].Entity != e2 || renderables[1].Entity != e1 {
		t.Fatalf("unexpected render order: got [%d, %d], want [%d, %d]", renderables[0].Entity, renderables[1].Entity, e2, e1)
	}
}

func TestTransformPreservesCenterCoordinatesAndRadians(t *testing.T) {
	w := NewWorld()
	e := w.CreateEntity()

	want := Transform{
		X:               480.5,
		Y:               270.25,
		RotationRadians: math.Pi / 2,
	}

	if err := w.SetTransform(e, want); err != nil {
		t.Fatalf("set transform: %v", err)
	}

	got, ok := w.Transform(e)
	if !ok {
		t.Fatalf("expected transform to exist")
	}

	if got != want {
		t.Fatalf("transform mismatch: got %+v want %+v", got, want)
	}
}

func TestTankMoversReturnsDriveAndDefaultIntent(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()

	if err := w.SetTransform(e1, Transform{X: 100, Y: 200, RotationRadians: 0}); err != nil {
		t.Fatalf("set transform e1: %v", err)
	}
	if err := w.SetTankDrive(e1, TankDrive{MoveSpeed: 180, TurnSpeedRadiansPS: math.Pi}); err != nil {
		t.Fatalf("set drive e1: %v", err)
	}

	if err := w.SetTransform(e2, Transform{X: 1, Y: 2, RotationRadians: math.Pi / 2}); err != nil {
		t.Fatalf("set transform e2: %v", err)
	}
	if err := w.SetTankDrive(e2, TankDrive{MoveSpeed: 90, TurnSpeedRadiansPS: math.Pi / 2}); err != nil {
		t.Fatalf("set drive e2: %v", err)
	}
	if err := w.SetTankAcceleration(e2, TankAcceleration{ForwardUnitsPerSecondSquared: 300}); err != nil {
		t.Fatalf("set acceleration e2: %v", err)
	}
	if err := w.SetTankDeceleration(e2, TankDeceleration{CoastUnitsPerSecondSquared: 140, BrakeUnitsPerSecondSquared: 360}); err != nil {
		t.Fatalf("set deceleration e2: %v", err)
	}
	if err := w.SetTankVelocity(e2, TankVelocity{ForwardUnitsPerSecond: 12}); err != nil {
		t.Fatalf("set velocity e2: %v", err)
	}
	if err := w.SetTankControlIntent(e2, TankControlIntent{Throttle: 0.5, Turn: -1}); err != nil {
		t.Fatalf("set intent e2: %v", err)
	}

	movers := w.TankMovers()
	if len(movers) != 2 {
		t.Fatalf("expected 2 tank movers, got %d", len(movers))
	}

	if movers[0].Entity != e1 || movers[1].Entity != e2 {
		t.Fatalf("unexpected mover order: got [%d, %d], want [%d, %d]", movers[0].Entity, movers[1].Entity, e1, e2)
	}

	if movers[0].Intent != (TankControlIntent{}) {
		t.Fatalf("missing intent should default to zero intent, got %+v", movers[0].Intent)
	}
	if movers[0].Acceleration != (TankAcceleration{}) {
		t.Fatalf("missing acceleration should default to zero acceleration, got %+v", movers[0].Acceleration)
	}
	if movers[0].Deceleration != (TankDeceleration{}) {
		t.Fatalf("missing deceleration should default to zero deceleration, got %+v", movers[0].Deceleration)
	}
	if movers[0].Velocity != (TankVelocity{}) {
		t.Fatalf("missing velocity should default to zero velocity, got %+v", movers[0].Velocity)
	}

	if movers[1].Intent != (TankControlIntent{Throttle: 0.5, Turn: -1}) {
		t.Fatalf("unexpected e2 intent: %+v", movers[1].Intent)
	}
	if movers[1].Acceleration != (TankAcceleration{ForwardUnitsPerSecondSquared: 300}) {
		t.Fatalf("unexpected e2 acceleration: %+v", movers[1].Acceleration)
	}
	if movers[1].Deceleration != (TankDeceleration{CoastUnitsPerSecondSquared: 140, BrakeUnitsPerSecondSquared: 360}) {
		t.Fatalf("unexpected e2 deceleration: %+v", movers[1].Deceleration)
	}
	if movers[1].Velocity != (TankVelocity{ForwardUnitsPerSecond: 12}) {
		t.Fatalf("unexpected e2 velocity: %+v", movers[1].Velocity)
	}
}

func TestColliderQueriesReturnExpectedBodies(t *testing.T) {
	w := NewWorld()

	tank := w.CreateEntity()
	if err := w.SetTransform(tank, Transform{X: 10, Y: 20}); err != nil {
		t.Fatalf("set tank transform: %v", err)
	}
	if err := w.SetTankDrive(tank, TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}); err != nil {
		t.Fatalf("set tank drive: %v", err)
	}
	if err := w.SetCircleCollider(tank, CircleCollider{Radius: 9}); err != nil {
		t.Fatalf("set tank collider: %v", err)
	}

	obstacle := w.CreateEntity()
	if err := w.SetTransform(obstacle, Transform{X: 100, Y: 120}); err != nil {
		t.Fatalf("set obstacle transform: %v", err)
	}
	if err := w.SetCircleCollider(obstacle, CircleCollider{Radius: 14}); err != nil {
		t.Fatalf("set obstacle collider: %v", err)
	}
	if err := w.SetSolidObstacle(obstacle); err != nil {
		t.Fatalf("set obstacle solid: %v", err)
	}

	tankBodies := w.TankColliderBodies()
	if len(tankBodies) != 1 {
		t.Fatalf("expected 1 tank collider body, got %d", len(tankBodies))
	}
	if tankBodies[0].Entity != tank {
		t.Fatalf("unexpected tank collider entity: got %d want %d", tankBodies[0].Entity, tank)
	}

	obstacleBodies := w.SolidObstacleColliderBodies()
	if len(obstacleBodies) != 1 {
		t.Fatalf("expected 1 obstacle collider body, got %d", len(obstacleBodies))
	}
	if obstacleBodies[0].Entity != obstacle {
		t.Fatalf("unexpected obstacle collider entity: got %d want %d", obstacleBodies[0].Entity, obstacle)
	}
}
