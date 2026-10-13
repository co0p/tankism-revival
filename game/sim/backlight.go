package sim

import "math"

func TankBackLightPosition(x, y, rotationRadians, distance float64) (float64, float64) {
	forwardX := math.Sin(rotationRadians)
	forwardY := -math.Cos(rotationRadians)
	return x - forwardX*distance, y - forwardY*distance
}
