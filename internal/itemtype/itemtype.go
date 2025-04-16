package itemtype

import "fmt"

type ItemCategory int

const (
	CategoryUndefined ItemCategory = iota
	CategoryVegetable
	CategoryWeopon
	CategorySoup // Or could also be named buff
)

func (ic ItemCategory) String() string {
	switch ic {
	case CategoryUndefined:
		return "Undefined"
	case CategoryVegetable:
		return "Vegtable"
	case CategoryWeopon:
		return "Weopon"
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
)

// Saves meta information for each item type
// To help us connect Categories and ItemTypes
var itemInfo = map[ItemType]struct {
	DisplayName string
	Category    ItemCategory
}{
	Potato: {"Potato", CategoryVegetable},
	Carrot: {"Carrot", CategoryVegetable},
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

func (it ItemType) IsWeopon() bool {
	return it.Category() == CategoryWeopon
}

func (it ItemType) IsSoup() bool {
	return it.Category() == CategorySoup
}
