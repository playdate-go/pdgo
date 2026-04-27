// Exposure example for Playdate in Go
// Demonstrates sprite coverage detection, accelerometer input, and pixel counting
// Ported from Playdate SDK C_API/Examples/Exposure

package main

import (
	"math/bits"

	"github.com/playdate-go/pdgo"
)

const (
	TargetSize = 100
	CoverSize  = 60
	CoverCount = 100
)

var (
	pd *pdgo.PlaydateAPI

	// Sprites
	target     *pdgo.LCDSprite
	targetMask *pdgo.LCDBitmap
	covers     []*coverSprite

	// Images
	testImage *pdgo.LCDBitmap

	// Pattern for striped fill (8 bytes for 8x8 pattern)
	stripesPattern = []byte{0x0f, 0x1e, 0x3c, 0x78, 0xf0, 0xe1, 0xc3, 0x87}

	// Total pixel count for coverage calculation
	totalCount int

	// Accelerometer smoothing
	shiftX, shiftY float32
	smoothing      float32 = 0.9
)

type coverSprite struct {
	sprite      *pdgo.LCDSprite
	spriteImage *pdgo.LCDBitmap
	mask        *pdgo.LCDBitmap
	px, py      float32 // position offset
	pz          float32 // depth for parallax
}

// XorShift random state
var randState uint32 = 12345

func random() uint32 {
	x := randState
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	randState = x
	return x
}

func randFloat() float32 {
	return float32(random()) / float32(0xFFFFFFFF)
}

func main() {}

//export initGame
func initGame() {
	// Run as fast as possible for smooth animation
	pd.Display.SetRefreshRate(0)

	// Create the striped pattern bitmap for filling
	patternBitmap := createPatternBitmap()

	// Create target sprite (black striped circle)
	blackCircle := pd.Graphics.NewBitmap(TargetSize, TargetSize, pdgo.SolidWhite)
	pd.Graphics.PushContext(blackCircle)
	// Fill with pattern by using stencil
	drawPatternedEllipse(TargetSize/2, TargetSize/2, TargetSize, TargetSize, patternBitmap)
	pd.Graphics.DrawEllipse(0, 0, TargetSize, TargetSize, 1, 0, 360, pdgo.SolidBlack)
	pd.Graphics.PopContext()

	// Create target sprite
	target = pd.Sprite.NewSprite()
	pd.Sprite.SetImage(target, blackCircle, pdgo.BitmapUnflipped)
	pd.Sprite.SetCollideRect(target, pdgo.PDRect{X: 0, Y: 0, Width: TargetSize, Height: TargetSize})
	pd.Sprite.AddSprite(target)
	pd.Sprite.MoveTo(target, 200, 120)

	// Cache target mask for coverage calculation
	targetMask = pd.Graphics.GetBitmapMask(blackCircle)

	// Create white circle image for covers
	whiteCircle := pd.Graphics.NewBitmap(CoverSize, CoverSize, pdgo.SolidWhite)
	pd.Graphics.PushContext(whiteCircle)
	pd.Graphics.DrawEllipse(0, 0, CoverSize, CoverSize, 1, 0, 360, pdgo.SolidBlack)
	pd.Graphics.PopContext()

	// Create cover sprites
	covers = make([]*coverSprite, CoverCount)
	for i := 0; i < CoverCount; i++ {
		c := &coverSprite{}
		c.sprite = pd.Sprite.NewSprite()
		c.spriteImage = pd.Graphics.CopyBitmap(whiteCircle)
		pd.Sprite.SetImage(c.sprite, c.spriteImage, pdgo.BitmapUnflipped)
		pd.Sprite.SetCollideRect(c.sprite, pdgo.PDRect{X: 0, Y: 0, Width: CoverSize, Height: CoverSize})

		// Get and invert the mask
		originalMask := pd.Graphics.GetBitmapMask(c.spriteImage)
		if originalMask != nil {
			c.mask = invertBitmap(originalMask)
		}

		covers[i] = c
		pd.Sprite.AddSprite(c.sprite)
	}

	// Create test image for coverage calculation
	testImage = pd.Graphics.NewBitmap(TargetSize, TargetSize, pdgo.SolidBlack)

	// Calculate total pixels in target mask
	pd.Graphics.PushContext(testImage)
	if targetMask != nil {
		pd.Graphics.DrawBitmap(targetMask, 0, 0, pdgo.BitmapUnflipped)
	}
	pd.Graphics.PopContext()
	totalCount = popCount(testImage)

	// Randomize initial cover positions
	randomizeCovers()

	// Start accelerometer
	pd.System.SetPeripheralsEnabled(pdgo.PeripheralAccelerometer)

	// Set update callback
	pd.System.SetUpdateCallback(update)
}

