// pd_cgo.c - CGO runtime support for Playdate
// This file provides C wrappers that Go code calls via CGO

#include <stdint.h>
#include <stddef.h>
#include "pd_api.h"

// ============== Global PlaydateAPI pointer ==============

PlaydateAPI* pd = NULL;

#ifdef TARGET_PLAYDATE
// ============== Interrupt functions for TinyGo cortexm runtime ==============
// Only needed for device builds (ARM Cortex-M)

void EnableInterrupts(void) {
    __asm__ volatile ("cpsie i" ::: "memory");
}

void DisableInterrupts(void) {
    __asm__ volatile ("cpsid i" ::: "memory");
}

// ============== TinyGo Runtime Support ==============
// Only needed for TinyGo device builds

void* runtime__cgo_pd_realloc(void* ptr, size_t size) __asm__("runtime._cgo_pd_realloc");
void* runtime__cgo_pd_realloc(void* ptr, size_t size) {
    return pd ? pd->system->realloc(ptr, size) : NULL;
}

uint32_t runtime__cgo_pd_getCurrentTimeMS(void) __asm__("runtime._cgo_pd_getCurrentTimeMS");
uint32_t runtime__cgo_pd_getCurrentTimeMS(void) {
    return pd ? pd->system->getCurrentTimeMilliseconds() : 0;
}

void runtime__cgo_pd_logToConsole(const char* msg) __asm__("runtime._cgo_pd_logToConsole");
void runtime__cgo_pd_logToConsole(const char* msg) {
    if (pd) pd->system->logToConsole("%s", msg);
}
#endif // TARGET_PLAYDATE

// ============== System API ==============

void pd_sys_log(const char* msg) {
    if (pd) pd->system->logToConsole("%s", msg);
}

void pd_sys_error(const char* msg) {
    if (pd) pd->system->error("%s", msg);
}

void pd_sys_drawFPS(int x, int y) {
    if (pd) pd->system->drawFPS(x, y);
}

uint32_t pd_sys_getCurrentTimeMS(void) {
    return pd ? pd->system->getCurrentTimeMilliseconds() : 0;
}

uint32_t pd_sys_getSecondsSinceEpoch(unsigned int* ms) {
    if (pd) return pd->system->getSecondsSinceEpoch(ms);
    return 0;
}

void pd_sys_getButtonState(uint32_t* current, uint32_t* pushed, uint32_t* released) {
    if (pd) {
        PDButtons c, p, r;
        pd->system->getButtonState(&c, &p, &r);
        if (current) *current = c;
        if (pushed) *pushed = p;
        if (released) *released = r;
    }
}

void pd_sys_setPeripheralsEnabled(uint32_t mask) {
    if (pd) pd->system->setPeripheralsEnabled(mask);
}

void pd_sys_getAccelerometer(float* x, float* y, float* z) {
    if (pd) pd->system->getAccelerometer(x, y, z);
}

float pd_sys_getCrankChange(void) {
    return pd ? pd->system->getCrankChange() : 0;
}

float pd_sys_getCrankAngle(void) {
    return pd ? pd->system->getCrankAngle() : 0;
}

int pd_sys_isCrankDocked(void) {
    return pd ? pd->system->isCrankDocked() : 0;
}

int pd_sys_setCrankSoundsDisabled(int disabled) {
    return pd ? pd->system->setCrankSoundsDisabled(disabled) : 0;
}

int pd_sys_getFlipped(void) {
    return pd ? pd->system->getFlipped() : 0;
}

void pd_sys_setAutoLockDisabled(int disabled) {
    if (pd) pd->system->setAutoLockDisabled(disabled);
}

int pd_sys_getLanguage(void) {
    return pd ? (int)pd->system->getLanguage() : 0;
}

float pd_sys_getBatteryPercentage(void) {
    return pd ? pd->system->getBatteryPercentage() : 1.0f;
}

float pd_sys_getBatteryVoltage(void) {
    return pd ? pd->system->getBatteryVoltage() : 4.2f;
}

// Update callback support
extern int pdgo_update_trampoline(void);

