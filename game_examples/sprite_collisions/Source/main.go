// Sprite Collisions example for Playdate in Go
// Demonstrates sprite collision detection and raycasting
// Ported from Playdate SDK C_API/Examples/Sprite Collisions

package main

import (
	"math"
	"unsafe"

	"github.com/playdate-go/pdgo"
)

var (
	pd *pdgo.PlaydateAPI

	player1 *pdgo.LCDSprite
	player2 *pdgo.LCDSprite

	p1VelocityX float32 = 60.0
	p1VelocityY float32 = 50.0
	p2VelocityX float32 = 20.0
	p2VelocityY float32 = 30.0

	rayRotation float32 = 0.0

	dt float32 = 0.05

	// Store player pointers for comparison in callbacks
	player1Ptr unsafe.Pointer
	player2Ptr unsafe.Pointer

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
	pd.Display.SetRefreshRate(50)

	// Create two bouncing players
	player1 = createPlayer(190, 223, 20, 20)
	player2 = createPlayer(100, 40, 16, 16)

	// Store pointers for comparison in callbacks
	player1Ptr = player1.Ptr()
	player2Ptr = player2.Ptr()

	// Create border walls
	borderSize := float32(5)
	displayWidth := float32(400)
	displayHeight := float32(240)

	createBlock(0, 0, displayWidth, borderSize)                                              // top
	createBlock(0, borderSize, borderSize, displayHeight-borderSize*2)                       // left
	createBlock(displayWidth-borderSize, borderSize, borderSize, displayHeight-borderSize*2) // right
	createBlock(0, displayHeight-borderSize, displayWidth, borderSize)                       // bottom

	// Create random obstacle blocks
	for i := 0; i < 6; i++ {
		createBlock(
			float32(randInt(270)+50),
			float32(randInt(100)+50),
			float32(randInt(30)+10),
			float32(randInt(90)+10),
		)
	}
}

// createPlayer creates a bouncing player sprite
func createPlayer(x, y, w, h float32) *pdgo.LCDSprite {
	sprite := pd.Sprite.NewSprite()

	// Set bounds centered at x,y
	bounds := pdgo.PDRect{
		X:      x - w/2,
		Y:      y - h/2,
		Width:  w,
		Height: h,
	}
	pd.Sprite.SetBounds(sprite, bounds)

	// Set up collision rectangle (relative to bounds)
	collideRect := pdgo.PDRect{X: 0, Y: 0, Width: w, Height: h}
	pd.Sprite.SetCollideRect(sprite, collideRect)

	// Set collision response to bounce
	pd.Sprite.SetCollisionResponseFunction(sprite, func(_, _ *pdgo.LCDSprite) pdgo.SpriteCollisionResponseType {
		return pdgo.CollisionTypeBounce
	})

	// Set draw function
	pd.Sprite.SetDrawFunction(sprite, func(s *pdgo.LCDSprite, bounds, drawRect pdgo.PDRect) {
		pd.Graphics.FillRect(int(bounds.X), int(bounds.Y), int(bounds.Width), int(bounds.Height), pdgo.SolidBlack)
	})

	// Set update function
	pd.Sprite.SetUpdateFunction(sprite, func(s *pdgo.LCDSprite) {
		updatePlayer(s)
	})

	pd.Sprite.SetZIndex(sprite, 1000)
	pd.Sprite.AddSprite(sprite)

	return sprite
}

// createBlock creates a static block sprite
func createBlock(x, y, w, h float32) *pdgo.LCDSprite {
	sprite := pd.Sprite.NewSprite()

	bounds := pdgo.PDRect{X: x, Y: y, Width: w, Height: h}
	pd.Sprite.SetBounds(sprite, bounds)

	// Set collision rectangle
	collideRect := pdgo.PDRect{X: 0, Y: 0, Width: w, Height: h}
	pd.Sprite.SetCollideRect(sprite, collideRect)

	// Set draw function
	pd.Sprite.SetDrawFunction(sprite, func(s *pdgo.LCDSprite, bounds, drawRect pdgo.PDRect) {
		pd.Graphics.DrawRect(int(bounds.X), int(bounds.Y), int(bounds.Width), int(bounds.Height), pdgo.SolidBlack)
	})

	pd.Sprite.AddSprite(sprite)

	return sprite
}

// updatePlayer updates a player sprite with collision detection
func updatePlayer(sprite *pdgo.LCDSprite) {
	var dx, dy float32
	spritePtr := sprite.Ptr()

	// Determine which player and get velocity by comparing pointers
	if spritePtr == player1Ptr {
		dx = p1VelocityX * dt
		dy = p1VelocityY * dt
	} else {
		dx = p2VelocityX * dt
		dy = p2VelocityY * dt
	}

	if dx == 0 && dy == 0 {
		return
	}

	// Get current position
	x, y := pd.Sprite.GetPosition(sprite)

	// Move with collision detection
	collisions, _, _ := pd.Sprite.MoveWithCollisions(sprite, x+dx, y+dy)

	// Handle collisions - reverse velocity based on normal
	for _, c := range collisions {
		if c.Normal.X != 0 {
			if spritePtr == player1Ptr {
				p1VelocityX = -p1VelocityX
			} else {
				p2VelocityX = -p2VelocityX
			}
		}
		if c.Normal.Y != 0 {
			if spritePtr == player1Ptr {
				p1VelocityY = -p1VelocityY
			} else {
				p2VelocityY = -p2VelocityY
			}
		}
	}
}

// drawRays draws raycasting lines from a player
func drawRays(player *pdgo.LCDSprite) {
	startX, startY := pd.Sprite.GetPosition(player)

	for extraAngle := 0; extraAngle < 360; extraAngle += 30 {
		rads := float64(rayRotation+float32(extraAngle)) * math.Pi / 180.0
		endX := int(500*math.Cos(rads)) + int(startX)
		endY := int(500*math.Sin(rads)) + int(startY)

		// Query sprites along the ray
		results := pd.Sprite.QuerySpriteInfoAlongLine(startX, startY, float32(endX), float32(endY))

		// Draw collision points
		for _, info := range results {
			// Skip player sprites by comparing pointers
			infoPtr := info.Sprite.Ptr()
			if infoPtr == player1Ptr || infoPtr == player2Ptr {
				continue
			}

			// Draw entry point ellipse
			r := int(6 - 10*info.Ti1)
			if r > 0 {
				pd.Graphics.DrawEllipse(
					int(info.EntryPoint.X)-r, int(info.EntryPoint.Y)-r,
					r*2, r*2, 3,
					0, 360, pdgo.SolidBlack,
				)
			}

			// Draw exit point ellipse
			r = int(5 - 10*info.Ti1)
			if r > 0 {
				pd.Graphics.DrawEllipse(
					int(info.ExitPoint.X)-r, int(info.ExitPoint.Y)-r,
					r*2, r*2, 1,
					0, 360, pdgo.SolidBlack,
				)
			}
		}

		// Draw the ray line
		pd.Graphics.DrawLine(int(startX), int(startY), endX, endY, 1, pdgo.SolidBlack)
	}
}

// update is called every frame
func update() int {
	// Update and draw all sprites (calls update callbacks)
	pd.Sprite.UpdateAndDrawSprites()

	// Draw raycasting from both players
	drawRays(player1)
	drawRays(player2)

	// Draw FPS
	pd.System.DrawFPS(8, 8)

	// Rotate rays
	rayRotation += 0.5

	return 1
}

func main() {}
