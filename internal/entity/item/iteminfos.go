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
		IconName:    "potato_icon",
	},
	itemtype.Carrot: {
		DisplayName: "Carrot",
		Category:    itemtype.CategoryVegetable,
		Soup:        nil,
		IconName:    "carrot_icon",
	},
	itemtype.Onion: {
		DisplayName: "Onion",
		Category:    itemtype.CategoryVegetable,
		Soup:        nil,
		IconName:    "onion_icon",
	},
	itemtype.Spoon: {
		DisplayName: "Spoon",
		Category:    itemtype.CategoryWeapon,
		Soup:        nil,
		IconName:    "spoon_icon",
	},
	itemtype.ThrowingKnifes: {
		DisplayName: "Throwing Knifes",
		Category:    itemtype.CategoryWeapon,
		Soup:        nil,
		IconName:    "throwing_knifes_icon",
	},
	itemtype.RollingPin: {
		DisplayName: "Rolling Pin",
		Category:    itemtype.CategoryWeapon,
		Soup:        nil,
		IconName:    "rolling_pin_icon",
	},
	itemtype.Thermalmixer: {
		DisplayName: "Thermalmixer",
		Category:    itemtype.CategoryWeapon,
		Soup:        nil,
		IconName:    "thermalmixer_icon",
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
