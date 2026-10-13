package sim

import (
	"math"
	"slices"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"github.com/co0p/tankism-revival/lib/mathutil"
)

const (
	MaxTrackMarks = 512
	TrackSpacing  = 4.0
	TrackLifetime = 14.0
)

type TrackMark struct {
	X, Y, RotationRadians float64
	Age                   float64
}

type trackEmitter struct {
	previous ecscomponents.Transform
	distance float64
}

// Tracks stores ground marks separately from gameplay entities.
type Tracks struct {
	marks    []TrackMark
	emitters map[ecs.Entity]trackEmitter
}

func NewTracks() *Tracks {
	return &Tracks{marks: make([]TrackMark, 0, MaxTrackMarks), emitters: make(map[ecs.Entity]trackEmitter)}
}

func (t *Tracks) Marks() []TrackMark { return slices.Clone(t.marks) }

func (t *Tracks) Advance(deltaSeconds float64) {
	alive := t.marks[:0]
	for _, mark := range t.marks {
		mark.Age += math.Max(0, deltaSeconds)
		if mark.Age < TrackLifetime {
			alive = append(alive, mark)
		}
	}
	t.marks = alive
}

func (t *Tracks) add(mark TrackMark) {
	if len(t.marks) == MaxTrackMarks {
		copy(t.marks, t.marks[1:])
		t.marks = t.marks[:len(t.marks)-1]
	}
	t.marks = append(t.marks, mark)
}

// Update runs after collisions and clamping. Stationary turns leave no marks.
func (t *Tracks) Update(world *ecs.World, deltaSeconds float64) {
	t.Advance(deltaSeconds)
	active := make(map[ecs.Entity]bool)
	for _, tank := range world.TankMovers() {
		if health, ok := world.Health(tank.Entity); ok && (health.Current <= 0 || health.Max <= 0) {
			continue
		}
		active[tank.Entity] = true
		emitter, found := t.emitters[tank.Entity]
		current := tank.Transform
		if found {
			dx, dy := current.X-emitter.previous.X, current.Y-emitter.previous.Y
			distance := math.Hypot(dx, dy)
			forwardX, forwardY := forwardVector(current.RotationRadians)
			// Ignore teleports and perpendicular collision corrections.
			if distance > 0 && distance <= 64 && math.Abs(dx*forwardX+dy*forwardY) >= distance*0.5 {
				rotationDelta := mathutil.NormalizeRadians(current.RotationRadians - emitter.previous.RotationRadians)
				for step := TrackSpacing - emitter.distance; step <= distance; step += TrackSpacing {
					fraction := step / distance
					t.add(TrackMark{X: emitter.previous.X + dx*fraction, Y: emitter.previous.Y + dy*fraction, RotationRadians: emitter.previous.RotationRadians + rotationDelta*fraction})
				}
				emitter.distance = math.Mod(emitter.distance+distance, TrackSpacing)
			} else if distance > 64 {
				emitter.distance = 0
			}
		}
		emitter.previous = current
		t.emitters[tank.Entity] = emitter
	}
	for entity := range t.emitters {
		if !active[entity] {
			delete(t.emitters, entity)
		}
	}
}
