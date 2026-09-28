package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/co0p/tankism-revival/game"
	"github.com/co0p/tankism-revival/game/ecs"
	gameinput "github.com/co0p/tankism-revival/game/input"
	gamerender "github.com/co0p/tankism-revival/game/render"
	"github.com/co0p/tankism-revival/game/sim"
	"github.com/co0p/tankism-revival/game/worldseed"
	gameassets "github.com/co0p/tankism-revival/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	windowScale = 2
	windowTitle = "Tankism Revival"

	fixedDeltaSeconds = 1.0 / 60.0

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

	defaultPlayerWeaponCooldown            = 0.28
	defaultPlayerProjectileSpeed           = 360.0
	defaultPlayerProjectileLife            = 1.25
	defaultPlayerProjectileDamage          = 12.0
	defaultPlayerMuzzleOffset              = 22.0
	defaultPlayerProjectileRadius          = 4.0
	defaultPlayerProjectileLightRadius     = 56.0
	defaultPlayerProjectileLightIntensity  = 0.34
	defaultPlayerMuzzleFlashLifetime       = 0.07
	defaultPlayerMuzzleFlashOffset         = 16.0
	defaultPlayerMuzzleFlashLightRadius    = 64.0
	defaultPlayerMuzzleFlashLightIntensity = 0.52

	defaultEnemyWeaponCooldown                  = 0.74
	defaultEnemyProjectileSpeed                 = 300.0
	defaultEnemyProjectileLife                  = 1.15
	defaultEnemyProjectileDamage                = 6.0
	defaultEnemyMuzzleOffset                    = 22.0
	defaultEnemyProjectileRadius                = 4.0
	defaultEnemyProjectileLightRadius           = 48.0
	defaultEnemyProjectileLightIntensity        = 0.24
	defaultEnemyMuzzleFlashLifetime             = 0.07
	defaultEnemyMuzzleFlashOffset               = 16.0
	defaultEnemyMuzzleFlashLightRadius          = 56.0
	defaultEnemyMuzzleFlashLightIntensity       = 0.44
	defaultGrassTileSeed                  int64 = 42
	defaultObstacleSeed                   int64 = 31415

	defaultAmbientDarkness             = 0.70
	defaultPlayerLightRadius           = 170.0
	defaultPlayerLightIntensity        = 0.20
	defaultPlayerBackLightRadius       = 82.0
	defaultPlayerBackLightIntensity    = 0.10
	defaultPlayerBackLightDistance     = 12.0
	defaultObstacleBackLightRadius     = 58.0
	defaultObstacleBackLightIntensity  = 0.02
	defaultObstacleBackLightIntensityJ = 0.02
	defaultObstacleBackLightYOffset    = 8.0
	defaultObstacleLitChance           = 0.55

	obstacleCountPerType = 14
	obstacleEdgeMargin   = 64.0
	obstacleClearRadius  = 180.0

	playerCollisionRadius = 16.0
	enemyCollisionRadius  = 16.0
	barrelCollisionRadius = 11.0
	fenceCollisionRadius  = 14.0
	crateCollisionRadius  = 13.0

	playerMaxHealth = 100.0
	enemyMaxHealth  = 70.0
)

var (
	outsideMapColor = color.RGBA{R: 0x1c, G: 0x27, B: 0x1b, A: 0xff}
)

type App struct {
	world                 *ecs.World
	worldMap              game.Map
	camera                game.Camera
	playerEntity          ecs.Entity
	enemyEntity           ecs.Entity
	playerBackLightEntity ecs.Entity
	renderer              *gamerender.WorldRenderer
}

