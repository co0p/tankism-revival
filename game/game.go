package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	gameinput "github.com/co0p/tankism-revival/game/input"
	gamerender "github.com/co0p/tankism-revival/game/render"
	"github.com/co0p/tankism-revival/game/sim"
	"github.com/co0p/tankism-revival/game/worldseed"
	libaudio "github.com/co0p/tankism-revival/lib/audio"
	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/co0p/tankism-revival/lib/sprites"
	gameassets "github.com/co0p/tankism-revival/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	fixedDeltaSeconds = 1.0 / 60.0
	DefaultMapPath    = "resources/maps/mapism_default.tmj"

	defaultPlayerMoveSpeed          = 160
	defaultPlayerTurnSpeedRadiansPS = math.Pi
	defaultPlayerAcceleration       = 240
	defaultPlayerCoastDeceleration  = 180
	defaultPlayerBrakeDeceleration  = 420
	defaultEnemyMoveSpeed           = 105
	defaultEnemyTurnSpeedRadiansPS  = 1.9
	defaultEnemyAcceleration        = 150
	defaultEnemyCoastDeceleration   = 140
	defaultEnemyBrakeDeceleration   = 240
	defaultEnemyStopDistance        = 145.0

	defaultPlayerWeaponCooldown            = 0.40
	defaultPlayerProjectileSpeed           = 360.0
	defaultPlayerProjectileLife            = 1.70
	defaultPlayerProjectileDamage          = 12.0
	defaultPlayerMuzzleOffset              = 22.0
	defaultPlayerProjectileRadius          = 4.0
	defaultPlayerProjectileLightRadius     = 56.0
	defaultPlayerProjectileLightIntensity  = 0.34
	defaultPlayerMuzzleFlashLifetime       = 0.07
	defaultPlayerMuzzleFlashOffset         = 16.0
	defaultPlayerMuzzleFlashLightRadius    = 64.0
	defaultPlayerMuzzleFlashLightIntensity = 0.52
	defaultPlayerAccuracySpreadBase        = 0.045
	defaultPlayerAccuracyMinSpread         = 0.006
	defaultPlayerAccuracyMovingSpread      = 0.12
	defaultPlayerAccuracyStoppedMultiplier = 0.45
	defaultPlayerAccuracyRelaxedMultiplier = 0.85

	defaultEnemyWeaponCooldown            = 0.92
	defaultEnemyProjectileSpeed           = 300.0
	defaultEnemyProjectileLife            = 1.45
	defaultEnemyProjectileDamage          = 6.0
	defaultEnemyMuzzleOffset              = 22.0
	defaultEnemyProjectileRadius          = 4.0
	defaultEnemyProjectileLightRadius     = 48.0
	defaultEnemyProjectileLightIntensity  = 0.24
	defaultEnemyMuzzleFlashLifetime       = 0.07
	defaultEnemyMuzzleFlashOffset         = 16.0
	defaultEnemyMuzzleFlashLightRadius    = 56.0
	defaultEnemyMuzzleFlashLightIntensity = 0.44
	defaultEnemyAccuracySpreadBase        = 0.085
	defaultEnemyAccuracyMinSpread         = 0.010
	defaultEnemyAccuracyMovingSpread      = 0.18
	defaultEnemyAccuracyStoppedMultiplier = 0.50
	defaultEnemyAccuracyRelaxedMultiplier = 0.70
	defaultAmbientDarkness                = 0.70
	defaultPlayerLightRadius              = 170.0
	defaultPlayerLightIntensity           = 0.20
	defaultPlayerBackLightRadius          = 82.0
	defaultPlayerBackLightIntensity       = 0.10
	defaultPlayerBackLightDistance        = 12.0

	playerCollisionRadius       = 16.0
	enemyCollisionRadius        = 16.0
	barrelCollisionRadius       = 11.0
	fenceCollisionRadius        = 14.0
	crateCollisionRadius        = 13.0
	treeLargeCollisionRadius    = 12.0
	treeSmallCollisionRadius    = 8.0
	wireStraightCollisionRadius = 12.0
	wireCrookedCollisionRadius  = 16.0
	tankCrossCollisionRadius    = 12.0

	playerMaxHealth = 100.0
	enemyMaxHealth  = 70.0

	defaultAudioSampleRate = 44100

	defaultPlayerEngineIdleVolume    = 0.04
	defaultPlayerEngineMovingVolume  = 0.18
	defaultPlayerEngineAttackPerSec  = 1.8
	defaultPlayerEngineReleasePerSec = 1.0
	defaultEnemyEngineIdleVolume     = 0.04
	defaultEnemyEngineMovingVolume   = 0.18
	defaultEnemyEngineAttackPerSec   = 1.4
	defaultEnemyEngineReleasePerSec  = 0.8

	defaultExplosionLifetime       = 0.22
	defaultExplosionBurstLifetime  = 0.10
	defaultExplosionSmokeLifetime  = 0.48
	defaultExplosionPlumeLifetime  = 0.62
	defaultExplosionLightRadius    = 86.0
	defaultExplosionLightIntensity = 0.60
	defaultBurstLightRadius        = 118.0
	defaultBurstLightIntensity     = 0.92

	defaultGameOverShakeDurationSeconds = 0.32
	defaultGameOverShakeAmplitude       = 9.0
	defaultGameOverShakeFrequency       = 34.0
)

