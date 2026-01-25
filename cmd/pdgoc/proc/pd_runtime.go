package proc

const rawRuntimeC = `// Playdate C Runtime for TinyGo pdgo bindings
// This file provides C wrappers for all Playdate SDK functions

#include <stddef.h>
#include <stdint.h>

typedef enum { kEventInit = 0 } PDSystemEvent;
typedef int PDCallbackFunction(void* userdata);
typedef uint32_t LCDColor;
typedef int LCDSolidColor;
typedef int LCDBitmapDrawMode;
typedef int LCDBitmapFlip;
typedef int PDStringEncoding;
typedef int LCDLineCapStyle;
typedef int FileOptions;

struct LCDBitmap;
struct LCDFont;
struct LCDSprite;
struct LCDBitmapTable;
struct SDFile;
struct FilePlayer;
struct SamplePlayer;
struct AudioSample;

// ============== System API ==============
struct playdate_sys {
    void* (*realloc)(void* ptr, size_t size);
    int (*formatString)(char** ret, const char* fmt, ...);
    void (*logToConsole)(const char* fmt, ...);
    void (*error)(const char* fmt, ...);
    int (*getLanguage)(void);
    unsigned int (*getCurrentTimeMilliseconds)(void);
    unsigned int (*getSecondsSinceEpoch)(unsigned int* milliseconds);
    void (*drawFPS)(int x, int y);
    void (*setUpdateCallback)(PDCallbackFunction* update, void* userdata);
    void (*getButtonState)(uint32_t* current, uint32_t* pushed, uint32_t* released);
    void (*setPeripheralsEnabled)(uint32_t mask);
    void (*getAccelerometer)(float* outx, float* outy, float* outz);
    float (*getCrankChange)(void);
    float (*getCrankAngle)(void);
    int (*isCrankDocked)(void);
    int (*setCrankSoundsDisabled)(int flag);
    int (*getFlipped)(void);
    void (*setAutoLockDisabled)(int disable);
    void (*setMenuImage)(struct LCDBitmap* bitmap, int xOffset);
    void* addMenuItem;
    void* addCheckmarkMenuItem;
    void* addOptionsMenuItem;
    void (*removeAllMenuItems)(void);
    void* removeMenuItem;
    int (*getMenuItemValue)(void* menuItem);
    void (*setMenuItemValue)(void* menuItem, int value);
    const char* (*getMenuItemTitle)(void* menuItem);
    void (*setMenuItemTitle)(void* menuItem, const char* title);
    void* (*getMenuItemUserdata)(void* menuItem);
    void (*setMenuItemUserdata)(void* menuItem, void* ud);
    int (*getReduceFlashing)(void);
    float (*getElapsedTime)(void);
    void (*resetElapsedTime)(void);
    float (*getBatteryPercentage)(void);
    float (*getBatteryVoltage)(void);
    int32_t (*getTimezoneOffset)(void);
    int (*shouldDisplay24HourTime)(void);
    void (*convertEpochToDateTime)(uint32_t epoch, void* datetime);
    uint32_t (*convertDateTimeToEpoch)(void* datetime);
    void (*clearICache)(void);
};

// ============== File API ==============
struct playdate_file {
    const char* (*geterr)(void);
    int (*listfiles)(const char* path, void (*callback)(const char* path, void* userdata), void* userdata, int showhidden);
    int (*stat)(const char* path, void* stat);
    int (*mkdir)(const char* path);
    int (*unlink)(const char* name, int recursive);
    int (*rename)(const char* from, const char* to);
    struct SDFile* (*open)(const char* name, FileOptions mode);
    int (*close)(struct SDFile* file);
    int (*read)(struct SDFile* file, void* buf, unsigned int len);
    int (*write)(struct SDFile* file, const void* buf, unsigned int len);
    int (*flush)(struct SDFile* file);
    int (*tell)(struct SDFile* file);
    int (*seek)(struct SDFile* file, int pos, int whence);
};

// ============== Graphics API ==============
struct playdate_graphics {
    void* video;
    void (*clear)(LCDColor color);
    void (*setBackgroundColor)(LCDSolidColor color);
    void (*setStencil)(struct LCDBitmap* stencil);
    LCDBitmapDrawMode (*setDrawMode)(LCDBitmapDrawMode mode);
    void (*setDrawOffset)(int dx, int dy);
    void (*setClipRect)(int x, int y, int width, int height);
    void (*clearClipRect)(void);
    void (*setLineCapStyle)(LCDLineCapStyle endCapStyle);
    void (*setFont)(struct LCDFont* font);
    void (*setTextTracking)(int tracking);
    void (*pushContext)(struct LCDBitmap* target);
    void (*popContext)(void);
    void (*drawBitmap)(struct LCDBitmap* bitmap, int x, int y, LCDBitmapFlip flip);
    void (*tileBitmap)(struct LCDBitmap* bitmap, int x, int y, int width, int height, LCDBitmapFlip flip);
    void (*drawLine)(int x1, int y1, int x2, int y2, int width, LCDColor color);
    void (*fillTriangle)(int x1, int y1, int x2, int y2, int x3, int y3, LCDColor color);
    void (*drawRect)(int x, int y, int width, int height, LCDColor color);
    void (*fillRect)(int x, int y, int width, int height, LCDColor color);
    void (*drawEllipse)(int x, int y, int width, int height, int lineWidth, float startAngle, float endAngle, LCDColor color);
    void (*fillEllipse)(int x, int y, int width, int height, float startAngle, float endAngle, LCDColor color);
    void (*drawScaledBitmap)(struct LCDBitmap* bitmap, int x, int y, float xscale, float yscale);
    int (*drawText)(const void* text, size_t len, PDStringEncoding encoding, int x, int y);
    struct LCDBitmap* (*newBitmap)(int width, int height, LCDColor bgcolor);
    void (*freeBitmap)(struct LCDBitmap* bitmap);
    struct LCDBitmap* (*loadBitmap)(const char* path, const char** outerr);
    struct LCDBitmap* (*copyBitmap)(struct LCDBitmap* bitmap);
    void (*loadIntoBitmap)(const char* path, struct LCDBitmap* bitmap, const char** outerr);
    void (*getBitmapData)(struct LCDBitmap* bitmap, int* width, int* height, int* rowbytes, uint8_t** mask, uint8_t** data);
    void (*clearBitmap)(struct LCDBitmap* bitmap, LCDColor bgcolor);
    struct LCDBitmap* (*rotatedBitmap)(struct LCDBitmap* bitmap, float rotation, float xscale, float yscale, int* allocedSize);
    struct LCDBitmapTable* (*newBitmapTable)(int count, int width, int height);
    void (*freeBitmapTable)(struct LCDBitmapTable* table);
    struct LCDBitmapTable* (*loadBitmapTable)(const char* path, const char** outerr);
    void (*loadIntoBitmapTable)(const char* path, struct LCDBitmapTable* table, const char** outerr);
    struct LCDBitmap* (*getTableBitmap)(struct LCDBitmapTable* table, int idx);
    struct LCDFont* (*loadFont)(const char* path, const char** outErr);
    void* (*getFontPage)(struct LCDFont* font, uint32_t c);
    void* (*getPageGlyph)(void* page, uint32_t c, struct LCDBitmap** bitmap, int* advance);
    int (*getGlyphKerning)(void* glyph, uint32_t glyphcode, uint32_t nextcode);
    int (*getTextWidth)(struct LCDFont* font, const void* text, size_t len, PDStringEncoding encoding, int tracking);
    uint8_t* (*getFrame)(void);
    uint8_t* (*getDisplayFrame)(void);
    void* (*getDebugBitmap)(void);
    struct LCDBitmap* (*copyFrameBufferBitmap)(void);
    void (*markUpdatedRows)(int start, int end);
    void (*display)(void);
    void (*setColorToPattern)(LCDColor* color, struct LCDBitmap* bitmap, int x, int y);
    int (*checkMaskCollision)(struct LCDBitmap* bitmap1, int x1, int y1, LCDBitmapFlip flip1, struct LCDBitmap* bitmap2, int x2, int y2, LCDBitmapFlip flip2, void* rect);
    void (*setScreenClipRect)(int x, int y, int width, int height);
    void (*fillPolygon)(int nPoints, int* coords, LCDColor color, int fillrule);
    uint8_t (*getFontHeight)(struct LCDFont* font);
    struct LCDBitmap* (*getDisplayBufferBitmap)(void);
    void (*drawRotatedBitmap)(struct LCDBitmap* bitmap, int x, int y, float rotation, float centerx, float centery, float xscale, float yscale);
    void (*setTextLeading)(int lineHeightAdustment);
    int (*setBitmapMask)(struct LCDBitmap* bitmap, struct LCDBitmap* mask);
    struct LCDBitmap* (*getBitmapMask)(struct LCDBitmap* bitmap);
    void (*setStencilImage)(struct LCDBitmap* stencil, int tile);
    struct LCDFont* (*makeFontFromData)(void* data, int wide);
};

// ============== Display API ==============
struct playdate_display {
    int (*getWidth)(void);
    int (*getHeight)(void);
    void (*setRefreshRate)(float rate);
    void (*setInverted)(int flag);
    void (*setScale)(unsigned int s);
    void (*setMosaic)(unsigned int x, unsigned int y);
    void (*setFlipped)(int x, int y);
    void (*setOffset)(int x, int y);
    float (*getRefreshRate)(void);
};

// ============== Sprite API ==============
struct playdate_sprite {
    void (*setAlwaysRedraw)(int flag);
    void (*addDirtyRect)(int x, int y, int width, int height);
    void (*drawSprites)(void);
    void (*updateAndDrawSprites)(void);
    struct LCDSprite* (*newSprite)(void);
    void (*freeSprite)(struct LCDSprite* sprite);
    struct LCDSprite* (*copy)(struct LCDSprite* sprite);
    void (*addSprite)(struct LCDSprite* sprite);
    void (*removeSprite)(struct LCDSprite* sprite);
    void (*removeSprites)(struct LCDSprite** sprites, int count);
    void (*removeAllSprites)(void);
    int (*getSpriteCount)(void);
    void (*setBounds)(struct LCDSprite* sprite, float x, float y, float width, float height);
    void (*getBounds)(struct LCDSprite* sprite, float* x, float* y, float* width, float* height);
    void (*moveTo)(struct LCDSprite* sprite, float x, float y);
    void (*moveBy)(struct LCDSprite* sprite, float dx, float dy);
    void (*setImage)(struct LCDSprite* sprite, struct LCDBitmap* image, LCDBitmapFlip flip);
    struct LCDBitmap* (*getImage)(struct LCDSprite* sprite);
    void (*setSize)(struct LCDSprite* s, float width, float height);
    void (*setZIndex)(struct LCDSprite* sprite, int16_t zIndex);
    int16_t (*getZIndex)(struct LCDSprite* sprite);
    void (*setDrawMode)(struct LCDSprite* sprite, LCDBitmapDrawMode mode);
    void (*setImageFlip)(struct LCDSprite* sprite, LCDBitmapFlip flip);
    LCDBitmapFlip (*getImageFlip)(struct LCDSprite* sprite);
    void (*setStencil)(struct LCDSprite* sprite, struct LCDBitmap* stencil);
    void (*setClipRect)(struct LCDSprite* sprite, int x, int y, int width, int height);
    void (*clearClipRect)(struct LCDSprite* sprite);
    void (*setClipRectsInRange)(int x, int y, int width, int height, int startZ, int endZ);
    void (*clearClipRectsInRange)(int startZ, int endZ);
    void (*setUpdatesEnabled)(struct LCDSprite* sprite, int flag);
    int (*updatesEnabled)(struct LCDSprite* sprite);
    void (*setCollisionsEnabled)(struct LCDSprite* sprite, int flag);
    int (*collisionsEnabled)(struct LCDSprite* sprite);
    void (*setVisible)(struct LCDSprite* sprite, int flag);
    int (*isVisible)(struct LCDSprite* sprite);
    void (*setOpaque)(struct LCDSprite* sprite, int flag);
    void (*markDirty)(struct LCDSprite* sprite);
    void (*setTag)(struct LCDSprite* sprite, uint8_t tag);
    uint8_t (*getTag)(struct LCDSprite* sprite);
    void (*setIgnoresDrawOffset)(struct LCDSprite* sprite, int flag);
    void (*setUpdateFunction)(struct LCDSprite* sprite, void* func);
    void (*setDrawFunction)(struct LCDSprite* sprite, void* func);
    void (*getPosition)(struct LCDSprite* sprite, float* x, float* y);
    void (*resetCollisionWorld)(void);
    void (*setCollideRect)(struct LCDSprite* sprite, float x, float y, float width, float height);
    void (*getCollideRect)(struct LCDSprite* sprite, float* x, float* y, float* width, float* height);
    void (*clearCollideRect)(struct LCDSprite* sprite);
    void (*setCollisionResponseFunction)(struct LCDSprite* sprite, void* func);
    void* (*checkCollisions)(struct LCDSprite* sprite, float goalX, float goalY, float* actualX, float* actualY, int* len);
    void* (*moveWithCollisions)(struct LCDSprite* sprite, float goalX, float goalY, float* actualX, float* actualY, int* len);
    void* (*querySpritesAtPoint)(float x, float y, int* len);
    void* (*querySpritesInRect)(float x, float y, float width, float height, int* len);
    void* (*querySpritesAlongLine)(float x1, float y1, float x2, float y2, int* len);
    void* (*querySpriteInfoAlongLine)(float x1, float y1, float x2, float y2, int* len);
    void* (*overlappingSprites)(struct LCDSprite* sprite, int* len);
    void* (*allOverlappingSprites)(int* len);
    void (*setStencilPattern)(struct LCDSprite* sprite, uint8_t pattern[8]);
    void (*clearStencil)(struct LCDSprite* sprite);
    void (*setUserdata)(struct LCDSprite* sprite, void* userdata);
    void* (*getUserdata)(struct LCDSprite* sprite);
    void (*setStencilImage)(struct LCDSprite* sprite, struct LCDBitmap* stencil, int tile);
};

// ============== Sound API ==============
struct playdate_sound_fileplayer {
    struct FilePlayer* (*newPlayer)(void);
    void (*freePlayer)(struct FilePlayer* player);
    int (*loadIntoPlayer)(struct FilePlayer* player, const char* path);
    void (*setBufferLength)(struct FilePlayer* player, float bufferLen);
    int (*play)(struct FilePlayer* player, int repeat);
    int (*isPlaying)(struct FilePlayer* player);
    void (*pause)(struct FilePlayer* player);
    void (*stop)(struct FilePlayer* player);
    void (*setVolume)(struct FilePlayer* player, float left, float right);
    void (*getVolume)(struct FilePlayer* player, float* left, float* right);
    float (*getLength)(struct FilePlayer* player);
    void (*setOffset)(struct FilePlayer* player, float offset);
    float (*getOffset)(struct FilePlayer* player);
    void (*setRate)(struct FilePlayer* player, float rate);
    float (*getRate)(struct FilePlayer* player);
    void (*setLoopRange)(struct FilePlayer* player, float start, float end);
    int (*didUnderrun)(struct FilePlayer* player);
    void (*setFinishCallback)(struct FilePlayer* player, void* callback, void* userdata);
    void (*setLoopCallback)(struct FilePlayer* player, void* callback, void* userdata);
    void (*setStopOnUnderrun)(struct FilePlayer* player, int flag);
    void (*fadeVolume)(struct FilePlayer* player, float left, float right, int32_t len, void* finishCallback, void* userdata);
    void (*setMP3StreamSource)(struct FilePlayer* player, void* dataSource, void* userdata, float bufferLen);
};

struct playdate_sound_sample {
    struct AudioSample* (*newSampleBuffer)(int byteCount);
    int (*loadIntoSample)(struct AudioSample* sample, const char* path);
    struct AudioSample* (*load)(const char* path);
    struct AudioSample* (*newSampleFromData)(uint8_t* data, int format, uint32_t sampleRate, int byteCount, int shouldFreeData);
    void (*getData)(struct AudioSample* sample, uint8_t** data, int* format, uint32_t* sampleRate, uint32_t* byteLength);
    void (*freeSample)(struct AudioSample* sample);
    float (*getLength)(struct AudioSample* sample);
    int (*decompress)(struct AudioSample* sample);
};

struct playdate_sound_sampleplayer {
    struct SamplePlayer* (*newPlayer)(void);
    void (*freePlayer)(struct SamplePlayer* player);
    void (*setSample)(struct SamplePlayer* player, struct AudioSample* sample);
    int (*play)(struct SamplePlayer* player, int repeat, float rate);
    int (*isPlaying)(struct SamplePlayer* player);
    void (*stop)(struct SamplePlayer* player);
    void (*setVolume)(struct SamplePlayer* player, float left, float right);
    void (*getVolume)(struct SamplePlayer* player, float* left, float* right);
    float (*getLength)(struct SamplePlayer* player);
    void (*setOffset)(struct SamplePlayer* player, float offset);
    float (*getOffset)(struct SamplePlayer* player);
    void (*setRate)(struct SamplePlayer* player, float rate);
    float (*getRate)(struct SamplePlayer* player);
    void (*setPlayRange)(struct SamplePlayer* player, int start, int end);
    void (*setFinishCallback)(struct SamplePlayer* player, void* callback, void* userdata);
    void (*setLoopCallback)(struct SamplePlayer* player, void* callback, void* userdata);
    void (*setPaused)(struct SamplePlayer* player, int flag);
};

struct PDSynth;
typedef int SoundWaveform;

struct playdate_sound_synth {
    struct PDSynth* (*newSynth)(void);
    void (*freeSynth)(struct PDSynth* synth);
    void (*setWaveform)(struct PDSynth* synth, SoundWaveform wave);
    void* setGenerator;
    void (*setSample)(struct PDSynth* synth, struct AudioSample* sample, uint32_t sustainStart, uint32_t sustainEnd);
    void (*setAttackTime)(struct PDSynth* synth, float attack);
    void (*setDecayTime)(struct PDSynth* synth, float decay);
    void (*setSustainLevel)(struct PDSynth* synth, float level);
    void (*setReleaseTime)(struct PDSynth* synth, float release);
    void (*setTranspose)(struct PDSynth* synth, float halfSteps);
    void* setFrequencyModulator;
    void* getFrequencyModulator;
    void* setAmplitudeModulator;
    void* getAmplitudeModulator;
    int (*getParameterCount)(struct PDSynth* synth);
    int (*setParameter)(struct PDSynth* synth, int parameter, float value);
    void* setParameterModulator;
    void* getParameterModulator;
    void (*playNote)(struct PDSynth* synth, float freq, float vel, float len, uint32_t when);
    void (*playMIDINote)(struct PDSynth* synth, float note, float vel, float len, uint32_t when);
    void (*noteOff)(struct PDSynth* synth, uint32_t when);
    void (*stop)(struct PDSynth* synth);
    void (*setVolume)(struct PDSynth* synth, float left, float right);
    void (*getVolume)(struct PDSynth* synth, float* left, float* right);
    int (*isPlaying)(struct PDSynth* synth);
    void* getEnvelope;
    int (*setWavetable)(struct PDSynth* synth, struct AudioSample* sample, int log2size, int columns, int rows);
    struct PDSynth* (*copy)(struct PDSynth* synth);
};

struct playdate_sound {
    void* channel;
    const struct playdate_sound_fileplayer* fileplayer;
    const struct playdate_sound_sample* sample;
    const struct playdate_sound_sampleplayer* sampleplayer;
    const struct playdate_sound_synth* synth;
    void* sequence;
    void* effect;
    void* lfo;
    void* envelope;
    void* source;
    void* controlsignal;
    void* track;
    void* instrument;
    uint32_t (*getCurrentTime)(void);
    void* (*addSource)(void* callback, void* context, int stereo);
    void* getDefaultChannel;
    void* addChannel;
    void* removeChannel;
    void* setMicCallback;
    void (*getHeadphoneState)(int* headphone, int* headsetmic, void (*changeCallback)(int headphone, int mic));
    void (*setOutputsActive)(int headphone, int speaker);
    void* removeSource;
    void* signal;
};

// ============== Main API ==============
typedef struct PlaydateAPI {
    const struct playdate_sys* system;
    const struct playdate_file* file;
    const struct playdate_graphics* graphics;
    const struct playdate_sprite* sprite;
    const struct playdate_display* display;
    const struct playdate_sound* sound;
    const void* lua;
    const void* json;
    const void* scoreboards;
    const void* network;
} PlaydateAPI;

static PlaydateAPI* pd = 0;

// ============== TinyGo Runtime ==============
void* runtime__cgo_pd_realloc(void* ptr, size_t size) __asm__("runtime._cgo_pd_realloc");
void* runtime__cgo_pd_realloc(void* ptr, size_t size) { return pd ? pd->system->realloc(ptr, size) : 0; }

uint32_t runtime__cgo_pd_getCurrentTimeMS(void) __asm__("runtime._cgo_pd_getCurrentTimeMS");
uint32_t runtime__cgo_pd_getCurrentTimeMS(void) { return pd ? pd->system->getCurrentTimeMilliseconds() : 0; }

void runtime__cgo_pd_logToConsole(const char* msg) __asm__("runtime._cgo_pd_logToConsole");
void runtime__cgo_pd_logToConsole(const char* msg) { if (pd) pd->system->logToConsole("%s", msg); }

// ============== System API Wrappers ==============
void pd_sys_logToConsole(const char* msg) { if (pd) pd->system->logToConsole("%s", msg); }
void pd_sys_error(const char* msg) { if (pd) pd->system->error("%s", msg); }
void pd_sys_drawFPS(int x, int y) { if (pd) pd->system->drawFPS(x, y); }
uint32_t pd_sys_getCurrentTimeMilliseconds(void) { return pd ? pd->system->getCurrentTimeMilliseconds() : 0; }
uint32_t pd_sys_getSecondsSinceEpoch(uint32_t* ms) { return pd ? pd->system->getSecondsSinceEpoch((unsigned int*)ms) : 0; }
void pd_sys_getButtonState(uint32_t* c, uint32_t* p, uint32_t* r) { if (pd) pd->system->getButtonState(c, p, r); }
void pd_sys_setPeripheralsEnabled(uint32_t mask) { if (pd) pd->system->setPeripheralsEnabled(mask); }
void pd_sys_getAccelerometer(float* x, float* y, float* z) { if (pd) pd->system->getAccelerometer(x, y, z); }
float pd_sys_getCrankChange(void) { return pd ? pd->system->getCrankChange() : 0; }
float pd_sys_getCrankAngle(void) { return pd ? pd->system->getCrankAngle() : 0; }
int pd_sys_isCrankDocked(void) { return pd ? pd->system->isCrankDocked() : 0; }
int pd_sys_setCrankSoundsDisabled(int flag) { return pd ? pd->system->setCrankSoundsDisabled(flag) : 0; }
int pd_sys_getFlipped(void) { return pd ? pd->system->getFlipped() : 0; }
void pd_sys_setAutoLockDisabled(int flag) { if (pd) pd->system->setAutoLockDisabled(flag); }
int pd_sys_getLanguage(void) { return pd ? pd->system->getLanguage() : 0; }
float pd_sys_getBatteryPercentage(void) { return pd ? pd->system->getBatteryPercentage() : 1.0f; }
float pd_sys_getBatteryVoltage(void) { return pd ? pd->system->getBatteryVoltage() : 4.2f; }

// ============== File API Wrappers ==============
struct SDFile* pd_file_open(const char* path, int mode) { return pd ? pd->file->open(path, mode) : 0; }
int pd_file_close(struct SDFile* file) { return pd ? pd->file->close(file) : -1; }
int pd_file_read(struct SDFile* file, void* buf, unsigned int len) { return pd ? pd->file->read(file, buf, len) : -1; }
int pd_file_write(struct SDFile* file, const void* buf, unsigned int len) { return pd ? pd->file->write(file, buf, len) : -1; }
int pd_file_flush(struct SDFile* file) { return pd ? pd->file->flush(file) : -1; }
int pd_file_tell(struct SDFile* file) { return pd ? pd->file->tell(file) : -1; }
int pd_file_seek(struct SDFile* file, int pos, int whence) { return pd ? pd->file->seek(file, pos, whence) : -1; }
int pd_file_stat(const char* path, int* isDir, int* size, int* mtime) { 
    if (!pd) return -1;
    struct { int isdir; unsigned int size; int m_year, m_month, m_day, m_hour, m_minute, m_second; } st;
    int r = pd->file->stat(path, &st);
    if (r == 0) { *isDir = st.isdir; *size = st.size; *mtime = 0; }
    return r;
}
int pd_file_mkdir(const char* path) { return pd ? pd->file->mkdir(path) : -1; }
int pd_file_unlink(const char* path, int recursive) { return pd ? pd->file->unlink(path, recursive) : -1; }
int pd_file_rename(const char* from, const char* to) { return pd ? pd->file->rename(from, to) : -1; }

// ============== Graphics API Wrappers ==============
void pd_gfx_clear(uint32_t color) { if (pd) pd->graphics->clear(color); }
void pd_gfx_setBackgroundColor(int color) { if (pd) pd->graphics->setBackgroundColor(color); }
int pd_gfx_setDrawMode(int mode) { return pd ? pd->graphics->setDrawMode(mode) : 0; }
void pd_gfx_setDrawOffset(int dx, int dy) { if (pd) pd->graphics->setDrawOffset(dx, dy); }
void pd_gfx_setClipRect(int x, int y, int w, int h) { if (pd) pd->graphics->setClipRect(x, y, w, h); }
void pd_gfx_clearClipRect(void) { if (pd) pd->graphics->clearClipRect(); }
void pd_gfx_setLineCapStyle(int style) { if (pd) pd->graphics->setLineCapStyle(style); }
void pd_gfx_setFont(struct LCDFont* font) { if (pd) pd->graphics->setFont(font); }
void pd_gfx_setTextTracking(int tracking) { if (pd) pd->graphics->setTextTracking(tracking); }
void pd_gfx_pushContext(struct LCDBitmap* target) { if (pd) pd->graphics->pushContext(target); }
void pd_gfx_popContext(void) { if (pd) pd->graphics->popContext(); }
void pd_gfx_fillRect(int x, int y, int w, int h, uint32_t color) { if (pd) pd->graphics->fillRect(x, y, w, h, color); }
void pd_gfx_drawRect(int x, int y, int w, int h, uint32_t color) { if (pd) pd->graphics->drawRect(x, y, w, h, color); }
void pd_gfx_drawLine(int x1, int y1, int x2, int y2, int width, uint32_t color) { if (pd) pd->graphics->drawLine(x1, y1, x2, y2, width, color); }
void pd_gfx_fillTriangle(int x1, int y1, int x2, int y2, int x3, int y3, uint32_t color) { if (pd) pd->graphics->fillTriangle(x1, y1, x2, y2, x3, y3, color); }
void pd_gfx_drawEllipse(int x, int y, int w, int h, int lineWidth, float startAngle, float endAngle, uint32_t color) { if (pd) pd->graphics->drawEllipse(x, y, w, h, lineWidth, startAngle, endAngle, color); }
void pd_gfx_fillEllipse(int x, int y, int w, int h, float startAngle, float endAngle, uint32_t color) { if (pd) pd->graphics->fillEllipse(x, y, w, h, startAngle, endAngle, color); }
int pd_gfx_drawText(const char* text, int len, int enc, int x, int y) { return pd ? pd->graphics->drawText(text, len, enc, x, y) : 0; }
int pd_gfx_getTextWidth(struct LCDFont* font, const char* text, int len, int enc, int tracking) { return pd ? pd->graphics->getTextWidth(font, text, len, enc, tracking) : 0; }
struct LCDBitmap* pd_gfx_newBitmap(int w, int h, uint32_t bgcolor) { return pd ? pd->graphics->newBitmap(w, h, bgcolor) : 0; }
void pd_gfx_freeBitmap(struct LCDBitmap* bmp) { if (pd) pd->graphics->freeBitmap(bmp); }
struct LCDBitmap* pd_gfx_loadBitmap(const char* path) { if (!pd) return 0; const char* e = 0; return pd->graphics->loadBitmap(path, &e); }
struct LCDBitmap* pd_gfx_copyBitmap(struct LCDBitmap* bmp) { return pd ? pd->graphics->copyBitmap(bmp) : 0; }
void pd_gfx_drawBitmap(struct LCDBitmap* bmp, int x, int y, int flip) { if (pd) pd->graphics->drawBitmap(bmp, x, y, flip); }
void pd_gfx_tileBitmap(struct LCDBitmap* bmp, int x, int y, int w, int h, int flip) { if (pd) pd->graphics->tileBitmap(bmp, x, y, w, h, flip); }
void pd_gfx_drawScaledBitmap(struct LCDBitmap* bmp, int x, int y, float xs, float ys) { if (pd) pd->graphics->drawScaledBitmap(bmp, x, y, xs, ys); }
void pd_gfx_drawRotatedBitmap(struct LCDBitmap* bmp, int x, int y, float rot, float cx, float cy, float xs, float ys) { if (pd) pd->graphics->drawRotatedBitmap(bmp, x, y, rot, cx, cy, xs, ys); }
void pd_gfx_getBitmapData(struct LCDBitmap* bmp, int* w, int* h, int* rb, uint8_t** mask, uint8_t** data) { if (pd) pd->graphics->getBitmapData(bmp, w, h, rb, mask, data); }
void pd_gfx_clearBitmap(struct LCDBitmap* bmp, uint32_t bgcolor) { if (pd) pd->graphics->clearBitmap(bmp, bgcolor); }
struct LCDBitmapTable* pd_gfx_newBitmapTable(int count, int w, int h) { return pd ? pd->graphics->newBitmapTable(count, w, h) : 0; }
void pd_gfx_freeBitmapTable(struct LCDBitmapTable* table) { if (pd) pd->graphics->freeBitmapTable(table); }
struct LCDBitmapTable* pd_gfx_loadBitmapTable(const char* path) { if (!pd) return 0; const char* e = 0; return pd->graphics->loadBitmapTable(path, &e); }
struct LCDBitmap* pd_gfx_getTableBitmap(struct LCDBitmapTable* table, int idx) { return pd ? pd->graphics->getTableBitmap(table, idx) : 0; }
struct LCDFont* pd_gfx_loadFont(const char* path) { if (!pd) return 0; const char* e = 0; return pd->graphics->loadFont(path, &e); }
uint8_t* pd_gfx_getFrame(void) { return pd ? pd->graphics->getFrame() : 0; }
uint8_t* pd_gfx_getDisplayFrame(void) { return pd ? pd->graphics->getDisplayFrame() : 0; }
void pd_gfx_markUpdatedRows(int start, int end) { if (pd) pd->graphics->markUpdatedRows(start, end); }
void pd_gfx_display(void) { if (pd) pd->graphics->display(); }

// ============== Display API Wrappers ==============
int pd_display_getWidth(void) { return pd ? pd->display->getWidth() : 400; }
int pd_display_getHeight(void) { return pd ? pd->display->getHeight() : 240; }
void pd_display_setRefreshRate(float rate) { if (pd) pd->display->setRefreshRate(rate); }
float pd_display_getRefreshRate(void) { return pd ? pd->display->getRefreshRate() : 30.0f; }
void pd_display_setInverted(int flag) { if (pd) pd->display->setInverted(flag); }
void pd_display_setScale(unsigned int s) { if (pd) pd->display->setScale(s); }
void pd_display_setMosaic(unsigned int x, unsigned int y) { if (pd) pd->display->setMosaic(x, y); }
void pd_display_setFlipped(int x, int y) { if (pd) pd->display->setFlipped(x, y); }
void pd_display_setOffset(int x, int y) { if (pd) pd->display->setOffset(x, y); }

// ============== Sprite API Wrappers ==============
void pd_sprite_setAlwaysRedraw(int flag) { if (pd) pd->sprite->setAlwaysRedraw(flag); }
void pd_sprite_updateAndDrawSprites(void) { if (pd) pd->sprite->updateAndDrawSprites(); }
struct LCDSprite* pd_sprite_newSprite(void) { return pd ? pd->sprite->newSprite() : 0; }
void pd_sprite_freeSprite(struct LCDSprite* s) { if (pd) pd->sprite->freeSprite(s); }
void pd_sprite_addSprite(struct LCDSprite* s) { if (pd) pd->sprite->addSprite(s); }
void pd_sprite_removeSprite(struct LCDSprite* s) { if (pd) pd->sprite->removeSprite(s); }
void pd_sprite_removeAllSprites(void) { if (pd) pd->sprite->removeAllSprites(); }
int pd_sprite_getSpriteCount(void) { return pd ? pd->sprite->getSpriteCount() : 0; }
void pd_sprite_setBounds(struct LCDSprite* s, float x, float y, float w, float h) { if (pd) pd->sprite->setBounds(s, x, y, w, h); }
void pd_sprite_getBounds(struct LCDSprite* s, float* x, float* y, float* w, float* h) { if (pd) pd->sprite->getBounds(s, x, y, w, h); }
void pd_sprite_moveTo(struct LCDSprite* s, float x, float y) { if (pd) pd->sprite->moveTo(s, x, y); }
void pd_sprite_moveBy(struct LCDSprite* s, float dx, float dy) { if (pd) pd->sprite->moveBy(s, dx, dy); }
void pd_sprite_setImage(struct LCDSprite* s, struct LCDBitmap* img, int flip) { if (pd) pd->sprite->setImage(s, img, flip); }
struct LCDBitmap* pd_sprite_getImage(struct LCDSprite* s) { return pd ? pd->sprite->getImage(s) : 0; }
void pd_sprite_setZIndex(struct LCDSprite* s, int16_t z) { if (pd) pd->sprite->setZIndex(s, z); }
int16_t pd_sprite_getZIndex(struct LCDSprite* s) { return pd ? pd->sprite->getZIndex(s) : 0; }
void pd_sprite_setTag(struct LCDSprite* s, uint8_t tag) { if (pd) pd->sprite->setTag(s, tag); }
uint8_t pd_sprite_getTag(struct LCDSprite* s) { return pd ? pd->sprite->getTag(s) : 0; }
void pd_sprite_setVisible(struct LCDSprite* s, int flag) { if (pd) pd->sprite->setVisible(s, flag); }
int pd_sprite_isVisible(struct LCDSprite* s) { return pd ? pd->sprite->isVisible(s) : 0; }
void pd_sprite_setOpaque(struct LCDSprite* s, int flag) { if (pd) pd->sprite->setOpaque(s, flag); }
void pd_sprite_setDrawMode(struct LCDSprite* s, int mode) { if (pd) pd->sprite->setDrawMode(s, mode); }
void pd_sprite_setImageFlip(struct LCDSprite* s, int flip) { if (pd) pd->sprite->setImageFlip(s, flip); }
int pd_sprite_getImageFlip(struct LCDSprite* s) { return pd ? pd->sprite->getImageFlip(s) : 0; }
void pd_sprite_setUpdatesEnabled(struct LCDSprite* s, int flag) { if (pd) pd->sprite->setUpdatesEnabled(s, flag); }
void pd_sprite_getPosition(struct LCDSprite* s, float* x, float* y) { if (pd) pd->sprite->getPosition(s, x, y); }
void pd_sprite_resetCollisionWorld(void) { if (pd) pd->sprite->resetCollisionWorld(); }
void pd_sprite_setCollideRect(struct LCDSprite* s, float x, float y, float w, float h) { if (pd) pd->sprite->setCollideRect(s, x, y, w, h); }
void pd_sprite_getCollideRect(struct LCDSprite* s, float* x, float* y, float* w, float* h) { if (pd) pd->sprite->getCollideRect(s, x, y, w, h); }
void pd_sprite_clearCollideRect(struct LCDSprite* s) { if (pd) pd->sprite->clearCollideRect(s); }
void pd_sprite_setCollisionsEnabled(struct LCDSprite* s, int flag) { if (pd) pd->sprite->setCollisionsEnabled(s, flag); }
void* pd_sprite_moveWithCollisions(struct LCDSprite* s, float gx, float gy, float* ax, float* ay, int* len) { return pd ? pd->sprite->moveWithCollisions(s, gx, gy, ax, ay, len) : 0; }
void* pd_sprite_checkCollisions(struct LCDSprite* s, int* len) { float x, y; return pd ? pd->sprite->checkCollisions(s, 0, 0, &x, &y, len) : 0; }
void* pd_sprite_querySpritesAtPoint(float x, float y, int* len) { return pd ? pd->sprite->querySpritesAtPoint(x, y, len) : 0; }
void* pd_sprite_querySpritesInRect(float x, float y, float w, float h, int* len) { return pd ? pd->sprite->querySpritesInRect(x, y, w, h, len) : 0; }
void* pd_sprite_querySpritesAlongLine(float x1, float y1, float x2, float y2, int* len) { return pd ? pd->sprite->querySpritesAlongLine(x1, y1, x2, y2, len) : 0; }
void* pd_sprite_allOverlappingSprites(int* len) { return pd ? pd->sprite->allOverlappingSprites(len) : 0; }
void pd_sprite_drawSprites(void) { if (pd) pd->sprite->drawSprites(); }
void pd_sprite_markDirty(struct LCDSprite* s) { if (pd) pd->sprite->markDirty(s); }

// ============== Sprite Callback Trampolines ==============
// Go trampoline functions (exported from TinyGo)
extern void pdgo_sprite_update_trampoline(uintptr_t spritePtr) __asm__("pdgo_sprite_update_trampoline");
extern void pdgo_sprite_draw_trampoline(uintptr_t spritePtr, float bx, float by, float bw, float bh, float dx, float dy, float dw, float dh) __asm__("pdgo_sprite_draw_trampoline");
extern int32_t pdgo_sprite_collision_trampoline(uintptr_t spritePtr, uintptr_t otherPtr) __asm__("pdgo_sprite_collision_trampoline");

// C callback wrappers that call Go trampolines
static void sprite_update_callback_wrapper(struct LCDSprite* sprite) {
    pdgo_sprite_update_trampoline((uintptr_t)sprite);
}

typedef struct { float x; float y; float width; float height; } PDRectLocal;

static void sprite_draw_callback_wrapper(struct LCDSprite* sprite, PDRectLocal bounds, PDRectLocal drawrect) {
    pdgo_sprite_draw_trampoline((uintptr_t)sprite, bounds.x, bounds.y, bounds.width, bounds.height, drawrect.x, drawrect.y, drawrect.width, drawrect.height);
}

typedef int SpriteCollisionResponseTypeLocal;

static SpriteCollisionResponseTypeLocal sprite_collision_callback_wrapper(struct LCDSprite* sprite, struct LCDSprite* other) {
    return (SpriteCollisionResponseTypeLocal)pdgo_sprite_collision_trampoline((uintptr_t)sprite, (uintptr_t)other);
}

// Functions to register callbacks with sprites
void pd_sprite_setUpdateFunction(struct LCDSprite* s, int hasCallback) {
    if (pd && s) {
        if (hasCallback) {
            pd->sprite->setUpdateFunction(s, sprite_update_callback_wrapper);
        } else {
            pd->sprite->setUpdateFunction(s, 0);
        }
    }
}

void pd_sprite_setDrawFunction(struct LCDSprite* s, int hasCallback) {
    if (pd && s) {
        if (hasCallback) {
            pd->sprite->setDrawFunction(s, (void*)sprite_draw_callback_wrapper);
        } else {
            pd->sprite->setDrawFunction(s, 0);
        }
    }
}

void pd_sprite_setCollisionResponseFunction(struct LCDSprite* s, int hasCallback) {
    if (pd && s) {
        if (hasCallback) {
            pd->sprite->setCollisionResponseFunction(s, (void*)sprite_collision_callback_wrapper);
        } else {
            pd->sprite->setCollisionResponseFunction(s, 0);
        }
    }
}
// ============== Sound API Wrappers ==============
struct FilePlayer* pd_sound_newFilePlayer(void) { return pd && pd->sound && pd->sound->fileplayer ? pd->sound->fileplayer->newPlayer() : 0; }
void pd_sound_freeFilePlayer(struct FilePlayer* p) { if (pd && pd->sound && pd->sound->fileplayer) pd->sound->fileplayer->freePlayer(p); }
int pd_sound_loadIntoFilePlayer(struct FilePlayer* p, const char* path) { return pd && pd->sound && pd->sound->fileplayer ? pd->sound->fileplayer->loadIntoPlayer(p, path) : 0; }
void pd_sound_playFilePlayer(struct FilePlayer* p, int repeat) { if (pd && pd->sound && pd->sound->fileplayer) pd->sound->fileplayer->play(p, repeat); }
void pd_sound_stopFilePlayer(struct FilePlayer* p) { if (pd && pd->sound && pd->sound->fileplayer) pd->sound->fileplayer->stop(p); }
void pd_sound_pauseFilePlayer(struct FilePlayer* p) { if (pd && pd->sound && pd->sound->fileplayer) pd->sound->fileplayer->pause(p); }
int pd_sound_isFilePlayerPlaying(struct FilePlayer* p) { return pd && pd->sound && pd->sound->fileplayer ? pd->sound->fileplayer->isPlaying(p) : 0; }
void pd_sound_setFilePlayerVolume(struct FilePlayer* p, float l, float r) { if (pd && pd->sound && pd->sound->fileplayer) pd->sound->fileplayer->setVolume(p, l, r); }
void pd_sound_getFilePlayerVolume(struct FilePlayer* p, float* l, float* r) { if (pd && pd->sound && pd->sound->fileplayer) pd->sound->fileplayer->getVolume(p, l, r); }
float pd_sound_getFilePlayerLength(struct FilePlayer* p) { return pd && pd->sound && pd->sound->fileplayer ? pd->sound->fileplayer->getLength(p) : 0; }
void pd_sound_setFilePlayerOffset(struct FilePlayer* p, float o) { if (pd && pd->sound && pd->sound->fileplayer) pd->sound->fileplayer->setOffset(p, o); }
float pd_sound_getFilePlayerOffset(struct FilePlayer* p) { return pd && pd->sound && pd->sound->fileplayer ? pd->sound->fileplayer->getOffset(p) : 0; }
void pd_sound_setFilePlayerRate(struct FilePlayer* p, float r) { if (pd && pd->sound && pd->sound->fileplayer) pd->sound->fileplayer->setRate(p, r); }

struct SamplePlayer* pd_sound_newSamplePlayer(void) { return pd && pd->sound && pd->sound->sampleplayer ? pd->sound->sampleplayer->newPlayer() : 0; }
void pd_sound_freeSamplePlayer(struct SamplePlayer* p) { if (pd && pd->sound && pd->sound->sampleplayer) pd->sound->sampleplayer->freePlayer(p); }
void pd_sound_setSamplePlayerSample(struct SamplePlayer* p, struct AudioSample* s) { if (pd && pd->sound && pd->sound->sampleplayer) pd->sound->sampleplayer->setSample(p, s); }
void pd_sound_playSamplePlayer(struct SamplePlayer* p, int repeat, float rate) { if (pd && pd->sound && pd->sound->sampleplayer) pd->sound->sampleplayer->play(p, repeat, rate); }
void pd_sound_stopSamplePlayer(struct SamplePlayer* p) { if (pd && pd->sound && pd->sound->sampleplayer) pd->sound->sampleplayer->stop(p); }
int pd_sound_isSamplePlayerPlaying(struct SamplePlayer* p) { return pd && pd->sound && pd->sound->sampleplayer ? pd->sound->sampleplayer->isPlaying(p) : 0; }
void pd_sound_setSamplePlayerVolume(struct SamplePlayer* p, float l, float r) { if (pd && pd->sound && pd->sound->sampleplayer) pd->sound->sampleplayer->setVolume(p, l, r); }

struct AudioSample* pd_sound_newSample(int len) { return pd && pd->sound && pd->sound->sample ? pd->sound->sample->newSampleBuffer(len) : 0; }
struct AudioSample* pd_sound_loadSample(const char* path) { return pd && pd->sound && pd->sound->sample ? pd->sound->sample->load(path) : 0; }
void pd_sound_freeSample(struct AudioSample* s) { if (pd && pd->sound && pd->sound->sample) pd->sound->sample->freeSample(s); }

void pd_sound_getHeadphoneState(int* h, int* m) { if (pd && pd->sound) pd->sound->getHeadphoneState(h, m, 0); }
void pd_sound_setOutputsActive(int h, int s) { if (pd && pd->sound) pd->sound->setOutputsActive(h, s); }

// ============== Synth API Wrappers ==============
struct PDSynth* pd_sound_synth_new(void) { return pd && pd->sound && pd->sound->synth ? pd->sound->synth->newSynth() : 0; }
void pd_sound_synth_free(struct PDSynth* s) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->freeSynth(s); }
void pd_sound_synth_setWaveform(struct PDSynth* s, int wave) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->setWaveform(s, wave); }
void pd_sound_synth_setAttackTime(struct PDSynth* s, float t) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->setAttackTime(s, t); }
void pd_sound_synth_setDecayTime(struct PDSynth* s, float t) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->setDecayTime(s, t); }
void pd_sound_synth_setSustainLevel(struct PDSynth* s, float l) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->setSustainLevel(s, l); }
void pd_sound_synth_setReleaseTime(struct PDSynth* s, float t) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->setReleaseTime(s, t); }
void pd_sound_synth_setTranspose(struct PDSynth* s, float h) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->setTranspose(s, h); }
void pd_sound_synth_playNote(struct PDSynth* s, float freq, float vel, float len, uint32_t when) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->playNote(s, freq, vel, len, when); }
void pd_sound_synth_playMIDINote(struct PDSynth* s, float note, float vel, float len, uint32_t when) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->playMIDINote(s, note, vel, len, when); }
void pd_sound_synth_noteOff(struct PDSynth* s, uint32_t when) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->noteOff(s, when); }
void pd_sound_synth_stop(struct PDSynth* s) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->stop(s); }
void pd_sound_synth_setVolume(struct PDSynth* s, float l, float r) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->setVolume(s, l, r); }
void pd_sound_synth_getVolume(struct PDSynth* s, float* l, float* r) { if (pd && pd->sound && pd->sound->synth && s) pd->sound->synth->getVolume(s, l, r); }
int pd_sound_synth_isPlaying(struct PDSynth* s) { return pd && pd->sound && pd->sound->synth && s ? pd->sound->synth->isPlaying(s) : 0; }

// ============== Event Handler ==============
extern int updateCallback(void* userdata);
extern int eventHandler(PlaydateAPI* playdate, PDSystemEvent event, uint32_t arg);

int eventHandlerShim(PlaydateAPI* playdate, PDSystemEvent event, uint32_t arg) {
    if (event == kEventInit) {
        pd = playdate;
        // TinyGo runtime initializes automatically on first Go function call
        int result = eventHandler(playdate, event, arg);
        pd->system->setUpdateCallback(updateCallback, 0);
        return result;
    }
    return eventHandler(playdate, event, arg);
}`
