package gamepalette

import gameassets "github.com/co0p/tankism-revival/resources"

type Layer string

const (
	LayerGround Layer = "ground"
	LayerProps  Layer = "props"
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
	Path  string
	Layer Layer
}

var Sprites = []Sprite{
	{ID: SpriteTileGrass1, Path: gameassets.TileGrass1Path, Layer: LayerGround},
	{ID: SpriteTileGrass2, Path: gameassets.TileGrass2Path, Layer: LayerGround},

	{ID: SpriteObstacleBarrelGrayTop, Path: gameassets.ObstacleBarrelGrayTopPath, Layer: LayerProps},
	{ID: SpriteObstacleFenceRed, Path: gameassets.ObstacleFenceRedPath, Layer: LayerProps},
	{ID: SpriteObstacleCrateWood, Path: gameassets.ObstacleCrateWoodPath, Layer: LayerProps},

	{ID: SpriteTankRedSmall, Path: gameassets.PlayerTankRedSmallPath, Layer: LayerProps},
	{ID: SpriteTankBlueSmall, Path: gameassets.EnemyTankBlueSmallPath, Layer: LayerProps},
	{ID: SpriteProjectileRed, Path: gameassets.ProjectileRedPath, Layer: LayerProps},
	{ID: SpriteProjectileBlue, Path: gameassets.ProjectileBluePath, Layer: LayerProps},
	{ID: SpriteShotRed, Path: gameassets.ShotRedPath, Layer: LayerProps},
	{ID: SpriteShotOrange, Path: gameassets.ShotOrangePath, Layer: LayerProps},

	{ID: SpriteObstacleTreeGreenLarge, Path: gameassets.ObstacleTreeGreenLargePath, Layer: LayerProps},
	{ID: SpriteObstacleTreeGreenSmall, Path: gameassets.ObstacleTreeGreenSmallPath, Layer: LayerProps},
	{ID: SpriteObstacleTreeBrownLarge, Path: gameassets.ObstacleTreeBrownLargePath, Layer: LayerProps},
	{ID: SpriteObstacleTreeBrownSmall, Path: gameassets.ObstacleTreeBrownSmallPath, Layer: LayerProps},
	{ID: SpriteObstacleWireStraight, Path: gameassets.ObstacleWireStraightPath, Layer: LayerProps},
	{ID: SpriteObstacleWireCrooked, Path: gameassets.ObstacleWireCrookedPath, Layer: LayerProps},
	{ID: SpriteObstacleTankCrossMetal, Path: gameassets.ObstacleTankCrossMetalPath, Layer: LayerProps},
}

func SpriteByGID(gid int) (Sprite, bool) {
	if gid <= 0 || gid > len(Sprites) {
		return Sprite{}, false
	}

	return Sprites[gid-1], true
}
