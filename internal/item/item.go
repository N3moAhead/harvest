package item

import (
	"image/color"
	"time"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/inventory"
	"github.com/N3moAhead/harvest/internal/itemtype"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Item struct {
	entity.Entity
	Type itemtype.ItemType
}

func (i *Item) Update(player *player.Player, inventory *inventory.Inventory) (removeItem bool) {
	diff := player.Pos.Sub(i.Pos)
	len := diff.Len() // the distance from player to item

	if len < config.PLAYER_PICKUP_RADIUS {
		info := itemtype.ItemInfo(i.Type)
		if info.Category == itemtype.CategorySoup {
			inventory.AddSoup(info.Buff)
			expiry := time.Now().Add(component.BuffDefs[info.Buff].Duration)
			player.ExtendOrAddBuff(component.Buff{
				Type:      info.Buff,
				Level:     1,
				ExpiresAt: expiry,
			})
		} else {
			// Picking up the item into the inventory
			inventory.AddVegtable(i.Type)
		}
		return true
	}
	if len < player.MagnetRadius {
		dir := diff.Normalize() // Direction towards the player
		// Calculating movement towards the player with the correct speed
		moveStep := dir.Mul(config.PLAYER_MAGNET_ATTRACTION_SPEED)
		// Avoid overshooting the target
		if moveStep.Len() > len {
			i.Pos = player.Pos
		} else {
			i.Pos = i.Pos.Add(moveStep)
		}

	}
	return false
}

func (i *Item) Draw(screen *ebiten.Image, mapOffsetX float64, mapOffsetY float64) {
	var itemColor color.RGBA
	switch {
	case i.Type == itemtype.Carrot:
		itemColor = color.RGBA{R: 230, G: 126, B: 34, A: 255}
	case i.Type == itemtype.Potato:
		itemColor = color.RGBA{R: 183, G: 146, B: 104, A: 255}
	case i.Type == itemtype.DamageSoup:
		itemColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	case i.Type == itemtype.MagnetRadiusSoup:
		itemColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	default:
		itemColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	vector.DrawFilledRect(
		screen,
		float32(i.Pos.X)-float32(mapOffsetX)-4.0,
		float32(i.Pos.Y)-float32(mapOffsetY)-4.0,
		8.0,
		8.0,
		itemColor,
		true,
	)
}

func newItemBase(posX float64, posY float64) *Item {
	baseClass := entity.NewEntity(posX, posY)
	return &Item{
		Entity: *baseClass,
		Type:   itemtype.Undefined,
	}
}
