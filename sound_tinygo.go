//go:build tinygo

// TinyGo implementation of Sound API

package pdgo

// SoundSource represents a sound source
type SoundSource struct {
	ptr uintptr
}

// FilePlayer plays audio files
type FilePlayer struct {
	ptr uintptr
}

// SamplePlayer plays audio samples
type SamplePlayer struct {
	ptr uintptr
}

// AudioSample represents an audio sample
type AudioSample struct {
	ptr uintptr
}

// PDSynth represents a synthesizer
type PDSynth struct {
	ptr uintptr
}

// SoundChannel represents a sound channel
type SoundChannel struct {
	ptr uintptr
}

// SoundEffect represents a sound effect
type SoundEffect struct {
	ptr uintptr
}

// SoundSequence represents a sound sequence
type SoundSequence struct {
	ptr uintptr
}

// SequenceTrack represents a track in a sequence
type SequenceTrack struct {
	ptr uintptr
}

// PDSynthInstrument represents a synth instrument
type PDSynthInstrument struct {
	ptr uintptr
}

// ControlSignal represents a control signal
type ControlSignal struct {
	ptr uintptr
}

// SoundWaveform represents synth waveforms
type SoundWaveform int32

const (
	WaveformSquare    SoundWaveform = 0
	WaveformTriangle  SoundWaveform = 1
	WaveformSine      SoundWaveform = 2
	WaveformNoise     SoundWaveform = 3
	WaveformSawtooth  SoundWaveform = 4
	WaveformPOPhase   SoundWaveform = 5
	WaveformPODigital SoundWaveform = 6
	WaveformPOVosim   SoundWaveform = 7
)

// MIDINote represents a MIDI note number
type MIDINote float32

// Sound provides access to sound functions
type Sound struct {
	Channel    *ChannelAPI
	Sample     *SampleAPI
	Synth      *SynthAPI
	Sequence   *SequenceAPI
	Track      *TrackAPI
	Instrument *InstrumentAPI
}

func newSound() *Sound {
	return &Sound{
		Channel:    &ChannelAPI{},
		Sample:     &SampleAPI{},
		Synth:      &SynthAPI{},
		Sequence:   &SequenceAPI{},
		Track:      &TrackAPI{},
		Instrument: &InstrumentAPI{},
	}
}

// GetDefaultChannel returns the default sound channel
func (s *Sound) GetDefaultChannel() *SoundChannel {
	if bridgeSoundGetDefaultChannel != nil {
		ptr := bridgeSoundGetDefaultChannel()
		if ptr != 0 {
			return &SoundChannel{ptr: ptr}
		}
	}
	return nil
}

// SynthAPI wraps synth functions
type SynthAPI struct{}

// NewSynth creates a new synth
func (sy *SynthAPI) NewSynth() *PDSynth {
	if bridgeSoundSynthNew != nil {
		ptr := bridgeSoundSynthNew()
		if ptr != 0 {
			return &PDSynth{ptr: ptr}
		}
	}
	return nil
}

// FreeSynth frees a synth
func (sy *SynthAPI) FreeSynth(synth *PDSynth) {
	if bridgeSoundSynthFree != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthFree(synth.ptr)
		synth.ptr = 0
	}
}

// SetWaveform sets the waveform
func (sy *SynthAPI) SetWaveform(synth *PDSynth, wave SoundWaveform) {
	if bridgeSoundSynthSetWaveform != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthSetWaveform(synth.ptr, int32(wave))
	}
}

// SetAttackTime sets the attack time
func (sy *SynthAPI) SetAttackTime(synth *PDSynth, attack float32) {
	if bridgeSoundSynthSetAttack != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthSetAttack(synth.ptr, attack)
	}
}

// SetDecayTime sets the decay time
func (sy *SynthAPI) SetDecayTime(synth *PDSynth, decay float32) {
	if bridgeSoundSynthSetDecay != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthSetDecay(synth.ptr, decay)
	}
}

// SetSustainLevel sets the sustain level
func (sy *SynthAPI) SetSustainLevel(synth *PDSynth, sustain float32) {
	if bridgeSoundSynthSetSustain != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthSetSustain(synth.ptr, sustain)
	}
}

// SetReleaseTime sets the release time
func (sy *SynthAPI) SetReleaseTime(synth *PDSynth, release float32) {
	if bridgeSoundSynthSetRelease != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthSetRelease(synth.ptr, release)
	}
}

// SetTranspose sets the transpose
func (sy *SynthAPI) SetTranspose(synth *PDSynth, halfSteps float32) {
	if bridgeSoundSynthSetTranspose != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthSetTranspose(synth.ptr, halfSteps)
	}
}

