package sprites

import (
	"fmt"
	"image"
	"image/png"
	"io/fs"
)

// ExtractCircles separates disconnected shapes in a sprite. Alpha below 16/255
// is ignored: the source has nearly invisible export residue between circles.
// The scan order is stable so particle variant IDs are repeatable.
func ExtractCircles(source image.Image) []*image.NRGBA {
	bounds := source.Bounds()
	seen := make(map[image.Point]bool)
	var circles []*image.NRGBA
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			start := image.Pt(x, y)
			_, _, _, alpha := source.At(x, y).RGBA()
			if seen[start] || alpha < 16*257 {
				continue
			}
			queue := []image.Point{start}
			seen[start] = true
			box := image.Rect(x, y, x+1, y+1)
			for i := 0; i < len(queue); i++ {
				p := queue[i]
				box = box.Union(image.Rect(p.X, p.Y, p.X+1, p.Y+1))
				for _, offset := range []image.Point{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}} {
					next := p.Add(offset)
					if !next.In(bounds) || seen[next] {
						continue
					}
					_, _, _, a := source.At(next.X, next.Y).RGBA()
					if a >= 16*257 {
						seen[next] = true
						queue = append(queue, next)
					}
				}
			}
			circle := image.NewNRGBA(image.Rect(0, 0, box.Dx(), box.Dy()))
			for _, p := range queue {
				circle.Set(p.X-box.Min.X, p.Y-box.Min.Y, source.At(p.X, p.Y))
			}
			circles = append(circles, circle)
		}
	}
	return circles
}

func LoadCircles(assetsFS fs.FS, path string) ([]*image.NRGBA, error) {
	file, err := assetsFS.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open smoke circles: %w", err)
	}
	defer file.Close()
	source, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode smoke circles: %w", err)
	}
	circles := ExtractCircles(source)
	if len(circles) == 0 {
		return nil, fmt.Errorf("smoke sprite %q contains no circles", path)
	}
	return circles, nil
}
