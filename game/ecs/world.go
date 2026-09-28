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
	SpriteID string
	Layer    int
}

type TankControlIntent struct {
	Throttle float64
	Turn     float64
}

type TankDrive struct {
	MoveSpeed          float64
	TurnSpeedRadiansPS float64
}

type Renderable struct {
	Entity    Entity
	Transform Transform
	Sprite    SpriteRenderer
}

type TankMover struct {
	Entity    Entity
	Transform Transform
	Drive     TankDrive
	Intent    TankControlIntent
}

type World struct {
	nextEntityID Entity
	alive        map[Entity]struct{}
	transforms   map[Entity]Transform
	sprites      map[Entity]SpriteRenderer
	intents      map[Entity]TankControlIntent
	drives       map[Entity]TankDrive
}

func NewWorld() *World {
	return &World{
		nextEntityID: 1,
		alive:        make(map[Entity]struct{}),
		transforms:   make(map[Entity]Transform),
		sprites:      make(map[Entity]SpriteRenderer),
		intents:      make(map[Entity]TankControlIntent),
		drives:       make(map[Entity]TankDrive),
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
		if a.Sprite.Layer < b.Sprite.Layer {
			return -1
		}
		if a.Sprite.Layer > b.Sprite.Layer {
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

		intent := w.intents[entity]

		movers = append(movers, TankMover{
			Entity:    entity,
			Transform: transform,
			Drive:     drive,
			Intent:    intent,
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
