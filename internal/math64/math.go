package math64

import "math"

const (
	Pi          = math.Pi
	degRadRatio = Pi / 180.0
	radDegRatio = 180.0 / Pi
)

// DegToRad converts a scalar from degrees to radians.
func DegToRad(degrees float64) float64 {
	return degrees * degRadRatio
}

// RadToDeg converts a scalar from radians to degrees.
func RadToDeg(radians float64) float64 {
	return radians * radDegRatio
}
