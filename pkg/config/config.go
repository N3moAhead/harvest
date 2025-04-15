package config

// Currently just fixed constant values. PLS do not overuse it.
// TODO: Implement functions to read the config from a toml or json file
const (
	/// --- Window Settings ---
	SCREEN_WIDTH  = 700
	SCREEN_HEIGHT = 700
	/// --- Camera Settings ---
	CAMERA_SPEED = 4.0
	/// --- World Settings ---
	TILE_SIZE       = 32 // TILE_SIZE is the size of a tile in pixels.
	HEIGHT_IN_TILES = 40 // The number of tiles in the X direction.
	WIDTH_IN_TILES  = 40 // The number of tiles in the Y direction.
	/// --- Player Settings ---
	INITIAL_PLAYER_SPEED           = 6.0 // Initial Player Speed
	INITIAL_PLAYER_MAGNET_RADIUS   = 150.0
	PLAYER_PICKUP_RADIUS           = 50.0 // The radius in which items will be picked up into the players inventory
	PLAYER_MAGNET_ATTRACTION_SPEED = 10.0 // Determines how fast items move towards the player
	/// --- Audio Settings ---
	AUDIO_SAMPLE_RATE = 44100
)