var (
	outsideMapColor = color.RGBA{R: 0x1c, G: 0x27, B: 0x1b, A: 0xff}
)

type Game struct {
	world                 *ecs.World
	mapEntity             ecs.Entity
	camera                Camera
	playerEntity          ecs.Entity
	enemyEntities         []ecs.Entity
	playerBackLightEntity ecs.Entity
	renderer              *gamerender.WorldRenderer
	soundMixer            *libaudio.EffectMixer
	activeEngineLoops     map[string]struct{}
	gameOver              bool
	gameOverShakeSeconds  float64
	gameOverShakePhase    float64
	gameOverShakeX        float64
	gameOverShakeY        float64
}

func NewGame() (*Game, error) {
	return NewGameWithMap(DefaultMapPath)
}

func NewGameWithMap(mapPath string) (*Game, error) {
	if mapPath == "" {
		mapPath = DefaultMapPath
	}

	spriteAssets := make([]sprites.Asset, 0, len(gamepalette.Sprites)+5)
	for _, spriteDef := range gamepalette.Sprites {
		spriteAssets = append(spriteAssets, sprites.Asset{ID: spriteDef.ID, Path: spriteDef.Path})
	}
	spriteAssets = append(spriteAssets,
		sprites.Asset{ID: SpriteExplosionCore, Path: gameassets.Explosion3Path},
		sprites.Asset{ID: SpriteExplosionBurst, Path: gameassets.Explosion1Path},
		sprites.Asset{ID: SpriteExplosionSmoke, Path: gameassets.ExplosionSmoke3Path},
		sprites.Asset{ID: SpriteExplosionPlume, Path: gameassets.ExplosionSmoke1Path},
		sprites.Asset{ID: SpriteEnemyWreckHull, Path: gameassets.TankBodyDarkPath},
	)

	spriteImages, err := sprites.Load(gameassets.FS, spriteAssets)
	if err != nil {
		return nil, err
	}

	soundBank, err := libaudio.LoadSoundBank(gameassets.FS, defaultAudioSampleRate, []libaudio.SoundAsset{
		{ID: SoundCannonFire, Path: gameassets.CannonFireSFXPath},
		{ID: SoundEngineHeavyLoop, Path: gameassets.EngineLoopSFXPath},
	})
	if err != nil {
		return nil, err
	}

	soundMixer := libaudio.NewEffectMixer(defaultAudioSampleRate, soundBank)

	renderer := gamerender.NewWorldRenderer(spriteImages, gamerender.Config{
		OutsideMapColor:     outsideMapColor,
		LightSpriteSize:     256,
		LightIntensityScale: 0.65,
		HealthBarWidth:      34,
		HealthBarHeight:     4,
		HealthBarYOffset:    30,
		HealthBarBorderSize: 1,
		ForegroundMinZIndex: ZIndexObstacles,
	})

	world := ecs.NewWorld()
	seedResult, err := worldseed.SeedFromTiledMap(world, worldseed.TiledSeedConfig{
		Path:            mapPath,
		GroundZIndex:    ZIndexMap,
		PropsZIndex:     ZIndexObstacles,
		PropsTallZIndex: ZIndexPropsTall,
		ColliderRadiusBySprite: map[string]float64{
			SpriteObstacleBarrelGrayTop:  barrelCollisionRadius,
			SpriteObstacleFenceRed:       fenceCollisionRadius,
			SpriteObstacleCrateWood:      crateCollisionRadius,
			SpriteObstacleTreeGreenLarge: treeLargeCollisionRadius,
			SpriteObstacleTreeGreenSmall: treeSmallCollisionRadius,
			SpriteObstacleTreeBrownLarge: treeLargeCollisionRadius,
			SpriteObstacleTreeBrownSmall: treeSmallCollisionRadius,
			SpriteObstacleWireStraight:   wireStraightCollisionRadius,
			SpriteObstacleWireCrooked:    wireCrookedCollisionRadius,
			SpriteObstacleTankCrossMetal: tankCrossCollisionRadius,
		},
		DefaultEnemyOffsetX: 280,
		DefaultEnemyOffsetY: -220,
	})
	if err != nil {
		return nil, fmt.Errorf("seed map from %q: %w", mapPath, err)
	}

	worldMap := Map{Width: seedResult.MapWidth, Height: seedResult.MapHeight}
	mapEntity := world.CreateEntity()
	if err := world.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: worldMap.Width, Height: worldMap.Height}); err != nil {
		return nil, fmt.Errorf("set map bounds: %w", err)
	}

	player := world.CreateEntity()
	spawnX := seedResult.PlayerSpawnX
	spawnY := seedResult.PlayerSpawnY

	if err := world.SetTransform(player, ecscomponents.Transform{
		X:               spawnX,
		Y:               spawnY,
		RotationRadians: 0,
	}); err != nil {
		return nil, fmt.Errorf("set player transform: %w", err)
	}

	if err := world.SetSpriteRenderer(player, ecscomponents.SpriteRenderer{
		SpriteID:              SpriteTankRedSmall,
		ZIndex:                ZIndexActors,
		RotationOffsetRadians: math.Pi,
	}); err != nil {
		return nil, fmt.Errorf("set player sprite renderer: %w", err)
	}

	if err := world.SetTankDrive(player, ecscomponents.TankDrive{
		MoveSpeed:          defaultPlayerMoveSpeed,
		TurnSpeedRadiansPS: defaultPlayerTurnSpeedRadiansPS,
	}); err != nil {
		return nil, fmt.Errorf("set player tank drive: %w", err)
	}

	if err := world.SetTankAcceleration(player, ecscomponents.TankAcceleration{
		ForwardUnitsPerSecondSquared: defaultPlayerAcceleration,
	}); err != nil {
		return nil, fmt.Errorf("set player tank acceleration: %w", err)
	}

	if err := world.SetTankDeceleration(player, ecscomponents.TankDeceleration{
		CoastUnitsPerSecondSquared: defaultPlayerCoastDeceleration,
		BrakeUnitsPerSecondSquared: defaultPlayerBrakeDeceleration,
	}); err != nil {
		return nil, fmt.Errorf("set player tank deceleration: %w", err)
	}

	if err := world.SetTankVelocity(player, ecscomponents.TankVelocity{}); err != nil {
		return nil, fmt.Errorf("set player tank velocity: %w", err)
	}

	if err := world.SetTankControlIntent(player, ecscomponents.TankControlIntent{}); err != nil {
		return nil, fmt.Errorf("set player tank control intent: %w", err)
	}

	if err := world.SetCircleCollider(player, ecscomponents.CircleCollider{Radius: playerCollisionRadius}); err != nil {
		return nil, fmt.Errorf("set player collider: %w", err)
	}
	if err := world.SetHealth(player, ecscomponents.Health{Current: playerMaxHealth, Max: playerMaxHealth}); err != nil {
		return nil, fmt.Errorf("set player health: %w", err)
	}

	if err := world.SetPointLight(player, ecscomponents.PointLight{
		Radius:    defaultPlayerLightRadius,
		Intensity: defaultPlayerLightIntensity,
		Color:     ecscomponents.LightColor{R: 1, G: 0.93, B: 0.82},
	}); err != nil {
		return nil, fmt.Errorf("set player point light: %w", err)
	}
	if err := world.SetWeapon(player, ecscomponents.Weapon{
		CooldownSeconds:           defaultPlayerWeaponCooldown,
		ProjectileSpeed:           defaultPlayerProjectileSpeed,
		ProjectileLifetime:        defaultPlayerProjectileLife,
		ProjectileDamage:          defaultPlayerProjectileDamage,
		MuzzleOffset:              defaultPlayerMuzzleOffset,
		AccuracySpreadRadians:     defaultPlayerAccuracySpreadBase,
		AccuracyMinSpreadRadians:  defaultPlayerAccuracyMinSpread,
		AccuracyMovingSpread:      defaultPlayerAccuracyMovingSpread,
		AccuracyStoppedMultiplier: defaultPlayerAccuracyStoppedMultiplier,
		AccuracyRelaxedMultiplier: defaultPlayerAccuracyRelaxedMultiplier,
		ProjectileSpriteID:        SpriteProjectileRed,
		ProjectileRadius:          defaultPlayerProjectileRadius,
		ProjectileLightRadius:     defaultPlayerProjectileLightRadius,
		ProjectileLightIntensity:  defaultPlayerProjectileLightIntensity,
		ProjectileLightColor:      ecscomponents.LightColor{R: 1.00, G: 0.62, B: 0.38},
		MuzzleFlashSpriteID:       SpriteShotRed,
		MuzzleFlashLifetime:       defaultPlayerMuzzleFlashLifetime,
		MuzzleFlashOffset:         defaultPlayerMuzzleFlashOffset,
		MuzzleFlashLightRadius:    defaultPlayerMuzzleFlashLightRadius,
		MuzzleFlashLightIntensity: defaultPlayerMuzzleFlashLightIntensity,
		MuzzleFlashLightColor:     ecscomponents.LightColor{R: 1.00, G: 0.68, B: 0.45},
	}); err != nil {
		return nil, fmt.Errorf("set player weapon: %w", err)
	}
	if err := world.SetAudioPlayer(player, ecscomponents.AudioPlayer{ShootSoundID: SoundCannonFire, Volume: 0.55}); err != nil {
		return nil, fmt.Errorf("set player audio player: %w", err)
	}
	if err := world.SetSoundVolumeByDistance(player, ecscomponents.SoundVolumeByDistance{MaxDistance: 900, MinFactor: 0.25}); err != nil {
		return nil, fmt.Errorf("set player sound volume by distance: %w", err)
	}
	if err := world.SetEngineAudio(player, ecscomponents.EngineAudio{
		SoundID:       SoundEngineHeavyLoop,
		IdleVolume:    defaultPlayerEngineIdleVolume,
		MovingVolume:  defaultPlayerEngineMovingVolume,
		AttackPerSec:  defaultPlayerEngineAttackPerSec,
		ReleasePerSec: defaultPlayerEngineReleasePerSec,
		CurrentVolume: defaultPlayerEngineIdleVolume,
	}); err != nil {
		return nil, fmt.Errorf("set player engine audio: %w", err)
	}

	playerBackLight := world.CreateEntity()
	backLightX, backLightY := sim.TankBackLightPosition(spawnX, spawnY, 0, defaultPlayerBackLightDistance)
	if err := world.SetTransform(playerBackLight, ecscomponents.Transform{
		X:               backLightX,
		Y:               backLightY,
		RotationRadians: 0,
	}); err != nil {
		return nil, fmt.Errorf("set player back light transform: %w", err)
	}
	if err := world.SetPointLight(playerBackLight, ecscomponents.PointLight{
		Radius:    defaultPlayerBackLightRadius,
		Intensity: defaultPlayerBackLightIntensity,
		Color:     ecscomponents.LightColor{R: 0.75, G: 0.80, B: 0.95},
	}); err != nil {
		return nil, fmt.Errorf("set player back light: %w", err)
	}

	ambient := world.CreateEntity()
	if err := world.SetAmbientLight(ambient, ecscomponents.AmbientLight{Darkness: defaultAmbientDarkness}); err != nil {
		return nil, fmt.Errorf("set ambient light: %w", err)
	}

	spawnEnemy := func(spawnX, spawnY float64) (ecs.Entity, error) {
		enemy := world.CreateEntity()
		if err := world.SetTransform(enemy, ecscomponents.Transform{
			X:               spawnX,
			Y:               spawnY,
			RotationRadians: math.Pi,
		}); err != nil {
			return 0, fmt.Errorf("set enemy transform: %w", err)
		}
		if err := world.SetSpriteRenderer(enemy, ecscomponents.SpriteRenderer{
			SpriteID:              SpriteTankBlueSmall,
			ZIndex:                ZIndexActors,
			RotationOffsetRadians: math.Pi,
		}); err != nil {
			return 0, fmt.Errorf("set enemy sprite renderer: %w", err)
		}
		if err := world.SetTankDrive(enemy, ecscomponents.TankDrive{
			MoveSpeed:          defaultEnemyMoveSpeed,
			TurnSpeedRadiansPS: defaultEnemyTurnSpeedRadiansPS,
		}); err != nil {
			return 0, fmt.Errorf("set enemy tank drive: %w", err)
		}
		if err := world.SetTankAcceleration(enemy, ecscomponents.TankAcceleration{
			ForwardUnitsPerSecondSquared: defaultEnemyAcceleration,
		}); err != nil {
			return 0, fmt.Errorf("set enemy tank acceleration: %w", err)
		}
		if err := world.SetTankDeceleration(enemy, ecscomponents.TankDeceleration{
			CoastUnitsPerSecondSquared: defaultEnemyCoastDeceleration,
			BrakeUnitsPerSecondSquared: defaultEnemyBrakeDeceleration,
		}); err != nil {
			return 0, fmt.Errorf("set enemy tank deceleration: %w", err)
		}
		if err := world.SetTankVelocity(enemy, ecscomponents.TankVelocity{}); err != nil {
			return 0, fmt.Errorf("set enemy tank velocity: %w", err)
		}
		if err := world.SetTankControlIntent(enemy, ecscomponents.TankControlIntent{}); err != nil {
			return 0, fmt.Errorf("set enemy tank control intent: %w", err)
		}
		if err := world.SetEnemyAI(enemy, sim.NewEnemyAI(enemy)); err != nil {
			return 0, fmt.Errorf("set enemy ai: %w", err)
		}
		if err := world.SetCircleCollider(enemy, ecscomponents.CircleCollider{Radius: enemyCollisionRadius}); err != nil {
			return 0, fmt.Errorf("set enemy collider: %w", err)
		}
		if err := world.SetHealth(enemy, ecscomponents.Health{Current: enemyMaxHealth, Max: enemyMaxHealth}); err != nil {
			return 0, fmt.Errorf("set enemy health: %w", err)
		}
		if err := world.SetPointLight(enemy, ecscomponents.PointLight{
			Radius:    120,
			Intensity: 0.10,
			Color:     ecscomponents.LightColor{R: 0.7, G: 0.82, B: 1.0},
		}); err != nil {
			return 0, fmt.Errorf("set enemy point light: %w", err)
		}
		if err := world.SetWeapon(enemy, ecscomponents.Weapon{
			CooldownSeconds:           defaultEnemyWeaponCooldown,
			ProjectileSpeed:           defaultEnemyProjectileSpeed,
			ProjectileLifetime:        defaultEnemyProjectileLife,
			ProjectileDamage:          defaultEnemyProjectileDamage,
			MuzzleOffset:              defaultEnemyMuzzleOffset,
			AccuracySpreadRadians:     defaultEnemyAccuracySpreadBase,
			AccuracyMinSpreadRadians:  defaultEnemyAccuracyMinSpread,
			AccuracyMovingSpread:      defaultEnemyAccuracyMovingSpread,
			AccuracyStoppedMultiplier: defaultEnemyAccuracyStoppedMultiplier,
			AccuracyRelaxedMultiplier: defaultEnemyAccuracyRelaxedMultiplier,
			ProjectileSpriteID:        SpriteProjectileBlue,
			ProjectileRadius:          defaultEnemyProjectileRadius,
			ProjectileLightRadius:     defaultEnemyProjectileLightRadius,
			ProjectileLightIntensity:  defaultEnemyProjectileLightIntensity,
			ProjectileLightColor:      ecscomponents.LightColor{R: 0.74, G: 0.86, B: 1.00},
			MuzzleFlashSpriteID:       SpriteShotOrange,
			MuzzleFlashLifetime:       defaultEnemyMuzzleFlashLifetime,
			MuzzleFlashOffset:         defaultEnemyMuzzleFlashOffset,
			MuzzleFlashLightRadius:    defaultEnemyMuzzleFlashLightRadius,
			MuzzleFlashLightIntensity: defaultEnemyMuzzleFlashLightIntensity,
			MuzzleFlashLightColor:     ecscomponents.LightColor{R: 0.96, G: 0.80, B: 0.56},
		}); err != nil {
			return 0, fmt.Errorf("set enemy weapon: %w", err)
		}
		if err := world.SetAudioPlayer(enemy, ecscomponents.AudioPlayer{ShootSoundID: SoundCannonFire, Volume: 0.45}); err != nil {
			return 0, fmt.Errorf("set enemy audio player: %w", err)
		}
		if err := world.SetSoundVolumeByDistance(enemy, ecscomponents.SoundVolumeByDistance{MaxDistance: 900, MinFactor: 0.10}); err != nil {
			return 0, fmt.Errorf("set enemy sound volume by distance: %w", err)
		}
		if err := world.SetEngineAudio(enemy, ecscomponents.EngineAudio{
			SoundID:       SoundEngineHeavyLoop,
			IdleVolume:    defaultEnemyEngineIdleVolume,
			MovingVolume:  defaultEnemyEngineMovingVolume,
			AttackPerSec:  defaultEnemyEngineAttackPerSec,
			ReleasePerSec: defaultEnemyEngineReleasePerSec,
			CurrentVolume: defaultEnemyEngineIdleVolume,
		}); err != nil {
			return 0, fmt.Errorf("set enemy engine audio: %w", err)
		}

		return enemy, nil
	}

	enemySpawns := seedResult.EnemySpawns
	if len(enemySpawns) == 0 {
		enemySpawns = []worldseed.SpawnPoint{{X: seedResult.EnemySpawnX, Y: seedResult.EnemySpawnY}}
	}

	enemyEntities := make([]ecs.Entity, 0, len(enemySpawns))
	for _, spawn := range enemySpawns {
		enemySpawnX, enemySpawnY := worldMap.ClampPosition(spawn.X, spawn.Y)
		enemy, err := spawnEnemy(enemySpawnX, enemySpawnY)
		if err != nil {
			return nil, err
		}
		enemyEntities = append(enemyEntities, enemy)
	}

	camera := NewCamera(float64(LogicalWidth), float64(LogicalHeight))
	camera.FollowTarget(worldMap.Width, worldMap.Height, spawnX, spawnY)

	return &Game{
		world:                 world,
		mapEntity:             mapEntity,
		camera:                camera,
		playerEntity:          player,
		enemyEntities:         enemyEntities,
		playerBackLightEntity: playerBackLight,
		renderer:              renderer,
		soundMixer:            soundMixer,
		activeEngineLoops:     make(map[string]struct{}),
	}, nil
}

