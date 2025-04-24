package component

import (
	"fmt"
	"math"
)

// Vector2D represents a 2-dimensional vector with float64 components.
type Vector2D struct {
	X, Y float64
}

func (v Vector2D) String() string {
	return fmt.Sprintf("Vector2D{X: %f, Y: %f}", v.X, v.Y)
}

// NewVector2D constructs a new Vector2D with the given x and y components.
func NewVector2D(x, y float64) Vector2D {
	return Vector2D{x, y}
}

// Add returns the vector sum of v and other (v + other).
// It does not modify the original vector v.
func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{v.X + other.X, v.Y + other.Y}
}

// Sub returns the vector difference of v and other (v - other).
// It does not modify the original vector v.
func (v Vector2D) Sub(other Vector2D) Vector2D {
	return Vector2D{v.X - other.X, v.Y - other.Y}
}

// Mul returns the vector v scaled by the given scalar factor (v * scalar).
// It does not modify the original vector v.
func (v Vector2D) Mul(scalar float64) Vector2D {
	return Vector2D{v.X * scalar, v.Y * scalar}
}

// LengthSq returns the squared magnitude (length) of the vector (v.X*v.X + v.Y*v.Y).
// This is computationally cheaper than Len() as it avoids the square root calculation.
// Useful for comparing vector lengths.
func (v Vector2D) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Len returns the magnitude (length) of the vector (sqrt(v.X*v.X + v.Y*v.Y)).
// If you only need to compare lengths, use LengthSq() for better performance.
func (v Vector2D) Len() float64 {
	// Directly use LengthSq() for clarity and potential minor optimization reuse.
	return math.Sqrt(v.LengthSq())
}

// Dot returns the dot product (scalar product) of vectors v and other.
// The dot product is calculated as v.X*other.X + v.Y*other.Y.
// Useful for calculating angles between vectors or projecting one vector onto another.
func (v Vector2D) Dot(other Vector2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Normalize returns a unit vector (a vector with length 1) pointing in the same direction as v.
// If the original vector v has a length of 0, it returns a zero vector {0, 0}.
// It does not modify the original vector v.
func (v Vector2D) Normalize() Vector2D {
	lenSq := v.LengthSq()
	if lenSq == 0 {
		// Cannot normalize a zero-length vector, return zero vector.
		return Vector2D{0, 0}
	}
	len := math.Sqrt(lenSq) // Calculate length only if non-zero
	return Vector2D{v.X / len, v.Y / len}
}
