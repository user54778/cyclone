package physics

/*
import (
	"math"

	"github.com/user54778/cyclone/internal/math64"
	phy "github.com/user54778/cyclone/internal/physics"
)

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

func (u *UpliftForceGenerator) UpdateForce(particle *phy.Particle, _ float64) {
	// When the force generator is asked to apply its force,
	// it should test the X-Z coordinate of the object against the origin. If this coordinate
	// is within *a given distance of the origin*, then the uplift should be applied.
	// Otherwise, there is no force.
	distance := u.calcDistance(particle)

	if distance <= u.Radius {
		forceMag := u.Force * (1 - distance/u.Radius) // Represent the fall off of force based on distance from the center.

		uplift := math64.NewVector3(0, forceMag, 0) // Upward force to apply to the particle

		particle.AddForce(uplift)
	}
}

// calcDistance calculates the distance between the particle and the center point using only
// the XZ coordinate plane.
func (u *UpliftForceGenerator) calcDistance(particle *phy.Particle) float64 {
	dx := particle.Position.X - u.Center.X
	dz := particle.Position.Z - u.Center.Z
	return math.Sqrt(dx*dx + dz*dz)
}
*/
