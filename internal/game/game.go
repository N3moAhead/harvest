package game

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/enemy"
	"github.com/N3moAhead/harvest/internal/inventory"
	"github.com/N3moAhead/harvest/internal/item"
	"github.com/N3moAhead/harvest/internal/itemtype"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/internal/weapon"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// --- Types ---

type Game struct {
	Player               *player.Player
	World                *world.World
	Enemies              []enemy.EnemyInterface
	Spawner              *world.EnemySpawner
	previousSpacePressed bool // TODO remove this later, just for testing
	items                []*item.Item
	inventory            *inventory.Inventory
}

func (g *Game) Update() error {
	// --- Delta Time Update ---
	dt := 1.0 / float64(ebiten.TPS())
	dtDuration := time.Second / time.Duration(ebiten.TPS())

	// --- Check for Exit ---
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Game Quit!")
	}

	// --- Player Input & Movement ---
	moveDir := component.Vector2D{X: 0, Y: 0}
	moved := false
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		moveDir.Y -= 1
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		moveDir.Y += 1
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveDir.X -= 1
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveDir.X += 1
		moved = true
	}
	if moveDir.Y != 0 && moveDir.X != 0 {
		moveDir = moveDir.Normalize()
		moved = true
	}

	// If the player moved update the facingDirection and the player position
	if moved {
		normalizedMoveDirection := moveDir.Normalize()
		g.Player.Pos = g.Player.Pos.Add(moveDir.Mul(g.Player.Speed))
		g.Player.FacingDirection = normalizedMoveDirection
	}

	/// --- SFX TEST && ENEMY TEST PLS REMOVE LATER IN THE GAME ---
	// For Testing pressing the space button will play a lazer sound
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	if spacePressed && !g.previousSpacePressed {
		// Circle Pattern
		newEnemies := g.Spawner.SpawnCircle("carrot", g.Player, 150, 8)
		fmt.Println("New Enemies Spawned:", newEnemies)

		g.Enemies = append(g.Enemies, newEnemies...)

		// ZigZag Pattern
		// padding := component.Vector2D{X: 50, Y: 0}
		// newEnemies = g.Spawner.SpawnZigZag("carrot", g.Player.Pos.Add(padding), 6, 40, 30)
		// g.Enemies = append(g.Enemies, newEnemies...)

		// Line Pattern
		// newEnemies = g.Spawner.SpawnLine("carrot", g.Player.Pos.Add(padding), 5, 30, 0)

		// Random Pattern
		// newEnemies = g.Spawner.SpawnMoreRandom(10, "carrot")

		// Appends new enemies to the world
		// g.Enemies = append(g.Enemies, newEnemies...)

		// for _, enemy := range newEnemies {
		// 	g.Enemies = append(g.Enemies, enemy)
		// }
	}
	g.previousSpacePressed = spacePressed

	/// --- Update Items on the Ground ---
	// TODO maybe we should move this code out of the game.go file to keep it clean
	n := 0
	for i := range g.items {
		item := g.items[i]
		// The update function also puts collected items into the inventory
		itemPickedUp := item.Update(g.Player)
		if itemPickedUp {
			// Add picked up items into the inventory
			switch item.Type.Category() {
			case itemtype.CategoryVegetable:
				g.inventory.AddVegtable(item.Type)
				break
			case itemtype.CategoryWeapon:
				switch item.Type {
				case itemtype.Spoon:
					newWeapon := weapon.NewSpoon()
					added := g.inventory.AddWeapon(newWeapon)
					if !added {
						fmt.Printf("Inventory is full or weapon '%s' already exists\n", newWeapon.Name())
					} else {
						fmt.Printf("Weapon '%s' added to Inventory\n", newWeapon.Name())
					}
					break
				default:
					fmt.Printf("Warning: Unknown weapon type: %s", item.Type.String())
					break
				}
				break
			default:
				panic(fmt.Errorf("Unhandeld item category: %s in items update", item.Type.Category().String()))
			}
		} else {
			// Remove items after the player picked them up
			if n != i {
				g.items[n] = item
			}
			n++
		}
	}
	g.items = g.items[:n]

	/// --- Update the Weopons ---
	// TODO check for copy usage
	for _, weapon := range g.inventory.Weapons {
		if weapon != nil {
			weapon.Update(g.Player, g.Enemies, dtDuration)
		}
	}

	// TODO Testing spawning items pls remove for production!
	// Pressing K will spawn items
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		for range 50 {
			posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
			posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
			g.items = append(g.items, item.NewCarrot(posX, posY))
		}
		for range 50 {
			posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
			posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
			g.items = append(g.items, item.NewPotato(posX, posY))
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyB) {
		// to test speed
		// for range 3 {
		// 	posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
		// 	posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
		// 	g.items = append(g.items, item.NewSoup(posX, posY, itemtype.SpeedSoup))
		// }
		// to test magnet
		for range 3 {
			posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
			posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
			g.items = append(g.items, item.NewSoup(posX, posY, itemtype.MagnetRadiusSoup))
		}
	}

	// Testing sfx Remove for production
	// Pressing L will play a lazer sound
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		laserSfx, ok := assets.AssetStore.GetSFXData("laser")
		if ok {
			sfxPlayer := assets.AudioContext.NewPlayerFromBytes(laserSfx)
			sfxPlayer.Play()
		}
	}

	// --- World ---
	// TODO remove this code just for testing you can display fancy camera movement
	// Pressing J will move the camera to the top left
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		g.World.Update(component.Vector2D{X: 0.0, Y: 0.0}, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, dt)
	} else {
		g.World.Update(g.Player.Pos, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, dt)
	}

	// --- Enemies ---
	for _, e := range g.Enemies {
		e.Update(g.Player, dt)
	}

	g.Player.Update(dt, g.inventory)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// --- Drawing the Map ---
	g.World.Draw(screen)

	mapOffsetX, mapOffsetY := g.World.GetCameraPosition()

	// --- Drawing all Items ---
	// Warning currently it's not getting checked if an item is on screen or not
	for _, item := range g.items {
		item.Draw(screen, mapOffsetX, mapOffsetY)
	}

	// --- Drawing the Player ---
	g.Player.Draw(screen, mapOffsetX, mapOffsetY)

	// --- Drawing the Enemies ---
	for _, e := range g.Enemies {
		if e.IsAlive() {
			e.Draw(screen, mapOffsetX, mapOffsetY)
		}
	}

	// -- Drawing Weapon Effects ---
	for _, w := range g.inventory.Weapons {
		if w != nil {
			w.Draw(screen, g.Player, mapOffsetX, mapOffsetY)
		}
	}

	// --- Drawing the HUD ---
	fpsText := fmt.Sprintf("FPS: %.1f ", ebiten.ActualFPS())
	ebitenutil.DebugPrintAt(screen, fpsText+fmt.Sprintf("HP: %d / %d\n", int(g.Player.Health.HP), int(g.Player.Health.MaxHP)), 10, 5)
	g.inventory.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.SCREEN_WIDTH, config.SCREEN_HEIGHT
}

// --- Internal ---

func init() {
	ebiten.SetWindowSize(config.SCREEN_WIDTH, config.SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Harvest by Wurzelwerk")
	ebiten.SetTPS(60)
	ebiten.SetFullscreen(true)
}

// --- Public ---

func NewGame() *Game {
	p := player.NewPlayer()
	w := world.NewWorld(config.WIDTH_IN_TILES, config.HEIGHT_IN_TILES)
	s := world.NewEnemySpawner()
	i := inventory.NewInventory()
	items := []*item.Item{
		item.NewSpoon(50, 50),
	}

	// register enemy factories
	s.RegisterFactory(enemy.TypeCarrot.String(), func(pos component.Vector2D) enemy.EnemyInterface {
		return enemy.NewCarrotEnemy(pos)
	})
	g := &Game{
		Player:    p,
		World:     w,
		Enemies:   []enemy.EnemyInterface{},
		Spawner:   s,
		inventory: i,
		items:     items,
	}

	return g
}
