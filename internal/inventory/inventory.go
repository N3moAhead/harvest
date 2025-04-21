package inventory

import (
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/itemtype"
)

type Inventory struct {
	Vegtables map[itemtype.ItemType]int  // Mapping the item type to the amount
	Soups     map[component.BuffType]int // same here, but for soups
}

func (i *Inventory) AddVegtable(itemType itemtype.ItemType) {
	i.Vegtables[itemType]++
}

func (i *Inventory) AddSoup(buffType component.BuffType) {
	i.Soups[buffType]++
}

func (i *Inventory) RemoveSoup(buffType component.BuffType) {
	i.Soups[buffType]--
	if i.Soups[buffType] <= 0 {
		delete(i.Soups, buffType)
	}
}

func NewInventory() *Inventory {
	return &Inventory{
		Vegtables: make(map[itemtype.ItemType]int),
		Soups:     make(map[component.BuffType]int),
	}
}
