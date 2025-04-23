package assets

import (
	"bytes"
	"fmt"

	"github.com/N3moAhead/harvest/pkg/config"
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
		"spoon_slash": "assets/images/weapons/spoon/spoon_slash.png",
	}
	sfxToLoad := map[string]string{
		"laser": "assets/audio/sfx/laserTest.wav",
	}
	musicToLoad := map[string]string{
		"menu": "assets/audio/music/8bitMenuMusic.mp3",
	}

	err := AssetStore.Load(imagesToLoad, sfxToLoad, musicToLoad, config.AUDIO_SAMPLE_RATE)
	if err != nil {
		panic(err)
	}

	// TODO REMOVE or change this section
	// This here should just be a test to test running music :)
	music, ok := AssetStore.GetMusicData("menu")
	if ok {
		musicBytesReader := bytes.NewReader(music)
		loop := audio.NewInfiniteLoop(musicBytesReader, int64(len(music)))

		MusicPlayer, err = AudioContext.NewPlayer(loop)
		if err == nil {
			MusicPlayer.Play()
		} else {
			err = fmt.Errorf("Musikplayer konnte nicht erstellt werden: %v\n", err)
			panic(err)
		}
	}
}
