package item

import (
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
)

/// --- Vegtables ---

func NewCarrot(posX float64, posY float64) *Item {
	carrotIcon, _ := assets.AssetStore.GetImage("carrot_icon")
	newItem := newItemBase(posX, posY, carrotIcon)
	newItem.Type = itemtype.Carrot
	return newItem
}

func NewPotato(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY, nil)
	newItem.Type = itemtype.Potato
	return newItem
}

/// --- Weapons ---

func NewSpoon(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY, nil)
	newItem.Type = itemtype.Spoon
	return newItem
}

/// --- Soups ---

func NewSoup(x, y float64, typeBuff itemtype.ItemType) *Item {
	newItem := newItemBase(x, y, nil)
	newItem.Type = typeBuff
	return newItem
}
