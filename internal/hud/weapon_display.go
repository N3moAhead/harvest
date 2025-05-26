package hud

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/pkg/ui"
)

type WeaponDisplay struct {
	ui.Container
	inv *inventory.Inventory
}

func NewWeaponDisplay(x, y float64, invRef *inventory.Inventory) *WeaponDisplay {
	containerOptions := &ui.ContainerOptions{
		Direction: ui.Row,
		Gap:       0,
	}
	newDisplay := &WeaponDisplay{
		Container: *ui.NewContainer(x, y, containerOptions),
		inv:       invRef,
	}

	newDisplay.Width = config.MAX_WEAPONS * config.ITEM_FRAME_SIZE
	newDisplay.Height = config.ITEM_FRAME_SIZE

	for range config.MAX_WEAPONS {
		// The position for the frame does not matter because the
		// container will set it for us
		newWeaponFrame := NewWeaponFrame(10, 10, invRef)
		newDisplay.AddChild(newWeaponFrame)
	}

	return newDisplay
}

func (w *WeaponDisplay) Update(input *ui.InputState) {
	currentWeapon := 0
	for _, child := range w.Children {
		if weaponFrame, ok := child.(WeaponFrameInterface); ok {
			weapon := w.inv.Weapons[currentWeapon]
			if weapon != nil {
				fmt.Println(weapon.GetType())
				weaponFrame.UpdateWeaponFrameValues(weapon.GetType(), weapon.Description())
			} else {
				weaponFrame.UpdateWeaponFrameValues(itemtype.Undefined, "")
			}
			currentWeapon++
		}
	}
	w.Container.Update(input)
}

var _ ui.UIElement = (*WeaponDisplay)(nil)
