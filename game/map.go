package game

type Map struct {
	Width  float64
	Height float64
}

func DefaultMap() Map {
	return Map{
		Width:  2400,
		Height: 1800,
	}
}

func (m Map) ClampPosition(x, y float64) (float64, float64) {
	return clamp(x, 0, m.Width), clamp(y, 0, m.Height)
}

func (m Map) IsLargerThanViewport(viewWidth, viewHeight float64) bool {
	return m.Width > viewWidth && m.Height > viewHeight
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}
