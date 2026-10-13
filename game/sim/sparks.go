package sim

import (
	"math"
	"slices"

	"github.com/co0p/tankism-revival/game/ecs"
)

const (
	MaxSparkParticles = 192
	ImpactSparkCount  = 8
	sparkDrag         = 7.0
)

type SparkParticle struct {
	X, Y                 float64
	VelocityX, VelocityY float64
	Age, Lifetime        float64
}

// SparkParticles contains bounded, deterministic presentation data.
type SparkParticles struct {
	particles []SparkParticle
}

func NewSparkParticles() *SparkParticles {
	return &SparkParticles{particles: make([]SparkParticle, 0, MaxSparkParticles)}
}

func (s *SparkParticles) Particles() []SparkParticle {
	return slices.Clone(s.particles)
}

func (s *SparkParticles) Burst(projectile ecs.ProjectileState) {
	count := min(ImpactSparkCount, MaxSparkParticles-len(s.particles))
	for i := 0; i < count; i++ {
		roll := deterministicShotUnit(projectile.Entity, uint64(i))
		angle := (float64(i) + roll*0.6) * 2 * math.Pi / ImpactSparkCount
		speed := 65 + roll*90
		s.particles = append(s.particles, SparkParticle{
			X: projectile.Transform.X, Y: projectile.Transform.Y,
			VelocityX: math.Cos(angle) * speed,
			VelocityY: math.Sin(angle) * speed,
			Lifetime:  0.18 + roll*0.14,
		})
	}
}

// Advance runs before impacts so new bursts are visible on their spawning tick.
func (s *SparkParticles) Advance(deltaSeconds float64) {
	dt := math.Max(0, deltaSeconds)
	decay := math.Exp(-sparkDrag * dt)
	distanceFactor := (1 - decay) / sparkDrag
	alive := s.particles[:0]
	for _, p := range s.particles {
		p.Age += dt
		if p.Age >= p.Lifetime {
			continue
		}
		p.X += p.VelocityX * distanceFactor
		p.Y += p.VelocityY * distanceFactor
		p.VelocityX *= decay
		p.VelocityY *= decay
		alive = append(alive, p)
	}
	s.particles = alive
}
