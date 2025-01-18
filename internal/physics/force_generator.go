package physics

import (
	"math"

	"github.com/user54778/cyclone/internal/math64"
)

// TODO: Add tests for add and removal of forces to the registry.

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

// GravityGenerator represents a gravity force generator, which only requires the vector in which the gravity acts in.
type GravityGenerator struct {
	Gravity math64.Vector3
}

func NewGravityGenerator(gravity math64.Vector3) *GravityGenerator {
	return &GravityGenerator{
		Gravity: gravity,
	}
}

func (g *GravityGenerator) UpdateForce(particle *Particle, _ float64) {
	if !particle.HasFiniteMass() {
		return
	}

	particle.AddForce(g.Gravity.ScaleCopy(particle.Mass()))
}

// DragGenerator is a model to represent a drag force applied to a point mass,
// where k1 and k2 are two constants that characterize how *strong* the drag force is,
// named drag coefficients.
type DragGenerator struct {
	K1 float64
	K2 float64
}

func NewDragGenerator(k1, k2 float64) *DragGenerator {
	return &DragGenerator{
		K1: k1,
		K2: k2,
	}
}

// UpdateForce updates the drag force on a particle.
//
// The force acts in the opposite direction to the *velocity* of the object, with a strength
// that depends on *both* the *speed* of the object and the speed squared.
// The k2 value will grow *faster* at higher speeds-this is why cars don't accelerate infinitely,
// as for every doubling of speed, the *drag* nearly *quadruples*.
func (d *DragGenerator) UpdateForce(particle *Particle, _ float64) {
	// F_drag = -norm(vel(particle))*(k1*norm(vel(particle)) + k2*norm(vel(particle))^2)
	force := particle.Velocity

	// Calculate the total drag coefficient
	dragCoeff := force.Magnitude()
	dragCoeff = d.K1*dragCoeff + d.K2*dragCoeff

	// Calculate the final force and apply it
	force = force.Normalize()
	force = force.ScaleCopy(dragCoeff)
	particle.AddForce(force)
}

// UpliftForceGenerator represents an uplift force on a particle. An uplift force is simply
// "any upward pressure applied to a structure (particle) that has the *potential* to raise it relative to its surroundings."
//
// Only the XZ plane is considered for uplift above a particular point.
type UpliftForceGenerator struct {
	Center math64.Vector3
	Radius float64 // Radius of circular area of effect on the X-Z plane
	Force  float64 // *Magnitude* of uplift force.
}

func NewUpliftForceGenerator(center math64.Vector3, radius, force float64) *UpliftForceGenerator {
	return &UpliftForceGenerator{
		Center: center,
		Radius: radius,
		Force:  force,
	}
}

// UpdateForce updates the UpliftForceGenerator force based on the distance of the particle to the origin of the xz plane.
func (u *UpliftForceGenerator) UpdateForce(particle *Particle, _ float64) {
	// When the force generator is asked to apply its force,
	// it should test the X-Z coordinate of the object against the origin. If this coordinate
	// is within *a given distance of the origin*, then the uplift should be applied.
	// Otherwise, there is no force.
	distance := u.calcDistance(particle)

	// Apply the force if the distance is within the radius of the xz-plane.
	if distance <= u.Radius {
		forceMag := u.Force * (1 - distance/u.Radius) // Represent the fall off of force based on distance from the center.

		uplift := math64.NewVector3(0, forceMag, 0) // Upward force to apply to the particle

		particle.AddForce(uplift)
	}
}

// calcDistance calculates the distance between the particle and the center point using only
// the XZ coordinate plane.
func (u *UpliftForceGenerator) calcDistance(particle *Particle) float64 {
	dx := particle.Position.X - u.Center.X
	dz := particle.Position.Z - u.Center.Z
	return math.Sqrt(dx*dx + dz*dz)
}

// AirBrakeForceGenerator represents an air brake, where the single parameter determines whether it should be
// "On" or "Off".
type AirBrakeForceGenerator struct {
	On          bool
	NormalDrag  *DragGenerator
	BrakingDrag *DragGenerator
}

func NewAirBrakeForceGenerator(normK1, normK2, brakeK1, brakeK2 float64) *AirBrakeForceGenerator {
	return &AirBrakeForceGenerator{
		On:          false,
		NormalDrag:  NewDragGenerator(normK1, normK2),
		BrakingDrag: NewDragGenerator(brakeK1, brakeK2),
	}
}

func (a *AirBrakeForceGenerator) UpdateForce(particle *Particle, duration float64) {
	// FIXME: not being added to the particle
	if a.On {
		// a.BrakingDrag.UpdateForce(particle, duration)
	} else {
		// a.NormalDrag.UpdateForce(particle, duration)
	}
}
