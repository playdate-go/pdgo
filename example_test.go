// Example usage of pdgo package
// This file shows how to use the Go bindings for Playdate SDK
//
// Note: This is documentation-only code. Actual Playdate games need to be
// cross-compiled for the ARM Cortex-M7 processor.

package pdgo_test

import (
	"pdgo"
	"unsafe"
)

// Example shows basic usage of the pdgo package
func Example() {
	// This is a conceptual example showing how to structure a Playdate game in Go

	// The eventHandler would be called from C like this:
	// int eventHandler(PlaydateAPI* pd, PDSystemEvent event, uint32_t arg)

	// In your actual game, you would export an eventHandler function:
	/*
		//export eventHandler
		func eventHandler(pdAPI unsafe.Pointer, event C.PDSystemEvent, arg C.uint32_t) C.int {
			switch pdgo.PDSystemEvent(event) {
			case pdgo.EventInit:
				// Initialize the API
				pd := pdgo.Init(pdAPI)

				// Set up the update callback
				pd.System.SetUpdateCallback(update)

				// Log to console
				pd.System.LogToConsole("Hello from Go!")

			case pdgo.EventTerminate:
				// Cleanup
			}
			return 0
		}
	*/
}

// ExampleGraphics shows how to use graphics functions
func ExampleGraphics() {
	// Assuming pd is initialized
	pd := pdgo.GetAPI()
	if pd == nil {
		return
	}

	gfx := pd.Graphics

	// Clear the screen
	gfx.Clear(pdgo.NewColorFromSolid(pdgo.ColorWhite))

	// Draw a rectangle
	gfx.FillRect(10, 10, 100, 50, pdgo.NewColorFromSolid(pdgo.ColorBlack))

	// Draw text
	gfx.DrawText("Hello Playdate!", 50, 100)

	// Draw a line
	gfx.DrawLine(0, 0, 400, 240, 2, pdgo.NewColorFromSolid(pdgo.ColorBlack))

	// Draw an ellipse
	gfx.DrawEllipse(200, 120, 80, 60, 2, 0, 360, pdgo.NewColorFromSolid(pdgo.ColorBlack))
}

// ExampleSprites shows how to use sprites
func ExampleSprites() {
	pd := pdgo.GetAPI()
	if pd == nil {
		return
	}

	spr := pd.Sprite
	gfx := pd.Graphics

	// Create a new sprite
	sprite := spr.NewSprite()
	if sprite == nil {
		return
	}
	defer spr.FreeSprite(sprite)

	// Load an image
	bitmap, err := gfx.LoadBitmap("images/player.png")
	if err != nil {
		pd.System.LogToConsole("Failed to load bitmap")
		return
	}
	defer gfx.FreeBitmap(bitmap)

	// Set the sprite's image
	spr.SetImage(sprite, bitmap, pdgo.BitmapUnflipped)

	// Set position
	spr.MoveTo(sprite, 200, 120)

	// Set collision rect
	spr.SetCollideRect(sprite, pdgo.PDRect{X: 0, Y: 0, Width: 32, Height: 32})

	// Set collision response
	spr.SetCollisionResponseFunction(sprite, func(s, other *pdgo.LCDSprite) pdgo.SpriteCollisionResponseType {
		return pdgo.CollisionTypeSlide
	})

	// Add to display list
	spr.AddSprite(sprite)

	// In update function:
	// spr.UpdateAndDrawSprites()
}

