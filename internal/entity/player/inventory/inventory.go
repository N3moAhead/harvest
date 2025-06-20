package inventory

import (
	"fmt"
	"sort"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/toast"
	"github.com/N3moAhead/harvest/internal/weapon"
	"github.com/hajimehoshi/ebiten/v2"
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
			if ok := existingWeapon.LevelUp(); ok {
				toast.AddToast(fmt.Sprintf("'%s' updated to level %d", existingWeapon.Name(), existingWeapon.Level()))
				fmt.Printf("Weapon '%s' updated to level %d", existingWeapon.Name(), existingWeapon.Level())
			} else {
				fmt.Printf("Weapon '%s' is already at max level %d", existingWeapon.Name(), existingWeapon.MaxLevel())
			}
			return false
		}
	}

	for i := 0; i < inv.MaxWeapons; i++ {
		if inv.Weapons[i] == nil {
			inv.Weapons[i] = newWeapon
			// TODO remove debugging statement
			fmt.Printf("Added weapon '%s' to slot %d \n", newWeapon.Name(), i+1)
			toast.AddToast(fmt.Sprintf("%s collected!", newWeapon.Name()))
			break
		}
	}

	return true
}

func (i *Inventory) Draw(screen *ebiten.Image) {
	// Add a soup display similiar to the vegtable display
	offset := 0.0
	soupTypes := make([]itemtype.ItemType, 0)
	for k, _ := range i.Soups {
		soupTypes = append(soupTypes, k)
	}
	sort.Sort(itemtype.ByItemType(soupTypes))
	for _, soupKey := range soupTypes {
		amount := i.Soups[soupKey]
		drawItemDisplay(screen, soupKey, amount, offset)
		offset += 18
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
