package scene

import (
	"fmt"
	"image/color"

	"github.com/co0p/tankism-revival/lib/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	defaultLoadingDurationSeconds = 3.0
	defaultLoadingBarWidth        = 200.0
	defaultLoadingBarHeight       = 16.0
)

// LoadingScene displays a fake progress bar for a fixed duration, then transitions to the menu.
type LoadingScene struct {
	durationSeconds float64
	elapsedSeconds  float64
	logicalWidth    int
	logicalHeight   int
	textRenderer    *fonts.Renderer
}

// NewLoadingScene creates a loading scene.
// Context must contain a font loader at FontLoaderContextKey.
func NewLoadingScene(context map[string]interface{}) (Scene, error) {
	loader, ok := context[FontLoaderContextKey].(*fonts.Loader)
	if !ok || loader == nil {
		return nil, fmt.Errorf("scene/loading: font loader missing from context")
	}
	face, err := loader.LoadKenneyBlocksFace(20)
	if err != nil {
		return nil, fmt.Errorf("scene/loading: load font face: %w", err)
	}
	renderer := fonts.NewRenderer(face)

	return &LoadingScene{
		durationSeconds: defaultLoadingDurationSeconds,
		elapsedSeconds:  0,
		logicalWidth:    960,
		logicalHeight:   540,
		textRenderer:    renderer,
	}, nil
}

// Update advances the loading progress and transitions when complete.
func (ls *LoadingScene) Update(deltaSeconds float64) SceneResult {
	ls.elapsedSeconds += deltaSeconds
	if ls.elapsedSeconds >= ls.durationSeconds {
		return NewTransitionToScene(SceneIDMenu)
	}
	return StayInScene{}
}

// Draw renders the loading screen with a progress bar and title.
func (ls *LoadingScene) Draw(screen *ebiten.Image) {
	// Fill background with dark color
	screen.Fill(ColorDarkBg)

	// Calculate progress (0.0 to 1.0)
	progress := ls.elapsedSeconds / ls.durationSeconds
	if progress > 1.0 {
		progress = 1.0
	}

	// Draw title with icon
	titleX := ls.logicalWidth / 2
	titleY := ls.logicalHeight/2 - 80
	titleText := LoadingIcon + " Tankism Revival"
	ls.drawLoadingText(screen, titleText, titleX-60, titleY, ColorTitle)

	// Draw progress bar background
	barX := float64(ls.logicalWidth)/2 - defaultLoadingBarWidth/2
	barY := float64(ls.logicalHeight) / 2
	vector.StrokeRect(screen, float32(barX), float32(barY), float32(defaultLoadingBarWidth), float32(defaultLoadingBarHeight), 2, ColorGreen, true)

	// Draw progress bar fill
	fillWidth := defaultLoadingBarWidth * progress
	vector.FillRect(screen, float32(barX), float32(barY), float32(fillWidth), float32(defaultLoadingBarHeight), ColorGreen, true)

	// Draw percentage text
	percentageInt := int(progress * 100)
	percentageText := "Loading... " + string(rune(percentageInt/10%10+48)) + string(rune(percentageInt%10+48)) + "%"
	ls.drawLoadingText(screen, percentageText, ls.logicalWidth/2-40, ls.logicalHeight/2+40, ColorTitle)
}

// drawLoadingText renders text at (x, y) with the given color using the Kenney Blocks font.
func (ls *LoadingScene) drawLoadingText(screen *ebiten.Image, text string, x, y int, clr color.Color) {
	if err := ls.textRenderer.DrawText(screen, text, x, y, clr); err != nil {
		panic(fmt.Sprintf("scene/loading: DrawText: %v", err))
	}
}
