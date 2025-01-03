package physics

import "github.com/user54778/cyclone/internal/math64"

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
}
