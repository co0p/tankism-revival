package game

import (
	"github.com/co0p/tankism-revival/game/ecs"
	"github.com/co0p/tankism-revival/game/sim"
	"testing"
)

// TestGameOverExpiresGroundImpacts verifies that ground-impact effects expire
// after GroundSplashLifetime seconds of simulation, regardless of game state.
func TestGameOverExpiresGroundImpacts(t *testing.T) {
	ground := sim.NewGroundImpacts(8)
	ground.Burst(ecs.ProjectileState{Entity: 1})
	ground.Advance(sim.GroundSplashLifetime)
	if len(ground.Splashes()) != 0 || len(ground.Dust()) != 0 {
		t.Fatal("ground effects should expire after GroundSplashLifetime seconds")
	}
}
