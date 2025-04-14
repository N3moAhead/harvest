package world

import (
	"math"
	"math/rand"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/enemy"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/pkg/config"
)

type EnemySpawner struct{}

func NewEnemySpawner() *EnemySpawner {
	return &EnemySpawner{}
}

// TODO maybe Factory Method for different enemies?
func (s *EnemySpawner) SpawnCarrotEnemy(pos component.Vector2D) enemy.EnemyInterface {
	return enemy.NewCarrotEnemy(pos)
}

func (s *EnemySpawner) SpawnCarrotEnemyRandom() enemy.EnemyInterface {
	pos := component.NewVector2D(rand.Float64()*config.SCREEN_WIDTH, rand.Float64()*config.SCREEN_HEIGHT)
	return s.SpawnCarrotEnemy(pos)
}

// === Spawn Patterns ===

// 1. Circle Pattern (um Spieler herum)
func (s *EnemySpawner) SpawnCircle(p *player.Player, radius float64, count int) []*enemy.CarrotEnemy {
	positions := make([]component.Vector2D, count)
	angleStep := 2 * math.Pi / float64(count)

	for i := 0; i < count; i++ {
		angle := float64(i) * angleStep
		x := p.Pos.X + math.Cos(angle)*radius
		y := p.Pos.Y + math.Sin(angle)*radius
		positions[i] = component.NewVector2D(x, y)
	}
	return s.spawnAtPositions(positions)
}

// 2. ZigZag Pattern
func (s *EnemySpawner) SpawnZigZag(start component.Vector2D, count int, stepX, stepY float64) []*enemy.CarrotEnemy {
	positions := make([]component.Vector2D, count)

	for i := 0; i < count; i++ {
		offsetX := float64(i) * stepX
		offsetY := float64(i%2) * stepY
		positions[i] = component.NewVector2D(start.X+offsetX, start.Y+offsetY)
	}
	return s.spawnAtPositions(positions)
}

// 3. Line Pattern
func (s *EnemySpawner) SpawnLine(start component.Vector2D, count int, stepX, stepY float64) []*enemy.CarrotEnemy {
	positions := make([]component.Vector2D, count)

	for i := 0; i < count; i++ {
		offsetX := float64(i) * stepX
		offsetY := float64(i) * stepY
		positions[i] = component.NewVector2D(start.X+offsetX, start.Y+offsetY)
	}
	return s.spawnAtPositions(positions)
}

// 4. Random Pattern
func (s *EnemySpawner) SpawnRandom(count int) []*enemy.CarrotEnemy {
	positions := make([]component.Vector2D, count)

	for i := 0; i < count; i++ {
		x := rand.Float64() * config.SCREEN_WIDTH
		y := rand.Float64() * config.SCREEN_HEIGHT
		positions[i] = component.NewVector2D(x, y)
	}
	return s.spawnAtPositions(positions)
}

// === Gemeinsame Spawnhilfe ===

func (s *EnemySpawner) spawnAtPositions(positions []component.Vector2D) []*enemy.CarrotEnemy {
	enemies := make([]*enemy.CarrotEnemy, len(positions))

	for i, pos := range positions {
		enemies[i] = s.SpawnCarrotEnemy(pos).(*enemy.CarrotEnemy)
	}
	return enemies
}
