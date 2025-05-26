package hud

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/hajimehoshi/ebiten/v2"
)

func getItemInfo(itemType itemtype.ItemType) item.ItemInfo {
	itemInfo, ok := item.ItemInfos[itemType]
	if !ok {
		fmt.Println("Warning: Could not get itemInfo for: ", itemType.String())
	}
	return itemInfo
}

func getItemIcon(itemInfo item.ItemInfo) *ebiten.Image {
	var icon *ebiten.Image = nil
	if itemIcon, ok := assets.AssetStore.GetImage(itemInfo.IconName); ok {
		icon = itemIcon
	} else {
		if noIcon, ok := assets.AssetStore.GetImage("no_icon"); ok {
			icon = noIcon
		} else {
			fmt.Println("Error Could not load 'no_icon' in newItemBase")
		}
	}
	return icon
}
