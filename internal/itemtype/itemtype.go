package itemtype

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/component"
)

type ItemCategory int

const (
	CategoryUndefined ItemCategory = iota
	CategoryVegetable
	CategoryWeapon
	CategorySoup // Or could also be named buff
)

func (ic ItemCategory) String() string {
	switch ic {
	case CategoryUndefined:
		return "Undefined"
	case CategoryVegetable:
		return "Vegtable"
	case CategoryWeapon:
		return "Weapon"
	case CategorySoup:
		return "Soup"
	default:
		return "Unknown"
	}
}

type ItemType int

const (
	Undefined ItemType = iota
	Potato
	Carrot
	DamageSoup
	MagnetRadiusSoup
	SpeedSoup
)

// Saves meta information for each item type
// To help us connect Categories and ItemTypes
var itemInfo = map[ItemType]struct {
	DisplayName string
	Category    ItemCategory
	Buff        component.BuffType
}{
	Potato:           {"Potato", CategoryVegetable, 0},
	Carrot:           {"Carrot", CategoryVegetable, 0},
	DamageSoup:       {"Damge enhancing Soup", CategorySoup, component.Damage},
	MagnetRadiusSoup: {"Soup", CategorySoup, component.MagnetRadius},
	SpeedSoup:        {"Soup", CategorySoup, component.Speed},
}

func ItemInfo(t ItemType) (info struct {
	DisplayName string
	Category    ItemCategory
	Buff        component.BuffType
}) {
	if info, ok := itemInfo[t]; ok {
		return info
	}
	// default value if not found??
	return struct {
		DisplayName string
		Category    ItemCategory
		Buff        component.BuffType
	}{CategoryUndefined.String(), CategoryUndefined, 0}
}

func (it ItemType) String() string {
	if info, ok := itemInfo[it]; ok {
		return info.DisplayName
	}
	return fmt.Sprintf("ItemType(%d)", it)
}

func (it ItemType) Category() ItemCategory {
	if info, ok := itemInfo[it]; ok {
		return info.Category
	}
	return CategoryUndefined
}

func (it ItemType) IsVegtable() bool {
	return it.Category() == CategoryVegetable
}

func (it ItemType) IsWeapon() bool {
	return it.Category() == CategoryWeapon
}

func (it ItemType) IsSoup() bool {
	return it.Category() == CategorySoup
}
