// Go Logo example for Playdate using pdgo API
// Displays Go mascot logo on screen

package main

import (
	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI
var logo *pdgo.LCDBitmap

func initGame() {
	var err error
	logo, err = pd.Graphics.LoadBitmap("assets/go-logo")
	if err != nil {
		pd.System.Error("Failed to load go-logo image!")
	}
}

func update() int {
	// Clear screen with white
	pd.Graphics.Clear(pdgo.NewColorFromSolid(pdgo.ColorWhite))

	// Draw logo centered on screen
	if logo != nil {
		// Get bitmap dimensions
		data := pd.Graphics.GetBitmapData(logo)
		if data != nil {
			// Center the image
			x := (pdgo.LCDColumns - data.Width) / 2
			y := (pdgo.LCDRows - data.Height) / 2
			pd.Graphics.DrawBitmap(logo, x, y, pdgo.BitmapUnflipped)
		}
	}

	// FPS counter
	pd.System.DrawFPS(0, 0)

	return 1
}

func main() {}
