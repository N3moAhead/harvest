package gamescene

import (
	"fmt"
	"math/rand"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/toast"
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
			g.Score += 1 // Picking up items increases the score
			// Add picked up items into the inventory
			switch gItem.CategoryOf() {
			case itemtype.CategoryVegetable:
				g.inventory.AddVegtable(gItem.Type)
			case itemtype.CategorySoup:
				g.Score += 10000
				g.inventory.AddSoup(gItem.Type)
				soup := gItem.RetrieveItemInfo().Soup
				toast.AddToast(fmt.Sprintf("%s collected! +10.000 Score", soup.Type.String()))
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
	worldWidth := config.WIDTH_IN_TILES * config.TILE_SIZE
	worldHeight := config.HEIGHT_IN_TILES * config.TILE_SIZE
	items := []*item.Item{
		item.NewSpoon(
			(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
			(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2-50,
		),
	}
	items = append(items, getWeaponsAtRandomPositions(worldWidth, worldHeight, item.NewSpoon, 2)...)
	items = append(items, getWeaponsAtRandomPositions(worldWidth, worldHeight, item.NewThrowingKnifes, 3)...)
	items = append(items, getWeaponsAtRandomPositions(worldWidth, worldHeight, item.NewRollingPin, 3)...)
	items = append(items, getWeaponsAtRandomPositions(worldWidth, worldHeight, item.NewThermalmixer, 3)...)

	return items
}

func getWeaponsAtRandomPositions(worldWidth, worldHeight int, create func(x, y float64) *item.Item, amount int) []*item.Item {
	var items []*item.Item = make([]*item.Item, 0)
	for range amount {
		x := rand.Float64() * float64(worldWidth)
		y := rand.Float64() * float64(worldHeight)
		items = append(items, create(x, y))
	}
	return items
}