func NewApp() (*App, error) {
	playerSprite, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.PlayerTankRedSmallPath)
	if err != nil {
		return nil, fmt.Errorf("load player sprite %q: %w", gameassets.PlayerTankRedSmallPath, err)
	}
	enemySprite, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.EnemyTankBlueSmallPath)
	if err != nil {
		return nil, fmt.Errorf("load enemy sprite %q: %w", gameassets.EnemyTankBlueSmallPath, err)
	}
	grassTile1, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.TileGrass1Path)
	if err != nil {
		return nil, fmt.Errorf("load grass tile %q: %w", gameassets.TileGrass1Path, err)
	}
	grassTile2, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.TileGrass2Path)
	if err != nil {
		return nil, fmt.Errorf("load grass tile %q: %w", gameassets.TileGrass2Path, err)
	}
	obstacleBarrel, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.ObstacleBarrelRedTopPath)
	if err != nil {
		return nil, fmt.Errorf("load obstacle sprite %q: %w", gameassets.ObstacleBarrelRedTopPath, err)
	}
	obstacleFence, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.ObstacleFenceRedPath)
	if err != nil {
		return nil, fmt.Errorf("load obstacle sprite %q: %w", gameassets.ObstacleFenceRedPath, err)
	}
	obstacleCrate, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.ObstacleCrateWoodPath)
	if err != nil {
		return nil, fmt.Errorf("load obstacle sprite %q: %w", gameassets.ObstacleCrateWoodPath, err)
	}
	projectileRed, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.ProjectileRedPath)
	if err != nil {
		return nil, fmt.Errorf("load projectile sprite %q: %w", gameassets.ProjectileRedPath, err)
	}
	projectileBlue, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.ProjectileBluePath)
	if err != nil {
		return nil, fmt.Errorf("load projectile sprite %q: %w", gameassets.ProjectileBluePath, err)
	}
	shotRed, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.ShotRedPath)
	if err != nil {
		return nil, fmt.Errorf("load shot sprite %q: %w", gameassets.ShotRedPath, err)
	}
	shotOrange, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.ShotOrangePath)
	if err != nil {
		return nil, fmt.Errorf("load shot sprite %q: %w", gameassets.ShotOrangePath, err)
	}

	sprites := map[string]*ebiten.Image{
		game.SpriteTankRedSmall:         playerSprite,
		game.SpriteTankBlueSmall:        enemySprite,
		game.SpriteTileGrass1:           grassTile1,
		game.SpriteTileGrass2:           grassTile2,
		game.SpriteObstacleBarrelRedTop: obstacleBarrel,
		game.SpriteObstacleFenceRed:     obstacleFence,
		game.SpriteObstacleCrateWood:    obstacleCrate,
		game.SpriteProjectileRed:        projectileRed,
		game.SpriteProjectileBlue:       projectileBlue,
		game.SpriteShotRed:              shotRed,
		game.SpriteShotOrange:           shotOrange,
	}

	renderer := gamerender.NewWorldRenderer(sprites, gamerender.Config{
		OutsideMapColor:     outsideMapColor,
		LightSpriteSize:     256,
		LightIntensityScale: 0.65,
		HealthBarWidth:      34,
		HealthBarHeight:     4,
		HealthBarYOffset:    30,
		HealthBarBorderSize: 1,
	})

	world := ecs.NewWorld()
	worldMap := game.DefaultMap()

	tileBounds := grassTile1.Bounds()
	tileWidth := tileBounds.Dx()
	tileHeight := tileBounds.Dy()
	if tileWidth <= 0 || tileHeight <= 0 {
		return nil, fmt.Errorf("invalid grass tile size: %dx%d", tileWidth, tileHeight)
	}

	if err := worldseed.SeedGrassMap(world, worldMap, worldseed.GrassConfig{
		TileWidth:       tileWidth,
		TileHeight:      tileHeight,
		Seed:            defaultGrassTileSeed,
		PrimarySpriteID: game.SpriteTileGrass1,
		AltSpriteID:     game.SpriteTileGrass2,
		ZIndex:          game.ZIndexMap,
	}); err != nil {
		return nil, fmt.Errorf("seed grass map: %w", err)
	}

	player := world.CreateEntity()
	spawnX := worldMap.Width / 2
	spawnY := worldMap.Height / 2

	if err := world.SetTransform(player, ecs.Transform{
		X:               spawnX,
		Y:               spawnY,
		RotationRadians: 0,
	}); err != nil {
		return nil, fmt.Errorf("set player transform: %w", err)
	}

	if err := world.SetSpriteRenderer(player, ecs.SpriteRenderer{
		SpriteID:              game.SpriteTankRedSmall,
		ZIndex:                game.ZIndexActors,
		RotationOffsetRadians: math.Pi,
	}); err != nil {
		return nil, fmt.Errorf("set player sprite renderer: %w", err)
	}

	if err := world.SetTankDrive(player, ecs.TankDrive{
		MoveSpeed:          defaultPlayerMoveSpeed,
		TurnSpeedRadiansPS: defaultPlayerTurnSpeedRadiansPS,
	}); err != nil {
		return nil, fmt.Errorf("set player tank drive: %w", err)
	}

	if err := world.SetTankAcceleration(player, ecs.TankAcceleration{
		ForwardUnitsPerSecondSquared: defaultPlayerAcceleration,
	}); err != nil {
		return nil, fmt.Errorf("set player tank acceleration: %w", err)
	}

	if err := world.SetTankDeceleration(player, ecs.TankDeceleration{
		CoastUnitsPerSecondSquared: defaultPlayerCoastDeceleration,
		BrakeUnitsPerSecondSquared: defaultPlayerBrakeDeceleration,
	}); err != nil {
		return nil, fmt.Errorf("set player tank deceleration: %w", err)
	}

	if err := world.SetTankVelocity(player, ecs.TankVelocity{}); err != nil {
		return nil, fmt.Errorf("set player tank velocity: %w", err)
	}

	if err := world.SetTankControlIntent(player, ecs.TankControlIntent{}); err != nil {
		return nil, fmt.Errorf("set player tank control intent: %w", err)
	}

	if err := world.SetCircleCollider(player, ecs.CircleCollider{Radius: playerCollisionRadius}); err != nil {
		return nil, fmt.Errorf("set player collider: %w", err)
	}
	if err := world.SetHealth(player, ecs.Health{Current: playerMaxHealth, Max: playerMaxHealth}); err != nil {
		return nil, fmt.Errorf("set player health: %w", err)
	}

	if err := world.SetPointLight(player, ecs.PointLight{
		Radius:    defaultPlayerLightRadius,
		Intensity: defaultPlayerLightIntensity,
		Color:     ecs.LightColor{R: 1, G: 0.93, B: 0.82},
	}); err != nil {
		return nil, fmt.Errorf("set player point light: %w", err)
	}
	if err := world.SetWeapon(player, ecs.Weapon{
		CooldownSeconds:           defaultPlayerWeaponCooldown,
		ProjectileSpeed:           defaultPlayerProjectileSpeed,
		ProjectileLifetime:        defaultPlayerProjectileLife,
		ProjectileDamage:          defaultPlayerProjectileDamage,
		MuzzleOffset:              defaultPlayerMuzzleOffset,
		ProjectileSpriteID:        game.SpriteProjectileRed,
		ProjectileRadius:          defaultPlayerProjectileRadius,
		ProjectileLightRadius:     defaultPlayerProjectileLightRadius,
		ProjectileLightIntensity:  defaultPlayerProjectileLightIntensity,
		ProjectileLightColor:      ecs.LightColor{R: 1.00, G: 0.62, B: 0.38},
		MuzzleFlashSpriteID:       game.SpriteShotRed,
		MuzzleFlashLifetime:       defaultPlayerMuzzleFlashLifetime,
		MuzzleFlashOffset:         defaultPlayerMuzzleFlashOffset,
		MuzzleFlashLightRadius:    defaultPlayerMuzzleFlashLightRadius,
		MuzzleFlashLightIntensity: defaultPlayerMuzzleFlashLightIntensity,
		MuzzleFlashLightColor:     ecs.LightColor{R: 1.00, G: 0.68, B: 0.45},
	}); err != nil {
		return nil, fmt.Errorf("set player weapon: %w", err)
	}

	playerBackLight := world.CreateEntity()
	backLightX, backLightY := tankBackLightPosition(spawnX, spawnY, 0)
	if err := world.SetTransform(playerBackLight, ecs.Transform{
		X:               backLightX,
		Y:               backLightY,
		RotationRadians: 0,
	}); err != nil {
		return nil, fmt.Errorf("set player back light transform: %w", err)
	}
	if err := world.SetPointLight(playerBackLight, ecs.PointLight{
		Radius:    defaultPlayerBackLightRadius,
		Intensity: defaultPlayerBackLightIntensity,
		Color:     ecs.LightColor{R: 0.75, G: 0.80, B: 0.95},
	}); err != nil {
		return nil, fmt.Errorf("set player back light: %w", err)
	}

	ambient := world.CreateEntity()
	if err := world.SetAmbientLight(ambient, ecs.AmbientLight{Darkness: defaultAmbientDarkness}); err != nil {
		return nil, fmt.Errorf("set ambient light: %w", err)
	}

	enemy := world.CreateEntity()
	enemySpawnX, enemySpawnY := worldMap.ClampPosition(spawnX+280, spawnY-220)
	if err := world.SetTransform(enemy, ecs.Transform{
		X:               enemySpawnX,
		Y:               enemySpawnY,
		RotationRadians: math.Pi,
	}); err != nil {
		return nil, fmt.Errorf("set enemy transform: %w", err)
	}
	if err := world.SetSpriteRenderer(enemy, ecs.SpriteRenderer{
		SpriteID:              game.SpriteTankBlueSmall,
		ZIndex:                game.ZIndexActors,
		RotationOffsetRadians: math.Pi,
	}); err != nil {
		return nil, fmt.Errorf("set enemy sprite renderer: %w", err)
	}
	if err := world.SetTankDrive(enemy, ecs.TankDrive{
		MoveSpeed:          defaultEnemyMoveSpeed,
		TurnSpeedRadiansPS: defaultEnemyTurnSpeedRadiansPS,
	}); err != nil {
		return nil, fmt.Errorf("set enemy tank drive: %w", err)
	}
	if err := world.SetTankAcceleration(enemy, ecs.TankAcceleration{
		ForwardUnitsPerSecondSquared: defaultEnemyAcceleration,
	}); err != nil {
		return nil, fmt.Errorf("set enemy tank acceleration: %w", err)
	}
	if err := world.SetTankDeceleration(enemy, ecs.TankDeceleration{
		CoastUnitsPerSecondSquared: defaultEnemyCoastDeceleration,
		BrakeUnitsPerSecondSquared: defaultEnemyBrakeDeceleration,
	}); err != nil {
		return nil, fmt.Errorf("set enemy tank deceleration: %w", err)
	}
	if err := world.SetTankVelocity(enemy, ecs.TankVelocity{}); err != nil {
		return nil, fmt.Errorf("set enemy tank velocity: %w", err)
	}
	if err := world.SetTankControlIntent(enemy, ecs.TankControlIntent{}); err != nil {
		return nil, fmt.Errorf("set enemy tank control intent: %w", err)
	}
	if err := world.SetCircleCollider(enemy, ecs.CircleCollider{Radius: enemyCollisionRadius}); err != nil {
		return nil, fmt.Errorf("set enemy collider: %w", err)
	}
	if err := world.SetHealth(enemy, ecs.Health{Current: enemyMaxHealth, Max: enemyMaxHealth}); err != nil {
		return nil, fmt.Errorf("set enemy health: %w", err)
	}
	if err := world.SetPointLight(enemy, ecs.PointLight{
		Radius:    120,
		Intensity: 0.10,
		Color:     ecs.LightColor{R: 0.7, G: 0.82, B: 1.0},
	}); err != nil {
		return nil, fmt.Errorf("set enemy point light: %w", err)
	}
	if err := world.SetWeapon(enemy, ecs.Weapon{
		CooldownSeconds:           defaultEnemyWeaponCooldown,
		ProjectileSpeed:           defaultEnemyProjectileSpeed,
		ProjectileLifetime:        defaultEnemyProjectileLife,
		ProjectileDamage:          defaultEnemyProjectileDamage,
		MuzzleOffset:              defaultEnemyMuzzleOffset,
		ProjectileSpriteID:        game.SpriteProjectileBlue,
		ProjectileRadius:          defaultEnemyProjectileRadius,
		ProjectileLightRadius:     defaultEnemyProjectileLightRadius,
		ProjectileLightIntensity:  defaultEnemyProjectileLightIntensity,
		ProjectileLightColor:      ecs.LightColor{R: 0.74, G: 0.86, B: 1.00},
		MuzzleFlashSpriteID:       game.SpriteShotOrange,
		MuzzleFlashLifetime:       defaultEnemyMuzzleFlashLifetime,
		MuzzleFlashOffset:         defaultEnemyMuzzleFlashOffset,
		MuzzleFlashLightRadius:    defaultEnemyMuzzleFlashLightRadius,
		MuzzleFlashLightIntensity: defaultEnemyMuzzleFlashLightIntensity,
		MuzzleFlashLightColor:     ecs.LightColor{R: 0.96, G: 0.80, B: 0.56},
	}); err != nil {
		return nil, fmt.Errorf("set enemy weapon: %w", err)
	}

	if err := worldseed.SeedObstacles(world, worldMap, worldseed.ObstacleConfig{
		Seed:              defaultObstacleSeed,
		CountPerType:      obstacleCountPerType,
		PlacementAttempts: 200,
		EdgeMargin:        obstacleEdgeMargin,
		ClearRadius:       obstacleClearRadius,
		SpawnX:            spawnX,
		SpawnY:            spawnY,
		EnemySpawnX:       enemySpawnX,
		EnemySpawnY:       enemySpawnY,
		Sprites: []string{
			game.SpriteObstacleBarrelRedTop,
			game.SpriteObstacleFenceRed,
			game.SpriteObstacleCrateWood,
		},
		ColliderRadiusBySprite: map[string]float64{
			game.SpriteObstacleBarrelRedTop: barrelCollisionRadius,
			game.SpriteObstacleFenceRed:     fenceCollisionRadius,
			game.SpriteObstacleCrateWood:    crateCollisionRadius,
		},
		ZIndex:                   game.ZIndexObstacles,
		BackLightChance:          defaultObstacleLitChance,
		BackLightRadius:          defaultObstacleBackLightRadius,
		BackLightIntensity:       defaultObstacleBackLightIntensity,
		BackLightIntensityJitter: defaultObstacleBackLightIntensityJ,
		BackLightYOffset:         defaultObstacleBackLightYOffset,
		BackLightColorBase:       ecs.LightColor{R: 0.80, G: 0.78, B: 0.72},
		BackLightColorJitter:     ecs.LightColor{R: 0.08, G: 0.08, B: 0.06},
	}); err != nil {
		return nil, fmt.Errorf("seed obstacles: %w", err)
	}

	camera := game.NewCamera(float64(game.LogicalWidth), float64(game.LogicalHeight))
	camera.FollowTarget(worldMap, spawnX, spawnY)

	return &App{
		world:                 world,
		worldMap:              worldMap,
		camera:                camera,
		playerEntity:          player,
		enemyEntity:           enemy,
		playerBackLightEntity: playerBackLight,
		renderer:              renderer,
	}, nil
}

