// internal/component/buff.go
package component // change package?

import "time"

type SoupType int // TODO adjust naming either to BUFF or SOUPS

const (
	SpeedSoup SoupType = iota
	MagnetSoup
	DamageSoup  // example
	DefenseSoup // example
	// …
)

type Soup struct {
	Type         SoupType // maybe set to itemtype hmm :/
	BuffPerLevel float32
	// Level        int
	Duration  time.Duration
	ExpiresAt time.Time
}
