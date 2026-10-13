package render

import "image"

// effectVisible includes partially visible effects near the viewport edges.
func effectVisible(view image.Rectangle, x, y, radius float64) bool {
	return radius > 0 && x+radius > float64(view.Min.X) && x-radius < float64(view.Max.X) && y+radius > float64(view.Min.Y) && y-radius < float64(view.Max.Y)
}
