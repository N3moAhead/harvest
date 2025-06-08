package toast

import (
	"fmt"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

var toasts []*Toast = make([]*Toast, 0)

func AddToast(txt string) {
	fnt, ok := assets.AssetStore.GetFont("micro")
	if ok {
		newToast := newToast(txt, fnt, config.DEFAULT_TOAST_DURATION)
		toasts = append(toasts, newToast)
	} else {
		fmt.Println("Warning: Could not load font in AddToast")
	}
}

func AddCustomToast(txt string, fnt font.Face, duration time.Duration) {
	newToast := newToast(txt, fnt, config.DEFAULT_TOAST_DURATION)
	toasts = append(toasts, newToast)
}

func UpdateToasts() {
	n := 0
	for i, toast := range toasts {
		isAlive := toast.Update()
		if isAlive {
			if i != n {
				toasts[n] = toast
			}
			n++
		}
	}
	toasts = toasts[:n]
}

func DrawToasts(screen *ebiten.Image) {
	paddingTop := 0
	for _, toast := range toasts {
		toast.Draw(screen, paddingTop)
		paddingTop += 30
	}
}
