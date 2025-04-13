package game

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

var (
	assetStore   *assets.Store
	audioContext *audio.Context
	musicPlayer  *audio.Player
)

// --- Types ---

type Game struct {
	Player *player.Player
	World  *world.World
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

	/// --- SFX TEST PLS REMOVE LATER IN THE GAME ---
	// For Testing pressing the space button will play a lazer sound
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		laserSfx, ok := assetStore.GetSFXData("laser")
		if ok {
			sfxPlayer := audioContext.NewPlayerFromBytes(laserSfx)
			sfxPlayer.Play()
		}
	}

	// --- World ---
	g.World.Update(g.Player.Pos, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, dt)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// --- Drawing the Map ---
	g.World.Draw(screen)

	// --- Drawing the Player ---
	mapOffsetX, mapOffsetY := g.World.GetCameraPosition()
	g.Player.Draw(screen, assetStore, mapOffsetX, mapOffsetY)
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
	g := &Game{
		Player: p,
		World:  w,
	}
	return g
}
