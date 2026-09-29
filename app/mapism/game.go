package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"

	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/co0p/tankism-revival/lib/sprites"
	"github.com/co0p/tankism-revival/lib/tiled"
	gameassets "github.com/co0p/tankism-revival/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	LogicalWidth              = 1280
	LogicalHeight             = 720
	DefaultMapWidth           = 40
	DefaultMapHeight          = 24
	DefaultTileWidth          = 32
	DefaultTileHeight         = 32
	DefaultMapPath            = "resources/maps/mapism_default.tmj"
	defaultResizeStep         = 1
	rotationStepDegree        = 1
	rotationRepeatDelayFrames = 12
	rotationRepeatEveryFrames = 2
	uiMargin                  = 12
	palettePanelWidth         = 300
	paletteHeaderHeight       = 34
	paletteEntryHeight        = 56
	paletteEntryGap           = 6
)

var (
	bgColor            = color.RGBA{R: 0x12, G: 0x16, B: 0x1d, A: 0xff}
	mapBgColor         = color.RGBA{R: 0x18, G: 0x1f, B: 0x28, A: 0xff}
	gridLineColor      = color.RGBA{R: 0x2e, G: 0x36, B: 0x42, A: 0xff}
	panelBgColor       = color.RGBA{R: 0x11, G: 0x14, B: 0x1a, A: 0xff}
	panelBorderColor   = color.RGBA{R: 0x3a, G: 0x44, B: 0x51, A: 0xff}
	entryBgColor       = color.RGBA{R: 0x1c, G: 0x23, B: 0x2d, A: 0xff}
	entrySelectedColor = color.RGBA{R: 0x38, G: 0x66, B: 0x99, A: 0xff}
	playerSpawnColor   = color.RGBA{R: 0x40, G: 0xe0, B: 0x72, A: 0xff}
	enemySpawnColor    = color.RGBA{R: 0xe0, G: 0x59, B: 0x59, A: 0xff}
)

type Game struct {
	mapPath          string
	doc              tiled.Map
	status           string
	showGrid         bool
	palette          []paletteEntry
	selectedPalette  int
	paletteScrollY   int
	selectedRotation int
}

