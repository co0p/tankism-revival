package sim

import (
	"reflect"
	"testing"
)

func TestDefaultPipelineSpecsOrder(t *testing.T) {
	specs := DefaultPipelineSpecs()
	got := make([]string, len(specs))
	for i, spec := range specs {
		got[i] = spec.Name
	}
	want := []string{
		"ai_enemy_markov_astar_intent",
		"tank_movement_integration",
		"engine_audio_volume",
		"weapon_cooldowns_and_projectile_spawn",
		"sound_volume_by_proximity",
		"ground_impact_effects",
		"projectile_integration",
		"hit_flash_expiration",
		"spark_particle_integration",
		"projectile_impacts",
		"tank_destruction_effects",
		"tank_tank_collision_resolution",
		"tank_obstacle_collision_resolution",
		"lifetime_expiration",
		"player_clamp_to_map_bounds",
		"enemy_clamp_to_map_bounds",
		"tank_tracks",
		"tank_smoke_particles",
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
		t.Fatal("expected non-empty pipeline specs")
	}
	for _, spec := range specs {
		if spec.Name == "" {
			t.Fatal("system spec name must be set")
		}
		if len(spec.Reads) == 0 {
			t.Fatalf("system %q must declare read dependencies", spec.Name)
		}
		if len(spec.Writes) == 0 {
			t.Fatalf("system %q must declare write dependencies", spec.Name)
		}
	}
}

func TestDefaultPipelineSpecsReturnsCopies(t *testing.T) {
	specs := DefaultPipelineSpecs()
	specs[0].Name = "changed"
	specs[0].Reads[0] = "changed"
	specs[0].Writes[0] = "changed"
	if !reflect.DeepEqual(DefaultPipelineSpecs()[0], enemyAISpec()) {
		t.Fatal("spec query exposed mutable pipeline state")
	}
}

func TestDefaultPipelineSystemsHaveRunCallbacks(t *testing.T) {
	for _, system := range defaultPipeline {
		if system.run == nil {
			t.Fatalf("system %q has no run callback", system.spec.Name)
		}
	}
}