func createPatternBitmap() *pdgo.LCDBitmap {
	// Create an 8x8 bitmap with the stripe pattern
	bitmap := pd.Graphics.NewBitmap(8, 8, pdgo.SolidWhite)
	data := pd.Graphics.GetBitmapData(bitmap)
	if data != nil && data.Data != nil {
		// The pattern is stored as 8 bytes, one per row
		// Each byte represents 8 pixels (1 bit per pixel)
		for y := 0; y < 8 && y < len(stripesPattern); y++ {
			rowStart := y * data.RowBytes
			if rowStart < len(data.Data) {
				data.Data[rowStart] = stripesPattern[y]
			}
		}
	}
	return bitmap
}

func drawPatternedEllipse(cx, cy, w, h int, pattern *pdgo.LCDBitmap) {
	// Use the pattern as a stencil to fill the ellipse
	pd.Graphics.SetStencilImage(pattern, true)
	pd.Graphics.FillEllipse(cx-w/2, cy-h/2, w, h, 0, 360, pdgo.SolidBlack)
	// Clear stencil
	pd.Graphics.SetStencilImage(nil, false)
}

func invertBitmap(bitmap *pdgo.LCDBitmap) *pdgo.LCDBitmap {
	if bitmap == nil {
		return nil
	}

	data := pd.Graphics.GetBitmapData(bitmap)
	if data == nil || data.Data == nil {
		return nil
	}

	// Create a new bitmap and invert its data
	inverted := pd.Graphics.NewBitmap(data.Width, data.Height, pdgo.SolidWhite)
	invData := pd.Graphics.GetBitmapData(inverted)

	if invData != nil && invData.Data != nil {
		for i := 0; i < len(data.Data) && i < len(invData.Data); i++ {
			invData.Data[i] = ^data.Data[i]
		}
	}

	return inverted
}

func randomizeCovers() {
	for _, c := range covers {
		c.px = 4*randFloat() - 2
		c.py = 4*randFloat() - 2
		c.pz = 3*randFloat() + 0.2
		pd.Sprite.SetZIndex(c.sprite, int16(c.pz*1000))
	}
}

