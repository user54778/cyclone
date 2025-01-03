package math64

import (
	"fmt"
	"math"
)

// Vector3 represents a vector in the 3D cartesian vector space.
type Vector3 struct {
	X, Y, Z float64
}

// NewVector3 creates a Vector3 with given parameters.
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{
		X: x,
		Y: y,
		Z: z,
	}
}

// NewZeroVector3 creates and returns a zeroed Vector3.
func NewZeroVector3() Vector3 {
	return Vector3{
		X: 0,
		Y: 0,
		Z: 0,
	}
}

// Multiplies a Vector3 by a scalar k.
func (v *Vector3) Scale(k float64) {
	v.X *= k
	v.Y *= k
	v.Z *= k
}

// Returns a copy of the vector scaled by k.
func (v Vector3) ScaleCopy(k float64) Vector3 {
	return Vector3{
		X: v.X * k,
		Y: v.Y * k,
		Z: v.Z * k,
	}
}

// Add directly adds the components of s to v.
func (v *Vector3) Add(s Vector3) {
	v.X += s.X
	v.Y += s.Y
	v.Z += s.Z
}

// AddCopy returns a Vector3 of components of s added to v.
func (v Vector3) AddCopy(s Vector3) Vector3 {
	return Vector3{
		X: v.X + s.X,
		Y: v.Y + s.Y,
		Z: v.Z + s.Z,
	}
}

// Sub directly adds the components of s to v.
func (v *Vector3) Sub(s Vector3) {
	v.X -= s.X
	v.Y -= s.Y
	v.Z -= s.Z
}

// SubCopy returns a Vector3 of components of s added to v.
func (v Vector3) SubCopy(s Vector3) Vector3 {
	return Vector3{
		X: v.X - s.X,
		Y: v.Y - s.Y,
		Z: v.Z - s.Z,
	}
}

// Invert flips all components of a Vector3.
func (v Vector3) Invert() Vector3 {
	return Vector3{
		X: -v.X,
		Y: -v.Y,
		Z: -v.Z,
	}
}

// AddScaledVector adds the components of s to v, scaled by k.
func (v *Vector3) AddScaledVector(s Vector3, k float64) {
	v.X += s.X * k
	v.Y += s.Y * k
	v.Z += s.Z * k
}

// ComponentProduct modifies v.
func (v *Vector3) Component(s Vector3) {
	v.X *= s.X
	v.Y *= s.Y
	v.Z *= s.Z
}

func (v Vector3) ComponentCopy(s Vector3) Vector3 {
	return Vector3{
		X: v.X * s.X,
		Y: v.Y * s.Y,
		Z: v.Z * s.Z,
	}
}

func (v Vector3) Dot(s Vector3) float64 {
	return v.X*s.X + v.Y*s.Y + v.Z*s.Z
}

func (v Vector3) Cross(s Vector3) Vector3 {
	epsilon := 1e-9
	cross := Vector3{
		X: v.Y*s.Z - v.Z*s.Y,
		Y: v.Z*s.X - v.X*s.Z,
		Z: v.X*s.Y - v.Y*s.X,
	}

	// Avoid floating-point precision errors.
	if math.Abs(cross.X) < epsilon {
		cross.X = 0.0
	}
	if math.Abs(cross.Y) < epsilon {
		cross.Y = 0.0
	}
	if math.Abs(cross.Z) < epsilon {
		cross.Z = 0.0
	}

	return cross
}

// Magnitude computes the magnitude of a Vector3.
func (v Vector3) Magnitude() float64 {
	return math.Sqrt(v.lengthSquared())
}

// lengthSquared computes the squared magnitude of a Vector3.
func (v Vector3) lengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Normalize resizes a Vector3 with a unit length of 1, i.e., turns it into a unit vector.
func (v *Vector3) Normalize() {
	n := v.Magnitude()
	if n > 0 {
		v.X *= 1.0 / n
		v.Y *= 1.0 / n
		v.Z *= 1.0 / n
	}
}

// makeOrthonormalBasis offers a primitive orthogonalization algorithm for three vectors.
func makeOrthonormalBasis(a, b, c *Vector3) error {
	a.Normalize()
	*c = a.Cross(*b)

	// NOTE: a and b must NOT be parallel
	if c.lengthSquared() == 0 {
		return fmt.Errorf("a and b are parallel")
	}

	c.Normalize()
	// ensure b is orthogonal to a, as a and c are already orthogonal
	// and normalized.
	*b = c.Cross(*a)
	return nil
}
