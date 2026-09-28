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
	Fire     bool
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

type Weapon struct {
	CooldownSeconds           float64
	CooldownRemaining         float64
	ProjectileSpeed           float64
	ProjectileLifetime        float64
	ProjectileDamage          float64
	MuzzleOffset              float64
	ProjectileSpriteID        string
	ProjectileRadius          float64
	ProjectileLightRadius     float64
	ProjectileLightIntensity  float64
	ProjectileLightColor      LightColor
	MuzzleFlashSpriteID       string
	MuzzleFlashLifetime       float64
	MuzzleFlashLightRadius    float64
	MuzzleFlashLightIntensity float64
	MuzzleFlashLightColor     LightColor
}

type Projectile struct {
	Owner             Entity
	VelocityX         float64
	VelocityY         float64
	RemainingLifetime float64
	Damage            float64
}

type Lifetime struct {
	RemainingSeconds float64
}

type Health struct {
	Current float64
	Max     float64
}

type LightColor struct {
	R float64
	G float64
	B float64
}

var LightColorWhite = LightColor{R: 1, G: 1, B: 1}

type AmbientLight struct {
	Darkness float64
}

type PointLight struct {
	Radius    float64
	Intensity float64
	Color     LightColor
}

type CircleCollider struct {
	Radius float64
}

type ColliderBody struct {
	Entity    Entity
	Transform Transform
	Collider  CircleCollider
}

type PointLightSource struct {
	Entity    Entity
	Transform Transform
	Light     PointLight
}

type ArmedTank struct {
	Entity    Entity
	Transform Transform
	Intent    TankControlIntent
	Weapon    Weapon
}

type ProjectileState struct {
	Entity     Entity
	Transform  Transform
	Projectile Projectile
	Collider   CircleCollider
}

type TankHealthState struct {
	Entity    Entity
	Transform Transform
	Health    Health
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
	health       map[Entity]Health
	weapons      map[Entity]Weapon
	projectiles  map[Entity]Projectile
	lifetimes    map[Entity]Lifetime
	colliders    map[Entity]CircleCollider
	solid        map[Entity]struct{}
	ambient      map[Entity]AmbientLight
	pointLights  map[Entity]PointLight
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
		health:       make(map[Entity]Health),
		weapons:      make(map[Entity]Weapon),
		projectiles:  make(map[Entity]Projectile),
		lifetimes:    make(map[Entity]Lifetime),
		colliders:    make(map[Entity]CircleCollider),
		solid:        make(map[Entity]struct{}),
		ambient:      make(map[Entity]AmbientLight),
		pointLights:  make(map[Entity]PointLight),
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
	delete(w.health, entity)
	delete(w.weapons, entity)
	delete(w.projectiles, entity)
	delete(w.lifetimes, entity)
	delete(w.colliders, entity)
	delete(w.solid, entity)
	delete(w.ambient, entity)
	delete(w.pointLights, entity)

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

func (w *World) SetHealth(entity Entity, health Health) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.health[entity] = health
	return nil
}

func (w *World) Health(entity Entity) (Health, bool) {
	health, ok := w.health[entity]
	return health, ok
}

func (w *World) RemoveHealth(entity Entity) {
	delete(w.health, entity)
}

func (w *World) SetWeapon(entity Entity, weapon Weapon) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.weapons[entity] = weapon
	return nil
}

func (w *World) Weapon(entity Entity) (Weapon, bool) {
	weapon, ok := w.weapons[entity]
	return weapon, ok
}

func (w *World) RemoveWeapon(entity Entity) {
	delete(w.weapons, entity)
}

func (w *World) SetProjectile(entity Entity, projectile Projectile) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.projectiles[entity] = projectile
	return nil
}

func (w *World) Projectile(entity Entity) (Projectile, bool) {
	projectile, ok := w.projectiles[entity]
	return projectile, ok
}

func (w *World) RemoveProjectile(entity Entity) {
	delete(w.projectiles, entity)
}

func (w *World) SetLifetime(entity Entity, lifetime Lifetime) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.lifetimes[entity] = lifetime
	return nil
}

func (w *World) Lifetime(entity Entity) (Lifetime, bool) {
	lifetime, ok := w.lifetimes[entity]
	return lifetime, ok
}

func (w *World) RemoveLifetime(entity Entity) {
	delete(w.lifetimes, entity)
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

func (w *World) SetAmbientLight(entity Entity, light AmbientLight) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.ambient[entity] = light
	return nil
}

func (w *World) AmbientLight(entity Entity) (AmbientLight, bool) {
	light, ok := w.ambient[entity]
	return light, ok
}

func (w *World) RemoveAmbientLight(entity Entity) {
	delete(w.ambient, entity)
}

func (w *World) SetPointLight(entity Entity, light PointLight) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	w.pointLights[entity] = light
	return nil
}

func (w *World) PointLight(entity Entity) (PointLight, bool) {
	light, ok := w.pointLights[entity]
	return light, ok
}

func (w *World) RemovePointLight(entity Entity) {
	delete(w.pointLights, entity)
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

func (w *World) AmbientLights() []AmbientLight {
	lights := make([]AmbientLight, 0, len(w.ambient))

	for _, light := range w.ambient {
		lights = append(lights, light)
	}

	return lights
}

func (w *World) PointLightSources() []PointLightSource {
	sources := make([]PointLightSource, 0, len(w.pointLights))

	for entity, light := range w.pointLights {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		sources = append(sources, PointLightSource{
			Entity:    entity,
			Transform: transform,
			Light:     light,
		})
	}

	slices.SortFunc(sources, func(a, b PointLightSource) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return sources
}

func (w *World) ArmedTanks() []ArmedTank {
	armed := make([]ArmedTank, 0, len(w.weapons))

	for entity, weapon := range w.weapons {
		_, hasDrive := w.drives[entity]
		if !hasDrive {
			continue
		}

		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		intent := w.intents[entity]

		armed = append(armed, ArmedTank{
			Entity:    entity,
			Transform: transform,
			Intent:    intent,
			Weapon:    weapon,
		})
	}

	slices.SortFunc(armed, func(a, b ArmedTank) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return armed
}

func (w *World) ProjectileStates() []ProjectileState {
	states := make([]ProjectileState, 0, len(w.projectiles))

	for entity, projectile := range w.projectiles {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		collider, hasCollider := w.colliders[entity]
		if !hasCollider {
			continue
		}

		states = append(states, ProjectileState{
			Entity:     entity,
			Transform:  transform,
			Projectile: projectile,
			Collider:   collider,
		})
	}

	slices.SortFunc(states, func(a, b ProjectileState) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return states
}

func (w *World) LifetimeStates() map[Entity]Lifetime {
	states := make(map[Entity]Lifetime, len(w.lifetimes))
	for entity, lifetime := range w.lifetimes {
		states[entity] = lifetime
	}

	return states
}

func (w *World) TankHealthStates() []TankHealthState {
	states := make([]TankHealthState, 0, len(w.health))

	for entity := range w.drives {
		health, hasHealth := w.health[entity]
		if !hasHealth {
			continue
		}

		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		states = append(states, TankHealthState{
			Entity:    entity,
			Transform: transform,
			Health:    health,
		})
	}

	slices.SortFunc(states, func(a, b TankHealthState) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return states
}