// PlayNote plays a note
func (sy *SynthAPI) PlayNote(synth *PDSynth, freq, vel, length float32, when uint32) {
	if bridgeSoundSynthPlayNote != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthPlayNote(synth.ptr, freq, vel, length, when)
	}
}

// PlayMIDINote plays a MIDI note
func (sy *SynthAPI) PlayMIDINote(synth *PDSynth, note MIDINote, vel, length float32, when uint32) {
	if bridgeSoundSynthPlayMIDINote != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthPlayMIDINote(synth.ptr, float32(note), vel, length, when)
	}
}

// NoteOff releases a note
func (sy *SynthAPI) NoteOff(synth *PDSynth, when uint32) {
	if bridgeSoundSynthNoteOff != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthNoteOff(synth.ptr, when)
	}
}

// Stop stops the synth
func (sy *SynthAPI) Stop(synth *PDSynth) {
	if bridgeSoundSynthStop != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthStop(synth.ptr)
	}
}

// SetVolume sets the volume
func (sy *SynthAPI) SetVolume(synth *PDSynth, left, right float32) {
	if bridgeSoundSynthSetVolume != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthSetVolume(synth.ptr, left, right)
	}
}

// GetVolume returns the volume
func (sy *SynthAPI) GetVolume(synth *PDSynth) (left, right float32) {
	if bridgeSoundSynthGetVolume != nil && synth != nil && synth.ptr != 0 {
		bridgeSoundSynthGetVolume(synth.ptr, &left, &right)
	}
	return
}

// IsPlaying returns whether the synth is playing
func (sy *SynthAPI) IsPlaying(synth *PDSynth) bool {
	if bridgeSoundSynthIsPlaying != nil && synth != nil && synth.ptr != 0 {
		return bridgeSoundSynthIsPlaying(synth.ptr) != 0
	}
	return false
}

// SetSample sets the sample for the synth
func (sy *SynthAPI) SetSample(synth *PDSynth, sample *AudioSample, sustainStart, sustainEnd uint32) {
	if bridgeSoundSynthSetSample != nil && synth != nil && synth.ptr != 0 {
		var samplePtr uintptr
		if sample != nil {
			samplePtr = sample.ptr
		}
		bridgeSoundSynthSetSample(synth.ptr, samplePtr, sustainStart, sustainEnd)
	}
}

// Copy copies a synth
func (sy *SynthAPI) Copy(synth *PDSynth) *PDSynth {
	if bridgeSoundSynthCopy != nil && synth != nil && synth.ptr != 0 {
		ptr := bridgeSoundSynthCopy(synth.ptr)
		if ptr != 0 {
			return &PDSynth{ptr: ptr}
		}
	}
	return nil
}

// ============== ChannelAPI ==============

// ChannelAPI wraps channel functions
type ChannelAPI struct{}

// AddInstrumentAsSource adds an instrument as a source to a channel
func (c *ChannelAPI) AddInstrumentAsSource(channel *SoundChannel, inst *PDSynthInstrument) bool {
	if bridgeSoundChannelAddInstrument != nil && channel != nil && channel.ptr != 0 && inst != nil && inst.ptr != 0 {
		return bridgeSoundChannelAddInstrument(channel.ptr, inst.ptr) != 0
	}
	return false
}

// ============== SampleAPI ==============

// SampleAPI wraps sample functions
type SampleAPI struct{}

// Load loads an audio sample from file
func (sa *SampleAPI) Load(path string) *AudioSample {
	if bridgeSoundLoadSample != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		ptr := bridgeSoundLoadSample(&buf[0])
		if ptr != 0 {
			return &AudioSample{ptr: ptr}
		}
	}
	return nil
}

// ============== SequenceAPI ==============

// SequenceAPI wraps sequence functions
type SequenceAPI struct{}

// NewSequence creates a new sequence
func (s *SequenceAPI) NewSequence() *SoundSequence {
	if bridgeSoundSequenceNew != nil {
		ptr := bridgeSoundSequenceNew()
		if ptr != 0 {
			return &SoundSequence{ptr: ptr}
		}
	}
	return nil
}

// LoadMIDIFile loads a MIDI file into the sequence
func (s *SequenceAPI) LoadMIDIFile(seq *SoundSequence, path string) error {
	if bridgeSoundSequenceLoadMIDI != nil && seq != nil && seq.ptr != 0 {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		if bridgeSoundSequenceLoadMIDI(seq.ptr, &buf[0]) != 0 {
			return nil
		}
	}
	return &soundError{op: "loadMIDI", path: path}
}

