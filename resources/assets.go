package resources

import "embed"

const (
	PlayerTankRedSmallPath = "images/kenney_topdown_tanks_redux/tank_red_small.png"
	EnemyTankBlueSmallPath = "images/kenney_topdown_tanks_redux/tank_blue_small.png"
	TileGrass1Path         = "images/kenney_topdown_tanks_redux/tile_grass_1.png"
	TileGrass2Path         = "images/kenney_topdown_tanks_redux/tile_grass_2.png"

	ObstacleBarrelGrayTopPath  = "images/kenney_topdown_tanks_redux/obstacle_barrel_gray_top.png"
	ObstacleFenceRedPath       = "images/kenney_topdown_tanks_redux/obstacle_fence_red.png"
	ObstacleCrateWoodPath      = "images/kenney_topdown_tanks_redux/obstacle_crate_wood.png"
	ObstacleTreeGreenLargePath = "images/kenney_topdown_tanks_redux/obstacle_tree_green_large.png"
	ObstacleTreeGreenSmallPath = "images/kenney_topdown_tanks_redux/obstacle_tree_green_small.png"
	ObstacleTreeBrownLargePath = "images/kenney_topdown_tanks_redux/obstacle_tree_brown_large.png"
	ObstacleTreeBrownSmallPath = "images/kenney_topdown_tanks_redux/obstacle_tree_brown_small.png"
	ObstacleWireStraightPath   = "images/kenney_topdown_tanks_redux/obstacle_wire_straight.png"
	ObstacleWireCrookedPath    = "images/kenney_topdown_tanks_redux/obstacle_wire_crooked.png"
	ObstacleTankCrossMetalPath = "images/kenney_topdown_tanks_redux/obstacle_tank_cross_metal.png"

	// Kept for compatibility with previous naming; now points to the gray barrel sprite.
	ObstacleBarrelRedTopPath = ObstacleBarrelGrayTopPath

	ProjectileRedPath  = "images/kenney_topdown_tanks_redux/projectile_red.png"
	ProjectileBluePath = "images/kenney_topdown_tanks_redux/projectile_blue.png"

	ShotRedPath    = "images/kenney_topdown_tanks_redux/shot_red.png"
	ShotOrangePath = "images/kenney_topdown_tanks_redux/shot_orange.png"

	Explosion3Path      = "images/kenney_topdown_tanks_redux/explosion_3.png"
	ExplosionSmoke3Path = "images/kenney_topdown_tanks_redux/explosion_smoke_3.png"
	Explosion1Path      = "images/kenney_topdown_tanks_redux/explosion_1.png"
	ExplosionSmoke1Path = "images/kenney_topdown_tanks_redux/explosion_smoke_1.png"
	TankBodyDarkPath    = "images/kenney_topdown_tanks_redux/tank_body_dark.png"
	SmokeCirclesPath    = "images/kenney_topdown_tanks_redux/smoke_circles.png"
	TracksSmallPath     = "images/kenney_topdown_tanks_redux/tracks_small.png"
	GroundSplashPath    = "images/kenney_topdown_tanks_redux/ground_splash.png"

	CannonFireSFXPath = "audio/opengameart_top_down_shooter/cannon_fire.ogg"
	EngineLoopSFXPath = "audio/opengameart_top_down_shooter/engine_heavy_loop.ogg"
)

// FS contains embedded runtime assets needed by the executable.
//
//go:embed images/kenney_topdown_tanks_redux/*.png audio/opengameart_top_down_shooter/*.ogg
var FS embed.FS
