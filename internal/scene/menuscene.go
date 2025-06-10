package scene

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/hud"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type MenuScene struct {
	BaseScene
	uiManager    *ui.UIManager
	icon         *ebiten.Image
	world        *world.World
	targetPos    component.Vector2D
	currentAngle float64
	isRunning    bool
	angularSpeed float64
}

func NewMenuScene(setExitGame func(), stats PlayerStats) *MenuScene {
	icon, ok := assets.AssetStore.GetImage("menu-icon")
	if !ok {
		panic("menu-icon nicht im AssetStore gefunden")
	}

	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Unable to load font in new base scene")
	}
	microFont, ok := assets.AssetStore.GetFont("micro")
	if !ok {
		panic("Unable to load font in new base scene")
	}
	newUiManager := ui.NewUIManager()
	newMenuScene := &MenuScene{
		BaseScene:    *NewBaseScene(),
		uiManager:    newUiManager,
		icon:         icon,
		isRunning:    true,
		world:        world.NewWorld(400, 400),
		angularSpeed: 0.1,
		targetPos: component.NewVector2D(float64((config.WIDTH_IN_TILES*config.TILE_SIZE)/2),
			float64((config.HEIGHT_IN_TILES*config.TILE_SIZE)/2.0)),
		currentAngle: 0.0,
	}

	startBtn := ui.NewButton(0, 0, 150, 40, "Start", fontFace, func() { newMenuScene.SetIsRunning(false) })
	endGameBtn := ui.NewButton(0, 0, 150, 40, "Exit", fontFace, setExitGame)
	container := ui.NewContainer((config.SCREEN_WIDTH-150)/2, 350, &ui.ContainerOptions{
		Direction: ui.Col,
		Gap:       10,
	})
	container.AddChild(startBtn)
	container.AddChild(endGameBtn)
	newUiManager.AddElement(container)
	highScoreDisplay := hud.NewScoreDisplay(&stats.highScore, "Highscore")
	newUiManager.AddElement(highScoreDisplay)

	statsContainer := ui.NewContainer(10, 10, &ui.ContainerOptions{
		Direction: ui.Col,
		Gap:       10,
	})
	levelDisplay := ui.NewLabel(0, 0, fmt.Sprintf("Player Level: %d", stats.playerLevel), microFont, color.White)
	statsContainer.AddChild(levelDisplay)
	newUiManager.AddElement(statsContainer)

	music, ok := assets.AssetStore.GetMusicData("menu")
	if ok {
		musicBytesReader := bytes.NewReader(music)
		loop := audio.NewInfiniteLoop(musicBytesReader, int64(len(music)))
		if assets.MusicPlayer != nil {
			assets.MusicPlayer.Close()
			assets.MusicPlayer, _ = assets.AudioContext.NewPlayer(loop)
		}
		assets.MusicPlayer, _ = assets.AudioContext.NewPlayer(loop)
		assets.MusicPlayer.Play()
	} else {
		fmt.Println("Warning: could not load menu music")
	}

	return newMenuScene
}

func (l *MenuScene) Draw(screen *ebiten.Image) {
	screenWidth := screen.Bounds().Dx()
	iconWidth := l.icon.Bounds().Dx()
	scale := 0.14 // Scale the icon
	scaledW := float64(iconWidth) * scale

	l.world.Draw(screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		(float64(screenWidth)-scaledW)/2, // horziontal center
		30,                               // 30px from top
	)
	screen.DrawImage(l.icon, op)

	l.uiManager.Draw(screen)
}

func (l *MenuScene) Update() error {
	dt := 1.0 / float64(ebiten.TPS())

	circleCenterX := float64((config.WIDTH_IN_TILES * config.TILE_SIZE) / 2)
	circleCenterY := float64((config.HEIGHT_IN_TILES * config.TILE_SIZE) / 2.0)

	radius := 650.0

	l.currentAngle += l.angularSpeed * dt

	if l.currentAngle >= 2*math.Pi {
		l.currentAngle -= 2 * math.Pi
	}
	if l.currentAngle < 0 {
		l.currentAngle += 2 * math.Pi
	}

	l.targetPos.X = circleCenterX + radius*math.Cos(l.currentAngle)
	l.targetPos.Y = circleCenterY + radius*math.Sin(l.currentAngle)

	l.world.Update(l.targetPos, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, dt)
	l.uiManager.Update()
	return nil
}

func (m *MenuScene) IsRunning() bool {
	return m.isRunning
}

func (m *MenuScene) SetIsRunning(r bool) {
	m.isRunning = r
}

var _ Scene = (*MenuScene)(nil)