static int update_callback_wrapper(void* userdata) {
    (void)userdata;
    return pdgo_update_trampoline();
}

void pd_sys_setUpdateCallback(void) {
    if (pd) pd->system->setUpdateCallback(update_callback_wrapper, NULL);
}

// ============== Display API ==============

int pd_display_getWidth(void) {
    return pd ? pd->display->getWidth() : 400;
}

int pd_display_getHeight(void) {
    return pd ? pd->display->getHeight() : 240;
}

void pd_display_setRefreshRate(float rate) {
    if (pd) pd->display->setRefreshRate(rate);
}

float pd_display_getRefreshRate(void) {
    return pd ? pd->display->getRefreshRate() : 30.0f;
}

void pd_display_setInverted(int inverted) {
    if (pd) pd->display->setInverted(inverted);
}

void pd_display_setScale(uint32_t scale) {
    if (pd) pd->display->setScale(scale);
}

void pd_display_setMosaic(uint32_t x, uint32_t y) {
    if (pd) pd->display->setMosaic(x, y);
}

void pd_display_setFlipped(int x, int y) {
    if (pd) pd->display->setFlipped(x, y);
}

void pd_display_setOffset(int x, int y) {
    if (pd) pd->display->setOffset(x, y);
}

// ============== Graphics API ==============

void pd_gfx_clear(uint32_t color) {
    if (pd) pd->graphics->clear(color);
}

void pd_gfx_setBackgroundColor(int color) {
    if (pd) pd->graphics->setBackgroundColor(color);
}

int pd_gfx_setDrawMode(int mode) {
    return pd ? pd->graphics->setDrawMode(mode) : mode;
}

void pd_gfx_setDrawOffset(int dx, int dy) {
    if (pd) pd->graphics->setDrawOffset(dx, dy);
}

void pd_gfx_setClipRect(int x, int y, int w, int h) {
    if (pd) pd->graphics->setClipRect(x, y, w, h);
}

void pd_gfx_clearClipRect(void) {
    if (pd) pd->graphics->clearClipRect();
}

void pd_gfx_setLineCapStyle(int style) {
    if (pd) pd->graphics->setLineCapStyle(style);
}

void pd_gfx_setFont(void* font) {
    if (pd) pd->graphics->setFont(font);
}

void pd_gfx_setTextTracking(int tracking) {
    if (pd) pd->graphics->setTextTracking(tracking);
}

void pd_gfx_pushContext(void* target) {
    if (pd) pd->graphics->pushContext(target);
}

void pd_gfx_popContext(void) {
    if (pd) pd->graphics->popContext();
}

// Drawing primitives
void pd_gfx_fillRect(int x, int y, int w, int h, uint32_t color) {
    if (pd) pd->graphics->fillRect(x, y, w, h, color);
}

void pd_gfx_drawRect(int x, int y, int w, int h, uint32_t color) {
    if (pd) pd->graphics->drawRect(x, y, w, h, color);
}

void pd_gfx_drawLine(int x1, int y1, int x2, int y2, int width, uint32_t color) {
    if (pd) pd->graphics->drawLine(x1, y1, x2, y2, width, color);
}

void pd_gfx_fillTriangle(int x1, int y1, int x2, int y2, int x3, int y3, uint32_t color) {
    if (pd) pd->graphics->fillTriangle(x1, y1, x2, y2, x3, y3, color);
}

void pd_gfx_drawEllipse(int x, int y, int w, int h, int lineWidth, float startAngle, float endAngle, uint32_t color) {
    if (pd) pd->graphics->drawEllipse(x, y, w, h, lineWidth, startAngle, endAngle, color);
}

void pd_gfx_fillEllipse(int x, int y, int w, int h, float startAngle, float endAngle, uint32_t color) {
    if (pd) pd->graphics->fillEllipse(x, y, w, h, startAngle, endAngle, color);
}

// Text
int pd_gfx_drawText(const char* text, int len, int encoding, int x, int y) {
    if (pd) return pd->graphics->drawText(text, len, encoding, x, y);
    return 0;
}

