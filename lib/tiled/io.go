package tiled

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func LoadJSON(r io.Reader) (Map, error) {
	var m Map
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return Map{}, fmt.Errorf("decode tiled map JSON: %w", err)
	}

	if err := m.EnsureEditorLayers(); err != nil {
		return Map{}, err
	}

	return m, nil
}

func SaveJSON(w io.Writer, m Map) error {
	if err := m.EnsureEditorLayers(); err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m); err != nil {
		return fmt.Errorf("encode tiled map JSON: %w", err)
	}

	return nil
}

func ReadFile(path string) (Map, error) {
	file, err := os.Open(path)
	if err != nil {
		return Map{}, fmt.Errorf("open map file %q: %w", path, err)
	}
	defer file.Close()

	return LoadJSON(file)
}

func WriteFile(path string, m Map) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create map directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create map file %q: %w", path, err)
	}
	defer file.Close()

	if err := SaveJSON(file, m); err != nil {
		return fmt.Errorf("write map file %q: %w", path, err)
	}

	return nil
}
