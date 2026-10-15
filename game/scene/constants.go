package scene

import "image/color"

// Context keys for scene initialization
const (
	FontLoaderContextKey      = "fontLoader"
	BlurShaderContextKey      = "blurShader"
	BackgroundImageContextKey = "backgroundImage"
)

// UI Constants for scene rendering

// Icon characters
const (
	SelectedIcon = ">" // Icon for selected menu items
	SpacingIcon  = " " // Spacing for unselected items
	LoadingIcon  = "★" // Loading indicator icon
)

// Font sizes (placeholders for future Kenney Blocks font)
const (
	FontSizeSmall  = 12
	FontSizeMedium = 16
	FontSizeLarge  = 24
	FontSizeHuge   = 32
)

// Color palette for UI
var (
	// Primary colors
	ColorGreen  = color.RGBA{R: 0x4a, G: 0xd6, B: 0x3a, A: 0xff} // Bright green for selection/success
	ColorRed    = color.RGBA{R: 0xff, G: 0x40, B: 0x40, A: 0xff} // Red for defeat/error
	ColorGray   = color.RGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff} // Gray for unselected/disabled
	ColorWhite  = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff} // White for titles
	ColorDarkBg = color.RGBA{R: 0x1c, G: 0x27, B: 0x1b, A: 0xff} // Dark game background

	// Semantic colors
	ColorSelected   = ColorGreen
	ColorUnselected = ColorGray
	ColorVictory    = ColorGreen
	ColorDefeat     = ColorRed
	ColorTitle      = ColorWhite

	// UI overlay colors
	ColorOverlayLight = color.RGBA{R: 0, G: 0, B: 0, A: 128} // Semi-transparent overlay
	ColorOverlayDark  = color.RGBA{R: 0, G: 0, B: 0, A: 192} // Darker overlay
)