int pd_gfx_getTextWidth(void* font, const char* text, int len, int encoding, int tracking) {
    if (pd) return pd->graphics->getTextWidth(font, text, len, encoding, tracking);
    return 0;
}

void* pd_gfx_loadFont(const char* path, const char** err) {
    if (pd) return pd->graphics->loadFont(path, err);
    return NULL;
}

// Bitmap
void* pd_gfx_newBitmap(int w, int h, uint32_t bgcolor) {
    if (pd) return pd->graphics->newBitmap(w, h, bgcolor);
    return NULL;
}

void pd_gfx_freeBitmap(void* bitmap) {
    if (pd) pd->graphics->freeBitmap(bitmap);
}

void* pd_gfx_loadBitmap(const char* path, const char** err) {
    if (pd) return pd->graphics->loadBitmap(path, err);
    return NULL;
}

void* pd_gfx_copyBitmap(void* bitmap) {
    if (pd) return pd->graphics->copyBitmap(bitmap);
    return NULL;
}

void pd_gfx_drawBitmap(void* bitmap, int x, int y, int flip) {
    if (pd) pd->graphics->drawBitmap(bitmap, x, y, flip);
}

void pd_gfx_tileBitmap(void* bitmap, int x, int y, int w, int h, int flip) {
    if (pd) pd->graphics->tileBitmap(bitmap, x, y, w, h, flip);
}

void pd_gfx_drawScaledBitmap(void* bitmap, int x, int y, float xscale, float yscale) {
    if (pd) pd->graphics->drawScaledBitmap(bitmap, x, y, xscale, yscale);
}

void pd_gfx_drawRotatedBitmap(void* bitmap, int x, int y, float rotation, float cx, float cy, float xscale, float yscale) {
    if (pd) pd->graphics->drawRotatedBitmap(bitmap, x, y, rotation, cx, cy, xscale, yscale);
}

void pd_gfx_getBitmapData(void* bitmap, int* w, int* h, int* rowbytes, uint8_t** mask, uint8_t** data) {
    if (pd) pd->graphics->getBitmapData(bitmap, w, h, rowbytes, mask, data);
}

void pd_gfx_clearBitmap(void* bitmap, uint32_t bgcolor) {
    if (pd) pd->graphics->clearBitmap(bitmap, bgcolor);
}

// BitmapTable
void* pd_gfx_newBitmapTable(int count, int w, int h) {
    if (pd) return pd->graphics->newBitmapTable(count, w, h);
    return NULL;
}

void pd_gfx_freeBitmapTable(void* table) {
    if (pd) pd->graphics->freeBitmapTable(table);
}

void* pd_gfx_loadBitmapTable(const char* path, const char** err) {
    if (pd) return pd->graphics->loadBitmapTable(path, err);
    return NULL;
}

void* pd_gfx_getTableBitmap(void* table, int idx) {
    if (pd) return pd->graphics->getTableBitmap(table, idx);
    return NULL;
}

// Frame buffer
uint8_t* pd_gfx_getFrame(void) {
    if (pd) return pd->graphics->getFrame();
    return NULL;
}

uint8_t* pd_gfx_getDisplayFrame(void) {
    if (pd) return pd->graphics->getDisplayFrame();
    return NULL;
}

void pd_gfx_markUpdatedRows(int start, int end) {
    if (pd) pd->graphics->markUpdatedRows(start, end);
}

void pd_gfx_display(void) {
    if (pd) pd->graphics->display();
}

void* pd_gfx_getDisplayBufferBitmap(void) {
    if (pd) return pd->graphics->getDisplayBufferBitmap();
    return NULL;
}

// ============== Sound API ==============

void* pd_sound_getDefaultChannel(void) {
    if (pd) return pd->sound->getDefaultChannel();
    return NULL;
}

int pd_sound_channel_addSource(void* channel, void* source) {
    if (pd) return pd->sound->channel->addSource(channel, source);
    return 0;
}

// Synth
void* pd_sound_synth_new(void) {
    if (pd) return pd->sound->synth->newSynth();
    return NULL;
}

void pd_sound_synth_free(void* synth) {
    if (pd) pd->sound->synth->freeSynth(synth);
}

