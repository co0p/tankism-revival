package ecs

import (
	"errors"
	"slices"

	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

var ErrEntityNotAlive = errors.New("entity is not alive")

type Entity = ecscomponents.Entity

type Component interface {
	ecscomponents.Transform |
		ecscomponents.SpriteRenderer |
		ecscomponents.TankControlIntent |
		ecscomponents.TankDrive |
		ecscomponents.TankAcceleration |
		ecscomponents.TankDeceleration |
		ecscomponents.TankVelocity |
		ecscomponents.Health |
		ecscomponents.Weapon |
		ecscomponents.Projectile |
		ecscomponents.Lifetime |
		ecscomponents.CircleCollider |
		ecscomponents.SolidObstacle |
		ecscomponents.MapBounds |
		ecscomponents.AmbientLight |
		ecscomponents.PointLight |
		ecscomponents.AudioPlayer |
		ecscomponents.PlaySoundIntent |
		ecscomponents.SoundVolumeByDistance |
		ecscomponents.EngineAudio
}

type ColliderBody struct {
	Entity    Entity
	Transform ecscomponents.Transform
	Collider  ecscomponents.CircleCollider
}

type PointLightSource struct {
	Entity    Entity
	Transform ecscomponents.Transform
	Light     ecscomponents.PointLight
}

type SoundIntent struct {
	Entity Entity
	Intent ecscomponents.PlaySoundIntent
}

type EngineAudioEmitter struct {
	Entity      Entity
	Drive       ecscomponents.TankDrive
	Velocity    ecscomponents.TankVelocity
	EngineAudio ecscomponents.EngineAudio
}

type ArmedTank struct {
	Entity    Entity
	Transform ecscomponents.Transform
	Intent    ecscomponents.TankControlIntent
	Weapon    ecscomponents.Weapon
}

type ProjectileState struct {
	Entity     Entity
	Transform  ecscomponents.Transform
	Projectile ecscomponents.Projectile
	Collider   ecscomponents.CircleCollider
}

type TankHealthState struct {
	Entity    Entity
	Transform ecscomponents.Transform
	Health    ecscomponents.Health
}

type Renderable struct {
	Entity    Entity
	Transform ecscomponents.Transform
	Sprite    ecscomponents.SpriteRenderer
}

type TankMover struct {
	Entity       Entity
	Transform    ecscomponents.Transform
	Drive        ecscomponents.TankDrive
	Acceleration ecscomponents.TankAcceleration
	Deceleration ecscomponents.TankDeceleration
	Velocity     ecscomponents.TankVelocity
	Intent       ecscomponents.TankControlIntent
}

type World struct {
	nextEntityID Entity
	alive        map[Entity]struct{}
	transforms   map[Entity]ecscomponents.Transform
	sprites      map[Entity]ecscomponents.SpriteRenderer
	intents      map[Entity]ecscomponents.TankControlIntent
	drives       map[Entity]ecscomponents.TankDrive
	acceleration map[Entity]ecscomponents.TankAcceleration
	deceleration map[Entity]ecscomponents.TankDeceleration
	velocity     map[Entity]ecscomponents.TankVelocity
	health       map[Entity]ecscomponents.Health
	weapons      map[Entity]ecscomponents.Weapon
	projectiles  map[Entity]ecscomponents.Projectile
	lifetimes    map[Entity]ecscomponents.Lifetime
	colliders    map[Entity]ecscomponents.CircleCollider
	solid        map[Entity]ecscomponents.SolidObstacle
	mapBounds    map[Entity]ecscomponents.MapBounds
	ambient      map[Entity]ecscomponents.AmbientLight
	pointLights  map[Entity]ecscomponents.PointLight
	audioPlayers map[Entity]ecscomponents.AudioPlayer
	soundIntents map[Entity]ecscomponents.PlaySoundIntent
	soundVolumes map[Entity]ecscomponents.SoundVolumeByDistance
	engineAudio  map[Entity]ecscomponents.EngineAudio
}

func NewWorld() *World {
	return &World{
		nextEntityID: 1,
		alive:        make(map[Entity]struct{}),
		transforms:   make(map[Entity]ecscomponents.Transform),
		sprites:      make(map[Entity]ecscomponents.SpriteRenderer),
		intents:      make(map[Entity]ecscomponents.TankControlIntent),
		drives:       make(map[Entity]ecscomponents.TankDrive),
		acceleration: make(map[Entity]ecscomponents.TankAcceleration),
		deceleration: make(map[Entity]ecscomponents.TankDeceleration),
		velocity:     make(map[Entity]ecscomponents.TankVelocity),
		health:       make(map[Entity]ecscomponents.Health),
		weapons:      make(map[Entity]ecscomponents.Weapon),
		projectiles:  make(map[Entity]ecscomponents.Projectile),
		lifetimes:    make(map[Entity]ecscomponents.Lifetime),
		colliders:    make(map[Entity]ecscomponents.CircleCollider),
		solid:        make(map[Entity]ecscomponents.SolidObstacle),
		mapBounds:    make(map[Entity]ecscomponents.MapBounds),
		ambient:      make(map[Entity]ecscomponents.AmbientLight),
		pointLights:  make(map[Entity]ecscomponents.PointLight),
		audioPlayers: make(map[Entity]ecscomponents.AudioPlayer),
		soundIntents: make(map[Entity]ecscomponents.PlaySoundIntent),
		soundVolumes: make(map[Entity]ecscomponents.SoundVolumeByDistance),
		engineAudio:  make(map[Entity]ecscomponents.EngineAudio),
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
	delete(w.mapBounds, entity)
	delete(w.ambient, entity)
	delete(w.pointLights, entity)
	delete(w.audioPlayers, entity)
	delete(w.soundIntents, entity)
	delete(w.soundVolumes, entity)
	delete(w.engineAudio, entity)

	return true
}

func AttachComponent[T Component](w *World, entity Entity, component T) error {
	if !w.IsAlive(entity) {
		return ErrEntityNotAlive
	}

	switch c := any(component).(type) {
	case ecscomponents.Transform:
		w.transforms[entity] = c
	case ecscomponents.SpriteRenderer:
		w.sprites[entity] = c
	case ecscomponents.TankControlIntent:
		w.intents[entity] = c
	case ecscomponents.TankDrive:
		w.drives[entity] = c
	case ecscomponents.TankAcceleration:
		w.acceleration[entity] = c
	case ecscomponents.TankDeceleration:
		w.deceleration[entity] = c
	case ecscomponents.TankVelocity:
		w.velocity[entity] = c
	case ecscomponents.Health:
		w.health[entity] = c
	case ecscomponents.Weapon:
		w.weapons[entity] = c
	case ecscomponents.Projectile:
		w.projectiles[entity] = c
	case ecscomponents.Lifetime:
		w.lifetimes[entity] = c
	case ecscomponents.CircleCollider:
		w.colliders[entity] = c
	case ecscomponents.SolidObstacle:
		w.solid[entity] = c
	case ecscomponents.MapBounds:
		w.mapBounds[entity] = c
	case ecscomponents.AmbientLight:
		w.ambient[entity] = c
	case ecscomponents.PointLight:
		w.pointLights[entity] = c
	case ecscomponents.AudioPlayer:
		w.audioPlayers[entity] = c
	case ecscomponents.PlaySoundIntent:
		w.soundIntents[entity] = c
	case ecscomponents.SoundVolumeByDistance:
		w.soundVolumes[entity] = c
	case ecscomponents.EngineAudio:
		w.engineAudio[entity] = c
	}

	return nil
}

func (w *World) SetTransform(entity Entity, transform ecscomponents.Transform) error {
	return AttachComponent(w, entity, transform)
}

func (w *World) Transform(entity Entity) (ecscomponents.Transform, bool) {
	transform, ok := w.transforms[entity]
	return transform, ok
}

func (w *World) RemoveTransform(entity Entity) {
	delete(w.transforms, entity)
}

func (w *World) SetSpriteRenderer(entity Entity, sprite ecscomponents.SpriteRenderer) error {
	return AttachComponent(w, entity, sprite)
}

func (w *World) SpriteRenderer(entity Entity) (ecscomponents.SpriteRenderer, bool) {
	sprite, ok := w.sprites[entity]
	return sprite, ok
}

func (w *World) RemoveSpriteRenderer(entity Entity) {
	delete(w.sprites, entity)
}

func (w *World) SetTankControlIntent(entity Entity, intent ecscomponents.TankControlIntent) error {
	return AttachComponent(w, entity, intent)
}

func (w *World) TankControlIntent(entity Entity) (ecscomponents.TankControlIntent, bool) {
	intent, ok := w.intents[entity]
	return intent, ok
}

func (w *World) RemoveTankControlIntent(entity Entity) {
	delete(w.intents, entity)
}

func (w *World) SetTankDrive(entity Entity, drive ecscomponents.TankDrive) error {
	return AttachComponent(w, entity, drive)
}

func (w *World) TankDrive(entity Entity) (ecscomponents.TankDrive, bool) {
	drive, ok := w.drives[entity]
	return drive, ok
}

func (w *World) RemoveTankDrive(entity Entity) {
	delete(w.drives, entity)
}

func (w *World) SetTankAcceleration(entity Entity, acceleration ecscomponents.TankAcceleration) error {
	return AttachComponent(w, entity, acceleration)
}

func (w *World) TankAcceleration(entity Entity) (ecscomponents.TankAcceleration, bool) {
	acceleration, ok := w.acceleration[entity]
	return acceleration, ok
}

func (w *World) RemoveTankAcceleration(entity Entity) {
	delete(w.acceleration, entity)
}

func (w *World) SetTankDeceleration(entity Entity, deceleration ecscomponents.TankDeceleration) error {
	return AttachComponent(w, entity, deceleration)
}

func (w *World) TankDeceleration(entity Entity) (ecscomponents.TankDeceleration, bool) {
	deceleration, ok := w.deceleration[entity]
	return deceleration, ok
}

func (w *World) RemoveTankDeceleration(entity Entity) {
	delete(w.deceleration, entity)
}

func (w *World) SetTankVelocity(entity Entity, velocity ecscomponents.TankVelocity) error {
	return AttachComponent(w, entity, velocity)
}

func (w *World) TankVelocity(entity Entity) (ecscomponents.TankVelocity, bool) {
	velocity, ok := w.velocity[entity]
	return velocity, ok
}

func (w *World) RemoveTankVelocity(entity Entity) {
	delete(w.velocity, entity)
}

func (w *World) SetHealth(entity Entity, health ecscomponents.Health) error {
	return AttachComponent(w, entity, health)
}

func (w *World) Health(entity Entity) (ecscomponents.Health, bool) {
	health, ok := w.health[entity]
	return health, ok
}

func (w *World) RemoveHealth(entity Entity) {
	delete(w.health, entity)
}

func (w *World) SetWeapon(entity Entity, weapon ecscomponents.Weapon) error {
	return AttachComponent(w, entity, weapon)
}

func (w *World) Weapon(entity Entity) (ecscomponents.Weapon, bool) {
	weapon, ok := w.weapons[entity]
	return weapon, ok
}

func (w *World) RemoveWeapon(entity Entity) {
	delete(w.weapons, entity)
}

func (w *World) SetProjectile(entity Entity, projectile ecscomponents.Projectile) error {
	return AttachComponent(w, entity, projectile)
}

func (w *World) Projectile(entity Entity) (ecscomponents.Projectile, bool) {
	projectile, ok := w.projectiles[entity]
	return projectile, ok
}

func (w *World) RemoveProjectile(entity Entity) {
	delete(w.projectiles, entity)
}

func (w *World) SetLifetime(entity Entity, lifetime ecscomponents.Lifetime) error {
	return AttachComponent(w, entity, lifetime)
}

func (w *World) Lifetime(entity Entity) (ecscomponents.Lifetime, bool) {
	lifetime, ok := w.lifetimes[entity]
	return lifetime, ok
}

func (w *World) RemoveLifetime(entity Entity) {
	delete(w.lifetimes, entity)
}

func (w *World) SetCircleCollider(entity Entity, collider ecscomponents.CircleCollider) error {
	return AttachComponent(w, entity, collider)
}

func (w *World) CircleCollider(entity Entity) (ecscomponents.CircleCollider, bool) {
	collider, ok := w.colliders[entity]
	return collider, ok
}

func (w *World) RemoveCircleCollider(entity Entity) {
	delete(w.colliders, entity)
}

func (w *World) SetSolidObstacle(entity Entity) error {
	return AttachComponent(w, entity, ecscomponents.SolidObstacle{})
}

func (w *World) IsSolidObstacle(entity Entity) bool {
	_, ok := w.solid[entity]
	return ok
}

func (w *World) RemoveSolidObstacle(entity Entity) {
	delete(w.solid, entity)
}

func (w *World) SetMapBounds(entity Entity, bounds ecscomponents.MapBounds) error {
	return AttachComponent(w, entity, bounds)
}

func (w *World) MapBounds(entity Entity) (ecscomponents.MapBounds, bool) {
	bounds, ok := w.mapBounds[entity]
	return bounds, ok
}

func (w *World) RemoveMapBounds(entity Entity) {
	delete(w.mapBounds, entity)
}

func (w *World) SetAmbientLight(entity Entity, light ecscomponents.AmbientLight) error {
	return AttachComponent(w, entity, light)
}

func (w *World) AmbientLight(entity Entity) (ecscomponents.AmbientLight, bool) {
	light, ok := w.ambient[entity]
	return light, ok
}

func (w *World) RemoveAmbientLight(entity Entity) {
	delete(w.ambient, entity)
}

func (w *World) SetPointLight(entity Entity, light ecscomponents.PointLight) error {
	return AttachComponent(w, entity, light)
}

func (w *World) PointLight(entity Entity) (ecscomponents.PointLight, bool) {
	light, ok := w.pointLights[entity]
	return light, ok
}

func (w *World) RemovePointLight(entity Entity) {
	delete(w.pointLights, entity)
}

func (w *World) SetAudioPlayer(entity Entity, player ecscomponents.AudioPlayer) error {
	return AttachComponent(w, entity, player)
}

func (w *World) AudioPlayer(entity Entity) (ecscomponents.AudioPlayer, bool) {
	player, ok := w.audioPlayers[entity]
	return player, ok
}

func (w *World) RemoveAudioPlayer(entity Entity) {
	delete(w.audioPlayers, entity)
}

func (w *World) SetPlaySoundIntent(entity Entity, intent ecscomponents.PlaySoundIntent) error {
	return AttachComponent(w, entity, intent)
}

func (w *World) PlaySoundIntent(entity Entity) (ecscomponents.PlaySoundIntent, bool) {
	intent, ok := w.soundIntents[entity]
	return intent, ok
}

func (w *World) RemovePlaySoundIntent(entity Entity) {
	delete(w.soundIntents, entity)
}

func (w *World) SoundIntents() []SoundIntent {
	intents := make([]SoundIntent, 0, len(w.soundIntents))

	for entity, intent := range w.soundIntents {
		intents = append(intents, SoundIntent{
			Entity: entity,
			Intent: intent,
		})
	}

	slices.SortFunc(intents, func(a, b SoundIntent) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return intents
}

func (w *World) SetSoundVolumeByDistance(entity Entity, volume ecscomponents.SoundVolumeByDistance) error {
	return AttachComponent(w, entity, volume)
}

func (w *World) SoundVolumeByDistance(entity Entity) (ecscomponents.SoundVolumeByDistance, bool) {
	volume, ok := w.soundVolumes[entity]
	return volume, ok
}

func (w *World) RemoveSoundVolumeByDistance(entity Entity) {
	delete(w.soundVolumes, entity)
}

func (w *World) SetEngineAudio(entity Entity, engineAudio ecscomponents.EngineAudio) error {
	return AttachComponent(w, entity, engineAudio)
}

func (w *World) EngineAudio(entity Entity) (ecscomponents.EngineAudio, bool) {
	engineAudio, ok := w.engineAudio[entity]
	return engineAudio, ok
}

func (w *World) RemoveEngineAudio(entity Entity) {
	delete(w.engineAudio, entity)
}

func (w *World) EngineAudioEmitters() []EngineAudioEmitter {
	emitters := make([]EngineAudioEmitter, 0, len(w.engineAudio))

	for entity, engineAudio := range w.engineAudio {
		drive, hasDrive := w.drives[entity]
		if !hasDrive {
			continue
		}

		velocity, hasVelocity := w.velocity[entity]
		if !hasVelocity {
			continue
		}

		emitters = append(emitters, EngineAudioEmitter{
			Entity:      entity,
			Drive:       drive,
			Velocity:    velocity,
			EngineAudio: engineAudio,
		})
	}

	slices.SortFunc(emitters, func(a, b EngineAudioEmitter) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return emitters
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

func (w *World) AmbientLights() []ecscomponents.AmbientLight {
	lights := make([]ecscomponents.AmbientLight, 0, len(w.ambient))

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

func (w *World) LifetimeStates() map[Entity]ecscomponents.Lifetime {
	states := make(map[Entity]ecscomponents.Lifetime, len(w.lifetimes))
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
