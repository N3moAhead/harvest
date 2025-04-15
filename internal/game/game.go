package game

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/inventory"
	"github.com/N3moAhead/harvest/internal/item"
	"github.com/N3moAhead/harvest/internal/itemtype"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	assetStore   *assets.Store
	audioContext *audio.Context
	musicPlayer  *audio.Player
)

// --- Types ---

type Game struct {
	Player    *player.Player
	World     *world.World
	items     []*item.Item
	inventory *inventory.Inventory
}

func (g *Game) Update() error {
	// --- Delta Time Update ---
	dt := 1.0 / float64(ebiten.TPS())

	// --- Check for Exit ---
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Game Quit!")
	}

	// --- Player Input & Movement ---
	moveDir := component.Vector2D{X: 0, Y: 0}
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
	if moveDir.Y != 0 && moveDir.X != 0 {
		moveDir = moveDir.Normalize()
	}
	g.Player.Pos = g.Player.Pos.Add(moveDir.Mul(g.Player.Speed))

	for i := range len(g.items) {
		item := g.items[i]
		removeItem := item.Update(g.Player, g.inventory)
		if removeItem {

		}
	}

	/// --- Update Items on the Ground ---
	// TODO maybe we should move this code out of the game.go file to keep it clean
	n := 0
	for i := range g.items {
		item := g.items[i]
		// The update function also puts collected items into the inventory
		removeItem := item.Update(g.Player, g.inventory)
		// Remove items after
		// the player picked them up
		if !removeItem {
			if n != i {
				g.items[n] = item
			}
			n++
		}
	}
	g.items = g.items[:n]

	// TODO Testing spawning items pls remove for production!
	// Pressing K will spawn items
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		for range 10 {
			posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
			posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
			g.items = append(g.items, item.NewCarrot(posX, posY))
		}
		for range 10 {
			posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
			posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
			g.items = append(g.items, item.NewPotato(posX, posY))
		}
	}

	// Testing sfx Remove for production
	// Pressing L will play a lazer sound
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		laserSfx, ok := assetStore.GetSFXData("laser")
		if ok {
			sfxPlayer := audioContext.NewPlayerFromBytes(laserSfx)
			sfxPlayer.Play()
		}
	}

	// --- World ---
	// TODO remove this code just for testing you can display fancy camera movement
	// Pressing space will move the camera to the top left
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.World.Update(component.Vector2D{X: 0.0, Y: 0.0}, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, dt)
	} else {
		g.World.Update(g.Player.Pos, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, dt)
	}

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
	g.Player.Draw(screen, assetStore, mapOffsetX, mapOffsetY)

	// --- Drawing the HUD ---
	if amount, ok := g.inventory.Vegtables[itemtype.Potato]; ok {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Item: %s, Amount: %d\n\n", itemtype.Potato.String(), amount), 10, 10)
	}
	if amount, ok := g.inventory.Vegtables[itemtype.Carrot]; ok {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Item: %s, Amount: %d\n\n", itemtype.Carrot.String(), amount), 10, 30)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.SCREEN_WIDTH, config.SCREEN_HEIGHT
}

// --- Internal ---

func init() {
	ebiten.SetWindowSize(config.SCREEN_WIDTH, config.SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Harvest by Wurzelwerk")
	ebiten.SetTPS(60)

	// TODO Move all this assetStore init stuff into
	// a seperate file to keep the game.go file clean

	// A new Audio Context
	audioContext = audio.NewContext(config.AUDIO_SAMPLE_RATE)
	// Initing the asset store
	assetStore = assets.NewStore()

	// Always image name to path
	imagesToLoad := map[string]string{
		"player": "assets/images/CookTestImage.png",
	}
	sfxToLoad := map[string]string{
		"laser": "assets/audio/sfx/laserTest.wav",
	}
	musicToLoad := map[string]string{
		"menu": "assets/audio/music/8bitMenuMusic.mp3",
	}

	err := assetStore.Load(imagesToLoad, sfxToLoad, musicToLoad, config.AUDIO_SAMPLE_RATE)
	if err != nil {
		panic(err)
	}

	// TODO REMOVE or change this section
	// This here should just be a test to test running music :)
	music, ok := assetStore.GetMusicData("menu")
	if ok {
		musicBytesReader := bytes.NewReader(music)
		loop := audio.NewInfiniteLoop(musicBytesReader, int64(len(music)))

		musicPlayer, err = audioContext.NewPlayer(loop)
		if err == nil {
			musicPlayer.Play()
		} else {
			err = fmt.Errorf("Musikplayer konnte nicht erstellt werden: %v\n", err)
			panic(err)
		}
	}
}

// --- Public ---

func NewGame() *Game {
	p := player.NewPlayer()
	w := world.NewWorld(config.WIDTH_IN_TILES, config.HEIGHT_IN_TILES)
	i := inventory.NewInventory()
	g := &Game{
		Player:    p,
		World:     w,
		inventory: i,
	}
	return g
}