void pd_sound_synth_setWaveform(void* synth, int waveform) {
    if (pd) pd->sound->synth->setWaveform(synth, waveform);
}

void pd_sound_synth_setAttackTime(void* synth, float attack) {
    if (pd) pd->sound->synth->setAttackTime(synth, attack);
}

void pd_sound_synth_setDecayTime(void* synth, float decay) {
    if (pd) pd->sound->synth->setDecayTime(synth, decay);
}

void pd_sound_synth_setSustainLevel(void* synth, float sustain) {
    if (pd) pd->sound->synth->setSustainLevel(synth, sustain);
}

void pd_sound_synth_setReleaseTime(void* synth, float release) {
    if (pd) pd->sound->synth->setReleaseTime(synth, release);
}

void pd_sound_synth_setTranspose(void* synth, float halfSteps) {
    if (pd) pd->sound->synth->setTranspose(synth, halfSteps);
}

void pd_sound_synth_playNote(void* synth, float freq, float vel, float len, uint32_t when) {
    if (pd) pd->sound->synth->playNote(synth, freq, vel, len, when);
}

void pd_sound_synth_playMIDINote(void* synth, float note, float vel, float len, uint32_t when) {
    if (pd) pd->sound->synth->playMIDINote(synth, note, vel, len, when);
}

void pd_sound_synth_noteOff(void* synth, uint32_t when) {
    if (pd) pd->sound->synth->noteOff(synth, when);
}

void pd_sound_synth_stop(void* synth) {
    if (pd) pd->sound->synth->stop(synth);
}

void pd_sound_synth_setVolume(void* synth, float left, float right) {
    if (pd) pd->sound->synth->setVolume(synth, left, right);
}

void pd_sound_synth_getVolume(void* synth, float* left, float* right) {
    if (pd) pd->sound->synth->getVolume(synth, left, right);
}

int pd_sound_synth_isPlaying(void* synth) {
    if (pd) return pd->sound->synth->isPlaying(synth);
    return 0;
}

void pd_sound_synth_setSample(void* synth, void* sample, uint32_t sustainStart, uint32_t sustainEnd) {
    if (pd) pd->sound->synth->setSample(synth, sample, sustainStart, sustainEnd);
}

void* pd_sound_synth_copy(void* synth) {
    if (pd) return pd->sound->synth->copy(synth);
    return NULL;
}

// Sample
void* pd_sound_sample_new(int len) {
    if (pd) return pd->sound->sample->newSampleBuffer(len);
    return NULL;
}

void* pd_sound_sample_load(const char* path) {
    if (pd) return pd->sound->sample->load(path);
    return NULL;
}

void pd_sound_sample_free(void* sample) {
    if (pd) pd->sound->sample->freeSample(sample);
}

// FilePlayer
void* pd_sound_fileplayer_new(void) {
    if (pd) return pd->sound->fileplayer->newPlayer();
    return NULL;
}

void pd_sound_fileplayer_free(void* player) {
    if (pd) pd->sound->fileplayer->freePlayer(player);
}

int pd_sound_fileplayer_load(void* player, const char* path) {
    if (pd) return pd->sound->fileplayer->loadIntoPlayer(player, path);
    return 0;
}

int pd_sound_fileplayer_play(void* player, int repeat) {
    if (pd) return pd->sound->fileplayer->play(player, repeat);
    return 0;
}

void pd_sound_fileplayer_stop(void* player) {
    if (pd) pd->sound->fileplayer->stop(player);
}

void pd_sound_fileplayer_pause(void* player) {
    if (pd) pd->sound->fileplayer->pause(player);
}

int pd_sound_fileplayer_isPlaying(void* player) {
    if (pd) return pd->sound->fileplayer->isPlaying(player);
    return 0;
}

void pd_sound_fileplayer_setVolume(void* player, float left, float right) {
    if (pd) pd->sound->fileplayer->setVolume(player, left, right);
}

void pd_sound_fileplayer_getVolume(void* player, float* left, float* right) {
    if (pd) pd->sound->fileplayer->getVolume(player, left, right);
}