// GetTrackCount returns the track count
func (s *SequenceAPI) GetTrackCount(seq *SoundSequence) int {
	if bridgeSoundSequenceGetTrackCount != nil && seq != nil && seq.ptr != 0 {
		return int(bridgeSoundSequenceGetTrackCount(seq.ptr))
	}
	return 0
}

// GetTrackAtIndex returns a track at an index
func (s *SequenceAPI) GetTrackAtIndex(seq *SoundSequence, index uint) *SequenceTrack {
	if bridgeSoundSequenceGetTrackAtIndex != nil && seq != nil && seq.ptr != 0 {
		ptr := bridgeSoundSequenceGetTrackAtIndex(seq.ptr, uint32(index))
		if ptr != 0 {
			return &SequenceTrack{ptr: ptr}
		}
	}
	return nil
}

// Play plays the sequence
func (s *SequenceAPI) Play(seq *SoundSequence) {
	if bridgeSoundSequencePlay != nil && seq != nil && seq.ptr != 0 {
		bridgeSoundSequencePlay(seq.ptr)
	}
}

// Stop stops the sequence
func (s *SequenceAPI) Stop(seq *SoundSequence) {
	if bridgeSoundSequenceStop != nil && seq != nil && seq.ptr != 0 {
		bridgeSoundSequenceStop(seq.ptr)
	}
}

// GetCurrentStep returns the current step in the sequence
func (s *SequenceAPI) GetCurrentStep(seq *SoundSequence) int {
	if bridgeSoundSequenceGetCurrentStep != nil && seq != nil && seq.ptr != 0 {
		return int(bridgeSoundSequenceGetCurrentStep(seq.ptr))
	}
	return 0
}

// ============== TrackAPI ==============

// TrackAPI wraps track functions
type TrackAPI struct{}

// SetInstrument sets the instrument for a track
func (t *TrackAPI) SetInstrument(track *SequenceTrack, inst *PDSynthInstrument) {
	if bridgeSoundTrackSetInstrument != nil && track != nil && track.ptr != 0 {
		var instPtr uintptr
		if inst != nil {
			instPtr = inst.ptr
		}
		bridgeSoundTrackSetInstrument(track.ptr, instPtr)
	}
}

// GetPolyphony returns the maximum polyphony for the track
func (t *TrackAPI) GetPolyphony(track *SequenceTrack) int {
	if bridgeSoundTrackGetPolyphony != nil && track != nil && track.ptr != 0 {
		return int(bridgeSoundTrackGetPolyphony(track.ptr))
	}
	return 0
}

// GetIndexForStep returns the internal index for the step
func (t *TrackAPI) GetIndexForStep(track *SequenceTrack, step uint32) int {
	if bridgeSoundTrackGetIndexForStep != nil && track != nil && track.ptr != 0 {
		return int(bridgeSoundTrackGetIndexForStep(track.ptr, step))
	}
	return 0
}

// GetNoteAtIndex returns note information at the given index
func (t *TrackAPI) GetNoteAtIndex(track *SequenceTrack, index int) (step, length uint32, note MIDINote, velocity float32, ok bool) {
	if bridgeSoundTrackGetNoteAtIndex != nil && track != nil && track.ptr != 0 {
		var s, l uint32
		var n float32
		var v float32
		result := bridgeSoundTrackGetNoteAtIndex(track.ptr, int32(index), &s, &l, &n, &v)
		if result != 0 {
			return s, l, MIDINote(n), v, true
		}
	}
	return 0, 0, 0, 0, false
}

// ============== InstrumentAPI ==============

// InstrumentAPI wraps instrument functions
type InstrumentAPI struct{}

// NewInstrument creates a new instrument
func (i *InstrumentAPI) NewInstrument() *PDSynthInstrument {
	if bridgeSoundInstrumentNew != nil {
		ptr := bridgeSoundInstrumentNew()
		if ptr != 0 {
			return &PDSynthInstrument{ptr: ptr}
		}
	}
	return nil
}

// SetVolume sets the instrument volume
func (i *InstrumentAPI) SetVolume(inst *PDSynthInstrument, left, right float32) {
	if bridgeSoundInstrumentSetVolume != nil && inst != nil && inst.ptr != 0 {
		bridgeSoundInstrumentSetVolume(inst.ptr, left, right)
	}
}