func (g *Game) Update() error {
	return g.updateGame(fixedDeltaSeconds)
}

func (g *Game) updateGame(deltaSeconds float64) error {
	if g.gameOver {
		g.updateGameOverShake(deltaSeconds)
		return nil
	}

	intent := gameinput.ReadPlayerTankIntent()
	if err := g.world.SetTankControlIntent(g.playerEntity, intent); err != nil {
		return fmt.Errorf("set player intent: %w", err)
	}

	if err := sim.RunDefaultPipeline(g.world, deltaSeconds, sim.PipelineConfig{
		PlayerEntity:        g.playerEntity,
		MapEntity:           g.mapEntity,
		AudioListenerEntity: g.playerEntity,
		EnemyStopDistance:   defaultEnemyStopDistance,
		ProjectileZIndex:    ZIndexProjectiles,
		MuzzleFlashZIndex:   ZIndexEffects,
		TankDestruction: sim.TankDestructionConfig{
			ExplosionSpriteID:        SpriteExplosionCore,
			ExplosionBurstSpriteID:   SpriteExplosionBurst,
			ExplosionSmokeSpriteID:   SpriteExplosionSmoke,
			ExplosionPlumeSpriteID:   SpriteExplosionPlume,
			EnemyWreckSpriteID:       SpriteEnemyWreckHull,
			ExplosionZIndex:          ZIndexEffects,
			EnemyWreckZIndex:         ZIndexActors,
			ExplosionLifetime:        defaultExplosionLifetime,
			ExplosionBurstLifetime:   defaultExplosionBurstLifetime,
			ExplosionSmokeLifetime:   defaultExplosionSmokeLifetime,
			ExplosionPlumeLifetime:   defaultExplosionPlumeLifetime,
			ExplosionLightRadius:     defaultExplosionLightRadius,
			ExplosionLightIntensity:  defaultExplosionLightIntensity,
			BurstLightRadius:         defaultBurstLightRadius,
			BurstLightIntensity:      defaultBurstLightIntensity,
			ExplosionLightColor:      ecscomponents.LightColor{R: 1.0, G: 0.62, B: 0.36},
			EnemyWreckColliderRadius: enemyCollisionRadius,
		},
		Camera:                  &g.camera,
		PlayerBackLightEntity:   g.playerBackLightEntity,
		PlayerBackLightDistance: defaultPlayerBackLightDistance,
	}); err != nil {
		return fmt.Errorf("run simulation pipeline: %w", err)
	}

	if g.playerDestroyed() {
		g.setGameOver()
		return nil
	}

	if err := g.syncEngineLoops(); err != nil {
		return fmt.Errorf("sync engine loops: %w", err)
	}

	if err := g.playQueuedSounds(); err != nil {
		return fmt.Errorf("play queued sounds: %w", err)
	}

	return nil
}

