package item

import "github.com/N3moAhead/harvest/internal/entity/item/itemtype"

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

/// --- Weapons ---

func NewSpoon(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY)
	newItem.Type = itemtype.Spoon
	return newItem
}

/// --- Soups ---

func NewSoup(x, y float64, typeBuff itemtype.ItemType) *Item {
	newItem := newItemBase(x, y)
	newItem.Type = typeBuff
	return newItem
}
