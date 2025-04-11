package component

import "math"

type Vector2D struct {
	X, Y float64
}

func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{v.X + other.X, v.Y + other.Y}
}

func (v Vector2D) Sub(other Vector2D) Vector2D {
	return Vector2D{v.X - other.X, v.Y - other.Y}
}

func (v Vector2D) Mul(scalar float64) Vector2D {
	return Vector2D{v.X * scalar, v.Y * scalar}
}

func (v Vector2D) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2D) Normalize() Vector2D {
	len := v.Len()
	if len == 0 {
		return Vector2D{0, 0}
	}
	return Vector2D{v.X / len, v.Y / len}
}

func NewVector2D(x, y float64) Vector2D {
	return Vector2D{x, y}
}
