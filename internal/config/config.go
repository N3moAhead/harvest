package config

// Currently just fixed constant values. PLS do not overuse it.
// TODO: Implement functions to read the config from a toml or json file
const (
	/// --- Window Settings ---
	SCREEN_WIDTH  = 896
	SCREEN_HEIGHT = 504
	/// --- Camera Settings ---
	CAMERA_SPEED = 6.0
	/// --- World Settings ---
	TILE_SIZE       = 16  // TILE_SIZE is the size of a tile in pixels.
	HEIGHT_IN_TILES = 200 // The number of tiles in the X direction.
	WIDTH_IN_TILES  = 200 // The number of tiles in the Y direction.
	/// --- Player Settings ---
	INITIAL_PLAYER_SPEED           = 6.0 // Initial Player Speed
	INITIAL_PLAYER_MAGNET_RADIUS   = 50.0
	PLAYER_PICKUP_RADIUS           = 5.0   // The radius in which items will be picked up into the players inventory
	PLAYER_MAGNET_ATTRACTION_SPEED = 7.0   // Determines how fast items move towards the player
	PLAYER_INTERACT_RADIUS         = 20.0  // The radius in which the player can interact with cookstations, NPCs, etc.
	SHOW_RECIPE_RANGE              = 200.0 // The range in which the player can see the recipe of a cookstation
	/// --- Audio Settings ---
	AUDIO_SAMPLE_RATE = 44100
	/// --- Inventory Settings ---
	MAX_WEAPONS = 5
	/// --- Enemy Settings ---
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

	/// --- Enemy Spawning Settings ---
	BASE_SPAWN_INTERVAL_SEC = 2.0   // Base interval (seconds) for spawning enemies
	BASE_COUNT_PER_BATCH    = 1     // Base count of enemies per batch
	MIX_START_SEC           = 120.0 // Time when mixing starts
)
