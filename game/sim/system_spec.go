package sim

import ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"

// SystemSpec declares the data read and written by a pipeline step.
type SystemSpec struct {
	Name   string
	Reads  []ecscomponents.Kind
	Writes []ecscomponents.Kind
}
