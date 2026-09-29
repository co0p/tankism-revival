package tiled

import "math"

const (
	tiledFlipHorizontal = uint32(0x80000000)
	tiledFlipVertical   = uint32(0x40000000)
	tiledFlipDiagonal   = uint32(0x20000000)
	tiledFlipAll        = tiledFlipHorizontal | tiledFlipVertical | tiledFlipDiagonal
)

func NormalizeRotationDegrees(rotation int) int {
	rotation %= 360
	if rotation < 0 {
		rotation += 360
	}

	return rotation
}

func QuantizeOrthogonalRotationDegrees(rotation int) int {
	normalized := NormalizeRotationDegrees(rotation)
	steps := int(math.Round(float64(normalized) / 90.0))
	return NormalizeRotationDegrees(steps * 90)
}

func EncodeRotatedGID(baseGID int, rotationDegrees int) int {
	rotationDegrees = QuantizeOrthogonalRotationDegrees(rotationDegrees)
	if baseGID <= 0 {
		return 0
	}

	raw := uint32(baseGID)

	switch rotationDegrees {
	case 90:
		raw |= tiledFlipHorizontal | tiledFlipDiagonal
	case 180:
		raw |= tiledFlipHorizontal | tiledFlipVertical
	case 270:
		raw |= tiledFlipVertical | tiledFlipDiagonal
	}

	return int(raw)
}

func DecodeRotatedGID(rawGID int) (baseGID int, rotationDegrees int) {
	if rawGID <= 0 {
		return 0, 0
	}

	raw := uint32(rawGID)
	flags := raw & tiledFlipAll
	base := raw & ^tiledFlipAll

	switch flags {
	case tiledFlipHorizontal | tiledFlipDiagonal:
		rotationDegrees = 90
	case tiledFlipHorizontal | tiledFlipVertical:
		rotationDegrees = 180
	case tiledFlipVertical | tiledFlipDiagonal:
		rotationDegrees = 270
	default:
		rotationDegrees = 0
	}

	return int(base), rotationDegrees
}
