package sim

import (
	"math"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestNewEnemyAIIsDeterministicForEntity(t *testing.T) {
	a := NewEnemyAI(42)
	b := NewEnemyAI(42)

	if a.Behavior != b.Behavior ||
		a.FavoredBehavior != b.FavoredBehavior ||
		a.DecisionCooldownSeconds != b.DecisionCooldownSeconds ||
		a.RNGState != b.RNGState ||
		len(a.Path) != len(b.Path) ||
		a.PathIndex != b.PathIndex {
		t.Fatalf("expected deterministic enemy ai seed for same entity")
	}

	if !isBehaviorValid(a.Behavior) {
		t.Fatalf("expected valid initial behavior, got %v", a.Behavior)
	}
	if !isBehaviorValid(a.FavoredBehavior) {
		t.Fatalf("expected valid favored behavior, got %v", a.FavoredBehavior)
	}
	if a.RNGState == 0 {
		t.Fatalf("expected non-zero RNG state")
	}
}

func TestUpdateEnemyTankIntentAggressiveFireWhenAlignedAndClear(t *testing.T) {
	w, mapEntity, enemy, player := newEnemyAITestWorld(t)

	must(t, w.SetTransform(enemy, ecscomponents.Transform{X: 160, Y: 200, RotationRadians: 0}))
	must(t, w.SetTransform(player, ecscomponents.Transform{X: 160, Y: 60, RotationRadians: 0}))
	must(t, w.SetEnemyAI(enemy, ecscomponents.EnemyAI{
		Behavior:                ecscomponents.EnemyBehaviorAggressive,
		FavoredBehavior:         ecscomponents.EnemyBehaviorAggressive,
		DecisionCooldownSeconds: 10,
		RNGState:                7,
	}))

	if err := UpdateEnemyTankIntent(w, enemy, player, mapEntity, 10, 1.0/60.0); err != nil {
		t.Fatalf("update enemy intent: %v", err)
	}

	intent, ok := w.TankControlIntent(enemy)
	if !ok {
		t.Fatalf("expected enemy intent")
	}

	approxEnemyAIEqual(t, intent.Turn, 0, 1e-9)
	approxEnemyAIEqual(t, intent.Throttle, 1, 1e-9)
	if !intent.Fire {
		t.Fatalf("expected aggressive enemy to fire when aligned and clear")
	}
}

func TestUpdateEnemyTankIntentStopsWithinStopDistance(t *testing.T) {
	w, mapEntity, enemy, player := newEnemyAITestWorld(t)

	must(t, w.SetTransform(enemy, ecscomponents.Transform{X: 160, Y: 200, RotationRadians: 0}))
	must(t, w.SetTransform(player, ecscomponents.Transform{X: 160, Y: 180, RotationRadians: 0}))
	must(t, w.SetEnemyAI(enemy, ecscomponents.EnemyAI{
		Behavior:                ecscomponents.EnemyBehaviorAggressive,
		FavoredBehavior:         ecscomponents.EnemyBehaviorAggressive,
		DecisionCooldownSeconds: 10,
		RNGState:                9,
	}))

	if err := UpdateEnemyTankIntent(w, enemy, player, mapEntity, 40, 1.0/60.0); err != nil {
		t.Fatalf("update enemy intent: %v", err)
	}

	intent, ok := w.TankControlIntent(enemy)
	if !ok {
		t.Fatalf("expected enemy intent")
	}

	approxEnemyAIEqual(t, intent.Throttle, 0, 1e-9)
}

func TestUpdateEnemyTankIntentBlockedLOSPreventsFire(t *testing.T) {
	w, mapEntity, enemy, player := newEnemyAITestWorld(t)

	must(t, w.SetTransform(enemy, ecscomponents.Transform{X: 160, Y: 220, RotationRadians: 0}))
	must(t, w.SetTransform(player, ecscomponents.Transform{X: 160, Y: 40, RotationRadians: 0}))
	must(t, w.SetEnemyAI(enemy, ecscomponents.EnemyAI{
		Behavior:                ecscomponents.EnemyBehaviorAggressive,
		FavoredBehavior:         ecscomponents.EnemyBehaviorAggressive,
		DecisionCooldownSeconds: 10,
		RNGState:                11,
	}))

	obstacle := w.CreateEntity()
	must(t, w.SetTransform(obstacle, ecscomponents.Transform{X: 160, Y: 130}))
	must(t, w.SetCircleCollider(obstacle, ecscomponents.CircleCollider{Radius: 20}))
	must(t, w.SetSolidObstacle(obstacle))

	if err := UpdateEnemyTankIntent(w, enemy, player, mapEntity, 10, 1.0/60.0); err != nil {
		t.Fatalf("update enemy intent: %v", err)
	}

	intent, ok := w.TankControlIntent(enemy)
	if !ok {
		t.Fatalf("expected enemy intent")
	}

	if intent.Fire {
		t.Fatalf("expected blocked LOS to prevent firing")
	}
}

func TestFindPathAStarBuildsPathAroundObstacle(t *testing.T) {
	nav := buildNavGrid(
		ecscomponents.MapBounds{Width: 160, Height: 160},
		[]ecs.ColliderBody{{
			Transform: ecscomponents.Transform{X: 80, Y: 80},
			Collider:  ecscomponents.CircleCollider{Radius: 8},
		}},
		16,
	)

	start := nav.worldToCell(16, 80)
	goal := nav.worldToCell(144, 80)
	path := findPathAStar(nav, start, goal)
	if len(path) < 2 {
		t.Fatalf("expected non-empty multi-step path around obstacle")
	}

	for _, step := range path {
		if step.col == 2 && step.row == 2 {
			t.Fatalf("path should not pass through blocked center cell")
		}
	}
}

func TestUpdateEnemyTankIntentCautiousDoesNotFireOutsideDecisionTick(t *testing.T) {
	w, mapEntity, enemy, player := newEnemyAITestWorld(t)

	must(t, w.SetTransform(enemy, ecscomponents.Transform{X: 160, Y: 200, RotationRadians: 0}))
	must(t, w.SetTransform(player, ecscomponents.Transform{X: 160, Y: 40, RotationRadians: 0}))
	must(t, w.SetEnemyAI(enemy, ecscomponents.EnemyAI{
		Behavior:                ecscomponents.EnemyBehaviorCautious,
		FavoredBehavior:         ecscomponents.EnemyBehaviorCautious,
		DecisionCooldownSeconds: 10,
		RNGState:                21,
	}))

	if err := UpdateEnemyTankIntent(w, enemy, player, mapEntity, 10, 1.0/60.0); err != nil {
		t.Fatalf("update enemy intent: %v", err)
	}

	intent, ok := w.TankControlIntent(enemy)
	if !ok {
		t.Fatalf("expected enemy intent")
	}

	if intent.Fire {
		t.Fatalf("expected cautious behavior not to fire outside decision tick")
	}
}

func TestUpdateEnemyTankIntentsUpdatesAllEnemyEntities(t *testing.T) {
	w := ecs.NewWorld()
	mapEntity := w.CreateEntity()
	player := w.CreateEntity()
	enemyA := w.CreateEntity()
	enemyB := w.CreateEntity()

	must(t, w.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: 320, Height: 320}))
	must(t, w.SetTransform(player, ecscomponents.Transform{X: 160, Y: 40, RotationRadians: 0}))

	for _, enemy := range []ecs.Entity{enemyA, enemyB} {
		must(t, w.SetTransform(enemy, ecscomponents.Transform{X: 160, Y: 200, RotationRadians: 0}))
		must(t, w.SetTankDrive(enemy, ecscomponents.TankDrive{MoveSpeed: 100, TurnSpeedRadiansPS: 1}))
		must(t, w.SetTankControlIntent(enemy, ecscomponents.TankControlIntent{}))
		must(t, w.SetHealth(enemy, ecscomponents.Health{Current: 70, Max: 70}))
		must(t, w.SetWeapon(enemy, ecscomponents.Weapon{CooldownSeconds: 1, CooldownRemaining: 0}))
		must(t, w.SetCircleCollider(enemy, ecscomponents.CircleCollider{Radius: 16}))
		must(t, w.SetEnemyAI(enemy, ecscomponents.EnemyAI{
			Behavior:                ecscomponents.EnemyBehaviorAggressive,
			FavoredBehavior:         ecscomponents.EnemyBehaviorAggressive,
			DecisionCooldownSeconds: 10,
			RNGState:                uint64(enemy) + 1,
		}))
	}

	if err := UpdateEnemyTankIntents(w, player, mapEntity, 10, 1.0/60.0); err != nil {
		t.Fatalf("update enemy intents: %v", err)
	}

	intentA, ok := w.TankControlIntent(enemyA)
	if !ok {
		t.Fatalf("expected intent for enemyA")
	}
	intentB, ok := w.TankControlIntent(enemyB)
	if !ok {
		t.Fatalf("expected intent for enemyB")
	}

	if intentA.Throttle == 0 || intentB.Throttle == 0 {
		t.Fatalf("expected both enemy intents to be updated with movement")
	}
}

