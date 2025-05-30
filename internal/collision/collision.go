package collision

import (
	"math"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
)

// findEnemiesInArc searches for enemies within an arc.
// center: The center of the attack (player position).
// radius: The range of the attack.
// direction: Normalized vector of the attack direction.
// angle: Total angle of the sector (e.g., math.Pi for 180 degrees).
// enemies: The list of enemies to check.
func FindEnemiesInArc(
	center component.Vector2D,
	radius float64,
	direction component.Vector2D,
	angle float64,
	enemies []enemy.EnemyInterface,
) []enemy.EnemyInterface {
	var foundEnemies []enemy.EnemyInterface
	halfAngle := angle / 2.0
	radiusSq := radius * radius

	if direction.LengthSq() == 0 {
		direction = component.Vector2D{X: 0, Y: -1} // Defaults to the upright position
	} else if direction.LengthSq() != 1.0 {
		dirNorm := direction.Normalize()
		direction = dirNorm
	}

	for _, enemy := range enemies {
		if enemy == nil || !enemy.IsAlive() { // Check if enemy is valid and alive
			continue
		}
		enemyPos := enemy.GetPosition()
		vecToEnemy := enemyPos.Sub(center) // Vector from center to enemy

		// 1. Distance check (squared, to avoid sqrt)
		distSq := vecToEnemy.LengthSq()
		if distSq > radiusSq || distSq == 0 {
			continue // Too far away or exactly at the center
		}

		vecToEnemyNorm := vecToEnemy.Normalize()

		// 2. Angle check
		// Calculate the cosine of the angle between the direction and the vector to the enemy
		// Dot product: a · b = |a| * |b| * cos(theta)
		// Since direction and vecToEnemyNorm are normalized (|a|=|b|=1), dot = cos(theta)
		dotProduct := direction.Dot(vecToEnemyNorm)

		// Acos gives the angle back. Check if it is within the allowed range.
		// Clamp the DotProduct to [-1, 1] due to possible float inaccuracies
		if dotProduct > 1.0 {
			dotProduct = 1.0
		}
		if dotProduct < -1.0 {
			dotProduct = -1.0
		}

		angleToEnemy := math.Acos(dotProduct)

		if angleToEnemy <= halfAngle {
			// Enemy is in sector
			foundEnemies = append(foundEnemies, enemy)
		}
	}
	return foundEnemies
}

func FindEnemiesInCircle(center component.Vector2D, radius float64, enemies []enemy.EnemyInterface) []enemy.EnemyInterface {
	var foundEnemies []enemy.EnemyInterface
	radiusSq := radius * radius

	for _, enemy := range enemies {
		distSq := center.Sub(enemy.GetPosition()).LengthSq()
		if distSq < radiusSq {
			foundEnemies = append(foundEnemies, enemy)
		}
	}

	return foundEnemies
}
