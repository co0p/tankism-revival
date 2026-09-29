package ecs

import (
	"slices"

	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
)

func (w *World) Renderables() []Renderable {
	renderables := make([]Renderable, 0, len(w.sprites))

	for entity, sprite := range w.sprites {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		renderables = append(renderables, Renderable{
			Entity:    entity,
			Transform: transform,
			Sprite:    sprite,
		})
	}

	slices.SortFunc(renderables, func(a, b Renderable) int {
		if a.Sprite.ZIndex < b.Sprite.ZIndex {
			return -1
		}
		if a.Sprite.ZIndex > b.Sprite.ZIndex {
			return 1
		}

		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return renderables
}

func (w *World) TankMovers() []TankMover {
	movers := make([]TankMover, 0, len(w.drives))

	for entity, drive := range w.drives {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		acceleration := w.acceleration[entity]
		deceleration := w.deceleration[entity]
		velocity := w.velocity[entity]
		intent := w.intents[entity]

		movers = append(movers, TankMover{
			Entity:       entity,
			Transform:    transform,
			Drive:        drive,
			Acceleration: acceleration,
			Deceleration: deceleration,
			Velocity:     velocity,
			Intent:       intent,
		})
	}

	slices.SortFunc(movers, func(a, b TankMover) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return movers
}

func (w *World) TankColliderBodies() []ColliderBody {
	bodies := make([]ColliderBody, 0, len(w.drives))

	for entity := range w.drives {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		collider, hasCollider := w.colliders[entity]
		if !hasCollider {
			continue
		}

		bodies = append(bodies, ColliderBody{
			Entity:    entity,
			Transform: transform,
			Collider:  collider,
		})
	}

	slices.SortFunc(bodies, func(a, b ColliderBody) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return bodies
}

func (w *World) SolidObstacleColliderBodies() []ColliderBody {
	bodies := make([]ColliderBody, 0, len(w.solid))

	for entity := range w.solid {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		collider, hasCollider := w.colliders[entity]
		if !hasCollider {
			continue
		}

		bodies = append(bodies, ColliderBody{
			Entity:    entity,
			Transform: transform,
			Collider:  collider,
		})
	}

	slices.SortFunc(bodies, func(a, b ColliderBody) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return bodies
}

func (w *World) AmbientLights() []ecscomponents.AmbientLight {
	lights := make([]ecscomponents.AmbientLight, 0, len(w.ambient))

	for _, light := range w.ambient {
		lights = append(lights, light)
	}

	return lights
}

func (w *World) PointLightSources() []PointLightSource {
	sources := make([]PointLightSource, 0, len(w.pointLights))

	for entity, light := range w.pointLights {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		sources = append(sources, PointLightSource{
			Entity:    entity,
			Transform: transform,
			Light:     light,
		})
	}

	slices.SortFunc(sources, func(a, b PointLightSource) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return sources
}

func (w *World) ArmedTanks() []ArmedTank {
	armed := make([]ArmedTank, 0, len(w.weapons))

	for entity, weapon := range w.weapons {
		_, hasDrive := w.drives[entity]
		if !hasDrive {
			continue
		}

		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		intent := w.intents[entity]

		armed = append(armed, ArmedTank{
			Entity:    entity,
			Transform: transform,
			Intent:    intent,
			Weapon:    weapon,
		})
	}

	slices.SortFunc(armed, func(a, b ArmedTank) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return armed
}

func (w *World) ProjectileStates() []ProjectileState {
	states := make([]ProjectileState, 0, len(w.projectiles))

	for entity, projectile := range w.projectiles {
		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		collider, hasCollider := w.colliders[entity]
		if !hasCollider {
			continue
		}

		states = append(states, ProjectileState{
			Entity:     entity,
			Transform:  transform,
			Projectile: projectile,
			Collider:   collider,
		})
	}

	slices.SortFunc(states, func(a, b ProjectileState) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return states
}

func (w *World) LifetimeStates() map[Entity]ecscomponents.Lifetime {
	states := make(map[Entity]ecscomponents.Lifetime, len(w.lifetimes))
	for entity, lifetime := range w.lifetimes {
		states[entity] = lifetime
	}

	return states
}

func (w *World) TankHealthStates() []TankHealthState {
	states := make([]TankHealthState, 0, len(w.health))

	for entity := range w.drives {
		health, hasHealth := w.health[entity]
		if !hasHealth {
			continue
		}

		transform, hasTransform := w.transforms[entity]
		if !hasTransform {
			continue
		}

		states = append(states, TankHealthState{
			Entity:    entity,
			Transform: transform,
			Health:    health,
		})
	}

	slices.SortFunc(states, func(a, b TankHealthState) int {
		if a.Entity < b.Entity {
			return -1
		}
		if a.Entity > b.Entity {
			return 1
		}

		return 0
	})

	return states
}