float pd_sound_fileplayer_getLength(void* player) {
    if (pd) return pd->sound->fileplayer->getLength(player);
    return 0;
}

void pd_sound_fileplayer_setOffset(void* player, float offset) {
    if (pd) pd->sound->fileplayer->setOffset(player, offset);
}

float pd_sound_fileplayer_getOffset(void* player) {
    if (pd) return pd->sound->fileplayer->getOffset(player);
    return 0;
}

void pd_sound_fileplayer_setRate(void* player, float rate) {
    if (pd) pd->sound->fileplayer->setRate(player, rate);
}

// SamplePlayer
void* pd_sound_sampleplayer_new(void) {
    if (pd) return pd->sound->sampleplayer->newPlayer();
    return NULL;
}

void pd_sound_sampleplayer_free(void* player) {
    if (pd) pd->sound->sampleplayer->freePlayer(player);
}

void pd_sound_sampleplayer_setSample(void* player, void* sample) {
    if (pd) pd->sound->sampleplayer->setSample(player, sample);
}

int pd_sound_sampleplayer_play(void* player, int repeat, float rate) {
    if (pd) return pd->sound->sampleplayer->play(player, repeat, rate);
    return 0;
}

void pd_sound_sampleplayer_stop(void* player) {
    if (pd) pd->sound->sampleplayer->stop(player);
}

int pd_sound_sampleplayer_isPlaying(void* player) {
    if (pd) return pd->sound->sampleplayer->isPlaying(player);
    return 0;
}

void pd_sound_sampleplayer_setVolume(void* player, float left, float right) {
    if (pd) pd->sound->sampleplayer->setVolume(player, left, right);
}

// Sequence
void* pd_sound_sequence_new(void) {
    if (pd) return pd->sound->sequence->newSequence();
    return NULL;
}

void pd_sound_sequence_free(void* seq) {
    if (pd) pd->sound->sequence->freeSequence(seq);
}

int pd_sound_sequence_loadMIDI(void* seq, const char* path) {
    if (pd) return pd->sound->sequence->loadMIDIFile(seq, path);
    return 0;
}

int pd_sound_sequence_getTrackCount(void* seq) {
    if (pd) return pd->sound->sequence->getTrackCount(seq);
    return 0;
}

void* pd_sound_sequence_getTrackAtIndex(void* seq, int idx) {
    if (pd) return pd->sound->sequence->getTrackAtIndex(seq, idx);
    return NULL;
}

void pd_sound_sequence_play(void* seq, void* finishCallback, void* userdata) {
    if (pd) pd->sound->sequence->play(seq, finishCallback, userdata);
}

void pd_sound_sequence_stop(void* seq) {
    if (pd) pd->sound->sequence->stop(seq);
}

int pd_sound_sequence_getCurrentStep(void* seq, int* timeOffset) {
    if (pd) return pd->sound->sequence->getCurrentStep(seq, timeOffset);
    return 0;
}

// Track
void pd_sound_track_setInstrument(void* track, void* inst) {
    if (pd) pd->sound->track->setInstrument(track, inst);
}

int pd_sound_track_getPolyphony(void* track) {
    if (pd) return pd->sound->track->getPolyphony(track);
    return 0;
}

int pd_sound_track_getIndexForStep(void* track, uint32_t step) {
    if (pd) return pd->sound->track->getIndexForStep(track, step);
    return 0;
}

int pd_sound_track_getNoteAtIndex(void* track, int index, uint32_t* step, uint32_t* len, float* note, float* vel) {
    if (pd) return pd->sound->track->getNoteAtIndex(track, index, step, len, note, vel);
    return 0;
}

// Instrument
void* pd_sound_instrument_new(void) {
    if (pd) return pd->sound->instrument->newInstrument();
    return NULL;
}

void pd_sound_instrument_free(void* inst) {
    if (pd) pd->sound->instrument->freeInstrument(inst);
}

void pd_sound_instrument_setVolume(void* inst, float left, float right) {
    if (pd) pd->sound->instrument->setVolume(inst, left, right);
}

