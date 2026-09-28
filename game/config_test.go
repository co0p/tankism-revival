package game

import "testing"

func TestLogicalAspectRatioIsSixteenByNine(t *testing.T) {
	got := LogicalAspectRatio()
	want := 16.0 / 9.0

	if got != want {
		t.Fatalf("logical aspect ratio mismatch: got %f want %f", got, want)
	}
}