// ExampleSound shows how to use sound
func ExampleSound() {
	pd := pdgo.GetAPI()
	if pd == nil {
		return
	}

	snd := pd.Sound

	// File player for music
	player := snd.FilePlayer.NewPlayer()
	if player == nil {
		return
	}
	defer snd.FilePlayer.FreePlayer(player)

	err := snd.FilePlayer.LoadIntoPlayer(player, "audio/music.pda")
	if err != nil {
		pd.System.LogToConsole("Failed to load music")
		return
	}

	snd.FilePlayer.SetVolume(player, 0.8, 0.8)
	snd.FilePlayer.Play(player, 0) // 0 = loop forever

	// Sample player for sound effects
	sample := snd.Sample.Load("audio/jump.pda")
	if sample == nil {
		return
	}
	defer snd.Sample.FreeSample(sample)

	samplePlayer := snd.SamplePlayer.NewPlayer()
	if samplePlayer == nil {
		return
	}
	defer snd.SamplePlayer.FreePlayer(samplePlayer)

	snd.SamplePlayer.SetSample(samplePlayer, sample)
	snd.SamplePlayer.Play(samplePlayer, 1, 1.0) // play once at normal rate

	// Synth for procedural audio
	synth := snd.Synth.NewSynth()
	if synth == nil {
		return
	}
	defer snd.Synth.FreeSynth(synth)

	snd.Synth.SetWaveform(synth, pdgo.WaveformSquare)
	snd.Synth.SetAttackTime(synth, 0.01)
	snd.Synth.SetDecayTime(synth, 0.1)
	snd.Synth.SetSustainLevel(synth, 0.5)
	snd.Synth.SetReleaseTime(synth, 0.2)

	// Play a note
	snd.Synth.PlayMIDINote(synth, pdgo.MIDINote(pdgo.NoteC4), 1.0, 0.5, 0)
}

// ExampleInput shows how to handle input
func ExampleInput() {
	pd := pdgo.GetAPI()
	if pd == nil {
		return
	}

	sys := pd.System

	// Enable accelerometer
	sys.SetPeripheralsEnabled(pdgo.PeripheralAccelerometer)

	// In update callback:
	updateCallback := func() int {
		// Get button state
		current, pushed, released := sys.GetButtonState()

		if pushed&pdgo.ButtonA != 0 {
			sys.LogToConsole("A button pressed!")
		}

		if released&pdgo.ButtonB != 0 {
			sys.LogToConsole("B button released!")
		}

		if current&pdgo.ButtonUp != 0 {
			// Moving up while held
		}

		// Get crank
		crankChange := sys.GetCrankChange()
		if crankChange != 0 {
			// Crank moved
		}

		if sys.IsCrankDocked() {
			// Crank is docked
		}

		// Get accelerometer
		ax, ay, az := sys.GetAccelerometer()
		_ = ax
		_ = ay
		_ = az

		return 1 // continue running
	}

	sys.SetUpdateCallback(updateCallback)
}

// ExampleFile shows how to use file operations
func ExampleFile() {
	pd := pdgo.GetAPI()
	if pd == nil {
		return
	}

	file := pd.File

	// List files
	files, err := file.ListFiles("data/", false)
	if err != nil {
		pd.System.LogToConsole("Failed to list files")
		return
	}

	for _, f := range files {
		pd.System.LogToConsole("File: " + f)
	}

	// Read a file
	f, err := file.Open("data/save.dat", pdgo.FileRead)
	if err != nil {
		pd.System.LogToConsole("Failed to open file")
		return
	}
	defer f.Close()

	data, err := f.ReadAll()
	if err != nil {
		pd.System.LogToConsole("Failed to read file")
		return
	}
	_ = data

	// Write a file
	wf, err := file.Open("data/save.dat", pdgo.FileWrite)
	if err != nil {
		pd.System.LogToConsole("Failed to create file")
		return
	}
	defer wf.Close()

	wf.WriteString("Hello, Playdate!")
	wf.Flush()
}

// update is an example update callback
func update() int {
	pd := pdgo.GetAPI()
	if pd == nil {
		return 0
	}

	// Clear screen
	pd.Graphics.Clear(pdgo.NewColorFromSolid(pdgo.ColorWhite))

	// Draw FPS
	pd.System.DrawFPS(0, 0)

	// Update and draw sprites
	pd.Sprite.UpdateAndDrawSprites()

	return 1 // Continue running
}

// Helper to make examples compile
var _ = unsafe.Pointer(nil)
