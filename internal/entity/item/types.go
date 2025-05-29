package item

import (
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
)

/// --- Vegtables ---

func NewCarrot(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY, itemtype.Carrot)
	return newItem
}

func NewPotato(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY, itemtype.Potato)
	return newItem
}

/// --- Weapons ---

func NewSpoon(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY, itemtype.Spoon)
	return newItem
}

func NewRollingPin(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY, itemtype.RollingPin)
	return newItem
}

func NewThermalmixer(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY, itemtype.Thermalmixer)
	return newItem
}

/// --- Soups ---

func NewSoup(x, y float64, typeBuff itemtype.ItemType) *Item {
	newItem := newItemBase(x, y, typeBuff)
	return newItem
}
