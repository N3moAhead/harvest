package assets

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Store struct {
	images map[string]*ebiten.Image
	sfx    map[string][]byte
	music  map[string][]byte
	// fonts  map[string]font.Face TODO maybe we will have to add fonts
}

func NewStore() *Store {
	return &Store{
		images: make(map[string]*ebiten.Image),
		sfx:    make(map[string][]byte),
		music:  make(map[string][]byte),
	}
}

func (s *Store) Load(imageFiles map[string]string) error {
	var err error

	for name, path := range imageFiles {
		if _, exists := s.images[name]; exists {
			// We will maybe load some graphics before the Load() will be called.
			// For example to use them on the loading screen or something so they will
			// be skipped here!
			continue
		}

		img, _, loadErr := ebitenutil.NewImageFromFile(path)
		if loadErr != nil {
			err = fmt.Errorf("failed to load image %s: %w", name, loadErr)
			continue
		}
		s.images[name] = img
		fmt.Printf("Image '%s' loaded\n", name)
	}

	return err
}

// // A function to load audio files
// func loadAudioFile(path string, sampleRate int) ([]byte, error) {

// }

// --- GETTERS ---

// func (s *Store) GetMusicData(name string) (music []byte, musicFound bool) {
// }

// func (s *Store) GetSFXData(name string) (sfx []byte, sfxFound bool) {

// }

func (s *Store) GetImage(name string) (image *ebiten.Image, imageFound bool) {
	img, ok := s.images[name]
	return img, ok
}
