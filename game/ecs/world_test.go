package ecs

import (
	"errors"
	"math"
	"testing"
)

func TestCreateEntityAllocatesUniqueAliveIDs(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()

	if e1 == 0 {
		t.Fatalf("first entity id should be non-zero")
	}

	if e2 <= e1 {
		t.Fatalf("entity ids should be monotonic: e1=%d e2=%d", e1, e2)
	}

	if !w.IsAlive(e1) || !w.IsAlive(e2) {
		t.Fatalf("new entities must be alive")
	}
}

func TestDestroyEntityRemovesAllComponents(t *testing.T) {
	w := NewWorld()
	e := w.CreateEntity()

	if err := w.SetTransform(e, Transform{X: 10, Y: 20, RotationRadians: math.Pi / 2}); err != nil {
		t.Fatalf("set transform: %v", err)
	}
	if err := w.SetSpriteRenderer(e, SpriteRenderer{SpriteID: "tank_red_small", ZIndex: 1}); err != nil {
		t.Fatalf("set sprite renderer: %v", err)
	}
	if err := w.SetTankControlIntent(e, TankControlIntent{Throttle: 1, Turn: -0.5}); err != nil {
		t.Fatalf("set tank control intent: %v", err)
	}
	if err := w.SetTankDrive(e, TankDrive{MoveSpeed: 100, TurnSpeedRadiansPS: math.Pi}); err != nil {
		t.Fatalf("set tank drive: %v", err)
	}
	if err := w.SetTankAcceleration(e, TankAcceleration{ForwardUnitsPerSecondSquared: 200}); err != nil {
		t.Fatalf("set tank acceleration: %v", err)
	}
	if err := w.SetTankDeceleration(e, TankDeceleration{CoastUnitsPerSecondSquared: 120, BrakeUnitsPerSecondSquared: 320}); err != nil {
		t.Fatalf("set tank deceleration: %v", err)
	}
	if err := w.SetTankVelocity(e, TankVelocity{ForwardUnitsPerSecond: 45}); err != nil {
		t.Fatalf("set tank velocity: %v", err)
	}
	if err := w.SetHealth(e, Health{Current: 80, Max: 100}); err != nil {
		t.Fatalf("set health: %v", err)
	}
	if err := w.SetWeapon(e, Weapon{CooldownSeconds: 0.3, ProjectileSpeed: 240, ProjectileLifetime: 1.2, ProjectileDamage: 10}); err != nil {
		t.Fatalf("set weapon: %v", err)
	}
	if err := w.SetProjectile(e, Projectile{Owner: e, VelocityX: 0, VelocityY: -1, RemainingLifetime: 1, Damage: 3}); err != nil {
		t.Fatalf("set projectile: %v", err)
	}
	if err := w.SetLifetime(e, Lifetime{RemainingSeconds: 0.1}); err != nil {
		t.Fatalf("set lifetime: %v", err)
	}
	if err := w.SetCircleCollider(e, CircleCollider{Radius: 10}); err != nil {
		t.Fatalf("set circle collider: %v", err)
	}
	if err := w.SetSolidObstacle(e); err != nil {
		t.Fatalf("set solid obstacle: %v", err)
	}
	if err := w.SetAmbientLight(e, AmbientLight{Darkness: 0.5}); err != nil {
		t.Fatalf("set ambient light: %v", err)
	}
	if err := w.SetPointLight(e, PointLight{Radius: 120, Intensity: 1, Color: LightColorWhite}); err != nil {
		t.Fatalf("set point light: %v", err)
	}

	if ok := w.DestroyEntity(e); !ok {
		t.Fatalf("destroy entity should succeed")
	}

	if w.IsAlive(e) {
		t.Fatalf("entity should not be alive after destroy")
	}

	if _, ok := w.Transform(e); ok {
		t.Fatalf("transform should be removed on destroy")
	}

	if _, ok := w.SpriteRenderer(e); ok {
		t.Fatalf("sprite renderer should be removed on destroy")
	}

	if _, ok := w.TankControlIntent(e); ok {
		t.Fatalf("tank control intent should be removed on destroy")
	}

	if _, ok := w.TankDrive(e); ok {
		t.Fatalf("tank drive should be removed on destroy")
	}

	if _, ok := w.TankAcceleration(e); ok {
		t.Fatalf("tank acceleration should be removed on destroy")
	}

	if _, ok := w.TankDeceleration(e); ok {
		t.Fatalf("tank deceleration should be removed on destroy")
	}

	if _, ok := w.TankVelocity(e); ok {
		t.Fatalf("tank velocity should be removed on destroy")
	}

	if _, ok := w.Health(e); ok {
		t.Fatalf("health should be removed on destroy")
	}

	if _, ok := w.Weapon(e); ok {
		t.Fatalf("weapon should be removed on destroy")
	}

	if _, ok := w.Projectile(e); ok {
		t.Fatalf("projectile should be removed on destroy")
	}

	if _, ok := w.Lifetime(e); ok {
		t.Fatalf("lifetime should be removed on destroy")
	}

	if _, ok := w.CircleCollider(e); ok {
		t.Fatalf("circle collider should be removed on destroy")
	}

	if w.IsSolidObstacle(e) {
		t.Fatalf("solid obstacle marker should be removed on destroy")
	}

	if _, ok := w.AmbientLight(e); ok {
		t.Fatalf("ambient light should be removed on destroy")
	}

	if _, ok := w.PointLight(e); ok {
		t.Fatalf("point light should be removed on destroy")
	}

	if ok := w.DestroyEntity(e); ok {
		t.Fatalf("destroying an already-destroyed entity should return false")
	}
}

