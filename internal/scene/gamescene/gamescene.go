package gamescene

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/cooking"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/internal/input"
	"github.com/N3moAhead/harvest/internal/toast"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Score interface {
	GetScore() int
}

type GameScene struct {
	lastSpawnTime            time.Time // last spawn batches
	waveDefinitions          []WaveDefinition
	currentWaveIndex         int
	lastWaveStartTime        time.Time
	gameStartTime            time.Time
	Player                   *player.Player
	World                    *world.World
	Enemies                  []enemy.EnemyInterface
	Spawner                  *world.EnemySpawner
	items                    []*item.Item
	inventory                *inventory.Inventory
	hud                      *ui.UIManager
	gameOverlay              *ui.UIManager
	isRunning                bool
	isPaused                 bool
	cookStations             []*cooking.CookStation
	startTime                time.Time // game start
	lastEnemySpawnTime       time.Time // last spawn batches
	lastCookStationSpawnTime time.Time
	Score                    int
}

func NewGameScene(backToMenu func(), playerLevel uint) *GameScene {
	newGameScene := &GameScene{
		Player:             player.NewPlayer(playerLevel),
		World:              world.NewWorld(config.WIDTH_IN_TILES, config.HEIGHT_IN_TILES),
		Enemies:            []enemy.EnemyInterface{},
		Spawner:            initEnemySpawner(),
		inventory:          inventory.NewInventory(),
		items:              initItems(),
		hud:                nil,
		isRunning:          true,
		cookStations:       []*cooking.CookStation{},
		startTime:          time.Now(),
		lastEnemySpawnTime: time.Now(),
		Score:              0,
	}
	newGameScene.initializeWaves()
	newGameScene.hud = initHUD(newGameScene)
	newGameScene.gameOverlay = initGameOverlay(newGameScene, backToMenu)

	/// Init Game Music
	game_music, ok := assets.AssetStore.GetMusicData("game")
	if ok {
		musicBytesReader := bytes.NewReader(game_music)
		loop := audio.NewInfiniteLoop(musicBytesReader, int64(len(game_music)))
		if assets.MusicPlayer != nil {
			assets.MusicPlayer.Close()
			assets.MusicPlayer, _ = assets.AudioContext.NewPlayer(loop)
		}
		assets.MusicPlayer, _ = assets.AudioContext.NewPlayer(loop)
		assets.MusicPlayer.Play()
	} else {
		fmt.Println("Warning: could not load game music")
	}

	return newGameScene
}

func (g *GameScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.SCREEN_WIDTH, config.SCREEN_HEIGHT
}

func (g *GameScene) IsRunning() bool {
	return g.isRunning
}

func (g *GameScene) GetScore() int {
	return g.Score
}

func (g *GameScene) SetIsRunning(running bool) {
	g.isRunning = running
}

func (g *GameScene) Update() error {
	/// --- Get User Input ---
	inputState := input.GetInputState()

	/// --- UI Update ---
	// The ui is always getting updated everything else can be paused.
	updateUI(g)

	// Pause on Escape
	if inputState.Esc {
		g.isPaused = true
	}

	// If the game is paused stop the update right here
	if g.isPaused {
		return nil
	}

	// --- Time Update ---
	dt := 1.0 / float64(ebiten.TPS())
	dtDuration := time.Second / time.Duration(ebiten.TPS())
	elapsed := float32(time.Since(g.startTime).Milliseconds())

	/// --- Update Player ---
	g.Player.Update(inputState, dt, g.inventory)

	/// --- Update Items on the Ground ---
	updateItems(g)

	/// --- Update the Weapons ---
	for _, weapon := range g.inventory.Weapons {
		if weapon != nil {
			weapon.Update(g.Player, g.Enemies, dtDuration)
		}
	}

	// TODO remove this code for production!
	// But keep it until cooking stations can spawn autonomisly
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

	/// --- Update World ---
	g.World.Update(
		g.Player.Pos,
		config.SCREEN_WIDTH,
		config.SCREEN_HEIGHT,
		dt,
	)

	/// --- Toast ---
	toast.UpdateToasts()

	/// --- Update Enemies ---
	updateEnemies(g, dt, elapsed)

	/// --- Update Cooking Stations ---
	updateCookStations(g, dt, g.inventory, elapsed)

	/// --- Check if player died ---
	if !g.Player.Alive() {
		deathSound, ok := assets.AssetStore.GetSFXData("player_death_sound")
		if ok {
			sfxPlayer := assets.AudioContext.NewPlayerFromBytes(deathSound)
			sfxPlayer.Play()
		}
		if assets.MusicPlayer.IsPlaying() {
			assets.MusicPlayer.Pause()
		}
		g.SetIsRunning(false)
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	/// --- Drawing the Map ---
	g.World.Draw(screen)

	mapOffsetX, mapOffsetY := g.World.GetCameraPosition()

	/// --- Drawing all Items ---
	drawItems(g, screen, mapOffsetX, mapOffsetY)

	/// --- Drawing the Player ---
	g.Player.Draw(screen, mapOffsetX, mapOffsetY)

	/// --- Drawing the Enemies ---
	drawEnemies(g, screen, mapOffsetX, mapOffsetY)

	/// --- Drawing Weapon Effects ---
	for _, w := range g.inventory.Weapons {
		if w != nil {
			w.Draw(screen, g.Player, mapOffsetX, mapOffsetY)
		}
	}

	/// --- Drawing Cooking Stations ---
	for _, cookStation := range g.cookStations {
		cookStation.Draw(screen, mapOffsetX, mapOffsetY)
	}

	/// --- Toasts ---
	toast.DrawToasts(screen)

	/// --- Drawing HUD ---
	drawUI(g, screen)
}
