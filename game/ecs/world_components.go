package ecs

import (
	"slices"

	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func (w *World) SetTransform(entity Entity, transform ecscomponents.Transform) error {
	return AttachComponent(w, entity, transform)
}

func (w *World) Transform(entity Entity) (ecscomponents.Transform, bool) {
	transform, ok := w.transforms[entity]
	return transform, ok
}

func (w *World) SetSpriteRenderer(entity Entity, sprite ecscomponents.SpriteRenderer) error {
	return AttachComponent(w, entity, sprite)
}

func (w *World) SpriteRenderer(entity Entity) (ecscomponents.SpriteRenderer, bool) {
	sprite, ok := w.sprites[entity]
	return sprite, ok
}

func (w *World) SetTankControlIntent(entity Entity, intent ecscomponents.TankControlIntent) error {
	return AttachComponent(w, entity, intent)
}

func (w *World) TankControlIntent(entity Entity) (ecscomponents.TankControlIntent, bool) {
	intent, ok := w.intents[entity]
	return intent, ok
}

func (w *World) SetTankDrive(entity Entity, drive ecscomponents.TankDrive) error {
	return AttachComponent(w, entity, drive)
}

func (w *World) TankDrive(entity Entity) (ecscomponents.TankDrive, bool) {
	drive, ok := w.drives[entity]
	return drive, ok
}

func (w *World) SetTankAcceleration(entity Entity, acceleration ecscomponents.TankAcceleration) error {
	return AttachComponent(w, entity, acceleration)
}

func (w *World) TankAcceleration(entity Entity) (ecscomponents.TankAcceleration, bool) {
	acceleration, ok := w.acceleration[entity]
	return acceleration, ok
}

func (w *World) SetTankDeceleration(entity Entity, deceleration ecscomponents.TankDeceleration) error {
	return AttachComponent(w, entity, deceleration)
}

func (w *World) TankDeceleration(entity Entity) (ecscomponents.TankDeceleration, bool) {
	deceleration, ok := w.deceleration[entity]
	return deceleration, ok
}

func (w *World) SetTankVelocity(entity Entity, velocity ecscomponents.TankVelocity) error {
	return AttachComponent(w, entity, velocity)
}

func (w *World) TankVelocity(entity Entity) (ecscomponents.TankVelocity, bool) {
	velocity, ok := w.velocity[entity]
	return velocity, ok
}

func (w *World) SetEnemyAI(entity Entity, enemyAI ecscomponents.EnemyAI) error {
	return AttachComponent(w, entity, enemyAI)
}

func (w *World) EnemyAI(entity Entity) (ecscomponents.EnemyAI, bool) {
	enemyAI, ok := w.enemyAI[entity]
	return enemyAI, ok
}

func (w *World) SetHealth(entity Entity, health ecscomponents.Health) error {
	return AttachComponent(w, entity, health)
}

func (w *World) Health(entity Entity) (ecscomponents.Health, bool) {
	health, ok := w.health[entity]
	return health, ok
}

func (w *World) SetWeapon(entity Entity, weapon ecscomponents.Weapon) error {
	return AttachComponent(w, entity, weapon)
}

func (w *World) Weapon(entity Entity) (ecscomponents.Weapon, bool) {
	weapon, ok := w.weapons[entity]
	return weapon, ok
}

func (w *World) SetProjectile(entity Entity, projectile ecscomponents.Projectile) error {
	return AttachComponent(w, entity, projectile)
}

func (w *World) Projectile(entity Entity) (ecscomponents.Projectile, bool) {
	projectile, ok := w.projectiles[entity]
	return projectile, ok
}

func (w *World) SetLifetime(entity Entity, lifetime ecscomponents.Lifetime) error {
	return AttachComponent(w, entity, lifetime)
}

func (w *World) Lifetime(entity Entity) (ecscomponents.Lifetime, bool) {
	lifetime, ok := w.lifetimes[entity]
	return lifetime, ok
}

func (w *World) SetCircleCollider(entity Entity, collider ecscomponents.CircleCollider) error {
	return AttachComponent(w, entity, collider)
}

func (w *World) CircleCollider(entity Entity) (ecscomponents.CircleCollider, bool) {
	collider, ok := w.colliders[entity]
	return collider, ok
}

func (w *World) SetSolidObstacle(entity Entity) error {
	return AttachComponent(w, entity, ecscomponents.SolidObstacle{})
}

func (w *World) IsSolidObstacle(entity Entity) bool {
	_, ok := w.solid[entity]
	return ok
}

func (w *World) SetMapBounds(entity Entity, bounds ecscomponents.MapBounds) error {
	return AttachComponent(w, entity, bounds)
}

func (w *World) MapBounds(entity Entity) (ecscomponents.MapBounds, bool) {
	bounds, ok := w.mapBounds[entity]
	return bounds, ok
}

func (w *World) SetAmbientLight(entity Entity, light ecscomponents.AmbientLight) error {
	return AttachComponent(w, entity, light)
}

func (w *World) AmbientLight(entity Entity) (ecscomponents.AmbientLight, bool) {
	light, ok := w.ambient[entity]
	return light, ok
}

func (w *World) SetPointLight(entity Entity, light ecscomponents.PointLight) error {
	return AttachComponent(w, entity, light)
}

func (w *World) PointLight(entity Entity) (ecscomponents.PointLight, bool) {
	light, ok := w.pointLights[entity]
	return light, ok
}

func (w *World) SetAudioPlayer(entity Entity, player ecscomponents.AudioPlayer) error {
	return AttachComponent(w, entity, player)
}

func (w *World) AudioPlayer(entity Entity) (ecscomponents.AudioPlayer, bool) {
	player, ok := w.audioPlayers[entity]
	return player, ok
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

func (w *World) SetEngineAudio(entity Entity, engineAudio ecscomponents.EngineAudio) error {
	return AttachComponent(w, entity, engineAudio)
}

func (w *World) EngineAudio(entity Entity) (ecscomponents.EngineAudio, bool) {
	engineAudio, ok := w.engineAudio[entity]
	return engineAudio, ok
}

func (w *World) EngineAudioEmitters() []EngineAudioEmitter {
	emitters := make([]EngineAudioEmitter, 0, len(w.engineAudio))

	for entity, engineAudio := range w.engineAudio {
		drive, hasDrive := w.drives[entity]
		if !hasDrive {
			continue
		}

		velocity := w.velocity[entity]

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
