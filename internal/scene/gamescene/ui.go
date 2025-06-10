package gamescene

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/hud"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func updateUI(g *GameScene) {
	if g.isPaused {
		g.gameOverlay.Update()
	} else {
		g.hud.Update()
	}
}

func drawUI(g *GameScene, screen *ebiten.Image) {
	if g.isPaused {
		g.gameOverlay.Draw(screen)
	} else {
		// --- Drawing the HUD ---
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d / %d\n", int(g.Player.Health.HP), int(g.Player.Health.MaxHP)), 10, config.SCREEN_HEIGHT-20)
		g.inventory.Draw(screen)
		g.hud.Draw(screen)
	}
}

func initHUD(g *GameScene) *ui.UIManager {
	newHUD := ui.NewUIManager()

	inventoryDisplay := hud.NewInventoryDisplay(10, 10, g.inventory)
	weaponDisplay := hud.NewWeaponDisplay(40, 10, g.inventory)
	frameContainer := ui.NewContainer(5, 5, &ui.ContainerOptions{
		Direction: ui.Row,
		Gap:       10,
	})
	frameContainer.AddChild(inventoryDisplay)
	frameContainer.AddChild(weaponDisplay)
	newHUD.AddElement(frameContainer)

	scoreDisplay := hud.NewScoreDisplay(&g.Score, "Score")

	newHUD.AddElement(scoreDisplay)

	return newHUD
}

func initGameOverlay(g *GameScene, backToMenu func()) *ui.UIManager {
	newGameOverlay := ui.NewUIManager()

	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Font Face could not be loaded")
	}

	// Container
	elementWidth := 300.0
	containerDrawX := (float64(config.SCREEN_WIDTH) - elementWidth) / 2
	container := ui.NewContainer(containerDrawX, 200, &ui.ContainerOptions{
		Direction: ui.Col,
		Gap:       10,
	})
	// Buttons
	resumeBtn := ui.NewButton(0, 0, elementWidth, 50, "Resume", fontFace, func() { g.isPaused = false })
	exitBtn := ui.NewButton(0, 0, elementWidth, 50, "Exit Game", fontFace, func() { backToMenu() })

	container.AddChild(resumeBtn)
	container.AddChild(exitBtn)
	newGameOverlay.AddElement(container)

	return newGameOverlay
}
