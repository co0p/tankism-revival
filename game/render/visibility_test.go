package render

import (
	"image"
	"testing"
)

func TestEffectVisible(t *testing.T) {
	view := image.Rect(0, 0, 320, 180)
	for _, test := range []struct {
		x, y, radius float64
		want         bool
	}{
		{100, 100, 9, true}, {-5, 50, 9, true}, {324, 50, 9, true}, {50, -5, 9, true}, {50, 184, 9, true},
		{-10, 50, 9, false}, {330, 50, 9, false}, {50, -10, 9, false}, {50, 190, 9, false},
	} {
		if got := effectVisible(view, test.x, test.y, test.radius); got != test.want {
			t.Fatalf("visibility at (%v,%v): got %v want %v", test.x, test.y, got, test.want)
		}
	}
}