// AddVoice adds a voice to the instrument
func (i *InstrumentAPI) AddVoice(inst *PDSynthInstrument, synth *PDSynth, rangeStart, rangeEnd MIDINote, transpose float32) bool {
	if bridgeSoundInstrumentAddVoice != nil && inst != nil && inst.ptr != 0 && synth != nil && synth.ptr != 0 {
		return bridgeSoundInstrumentAddVoice(inst.ptr, synth.ptr, float32(rangeStart), float32(rangeEnd), transpose) != 0
	}
	return false
}

// ============== FilePlayer ==============

// NewFilePlayer creates a new file player
func (s *Sound) NewFilePlayer() *FilePlayer {
	if bridgeSoundNewFilePlayer != nil {
		ptr := bridgeSoundNewFilePlayer()
		if ptr != 0 {
			return &FilePlayer{ptr: ptr}
		}
	}
	return nil
}

// FreeFilePlayer frees a file player
func (s *Sound) FreeFilePlayer(player *FilePlayer) {
	if bridgeSoundFreeFilePlayer != nil && player != nil && player.ptr != 0 {
		bridgeSoundFreeFilePlayer(player.ptr)
		player.ptr = 0
	}
}

// LoadIntoFilePlayer loads audio into file player
func (s *Sound) LoadIntoFilePlayer(player *FilePlayer, path string) error {
	if bridgeSoundLoadIntoFilePlayer != nil && player != nil && player.ptr != 0 {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		if bridgeSoundLoadIntoFilePlayer(player.ptr, &buf[0]) != 0 {
			return nil
		}
	}
	return &soundError{op: "load", path: path}
}

// PlayFilePlayer plays audio
func (s *Sound) PlayFilePlayer(player *FilePlayer, repeat int) {
	if bridgeSoundPlayFilePlayer != nil && player != nil && player.ptr != 0 {
		bridgeSoundPlayFilePlayer(player.ptr, int32(repeat))
	}
}

// StopFilePlayer stops playback
func (s *Sound) StopFilePlayer(player *FilePlayer) {
	if bridgeSoundStopFilePlayer != nil && player != nil && player.ptr != 0 {
		bridgeSoundStopFilePlayer(player.ptr)
	}
}

// PauseFilePlayer pauses playback
func (s *Sound) PauseFilePlayer(player *FilePlayer) {
	if bridgeSoundPauseFilePlayer != nil && player != nil && player.ptr != 0 {
		bridgeSoundPauseFilePlayer(player.ptr)
	}
}

// IsFilePlayerPlaying returns true if playing
func (s *Sound) IsFilePlayerPlaying(player *FilePlayer) bool {
	if bridgeSoundIsFilePlayerPlaying != nil && player != nil && player.ptr != 0 {
		return bridgeSoundIsFilePlayerPlaying(player.ptr) != 0
	}
	return false
}

// SetFilePlayerVolume sets volume (0.0 - 1.0)
func (s *Sound) SetFilePlayerVolume(player *FilePlayer, left, right float32) {
	if bridgeSoundSetFilePlayerVolume != nil && player != nil && player.ptr != 0 {
		bridgeSoundSetFilePlayerVolume(player.ptr, left, right)
	}
}

// GetFilePlayerVolume gets volume
func (s *Sound) GetFilePlayerVolume(player *FilePlayer) (left, right float32) {
	if bridgeSoundGetFilePlayerVolume != nil && player != nil && player.ptr != 0 {
		bridgeSoundGetFilePlayerVolume(player.ptr, &left, &right)
	}
	return
}

// GetFilePlayerLength returns length in seconds
func (s *Sound) GetFilePlayerLength(player *FilePlayer) float32 {
	if bridgeSoundGetFilePlayerLength != nil && player != nil && player.ptr != 0 {
		return bridgeSoundGetFilePlayerLength(player.ptr)
	}
	return 0
}

// SetFilePlayerOffset sets playback offset
func (s *Sound) SetFilePlayerOffset(player *FilePlayer, offset float32) {
	if bridgeSoundSetFilePlayerOffset != nil && player != nil && player.ptr != 0 {
		bridgeSoundSetFilePlayerOffset(player.ptr, offset)
	}
}

// GetFilePlayerOffset gets playback offset
func (s *Sound) GetFilePlayerOffset(player *FilePlayer) float32 {
	if bridgeSoundGetFilePlayerOffset != nil && player != nil && player.ptr != 0 {
		return bridgeSoundGetFilePlayerOffset(player.ptr)
	}
	return 0
}

// SetFilePlayerRate sets playback rate
func (s *Sound) SetFilePlayerRate(player *FilePlayer, rate float32) {
	if bridgeSoundSetFilePlayerRate != nil && player != nil && player.ptr != 0 {
		bridgeSoundSetFilePlayerRate(player.ptr, rate)
	}
}

