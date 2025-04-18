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
	AUDIO_SAMPLE_RATE    = 44100
	// Enemy: Carrot
	CARROT_SPEED           = 50.0
	CARROT_HEALTH          = 3
	CARROT_DAMAGE          = 5
	CARROT_ATTACK_COOLDOWN = 1.0
	CARROT_ATTACK_RANGE    = 20.0
	CARROT_ATTACK_START    = 0.0
	// Carrot Style
	CARROT_WIDTH   = 16
	CARROT_HEIGHT  = 16
	CARROT_COLOR_R = 255
	CARROT_COLOR_G = 128
	CARROT_COLOR_B = 0
	CARROT_COLOR_A = 255
)
