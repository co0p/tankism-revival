package worldseed

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"github.com/co0p/tankism-revival/lib/gamepalette"
	"github.com/co0p/tankism-revival/lib/mathutil"
	"github.com/co0p/tankism-revival/lib/tiled"
)

type TiledSeedConfig struct {
	Path string

	GroundZIndex    int
	PropsZIndex     int
	PropsTallZIndex int

	ColliderRadiusBySprite map[string]float64

	DefaultEnemyOffsetX float64
	DefaultEnemyOffsetY float64
}

type TiledSeedResult struct {
	MapWidth  float64
	MapHeight float64

	PlayerSpawnX float64
	PlayerSpawnY float64
	EnemySpawnX  float64
	EnemySpawnY  float64
}

func SeedFromTiledMap(world *ecs.World, cfg TiledSeedConfig) (TiledSeedResult, error) {
	if cfg.PropsTallZIndex == 0 {
		cfg.PropsTallZIndex = cfg.PropsZIndex
	}

	doc, err := tiled.ReadFile(cfg.Path)
	if err != nil {
		return TiledSeedResult{}, err
	}

	tileWidth := float64(doc.TileWidth)
	tileHeight := float64(doc.TileHeight)
	mapWidth := float64(doc.Width) * tileWidth
	mapHeight := float64(doc.Height) * tileHeight

	if err := seedGroundTiles(world, doc, cfg.GroundZIndex); err != nil {
		return TiledSeedResult{}, err
	}
	if err := seedObstacleTiles(world, doc, cfg.PropsZIndex, cfg.PropsTallZIndex, cfg.ColliderRadiusBySprite); err != nil {
		return TiledSeedResult{}, err
	}
	if err := seedPropsObjects(world, doc, tiled.PropsLayerName, cfg.PropsZIndex, cfg.PropsTallZIndex, cfg.ColliderRadiusBySprite); err != nil {
		return TiledSeedResult{}, err
	}
	if err := seedPropsObjects(world, doc, tiled.PropsTallLayerName, cfg.PropsTallZIndex, cfg.PropsTallZIndex, cfg.ColliderRadiusBySprite); err != nil {
		return TiledSeedResult{}, err
	}

	playerSpawnX := mapWidth / 2
	playerSpawnY := mapHeight / 2

	if spawn, ok := doc.Spawn(tiled.PlayerSpawnObjectName); ok {
		playerSpawnX = mathutil.Clamp(spawn.X, 0, mapWidth)
		playerSpawnY = mathutil.Clamp(spawn.Y, 0, mapHeight)
	}

	enemySpawnX := mathutil.Clamp(playerSpawnX+cfg.DefaultEnemyOffsetX, 0, mapWidth)
	enemySpawnY := mathutil.Clamp(playerSpawnY+cfg.DefaultEnemyOffsetY, 0, mapHeight)
	if spawn, ok := doc.Spawn(tiled.EnemySpawnObjectName); ok {
		enemySpawnX = mathutil.Clamp(spawn.X, 0, mapWidth)
		enemySpawnY = mathutil.Clamp(spawn.Y, 0, mapHeight)
	}

	return TiledSeedResult{
		MapWidth:     mapWidth,
		MapHeight:    mapHeight,
		PlayerSpawnX: playerSpawnX,
		PlayerSpawnY: playerSpawnY,
		EnemySpawnX:  enemySpawnX,
		EnemySpawnY:  enemySpawnY,
	}, nil
}

