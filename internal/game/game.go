package game

import (
	"errors"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
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

	// --- World ---
	g.World.Update(g.Player.Pos, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, dt)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// --- Drawing the Map ---
	g.World.Draw(screen)

	// --- Drawing the Player ---
	cameraPosX, cameraPosY := g.World.GetCameraPosition()
	g.Player.Draw(screen, cameraPosX, cameraPosY)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.SCREEN_WIDTH, config.SCREEN_HEIGHT
}

// --- Internal ---

func init() {
	ebiten.SetWindowSize(config.SCREEN_WIDTH, config.SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Harvest by Wurzelwerk")
	ebiten.SetTPS(60)
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
