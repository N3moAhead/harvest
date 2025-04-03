package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth          = 800
	screenHeight         = 600
	playerSize           = 20
	playerSpeed          = 3.0
	playerMaxHP          = 100
	enemySize            = 15
	enemySpeedBase       = 10
	enemyHPBase          = 10
	enemyDamage          = 10
	projectileSize       = 8
	projectileSpeed      = 20.0
	projectileDamageBase = 10
	projectileLifespan   = 5   // seconds
	fireRateBase         = 0.1 // seconds between shots
	xpGemSize            = 6
	xpGemValue           = 15
	xpPickupRadius       = 80.0 // Radius around player to attract gems
	levelUpBaseXP        = 100
	levelUpFactor        = 1.5

	enemySpawnIntervalBase = 1.0  // seconds
	enemySpawnIncreaseRate = 0.98 // Multiplier each spawn cycle
	enemySpawnMaxRate      = 0.1  // Minimum spawn interval
)

// --- Vector Math Helper ---

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

// --- Game Objects ---

type Player struct {
	Pos           Vector2D
	HP            int
	MaxHP         int
	XP            int
	Level         int
	XPToNextLevel int
	Speed         float64
}

type Enemy struct {
	Pos    Vector2D
	HP     int
	MaxHP  int
	Speed  float64
	Damage int
}

type Projectile struct {
	Pos      Vector2D
	Vel      Vector2D
	Damage   int
	Lifespan float64 // Remaining time in seconds
}

type XPGem struct {
	Pos   Vector2D
	Value int
}

// --- Game State ---

type Game struct {
	player               Player
	enemies              []*Enemy
	projectiles          []*Projectile
	xpGems               []*XPGem
	rng                  *rand.Rand
	fireCooldown         float64 // Time until next shot can be fired
	enemySpawnTimer      float64 // Time until next enemy spawns
	currentSpawnInterval float64
	projectileDamage     int     // Current damage based on level
	fireRate             float64 // Current fire rate based on level
	gameOver             bool
	gameTime             float64 // Total time elapsed
}

// --- Game Initialization ---

func NewGame() *Game {
	player := Player{
		Pos:           Vector2D{X: screenWidth / 2, Y: screenHeight / 2},
		HP:            playerMaxHP,
		MaxHP:         playerMaxHP,
		XP:            0,
		Level:         1,
		XPToNextLevel: levelUpBaseXP,
		Speed:         playerSpeed,
	}

	g := &Game{
		player:               player,
		enemies:              []*Enemy{},
		projectiles:          []*Projectile{},
		xpGems:               []*XPGem{},
		rng:                  rand.New(rand.NewSource(time.Now().UnixNano())),
		fireCooldown:         0,
		enemySpawnTimer:      enemySpawnIntervalBase,
		currentSpawnInterval: enemySpawnIntervalBase,
		projectileDamage:     projectileDamageBase,
		fireRate:             fireRateBase,
		gameOver:             false,
		gameTime:             0,
	}
	return g
}

