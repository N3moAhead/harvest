package hud

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type WeaponFrameInterface interface {
	UpdateWeaponFrameValues(itemType itemtype.ItemType, description string, level int)
}

type WeaponFrame struct {
	ui.BaseElement
	Inv          *inventory.Inventory
	ItemFrameImg *ebiten.Image
	itemType     itemtype.ItemType
	isHovered    bool
	description  string
	level        int
}

func NewWeaponFrame(x, y float64, inv *inventory.Inventory) *WeaponFrame {
	weaponFrame, ok := assets.AssetStore.GetImage("weapon_item_frame")
	if !ok {
		fmt.Println("Warning: Unable to laod weapon_item_frame in NewWeaponFrame")
	}
	return &WeaponFrame{
		BaseElement:  *ui.NewBaseElement(x, y, config.ITEM_FRAME_SIZE, config.ITEM_FRAME_SIZE),
		Inv:          inv,
		ItemFrameImg: weaponFrame,
		itemType:     itemtype.Undefined,
		description:  "",
	}
}

func (v *WeaponFrame) Update(input *ui.InputState) {
	v.isHovered = v.IsMouseOver(input.MouseX, input.MouseY)
	v.BaseElement.Update(input)
}

func (v *WeaponFrame) UpdateWeaponFrameValues(itemType itemtype.ItemType, description string, level int) {
	if itemType != itemtype.Undefined {
		v.itemType = itemType
		v.description = description
		v.level = level
	} else {
		v.itemType = itemtype.Undefined
		v.description = ""
		v.level = 0
	}
}

func (v *WeaponFrame) Draw(screen *ebiten.Image) {
	// Drawing the frame
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(v.X, v.Y)
	screen.DrawImage(v.ItemFrameImg, op)
	if v.itemType != itemtype.Undefined {
		v.drawItemDisplay(screen)
		if v.isHovered {
			ebitenutil.DebugPrintAt(screen, v.description, int(v.X), int(v.Y+v.Height))
		}
	}
}

func (v *WeaponFrame) drawItemDisplay(screen *ebiten.Image) {
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("lvl. %d", v.level), int(v.X+4), int(drawY+8))
}

var _ ui.UIElement = (*WeaponFrame)(nil)
var _ WeaponFrameInterface = (*WeaponFrame)(nil)
