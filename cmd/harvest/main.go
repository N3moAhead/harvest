package main

import (
	"log"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/scene"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	ebiten.SetWindowSize(config.SCREEN_WIDTH, config.SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Harvest by Wurzelwerk")
	ebiten.SetTPS(60)
	ebiten.SetFullscreen(true)
}

func main() {
	newSceneManger := scene.NewSceneManager()
	if err := ebiten.RunGame(newSceneManger); err != nil {
		log.Fatal(err)
	}
}
