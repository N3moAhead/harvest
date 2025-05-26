package item

import (
	"fmt"
	"image/color"

	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Item struct {
	entity.Entity
	Type itemtype.ItemType
}

func (i *Item) Update(player *player.Player) (itemPickedUp bool) {
	diff := player.Pos.Sub(i.Pos)
	len := diff.Len() // the distance from player to item

	if len < config.PLAYER_PICKUP_RADIUS {
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
	case i.Type == itemtype.Spoon:
		itemColor = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	case i.Type == itemtype.RollingPin:
		itemColor = color.RGBA{R: 145, G: 128, B: 128, A: 255}
	case i.Type == itemtype.DamageSoup:
		itemColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	case i.Type == itemtype.MagnetRadiusSoup:
		itemColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	case i.Type == itemtype.SpeedSoup:
		itemColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
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

func (i *Item) RetrieveItemInfo() ItemInfo {
	if info, ok := ItemInfos[i.Type]; ok {
		return info
	}
	return ItemInfo{
		DisplayName: fmt.Sprintf("ItemType(%d)", i.Type),
		Category:    itemtype.CategoryUndefined,
		Soup:        nil,
	}
}

func (i *Item) DisplayName() string {
	return i.RetrieveItemInfo().DisplayName
}

func (i *Item) CategoryOf() itemtype.ItemCategory {
	return i.RetrieveItemInfo().Category
}

func (i *Item) IsVegetable() bool {
	return i.CategoryOf() == itemtype.CategoryVegetable
}
func (i *Item) IsWeapon() bool {
	return i.CategoryOf() == itemtype.CategoryWeapon
}

func (i *Item) IsSoup() bool {
	return i.CategoryOf() == itemtype.CategorySoup
}
