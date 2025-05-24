package world

import (
	"math"
	"math/rand"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/N3moAhead/harvest/pkg/config"
)

// Factory for different enemies
type EnemyFactoryFunc func(pos component.Vector2D) enemy.EnemyInterface

type EnemySpawner struct {
	factories map[string]EnemyFactoryFunc
}

func NewEnemySpawner() *EnemySpawner {
	return &EnemySpawner{
		factories: make(map[string]EnemyFactoryFunc),
	}
}

func (s *EnemySpawner) RegisterFactory(enemyType string, factory EnemyFactoryFunc) {
	s.factories[enemyType] = factory
}

func (s *EnemySpawner) Spawn(enemyType string, pos component.Vector2D) enemy.EnemyInterface {
	if factory, ok := s.factories[enemyType]; ok {
		return factory(pos)
	}
	return nil
}

// Random Position
func (s *EnemySpawner) SpawnRandom(enemyType string) enemy.EnemyInterface {
	pos := component.NewVector2D(
		rand.Float64()*config.SCREEN_WIDTH,
		rand.Float64()*config.SCREEN_HEIGHT,
	)
	return s.Spawn(enemyType, pos)
}

// Spawn at multiple positions
func (s *EnemySpawner) SpawnAtPositions(enemyType string, positions []component.Vector2D) []enemy.EnemyInterface {
	enemies := make([]enemy.EnemyInterface, 0, len(positions))
	for _, pos := range positions {
		e := s.Spawn(enemyType, pos)
		if e != nil {
			enemies = append(enemies, e)
		}
	}
	return enemies
}

// === Spawn Patterns ===

// 1. Circle Pattern
func (s *EnemySpawner) SpawnCircle(enemyType string, p *player.Player, radius float64, count int) []enemy.EnemyInterface {
	positions := make([]component.Vector2D, count)
	angleStep := 2 * math.Pi / float64(count)
	for i := 0; i < count; i++ {
		angle := float64(i) * angleStep
		x := p.Pos.X + math.Cos(angle)*radius
		y := p.Pos.Y + math.Sin(angle)*radius
		positions[i] = component.NewVector2D(x, y)
	}
	return s.SpawnAtPositions(enemyType, positions)
}

// 2. ZigZag Pattern
func (s *EnemySpawner) SpawnZigZag(enemyType string, start component.Vector2D, count int, stepX, stepY float64) []enemy.EnemyInterface {
	positions := make([]component.Vector2D, count)
	for i := 0; i < count; i++ {
		offsetY := float64(i%2) * stepY
		positions[i] = component.NewVector2D(start.X+stepX*float64(i), start.Y+offsetY)
	}
	return s.SpawnAtPositions(enemyType, positions)
}

// 3. Line Pattern
func (s *EnemySpawner) SpawnLine(enemyType string, start component.Vector2D, count int, stepX, stepY float64) []enemy.EnemyInterface {
	positions := make([]component.Vector2D, count)
	for i := 0; i < count; i++ {
		positions[i] = component.NewVector2D(start.X+stepX*float64(i), start.Y+stepY*float64(i))
	}
	return s.SpawnAtPositions(enemyType, positions)
}

// 4. Random Pattern
func (s *EnemySpawner) SpawnMoreRandom(count int, enemyType string) []enemy.EnemyInterface {
	positions := make([]component.Vector2D, count)

	for i := 0; i < count; i++ {
		x := rand.Float64() * config.SCREEN_WIDTH
		y := rand.Float64() * config.SCREEN_HEIGHT
		positions[i] = component.NewVector2D(x, y)
	}
	return s.SpawnAtPositions(enemyType, positions)
}
