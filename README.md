# harvest

harvest/
├── main.go                 # Haupteinstiegspunkt der Anwendung
├── go.mod                  # Go Modul Definition
├── go.sum                  # Abhängigkeits-Checksummen
│
├── assets/                 # Verzeichnis für alle Spiel-Assets
│   ├── images/             # Bilder (Spritesheets, Hintergründe, UI-Elemente)
│   │   ├── player/
│   │   ├── enemies/
│   │   │   ├── tomato.png
│   │   │   ├── broccoli.png
│   │   │   └── ...
│   │   ├── weapons/
│   │   ├── projectiles/
│   │   ├── items/
│   │   └── ui/
│   ├── audio/              # Soundeffekte und Musik
│   │   ├── sfx/
│   │   └── music/
│   └── fonts/              # Schriftarten
│
├── cmd/                    # Hauptanwendungen (falls mehr als eine benötigt wird, z.B. Server/Client oder Tools)
│   └── harvest/
│       └── main.go         # (Alternative zum main.go im Root, oft für größere Projekte genutzt)
│
├── internal/               # Private Pakete, die nur innerhalb dieses Projekts verwendet werden sollen
│   ├── game/               # Kern-Spiellogik und -Zustand
│   │   ├── game.go         # Haupt-Game-Struktur (implementiert ebiten.Game), Update/Draw Loop
│   │   └── state.go        # Verwaltung des Spielzustands (Laufen, Pause, Game Over, Level Up)
│   │
│   ├── world/              # Verwaltung der Spielwelt und Entitäten darin
│   │   ├── world.go        # Enthält alle aktiven Entitäten (Player, Enemies, Projectiles, Items)
│   │   ├── spawner.go      # Logik zum Spawnen von Gegnern (Wellen, Timing, Positionierung)
│   │   └── camera.go       # Spielfeld-Kamera (falls benötigt, folgt meist dem Spieler)
│   │
│   ├── entity/             # Basisdefinitionen für alle "Dinge" im Spiel
│   │   └── entity.go       # Basis-Interface oder Struct für alle Entitäten (Position, Größe, etc.)
│   │
│   ├── player/             # Alles rund um den Spielercharakter
│   │   ├── player.go       # Spieler-Struct, Bewegung, Stats (Leben, Speed, etc.), Waffen-Handling
│   │   └── input.go        # Verarbeitung von Spieler-Input (Bewegung)
│   │
│   ├── enemy/              # Alles rund um die Gegner (Gemüse)
│   │   ├── enemy.go        # Basis-Interface/Struct für Gegner, gemeinsame Logik
│   │   ├── types.go        # Definitionen spezifischer Gegner-Typen (Tomate, Brokkoli, etc.)
│   │   └── ai.go           # KI-Verhalten (z.B. Spieler verfolgen)
│   │
│   ├── weapon/             # Waffen und Angriffslogik
│   │   ├── weapon.go       # Basis-Interface/Struct für Waffen, Cooldowns
│   │   └── types.go        # Definitionen spezifischer Waffen (Knoblauch-Aura, Karotten-Werfer, etc.)
│   │
│   ├── projectile/         # Geschosse, die von Waffen erzeugt werden
│   │   ├── projectile.go   # Basis-Interface/Struct für Geschosse, Bewegung, Lebensdauer
│   │   └── types.go        # Definitionen spezifischer Geschoss-Typen
│   │
│   ├── item/               # Aufsammelbare Gegenstände
│   │   ├── item.go         # Basis-Interface/Struct für Items (XP-Kristalle, Heilung, Truhen)
│   │   └── types.go        # Definitionen spezifischer Item-Typen und deren Effekte
│   │
│   ├── ui/                 # User Interface Elemente
│   │   ├── ui.go           # Haupt-UI-Manager (zeichnet HUD, Menüs)
│   │   ├── hud.go          # Head-Up Display (Lebensanzeige, XP-Balken, Timer, Waffenicons)
│   │   ├── levelup.go      # Level-Up Auswahlbildschirm
│   │   └── gameover.go     # Game Over Bildschirm
│   │
│   ├── component/          # (Optional) Entity-Component-System (ECS) Bausteine, falls ihr ECS nutzt
│   │   ├── position.go
│   │   ├── velocity.go
│   │   ├── health.go
│   │   └── ...
│   │
│   ├── system/             # (Optional) ECS Systeme
│   │   ├── movement.go
│   │   ├── collision.go
│   │   ├── rendering.go
│   │   └── ...
│
├── pkg/                    # Öffentliche Bibliotheken (falls ihr Code habt, der auch von anderen Projekten genutzt werden könnte)
│   ├── assets/             # Ladefunktionen und Verwaltung für Assets (Bilder, Audio)
│   │   └── loader.go
│   ├── config/             # Laden und Verwalten von Konfigurationen (z.B. aus Datei)
│   │   └── config.go
│   ├── collision/          # Kollisionserkennungs-Helfer (z.B. Rect vs Rect, Circle vs Rect)
│   │   └── collision.go
│   ├── vector/             # Einfache 2D-Vektor-Mathematik
│   │   └── vector.go
│   └── animation/          # Hilfsfunktionen für Sprite-Animationen
│       └── animation.go
│
└── config.toml             # (Optional) Konfigurationsdatei (Fenstergröße, Soundlautstärke etc.)