func NewGame(mapPath string, mapWidth, mapHeight int) (*Game, error) {
	if mapPath == "" {
		mapPath = DefaultMapPath
	}
	if mapWidth <= 0 {
		mapWidth = DefaultMapWidth
	}
	if mapHeight <= 0 {
		mapHeight = DefaultMapHeight
	}

	doc, status, err := loadOrCreate(mapPath, mapWidth, mapHeight)
	if err != nil {
		return nil, err
	}

	palette, err := loadPaletteEntries()
	if err != nil {
		return nil, err
	}

	g := &Game{
		mapPath:          mapPath,
		doc:              doc,
		status:           status,
		showGrid:         true,
		palette:          palette,
		selectedPalette:  0,
		paletteScrollY:   0,
		selectedRotation: 0,
	}

	if err := g.syncTileset(); err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) Update() error {
	g.handlePaletteWheelSelection()
	if err := g.handleMapResizeShortcuts(); err != nil {
		return err
	}
	g.handleRotationInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.showGrid = !g.showGrid
		if g.showGrid {
			g.status = "grid enabled"
		} else {
			g.status = "grid hidden"
		}
	}

	ctrl := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)

	if ctrl && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if err := g.syncTileset(); err != nil {
			return err
		}
		if err := tiled.WriteFile(g.mapPath, g.doc); err != nil {
			return fmt.Errorf("save map: %w", err)
		}
		g.status = fmt.Sprintf("saved map: %s", g.mapPath)
	}

	if ctrl && inpututil.IsKeyJustPressed(ebiten.KeyL) {
		doc, err := tiled.ReadFile(g.mapPath)
		if err != nil {
			return fmt.Errorf("load map: %w", err)
		}
		g.doc = doc
		if err := g.syncTileset(); err != nil {
			return err
		}
		g.status = fmt.Sprintf("loaded map: %s", g.mapPath)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		if err := g.doc.ClearTileLayer(tiled.GroundLayerName); err != nil {
			return err
		}
		if err := g.doc.ClearPropsObjects(); err != nil {
			return err
		}
		g.status = "cleared ground tiles and props objects"
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if err := g.setSpawnAtCursor(tiled.PlayerSpawnObjectName); err != nil {
			g.status = err.Error()
		} else {
			g.status = "set player spawn"
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		if err := g.setSpawnAtCursor(tiled.EnemySpawnObjectName); err != nil {
			g.status = err.Error()
		} else {
			g.status = "set enemy spawn"
		}
	}

	leftPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	leftJustPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	rightPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	rightJustPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

	if leftJustPressed {
		if g.selectPaletteEntryAtCursor() {
			selected := g.palette[g.selectedPalette]
			g.status = fmt.Sprintf("selected sprite: %s (layer: %s, rot: %d°)", selected.ID, selected.Layer, g.selectedRotation)
			return nil
		}
	}

	selected := g.palette[g.selectedPalette]
	if selected.Layer == gamepalette.LayerGround && leftPressed {
		if err := g.paintGroundTileAtCursor(); err != nil {
			g.status = err.Error()
		}
	}
	if selected.Layer == gamepalette.LayerProps && leftJustPressed {
		placed, err := g.placePropObjectAtCursor()
		if err != nil {
			g.status = err.Error()
		} else if placed {
			g.status = fmt.Sprintf("placed prop: %s", selected.ID)
		}
	}

	if selected.Layer == gamepalette.LayerGround && rightPressed {
		if err := g.eraseGroundTileAtCursor(); err != nil {
			g.status = err.Error()
		}
	}
	if selected.Layer == gamepalette.LayerProps && rightJustPressed {
		removed, err := g.erasePropObjectAtCursor()
		if err != nil {
			g.status = err.Error()
		} else if removed {
			g.status = "removed nearest prop"
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)

	cellPx, originX, originY := g.viewMetrics()
	drawWidth := cellPx * g.doc.Width
	drawHeight := cellPx * g.doc.Height

	ebitenutil.DrawRect(screen, float64(originX), float64(originY), float64(drawWidth), float64(drawHeight), mapBgColor)

	g.drawTileLayer(screen, tiled.GroundLayerName, cellPx, originX, originY)
	g.drawPropsObjects(screen, cellPx, originX, originY)

	if g.showGrid {
		for col := 0; col <= g.doc.Width; col++ {
			x := originX + col*cellPx
			ebitenutil.DrawLine(screen, float64(x), float64(originY), float64(x), float64(originY+drawHeight), gridLineColor)
		}

		for row := 0; row <= g.doc.Height; row++ {
			y := originY + row*cellPx
			ebitenutil.DrawLine(screen, float64(originX), float64(y), float64(originX+drawWidth), float64(y), gridLineColor)
		}
	}

	g.drawSpawn(screen, tiled.PlayerSpawnObjectName, playerSpawnColor, cellPx, originX, originY)
	g.drawSpawn(screen, tiled.EnemySpawnObjectName, enemySpawnColor, cellPx, originX, originY)
	g.drawCursorPreview(screen, cellPx, originX, originY)
	g.drawPalettePanel(screen)

	hoverInfo := ""
	if cx, cy, ok := g.hoveredCell(); ok {
		groundRawGID, _ := g.doc.TileGIDAtLayer(tiled.GroundLayerName, cx, cy)
		groundBaseGID, groundRotation := tiled.DecodeRotatedGID(groundRawGID)
		hoverInfo = fmt.Sprintf(" | cell: (%d,%d) ground=%d@%d°", cx, cy, groundBaseGID, groundRotation)
	}

	propsCount := len(g.doc.PropsObjects())

	selected := g.palette[g.selectedPalette]
	message := fmt.Sprintf(
		"Mapism (%s)\nScroll wheel over palette: select sprite | Left drag: paint ground / left click: place props | Right drag / right click: erase\nP: set player spawn | E: set enemy spawn | N: clear ground+props | Ctrl+Shift+Arrows: resize map\nR/Shift+R: rotate ±1° (hold; ground snaps to 90°) | Ctrl+S: save | Ctrl+L: reload | G: toggle grid%s\nselected: %s (gid=%d layer=%s rot=%d°) | map: %dx%d | props: %d | status: %s",
		g.mapPath,
		hoverInfo,
		selected.ID,
		selected.GID,
		selected.Layer,
		g.selectedRotation,
		g.doc.Width,
		g.doc.Height,
		propsCount,
		g.status,
	)

	ebitenutil.DebugPrintAt(screen, message, 12, 12)
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return LogicalWidth, LogicalHeight
}

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

	if err := g.doc.AddPropsObject(
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
	return g.doc.RemoveNearestPropsObject(worldX, worldY, maxDistance)
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
	cellPx = min(cellByWidth, cellByHeight)
	if cellPx < 2 {
		cellPx = 2
	}

	drawWidth := cellPx * g.doc.Width
	drawHeight := cellPx * g.doc.Height
	originX = uiMargin + (availableWidth-drawWidth)/2
	originY = uiMargin + (availableHeight-drawHeight)/2

	return cellPx, originX, originY
}

func (g *Game) drawSpawn(screen *ebiten.Image, spawnName string, markerColor color.Color, cellPx, originX, originY int) {
	spawn, ok := g.doc.Spawn(spawnName)
	if !ok || g.doc.TileWidth <= 0 || g.doc.TileHeight <= 0 {
		return
	}

	screenX := originX + int((spawn.X/float64(g.doc.TileWidth))*float64(cellPx))
	screenY := originY + int((spawn.Y/float64(g.doc.TileHeight))*float64(cellPx))

	size := max(5, cellPx/3)
	ebitenutil.DrawRect(
		screen,
		float64(screenX-size/2),
		float64(screenY-size/2),
		float64(size),
		float64(size),
		markerColor,
	)
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

func (g *Game) drawPropsObjects(screen *ebiten.Image, cellPx, originX, originY int) {
	for _, obj := range g.doc.PropsObjects() {
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
	if selected.Layer != gamepalette.LayerProps || selected.Image == nil {
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
				scale := minFloat(float64(thumbSize)/float64(bounds.Dx()), float64(thumbSize)/float64(bounds.Dy()))

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
		g.selectedPalette = clampInt(g.selectedPalette+delta, 0, len(g.palette)-1)
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

	maxScroll := max(0, len(g.palette)*entryStride-listH)
	g.paletteScrollY = clampInt(g.paletteScrollY, 0, maxScroll)
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

func (g *Game) syncTileset() error {
	if len(g.palette) == 0 {
		return fmt.Errorf("palette is empty")
	}

	mapDir := filepath.Dir(g.mapPath)
	tiles := make([]tiled.TilesetTile, 0, len(g.palette))

	for i, entry := range g.palette {
		bounds := entry.Image.Bounds()
		imagePath := filepath.Join("resources", filepath.FromSlash(entry.AssetPath))
		relPath, err := filepath.Rel(mapDir, imagePath)
		if err != nil {
			relPath = imagePath
		}

		tiles = append(tiles, tiled.TilesetTile{
			ID:          i,
			Image:       filepath.ToSlash(relPath),
			ImageWidth:  bounds.Dx(),
			ImageHeight: bounds.Dy(),
		})
	}

	g.doc.Tilesets = []tiled.Tileset{
		{
			FirstGID:   1,
			Name:       "mapism_palette",
			TileWidth:  g.doc.TileWidth,
			TileHeight: g.doc.TileHeight,
			TileCount:  len(g.palette),
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

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}

	return value
}
