package gamepalette

import gameassets "github.com/co0p/tankism-revival/resources"

type Layer string

const (
	LayerGround    Layer = "ground"
	LayerProps     Layer = "props"
	LayerPropsTall Layer = "props_tall"
)

const (
	SpriteTankRedSmall  = "tank_red_small"
	SpriteTankBlueSmall = "tank_blue_small"
	SpriteTileGrass1    = "tile_grass_1"
	SpriteTileGrass2    = "tile_grass_2"

	SpriteObstacleBarrelGrayTop  = "obstacle_barrel_gray_top"
	SpriteObstacleFenceRed       = "obstacle_fence_red"
	SpriteObstacleCrateWood      = "obstacle_crate_wood"
	SpriteObstacleTreeGreenLarge = "obstacle_tree_green_large"
	SpriteObstacleTreeGreenSmall = "obstacle_tree_green_small"
	SpriteObstacleTreeBrownLarge = "obstacle_tree_brown_large"
	SpriteObstacleTreeBrownSmall = "obstacle_tree_brown_small"
	SpriteObstacleWireStraight   = "obstacle_wire_straight"
	SpriteObstacleWireCrooked    = "obstacle_wire_crooked"
	SpriteObstacleTankCrossMetal = "obstacle_tank_cross_metal"

	// Kept for compatibility with previous naming.
	SpriteObstacleBarrelRedTop = SpriteObstacleBarrelGrayTop

	SpriteProjectileRed  = "projectile_red"
	SpriteProjectileBlue = "projectile_blue"

	SpriteShotRed    = "shot_red"
	SpriteShotOrange = "shot_orange"
)

type Sprite struct {
	ID    string
	GID   int
	Path  string
	Layer Layer
}

const (
	GIDTileGrass1 = 1
	GIDTileGrass2 = 2

	GIDObstacleBarrelGrayTop = 3
	GIDObstacleFenceRed      = 4
	GIDObstacleCrateWood     = 5

	GIDTankRedSmall  = 6
	GIDTankBlueSmall = 7

	GIDProjectileRed  = 8
	GIDProjectileBlue = 9

	GIDShotRed    = 10
	GIDShotOrange = 11

	GIDObstacleTreeGreenLarge = 12
	GIDObstacleTreeGreenSmall = 13
	GIDObstacleTreeBrownLarge = 14
	GIDObstacleTreeBrownSmall = 15
	GIDObstacleWireStraight   = 16
	GIDObstacleWireCrooked    = 17
	GIDObstacleTankCrossMetal = 18
)

var Sprites = []Sprite{
	{ID: SpriteTileGrass1, GID: GIDTileGrass1, Path: gameassets.TileGrass1Path, Layer: LayerGround},
	{ID: SpriteTileGrass2, GID: GIDTileGrass2, Path: gameassets.TileGrass2Path, Layer: LayerGround},

	{ID: SpriteObstacleBarrelGrayTop, GID: GIDObstacleBarrelGrayTop, Path: gameassets.ObstacleBarrelGrayTopPath, Layer: LayerProps},
	{ID: SpriteObstacleFenceRed, GID: GIDObstacleFenceRed, Path: gameassets.ObstacleFenceRedPath, Layer: LayerProps},
	{ID: SpriteObstacleCrateWood, GID: GIDObstacleCrateWood, Path: gameassets.ObstacleCrateWoodPath, Layer: LayerProps},

	{ID: SpriteTankRedSmall, GID: GIDTankRedSmall, Path: gameassets.PlayerTankRedSmallPath, Layer: LayerProps},
	{ID: SpriteTankBlueSmall, GID: GIDTankBlueSmall, Path: gameassets.EnemyTankBlueSmallPath, Layer: LayerProps},
	{ID: SpriteProjectileRed, GID: GIDProjectileRed, Path: gameassets.ProjectileRedPath, Layer: LayerProps},
	{ID: SpriteProjectileBlue, GID: GIDProjectileBlue, Path: gameassets.ProjectileBluePath, Layer: LayerProps},
	{ID: SpriteShotRed, GID: GIDShotRed, Path: gameassets.ShotRedPath, Layer: LayerProps},
	{ID: SpriteShotOrange, GID: GIDShotOrange, Path: gameassets.ShotOrangePath, Layer: LayerProps},

	{ID: SpriteObstacleTreeGreenLarge, GID: GIDObstacleTreeGreenLarge, Path: gameassets.ObstacleTreeGreenLargePath, Layer: LayerPropsTall},
	{ID: SpriteObstacleTreeGreenSmall, GID: GIDObstacleTreeGreenSmall, Path: gameassets.ObstacleTreeGreenSmallPath, Layer: LayerPropsTall},
	{ID: SpriteObstacleTreeBrownLarge, GID: GIDObstacleTreeBrownLarge, Path: gameassets.ObstacleTreeBrownLargePath, Layer: LayerPropsTall},
	{ID: SpriteObstacleTreeBrownSmall, GID: GIDObstacleTreeBrownSmall, Path: gameassets.ObstacleTreeBrownSmallPath, Layer: LayerPropsTall},
	{ID: SpriteObstacleWireStraight, GID: GIDObstacleWireStraight, Path: gameassets.ObstacleWireStraightPath, Layer: LayerProps},
	{ID: SpriteObstacleWireCrooked, GID: GIDObstacleWireCrooked, Path: gameassets.ObstacleWireCrookedPath, Layer: LayerProps},
	{ID: SpriteObstacleTankCrossMetal, GID: GIDObstacleTankCrossMetal, Path: gameassets.ObstacleTankCrossMetalPath, Layer: LayerProps},
}

var spriteByGID = func() map[int]Sprite {
	indexed := make(map[int]Sprite, len(Sprites))
	for _, sprite := range Sprites {
		indexed[sprite.GID] = sprite
	}

	return indexed
}()

func MaxGID() int {
	max := 0
	for _, sprite := range Sprites {
		if sprite.GID > max {
			max = sprite.GID
		}
	}

	return max
}

func SpriteByGID(gid int) (Sprite, bool) {
	sprite, ok := spriteByGID[gid]
	return sprite, ok
}
