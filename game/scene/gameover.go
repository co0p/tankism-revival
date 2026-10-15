package scene

import (
	"fmt"
	"image/color"

	"github.com/co0p/tankism-revival/lib/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	gameOverItemCount    = 2
	gameOverIndexRestart = 0
	gameOverIndexMenu    = 1
	outcomeContextKey    = "outcome"
	outcomeVictory       = "victory"
	outcomeDefeat        = "defeat"
)

// GameOverScene displays the outcome (victory or defeat) and offers restart or menu options.
type GameOverScene struct {
	outcome         string
	selectedIndex   int
	logicalWidth    int
	logicalHeight   int
	textRenderer    *fonts.Renderer
	blurShader      *ebiten.Shader
	backgroundImage *ebiten.Image // raw snapshot from PlayScene
	blurredBG       *ebiten.Image // cached blurred result
}

var gameOverItems = []string{"Restart", "Main Menu"}

// NewGameOverScene creates a game-over scene.
// Context must contain "outcome" = "victory" or "defeat" and a font loader at FontLoaderContextKey.
func NewGameOverScene(context map[string]interface{}) (Scene, error) {
	outcome := outcomeDefeat
	if context != nil {
		if o, ok := context[outcomeContextKey].(string); ok {
			outcome = o
		}
	}

	loader, ok := context[FontLoaderContextKey].(*fonts.Loader)
	if !ok || loader == nil {
		return nil, fmt.Errorf("scene/gameover: font loader missing from context")
	}
	face, err := loader.LoadKenneyBlocksFace(24)
	if err != nil {
		return nil, fmt.Errorf("scene/gameover: load font face: %w", err)
	}
	renderer := fonts.NewRenderer(face)

	blurShader, ok := context[BlurShaderContextKey].(*ebiten.Shader)
	if !ok || blurShader == nil {
		return nil, fmt.Errorf("scene/gameover: blur shader missing from context")
	}

	var bgImage *ebiten.Image
	if bg, ok := context[BackgroundImageContextKey].(*ebiten.Image); ok {
		bgImage = bg
	}

	return &GameOverScene{
		outcome:         outcome,
		selectedIndex:   gameOverIndexRestart,
		logicalWidth:    960,
		logicalHeight:   540,
		textRenderer:    renderer,
		blurShader:      blurShader,
		backgroundImage: bgImage,
	}, nil
}

// Update handles menu navigation and selection.
func (gos *GameOverScene) Update(deltaSeconds float64) SceneResult {
	// Navigate with arrow keys
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		gos.selectedIndex = (gos.selectedIndex - 1 + gameOverItemCount) % gameOverItemCount
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		gos.selectedIndex = (gos.selectedIndex + 1) % gameOverItemCount
	}

	// Select with Enter or Space
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch gos.selectedIndex {
		case gameOverIndexRestart:
			return NewTransitionToScene(SceneIDPlay)
		case gameOverIndexMenu:
			t := NewTransitionToScene(SceneIDMenu)
			if gos.backgroundImage != nil {
				t = t.WithContext(BackgroundImageContextKey, gos.backgroundImage)
			}
			return t
		}
	}

	return StayInScene{}
}

// Draw renders the game-over screen with outcome and options.
func (gos *GameOverScene) Draw(screen *ebiten.Image) {
	w, h := gos.logicalWidth, gos.logicalHeight

	// Build blurred background on first draw; cached for all subsequent frames.
	if gos.blurredBG == nil {
		gos.blurredBG = buildBlurredBG(gos.backgroundImage, gos.blurShader, w, h)
	}
	screen.DrawImage(gos.blurredBG, nil)

	// Dark overlay improves text legibility.
	overlay := ebiten.NewImage(w, h)
	overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 192})
	screen.DrawImage(overlay, nil)
	overlay.Dispose()

	// Outcome title coloured by result.
	outcomeTitle := "DEFEAT"
	outcomeColor := ColorDefeat
	if gos.outcome == outcomeVictory {
		outcomeTitle = "VICTORY"
		outcomeColor = ColorVictory
	}
	titleX := w/2 - 30
	titleY := h/2 - 120
	gos.drawOutcomeText(screen, outcomeTitle, titleX, titleY, outcomeColor)

	// Menu items.
	menuStartX := w/2 - 40
	menuStartY := h / 2

	for i, item := range gameOverItems {
		y := menuStartY + i*40
		x := menuStartX

		var icon string
		itemColor := ColorUnselected
		if i == gos.selectedIndex {
			icon = SelectedIcon
			itemColor = ColorSelected
		} else {
			icon = SpacingIcon
		}

		displayText := icon + " " + item
		gos.drawOutcomeText(screen, displayText, x, y, itemColor)
	}

	hintY := h - 60
	gos.drawOutcomeText(screen, "Use Up/Down to navigate, Enter to select", 180, hintY, ColorUnselected)
}

// drawOutcomeText renders text at (x, y) with the given color using the Kenney Blocks font.
func (gos *GameOverScene) drawOutcomeText(screen *ebiten.Image, text string, x, y int, clr color.Color) {
	if err := gos.textRenderer.DrawText(screen, text, x, y, clr); err != nil {
		panic(fmt.Sprintf("scene/gameover: DrawText: %v", err))
	}
}
