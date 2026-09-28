package resources

import "embed"

const PlayerTankRedSmallPath = "images/kenney_topdown_tanks_redux/tank_red_small.png"

// FS contains embedded runtime assets needed by the executable.
//
//go:embed images/kenney_topdown_tanks_redux/tank_red_small.png
var FS embed.FS
