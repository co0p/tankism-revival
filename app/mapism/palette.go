package main

import (
	"fmt"

	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/hajimehoshi/ebiten/v2"
)

type paletteEntry struct {
	ID        string
	AssetPath string
	Layer     gamepalette.Layer
	GID       int
	Image     *ebiten.Image
}

func buildPaletteEntries(loaded map[string]*ebiten.Image) ([]paletteEntry, error) {
	entries := make([]paletteEntry, 0, len(gamepalette.Sprites))

	for i, spriteDef := range gamepalette.Sprites {
		image, ok := loaded[spriteDef.ID]
		if !ok {
			return nil, fmt.Errorf("palette sprite %q was not loaded", spriteDef.ID)
		}

		entries = append(entries, paletteEntry{
			ID:        spriteDef.ID,
			AssetPath: spriteDef.Path,
			Layer:     spriteDef.Layer,
			GID:       i + 1,
			Image:     image,
		})
	}

	return entries, nil
}
