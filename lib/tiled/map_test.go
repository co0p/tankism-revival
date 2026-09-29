package tiled

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewOrthogonalMapIncludesEditorLayers(t *testing.T) {
	m := NewOrthogonalMap(20, 15, 32, 32)

	if m.Width != 20 || m.Height != 15 {
		t.Fatalf("unexpected map size: %dx%d", m.Width, m.Height)
	}

	if _, ok := m.Spawn(PlayerSpawnObjectName); !ok {
		t.Fatalf("player spawn should exist")
	}
	if _, ok := m.Spawn(EnemySpawnObjectName); !ok {
		t.Fatalf("enemy spawn should exist")
	}

	if gid, ok := m.TileGIDAtLayer(GroundLayerName, 0, 0); !ok || gid != 0 {
		t.Fatalf("ground layer should exist and start empty")
	}
	if gid, ok := m.TileGIDAtLayer(ObstacleLayerName, 0, 0); !ok || gid != 0 {
		t.Fatalf("obstacle layer should exist and start empty")
	}

	if props := m.PropsObjects(); len(props) != 0 {
		t.Fatalf("props layer should start empty")
	}
}

func TestReadWriteRoundTrip(t *testing.T) {
	dir := t.TempDir()
	mapPath := filepath.Join(dir, "test_map.tmj")

	m := NewOrthogonalMap(12, 8, 32, 32)
	if err := m.SetTileGIDOnLayer(GroundLayerName, 2, 3, 5); err != nil {
		t.Fatalf("set ground tile gid: %v", err)
	}
	if err := m.AddPropsObject(7, 120, 80, 20, 20, 90); err != nil {
		t.Fatalf("add props object: %v", err)
	}
	if err := m.SetSpawn(PlayerSpawnObjectName, 160, 224); err != nil {
		t.Fatalf("set player spawn: %v", err)
	}

	if err := WriteFile(mapPath, m); err != nil {
		t.Fatalf("write map: %v", err)
	}

	reloaded, err := ReadFile(mapPath)
	if err != nil {
		t.Fatalf("read map: %v", err)
	}

	gid, ok := reloaded.TileGIDAtLayer(GroundLayerName, 2, 3)
	if !ok || gid != 5 {
		t.Fatalf("expected ground gid=5 after reload, got ok=%t gid=%d", ok, gid)
	}

	props := reloaded.PropsObjects()
	if len(props) != 1 {
		t.Fatalf("expected one props object after reload, got %d", len(props))
	}
	if props[0].GID != 7 {
		t.Fatalf("unexpected props gid: %d", props[0].GID)
	}
	if props[0].Rotation != 90 {
		t.Fatalf("expected props rotation=90, got %f", props[0].Rotation)
	}

	spawn, ok := reloaded.Spawn(PlayerSpawnObjectName)
	if !ok {
		t.Fatalf("expected player spawn after reload")
	}
	if spawn.X != 160 || spawn.Y != 224 {
		t.Fatalf("unexpected player spawn: (%f,%f)", spawn.X, spawn.Y)
	}
}

func TestResizePreservesTileDataOverlap(t *testing.T) {
	m := NewOrthogonalMap(3, 2, 32, 32)
	if err := m.SetTileGIDOnLayer(GroundLayerName, 1, 1, 9); err != nil {
		t.Fatalf("set ground tile: %v", err)
	}

	if err := m.Resize(5, 4); err != nil {
		t.Fatalf("resize larger: %v", err)
	}
	gid, ok := m.TileGIDAtLayer(GroundLayerName, 1, 1)
	if !ok || gid != 9 {
		t.Fatalf("expected preserved gid after enlarge, got ok=%t gid=%d", ok, gid)
	}

	if err := m.Resize(2, 2); err != nil {
		t.Fatalf("resize smaller: %v", err)
	}
	gid, ok = m.TileGIDAtLayer(GroundLayerName, 1, 1)
	if !ok || gid != 9 {
		t.Fatalf("expected preserved overlap gid after shrink, got ok=%t gid=%d", ok, gid)
	}
}

func TestPropsObjectsCanBeAddedAndRemoved(t *testing.T) {
	m := NewOrthogonalMap(10, 10, 32, 32)

	if err := m.AddPropsObject(3, 64, 64, 16, 16, 0); err != nil {
		t.Fatalf("add props 1: %v", err)
	}
	if err := m.AddPropsObject(4, 200, 200, 16, 16, 0); err != nil {
		t.Fatalf("add props 2: %v", err)
	}

	removed, err := m.RemoveNearestPropsObject(70, 70, 20)
	if err != nil {
		t.Fatalf("remove nearest: %v", err)
	}
	if !removed {
		t.Fatalf("expected nearest props object to be removed")
	}

	props := m.PropsObjects()
	if len(props) != 1 {
		t.Fatalf("expected one props object left, got %d", len(props))
	}
	if props[0].GID != 4 {
		t.Fatalf("expected remaining props gid=4, got %d", props[0].GID)
	}
}

func TestLoadJSONAutoAddsEditorLayers(t *testing.T) {
	dir := t.TempDir()
	mapPath := filepath.Join(dir, "minimal.tmj")

	minimal := `{"height":4,"width":6,"tileheight":32,"tilewidth":32,"layers":[],"nextlayerid":1,"nextobjectid":1,"orientation":"orthogonal","renderorder":"right-down","type":"map","version":"1.10","tiledversion":"1.11.0"}`
	if err := os.WriteFile(mapPath, []byte(minimal), 0o644); err != nil {
		t.Fatalf("write minimal map: %v", err)
	}

	loaded, err := ReadFile(mapPath)
	if err != nil {
		t.Fatalf("read minimal map: %v", err)
	}

	if _, ok := loaded.Spawn(PlayerSpawnObjectName); !ok {
		t.Fatalf("expected auto-added player spawn")
	}
	if _, ok := loaded.TileGIDAtLayer(GroundLayerName, 0, 0); !ok {
		t.Fatalf("expected auto-added ground layer")
	}
	if len(loaded.PropsObjects()) != 0 {
		t.Fatalf("expected auto-added props layer to start empty")
	}
}
