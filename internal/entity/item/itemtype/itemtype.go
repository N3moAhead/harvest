package itemtype

type ItemCategory int

const (
	CategoryUndefined ItemCategory = iota
	CategoryVegetable
	CategoryWeapon
	CategorySoup
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
	Cabbage
	Onion
	Leek
	Radish
	Spoon
	ThrowingKnifes
	RollingPin
	Thermalmixer
	DamageSoup
	MagnetRadiusSoup
	SpeedSoup
	MaxItemType // This should always be the last item type
)

func (it ItemType) String() string {
	switch it {
	case Undefined:
		return "Undefined"
	case Potato:
		return "Potato"
	case Carrot:
		return "Carrot"
	case Cabbage:
		return "Cabbage"
	case Onion:
		return "Onion"
	case Leek:
		return "Leek"
	case Radish:
		return "Radish"
	case Spoon:
		return "Spoon"
	case RollingPin:
		return "Rolling Pin"
	case Thermalmixer:
		return "Thermalmixer"
	case DamageSoup:
		return "Damage Soup"
	case MagnetRadiusSoup:
		return "Magnet Radius Soup"
	case ThrowingKnifes:
		return "Throwing Knifes"
	case SpeedSoup:
		return "Speed Soup"
	default:
		return "Unknown"
	}
}

func (it ItemType) Category() ItemCategory {
	switch it {
	case Potato, Carrot, Onion, Leek, Cabbage, Radish:
		return CategoryVegetable
	case RollingPin, ThrowingKnifes, Spoon, Thermalmixer:
		return CategoryWeapon
	case DamageSoup, MagnetRadiusSoup, SpeedSoup:
		return CategorySoup
	default:
		return CategoryUndefined
	}
}

func GetItemTypesByCategory(category ItemCategory) []ItemType {
	var items []ItemType
	for i := Undefined + 1; i < MaxItemType; i++ {
		item := ItemType(i)
		if item.Category() == category {
			items = append(items, item)
		}
	}
	return items
}

// Implementing the sort interface for ItemType
type ByItemType []ItemType

func (a ByItemType) Len() int           { return len(a) }
func (a ByItemType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByItemType) Less(i, j int) bool { return a[i] < a[j] }
