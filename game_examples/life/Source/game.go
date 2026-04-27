package main

import "github.com/playdate-go/pdgo"

// Game represents the Game of Life simulation state
type Game struct {
	pd          *pdgo.PlaydateAPI
	initialized bool
	rngState    uint32
}

// NewGame creates a new Game instance
func NewGame(pd *pdgo.PlaydateAPI) *Game {
	return &Game{
		pd:          pd,
		initialized: false,
		rngState:    12345,
	}
}

// Init initializes the game
func (g *Game) Init() {
	// Run as fast as possible (0 = unlimited)
	g.pd.Display.SetRefreshRate(0)

	// Seed RNG with current time
	seconds, ms := g.pd.System.GetSecondsSinceEpoch()
	g.rngState = uint32(seconds) ^ uint32(ms<<16)
	if g.rngState == 0 {
		g.rngState = 12345
	}
}

// random returns a pseudo-random uint32 using XorShift
func (g *Game) random() uint32 {
	x := g.rngState
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	g.rngState = x
	return x
}

// Update is the main game loop - called every frame
func (g *Game) Update() int {
	// Initialize on first frame when framebuffer is ready
	if !g.initialized {
		g.Randomize()
		g.initialized = true
	}

	// Check for button press to randomize
	_, pushed, _ := g.pd.System.GetButtonState()
	if pushed&pdgo.ButtonA != 0 {
		g.Randomize()
	}

	// Get frame buffers
	nextFrame := g.pd.Graphics.GetFrame()    // Working buffer
	frame := g.pd.Graphics.GetDisplayFrame() // Buffer currently on screen

	if frame == nil || nextFrame == nil {
		return 1
	}

	// Simulate one generation
	g.Step(frame, nextFrame)

	// We modified framebuffer directly, tell the system about it
	g.pd.Graphics.MarkUpdatedRows(0, pdgo.LCDRows-1)

	return 1
}

// Randomize fills the frame buffer with random data
func (g *Game) Randomize() {
	frame := g.pd.Graphics.GetDisplayFrame()
	if frame == nil {
		return
	}

	for y := 0; y < pdgo.LCDRows; y++ {
		rowStart := y * pdgo.LCDRowSize
		for x := 0; x < pdgo.LCDColumns/8; x++ {
			frame[rowStart+x] = byte(g.random())
		}
	}
}
