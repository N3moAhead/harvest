package soups

import (
	"time"

	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
)

type Soup struct {
	Type         itemtype.ItemType
	BuffPerLevel float32
	// Level        int
	Duration  time.Duration
	ExpiresAt time.Time
}

var Definitions = map[itemtype.ItemType]*Soup{
	itemtype.DamageSoup: {
		Type:         itemtype.DamageSoup,
		BuffPerLevel: 5,
		Duration:     5 * time.Second,
	},
	itemtype.MagnetRadiusSoup: {
		Type:         itemtype.MagnetRadiusSoup,
		BuffPerLevel: 200,
		Duration:     2 * time.Second,
	},
	itemtype.SpeedSoup: {
		Type:         itemtype.SpeedSoup,
		BuffPerLevel: 20,
		Duration:     2 * time.Second,
	},
}
