package item

import (
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/itemtype"
)

func NewCarrot(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY)
	newItem.Type = itemtype.Carrot
	return newItem
}

func NewPotato(posX float64, posY float64) *Item {
	newItem := newItemBase(posX, posY)
	newItem.Type = itemtype.Potato
	return newItem
}

func NewSoup(x, y float64, buff component.BuffType) *Item {
	itm := newItemBase(x, y)
	switch buff {
	case component.MagnetRadius:
		itm.Type = itemtype.MagnetRadiusSoup
	case component.Damage:
		itm.Type = itemtype.DamageSoup
	case component.Speed:
		itm.Type = itemtype.SpeedSoup
	default:
		itm.Type = itemtype.Undefined
	}
	return itm
}

// func NewSoup(posX, posY float64, buffType component.BuffType, level int) *Item {
// 	itm := newItemBase(posX, posY)
// 	itm.Type = itemtype.Soup
// 	def := component.BuffDefs[buffType]
// 	itm.Buff = &component.Buff{
// 		Type:      buffType,
// 		Level:     level,
// 		ExpiresAt: time.Now().Add(def.Duration),
// 	}
// 	return itm
// }
