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

type FontConfig struct {
	Path string
	Size int
}

func LoadAllAssets() {
	// Always image name to path
	imagesToLoad := map[string]string{
		"player":             "assets/images/CookTestImage.png",
		"spoon_slash":        "assets/images/weapons/spoon/spoon_slash3.png",
		"rolling_pin_roll":   "assets/images/weapons/rolling_pin/rolling_pin_roll.png",
		"carrot":             "assets/images/carrot.png",
		"cook_station":       "assets/images/cookstation.png",
		"potato":             "assets/images/potato.png",
		"cabbage":            "assets/images/cabbage.png",
		"cabbage_icon":       "assets/images/icons/cabbage_icon.png",
		"onion":              "assets/images/onion.png",
		"onion_icon":         "assets/images/icons/onion_icon.png",
		"menu-icon":          "assets/images/menu_icon.png",
		"leek":               "assets/images/leek.png",
		"leek_icon":          "assets/images/icons/leek_icon.png",
		"radish":             "assets/images/radish.png",
		"radish_icon":        "assets/images/icons/radish_icon.png",
		"thermalmixer_slash": "assets/images/weapons/thermalmixer/thermalmixer_slash.png",
		// Map Tiles: (t stands for tile; f stands for floor; d stands for decor)
		"tf_grass_middle":      "assets/images/world/grass.png",
		"tf_grass_middle2":     "assets/images/world/grass2.png",
		"outdoor_decor_sprite": "assets/images/world/decorations.png",
		// Icons
		"carrot_icon":       "assets/images/icons/carrot_icon.png",
		"potato_icon":       "assets/images/icons/potato_icon.png",
		"spoon_icon":        "assets/images/icons/spoon_icon.png",
		"thermalmixer_icon": "assets/images/icons/thermalmixer_icon.png",
		"no_icon":           "assets/images/icons/no_icon.png",
		"rolling_pin_icon":  "assets/images/icons/rolling_pin_icon.png",
		// Hud
		"vegtable_item_frame": "assets/images/hud/hud_item_frame.png",
		"soup_item_frame":     "assets/images/hud/hud_item_frame2.png",
		"weapon_item_frame":   "assets/images/hud/hud_item_frame3.png",
		// Weapons
		"throwing_knifes_icon": "assets/images/icons/throwing_knifes_icon.png",
		"knife_projectile":     "assets/images/weapons/throwing_knifes/knife_projectile.png",
	}
	sfxToLoad := map[string]string{
		"laser":              "assets/audio/sfx/laserTest.wav",
		"spoon_slash":        "assets/audio/sfx/spoon_slash.mp3",
		"rolling_pin_roll":   "assets/audio/sfx/rolling_pin_roll.mp3",
		"thermalmixer_slash": "assets/audio/sfx/thermalmixer_slash.mp3",
		// TODO Add correct thermalmixer sound
		"game_loads_sound":   "assets/audio/sfx/game_loads_sound.wav",
		"player_hit_sound":   "assets/audio/sfx/player_hit_sound.wav",
		"player_death_sound": "assets/audio/sfx/player_death_sound.mp3",
		"knife_throw":        "assets/audio/sfx/knife_throw.wav",
		"knife_throw_impact": "assets/audio/sfx/knife_throw_impact.wav",
	}
	musicToLoad := map[string]string{
		"menu": "assets/audio/music/8bitMenuMusic.mp3",
	}

	fontsToLoad := map[string]FontConfig{
		"2p":    FontConfig{Path: "assets/fonts/PressStart2P-Regular.ttf", Size: 24},
		"micro": FontConfig{Path: "assets/fonts/micro.ttf", Size: 20},
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
	initFontsToLoad := map[string]FontConfig{
		"2p": FontConfig{Path: "assets/fonts/PressStart2P-Regular.ttf", Size: 24},
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
