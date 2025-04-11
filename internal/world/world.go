package world

import (
	"image/color"
	"math"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TileType int

const (
	TileGround TileType = iota
	TileWater
	TileWall
)

type World struct {
	tiles        [][]TileType
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

	tiles := make([][]TileType, heightInTiles)

	// Currently its just a small brown map with a small lake
	for y := range tiles {
		tiles[y] = make([]TileType, widthInTiles)
		for x := range tiles[y] {
			if x == 0 || y == 0 || x == widthInTiles-1 || y == heightInTiles-1 {
				tiles[y][x] = TileWall
			} else if x > 10 && x < 15 && y > 5 && y < 10 {
				tiles[y][x] = TileWater
			} else {
				tiles[y][x] = TileGround
			}
		}
	}

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

	// Avoid "overshooting" at high frame rate/speed - move at most the difference
	if moveStep.Len() > diff.Len() {
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
			tileType := m.tiles[y][x]

			worldX := float64(x * m.tileWidth)
			worldY := float64(y * m.tileHeight)

			screenX := float32(worldX - camX)
			screenY := float32(worldY - camY)

			var tileColor color.RGBA
			switch tileType {
			case TileGround:
				tileColor = color.RGBA{R: 139, G: 69, B: 19, A: 255} // Brown
			case TileWater:
				tileColor = color.RGBA{R: 0, G: 0, B: 200, A: 255} // Blue
			case TileWall:
				tileColor = color.RGBA{R: 100, G: 100, B: 100, A: 255} // Gray
			default:
				tileColor = color.RGBA{R: 255, G: 0, B: 255, A: 255} // Magenta
			}

			vector.DrawFilledRect(
				screen,
				screenX,
				screenY,
				float32(m.tileWidth),
				float32(m.tileHeight),
				tileColor,
				false,
			)
		}
	}
}

func (w *World) GetCameraPosition() (float64, float64) {
	return w.cameraPos.X, w.cameraPos.Y
}
