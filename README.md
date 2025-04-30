```ascii
  _    _                           _
 | |  | |                         | |
 | |__| | __ _ _ ____   _____  ___| |_
 |  __  |/ _` | '__\ \ / / _ \/ __| __|
 | |  | | (_| | |   \ V /  __/\__ \ |_
 |_|  |_|\__,_|_|    \_/ \___||___/\__|
```

## Project Overview

Harvest is a 2D rogue-like survival game developed in Go using the [Ebitengine](https://ebitengine.org/) library. The player controls a character defending against waves of opponents. This project was initially developed as part of a university software engineering course.

## Status

This project is currently under active development. Core mechanics are implemented, but further features, balancing, and polish are planned.

## Installation & Setup

1.  **Clone the repository:**
    ```bash
    # SSH
    git clone git@github.com:N3moAhead/harvest.git
    # HTTPS
    git clone https://github.com/N3moAhead/harvest.git

    cd harvest
    ```

2.  **Download Dependencies:**
    ```bash
    go mod tidy
    ```

## Running the Game

You can run the game using the provided Makefile (primarily for Linux/macOS) or standard Go commands.

### Using Make (Linux/macOS)

The `Makefile` provides convenient commands:

*   `make start`: Formats code (`go fmt`), checks for potential issues (`go vet`), builds the executable (`./harvest`), and runs it.
*   `make dev`: Uses `air` for live reloading during development. The application automatically rebuilds and restarts when code changes are detected. (Requires `air` to be installed).
*   `make clean`: Removes the compiled executable (`./harvest`).
*   `make build`: Just builds the executable file (`./harvest`).

Example:
```bash
make start
```

### Manual Compilation & Execution

If you don't have `make` installed or prefer standard Go commands:

**On Linux / macOS:**

1.  **Build:**
    ```bash
    go build -o harvest ./cmd/harvest/main.go
    ```
2.  **Run:**
    ```bash
    ./harvest
    ```

**On Windows:**

1.  **Build:**
    ```bash
    go build -o harvest.exe ./cmd/harvest/main.go
    ```
    *(Note: Windows uses backslashes `\` in paths in the console, but Go commands often understand forward slashes `/` as well)*
2.  **Run:**
    ```bash
    .\harvest.exe
    ```

## Project Structure

The project structure follows common Go conventions and is organized for clarity within an Ebitengine context:

```
.
├── assets          # Static game assets (images, audio, fonts)
│   ├── audio       # Sound files (music, sound effects)
│   ├── fonts       # Font files
│   └── images      # Image files (sprites, textures)
│
├── cmd             # Main applications (entry points)
│   └── harvest     # Specific entry point for the Harvest game
│       └── main.go # Main function (initialization and start of the game loop)
│
├── internal        # Internal core game code (not intended for external use)
│   ├── assets      # Runtime asset management (loading, caching)
│   ├── collision   # Logic for collision detection
│   ├── component   # Reusable data components for entities (e.g., position, health)
│   ├── enemy       # Definition and behavior of enemy entities
│   ├── entity      # Base definitions for game entities
│   ├── game        # Main game state management and implementation of the ebitengine.Game interface
│   ├── inventory   # Logic for the player inventory
│   ├── item        # Definition of in-game items
│   ├── itemtype    # Categorization and typing of items
│   ├── player      # Definition and behavior of the player entity
│   ├── weapon      # Definition and behavior of weapons
│   └── world       # Management of the game world, camera movement and spawning of entities
│
├── pkg             # Public libraries (if any, here for configuration)
│   └── config      # Loading and managing game configurations
│
├── go.mod          # Go module definition (dependencies)
├── go.sum          # Dependency checksums
├── LICENSE         # Project license file (Please add one!)
├── Makefile        # Automation of build and execution commands
└── README.md       # This file
```

### Explanation of Key Directories

*   **`cmd/harvest/main.go`**: The entry point of the application. This is where the Ebitengine window is initialized and the main game structure (`internal/game.Game`) is passed to `ebitengine.RunGame`.
*   **`internal/game`**: Contains the central `Game` struct that implements the `ebitengine.Game` interface (`Update`, `Draw`, `Layout`). It coordinates the various subsystems like world, player, enemies, and assets.
*   **`internal/world`**: Manages the state of the game world, including the camera focus and the tile Management
*   **`internal/entity`, `internal/component`, `internal/player`, `internal/enemy`**: Implement a form of Entity-Component-System (ECS) or a similar architecture. `entity` defines the base, `component` the reusable data blocks, and `player`/`enemy` specialize the behavior for specific entity types.
*   **`internal/assets`**: Responsible for loading and providing game assets (images, sounds) at runtime, often using Ebitengine's helper functions.
*   **`assets`**: Contains the raw asset data, which are loaded at runtime by the `internal/assets` package.
*   **`pkg/config`**: Contains reusable code for loading configuration, potentially usable by other projects (though often kept internal if specific to the game).

This structure promotes separation of concerns, making the codebase easier to navigate, maintain, and test.

## License

The project is licensed under the MIT License.
