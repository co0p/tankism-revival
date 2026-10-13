package sim

import "math"

const (
	PlayerHitShakeDurationSeconds = 0.16
	PlayerHitShakeAmplitude       = 3.0
	playerHitShakeFrequency       = 42.0
)

// PlayerHitShake is presentation-only state. Its offsets never move world entities.
type PlayerHitShake struct {
	RemainingSeconds float64
	phase            float64
	X, Y             float64
}

func (s *PlayerHitShake) Update(deltaSeconds float64, tookDamage bool) {
	if tookDamage {
		s.RemainingSeconds = PlayerHitShakeDurationSeconds
		s.phase = 0
	} else {
		s.RemainingSeconds = math.Max(0, s.RemainingSeconds-math.Max(0, deltaSeconds))
	}
	if s.RemainingSeconds <= 0 {
		s.X, s.Y = 0, 0
		return
	}

	s.phase += math.Max(0, deltaSeconds) * playerHitShakeFrequency
	envelope := s.RemainingSeconds / PlayerHitShakeDurationSeconds
	s.X = math.Sin(s.phase*2.1) * PlayerHitShakeAmplitude * envelope
	s.Y = math.Cos(s.phase*1.7) * PlayerHitShakeAmplitude * 0.8 * envelope
}
