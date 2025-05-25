package item

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
)

type Item struct {
	entity.Entity
	Type itemtype.ItemType
	Icon *ebiten.Image
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
	scaleFactor := config.ICON_ON_MAP_RENDER_SIZE / config.ICON_SIZE
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleFactor, scaleFactor)
	scaledIconHalfWidth := (config.ICON_SIZE * scaleFactor) / 2.0
	scaledIconHalfHeight := (config.ICON_SIZE * scaleFactor) / 2.0
	drawX := i.Pos.X - mapOffsetX - scaledIconHalfWidth
	drawY := i.Pos.Y - mapOffsetY - scaledIconHalfHeight
	op.GeoM.Translate(drawX, drawY)
	screen.DrawImage(i.Icon, op)
}

func newItemBase(posX float64, posY float64, icon *ebiten.Image) *Item {
	if icon == nil {
		if noIcon, ok := assets.AssetStore.GetImage("no_icon"); ok {
			icon = noIcon
		} else {
			fmt.Println("Error Could not load 'no_icon' in newItemBase")
		}
	}
	baseClass := entity.NewEntity(posX, posY)
	return &Item{
		Entity: *baseClass,
		Type:   itemtype.Undefined,
		Icon:   icon,
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
