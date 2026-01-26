//go:build !tinygo

package pdgo

/*
#cgo CFLAGS: -DTARGET_EXTENSION=1
#include "pd_api.h"
#include <stdlib.h>

// Sound API helper functions

// Channel functions
static SoundChannel* sound_channel_newChannel(const struct playdate_sound* snd) {
    return snd->channel->newChannel();
}

static void sound_channel_freeChannel(const struct playdate_sound* snd, SoundChannel* channel) {
    snd->channel->freeChannel(channel);
}

static int sound_channel_addSource(const struct playdate_sound* snd, SoundChannel* channel, SoundSource* source) {
    return snd->channel->addSource(channel, source);
}

static int sound_channel_removeSource(const struct playdate_sound* snd, SoundChannel* channel, SoundSource* source) {
    return snd->channel->removeSource(channel, source);
}

static int sound_channel_addEffect(const struct playdate_sound* snd, SoundChannel* channel, SoundEffect* effect) {
    return snd->channel->addEffect(channel, effect);
}

static int sound_channel_removeEffect(const struct playdate_sound* snd, SoundChannel* channel, SoundEffect* effect) {
    return snd->channel->removeEffect(channel, effect);
}

static void sound_channel_setVolume(const struct playdate_sound* snd, SoundChannel* channel, float volume) {
    snd->channel->setVolume(channel, volume);
}

static float sound_channel_getVolume(const struct playdate_sound* snd, SoundChannel* channel) {
    return snd->channel->getVolume(channel);
}

static void sound_channel_setPan(const struct playdate_sound* snd, SoundChannel* channel, float pan) {
    snd->channel->setPan(channel, pan);
}

// FilePlayer functions
static FilePlayer* sound_fileplayer_newPlayer(const struct playdate_sound* snd) {
    return snd->fileplayer->newPlayer();
}

static void sound_fileplayer_freePlayer(const struct playdate_sound* snd, FilePlayer* player) {
    snd->fileplayer->freePlayer(player);
}

static int sound_fileplayer_loadIntoPlayer(const struct playdate_sound* snd, FilePlayer* player, const char* path) {
    return snd->fileplayer->loadIntoPlayer(player, path);
}

static void sound_fileplayer_setBufferLength(const struct playdate_sound* snd, FilePlayer* player, float bufferLen) {
    snd->fileplayer->setBufferLength(player, bufferLen);
}

static int sound_fileplayer_play(const struct playdate_sound* snd, FilePlayer* player, int repeat) {
    return snd->fileplayer->play(player, repeat);
}

static int sound_fileplayer_isPlaying(const struct playdate_sound* snd, FilePlayer* player) {
    return snd->fileplayer->isPlaying(player);
}

static void sound_fileplayer_pause(const struct playdate_sound* snd, FilePlayer* player) {
    snd->fileplayer->pause(player);
}

static void sound_fileplayer_stop(const struct playdate_sound* snd, FilePlayer* player) {
    snd->fileplayer->stop(player);
}

static void sound_fileplayer_setVolume(const struct playdate_sound* snd, FilePlayer* player, float left, float right) {
    snd->fileplayer->setVolume(player, left, right);
}

static void sound_fileplayer_getVolume(const struct playdate_sound* snd, FilePlayer* player, float* left, float* right) {
    snd->fileplayer->getVolume(player, left, right);
}

static float sound_fileplayer_getLength(const struct playdate_sound* snd, FilePlayer* player) {
    return snd->fileplayer->getLength(player);
}

static void sound_fileplayer_setOffset(const struct playdate_sound* snd, FilePlayer* player, float offset) {
    snd->fileplayer->setOffset(player, offset);
}

static float sound_fileplayer_getOffset(const struct playdate_sound* snd, FilePlayer* player) {
    return snd->fileplayer->getOffset(player);
}

static void sound_fileplayer_setRate(const struct playdate_sound* snd, FilePlayer* player, float rate) {
    snd->fileplayer->setRate(player, rate);
}

static float sound_fileplayer_getRate(const struct playdate_sound* snd, FilePlayer* player) {
    return snd->fileplayer->getRate(player);
}

static void sound_fileplayer_setLoopRange(const struct playdate_sound* snd, FilePlayer* player, float start, float end) {
    snd->fileplayer->setLoopRange(player, start, end);
}

static int sound_fileplayer_didUnderrun(const struct playdate_sound* snd, FilePlayer* player) {
    return snd->fileplayer->didUnderrun(player);
}

static void sound_fileplayer_setStopOnUnderrun(const struct playdate_sound* snd, FilePlayer* player, int flag) {
    snd->fileplayer->setStopOnUnderrun(player, flag);
}

// Sample functions
static AudioSample* sound_sample_newSampleBuffer(const struct playdate_sound* snd, int byteCount) {
    return snd->sample->newSampleBuffer(byteCount);
}

static int sound_sample_loadIntoSample(const struct playdate_sound* snd, AudioSample* sample, const char* path) {
    return snd->sample->loadIntoSample(sample, path);
}

static AudioSample* sound_sample_load(const struct playdate_sound* snd, const char* path) {
    return snd->sample->load(path);
}

static AudioSample* sound_sample_newSampleFromData(const struct playdate_sound* snd, uint8_t* data, SoundFormat format, uint32_t sampleRate, int byteCount, int shouldFreeData) {
    return snd->sample->newSampleFromData(data, format, sampleRate, byteCount, shouldFreeData);
}

static void sound_sample_getData(const struct playdate_sound* snd, AudioSample* sample, uint8_t** data, SoundFormat* format, uint32_t* sampleRate, uint32_t* bytelength) {
    snd->sample->getData(sample, data, format, sampleRate, bytelength);
}

static void sound_sample_freeSample(const struct playdate_sound* snd, AudioSample* sample) {
    snd->sample->freeSample(sample);
}

static float sound_sample_getLength(const struct playdate_sound* snd, AudioSample* sample) {
    return snd->sample->getLength(sample);
}

static int sound_sample_decompress(const struct playdate_sound* snd, AudioSample* sample) {
    return snd->sample->decompress(sample);
}

// SamplePlayer functions
static SamplePlayer* sound_sampleplayer_newPlayer(const struct playdate_sound* snd) {
    return snd->sampleplayer->newPlayer();
}

static void sound_sampleplayer_freePlayer(const struct playdate_sound* snd, SamplePlayer* player) {
    snd->sampleplayer->freePlayer(player);
}

static void sound_sampleplayer_setSample(const struct playdate_sound* snd, SamplePlayer* player, AudioSample* sample) {
    snd->sampleplayer->setSample(player, sample);
}

static int sound_sampleplayer_play(const struct playdate_sound* snd, SamplePlayer* player, int repeat, float rate) {
    return snd->sampleplayer->play(player, repeat, rate);
}

static int sound_sampleplayer_isPlaying(const struct playdate_sound* snd, SamplePlayer* player) {
    return snd->sampleplayer->isPlaying(player);
}

static void sound_sampleplayer_stop(const struct playdate_sound* snd, SamplePlayer* player) {
    snd->sampleplayer->stop(player);
}

static void sound_sampleplayer_setVolume(const struct playdate_sound* snd, SamplePlayer* player, float left, float right) {
    snd->sampleplayer->setVolume(player, left, right);
}

static void sound_sampleplayer_getVolume(const struct playdate_sound* snd, SamplePlayer* player, float* left, float* right) {
    snd->sampleplayer->getVolume(player, left, right);
}

static float sound_sampleplayer_getLength(const struct playdate_sound* snd, SamplePlayer* player) {
    return snd->sampleplayer->getLength(player);
}

static void sound_sampleplayer_setOffset(const struct playdate_sound* snd, SamplePlayer* player, float offset) {
    snd->sampleplayer->setOffset(player, offset);
}

static float sound_sampleplayer_getOffset(const struct playdate_sound* snd, SamplePlayer* player) {
    return snd->sampleplayer->getOffset(player);
}

static void sound_sampleplayer_setRate(const struct playdate_sound* snd, SamplePlayer* player, float rate) {
    snd->sampleplayer->setRate(player, rate);
}

static float sound_sampleplayer_getRate(const struct playdate_sound* snd, SamplePlayer* player) {
    return snd->sampleplayer->getRate(player);
}

static void sound_sampleplayer_setPlayRange(const struct playdate_sound* snd, SamplePlayer* player, int start, int end) {
    snd->sampleplayer->setPlayRange(player, start, end);
}

static void sound_sampleplayer_setPaused(const struct playdate_sound* snd, SamplePlayer* player, int flag) {
    snd->sampleplayer->setPaused(player, flag);
}

// Synth functions
static PDSynth* sound_synth_newSynth(const struct playdate_sound* snd) {
    return snd->synth->newSynth();
}

static void sound_synth_freeSynth(const struct playdate_sound* snd, PDSynth* synth) {
    snd->synth->freeSynth(synth);
}

static void sound_synth_setWaveform(const struct playdate_sound* snd, PDSynth* synth, SoundWaveform wave) {
    snd->synth->setWaveform(synth, wave);
}

static void sound_synth_setSample(const struct playdate_sound* snd, PDSynth* synth, AudioSample* sample, uint32_t sustainStart, uint32_t sustainEnd) {
    snd->synth->setSample(synth, sample, sustainStart, sustainEnd);
}

static void sound_synth_setAttackTime(const struct playdate_sound* snd, PDSynth* synth, float attack) {
    snd->synth->setAttackTime(synth, attack);
}

static void sound_synth_setDecayTime(const struct playdate_sound* snd, PDSynth* synth, float decay) {
    snd->synth->setDecayTime(synth, decay);
}

static void sound_synth_setSustainLevel(const struct playdate_sound* snd, PDSynth* synth, float sustain) {
    snd->synth->setSustainLevel(synth, sustain);
}

static void sound_synth_setReleaseTime(const struct playdate_sound* snd, PDSynth* synth, float release) {
    snd->synth->setReleaseTime(synth, release);
}

static void sound_synth_setTranspose(const struct playdate_sound* snd, PDSynth* synth, float halfSteps) {
    snd->synth->setTranspose(synth, halfSteps);
}

static void sound_synth_setFrequencyModulator(const struct playdate_sound* snd, PDSynth* synth, PDSynthSignalValue* mod) {
    snd->synth->setFrequencyModulator(synth, mod);
}

static PDSynthSignalValue* sound_synth_getFrequencyModulator(const struct playdate_sound* snd, PDSynth* synth) {
    return snd->synth->getFrequencyModulator(synth);
}

static void sound_synth_setAmplitudeModulator(const struct playdate_sound* snd, PDSynth* synth, PDSynthSignalValue* mod) {
    snd->synth->setAmplitudeModulator(synth, mod);
}

static PDSynthSignalValue* sound_synth_getAmplitudeModulator(const struct playdate_sound* snd, PDSynth* synth) {
    return snd->synth->getAmplitudeModulator(synth);
}

static int sound_synth_getParameterCount(const struct playdate_sound* snd, PDSynth* synth) {
    return snd->synth->getParameterCount(synth);
}

static int sound_synth_setParameter(const struct playdate_sound* snd, PDSynth* synth, int parameter, float value) {
    return snd->synth->setParameter(synth, parameter, value);
}

static void sound_synth_setParameterModulator(const struct playdate_sound* snd, PDSynth* synth, int parameter, PDSynthSignalValue* mod) {
    snd->synth->setParameterModulator(synth, parameter, mod);
}

static PDSynthSignalValue* sound_synth_getParameterModulator(const struct playdate_sound* snd, PDSynth* synth, int parameter) {
    return snd->synth->getParameterModulator(synth, parameter);
}

static void sound_synth_playNote(const struct playdate_sound* snd, PDSynth* synth, float freq, float vel, float len, uint32_t when) {
    snd->synth->playNote(synth, freq, vel, len, when);
}

static void sound_synth_playMIDINote(const struct playdate_sound* snd, PDSynth* synth, MIDINote note, float vel, float len, uint32_t when) {
    snd->synth->playMIDINote(synth, note, vel, len, when);
}

static void sound_synth_noteOff(const struct playdate_sound* snd, PDSynth* synth, uint32_t when) {
    snd->synth->noteOff(synth, when);
}

static void sound_synth_stop(const struct playdate_sound* snd, PDSynth* synth) {
    snd->synth->stop(synth);
}

static void sound_synth_setVolume(const struct playdate_sound* snd, PDSynth* synth, float left, float right) {
    snd->synth->setVolume(synth, left, right);
}

static void sound_synth_getVolume(const struct playdate_sound* snd, PDSynth* synth, float* left, float* right) {
    snd->synth->getVolume(synth, left, right);
}

static int sound_synth_isPlaying(const struct playdate_sound* snd, PDSynth* synth) {
    return snd->synth->isPlaying(synth);
}

static PDSynthEnvelope* sound_synth_getEnvelope(const struct playdate_sound* snd, PDSynth* synth) {
    return snd->synth->getEnvelope(synth);
}

static int sound_synth_setWavetable(const struct playdate_sound* snd, PDSynth* synth, AudioSample* sample, int log2size, int columns, int rows) {
    return snd->synth->setWavetable(synth, sample, log2size, columns, rows);
}

static PDSynth* sound_synth_copy(const struct playdate_sound* snd, PDSynth* synth) {
    return snd->synth->copy(synth);
}

static void sound_synth_clearEnvelope(const struct playdate_sound* snd, PDSynth* synth) {
    snd->synth->clearEnvelope(synth);
}

// LFO functions
static PDSynthLFO* sound_lfo_newLFO(const struct playdate_sound* snd, LFOType type) {
    return snd->lfo->newLFO(type);
}

static void sound_lfo_freeLFO(const struct playdate_sound* snd, PDSynthLFO* lfo) {
    snd->lfo->freeLFO(lfo);
}

static void sound_lfo_setType(const struct playdate_sound* snd, PDSynthLFO* lfo, LFOType type) {
    snd->lfo->setType(lfo, type);
}

static void sound_lfo_setRate(const struct playdate_sound* snd, PDSynthLFO* lfo, float rate) {
    snd->lfo->setRate(lfo, rate);
}

static void sound_lfo_setPhase(const struct playdate_sound* snd, PDSynthLFO* lfo, float phase) {
    snd->lfo->setPhase(lfo, phase);
}

static void sound_lfo_setCenter(const struct playdate_sound* snd, PDSynthLFO* lfo, float center) {
    snd->lfo->setCenter(lfo, center);
}

static void sound_lfo_setDepth(const struct playdate_sound* snd, PDSynthLFO* lfo, float depth) {
    snd->lfo->setDepth(lfo, depth);
}

static void sound_lfo_setArpeggiation(const struct playdate_sound* snd, PDSynthLFO* lfo, int nSteps, float* steps) {
    snd->lfo->setArpeggiation(lfo, nSteps, steps);
}

static void sound_lfo_setDelay(const struct playdate_sound* snd, PDSynthLFO* lfo, float holdoff, float ramptime) {
    snd->lfo->setDelay(lfo, holdoff, ramptime);
}

static void sound_lfo_setRetrigger(const struct playdate_sound* snd, PDSynthLFO* lfo, int flag) {
    snd->lfo->setRetrigger(lfo, flag);
}

static float sound_lfo_getValue(const struct playdate_sound* snd, PDSynthLFO* lfo) {
    return snd->lfo->getValue(lfo);
}

static void sound_lfo_setGlobal(const struct playdate_sound* snd, PDSynthLFO* lfo, int global) {
    snd->lfo->setGlobal(lfo, global);
}

static void sound_lfo_setStartPhase(const struct playdate_sound* snd, PDSynthLFO* lfo, float phase) {
    snd->lfo->setStartPhase(lfo, phase);
}

// Envelope functions
static PDSynthEnvelope* sound_envelope_newEnvelope(const struct playdate_sound* snd, float attack, float decay, float sustain, float release) {
    return snd->envelope->newEnvelope(attack, decay, sustain, release);
}

static void sound_envelope_freeEnvelope(const struct playdate_sound* snd, PDSynthEnvelope* env) {
    snd->envelope->freeEnvelope(env);
}

static void sound_envelope_setAttack(const struct playdate_sound* snd, PDSynthEnvelope* env, float attack) {
    snd->envelope->setAttack(env, attack);
}

static void sound_envelope_setDecay(const struct playdate_sound* snd, PDSynthEnvelope* env, float decay) {
    snd->envelope->setDecay(env, decay);
}

static void sound_envelope_setSustain(const struct playdate_sound* snd, PDSynthEnvelope* env, float sustain) {
    snd->envelope->setSustain(env, sustain);
}

static void sound_envelope_setRelease(const struct playdate_sound* snd, PDSynthEnvelope* env, float release) {
    snd->envelope->setRelease(env, release);
}

static void sound_envelope_setLegato(const struct playdate_sound* snd, PDSynthEnvelope* env, int flag) {
    snd->envelope->setLegato(env, flag);
}

static void sound_envelope_setRetrigger(const struct playdate_sound* snd, PDSynthEnvelope* env, int flag) {
    snd->envelope->setRetrigger(env, flag);
}

static float sound_envelope_getValue(const struct playdate_sound* snd, PDSynthEnvelope* env) {
    return snd->envelope->getValue(env);
}

static void sound_envelope_setCurvature(const struct playdate_sound* snd, PDSynthEnvelope* env, float amount) {
    snd->envelope->setCurvature(env, amount);
}

static void sound_envelope_setVelocitySensitivity(const struct playdate_sound* snd, PDSynthEnvelope* env, float velsens) {
    snd->envelope->setVelocitySensitivity(env, velsens);
}

static void sound_envelope_setRateScaling(const struct playdate_sound* snd, PDSynthEnvelope* env, float scaling, MIDINote start, MIDINote end) {
    snd->envelope->setRateScaling(env, scaling, start, end);
}

// Sequence functions
static SoundSequence* sound_sequence_newSequence(const struct playdate_sound* snd) {
    return snd->sequence->newSequence();
}

static void sound_sequence_freeSequence(const struct playdate_sound* snd, SoundSequence* sequence) {
    snd->sequence->freeSequence(sequence);
}

static int sound_sequence_loadMIDIFile(const struct playdate_sound* snd, SoundSequence* seq, const char* path) {
    return snd->sequence->loadMIDIFile(seq, path);
}

static uint32_t sound_sequence_getTime(const struct playdate_sound* snd, SoundSequence* seq) {
    return snd->sequence->getTime(seq);
}

static void sound_sequence_setTime(const struct playdate_sound* snd, SoundSequence* seq, uint32_t time) {
    snd->sequence->setTime(seq, time);
}

static void sound_sequence_setLoops(const struct playdate_sound* snd, SoundSequence* seq, int loopstart, int loopend, int loops) {
    snd->sequence->setLoops(seq, loopstart, loopend, loops);
}

static void sound_sequence_setTempo(const struct playdate_sound* snd, SoundSequence* seq, float stepsPerSecond) {
    snd->sequence->setTempo(seq, stepsPerSecond);
}

static float sound_sequence_getTempo(const struct playdate_sound* snd, SoundSequence* seq) {
    return snd->sequence->getTempo(seq);
}

static int sound_sequence_getTrackCount(const struct playdate_sound* snd, SoundSequence* seq) {
    return snd->sequence->getTrackCount(seq);
}

static SequenceTrack* sound_sequence_addTrack(const struct playdate_sound* snd, SoundSequence* seq) {
    return snd->sequence->addTrack(seq);
}

static SequenceTrack* sound_sequence_getTrackAtIndex(const struct playdate_sound* snd, SoundSequence* seq, unsigned int track) {
    return snd->sequence->getTrackAtIndex(seq, track);
}

static void sound_sequence_setTrackAtIndex(const struct playdate_sound* snd, SoundSequence* seq, SequenceTrack* track, unsigned int idx) {
    snd->sequence->setTrackAtIndex(seq, track, idx);
}

static void sound_sequence_allNotesOff(const struct playdate_sound* snd, SoundSequence* seq) {
    snd->sequence->allNotesOff(seq);
}

static int sound_sequence_isPlaying(const struct playdate_sound* snd, SoundSequence* seq) {
    return snd->sequence->isPlaying(seq);
}

static uint32_t sound_sequence_getLength(const struct playdate_sound* snd, SoundSequence* seq) {
    return snd->sequence->getLength(seq);
}

static void sound_sequence_play(const struct playdate_sound* snd, SoundSequence* seq, SequenceFinishedCallback finishCallback, void* userdata) {
    snd->sequence->play(seq, finishCallback, userdata);
}

static void sound_sequence_stop(const struct playdate_sound* snd, SoundSequence* seq) {
    snd->sequence->stop(seq);
}

static int sound_sequence_getCurrentStep(const struct playdate_sound* snd, SoundSequence* seq, int* timeOffset) {
    return snd->sequence->getCurrentStep(seq, timeOffset);
}

static void sound_sequence_setCurrentStep(const struct playdate_sound* snd, SoundSequence* seq, int step, int timeOffset, int playNotes) {
    snd->sequence->setCurrentStep(seq, step, timeOffset, playNotes);
}

// Track functions
static SequenceTrack* sound_track_newTrack(const struct playdate_sound* snd) {
    return snd->track->newTrack();
}

static void sound_track_freeTrack(const struct playdate_sound* snd, SequenceTrack* track) {
    snd->track->freeTrack(track);
}

static void sound_track_setInstrument(const struct playdate_sound* snd, SequenceTrack* track, PDSynthInstrument* inst) {
    snd->track->setInstrument(track, inst);
}

static PDSynthInstrument* sound_track_getInstrument(const struct playdate_sound* snd, SequenceTrack* track) {
    return snd->track->getInstrument(track);
}

static void sound_track_addNoteEvent(const struct playdate_sound* snd, SequenceTrack* track, uint32_t step, uint32_t len, MIDINote note, float velocity) {
    snd->track->addNoteEvent(track, step, len, note, velocity);
}

static void sound_track_removeNoteEvent(const struct playdate_sound* snd, SequenceTrack* track, uint32_t step, MIDINote note) {
    snd->track->removeNoteEvent(track, step, note);
}

static void sound_track_clearNotes(const struct playdate_sound* snd, SequenceTrack* track) {
    snd->track->clearNotes(track);
}

static int sound_track_getControlSignalCount(const struct playdate_sound* snd, SequenceTrack* track) {
    return snd->track->getControlSignalCount(track);
}

static ControlSignal* sound_track_getControlSignal(const struct playdate_sound* snd, SequenceTrack* track, int idx) {
    return snd->track->getControlSignal(track, idx);
}

static void sound_track_clearControlEvents(const struct playdate_sound* snd, SequenceTrack* track) {
    snd->track->clearControlEvents(track);
}

static int sound_track_getPolyphony(const struct playdate_sound* snd, SequenceTrack* track) {
    return snd->track->getPolyphony(track);
}

static int sound_track_activeVoiceCount(const struct playdate_sound* snd, SequenceTrack* track) {
    return snd->track->activeVoiceCount(track);
}

static void sound_track_setMuted(const struct playdate_sound* snd, SequenceTrack* track, int mute) {
    snd->track->setMuted(track, mute);
}

static uint32_t sound_track_getLength(const struct playdate_sound* snd, SequenceTrack* track) {
    return snd->track->getLength(track);
}

static int sound_track_getIndexForStep(const struct playdate_sound* snd, SequenceTrack* track, uint32_t step) {
    return snd->track->getIndexForStep(track, step);
}

static int sound_track_getNoteAtIndex(const struct playdate_sound* snd, SequenceTrack* track, int index, uint32_t* outStep, uint32_t* outLen, MIDINote* outNote, float* outVelocity) {
    return snd->track->getNoteAtIndex(track, index, outStep, outLen, outNote, outVelocity);
}

static ControlSignal* sound_track_getSignalForController(const struct playdate_sound* snd, SequenceTrack* track, int controller, int create) {
    return snd->track->getSignalForController(track, controller, create);
}

// Instrument functions
static PDSynthInstrument* sound_instrument_newInstrument(const struct playdate_sound* snd) {
    return snd->instrument->newInstrument();
}

static void sound_instrument_freeInstrument(const struct playdate_sound* snd, PDSynthInstrument* inst) {
    snd->instrument->freeInstrument(inst);
}

static int sound_instrument_addVoice(const struct playdate_sound* snd, PDSynthInstrument* inst, PDSynth* synth, MIDINote rangeStart, MIDINote rangeEnd, float transpose) {
    return snd->instrument->addVoice(inst, synth, rangeStart, rangeEnd, transpose);
}

static PDSynth* sound_instrument_playNote(const struct playdate_sound* snd, PDSynthInstrument* inst, float frequency, float vel, float len, uint32_t when) {
    return snd->instrument->playNote(inst, frequency, vel, len, when);
}

static PDSynth* sound_instrument_playMIDINote(const struct playdate_sound* snd, PDSynthInstrument* inst, MIDINote note, float vel, float len, uint32_t when) {
    return snd->instrument->playMIDINote(inst, note, vel, len, when);
}

static void sound_instrument_setPitchBend(const struct playdate_sound* snd, PDSynthInstrument* inst, float bend) {
    snd->instrument->setPitchBend(inst, bend);
}

static void sound_instrument_setPitchBendRange(const struct playdate_sound* snd, PDSynthInstrument* inst, float halfSteps) {
    snd->instrument->setPitchBendRange(inst, halfSteps);
}

static void sound_instrument_setTranspose(const struct playdate_sound* snd, PDSynthInstrument* inst, float halfSteps) {
    snd->instrument->setTranspose(inst, halfSteps);
}

static void sound_instrument_noteOff(const struct playdate_sound* snd, PDSynthInstrument* inst, MIDINote note, uint32_t when) {
    snd->instrument->noteOff(inst, note, when);
}

static void sound_instrument_allNotesOff(const struct playdate_sound* snd, PDSynthInstrument* inst, uint32_t when) {
    snd->instrument->allNotesOff(inst, when);
}

static void sound_instrument_setVolume(const struct playdate_sound* snd, PDSynthInstrument* inst, float left, float right) {
    snd->instrument->setVolume(inst, left, right);
}

static void sound_instrument_getVolume(const struct playdate_sound* snd, PDSynthInstrument* inst, float* left, float* right) {
    snd->instrument->getVolume(inst, left, right);
}

static int sound_instrument_activeVoiceCount(const struct playdate_sound* snd, PDSynthInstrument* inst) {
    return snd->instrument->activeVoiceCount(inst);
}

// Global sound functions
static uint32_t sound_getCurrentTime(const struct playdate_sound* snd) {
    return snd->getCurrentTime();
}

static SoundChannel* sound_getDefaultChannel(const struct playdate_sound* snd) {
    return snd->getDefaultChannel();
}

static int sound_addChannel(const struct playdate_sound* snd, SoundChannel* channel) {
    return snd->addChannel(channel);
}

static int sound_removeChannel(const struct playdate_sound* snd, SoundChannel* channel) {
    return snd->removeChannel(channel);
}

static void sound_getHeadphoneState(const struct playdate_sound* snd, int* headphone, int* headsetmic) {
    snd->getHeadphoneState(headphone, headsetmic, NULL);
}

static void sound_setOutputsActive(const struct playdate_sound* snd, int headphone, int speaker) {
    snd->setOutputsActive(headphone, speaker);
}

static int sound_removeSource(const struct playdate_sound* snd, SoundSource* source) {
    return snd->removeSource(source);
}

static const char* sound_getError(const struct playdate_sound* snd) {
    return snd->getError();
}
*/
import "C"
import (
	"errors"
	"math"
)

