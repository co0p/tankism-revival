package main

import (
	"fmt"

	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/co0p/tankism-revival/lib/mathutil"
	"github.com/co0p/tankism-revival/lib/tiled"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) paintGroundTileAtCursor() error {
	cx, cy, ok := g.hoveredCell()
	if !ok {
		return nil
	}

	selected := g.palette[g.selectedPalette]
	rotatedGID := tiled.EncodeRotatedGID(selected.GID, g.selectedRotation)
	return g.doc.SetTileGIDOnLayer(tiled.GroundLayerName, cx, cy, rotatedGID)
}

func (g *Game) eraseGroundTileAtCursor() error {
	cx, cy, ok := g.hoveredCell()
	if !ok {
		return nil
	}
	return g.doc.SetTileGIDOnLayer(tiled.GroundLayerName, cx, cy, 0)
}

func (g *Game) placePropObjectAtCursor() (bool, error) {
	worldX, worldY, ok := g.cursorWorldPosition()
	if !ok {
		return false, nil
	}

	selected := g.palette[g.selectedPalette]
	bounds := selected.Image.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		return false, fmt.Errorf("selected sprite has invalid dimensions")
	}

	layerName := propsLayerNameForPaletteLayer(selected.Layer)
	if err := g.doc.AddPropsObjectOnLayer(
		layerName,
		selected.GID,
		worldX,
		worldY,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
		float64(g.selectedRotation),
	); err != nil {
		return false, err
	}

	return true, nil
}

func (g *Game) erasePropObjectAtCursor() (bool, error) {
	worldX, worldY, ok := g.cursorWorldPosition()
	if !ok {
		return false, nil
	}

	maxDistance := float64(g.doc.TileWidth) * 0.75
	selected := g.palette[g.selectedPalette]
	layerName := propsLayerNameForPaletteLayer(selected.Layer)

	removed, err := g.doc.RemoveNearestPropsObjectOnLayer(layerName, worldX, worldY, maxDistance)
	if err != nil {
		return false, err
	}
	if removed {
		return true, nil
	}

	if layerName == tiled.PropsLayerName {
		return g.doc.RemoveNearestPropsObjectOnLayer(tiled.PropsTallLayerName, worldX, worldY, maxDistance)
	}

	return g.doc.RemoveNearestPropsObjectOnLayer(tiled.PropsLayerName, worldX, worldY, maxDistance)
}

func propsLayerNameForPaletteLayer(layer gamepalette.Layer) string {
	if layer == gamepalette.LayerPropsTall {
		return tiled.PropsTallLayerName
	}

	return tiled.PropsLayerName
}

func (g *Game) setSpawnAtCursor(name string) error {
	cx, cy, ok := g.hoveredCell()
	if !ok {
		return fmt.Errorf("cursor is outside the map")
	}

	worldX := float64(cx*g.doc.TileWidth) + float64(g.doc.TileWidth)/2
	worldY := float64(cy*g.doc.TileHeight) + float64(g.doc.TileHeight)/2
	return g.doc.SetSpawn(name, worldX, worldY)
}

func (g *Game) hoveredCell() (cellX, cellY int, ok bool) {
	cellPx, originX, originY := g.viewMetrics()
	if cellPx <= 0 {
		return 0, 0, false
	}

	mx, my := ebiten.CursorPosition()
	localX := mx - originX
	localY := my - originY

	if localX < 0 || localY < 0 {
		return 0, 0, false
	}

	cellX = localX / cellPx
	cellY = localY / cellPx
	if cellX < 0 || cellY < 0 || cellX >= g.doc.Width || cellY >= g.doc.Height {
		return 0, 0, false
	}

	return cellX, cellY, true
}

func (g *Game) cursorWorldPosition() (worldX, worldY float64, ok bool) {
	cellPx, originX, originY := g.viewMetrics()
	if cellPx <= 0 || g.doc.TileWidth <= 0 || g.doc.TileHeight <= 0 {
		return 0, 0, false
	}

	mx, my := ebiten.CursorPosition()
	localX := mx - originX
	localY := my - originY
	if localX < 0 || localY < 0 {
		return 0, 0, false
	}

	maxX := g.doc.Width * cellPx
	maxY := g.doc.Height * cellPx
	if localX >= maxX || localY >= maxY {
		return 0, 0, false
	}

	worldX = float64(localX) * float64(g.doc.TileWidth) / float64(cellPx)
	worldY = float64(localY) * float64(g.doc.TileHeight) / float64(cellPx)
	return worldX, worldY, true
}

func (g *Game) viewMetrics() (cellPx, originX, originY int) {
	if g.doc.Width <= 0 || g.doc.Height <= 0 {
		return 1, 0, 0
	}

	panelX, _, _, _ := g.paletteRect()
	availableWidth := panelX - uiMargin*2
	availableHeight := LogicalHeight - uiMargin*2

	cellByWidth := availableWidth / g.doc.Width
	cellByHeight := availableHeight / g.doc.Height
	cellPx = mathutil.MinInt(cellByWidth, cellByHeight)
	if cellPx < 2 {
		cellPx = 2
	}

	drawWidth := cellPx * g.doc.Width
	drawHeight := cellPx * g.doc.Height
	originX = uiMargin + (availableWidth-drawWidth)/2
	originY = uiMargin + (availableHeight-drawHeight)/2

	return cellPx, originX, originY
}
