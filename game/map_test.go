package game

import "testing"

func TestDefaultMapIsLargerThanViewport(t *testing.T) {
	worldMap := DefaultMap()

	if !worldMap.IsLargerThanViewport(LogicalWidth, LogicalHeight) {
		t.Fatalf("default map must be larger than viewport: map=(%fx%f) viewport=(%dx%d)", worldMap.Width, worldMap.Height, LogicalWidth, LogicalHeight)
	}
}

func TestMapClampPosition(t *testing.T) {
	worldMap := Map{Width: 100, Height: 50}

	x, y := worldMap.ClampPosition(-10, 70)
	if x != 0 || y != 50 {
		t.Fatalf("expected clamped position (0,50), got (%f,%f)", x, y)
	}

	x, y = worldMap.ClampPosition(30, 20)
	if x != 30 || y != 20 {
		t.Fatalf("expected unchanged position (30,20), got (%f,%f)", x, y)
	}
}