int pd_sound_instrument_addVoice(void* inst, void* synth, float rangeStart, float rangeEnd, float transpose) {
    if (pd) return pd->sound->instrument->addVoice(inst, synth, rangeStart, rangeEnd, transpose);
    return 0;
}

// Channel with instrument (cast to SoundSource)
int pd_sound_channel_addInstrument(void* channel, void* inst) {
    if (pd) return pd->sound->channel->addSource(channel, (SoundSource*)inst);
    return 0;
}

// Global sound
void pd_sound_getHeadphoneState(int* headphone, int* mic) {
    if (pd) pd->sound->getHeadphoneState(headphone, mic, NULL);
}

void pd_sound_setOutputsActive(int headphone, int speaker) {
    if (pd) pd->sound->setOutputsActive(headphone, speaker);
}

// ============== File API ==============

void* pd_file_open(const char* path, int mode) {
    if (pd) return pd->file->open(path, mode);
    return NULL;
}

int pd_file_close(void* file) {
    if (pd) return pd->file->close(file);
    return -1;
}

int pd_file_read(void* file, void* buf, uint32_t len) {
    if (pd) return pd->file->read(file, buf, len);
    return -1;
}

int pd_file_write(void* file, const void* buf, uint32_t len) {
    if (pd) return pd->file->write(file, buf, len);
    return -1;
}

int pd_file_flush(void* file) {
    if (pd) return pd->file->flush(file);
    return -1;
}

int pd_file_tell(void* file) {
    if (pd) return pd->file->tell(file);
    return -1;
}

int pd_file_seek(void* file, int pos, int whence) {
    if (pd) return pd->file->seek(file, pos, whence);
    return -1;
}

int pd_file_mkdir(const char* path) {
    if (pd) return pd->file->mkdir(path);
    return -1;
}

int pd_file_unlink(const char* path, int recursive) {
    if (pd) return pd->file->unlink(path, recursive);
    return -1;
}

int pd_file_rename(const char* from, const char* to) {
    if (pd) return pd->file->rename(from, to);
    return -1;
}

int pd_file_stat(const char* path, FileStat* stat) {
    if (pd) return pd->file->stat(path, stat);
    return -1;
}

// ============== Sprite API ==============

void* pd_sprite_new(void) {
    if (pd) return pd->sprite->newSprite();
    return NULL;
}

void pd_sprite_free(void* sprite) {
    if (pd) pd->sprite->freeSprite(sprite);
}

void pd_sprite_add(void* sprite) {
    if (pd) pd->sprite->addSprite(sprite);
}

void pd_sprite_remove(void* sprite) {
    if (pd) pd->sprite->removeSprite(sprite);
}

void pd_sprite_removeAll(void) {
    if (pd) pd->sprite->removeAllSprites();
}

int pd_sprite_getCount(void) {
    if (pd) return pd->sprite->getSpriteCount();
    return 0;
}

void pd_sprite_setImage(void* sprite, void* image, int flip) {
    if (pd) pd->sprite->setImage(sprite, image, flip);
}

void* pd_sprite_getImage(void* sprite) {
    if (pd) return pd->sprite->getImage(sprite);
    return NULL;
}

void pd_sprite_setBounds(void* sprite, float x, float y, float w, float h) {
    PDRect r = { x, y, w, h };
    if (pd) pd->sprite->setBounds(sprite, r);
}

void pd_sprite_getBounds(void* sprite, float* x, float* y, float* w, float* h) {
    if (pd) {
        PDRect r = pd->sprite->getBounds(sprite);
        if (x) *x = r.x;
        if (y) *y = r.y;
        if (w) *w = r.width;
        if (h) *h = r.height;
    }
}

void pd_sprite_moveTo(void* sprite, float x, float y) {
    if (pd) pd->sprite->moveTo(sprite, x, y);
}

void pd_sprite_moveBy(void* sprite, float dx, float dy) {
    if (pd) pd->sprite->moveBy(sprite, dx, dy);
}

void pd_sprite_getPosition(void* sprite, float* x, float* y) {
    if (pd) pd->sprite->getPosition(sprite, x, y);
}

void pd_sprite_setZIndex(void* sprite, int16_t z) {
    if (pd) pd->sprite->setZIndex(sprite, z);
}