// Audio constants
const (
	AudioFramesPerCycle = 512
	NoteC4              = 60
)

// SoundFormat represents audio format
type SoundFormat int

const (
	Sound8bitMono    SoundFormat = C.kSound8bitMono
	Sound8bitStereo  SoundFormat = C.kSound8bitStereo
	Sound16bitMono   SoundFormat = C.kSound16bitMono
	Sound16bitStereo SoundFormat = C.kSound16bitStereo
	SoundADPCMMono   SoundFormat = C.kSoundADPCMMono
	SoundADPCMStereo SoundFormat = C.kSoundADPCMStereo
)

// IsStereo returns true if the format is stereo
func (f SoundFormat) IsStereo() bool {
	return (f & 1) != 0
}

// Is16bit returns true if the format is 16-bit
func (f SoundFormat) Is16bit() bool {
	return f >= Sound16bitMono
}

// SoundWaveform represents synth waveforms
type SoundWaveform int

const (
	WaveformSquare    SoundWaveform = C.kWaveformSquare
	WaveformTriangle  SoundWaveform = C.kWaveformTriangle
	WaveformSine      SoundWaveform = C.kWaveformSine
	WaveformNoise     SoundWaveform = C.kWaveformNoise
	WaveformSawtooth  SoundWaveform = C.kWaveformSawtooth
	WaveformPOPhase   SoundWaveform = C.kWaveformPOPhase
	WaveformPODigital SoundWaveform = C.kWaveformPODigital
	WaveformPOVosim   SoundWaveform = C.kWaveformPOVosim
)

