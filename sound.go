// pdgo Sound API - unified CGO implementation

package pdgo

/*
#include <stdint.h>

// Sound - Channel
void* pd_sound_getDefaultChannel(void);
int pd_sound_channel_addSource(void* channel, void* source);
int pd_sound_channel_addInstrument(void* channel, void* inst);

// Sound - Synth
void* pd_sound_synth_new(void);
void pd_sound_synth_free(void* synth);
void pd_sound_synth_setWaveform(void* synth, int waveform);
void pd_sound_synth_setAttackTime(void* synth, float attack);
void pd_sound_synth_setDecayTime(void* synth, float decay);
void pd_sound_synth_setSustainLevel(void* synth, float sustain);
void pd_sound_synth_setReleaseTime(void* synth, float release);
void pd_sound_synth_setTranspose(void* synth, float halfSteps);
void pd_sound_synth_playNote(void* synth, float freq, float vel, float len, uint32_t when);
void pd_sound_synth_playMIDINote(void* synth, float note, float vel, float len, uint32_t when);
void pd_sound_synth_noteOff(void* synth, uint32_t when);
void pd_sound_synth_stop(void* synth);
void pd_sound_synth_setVolume(void* synth, float left, float right);
void pd_sound_synth_getVolume(void* synth, float* left, float* right);
int pd_sound_synth_isPlaying(void* synth);
void pd_sound_synth_setSample(void* synth, void* sample, uint32_t sustainStart, uint32_t sustainEnd);
void* pd_sound_synth_copy(void* synth);

// Sound - Sample
void* pd_sound_sample_new(int len);
void* pd_sound_sample_load(const char* path);
void pd_sound_sample_free(void* sample);

// Sound - FilePlayer
void* pd_sound_fileplayer_new(void);
void pd_sound_fileplayer_free(void* player);
int pd_sound_fileplayer_load(void* player, const char* path);
int pd_sound_fileplayer_play(void* player, int repeat);
void pd_sound_fileplayer_stop(void* player);
void pd_sound_fileplayer_pause(void* player);
int pd_sound_fileplayer_isPlaying(void* player);
void pd_sound_fileplayer_setVolume(void* player, float left, float right);
void pd_sound_fileplayer_getVolume(void* player, float* left, float* right);
float pd_sound_fileplayer_getLength(void* player);
void pd_sound_fileplayer_setOffset(void* player, float offset);
float pd_sound_fileplayer_getOffset(void* player);
void pd_sound_fileplayer_setRate(void* player, float rate);

// Sound - SamplePlayer
void* pd_sound_sampleplayer_new(void);
void pd_sound_sampleplayer_free(void* player);
void pd_sound_sampleplayer_setSample(void* player, void* sample);
int pd_sound_sampleplayer_play(void* player, int repeat, float rate);
void pd_sound_sampleplayer_stop(void* player);
int pd_sound_sampleplayer_isPlaying(void* player);
void pd_sound_sampleplayer_setVolume(void* player, float left, float right);

// Sound - Sequence
void* pd_sound_sequence_new(void);
void pd_sound_sequence_free(void* seq);
int pd_sound_sequence_loadMIDI(void* seq, const char* path);
int pd_sound_sequence_getTrackCount(void* seq);
void* pd_sound_sequence_getTrackAtIndex(void* seq, int idx);
void pd_sound_sequence_play(void* seq, void* finishCallback, void* userdata);
void pd_sound_sequence_stop(void* seq);
int pd_sound_sequence_getCurrentStep(void* seq, int* timeOffset);

// Sound - Track
void pd_sound_track_setInstrument(void* track, void* inst);
int pd_sound_track_getPolyphony(void* track);
int pd_sound_track_getIndexForStep(void* track, uint32_t step);
int pd_sound_track_getNoteAtIndex(void* track, int index, uint32_t* step, uint32_t* len, float* note, float* vel);

// Sound - Instrument
void* pd_sound_instrument_new(void);
void pd_sound_instrument_free(void* inst);
void pd_sound_instrument_setVolume(void* inst, float left, float right);
int pd_sound_instrument_addVoice(void* inst, void* synth, float rangeStart, float rangeEnd, float transpose);

// Sound - Global
void pd_sound_getHeadphoneState(int* headphone, int* mic);
void pd_sound_setOutputsActive(int headphone, int speaker);
*/
import "C"
import "unsafe"