func TestSetComponentOnDeadEntityFails(t *testing.T) {
	w := NewWorld()

	err := w.SetTransform(999, Transform{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set transform on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetSpriteRenderer(999, SpriteRenderer{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set sprite renderer on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankControlIntent(999, TankControlIntent{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank control intent on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankDrive(999, TankDrive{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank drive on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankAcceleration(999, TankAcceleration{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank acceleration on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankDeceleration(999, TankDeceleration{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank deceleration on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetTankVelocity(999, TankVelocity{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set tank velocity on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetHealth(999, Health{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set health on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetWeapon(999, Weapon{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set weapon on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetProjectile(999, Projectile{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set projectile on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetLifetime(999, Lifetime{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set lifetime on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetCircleCollider(999, CircleCollider{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set circle collider on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetSolidObstacle(999)
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set solid obstacle on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetAmbientLight(999, AmbientLight{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set ambient light on dead entity should return ErrEntityNotAlive, got: %v", err)
	}

	err = w.SetPointLight(999, PointLight{})
	if !errors.Is(err, ErrEntityNotAlive) {
		t.Fatalf("set point light on dead entity should return ErrEntityNotAlive, got: %v", err)
	}
}

func TestRenderablesReturnsEntitiesWithTransformAndSpriteSortedByZIndexThenEntity(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()
	e3 := w.CreateEntity()

	if err := w.SetTransform(e1, Transform{X: 100, Y: 200, RotationRadians: math.Pi}); err != nil {
		t.Fatalf("set transform e1: %v", err)
	}
	if err := w.SetSpriteRenderer(e1, SpriteRenderer{SpriteID: "red", ZIndex: 5}); err != nil {
		t.Fatalf("set sprite e1: %v", err)
	}

	if err := w.SetTransform(e2, Transform{X: 1, Y: 2, RotationRadians: 0}); err != nil {
		t.Fatalf("set transform e2: %v", err)
	}
	if err := w.SetSpriteRenderer(e2, SpriteRenderer{SpriteID: "blue", ZIndex: 1}); err != nil {
		t.Fatalf("set sprite e2: %v", err)
	}

	if err := w.SetTransform(e3, Transform{X: 5, Y: 6, RotationRadians: math.Pi / 4}); err != nil {
		t.Fatalf("set transform e3: %v", err)
	}

	renderables := w.Renderables()

	if len(renderables) != 2 {
		t.Fatalf("expected 2 renderables, got %d", len(renderables))
	}

	if renderables[0].Entity != e2 || renderables[1].Entity != e1 {
		t.Fatalf("unexpected render order: got [%d, %d], want [%d, %d]", renderables[0].Entity, renderables[1].Entity, e2, e1)
	}
}

func TestTransformPreservesCenterCoordinatesAndRadians(t *testing.T) {
	w := NewWorld()
	e := w.CreateEntity()

	want := Transform{
		X:               480.5,
		Y:               270.25,
		RotationRadians: math.Pi / 2,
	}

	if err := w.SetTransform(e, want); err != nil {
		t.Fatalf("set transform: %v", err)
	}

	got, ok := w.Transform(e)
	if !ok {
		t.Fatalf("expected transform to exist")
	}

	if got != want {
		t.Fatalf("transform mismatch: got %+v want %+v", got, want)
	}
}

func TestTankMoversReturnsDriveAndDefaultIntent(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()

	if err := w.SetTransform(e1, Transform{X: 100, Y: 200, RotationRadians: 0}); err != nil {
		t.Fatalf("set transform e1: %v", err)
	}
	if err := w.SetTankDrive(e1, TankDrive{MoveSpeed: 180, TurnSpeedRadiansPS: math.Pi}); err != nil {
		t.Fatalf("set drive e1: %v", err)
	}

	if err := w.SetTransform(e2, Transform{X: 1, Y: 2, RotationRadians: math.Pi / 2}); err != nil {
		t.Fatalf("set transform e2: %v", err)
	}
	if err := w.SetTankDrive(e2, TankDrive{MoveSpeed: 90, TurnSpeedRadiansPS: math.Pi / 2}); err != nil {
		t.Fatalf("set drive e2: %v", err)
	}
	if err := w.SetTankAcceleration(e2, TankAcceleration{ForwardUnitsPerSecondSquared: 300}); err != nil {
		t.Fatalf("set acceleration e2: %v", err)
	}
	if err := w.SetTankDeceleration(e2, TankDeceleration{CoastUnitsPerSecondSquared: 140, BrakeUnitsPerSecondSquared: 360}); err != nil {
		t.Fatalf("set deceleration e2: %v", err)
	}
	if err := w.SetTankVelocity(e2, TankVelocity{ForwardUnitsPerSecond: 12}); err != nil {
		t.Fatalf("set velocity e2: %v", err)
	}
	if err := w.SetTankControlIntent(e2, TankControlIntent{Throttle: 0.5, Turn: -1}); err != nil {
		t.Fatalf("set intent e2: %v", err)
	}

	movers := w.TankMovers()
	if len(movers) != 2 {
		t.Fatalf("expected 2 tank movers, got %d", len(movers))
	}

	if movers[0].Entity != e1 || movers[1].Entity != e2 {
		t.Fatalf("unexpected mover order: got [%d, %d], want [%d, %d]", movers[0].Entity, movers[1].Entity, e1, e2)
	}

	if movers[0].Intent != (TankControlIntent{}) {
		t.Fatalf("missing intent should default to zero intent, got %+v", movers[0].Intent)
	}
	if movers[0].Acceleration != (TankAcceleration{}) {
		t.Fatalf("missing acceleration should default to zero acceleration, got %+v", movers[0].Acceleration)
	}
	if movers[0].Deceleration != (TankDeceleration{}) {
		t.Fatalf("missing deceleration should default to zero deceleration, got %+v", movers[0].Deceleration)
	}
	if movers[0].Velocity != (TankVelocity{}) {
		t.Fatalf("missing velocity should default to zero velocity, got %+v", movers[0].Velocity)
	}

	if movers[1].Intent != (TankControlIntent{Throttle: 0.5, Turn: -1}) {
		t.Fatalf("unexpected e2 intent: %+v", movers[1].Intent)
	}
	if movers[1].Acceleration != (TankAcceleration{ForwardUnitsPerSecondSquared: 300}) {
		t.Fatalf("unexpected e2 acceleration: %+v", movers[1].Acceleration)
	}
	if movers[1].Deceleration != (TankDeceleration{CoastUnitsPerSecondSquared: 140, BrakeUnitsPerSecondSquared: 360}) {
		t.Fatalf("unexpected e2 deceleration: %+v", movers[1].Deceleration)
	}
	if movers[1].Velocity != (TankVelocity{ForwardUnitsPerSecond: 12}) {
		t.Fatalf("unexpected e2 velocity: %+v", movers[1].Velocity)
	}
}

func TestColliderQueriesReturnExpectedBodies(t *testing.T) {
	w := NewWorld()

	tank := w.CreateEntity()
	if err := w.SetTransform(tank, Transform{X: 10, Y: 20}); err != nil {
		t.Fatalf("set tank transform: %v", err)
	}
	if err := w.SetTankDrive(tank, TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}); err != nil {
		t.Fatalf("set tank drive: %v", err)
	}
	if err := w.SetCircleCollider(tank, CircleCollider{Radius: 9}); err != nil {
		t.Fatalf("set tank collider: %v", err)
	}

	obstacle := w.CreateEntity()
	if err := w.SetTransform(obstacle, Transform{X: 100, Y: 120}); err != nil {
		t.Fatalf("set obstacle transform: %v", err)
	}
	if err := w.SetCircleCollider(obstacle, CircleCollider{Radius: 14}); err != nil {
		t.Fatalf("set obstacle collider: %v", err)
	}
	if err := w.SetSolidObstacle(obstacle); err != nil {
		t.Fatalf("set obstacle solid: %v", err)
	}

	tankBodies := w.TankColliderBodies()
	if len(tankBodies) != 1 {
		t.Fatalf("expected 1 tank collider body, got %d", len(tankBodies))
	}
	if tankBodies[0].Entity != tank {
		t.Fatalf("unexpected tank collider entity: got %d want %d", tankBodies[0].Entity, tank)
	}

	obstacleBodies := w.SolidObstacleColliderBodies()
	if len(obstacleBodies) != 1 {
		t.Fatalf("expected 1 obstacle collider body, got %d", len(obstacleBodies))
	}
	if obstacleBodies[0].Entity != obstacle {
		t.Fatalf("unexpected obstacle collider entity: got %d want %d", obstacleBodies[0].Entity, obstacle)
	}
}

func TestPointLightSourcesReturnsEntitiesWithTransformAndPointLight(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()

	if err := w.SetTransform(e1, Transform{X: 50, Y: 60}); err != nil {
		t.Fatalf("set transform e1: %v", err)
	}
	if err := w.SetPointLight(e1, PointLight{Radius: 120, Intensity: 0.8, Color: LightColor{R: 1, G: 0.9, B: 0.7}}); err != nil {
		t.Fatalf("set point light e1: %v", err)
	}

	if err := w.SetPointLight(e2, PointLight{Radius: 90, Intensity: 1, Color: LightColorWhite}); err != nil {
		t.Fatalf("set point light e2: %v", err)
	}

	sources := w.PointLightSources()
	if len(sources) != 1 {
		t.Fatalf("expected 1 point light source with transform, got %d", len(sources))
	}

	if sources[0].Entity != e1 {
		t.Fatalf("unexpected source entity: got %d want %d", sources[0].Entity, e1)
	}
}

func TestTankHealthStatesReturnsOnlyTanksWithHealthAndTransform(t *testing.T) {
	w := NewWorld()

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()
	e3 := w.CreateEntity()

	if err := w.SetTransform(e1, Transform{X: 10, Y: 20}); err != nil {
		t.Fatalf("set transform e1: %v", err)
	}
	if err := w.SetTankDrive(e1, TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}); err != nil {
		t.Fatalf("set drive e1: %v", err)
	}
	if err := w.SetHealth(e1, Health{Current: 90, Max: 100}); err != nil {
		t.Fatalf("set health e1: %v", err)
	}

	if err := w.SetTankDrive(e2, TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}); err != nil {
		t.Fatalf("set drive e2: %v", err)
	}
	if err := w.SetHealth(e2, Health{Current: 50, Max: 100}); err != nil {
		t.Fatalf("set health e2: %v", err)
	}

	if err := w.SetTransform(e3, Transform{X: 30, Y: 40}); err != nil {
		t.Fatalf("set transform e3: %v", err)
	}
	if err := w.SetHealth(e3, Health{Current: 10, Max: 100}); err != nil {
		t.Fatalf("set health e3: %v", err)
	}

	states := w.TankHealthStates()
	if len(states) != 1 {
		t.Fatalf("expected 1 tank health state, got %d", len(states))
	}

	if states[0].Entity != e1 {
		t.Fatalf("unexpected tank health entity: got %d want %d", states[0].Entity, e1)
	}
}

func TestArmedTanksAndProjectileStatesQueries(t *testing.T) {
	w := NewWorld()

	shooter := w.CreateEntity()
	if err := w.SetTransform(shooter, Transform{X: 5, Y: 6}); err != nil {
		t.Fatalf("set shooter transform: %v", err)
	}
	if err := w.SetTankDrive(shooter, TankDrive{MoveSpeed: 1, TurnSpeedRadiansPS: 1}); err != nil {
		t.Fatalf("set shooter drive: %v", err)
	}
	if err := w.SetTankControlIntent(shooter, TankControlIntent{Fire: true}); err != nil {
		t.Fatalf("set shooter intent: %v", err)
	}
	if err := w.SetWeapon(shooter, Weapon{CooldownSeconds: 0.4, ProjectileSpeed: 300, ProjectileLifetime: 1, ProjectileDamage: 9, ProjectileSpriteID: "p", ProjectileRadius: 2}); err != nil {
		t.Fatalf("set shooter weapon: %v", err)
	}

	armed := w.ArmedTanks()
	if len(armed) != 1 {
		t.Fatalf("expected 1 armed tank, got %d", len(armed))
	}
	if armed[0].Entity != shooter {
		t.Fatalf("unexpected armed tank entity: got %d want %d", armed[0].Entity, shooter)
	}

	projectileEntity := w.CreateEntity()
	if err := w.SetTransform(projectileEntity, Transform{X: 1, Y: 2}); err != nil {
		t.Fatalf("set projectile transform: %v", err)
	}
	if err := w.SetCircleCollider(projectileEntity, CircleCollider{Radius: 2}); err != nil {
		t.Fatalf("set projectile collider: %v", err)
	}
	if err := w.SetProjectile(projectileEntity, Projectile{Owner: shooter, VelocityX: 0, VelocityY: -1, RemainingLifetime: 1, Damage: 3}); err != nil {
		t.Fatalf("set projectile: %v", err)
	}

	projectiles := w.ProjectileStates()
	if len(projectiles) != 1 {
		t.Fatalf("expected 1 projectile state, got %d", len(projectiles))
	}
	if projectiles[0].Entity != projectileEntity {
		t.Fatalf("unexpected projectile entity: got %d want %d", projectiles[0].Entity, projectileEntity)
	}
}
