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
	pd.System.LogToConsole("Bach MIDI Player starting...")

	// --- Sequence ---
	seq = pd.Sound.Sequence.NewSequence()
	if err := pd.Sound.Sequence.LoadMIDIFile(seq, "bach.mid"); err != nil {
		pd.System.LogToConsole("Failed to load bach.mid")
		return
	}

	// --- Instrument ---
	inst := pd.Sound.Instrument.NewInstrument()
	pd.Sound.Instrument.SetVolume(inst, 0.2, 0.2)

	defaultChannel := pd.Sound.GetDefaultChannel()
	pd.Sound.Channel.AddInstrumentAsSource(defaultChannel, inst)

	// --- Base synth template ---
	baseSynth := pd.Sound.Synth.NewSynth()

	// Try sample first
	if piano := pd.Sound.Sample.Load("piano"); piano != nil {
		pd.Sound.Synth.SetSample(baseSynth, piano, 0, 0)
		pd.System.LogToConsole("Using piano sample")
	} else {
		pd.Sound.Synth.SetWaveform(baseSynth, pdgo.WaveformSine)
		pd.System.LogToConsole("Using sine fallback")
	}

	// --- Tracks & voices ---
	trackCount := pd.Sound.Sequence.GetTrackCount(seq)

	for i := 0; i < trackCount; i++ {
		track := pd.Sound.Sequence.GetTrackAtIndex(seq, uint(i))
		if track == nil {
			continue
		}

		pd.Sound.Track.SetInstrument(track, inst)

		poly := pd.Sound.Track.GetPolyphony(track)
		if poly < 1 {
			poly = 1
		}

		for v := 0; v < poly; v++ {
			s := pd.Sound.Synth.NewSynth()

			if piano := pd.Sound.Sample.Load("piano"); piano != nil {
				pd.Sound.Synth.SetSample(s, piano, 0, 0)
			} else {
				pd.Sound.Synth.SetWaveform(s, pdgo.WaveformSine)
			}

			pd.Sound.Instrument.AddVoice(inst, s, 0, 127, 0)
		}
	}

	// --- Play ---
	pd.Sound.Sequence.Play(seq)
	lastStep = -1
}

func update() int {
	gfx := pd.Graphics

	step := pd.Sound.Sequence.GetCurrentStep(seq)

	if step > lastStep {
		// Scroll left
		displayBitmap := gfx.GetDisplayBufferBitmap()
		if displayBitmap != nil {
			gfx.DrawBitmap(displayBitmap, -1, 0, pdgo.BitmapUnflipped)
		}

		// Clear right column
		gfx.FillRect(399, 0, 1, 240, pdgo.SolidWhite)

		// Draw notes
		trackCount := pd.Sound.Sequence.GetTrackCount(seq)
		for i := 0; i < trackCount; i++ {
			track := pd.Sound.Sequence.GetTrackAtIndex(seq, uint(i))
			if track == nil {
				continue
			}

			idx := pd.Sound.Track.GetIndexForStep(track, uint32(lastStep+1))
			for {
				noteStep, _, note, _, ok := pd.Sound.Track.GetNoteAtIndex(track, idx)
				if !ok || int(noteStep) > step {
					break
				}

				y := 240 - 3*(int(note)-20)
				if y >= 0 && y < 240 {
					gfx.FillRect(399, y, 1, 3, pdgo.SolidBlack)
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