func (a *App) Update() error {
	intent := gameinput.ReadPlayerTankIntent()
	if err := a.world.SetTankControlIntent(a.playerEntity, intent); err != nil {
		return fmt.Errorf("set player intent: %w", err)
	}

	if err := sim.UpdateChaserIntent(a.world, a.enemyEntity, a.playerEntity, defaultEnemyStopDistance); err != nil {
		return fmt.Errorf("update enemy intent: %w", err)
	}

	sim.IntegrateTankMovement(a.world, fixedDeltaSeconds)

	if err := sim.UpdateWeaponsAndSpawnProjectiles(a.world, fixedDeltaSeconds, game.ZIndexProjectiles, game.ZIndexEffects); err != nil {
		return fmt.Errorf("update weapons and spawn projectiles: %w", err)
	}

	if err := sim.IntegrateProjectiles(a.world, fixedDeltaSeconds); err != nil {
		return fmt.Errorf("integrate projectiles: %w", err)
	}

	if err := sim.ResolveProjectileImpacts(a.world); err != nil {
		return fmt.Errorf("resolve projectile impacts: %w", err)
	}

	if err := sim.ResolveTankObstacleCollisions(a.world); err != nil {
		return fmt.Errorf("resolve tank obstacle collisions: %w", err)
	}

	if err := sim.UpdateLifetimes(a.world, fixedDeltaSeconds); err != nil {
		return fmt.Errorf("update lifetimes: %w", err)
	}

	playerTransform, ok := a.world.Transform(a.playerEntity)
	if !ok {
		return fmt.Errorf("player transform missing")
	}

	clampedX, clampedY := a.worldMap.ClampPosition(playerTransform.X, playerTransform.Y)
	if clampedX != playerTransform.X || clampedY != playerTransform.Y {
		playerTransform.X = clampedX
		playerTransform.Y = clampedY

		if err := a.world.SetTransform(a.playerEntity, playerTransform); err != nil {
			return fmt.Errorf("clamp player transform: %w", err)
		}
	}

	enemyTransform, ok := a.world.Transform(a.enemyEntity)
	if ok {
		enemyClampedX, enemyClampedY := a.worldMap.ClampPosition(enemyTransform.X, enemyTransform.Y)
		if enemyClampedX != enemyTransform.X || enemyClampedY != enemyTransform.Y {
			enemyTransform.X = enemyClampedX
			enemyTransform.Y = enemyClampedY
			if err := a.world.SetTransform(a.enemyEntity, enemyTransform); err != nil {
				return fmt.Errorf("clamp enemy transform: %w", err)
			}
		}
	}

	a.camera.FollowTarget(a.worldMap, playerTransform.X, playerTransform.Y)

	backLightX, backLightY := tankBackLightPosition(
		playerTransform.X,
		playerTransform.Y,
		playerTransform.RotationRadians,
	)
	if err := a.world.SetTransform(a.playerBackLightEntity, ecs.Transform{
		X:               backLightX,
		Y:               backLightY,
		RotationRadians: 0,
	}); err != nil {
		return fmt.Errorf("update player back light transform: %w", err)
	}

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	a.renderer.Draw(screen, a.world, a.camera)
}

func (a *App) Layout(_, _ int) (screenWidth, screenHeight int) {
	return game.LogicalWidth, game.LogicalHeight
}

func tankBackLightPosition(x, y, rotationRadians float64) (float64, float64) {
	forwardX := math.Sin(rotationRadians)
	forwardY := -math.Cos(rotationRadians)

	return x - forwardX*defaultPlayerBackLightDistance, y - forwardY*defaultPlayerBackLightDistance
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowSize(game.LogicalWidth*windowScale, game.LogicalHeight*windowScale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(fmt.Errorf("run game: %w", err))
	}
}
