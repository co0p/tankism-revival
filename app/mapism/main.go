package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowTitle = "Mapism — Tankism Map Editor"
)

func main() {
	mapPath := flag.String("map", DefaultMapPath, "path to Tiled map JSON (.tmj)")
	mapWidth := flag.Int("width", DefaultMapWidth, "map width in tiles when creating a new map")
	mapHeight := flag.Int("height", DefaultMapHeight, "map height in tiles when creating a new map")
	flag.Parse()

	editorGame, err := NewGame(*mapPath, *mapWidth, *mapHeight)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowSize(LogicalWidth, LogicalHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(editorGame); err != nil {
		log.Fatal(fmt.Errorf("run map editor: %w", err))
	}
}
