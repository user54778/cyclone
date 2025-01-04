package physics

import (
	"math"

	"github.com/user54778/cyclone/internal/math64"
)

// Particle is the simplest object that can be simulated
// in a physics system.
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
	// This field is set with mass/inverseMass setters and private to help with this.
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
func (p Particle) Mass() float64 {
	if p.inverseMass == 0 {
		return math.Inf(1) // NOTE: Mass is positively infinite.
	} else {
		return 1.0 / p.inverseMass
	}
}

// inverseMass accesses the inverseMass directly.
func (p Particle) InverseMass() float64 {
	return p.inverseMass
}
