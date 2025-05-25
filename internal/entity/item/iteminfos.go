package item

import (
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/soups"
)

type ItemInfo struct {
	DisplayName string
	Category    itemtype.ItemCategory
	Soup        *soups.Soup
	IconName    string
}

// Saves detailed information for each item type
var ItemInfos = map[itemtype.ItemType]ItemInfo{
	itemtype.Potato: {
		DisplayName: "Potato",
		Category:    itemtype.CategoryVegetable,
		Soup:        nil,
		IconName:    "",
	},
	itemtype.Carrot: {
		DisplayName: "Carrot",
		Category:    itemtype.CategoryVegetable,
		Soup:        nil,
		IconName:    "carrot_icon",
	},
	itemtype.Spoon: {
		DisplayName: "Spoon",
		Category:    itemtype.CategoryWeapon,
		Soup:        nil,
		IconName:    "",
	},
	itemtype.DamageSoup: {
		DisplayName: "Damage Soup",
		Category:    itemtype.CategorySoup,
		Soup:        soups.Definitions[itemtype.DamageSoup],
		IconName:    "",
	},
	itemtype.MagnetRadiusSoup: {
		DisplayName: "Magnet Soup",
		Category:    itemtype.CategorySoup,
		Soup:        soups.Definitions[itemtype.MagnetRadiusSoup],
		IconName:    "",
	},
	itemtype.SpeedSoup: {
		DisplayName: "Speed Soup",
		Category:    itemtype.CategorySoup,
		Soup:        soups.Definitions[itemtype.SpeedSoup],
		IconName:    "",
	},
}
