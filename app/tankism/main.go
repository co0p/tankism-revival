package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

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

	defaultPlayerMoveSpeed          = 160
	defaultPlayerTurnSpeedRadiansPS = math.Pi
)

type App struct {
	world        *ecs.World
	playerEntity ecs.Entity
	sprites      map[string]*ebiten.Image
}

func NewApp() (*App, error) {
	playerSprite, _, err := ebitenutil.NewImageFromFileSystem(gameassets.FS, gameassets.PlayerTankRedSmallPath)
	if err != nil {
		return nil, fmt.Errorf("load player sprite %q: %w", gameassets.PlayerTankRedSmallPath, err)
	}

	world := ecs.NewWorld()

	player := world.CreateEntity()
	if err := world.SetTransform(player, ecs.Transform{
		X:               float64(game.LogicalWidth) / 2,
		Y:               float64(game.LogicalHeight) / 2,
		RotationRadians: 0,
	}); err != nil {
		return nil, fmt.Errorf("set player transform: %w", err)
	}

	if err := world.SetSpriteRenderer(player, ecs.SpriteRenderer{
		SpriteID: game.SpriteTankRedSmall,
		Layer:    0,
	}); err != nil {
		return nil, fmt.Errorf("set player sprite renderer: %w", err)
	}

	if err := world.SetTankDrive(player, ecs.TankDrive{
		MoveSpeed:          defaultPlayerMoveSpeed,
		TurnSpeedRadiansPS: defaultPlayerTurnSpeedRadiansPS,
	}); err != nil {
		return nil, fmt.Errorf("set player tank drive: %w", err)
	}

	if err := world.SetTankControlIntent(player, ecs.TankControlIntent{}); err != nil {
		return nil, fmt.Errorf("set player tank control intent: %w", err)
	}

	return &App{
		world:        world,
		playerEntity: player,
		sprites: map[string]*ebiten.Image{
			game.SpriteTankRedSmall: playerSprite,
		},
	}, nil
}

func (a *App) Update() error {
	intent := readPlayerIntent()
	if err := a.world.SetTankControlIntent(a.playerEntity, intent); err != nil {
		return fmt.Errorf("set player intent: %w", err)
	}

	sim.IntegrateTankMovement(a.world, fixedDeltaSeconds)

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0x24, G: 0x2f, B: 0x2d, A: 0xff})

	for _, renderable := range a.world.Renderables() {
		spriteImage, ok := a.sprites[renderable.Sprite.SpriteID]
		if !ok {
			continue
		}

		bounds := spriteImage.Bounds()
		halfWidth := float64(bounds.Dx()) / 2
		halfHeight := float64(bounds.Dy()) / 2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-halfWidth, -halfHeight)
		op.GeoM.Rotate(renderable.Transform.RotationRadians)
		op.GeoM.Translate(renderable.Transform.X, renderable.Transform.Y)

		screen.DrawImage(spriteImage, op)
	}
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
