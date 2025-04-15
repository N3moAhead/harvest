package item

import "github.com/N3moAhead/harvest/internal/itemtype"

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