func (g *Game) syncEngineLoops() error {
	if g.soundMixer == nil {
		return nil
	}

	activeLoops := make(map[string]struct{})

	for _, emitter := range g.world.EngineAudioEmitters() {
		engineAudio := emitter.EngineAudio
		if engineAudio.SoundID == "" {
			continue
		}

		loopID := fmt.Sprintf("engine_%d", emitter.Entity)
		activeLoops[loopID] = struct{}{}

		if err := g.soundMixer.SetLoop(loopID, engineAudio.SoundID, engineAudio.CurrentVolume); err != nil {
			return fmt.Errorf("set engine loop for entity %d: %w", emitter.Entity, err)
		}
	}

	for loopID := range g.activeEngineLoops {
		if _, keep := activeLoops[loopID]; keep {
			continue
		}

		g.soundMixer.RemoveLoop(loopID)
	}

	g.activeEngineLoops = activeLoops
	return nil
}

func (g *Game) playQueuedSounds() error {
	if g.soundMixer == nil {
		return nil
	}

	for _, queued := range g.world.SoundIntents() {
		if queued.Intent.SoundID == "" {
			g.world.RemovePlaySoundIntent(queued.Entity)
			continue
		}
		if queued.Intent.Volume <= 0 {
			g.world.RemovePlaySoundIntent(queued.Entity)
			continue
		}

		if err := g.soundMixer.Play(queued.Intent.SoundID, queued.Intent.Volume); err != nil {
			return fmt.Errorf("play sound %q: %w", queued.Intent.SoundID, err)
		}

		g.world.RemovePlaySoundIntent(queued.Entity)
	}

	g.soundMixer.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.updateDraw(screen, 1.0)
}