// LFOType represents LFO types
type LFOType int

const (
	LFOTypeSquare        LFOType = C.kLFOTypeSquare
	LFOTypeTriangle      LFOType = C.kLFOTypeTriangle
	LFOTypeSine          LFOType = C.kLFOTypeSine
	LFOTypeSampleAndHold LFOType = C.kLFOTypeSampleAndHold
	LFOTypeSawtoothUp    LFOType = C.kLFOTypeSawtoothUp
	LFOTypeSawtoothDown  LFOType = C.kLFOTypeSawtoothDown
	LFOTypeArpeggiator   LFOType = C.kLFOTypeArpeggiator
	LFOTypeFunction      LFOType = C.kLFOTypeFunction
)

// MIDINote represents a MIDI note
type MIDINote float32

// NoteToFrequency converts a MIDI note to frequency
func NoteToFrequency(note MIDINote) float32 {
	return float32(440 * math.Pow(2.0, float64(note-69)/12.0))
}

// FrequencyToNote converts frequency to MIDI note
func FrequencyToNote(freq float32) MIDINote {
	return MIDINote(12*math.Log2(float64(freq)) - 36.376316562)
}

// Sound type wrappers
type (
	SoundSource        struct{ ptr *C.SoundSource }
	SoundChannel       struct{ ptr *C.SoundChannel }
	SoundEffect        struct{ ptr *C.SoundEffect }
	FilePlayerC        struct{ ptr *C.FilePlayer }
	AudioSample        struct{ ptr *C.AudioSample }
	SamplePlayer       struct{ ptr *C.SamplePlayer }
	PDSynth            struct{ ptr *C.PDSynth }
	PDSynthLFO         struct{ ptr *C.PDSynthLFO }
	PDSynthEnvelope    struct{ ptr *C.PDSynthEnvelope }
	PDSynthSignalValue struct{ ptr *C.PDSynthSignalValue }
	SoundSequence      struct{ ptr *C.SoundSequence }
	SequenceTrack      struct{ ptr *C.SequenceTrack }
	PDSynthInstrument  struct{ ptr *C.PDSynthInstrument }
	ControlSignal      struct{ ptr *C.ControlSignal }
)

