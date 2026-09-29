package sim

import (
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestResolveTankDestructionsEnemySpawnsWreckAndExplosion(t *testing.T) {
	w := ecs.NewWorld()
	enemy := w.CreateEntity()

	must(t, w.SetTransform(enemy, ecscomponents.Transform{X: 100, Y: 120, RotationRadians: 0}))
	must(t, w.SetTankDrive(enemy, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetHealth(enemy, ecscomponents.Health{Current: 0, Max: 70}))
	must(t, w.SetEnemyAI(enemy, ecscomponents.EnemyAI{Behavior: ecscomponents.EnemyBehaviorAggressive}))

	cfg := TankDestructionConfig{
		ExplosionSpriteID:        "explosion_core",
		ExplosionSmokeSpriteID:   "explosion_smoke",
		EnemyWreckSpriteID:       "enemy_wreck",
		ExplosionZIndex:          120,
		EnemyWreckZIndex:         100,
		ExplosionLifetime:        0.20,
		ExplosionSmokeLifetime:   0.45,
		ExplosionLightRadius:     64,
		ExplosionLightIntensity:  0.7,
		ExplosionLightColor:      ecscomponents.LightColor{R: 1, G: 0.7, B: 0.4},
		EnemyWreckColliderRadius: 16,
	}

	if err := ResolveTankDestructions(w, 0, cfg); err != nil {
		t.Fatalf("resolve tank destructions: %v", err)
	}

	if w.IsAlive(enemy) {
		t.Fatalf("enemy entity should be destroyed")
	}

	if got := countRenderablesBySpriteID(w, "enemy_wreck"); got != 1 {
		t.Fatalf("expected one enemy wreck, got %d", got)
	}
	if got := countRenderablesBySpriteID(w, "explosion_core"); got != 1 {
		t.Fatalf("expected one explosion core sprite, got %d", got)
	}
	if got := countRenderablesBySpriteID(w, "explosion_smoke"); got != 1 {
		t.Fatalf("expected one explosion smoke sprite, got %d", got)
	}

	if len(w.PointLightSources()) != 1 {
		t.Fatalf("expected one explosion light source")
	}

	if lifetimes := w.LifetimeStates(); len(lifetimes) != 2 {
		t.Fatalf("expected two lifetime explosion entities, got %d", len(lifetimes))
	}

	solidBodies := w.SolidObstacleColliderBodies()
	if len(solidBodies) != 1 {
		t.Fatalf("expected one solid body for wreckage, got %d", len(solidBodies))
	}
	if solidBodies[0].Collider.Radius != 16 {
		t.Fatalf("expected wreckage collider radius 16, got %f", solidBodies[0].Collider.Radius)
	}
}

func TestResolveTankDestructionsPlayerSpawnsExplosionOnceAndMarksDestroyed(t *testing.T) {
	w := ecs.NewWorld()
	player := w.CreateEntity()

	must(t, w.SetTransform(player, ecscomponents.Transform{X: 20, Y: 40, RotationRadians: 0}))
	must(t, w.SetTankDrive(player, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetHealth(player, ecscomponents.Health{Current: 0, Max: 100}))

	cfg := TankDestructionConfig{
		ExplosionSpriteID:      "explosion_core",
		ExplosionSmokeSpriteID: "explosion_smoke",
		ExplosionZIndex:        120,
		ExplosionLifetime:      0.20,
		ExplosionSmokeLifetime: 0.45,
	}

	if err := ResolveTankDestructions(w, player, cfg); err != nil {
		t.Fatalf("resolve tank destructions first pass: %v", err)
	}

	health, ok := w.Health(player)
	if !ok {
		t.Fatalf("expected player health to remain attached")
	}
	if health.Max != 0 {
		t.Fatalf("expected player health max to be zero after destruction mark, got %f", health.Max)
	}

	if got := countRenderablesBySpriteID(w, "explosion_core"); got != 1 {
		t.Fatalf("expected one explosion core after first pass, got %d", got)
	}

	if err := ResolveTankDestructions(w, player, cfg); err != nil {
		t.Fatalf("resolve tank destructions second pass: %v", err)
	}

	if got := countRenderablesBySpriteID(w, "explosion_core"); got != 1 {
		t.Fatalf("expected player explosion not to respawn on second pass, got %d", got)
	}
}

func TestResolveTankDestructionsExplosionIntensityLayers(t *testing.T) {
	w := ecs.NewWorld()
	enemy := w.CreateEntity()

	must(t, w.SetTransform(enemy, ecscomponents.Transform{X: 40, Y: 50, RotationRadians: 0}))
	must(t, w.SetTankDrive(enemy, ecscomponents.TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}))
	must(t, w.SetHealth(enemy, ecscomponents.Health{Current: 0, Max: 50}))
	must(t, w.SetEnemyAI(enemy, ecscomponents.EnemyAI{Behavior: ecscomponents.EnemyBehaviorAggressive}))

	cfg := TankDestructionConfig{
		ExplosionSpriteID:       "explosion_core",
		ExplosionBurstSpriteID:  "explosion_burst",
		ExplosionSmokeSpriteID:  "explosion_smoke",
		ExplosionPlumeSpriteID:  "explosion_plume",
		EnemyWreckSpriteID:      "enemy_wreck",
		ExplosionZIndex:         120,
		EnemyWreckZIndex:        100,
		ExplosionLifetime:       0.20,
		ExplosionBurstLifetime:  0.10,
		ExplosionSmokeLifetime:  0.45,
		ExplosionPlumeLifetime:  0.62,
		ExplosionLightRadius:    64,
		ExplosionLightIntensity: 0.7,
		BurstLightRadius:        92,
		BurstLightIntensity:     0.95,
		ExplosionLightColor:     ecscomponents.LightColor{R: 1, G: 0.7, B: 0.4},
	}

	if err := ResolveTankDestructions(w, 0, cfg); err != nil {
		t.Fatalf("resolve tank destructions: %v", err)
	}

	if got := countRenderablesBySpriteID(w, "explosion_core"); got != 1 {
		t.Fatalf("expected one explosion core, got %d", got)
	}
	if got := countRenderablesBySpriteID(w, "explosion_burst"); got != 1 {
		t.Fatalf("expected one explosion burst, got %d", got)
	}
	if got := countRenderablesBySpriteID(w, "explosion_smoke"); got != 1 {
		t.Fatalf("expected one explosion smoke, got %d", got)
	}
	if got := countRenderablesBySpriteID(w, "explosion_plume"); got != 1 {
		t.Fatalf("expected one explosion plume, got %d", got)
	}

	if len(w.PointLightSources()) != 2 {
		t.Fatalf("expected two explosion light sources")
	}
	if lifetimes := w.LifetimeStates(); len(lifetimes) != 4 {
		t.Fatalf("expected four layered explosion lifetime entities, got %d", len(lifetimes))
	}
}

func countRenderablesBySpriteID(world *ecs.World, spriteID string) int {
	count := 0
	for _, renderable := range world.Renderables() {
		if renderable.Sprite.SpriteID == spriteID {
			count++
		}
	}

	return count
}
