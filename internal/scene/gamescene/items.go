package gamescene

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/weapon"
	"github.com/hajimehoshi/ebiten/v2"
)

func updateItems(g *GameScene) {
	n := 0
	for i := range g.items {
		gItem := g.items[i]
		// The update function also puts collected items into the inventory
		itemPickedUp := gItem.Update(g.Player)
		if itemPickedUp {
			// Add picked up items into the inventory
			switch gItem.CategoryOf() {
			case itemtype.CategoryVegetable:
				g.inventory.AddVegtable(gItem.Type)
			case itemtype.CategorySoup:
				g.inventory.AddSoup(gItem.Type)
				soup := gItem.RetrieveItemInfo().Soup
				g.Player.ExtendOrAddSoup(soup)
			case itemtype.CategoryWeapon:
				switch gItem.Type {
				case itemtype.Spoon:
					newWeapon := weapon.NewSpoon()
					added := g.inventory.AddWeapon(newWeapon)
					if !added {
						fmt.Printf("Inventory is full or weapon '%s' already exists\n", newWeapon.Name())
					} else {
						fmt.Printf("Weapon '%s' added to Inventory\n", newWeapon.Name())
					}
				case itemtype.ThrowingKnifes:
					newWeapon := weapon.NewThrowingKnife()
					added := g.inventory.AddWeapon(newWeapon)
					if !added {
						fmt.Printf("Inventory is full or weapon '%s' already exists\n", newWeapon.Name())
					} else {
						fmt.Printf("Weapon '%s' added to Inventory\n", newWeapon.Name())
					}
				case itemtype.RollingPin:
					newWeapon := weapon.NewRollingPin()
					added := g.inventory.AddWeapon(newWeapon)
					if !added {
						fmt.Printf("Inventory is full or weapon '%s' already exists\n", newWeapon.Name())
					} else {
						fmt.Printf("Weapon '%s' added to Inventory\n", newWeapon.Name())
					}
				case itemtype.Thermalmixer:
					newWeapon := weapon.NewThermalmixer()
					added := g.inventory.AddWeapon(newWeapon)
					if !added {
						fmt.Printf("Inventory is full or weapon '%s' already exists\n", newWeapon.Name())
					} else {
						fmt.Printf("Weapon '%s' added to Inventory\n", newWeapon.Name())
					}
				default:
					fmt.Printf("Warning: Unknown weapon type: %s", gItem.DisplayName())
				}
			default:
				panic(fmt.Errorf("unhandeld item category: %s in items update", gItem.CategoryOf().String()))
			}
		} else {
			// Remove items after the player picked them up
			if n != i {
				g.items[n] = gItem
			}
			n++
		}
	}
	g.items = g.items[:n]
}

func drawItems(g *GameScene, screen *ebiten.Image, mapOffsetX, mapOffsetY float64) {
	for _, item := range g.items {
		item.Draw(screen, mapOffsetX, mapOffsetY)
	}
}

func initItems() []*item.Item {
	items := []*item.Item{
		item.NewSpoon(
			(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
			(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2-50,
		),
		item.NewThrowingKnifes(
			(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
			(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2+80,
		),
		item.NewThrowingKnifes(
			(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
			(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2+160,
		),
		item.NewThrowingKnifes(
			(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
			(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2+240,
		),
		item.NewRollingPin(
			(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
			(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2-80,
		),
		item.NewThermalmixer(
			(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
			(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2-150,
		),
	}

	return items
}
