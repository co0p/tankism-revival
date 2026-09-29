package tiled

import "testing"

func TestQuantizeOrthogonalRotationDegrees(t *testing.T) {
	tests := []struct {
		rotation int
		expect   int
	}{
		{rotation: 0, expect: 0},
		{rotation: 1, expect: 0},
		{rotation: 44, expect: 0},
		{rotation: 45, expect: 90},
		{rotation: 134, expect: 90},
		{rotation: 135, expect: 180},
		{rotation: 224, expect: 180},
		{rotation: 225, expect: 270},
		{rotation: 314, expect: 270},
		{rotation: 315, expect: 0},
		{rotation: -1, expect: 0},
		{rotation: -136, expect: 180},
	}

	for _, tc := range tests {
		got := QuantizeOrthogonalRotationDegrees(tc.rotation)
		if got != tc.expect {
			t.Fatalf("rotation %d: expected %d, got %d", tc.rotation, tc.expect, got)
		}
	}
}

func TestEncodeRotatedGIDQuantizesRotation(t *testing.T) {
	base := 7

	encoded := EncodeRotatedGID(base, 89)
	decodedBase, decodedRotation := DecodeRotatedGID(encoded)
	if decodedBase != base {
		t.Fatalf("expected base gid %d, got %d", base, decodedBase)
	}
	if decodedRotation != 90 {
		t.Fatalf("expected quantized rotation 90, got %d", decodedRotation)
	}
}