// Sound wraps the playdate_sound API
type Sound struct {
	ptr          *C.struct_playdate_sound
	Channel      *ChannelAPI
	FilePlayer   *FilePlayerAPI
	Sample       *SampleAPI
	SamplePlayer *SamplePlayerAPI
	Synth        *SynthAPI
	LFO          *LFOAPI
	Envelope     *EnvelopeAPI
	Sequence     *SequenceAPI
	Track        *TrackAPI
	Instrument   *InstrumentAPI
}

func newSound(ptr *C.struct_playdate_sound) *Sound {
	s := &Sound{ptr: ptr}
	s.Channel = &ChannelAPI{snd: ptr}
	s.FilePlayer = &FilePlayerAPI{snd: ptr}
	s.Sample = &SampleAPI{snd: ptr}
	s.SamplePlayer = &SamplePlayerAPI{snd: ptr}
	s.Synth = &SynthAPI{snd: ptr}
	s.LFO = &LFOAPI{snd: ptr}
	s.Envelope = &EnvelopeAPI{snd: ptr}
	s.Sequence = &SequenceAPI{snd: ptr}
	s.Track = &TrackAPI{snd: ptr}
	s.Instrument = &InstrumentAPI{snd: ptr}
	return s
}

// GetCurrentTime returns the current sound time in samples
func (s *Sound) GetCurrentTime() uint32 {
	return uint32(C.sound_getCurrentTime(s.ptr))
}

// GetDefaultChannel returns the default sound channel
func (s *Sound) GetDefaultChannel() *SoundChannel {
	ptr := C.sound_getDefaultChannel(s.ptr)
	if ptr == nil {
		return nil
	}
	return &SoundChannel{ptr: ptr}
}

// AddChannel adds a channel to the sound system
func (s *Sound) AddChannel(channel *SoundChannel) bool {
	if channel == nil {
		return false
	}
	return C.sound_addChannel(s.ptr, channel.ptr) != 0
}

// RemoveChannel removes a channel from the sound system
func (s *Sound) RemoveChannel(channel *SoundChannel) bool {
	if channel == nil {
		return false
	}
	return C.sound_removeChannel(s.ptr, channel.ptr) != 0
}

// GetHeadphoneState returns headphone and headset mic state
func (s *Sound) GetHeadphoneState() (headphone, headsetMic bool) {
	var h, m C.int
	C.sound_getHeadphoneState(s.ptr, &h, &m)
	return h != 0, m != 0
}

// SetOutputsActive sets which outputs are active
func (s *Sound) SetOutputsActive(headphone, speaker bool) {
	h, sp := 0, 0
	if headphone {
		h = 1
	}
	if speaker {
		sp = 1
	}
	C.sound_setOutputsActive(s.ptr, C.int(h), C.int(sp))
}

// RemoveSource removes a sound source
func (s *Sound) RemoveSource(source *SoundSource) bool {
	if source == nil {
		return false
	}
	return C.sound_removeSource(s.ptr, source.ptr) != 0
}

// GetError returns the last sound error
func (s *Sound) GetError() string {
	err := C.sound_getError(s.ptr)
	if err == nil {
		return ""
	}
	return goString(err)
}

// ChannelAPI wraps channel functions
type ChannelAPI struct {
	snd *C.struct_playdate_sound
}

// NewChannel creates a new sound channel
func (c *ChannelAPI) NewChannel() *SoundChannel {
	ptr := C.sound_channel_newChannel(c.snd)
	if ptr == nil {
		return nil
	}
	return &SoundChannel{ptr: ptr}
}

// FreeChannel frees a sound channel
func (c *ChannelAPI) FreeChannel(channel *SoundChannel) {
	if channel != nil && channel.ptr != nil {
		C.sound_channel_freeChannel(c.snd, channel.ptr)
		channel.ptr = nil
	}
}

// AddSource adds a source to a channel
func (c *ChannelAPI) AddSource(channel *SoundChannel, source *SoundSource) bool {
	if channel == nil || source == nil {
		return false
	}
	return C.sound_channel_addSource(c.snd, channel.ptr, source.ptr) != 0
}

// AddInstrumentAsSource adds an instrument as a source to a channel
func (c *ChannelAPI) AddInstrumentAsSource(channel *SoundChannel, inst *PDSynthInstrument) bool {
	if channel == nil || inst == nil {
		return false
	}
	return C.sound_channel_addSource(c.snd, channel.ptr, (*C.SoundSource)(inst.ptr)) != 0
}

// RemoveSource removes a source from a channel
func (c *ChannelAPI) RemoveSource(channel *SoundChannel, source *SoundSource) bool {
	if channel == nil || source == nil {
		return false
	}
	return C.sound_channel_removeSource(c.snd, channel.ptr, source.ptr) != 0
}

// AddEffect adds an effect to a channel
func (c *ChannelAPI) AddEffect(channel *SoundChannel, effect *SoundEffect) bool {
	if channel == nil || effect == nil {
		return false
	}
	return C.sound_channel_addEffect(c.snd, channel.ptr, effect.ptr) != 0
}

// RemoveEffect removes an effect from a channel
func (c *ChannelAPI) RemoveEffect(channel *SoundChannel, effect *SoundEffect) bool {
	if channel == nil || effect == nil {
		return false
	}
	return C.sound_channel_removeEffect(c.snd, channel.ptr, effect.ptr) != 0
}

