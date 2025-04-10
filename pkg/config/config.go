package config

// Currently just fixed constant values. PLS do not overuse it.
// TODO: Implement functions to read the config from a toml or json file
const (
	SCREEN_WIDTH  = 700
	SCREEN_HEIGHT = 700
	CAMERA_SPEED  = 4.0
	// TILE_SIZE is the size of a tile in pixels.
	TILE_SIZE = 32
	// The number of tiles in the X direction.
	HEIGHT_IN_TILES = 40
	// The number of tiles in the Y direction.
	WIDTH_IN_TILES = 40
	// Initial Player Speed
	INITIAL_PLAYER_SPEED = 6.0
)
