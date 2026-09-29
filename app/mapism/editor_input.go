package main

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/lib/mathutil"
	"github.com/co0p/tankism-revival/lib/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) handleMapResizeShortcuts() error {
	ctrl := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)
	shift := ebiten.IsKeyPressed(ebiten.KeyShift)
	if !ctrl || !shift {
		return nil
	}

	newWidth := g.doc.Width
	newHeight := g.doc.Height

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		newWidth += defaultResizeStep
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		newWidth -= defaultResizeStep
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		newHeight += defaultResizeStep
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		newHeight -= defaultResizeStep
	}

	if newWidth == g.doc.Width && newHeight == g.doc.Height {
		return nil
	}

	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}

	if err := g.doc.Resize(newWidth, newHeight); err != nil {
		return err
	}

	g.status = fmt.Sprintf("resized map: %dx%d", g.doc.Width, g.doc.Height)
	return nil
}

func (g *Game) handlePaletteWheelSelection() {
	if !g.cursorInPalettePanel() {
		return
	}

	_, wheelY := ebiten.Wheel()
	if wheelY == 0 {
		return
	}

	steps := int(math.Ceil(math.Abs(wheelY)))
	if steps < 1 {
		steps = 1
	}

	delta := 1
	if wheelY > 0 {
		delta = -1
	}

	for range steps {
		g.selectedPalette = mathutil.ClampInt(g.selectedPalette+delta, 0, len(g.palette)-1)
	}

	g.ensurePaletteSelectionVisible()
	selected := g.palette[g.selectedPalette]
	g.status = fmt.Sprintf("selected sprite: %s (layer: %s, rot: %d°)", selected.ID, selected.Layer, g.selectedRotation)
}

func (g *Game) handleRotationInput() {
	pressedFrames := inpututil.KeyPressDuration(ebiten.KeyR)
	if pressedFrames <= 0 {
		return
	}

	if pressedFrames > 1 {
		if pressedFrames <= rotationRepeatDelayFrames || pressedFrames%rotationRepeatEveryFrames != 0 {
			return
		}
	}

	delta := rotationStepDegree
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		delta = -rotationStepDegree
	}

	g.selectedRotation = tiled.NormalizeRotationDegrees(g.selectedRotation + delta)
	selected := g.palette[g.selectedPalette]
	direction := "clockwise"
	if delta < 0 {
		direction = "counter-clockwise"
	}
	g.status = fmt.Sprintf("rotation set to %d° (%s) for %s", g.selectedRotation, direction, selected.ID)
}

func (g *Game) selectPaletteEntryAtCursor() bool {
	mx, my := ebiten.CursorPosition()
	idx, ok := g.paletteEntryAt(mx, my)
	if !ok {
		return false
	}

	g.selectedPalette = idx
	g.ensurePaletteSelectionVisible()
	return true
}

func (g *Game) paletteEntryAt(screenX, screenY int) (int, bool) {
	panelX, panelY, panelW, panelH := g.paletteRect()
	if screenX < panelX || screenX >= panelX+panelW || screenY < panelY || screenY >= panelY+panelH {
		return 0, false
	}

	listX := panelX + 10
	listY := panelY + paletteHeaderHeight
	listW := panelW - 20
	listH := panelH - paletteHeaderHeight - 10

	if screenX < listX || screenX >= listX+listW || screenY < listY || screenY >= listY+listH {
		return 0, false
	}

	for i := range g.palette {
		entryY := listY + i*(paletteEntryHeight+paletteEntryGap) - g.paletteScrollY
		if screenY >= entryY && screenY < entryY+paletteEntryHeight {
			return i, true
		}
	}

	return 0, false
}

func (g *Game) ensurePaletteSelectionVisible() {
	_, _, _, panelH := g.paletteRect()
	listH := panelH - paletteHeaderHeight - 10
	if listH <= 0 {
		g.paletteScrollY = 0
		return
	}

	entryStride := paletteEntryHeight + paletteEntryGap
	entryTop := g.selectedPalette * entryStride
	entryBottom := entryTop + paletteEntryHeight

	viewTop := g.paletteScrollY
	viewBottom := g.paletteScrollY + listH

	if entryTop < viewTop {
		g.paletteScrollY = entryTop
	} else if entryBottom > viewBottom {
		g.paletteScrollY = entryBottom - listH
	}

	maxScroll := mathutil.MaxInt(0, len(g.palette)*entryStride-listH)
	g.paletteScrollY = mathutil.ClampInt(g.paletteScrollY, 0, maxScroll)
}

func (g *Game) paletteByGID(gid int) *paletteEntry {
	for i := range g.palette {
		if g.palette[i].GID == gid {
			return &g.palette[i]
		}
	}

	return nil
}

func (g *Game) cursorInPalettePanel() bool {
	mx, my := ebiten.CursorPosition()
	panelX, panelY, panelW, panelH := g.paletteRect()
	return mx >= panelX && mx < panelX+panelW && my >= panelY && my < panelY+panelH
}

func (g *Game) paletteRect() (x, y, w, h int) {
	return LogicalWidth - palettePanelWidth - uiMargin, uiMargin, palettePanelWidth, LogicalHeight - uiMargin*2
}