// SetVolume sets the channel volume
func (c *ChannelAPI) SetVolume(channel *SoundChannel, volume float32) {
	if channel != nil {
		C.sound_channel_setVolume(c.snd, channel.ptr, C.float(volume))
	}
}

// GetVolume returns the channel volume
func (c *ChannelAPI) GetVolume(channel *SoundChannel) float32 {
	if channel == nil {
		return 0
	}
	return float32(C.sound_channel_getVolume(c.snd, channel.ptr))
}

// SetPan sets the channel pan
func (c *ChannelAPI) SetPan(channel *SoundChannel, pan float32) {
	if channel != nil {
		C.sound_channel_setPan(c.snd, channel.ptr, C.float(pan))
	}
}

// FilePlayerAPI wraps file player functions
type FilePlayerAPI struct {
	snd *C.struct_playdate_sound
}

// NewPlayer creates a new file player
func (f *FilePlayerAPI) NewPlayer() *FilePlayerC {
	ptr := C.sound_fileplayer_newPlayer(f.snd)
	if ptr == nil {
		return nil
	}
	return &FilePlayerC{ptr: ptr}
}

// FreePlayer frees a file player
func (f *FilePlayerAPI) FreePlayer(player *FilePlayerC) {
	if player != nil && player.ptr != nil {
		C.sound_fileplayer_freePlayer(f.snd, player.ptr)
		player.ptr = nil
	}
}

// LoadIntoPlayer loads a file into a player
func (f *FilePlayerAPI) LoadIntoPlayer(player *FilePlayerC, path string) error {
	if player == nil {
		return errors.New("player is nil")
	}
	cpath := cString(path)
	defer freeCString(cpath)

	if C.sound_fileplayer_loadIntoPlayer(f.snd, player.ptr, cpath) == 0 {
		return errors.New("failed to load file")
	}
	return nil
}

// SetBufferLength sets the buffer length
func (f *FilePlayerAPI) SetBufferLength(player *FilePlayerC, bufferLen float32) {
	if player != nil {
		C.sound_fileplayer_setBufferLength(f.snd, player.ptr, C.float(bufferLen))
	}
}

// Play plays the file
func (f *FilePlayerAPI) Play(player *FilePlayerC, repeat int) bool {
	if player == nil {
		return false
	}
	return C.sound_fileplayer_play(f.snd, player.ptr, C.int(repeat)) != 0
}

// IsPlaying returns whether the player is playing
func (f *FilePlayerAPI) IsPlaying(player *FilePlayerC) bool {
	if player == nil {
		return false
	}
	return C.sound_fileplayer_isPlaying(f.snd, player.ptr) != 0
}

// Pause pauses the player
func (f *FilePlayerAPI) Pause(player *FilePlayerC) {
	if player != nil {
		C.sound_fileplayer_pause(f.snd, player.ptr)
	}
}

// Stop stops the player
func (f *FilePlayerAPI) Stop(player *FilePlayerC) {
	if player != nil {
		C.sound_fileplayer_stop(f.snd, player.ptr)
	}
}

// SetVolume sets the volume
func (f *FilePlayerAPI) SetVolume(player *FilePlayerC, left, right float32) {
	if player != nil {
		C.sound_fileplayer_setVolume(f.snd, player.ptr, C.float(left), C.float(right))
	}
}

// GetVolume returns the volume
func (f *FilePlayerAPI) GetVolume(player *FilePlayerC) (left, right float32) {
	if player == nil {
		return 0, 0
	}
	var l, r C.float
	C.sound_fileplayer_getVolume(f.snd, player.ptr, &l, &r)
	return float32(l), float32(r)
}

// GetLength returns the length in seconds
func (f *FilePlayerAPI) GetLength(player *FilePlayerC) float32 {
	if player == nil {
		return 0
	}
	return float32(C.sound_fileplayer_getLength(f.snd, player.ptr))
}

// SetOffset sets the playback offset
func (f *FilePlayerAPI) SetOffset(player *FilePlayerC, offset float32) {
	if player != nil {
		C.sound_fileplayer_setOffset(f.snd, player.ptr, C.float(offset))
	}
}

// GetOffset returns the playback offset
func (f *FilePlayerAPI) GetOffset(player *FilePlayerC) float32 {
	if player == nil {
		return 0
	}
	return float32(C.sound_fileplayer_getOffset(f.snd, player.ptr))
}

// SetRate sets the playback rate
func (f *FilePlayerAPI) SetRate(player *FilePlayerC, rate float32) {
	if player != nil {
		C.sound_fileplayer_setRate(f.snd, player.ptr, C.float(rate))
	}
}

// GetRate returns the playback rate
func (f *FilePlayerAPI) GetRate(player *FilePlayerC) float32 {
	if player == nil {
		return 0
	}
	return float32(C.sound_fileplayer_getRate(f.snd, player.ptr))
}

// SetLoopRange sets the loop range
func (f *FilePlayerAPI) SetLoopRange(player *FilePlayerC, start, end float32) {
	if player != nil {
		C.sound_fileplayer_setLoopRange(f.snd, player.ptr, C.float(start), C.float(end))
	}
}

// DidUnderrun returns whether an underrun occurred
func (f *FilePlayerAPI) DidUnderrun(player *FilePlayerC) bool {
	if player == nil {
		return false
	}
	return C.sound_fileplayer_didUnderrun(f.snd, player.ptr) != 0
}

// SetStopOnUnderrun sets whether to stop on underrun
func (f *FilePlayerAPI) SetStopOnUnderrun(player *FilePlayerC, stop bool) {
	if player != nil {
		flag := 0
		if stop {
			flag = 1
		}
		C.sound_fileplayer_setStopOnUnderrun(f.snd, player.ptr, C.int(flag))
	}
}

// SampleAPI wraps sample functions
type SampleAPI struct {
	snd *C.struct_playdate_sound
}

// NewSampleBuffer creates a new sample buffer
func (sa *SampleAPI) NewSampleBuffer(byteCount int) *AudioSample {
	ptr := C.sound_sample_newSampleBuffer(sa.snd, C.int(byteCount))
	if ptr == nil {
		return nil
	}
	return &AudioSample{ptr: ptr}
}

// Load loads a sample from a file
func (sa *SampleAPI) Load(path string) *AudioSample {
	cpath := cString(path)
	defer freeCString(cpath)

	ptr := C.sound_sample_load(sa.snd, cpath)
	if ptr == nil {
		return nil
	}
	return &AudioSample{ptr: ptr}
}

// LoadIntoSample loads a sample into an existing buffer
func (sa *SampleAPI) LoadIntoSample(sample *AudioSample, path string) error {
	if sample == nil {
		return errors.New("sample is nil")
	}
	cpath := cString(path)
	defer freeCString(cpath)

	if C.sound_sample_loadIntoSample(sa.snd, sample.ptr, cpath) == 0 {
		return errors.New("failed to load sample")
	}
	return nil
}

// FreeSample frees a sample
func (sa *SampleAPI) FreeSample(sample *AudioSample) {
	if sample != nil && sample.ptr != nil {
		C.sound_sample_freeSample(sa.snd, sample.ptr)
		sample.ptr = nil
	}
}

// GetLength returns the sample length in seconds
func (sa *SampleAPI) GetLength(sample *AudioSample) float32 {
	if sample == nil {
		return 0
	}
	return float32(C.sound_sample_getLength(sa.snd, sample.ptr))
}

// Decompress decompresses a sample
func (sa *SampleAPI) Decompress(sample *AudioSample) bool {
	if sample == nil {
		return false
	}
	return C.sound_sample_decompress(sa.snd, sample.ptr) != 0
}

// SamplePlayerAPI wraps sample player functions
type SamplePlayerAPI struct {
	snd *C.struct_playdate_sound
}

// NewPlayer creates a new sample player
func (sp *SamplePlayerAPI) NewPlayer() *SamplePlayer {
	ptr := C.sound_sampleplayer_newPlayer(sp.snd)
	if ptr == nil {
		return nil
	}
	return &SamplePlayer{ptr: ptr}
}

// FreePlayer frees a sample player
func (sp *SamplePlayerAPI) FreePlayer(player *SamplePlayer) {
	if player != nil && player.ptr != nil {
		C.sound_sampleplayer_freePlayer(sp.snd, player.ptr)
		player.ptr = nil
	}
}

// SetSample sets the sample for a player
func (sp *SamplePlayerAPI) SetSample(player *SamplePlayer, sample *AudioSample) {
	if player == nil {
		return
	}
	var s *C.AudioSample
	if sample != nil {
		s = sample.ptr
	}
	C.sound_sampleplayer_setSample(sp.snd, player.ptr, s)
}

// Play plays the sample
func (sp *SamplePlayerAPI) Play(player *SamplePlayer, repeat int, rate float32) bool {
	if player == nil {
		return false
	}
	return C.sound_sampleplayer_play(sp.snd, player.ptr, C.int(repeat), C.float(rate)) != 0
}

// IsPlaying returns whether the player is playing
func (sp *SamplePlayerAPI) IsPlaying(player *SamplePlayer) bool {
	if player == nil {
		return false
	}
	return C.sound_sampleplayer_isPlaying(sp.snd, player.ptr) != 0
}

