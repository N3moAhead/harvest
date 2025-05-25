package hud

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ItemDisplayInterface interface {
	UpdateItemDisplay(itemType itemtype.ItemType)
}

// The ItemDisplay is an ui element that can display the
// amount of an item in the inventory.
// WARNING: Has to be child of vegtable inventory! Otherwise make sure
// to call the extra function UpdateValues correctly!
type ItemDisplay struct {
	ui.BaseElement
	Inv          *inventory.Inventory
	ItemFrameImg *ebiten.Image
	amount       int
	itemType     itemtype.ItemType
}

func NewItemDisplay(x, y float64, inv *inventory.Inventory, frameImage *ebiten.Image) *ItemDisplay {
	return &ItemDisplay{
		BaseElement:  *ui.NewBaseElement(x, y, config.ITEM_FRAME_SIZE, config.ITEM_FRAME_SIZE),
		Inv:          inv,
		ItemFrameImg: frameImage,
		amount:       0,
		itemType:     itemtype.Undefined,
	}
}

func (v *ItemDisplay) Update(input *ui.InputState) {
	v.BaseElement.Update(input)
}

func (v *ItemDisplay) UpdateItemDisplay(itemType itemtype.ItemType) {
	if itemType != itemtype.Undefined {
		v.itemType = itemType
		switch itemType.Category() {
		case itemtype.CategoryVegetable:
			v.amount = v.Inv.Vegetables[itemType]
		case itemtype.CategorySoup:
			v.amount = v.Inv.Soups[itemType]
		default:
			fmt.Println("Warning: Unhandeld Category in function UpdateItemDisplay")
		}
	} else {
		v.itemType = itemtype.Undefined
		v.amount = 0
	}
}

func (v *ItemDisplay) Draw(screen *ebiten.Image) {
	// Drawing the frame
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(v.X, v.Y)
	screen.DrawImage(v.ItemFrameImg, op)
	if v.itemType != itemtype.Undefined && v.amount != 0 {
		v.drawItemDisplay(screen)
	}
}

func (v *ItemDisplay) drawItemDisplay(screen *ebiten.Image) {
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

var _ ui.UIElement = (*ItemDisplay)(nil)
var _ ItemDisplayInterface = (*ItemDisplay)(nil)
