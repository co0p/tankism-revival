package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/co0p/tankism-revival/game"
	"github.com/co0p/tankism-revival/game/ecs"
	"github.com/co0p/tankism-revival/game/sim"
	gameassets "github.com/co0p/tankism-revival/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	windowScale = 2
	windowTitle = "Tankism Revival"

	fixedDeltaSeconds = 1.0 / 60.0

	defaultPlayerMoveSpeed                = 160
	defaultPlayerTurnSpeedRadiansPS       = math.Pi
	defaultPlayerAcceleration             = 240
	defaultPlayerCoastDeceleration        = 180
	defaultPlayerBrakeDeceleration        = 420
	defaultEnemyMoveSpeed                 = 145
	defaultEnemyTurnSpeedRadiansPS        = 2.5
	defaultEnemyAcceleration              = 210
	defaultEnemyCoastDeceleration         = 170
	defaultEnemyBrakeDeceleration         = 320
	defaultEnemyStopDistance              = 110.0
	defaultGrassTileSeed            int64 = 42
	defaultObstacleSeed             int64 = 31415

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
	defaultLightSpriteSize             = 256

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

	healthBarWidth      = 34.0
	healthBarHeight     = 4.0
	healthBarYOffset    = 30.0
	healthBarBorderSize = 1.0
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
	sprites               map[string]*ebiten.Image
	lightSprite           *ebiten.Image
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
	lightSprite := newRadialLightSprite(defaultLightSpriteSize)

	world := ecs.NewWorld()
	worldMap := game.DefaultMap()

	tileBounds := grassTile1.Bounds()
	tileWidth := tileBounds.Dx()
	tileHeight := tileBounds.Dy()
	if tileWidth <= 0 || tileHeight <= 0 {
		return nil, fmt.Errorf("invalid grass tile size: %dx%d", tileWidth, tileHeight)
	}

	if err := seedGrassMap(world, worldMap, tileWidth, tileHeight, defaultGrassTileSeed); err != nil {
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

	if err := seedObstacles(world, worldMap, spawnX, spawnY, enemySpawnX, enemySpawnY, defaultObstacleSeed); err != nil {
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
		lightSprite:           lightSprite,
		sprites: map[string]*ebiten.Image{
			game.SpriteTankRedSmall:         playerSprite,
			game.SpriteTankBlueSmall:        enemySprite,
			game.SpriteTileGrass1:           grassTile1,
			game.SpriteTileGrass2:           grassTile2,
			game.SpriteObstacleBarrelRedTop: obstacleBarrel,
			game.SpriteObstacleFenceRed:     obstacleFence,
			game.SpriteObstacleCrateWood:    obstacleCrate,
		},
	}, nil
}

