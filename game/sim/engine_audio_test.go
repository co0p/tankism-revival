package sim

import (
	"math"
	"testing"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestUpdateEngineAudioVolumes(t *testing.T) {
	w := ecs.NewWorld()
	tank := w.CreateEntity()

	must(t, w.SetTransform(tank, ecscomponents.Transform{}))
	must(t, w.SetTankDrive(tank, ecscomponents.TankDrive{MoveSpeed: 100}))
	must(t, w.SetTankAcceleration(tank, ecscomponents.TankAcceleration{ForwardUnitsPerSecondSquared: 200}))
	must(t, w.SetTankDeceleration(tank, ecscomponents.TankDeceleration{CoastUnitsPerSecondSquared: 100, BrakeUnitsPerSecondSquared: 200}))
	must(t, w.SetTankVelocity(tank, ecscomponents.TankVelocity{ForwardUnitsPerSecond: 0}))
	must(t, w.SetTankControlIntent(tank, ecscomponents.TankControlIntent{}))
	must(t, w.SetEngineAudio(tank, ecscomponents.EngineAudio{SoundID: "engine_heavy_loop", IdleVolume: 0.08, MovingVolume: 0.32}))

	if err := UpdateEngineAudioVolumes(w, 1.0/60.0); err != nil {
		t.Fatalf("update engine audio volumes: %v", err)
	}

	engineAudio, ok := w.EngineAudio(tank)
	if !ok {
		t.Fatalf("expected engine audio on tank")
	}
	if math.Abs(engineAudio.CurrentVolume-0.08) > 0.0001 {
		t.Fatalf("idle volume mismatch: got %f want %f", engineAudio.CurrentVolume, 0.08)
	}

	must(t, w.SetTankVelocity(tank, ecscomponents.TankVelocity{ForwardUnitsPerSecond: 50}))
	if err := UpdateEngineAudioVolumes(w, 1.0/60.0); err != nil {
		t.Fatalf("update engine audio volumes at half speed: %v", err)
	}
	engineAudio, _ = w.EngineAudio(tank)
	if math.Abs(engineAudio.CurrentVolume-0.20) > 0.0001 {
		t.Fatalf("half-speed volume mismatch: got %f want %f", engineAudio.CurrentVolume, 0.20)
	}

	must(t, w.SetTankVelocity(tank, ecscomponents.TankVelocity{ForwardUnitsPerSecond: -140}))
	if err := UpdateEngineAudioVolumes(w, 1.0/60.0); err != nil {
		t.Fatalf("update engine audio volumes above max speed: %v", err)
	}
	engineAudio, _ = w.EngineAudio(tank)
	if math.Abs(engineAudio.CurrentVolume-0.32) > 0.0001 {
		t.Fatalf("max volume mismatch: got %f want %f", engineAudio.CurrentVolume, 0.32)
	}
}

func TestUpdateEngineAudioVolumesSmoothsAttackAndRelease(t *testing.T) {
	w := ecs.NewWorld()
	tank := w.CreateEntity()

	must(t, w.SetTransform(tank, ecscomponents.Transform{}))
	must(t, w.SetTankDrive(tank, ecscomponents.TankDrive{MoveSpeed: 100}))
	must(t, w.SetTankAcceleration(tank, ecscomponents.TankAcceleration{ForwardUnitsPerSecondSquared: 200}))
	must(t, w.SetTankDeceleration(tank, ecscomponents.TankDeceleration{CoastUnitsPerSecondSquared: 100, BrakeUnitsPerSecondSquared: 200}))
	must(t, w.SetTankVelocity(tank, ecscomponents.TankVelocity{ForwardUnitsPerSecond: 100}))
	must(t, w.SetTankControlIntent(tank, ecscomponents.TankControlIntent{}))
	must(t, w.SetEngineAudio(tank, ecscomponents.EngineAudio{
		SoundID:       "engine_heavy_loop",
		IdleVolume:    0.10,
		MovingVolume:  0.30,
		AttackPerSec:  1.20,
		ReleasePerSec: 0.60,
		CurrentVolume: 0.10,
	}))

	if err := UpdateEngineAudioVolumes(w, 0.1); err != nil {
		t.Fatalf("update engine audio attack: %v", err)
	}

	engineAudio, ok := w.EngineAudio(tank)
	if !ok {
		t.Fatalf("expected engine audio on tank")
	}
	if math.Abs(engineAudio.CurrentVolume-0.22) > 0.0001 {
		t.Fatalf("attack smoothing mismatch: got %f want %f", engineAudio.CurrentVolume, 0.22)
	}

	must(t, w.SetTankVelocity(tank, ecscomponents.TankVelocity{ForwardUnitsPerSecond: 0}))
	if err := UpdateEngineAudioVolumes(w, 0.1); err != nil {
		t.Fatalf("update engine audio release: %v", err)
	}

	engineAudio, _ = w.EngineAudio(tank)
	if math.Abs(engineAudio.CurrentVolume-0.16) > 0.0001 {
		t.Fatalf("release smoothing mismatch: got %f want %f", engineAudio.CurrentVolume, 0.16)
	}
}
