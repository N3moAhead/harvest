package weapon

import (
	"math"

	"github.com/N3moAhead/harvest/internal/component"
)

func calculateSpreadDirections(baseDirection component.Vector2D, numProjectiles int, totalSpreadAngleDegrees float64) []component.Vector2D {
	if numProjectiles <= 0 {
		return []component.Vector2D{}
	}

	normalizedBaseDir := baseDirection.Normalize()
	directions := make([]component.Vector2D, numProjectiles)

	if numProjectiles == 1 {
		directions[0] = normalizedBaseDir
		return directions
	}

	totalSpreadAngleRad := totalSpreadAngleDegrees * math.Pi / 180.0
	angleStepRad := totalSpreadAngleRad / float64(numProjectiles-1)

	currentAngleRad := -totalSpreadAngleRad / 2.0

	for i := 0; i < numProjectiles; i++ {
		directions[i] = normalizedBaseDir.Rotate(currentAngleRad)
		currentAngleRad += angleStepRad
	}

	return directions
}