int16_t pd_sprite_getZIndex(void* sprite) {
    if (pd) return pd->sprite->getZIndex(sprite);
    return 0;
}

void pd_sprite_setTag(void* sprite, uint8_t tag) {
    if (pd) pd->sprite->setTag(sprite, tag);
}

uint8_t pd_sprite_getTag(void* sprite) {
    if (pd) return pd->sprite->getTag(sprite);
    return 0;
}

void pd_sprite_setVisible(void* sprite, int visible) {
    if (pd) pd->sprite->setVisible(sprite, visible);
}

int pd_sprite_isVisible(void* sprite) {
    if (pd) return pd->sprite->isVisible(sprite);
    return 0;
}

void pd_sprite_setOpaque(void* sprite, int opaque) {
    if (pd) pd->sprite->setOpaque(sprite, opaque);
}

void pd_sprite_setDrawMode(void* sprite, int mode) {
    if (pd) pd->sprite->setDrawMode(sprite, mode);
}

void pd_sprite_setImageFlip(void* sprite, int flip) {
    if (pd) pd->sprite->setImageFlip(sprite, flip);
}

int pd_sprite_getImageFlip(void* sprite) {
    if (pd) return pd->sprite->getImageFlip(sprite);
    return 0;
}

void pd_sprite_setUpdatesEnabled(void* sprite, int enabled) {
    if (pd) pd->sprite->setUpdatesEnabled(sprite, enabled);
}

void pd_sprite_markDirty(void* sprite) {
    if (pd) pd->sprite->markDirty(sprite);
}

void pd_sprite_drawSprites(void) {
    if (pd) pd->sprite->drawSprites();
}

void pd_sprite_updateAndDrawSprites(void) {
    if (pd) pd->sprite->updateAndDrawSprites();
}

void pd_sprite_setAlwaysRedraw(int always) {
    if (pd) pd->sprite->setAlwaysRedraw(always);
}

void pd_sprite_setCollideRect(void* sprite, float x, float y, float w, float h) {
    PDRect r = { x, y, w, h };
    if (pd) pd->sprite->setCollideRect(sprite, r);
}

void pd_sprite_getCollideRect(void* sprite, float* x, float* y, float* w, float* h) {
    if (pd) {
        PDRect r = pd->sprite->getCollideRect(sprite);
        if (x) *x = r.x;
        if (y) *y = r.y;
        if (w) *w = r.width;
        if (h) *h = r.height;
    }
}

void pd_sprite_clearCollideRect(void* sprite) {
    if (pd) pd->sprite->clearCollideRect(sprite);
}

void pd_sprite_setCollisionsEnabled(void* sprite, int enabled) {
    if (pd) pd->sprite->setCollisionsEnabled(sprite, enabled);
}

void pd_sprite_resetCollisionWorld(void) {
    if (pd) pd->sprite->resetCollisionWorld();
}

// ============== Lua API ==============

int pd_lua_getArgCount(void) {
    if (pd) return pd->lua->getArgCount();
    return 0;
}

int pd_lua_getArgType(int pos, const char** outClass) {
    if (pd) return pd->lua->getArgType(pos, outClass);
    return 0;
}

int pd_lua_argIsNil(int pos) {
    if (pd) return pd->lua->argIsNil(pos);
    return 1;
}

int pd_lua_getArgBool(int pos) {
    if (pd) return pd->lua->getArgBool(pos);
    return 0;
}

int pd_lua_getArgInt(int pos) {
    if (pd) return pd->lua->getArgInt(pos);
    return 0;
}

float pd_lua_getArgFloat(int pos) {
    if (pd) return pd->lua->getArgFloat(pos);
    return 0;
}

const char* pd_lua_getArgString(int pos) {
    if (pd) return pd->lua->getArgString(pos);
    return NULL;
}

void pd_lua_pushNil(void) {
    if (pd) pd->lua->pushNil();
}

void pd_lua_pushBool(int val) {
    if (pd) pd->lua->pushBool(val);
}

