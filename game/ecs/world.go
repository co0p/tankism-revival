package ecs

import (
	"errors"
	"slices"
)

var ErrEntityNotAlive = errors.New("entity is not alive")

type Entity uint64

type Transform struct {
	X               float64
	Y               float64
	RotationRadians float64
}

type SpriteRenderer struct {
	SpriteID              string
	ZIndex                int
	RotationOffsetRadians float64
	FlipX                 bool
	FlipY                 bool
}

type TankControlIntent struct {
	Throttle float64
	Turn     float64
}

type TankDrive struct {
	MoveSpeed          float64
	TurnSpeedRadiansPS float64
}

type TankAcceleration struct {
	ForwardUnitsPerSecondSquared float64
}

type TankDeceleration struct {
	CoastUnitsPerSecondSquared float64
	BrakeUnitsPerSecondSquared float64
}

type TankVelocity struct {
	ForwardUnitsPerSecond float64
}

type CircleCollider struct {
	Radius float64
}

type ColliderBody struct {
	Entity    Entity
	Transform Transform
	Collider  CircleCollider
}

type Renderable struct {
	Entity    Entity
	Transform Transform
	Sprite    SpriteRenderer
}

type TankMover struct {
	Entity       Entity
	Transform    Transform
	Drive        TankDrive
	Acceleration TankAcceleration
	Deceleration TankDeceleration
	Velocity     TankVelocity
	Intent       TankControlIntent
}

type World struct {
	nextEntityID Entity
	alive        map[Entity]struct{}
	transforms   map[Entity]Transform
	sprites      map[Entity]SpriteRenderer
	intents      map[Entity]TankControlIntent
	drives       map[Entity]TankDrive
	acceleration map[Entity]TankAcceleration
	deceleration map[Entity]TankDeceleration
	velocity     map[Entity]TankVelocity
	colliders    map[Entity]CircleCollider
	solid        map[Entity]struct{}
}

func NewWorld() *World {
	return &World{
		nextEntityID: 1,
		alive:        make(map[Entity]struct{}),
		transforms:   make(map[Entity]Transform),
		sprites:      make(map[Entity]SpriteRenderer),
		intents:      make(map[Entity]TankControlIntent),
		drives:       make(map[Entity]TankDrive),
		acceleration: make(map[Entity]TankAcceleration),
		deceleration: make(map[Entity]TankDeceleration),
		velocity:     make(map[Entity]TankVelocity),
		colliders:    make(map[Entity]CircleCollider),
		solid:        make(map[Entity]struct{}),
	}
}

func (w *World) CreateEntity() Entity {
	entity := w.nextEntityID
	w.nextEntityID++
	w.alive[entity] = struct{}{}

	return entity
}

func (w *World) IsAlive(entity Entity) bool {
	_, ok := w.alive[entity]
	return ok
}

func (w *World) DestroyEntity(entity Entity) bool {
	if !w.IsAlive(entity) {
		return false
	}

	delete(w.alive, entity)
	delete(w.transforms, entity)
	delete(w.sprites, entity)
	delete(w.intents, entity)
	delete(w.drives, entity)
	delete(w.acceleration, entity)
	delete(w.deceleration, entity)
	delete(w.velocity, entity)
	delete(w.colliders, entity)
	delete(w.solid, entity)

	return true
}

func (w *World) SetTransform(entity Entity, transform Transform) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.transforms[entity] = transform
	return nil
}

func (w *World) Transform(entity Entity) (Transform, bool) {
	transform, ok := w.transforms[entity]
	return transform, ok
}

func (w *World) RemoveTransform(entity Entity) {
	delete(w.transforms, entity)
}

func (w *World) SetSpriteRenderer(entity Entity, sprite SpriteRenderer) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.sprites[entity] = sprite
	return nil
}

func (w *World) SpriteRenderer(entity Entity) (SpriteRenderer, bool) {
	sprite, ok := w.sprites[entity]
	return sprite, ok
}

func (w *World) RemoveSpriteRenderer(entity Entity) {
	delete(w.sprites, entity)
}

func (w *World) SetTankControlIntent(entity Entity, intent TankControlIntent) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.intents[entity] = intent
	return nil
}

func (w *World) TankControlIntent(entity Entity) (TankControlIntent, bool) {
	intent, ok := w.intents[entity]
	return intent, ok
}

func (w *World) RemoveTankControlIntent(entity Entity) {
	delete(w.intents, entity)
}

