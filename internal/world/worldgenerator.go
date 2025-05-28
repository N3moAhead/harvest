package world

import (
	"fmt"
	"math/rand"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	flowerAmount = int((config.HEIGHT_IN_TILES * config.WIDTH_IN_TILES) * 0.08)
)

func NewMap() [][]Tile {
	tiles := make([][]Tile, config.HEIGHT_IN_TILES)

	grassMiddle1 := loadTileImage("tf_grass_middle")
	grassMiddle2 := loadTileImage("tf_grass_middle2")

	// Create a splice of all grass images
	plants := []*ebiten.Image{}

	cols, rows := 7, 4
	for row := range rows {
		for col := range cols {
			if row == 3 && col >= 4 {
				continue
			}
			name := "td_decor_" + fmt.Sprint(row) + "_" + fmt.Sprint(col)
			plantImage := loadTileImage(name)
			plants = append(plants, plantImage)
		}
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
