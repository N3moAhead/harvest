package assets

import (
	"fmt"
	"image"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/pkg/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

var (
	AssetStore   *Store
	AudioContext *audio.Context
	MusicPlayer  *audio.Player
)

func init() {
	// A new Audio Context
	AudioContext = audio.NewContext(config.AUDIO_SAMPLE_RATE)
	// Initing the asset store
	AssetStore = NewStore()

	// Always image name to path
	imagesToLoad := map[string]string{
		"player":      "assets/images/CookTestImage.png",
		"spoon_slash": "assets/images/weapons/spoon/spoon_slash3.png",
		"carrot":      "assets/images/carrot.png",
		// Map Tiles: (t stands for tile; f stands for floor; d stands for decor)
		"tf_grass_middle":      "assets/images/world/Grass_Middle.png",
		"outdoor_decor_sprite": "assets/images/world/outdoor_decor.png",
	}
	sfxToLoad := map[string]string{
		"laser":       "assets/audio/sfx/laserTest.wav",
		"spoon_slash": "assets/audio/sfx/spoon_slash.mp3",
	}
	// TODO Renable music
	musicToLoad := map[string]string{
		"menu": "assets/audio/music/8bitMenuMusic.mp3",
	}

	fontsToLoad := map[string]string{
		"2p": "assets/fonts/PressStart2P-Regular.ttf",
	}

	err := AssetStore.Load(imagesToLoad, sfxToLoad, fontsToLoad, musicToLoad, config.AUDIO_SAMPLE_RATE)
	if err != nil {
		panic(err)
	}

	addDecorsFromSpritesheet()

	// TODO Renable music
	// TODO REMOVE or change this section
	// This here should just be a test to test running music :)
	// music, ok := AssetStore.GetMusicData("menu")
	// if ok {
	// 	musicBytesReader := bytes.NewReader(music)
	// 	loop := audio.NewInfiniteLoop(musicBytesReader, int64(len(music)))

	// 	MusicPlayer, err = AudioContext.NewPlayer(loop)
	// 	if err == nil {
	// 		MusicPlayer.Play()
	// 	} else {
	// 		err = fmt.Errorf("Musikplayer konnte nicht erstellt werden: %v\n", err)
	// 		panic(err)
	// 	}
	// }
}

func addDecorsFromSpritesheet() {
	decorSprite, ok := AssetStore.GetImage("outdoor_decor_sprite")
	if ok {
		ts := config.TILE_SIZE
		addTileFromSprite("td_grass1", decorSprite, 0, 0)
		addTileFromSprite("td_grass2", decorSprite, 1*ts, 0)
		addTileFromSprite("td_grass3", decorSprite, 2*ts, 0)
		addTileFromSprite("td_flowerGrass1", decorSprite, 0, 1*ts)
		addTileFromSprite("td_flowerGrass2", decorSprite, 1*ts, 1*ts)
		addTileFromSprite("td_flowerGrass3", decorSprite, 2*ts, 1*ts)
		addTileFromSprite("td_mushroom", decorSprite, 2*ts, 7*ts)
		addTileFromSprite("td_red_flower1", decorSprite, 0, 8*ts)
		addTileFromSprite("td_red_flower2", decorSprite, 0, 9*ts)
		addTileFromSprite("td_red_flower3", decorSprite, 0, 10*ts)
		addTileFromSprite("td_red_flower4", decorSprite, 0, 11*ts)
		addTileFromSprite("td_yellow_flower1", decorSprite, 1*ts, 8*ts)
		addTileFromSprite("td_yellow_flower2", decorSprite, 1*ts, 9*ts)
		addTileFromSprite("td_yellow_flower3", decorSprite, 1*ts, 10*ts)
		addTileFromSprite("td_yellow_flower4", decorSprite, 1*ts, 11*ts)
	}
}

// Adds a tile from a spritesheet the AssetStore
// sx, sy are the position of the top left corner of the sub image on the spritesheet
// for the width and height config.TILE_SIZE will be used
func addTileFromSprite(name string, source *ebiten.Image, sx int, sy int) {
	fmt.Printf("Name: %s, sx %d, sy %d\n", name, sx, sy)
	subImg := util.GetSubImage(source, image.Rect(sx, sy, sx+config.TILE_SIZE, sy+config.TILE_SIZE))
	AssetStore.addImageToStore(name, subImg)
}
