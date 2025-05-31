package scene

import (
	"errors"
	"fmt"

	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/scene/gamescene"
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneId string

const (
	LOADING_SCENE SceneId = "game_loading_scene"
	MENU_SCENE    SceneId = "menu_scene"
	GAME_SCENE    SceneId = "game_scene"
	SCORE_SCENE   SceneId = "game_over_scene"
)

type SceneManager struct {
	currentScene SceneId
	loadingScene Scene
	menuScene    Scene
	gameScene    Scene
	scoreScene   Scene
}

func NewSceneManager() *SceneManager {
	return &SceneManager{
		currentScene: LOADING_SCENE,
		loadingScene: NewLoadingScene(),
	}
}

func (s *SceneManager) getCurrentScene() Scene {
	switch s.currentScene {
	case LOADING_SCENE:
		return s.loadingScene
	case MENU_SCENE:
		return s.menuScene
	case GAME_SCENE:
		return s.gameScene
	case SCORE_SCENE:
		return s.scoreScene
	default:
		fmt.Println(s.currentScene)
		panic("Could not load currentScene in scene manger.")
	}
}

func (s *SceneManager) setNextScene() {
	// Loading -> Menu -> Game -> Score -> Menu
	// Every time we switch screens we have to
	// create a new instance of the scene to get a new cleared one
	switch s.currentScene {
	case LOADING_SCENE:
		fmt.Println("Switched Scene to Menu")
		s.menuScene = NewMenuScene()
		s.currentScene = MENU_SCENE
	case MENU_SCENE:
		fmt.Println("Switched Scene to Game")
		s.gameScene = gamescene.NewGameScene()
		s.currentScene = GAME_SCENE
	case GAME_SCENE:
		fmt.Println("Switched Scene to Score")
		s.scoreScene = NewScoreScene()
		s.currentScene = SCORE_SCENE
	case SCORE_SCENE:
		fmt.Println("Switched Scene to Menu")
		s.menuScene = NewMenuScene()
		s.currentScene = MENU_SCENE
	}
}

func (s *SceneManager) Update() error {
	// --- Check for Exit ---
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Game Quit!")
	}
	scene := s.getCurrentScene()
	if !scene.IsRunning() {
		s.setNextScene()
	}
	err := scene.Update()
	return err
}

func (s *SceneManager) Draw(screen *ebiten.Image) {
	scene := s.getCurrentScene()
	scene.Draw(screen)
}

func (s *SceneManager) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.SCREEN_WIDTH, config.SCREEN_HEIGHT
}

var _ ebiten.Game = (*SceneManager)(nil)
