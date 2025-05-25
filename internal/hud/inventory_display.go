package hud

import (
	"fmt"
	"sort"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/pkg/ui"
)

type InventoryDisplay struct {
	ui.Container
	inv *inventory.Inventory
}

func NewInventoryDisplay(x, y float64, invRef *inventory.Inventory) *InventoryDisplay {
	containerOptions := &ui.ContainerOptions{
		Direction: ui.Col,
		Gap:       0,
	}
	inventoryDisplay := &InventoryDisplay{
		Container: *ui.NewContainer(x, y, containerOptions),
		inv:       invRef,
	}

	inventoryDisplay.Width = config.ITEM_FRAME_SIZE
	inventoryDisplay.Height = config.ITEM_FRAME_SIZE * config.VEGTABLE_TYPE_AMOUNT * config.SOUP_TYPE_AMOUNT

	for range config.VEGTABLE_TYPE_AMOUNT {
		// The position does not matter because the container will set it for us
		// after adding the vegtable dispaly to it
		vegtableFrame, ok := assets.AssetStore.GetImage("vegtable_item_frame")
		if !ok {
			fmt.Println("Warning: Could not load vegtable_item_frame in NewInventoryDisplay")
		}
		newVegtableDisplay := NewItemDisplay(10, 10, invRef, vegtableFrame)
		inventoryDisplay.AddChild(newVegtableDisplay)
	}

	for range config.SOUP_TYPE_AMOUNT {
		soupFrame, ok := assets.AssetStore.GetImage("soup_item_frame")
		if !ok {
			fmt.Println("Warning: Could not laod soup_item_frame in NewInventoryDisplay")
		}
		newSoupDisplay := NewItemDisplay(10, 10, invRef, soupFrame)
		inventoryDisplay.AddChild(newSoupDisplay)
	}

	return inventoryDisplay
}

func (v *InventoryDisplay) Update(input *ui.InputState) {
	currentChild := 0
	vegtableTypes := make([]itemtype.ItemType, 0)
	for k, _ := range v.inv.Vegetables {
		vegtableTypes = append(vegtableTypes, k)
	}
	sort.Sort(itemtype.ByItemType(vegtableTypes))
	for _, vegtableKey := range vegtableTypes {
		child := v.Children[currentChild]
		if vegtableDisplay, ok := child.(ItemDisplayInterface); ok {
			vegtableDisplay.UpdateItemDisplay(vegtableKey)
		}
		currentChild++
	}
	for currentChild < config.VEGTABLE_TYPE_AMOUNT {
		child := v.Children[currentChild]
		if vegtableDisplay, ok := child.(ItemDisplayInterface); ok {
			vegtableDisplay.UpdateItemDisplay(itemtype.Undefined)
		}
		currentChild++
	}
	soupTypes := make([]itemtype.ItemType, 0)
	for k, _ := range v.inv.Soups {
		soupTypes = append(soupTypes, k)
	}
	sort.Sort(itemtype.ByItemType(soupTypes))
	for _, soupType := range soupTypes {
		child := v.Children[currentChild]
		if soupDisplay, ok := child.(ItemDisplayInterface); ok {
			soupDisplay.UpdateItemDisplay(soupType)
		}
		currentChild++
	}
	for currentChild < config.VEGTABLE_TYPE_AMOUNT+config.SOUP_TYPE_AMOUNT {
		child := v.Children[currentChild]
		if soupDisplay, ok := child.(ItemDisplayInterface); ok {
			soupDisplay.UpdateItemDisplay(itemtype.Undefined)
		}
		currentChild++
	}
	v.Container.Update(input)
}

var _ ui.UIElement = (*InventoryDisplay)(nil)
