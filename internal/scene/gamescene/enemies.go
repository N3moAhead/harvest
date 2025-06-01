package gamescene

import (
	"math"
	"math/rand"
	"time"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/util"
	"github.com/hajimehoshi/ebiten/v2"
)

func updateEnemies(g *GameScene, dt float64, elapsed float32) {
	// SPAWN ENEMIES, based on elapsed time
	elapsedMs := float64(time.Since(g.startTime).Milliseconds())
	elapsedSec := elapsedMs / 1000.0
	difficulty := 1.0 + math.Sqrt(elapsedSec)/10.0 // increase difficulty over time, use square root to make it slower at the beginning
	// difficulty := 1.0 + elapsedSec/60.0 // 60.0 seconds is too short

	// intervalSec := baseIntervalSec / difficulty
	// decrease spawning interval/duration based on difficulty/ time
	intervalSec := config.BASE_SPAWN_INTERVAL_SEC / difficulty                  // decrease spawning interval/duration based on difficulty/ time
	count := int(math.Ceil(float64(config.BASE_COUNT_PER_BATCH) * difficulty))  // increase count of batches based on difficulty (number of pools --> per pool multiple enemies)
	mixProgress := util.Clamp((elapsedSec-config.MIX_START_SEC)/10.0, 0.0, 1.0) // mix progress from 0 to 1, after 120 seconds it will be 1.0

	// fmt.Printf("Elapsed Time: %.2f seconds, Mix Progress: %.2f, Interval: %.2f seconds, Count: %d, Difficulty: %.2f \n", elapsedSec, mixProgress, intervalSec, count, difficulty)
	if time.Since(g.lastSpawnTime).Seconds() >= intervalSec {
		g.lastSpawnTime = time.Now()
		spawnBatch(g, count, mixProgress)
	}

	for _, e := range g.Enemies {
		// e.Update(g.Player, dt)
		wasAlive := e.IsAlive()
		e.Update(g.Player, dt)
		if wasAlive && !e.IsAlive() {
			// enemy just died: generate drops
			elapsedMinutes := elapsed / 60000.0 // convert milliseconds to minutes
			drops := e.TryDrop(elapsedMinutes)
			for i := range drops {
				g.items = append(g.items, &drops[i])
			}
		}
	}
}

func drawEnemies(g *GameScene, screen *ebiten.Image, mapOffsetX, mapOffsetY float64) {
	for _, e := range g.Enemies {
		if e.IsAlive() {
			e.Draw(screen, mapOffsetX, mapOffsetY)
		}
	}
}

func spawnBatch(g *GameScene, count int, mixProgress float64) {
	pool := make([]string, count)
	for i := range pool {
		pool[i] = enemy.RandomEnemyType().String()
		// maybe if only use `types`:
		// pool[i] = types[rand.Intn(len(types))]
	}

	// completly mix pool if mixProgress is high enough
	if mixProgress > 0.8 {
		rand.Shuffle(len(pool), func(i, j int) {
			pool[i], pool[j] = pool[j], pool[i]
		})
	}

	for _, t := range pool {
		// fmt.Printf("Spawning %s enemy\n", t)
		if mixProgress < 0.3 {
			mapOffsetX, mapOffsetY := g.World.GetCameraPosition()
			g.Enemies = append(g.Enemies, g.Spawner.SpawnRandomInView(t, mapOffsetX, mapOffsetY))
		} else {
			switch rand.Intn(3) {
			case 0:
				g.Enemies = append(g.Enemies, g.Spawner.SpawnCircle(t, g.Player, 150, 6)...)
			case 1:
				g.Enemies = append(g.Enemies, g.Spawner.SpawnZigZag(t, g.Player.Pos, 5, 50, 20)...)
			default:
				g.Enemies = append(g.Enemies, g.Spawner.SpawnLine(t, g.Player.Pos, 5, 40, 10)...)
			}
		}
	}
}

func initEnemySpawner() *world.EnemySpawner {
	s := world.NewEnemySpawner()

	// register enemy factories
	s.RegisterFactory(enemy.TypeCarrot.String(), func(pos component.Vector2D) enemy.EnemyInterface {
		return enemy.NewCarrotEnemy(pos)
	})
	s.RegisterFactory(enemy.TypePotato.String(), func(pos component.Vector2D) enemy.EnemyInterface {
		return enemy.NewPotatoEnemy(pos)
	})
	s.RegisterFactory(enemy.TypeOnion.String(), func(pos component.Vector2D) enemy.EnemyInterface {
		return enemy.NewOnionEnemy(pos)
	})
	s.RegisterFactory(enemy.TypeCabbage.String(), func(pos component.Vector2D) enemy.EnemyInterface {
		return enemy.NewCabbageEnemy(pos)
	})
	s.RegisterFactory(enemy.TypeLeek.String(), func(pos component.Vector2D) enemy.EnemyInterface {
		return enemy.NewLeekEnemy(pos)
	})

	return s
}
