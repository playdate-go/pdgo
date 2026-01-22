package main

import (
	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI

var x, y int32 = 200, 120
var dx, dy int32 = 2, 3
var size int32 = 20

func update() int {
	pd.Graphics.Clear(pdgo.NewColorFromSolid(pdgo.ColorWhite))

	pd.Graphics.FillRect(int(x), int(y), int(size), int(size), pdgo.NewColorFromSolid(pdgo.ColorBlack))

	x += dx
	y += dy

	if x < 0 || x > int32(pdgo.LCDColumns)-size {
		dx = -dx
	}
	if y < 0 || y > int32(pdgo.LCDRows)-size {
		dy = -dy
	}

	pd.System.DrawFPS(0, 0)

	return 1
}

func initGame() {}

func main() {}
