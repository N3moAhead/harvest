package util

import (
	"math/rand/v2"

	"github.com/N3moAhead/harvest/internal/config"
)

// clamp limits v to [min, max]
func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// function to get a random position within the given position range
func GetRandomPositionInView(posX, camY float64) (float64, float64) {
	x := posX + rand.Float64()*config.SCREEN_WIDTH
	y := camY + rand.Float64()*config.SCREEN_HEIGHT
	return x, y
}
