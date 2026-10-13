package render

import (
	"image"
	"image/color"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"github.com/co0p/tankism-revival/game/sim"
	"github.com/co0p/tankism-revival/lib/mathutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	defaultLightSpriteSize     = 256
	defaultLightIntensityScale = 0.65

	defaultHealthBarWidth      = 34.0
	defaultHealthBarHeight     = 4.0
	defaultHealthBarYOffset    = 30.0
	defaultHealthBarBorderSize = 1.0

	defaultForegroundMinZIndex = 50
)

type Camera interface {
	WorldToScreen(worldX, worldY float64) (float64, float64)
}

type Config struct {
	GroundSplashSprite *ebiten.Image
	TrackSprite        *ebiten.Image
	SmokeSprites       []*ebiten.Image
	OutsideMapColor    color.RGBA

	LightSpriteSize     int
	LightIntensityScale float64

	HealthBarWidth      float64
	HealthBarHeight     float64
	HealthBarYOffset    float64
	HealthBarBorderSize float64

	ForegroundMinZIndex int
}

type WorldRenderer struct {
	groundSplashSprite *ebiten.Image
	trackRow           *ebiten.Image
	smokeSprites       []*ebiten.Image
	sprites            map[string]*ebiten.Image

	outsideMapColor color.RGBA
	lightSprite     *ebiten.Image

	lightIntensityScale float64
	foregroundMinZIndex int

	healthBarWidth      float64
	healthBarHeight     float64
	healthBarYOffset    float64
	healthBarBorderSize float64
}

func NewWorldRenderer(sprites map[string]*ebiten.Image, cfg Config) *WorldRenderer {
	if cfg.LightSpriteSize <= 0 {
		cfg.LightSpriteSize = defaultLightSpriteSize
	}
	if cfg.LightIntensityScale <= 0 {
		cfg.LightIntensityScale = defaultLightIntensityScale
	}
	if cfg.HealthBarWidth <= 0 {
		cfg.HealthBarWidth = defaultHealthBarWidth
	}
	if cfg.HealthBarHeight <= 0 {
		cfg.HealthBarHeight = defaultHealthBarHeight
	}
	if cfg.HealthBarYOffset == 0 {
		cfg.HealthBarYOffset = defaultHealthBarYOffset
	}
	if cfg.HealthBarBorderSize <= 0 {
		cfg.HealthBarBorderSize = defaultHealthBarBorderSize
	}
	if cfg.ForegroundMinZIndex <= 0 {
		cfg.ForegroundMinZIndex = defaultForegroundMinZIndex
	}

	spriteMap := make(map[string]*ebiten.Image, len(sprites))
	for spriteID, img := range sprites {
		spriteMap[spriteID] = img
	}

	var trackRow *ebiten.Image
	if cfg.TrackSprite != nil {
		bounds := cfg.TrackSprite.Bounds()
		// One paired row from the retina sprite, rather than the full straight strip.
		trackRow = cfg.TrackSprite.SubImage(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+8)).(*ebiten.Image)
	}
	return &WorldRenderer{
		groundSplashSprite:  cfg.GroundSplashSprite,
		trackRow:            trackRow,
		smokeSprites:        cfg.SmokeSprites,
		sprites:             spriteMap,
		outsideMapColor:     cfg.OutsideMapColor,
		lightSprite:         newRadialLightSprite(cfg.LightSpriteSize),
		lightIntensityScale: cfg.LightIntensityScale,
		foregroundMinZIndex: cfg.ForegroundMinZIndex,
		healthBarWidth:      cfg.HealthBarWidth,
		healthBarHeight:     cfg.HealthBarHeight,
		healthBarYOffset:    cfg.HealthBarYOffset,
		healthBarBorderSize: cfg.HealthBarBorderSize,
	}
}

func (r *WorldRenderer) Draw(screen *ebiten.Image, world *ecs.World, camera Camera) {
	r.DrawWithTracks(screen, world, camera, nil)
}

