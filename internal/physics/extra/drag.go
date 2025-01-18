package physics

/*
// DragGenerator is a model to represent a drag force applied to a point mass,
// where k1 and k2 are two constants that characterize how *strong* the drag force is,
// named drag coefficients.
//
// F_drag = -norm(vel(particle))*(k1*norm(vel(particle)) + k2*norm(vel(particle))^2)
// The force acts in the opposite direction to the *velocity* of the object, with a strength
// that depends on *both* the *speed* of the object and the speed squared.
// The k2 value will grow *faster* at higher speeds-this is why cars don't accelerate infinitely,
// as for every doubling of speed, the *drag* nearly *quadruples*.
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
*/