// --- Game Logic ---

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			// Restart game
			newG := NewGame()
			*g = *newG // Overwrite current game state
		}
		return nil
	}

	deltaTime := 1.0 / float64(ebiten.TPS())
	g.gameTime += deltaTime

	// --- Player Input & Movement ---
	moveDir := Vector2D{0, 0}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		moveDir.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		moveDir.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveDir.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveDir.X += 1
	}

	moveDir = moveDir.Normalize() // Prevent faster diagonal movement
	g.player.Pos = g.player.Pos.Add(moveDir.Mul(g.player.Speed))

	// Clamp player position to screen bounds
	g.player.Pos.X = math.Max(playerSize/2, math.Min(screenWidth-playerSize/2, g.player.Pos.X))
	g.player.Pos.Y = math.Max(playerSize/2, math.Min(screenHeight-playerSize/2, g.player.Pos.Y))

	// --- Automatic Firing ---
	g.fireCooldown -= deltaTime
	if g.fireCooldown <= 0 {
		g.fireProjectile()
		g.fireCooldown = g.fireRate // Use current fire rate
	}

	// --- Projectile Update ---
	nextProjectiles := g.projectiles[:0] // Re-use slice memory
	for _, p := range g.projectiles {
		p.Pos = p.Pos.Add(p.Vel.Mul(deltaTime))
		p.Lifespan -= deltaTime

		// Keep projectile if lifespan > 0 and within screen bounds (with some margin)
		if p.Lifespan > 0 && p.Pos.X > -projectileSize && p.Pos.X < screenWidth+projectileSize && p.Pos.Y > -projectileSize && p.Pos.Y < screenHeight+projectileSize {
			nextProjectiles = append(nextProjectiles, p)
		}
	}
	g.projectiles = nextProjectiles

	// --- Enemy Spawning ---
	g.enemySpawnTimer -= deltaTime
	if g.enemySpawnTimer <= 0 {
		g.spawnEnemy()
		// Decrease spawn interval over time, but not below max rate
		g.currentSpawnInterval *= enemySpawnIncreaseRate
		if g.currentSpawnInterval < enemySpawnMaxRate {
			g.currentSpawnInterval = enemySpawnMaxRate
		}
		g.enemySpawnTimer = g.currentSpawnInterval
	}

	// --- Enemy Update ---
	nextEnemies := g.enemies[:0] // Re-use slice memory
	for _, e := range g.enemies {
		// Move towards player
		dirToPlayer := g.player.Pos.Sub(e.Pos).Normalize()
		e.Pos = e.Pos.Add(dirToPlayer.Mul(e.Speed * deltaTime))

		// Check collision with player
		if checkCollision(g.player.Pos, playerSize, e.Pos, enemySize) {
			g.player.HP -= e.Damage
			if g.player.HP <= 0 {
				g.player.HP = 0
				g.gameOver = true
				// Don't process further updates if game over
				return nil
			}
			// Simple knockback (optional, can be glitchy without physics)
			// knockbackDir := e.Pos.Sub(g.player.Pos).Normalize()
			// e.Pos = e.Pos.Add(knockbackDir.Mul(5.0))

			// Mark enemy for removal after dealing damage (or let it persist)
			// For this simple version, we'll let it keep trying to damage
			// but typically you'd add a damage cooldown.
		}

		// Check collision with projectiles
		hit := false
		nextProjectiles = g.projectiles[:0] // Re-use slice memory for projectiles again
		for _, p := range g.projectiles {
			if !hit && checkCollision(e.Pos, enemySize, p.Pos, projectileSize) {
				e.HP -= p.Damage
				hit = true // Projectile hits only one enemy and disappears
				// Don't add this projectile to nextProjectiles
			} else {
				nextProjectiles = append(nextProjectiles, p) // Keep projectile if no hit
			}
		}
		g.projectiles = nextProjectiles // Update projectile list after checking this enemy

		if e.HP > 0 {
			nextEnemies = append(nextEnemies, e) // Keep enemy alive
		} else {
			// Enemy died - spawn XP gem
			g.spawnXPGem(e.Pos)
		}
	}
	g.enemies = nextEnemies

	// --- XP Gem Update ---
	nextXpGems := g.xpGems[:0] // Re-use slice memory
	for _, gem := range g.xpGems {
		distToPlayer := g.player.Pos.Sub(gem.Pos).Len()
		collected := false

		if distToPlayer < xpPickupRadius {
			// Move towards player if close
			dirToPlayer := g.player.Pos.Sub(gem.Pos).Normalize()
			gem.Pos = gem.Pos.Add(dirToPlayer.Mul(playerSpeed * 1.5 * deltaTime)) // Move faster than player

			// Check for actual collection (very close)
			if distToPlayer < (playerSize/2 + xpGemSize/2) {
				g.player.XP += gem.Value
				collected = true
				g.checkLevelUp()
			}
		}

		if !collected {
			nextXpGems = append(nextXpGems, gem)
		}
	}
	g.xpGems = nextXpGems

	return nil
}

// --- Helper Methods ---

func (g *Game) fireProjectile() {
	// Fire straight up for simplicity
	direction := Vector2D{X: 0, Y: -1}

	projectile := &Projectile{
		Pos:      g.player.Pos,
		Vel:      direction.Mul(projectileSpeed),
		Damage:   g.projectileDamage, // Use current damage
		Lifespan: projectileLifespan,
	}
	g.projectiles = append(g.projectiles, projectile)
}

func (g *Game) spawnEnemy() {
	// Spawn outside the screen bounds
	side := g.rng.Intn(4)
	pos := Vector2D{}
	margin := 50.0 // Spawn distance off-screen

	switch side {
	case 0: // Top
		pos.X = g.rng.Float64() * screenWidth
		pos.Y = -margin
	case 1: // Bottom
		pos.X = g.rng.Float64() * screenWidth
		pos.Y = screenHeight + margin
	case 2: // Left
		pos.X = -margin
		pos.Y = g.rng.Float64() * screenHeight
	case 3: // Right
		pos.X = screenWidth + margin
		pos.Y = g.rng.Float64() * screenHeight
	}

	// Increase enemy stats slightly based on game time/level (optional)
	timeFactor := 1.0 + (g.gameTime / 60.0) // Increase stats every minute
	enemyHP := int(float64(enemyHPBase) * timeFactor)
	enemySpeed := enemySpeedBase * timeFactor

	enemy := &Enemy{
		Pos:    pos,
		HP:     enemyHP,
		MaxHP:  enemyHP,
		Speed:  enemySpeed,
		Damage: enemyDamage,
	}
	g.enemies = append(g.enemies, enemy)
}

func (g *Game) spawnXPGem(pos Vector2D) {
	gem := &XPGem{
		Pos:   pos,
		Value: xpGemValue,
	}
	g.xpGems = append(g.xpGems, gem)
}

