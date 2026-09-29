package game

import "github.com/co0p/tankism-revival/lib/mathutil"

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
	return mathutil.Clamp(x, 0, m.Width), mathutil.Clamp(y, 0, m.Height)
}

func (m Map) IsLargerThanViewport(viewWidth, viewHeight float64) bool {
	return m.Width > viewWidth && m.Height > viewHeight
}
