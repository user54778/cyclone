package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/user54778/cyclone/internal/math64"
	"github.com/user54778/cyclone/internal/physics"
)

/*
* In the text we looked at two ways to represent damping: Equations 3.5 and 3.6.
* Implement a small test program that repeatedly simulates a pair of particles moving under
* gravity for a fixed duration (1 s, for example). One particle should use Equation 3.5
* and the other 3.6. Use random durations for the frame (within some small margin)
* to simulate a variable frame rate. How much difference, on average, is there between
* the velocities of the two particles at the end of each simulation?
 */

// Equation 3.5: p_vel_prime = p_vel*d + p_acc*t -> integrator w/out proportion velocity retained each second.
// Equation 3.6: Integrate().

// Repeatedly simulate pair of particles moving under gravity for 1 second.
func simulate(duration, minFrameTime, maxFrameTime float64, frameCount int, gravity math64.Vector3) float64 {
	// random object
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	p1 := physics.NewParticleMass(math64.NewVector3(1, 0, 0), math64.NewVector3(0, 1, 0),
		gravity, 0.99, 1.0)
	p2 := physics.NewParticleMass(math64.NewVector3(1, 0, 0), math64.NewVector3(0, 1, 0),
		gravity, 0.99, 1.0)

	elapsedTime := 0.0
	var totalDiff float64

	// Iterate for duration seconds
	for elapsedTime < duration {
		// Use random duration for the frame to simulate variable frame rate
		// This generates a random float between 0.0 and 1.0, scales it by the WIDTH of the range,
		// and shifts it to start from minFrameTime.
		frameTime := minFrameTime + r.Float64()*(maxFrameTime-minFrameTime)

		// Don't go over duration
		if elapsedTime+frameTime > duration {
			frameTime = duration - elapsedTime
		}

		// Integration
		err1 := p1.IntegrateNoTimeScale(frameTime)
		err2 := p2.Integrate(frameTime)
		if err1 != nil || err2 != nil {
			fmt.Println("error during integration:", err1, err2)
			continue
		}

		// NOTE: When variability of frame time is removed, there is still a velocity difference between the too.
		velocityDiff := p1.Velocity.SubCopy(p2.Velocity).Magnitude()
		// NOTE: When variability suddenly changes, we should see significant differences in velocity since the particle will experience more drag.
		totalDiff += velocityDiff

		// Iterate
		elapsedTime += frameTime
	}

	// Return the average difference between the velocities of the two particles.
	return totalDiff / float64(frameCount)
}

func main() {
	t := simulate(5.0, 0.01, 0.03, 165, math64.NewVector3(0, 15, 0))
	fmt.Printf("%#v\n", t)
}
