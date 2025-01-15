package physics

// ForceGenerator defines an interface for objects that can apply forces to one or more particles.
//
// If a force generator wants to apply a force, it can call the addForce method to the object that is
// passed in.
type ForceGenerator interface {
	// UpdateForce takes in the duration of the frame for which the force is needed, and a
	// pointer to the particle requesting the force.
	UpdateForce(particle *Particle, duration float64)
}

// ForceRegistry acts as a central registry of particles and force generators, holding
// a registry type in a slice.
//
// We need to also be able to register *which* force generators affect which particles. Decouple
// from the Particle object itself and have a central registry of particles and force generators.
type ForceRegistry struct {
	registrations []registry
}

// registry keeps track of one forceGenerator and the particle it is applied to.
type registry struct {
	particle *Particle
	fg       ForceGenerator
}

// AddForce registers the given force generator to apply to the given particle.
// It directly modifies the particle and fg (forceGenerator) passed in.
func (r *ForceRegistry) AddForce(particle *Particle, fg ForceGenerator) {
	r.registrations = append(r.registrations, registry{
		particle: particle,
		fg:       fg,
	})
}

// RemoveForce removes a given registered pair from the registry. If the pair is *not*
// registered, this method will do nothing.
func (r *ForceRegistry) RemoveForce(particle *Particle, fg ForceGenerator) {
	for i, reg := range r.registrations {
		if reg.particle == particle && reg.fg == fg {
			r.registrations = removeCopy(r.registrations, i)
			return // NOTE: Exit after removing the first match. This is very important to do here.
		}
	}
}

// Clear removes all force generator registrations from the registry *slice*,
// however, does *not* remove the particles or force generators themselves.
func (r *ForceRegistry) Clear() {
	r.registrations = nil
}

// UpdateForces calls all the force generators to update the forces of their
// corresponding particles.
func (r *ForceRegistry) UpdateForces(duration float64) {
	for _, reg := range r.registrations {
		reg.fg.UpdateForce(reg.particle, duration) // Notice how it calls the Interface function? Neat.
	}
}

// removeCopy is a helper function to remove an element from the underlying registry
// slice.
func removeCopy(registry []registry, i int) []registry {
	copy(registry[i:], registry[i+1:])
	return registry[:len(registry)-1]
}
