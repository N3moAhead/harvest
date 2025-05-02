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
	DamageSoup
	MagnetRadiusSoup
	SpeedSoup
)
