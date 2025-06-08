package gamescene

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/toast"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	waveIntervalSeconds = 20
	totalWaves          = 20
)

type WaveDefinition struct {
	EnemyTypes []enemy.EnemyType
	Count      int
}

func (g *GameScene) initializeWaves() {
	g.waveDefinitions = []WaveDefinition{
		// Wave 1
		{EnemyTypes: []enemy.EnemyType{enemy.TypeOnion}, Count: 140},
		// Wave 2
		{EnemyTypes: []enemy.EnemyType{enemy.TypeLeek}, Count: 145},
		// Wave 3
		{EnemyTypes: []enemy.EnemyType{enemy.TypeCarrot}, Count: 150},
		// Wave 4
		{EnemyTypes: []enemy.EnemyType{enemy.TypeRadish}, Count: 160},
		// Wave 5
		{EnemyTypes: []enemy.EnemyType{enemy.TypeCabbage}, Count: 170},
		// Wave 6
		{EnemyTypes: []enemy.EnemyType{enemy.TypePotato}, Count: 190},
		// Wave 7
		{EnemyTypes: []enemy.EnemyType{enemy.TypeOnion, enemy.TypeLeek}, Count: 200},
		// Wave 8
		{EnemyTypes: []enemy.EnemyType{enemy.TypeCarrot, enemy.TypeRadish}, Count: 210},
		// Wave 9
		{EnemyTypes: []enemy.EnemyType{enemy.TypeCabbage, enemy.TypePotato}, Count: 230},
		// Wave 10
		{EnemyTypes: []enemy.EnemyType{enemy.TypeOnion, enemy.TypeLeek, enemy.TypeCarrot}, Count: 260},
		// Wave 11
		{EnemyTypes: []enemy.EnemyType{enemy.TypeRadish, enemy.TypeCabbage, enemy.TypePotato}, Count: 300},
		// Wave 12
		{EnemyTypes: []enemy.EnemyType{enemy.TypeOnion, enemy.TypeCarrot, enemy.TypePotato}, Count: 330},
		// Wave 13
		{EnemyTypes: []enemy.EnemyType{enemy.TypeCabbage, enemy.TypeOnion, enemy.TypeLeek}, Count: 500},
		// Wave 14
		{EnemyTypes: []enemy.EnemyType{enemy.TypeCarrot, enemy.TypeRadish, enemy.TypePotato, enemy.TypeCabbage}, Count: 550},
		// Wave 15
		{EnemyTypes: []enemy.EnemyType{enemy.TypeOnion, enemy.TypeLeek, enemy.TypeRadish, enemy.TypePotato}, Count: 600},
		// Wave 16
		{EnemyTypes: []enemy.EnemyType{enemy.TypeCarrot, enemy.TypePotato, enemy.TypeCabbage, enemy.TypeOnion}, Count: 650},
		// Wave 17
		{EnemyTypes: []enemy.EnemyType{enemy.TypeLeek, enemy.TypeCarrot, enemy.TypePotato, enemy.TypeCabbage}, Count: 700},
		// Wave 18
		{EnemyTypes: []enemy.EnemyType{enemy.TypeOnion, enemy.TypeLeek, enemy.TypeCarrot, enemy.TypeRadish, enemy.TypePotato}, Count: 750},
		// Wave 19
		{EnemyTypes: []enemy.EnemyType{enemy.TypeCarrot, enemy.TypePotato, enemy.TypeCabbage, enemy.TypeOnion, enemy.TypeLeek}, Count: 800},
		// Wave 20
		{EnemyTypes: []enemy.EnemyType{enemy.TypeOnion, enemy.TypeLeek, enemy.TypeCarrot, enemy.TypeRadish, enemy.TypeCabbage, enemy.TypePotato}, Count: 1500},
	}
	g.currentWaveIndex = -1
	// g.lastWaveStartTime will be set when the first wave starts
}