func seedGroundTiles(world *ecs.World, doc tiled.Map, zIndex int) error {
	for row := 0; row < doc.Height; row++ {
		for col := 0; col < doc.Width; col++ {
			rawGID, ok := doc.TileGIDAtLayer(tiled.GroundLayerName, col, row)
			if !ok || rawGID == 0 {
				continue
			}

			baseGID, rotationDeg := tiled.DecodeRotatedGID(rawGID)
			spriteDef, ok := gamepalette.SpriteByGID(baseGID)
			if !ok {
				continue
			}

			e := world.CreateEntity()
			x := (float64(col) + 0.5) * float64(doc.TileWidth)
			y := (float64(row) + 0.5) * float64(doc.TileHeight)

			if err := world.SetTransform(e, ecscomponents.Transform{X: x, Y: y, RotationRadians: 0}); err != nil {
				return fmt.Errorf("set ground tile transform: %w", err)
			}
			if err := world.SetSpriteRenderer(e, ecscomponents.SpriteRenderer{
				SpriteID:              spriteDef.ID,
				ZIndex:                zIndex,
				RotationOffsetRadians: float64(rotationDeg) * (math.Pi / 180),
			}); err != nil {
				return fmt.Errorf("set ground tile sprite: %w", err)
			}
		}
	}

	return nil
}

func seedObstacleTiles(world *ecs.World, doc tiled.Map, zIndex int, tallPropsZIndex int, colliderRadiusBySprite map[string]float64) error {
	for row := 0; row < doc.Height; row++ {
		for col := 0; col < doc.Width; col++ {
			rawGID, ok := doc.TileGIDAtLayer(tiled.ObstacleLayerName, col, row)
			if !ok || rawGID == 0 {
				continue
			}

			baseGID, rotationDeg := tiled.DecodeRotatedGID(rawGID)
			spriteDef, ok := gamepalette.SpriteByGID(baseGID)
			if !ok {
				continue
			}

			x := (float64(col) + 0.5) * float64(doc.TileWidth)
			y := (float64(row) + 0.5) * float64(doc.TileHeight)
			spriteZIndex := zIndexForSpriteLayer(spriteDef.Layer, zIndex, tallPropsZIndex)
			if err := spawnPropEntity(world, spriteDef.ID, spriteZIndex, x, y, float64(rotationDeg), colliderRadiusBySprite); err != nil {
				return err
			}
		}
	}

	return nil
}

func seedPropsObjects(world *ecs.World, doc tiled.Map, layerName string, zIndex int, tallPropsZIndex int, colliderRadiusBySprite map[string]float64) error {
	for _, object := range doc.PropsObjectsOnLayer(layerName) {
		baseGID, tiledRotation := tiled.DecodeRotatedGID(object.GID)
		if baseGID <= 0 {
			continue
		}

		spriteDef, ok := gamepalette.SpriteByGID(baseGID)
		if !ok {
			continue
		}

		rotation := object.Rotation + float64(tiledRotation)
		spriteZIndex := zIndexForSpriteLayer(spriteDef.Layer, zIndex, tallPropsZIndex)
		if err := spawnPropEntity(world, spriteDef.ID, spriteZIndex, object.X, object.Y, rotation, colliderRadiusBySprite); err != nil {
			return err
		}
	}

	return nil
}

func zIndexForSpriteLayer(layer gamepalette.Layer, defaultZIndex, tallPropsZIndex int) int {
	if layer == gamepalette.LayerPropsTall {
		return tallPropsZIndex
	}

	return defaultZIndex
}

func spawnPropEntity(world *ecs.World, spriteID string, zIndex int, x, y, rotationDeg float64, colliderRadiusBySprite map[string]float64) error {
	e := world.CreateEntity()
	if err := world.SetTransform(e, ecscomponents.Transform{X: x, Y: y, RotationRadians: rotationDeg * (math.Pi / 180)}); err != nil {
		return fmt.Errorf("set prop transform: %w", err)
	}
	if err := world.SetSpriteRenderer(e, ecscomponents.SpriteRenderer{
		SpriteID:              spriteID,
		ZIndex:                zIndex,
		RotationOffsetRadians: 0,
	}); err != nil {
		return fmt.Errorf("set prop sprite: %w", err)
	}

	radius, ok := colliderRadiusBySprite[spriteID]
	if !ok {
		return nil
	}

	if err := world.SetCircleCollider(e, ecscomponents.CircleCollider{Radius: radius}); err != nil {
		return fmt.Errorf("set prop collider: %w", err)
	}
	if err := world.SetSolidObstacle(e); err != nil {
		return fmt.Errorf("set prop solid obstacle: %w", err)
	}

	return nil
}
