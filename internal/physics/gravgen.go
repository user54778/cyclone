package physics

import "github.com/user54778/cyclone/internal/math64"

type GravityGenerator struct {
	Gravity math64.Vector3
}

func NewGravityGenerator(gravity math64.Vector3) *GravityGenerator {
	return &GravityGenerator{
		Gravity: gravity,
	}
}

func (g *GravityGenerator) UpdateForce(particle *Particle, duration float64) {
	if !particle.HasFiniteMass() {
		return
	}

	particle.AddForce(g.Gravity.ScaleCopy(particle.Mass()))
}