func updateEnemies(g *GameScene, dt float64, elapsed float32) {
	if g.currentWaveIndex < totalWaves-1 {
		if g.currentWaveIndex == -1 || time.Since(g.lastWaveStartTime).Seconds() >= waveIntervalSeconds {
			g.currentWaveIndex++
			g.lastWaveStartTime = time.Now()
			spawnWaveEnemies(g)
			font, ok := assets.AssetStore.GetFont("2p")
			if ok {
				toast.AddCustomToast(fmt.Sprintf("Wave %d!", g.currentWaveIndex+1), font, 3*time.Second)
			}
		}
	} else {
		// TODO: Implement an endless mode thats total brutal chaos
	}

	resolveEnemyOverlaps(g)

	// Update existing enemies
	for i := len(g.Enemies) - 1; i >= 0; i-- {
		e := g.Enemies[i]
		wasAlive := e.IsAlive()
		e.Update(g.Player, dt)

		if wasAlive && !e.IsAlive() {
			elapsedMinutes := float64(elapsed) / 60000.0
			drops := e.TryDrop(float32(elapsedMinutes))
			// TODO each enemy should increase the score by a diffrent amount
			g.Score += 10
			for j := range drops {
				g.items = append(g.items, &drops[j])
			}
			// Remove dead enemy
			g.Enemies = append(g.Enemies[:i], g.Enemies[i+1:]...)
		}
	}
}

func resolveEnemyOverlaps(g *GameScene) {
	const SEPARATION_RADIUS float64 = config.ENEMY_SEPERATION_RADIUS
	const SEPARATION_RADIUS_SQ = SEPARATION_RADIUS * SEPARATION_RADIUS

	numEnemies := len(g.Enemies)
	if numEnemies < 2 {
		return
	}

	for i := 0; i < numEnemies; i++ {
		for j := i + 1; j < numEnemies; j++ {
			enemy1 := g.Enemies[i]
			enemy2 := g.Enemies[j]
			pos1 := enemy1.GetPosition()
			pos2 := enemy2.GetPosition()
			collisionVector := pos1.Sub(pos2)
			distSq := collisionVector.LengthSq()
			if distSq < SEPARATION_RADIUS_SQ && distSq > 0 {
				distance := math.Sqrt(distSq)
				overlap := SEPARATION_RADIUS - distance
				pushVector := collisionVector.Normalize().Mul(overlap * 0.5)
				enemy1.SetPosition(pos1.Add(pushVector))
				enemy2.SetPosition(pos2.Sub(pushVector))
			} else if distSq == 0 {
				nudge := component.NewVector2D(0.1, 0)
				enemy1.SetPosition(pos1.Add(nudge))
			}
		}
	}
}

