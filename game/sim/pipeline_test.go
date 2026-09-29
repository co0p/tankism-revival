package sim

import (
	"reflect"
	"testing"

	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func TestDefaultPipelineSpecsOrder(t *testing.T) {
	specs := DefaultPipelineSpecs()

	got := make([]string, len(specs))
	for i, spec := range specs {
		got[i] = spec.Name
	}

	want := []string{
		"ai_chaser_intent",
		"tank_movement_integration",
		"engine_audio_volume",
		"weapon_cooldowns_and_projectile_spawn",
		"sound_volume_by_proximity",
		"projectile_integration",
		"projectile_impacts",
		"tank_obstacle_collision_resolution",
		"lifetime_expiration",
		"player_clamp_to_map_bounds",
		"enemy_clamp_to_map_bounds",
		"camera_follow_player",
		"player_backlight_follow",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("pipeline order mismatch:\n got: %v\nwant: %v", got, want)
	}
}

func TestDefaultPipelineSpecsDeclareDependencies(t *testing.T) {
	specs := DefaultPipelineSpecs()
	if len(specs) == 0 {
		t.Fatalf("expected non-empty pipeline specs")
	}

	for _, spec := range specs {
		if spec.Name == "" {
			t.Fatalf("system spec name must be set")
		}
		if len(spec.Reads) == 0 {
			t.Fatalf("system %q must declare read dependencies", spec.Name)
		}
		if len(spec.Writes) == 0 {
			t.Fatalf("system %q must declare write dependencies", spec.Name)
		}
	}

	assertDependsOn := func(name string, reads []ecscomponents.Kind, writes []ecscomponents.Kind) {
		t.Helper()

		for _, spec := range specs {
			if spec.Name != name {
				continue
			}

			if !reflect.DeepEqual(spec.Reads, reads) {
				t.Fatalf("system %q read deps mismatch:\n got: %v\nwant: %v", name, spec.Reads, reads)
			}
			if !reflect.DeepEqual(spec.Writes, writes) {
				t.Fatalf("system %q write deps mismatch:\n got: %v\nwant: %v", name, spec.Writes, writes)
			}
			return
		}

		t.Fatalf("system %q not found", name)
	}

	assertDependsOn(
		"tank_movement_integration",
		[]ecscomponents.Kind{
			ecscomponents.KindTransform,
			ecscomponents.KindTankDrive,
			ecscomponents.KindTankAcceleration,
			ecscomponents.KindTankDeceleration,
			ecscomponents.KindTankVelocity,
			ecscomponents.KindTankControlIntent,
		},
		[]ecscomponents.Kind{
			ecscomponents.KindTransform,
			ecscomponents.KindTankVelocity,
		},
	)

	assertDependsOn(
		"engine_audio_volume",
		[]ecscomponents.Kind{
			ecscomponents.KindTankDrive,
			ecscomponents.KindTankVelocity,
			ecscomponents.KindEngineAudio,
		},
		[]ecscomponents.Kind{
			ecscomponents.KindEngineAudio,
		},
	)

	assertDependsOn(
		"weapon_cooldowns_and_projectile_spawn",
		[]ecscomponents.Kind{
			ecscomponents.KindTransform,
			ecscomponents.KindTankDrive,
			ecscomponents.KindTankControlIntent,
			ecscomponents.KindWeapon,
			ecscomponents.KindAudioPlayer,
		},
		[]ecscomponents.Kind{
			ecscomponents.KindWeapon,
			ecscomponents.KindTransform,
			ecscomponents.KindSpriteRenderer,
			ecscomponents.KindCircleCollider,
			ecscomponents.KindProjectile,
			ecscomponents.KindLifetime,
			ecscomponents.KindPointLight,
			ecscomponents.KindPlaySoundIntent,
		},
	)

	assertDependsOn(
		"sound_volume_by_proximity",
		[]ecscomponents.Kind{
			ecscomponents.KindPlaySoundIntent,
			ecscomponents.KindTransform,
			ecscomponents.KindSoundVolumeByDistance,
		},
		[]ecscomponents.Kind{
			ecscomponents.KindPlaySoundIntent,
		},
	)

	assertDependsOn(
		"projectile_impacts",
		[]ecscomponents.Kind{
			ecscomponents.KindProjectile,
			ecscomponents.KindTransform,
			ecscomponents.KindCircleCollider,
			ecscomponents.KindSolidObstacle,
			ecscomponents.KindTankDrive,
			ecscomponents.KindHealth,
		},
		[]ecscomponents.Kind{
			ecscomponents.KindHealth,
			ecscomponents.KindProjectile,
		},
	)

	assertDependsOn(
		"camera_follow_player",
		[]ecscomponents.Kind{
			ecscomponents.KindTransform,
			ecscomponents.KindMapBounds,
		},
		[]ecscomponents.Kind{
			ecscomponents.KindCamera,
		},
	)
}
