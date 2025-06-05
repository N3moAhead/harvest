package config

import "time"

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
	/// --- HUD Settings ---
	VEGTABLE_TYPE_AMOUNT = 5  // The amount of diffrent vegtable types
	SOUP_TYPE_AMOUNT     = 2  // The amount of diffrent soup types
	ITEM_FRAME_SIZE      = 48 // The size in pixels of an item frame
	/// --- Icon Settings ---
	ICON_SIZE               = 16.0 // The size in pixels of icon assets
	ICON_ON_MAP_RENDER_SIZE = 16.0 // The size in pixels on how large an item icon should be rendered
	/// --- Enemy Settings ---
	DEFAULT_ENEMY_ASSET_SIZE = 32.0 // THe size in pixels of default enemies
	// Enemy: Carrot
	CARROT_SPEED                  = 50.0
	CARROT_HEALTH                 = 2
	CARROT_DAMAGE                 = 5
	CARROT_ATTACK_COOLDOWN        = 1.0
	CARROT_ATTACK_RANGE           = 20.0
	CARROT_ATTACK_START           = 0.0
	CARROT_DROP_PROB              = 0.8 // 80% chance to drop an item
	CARROT_DROP_AMOUNT            = 1   // Drops 1 item
	CARROT_DROP_AMOUNT_PER_MINUTE = 0.1 // Drops additional 0.1 items per minute
	// CARROT_DROP_AMOUNT_PER_MINUTE = 0.4 // Drops additional 0.1 items per minute
	// Carrot Style
	CARROT_COLOR_R = 255
	CARROT_COLOR_G = 128
	CARROT_COLOR_B = 0
	CARROT_COLOR_A = 255
	// Enemy: Cabbage
	CABBAGE_SPEED                  = 30
	CABBAGE_HEALTH                 = 4
	CABBAGE_DAMAGE                 = 1
	CABBAGE_ATTACK_COOLDOWN        = 1.5
	CABBAGE_ATTACK_RANGE           = 40.0
	CABBAGE_DROP_PROB              = 0.8
	CABBAGE_DROP_AMOUNT            = 1
	CABBAGE_DROP_AMOUNT_PER_MINUTE = 0.1
	// Enemy: Onion
	ONION_SPEED                  = 40
	ONION_HEALTH                 = 1
	ONION_DAMAGE                 = 3
	ONION_ATTACK_COOLDOWN        = 1.0
	ONION_ATTACK_RANGE           = 15.0
	ONION_DROP_PROB              = 0.8
	ONION_DROP_AMOUNT            = 1
	ONION_DROP_AMOUNT_PER_MINUTE = 0.1
	// Enemy: Leek
	LEEK_SPEED                  = 45
	LEEK_HEALTH                 = 3
	LEEK_DAMAGE                 = 7
	LEEK_ATTACK_COOLDOWN        = 2.0
	LEEK_ATTACK_RANGE           = 25.0
	LEEK_DROP_PROB              = 0.6
	LEEK_DROP_AMOUNT            = 1
	LEEK_DROP_AMOUNT_PER_MINUTE = 0.1

	/// --- Enemy Spawning Settings ---
	BASE_SPAWN_INTERVAL_SEC = 2.0   // Base interval (seconds) for spawning enemies
	BASE_COUNT_PER_BATCH    = 1     // Base count of enemies per batch
	MIX_START_SEC           = 120.0 // Time when mixing starts
	// Enemy: Potato
	POTATO_SPEED                  = 25.0
	POTATO_HEALTH                 = 6
	POTATO_DAMAGE                 = 10
	POTATO_ATTACK_COOLDOWN        = 3.0
	POTATO_ATTACK_RANGE           = 40.0
	POTATO_ATTACK_START           = 0.0
	POTATO_DROP_PROB              = 0.8 // 80% chance to drop an item
	POTATO_DROP_AMOUNT            = 1   // Drops 1 item
	POTATO_DROP_AMOUNT_PER_MINUTE = 0.1 // Drops additional 0.1 items per minute
	// Potato Style
	POTATO_WIDTH   = 16
	POTATO_HEIGHT  = 16
	POTATO_COLOR_R = 183
	POTATO_COLOR_G = 146
	POTATO_COLOR_B = 104
	POTATO_COLOR_A = 255
	///--- TOASTS ---
	DEFAULT_TOAST_DURATION = 2 * time.Second
	TOAST_GAP              = 10
)
