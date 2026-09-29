package worldseed

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

type GrassConfig struct {
	TileWidth       int
	TileHeight      int
	Seed            int64
	PrimarySpriteID string
	AltSpriteID     string
	ZIndex          int
}

func SeedGrassMap(world *ecs.World, mapWidth, mapHeight float64, cfg GrassConfig) error {
	if cfg.TileWidth <= 0 || cfg.TileHeight <= 0 {
		return fmt.Errorf("invalid tile size %dx%d", cfg.TileWidth, cfg.TileHeight)
	}
	if cfg.PrimarySpriteID == "" || cfg.AltSpriteID == "" {
		return fmt.Errorf("grass sprite IDs are required")
	}

	rng := rand.New(rand.NewSource(cfg.Seed))

	columns := int(math.Ceil(mapWidth / float64(cfg.TileWidth)))
	rows := int(math.Ceil(mapHeight / float64(cfg.TileHeight)))

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			entity := world.CreateEntity()

			worldX := (float64(col) + 0.5) * float64(cfg.TileWidth)
			worldY := (float64(row) + 0.5) * float64(cfg.TileHeight)

			if err := world.SetTransform(entity, ecscomponents.Transform{X: worldX, Y: worldY, RotationRadians: 0}); err != nil {
				return err
			}

			spriteID := cfg.PrimarySpriteID
			if rng.Intn(2) == 1 {
				spriteID = cfg.AltSpriteID
			}

			rotationOffset := float64(rng.Intn(4)) * (math.Pi / 2)
			flipX := rng.Intn(2) == 1
			flipY := rng.Intn(2) == 1

			if err := world.SetSpriteRenderer(entity, ecscomponents.SpriteRenderer{
				SpriteID:              spriteID,
				ZIndex:                cfg.ZIndex,
				RotationOffsetRadians: rotationOffset,
				FlipX:                 flipX,
				FlipY:                 flipY,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

type ObstacleConfig struct {
	Seed              int64
	CountPerType      int
	PlacementAttempts int

	EdgeMargin  float64
	ClearRadius float64

	SpawnX, SpawnY           float64
	EnemySpawnX, EnemySpawnY float64

	Sprites                []string
	ColliderRadiusBySprite map[string]float64
	ZIndex                 int

	BackLightChance          float64
	BackLightRadius          float64
	BackLightIntensity       float64
	BackLightIntensityJitter float64
	BackLightYOffset         float64
	BackLightColorBase       ecscomponents.LightColor
	BackLightColorJitter     ecscomponents.LightColor
}

func SeedObstacles(world *ecs.World, mapWidth, mapHeight float64, cfg ObstacleConfig) error {
	if cfg.CountPerType <= 0 {
		return nil
	}
	if cfg.PlacementAttempts <= 0 {
		cfg.PlacementAttempts = 200
	}
	if len(cfg.Sprites) == 0 {
		return fmt.Errorf("obstacle sprite list cannot be empty")
	}

	rng := rand.New(rand.NewSource(cfg.Seed))
	clearRadiusSquared := cfg.ClearRadius * cfg.ClearRadius

	placeObstacle := func(spriteID string) error {
		radius, ok := cfg.ColliderRadiusBySprite[spriteID]
		if !ok {
			return fmt.Errorf("missing collider radius for sprite %q", spriteID)
		}

		for range cfg.PlacementAttempts {
			x := cfg.EdgeMargin + rng.Float64()*(mapWidth-2*cfg.EdgeMargin)
			y := cfg.EdgeMargin + rng.Float64()*(mapHeight-2*cfg.EdgeMargin)

			dx := x - cfg.SpawnX
			dy := y - cfg.SpawnY
			if dx*dx+dy*dy < clearRadiusSquared {
				continue
			}

			edx := x - cfg.EnemySpawnX
			edy := y - cfg.EnemySpawnY
			if edx*edx+edy*edy < clearRadiusSquared {
				continue
			}

			entity := world.CreateEntity()
			if err := world.SetTransform(entity, ecscomponents.Transform{X: x, Y: y, RotationRadians: 0}); err != nil {
				return err
			}

			rotationOffset := float64(rng.Intn(4)) * (math.Pi / 2)
			flipX := rng.Intn(2) == 1
			flipY := rng.Intn(2) == 1

			if err := world.SetSpriteRenderer(entity, ecscomponents.SpriteRenderer{
				SpriteID:              spriteID,
				ZIndex:                cfg.ZIndex,
				RotationOffsetRadians: rotationOffset,
				FlipX:                 flipX,
				FlipY:                 flipY,
			}); err != nil {
				return err
			}

			if err := world.SetCircleCollider(entity, ecscomponents.CircleCollider{Radius: radius}); err != nil {
				return err
			}

			if err := world.SetSolidObstacle(entity); err != nil {
				return err
			}

			if rng.Float64() < cfg.BackLightChance {
				lightIntensity := cfg.BackLightIntensity + rng.Float64()*cfg.BackLightIntensityJitter
				lightEntity := world.CreateEntity()
				if err := world.SetTransform(lightEntity, ecscomponents.Transform{
					X:               x,
					Y:               y + cfg.BackLightYOffset,
					RotationRadians: 0,
				}); err != nil {
					return err
				}

				if err := world.SetPointLight(lightEntity, ecscomponents.PointLight{
					Radius:    cfg.BackLightRadius,
					Intensity: lightIntensity,
					Color: ecscomponents.LightColor{
						R: cfg.BackLightColorBase.R + rng.Float64()*cfg.BackLightColorJitter.R,
						G: cfg.BackLightColorBase.G + rng.Float64()*cfg.BackLightColorJitter.G,
						B: cfg.BackLightColorBase.B + rng.Float64()*cfg.BackLightColorJitter.B,
					},
				}); err != nil {
					return err
				}
			}

			return nil
		}

		return fmt.Errorf("could not place obstacle %q after %d attempts", spriteID, cfg.PlacementAttempts)
	}

	for range cfg.CountPerType {
		for _, spriteID := range cfg.Sprites {
			if err := placeObstacle(spriteID); err != nil {
				return err
			}
		}
	}

	return nil
}
