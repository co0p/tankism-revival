package main

import (
	"fmt"
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
	defaultGrassTileSeed            int64 = 42
	defaultObstacleSeed             int64 = 31415

	obstacleCountPerType = 14
	obstacleEdgeMargin   = 64.0
	obstacleClearRadius  = 180.0

	playerCollisionRadius = 16.0
	barrelCollisionRadius = 11.0
	fenceCollisionRadius  = 14.0
	crateCollisionRadius  = 13.0
)

var (
	outsideMapColor = color.RGBA{R: 0x1c, G: 0x27, B: 0x1b, A: 0xff}
)

type App struct {
	world        *ecs.World
	worldMap     game.Map
	camera       game.Camera
	playerEntity ecs.Entity
	sprites      map[string]*ebiten.Image
}

func NewApp() (*App, error) {
	playerSprite, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.PlayerTankRedSmallPath)
	if err != nil {
		return nil, fmt.Errorf("load player sprite %q: %w", gameassets.PlayerTankRedSmallPath, err)
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

	if err := seedObstacles(world, worldMap, spawnX, spawnY, defaultObstacleSeed); err != nil {
		return nil, fmt.Errorf("seed obstacles: %w", err)
	}

	camera := game.NewCamera(float64(game.LogicalWidth), float64(game.LogicalHeight))
	camera.FollowTarget(worldMap, spawnX, spawnY)

	return &App{
		world:        world,
		worldMap:     worldMap,
		camera:       camera,
		playerEntity: player,
		sprites: map[string]*ebiten.Image{
			game.SpriteTankRedSmall:         playerSprite,
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

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(outsideMapColor)

	for _, renderable := range a.world.Renderables() {
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

func seedObstacles(world *ecs.World, worldMap game.Map, spawnX, spawnY float64, seed int64) error {
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
