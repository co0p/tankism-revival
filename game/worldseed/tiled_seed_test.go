package worldseed

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/co0p/tankism-revival/lib/tiled"
)

func TestSeedFromTiledMapSeedsLayersAndSpawns(t *testing.T) {
	tempMapPath := filepath.Join(t.TempDir(), "seed_map.tmj")

	doc := tiled.NewOrthogonalMap(4, 3, 32, 32)
	if err := doc.SetTileGIDOnLayer(tiled.GroundLayerName, 1, 1, tiled.EncodeRotatedGID(1, 90)); err != nil {
		t.Fatalf("set ground gid: %v", err)
	}
	if err := doc.SetTileGIDOnLayer(tiled.ObstacleLayerName, 2, 1, 3); err != nil {
		t.Fatalf("set obstacle gid: %v", err)
	}
	if err := doc.AddPropsObject(5, 100, 70, 28, 28, 45); err != nil {
		t.Fatalf("add props object: %v", err)
	}
	if err := doc.SetSpawn(tiled.PlayerSpawnObjectName, 64, 80); err != nil {
		t.Fatalf("set player spawn: %v", err)
	}
	if err := doc.SetSpawn(tiled.EnemySpawnObjectName, 112, 24); err != nil {
		t.Fatalf("set enemy spawn: %v", err)
	}

	if err := tiled.WriteFile(tempMapPath, doc); err != nil {
		t.Fatalf("write map file: %v", err)
	}

	world := ecs.NewWorld()
	result, err := SeedFromTiledMap(world, TiledSeedConfig{
		Path:         tempMapPath,
		GroundZIndex: 10,
		PropsZIndex:  20,
		ColliderRadiusBySprite: map[string]float64{
			gamepalette.SpriteObstacleBarrelGrayTop: 11,
			gamepalette.SpriteObstacleCrateWood:     13,
		},
		DefaultEnemyOffsetX: 50,
		DefaultEnemyOffsetY: -40,
	})
	if err != nil {
		t.Fatalf("seed from tiled map: %v", err)
	}

	if result.MapWidth != 128 || result.MapHeight != 96 {
		t.Fatalf("unexpected map size result: %.0fx%.0f", result.MapWidth, result.MapHeight)
	}
	if result.PlayerSpawnX != 64 || result.PlayerSpawnY != 80 {
		t.Fatalf("unexpected player spawn: (%f,%f)", result.PlayerSpawnX, result.PlayerSpawnY)
	}
	if result.EnemySpawnX != 112 || result.EnemySpawnY != 24 {
		t.Fatalf("unexpected enemy spawn: (%f,%f)", result.EnemySpawnX, result.EnemySpawnY)
	}

	renderables := world.Renderables()
	if len(renderables) != 3 {
		t.Fatalf("expected 3 renderables from map layers, got %d", len(renderables))
	}

	groundFound := false
	rotatedPropFound := false
	for _, r := range renderables {
		if r.Sprite.SpriteID == gamepalette.SpriteTileGrass1 {
			groundFound = true
			if math.Abs(r.Transform.X-48) > 0.001 || math.Abs(r.Transform.Y-48) > 0.001 {
				t.Fatalf("unexpected ground transform: (%f,%f)", r.Transform.X, r.Transform.Y)
			}
			if math.Abs(r.Sprite.RotationOffsetRadians-math.Pi/2) > 0.001 {
				t.Fatalf("expected ground 90° rotation offset, got %f", r.Sprite.RotationOffsetRadians)
			}
		}

		if r.Sprite.SpriteID == gamepalette.SpriteObstacleCrateWood {
			if math.Abs(r.Transform.X-100) < 0.001 && math.Abs(r.Transform.Y-70) < 0.001 {
				rotatedPropFound = true
				if math.Abs(r.Transform.RotationRadians-math.Pi/4) > 0.001 {
					t.Fatalf("expected prop rotation 45°, got %f", r.Transform.RotationRadians)
				}
			}
		}
	}

	if !groundFound {
		t.Fatalf("expected seeded ground tile renderable")
	}
	if !rotatedPropFound {
		t.Fatalf("expected seeded rotated props object")
	}

	solidBodies := world.SolidObstacleColliderBodies()
	if len(solidBodies) != 2 {
		t.Fatalf("expected 2 solid obstacle colliders, got %d", len(solidBodies))
	}
}

func TestSeedFromTiledMapClampsSpawnsToBounds(t *testing.T) {
	tempMapPath := filepath.Join(t.TempDir(), "seed_map_clamp.tmj")

	doc := tiled.NewOrthogonalMap(2, 2, 32, 32)
	if err := doc.SetSpawn(tiled.PlayerSpawnObjectName, -20, 200); err != nil {
		t.Fatalf("set player spawn: %v", err)
	}
	if err := doc.SetSpawn(tiled.EnemySpawnObjectName, 1000, -50); err != nil {
		t.Fatalf("set enemy spawn: %v", err)
	}
	if err := tiled.WriteFile(tempMapPath, doc); err != nil {
		t.Fatalf("write map file: %v", err)
	}

	world := ecs.NewWorld()
	result, err := SeedFromTiledMap(world, TiledSeedConfig{Path: tempMapPath})
	if err != nil {
		t.Fatalf("seed from tiled map: %v", err)
	}

	if result.MapWidth != 64 || result.MapHeight != 64 {
		t.Fatalf("unexpected map dimensions: %.0fx%.0f", result.MapWidth, result.MapHeight)
	}
	if result.PlayerSpawnX != 0 || result.PlayerSpawnY != 64 {
		t.Fatalf("player spawn should be clamped to bounds, got (%f,%f)", result.PlayerSpawnX, result.PlayerSpawnY)
	}
	if result.EnemySpawnX != 64 || result.EnemySpawnY != 0 {
		t.Fatalf("enemy spawn should be clamped to bounds, got (%f,%f)", result.EnemySpawnX, result.EnemySpawnY)
	}
}
