package main

import (
	"fmt"
	"image/color"

	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/co0p/tankism-revival/lib/tiled"
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
		if err := g.doc.ClearTallPropsObjects(); err != nil {
			return err
		}
		g.status = "cleared ground tiles and props layers"
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if err := g.setSpawnAtCursor(tiled.PlayerSpawnObjectName); err != nil {
			g.status = err.Error()
		} else {
			g.status = "set player spawn"
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			removed, err := g.removeNearestSpawnAtCursor(tiled.EnemySpawnObjectName)
			if err != nil {
				g.status = err.Error()
			} else if removed {
				g.status = "removed nearest enemy spawn"
			} else {
				g.status = "no enemy spawn near cursor"
			}
		} else {
			if err := g.addSpawnAtCursor(tiled.EnemySpawnObjectName); err != nil {
				g.status = err.Error()
			} else {
				count := len(g.doc.Spawns(tiled.EnemySpawnObjectName))
				g.status = fmt.Sprintf("added enemy spawn (%d total)", count)
			}
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
	if (selected.Layer == gamepalette.LayerProps || selected.Layer == gamepalette.LayerPropsTall) && leftJustPressed {
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
	if (selected.Layer == gamepalette.LayerProps || selected.Layer == gamepalette.LayerPropsTall) && rightJustPressed {
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
	g.drawPropsObjectsOnLayer(screen, tiled.PropsLayerName, cellPx, originX, originY)
	g.drawPropsObjectsOnLayer(screen, tiled.PropsTallLayerName, cellPx, originX, originY)

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

	g.drawSpawns(screen, tiled.PlayerSpawnObjectName, playerSpawnColor, cellPx, originX, originY)
	g.drawSpawns(screen, tiled.EnemySpawnObjectName, enemySpawnColor, cellPx, originX, originY)
	g.drawCursorPreview(screen, cellPx, originX, originY)
	g.drawPalettePanel(screen)

	hoverInfo := ""
	if cx, cy, ok := g.hoveredCell(); ok {
		groundRawGID, _ := g.doc.TileGIDAtLayer(tiled.GroundLayerName, cx, cy)
		groundBaseGID, groundRotation := tiled.DecodeRotatedGID(groundRawGID)
		hoverInfo = fmt.Sprintf(" | cell: (%d,%d) ground=%d@%d°", cx, cy, groundBaseGID, groundRotation)
	}

	propsCount := len(g.doc.AllPropsObjects())
	enemySpawnsCount := len(g.doc.Spawns(tiled.EnemySpawnObjectName))

	selected := g.palette[g.selectedPalette]
	message := fmt.Sprintf(
		"Mapism (%s)\nScroll wheel over palette: select sprite | Left drag: paint ground / left click: place props | Right drag / right click: erase\nP: set player spawn | E: add enemy spawn | Shift+E: remove nearest enemy spawn | N: clear ground+props layers | Ctrl+Shift+Arrows: resize map\nR/Shift+R: rotate ±1° (hold; ground snaps to 90°) | Ctrl+S: save | Ctrl+L: reload | G: toggle grid%s\nselected: %s (gid=%d layer=%s rot=%d°) | map: %dx%d | props: %d | enemy spawns: %d | status: %s",
		g.mapPath,
		hoverInfo,
		selected.ID,
		selected.GID,
		selected.Layer,
		g.selectedRotation,
		g.doc.Width,
		g.doc.Height,
		propsCount,
		enemySpawnsCount,
		g.status,
	)

	ebitenutil.DebugPrintAt(screen, message, 12, 12)
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return LogicalWidth, LogicalHeight
}
