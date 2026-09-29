package mathutil

import (
	"math"
	"testing"
)

func TestClamp(t *testing.T) {
	if got := Clamp(-1, 0, 10); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
	if got := Clamp(11, 0, 10); got != 10 {
		t.Fatalf("expected 10, got %f", got)
	}
	if got := Clamp(4.5, 0, 10); got != 4.5 {
		t.Fatalf("expected 4.5, got %f", got)
	}
}

func TestClamp01AndClampUnit(t *testing.T) {
	if got := Clamp01(-0.5); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
	if got := Clamp01(1.5); got != 1 {
		t.Fatalf("expected 1, got %f", got)
	}
	if got := ClampUnit(-2); got != -1 {
		t.Fatalf("expected -1, got %f", got)
	}
	if got := ClampUnit(2); got != 1 {
		t.Fatalf("expected 1, got %f", got)
	}
}

func TestMoveTowards(t *testing.T) {
	if got := MoveTowards(0, 10, 3); got != 3 {
		t.Fatalf("expected 3, got %f", got)
	}
	if got := MoveTowards(10, 0, 4); got != 6 {
		t.Fatalf("expected 6, got %f", got)
	}
	if got := MoveTowards(1, 2, 0); got != 2 {
		t.Fatalf("expected target when maxDelta<=0, got %f", got)
	}
}

func TestNormalizeRadians(t *testing.T) {
	if got := NormalizeRadians(-math.Pi / 2); got <= 0 || got >= 2*math.Pi {
		t.Fatalf("expected wrapped angle in [0,2pi), got %f", got)
	}
	if got := NormalizeRadians(5 * math.Pi); math.Abs(got-math.Pi) > 1e-9 {
		t.Fatalf("expected pi, got %f", got)
	}
}
