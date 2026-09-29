package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/co0p/tankism-revival/lib/sprites"
	"github.com/co0p/tankism-revival/lib/tiled"
	gameassets "github.com/co0p/tankism-revival/resources"
)

func (g *Game) syncTileset() error {
	if len(g.palette) == 0 {
		return fmt.Errorf("palette is empty")
	}

	mapDir := filepath.Dir(g.mapPath)
	tiles := make([]tiled.TilesetTile, 0, len(g.palette))

	for _, entry := range g.palette {
		bounds := entry.Image.Bounds()
		imagePath := filepath.Join("resources", filepath.FromSlash(entry.AssetPath))
		relPath, err := filepath.Rel(mapDir, imagePath)
		if err != nil {
			relPath = imagePath
		}

		tiles = append(tiles, tiled.TilesetTile{
			ID:          entry.GID - 1,
			Image:       filepath.ToSlash(relPath),
			ImageWidth:  bounds.Dx(),
			ImageHeight: bounds.Dy(),
		})
	}

	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].ID < tiles[j].ID
	})

	g.doc.Tilesets = []tiled.Tileset{
		{
			FirstGID:   1,
			Name:       "mapism_palette",
			TileWidth:  g.doc.TileWidth,
			TileHeight: g.doc.TileHeight,
			TileCount:  gamepalette.MaxGID(),
			Columns:    0,
			Tiles:      tiles,
		},
	}

	return nil
}

func loadPaletteEntries() ([]paletteEntry, error) {
	assets := make([]sprites.Asset, 0, len(gamepalette.Sprites))
	for _, spriteDef := range gamepalette.Sprites {
		assets = append(assets, sprites.Asset{ID: spriteDef.ID, Path: spriteDef.Path})
	}

	loaded, err := sprites.Load(gameassets.FS, assets)
	if err != nil {
		return nil, err
	}

	return buildPaletteEntries(loaded)
}

func loadOrCreate(mapPath string, mapWidth, mapHeight int) (tiled.Map, string, error) {
	if _, err := os.Stat(mapPath); err == nil {
		doc, readErr := tiled.ReadFile(mapPath)
		if readErr != nil {
			return tiled.Map{}, "", readErr
		}

		return doc, fmt.Sprintf("loaded map: %s", mapPath), nil
	} else if !os.IsNotExist(err) {
		return tiled.Map{}, "", err
	}

	doc := tiled.NewOrthogonalMap(mapWidth, mapHeight, DefaultTileWidth, DefaultTileHeight)
	return doc, fmt.Sprintf("new map (unsaved): %s", mapPath), nil
}
