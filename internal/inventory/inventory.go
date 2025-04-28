package inventory

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/itemtype"
	"github.com/N3moAhead/harvest/internal/weapon"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Inventory struct {
	Vegtables  map[itemtype.ItemType]int  // Mapping the item type to the amount
	Soups      map[component.SoupType]int // same here, but for soups
	Weapons    []weapon.Weapon
	MaxWeapons int
}

func (i *Inventory) AddVegtable(itemType itemtype.ItemType) {
	i.Vegtables[itemType]++
}

func (inv *Inventory) AddWeapon(newWeapon weapon.Weapon) (didWork bool) {
	// If the weapon already exists we level up
	for _, existingWeapon := range inv.Weapons {
		if existingWeapon != nil && existingWeapon.Name() == newWeapon.Name() {
			// TODO level up the weapon... still a design decision
			return false
		}
	}

	for i := 0; i < inv.MaxWeapons; i++ {
		if inv.Weapons[i] == nil {
			inv.Weapons[i] = newWeapon
			// TODO remove debugging statement
			fmt.Printf("Added weapon '%s' to slot %d \n", newWeapon.Name(), i+1)
			break
		}
	}

	return true
}

func (i *Inventory) Draw(screen *ebiten.Image) {
	if amount, ok := i.Vegtables[itemtype.Potato]; ok {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Item: %s, Amount: %d\n\n", itemtype.Potato.String(), amount), 10, 20)
	}
	if amount, ok := i.Vegtables[itemtype.Carrot]; ok {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Item: %s, Amount: %d\n\n", itemtype.Carrot.String(), amount), 10, 35)
	}
}

func (i *Inventory) AddSoup(soupType component.SoupType) {
	i.Soups[soupType]++
}

func (i *Inventory) RemoveSoup(soupType component.SoupType) {
	i.Soups[soupType]--
	if i.Soups[soupType] <= 0 {
		delete(i.Soups, soupType)
	}
}

func NewInventory() *Inventory {
	return &Inventory{
		Vegtables:  make(map[itemtype.ItemType]int),
		Soups:      make(map[component.SoupType]int),
		Weapons:    make([]weapon.Weapon, config.MAX_WEAPONS),
		MaxWeapons: config.MAX_WEAPONS,
	}
}
