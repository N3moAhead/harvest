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
	Spoon
	RollingPin
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
	case Spoon:
		return "Spoon"
	case RollingPin:
		return "Rolling Pin"
	case DamageSoup:
		return "Damage Soup"
	case MagnetRadiusSoup:
		return "Magnet Radius Soup"
	case SpeedSoup:
		return "Speed Soup"
	default:
		return "Unknown"
	}
}
