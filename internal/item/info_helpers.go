package item

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/itemtype"
)

// RetrieveItemInfo returns the ItemInfo for a given item type
func RetrieveItemInfo(it itemtype.ItemType) ItemInfo {
	if info, ok := ItemInfos[it]; ok {
		return info
	}
	return ItemInfo{
		DisplayName: fmt.Sprintf("ItemType(%d)", it),
		Category:    itemtype.CategoryUndefined,
		Soup:        nil,
	}
}

func DisplayName(it itemtype.ItemType) string {
	return RetrieveItemInfo(it).DisplayName
}

func CategoryOf(it itemtype.ItemType) itemtype.ItemCategory {
	return RetrieveItemInfo(it).Category
}

func IsVegetable(it itemtype.ItemType) bool {
	return CategoryOf(it) == itemtype.CategoryVegetable
}

func IsWeapon(it itemtype.ItemType) bool {
	return CategoryOf(it) == itemtype.CategoryWeapon
}

func IsSoup(it itemtype.ItemType) bool {
	return CategoryOf(it) == itemtype.CategorySoup
}
