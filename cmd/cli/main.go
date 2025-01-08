package main

import (
	"fmt"

	"github.com/user54778/cyclone/internal/math64"
)

func main() {
	u := math64.NewVector3(0, 1, 1)
	v := math64.NewVector3(0, -1, 0)
	fmt.Println(u.Cross(v))

	/*
		m := physics.NewParticleMass(u, v, math64.NewVector3(0, 0, 0), 0.999, 2.0)
		fmt.Println(m)
		fmt.Println(m.Mass())
		fmt.Println(m.InverseMass())
	*/
	/*
		p := physics.Particle{
			Position:     u,
			Velocity:     v,
			Acceleration: math64.Vector3{X: 0, Y: -15, Z: 0},
			Damping:      0.999,
		}
		p.SetInverseMass(1.0)
		fmt.Printf("%#v\n", p)
	*/
	/*
		p := physics.NewParticleInverseMass(u, v, math64.NewVector3(0, -15.0, 0), 0.99, 1.0)
		fmt.Printf("%#v\n", p)

		err := p.Integrate(float64(0.016))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Position: %+v\n", p.Position)
		fmt.Printf("Velocity: %+v\n", p.Velocity)
	*/

	// 3.4
	/*
		p := math64.NewVector3(1.0, 2.0, 3.0)
		pv := math64.NewVector3(1.0, -1.0, 2.0)
		pa := math64.NewVector3(0, 1.0, -1.0)
		fmt.Println(p, pv, pa)

		pv = pv.ScaleCopy(5.0)
		pa = pa.ScaleCopy((math.Pow(5, 2.0)) / 2)

		p_prime := p.AddCopy(pv).AddCopy(pa)

		fmt.Println(p, pv, pa)
		fmt.Println(p_prime)
	*/
}
