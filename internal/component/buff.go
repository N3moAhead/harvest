// internal/component/buff.go
package component // change package?

import "time"

type BuffType int // TODO adjust naming either to BUFF or SOUPS

const (
	Speed BuffType = iota
	MagnetRadius
	Damage  // example
	Defense // example
	// …
)

type Buff struct {
	Type      BuffType
	Level     int
	ExpiresAt time.Time
}

type BuffDefinition struct {
	BuffPerLevel float32
	Duration     time.Duration
}

var BuffDefs = map[BuffType]BuffDefinition{
	Speed: {
		BuffPerLevel: 0.5,
		Duration:     2 * time.Second,
	},
	MagnetRadius: {
		BuffPerLevel: 0.5,
		Duration:     2 * time.Second,
	},
	Damage: { // example
		BuffPerLevel: 1,
		Duration:     5 * time.Second,
	},
}
