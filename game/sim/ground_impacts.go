package sim

import (
	"github.com/co0p/tankism-revival/game/ecs"
	"math"
	"slices"
)

const (
	MaxGroundSplashes    = 128
	MaxDustParticles     = 192
	GroundSplashLifetime = 8.0
	GroundSplashSize     = 18.0
	GroundDustCount      = 6
)

type GroundSplash struct{ X, Y, RotationRadians, Age float64 }
type DustParticle struct {
	X, Y, VelocityX, VelocityY, Age, Lifetime float64
	Variant                                   int
}
type GroundImpacts struct {
	splashes []GroundSplash
	dust     []DustParticle
	variants int
}

func NewGroundImpacts(variants int) *GroundImpacts {
	return &GroundImpacts{splashes: make([]GroundSplash, 0, MaxGroundSplashes), dust: make([]DustParticle, 0, MaxDustParticles), variants: variants}
}
func (g *GroundImpacts) Splashes() []GroundSplash { return slices.Clone(g.splashes) }
func (g *GroundImpacts) Dust() []DustParticle     { return slices.Clone(g.dust) }

func (g *GroundImpacts) Burst(projectile ecs.ProjectileState) {
	if len(g.splashes) == MaxGroundSplashes {
		copy(g.splashes, g.splashes[1:])
		g.splashes = g.splashes[:len(g.splashes)-1]
	}
	g.splashes = append(g.splashes, GroundSplash{X: projectile.Transform.X, Y: projectile.Transform.Y, RotationRadians: deterministicShotUnit(projectile.Entity, 0) * 2 * math.Pi})
	if g.variants <= 0 {
		return
	}
	count := min(GroundDustCount, MaxDustParticles-len(g.dust))
	for i := 0; i < count; i++ {
		roll := deterministicShotUnit(projectile.Entity, uint64(i+1))
		angle := (float64(i) + roll) * 2 * math.Pi / GroundDustCount
		speed := 15 + roll*22
		g.dust = append(g.dust, DustParticle{X: projectile.Transform.X, Y: projectile.Transform.Y, VelocityX: math.Cos(angle) * speed, VelocityY: math.Sin(angle)*speed - 5, Lifetime: 0.4 + roll*0.25, Variant: i % g.variants})
	}
}

func (g *GroundImpacts) Advance(deltaSeconds float64) {
	dt := math.Max(0, deltaSeconds)
	splashes := g.splashes[:0]
	for _, splash := range g.splashes {
		splash.Age += dt
		if splash.Age < GroundSplashLifetime {
			splashes = append(splashes, splash)
		}
	}
	g.splashes = splashes
	decay := math.Exp(-4 * dt)
	dust := g.dust[:0]
	for _, p := range g.dust {
		p.Age += dt
		if p.Age >= p.Lifetime {
			continue
		}
		p.X += p.VelocityX * (1 - decay) / 4
		p.Y += p.VelocityY * (1 - decay) / 4
		p.VelocityX *= decay
		p.VelocityY *= decay
		dust = append(dust, p)
	}
	g.dust = dust
}