void pd_lua_pushInt(int val) {
    if (pd) pd->lua->pushInt(val);
}

void pd_lua_pushFloat(float val) {
    if (pd) pd->lua->pushFloat(val);
}

void pd_lua_pushString(const char* str) {
    if (pd) pd->lua->pushString(str);
}

// ============== JSON API ==============

// JSON parsing is callback-based, simplified version
int pd_json_decode(const char* str) {
    // JSON API requires callbacks, return 0 for now
    (void)str;
    return 0;
}

// ============== Scoreboards API ==============

void pd_scoreboards_addScore(const char* boardID, uint32_t value, void* callback) {
    if (pd && pd->scoreboards) pd->scoreboards->addScore(boardID, value, (AddScoreCallback)callback);
}

void pd_scoreboards_getPersonalBest(const char* boardID, void* callback) {
    if (pd && pd->scoreboards) pd->scoreboards->getPersonalBest(boardID, (PersonalBestCallback)callback);
}

void pd_scoreboards_freeScore(void* score) {
    if (pd && pd->scoreboards) pd->scoreboards->freeScore(score);
}

void pd_scoreboards_getScoreboards(void* callback) {
    if (pd && pd->scoreboards) pd->scoreboards->getScoreboards((BoardsListCallback)callback);
}

void pd_scoreboards_freeBoardsList(void* boards) {
    if (pd && pd->scoreboards) pd->scoreboards->freeBoardsList(boards);
}

void pd_scoreboards_getScores(const char* boardID, void* callback) {
    if (pd && pd->scoreboards) pd->scoreboards->getScores(boardID, (ScoresCallback)callback);
}

void pd_scoreboards_freeScoresList(void* scores) {
    if (pd && pd->scoreboards) pd->scoreboards->freeScoresList(scores);
}

// ============== Sprite Callbacks ==============

// Go trampolines - called from C wrappers
extern void pdgo_sprite_update_trampoline(void* sprite);
extern void pdgo_sprite_draw_trampoline(void* sprite, float bx, float by, float bw, float bh, float dx, float dy, float dw, float dh);
extern int pdgo_sprite_collision_trampoline(void* sprite, void* other);

// C wrapper for update callback
static void sprite_update_wrapper(LCDSprite* sprite) {
    pdgo_sprite_update_trampoline(sprite);
}

// C wrapper for draw callback
static void sprite_draw_wrapper(LCDSprite* sprite, PDRect bounds, PDRect drawRect) {
    pdgo_sprite_draw_trampoline(sprite, bounds.x, bounds.y, bounds.width, bounds.height,
                                 drawRect.x, drawRect.y, drawRect.width, drawRect.height);
}

// C wrapper for collision response callback
static SpriteCollisionResponseType sprite_collision_wrapper(LCDSprite* sprite, LCDSprite* other) {
    return (SpriteCollisionResponseType)pdgo_sprite_collision_trampoline(sprite, other);
}

void pd_sprite_setUpdateFunction(void* s) {
    if (pd && s) pd->sprite->setUpdateFunction(s, sprite_update_wrapper);
}

void pd_sprite_setDrawFunction(void* s) {
    if (pd && s) pd->sprite->setDrawFunction(s, sprite_draw_wrapper);
}

void pd_sprite_setCollisionResponseFunction(void* s) {
    if (pd && s) pd->sprite->setCollisionResponseFunction(s, sprite_collision_wrapper);
}

// ============== Entry Points ==============

// For simulator builds, the eventHandler is exported from Go (main_cgo.go)
// This function just sets the pd pointer
void pd_set_api(void* playdate) {
    pd = (PlaydateAPI*)playdate;
}

#ifdef TARGET_PLAYDATE
// Device entry point - calls Go init
// Update callback is set via pdgo.SetUpdateCallback() in go_init
extern void runtime_init(void);
extern void go_init(void* playdate);

int eventHandler(PlaydateAPI* playdate, PDSystemEvent event, uint32_t arg) {
    (void)arg;
    if (event == kEventInit) {
        pd = playdate;
        runtime_init();
        go_init(pd);
    }
    return 0;
}
#endif // TARGET_PLAYDATE