func (r *WorldRenderer) DrawWithTracks(screen *ebiten.Image, world *ecs.World, camera Camera, marks []sim.TrackMark) {
	r.DrawWithGroundEffects(screen, world, camera, marks, nil)
}

func (r *WorldRenderer) DrawWithGroundEffects(screen *ebiten.Image, world *ecs.World, camera Camera, marks []sim.TrackMark, splashes []sim.GroundSplash) {
	screen.Fill(r.outsideMapColor)

	r.drawRenderables(screen, world, camera, func(renderable ecs.Renderable) bool { return renderable.Sprite.ZIndex < r.foregroundMinZIndex })
	r.drawTracks(screen, marks, camera)
	r.drawGroundSplashes(screen, splashes, camera)
	r.drawRenderables(screen, world, camera, func(renderable ecs.Renderable) bool { return renderable.Sprite.ZIndex >= r.foregroundMinZIndex })
	r.drawLighting(screen, world, camera)
	r.drawRenderables(screen, world, camera, func(renderable ecs.Renderable) bool {
		return renderable.Sprite.ZIndex >= r.foregroundMinZIndex
	})
	r.drawHitRings(screen, world, camera)
	r.drawTankHealthBars(screen, world, camera)
}

func (r *WorldRenderer) drawGroundSplashes(screen *ebiten.Image, splashes []sim.GroundSplash, camera Camera) {
	if r.groundSplashSprite == nil {
		return
	}
	bounds := r.groundSplashSprite.Bounds()
	for _, splash := range splashes {
		x, y := camera.WorldToScreen(splash.X, splash.Y)
		// Rotation can extend the corners beyond the sprite's half-width.
		if !effectVisible(screen.Bounds(), x, y, sim.GroundSplashSize*math.Sqrt2/2) {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		op.GeoM.Scale(sim.GroundSplashSize/float64(bounds.Dx()), sim.GroundSplashSize/float64(bounds.Dy()))
		op.GeoM.Rotate(splash.RotationRadians)
		op.GeoM.Translate(x, y)
		op.ColorScale.ScaleAlpha(float32(0.8 * (1 - mathutil.Clamp01(splash.Age/sim.GroundSplashLifetime))))
		screen.DrawImage(r.groundSplashSprite, op)
	}
}

func (r *WorldRenderer) DrawDust(screen *ebiten.Image, particles []sim.DustParticle, camera Camera) {
	for _, p := range particles {
		if p.Variant < 0 || p.Variant >= len(r.smokeSprites) || p.Lifetime <= 0 || p.Age >= p.Lifetime {
			continue
		}
		progress := mathutil.Clamp01(p.Age / p.Lifetime)
		size := 4 + progress*12
		x, y := camera.WorldToScreen(p.X, p.Y)
		if !effectVisible(screen.Bounds(), x, y, size/2) {
			continue
		}
		img := r.smokeSprites[p.Variant]
		bounds := img.Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		op.GeoM.Scale(size/float64(bounds.Dx()), size/float64(bounds.Dy()))
		op.GeoM.Translate(x, y)
		alpha := 0.65 * (1 - progress)
		op.ColorScale.Scale(float32(alpha), float32(alpha*0.85), float32(alpha*0.6), float32(alpha))
		screen.DrawImage(img, op)
	}
}

func (r *WorldRenderer) drawTracks(screen *ebiten.Image, marks []sim.TrackMark, camera Camera) {
	if r.trackRow == nil {
		return
	}
	bounds := r.trackRow.Bounds()
	for _, mark := range marks {
		remaining := 1 - mathutil.Clamp01(mark.Age/sim.TrackLifetime)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		op.GeoM.Scale(30/float64(bounds.Dx()), 0.5)
		op.GeoM.Rotate(mark.RotationRadians)
		x, y := camera.WorldToScreen(mark.X, mark.Y)
		op.GeoM.Translate(x, y)
		op.ColorScale.ScaleAlpha(float32(0.55 * remaining))
		screen.DrawImage(r.trackRow, op)
	}
}

func (r *WorldRenderer) DrawSmoke(screen *ebiten.Image, particles []sim.SmokeParticle, camera Camera) {
	for _, p := range particles {
		if p.Variant < 0 || p.Variant >= len(r.smokeSprites) || p.Lifetime <= 0 {
			continue
		}
		img := r.smokeSprites[p.Variant]
		progress := mathutil.Clamp01(p.Age / p.Lifetime)
		size := p.StartSize + (p.EndSize-p.StartSize)*progress
		// Brief fade-in keeps each newly emitted circle from popping into view.
		alpha := p.Opacity * math.Min(1, progress/0.1) * (1 - progress)
		bounds := img.Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		op.GeoM.Scale(size/float64(bounds.Dx()), size/float64(bounds.Dy()))
		x, y := camera.WorldToScreen(p.X, p.Y)
		op.GeoM.Translate(x, y)
		op.ColorScale.ScaleAlpha(float32(alpha))
		screen.DrawImage(img, op)
	}
}

func (r *WorldRenderer) DrawSparks(screen *ebiten.Image, particles []sim.SparkParticle, camera Camera) {
	for _, p := range particles {
		if p.Lifetime <= 0 || p.Age >= p.Lifetime {
			continue
		}
		remaining := 1 - mathutil.Clamp01(p.Age/p.Lifetime)
		x, y := camera.WorldToScreen(p.X, p.Y)
		// Short trails follow each spark's velocity, shrinking as it slows.
		tailX := x - p.VelocityX*0.025
		tailY := y - p.VelocityY*0.025
		vector.StrokeLine(screen, float32(tailX), float32(tailY), float32(x), float32(y), 1.5,
			color.NRGBA{R: 255, G: uint8(110 + 110*remaining), B: uint8(30 + 120*remaining), A: uint8(255 * remaining)}, true)
	}
}

func (r *WorldRenderer) drawRenderables(screen *ebiten.Image, world *ecs.World, camera Camera, include func(ecs.Renderable) bool) {
	for _, renderable := range world.Renderables() {
		if include != nil && !include(renderable) {
			continue
		}

		spriteImage, ok := r.sprites[renderable.Sprite.SpriteID]
		if !ok {
			continue
		}

		bounds := spriteImage.Bounds()
		halfWidth := float64(bounds.Dx()) / 2
		halfHeight := float64(bounds.Dy()) / 2

		op := &ebiten.DrawImageOptions{}
		scaleX := 1.0
		scaleY := 1.0
		if renderable.Sprite.FlipX {
			scaleX = -1
		}
		if renderable.Sprite.FlipY {
			scaleY = -1
		}

		op.GeoM.Translate(-halfWidth, -halfHeight)
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Rotate(renderable.Transform.RotationRadians + renderable.Sprite.RotationOffsetRadians)

		screenX, screenY := camera.WorldToScreen(renderable.Transform.X, renderable.Transform.Y)
		op.GeoM.Translate(screenX, screenY)

		screen.DrawImage(spriteImage, op)
	}
}

func (r *WorldRenderer) drawHitRings(screen *ebiten.Image, world *ecs.World, camera Camera) {
	for _, flash := range world.HitFlashStates() {
		if flash.RemainingSeconds <= 0 || flash.DurationSeconds <= 0 {
			continue
		}
		remaining := mathutil.Clamp01(flash.RemainingSeconds / flash.DurationSeconds)
		radius := 5 + (1-remaining)*9
		x, y := camera.WorldToScreen(flash.X, flash.Y)
		vector.StrokeCircle(screen, float32(x), float32(y), float32(radius), 2,
			color.NRGBA{R: 255, G: 240, B: 180, A: uint8(255 * remaining)}, true)
	}
}

func (r *WorldRenderer) drawLighting(screen *ebiten.Image, world *ecs.World, camera Camera) {
	ambientDarkness := combinedAmbientDarkness(world.AmbientLights())
	if ambientDarkness > 0 {
		overlayAlpha := uint8(mathutil.Clamp01(ambientDarkness) * 255)
		bounds := screen.Bounds()
		viewWidth := float64(bounds.Dx())
		viewHeight := float64(bounds.Dy())

		ebitenutil.DrawRect(
			screen,
			0,
			0,
			viewWidth,
			viewHeight,
			color.RGBA{R: 0, G: 0, B: 0, A: overlayAlpha},
		)
	}

	lightSize := r.lightSprite.Bounds().Dx()
	if lightSize <= 0 {
		return
	}

	for _, source := range world.PointLightSources() {
		if source.Light.Radius <= 0 || source.Light.Intensity <= 0 {
			continue
		}

		screenX, screenY := camera.WorldToScreen(source.Transform.X, source.Transform.Y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(lightSize)/2, -float64(lightSize)/2)

		scale := (source.Light.Radius * 2) / float64(lightSize)
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(screenX, screenY)

		intensity := mathutil.Clamp01(source.Light.Intensity) * r.lightIntensityScale
		if intensity == 0 {
			continue
		}
		op.ColorScale.Scale(
			float32(source.Light.Color.R*intensity),
			float32(source.Light.Color.G*intensity),
			float32(source.Light.Color.B*intensity),
			float32(intensity),
		)
		op.Blend = ebiten.BlendLighter

		screen.DrawImage(r.lightSprite, op)
	}
}

func (r *WorldRenderer) drawTankHealthBars(screen *ebiten.Image, world *ecs.World, camera Camera) {
	for _, state := range world.TankHealthStates() {
		if state.Health.Max <= 0 {
			continue
		}

		ratio := mathutil.Clamp01(state.Health.Current / state.Health.Max)

		screenX, screenY := camera.WorldToScreen(state.Transform.X, state.Transform.Y)
		barX := screenX - r.healthBarWidth/2
		barY := screenY + r.healthBarYOffset

		ebitenutil.DrawRect(
			screen,
			barX-r.healthBarBorderSize,
			barY-r.healthBarBorderSize,
			r.healthBarWidth+2*r.healthBarBorderSize,
			r.healthBarHeight+2*r.healthBarBorderSize,
			color.RGBA{R: 0x12, G: 0x12, B: 0x12, A: 0xd8},
		)
		ebitenutil.DrawRect(
			screen,
			barX,
			barY,
			r.healthBarWidth,
			r.healthBarHeight,
			color.RGBA{R: 0x3a, G: 0x3a, B: 0x3a, A: 0xe8},
		)

		fillWidth := r.healthBarWidth * ratio
		if fillWidth > 0 {
			ebitenutil.DrawRect(
				screen,
				barX,
				barY,
				fillWidth,
				r.healthBarHeight,
				healthColor(ratio),
			)
		}
	}
}

func healthColor(ratio float64) color.RGBA {
	if ratio > 0.60 {
		return color.RGBA{R: 0x55, G: 0xcf, B: 0x58, A: 0xff}
	}
	if ratio > 0.30 {
		return color.RGBA{R: 0xe0, G: 0xb8, B: 0x3d, A: 0xff}
	}

	return color.RGBA{R: 0xdf, G: 0x54, B: 0x54, A: 0xff}
}

func combinedAmbientDarkness(lights []ecscomponents.AmbientLight) float64 {
	if len(lights) == 0 {
		return 0
	}

	darkness := 0.0
	for _, light := range lights {
		if light.Darkness > darkness {
			darkness = light.Darkness
		}
	}

	return mathutil.Clamp01(darkness)
}

func newRadialLightSprite(diameter int) *ebiten.Image {
	if diameter < 2 {
		diameter = 2
	}

	img := image.NewNRGBA(image.Rect(0, 0, diameter, diameter))
	radius := float64(diameter) / 2
	cx := radius - 0.5
	cy := radius - 0.5

	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			distance := math.Sqrt(dx*dx+dy*dy) / radius
			if distance >= 1 {
				continue
			}

			falloff := 1 - distance
			alpha := uint8(mathutil.Clamp01(math.Pow(falloff, 3.0)) * 190)
			img.Set(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: alpha})
		}
	}

	return ebiten.NewImageFromImage(img)
}
