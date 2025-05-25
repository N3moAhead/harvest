package hud

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type VegtableDisplayInterface interface {
	UpdateVegtableDisplay(itemType itemtype.ItemType)
}

// The Vegtable Display is an ui element that can display the
// amount of an item in the inventory.
// WARNING: Has to be child of vegtable inventory! Otherwise make sure
// to call the extra function UpdateValues correctly!
type VegtableDisplay struct {
	ui.BaseElement
	Inv          *inventory.Inventory
	ItemFrameImg *ebiten.Image
	amount       int
	itemType     itemtype.ItemType
}

func NewVegtableDisplay(x, y float64, inv *inventory.Inventory) *VegtableDisplay {
	itemFrameImg, ok := assets.AssetStore.GetImage("item_frame")
	if !ok {
		fmt.Println("Warning: Could not load image item_frame in NewVegtableDisplay")
	}
	return &VegtableDisplay{
		BaseElement:  *ui.NewBaseElement(x, y, config.ITEM_FRAME_SIZE, config.ITEM_FRAME_SIZE),
		Inv:          inv,
		ItemFrameImg: itemFrameImg,
		amount:       0,
		itemType:     itemtype.Undefined,
	}
}

func (v *VegtableDisplay) Update(input *ui.InputState) {
	v.BaseElement.Update(input)
}

func (v *VegtableDisplay) UpdateVegtableDisplay(itemType itemtype.ItemType) {
	if itemType != itemtype.Undefined {
		v.itemType = itemType
		v.amount = v.Inv.Vegetables[itemType]
	} else {
		v.itemType = itemtype.Undefined
		v.amount = 0
	}
}

func (v *VegtableDisplay) Draw(screen *ebiten.Image) {
	// Drawing the frame
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(v.X, v.Y)
	screen.DrawImage(v.ItemFrameImg, op)
	if v.itemType != itemtype.Undefined && v.amount != 0 {
		v.drawItemDisplay(screen)
	}
}

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

func (v *VegtableDisplay) drawItemDisplay(screen *ebiten.Image) {
	itemInfo := getItemInfo(v.itemType)
	itemIcon := getItemIcon(itemInfo)
	bounds := itemIcon.Bounds()
	offsetX := (float64(config.ITEM_FRAME_SIZE - bounds.Dx())) / 2
	offsetY := (float64(config.ITEM_FRAME_SIZE - bounds.Dy())) / 2
	drawX := v.X + offsetX
	drawY := v.Y + offsetY
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(drawX, drawY)
	screen.DrawImage(itemIcon, op)
	// TODO replace the debug print with a real font!!
	// Adding some padding here to move the the number close to the bottom right
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", v.amount), int(drawX+7), int(drawY+5))
}

var _ ui.UIElement = (*VegtableDisplay)(nil)
var _ VegtableDisplayInterface = (*VegtableDisplay)(nil)
