package itemtype

import (
	"fmt"
	"time"

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

type ItemInfo struct {
	DisplayName string
	Category    ItemCategory
	Buff        *component.Buff
}

// Saves meta information for each item type
// To help us connect Categories and ItemTypes
var ItemInfos = map[ItemType]ItemInfo{
	Potato: {
		DisplayName: "Potato",
		Category:    CategoryVegetable,
		Buff:        nil,
	},
	Carrot: {
		DisplayName: "Carrot",
		Category:    CategoryVegetable,
		Buff:        nil,
	},
	DamageSoup: {
		DisplayName: "Damage Soup",
		Category:    CategorySoup,
		Buff: &component.Buff{
			Type:         component.DamageBuff,
			BuffPerLevel: 5,
			Duration:     5 * time.Second,
		},
	},
	MagnetRadiusSoup: {
		DisplayName: "Magnet Soup",
		Category:    CategorySoup,
		Buff: &component.Buff{
			Type:         component.MagnetRadiusBuff,
			BuffPerLevel: 200,
			Duration:     2 * time.Second,
		},
	},
	SpeedSoup: {
		DisplayName: "Speed Soup",
		Category:    CategorySoup,
		Buff: &component.Buff{
			Type:         component.SpeedBuff,
			BuffPerLevel: 20,
			Duration:     2 * time.Second,
		},
	},
}

func RetrieveItemInfo(t ItemType) (info ItemInfo) {
	if info, ok := ItemInfos[t]; ok {
		return info
	}
	return ItemInfo{
		DisplayName: CategoryUndefined.String(),
		Category:    CategoryUndefined,
		Buff:        nil,
	}
}

func (it ItemType) String() string {
	if info, ok := ItemInfos[it]; ok {
		return info.DisplayName
	}
	return fmt.Sprintf("ItemType(%d)", it)
}

func (it ItemType) Category() ItemCategory {
	if info, ok := ItemInfos[it]; ok {
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
