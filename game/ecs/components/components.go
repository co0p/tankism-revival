package components

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
	MuzzleFlashOffset         float64
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

type AudioPlayer struct {
	ShootSoundID string
	Volume       float64
}

type PlaySoundIntent struct {
	SoundID string
	Volume  float64
}

type SoundVolumeByDistance struct {
	MaxDistance float64
	MinFactor   float64
}

type EngineAudio struct {
	SoundID       string
	IdleVolume    float64
	MovingVolume  float64
	AttackPerSec  float64
	ReleasePerSec float64
	CurrentVolume float64
}

type CircleCollider struct {
	Radius float64
}

type SolidObstacle struct{}

type MapBounds struct {
	Width  float64
	Height float64
}
