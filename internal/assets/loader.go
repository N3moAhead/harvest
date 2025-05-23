package assets

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Store struct {
	images map[string]*ebiten.Image
	sfx    map[string][]byte
	music  map[string][]byte
	fonts  map[string]font.Face
}

func NewStore() *Store {
	return &Store{
		images: make(map[string]*ebiten.Image),
		sfx:    make(map[string][]byte),
		music:  make(map[string][]byte),
		fonts:  make(map[string]font.Face),
	}
}

// --- GETTERS ---

func (s *Store) GetMusicData(name string) (music []byte, musicFound bool) {
	music, ok := s.music[name]
	return music, ok
}

func (s *Store) GetSFXData(name string) (sfx []byte, sfxFound bool) {
	sfx, ok := s.sfx[name]
	return sfx, ok
}

func (s *Store) GetImage(name string) (image *ebiten.Image, imageFound bool) {
	img, ok := s.images[name]
	return img, ok
}

func (s *Store) GetFont(name string) (fontFace font.Face, fontFound bool) {
	fontFace, ok := s.fonts[name]
	return fontFace, ok
}

func (s *Store) Load(
	imageFiles map[string]string,
	sfxFiles map[string]string,
	fontFiles map[string]string,
	musicFiles map[string]string,
	audioSampleRate int,
) error {
	var err error

	fmt.Println("Loading image files...")
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
		fmt.Printf("Image '%s' loaded: %s\n", name, path)
	}

	fmt.Println("Loading sfx files...")
	for name, path := range sfxFiles {
		if _, exists := s.sfx[name]; exists {
			continue
		}
		data, loadErr := loadAudioFile(path, audioSampleRate)
		if loadErr != nil {
			err = fmt.Errorf("Error while loading the sfx'%s' (%s): %w", name, path, loadErr)
			fmt.Println(err)
			continue
		}
		s.sfx[name] = data
		fmt.Printf("SFX '%s' loaded: %s\n", name, path)
	}

	fmt.Println("Loading font files...")
	for name, path := range fontFiles {
		if _, exists := s.fonts[name]; exists {
			continue
		}
		font, loadErr := loadFontFile(path, 24, 72, font.HintingFull)
		if loadErr != nil {
			err = fmt.Errorf("Error while loading the font'%s' (%s): %w", name, path, loadErr)
			fmt.Println(err)
			continue
		}
		s.fonts[name] = font
		fmt.Printf("Font '%s' loaded: %s\n", name, path)
	}

	// Loading music is very slow because currently its getting loaded
	// fully into memory. We should consider to stream parts of the music file only
	// when we need it to speed up loading times :)
	// I will just keep it for the moment causing to simplicity XP
	fmt.Println("Loading music files...")
	for name, path := range musicFiles {
		if _, exists := s.music[name]; exists {
			continue
		}
		data, loadErr := loadAudioFile(path, audioSampleRate)
		if loadErr != nil {
			err = fmt.Errorf("Error while loading the music'%s' (%s): %w", name, path, loadErr)
			fmt.Println(err)
			continue
		}
		s.music[name] = data
		fmt.Printf("Music '%s' loaded: %s\n", name, path)
	}

	return err
}

// A function to load audio files
func loadAudioFile(path string, sampleRate int) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Could not opn: %w", err)
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".wav":
		s, err := wav.DecodeWithSampleRate(sampleRate, f)
		if err != nil {
			return nil, fmt.Errorf("wav decoding gone wrong: %w", err)
		}
		// Read the whole decoded stream and save it to memory
		data, err := io.ReadAll(s)
		if err != nil {
			return nil, fmt.Errorf("Error while reading the decoded audio stream: %w", err)
		}
		return data, nil
	case ".mp3":
		s, err := mp3.DecodeWithSampleRate(sampleRate, f)
		if err != nil {
			return nil, fmt.Errorf("mp3 decoding gone wrong: %w", err)
		}
		// Read the whole decoded stream and save it to memory
		data, err := io.ReadAll(s)
		if err != nil {
			return nil, fmt.Errorf("Error while reading the decoded audio stream: %w", err)
		}
		return data, nil
	case ".ogg":
		s, err := vorbis.DecodeWithSampleRate(sampleRate, f)
		if err != nil {
			return nil, fmt.Errorf("ogg/vorbis decoding gone wrong: %w", err)
		}
		// Read the whole decoded stream and save it to memory
		data, err := io.ReadAll(s)
		if err != nil {
			return nil, fmt.Errorf("Error while reading the decoded audio stream: %w", err)
		}
		return data, nil
	default:
		return nil, fmt.Errorf("Unknown Audio Format: %s", ext)
	}
}

func loadFontFile(path string, size float64, dpi float64, hinting font.Hinting) (font.Face, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	myFontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: hinting,
	})
	if err != nil {
		return nil, err
	}
	return myFontFace, nil
}