func (g *Game) updateDraw(screen *ebiten.Image, interpolationAlpha float64) {
	_ = interpolationAlpha
	drawCamera := g.camera
	drawCamera.X += g.gameOverShakeX
	drawCamera.Y += g.gameOverShakeY
	g.renderer.Draw(screen, g.world, &drawCamera)

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER - player tank destroyed", LogicalWidth/2-130, LogicalHeight/2)
	}
}

func (g *Game) playerDestroyed() bool {
	if !g.world.IsAlive(g.playerEntity) {
		return true
	}

	health, ok := g.world.Health(g.playerEntity)
	if !ok {
		return true
	}

	return health.Current <= 0 || health.Max <= 0
}

func (g *Game) setGameOver() {
	if g.gameOver {
		return
	}

	g.gameOver = true
	g.gameOverShakeSeconds = defaultGameOverShakeDurationSeconds
	g.gameOverShakePhase = 0
	g.gameOverShakeX = 0
	g.gameOverShakeY = 0

	for _, queued := range g.world.SoundIntents() {
		g.world.RemovePlaySoundIntent(queued.Entity)
	}

	if g.soundMixer == nil {
		return
	}

	for loopID := range g.activeEngineLoops {
		g.soundMixer.RemoveLoop(loopID)
	}
	g.activeEngineLoops = make(map[string]struct{})
	g.soundMixer.Update()
}

func (g *Game) updateGameOverShake(deltaSeconds float64) {
	if g.gameOverShakeSeconds <= 0 {
		g.gameOverShakeX = 0
		g.gameOverShakeY = 0
		return
	}

	g.gameOverShakeSeconds -= deltaSeconds
	if g.gameOverShakeSeconds <= 0 {
		g.gameOverShakeX = 0
		g.gameOverShakeY = 0
		g.gameOverShakeSeconds = 0
		return
	}

	envelope := g.gameOverShakeSeconds / defaultGameOverShakeDurationSeconds
	g.gameOverShakePhase += deltaSeconds * defaultGameOverShakeFrequency
	g.gameOverShakeX = math.Sin(g.gameOverShakePhase*2.1) * defaultGameOverShakeAmplitude * envelope
	g.gameOverShakeY = math.Cos(g.gameOverShakePhase*1.7) * (defaultGameOverShakeAmplitude * 0.8) * envelope
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return LogicalWidth, LogicalHeight
}