Beschreibung der Pakete und Verzeichnisse:

gemuese-survivor/ (Root)

main.go: Der Haupteinstiegspunkt. Initialisiert Ebitengine, lädt Konfigurationen und Assets (über die entsprechenden Pakete), erstellt das Haupt-Game-Objekt aus internal/game und startet die Ebitengine-Schleife (ebiten.RunGame).

go.mod, go.sum: Go-Modul-Dateien zur Verwaltung von Abhängigkeiten.

assets/: Enthält alle statischen Spieldateien (Grafiken, Sounds, Schriften). Gut strukturiert, damit man Assets leicht findet. Wird nicht kompiliert, sondern zur Laufzeit vom Programm geladen (meist über pkg/assets).

config.toml (optional): Eine Konfigurationsdatei für Einstellungen, die nicht fest im Code stehen sollen.

cmd/gemuese-survivor/ (Optional)

Eine gängige Konvention, um den main Code zu isolieren, besonders wenn man plant, später weitere ausführbare Dateien (z.B. einen Level-Editor) hinzuzufügen. Der Inhalt wäre derselbe wie main.go im Root-Verzeichnis. Wenn ihr nur eine ausführbare Datei habt, ist main.go im Root auch völlig in Ordnung.

internal/

Dieses Verzeichnis signalisiert Go, dass die darin enthaltenen Pakete nur für den internen Gebrauch innerhalb des gemuese-survivor-Moduls bestimmt sind. Sie können nicht von externen Projekten importiert werden.

game/: Das Herzstück des Spiels.

game.go: Definiert die Hauptstruktur (z.B. Game), die das ebiten.Game-Interface implementiert (Update, Draw, Layout). Koordiniert die Logik, indem es Methoden der anderen Pakete (world, player, ui etc.) aufruft. Hält den globalen Spielzustand (Punkte, Level, Zeit).

state.go: Verwaltet verschiedene Spielzustände (z.B. Running, Paused, LevelUp, GameOver) und steuert den Übergang zwischen ihnen.

world/: Verantwortlich für den Spielbereich und die darin befindlichen dynamischen Objekte.

world.go: Enthält Listen/Slices aller aktiven Entitäten (Gegner, Projektile, Items) und des Spielers. Bietet Methoden zum Hinzufügen/Entfernen von Entitäten. Führt Kollisionsabfragen durch (ruft Helfer aus pkg/collision auf). Definiert vielleicht die Grenzen der Spielwelt.

spawner.go: Logik, die entscheidet, wann, wo und welche Gemüse-Gegner erscheinen sollen, basierend auf der Spielzeit oder dem Spielerlevel.

camera.go: Falls die Welt größer als der Bildschirm ist, verwaltet die Kamera die Position des sichtbaren Ausschnitts (meist zentriert auf den Spieler).

entity/: (Optional, aber oft nützlich)

entity.go: Eine Basisdefinition, von der Player, Enemy, Projectile, Item erben oder die sie implementieren können. Könnte grundlegende Eigenschaften wie Position, Größe, eine eindeutige ID und vielleicht eine Update()- und Draw()-Methode definieren.

player/: Alles spezifisch für den Spieler.

player.go: Definiert die Spielerstruktur mit Attributen (HP, Geschwindigkeit, Sammelradius etc.), Methoden für Bewegung, Schadensaufnahme, Level-Up und das Verwalten der ausgerüsteten Waffen.

input.go: Liest Ebitengine-Input und übersetzt ihn in Spieleraktionen (z.B. Bewegungsvektor).

