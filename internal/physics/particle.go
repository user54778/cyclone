package physics

import (
	"fmt"
	"math"

	"github.com/user54778/cyclone/internal/math64"
)

// Particle is the simplest object that can be simulated
// in a physics system. It is a point mass object; an object
// with mass, but no size, that CAN move through space but has NO
// internal degrees of freedom (can't rotate).
type Particle struct {
	// Describes linear position of the Particle in the world space.
	Position math64.Vector3
	// Hold linear velocity of the Particle in the world space.
	Velocity math64.Vector3
	// Hold acceleration of the Particle. This can be used to set
	// acceleration due to gravity, or any other constant acceleration.
	Acceleration math64.Vector3
	// Damping is our solution to give a rough approximation for drag
	// to apply to our particle in accordance with Newton's First Law.
	Damping float64
	// inverseMass is more useful to hold since it makes integration simpler
	// and is more useful to have objects with infinite mass (i.e., walls, floors, etc)
	// than storing mass itself, which could (although shouldn't) have zero mass.
	// This field is equipped with mass/inverseMass setters and private to help with this.
	inverseMass float64
}

// NewParticleMass creates a Particle object where the MASS itself is passed in as a parameter.
func NewParticleMass(position, velocity, acceleration math64.Vector3, damping, mass float64) Particle {
	p := Particle{
		Position:     position,
		Velocity:     velocity,
		Acceleration: acceleration,
		Damping:      damping,
	}
	p.SetMass(mass)

	return p
}

// NewParticleInverseMass creates a Particle object where the INVERSE MASS is passed in as a parameter.
func NewParticleInverseMass(position, velocity, acceleration math64.Vector3, damping, inverseMass float64) Particle {
	p := Particle{
		Position:     position,
		Velocity:     velocity,
		Acceleration: acceleration,
		Damping:      damping,
	}
	p.SetInverseMass(inverseMass)

	return p
}

// SetMass is a helper to set the particle's mass, and calculates its inverse mass.
// Zero or negative mass is treated as infinite.
func (p *Particle) SetMass(mass float64) {
	if mass <= 0 {
		p.inverseMass = 0.0
	} else {
		p.inverseMass = 1.0 / mass
	}
}

// SetInverseMass sets the inverseMass directly.
// Zero or negative inverse will be treated as infinite.
func (p *Particle) SetInverseMass(inverseMass float64) {
	if inverseMass <= 0 {
		p.inverseMass = 0.0
	} else {
		p.inverseMass = inverseMass
	}
}

// Mass is used to access the mass of the particle directly.
func (p *Particle) Mass() float64 {
	if p.inverseMass == 0 {
		return math.Inf(1) // NOTE: Mass is positively infinite.
	} else {
		return 1.0 / p.inverseMass
	}
}

// inverseMass accesses the inverseMass directly.
func (p *Particle) InverseMass() float64 {
	return p.inverseMass
}

// KineticEnergy returns the kinetic energy of a particle, given by the
// equation: K = 1/2m*mag(v)^2.
func (p *Particle) KineticEnergy() float64 {
	return 0.5 * p.Mass() * p.Velocity.Magnitude() * p.Velocity.Magnitude()
}

// Integrate performs Newton-Euler integration, a linear approximation
// to the correct integral, to integrate a Particle forward in time by duration amount.
func (p *Particle) Integrate(duration float64) error {
	// TODO: Modify to use a logging error system.
	switch {
	case p.inverseMass <= 0.0:
		return fmt.Errorf("integration is not performed on infinite mass")
	case duration < 0.0:
		return fmt.Errorf("can not perform integration on a negative duration")
	}

	// NOTE: That I am using pointer methods for Vector operations; copying will result
	// in thousands of vectors not used due to how often this function will be called.

	// Update position based on velocity
	// NOTE: Go will automatically dereference p since p.Position is
	// an addressable object.
	p.Position.ScaleAdd(p.Velocity, duration)

	// Update velocity based on acceleration
	resultingAcceleration := p.Acceleration

	// Update linear velocity from the acceleration.
	p.Velocity.ScaleAdd(resultingAcceleration, duration)

	// Impose drag. Match time scales by exponentiating time by drag, counteracting the effects
	// of the linearity of acceleration integration.
	// NOTE: Look at exponential decay for rate of change, and why it works here.
	dampingFactor := math.Pow(p.Damping, duration)
	p.Velocity.Scale(dampingFactor)

	return nil
}

// Used purely to measure difference when you don't incorporate time into drag.
func (p *Particle) IntegrateNoTimeScale(duration float64) error {
	switch {
	case p.inverseMass <= 0.0:
		return fmt.Errorf("integration is not performed on infinite mass")
	case duration < 0.0:
		return fmt.Errorf("can not perform integration on a negative duration")
	}

	p.Position.ScaleAdd(p.Velocity, duration)

	resultingAcceleration := p.Acceleration
	p.Velocity.ScaleAdd(resultingAcceleration, duration)
	p.Velocity.Scale(p.Damping)
	return nil
}