func (a *App) Update() error {
	intent := readPlayerIntent()
	if err := a.world.SetTankControlIntent(a.playerEntity, intent); err != nil {
		return fmt.Errorf("set player intent: %w", err)
	}

	if err := sim.UpdateChaserIntent(a.world, a.enemyEntity, a.playerEntity, defaultEnemyStopDistance); err != nil {
		return fmt.Errorf("update enemy intent: %w", err)
	}

	sim.IntegrateTankMovement(a.world, fixedDeltaSeconds)

	if err := sim.ResolveTankObstacleCollisions(a.world); err != nil {
		return fmt.Errorf("resolve tank obstacle collisions: %w", err)
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
	screen.Fill(outsideMapColor)
	a.drawRenderables(screen, nil)
	a.drawLighting(screen)
	a.drawRenderables(screen, func(r ecs.Renderable) bool {
		return r.Sprite.ZIndex >= game.ZIndexObstacles
	})
	a.drawTankHealthBars(screen)
}

func (a *App) drawRenderables(screen *ebiten.Image, include func(ecs.Renderable) bool) {
	for _, renderable := range a.world.Renderables() {
		if include != nil && !include(renderable) {
			continue
		}

		spriteImage, ok := a.sprites[renderable.Sprite.SpriteID]
		if !ok {
			continue
		}

		bounds := spriteImage.Bounds()
		halfWidth := float64(bounds.Dx()) / 2
		halfHeight := float64(bounds.Dy()) / 2

		op := &ebiten.DrawImageOptions{}
		scaleX := 1.0
		scaleY := 1.0
		if renderable.Sprite.FlipX {
			scaleX = -1
		}
		if renderable.Sprite.FlipY {
			scaleY = -1
		}

		op.GeoM.Translate(-halfWidth, -halfHeight)
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Rotate(renderable.Transform.RotationRadians + renderable.Sprite.RotationOffsetRadians)

		screenX, screenY := a.camera.WorldToScreen(renderable.Transform.X, renderable.Transform.Y)
		op.GeoM.Translate(screenX, screenY)

		screen.DrawImage(spriteImage, op)
	}
}

func (a *App) drawLighting(screen *ebiten.Image) {
	ambientDarkness := a.combinedAmbientDarkness()
	if ambientDarkness > 0 {
		overlayAlpha := uint8(clamp01(ambientDarkness) * 255)
		ebitenutil.DrawRect(
			screen,
			0,
			0,
			float64(game.LogicalWidth),
			float64(game.LogicalHeight),
			color.RGBA{R: 0, G: 0, B: 0, A: overlayAlpha},
		)
	}

	lightSize := a.lightSprite.Bounds().Dx()
	if lightSize <= 0 {
		return
	}

	for _, source := range a.world.PointLightSources() {
		if source.Light.Radius <= 0 || source.Light.Intensity <= 0 {
			continue
		}

		screenX, screenY := a.camera.WorldToScreen(source.Transform.X, source.Transform.Y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(lightSize)/2, -float64(lightSize)/2)

		scale := (source.Light.Radius * 2) / float64(lightSize)
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(screenX, screenY)

		intensity := clamp01(source.Light.Intensity) * 0.65
		if intensity == 0 {
			continue
		}
		op.ColorScale.Scale(
			float32(source.Light.Color.R*intensity),
			float32(source.Light.Color.G*intensity),
			float32(source.Light.Color.B*intensity),
			float32(intensity),
		)
		op.Blend = ebiten.BlendLighter

		screen.DrawImage(a.lightSprite, op)
	}
}

func (a *App) drawTankHealthBars(screen *ebiten.Image) {
	for _, state := range a.world.TankHealthStates() {
		if state.Health.Max <= 0 {
			continue
		}

		ratio := clamp01(state.Health.Current / state.Health.Max)

		screenX, screenY := a.camera.WorldToScreen(state.Transform.X, state.Transform.Y)
		barX := screenX - healthBarWidth/2
		barY := screenY + healthBarYOffset

		ebitenutil.DrawRect(
			screen,
			barX-healthBarBorderSize,
			barY-healthBarBorderSize,
			healthBarWidth+2*healthBarBorderSize,
			healthBarHeight+2*healthBarBorderSize,
			color.RGBA{R: 0x12, G: 0x12, B: 0x12, A: 0xd8},
		)
		ebitenutil.DrawRect(
			screen,
			barX,
			barY,
			healthBarWidth,
			healthBarHeight,
			color.RGBA{R: 0x3a, G: 0x3a, B: 0x3a, A: 0xe8},
		)

		fillWidth := healthBarWidth * ratio
		if fillWidth > 0 {
			ebitenutil.DrawRect(
				screen,
				barX,
				barY,
				fillWidth,
				healthBarHeight,
				healthColor(ratio),
			)
		}
	}
}

func healthColor(ratio float64) color.RGBA {
	if ratio > 0.60 {
		return color.RGBA{R: 0x55, G: 0xcf, B: 0x58, A: 0xff}
	}
	if ratio > 0.30 {
		return color.RGBA{R: 0xe0, G: 0xb8, B: 0x3d, A: 0xff}
	}

	return color.RGBA{R: 0xdf, G: 0x54, B: 0x54, A: 0xff}
}

func (a *App) combinedAmbientDarkness() float64 {
	lights := a.world.AmbientLights()
	if len(lights) == 0 {
		return 0
	}

	darkness := 0.0
	for _, light := range lights {
		if light.Darkness > darkness {
			darkness = light.Darkness
		}
	}

	return clamp01(darkness)
}

func newRadialLightSprite(diameter int) *ebiten.Image {
	if diameter < 2 {
		diameter = 2
	}

	img := image.NewNRGBA(image.Rect(0, 0, diameter, diameter))
	radius := float64(diameter) / 2
	cx := radius - 0.5
	cy := radius - 0.5

	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			distance := math.Sqrt(dx*dx+dy*dy) / radius
			if distance >= 1 {
				continue
			}

			falloff := 1 - distance
			alpha := uint8(clamp01(math.Pow(falloff, 3.0)) * 190)
			img.Set(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: alpha})
		}
	}

	return ebiten.NewImageFromImage(img)
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}

	return value
}