func TestMarkovTransitionWeightsShiftTowardFleeingOnLowHealth(t *testing.T) {
	high := markovTransitionWeights(ecscomponents.EnemyBehaviorAggressive, ecscomponents.EnemyBehaviorAggressive, 0.9, true, 0)
	low := markovTransitionWeights(ecscomponents.EnemyBehaviorAggressive, ecscomponents.EnemyBehaviorAggressive, 0.2, true, 0)

	if low.fleeing <= high.fleeing {
		t.Fatalf("expected low health to increase fleeing weight: low=%f high=%f", low.fleeing, high.fleeing)
	}
	if low.aggressive >= high.aggressive {
		t.Fatalf("expected low health to decrease aggressive weight: low=%f high=%f", low.aggressive, high.aggressive)
	}
}

func TestMarkovTransitionWeightsIncreaseAggressiveWithNearbyFriends(t *testing.T) {
	withoutFriends := markovTransitionWeights(ecscomponents.EnemyBehaviorCautious, ecscomponents.EnemyBehaviorCautious, 0.7, true, 0)
	withFriends := markovTransitionWeights(ecscomponents.EnemyBehaviorCautious, ecscomponents.EnemyBehaviorCautious, 0.7, true, 3)

	if withFriends.aggressive <= withoutFriends.aggressive {
		t.Fatalf("expected nearby friends to increase aggressive probability: without=%f with=%f", withoutFriends.aggressive, withFriends.aggressive)
	}
}

func newEnemyAITestWorld(t *testing.T) (*ecs.World, ecs.Entity, ecs.Entity, ecs.Entity) {
	t.Helper()

	w := ecs.NewWorld()
	mapEntity := w.CreateEntity()
	enemy := w.CreateEntity()
	player := w.CreateEntity()

	must(t, w.SetMapBounds(mapEntity, ecscomponents.MapBounds{Width: 320, Height: 320}))
	must(t, w.SetTankControlIntent(enemy, ecscomponents.TankControlIntent{}))
	must(t, w.SetHealth(enemy, ecscomponents.Health{Current: 70, Max: 70}))
	must(t, w.SetWeapon(enemy, ecscomponents.Weapon{CooldownSeconds: 1, CooldownRemaining: 0}))
	must(t, w.SetCircleCollider(enemy, ecscomponents.CircleCollider{Radius: 16}))

	return w, mapEntity, enemy, player
}

func approxEnemyAIEqual(t *testing.T, got, want, tol float64) {
	t.Helper()

	if math.Abs(got-want) > tol {
		t.Fatalf("mismatch: got %f want %f (tol %f)", got, want, tol)
	}
}
