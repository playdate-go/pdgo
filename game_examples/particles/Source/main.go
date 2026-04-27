// Particles example for Playdate in Go
// Demonstrates a particle system with falling snowflakes
// Ported from Playdate SDK C_API/Examples/Particles

package main

import (
	"github.com/playdate-go/pdgo"
)

var (
	pd *pdgo.PlaydateAPI

	flakes    [4]*pdgo.LCDBitmap
	particles *ParticleSystem
	font      *pdgo.LCDFont

	particleCount    = 20
	showInstructions = true

	testRect = pdgo.PDRect{X: 170, Y: 90, Width: 60, Height: 60}
)

// Particle represents a single snowflake particle
type Particle struct {
	x     float32
	y     int
	w     int
	h     int
	speed int
	drift float32
	pType int // which snowflake image to use
}

// ParticleSystem manages a collection of particles
type ParticleSystem struct {
	count     int
	particles []Particle
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

func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	return int(random() % uint32(n))
}

func randFloat(n int) float32 {
	return float32((randInt(n*10) - n*5)) / 10.0
}

// initGame is called once when the game starts
func initGame() {
	// Load font
	var err error
	font, err = pd.Graphics.LoadFont("font/namco-1x")
	if err != nil {
		pd.System.Error("Couldn't load font")
	}

	// Load snowflake images
	flakes[0], err = pd.Graphics.LoadBitmap("images/snowflake1")
	if err != nil {
		pd.System.Error("Couldn't load snowflake1")
	}
	flakes[1], err = pd.Graphics.LoadBitmap("images/snowflake2")
	if err != nil {
		pd.System.Error("Couldn't load snowflake2")
	}
	flakes[2], err = pd.Graphics.LoadBitmap("images/snowflake3")
	if err != nil {
		pd.System.Error("Couldn't load snowflake3")
	}
	flakes[3], err = pd.Graphics.LoadBitmap("images/snowflake4")
	if err != nil {
		pd.System.Error("Couldn't load snowflake4")
	}

	// Create particle system
	particles = &ParticleSystem{}
	particles.SetCount(particleCount)
}

// SetCount sets the number of particles in the system
func (ps *ParticleSystem) SetCount(count int) {
	// Preserve existing particles if increasing count
	oldParticles := ps.particles
	oldCount := ps.count
	startIndex := 0

	if oldCount > 0 && count > oldCount {
		startIndex = oldCount
	} else if oldCount > 0 && count < oldCount {
		startIndex = count
	}

	ps.count = count
	ps.particles = make([]Particle, count)

	// Initialize new particles
	for i := startIndex; i < count; i++ {
		ps.particles[i] = Particle{
			x:     float32(randInt(441) - 20),
			y:     randInt(241) - 240,
			w:     19,
			h:     21,
			speed: randInt(4) + 1,
			drift: randFloat(5),
			pType: randInt(4),
		}
	}

	// Copy existing particles
	for i := 0; i < startIndex && i < oldCount; i++ {
		ps.particles[i] = oldParticles[i]
	}
}

// Update updates all particle positions
func (ps *ParticleSystem) Update() {
	for i := 0; i < ps.count; i++ {
		p := &ps.particles[i]

		// Add some random drift
		p.drift += randFloat(5) / 10.0
		p.y += p.speed
		p.x += p.drift

		// Reset particle if it goes off screen
		if p.y > 240 {
			p.y = -22
			p.x = float32(randInt(441) - 20)
			p.drift = randFloat(5)
		}
	}
}

// Draw renders all particles
func (ps *ParticleSystem) Draw() {
	pd.Graphics.Clear(pdgo.SolidBlack)

	// Use inverted draw mode for snowflakes
	pd.Graphics.SetDrawMode(pdgo.DrawModeInverted)

	for i := 0; i < ps.count; i++ {
		p := &ps.particles[i]
		if flakes[p.pType] != nil {
			pd.Graphics.DrawBitmap(flakes[p.pType], int(p.x), p.y, pdgo.BitmapUnflipped)
		}
	}

	// Reset draw mode
	pd.Graphics.SetDrawMode(pdgo.DrawModeCopy)
}

// ParticleCountInRect returns count of particles in a rectangle
func (ps *ParticleSystem) ParticleCountInRect(x, y, w, h int) int {
	count := 0
	for i := 0; i < ps.count; i++ {
		p := &ps.particles[i]
		// Check if particle intersects with rect
		if !(p.y >= y+h || p.y+p.h <= y || int(p.x) >= x+w || int(p.x)+p.w <= x) {
			count++
		}
	}
	return count
}

// update is called every frame
func update() int {
	// Handle button input
	_, pushed, _ := pd.System.GetButtonState()

	if pushed&pdgo.ButtonUp != 0 {
		particleCount += 40
		particles.SetCount(particleCount)
		showInstructions = false
	}

	if pushed&pdgo.ButtonDown != 0 {
		particleCount -= 40
		if particleCount < 20 {
			particleCount = 20
		}
		particles.SetCount(particleCount)
		showInstructions = false
	}

	// Update and draw particles
	particles.Update()
	particles.Draw()

	// Draw test rectangle (white outline)
	pd.Graphics.DrawRect(int(testRect.X), int(testRect.Y), int(testRect.Width), int(testRect.Height), pdgo.SolidWhite)

	// Draw UI text
	if font != nil {
		pd.Graphics.SetFont(font)
	}
	pd.Graphics.SetTextTracking(-1)

	// Draw particle count
	pd.Graphics.DrawText("count: "+itoa(particleCount), 4, 2)

	// Draw count in test rectangle
	boxCount := particles.ParticleCountInRect(int(testRect.X), int(testRect.Y), int(testRect.Width), int(testRect.Height))
	pd.Graphics.DrawText("  box: "+itoa(boxCount), 4, 12)

	// Draw instructions
	if showInstructions {
		pd.Graphics.DrawText("Press UP to add more snowflakes, DOWN to remove", 10, 60)
	}

	// Draw FPS
	pd.System.DrawFPS(2, 224)

	return 1
}

// Simple integer to string conversion
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var negative bool
	if n < 0 {
		negative = true
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}

func main() {}
