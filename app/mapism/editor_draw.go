package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/co0p/tankism-revival/lib/mathutil"
	"github.com/co0p/tankism-revival/lib/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) drawSpawns(screen *ebiten.Image, spawnName string, markerColor color.Color, cellPx, originX, originY int) {
	spawns := g.doc.Spawns(spawnName)
	if len(spawns) == 0 || g.doc.TileWidth <= 0 || g.doc.TileHeight <= 0 {
		return
	}

	size := mathutil.MaxInt(5, cellPx/3)
	for _, spawn := range spawns {
		screenX := originX + int((spawn.X/float64(g.doc.TileWidth))*float64(cellPx))
		screenY := originY + int((spawn.Y/float64(g.doc.TileHeight))*float64(cellPx))

		ebitenutil.DrawRect(
			screen,
			float64(screenX-size/2),
			float64(screenY-size/2),
			float64(size),
			float64(size),
			markerColor,
		)
	}
}

func (g *Game) drawTileLayer(screen *ebiten.Image, layerName string, cellPx, originX, originY int) {
	for row := 0; row < g.doc.Height; row++ {
		for col := 0; col < g.doc.Width; col++ {
			gid, ok := g.doc.TileGIDAtLayer(layerName, col, row)
			if !ok || gid == 0 {
				continue
			}

			x := originX + col*cellPx
			y := originY + row*cellPx
			g.drawTile(screen, gid, x, y, cellPx)
		}
	}
}

func (g *Game) drawPropsObjectsOnLayer(screen *ebiten.Image, layerName string, cellPx, originX, originY int) {
	for _, obj := range g.doc.PropsObjectsOnLayer(layerName) {
		entry := g.paletteByGID(obj.GID)
		if entry == nil || entry.Image == nil {
			continue
		}

		imgBounds := entry.Image.Bounds()
		if imgBounds.Dx() <= 0 || imgBounds.Dy() <= 0 {
			continue
		}

		objectWidth := obj.Width
		if objectWidth <= 0 {
			objectWidth = float64(imgBounds.Dx())
		}
		objectHeight := obj.Height
		if objectHeight <= 0 {
			objectHeight = float64(imgBounds.Dy())
		}

		screenCenterX := originX + int((obj.X/float64(g.doc.TileWidth))*float64(cellPx))
		screenCenterY := originY + int((obj.Y/float64(g.doc.TileHeight))*float64(cellPx))
		screenWidth := int((objectWidth / float64(g.doc.TileWidth)) * float64(cellPx))
		screenHeight := int((objectHeight / float64(g.doc.TileHeight)) * float64(cellPx))
		if screenWidth <= 0 {
			screenWidth = 1
		}
		if screenHeight <= 0 {
			screenHeight = 1
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(imgBounds.Dx())/2, -float64(imgBounds.Dy())/2)
		op.GeoM.Scale(
			float64(screenWidth)/float64(imgBounds.Dx()),
			float64(screenHeight)/float64(imgBounds.Dy()),
		)
		op.GeoM.Rotate(obj.Rotation * (math.Pi / 180))
		op.GeoM.Translate(float64(screenCenterX), float64(screenCenterY))
		screen.DrawImage(entry.Image, op)
	}
}

func (g *Game) drawCursorPreview(screen *ebiten.Image, cellPx, originX, originY int) {
	selected := g.palette[g.selectedPalette]
	if (selected.Layer != gamepalette.LayerProps && selected.Layer != gamepalette.LayerPropsTall) || selected.Image == nil {
		return
	}

	worldX, worldY, ok := g.cursorWorldPosition()
	if !ok || g.doc.TileWidth <= 0 || g.doc.TileHeight <= 0 {
		return
	}

	imgBounds := selected.Image.Bounds()
	if imgBounds.Dx() <= 0 || imgBounds.Dy() <= 0 {
		return
	}

	screenCenterX := originX + int((worldX/float64(g.doc.TileWidth))*float64(cellPx))
	screenCenterY := originY + int((worldY/float64(g.doc.TileHeight))*float64(cellPx))
	screenWidth := int((float64(imgBounds.Dx()) / float64(g.doc.TileWidth)) * float64(cellPx))
	screenHeight := int((float64(imgBounds.Dy()) / float64(g.doc.TileHeight)) * float64(cellPx))
	if screenWidth <= 0 {
		screenWidth = 1
	}
	if screenHeight <= 0 {
		screenHeight = 1
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(imgBounds.Dx())/2, -float64(imgBounds.Dy())/2)
	op.GeoM.Scale(
		float64(screenWidth)/float64(imgBounds.Dx()),
		float64(screenHeight)/float64(imgBounds.Dy()),
	)
	op.GeoM.Rotate(float64(g.selectedRotation) * (math.Pi / 180))
	op.GeoM.Translate(float64(screenCenterX), float64(screenCenterY))
	op.ColorScale.ScaleAlpha(0.65)
	screen.DrawImage(selected.Image, op)
}

