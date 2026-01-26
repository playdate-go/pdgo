package main

import (
	"github.com/playdate-go/pdgo"
)

var (
	pd       *pdgo.PlaydateAPI
	seq      *pdgo.SoundSequence
	lastStep int
)

func initGame() {
	// Create and load MIDI sequence
	seq = pd.Sound.Sequence.NewSequence()
	if err := pd.Sound.Sequence.LoadMIDIFile(seq, "bach.mid"); err != nil {
		pd.System.LogToConsole("Failed to load bach.mid")
		return
	}

	// Create instrument for playback
	inst := pd.Sound.Instrument.NewInstrument()
	pd.Sound.Instrument.SetVolume(inst, 0.2, 0.2)

	// Add instrument to default channel
	defaultChannel := pd.Sound.GetDefaultChannel()
	pd.Sound.Channel.AddInstrumentAsSource(defaultChannel, inst)

	// Get track count
	trackCount := pd.Sound.Sequence.GetTrackCount(seq)

	// Create synth with piano sample
	synth := pd.Sound.Synth.NewSynth()
	piano := pd.Sound.Sample.Load("piano")
	if piano != nil {
		pd.Sound.Synth.SetSample(synth, piano, 0, 0)
	} else {
		// Fallback to sine wave if no piano sample
		pd.Sound.Synth.SetWaveform(synth, pdgo.WaveformSine)
	}

	// Assign instrument to each track and add voices
	for i := 0; i < trackCount; i++ {
		track := pd.Sound.Sequence.GetTrackAtIndex(seq, uint(i))
		pd.Sound.Track.SetInstrument(track, inst)

		// Add voices based on track polyphony
		polyphony := pd.Sound.Track.GetPolyphony(track)
		for p := polyphony; p > 0; p-- {
			voice := pd.Sound.Synth.Copy(synth)
			pd.Sound.Instrument.AddVoice(inst, voice, 0, 127, 0)
		}
	}

	// Start playback
	pd.Sound.Sequence.Play(seq)

	lastStep = -1
}

func update() int {
	gfx := pd.Graphics

	step := pd.Sound.Sequence.GetCurrentStep(seq)

	if step > lastStep {
		// Scroll display left by 1 pixel
		displayBitmap := gfx.GetDisplayBufferBitmap()
		if displayBitmap != nil {
			gfx.DrawBitmap(displayBitmap, -1, 0, pdgo.BitmapUnflipped)
		}

		// Clear the rightmost column
		gfx.FillRect(399, 0, 1, 240, pdgo.NewColorFromSolid(pdgo.ColorWhite))

		// Draw notes for each track
		trackCount := pd.Sound.Sequence.GetTrackCount(seq)
		for i := 0; i < trackCount; i++ {
			track := pd.Sound.Sequence.GetTrackAtIndex(seq, uint(i))

			// Get notes that started since last step
			idx := pd.Sound.Track.GetIndexForStep(track, uint32(lastStep+1))
			for {
				noteStep, _, note, _, ok := pd.Sound.Track.GetNoteAtIndex(track, idx)
				if !ok || int(noteStep) > step {
					break
				}

				// Draw note as a small rectangle
				// Map MIDI note (typically 20-100) to screen Y position
				y := 240 - 3*(int(note)-20)
				if y >= 0 && y < 240 {
					gfx.FillRect(399, y, 1, 3, pdgo.NewColorFromSolid(pdgo.ColorBlack))
				}
				idx++
			}
		}
	}

	lastStep = step

	pd.System.DrawFPS(0, 0)
	return 1
}

func main() {}
