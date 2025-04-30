package itemtype

import (
	"fmt"
	"time"
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
	Spoon
	DamageSoup
	MagnetRadiusSoup
	SpeedSoup
)

type Soup struct {
	Type         ItemType // maybe set to itemtype hmm :/
	BuffPerLevel float32
	// Level        int
	Duration  time.Duration
	ExpiresAt time.Time
}

type ItemInfo struct {
	DisplayName string
	Category    ItemCategory
	Soup        *Soup
}

// Saves meta information for each item type
// To help us connect Categories and ItemTypes
var ItemInfos = map[ItemType]ItemInfo{
	Potato: {
		DisplayName: "Potato",
		Category:    CategoryVegetable,
		Soup:        nil,
	},
	Carrot: {
		DisplayName: "Carrot",
		Category:    CategoryVegetable,
		Soup:        nil,
	},
	Spoon: {
		DisplayName: "Spoon",
		Category:    CategoryVegetable,
		Soup:        nil,
	},
	DamageSoup: {
		DisplayName: "Damage Soup",
		Category:    CategorySoup,
		Soup: &Soup{
			Type:         DamageSoup,
			BuffPerLevel: 5,
			Duration:     5 * time.Second,
		},
	},
	MagnetRadiusSoup: {
		DisplayName: "Magnet Soup",
		Category:    CategorySoup,
		Soup: &Soup{
			Type:         MagnetRadiusSoup,
			BuffPerLevel: 200,
			Duration:     2 * time.Second,
		},
	},
	SpeedSoup: {
		DisplayName: "Speed Soup",
		Category:    CategorySoup,
		Soup: &Soup{
			Type:         SpeedSoup,
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
		Soup:        nil,
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