func (g *Game) drawTile(screen *ebiten.Image, gid, x, y, cellPx int) {
	baseGID, rotation := tiled.DecodeRotatedGID(gid)
	entry := g.paletteByGID(baseGID)
	if entry == nil || entry.Image == nil {
		ebitenutil.DrawRect(screen, float64(x), float64(y), float64(cellPx), float64(cellPx), entryBgColor)
		return
	}

	imageBounds := entry.Image.Bounds()
	if imageBounds.Dx() <= 0 || imageBounds.Dy() <= 0 {
		ebitenutil.DrawRect(screen, float64(x), float64(y), float64(cellPx), float64(cellPx), entryBgColor)
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(imageBounds.Dx())/2, -float64(imageBounds.Dy())/2)
	op.GeoM.Scale(float64(cellPx)/float64(imageBounds.Dx()), float64(cellPx)/float64(imageBounds.Dy()))
	op.GeoM.Rotate(float64(rotation) * (math.Pi / 180))
	op.GeoM.Translate(float64(x+cellPx/2), float64(y+cellPx/2))
	screen.DrawImage(entry.Image, op)
}

func (g *Game) drawPalettePanel(screen *ebiten.Image) {
	panelX, panelY, panelW, panelH := g.paletteRect()
	ebitenutil.DrawRect(screen, float64(panelX), float64(panelY), float64(panelW), float64(panelH), panelBgColor)

	ebitenutil.DrawLine(screen, float64(panelX), float64(panelY), float64(panelX+panelW), float64(panelY), panelBorderColor)
	ebitenutil.DrawLine(screen, float64(panelX), float64(panelY+panelH), float64(panelX+panelW), float64(panelY+panelH), panelBorderColor)
	ebitenutil.DrawLine(screen, float64(panelX), float64(panelY), float64(panelX), float64(panelY+panelH), panelBorderColor)
	ebitenutil.DrawLine(screen, float64(panelX+panelW), float64(panelY), float64(panelX+panelW), float64(panelY+panelH), panelBorderColor)

	ebitenutil.DebugPrintAt(screen, "Sprite Palette (wheel)", panelX+10, panelY+8)

	listX := panelX + 10
	listY := panelY + paletteHeaderHeight
	listW := panelW - 20
	listH := panelH - paletteHeaderHeight - 10

	for i, entry := range g.palette {
		entryY := listY + i*(paletteEntryHeight+paletteEntryGap) - g.paletteScrollY
		entryBottom := entryY + paletteEntryHeight
		if entryBottom < listY || entryY > listY+listH {
			continue
		}

		entryColor := entryBgColor
		if i == g.selectedPalette {
			entryColor = entrySelectedColor
		}

		ebitenutil.DrawRect(screen, float64(listX), float64(entryY), float64(listW), float64(paletteEntryHeight), entryColor)

		thumbSize := paletteEntryHeight - 12
		thumbX := listX + 6
		thumbY := entryY + 6
		if entry.Image != nil {
			bounds := entry.Image.Bounds()
			if bounds.Dx() > 0 && bounds.Dy() > 0 {
				scale := mathutil.MinFloat64(float64(thumbSize)/float64(bounds.Dx()), float64(thumbSize)/float64(bounds.Dy()))

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(scale, scale)
				drawW := int(float64(bounds.Dx()) * scale)
				drawH := int(float64(bounds.Dy()) * scale)
				op.GeoM.Translate(float64(thumbX+(thumbSize-drawW)/2), float64(thumbY+(thumbSize-drawH)/2))
				screen.DrawImage(entry.Image, op)
			}
		}

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s (gid %d, %s)", entry.ID, entry.GID, entry.Layer), thumbX+thumbSize+8, entryY+18)
	}
}
