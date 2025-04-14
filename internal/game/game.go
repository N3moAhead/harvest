package game

import (
	"errors"
	"math/rand/v2"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/enemy"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/internal/world"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
)

// --- Types ---

type Game struct {
	Player *player.Player
	World  *world.World
	Enemies []enemy.EnemyInterface
	previousSpacePressed bool // TODO remove this later, just for testing
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

	/// --- ENEMY TEST PLS REMOVE LATER IN THE GAME ---
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	if spacePressed  && !g.previousSpacePressed {
		pos := component.NewVector2D(rand.Float64()*500, rand.Float64()*500)
		// pos := component.NewVector2D(100, 100)
		e := enemy.NewCarrotEnemy(pos)
		g.Enemies = append(g.Enemies, e)
	}
	g.previousSpacePressed = spacePressed
	// --- World ---
	g.World.Update(g.Player.Pos, config.SCREEN_WIDTH, config.SCREEN_HEIGHT, dt)

	// --- Enemies ---
	for _, e := range g.Enemies {
		e.Update(g.Player, dt)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// --- Drawing the Map ---
	g.World.Draw(screen)

	// --- Drawing the Player ---
	cameraPosX, cameraPosY := g.World.GetCameraPosition()
	g.Player.Draw(screen, cameraPosX, cameraPosY)

	// --- Drawing the Enemies ---
	for _, e := range g.Enemies {
		if e.IsAlive() {
			e.Draw(screen, cameraPosX, cameraPosY)
		}
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
}

// --- Public ---

func NewGame() *Game {
	p := player.NewPlayer()
	w := world.NewWorld(config.WIDTH_IN_TILES, config.HEIGHT_IN_TILES)
	pos := component.NewVector2D(100, 100)// TODO remove this later, just for testing
	e := enemy.NewCarrotEnemy(pos)
	g := &Game{
		Player: p,
		World:  w,
		Enemies: []enemy.EnemyInterface{e},// TODO remove e later, just for testing
	}
	return g
}
