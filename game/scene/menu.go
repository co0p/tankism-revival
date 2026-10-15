package scene

import (
	"fmt"
	"image/color"

	"github.com/co0p/tankism-revival/lib/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	menuItemCount  = 2
	menuIndexPlay  = 0
	menuIndexExit  = 1
	menuBlurRadius = 8.0
)

var (
	menuItems = []string{"Single Player", "Exit"}
)

// MenuScene displays a main menu with options to play or exit.
type MenuScene struct {
	selectedIndex int
	logicalWidth  int
	logicalHeight int
	blurShader    *ebiten.Shader
	blurredBG     *ebiten.Image // cached blurred background; nil until first draw
	rawBackground *ebiten.Image // unblurred snapshot passed in via context
	textRenderer  *fonts.Renderer
}

// NewMenuScene creates a menu scene.
// Context must contain FontLoaderContextKey and BlurShaderContextKey.
// It may optionally contain BackgroundImageContextKey with a game-world snapshot.
func NewMenuScene(context map[string]interface{}) (Scene, error) {
	loader, ok := context[FontLoaderContextKey].(*fonts.Loader)
	if !ok || loader == nil {
		return nil, fmt.Errorf("scene/menu: font loader missing from context")
	}
	face, err := loader.LoadKenneyBlocksFace(24)
	if err != nil {
		return nil, fmt.Errorf("scene/menu: load font face: %w", err)
	}

	blurShader, ok := context[BlurShaderContextKey].(*ebiten.Shader)
	if !ok || blurShader == nil {
		return nil, fmt.Errorf("scene/menu: blur shader missing from context")
	}

	var rawBG *ebiten.Image
	if bg, ok := context[BackgroundImageContextKey].(*ebiten.Image); ok && bg != nil {
		rawBG = bg
	}

	return &MenuScene{
		selectedIndex: menuIndexPlay,
		logicalWidth:  960,
		logicalHeight: 540,
		blurShader:    blurShader,
		rawBackground: rawBG,
		textRenderer:  fonts.NewRenderer(face),
	}, nil
}

// Update handles menu navigation and selection.
func (ms *MenuScene) Update(deltaSeconds float64) SceneResult {
	// Navigate with arrow keys
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		ms.selectedIndex = (ms.selectedIndex - 1 + menuItemCount) % menuItemCount
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		ms.selectedIndex = (ms.selectedIndex + 1) % menuItemCount
	}

	// Select with Enter or Space
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch ms.selectedIndex {
		case menuIndexPlay:
			return NewTransitionToScene(SceneIDPlay)
		case menuIndexExit:
			return QuitGame{}
		}
	}

	// Escape acts as exit from the top-level menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return QuitGame{}
	}

	return StayInScene{}
}

// Draw renders the menu with title and options.
func (ms *MenuScene) Draw(screen *ebiten.Image) {
	w, h := ms.logicalWidth, ms.logicalHeight

	// Build blurred background on first draw (cached for all subsequent frames).
	if ms.blurredBG == nil {
		ms.blurredBG = ms.buildBlurredBackground(w, h)
	}
	screen.DrawImage(ms.blurredBG, nil)

	// Semi-transparent dark overlay improves text legibility over any background.
	overlay := ebiten.NewImage(w, h)
	overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 140})
	screen.DrawImage(overlay, nil)
	overlay.Dispose()

	// Title
	titleX := w/2 - 65
	titleY := h/2 - 120
	ms.drawColoredText(screen, "Tankism Revival", titleX, titleY, ColorTitle)

	// Menu items
	menuStartX := w/2 - 40
	menuStartY := h / 2

	for i, item := range menuItems {
		y := menuStartY + i*40
		x := menuStartX

		var icon string
		itemColor := ColorUnselected
		if i == ms.selectedIndex {
			icon = SelectedIcon
			itemColor = ColorSelected
		} else {
			icon = SpacingIcon
		}

		displayText := icon + " " + item
		ms.drawColoredText(screen, displayText, x, y, itemColor)
	}

	// Hint
	hintY := h - 60
	ms.drawColoredText(screen, "Use Up/Down to navigate, Enter to select", 180, hintY, ColorUnselected)
}

// buildBlurredBG returns a blurred copy of src, or a solid dark background when src is nil.
// The result is owned by the caller; call Dispose when no longer needed.
func buildBlurredBG(src *ebiten.Image, shader *ebiten.Shader, w, h int) *ebiten.Image {
	result := ebiten.NewImage(w, h)

	if src == nil {
		result.Fill(color.RGBA{R: 0x1c, G: 0x27, B: 0x1b, A: 0xff})
		return result
	}

	// Two-pass separable Gaussian blur: horizontal then vertical.
	tmp := ebiten.NewImage(w, h)
	defer tmp.Dispose()

	uniforms := map[string]any{
		"Radius": float32(menuBlurRadius),
		"Axis":   float32(0),
	}
	op := &ebiten.DrawRectShaderOptions{
		Images:   [4]*ebiten.Image{src},
		Uniforms: uniforms,
	}
	tmp.DrawRectShader(w, h, shader, op)

	uniforms["Axis"] = float32(1)
	op2 := &ebiten.DrawRectShaderOptions{
		Images:   [4]*ebiten.Image{tmp},
		Uniforms: uniforms,
	}
	result.DrawRectShader(w, h, shader, op2)

	return result
}

// buildBlurredBackground is the MenuScene-specific wrapper.
func (ms *MenuScene) buildBlurredBackground(w, h int) *ebiten.Image {
	return buildBlurredBG(ms.rawBackground, ms.blurShader, w, h)
}

// drawColoredText renders text at (x, y) with the given color using the Kenney Blocks font.
func (ms *MenuScene) drawColoredText(screen *ebiten.Image, text string, x, y int, clr color.Color) {
	if err := ms.textRenderer.DrawText(screen, text, x, y, clr); err != nil {
		panic(fmt.Sprintf("scene/menu: DrawText: %v", err))
	}
}