// ============== SamplePlayer ==============

// NewSamplePlayer creates a new sample player
func (s *Sound) NewSamplePlayer() *SamplePlayer {
	if bridgeSoundNewSamplePlayer != nil {
		ptr := bridgeSoundNewSamplePlayer()
		if ptr != 0 {
			return &SamplePlayer{ptr: ptr}
		}
	}
	return nil
}

// FreeSamplePlayer frees a sample player
func (s *Sound) FreeSamplePlayer(player *SamplePlayer) {
	if bridgeSoundFreeSamplePlayer != nil && player != nil && player.ptr != 0 {
		bridgeSoundFreeSamplePlayer(player.ptr)
		player.ptr = 0
	}
}

// SetSamplePlayerSample sets the sample to play
func (s *Sound) SetSamplePlayerSample(player *SamplePlayer, sample *AudioSample) {
	if bridgeSoundSetSamplePlayerSample != nil && player != nil && player.ptr != 0 && sample != nil {
		bridgeSoundSetSamplePlayerSample(player.ptr, sample.ptr)
	}
}

// PlaySamplePlayer plays the sample
func (s *Sound) PlaySamplePlayer(player *SamplePlayer, repeat int, rate float32) {
	if bridgeSoundPlaySamplePlayer != nil && player != nil && player.ptr != 0 {
		bridgeSoundPlaySamplePlayer(player.ptr, int32(repeat), rate)
	}
}

// StopSamplePlayer stops playback
func (s *Sound) StopSamplePlayer(player *SamplePlayer) {
	if bridgeSoundStopSamplePlayer != nil && player != nil && player.ptr != 0 {
		bridgeSoundStopSamplePlayer(player.ptr)
	}
}

// IsSamplePlayerPlaying returns true if playing
func (s *Sound) IsSamplePlayerPlaying(player *SamplePlayer) bool {
	if bridgeSoundIsSamplePlayerPlaying != nil && player != nil && player.ptr != 0 {
		return bridgeSoundIsSamplePlayerPlaying(player.ptr) != 0
	}
	return false
}

// SetSamplePlayerVolume sets volume
func (s *Sound) SetSamplePlayerVolume(player *SamplePlayer, left, right float32) {
	if bridgeSoundSetSamplePlayerVolume != nil && player != nil && player.ptr != 0 {
		bridgeSoundSetSamplePlayerVolume(player.ptr, left, right)
	}
}

// ============== AudioSample ==============

// NewAudioSample creates a new audio sample
func (s *Sound) NewAudioSample(length int) *AudioSample {
	if bridgeSoundNewSample != nil {
		ptr := bridgeSoundNewSample(int32(length))
		if ptr != 0 {
			return &AudioSample{ptr: ptr}
		}
	}
	return nil
}

// LoadAudioSample loads an audio sample from file
func (s *Sound) LoadAudioSample(path string) *AudioSample {
	if bridgeSoundLoadSample != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		ptr := bridgeSoundLoadSample(&buf[0])
		if ptr != 0 {
			return &AudioSample{ptr: ptr}
		}
	}
	return nil
}

// FreeAudioSample frees an audio sample
func (s *Sound) FreeAudioSample(sample *AudioSample) {
	if bridgeSoundFreeSample != nil && sample != nil && sample.ptr != 0 {
		bridgeSoundFreeSample(sample.ptr)
		sample.ptr = 0
	}
}

// ============== Global Sound ==============

// SetHeadphoneCallback sets headphone callback (stub)
func (s *Sound) SetHeadphoneCallback(callback func(headphone, mic int)) {
	// Not implemented for TinyGo
}

// GetHeadphoneState returns headphone state
func (s *Sound) GetHeadphoneState() (headphone, mic bool) {
	if bridgeSoundGetHeadphoneState != nil {
		var h, m int32
		bridgeSoundGetHeadphoneState(&h, &m)
		return h != 0, m != 0
	}
	return false, false
}

// SetOutputsActive sets which outputs are active
func (s *Sound) SetOutputsActive(headphone, speaker bool) {
	if bridgeSoundSetOutputsActive != nil {
		var h, sp int32
		if headphone {
			h = 1
		}
		if speaker {
			sp = 1
		}
		bridgeSoundSetOutputsActive(h, sp)
	}
}

type soundError struct {
	op   string
	path string
}

func (e *soundError) Error() string {
	if e.path != "" {
		return "sound " + e.op + " failed: " + e.path
	}
	return "sound " + e.op + " failed"
}
