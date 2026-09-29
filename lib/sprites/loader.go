package sprites

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Asset struct {
	ID   string
	Path string
}

func Load(assetsFS fs.FS, assets []Asset) (map[string]*ebiten.Image, error) {
	loaded := make(map[string]*ebiten.Image, len(assets))

	for _, asset := range assets {
		sprite, _, err := ebitenutil.NewImageFromFileSystem(assetsFS, asset.Path)
		if err != nil {
			return nil, fmt.Errorf("load sprite %q from %q: %w", asset.ID, asset.Path, err)
		}
		loaded[asset.ID] = sprite
	}

	return loaded, nil
}
