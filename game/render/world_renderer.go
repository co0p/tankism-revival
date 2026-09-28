package render

import (
	"image"
	"image/color"
	"math"

	"github.com/co0p/tankism-revival/game"
	"github.com/co0p/tankism-revival/game/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	defaultLightSpriteSize     = 256
	defaultLightIntensityScale = 0.65

	defaultHealthBarWidth      = 34.0
	defaultHealthBarHeight     = 4.0
	defaultHealthBarYOffset    = 30.0
	defaultHealthBarBorderSize = 1.0
)

type Config struct {
	OutsideMapColor color.RGBA

	LightSpriteSize     int
	LightIntensityScale float64

	HealthBarWidth      float64
	HealthBarHeight     float64
	HealthBarYOffset    float64
	HealthBarBorderSize float64
}

type WorldRenderer struct {
	sprites map[string]*ebiten.Image

	outsideMapColor color.RGBA
	lightSprite     *ebiten.Image

	lightIntensityScale float64

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

	spriteMap := make(map[string]*ebiten.Image, len(sprites))
	for spriteID, img := range sprites {
		spriteMap[spriteID] = img
	}

	return &WorldRenderer{
		sprites:             spriteMap,
		outsideMapColor:     cfg.OutsideMapColor,
		lightSprite:         newRadialLightSprite(cfg.LightSpriteSize),
		lightIntensityScale: cfg.LightIntensityScale,
		healthBarWidth:      cfg.HealthBarWidth,
		healthBarHeight:     cfg.HealthBarHeight,
		healthBarYOffset:    cfg.HealthBarYOffset,
		healthBarBorderSize: cfg.HealthBarBorderSize,
	}
}

func (r *WorldRenderer) Draw(screen *ebiten.Image, world *ecs.World, camera game.Camera) {
	screen.Fill(r.outsideMapColor)

	r.drawRenderables(screen, world, camera, nil)
	r.drawLighting(screen, world, camera)
	r.drawRenderables(screen, world, camera, func(renderable ecs.Renderable) bool {
		return renderable.Sprite.ZIndex >= game.ZIndexObstacles
	})
	r.drawTankHealthBars(screen, world, camera)
}

func (r *WorldRenderer) drawRenderables(screen *ebiten.Image, world *ecs.World, camera game.Camera, include func(ecs.Renderable) bool) {
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

func (r *WorldRenderer) drawLighting(screen *ebiten.Image, world *ecs.World, camera game.Camera) {
	ambientDarkness := combinedAmbientDarkness(world.AmbientLights())
	if ambientDarkness > 0 {
		overlayAlpha := uint8(clamp01(ambientDarkness) * 255)
		viewWidth := camera.ViewWidth
		viewHeight := camera.ViewHeight
		if viewWidth <= 0 || viewHeight <= 0 {
			bounds := screen.Bounds()
			viewWidth = float64(bounds.Dx())
			viewHeight = float64(bounds.Dy())
		}

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

		intensity := clamp01(source.Light.Intensity) * r.lightIntensityScale
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

func (r *WorldRenderer) drawTankHealthBars(screen *ebiten.Image, world *ecs.World, camera game.Camera) {
	for _, state := range world.TankHealthStates() {
		if state.Health.Max <= 0 {
			continue
		}

		ratio := clamp01(state.Health.Current / state.Health.Max)

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

func combinedAmbientDarkness(lights []ecs.AmbientLight) float64 {
	if len(lights) == 0 {
		return 0
	}

	darkness := 0.0
	for _, light := range lights {
		if light.Darkness > darkness {
			darkness = light.Darkness
		}
	}

	return clamp01(darkness)
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
			alpha := uint8(clamp01(math.Pow(falloff, 3.0)) * 190)
			img.Set(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: alpha})
		}
	}

	return ebiten.NewImageFromImage(img)
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}

	return value
}
