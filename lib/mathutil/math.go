package mathutil

import "math"

func Clamp(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}

	return value
}

func Clamp01(value float64) float64 {
	return Clamp(value, 0, 1)
}

func ClampUnit(value float64) float64 {
	return Clamp(value, -1, 1)
}

func ClampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}

	return value
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

func MoveTowards(current, target, maxDelta float64) float64 {
	if maxDelta <= 0 {
		return target
	}

	if target > current {
		return math.Min(current+maxDelta, target)
	}

	if target < current {
		return math.Max(current-maxDelta, target)
	}

	return target
}

func NormalizeRadians(radians float64) float64 {
	wrapped := math.Mod(radians, 2*math.Pi)
	if wrapped < 0 {
		wrapped += 2 * math.Pi
	}

	return wrapped
}