// Stop stops the player
func (sp *SamplePlayerAPI) Stop(player *SamplePlayer) {
	if player != nil {
		C.sound_sampleplayer_stop(sp.snd, player.ptr)
	}
}

// SetVolume sets the volume
func (sp *SamplePlayerAPI) SetVolume(player *SamplePlayer, left, right float32) {
	if player != nil {
		C.sound_sampleplayer_setVolume(sp.snd, player.ptr, C.float(left), C.float(right))
	}
}

// GetVolume returns the volume
func (sp *SamplePlayerAPI) GetVolume(player *SamplePlayer) (left, right float32) {
	if player == nil {
		return 0, 0
	}
	var l, r C.float
	C.sound_sampleplayer_getVolume(sp.snd, player.ptr, &l, &r)
	return float32(l), float32(r)
}

// GetLength returns the length
func (sp *SamplePlayerAPI) GetLength(player *SamplePlayer) float32 {
	if player == nil {
		return 0
	}
	return float32(C.sound_sampleplayer_getLength(sp.snd, player.ptr))
}

// SetOffset sets the offset
func (sp *SamplePlayerAPI) SetOffset(player *SamplePlayer, offset float32) {
	if player != nil {
		C.sound_sampleplayer_setOffset(sp.snd, player.ptr, C.float(offset))
	}
}

// GetOffset returns the offset
func (sp *SamplePlayerAPI) GetOffset(player *SamplePlayer) float32 {
	if player == nil {
		return 0
	}
	return float32(C.sound_sampleplayer_getOffset(sp.snd, player.ptr))
}

// SetRate sets the rate
func (sp *SamplePlayerAPI) SetRate(player *SamplePlayer, rate float32) {
	if player != nil {
		C.sound_sampleplayer_setRate(sp.snd, player.ptr, C.float(rate))
	}
}

// GetRate returns the rate
func (sp *SamplePlayerAPI) GetRate(player *SamplePlayer) float32 {
	if player == nil {
		return 0
	}
	return float32(C.sound_sampleplayer_getRate(sp.snd, player.ptr))
}

// SetPlayRange sets the play range
func (sp *SamplePlayerAPI) SetPlayRange(player *SamplePlayer, start, end int) {
	if player != nil {
		C.sound_sampleplayer_setPlayRange(sp.snd, player.ptr, C.int(start), C.int(end))
	}
}

// SetPaused sets the paused state
func (sp *SamplePlayerAPI) SetPaused(player *SamplePlayer, paused bool) {
	if player != nil {
		flag := 0
		if paused {
			flag = 1
		}
		C.sound_sampleplayer_setPaused(sp.snd, player.ptr, C.int(flag))
	}
}

// SynthAPI wraps synth functions
type SynthAPI struct {
	snd *C.struct_playdate_sound
}

// NewSynth creates a new synth
func (sy *SynthAPI) NewSynth() *PDSynth {
	ptr := C.sound_synth_newSynth(sy.snd)
	if ptr == nil {
		return nil
	}
	return &PDSynth{ptr: ptr}
}

// FreeSynth frees a synth
func (sy *SynthAPI) FreeSynth(synth *PDSynth) {
	if synth != nil && synth.ptr != nil {
		C.sound_synth_freeSynth(sy.snd, synth.ptr)
		synth.ptr = nil
	}
}

// SetWaveform sets the waveform
func (sy *SynthAPI) SetWaveform(synth *PDSynth, wave SoundWaveform) {
	if synth != nil {
		C.sound_synth_setWaveform(sy.snd, synth.ptr, C.SoundWaveform(wave))
	}
}

// SetSample sets the sample
func (sy *SynthAPI) SetSample(synth *PDSynth, sample *AudioSample, sustainStart, sustainEnd uint32) {
	if synth == nil {
		return
	}
	var s *C.AudioSample
	if sample != nil {
		s = sample.ptr
	}
	C.sound_synth_setSample(sy.snd, synth.ptr, s, C.uint32_t(sustainStart), C.uint32_t(sustainEnd))
}

// SetAttackTime sets the attack time
func (sy *SynthAPI) SetAttackTime(synth *PDSynth, attack float32) {
	if synth != nil {
		C.sound_synth_setAttackTime(sy.snd, synth.ptr, C.float(attack))
	}
}

// SetDecayTime sets the decay time
func (sy *SynthAPI) SetDecayTime(synth *PDSynth, decay float32) {
	if synth != nil {
		C.sound_synth_setDecayTime(sy.snd, synth.ptr, C.float(decay))
	}
}

// SetSustainLevel sets the sustain level
func (sy *SynthAPI) SetSustainLevel(synth *PDSynth, sustain float32) {
	if synth != nil {
		C.sound_synth_setSustainLevel(sy.snd, synth.ptr, C.float(sustain))
	}
}

// SetReleaseTime sets the release time
func (sy *SynthAPI) SetReleaseTime(synth *PDSynth, release float32) {
	if synth != nil {
		C.sound_synth_setReleaseTime(sy.snd, synth.ptr, C.float(release))
	}
}

// SetTranspose sets the transpose
func (sy *SynthAPI) SetTranspose(synth *PDSynth, halfSteps float32) {
	if synth != nil {
		C.sound_synth_setTranspose(sy.snd, synth.ptr, C.float(halfSteps))
	}
}

// PlayNote plays a note
func (sy *SynthAPI) PlayNote(synth *PDSynth, freq, vel, length float32, when uint32) {
	if synth != nil {
		C.sound_synth_playNote(sy.snd, synth.ptr, C.float(freq), C.float(vel), C.float(length), C.uint32_t(when))
	}
}

// PlayMIDINote plays a MIDI note
func (sy *SynthAPI) PlayMIDINote(synth *PDSynth, note MIDINote, vel, length float32, when uint32) {
	if synth != nil {
		C.sound_synth_playMIDINote(sy.snd, synth.ptr, C.MIDINote(note), C.float(vel), C.float(length), C.uint32_t(when))
	}
}

// NoteOff releases a note
func (sy *SynthAPI) NoteOff(synth *PDSynth, when uint32) {
	if synth != nil {
		C.sound_synth_noteOff(sy.snd, synth.ptr, C.uint32_t(when))
	}
}

// Stop stops the synth
func (sy *SynthAPI) Stop(synth *PDSynth) {
	if synth != nil {
		C.sound_synth_stop(sy.snd, synth.ptr)
	}
}

// SetVolume sets the volume
func (sy *SynthAPI) SetVolume(synth *PDSynth, left, right float32) {
	if synth != nil {
		C.sound_synth_setVolume(sy.snd, synth.ptr, C.float(left), C.float(right))
	}
}

// GetVolume returns the volume
func (sy *SynthAPI) GetVolume(synth *PDSynth) (left, right float32) {
	if synth == nil {
		return 0, 0
	}
	var l, r C.float
	C.sound_synth_getVolume(sy.snd, synth.ptr, &l, &r)
	return float32(l), float32(r)
}

// IsPlaying returns whether the synth is playing
func (sy *SynthAPI) IsPlaying(synth *PDSynth) bool {
	if synth == nil {
		return false
	}
	return C.sound_synth_isPlaying(sy.snd, synth.ptr) != 0
}

// Copy copies a synth
func (sy *SynthAPI) Copy(synth *PDSynth) *PDSynth {
	if synth == nil {
		return nil
	}
	ptr := C.sound_synth_copy(sy.snd, synth.ptr)
	if ptr == nil {
		return nil
	}
	return &PDSynth{ptr: ptr}
}

// LFOAPI wraps LFO functions
type LFOAPI struct {
	snd *C.struct_playdate_sound
}

// NewLFO creates a new LFO
func (l *LFOAPI) NewLFO(lfoType LFOType) *PDSynthLFO {
	ptr := C.sound_lfo_newLFO(l.snd, C.LFOType(lfoType))
	if ptr == nil {
		return nil
	}
	return &PDSynthLFO{ptr: ptr}
}

// FreeLFO frees an LFO
func (l *LFOAPI) FreeLFO(lfo *PDSynthLFO) {
	if lfo != nil && lfo.ptr != nil {
		C.sound_lfo_freeLFO(l.snd, lfo.ptr)
		lfo.ptr = nil
	}
}

// SetType sets the LFO type
func (l *LFOAPI) SetType(lfo *PDSynthLFO, lfoType LFOType) {
	if lfo != nil {
		C.sound_lfo_setType(l.snd, lfo.ptr, C.LFOType(lfoType))
	}
}

// SetRate sets the rate
func (l *LFOAPI) SetRate(lfo *PDSynthLFO, rate float32) {
	if lfo != nil {
		C.sound_lfo_setRate(l.snd, lfo.ptr, C.float(rate))
	}
}

// SetPhase sets the phase
func (l *LFOAPI) SetPhase(lfo *PDSynthLFO, phase float32) {
	if lfo != nil {
		C.sound_lfo_setPhase(l.snd, lfo.ptr, C.float(phase))
	}
}

// SetCenter sets the center
func (l *LFOAPI) SetCenter(lfo *PDSynthLFO, center float32) {
	if lfo != nil {
		C.sound_lfo_setCenter(l.snd, lfo.ptr, C.float(center))
	}
}

