package scene

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/cooking"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/internal/input"
	"github.com/N3moAhead/harvest/internal/weapon"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// --- Types ---

type GameScene struct {
	Player               *player.Player
	World                *world.World
	Enemies              []enemy.EnemyInterface
	Spawner              *world.EnemySpawner
	previousSpacePressed bool // TODO remove this later, just for testing
	items                []*item.Item
	inventory            *inventory.Inventory
	ui                   *ui.UIManager
	isRunning            bool
	cookStations         []*cooking.CookStation
}

func (g *GameScene) Update() error {
	// --- Delta Time Update ---
	dt := 1.0 / float64(ebiten.TPS())
	dtDuration := time.Second / time.Duration(ebiten.TPS())

	/// --- Get User Input ---
	inputState := input.GetInputState()

	// --- Player Update ---
	g.Player.Update(inputState, dt, g.inventory)

	/// --- UI Update ---
	g.ui.Update()

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
		gItem := g.items[i]
		// The update function also puts collected items into the inventory
		itemPickedUp := gItem.Update(g.Player)
		if itemPickedUp {
			// Add picked up items into the inventory
			switch gItem.CategoryOf() {
			case itemtype.CategoryVegetable:
				g.inventory.AddVegtable(gItem.Type)
			case itemtype.CategorySoup:
				g.inventory.AddSoup(gItem.Type)
				soup := gItem.RetrieveItemInfo().Soup
				g.Player.ExtendOrAddSoup(soup)
			case itemtype.CategoryWeapon:
				switch gItem.Type {
				case itemtype.Spoon:
					newWeapon := weapon.NewSpoon()
					added := g.inventory.AddWeapon(newWeapon)
					if !added {
						fmt.Printf("Inventory is full or weapon '%s' already exists\n", newWeapon.Name())
					} else {
						fmt.Printf("Weapon '%s' added to Inventory\n", newWeapon.Name())
					}
				default:
					fmt.Printf("Warning: Unknown weapon type: %s", gItem.DisplayName())
				}
			default:
				panic(fmt.Errorf("unhandeld item category: %s in items update", gItem.CategoryOf().String()))
			}
		} else {
			// Remove items after the player picked them up
			if n != i {
				g.items[n] = gItem
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
		for range 3 {
			posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
			posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
			g.items = append(g.items, item.NewSoup(posX, posY, itemtype.SpeedSoup))
		}
		// to test magnet
		// for range 3 {
		// 	posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
		// 	posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
		// 	g.items = append(g.items, item.NewSoup(posX, posY, itemtype.MagnetRadiusSoup))
		// }
	}

	if ebiten.IsKeyPressed(ebiten.KeyC) {
		for range 3 {
			posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
			posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
			g.cookStations = append(g.cookStations, cooking.NewCookStation(
				posX,
				posY,
				// cooking.RecipeDefinitions[itemtype.MagnetRadiusSoup],
				cooking.RecipeDefinitions[itemtype.SpeedSoup],
				1.0, // cost factor?
			))
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

	for _, cookStation := range g.cookStations {
		cookStation.Update(g.Player, g.inventory)
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
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

	for _, cookStation := range g.cookStations {
		cookStation.Draw(screen, mapOffsetX, mapOffsetY)
	}

	// --- Drawing the HUD ---
	fpsText := fmt.Sprintf("FPS: %.1f ", ebiten.ActualFPS())
	ebitenutil.DebugPrintAt(screen, fpsText+fmt.Sprintf("HP: %d / %d\n", int(g.Player.Health.HP), int(g.Player.Health.MaxHP)), 10, 5)
	g.inventory.Draw(screen)
	g.ui.Draw(screen)
}

func (g *GameScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.SCREEN_WIDTH, config.SCREEN_HEIGHT
}

func (g *GameScene) IsRunning() bool {
	return g.isRunning
}

func (g *GameScene) SetIsRunning(running bool) {
	g.isRunning = running
}

// --- Public ---

func NewGameScene() *GameScene {
	p := player.NewPlayer()
	w := world.NewWorld(config.WIDTH_IN_TILES, config.HEIGHT_IN_TILES)
	s := world.NewEnemySpawner()
	i := inventory.NewInventory()
	items := []*item.Item{
		item.NewSpoon(50, 50),
	}
	uiManager := ui.NewUIManager()
	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Font Face could not be loaded")
	}

	label1 := ui.NewLabel(50, 50, "Lachen", fontFace, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	label2 := ui.NewLabel(50, 80, "Weinen", fontFace, color.RGBA{R: 255, G: 255, B: 0, A: 255})
	label3 := ui.NewLabel(200, 400, "Tanzen", fontFace, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	label4 := ui.NewLabel(200, 200, "Welt", fontFace, color.RGBA{R: 0, G: 255, B: 255, A: 255})

	container2 := ui.NewContainer(300, 5, &ui.ContainerOptions{
		Direction: ui.Row,
		Gap:       10,
	})
	container2.AddChild(label1)
	container2.AddChild(label2)
	container2.AddChild(label3)
	container2.AddChild(label4)

	uiManager.AddElement(container2)

	// register enemy factories
	s.RegisterFactory(enemy.TypeCarrot.String(), func(pos component.Vector2D) enemy.EnemyInterface {
		return enemy.NewCarrotEnemy(pos)
	})

	newGameScene := &GameScene{
		Player:       p,
		World:        w,
		Enemies:      []enemy.EnemyInterface{},
		Spawner:      s,
		inventory:    i,
		items:        items,
		ui:           uiManager,
		isRunning:    true,
		cookStations: []*cooking.CookStation{},
	}

	nextSceneButton := ui.NewButton(300, 300, 250, 50, "Next", fontFace, func() { newGameScene.SetIsRunning(false) })

	container1 := ui.NewContainer(5, 150, &ui.ContainerOptions{
		Direction: ui.Col,
		Gap:       10,
	})
	container1.AddChild(nextSceneButton)
	uiManager.AddElement(container1)

	return newGameScene
}
