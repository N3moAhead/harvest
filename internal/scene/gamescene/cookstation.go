package gamescene

import (
	"time"

	"github.com/N3moAhead/harvest/internal/cooking"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/pkg/util"
)

func updateCookStations(g *GameScene, dt float64, inv *inventory.Inventory, elapsed float32) {
	elapsedSec := float64(time.Since(g.startTime).Seconds())
	// start at 15 seconds, decrease to 5 seconds after 30 minutes /1800 seconds
	rawInterval := 15.0 - (elapsedSec / 180.0)
	interval := util.Clamp(rawInterval, 3.0, 15.0)

	if time.Since(g.lastCookStationSpawnTime).Seconds() >= interval {
		g.lastCookStationSpawnTime = time.Now()
		spawnCookBatch(g, 1) // Spawn a single cook station every interval
	}
	for _, cs := range g.cookStations {
		scoreAddition := cs.Update(g.Player, g.inventory)
		g.Score += scoreAddition
	}
}
func spawnCookBatch(g *GameScene, count int) {
	for i := 0; i < count; i++ {
		cameraX, cameraY := g.World.GetCameraPosition()
		spawnX, spawnY := util.GetRandomPositionInView(cameraX, cameraY) // Get a random position in the view
		recipe := cooking.GetRandomRecipe()
		station := cooking.NewCookStation(spawnX, spawnY, recipe, 1.0) // TODO 1.0 is the cost factor, can be adjusted by elapsed time
		g.cookStations = append(g.cookStations, station)
	}
}
