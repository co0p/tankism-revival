package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/co0p/tankism-revival/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowScale = 2
	windowTitle = "Tankism Revival"
)

type App struct{}

func (a *App) Update() error {
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0x24, G: 0x2f, B: 0x2d, A: 0xff})
}

func (a *App) Layout(_, _ int) (screenWidth, screenHeight int) {
	return game.LogicalWidth, game.LogicalHeight
}

func main() {
	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowSize(game.LogicalWidth*windowScale, game.LogicalHeight*windowScale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&App{}); err != nil {
		log.Fatal(fmt.Errorf("run game: %w", err))
	}
}
