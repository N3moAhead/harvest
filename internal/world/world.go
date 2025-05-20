package world

import (
	"math"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	tiles        [][]Tile
	tileWidth    int
	tileHeight   int
	mapWidthPx   int
	mapHeightPx  int
	cameraPos    component.Vector2D // Current camera position
	cameraTarget component.Vector2D // Target camera position
	cameraSpeed  float64            // How fast the camera is moving
}

func NewWorld(widthInTiles, heightInTiles int) *World {
	tileW, tileH := config.TILE_SIZE, config.TILE_SIZE

	tiles := NewMap()

	m := &World{
		tiles:        tiles,
		tileWidth:    tileW,
		tileHeight:   tileH,
		mapWidthPx:   widthInTiles * tileW,
		mapHeightPx:  heightInTiles * tileH,
		cameraPos:    component.Vector2D{X: 0, Y: 0},
		cameraTarget: component.Vector2D{X: 0, Y: 0},
		// The camera speed. The camera is faster as higher this value gets
		cameraSpeed: config.CAMERA_SPEED,
	}
	return m
}

func (m *World) Update(targetWorldPos component.Vector2D, screenWidth, screenHeight int, dt float64) {
	// Calculate the target position of the camera
	// The camera should keep the player in the center of the screen.
	// Therefore, the top left corner of the camera = playerPos - half screen size.
	m.cameraTarget.X = targetWorldPos.X - float64(screenWidth)/2
	m.cameraTarget.Y = targetWorldPos.Y - float64(screenHeight)/2

	// Clamp the camera position to the map boundaries
	m.cameraTarget.X = math.Max(0, math.Min(m.cameraTarget.X, float64(m.mapWidthPx-screenWidth)))
	m.cameraTarget.Y = math.Max(0, math.Min(m.cameraTarget.Y, float64(m.mapHeightPx-screenHeight)))

	// Smoothly move the camera towards the target position (Linear Interpolation - Lerp)
	// Calculate the vector from the current position to the target position
	diff := m.cameraTarget.Sub(m.cameraPos)
	// Move the camera a portion of the way towards the target, depending on speed and dt
	// A simple Lerp factor would also be possible: factor := 0.1; m.cameraPos = m.cameraPos.Add(diff.Mul(factor))
	// But using dt makes it frame-independent:
	moveStep := diff.Mul(m.cameraSpeed * dt)

	moveStepLen := moveStep.Len()
	// Avoid "overshooting" at high frame rate/speed - move at most the difference
	if moveStepLen > diff.Len() || moveStepLen < 0.1 {
		m.cameraPos = m.cameraTarget
	} else {
		m.cameraPos = m.cameraPos.Add(moveStep)
	}

	// For Safety clamp the camera position to the map boundaries again maybe the value came a bit off
	m.cameraPos.X = math.Max(0, math.Min(m.cameraPos.X, float64(m.mapWidthPx-screenWidth)))
	m.cameraPos.Y = math.Max(0, math.Min(m.cameraPos.Y, float64(m.mapHeightPx-screenHeight)))
}

func (m *World) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	camX := m.cameraPos.X
	camY := m.cameraPos.Y

	// Determine the tile indices to draw
	// The plus 1 is a small buffer to ensure we draw enough tiles
	startTileX := int(math.Floor(camX/float64(m.tileWidth))) - 1
	startTileY := int(math.Floor(camY/float64(m.tileHeight))) - 1
	endTileX := int(math.Ceil((camX+float64(screenWidth))/float64(m.tileWidth))) + 1
	endTileY := int(math.Ceil((camY+float64(screenHeight))/float64(m.tileHeight))) + 1

	mapTilesWidth := len(m.tiles[0])
	mapTilesHeight := len(m.tiles)
	startTileX = max(0, startTileX)
	startTileY = max(0, startTileY)
	endTileX = min(mapTilesWidth, endTileX)
	endTileY = min(mapTilesHeight, endTileY)

	// Drawing the visible tiles
	for y := startTileY; y < endTileY; y++ {
		for x := startTileX; x < endTileX; x++ {
			tile := m.tiles[y][x]

			worldX := float64(x * m.tileWidth)
			worldY := float64(y * m.tileHeight)

			screenX := float32(worldX - camX)
			screenY := float32(worldY - camY)

			tile.Draw(screen, float64(screenX), float64(screenY))

			// vector.DrawFilledRect(
			// 	screen,
			// 	screenX,
			// 	screenY,
			// 	float32(m.tileWidth),
			// 	float32(m.tileHeight),
			// 	tileColor,
			// 	false,
			// )
		}
	}
}

func (w *World) GetCameraPosition() (float64, float64) {
	return w.cameraPos.X, w.cameraPos.Y
}
