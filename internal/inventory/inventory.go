package inventory

import "github.com/N3moAhead/harvest/internal/itemtype"

type Inventory struct {
	Vegtables map[itemtype.ItemType]int // Mapping the item type to the amount
}

func (i *Inventory) AddVegtable(itemType itemtype.ItemType) {
	i.Vegtables[itemType]++
}

func NewInventory() *Inventory {
	return &Inventory{
		Vegtables: make(map[itemtype.ItemType]int),
	}
}
