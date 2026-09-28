package resources

import "embed"

const (
	PlayerTankRedSmallPath = "images/kenney_topdown_tanks_redux/tank_red_small.png"
	EnemyTankBlueSmallPath = "images/kenney_topdown_tanks_redux/tank_blue_small.png"
	TileGrass1Path         = "images/kenney_topdown_tanks_redux/tile_grass_1.png"
	TileGrass2Path         = "images/kenney_topdown_tanks_redux/tile_grass_2.png"

	ObstacleBarrelRedTopPath = "images/kenney_topdown_tanks_redux/obstacle_barrel_red_top.png"
	ObstacleFenceRedPath     = "images/kenney_topdown_tanks_redux/obstacle_fence_red.png"
	ObstacleCrateWoodPath    = "images/kenney_topdown_tanks_redux/obstacle_crate_wood.png"
)

// FS contains embedded runtime assets needed by the executable.
//
//go:embed images/kenney_topdown_tanks_redux/*.png
var FS embed.FS