enemy/: Alles spezifisch für die Gemüse-Gegner.

enemy.go: Basis-Interface oder -Struct für alle Gegner. Gemeinsame Funktionen wie TakeDamage().

types.go: Konkrete Implementierungen für jeden Gegnertyp (z.B. Tomato, Broccoli) mit ihren spezifischen Werten (HP, Geschwindigkeit, Schaden) und vielleicht Sprite-Referenzen.

ai.go: Funktionen oder Methoden für das Gegnerverhalten (z.B. ChasePlayer, SimpleMove).

weapon/: Logik für die Waffen des Spielers.

weapon.go: Basis-Interface oder -Struct für Waffen. Definiert Attribute wie Schaden, Feuerrate, Reichweite, Cooldown. Enthält Logik zum Auslösen des Angriffs (z.B. Erzeugen von Projektilen).

types.go: Konkrete Implementierungen für jede Waffe (GarlicAura, CarrotLauncher, PeaShooter).

projectile/: Für Dinge, die von Waffen abgefeuert werden.

projectile.go: Basis-Interface oder -Struct für Geschosse. Definiert Bewegungsmuster, Lebensdauer, Schaden und Kollisionslogik mit Gegnern.

types.go: Konkrete Implementierungen für verschiedene Geschossarten.

item/: Für aufsammelbare Gegenstände.

item.go: Basis-Interface oder -Struct für Items. Definiert, wie sie sich verhalten (z.B. zum Spieler gezogen werden) und was passiert, wenn sie eingesammelt werden.

types.go: Konkrete Implementierungen (XpGem, HealthUp, TreasureChest).

ui/: Zeichnen von Informationen und Menüs auf dem Bildschirm.

ui.go: Koordiniert das Zeichnen der verschiedenen UI-Teile.

hud.go, levelup.go, gameover.go: Zeichnen die spezifischen UI-Elemente und Bildschirme. Greifen auf Daten aus game, player etc. zu, um aktuelle Werte anzuzeigen.

component/ & system/ (Optional - ECS): Wenn ihr einen Entity-Component-System-Ansatz verfolgt (wie z.B. mit der ecs-Bibliothek), wären hier die Komponenten (reine Daten) und Systeme (Logik, die auf Komponenten wirkt) untergebracht. Dies ist eine fortgeschrittenere Architektur.

pkg/

Dieses Verzeichnis enthält wiederverwendbare Pakete, die potenziell auch in anderen Projekten nützlich sein könnten (im Gegensatz zu internal/, das spezifisch für dieses Spiel ist).

assets/: Standardisierte Funktionen zum Laden von Bildern (*ebiten.Image), Sounds etc. aus dem assets-Verzeichnis. Kann Caching implementieren.

config/: Laden von Einstellungen aus einer Datei (z.B. config.toml oder config.json).

collision/: Allgemeine Funktionen zur Kollisionserkennung zwischen geometrischen Formen (Rechtecke sind in 2D-Spielen sehr häufig).

vector/: Grundlegende 2D-Vektoroperationen (Addition, Subtraktion, Skalierung, Normalisierung), die in vielen Bereichen (Bewegung, Physik) nützlich sind.

animation/: Helfer zum Verwalten von Frame-basierten Animationen aus Spritesheets.

Zusätzliche Hinweise:

Ebitengine Integration: Ebitengine wird hauptsächlich in main.go (Initialisierung), internal/game/game.go (Update/Draw Loop, Layout), internal/player/input.go (Input Handling) und pkg/assets/loader.go (Asset-Laden) sowie überall dort verwendet, wo gezeichnet wird (UI, Entitäten).

Abhängigkeiten: Versucht, die Abhängigkeiten zwischen den Paketen gering zu halten. Zum Beispiel sollte enemy nicht direkt von weapon abhängen. Interaktionen geschehen oft über das world-Paket (z.B. world prüft Kollisionen zwischen projectile und enemy).

Flexibilität: Diese Struktur ist ein Vorschlag. Passt sie an eure spezifischen Bedürfnisse an. Wenn ein Paket zu groß wird, teilt es auf. Wenn zwei Pakete sehr eng zusammenarbeiten, fasst sie vielleicht zusammen.

Start Small: Ihr müsst nicht sofort alle Pakete erstellen. Beginnt mit main, game, player und vielleicht assets und fügt weitere hinzu, wenn ihr sie braucht.

Viel Erfolg bei der Entwicklung eures Gemüse-Survivor-Klons!
