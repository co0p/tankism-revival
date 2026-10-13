package game

import (
	"github.com/co0p/tankism-revival/game/ecs"
	"github.com/co0p/tankism-revival/game/sim"
	"testing"
)

func TestGameOverExpiresGroundImpacts(t *testing.T) {
	ground := sim.NewGroundImpacts(8)
	ground.Burst(ecs.ProjectileState{Entity: 1})
	g := Game{gameOver: true, groundImpacts: ground}
	if err := g.updateGame(sim.GroundSplashLifetime); err != nil {
		t.Fatal(err)
	}
	if len(ground.Splashes()) != 0 || len(ground.Dust()) != 0 {
		t.Fatal("ground effects should expire after game over")
	}
}
