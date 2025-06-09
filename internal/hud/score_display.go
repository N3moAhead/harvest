package hud

import (
	"fmt"
	"image/color"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type ScoreDisplay struct {
	ui.Label
	ScorePointer *int
	ScoreText    string
}

func NewScoreDisplay(scorePointer *int, scoreText string) *ScoreDisplay {
	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		fmt.Println("Warning: Could not load NewFont in NewScoreDisplay")
	}

	newScoreDisplay := &ScoreDisplay{
		Label:        *ui.NewLabel(0, 0, fmt.Sprintf("%s: XXXXXX", scoreText), fontFace, color.White),
		ScorePointer: scorePointer,
		ScoreText:    scoreText,
	}

	drawX := config.SCREEN_WIDTH - newScoreDisplay.Width
	drawY := 20.0
	newScoreDisplay.X = drawX
	newScoreDisplay.Y = drawY

	return newScoreDisplay
}

func (sd *ScoreDisplay) Draw(screen *ebiten.Image) {
	if !sd.Visible {
		return
	}

	if sd.Font != nil {
		drawText := fmt.Sprintf("%s: %d", sd.ScoreText, *sd.ScorePointer)
		bounds := text.BoundString(sd.Font, drawText)
		textY := sd.Y + float64(sd.Font.Metrics().Ascent/64)
		textX := config.SCREEN_WIDTH - (bounds.Dx() + 10)
		text.Draw(screen, drawText, sd.Font, textX, int(textY), sd.Color)
	}
	sd.BaseElement.Draw(screen)
}

var _ ui.UIElement = (*ScoreDisplay)(nil)
