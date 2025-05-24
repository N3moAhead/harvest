package item

import (
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/soups"
)

type ItemInfo struct {
	DisplayName string
	Category    itemtype.ItemCategory
	Soup        *soups.Soup
}

// Saves detailed information for each item type
var ItemInfos = map[itemtype.ItemType]ItemInfo{
	itemtype.Potato: {
		DisplayName: "Potato",
		Category:    itemtype.CategoryVegetable,
		Soup:        nil,
	},
	itemtype.Carrot: {
		DisplayName: "Carrot",
		Category:    itemtype.CategoryVegetable,
		Soup:        nil,
	},
	itemtype.Spoon: {
		DisplayName: "Spoon",
		Category:    itemtype.CategoryWeapon,
		Soup:        nil,
	},
	itemtype.DamageSoup: {
		DisplayName: "Damage Soup",
		Category:    itemtype.CategorySoup,
		Soup:        soups.Definitions[itemtype.DamageSoup],
	},
	itemtype.MagnetRadiusSoup: {
		DisplayName: "Magnet Soup",
		Category:    itemtype.CategorySoup,
		Soup:        soups.Definitions[itemtype.MagnetRadiusSoup],
	},
	itemtype.SpeedSoup: {
		DisplayName: "Speed Soup",
		Category:    itemtype.CategorySoup,
		Soup:        soups.Definitions[itemtype.SpeedSoup],
	},
}