func update() int {
	// Read and smooth accelerometer data
	accelX, accelY, _ := pd.System.GetAccelerometer()
	shiftX = smoothing*shiftX + (1-smoothing)*accelX
	shiftY = smoothing*shiftY + (1-smoothing)*accelY

	// Check for A/B button to randomize
	_, pushed, _ := pd.System.GetButtonState()
	if pushed&(pdgo.ButtonA|pdgo.ButtonB) != 0 {
		randomizeCovers()
	}

	// Update cover sprite positions with parallax
	for _, c := range covers {
		// Perspective transform based on depth
		x := 200 + 200*(c.px+shiftX*c.pz)
		y := 120 + 120*(c.py+shiftY*c.pz)
		pd.Sprite.MoveTo(c.sprite, x, y)
	}

	// Update sprite scene
	pd.Sprite.UpdateAndDrawSprites()

	// Calculate coverage
	overlapping := pd.Sprite.OverlappingSprites(target)
	visible := float32(1.0)

	if len(overlapping) > 0 {
		// Reset test image
		pd.Graphics.PushContext(testImage)
		pd.Graphics.Clear(pdgo.SolidBlack)

		// Draw target mask
		if targetMask != nil {
			pd.Graphics.DrawBitmap(targetMask, 0, 0, pdgo.BitmapUnflipped)
		}

		// Draw inverted cover masks to show coverage
		oldMode := pd.Graphics.SetDrawMode(pdgo.DrawModeWhiteTransp)
		for _, s := range overlapping {
			// Find the cover sprite
			for _, c := range covers {
				if c.sprite.Ptr() == s.Ptr() && c.mask != nil {
					x, y := pd.Sprite.GetPosition(c.sprite)
					// Offset to draw relative to target position
					drawX := int(x - CoverSize/2 - (200 - TargetSize/2))
					drawY := int(y - CoverSize/2 - (120 - TargetSize/2))
					pd.Graphics.DrawBitmap(c.mask, drawX, drawY, pdgo.BitmapUnflipped)
					break
				}
			}
		}
		pd.Graphics.SetDrawMode(oldMode)
		pd.Graphics.PopContext()

		// Calculate visible percentage
		count := popCount(testImage)
		if totalCount > 0 {
			visible = float32(count) / float32(totalCount)
		}
	}

	// Draw coverage meter
	drawCoverageMeter(visible)

	// Draw test image preview
	pd.Graphics.DrawBitmap(testImage, 0, 140, pdgo.BitmapUnflipped)
	pd.Graphics.DrawLine(0, 139, 100, 139, 1, pdgo.SolidBlack)
	pd.Graphics.DrawLine(100, 139, 100, 240, 1, pdgo.SolidBlack)

	// Draw FPS
	pd.System.DrawFPS(0, 0)

	return 1
}

func drawCoverageMeter(visible float32) {
	// Meter position
	meterX := 380
	meterWidth := 20
	meterHeight := 240

	// Calculate filled height
	y := int(float32(meterHeight) * visible)

	// Clear top part (uncovered)
	pd.Graphics.FillRect(meterX, 0, meterWidth, meterHeight-y, pdgo.SolidWhite)

	// Fill bottom part with pattern (covered)
	// Draw stripes manually
	for row := meterHeight - y; row < meterHeight; row += 8 {
		for col := 0; col < meterWidth; col++ {
			patternRow := (row / 8) % 8
			if patternRow < len(stripesPattern) {
				bit := (stripesPattern[patternRow] >> (col % 8)) & 1
				if bit != 0 {
					pd.Graphics.FillRect(meterX+col, row, 1, 1, pdgo.SolidBlack)
				}
			}
		}
	}

	// Draw meter outline
	pd.Graphics.DrawLine(meterX, 0, meterX, meterHeight, 1, pdgo.SolidBlack)
}

// popCount counts the number of white pixels (set bits) in the bitmap
func popCount(bitmap *pdgo.LCDBitmap) int {
	if bitmap == nil {
		return 0
	}

	data := pd.Graphics.GetBitmapData(bitmap)
	if data == nil || data.Data == nil {
		return 0
	}

	count := 0
	height := data.Height
	rowbytes := data.RowBytes

	for y := 0; y < height; y++ {
		rowOffset := y * rowbytes
		// Process 4 bytes (32 bits) at a time
		for x := 0; x < rowbytes; x += 4 {
			// Pack 4 bytes into a uint32
			var val uint32
			for i := 0; i < 4 && (x+i) < rowbytes; i++ {
				idx := rowOffset + x + i
				if idx < len(data.Data) {
					val |= uint32(data.Data[idx]) << (i * 8)
				}
			}
			// Count set bits
			count += bits.OnesCount32(val)
		}
	}

	return count
}
