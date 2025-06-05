package player

import (
	"image/color"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/input"
	"github.com/N3moAhead/harvest/internal/soups"
	"github.com/N3moAhead/harvest/pkg/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type InventoryProvider interface {
	RemoveAllSoups(soupType itemtype.ItemType)
}

type Player struct {
	entity.Entity
	Soups           []soups.Soup
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

func (p *Player) ExtendOrAddSoup(soup *soups.Soup) {
	now := time.Now()

	for i := range p.Soups {
		if p.Soups[i].Type == soup.Type {
			if now.Before(p.Soups[i].ExpiresAt) {
				p.Soups[i].ExpiresAt = p.Soups[i].ExpiresAt.Add(soup.Duration)
			} else {
				p.Soups[i].ExpiresAt = now.Add(soup.Duration)
			}
			// p.Buffs[i].Level++
			return
		}
	}

	newSoup := soups.Soup{
		Type:         soup.Type,
		BuffPerLevel: soup.BuffPerLevel,
		Duration:     soup.Duration,
		ExpiresAt:    now.Add(soup.Duration),
	}
	p.Soups = append(p.Soups, newSoup)
}

func (p *Player) Update(inputState *input.InputState, dt float64, inventory InventoryProvider) { //TODO maybe add inventory to player struct?
	now := time.Now()

	// Update player position
	moveDir := component.Vector2D{X: 0, Y: 0}
	moved := false
	if inputState.Up {
		moveDir.Y -= 1
		moved = true
	}
	if inputState.Down {
		moveDir.Y += 1
		moved = true
	}
	if inputState.Left {
		moveDir.X -= 1
		moved = true
	}
	if inputState.Right {
		moveDir.X += 1
		moved = true
	}

	if moveDir.Y != 0 && moveDir.X != 0 {
		moveDir = moveDir.Normalize()
		moved = true
	}

	// If the player moved update the facingDirection and the player position
	if moved {
		p.Pos = p.Pos.Add(moveDir.Mul(p.Speed))
		p.Pos.X = util.Clamp(p.Pos.X, 16, config.WIDTH_IN_TILES*config.TILE_SIZE-16)
		p.Pos.Y = util.Clamp(p.Pos.Y, 16, config.HEIGHT_IN_TILES*config.TILE_SIZE-16)
		p.FacingDirection = moveDir
	}

	// Update soups
	activeSoups := p.Soups[:0]
	for _, soup := range p.Soups { // filter out expired buffs
		if now.Before(soup.ExpiresAt) {
			activeSoups = append(activeSoups, soup)
		} else {
			inventory.RemoveAllSoups(soup.Type)
		}
	}
	p.Soups = activeSoups

	// reset player stats to default/base values
	p.MagnetRadius = config.INITIAL_PLAYER_MAGNET_RADIUS
	p.Speed = config.INITIAL_PLAYER_SPEED

	for _, soup := range p.Soups {
		buffVal := float64(soup.BuffPerLevel)
		// buffVal := float64(def.BuffPerLevel) * float64(b.Level)
		switch soup.Type {
		case itemtype.MagnetRadiusSoup:
			p.MagnetRadius += buffVal
		case itemtype.SpeedSoup:
			p.Speed += buffVal
		}
	}
}

func (p *Player) Damage(amount float64) {
	hitSound, ok := assets.AssetStore.GetSFXData("player_hit_sound")
	if ok {
		sfxPlayer := assets.AudioContext.NewPlayerFromBytes(hitSound)
		sfxPlayer.Play()
	}
	p.Health.Damage(amount)
}

func (p *Player) Alive() bool {
	return p.Health.HP > 0
}

func (p *Player) HasSoup(soupType itemtype.ItemType) bool {
	for _, soup := range p.Soups {
		if soup.Type == soupType {
			return true
		}
	}
	return false
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
