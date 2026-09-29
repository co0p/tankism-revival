package game

type Camera struct {
	X          float64
	Y          float64
	ViewWidth  float64
	ViewHeight float64
}

func NewCamera(viewWidth, viewHeight float64) Camera {
	return Camera{
		ViewWidth:  viewWidth,
		ViewHeight: viewHeight,
	}
}

func (c *Camera) FollowTarget(mapWidth, mapHeight, targetX, targetY float64) {
	wantedX := targetX - c.ViewWidth/2
	wantedY := targetY - c.ViewHeight/2

	maxX := mapWidth - c.ViewWidth
	maxY := mapHeight - c.ViewHeight

	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}

	c.X = clamp(wantedX, 0, maxX)
	c.Y = clamp(wantedY, 0, maxY)
}

func (c Camera) WorldToScreen(worldX, worldY float64) (float64, float64) {
	return worldX - c.X, worldY - c.Y
}
