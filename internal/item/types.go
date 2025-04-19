package item

import "github.com/N3moAhead/harvest/internal/itemtype"

/// --- Vegtables ---

func NewCarrot(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY)
	newItem.Type = itemtype.Carrot
	return newItem
}

func NewPotato(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY)
	newItem.Type = itemtype.Potato
	return newItem
}

/// --- Weopons ---

func NewSpoon(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY)
	newItem.Type = itemtype.Spoon
	return newItem
}
