package sim

import (
	"math"
	"slices"

	"github.com/co0p/tankism-revival/game/ecs"
	"github.com/co0p/tankism-revival/lib/mathutil"
)

const (
	MaxSmokeParticles    = 256
	SmokeHealthThreshold = 0.75
	smokeMaxRate         = 22.0
	smokeLifetime        = 1.6
)

type SmokeParticle struct {
	X, Y                 float64
	VelocityX, VelocityY float64
	Age, Lifetime        float64
	StartSize, EndSize   float64
	Opacity              float64
	Variant              int
}

// SmokeParticles owns capped presentation data, not framework resources.
type SmokeParticles struct {
	particles []SmokeParticle
	emission  map[ecs.Entity]float64
	sequence  uint64
	variants  int
}

func NewSmokeParticles(variants int) *SmokeParticles {
	return &SmokeParticles{
		particles: make([]SmokeParticle, 0, MaxSmokeParticles),
		emission:  make(map[ecs.Entity]float64),
		variants:  variants,
	}
}

func (s *SmokeParticles) Particles() []SmokeParticle {
	return slices.Clone(s.particles)
}

func (s *SmokeParticles) Advance(deltaSeconds float64) {
	dt := math.Max(0, deltaSeconds)
	alive := s.particles[:0]
	for _, p := range s.particles {
		p.Age += dt
		if p.Age >= p.Lifetime {
			continue
		}
		p.X += p.VelocityX * dt
		p.Y += p.VelocityY * dt
		alive = append(alive, p)
	}
	s.particles = alive
}

func (s *SmokeParticles) Update(world *ecs.World, deltaSeconds float64) {
	s.Advance(deltaSeconds)
	active := make(map[ecs.Entity]bool)
	for _, tank := range world.TankHealthStates() {
		if tank.Health.Max <= 0 || tank.Health.Current <= 0 || s.variants <= 0 {
			continue
		}
		ratio := mathutil.Clamp01(tank.Health.Current / tank.Health.Max)
		if ratio >= SmokeHealthThreshold {
			continue
		}
		active[tank.Entity] = true
		severity := (SmokeHealthThreshold - ratio) / SmokeHealthThreshold
		rate := 3 + severity*(smokeMaxRate-3)
		s.emission[tank.Entity] += math.Max(0, deltaSeconds) * rate
		count := int(s.emission[tank.Entity])
		s.emission[tank.Entity] -= float64(count)
		count = min(count, MaxSmokeParticles-len(s.particles))
		for i := 0; i < count; i++ {
			roll := deterministicShotUnit(tank.Entity, s.sequence)
			s.sequence++
			forwardX, forwardY := forwardVector(tank.Transform.RotationRadians)
			s.particles = append(s.particles, SmokeParticle{
				X:         tank.Transform.X - forwardX*8 + (roll-0.5)*6,
				Y:         tank.Transform.Y - forwardY*8,
				VelocityX: 5 + (roll-0.5)*12,
				VelocityY: -15 - roll*12,
				Lifetime:  smokeLifetime + roll*0.4,
				StartSize: 5 + severity*4,
				EndSize:   16 + severity*10,
				Opacity:   0.35 + severity*0.3,
				Variant:   int(s.sequence % uint64(s.variants)),
			})
		}
	}
	for entity := range s.emission {
		if !active[entity] {
			delete(s.emission, entity)
		}
	}
}