// SoundSource represents a sound source
type SoundSource struct {
	ptr unsafe.Pointer
}

// FilePlayer plays audio files
type FilePlayer struct {
	ptr unsafe.Pointer
}

// SamplePlayer plays audio samples
type SamplePlayer struct {
	ptr unsafe.Pointer
}

// AudioSample represents an audio sample
type AudioSample struct {
	ptr unsafe.Pointer
}

// PDSynth represents a synthesizer
type PDSynth struct {
	ptr unsafe.Pointer
}

// SoundChannel represents a sound channel
type SoundChannel struct {
	ptr unsafe.Pointer
}

// SoundEffect represents a sound effect
type SoundEffect struct {
	ptr unsafe.Pointer
}

// SoundSequence represents a sound sequence
type SoundSequence struct {
	ptr unsafe.Pointer
}

// SequenceTrack represents a track in a sequence
type SequenceTrack struct {
	ptr unsafe.Pointer
}

// PDSynthInstrument represents a synth instrument
type PDSynthInstrument struct {
	ptr unsafe.Pointer
}

// ControlSignal represents a control signal
type ControlSignal struct {
	ptr unsafe.Pointer
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
	ptr := C.pd_sound_getDefaultChannel()
	if ptr != nil {
		return &SoundChannel{ptr: ptr}
	}
	return nil
}

// SynthAPI wraps synth functions
type SynthAPI struct{}

// NewSynth creates a new synth
func (sy *SynthAPI) NewSynth() *PDSynth {
	ptr := C.pd_sound_synth_new()
	if ptr != nil {
		return &PDSynth{ptr: ptr}
	}
	return nil
}

// FreeSynth frees a synth
func (sy *SynthAPI) FreeSynth(synth *PDSynth) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_free(synth.ptr)
		synth.ptr = nil
	}
}

// SetWaveform sets the waveform
func (sy *SynthAPI) SetWaveform(synth *PDSynth, wave SoundWaveform) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_setWaveform(synth.ptr, C.int(wave))
	}
}

// SetAttackTime sets the attack time
func (sy *SynthAPI) SetAttackTime(synth *PDSynth, attack float32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_setAttackTime(synth.ptr, C.float(attack))
	}
}

// SetDecayTime sets the decay time
func (sy *SynthAPI) SetDecayTime(synth *PDSynth, decay float32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_setDecayTime(synth.ptr, C.float(decay))
	}
}

// SetSustainLevel sets the sustain level
func (sy *SynthAPI) SetSustainLevel(synth *PDSynth, sustain float32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_setSustainLevel(synth.ptr, C.float(sustain))
	}
}

// SetReleaseTime sets the release time
func (sy *SynthAPI) SetReleaseTime(synth *PDSynth, release float32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_setReleaseTime(synth.ptr, C.float(release))
	}
}

// SetTranspose sets the transpose
func (sy *SynthAPI) SetTranspose(synth *PDSynth, halfSteps float32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_setTranspose(synth.ptr, C.float(halfSteps))
	}
}

// PlayNote plays a note
func (sy *SynthAPI) PlayNote(synth *PDSynth, freq, vel, length float32, when uint32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_playNote(synth.ptr, C.float(freq), C.float(vel), C.float(length), C.uint32_t(when))
	}
}

