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
	Soups           []component.Soup
	Speed           float64
	MagnetRadius    float64
	Health          component.Health
	FacingDirection component.Vector2D
}

// The player is currently just drawn as a rectangle.
// TODO: Draw the player with assets
func (p *Player) Draw(screen *ebiten.Image, mapOffsetX float64, mapOffsetY float64) {
	// TODO move the player rect size to the config or somewhere else
	rectSize := 32.0
	var halfRectSize float64 = rectSize / 2

	if playerImg, ok := assets.AssetStore.GetImage("player"); ok {
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

func (p *Player) GetFacingDirection() component.Vector2D {
	// Fallback if the facing direction is not defined
	if p.FacingDirection.LengthSq() == 0 {
		return component.Vector2D{X: 0, Y: -1}
	}
	// Its important that this direction is normalized
	// but im making sure it is when setting it in the game.go file
	return p.FacingDirection
}

func (p *Player) ExtendOrAddSoup(itemType itemtype.ItemType, inventory *inventory.Inventory) {
	info := itemtype.RetrieveItemInfo(itemType)
	def := info.Soup
	now := time.Now()
	if def == nil {
		return
	}

	inventory.AddSoup(info.Soup.Type)

	for i := range p.Soups {
		if p.Soups[i].Type == def.Type {
			if now.Before(p.Soups[i].ExpiresAt) {
				p.Soups[i].ExpiresAt = p.Soups[i].ExpiresAt.Add(def.Duration)
			} else {
				p.Soups[i].ExpiresAt = now.Add(def.Duration)
			}
			// p.Buffs[i].Level++
			return
		}
	}

	newSoup := component.Soup{
		Type:         def.Type,
		BuffPerLevel: def.BuffPerLevel,
		Duration:     def.Duration,
		ExpiresAt:    now.Add(def.Duration),
	}
	p.Soups = append(p.Soups, newSoup)
}

func (p *Player) Update(dt float64, inventory *inventory.Inventory) { //TODO maybe add inventory to player struct?
	now := time.Now()
	activeSoups := p.Soups[:0]
	for _, b := range p.Soups { // filter out expired buffs ⏲️
		if now.Before(b.ExpiresAt) {
			activeSoups = append(activeSoups, b)
		} else {
			inventory.RemoveSoup(b.Type)
		}
	}
	p.Soups = activeSoups

	// reset player stats to default/base values
	p.MagnetRadius = config.INITIAL_PLAYER_MAGNET_RADIUS
	p.Speed = config.INITIAL_PLAYER_SPEED

	for _, b := range p.Soups {
		buffVal := float64(b.BuffPerLevel)
		// buffVal := float64(def.BuffPerLevel) * float64(b.Level)
		switch b.Type {
		case component.MagnetSoup:
			p.MagnetRadius += buffVal
		case component.SpeedSoup:
			p.Speed += buffVal
		}
	}
}

// TODO: implement a LoadPlayer function to get the saved
// game state from the past
func NewPlayer() *Player {
	baseEntity := entity.NewEntity(
		(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
		(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2,
	)

	p := &Player{
		Entity:          *baseEntity,
		MagnetRadius:    config.INITIAL_PLAYER_MAGNET_RADIUS,
		Speed:           config.INITIAL_PLAYER_SPEED,
		Health:          component.NewHealth(100),
		FacingDirection: component.NewVector2D(0, -1), // Default looks up
	}
	return p
}