func (w *World) SetTankDrive(entity Entity, drive TankDrive) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.drives[entity] = drive
	return nil
}

func (w *World) TankDrive(entity Entity) (TankDrive, bool) {
	drive, ok := w.drives[entity]
	return drive, ok
}

func (w *World) RemoveTankDrive(entity Entity) {
	delete(w.drives, entity)
}

func (w *World) SetTankAcceleration(entity Entity, acceleration TankAcceleration) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.acceleration[entity] = acceleration
	return nil
}

func (w *World) TankAcceleration(entity Entity) (TankAcceleration, bool) {
	acceleration, ok := w.acceleration[entity]
	return acceleration, ok
}

func (w *World) RemoveTankAcceleration(entity Entity) {
	delete(w.acceleration, entity)
}

func (w *World) SetTankDeceleration(entity Entity, deceleration TankDeceleration) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.deceleration[entity] = deceleration
	return nil
}

func (w *World) TankDeceleration(entity Entity) (TankDeceleration, bool) {
	deceleration, ok := w.deceleration[entity]
	return deceleration, ok
}

func (w *World) RemoveTankDeceleration(entity Entity) {
	delete(w.deceleration, entity)
}

func (w *World) SetTankVelocity(entity Entity, velocity TankVelocity) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.velocity[entity] = velocity
	return nil
}

func (w *World) TankVelocity(entity Entity) (TankVelocity, bool) {
	velocity, ok := w.velocity[entity]
	return velocity, ok
}

func (w *World) RemoveTankVelocity(entity Entity) {
	delete(w.velocity, entity)
}

func (w *World) SetCircleCollider(entity Entity, collider CircleCollider) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.colliders[entity] = collider
	return nil
}

func (w *World) CircleCollider(entity Entity) (CircleCollider, bool) {
	collider, ok := w.colliders[entity]
	return collider, ok
}

func (w *World) RemoveCircleCollider(entity Entity) {
	delete(w.colliders, entity)
}

func (w *World) SetSolidObstacle(entity Entity) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.solid[entity] = struct{}{}
	return nil
}

func (w *World) IsSolidObstacle(entity Entity) bool {
	_, ok := w.solid[entity]
	return ok
}

func (w *World) RemoveSolidObstacle(entity Entity) {
	delete(w.solid, entity)
}

func (w *World) Renderables() []Renderable {
	renderables := make([]Renderable, 0, len(w.sprites))

	for entity, sprite := range w.sprites {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		renderables = append(renderables, Renderable{
			Entity:    entity,
			Transform: transform,
			Sprite:    sprite,
		})
	}

	slices.SortFunc(renderables, func(a, b Renderable) int {
		if a.Sprite.ZIndex < b.Sprite.ZIndex {
			return -1
		}
		if a.Sprite.ZIndex > b.Sprite.ZIndex {
			return 1
		}

		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return renderables
}

func (w *World) TankMovers() []TankMover {
	movers := make([]TankMover, 0, len(w.drives))

	for entity, drive := range w.drives {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		acceleration := w.acceleration[entity]
		deceleration := w.deceleration[entity]
		velocity := w.velocity[entity]
		intent := w.intents[entity]

		movers = append(movers, TankMover{
			Entity:       entity,
			Transform:    transform,
			Drive:        drive,
			Acceleration: acceleration,
			Deceleration: deceleration,
			Velocity:     velocity,
			Intent:       intent,
		})
	}

	slices.SortFunc(movers, func(a, b TankMover) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return movers
}

func (w *World) TankColliderBodies() []ColliderBody {
	bodies := make([]ColliderBody, 0, len(w.drives))

	for entity := range w.drives {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		collider, hasCollider := w.colliders[entity]
		if !hasCollider {
			continue
		}

		bodies = append(bodies, ColliderBody{
			Entity:    entity,
			Transform: transform,
			Collider:  collider,
		})
	}

	slices.SortFunc(bodies, func(a, b ColliderBody) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return bodies
}

func (w *World) SolidObstacleColliderBodies() []ColliderBody {
	bodies := make([]ColliderBody, 0, len(w.solid))

	for entity := range w.solid {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		collider, hasCollider := w.colliders[entity]
		if !hasCollider {
			continue
		}

		bodies = append(bodies, ColliderBody{
			Entity:    entity,
			Transform: transform,
			Collider:  collider,
		})
	}

	slices.SortFunc(bodies, func(a, b ColliderBody) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return bodies
}
