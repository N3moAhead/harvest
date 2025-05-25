package hud

import (
	"sort"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/pkg/ui"
)

type VegtableInventoryDisplay struct {
	ui.Container
	inv              *inventory.Inventory
	vegtableDisplays []VegtableDisplay
}

func NewVegtableInventoryDisplay(x, y float64, invRef *inventory.Inventory) *VegtableInventoryDisplay {
	containerOptions := &ui.ContainerOptions{
		Direction: ui.Col,
		Gap:       0,
	}
	vegtableDisplay := &VegtableInventoryDisplay{
		Container: *ui.NewContainer(x, y, containerOptions),
		inv:       invRef,
	}

	for range config.VEGTABLE_TYPE_AMOUNT {
		// The position does not matter because the container will set it for us
		// after adding the vegtable dispaly to it
		newVegtableDisplay := NewVegtableDisplay(10, 10, invRef)
		vegtableDisplay.AddChild(newVegtableDisplay)
	}

	return vegtableDisplay
}

func (v *VegtableInventoryDisplay) Update(input *ui.InputState) {
	currentChild := 0
	vegtableTypes := make([]itemtype.ItemType, 0)
	for k, _ := range v.inv.Vegetables {
		vegtableTypes = append(vegtableTypes, k)
	}
	sort.Sort(itemtype.ByItemType(vegtableTypes))
	for _, vegtableKey := range vegtableTypes {
		child := v.Children[currentChild]
		if vegtableDisplay, ok := child.(VegtableDisplayInterface); ok {
			vegtableDisplay.UpdateVegtableDisplay(vegtableKey)
		}
		currentChild++
	}
	for currentChild < len(v.Children) {
		child := v.Children[currentChild]
		if vegtableDisplay, ok := child.(VegtableDisplayInterface); ok {
			vegtableDisplay.UpdateVegtableDisplay(itemtype.Undefined)
		}
		currentChild++
	}
	v.Container.Update(input)
}

var _ ui.UIElement = (*VegtableInventoryDisplay)(nil)