// SetDepth sets the depth
func (l *LFOAPI) SetDepth(lfo *PDSynthLFO, depth float32) {
	if lfo != nil {
		C.sound_lfo_setDepth(l.snd, lfo.ptr, C.float(depth))
	}
}

// SetArpeggiation sets the arpeggiation
func (l *LFOAPI) SetArpeggiation(lfo *PDSynthLFO, steps []float32) {
	if lfo == nil || len(steps) == 0 {
		return
	}
	C.sound_lfo_setArpeggiation(l.snd, lfo.ptr, C.int(len(steps)), (*C.float)(&steps[0]))
}

// SetDelay sets the delay
func (l *LFOAPI) SetDelay(lfo *PDSynthLFO, holdoff, ramptime float32) {
	if lfo != nil {
		C.sound_lfo_setDelay(l.snd, lfo.ptr, C.float(holdoff), C.float(ramptime))
	}
}

// SetRetrigger sets retrigger
func (l *LFOAPI) SetRetrigger(lfo *PDSynthLFO, retrigger bool) {
	if lfo != nil {
		flag := 0
		if retrigger {
			flag = 1
		}
		C.sound_lfo_setRetrigger(l.snd, lfo.ptr, C.int(flag))
	}
}

// GetValue returns the current value
func (l *LFOAPI) GetValue(lfo *PDSynthLFO) float32 {
	if lfo == nil {
		return 0
	}
	return float32(C.sound_lfo_getValue(l.snd, lfo.ptr))
}

// SetGlobal sets whether the LFO is global
func (l *LFOAPI) SetGlobal(lfo *PDSynthLFO, global bool) {
	if lfo != nil {
		flag := 0
		if global {
			flag = 1
		}
		C.sound_lfo_setGlobal(l.snd, lfo.ptr, C.int(flag))
	}
}

// SetStartPhase sets the start phase
func (l *LFOAPI) SetStartPhase(lfo *PDSynthLFO, phase float32) {
	if lfo != nil {
		C.sound_lfo_setStartPhase(l.snd, lfo.ptr, C.float(phase))
	}
}

// EnvelopeAPI wraps envelope functions
type EnvelopeAPI struct {
	snd *C.struct_playdate_sound
}

// NewEnvelope creates a new envelope
func (e *EnvelopeAPI) NewEnvelope(attack, decay, sustain, release float32) *PDSynthEnvelope {
	ptr := C.sound_envelope_newEnvelope(e.snd, C.float(attack), C.float(decay), C.float(sustain), C.float(release))
	if ptr == nil {
		return nil
	}
	return &PDSynthEnvelope{ptr: ptr}
}

// FreeEnvelope frees an envelope
func (e *EnvelopeAPI) FreeEnvelope(env *PDSynthEnvelope) {
	if env != nil && env.ptr != nil {
		C.sound_envelope_freeEnvelope(e.snd, env.ptr)
		env.ptr = nil
	}
}

// SetAttack sets the attack
func (e *EnvelopeAPI) SetAttack(env *PDSynthEnvelope, attack float32) {
	if env != nil {
		C.sound_envelope_setAttack(e.snd, env.ptr, C.float(attack))
	}
}

// SetDecay sets the decay
func (e *EnvelopeAPI) SetDecay(env *PDSynthEnvelope, decay float32) {
	if env != nil {
		C.sound_envelope_setDecay(e.snd, env.ptr, C.float(decay))
	}
}

// SetSustain sets the sustain
func (e *EnvelopeAPI) SetSustain(env *PDSynthEnvelope, sustain float32) {
	if env != nil {
		C.sound_envelope_setSustain(e.snd, env.ptr, C.float(sustain))
	}
}

// SetRelease sets the release
func (e *EnvelopeAPI) SetRelease(env *PDSynthEnvelope, release float32) {
	if env != nil {
		C.sound_envelope_setRelease(e.snd, env.ptr, C.float(release))
	}
}

// SetLegato sets legato
func (e *EnvelopeAPI) SetLegato(env *PDSynthEnvelope, legato bool) {
	if env != nil {
		flag := 0
		if legato {
			flag = 1
		}
		C.sound_envelope_setLegato(e.snd, env.ptr, C.int(flag))
	}
}

// SetRetrigger sets retrigger
func (e *EnvelopeAPI) SetRetrigger(env *PDSynthEnvelope, retrigger bool) {
	if env != nil {
		flag := 0
		if retrigger {
			flag = 1
		}
		C.sound_envelope_setRetrigger(e.snd, env.ptr, C.int(flag))
	}
}

// GetValue returns the current value
func (e *EnvelopeAPI) GetValue(env *PDSynthEnvelope) float32 {
	if env == nil {
		return 0
	}
	return float32(C.sound_envelope_getValue(e.snd, env.ptr))
}

// SetCurvature sets the curvature
func (e *EnvelopeAPI) SetCurvature(env *PDSynthEnvelope, amount float32) {
	if env != nil {
		C.sound_envelope_setCurvature(e.snd, env.ptr, C.float(amount))
	}
}

// SetVelocitySensitivity sets velocity sensitivity
func (e *EnvelopeAPI) SetVelocitySensitivity(env *PDSynthEnvelope, velsens float32) {
	if env != nil {
		C.sound_envelope_setVelocitySensitivity(e.snd, env.ptr, C.float(velsens))
	}
}

// SetRateScaling sets rate scaling
func (e *EnvelopeAPI) SetRateScaling(env *PDSynthEnvelope, scaling float32, start, end MIDINote) {
	if env != nil {
		C.sound_envelope_setRateScaling(e.snd, env.ptr, C.float(scaling), C.MIDINote(start), C.MIDINote(end))
	}
}

// SequenceAPI wraps sequence functions
type SequenceAPI struct {
	snd *C.struct_playdate_sound
}

// NewSequence creates a new sequence
func (s *SequenceAPI) NewSequence() *SoundSequence {
	ptr := C.sound_sequence_newSequence(s.snd)
	if ptr == nil {
		return nil
	}
	return &SoundSequence{ptr: ptr}
}

// FreeSequence frees a sequence
func (s *SequenceAPI) FreeSequence(seq *SoundSequence) {
	if seq != nil && seq.ptr != nil {
		C.sound_sequence_freeSequence(s.snd, seq.ptr)
		seq.ptr = nil
	}
}

// LoadMIDIFile loads a MIDI file
func (s *SequenceAPI) LoadMIDIFile(seq *SoundSequence, path string) error {
	if seq == nil {
		return errors.New("sequence is nil")
	}
	cpath := cString(path)
	defer freeCString(cpath)

	if C.sound_sequence_loadMIDIFile(s.snd, seq.ptr, cpath) == 0 {
		return errors.New("failed to load MIDI file")
	}
	return nil
}

// GetTime returns the sequence time
func (s *SequenceAPI) GetTime(seq *SoundSequence) uint32 {
	if seq == nil {
		return 0
	}
	return uint32(C.sound_sequence_getTime(s.snd, seq.ptr))
}

// SetTime sets the sequence time
func (s *SequenceAPI) SetTime(seq *SoundSequence, time uint32) {
	if seq != nil {
		C.sound_sequence_setTime(s.snd, seq.ptr, C.uint32_t(time))
	}
}

// SetLoops sets the loop range
func (s *SequenceAPI) SetLoops(seq *SoundSequence, loopStart, loopEnd, loops int) {
	if seq != nil {
		C.sound_sequence_setLoops(s.snd, seq.ptr, C.int(loopStart), C.int(loopEnd), C.int(loops))
	}
}

// SetTempo sets the tempo
func (s *SequenceAPI) SetTempo(seq *SoundSequence, stepsPerSecond float32) {
	if seq != nil {
		C.sound_sequence_setTempo(s.snd, seq.ptr, C.float(stepsPerSecond))
	}
}

// GetTempo returns the tempo
func (s *SequenceAPI) GetTempo(seq *SoundSequence) float32 {
	if seq == nil {
		return 0
	}
	return float32(C.sound_sequence_getTempo(s.snd, seq.ptr))
}

// GetTrackCount returns the track count
func (s *SequenceAPI) GetTrackCount(seq *SoundSequence) int {
	if seq == nil {
		return 0
	}
	return int(C.sound_sequence_getTrackCount(s.snd, seq.ptr))
}

// AddTrack adds a track
func (s *SequenceAPI) AddTrack(seq *SoundSequence) *SequenceTrack {
	if seq == nil {
		return nil
	}
	ptr := C.sound_sequence_addTrack(s.snd, seq.ptr)
	if ptr == nil {
		return nil
	}
	return &SequenceTrack{ptr: ptr}
}

// GetTrackAtIndex returns a track at an index
func (s *SequenceAPI) GetTrackAtIndex(seq *SoundSequence, index uint) *SequenceTrack {
	if seq == nil {
		return nil
	}
	ptr := C.sound_sequence_getTrackAtIndex(s.snd, seq.ptr, C.uint(index))
	if ptr == nil {
		return nil
	}
	return &SequenceTrack{ptr: ptr}
}

