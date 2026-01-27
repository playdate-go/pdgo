// Hello World example for Playdate in Go
// Bouncing "Hello World!" text

package main

import (
	"github.com/playdate-go/pdgo"
)

const (
	textWidth  = 86
	textHeight = 16
)

var (
	pd   *pdgo.PlaydateAPI
	font *pdgo.LCDFont

	x  = (pdgo.LCDColumns - textWidth) / 2
	y  = (pdgo.LCDRows - textHeight) / 2
	dx = 1
	dy = 2
)

// initGame is called once when the game starts
func initGame() {
	// Load font
	var err error
	font, err = pd.Graphics.LoadFont("/System/Fonts/Asheville-Sans-14-Bold.pft")
	if err != nil {
		pd.System.Error("Couldn't load font")
	}
}

// update is called every frame
func update() int {
	pd.Graphics.Clear(pdgo.SolidWhite)

	if font != nil {
		pd.Graphics.SetFont(font)
	}
	pd.Graphics.DrawText("Hello World!", x, y)

	// Move text
	x += dx
	y += dy

	// Bounce off edges
	if x < 0 || x > pdgo.LCDColumns-textWidth {
		dx = -dx
	}

	if y < 0 || y > pdgo.LCDRows-textHeight {
		dy = -dy
	}

	// Draw FPS counter
	pd.System.DrawFPS(0, 0)

	return 1
}

func main() {}
