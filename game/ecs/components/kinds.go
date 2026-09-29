package components

type Kind string

const (
	KindTransform             Kind = "transform"
	KindSpriteRenderer        Kind = "sprite_renderer"
	KindTankControlIntent     Kind = "tank_control_intent"
	KindTankDrive             Kind = "tank_drive"
	KindTankAcceleration      Kind = "tank_acceleration"
	KindTankDeceleration      Kind = "tank_deceleration"
	KindTankVelocity          Kind = "tank_velocity"
	KindEnemyAI               Kind = "enemy_ai"
	KindWeapon                Kind = "weapon"
	KindProjectile            Kind = "projectile"
	KindLifetime              Kind = "lifetime"
	KindHealth                Kind = "health"
	KindCircleCollider        Kind = "circle_collider"
	KindSolidObstacle         Kind = "solid_obstacle"
	KindMapBounds             Kind = "map_bounds"
	KindAmbientLight          Kind = "ambient_light"
	KindPointLight            Kind = "point_light"
	KindAudioPlayer           Kind = "audio_player"
	KindPlaySoundIntent       Kind = "play_sound_intent"
	KindSoundVolumeByDistance Kind = "sound_volume_by_distance"
	KindEngineAudio           Kind = "engine_audio"

	// pipeline resource that is not an ECS component
	KindCamera Kind = "camera"
)
