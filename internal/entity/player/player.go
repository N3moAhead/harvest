package player

import (
	"fmt"
	"image/color"
	"time"

	"github.com/N3moAhead/harvest/internal/animation"
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
	animationStore  *animation.AnimationStore
}

const (
	IDLE       = "idle"
	UP         = "up"
	UP_RIGHT   = "up_right"
	RIGHT      = "right"
	DOWN_RIGHT = "down_right"
	DOWN       = "down"
	DOWN_LEFT  = "down_left"
	LEFT       = "left"
	UP_LEFT    = "up_left"
)

// The player is currently just drawn as a rectangle.
// TODO: Draw the player with assets
func (p *Player) Draw(screen *ebiten.Image, mapOffsetX float64, mapOffsetY float64) {
	// TODO move the player rect size to the config or somewhere else
	rectSize := 32.0
	var halfRectSize float64 = rectSize / 2

	playerImg := p.animationStore.GetImage()

	if playerImg != nil {
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
	p.animationStore.Update()

	// Update player position
	moveDir := component.Vector2D{X: 0, Y: 0}
	moved := false

	currentAnimation := IDLE
	if inputState.Up {
		moveDir.Y -= 1
		moved = true
		currentAnimation = UP
	}
	if inputState.Down {
		moveDir.Y += 1
		moved = true
		currentAnimation = DOWN
	}
	if inputState.Left {
		moveDir.X -= 1
		moved = true
		currentAnimation = LEFT
	}
	if inputState.Right {
		moveDir.X += 1
		moved = true
		currentAnimation = RIGHT
	}
	if inputState.Up && inputState.Right {
		currentAnimation = UP_RIGHT
	}
	if inputState.Up && inputState.Left {
		currentAnimation = UP_LEFT
	}
	if inputState.Down && inputState.Right {
		currentAnimation = DOWN_RIGHT
	}
	if inputState.Down && inputState.Left {
		currentAnimation = DOWN_LEFT
	}

	if moveDir.Y != 0 && moveDir.X != 0 {
		moveDir = moveDir.Normalize()
		moved = true
	}

	p.animationStore.SetCurrentAnimation(currentAnimation)

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

// TODO: implement a LoadPlayer function to get the saved
// game state from the past
func NewPlayer() *Player {
	baseEntity := entity.NewEntity(
		(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
		(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2,
	)

	store := animation.NewAnimationStore()
	playerImg, ok := assets.AssetStore.GetImage("player")
	playerIdleImg, ok := assets.AssetStore.GetImage("player_idle")
	if ok {
		up, err := animation.NewAnimation(playerImg, 32, 32, 0, 0, 4, 6, true)
		if err == nil {
			store.AddAnimation(UP, up)
		}
		upRight, err := animation.NewAnimation(playerImg, 32, 32, 0, 1*32, 4, 6, true)
		if err == nil {
			store.AddAnimation(UP_RIGHT, upRight)
		}
		right, err := animation.NewAnimation(playerImg, 32, 32, 0, 2*32, 4, 6, true)
		if err == nil {
			store.AddAnimation(RIGHT, right)
		}
		downRight, err := animation.NewAnimation(playerImg, 32, 32, 0, 3*32, 4, 6, true)
		if err == nil {
			store.AddAnimation(DOWN_RIGHT, downRight)
		}
		down, err := animation.NewAnimation(playerImg, 32, 32, 0, 4*32, 4, 6, true)
		if err == nil {
			store.AddAnimation(DOWN, down)
		}
		downLeft, err := animation.NewAnimation(playerImg, 32, 32, 0, 5*32, 4, 6, true)
		if err == nil {
			store.AddAnimation(DOWN_LEFT, downLeft)
		}
		left, err := animation.NewAnimation(playerImg, 32, 32, 0, 6*32, 4, 6, true)
		if err == nil {
			store.AddAnimation(LEFT, left)
		}
		upLeft, err := animation.NewAnimation(playerImg, 32, 32, 0, 7*32, 4, 6, true)
		if err == nil {
			store.AddAnimation(UP_LEFT, upLeft)
		}
		idle, err := animation.NewAnimation(playerIdleImg, 32, 32, 0, 0, 1, 6, true)
		if err == nil {
			store.AddAnimation(IDLE, idle)
		}
		store.SetCurrentAnimation(IDLE)
	} else {
		fmt.Println("Warning: Could not load player img in NewPlayer()")
	}

	p := &Player{
		Entity:          *baseEntity,
		MagnetRadius:    config.INITIAL_PLAYER_MAGNET_RADIUS,
		Speed:           config.INITIAL_PLAYER_SPEED,
		Health:          component.NewHealth(100),
		FacingDirection: component.NewVector2D(0, -1), // Default looks up
		animationStore:  store,
	}
	return p
}
