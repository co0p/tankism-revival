package shaders

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

// Load compiles and returns a Kage shader from the given filesystem path.
func Load(fsys fs.FS, path string) (*ebiten.Shader, error) {
	src, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("shaders: read %q: %w", path, err)
	}
	shader, err := ebiten.NewShader(src)
	if err != nil {
		return nil, fmt.Errorf("shaders: compile %q: %w", path, err)
	}
	return shader, nil
}
