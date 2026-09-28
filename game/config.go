package game

const (
	LogicalWidth  = 960
	LogicalHeight = 540
)

func LogicalAspectRatio() float64 {
	return float64(LogicalWidth) / float64(LogicalHeight)
}
