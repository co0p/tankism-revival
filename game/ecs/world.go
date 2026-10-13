package ecs

import (
	"errors"

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
		ecscomponents.EnemyAI |
		ecscomponents.Health |
		ecscomponents.HitFlash |
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

type EnemyAITank struct {
	Entity    Entity
	Transform ecscomponents.Transform
	Drive     ecscomponents.TankDrive
	EnemyAI   ecscomponents.EnemyAI
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
	enemyAI      map[Entity]ecscomponents.EnemyAI
	health       map[Entity]ecscomponents.Health
	hitFlashes   map[Entity]ecscomponents.HitFlash
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
		enemyAI:      make(map[Entity]ecscomponents.EnemyAI),
		health:       make(map[Entity]ecscomponents.Health),
		hitFlashes:   make(map[Entity]ecscomponents.HitFlash),
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
	delete(w.enemyAI, entity)
	delete(w.health, entity)
	delete(w.hitFlashes, entity)
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
	case ecscomponents.EnemyAI:
		w.enemyAI[entity] = c
	case ecscomponents.Health:
		w.health[entity] = c
	case ecscomponents.HitFlash:
		w.hitFlashes[entity] = c
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