// PlayMIDINote plays a MIDI note
func (sy *SynthAPI) PlayMIDINote(synth *PDSynth, note MIDINote, vel, length float32, when uint32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_playMIDINote(synth.ptr, C.float(note), C.float(vel), C.float(length), C.uint32_t(when))
	}
}

// NoteOff releases a note
func (sy *SynthAPI) NoteOff(synth *PDSynth, when uint32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_noteOff(synth.ptr, C.uint32_t(when))
	}
}

// Stop stops the synth
func (sy *SynthAPI) Stop(synth *PDSynth) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_stop(synth.ptr)
	}
}

// SetVolume sets the volume
func (sy *SynthAPI) SetVolume(synth *PDSynth, left, right float32) {
	if synth != nil && synth.ptr != nil {
		C.pd_sound_synth_setVolume(synth.ptr, C.float(left), C.float(right))
	}
}

// GetVolume returns the volume
func (sy *SynthAPI) GetVolume(synth *PDSynth) (left, right float32) {
	if synth != nil && synth.ptr != nil {
		var l, r C.float
		C.pd_sound_synth_getVolume(synth.ptr, &l, &r)
		return float32(l), float32(r)
	}
	return 0, 0
}

// IsPlaying returns whether the synth is playing
func (sy *SynthAPI) IsPlaying(synth *PDSynth) bool {
	if synth != nil && synth.ptr != nil {
		return C.pd_sound_synth_isPlaying(synth.ptr) != 0
	}
	return false
}

// SetSample sets the sample for the synth
func (sy *SynthAPI) SetSample(synth *PDSynth, sample *AudioSample, sustainStart, sustainEnd uint32) {
	if synth != nil && synth.ptr != nil {
		var samplePtr unsafe.Pointer
		if sample != nil {
			samplePtr = sample.ptr
		}
		C.pd_sound_synth_setSample(synth.ptr, samplePtr, C.uint32_t(sustainStart), C.uint32_t(sustainEnd))
	}
}