func (g *Game) checkLevelUp() {
	for g.player.XP >= g.player.XPToNextLevel {
		g.player.Level++
		g.player.XP -= g.player.XPToNextLevel
		g.player.XPToNextLevel = int(float64(g.player.XPToNextLevel) * levelUpFactor)

		// Apply level up bonuses (simple version)
		g.projectileDamage += 5 // Increase damage
		g.fireRate *= 0.95      // Increase fire speed (decrease interval)
		// Could also increase HP, speed, add new weapons etc. here
		g.player.HP = g.player.MaxHP // Heal on level up
		fmt.Printf("Level Up! Reached Level %d. Damage: %d, Fire Rate: %.2f\n", g.player.Level, g.projectileDamage, g.fireRate)
	}
}

func checkCollision(pos1 Vector2D, size1 float64, pos2 Vector2D, size2 float64) bool {
	// Simple Axis-Aligned Bounding Box (AABB) collision for squares/cubes centered at Pos
	halfSize1 := size1 / 2
	halfSize2 := size2 / 2

	return pos1.X-halfSize1 < pos2.X+halfSize2 &&
		pos1.X+halfSize1 > pos2.X-halfSize2 &&
		pos1.Y-halfSize1 < pos2.Y+halfSize2 &&
		pos1.Y+halfSize1 > pos2.Y-halfSize2
}

// --- Drawing ---

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255}) // Dark grey background

	// Draw XP Gems (draw first so they appear below others)
	for _, gem := range g.xpGems {
		drawSquare(screen, gem.Pos, xpGemSize, color.RGBA{255, 255, 0, 255}) // Yellow
	}

	// Draw Enemies
	for _, e := range g.enemies {
		// Optional: Draw HP bar above enemy
		hpRatio := float32(e.HP) / float32(e.MaxHP)
		vector.DrawFilledRect(screen, float32(e.Pos.X-enemySize/2), float32(e.Pos.Y-enemySize/2-5), float32(enemySize), 3, color.RGBA{255, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, float32(e.Pos.X-enemySize/2), float32(e.Pos.Y-enemySize/2-5), float32(enemySize)*hpRatio, 3, color.RGBA{0, 255, 0, 255}, false)
		drawSquare(screen, e.Pos, enemySize, color.RGBA{0, 0, 255, 255}) // Blue
	}

	// Draw Projectiles
	for _, p := range g.projectiles {
		drawSquare(screen, p.Pos, projectileSize, color.RGBA{0, 255, 0, 255}) // Green
	}

	// Draw Player
	drawSquare(screen, g.player.Pos, playerSize, color.RGBA{255, 0, 0, 255}) // Red

	// --- Draw UI ---
	// HP Bar
	hpRatio := float64(g.player.HP) / float64(g.player.MaxHP)
	hpBarWidth := 100.0
	hpBarHeight := 10.0
	vector.DrawFilledRect(screen, 10, 10, float32(hpBarWidth), float32(hpBarHeight), color.RGBA{100, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, 10, 10, float32(hpBarWidth*hpRatio), float32(hpBarHeight), color.RGBA{255, 0, 0, 255}, false)

	// XP Bar
	xpRatio := float64(g.player.XP) / float64(g.player.XPToNextLevel)
	xpBarWidth := 100.0
	xpBarHeight := 10.0
	vector.DrawFilledRect(screen, 10, float32(15+hpBarHeight), float32(xpBarWidth), float32(xpBarHeight), color.RGBA{100, 100, 0, 255}, false)
	vector.DrawFilledRect(screen, 10, float32(15+hpBarHeight), float32(xpBarWidth*xpRatio), float32(xpBarHeight), color.RGBA{255, 255, 0, 255}, false)

	// Text Info
	levelStr := fmt.Sprintf("Level: %d", g.player.Level)
	hpStr := fmt.Sprintf("HP: %d/%d", g.player.HP, g.player.MaxHP)
	xpStr := fmt.Sprintf("XP: %d/%d", g.player.XP, g.player.XPToNextLevel)
	enemyCountStr := fmt.Sprintf("Enemies: %d", len(g.enemies))
	timeStr := fmt.Sprintf("Time: %.1fs", g.gameTime)

	ebitenutil.DebugPrintAt(screen, hpStr, 120, 10)
	ebitenutil.DebugPrintAt(screen, levelStr, 10, 35)
	ebitenutil.DebugPrintAt(screen, xpStr, 10, 50)
	ebitenutil.DebugPrintAt(screen, enemyCountStr, 10, 65)
	ebitenutil.DebugPrintAt(screen, timeStr, 10, 80)

	// Game Over Message
	if g.gameOver {
		msg := "Game Over!\nPress 'R' to Restart"
		textX := (screenWidth - 10) / 2
		textY := (screenHeight - 10) / 2
		ebitenutil.DebugPrintAt(screen, msg, textX, textY)
	}
}

// Helper to draw a centered square
func drawSquare(screen *ebiten.Image, pos Vector2D, size float64, clr color.Color) {
	halfSize := float32(size / 2)
	vector.DrawFilledRect(screen, float32(pos.X)-halfSize, float32(pos.Y)-halfSize, float32(size), float32(size), clr, false)
}

// Layout defines the logical screen size. Scaling is handled by Ebitengine.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// --- Main Function ---

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Harvest Demo With Ebitengine in Golang")
	ebiten.SetTPS(60) // Target 60 ticks per second for updates

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
