package fonts

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// FontAssets provides access to embedded font files
type FontAssets interface {
	ReadFile(name string) ([]byte, error)
}

// Loader loads and caches font faces
type Loader struct {
	assets FontAssets
	faces  map[string]*text.GoTextFace
}

// NewLoader creates a font loader with the given asset source
func NewLoader(assets FontAssets) *Loader {
	return &Loader{
		assets: assets,
		faces:  make(map[string]*text.GoTextFace),
	}
}

// LoadFace loads a TrueType font and returns a text face at the specified size
func (l *Loader) LoadFace(fontPath string, size float64) (*text.GoTextFace, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("%s@%.1f", fontPath, size)
	if face, ok := l.faces[cacheKey]; ok {
		return face, nil
	}

	// Load font file
	data, err := l.assets.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("read font file %q: %w", fontPath, err)
	}

	// Create a GoTextFaceSource from the font data
	source, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create font source from %q: %w", fontPath, err)
	}

	// Create a face at the specified size
	face := &text.GoTextFace{
		Source: source,
		Size:   size,
	}

	// Cache the face
	l.faces[cacheKey] = face

	return face, nil
}

// LoadKenneyBlocksFace loads the Kenney Blocks font at the specified size
func (l *Loader) LoadKenneyBlocksFace(size float64) (*text.GoTextFace, error) {
	return l.LoadFace("fonts/Kenney Blocks.ttf", size)
}

// LoadKenneyFutureFace loads the Kenney Future font at the specified size
func (l *Loader) LoadKenneyFutureFace(size float64) (*text.GoTextFace, error) {
	return l.LoadFace("fonts/Kenney Future.ttf", size)
}
