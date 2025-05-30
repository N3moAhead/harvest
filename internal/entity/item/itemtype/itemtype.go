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
	Spoon
	ThrowingKnifes
	RollingPin
	Thermalmixer
	DamageSoup
	MagnetRadiusSoup
	SpeedSoup
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
	case Potato, Carrot, Onion, Leek, Cabbage:
		return CategoryVegetable
	case RollingPin, ThrowingKnifes, Spoon, Thermalmixer:
		return CategoryWeapon
	case DamageSoup, MagnetRadiusSoup, SpeedSoup:
		return CategorySoup
	default:
		return CategoryUndefined
	}
}

// Implementing the sort interface for ItemType
type ByItemType []ItemType

func (a ByItemType) Len() int           { return len(a) }
func (a ByItemType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByItemType) Less(i, j int) bool { return a[i] < a[j] }
