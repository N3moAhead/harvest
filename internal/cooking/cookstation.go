package cooking

import (
	"fmt"
	"image/color"

	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/N3moAhead/harvest/internal/entity/player/inventory"
	"github.com/N3moAhead/harvest/internal/soups"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Recipe struct {
	Soup        itemtype.ItemType
	Ingredients map[itemtype.ItemType]int
}

type CookStation struct {
	entity.Entity
	Recipe         Recipe
	Used           bool
	CostFactor     float64 // to scale the cost of ingredients with the game difficulty
	animationStore *animation.AnimationStore
}

var RecipeDefinitions = map[itemtype.ItemType]Recipe{
	itemtype.DamageSoup: {
		Soup: itemtype.DamageSoup,
		Ingredients: map[itemtype.ItemType]int{
			itemtype.Carrot: 2,
			itemtype.Potato: 1,
		},
	},
	itemtype.MagnetRadiusSoup: {
		Soup: itemtype.MagnetRadiusSoup,
		Ingredients: map[itemtype.ItemType]int{
			itemtype.Potato: 2,
		},
	},
	itemtype.SpeedSoup: {
		Soup: itemtype.SpeedSoup,
		Ingredients: map[itemtype.ItemType]int{
			itemtype.Potato: 2,
		},
	},
	// ...
}

func NewCookStation(x, y float64, recipe Recipe, costFactor float64) *CookStation {
	animationStore := animation.NewAnimationStore()
	cookStation, ok := assets.AssetStore.GetImage("cook_station")
	if ok {
		defaultAnimation, err := animation.NewAnimation(cookStation, 32, 32, 0, 27, 8, 6, true)
		// defaultAnimation, err := animation.NewAnimation(cookStation, 32, 32, 0, 32, 8, 6, false)
		// defaultAnimation, err := animation.NewAnimation(cookStation, 32, 32, 0, 0, 8, 6, false)
		if err == nil {
			animationStore.AddAnimation("default", defaultAnimation)
		}
		ok := animationStore.SetCurrentAnimation("default")
		if !ok {
			fmt.Println("Warning: Unable to start the spawning animation")
		}
	}

	baseEntity := entity.NewEntity(x, y)
	fmt.Printf("Creating CookStation at (%.2f, %.2f) with recipe %s\n", x, y, recipe.Soup.String())
	return &CookStation{
		Entity:         *baseEntity,
		Recipe:         recipe,
		CostFactor:     costFactor,
		animationStore: animationStore,
	}
}

func (cookStation *CookStation) Update(player *player.Player, inv *inventory.Inventory) {
	if cookStation.Used {
		return
	}
	cookStation.animationStore.Update()
	cookStation.animationStore.SetCurrentAnimation("default")
	diff := player.Pos.Sub(cookStation.Pos)
	if diff.Len() < config.PLAYER_PICKUP_RADIUS { // TODO maybe change to own e.g. PLAYER_INTERACT_RADIUS ?
		ok := true
		for t, amt := range cookStation.Recipe.Ingredients {
			if inv.Vegetables[t] < int(float64(amt)*cookStation.CostFactor) {
				ok = false
			}
		}
		if ok {
			// g.inventory.AddSoup(gItem.Type)
			// soup := gItem.RetrieveItemInfo().Soup
			// g.Player.ExtendOrAddSoup(soup)
			for t, amt := range cookStation.Recipe.Ingredients {
				inv.Vegetables[t] -= int(float64(amt) * cookStation.CostFactor)
			}
			inv.AddSoup(cookStation.Recipe.Soup)
			player.ExtendOrAddSoup(soups.Definitions[cookStation.Recipe.Soup])
			cookStation.Used = true
		}
	}
}

func (cs *CookStation) Draw(screen *ebiten.Image, camX, camY float64) {
	if cs.Used {
		return
	}
	frameImage := cs.animationStore.GetImage()
	if frameImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(cs.Pos.X-camX-16.0, cs.Pos.Y-camY-16.0)
		screen.DrawImage(frameImage, op)

	} else {
		cs.DefaultDraw(screen, camX, camY, config.CARROT_WIDTH, config.CARROT_HEIGHT,
			color.RGBA{R: 180, G: 13, B: 27, A: 255})
	}

	// TODO other way to draw text recept
	// recept sign?
	x := float32(cs.Pos.X - camX + 20)
	y := float32(cs.Pos.Y - camY - 20)
	textRecipe := cs.Recipe.Soup.String() + ": "
	for t, amt := range cs.Recipe.Ingredients {
		textRecipe += fmt.Sprintf("%s x%d\n", t.String(), int(float64(amt)*cs.CostFactor))
	}

	var fontFace font.Face = basicfont.Face7x13
	text.Draw(screen, textRecipe, fontFace, int(x), int(y), color.White)

}

func (cs *CookStation) DefaultDraw(screen *ebiten.Image, camX, camY float64, width int, height int, color color.RGBA) {
	x := float32(cs.Pos.X - camX)
	y := float32(cs.Pos.Y - camY)
	vector.DrawFilledRect(
		screen,
		x, y,
		float32(width), float32(height),
		color,
		false,
	)
}