// Copy copies a synth
func (sy *SynthAPI) Copy(synth *PDSynth) *PDSynth {
	if synth != nil && synth.ptr != nil {
		ptr := C.pd_sound_synth_copy(synth.ptr)
		if ptr != nil {
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
	if channel != nil && channel.ptr != nil && inst != nil && inst.ptr != nil {
		return C.pd_sound_channel_addInstrument(channel.ptr, inst.ptr) != 0
	}
	return false
}

// ============== SampleAPI ==============

// SampleAPI wraps sample functions
type SampleAPI struct{}

// Load loads an audio sample from file
func (sa *SampleAPI) Load(path string) *AudioSample {
	cpath := make([]byte, len(path)+1)
	copy(cpath, path)
	ptr := C.pd_sound_sample_load((*C.char)(unsafe.Pointer(&cpath[0])))
	if ptr != nil {
		return &AudioSample{ptr: ptr}
	}
	return nil
}

// ============== SequenceAPI ==============

// SequenceAPI wraps sequence functions
type SequenceAPI struct{}

// NewSequence creates a new sequence
func (s *SequenceAPI) NewSequence() *SoundSequence {
	ptr := C.pd_sound_sequence_new()
	if ptr != nil {
		return &SoundSequence{ptr: ptr}
	}
	return nil
}

// LoadMIDIFile loads a MIDI file into the sequence
func (s *SequenceAPI) LoadMIDIFile(seq *SoundSequence, path string) error {
	if seq != nil && seq.ptr != nil {
		cpath := make([]byte, len(path)+1)
		copy(cpath, path)
		if C.pd_sound_sequence_loadMIDI(seq.ptr, (*C.char)(unsafe.Pointer(&cpath[0]))) != 0 {
			return nil
		}
	}
	return &soundError{op: "loadMIDI", path: path}
}

// GetTrackCount returns the track count
func (s *SequenceAPI) GetTrackCount(seq *SoundSequence) int {
	if seq != nil && seq.ptr != nil {
		return int(C.pd_sound_sequence_getTrackCount(seq.ptr))
	}
	return 0
}

// GetTrackAtIndex returns a track at an index
func (s *SequenceAPI) GetTrackAtIndex(seq *SoundSequence, index uint) *SequenceTrack {
	if seq != nil && seq.ptr != nil {
		ptr := C.pd_sound_sequence_getTrackAtIndex(seq.ptr, C.int(index))
		if ptr != nil {
			return &SequenceTrack{ptr: ptr}
		}
	}
	return nil
}

// Play plays the sequence
func (s *SequenceAPI) Play(seq *SoundSequence) {
	if seq != nil && seq.ptr != nil {
		C.pd_sound_sequence_play(seq.ptr, nil, nil)
	}
}

// Stop stops the sequence
func (s *SequenceAPI) Stop(seq *SoundSequence) {
	if seq != nil && seq.ptr != nil {
		C.pd_sound_sequence_stop(seq.ptr)
	}
}

// GetCurrentStep returns the current step in the sequence
func (s *SequenceAPI) GetCurrentStep(seq *SoundSequence) int {
	if seq != nil && seq.ptr != nil {
		return int(C.pd_sound_sequence_getCurrentStep(seq.ptr, nil))
	}
	return 0
}

// ============== TrackAPI ==============

// TrackAPI wraps track functions
type TrackAPI struct{}

// SetInstrument sets the instrument for a track
func (t *TrackAPI) SetInstrument(track *SequenceTrack, inst *PDSynthInstrument) {
	if track != nil && track.ptr != nil {
		var instPtr unsafe.Pointer
		if inst != nil {
			instPtr = inst.ptr
		}
		C.pd_sound_track_setInstrument(track.ptr, instPtr)
	}
}

// GetPolyphony returns the maximum polyphony for the track
func (t *TrackAPI) GetPolyphony(track *SequenceTrack) int {
	if track != nil && track.ptr != nil {
		return int(C.pd_sound_track_getPolyphony(track.ptr))
	}
	return 0
}

// GetIndexForStep returns the internal index for the step
func (t *TrackAPI) GetIndexForStep(track *SequenceTrack, step uint32) int {
	if track != nil && track.ptr != nil {
		return int(C.pd_sound_track_getIndexForStep(track.ptr, C.uint32_t(step)))
	}
	return 0
}

// GetNoteAtIndex returns note information at the given index
func (t *TrackAPI) GetNoteAtIndex(track *SequenceTrack, index int) (step, length uint32, note MIDINote, velocity float32, ok bool) {
	if track != nil && track.ptr != nil {
		var s, l C.uint32_t
		var n, v C.float
		result := C.pd_sound_track_getNoteAtIndex(track.ptr, C.int(index), &s, &l, &n, &v)
		if result != 0 {
			return uint32(s), uint32(l), MIDINote(n), float32(v), true
		}
	}
	return 0, 0, 0, 0, false
}

// ============== InstrumentAPI ==============

// InstrumentAPI wraps instrument functions
type InstrumentAPI struct{}

// NewInstrument creates a new instrument
func (i *InstrumentAPI) NewInstrument() *PDSynthInstrument {
	ptr := C.pd_sound_instrument_new()
	if ptr != nil {
		return &PDSynthInstrument{ptr: ptr}
	}
	return nil
}

// SetVolume sets the instrument volume
func (i *InstrumentAPI) SetVolume(inst *PDSynthInstrument, left, right float32) {
	if inst != nil && inst.ptr != nil {
		C.pd_sound_instrument_setVolume(inst.ptr, C.float(left), C.float(right))
	}
}

// AddVoice adds a voice to the instrument
func (i *InstrumentAPI) AddVoice(inst *PDSynthInstrument, synth *PDSynth, rangeStart, rangeEnd MIDINote, transpose float32) bool {
	if inst != nil && inst.ptr != nil && synth != nil && synth.ptr != nil {
		return C.pd_sound_instrument_addVoice(inst.ptr, synth.ptr, C.float(rangeStart), C.float(rangeEnd), C.float(transpose)) != 0
	}
	return false
}

// ============== FilePlayer ==============

// NewFilePlayer creates a new file player
func (s *Sound) NewFilePlayer() *FilePlayer {
	ptr := C.pd_sound_fileplayer_new()
	if ptr != nil {
		return &FilePlayer{ptr: ptr}
	}
	return nil
}

// FreeFilePlayer frees a file player
func (s *Sound) FreeFilePlayer(player *FilePlayer) {
	if player != nil && player.ptr != nil {
		C.pd_sound_fileplayer_free(player.ptr)
		player.ptr = nil
	}
}

// LoadIntoFilePlayer loads audio into file player
func (s *Sound) LoadIntoFilePlayer(player *FilePlayer, path string) error {
	if player != nil && player.ptr != nil {
		cpath := make([]byte, len(path)+1)
		copy(cpath, path)
		if C.pd_sound_fileplayer_load(player.ptr, (*C.char)(unsafe.Pointer(&cpath[0]))) != 0 {
			return nil
		}
	}
	return &soundError{op: "load", path: path}
}

// PlayFilePlayer plays audio
func (s *Sound) PlayFilePlayer(player *FilePlayer, repeat int) {
	if player != nil && player.ptr != nil {
		C.pd_sound_fileplayer_play(player.ptr, C.int(repeat))
	}
}

// StopFilePlayer stops playback
func (s *Sound) StopFilePlayer(player *FilePlayer) {
	if player != nil && player.ptr != nil {
		C.pd_sound_fileplayer_stop(player.ptr)
	}
}

// PauseFilePlayer pauses playback
func (s *Sound) PauseFilePlayer(player *FilePlayer) {
	if player != nil && player.ptr != nil {
		C.pd_sound_fileplayer_pause(player.ptr)
	}
}

// IsFilePlayerPlaying returns true if playing
func (s *Sound) IsFilePlayerPlaying(player *FilePlayer) bool {
	if player != nil && player.ptr != nil {
		return C.pd_sound_fileplayer_isPlaying(player.ptr) != 0
	}
	return false
}

// SetFilePlayerVolume sets volume (0.0 - 1.0)
func (s *Sound) SetFilePlayerVolume(player *FilePlayer, left, right float32) {
	if player != nil && player.ptr != nil {
		C.pd_sound_fileplayer_setVolume(player.ptr, C.float(left), C.float(right))
	}
}

// GetFilePlayerVolume gets volume
func (s *Sound) GetFilePlayerVolume(player *FilePlayer) (left, right float32) {
	if player != nil && player.ptr != nil {
		var l, r C.float
		C.pd_sound_fileplayer_getVolume(player.ptr, &l, &r)
		return float32(l), float32(r)
	}
	return 0, 0
}

// GetFilePlayerLength returns length in seconds
func (s *Sound) GetFilePlayerLength(player *FilePlayer) float32 {
	if player != nil && player.ptr != nil {
		return float32(C.pd_sound_fileplayer_getLength(player.ptr))
	}
	return 0
}

// SetFilePlayerOffset sets playback offset
func (s *Sound) SetFilePlayerOffset(player *FilePlayer, offset float32) {
	if player != nil && player.ptr != nil {
		C.pd_sound_fileplayer_setOffset(player.ptr, C.float(offset))
	}
}

// GetFilePlayerOffset gets playback offset
func (s *Sound) GetFilePlayerOffset(player *FilePlayer) float32 {
	if player != nil && player.ptr != nil {
		return float32(C.pd_sound_fileplayer_getOffset(player.ptr))
	}
	return 0
}

// SetFilePlayerRate sets playback rate
func (s *Sound) SetFilePlayerRate(player *FilePlayer, rate float32) {
	if player != nil && player.ptr != nil {
		C.pd_sound_fileplayer_setRate(player.ptr, C.float(rate))
	}
}

// ============== SamplePlayer ==============

// NewSamplePlayer creates a new sample player
func (s *Sound) NewSamplePlayer() *SamplePlayer {
	ptr := C.pd_sound_sampleplayer_new()
	if ptr != nil {
		return &SamplePlayer{ptr: ptr}
	}
	return nil
}

// FreeSamplePlayer frees a sample player
func (s *Sound) FreeSamplePlayer(player *SamplePlayer) {
	if player != nil && player.ptr != nil {
		C.pd_sound_sampleplayer_free(player.ptr)
		player.ptr = nil
	}
}

// SetSamplePlayerSample sets the sample to play
func (s *Sound) SetSamplePlayerSample(player *SamplePlayer, sample *AudioSample) {
	if player != nil && player.ptr != nil && sample != nil {
		C.pd_sound_sampleplayer_setSample(player.ptr, sample.ptr)
	}
}

// PlaySamplePlayer plays the sample
func (s *Sound) PlaySamplePlayer(player *SamplePlayer, repeat int, rate float32) {
	if player != nil && player.ptr != nil {
		C.pd_sound_sampleplayer_play(player.ptr, C.int(repeat), C.float(rate))
	}
}

// StopSamplePlayer stops playback
func (s *Sound) StopSamplePlayer(player *SamplePlayer) {
	if player != nil && player.ptr != nil {
		C.pd_sound_sampleplayer_stop(player.ptr)
	}
}

// IsSamplePlayerPlaying returns true if playing
func (s *Sound) IsSamplePlayerPlaying(player *SamplePlayer) bool {
	if player != nil && player.ptr != nil {
		return C.pd_sound_sampleplayer_isPlaying(player.ptr) != 0
	}
	return false
}

// SetSamplePlayerVolume sets volume
func (s *Sound) SetSamplePlayerVolume(player *SamplePlayer, left, right float32) {
	if player != nil && player.ptr != nil {
		C.pd_sound_sampleplayer_setVolume(player.ptr, C.float(left), C.float(right))
	}
}

// ============== AudioSample ==============

// NewAudioSample creates a new audio sample
func (s *Sound) NewAudioSample(length int) *AudioSample {
	ptr := C.pd_sound_sample_new(C.int(length))
	if ptr != nil {
		return &AudioSample{ptr: ptr}
	}
	return nil
}

// LoadAudioSample loads an audio sample from file
func (s *Sound) LoadAudioSample(path string) *AudioSample {
	cpath := make([]byte, len(path)+1)
	copy(cpath, path)
	ptr := C.pd_sound_sample_load((*C.char)(unsafe.Pointer(&cpath[0])))
	if ptr != nil {
		return &AudioSample{ptr: ptr}
	}
	return nil
}

// FreeAudioSample frees an audio sample
func (s *Sound) FreeAudioSample(sample *AudioSample) {
	if sample != nil && sample.ptr != nil {
		C.pd_sound_sample_free(sample.ptr)
		sample.ptr = nil
	}
}

// ============== Global Sound ==============

// SetHeadphoneCallback sets headphone callback (stub)
func (s *Sound) SetHeadphoneCallback(callback func(headphone, mic int)) {
	// Not implemented
}

// GetHeadphoneState returns headphone state
func (s *Sound) GetHeadphoneState() (headphone, mic bool) {
	var h, m C.int
	C.pd_sound_getHeadphoneState(&h, &m)
	return h != 0, m != 0
}

// SetOutputsActive sets which outputs are active
func (s *Sound) SetOutputsActive(headphone, speaker bool) {
	var h, sp C.int
	if headphone {
		h = 1
	}
	if speaker {
		sp = 1
	}
	C.pd_sound_setOutputsActive(h, sp)
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
