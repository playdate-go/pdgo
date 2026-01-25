package proc

const rawBridgeTemplate = `//go:build tinygo

// Bridge template for TinyGo Playdate games
// Copy this file to your game's main package
//
// Usage:
// 1. Copy this file to your game directory
// 2. Call pdgo.RegisterBridge() in init() with the functions you need

package main

import (
	_ "unsafe"

	"github.com/playdate-go/pdgo"
)

// ============== System Functions ==============

//go:linkname sysLog pd_sys_logToConsole
func sysLog(msg *byte)

//go:linkname sysError pd_sys_error
func sysError(msg *byte)

//go:linkname sysDrawFPS pd_sys_drawFPS
func sysDrawFPS(x, y int32)

//go:linkname sysGetCurrentTimeMS pd_sys_getCurrentTimeMilliseconds
func sysGetCurrentTimeMS() uint32

//go:linkname sysGetSecondsSinceEpoch pd_sys_getSecondsSinceEpoch
func sysGetSecondsSinceEpoch(ms *uint32) uint32

//go:linkname sysGetButtonState pd_sys_getButtonState
func sysGetButtonState(current, pushed, released *uint32)

//go:linkname sysSetPeripheralsEnabled pd_sys_setPeripheralsEnabled
func sysSetPeripheralsEnabled(mask uint32)

//go:linkname sysGetAccelerometer pd_sys_getAccelerometer
func sysGetAccelerometer(x, y, z *float32)

//go:linkname sysGetCrankChange pd_sys_getCrankChange
func sysGetCrankChange() float32

//go:linkname sysGetCrankAngle pd_sys_getCrankAngle
func sysGetCrankAngle() float32

//go:linkname sysIsCrankDocked pd_sys_isCrankDocked
func sysIsCrankDocked() int32

//go:linkname sysSetCrankSoundsDisabled pd_sys_setCrankSoundsDisabled
func sysSetCrankSoundsDisabled(flag int32) int32

//go:linkname sysGetFlipped pd_sys_getFlipped
func sysGetFlipped() int32

//go:linkname sysSetAutoLockDisabled pd_sys_setAutoLockDisabled
func sysSetAutoLockDisabled(flag int32)

//go:linkname sysGetLanguage pd_sys_getLanguage
func sysGetLanguage() int32

//go:linkname sysGetBatteryPercentage pd_sys_getBatteryPercentage
func sysGetBatteryPercentage() float32

//go:linkname sysGetBatteryVoltage pd_sys_getBatteryVoltage
func sysGetBatteryVoltage() float32

// ============== Graphics Functions ==============

//go:linkname gfxClear pd_gfx_clear
func gfxClear(color uint32)

//go:linkname gfxSetBackgroundColor pd_gfx_setBackgroundColor
func gfxSetBackgroundColor(color int32)

//go:linkname gfxSetDrawMode pd_gfx_setDrawMode
func gfxSetDrawMode(mode int32) int32

//go:linkname gfxSetDrawOffset pd_gfx_setDrawOffset
func gfxSetDrawOffset(dx, dy int32)

//go:linkname gfxSetClipRect pd_gfx_setClipRect
func gfxSetClipRect(x, y, w, h int32)

//go:linkname gfxClearClipRect pd_gfx_clearClipRect
func gfxClearClipRect()

//go:linkname gfxSetLineCapStyle pd_gfx_setLineCapStyle
func gfxSetLineCapStyle(style int32)

//go:linkname gfxSetFont pd_gfx_setFont
func gfxSetFont(font uintptr)

//go:linkname gfxSetTextTracking pd_gfx_setTextTracking
func gfxSetTextTracking(tracking int32)

//go:linkname gfxPushContext pd_gfx_pushContext
func gfxPushContext(bitmap uintptr)

//go:linkname gfxPopContext pd_gfx_popContext
func gfxPopContext()

//go:linkname gfxFillRect pd_gfx_fillRect
func gfxFillRect(x, y, w, h int32, color uint32)

//go:linkname gfxDrawRect pd_gfx_drawRect
func gfxDrawRect(x, y, w, h int32, color uint32)

//go:linkname gfxDrawLine pd_gfx_drawLine
func gfxDrawLine(x1, y1, x2, y2, width int32, color uint32)

//go:linkname gfxFillTriangle pd_gfx_fillTriangle
func gfxFillTriangle(x1, y1, x2, y2, x3, y3 int32, color uint32)

//go:linkname gfxDrawEllipse pd_gfx_drawEllipse
func gfxDrawEllipse(x, y, w, h, lineWidth int32, startAngle, endAngle float32, color uint32)

//go:linkname gfxFillEllipse pd_gfx_fillEllipse
func gfxFillEllipse(x, y, w, h int32, startAngle, endAngle float32, color uint32)

//go:linkname gfxDrawText pd_gfx_drawText
func gfxDrawText(text *byte, len, enc, x, y int32) int32

//go:linkname gfxGetTextWidth pd_gfx_getTextWidth
func gfxGetTextWidth(font uintptr, text *byte, len, enc, tracking int32) int32

//go:linkname gfxNewBitmap pd_gfx_newBitmap
func gfxNewBitmap(w, h int32, bgcolor uint32) uintptr

//go:linkname gfxFreeBitmap pd_gfx_freeBitmap
func gfxFreeBitmap(bitmap uintptr)

//go:linkname gfxLoadBitmap pd_gfx_loadBitmap
func gfxLoadBitmap(path *byte) uintptr

//go:linkname gfxCopyBitmap pd_gfx_copyBitmap
func gfxCopyBitmap(bitmap uintptr) uintptr

//go:linkname gfxDrawBitmap pd_gfx_drawBitmap
func gfxDrawBitmap(bitmap uintptr, x, y, flip int32)

//go:linkname gfxTileBitmap pd_gfx_tileBitmap
func gfxTileBitmap(bitmap uintptr, x, y, w, h, flip int32)

//go:linkname gfxDrawScaledBitmap pd_gfx_drawScaledBitmap
func gfxDrawScaledBitmap(bitmap uintptr, x, y int32, xscale, yscale float32)

//go:linkname gfxDrawRotatedBitmap pd_gfx_drawRotatedBitmap
func gfxDrawRotatedBitmap(bitmap uintptr, x, y int32, rotation, cx, cy, xscale, yscale float32)

//go:linkname gfxGetBitmapData pd_gfx_getBitmapData
func gfxGetBitmapData(bitmap uintptr, w, h, rowbytes *int32, mask, data *uintptr)

//go:linkname gfxClearBitmap pd_gfx_clearBitmap
func gfxClearBitmap(bitmap uintptr, bgcolor uint32)

//go:linkname gfxNewBitmapTable pd_gfx_newBitmapTable
func gfxNewBitmapTable(count, w, h int32) uintptr

//go:linkname gfxFreeBitmapTable pd_gfx_freeBitmapTable
func gfxFreeBitmapTable(table uintptr)

//go:linkname gfxLoadBitmapTable pd_gfx_loadBitmapTable
func gfxLoadBitmapTable(path *byte) uintptr

//go:linkname gfxGetTableBitmap pd_gfx_getTableBitmap
func gfxGetTableBitmap(table uintptr, idx int32) uintptr

//go:linkname gfxLoadFont pd_gfx_loadFont
func gfxLoadFont(path *byte) uintptr

//go:linkname gfxGetFrame pd_gfx_getFrame
func gfxGetFrame() uintptr

//go:linkname gfxGetDisplayFrame pd_gfx_getDisplayFrame
func gfxGetDisplayFrame() uintptr

//go:linkname gfxMarkUpdatedRows pd_gfx_markUpdatedRows
func gfxMarkUpdatedRows(start, end int32)

//go:linkname gfxDisplay pd_gfx_display
func gfxDisplay()

// ============== Display Functions ==============

//go:linkname displayGetWidth pd_display_getWidth
func displayGetWidth() int32

//go:linkname displayGetHeight pd_display_getHeight
func displayGetHeight() int32

//go:linkname displaySetRefreshRate pd_display_setRefreshRate
func displaySetRefreshRate(rate float32)

//go:linkname displayGetRefreshRate pd_display_getRefreshRate
func displayGetRefreshRate() float32

//go:linkname displaySetInverted pd_display_setInverted
func displaySetInverted(flag int32)

//go:linkname displaySetScale pd_display_setScale
func displaySetScale(s uint32)

//go:linkname displaySetMosaic pd_display_setMosaic
func displaySetMosaic(x, y uint32)

//go:linkname displaySetFlipped pd_display_setFlipped
func displaySetFlipped(x, y int32)

//go:linkname displaySetOffset pd_display_setOffset
func displaySetOffset(x, y int32)

// ============== Sprite Functions ==============

//go:linkname spriteNewSprite pd_sprite_newSprite
func spriteNewSprite() uintptr

//go:linkname spriteFreeSprite pd_sprite_freeSprite
func spriteFreeSprite(sprite uintptr)

//go:linkname spriteAddSprite pd_sprite_addSprite
func spriteAddSprite(sprite uintptr)

//go:linkname spriteRemoveSprite pd_sprite_removeSprite
func spriteRemoveSprite(sprite uintptr)

//go:linkname spriteRemoveAllSprites pd_sprite_removeAllSprites
func spriteRemoveAllSprites()

//go:linkname spriteGetSpriteCount pd_sprite_getSpriteCount
func spriteGetSpriteCount() int32

//go:linkname spriteSetImage pd_sprite_setImage
func spriteSetImage(sprite, image uintptr, flip int32)

//go:linkname spriteGetImage pd_sprite_getImage
func spriteGetImage(sprite uintptr) uintptr

//go:linkname spriteSetBounds pd_sprite_setBounds
func spriteSetBounds(sprite uintptr, x, y, w, h float32)

//go:linkname spriteGetBounds pd_sprite_getBounds
func spriteGetBounds(sprite uintptr, x, y, w, h *float32)

//go:linkname spriteMoveTo pd_sprite_moveTo
func spriteMoveTo(sprite uintptr, x, y float32)

//go:linkname spriteMoveBy pd_sprite_moveBy
func spriteMoveBy(sprite uintptr, dx, dy float32)

//go:linkname spriteGetPosition pd_sprite_getPosition
func spriteGetPosition(sprite uintptr, x, y *float32)

//go:linkname spriteSetZIndex pd_sprite_setZIndex
func spriteSetZIndex(sprite uintptr, z int16)

//go:linkname spriteGetZIndex pd_sprite_getZIndex
func spriteGetZIndex(sprite uintptr) int16

//go:linkname spriteSetTag pd_sprite_setTag
func spriteSetTag(sprite uintptr, tag uint8)

//go:linkname spriteGetTag pd_sprite_getTag
func spriteGetTag(sprite uintptr) uint8

//go:linkname spriteSetVisible pd_sprite_setVisible
func spriteSetVisible(sprite uintptr, visible int32)

//go:linkname spriteIsVisible pd_sprite_isVisible
func spriteIsVisible(sprite uintptr) int32

//go:linkname spriteSetOpaque pd_sprite_setOpaque
func spriteSetOpaque(sprite uintptr, opaque int32)

//go:linkname spriteSetDrawMode pd_sprite_setDrawMode
func spriteSetDrawMode(sprite uintptr, mode int32)

//go:linkname spriteSetImageFlip pd_sprite_setImageFlip
func spriteSetImageFlip(sprite uintptr, flip int32)

//go:linkname spriteGetImageFlip pd_sprite_getImageFlip
func spriteGetImageFlip(sprite uintptr) int32

//go:linkname spriteSetUpdatesEnabled pd_sprite_setUpdatesEnabled
func spriteSetUpdatesEnabled(sprite uintptr, flag int32)

//go:linkname spriteSetCollideRect pd_sprite_setCollideRect
func spriteSetCollideRect(sprite uintptr, x, y, w, h float32)

//go:linkname spriteGetCollideRect pd_sprite_getCollideRect
func spriteGetCollideRect(sprite uintptr, x, y, w, h *float32)

//go:linkname spriteClearCollideRect pd_sprite_clearCollideRect
func spriteClearCollideRect(sprite uintptr)

//go:linkname spriteSetCollisionsEnabled pd_sprite_setCollisionsEnabled
func spriteSetCollisionsEnabled(sprite uintptr, flag int32)

//go:linkname spriteMoveWithCollisions pd_sprite_moveWithCollisions
func spriteMoveWithCollisions(sprite uintptr, goalX, goalY float32, actualX, actualY *float32, count *int32) uintptr

//go:linkname spriteCheckCollisions pd_sprite_checkCollisions
func spriteCheckCollisions(sprite uintptr, count *int32) uintptr

//go:linkname spriteUpdateAndDrawSprites pd_sprite_updateAndDrawSprites
func spriteUpdateAndDrawSprites()

//go:linkname spriteSetAlwaysRedraw pd_sprite_setAlwaysRedraw
func spriteSetAlwaysRedraw(flag int32)

//go:linkname spriteResetCollisionWorld pd_sprite_resetCollisionWorld
func spriteResetCollisionWorld()

//go:linkname spriteQuerySpritesAtPoint pd_sprite_querySpritesAtPoint
func spriteQuerySpritesAtPoint(x, y float32, count *int32) uintptr

//go:linkname spriteQuerySpritesInRect pd_sprite_querySpritesInRect
func spriteQuerySpritesInRect(x, y, w, h float32, count *int32) uintptr

//go:linkname spriteQuerySpritesAlongLine pd_sprite_querySpritesAlongLine
func spriteQuerySpritesAlongLine(x1, y1, x2, y2 float32, count *int32) uintptr

//go:linkname spriteAllOverlappingSprites pd_sprite_allOverlappingSprites
func spriteAllOverlappingSprites(count *int32) uintptr

//go:linkname spriteDrawSprites pd_sprite_drawSprites
func spriteDrawSprites()

//go:linkname spriteMarkDirty pd_sprite_markDirty
func spriteMarkDirty(sprite uintptr)

//go:linkname spriteSetUpdateFunction pd_sprite_setUpdateFunction
func spriteSetUpdateFunction(sprite uintptr, hasCallback int32)

//go:linkname spriteSetDrawFunction pd_sprite_setDrawFunction
func spriteSetDrawFunction(sprite uintptr, hasCallback int32)

//go:linkname spriteSetCollisionResponseFunction pd_sprite_setCollisionResponseFunction
func spriteSetCollisionResponseFunction(sprite uintptr, hasCallback int32)

// Wrapper functions to convert uintptr to int32 for C calls
func spriteSetUpdateFunctionWrapper(sprite, hasCallback uintptr) {
	spriteSetUpdateFunction(sprite, int32(hasCallback))
}

func spriteSetDrawFunctionWrapper(sprite, hasCallback uintptr) {
	spriteSetDrawFunction(sprite, int32(hasCallback))
}

func spriteSetCollisionResponseFunctionWrapper(sprite, hasCallback uintptr) {
	spriteSetCollisionResponseFunction(sprite, int32(hasCallback))
}

// ============== File Functions ==============

//go:linkname fileOpen pd_file_open
func fileOpen(path *byte, mode int32) uintptr

//go:linkname fileClose pd_file_close
func fileClose(file uintptr) int32

//go:linkname fileRead pd_file_read
func fileRead(file uintptr, buf *byte, len uint32) int32

//go:linkname fileWrite pd_file_write
func fileWrite(file uintptr, buf *byte, len uint32) int32

//go:linkname fileFlush pd_file_flush
func fileFlush(file uintptr) int32

//go:linkname fileTell pd_file_tell
func fileTell(file uintptr) int32

//go:linkname fileSeek pd_file_seek
func fileSeek(file uintptr, pos, whence int32) int32

//go:linkname fileStat pd_file_stat
func fileStat(path *byte, isDir, size, mtime *int32) int32

//go:linkname fileMkdir pd_file_mkdir
func fileMkdir(path *byte) int32

//go:linkname fileUnlink pd_file_unlink
func fileUnlink(path *byte, recursive int32) int32

//go:linkname fileRename pd_file_rename
func fileRename(from, to *byte) int32

// ============== Sound - FilePlayer Functions ==============

//go:linkname soundNewFilePlayer pd_sound_newFilePlayer
func soundNewFilePlayer() uintptr

//go:linkname soundFreeFilePlayer pd_sound_freeFilePlayer
func soundFreeFilePlayer(player uintptr)

//go:linkname soundLoadIntoFilePlayer pd_sound_loadIntoFilePlayer
func soundLoadIntoFilePlayer(player uintptr, path *byte) int32

//go:linkname soundPlayFilePlayer pd_sound_playFilePlayer
func soundPlayFilePlayer(player uintptr, repeat int32)

//go:linkname soundStopFilePlayer pd_sound_stopFilePlayer
func soundStopFilePlayer(player uintptr)

//go:linkname soundPauseFilePlayer pd_sound_pauseFilePlayer
func soundPauseFilePlayer(player uintptr)

//go:linkname soundIsFilePlayerPlaying pd_sound_isFilePlayerPlaying
func soundIsFilePlayerPlaying(player uintptr) int32

//go:linkname soundSetFilePlayerVolume pd_sound_setFilePlayerVolume
func soundSetFilePlayerVolume(player uintptr, left, right float32)

//go:linkname soundGetFilePlayerVolume pd_sound_getFilePlayerVolume
func soundGetFilePlayerVolume(player uintptr, left, right *float32)

//go:linkname soundGetFilePlayerLength pd_sound_getFilePlayerLength
func soundGetFilePlayerLength(player uintptr) float32

//go:linkname soundSetFilePlayerOffset pd_sound_setFilePlayerOffset
func soundSetFilePlayerOffset(player uintptr, offset float32)

//go:linkname soundGetFilePlayerOffset pd_sound_getFilePlayerOffset
func soundGetFilePlayerOffset(player uintptr) float32

//go:linkname soundSetFilePlayerRate pd_sound_setFilePlayerRate
func soundSetFilePlayerRate(player uintptr, rate float32)

// ============== Sound - SamplePlayer Functions ==============

//go:linkname soundNewSamplePlayer pd_sound_newSamplePlayer
func soundNewSamplePlayer() uintptr

//go:linkname soundFreeSamplePlayer pd_sound_freeSamplePlayer
func soundFreeSamplePlayer(player uintptr)

//go:linkname soundSetSamplePlayerSample pd_sound_setSamplePlayerSample
func soundSetSamplePlayerSample(player, sample uintptr)

//go:linkname soundPlaySamplePlayer pd_sound_playSamplePlayer
func soundPlaySamplePlayer(player uintptr, repeat int32, rate float32)

//go:linkname soundStopSamplePlayer pd_sound_stopSamplePlayer
func soundStopSamplePlayer(player uintptr)

//go:linkname soundIsSamplePlayerPlaying pd_sound_isSamplePlayerPlaying
func soundIsSamplePlayerPlaying(player uintptr) int32

//go:linkname soundSetSamplePlayerVolume pd_sound_setSamplePlayerVolume
func soundSetSamplePlayerVolume(player uintptr, left, right float32)

// ============== Sound - Sample Functions ==============

//go:linkname soundNewSample pd_sound_newSample
func soundNewSample(length int32) uintptr

//go:linkname soundLoadSample pd_sound_loadSample
func soundLoadSample(path *byte) uintptr

//go:linkname soundFreeSample pd_sound_freeSample
func soundFreeSample(sample uintptr)

// ============== Sound - Global Functions ==============

//go:linkname soundGetHeadphoneState pd_sound_getHeadphoneState
func soundGetHeadphoneState(headphone, mic *int32)

//go:linkname soundSetOutputsActive pd_sound_setOutputsActive
func soundSetOutputsActive(headphone, speaker int32)

// ============== Sound - Synth Functions ==============

//go:linkname soundSynthNew pd_sound_synth_new
func soundSynthNew() uintptr

//go:linkname soundSynthFree pd_sound_synth_free
func soundSynthFree(synth uintptr)

//go:linkname soundSynthSetWaveform pd_sound_synth_setWaveform
func soundSynthSetWaveform(synth uintptr, wave int32)

//go:linkname soundSynthSetAttackTime pd_sound_synth_setAttackTime
func soundSynthSetAttackTime(synth uintptr, attack float32)

//go:linkname soundSynthSetDecayTime pd_sound_synth_setDecayTime
func soundSynthSetDecayTime(synth uintptr, decay float32)

//go:linkname soundSynthSetSustainLevel pd_sound_synth_setSustainLevel
func soundSynthSetSustainLevel(synth uintptr, sustain float32)

//go:linkname soundSynthSetReleaseTime pd_sound_synth_setReleaseTime
func soundSynthSetReleaseTime(synth uintptr, release float32)

//go:linkname soundSynthSetTranspose pd_sound_synth_setTranspose
func soundSynthSetTranspose(synth uintptr, halfSteps float32)

//go:linkname soundSynthPlayNote pd_sound_synth_playNote
func soundSynthPlayNote(synth uintptr, freq, vel, length float32, when uint32)

//go:linkname soundSynthPlayMIDINote pd_sound_synth_playMIDINote
func soundSynthPlayMIDINote(synth uintptr, note, vel, length float32, when uint32)

//go:linkname soundSynthNoteOff pd_sound_synth_noteOff
func soundSynthNoteOff(synth uintptr, when uint32)

//go:linkname soundSynthStop pd_sound_synth_stop
func soundSynthStop(synth uintptr)

//go:linkname soundSynthSetVolume pd_sound_synth_setVolume
func soundSynthSetVolume(synth uintptr, left, right float32)

//go:linkname soundSynthGetVolume pd_sound_synth_getVolume
func soundSynthGetVolume(synth uintptr, left, right *float32)

//go:linkname soundSynthIsPlaying pd_sound_synth_isPlaying
func soundSynthIsPlaying(synth uintptr) int32

// ============== Register Bridge ==============

func init() {
	pdgo.RegisterBridge(pdgo.Bridge{
		// System
		SysLog:                    sysLog,
		SysError:                  sysError,
		SysDrawFPS:                sysDrawFPS,
		SysGetCurrentTimeMS:       sysGetCurrentTimeMS,
		SysGetSecondsSinceEpoch:   sysGetSecondsSinceEpoch,
		SysGetButtonState:         sysGetButtonState,
		SysSetPeripheralsEnabled:  sysSetPeripheralsEnabled,
		SysGetAccelerometer:       sysGetAccelerometer,
		SysGetCrankChange:         sysGetCrankChange,
		SysGetCrankAngle:          sysGetCrankAngle,
		SysIsCrankDocked:          sysIsCrankDocked,
		SysSetCrankSoundsDisabled: sysSetCrankSoundsDisabled,
		SysGetFlipped:             sysGetFlipped,
		SysSetAutoLockDisabled:    sysSetAutoLockDisabled,
		SysGetLanguage:            sysGetLanguage,
		SysGetBatteryPercentage:   sysGetBatteryPercentage,
		SysGetBatteryVoltage:      sysGetBatteryVoltage,

		// Graphics
		GfxClear:              gfxClear,
		GfxSetBackgroundColor: gfxSetBackgroundColor,
		GfxSetDrawMode:        gfxSetDrawMode,
		GfxSetDrawOffset:      gfxSetDrawOffset,
		GfxSetClipRect:        gfxSetClipRect,
		GfxClearClipRect:      gfxClearClipRect,
		GfxSetLineCapStyle:    gfxSetLineCapStyle,
		GfxSetFont:            gfxSetFont,
		GfxSetTextTracking:    gfxSetTextTracking,
		GfxPushContext:        gfxPushContext,
		GfxPopContext:         gfxPopContext,
		GfxFillRect:           gfxFillRect,
		GfxDrawRect:           gfxDrawRect,
		GfxDrawLine:           gfxDrawLine,
		GfxFillTriangle:       gfxFillTriangle,
		GfxDrawEllipse:        gfxDrawEllipse,
		GfxFillEllipse:        gfxFillEllipse,
		GfxDrawText:           gfxDrawText,
		GfxGetTextWidth:       gfxGetTextWidth,
		GfxNewBitmap:          gfxNewBitmap,
		GfxFreeBitmap:         gfxFreeBitmap,
		GfxLoadBitmap:         gfxLoadBitmap,
		GfxCopyBitmap:         gfxCopyBitmap,
		GfxDrawBitmap:         gfxDrawBitmap,
		GfxTileBitmap:         gfxTileBitmap,
		GfxDrawScaledBitmap:   gfxDrawScaledBitmap,
		GfxDrawRotatedBitmap:  gfxDrawRotatedBitmap,
		GfxGetBitmapData:      gfxGetBitmapData,
		GfxClearBitmap:        gfxClearBitmap,
		GfxNewBitmapTable:     gfxNewBitmapTable,
		GfxFreeBitmapTable:    gfxFreeBitmapTable,
		GfxLoadBitmapTable:    gfxLoadBitmapTable,
		GfxGetTableBitmap:     gfxGetTableBitmap,
		GfxLoadFont:           gfxLoadFont,
		GfxGetFrame:           gfxGetFrame,
		GfxGetDisplayFrame:    gfxGetDisplayFrame,
		GfxMarkUpdatedRows:    gfxMarkUpdatedRows,
		GfxDisplay:            gfxDisplay,

		// Display
		DisplayGetWidth:       displayGetWidth,
		DisplayGetHeight:      displayGetHeight,
		DisplaySetRefreshRate: displaySetRefreshRate,
		DisplayGetRefreshRate: displayGetRefreshRate,
		DisplaySetInverted:    displaySetInverted,
		DisplaySetScale:       displaySetScale,
		DisplaySetMosaic:      displaySetMosaic,
		DisplaySetFlipped:     displaySetFlipped,
		DisplaySetOffset:      displaySetOffset,

		// Sprite
		SpriteNewSprite:             spriteNewSprite,
		SpriteFreeSprite:            spriteFreeSprite,
		SpriteAddSprite:             spriteAddSprite,
		SpriteRemoveSprite:          spriteRemoveSprite,
		SpriteRemoveAllSprites:      spriteRemoveAllSprites,
		SpriteGetSpriteCount:        spriteGetSpriteCount,
		SpriteSetImage:              spriteSetImage,
		SpriteGetImage:              spriteGetImage,
		SpriteSetBounds:             spriteSetBounds,
		SpriteGetBounds:             spriteGetBounds,
		SpriteMoveTo:                spriteMoveTo,
		SpriteMoveBy:                spriteMoveBy,
		SpriteGetPosition:           spriteGetPosition,
		SpriteSetZIndex:             spriteSetZIndex,
		SpriteGetZIndex:             spriteGetZIndex,
		SpriteSetTag:                spriteSetTag,
		SpriteGetTag:                spriteGetTag,
		SpriteSetVisible:            spriteSetVisible,
		SpriteIsVisible:             spriteIsVisible,
		SpriteSetOpaque:             spriteSetOpaque,
		SpriteSetDrawMode:           spriteSetDrawMode,
		SpriteSetImageFlip:          spriteSetImageFlip,
		SpriteGetImageFlip:          spriteGetImageFlip,
		SpriteSetUpdatesEnabled:     spriteSetUpdatesEnabled,
		SpriteSetCollideRect:        spriteSetCollideRect,
		SpriteGetCollideRect:        spriteGetCollideRect,
		SpriteClearCollideRect:      spriteClearCollideRect,
		SpriteSetCollisionsEnabled:  spriteSetCollisionsEnabled,
		SpriteMoveWithCollisions:    spriteMoveWithCollisions,
		SpriteCheckCollisions:       spriteCheckCollisions,
		SpriteUpdateAndDrawSprites:  spriteUpdateAndDrawSprites,
		SpriteSetAlwaysRedraw:       spriteSetAlwaysRedraw,
		SpriteResetCollisionWorld:   spriteResetCollisionWorld,
		SpriteQuerySpritesAtPoint:   spriteQuerySpritesAtPoint,
		SpriteQuerySpritesInRect:    spriteQuerySpritesInRect,
		SpriteQuerySpritesAlongLine: spriteQuerySpritesAlongLine,
		SpriteAllOverlappingSprites: spriteAllOverlappingSprites,
		SpriteDrawSprites:                  spriteDrawSprites,
		SpriteMarkDirty:                    spriteMarkDirty,
		SpriteSetUpdateFunction:            spriteSetUpdateFunctionWrapper,
		SpriteSetDrawFunction:              spriteSetDrawFunctionWrapper,
		SpriteSetCollisionResponseFunction: spriteSetCollisionResponseFunctionWrapper,

		// File
		FileOpen:   fileOpen,
		FileClose:  fileClose,
		FileRead:   fileRead,
		FileWrite:  fileWrite,
		FileFlush:  fileFlush,
		FileTell:   fileTell,
		FileSeek:   fileSeek,
		FileStat:   fileStat,
		FileMkdir:  fileMkdir,
		FileUnlink: fileUnlink,
		FileRename: fileRename,

		// Sound - FilePlayer
		SoundNewFilePlayer:       soundNewFilePlayer,
		SoundFreeFilePlayer:      soundFreeFilePlayer,
		SoundLoadIntoFilePlayer:  soundLoadIntoFilePlayer,
		SoundPlayFilePlayer:      soundPlayFilePlayer,
		SoundStopFilePlayer:      soundStopFilePlayer,
		SoundPauseFilePlayer:     soundPauseFilePlayer,
		SoundIsFilePlayerPlaying: soundIsFilePlayerPlaying,
		SoundSetFilePlayerVolume: soundSetFilePlayerVolume,
		SoundGetFilePlayerVolume: soundGetFilePlayerVolume,
		SoundGetFilePlayerLength: soundGetFilePlayerLength,
		SoundSetFilePlayerOffset: soundSetFilePlayerOffset,
		SoundGetFilePlayerOffset: soundGetFilePlayerOffset,
		SoundSetFilePlayerRate:   soundSetFilePlayerRate,

		// Sound - SamplePlayer
		SoundNewSamplePlayer:       soundNewSamplePlayer,
		SoundFreeSamplePlayer:      soundFreeSamplePlayer,
		SoundSetSamplePlayerSample: soundSetSamplePlayerSample,
		SoundPlaySamplePlayer:      soundPlaySamplePlayer,
		SoundStopSamplePlayer:      soundStopSamplePlayer,
		SoundIsSamplePlayerPlaying: soundIsSamplePlayerPlaying,
		SoundSetSamplePlayerVolume: soundSetSamplePlayerVolume,

		// Sound - Sample
		SoundNewSample:  soundNewSample,
		SoundLoadSample: soundLoadSample,
		SoundFreeSample: soundFreeSample,

		// Sound - Global
		SoundGetHeadphoneState: soundGetHeadphoneState,
		SoundSetOutputsActive:  soundSetOutputsActive,

		// Sound - Synth
		SoundSynthNew:          soundSynthNew,
		SoundSynthFree:         soundSynthFree,
		SoundSynthSetWaveform:  soundSynthSetWaveform,
		SoundSynthSetAttack:    soundSynthSetAttackTime,
		SoundSynthSetDecay:     soundSynthSetDecayTime,
		SoundSynthSetSustain:   soundSynthSetSustainLevel,
		SoundSynthSetRelease:   soundSynthSetReleaseTime,
		SoundSynthSetTranspose: soundSynthSetTranspose,
		SoundSynthPlayNote:     soundSynthPlayNote,
		SoundSynthPlayMIDINote: soundSynthPlayMIDINote,
		SoundSynthNoteOff:      soundSynthNoteOff,
		SoundSynthStop:         soundSynthStop,
		SoundSynthSetVolume:    soundSynthSetVolume,
		SoundSynthGetVolume:    soundSynthGetVolume,
		SoundSynthIsPlaying:    soundSynthIsPlaying,
	})
}`
