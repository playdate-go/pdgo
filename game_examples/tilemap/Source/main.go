// Tilemap example for Playdate in Go
// Demonstrates using tilemaps with bitmap tables
// Ported from Playdate SDK C_API/Examples/Tilemap

package main

import (
	"github.com/playdate-go/pdgo"
)

var (
	pd      *pdgo.PlaydateAPI
	tilemap *pdgo.LCDTileMap

	// XorShift random state
	randState uint32 = 12345
)

// random returns a pseudo-random uint32 using XorShift
func random() uint32 {
	x := randState
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	randState = x
	return x
}

// randInt returns a random number in range [0, n)
func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	return int(random() % uint32(n))
}

// initGame is called once when the game starts
func initGame() {
	// Create a new tilemap
	tilemap = pd.Graphics.NewTilemap()

	// Set the tilemap size (20 tiles wide x 12 tiles high)
	// With 20x20 pixel tiles, this fills the 400x240 screen
	tilemap.SetSize(20, 12)

	// Load the bitmap table containing tile images
	// The "font" bitmap table contains character images
	table, err := pd.Graphics.LoadBitmapTable("font")
	if err != nil {
		pd.System.Error("Couldn't load font bitmap table")
		return
	}

	// Set the image table for the tilemap
	tilemap.SetImageTable(table)
}

// update is called every frame
func update() int {
	// Clear the screen
	pd.Graphics.Clear(pdgo.SolidWhite)

	// Randomly change tiles to create a chaotic effect
	// Set a random tile at a random position to a random image index
	tilemap.SetTileAtPosition(
		randInt(20),          // x: 0-19
		randInt(12),          // y: 0-11
		uint16(randInt(200)), // tile index: 0-199
	)

	// Draw the tilemap at the top-left corner
	tilemap.DrawAtPoint(0, 0)

	return 1
}

func main() {}
