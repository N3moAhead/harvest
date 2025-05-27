package world

import (
	"fmt"
	"math/rand"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	flowerAmount = int((config.HEIGHT_IN_TILES * config.WIDTH_IN_TILES) * 0.2)
)

func NewMap() [][]Tile {
	tiles := make([][]Tile, config.HEIGHT_IN_TILES)

	grassMiddle1 := loadTileImage("tf_grass_middle")
	grassMiddle2 := loadTileImage("tf_grass_middle2")
	grass1 := loadTileImage("td_grass1")
	grass2 := loadTileImage("td_grass2")
	grass3 := loadTileImage("td_grass3")
	flowerGrass1 := loadTileImage("td_flowerGrass1")
	flowerGrass2 := loadTileImage("td_flowerGrass2")
	flowerGrass3 := loadTileImage("td_flowerGrass3")
	mushroom := loadTileImage("td_mushroom")
	redFlower1 := loadTileImage("td_red_flower1")
	redFlower2 := loadTileImage("td_red_flower2")
	redFlower3 := loadTileImage("td_red_flower3")
	redFlower4 := loadTileImage("td_red_flower4")
	yellowFlower1 := loadTileImage("td_yellow_flower1")
	yellowFlower2 := loadTileImage("td_yellow_flower2")
	yellowFlower3 := loadTileImage("td_yellow_flower3")
	yellowFlower4 := loadTileImage("td_yellow_flower4")

	// Create a splice of all grass images
	plants := []*ebiten.Image{
		grass1,
		grass2,
		grass3,
		flowerGrass1,
		flowerGrass2,
		flowerGrass3,
		mushroom,
		redFlower1,
		redFlower2,
		redFlower3,
		redFlower4,
		yellowFlower1,
		yellowFlower2,
		yellowFlower3,
		yellowFlower4,
	}

	for y := range tiles {
		tiles[y] = make([]Tile, config.WIDTH_IN_TILES)
		for x := range tiles[y] {
			var floorImage *ebiten.Image = grassMiddle1
			if (y+x)%2 == 0 {
				floorImage = grassMiddle2
			}
			tiles[y][x] = Tile{
				Type:       GrassMiddle,
				FloorImage: floorImage,
				IsWalkable: true,
			}
		}
	}

	for range flowerAmount {
		tiles[rand.Intn(config.HEIGHT_IN_TILES)][rand.Intn(config.WIDTH_IN_TILES)].DecorImage = plants[rand.Intn(len(plants))]
	}

	return tiles
}

func loadTileImage(name string) *ebiten.Image {
	image, ok := assets.AssetStore.GetImage(name)
	if !ok {
		fmt.Printf("Warning could not load %s\n", name)
	}
	return image
}
