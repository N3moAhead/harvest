package inventory

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/weapon"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Inventory struct {
	Vegetables map[itemtype.ItemType]int // Mapping the item type to the amount
	Soups      map[itemtype.ItemType]int // same here, but for soups
	Weapons    []weapon.Weapon
	MaxWeapons int
}

func (i *Inventory) AddVegtable(itemType itemtype.ItemType) {
	i.Vegetables[itemType]++
}

func (i *Inventory) RemoveVegetable(itemType itemtype.ItemType) {
	if i.Vegetables[itemType] > 0 {
		i.Vegetables[itemType]--
		if i.Vegetables[itemType] <= 0 {
			delete(i.Vegetables, itemType)
		}
	}
}

func (i *Inventory) RemoveNVegetables(itemType itemtype.ItemType, amount int) {
	if i.Vegetables[itemType] >= amount {
		i.Vegetables[itemType] -= amount
		if i.Vegetables[itemType] <= 0 {
			delete(i.Vegetables, itemType)
		}
	} else {
		delete(i.Vegetables, itemType)
		fmt.Printf("Not enough %s in inventory to remove %d items\n", itemType.String(), amount)
	}
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
	if amount, ok := i.Vegetables[itemtype.Potato]; ok {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Item: %s, Amount: %d\n\n", itemtype.Potato.String(), amount), 10, 20)
	}
	if amount, ok := i.Vegetables[itemtype.Carrot]; ok {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Item: %s, Amount: %d\n\n", itemtype.Carrot.String(), amount), 10, 35)
	}
	if amount, ok := i.Soups[itemtype.MagnetRadiusSoup]; ok {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Buff: %s, Amount: %d\n\n", itemtype.MagnetRadiusSoup.String(), amount), 10, 50)
	}
	if amount, ok := i.Soups[itemtype.SpeedSoup]; ok {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Buff: %s, Amount: %d\n\n", itemtype.SpeedSoup.String(), amount), 10, 65)
	}
}

func (i *Inventory) AddSoup(soupType itemtype.ItemType) {
	i.Soups[soupType]++
}

func (i *Inventory) RemoveSoup(soupType itemtype.ItemType) {
	i.Soups[soupType]--
	if i.Soups[soupType] <= 0 {
		delete(i.Soups, soupType)
	}
}

func (i *Inventory) RemoveAllSoups(soupType itemtype.ItemType) {
	delete(i.Soups, soupType)
}

func NewInventory() *Inventory {
	return &Inventory{
		Vegetables: make(map[itemtype.ItemType]int),
		Soups:      make(map[itemtype.ItemType]int),
		Weapons:    make([]weapon.Weapon, config.MAX_WEAPONS),
		MaxWeapons: config.MAX_WEAPONS,
	}
}
