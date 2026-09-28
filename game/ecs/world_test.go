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
	if err := w.SetSpriteRenderer(e, SpriteRenderer{SpriteID: "tank_red_small", Layer: 1}); err != nil {
		t.Fatalf("set sprite renderer: %v", err)
	}
	if err := w.SetTankControlIntent(e, TankControlIntent{Throttle: 1, Turn: -0.5}); err != nil {
		t.Fatalf("set tank control intent: %v", err)
	}
	if err := w.SetTankDrive(e, TankDrive{MoveSpeed: 100, TurnSpeedRadiansPS: math.Pi}); err != nil {
		t.Fatalf("set tank drive: %v", err)
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
}

func TestRenderablesReturnsEntitiesWithTransformAndSpriteSortedByLayerThenEntity(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()
	e3 := w.CreateEntity()

	if err := w.SetTransform(e1, Transform{X: 100, Y: 200, RotationRadians: math.Pi}); err != nil {
		t.Fatalf("set transform e1: %v", err)
	}
	if err := w.SetSpriteRenderer(e1, SpriteRenderer{SpriteID: "red", Layer: 5}); err != nil {
		t.Fatalf("set sprite e1: %v", err)
	}

	if err := w.SetTransform(e2, Transform{X: 1, Y: 2, RotationRadians: 0}); err != nil {
		t.Fatalf("set transform e2: %v", err)
	}
	if err := w.SetSpriteRenderer(e2, SpriteRenderer{SpriteID: "blue", Layer: 1}); err != nil {
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

	if movers[1].Intent != (TankControlIntent{Throttle: 0.5, Turn: -1}) {
		t.Fatalf("unexpected e2 intent: %+v", movers[1].Intent)
	}
}
