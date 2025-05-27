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

func LoadAllAssets() {
	// Always image name to path
	imagesToLoad := map[string]string{
		"player":           "assets/images/CookTestImage.png",
		"spoon_slash":      "assets/images/weapons/spoon/spoon_slash3.png",
		"rolling_pin_roll": "assets/images/weapons/rolling_pin/rolling_pin_roll.png",
		"carrot":           "assets/images/carrot.png",
		"cook_station":     "assets/images/cookstation.png",
		"potato":           "assets/images/potato.png",
		"menu-icon":        "assets/images/menu_icon.png",
		// Map Tiles: (t stands for tile; f stands for floor; d stands for decor)
		"tf_grass_middle":      "assets/images/world/grass.png",
		"tf_grass_middle2":     "assets/images/world/grass2.png",
		"outdoor_decor_sprite": "assets/images/world/decorations.png",
		// Icons
		"carrot_icon": "assets/images/icons/carrot_icon.png",
		"potato_icon": "assets/images/icons/potato_icon.png",
		"spoon_icon":  "assets/images/icons/spoon_icon.png",
		"no_icon":     "assets/images/icons/no_icon.png",
		// Hud
		"vegtable_item_frame": "assets/images/hud/hud_item_frame.png",
		"soup_item_frame":     "assets/images/hud/hud_item_frame2.png",
		"weapon_item_frame":   "assets/images/hud/hud_item_frame3.png",
	}
	sfxToLoad := map[string]string{
		"laser":            "assets/audio/sfx/laserTest.wav",
		"spoon_slash":      "assets/audio/sfx/spoon_slash.mp3",
		"rolling_pin_roll": "assets/audio/sfx/rolling_pin_roll.mp3",
		"game_loads_sound": "assets/audio/sfx/game_loads_sound.wav",
	}
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

	addDecorsFromSpriteSheet()
}

func init() {
	// A new Audio Context
	AudioContext = audio.NewContext(config.AUDIO_SAMPLE_RATE)
	// Initing the asset store
	AssetStore = NewStore()

	// On init just load the needed stuff for the loading screen afterwards
	// Everything else can be loaded
	initImagesToLoad := map[string]string{}
	initSFXToLoad := map[string]string{
		"game_loads_sound": "assets/audio/sfx/game_loads_sound.wav",
	}
	initMusicToLoad := map[string]string{}
	initFontsToLoad := map[string]string{
		"2p": "assets/fonts/PressStart2P-Regular.ttf",
	}
	err := AssetStore.Load(initImagesToLoad, initSFXToLoad, initFontsToLoad, initMusicToLoad, config.AUDIO_SAMPLE_RATE)
	if err != nil {
		panic(err)
	}

}

func addDecorsFromSpriteSheet() {
	decorSprite, ok := AssetStore.GetImage("outdoor_decor_sprite")
	cols, rows := 7, 4
	ts := config.TILE_SIZE
	if ok {
		for row := range rows {
			for col := range cols {
				name := "td_decor_" + fmt.Sprint(row) + "_" + fmt.Sprint(col)
				addTileFromSprite(name, decorSprite, col*ts, row*ts)
			}
		}
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
