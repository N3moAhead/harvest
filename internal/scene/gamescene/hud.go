package gamescene

import (
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/hud"
	"github.com/N3moAhead/harvest/pkg/ui"
)

func initHUD(g *GameScene) *ui.UIManager {
	uiManager := ui.NewUIManager()

	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Font Face could not be loaded")
	}

	inventoryDisplay := hud.NewInventoryDisplay(10, 10, g.inventory)
	weaponDisplay := hud.NewWeaponDisplay(40, 10, g.inventory)
	frameContainer := ui.NewContainer(5, 5, &ui.ContainerOptions{
		Direction: ui.Row,
		Gap:       10,
	})
	frameContainer.AddChild(inventoryDisplay)
	frameContainer.AddChild(weaponDisplay)
	uiManager.AddElement(frameContainer)

	// TODO build a real in game overlay
	pauseButton := ui.NewButton(200, 200, 200, 50, "Toggle Pause", fontFace, func() {
		g.isPaused = !g.isPaused
	})
	uiManager.AddElement(pauseButton)

	return uiManager
}