func seedGrassMap(world *ecs.World, worldMap game.Map, tileWidth, tileHeight int, seed int64) error {
	rng := rand.New(rand.NewSource(seed))

	columns := int(math.Ceil(worldMap.Width / float64(tileWidth)))
	rows := int(math.Ceil(worldMap.Height / float64(tileHeight)))

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			entity := world.CreateEntity()

			worldX := (float64(col) + 0.5) * float64(tileWidth)
			worldY := (float64(row) + 0.5) * float64(tileHeight)

			if err := world.SetTransform(entity, ecs.Transform{X: worldX, Y: worldY, RotationRadians: 0}); err != nil {
				return err
			}

			spriteID := game.SpriteTileGrass1
			if rng.Intn(2) == 1 {
				spriteID = game.SpriteTileGrass2
			}

			rotationOffset := float64(rng.Intn(4)) * (math.Pi / 2)
			flipX := rng.Intn(2) == 1
			flipY := rng.Intn(2) == 1

			if err := world.SetSpriteRenderer(entity, ecs.SpriteRenderer{
				SpriteID:              spriteID,
				ZIndex:                game.ZIndexMap,
				RotationOffsetRadians: rotationOffset,
				FlipX:                 flipX,
				FlipY:                 flipY,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func seedObstacles(world *ecs.World, worldMap game.Map, spawnX, spawnY, enemySpawnX, enemySpawnY float64, seed int64) error {
	rng := rand.New(rand.NewSource(seed))

	colliderRadiusBySprite := map[string]float64{
		game.SpriteObstacleBarrelRedTop: barrelCollisionRadius,
		game.SpriteObstacleFenceRed:     fenceCollisionRadius,
		game.SpriteObstacleCrateWood:    crateCollisionRadius,
	}

	placeObstacle := func(spriteID string) error {
		radius, ok := colliderRadiusBySprite[spriteID]
		if !ok {
			return fmt.Errorf("missing collider radius for sprite %q", spriteID)
		}

		for range 200 {
			x := obstacleEdgeMargin + rng.Float64()*(worldMap.Width-2*obstacleEdgeMargin)
			y := obstacleEdgeMargin + rng.Float64()*(worldMap.Height-2*obstacleEdgeMargin)

			dx := x - spawnX
			dy := y - spawnY
			if dx*dx+dy*dy < obstacleClearRadius*obstacleClearRadius {
				continue
			}

			edx := x - enemySpawnX
			edy := y - enemySpawnY
			if edx*edx+edy*edy < obstacleClearRadius*obstacleClearRadius {
				continue
			}

			entity := world.CreateEntity()
			if err := world.SetTransform(entity, ecs.Transform{X: x, Y: y, RotationRadians: 0}); err != nil {
				return err
			}

			rotationOffset := float64(rng.Intn(4)) * (math.Pi / 2)
			flipX := rng.Intn(2) == 1
			flipY := rng.Intn(2) == 1

			if err := world.SetSpriteRenderer(entity, ecs.SpriteRenderer{
				SpriteID:              spriteID,
				ZIndex:                game.ZIndexObstacles,
				RotationOffsetRadians: rotationOffset,
				FlipX:                 flipX,
				FlipY:                 flipY,
			}); err != nil {
				return err
			}

			if err := world.SetCircleCollider(entity, ecs.CircleCollider{Radius: radius}); err != nil {
				return err
			}

			if err := world.SetSolidObstacle(entity); err != nil {
				return err
			}

			if rng.Float64() < defaultObstacleLitChance {
				lightIntensity := defaultObstacleBackLightIntensity + rng.Float64()*defaultObstacleBackLightIntensityJ
				lightEntity := world.CreateEntity()
				if err := world.SetTransform(lightEntity, ecs.Transform{
					X:               x,
					Y:               y + defaultObstacleBackLightYOffset,
					RotationRadians: 0,
				}); err != nil {
					return err
				}
				if err := world.SetPointLight(lightEntity, ecs.PointLight{
					Radius:    defaultObstacleBackLightRadius,
					Intensity: lightIntensity,
					Color: ecs.LightColor{
						R: 0.80 + rng.Float64()*0.08,
						G: 0.78 + rng.Float64()*0.08,
						B: 0.72 + rng.Float64()*0.06,
					},
				}); err != nil {
					return err
				}
			}

			return nil
		}

		return fmt.Errorf("could not place obstacle %q after many attempts", spriteID)
	}

	for range obstacleCountPerType {
		if err := placeObstacle(game.SpriteObstacleBarrelRedTop); err != nil {
			return err
		}
		if err := placeObstacle(game.SpriteObstacleFenceRed); err != nil {
			return err
		}
		if err := placeObstacle(game.SpriteObstacleCrateWood); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Layout(_, _ int) (screenWidth, screenHeight int) {
	return game.LogicalWidth, game.LogicalHeight
}

func readPlayerIntent() ecs.TankControlIntent {
	var intent ecs.TankControlIntent

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		intent.Throttle += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		intent.Throttle -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		intent.Turn -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		intent.Turn += 1
	}

	return intent
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
