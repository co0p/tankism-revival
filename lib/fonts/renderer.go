package fonts

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Renderer renders text with a loaded font face.
type Renderer struct {
	face *text.GoTextFace
}

// NewRenderer creates a text renderer. face must not be nil.
func NewRenderer(face *text.GoTextFace) *Renderer {
	if face == nil {
		panic("fonts.NewRenderer: face must not be nil")
	}
	return &Renderer{face: face}
}

// DrawText renders txt at (x, y) with the given color.
func (r *Renderer) DrawText(screen *ebiten.Image, txt string, x, y int, clr color.Color) error {
	if r.face == nil {
		return fmt.Errorf("fonts.Renderer: no font face set")
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)

	text.Draw(screen, txt, r.face, op)
	return nil
}

// Face returns the current font face.
func (r *Renderer) Face() *text.GoTextFace {
	return r.face
}
