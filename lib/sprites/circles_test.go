package sprites

import (
	"image"
	"image/color"
	"testing"
	"testing/fstest"

	"github.com/co0p/tankism-revival/resources"
)

func TestExtractCircles(t *testing.T) {
	source := image.NewNRGBA(image.Rect(0, 0, 8, 6))
	source.Set(1, 1, color.NRGBA{R: 120, A: 255})
	source.Set(2, 1, color.NRGBA{R: 130, A: 128})
	source.Set(6, 4, color.NRGBA{G: 140, A: 255})
	// Almost-transparent export residue must not connect separate circles.
	for y := 2; y <= 4; y++ {
		source.Set(2, y, color.NRGBA{A: 3})
	}
	for x := 3; x < 6; x++ {
		source.Set(x, 4, color.NRGBA{A: 3})
	}
	circles := ExtractCircles(source)
	if len(circles) != 2 || circles[0].Bounds() != image.Rect(0, 0, 2, 1) || circles[1].Bounds() != image.Rect(0, 0, 1, 1) {
		t.Fatal("expected two tightly cropped components in scan order")
	}
	if circles[0].NRGBAAt(1, 0) != source.NRGBAAt(2, 1) {
		t.Fatal("extraction did not preserve original color and alpha")
	}
}

func TestKenneySmokeCircles(t *testing.T) {
	circles, err := LoadCircles(resources.FS, resources.SmokeCirclesPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(circles) != 8 {
		t.Fatalf("expected the eight gray smoke circles, got %d", len(circles))
	}
	for _, circle := range circles {
		if circle.Bounds().Dx() < 5 || circle.Bounds().Dy() < 5 {
			t.Fatal("extracted a fragment rather than a smoke circle")
		}
	}
}

func TestLoadCirclesErrors(t *testing.T) {
	if _, err := LoadCircles(fstest.MapFS{}, "missing.png"); err == nil {
		t.Fatal("missing asset should fail")
	}
	if _, err := LoadCircles(fstest.MapFS{"bad.png": {Data: []byte("not a PNG")}}, "bad.png"); err == nil {
		t.Fatal("invalid image should fail")
	}
}