func spawnWaveEnemies(g *GameScene) {
	if g.currentWaveIndex < 0 || g.currentWaveIndex >= len(g.waveDefinitions) {
		return
	}

	waveDef := g.waveDefinitions[g.currentWaveIndex]
	numEnemyTypesInWave := len(waveDef.EnemyTypes)
	if numEnemyTypesInWave == 0 {
		return
	}

	countPerType := waveDef.Count / numEnemyTypesInWave
	if countPerType == 0 && waveDef.Count > 0 {
		countPerType = 1
	}

	for _, enemyTypeEnum := range waveDef.EnemyTypes {
		enemyTypeStr := enemyTypeEnum.String()

		if g.Player == nil || g.World == nil {
			for i := 0; i < countPerType; i++ {
				g.Enemies = append(g.Enemies, g.Spawner.SpawnRandom(enemyTypeStr))
			}
			continue
		}

		spawnPatternChoice := rand.Intn(4)

		if g.currentWaveIndex < 3 {
			spawnPatternChoice = 0
		}

		switch spawnPatternChoice {
		case 0: // Spawn Random Outside of View
			for i := 0; i < countPerType; i++ {
				spawnPos := getOffscreenSpawnPosition(g.Player.Pos, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, 100.0)
				newEnemy := g.Spawner.Spawn(enemyTypeStr, spawnPos)
				if newEnemy != nil {
					g.Enemies = append(g.Enemies, newEnemy)
				}
			}
		case 1: // Spawn Circle
			if countPerType > 0 {
				g.Enemies = append(g.Enemies, g.Spawner.SpawnCircle(enemyTypeStr, g.Player, 500+rand.Float64()*100, countPerType)...)
			}
		case 2: // Spawn ZigZag
			if countPerType > 0 {
				enemiesSpawned := 0
				for enemiesSpawned < countPerType {
					numToSpawnThisFormation := config.ENEMY_PER_SUB_FORMATION
					if countPerType-enemiesSpawned < config.ENEMY_PER_SUB_FORMATION {
						numToSpawnThisFormation = countPerType - enemiesSpawned
					}
					if numToSpawnThisFormation <= 0 {
						break
					}

					startPos := component.NewVector2D(g.Player.Pos.X+float64(rand.Intn(1500)-750), g.Player.Pos.Y+float64(rand.Intn(1500)-750))

					g.Enemies = append(g.Enemies, g.Spawner.SpawnZigZag(enemyTypeStr, startPos, numToSpawnThisFormation, 30+rand.Float64()*20, 15+rand.Float64()*10)...)
					enemiesSpawned += numToSpawnThisFormation
				}
			}
		default:
			if countPerType > 0 {
				enemiesSpawned := 0
				for enemiesSpawned < countPerType {
					numToSpawnThisFormation := config.ENEMY_PER_SUB_FORMATION
					if countPerType-enemiesSpawned < config.ENEMY_PER_SUB_FORMATION {
						numToSpawnThisFormation = countPerType - enemiesSpawned
					}
					if numToSpawnThisFormation <= 0 {
						break
					}

					startPos := component.NewVector2D(g.Player.Pos.X+float64(rand.Intn(1500)-750), g.Player.Pos.Y+float64(rand.Intn(1500)-750))

					g.Enemies = append(g.Enemies, g.Spawner.SpawnLine(enemyTypeStr, startPos, numToSpawnThisFormation, 25+rand.Float64()*15, 5+rand.Float64()*5)...)
					enemiesSpawned += numToSpawnThisFormation
				}
			}
		}
	}
}

// Helper function to get a spawn position just off-screen
func getOffscreenSpawnPosition(playerPos component.Vector2D, screenWidth, screenHeight, buffer float64) component.Vector2D {
	screenLeftEdge := playerPos.X - screenWidth/2
	screenRightEdge := playerPos.X + screenWidth/2
	screenTopEdge := playerPos.Y - screenHeight/2
	screenBottomEdge := playerPos.Y + screenHeight/2

	side := rand.Intn(4)
	var x, y float64

	randomizedOffset := buffer + (rand.Float64() * buffer)

	switch side {
	case 0: // Top
		x = screenLeftEdge + rand.Float64()*screenWidth
		y = screenTopEdge - randomizedOffset
	case 1: // Bottom
		x = screenLeftEdge + rand.Float64()*screenWidth
		y = screenBottomEdge + randomizedOffset
	case 2: // Left
		x = screenLeftEdge - randomizedOffset
		y = screenTopEdge + rand.Float64()*screenHeight
	default: // Right
		x = screenRightEdge + randomizedOffset
		y = screenTopEdge + rand.Float64()*screenHeight
	}
	return component.NewVector2D(x, y)
}

func drawEnemies(g *GameScene, screen *ebiten.Image, mapOffsetX, mapOffsetY float64) {
	for _, e := range g.Enemies {
		if e.IsAlive() {
			e.Draw(screen, mapOffsetX, mapOffsetY)
		}
	}
}

func initEnemySpawner() *world.EnemySpawner {
	s := world.NewEnemySpawner()

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
	s.RegisterFactory(enemy.TypeRadish.String(), func(pos component.Vector2D) enemy.EnemyInterface {
		return enemy.NewRadishEnemy(pos)
	})

	return s
}
