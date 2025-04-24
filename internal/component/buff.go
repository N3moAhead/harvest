// internal/component/buff.go
package component // change package?

import "time"

type BuffType int // TODO adjust naming either to BUFF or SOUPS

const (
	SpeedBuff BuffType = iota
	MagnetRadiusBuff
	DamageBuff  // example
	DefenseBuff // example
	// …
)

type Buff struct {
	Type         BuffType // maybe set to itemtype hmm :/
	BuffPerLevel float32
	// Level        int
	Duration  time.Duration
	ExpiresAt time.Time
}
