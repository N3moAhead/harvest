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
	// If set to true the game will end in the next update loop
	exitGame bool
	stats    PlayerStats
}

type PlayerStats struct {
	highScore        int
	lastGameScore    int // The score of the last played game
	lastGameXPEarned uint
	playerXP         uint // 10.000 Score Points = 1 XP
	playerLevel      uint // 10 XP => 1 Player Level
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
		fmt.Println("Warning: Could not load currentScene in scene manager. Falling back to menu!")
		s.setNextScene(MENU_SCENE)
		return s.getCurrentScene()
	}
}

func (s *SceneManager) setNextScene(scene SceneId) {
	// Its not possible to set the loading scene because
	// its only used at startup
	switch scene {
	case MENU_SCENE:
		fmt.Println("Switched Scene to Menu")
		s.menuScene = NewMenuScene(s.setExitGame, s.stats)
		s.currentScene = MENU_SCENE
	case GAME_SCENE:
		fmt.Println("Switched Scene to Game Scene")
		s.gameScene = gamescene.NewGameScene(func() { s.updateHighScore(); s.setNextScene(MENU_SCENE) }, s.stats.playerLevel)
		s.currentScene = GAME_SCENE
	case SCORE_SCENE:
		fmt.Println("Switched Scene to Score")
		s.scoreScene = NewScoreScene(s.stats)
		s.currentScene = SCORE_SCENE
	default:
		fmt.Println("Warning: Switched to menu, received undefined SceneId: ", scene)
		s.setNextScene(MENU_SCENE)
	}
}

// A scene ends and this functions returns a logical
// follow up scene in the following direction
// LoadingScene -> MenuScene -> GameScene -> ScoreScene -> MenuScene
func (s *SceneManager) determineFollowUpScene() SceneId {
	switch s.currentScene {
	case LOADING_SCENE:
		return MENU_SCENE
	case MENU_SCENE:
		return GAME_SCENE
	case GAME_SCENE:
		return SCORE_SCENE
	case SCORE_SCENE:
		return MENU_SCENE
	default:
		fmt.Println("Warning: Follow up -> MenuScene Could not determine follow up scene the scene:", s.currentScene)
		return MENU_SCENE
	}
}

func (s *SceneManager) updateHighScore() {
	if s.currentScene == GAME_SCENE {
		scene := s.getCurrentScene()
		if scoreScene, ok := scene.(gamescene.Score); ok {
			newScore := scoreScene.GetScore()
			s.stats.lastGameXPEarned = uint(newScore / 10000)
			s.stats.playerXP += s.stats.lastGameXPEarned
			s.stats.playerLevel = uint(s.stats.playerXP / 10)
			s.stats.lastGameScore = newScore
			if s.stats.highScore < newScore {
				s.stats.highScore = newScore
			}
		}
	}
}

func (s *SceneManager) Update() error {
	// --- Check for Exit ---
	if s.exitGame {
		return errors.New("Quitted Game")
	}

	scene := s.getCurrentScene()
	if !scene.IsRunning() {
		// If a game scene just ended this function will update the highscore
		s.updateHighScore()
		followUpScene := s.determineFollowUpScene()
		s.setNextScene(followUpScene)
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

func (s *SceneManager) setExitGame() {
	s.exitGame = true
}

var _ ebiten.Game = (*SceneManager)(nil)
