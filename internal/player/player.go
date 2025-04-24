package player

import (
	"image/color"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/inventory"
	"github.com/N3moAhead/harvest/internal/itemtype"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	entity.Entity
	Buffs        []component.Buff
	Speed        float64
	MagnetRadius float64
	Health       component.Health
}

// The player is currently just drawn as a rectangle.
// TODO: Draw the player with assets
func (p *Player) Draw(screen *ebiten.Image, assetStore *assets.Store, mapOffsetX float64, mapOffsetY float64) {
	// TODO move the player rect size to the config or somewhere else
	rectSize := 32.0
	var halfRectSize float64 = rectSize / 2

	if playerImg, ok := assetStore.GetImage("player"); ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.Pos.X-mapOffsetX-halfRectSize, p.Pos.Y-mapOffsetY-halfRectSize)
		screen.DrawImage(playerImg, op)
	} else {
		// Fallback if player image could not be loaded
		vector.DrawFilledRect(
			screen,
			float32(p.Pos.X)-float32(mapOffsetX)-float32(halfRectSize),
			float32(p.Pos.Y)-float32(mapOffsetY)-float32(halfRectSize),
			float32(rectSize),
			float32(rectSize),
			color.RGBA{R: 255, G: 255, B: 255, A: 255},
			true,
		)
	}
}

func (p *Player) ExtendOrAddBuff(itemType itemtype.ItemType, inventory *inventory.Inventory) {
	info := itemtype.RetrieveItemInfo(itemType)
	def := info.Buff
	now := time.Now()
	if def == nil {
		return
	}

	inventory.AddSoup(info.Buff.Type)

	for i := range p.Buffs {
		if p.Buffs[i].Type == def.Type {
			if now.Before(p.Buffs[i].ExpiresAt) {
				p.Buffs[i].ExpiresAt = p.Buffs[i].ExpiresAt.Add(def.Duration)
			} else {
				p.Buffs[i].ExpiresAt = now.Add(def.Duration)
			}
			// p.Buffs[i].Level++
			return
		}
	}

	newBuff := component.Buff{
		Type:         def.Type,
		BuffPerLevel: def.BuffPerLevel,
		Duration:     def.Duration,
		ExpiresAt:    now.Add(def.Duration),
	}
	p.Buffs = append(p.Buffs, newBuff)
}

func (p *Player) Update(dt float64, inventory *inventory.Inventory) { //TODO maybe add inventory to player struct?
	now := time.Now()
	aliveBuffs := p.Buffs[:0]
	for _, b := range p.Buffs { // filter out expired buffs ⏲️
		if now.Before(b.ExpiresAt) {
			aliveBuffs = append(aliveBuffs, b)
		} else {
			inventory.RemoveSoup(b.Type)
		}
	}
	p.Buffs = aliveBuffs

	// reset player stats to default/base values
	p.MagnetRadius = config.INITIAL_PLAYER_MAGNET_RADIUS
	p.Speed = config.INITIAL_PLAYER_SPEED

	for _, b := range p.Buffs {
		buffVal := float64(b.BuffPerLevel)
		// buffVal := float64(def.BuffPerLevel) * float64(b.Level)
		switch b.Type {
		case component.MagnetRadiusBuff:
			p.MagnetRadius += buffVal
		case component.SpeedBuff:
			p.Speed += buffVal
		}
	}
}

// TODO: implement a LoadPlayer function to get the saved
// game state from the past
func NewPlayer() *Player {
	baseEntity := entity.NewEntity(config.SCREEN_WIDTH/2, config.SCREEN_HEIGHT/2)

	p := &Player{
		Entity:       *baseEntity,
		MagnetRadius: config.INITIAL_PLAYER_MAGNET_RADIUS,
		Speed:        config.INITIAL_PLAYER_SPEED,
		Health:       component.NewHealth(100),
	}
	return p
}
