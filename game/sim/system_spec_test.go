package sim

import (
	"reflect"
	"testing"
)

// assertSystemSpec also verifies that the spec is wired into the pipeline.
func assertSystemSpec(t *testing.T, got, want SystemSpec) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("system spec mismatch:\n got: %v\nwant: %v", got, want)
	}
	for _, spec := range DefaultPipelineSpecs() {
		if spec.Name == want.Name {
			if !reflect.DeepEqual(spec, want) {
				t.Fatalf("pipeline spec mismatch for %q", want.Name)
			}
			return
		}
	}
	t.Fatalf("system %q missing from pipeline", want.Name)
}