// IsPlaying returns whether the sequence is playing
func (s *SequenceAPI) IsPlaying(seq *SoundSequence) bool {
	if seq == nil {
		return false
	}
	return C.sound_sequence_isPlaying(s.snd, seq.ptr) != 0
}

// GetLength returns the sequence length
func (s *SequenceAPI) GetLength(seq *SoundSequence) uint32 {
	if seq == nil {
		return 0
	}
	return uint32(C.sound_sequence_getLength(s.snd, seq.ptr))
}

// Play plays the sequence
func (s *SequenceAPI) Play(seq *SoundSequence) {
	if seq != nil {
		C.sound_sequence_play(s.snd, seq.ptr, nil, nil)
	}
}

// Stop stops the sequence
func (s *SequenceAPI) Stop(seq *SoundSequence) {
	if seq != nil {
		C.sound_sequence_stop(s.snd, seq.ptr)
	}
}

// AllNotesOff turns off all notes
func (s *SequenceAPI) AllNotesOff(seq *SoundSequence) {
	if seq != nil {
		C.sound_sequence_allNotesOff(s.snd, seq.ptr)
	}
}

// GetCurrentStep returns the current step in the sequence
func (s *SequenceAPI) GetCurrentStep(seq *SoundSequence) int {
	if seq == nil {
		return 0
	}
	return int(C.sound_sequence_getCurrentStep(s.snd, seq.ptr, nil))
}

// SetCurrentStep sets the current step in the sequence
func (s *SequenceAPI) SetCurrentStep(seq *SoundSequence, step, timeOffset int, playNotes bool) {
	if seq != nil {
		pn := 0
		if playNotes {
			pn = 1
		}
		C.sound_sequence_setCurrentStep(s.snd, seq.ptr, C.int(step), C.int(timeOffset), C.int(pn))
	}
}

// TrackAPI wraps track functions
type TrackAPI struct {
	snd *C.struct_playdate_sound
}

// NewTrack creates a new track
func (t *TrackAPI) NewTrack() *SequenceTrack {
	ptr := C.sound_track_newTrack(t.snd)
	if ptr == nil {
		return nil
	}
	return &SequenceTrack{ptr: ptr}
}

// FreeTrack frees a track
func (t *TrackAPI) FreeTrack(track *SequenceTrack) {
	if track != nil && track.ptr != nil {
		C.sound_track_freeTrack(t.snd, track.ptr)
		track.ptr = nil
	}
}

// SetInstrument sets the instrument for a track
func (t *TrackAPI) SetInstrument(track *SequenceTrack, inst *PDSynthInstrument) {
	if track == nil {
		return
	}
	var i *C.PDSynthInstrument
	if inst != nil {
		i = inst.ptr
	}
	C.sound_track_setInstrument(t.snd, track.ptr, i)
}

// GetInstrument returns the track's instrument
func (t *TrackAPI) GetInstrument(track *SequenceTrack) *PDSynthInstrument {
	if track == nil {
		return nil
	}
	ptr := C.sound_track_getInstrument(t.snd, track.ptr)
	if ptr == nil {
		return nil
	}
	return &PDSynthInstrument{ptr: ptr}
}

// AddNoteEvent adds a note event
func (t *TrackAPI) AddNoteEvent(track *SequenceTrack, step, length uint32, note MIDINote, velocity float32) {
	if track != nil {
		C.sound_track_addNoteEvent(t.snd, track.ptr, C.uint32_t(step), C.uint32_t(length), C.MIDINote(note), C.float(velocity))
	}
}

// RemoveNoteEvent removes a note event
func (t *TrackAPI) RemoveNoteEvent(track *SequenceTrack, step uint32, note MIDINote) {
	if track != nil {
		C.sound_track_removeNoteEvent(t.snd, track.ptr, C.uint32_t(step), C.MIDINote(note))
	}
}

// ClearNotes clears all notes
func (t *TrackAPI) ClearNotes(track *SequenceTrack) {
	if track != nil {
		C.sound_track_clearNotes(t.snd, track.ptr)
	}
}

// GetLength returns the track length
func (t *TrackAPI) GetLength(track *SequenceTrack) uint32 {
	if track == nil {
		return 0
	}
	return uint32(C.sound_track_getLength(t.snd, track.ptr))
}

// SetMuted sets whether the track is muted
func (t *TrackAPI) SetMuted(track *SequenceTrack, muted bool) {
	if track != nil {
		flag := 0
		if muted {
			flag = 1
		}
		C.sound_track_setMuted(t.snd, track.ptr, C.int(flag))
	}
}

// GetPolyphony returns the maximum polyphony for the track
func (t *TrackAPI) GetPolyphony(track *SequenceTrack) int {
	if track == nil {
		return 0
	}
	return int(C.sound_track_getPolyphony(t.snd, track.ptr))
}

// GetIndexForStep returns the internal index for the step
func (t *TrackAPI) GetIndexForStep(track *SequenceTrack, step uint32) int {
	if track == nil {
		return 0
	}
	return int(C.sound_track_getIndexForStep(t.snd, track.ptr, C.uint32_t(step)))
}

// GetNoteAtIndex returns note information at the given index
// Returns step, length, note, velocity, and whether the note exists
func (t *TrackAPI) GetNoteAtIndex(track *SequenceTrack, index int) (step, length uint32, note MIDINote, velocity float32, ok bool) {
	if track == nil {
		return 0, 0, 0, 0, false
	}
	var s, l C.uint32_t
	var n C.MIDINote
	var v C.float
	result := C.sound_track_getNoteAtIndex(t.snd, track.ptr, C.int(index), &s, &l, &n, &v)
	if result == 0 {
		return 0, 0, 0, 0, false
	}
	return uint32(s), uint32(l), MIDINote(n), float32(v), true
}

// InstrumentAPI wraps instrument functions
type InstrumentAPI struct {
	snd *C.struct_playdate_sound
}

// NewInstrument creates a new instrument
func (i *InstrumentAPI) NewInstrument() *PDSynthInstrument {
	ptr := C.sound_instrument_newInstrument(i.snd)
	if ptr == nil {
		return nil
	}
	return &PDSynthInstrument{ptr: ptr}
}

// FreeInstrument frees an instrument
func (i *InstrumentAPI) FreeInstrument(inst *PDSynthInstrument) {
	if inst != nil && inst.ptr != nil {
		C.sound_instrument_freeInstrument(i.snd, inst.ptr)
		inst.ptr = nil
	}
}

// AddVoice adds a voice to an instrument
func (i *InstrumentAPI) AddVoice(inst *PDSynthInstrument, synth *PDSynth, rangeStart, rangeEnd MIDINote, transpose float32) bool {
	if inst == nil || synth == nil {
		return false
	}
	return C.sound_instrument_addVoice(i.snd, inst.ptr, synth.ptr, C.MIDINote(rangeStart), C.MIDINote(rangeEnd), C.float(transpose)) != 0
}

// PlayNote plays a note
func (i *InstrumentAPI) PlayNote(inst *PDSynthInstrument, frequency, vel, length float32, when uint32) *PDSynth {
	if inst == nil {
		return nil
	}
	ptr := C.sound_instrument_playNote(i.snd, inst.ptr, C.float(frequency), C.float(vel), C.float(length), C.uint32_t(when))
	if ptr == nil {
		return nil
	}
	return &PDSynth{ptr: ptr}
}

// PlayMIDINote plays a MIDI note
func (i *InstrumentAPI) PlayMIDINote(inst *PDSynthInstrument, note MIDINote, vel, length float32, when uint32) *PDSynth {
	if inst == nil {
		return nil
	}
	ptr := C.sound_instrument_playMIDINote(i.snd, inst.ptr, C.MIDINote(note), C.float(vel), C.float(length), C.uint32_t(when))
	if ptr == nil {
		return nil
	}
	return &PDSynth{ptr: ptr}
}

// NoteOff turns off a note
func (i *InstrumentAPI) NoteOff(inst *PDSynthInstrument, note MIDINote, when uint32) {
	if inst != nil {
		C.sound_instrument_noteOff(i.snd, inst.ptr, C.MIDINote(note), C.uint32_t(when))
	}
}

// AllNotesOff turns off all notes
func (i *InstrumentAPI) AllNotesOff(inst *PDSynthInstrument, when uint32) {
	if inst != nil {
		C.sound_instrument_allNotesOff(i.snd, inst.ptr, C.uint32_t(when))
	}
}

// SetVolume sets the volume
func (i *InstrumentAPI) SetVolume(inst *PDSynthInstrument, left, right float32) {
	if inst != nil {
		C.sound_instrument_setVolume(i.snd, inst.ptr, C.float(left), C.float(right))
	}
}

// GetVolume returns the volume
func (i *InstrumentAPI) GetVolume(inst *PDSynthInstrument) (left, right float32) {
	if inst == nil {
		return 0, 0
	}
	var l, r C.float
	C.sound_instrument_getVolume(i.snd, inst.ptr, &l, &r)
	return float32(l), float32(r)
}

// ActiveVoiceCount returns the active voice count
func (i *InstrumentAPI) ActiveVoiceCount(inst *PDSynthInstrument) int {
	if inst == nil {
		return 0
	}
	return int(C.sound_instrument_activeVoiceCount(i.snd, inst.ptr))
}
