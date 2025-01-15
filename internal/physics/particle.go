// Package physics provides various methods to perform point
// mass physics.
package physics

import (
	"math"

	"github.com/user54778/cyclone/internal/math64"
	"github.com/user54778/cyclone/internal/physicslog"
)

// Particle is the simplest object that can be simulated
// in a physics system.
//
// It is a point mass object-an object with mass, but no size, that CAN move through space but has NO
// internal degrees of freedom (can't rotate).
type Particle struct {
	// Describes linear position of the Particle in the world space.
	Position math64.Vector3
	// Hold linear velocity of the Particle in the world space.
	Velocity math64.Vector3
	// Hold acceleration of the Particle. This can be used to set
	// acceleration due to gravity, or any other *constant* acceleration.
	Acceleration math64.Vector3
	// Damping is our solution to give a rough approximation for drag
	// to apply to our particle in accordance with Newton's First Law.
	Damping float64
	// Inverse Mass is more useful to hold since it makes integration simpler
	// and is more useful to have objects with infinite mass (i.e., walls, floors, etc)
	// than storing mass itself, which could (although shouldn't) have zero mass.
	// You can set this field with the Mass and InverseMass setters or set it directly.
	inverseMass float64
	// forceAccumulator accumulates every force to be applied at the next
	// iteration *only*. It is zeroed at each integration step.
	forceAccumulator math64.Vector3
}

// NewParticleMass creates a Particle object where the *mass* itself is passed in as a parameter.
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

// NewParticleInverseMass creates a Particle object where the *inverse mass* is passed in as a parameter.
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

func (p *Particle) HasFiniteMass() bool {
	return p.inverseMass > 0.0
}

// KineticEnergy returns the kinetic energy of a particle, given by the
// equation: K = 1/2m*mag(v)^2.
func (p *Particle) KineticEnergy() float64 {
	return 0.5 * p.Mass() * p.Velocity.Magnitude() * p.Velocity.Magnitude()
}

// AddForce adds force to the particle to be applied at the next iteration.
func (p *Particle) AddForce(force math64.Vector3) {
	p.forceAccumulator.Add(force) // NOTE: This directly adds to the particle's ForceAccumulator,
	// see if doing this in a functional style does not eat much memory.
}

// ClearForces sets the forceAccumulator to the zero value for a math64.Vector3.
func (p *Particle) ClearForces() {
	p.forceAccumulator = math64.Vector3{}
}

// Integrate updates the position and velocity of a point mass using equations for constant
// acceleration.
func (p *Particle) Integrate(duration float64) error {
	// TODO: Modify to use physicslog.
	switch {
	case p.inverseMass <= 0.0:
		// return fmt.Errorf("integration is not performed on infinite mass")
		return newPhysicsError("integration is not performed on infinite mass")
	case duration <= 0.0:
		// return fmt.Errorf("can not perform integration on a negative duration")
		return newPhysicsError("can not perform integration on a negative duration")
	}
	// NOTE: I am using pointer methods for Vector operations; copying will result
	// in thousands of vectors not used due to how often this function will be called.

	// Update position based on velocity
	// NOTE: Go will automatically dereference p since p.Position is
	// an addressable object.
	p.Position.ScaleAdd(p.Velocity, duration)

	// Update velocity based on acceleration
	resultingAcceleration := p.Acceleration
	// Apply the force's accumulated to the resulting acceleration.
	resultingAcceleration.ScaleAdd(p.forceAccumulator, p.inverseMass)
	// Update linear velocity from the acceleration.
	p.Velocity.ScaleAdd(resultingAcceleration, duration)

	// Impose drag. Match time scales by exponentiating time by drag, counteracting the effects
	// of the linearity of acceleration integration.
	dampingFactor := math.Pow(p.Damping, duration)
	p.Velocity.Scale(dampingFactor)

	// Clear the accumulated force after applying it to the particle.
	p.ClearForces()

	return nil
}

// PhysicsError represents specific errors relevant to our physics engine.
type PhysicsError struct {
	Message string
}

// Error is used to implement the Error interface.
func (e *PhysicsError) Error() string {
	return e.Message
}

// newPhysicsError creates a physics error object and logs it using PhysicsLogger.
func newPhysicsError(message string) error {
	err := &PhysicsError{
		Message: message,
	}

	logger := physicslog.NewPhysicsLogger(physicslog.LevelError)
	logger.LogError(message)

	return err
}

// ErrorCode represents different specific error codes the particle can throw out.
/*
type ErrorCode int

// FIXME: Look into alternate ways to better supply all the errors in the engine.
const (
	ErrInfiniteMass ErrorCode = iota
	ErrNegativeDuration
)

func (e ErrorCode) String() string {
	switch e {
	case ErrInfiniteMass:
		return "ErrInfiniteMass:"
	case ErrNegativeDuration:
		return "ErrNegativeDuration:"
	default:
		return ""
	}
}
*/

/*
// Deprecated: Only use Integrate() to perform integration. This should only ever be used
// to compare differences in velocity of the two functions when time is not incorporated in drag.
func (p *Particle) IntegrateNoTimeScale(duration float64) error {
	switch {
	case p.InverseMass <= 0.0:
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
*/
