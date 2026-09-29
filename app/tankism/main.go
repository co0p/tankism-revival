package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/co0p/tankism-revival/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowScale = 2
	windowTitle = "Tankism Revival"
)

func main() {
	mapPath := flag.String("map", game.DefaultMapPath, "path to Tiled map JSON (.tmj)")
	flag.Parse()

	gameInstance, err := game.NewGameWithMap(*mapPath)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowSize(game.LogicalWidth*windowScale, game.LogicalHeight*windowScale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(fmt.Errorf("run game: %w", err))
	}
}
