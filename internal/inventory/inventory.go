package inventory

import (
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/itemtype"
)

type Inventory struct {
	Vegtables map[itemtype.ItemType]int  // Mapping the item type to the amount
	Soups     map[component.SoupType]int // same here, but for soups
}

func (i *Inventory) AddVegtable(itemType itemtype.ItemType) {
	i.Vegtables[itemType]++
}

func (i *Inventory) AddSoup(soupType component.SoupType) {
	i.Soups[soupType]++
}

func (i *Inventory) RemoveSoup(soupType component.SoupType) {
	i.Soups[soupType]--
	if i.Soups[soupType] <= 0 {
		delete(i.Soups, soupType)
	}
}

func NewInventory() *Inventory {
	return &Inventory{
		Vegtables: make(map[itemtype.ItemType]int),
		Soups:     make(map[component.SoupType]int),
	}
}
