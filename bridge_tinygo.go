//go:build tinygo

// Bridge for TinyGo - allows calling C functions from pdgo package
// Main package must register the actual C function implementations via RegisterBridge()

package pdgo

import "unsafe"

// Helper for unsafe pointer conversion
func unsafePtr(p uintptr) unsafe.Pointer {
	return unsafe.Pointer(p)
}

// ============== System Bridge ==============

var (
	bridgeSysLog                    func(*byte)
	bridgeSysError                  func(*byte)
	bridgeSysDrawFPS                func(int32, int32)
	bridgeSysGetCurrentTimeMS       func() uint32
	bridgeSysGetSecondsSinceEpoch   func(*uint32) uint32
	bridgeSysGetButtonState         func(*uint32, *uint32, *uint32)
	bridgeSysSetPeripheralsEnabled  func(uint32)
	bridgeSysGetAccelerometer       func(*float32, *float32, *float32)
	bridgeSysGetCrankChange         func() float32
	bridgeSysGetCrankAngle          func() float32
	bridgeSysIsCrankDocked          func() int32
	bridgeSysSetCrankSoundsDisabled func(int32) int32
	bridgeSysGetFlipped             func() int32
	bridgeSysSetAutoLockDisabled    func(int32)
	bridgeSysGetLanguage            func() int32
	bridgeSysGetBatteryPercentage   func() float32
	bridgeSysGetBatteryVoltage      func() float32
)

// ============== Graphics Bridge ==============

var (
	bridgeGfxClear              func(uint32)
	bridgeGfxSetBackgroundColor func(int32)
	bridgeGfxSetDrawMode        func(int32) int32
	bridgeGfxSetDrawOffset      func(int32, int32)
	bridgeGfxSetClipRect        func(int32, int32, int32, int32)
	bridgeGfxClearClipRect      func()
	bridgeGfxSetLineCapStyle    func(int32)
	bridgeGfxSetFont            func(uintptr)
	bridgeGfxSetTextTracking    func(int32)
	bridgeGfxPushContext        func(uintptr)
	bridgeGfxPopContext         func()

	// Drawing
	bridgeGfxFillRect     func(int32, int32, int32, int32, uint32)
	bridgeGfxDrawRect     func(int32, int32, int32, int32, uint32)
	bridgeGfxDrawLine     func(int32, int32, int32, int32, int32, uint32)
	bridgeGfxFillTriangle func(int32, int32, int32, int32, int32, int32, uint32)
	bridgeGfxDrawEllipse  func(int32, int32, int32, int32, int32, float32, float32, uint32)
	bridgeGfxFillEllipse  func(int32, int32, int32, int32, float32, float32, uint32)
	bridgeGfxDrawText     func(*byte, int32, int32, int32, int32) int32
	bridgeGfxGetTextWidth func(uintptr, *byte, int32, int32, int32) int32

	// Bitmap
	bridgeGfxNewBitmap         func(int32, int32, uint32) uintptr
	bridgeGfxFreeBitmap        func(uintptr)
	bridgeGfxLoadBitmap        func(*byte) uintptr
	bridgeGfxCopyBitmap        func(uintptr) uintptr
	bridgeGfxDrawBitmap        func(uintptr, int32, int32, int32)
	bridgeGfxTileBitmap        func(uintptr, int32, int32, int32, int32, int32)
	bridgeGfxDrawScaledBitmap  func(uintptr, int32, int32, float32, float32)
	bridgeGfxDrawRotatedBitmap func(uintptr, int32, int32, float32, float32, float32, float32, float32)
	bridgeGfxGetBitmapData     func(uintptr, *int32, *int32, *int32, *uintptr, *uintptr)
	bridgeGfxClearBitmap       func(uintptr, uint32)

	// BitmapTable
	bridgeGfxNewBitmapTable  func(int32, int32, int32) uintptr
	bridgeGfxFreeBitmapTable func(uintptr)
	bridgeGfxLoadBitmapTable func(*byte) uintptr
	bridgeGfxGetTableBitmap  func(uintptr, int32) uintptr

	// Font
	bridgeGfxLoadFont func(*byte) uintptr

	// Frame
	bridgeGfxGetFrame        func() uintptr
	bridgeGfxGetDisplayFrame func() uintptr
	bridgeGfxMarkUpdatedRows func(int32, int32)
	bridgeGfxDisplay         func()
)

// ============== Display Bridge ==============

var (
	bridgeDisplayGetWidth       func() int32
	bridgeDisplayGetHeight      func() int32
	bridgeDisplaySetRefreshRate func(float32)
	bridgeDisplayGetRefreshRate func() float32
	bridgeDisplaySetInverted    func(int32)
	bridgeDisplaySetScale       func(uint32)
	bridgeDisplaySetMosaic      func(uint32, uint32)
	bridgeDisplaySetFlipped     func(int32, int32)
	bridgeDisplaySetOffset      func(int32, int32)
)

// ============== Sprite Bridge ==============

var (
	bridgeSpriteNewSprite                    func() uintptr
	bridgeSpriteFreeSprite                   func(uintptr)
	bridgeSpriteAddSprite                    func(uintptr)
	bridgeSpriteRemoveSprite                 func(uintptr)
	bridgeSpriteRemoveAllSprites             func()
	bridgeSpriteGetSpriteCount               func() int32
	bridgeSpriteSetImage                     func(uintptr, uintptr, int32)
	bridgeSpriteGetImage                     func(uintptr) uintptr
	bridgeSpriteSetBounds                    func(uintptr, float32, float32, float32, float32)
	bridgeSpriteGetBounds                    func(uintptr, *float32, *float32, *float32, *float32)
	bridgeSpriteMoveTo                       func(uintptr, float32, float32)
	bridgeSpriteMoveBy                       func(uintptr, float32, float32)
	bridgeSpriteGetPosition                  func(uintptr, *float32, *float32)
	bridgeSpriteSetZIndex                    func(uintptr, int16)
	bridgeSpriteGetZIndex                    func(uintptr) int16
	bridgeSpriteSetTag                       func(uintptr, uint8)
	bridgeSpriteGetTag                       func(uintptr) uint8
	bridgeSpriteSetVisible                   func(uintptr, int32)
	bridgeSpriteIsVisible                    func(uintptr) int32
	bridgeSpriteSetOpaque                    func(uintptr, int32)
	bridgeSpriteSetDrawMode                  func(uintptr, int32)
	bridgeSpriteSetImageFlip                 func(uintptr, int32)
	bridgeSpriteGetImageFlip                 func(uintptr) int32
	bridgeSpriteSetUpdatesEnabled            func(uintptr, int32)
	bridgeSpriteMarkDirty                    func(uintptr)
	bridgeSpriteSetCollideRect               func(uintptr, float32, float32, float32, float32)
	bridgeSpriteGetCollideRect               func(uintptr, *float32, *float32, *float32, *float32)
	bridgeSpriteClearCollideRect             func(uintptr)
	bridgeSpriteSetCollisionsEnabled         func(uintptr, int32)
	bridgeSpriteMoveWithCollisions           func(uintptr, float32, float32, *float32, *float32, *int32) uintptr
	bridgeSpriteCheckCollisions              func(uintptr, *int32) uintptr
	bridgeSpriteDrawSprites                  func()
	bridgeSpriteUpdateAndDrawSprites         func()
	bridgeSpriteSetAlwaysRedraw              func(int32)
	bridgeSpriteResetCollisionWorld          func()
	bridgeSpriteQuerySpritesAtPoint          func(float32, float32, *int32) uintptr
	bridgeSpriteQuerySpritesInRect           func(float32, float32, float32, float32, *int32) uintptr
	bridgeSpriteQuerySpritesAlongLine        func(float32, float32, float32, float32, *int32) uintptr
	bridgeSpriteAllOverlappingSprites        func(*int32) uintptr
	bridgeSpriteSetUpdateFunction            func(uintptr, uintptr)
	bridgeSpriteSetDrawFunction              func(uintptr, uintptr)
	bridgeSpriteSetCollisionResponseFunction func(uintptr, uintptr)
)

// ============== File Bridge ==============

var (
	bridgeFileOpen   func(*byte, int32) uintptr
	bridgeFileClose  func(uintptr) int32
	bridgeFileRead   func(uintptr, *byte, uint32) int32
	bridgeFileWrite  func(uintptr, *byte, uint32) int32
	bridgeFileFlush  func(uintptr) int32
	bridgeFileTell   func(uintptr) int32
	bridgeFileSeek   func(uintptr, int32, int32) int32
	bridgeFileStat   func(*byte, *int32, *int32, *int32) int32
	bridgeFileMkdir  func(*byte) int32
	bridgeFileUnlink func(*byte, int32) int32
	bridgeFileRename func(*byte, *byte) int32
)

// ============== Sound Bridge ==============

var (
	// FilePlayer
	bridgeSoundNewFilePlayer       func() uintptr
	bridgeSoundFreeFilePlayer      func(uintptr)
	bridgeSoundLoadIntoFilePlayer  func(uintptr, *byte) int32
	bridgeSoundPlayFilePlayer      func(uintptr, int32)
	bridgeSoundStopFilePlayer      func(uintptr)
	bridgeSoundPauseFilePlayer     func(uintptr)
	bridgeSoundIsFilePlayerPlaying func(uintptr) int32
	bridgeSoundSetFilePlayerVolume func(uintptr, float32, float32)
	bridgeSoundGetFilePlayerVolume func(uintptr, *float32, *float32)
	bridgeSoundGetFilePlayerLength func(uintptr) float32
	bridgeSoundSetFilePlayerOffset func(uintptr, float32)
	bridgeSoundGetFilePlayerOffset func(uintptr) float32
	bridgeSoundSetFilePlayerRate   func(uintptr, float32)

	// SamplePlayer
	bridgeSoundNewSamplePlayer       func() uintptr
	bridgeSoundFreeSamplePlayer      func(uintptr)
	bridgeSoundSetSamplePlayerSample func(uintptr, uintptr)
	bridgeSoundPlaySamplePlayer      func(uintptr, int32, float32)
	bridgeSoundStopSamplePlayer      func(uintptr)
	bridgeSoundIsSamplePlayerPlaying func(uintptr) int32
	bridgeSoundSetSamplePlayerVolume func(uintptr, float32, float32)

	// Sample
	bridgeSoundNewSample  func(int32) uintptr
	bridgeSoundLoadSample func(*byte) uintptr
	bridgeSoundFreeSample func(uintptr)

	// Global
	bridgeSoundGetHeadphoneState func(*int32, *int32)
	bridgeSoundSetOutputsActive  func(int32, int32)
)

// ============== Lua Bridge ==============

var (
	bridgeLuaGetArgCount  func() int32
	bridgeLuaGetArgType   func(int32) int32
	bridgeLuaArgIsNil     func(int32) int32
	bridgeLuaGetArgBool   func(int32) int32
	bridgeLuaGetArgInt    func(int32) int32
	bridgeLuaGetArgFloat  func(int32) float32
	bridgeLuaGetArgString func(int32) uintptr
	bridgeLuaPushNil      func()
	bridgeLuaPushBool     func(int32)
	bridgeLuaPushInt      func(int32)
	bridgeLuaPushFloat    func(float32)
	bridgeLuaPushString   func(*byte)
)

// ============== JSON Bridge ==============

var (
	bridgeJSONDecode func(*byte) int32
)

// ============== Scoreboards Bridge ==============

var (
	bridgeScoreboardsAddScore        func(*byte, uint32, uintptr)
	bridgeScoreboardsGetPersonalBest func(*byte, uintptr)
	bridgeScoreboardsFreeScore       func(uintptr)
	bridgeScoreboardsGetScoreboards  func(uintptr)
	bridgeScoreboardsGetScores       func(*byte, uintptr)
)

// ============== Bridge Registration ==============

// Bridge holds all function pointers to C runtime functions
type Bridge struct {
	// System
	SysLog                    func(*byte)
	SysError                  func(*byte)
	SysDrawFPS                func(int32, int32)
	SysGetCurrentTimeMS       func() uint32
	SysGetSecondsSinceEpoch   func(*uint32) uint32
	SysGetButtonState         func(*uint32, *uint32, *uint32)
	SysSetPeripheralsEnabled  func(uint32)
	SysGetAccelerometer       func(*float32, *float32, *float32)
	SysGetCrankChange         func() float32
	SysGetCrankAngle          func() float32
	SysIsCrankDocked          func() int32
	SysSetCrankSoundsDisabled func(int32) int32
	SysGetFlipped             func() int32
	SysSetAutoLockDisabled    func(int32)
	SysGetLanguage            func() int32
	SysGetBatteryPercentage   func() float32
	SysGetBatteryVoltage      func() float32

	// Graphics
	GfxClear              func(uint32)
	GfxSetBackgroundColor func(int32)
	GfxSetDrawMode        func(int32) int32
	GfxSetDrawOffset      func(int32, int32)
	GfxSetClipRect        func(int32, int32, int32, int32)
	GfxClearClipRect      func()
	GfxSetLineCapStyle    func(int32)
	GfxSetFont            func(uintptr)
	GfxSetTextTracking    func(int32)
	GfxPushContext        func(uintptr)
	GfxPopContext         func()
	GfxFillRect           func(int32, int32, int32, int32, uint32)
	GfxDrawRect           func(int32, int32, int32, int32, uint32)
	GfxDrawLine           func(int32, int32, int32, int32, int32, uint32)
	GfxFillTriangle       func(int32, int32, int32, int32, int32, int32, uint32)
	GfxDrawEllipse        func(int32, int32, int32, int32, int32, float32, float32, uint32)
	GfxFillEllipse        func(int32, int32, int32, int32, float32, float32, uint32)
	GfxDrawText           func(*byte, int32, int32, int32, int32) int32
	GfxGetTextWidth       func(uintptr, *byte, int32, int32, int32) int32
	GfxNewBitmap          func(int32, int32, uint32) uintptr
	GfxFreeBitmap         func(uintptr)
	GfxLoadBitmap         func(*byte) uintptr
	GfxCopyBitmap         func(uintptr) uintptr
	GfxDrawBitmap         func(uintptr, int32, int32, int32)
	GfxTileBitmap         func(uintptr, int32, int32, int32, int32, int32)
	GfxDrawScaledBitmap   func(uintptr, int32, int32, float32, float32)
	GfxDrawRotatedBitmap  func(uintptr, int32, int32, float32, float32, float32, float32, float32)
	GfxGetBitmapData      func(uintptr, *int32, *int32, *int32, *uintptr, *uintptr)
	GfxClearBitmap        func(uintptr, uint32)
	GfxNewBitmapTable     func(int32, int32, int32) uintptr
	GfxFreeBitmapTable    func(uintptr)
	GfxLoadBitmapTable    func(*byte) uintptr
	GfxGetTableBitmap     func(uintptr, int32) uintptr
	GfxLoadFont           func(*byte) uintptr
	GfxGetFrame           func() uintptr
	GfxGetDisplayFrame    func() uintptr
	GfxMarkUpdatedRows    func(int32, int32)
	GfxDisplay            func()

	// Display
	DisplayGetWidth       func() int32
	DisplayGetHeight      func() int32
	DisplaySetRefreshRate func(float32)
	DisplayGetRefreshRate func() float32
	DisplaySetInverted    func(int32)
	DisplaySetScale       func(uint32)
	DisplaySetMosaic      func(uint32, uint32)
	DisplaySetFlipped     func(int32, int32)
	DisplaySetOffset      func(int32, int32)

	// Sprite
	SpriteNewSprite                    func() uintptr
	SpriteFreeSprite                   func(uintptr)
	SpriteAddSprite                    func(uintptr)
	SpriteRemoveSprite                 func(uintptr)
	SpriteRemoveAllSprites             func()
	SpriteGetSpriteCount               func() int32
	SpriteSetImage                     func(uintptr, uintptr, int32)
	SpriteGetImage                     func(uintptr) uintptr
	SpriteSetBounds                    func(uintptr, float32, float32, float32, float32)
	SpriteGetBounds                    func(uintptr, *float32, *float32, *float32, *float32)
	SpriteMoveTo                       func(uintptr, float32, float32)
	SpriteMoveBy                       func(uintptr, float32, float32)
	SpriteGetPosition                  func(uintptr, *float32, *float32)
	SpriteSetZIndex                    func(uintptr, int16)
	SpriteGetZIndex                    func(uintptr) int16
	SpriteSetTag                       func(uintptr, uint8)
	SpriteGetTag                       func(uintptr) uint8
	SpriteSetVisible                   func(uintptr, int32)
	SpriteIsVisible                    func(uintptr) int32
	SpriteSetOpaque                    func(uintptr, int32)
	SpriteSetDrawMode                  func(uintptr, int32)
	SpriteSetImageFlip                 func(uintptr, int32)
	SpriteGetImageFlip                 func(uintptr) int32
	SpriteSetUpdatesEnabled            func(uintptr, int32)
	SpriteMarkDirty                    func(uintptr)
	SpriteSetCollideRect               func(uintptr, float32, float32, float32, float32)
	SpriteGetCollideRect               func(uintptr, *float32, *float32, *float32, *float32)
	SpriteClearCollideRect             func(uintptr)
	SpriteSetCollisionsEnabled         func(uintptr, int32)
	SpriteMoveWithCollisions           func(uintptr, float32, float32, *float32, *float32, *int32) uintptr
	SpriteCheckCollisions              func(uintptr, *int32) uintptr
	SpriteDrawSprites                  func()
	SpriteUpdateAndDrawSprites         func()
	SpriteSetAlwaysRedraw              func(int32)
	SpriteResetCollisionWorld          func()
	SpriteQuerySpritesAtPoint          func(float32, float32, *int32) uintptr
	SpriteQuerySpritesInRect           func(float32, float32, float32, float32, *int32) uintptr
	SpriteQuerySpritesAlongLine        func(float32, float32, float32, float32, *int32) uintptr
	SpriteAllOverlappingSprites        func(*int32) uintptr
	SpriteSetUpdateFunction            func(uintptr, uintptr)
	SpriteSetDrawFunction              func(uintptr, uintptr)
	SpriteSetCollisionResponseFunction func(uintptr, uintptr)

	// File
	FileOpen   func(*byte, int32) uintptr
	FileClose  func(uintptr) int32
	FileRead   func(uintptr, *byte, uint32) int32
	FileWrite  func(uintptr, *byte, uint32) int32
	FileFlush  func(uintptr) int32
	FileTell   func(uintptr) int32
	FileSeek   func(uintptr, int32, int32) int32
	FileStat   func(*byte, *int32, *int32, *int32) int32
	FileMkdir  func(*byte) int32
	FileUnlink func(*byte, int32) int32
	FileRename func(*byte, *byte) int32

	// Sound - FilePlayer
	SoundNewFilePlayer       func() uintptr
	SoundFreeFilePlayer      func(uintptr)
	SoundLoadIntoFilePlayer  func(uintptr, *byte) int32
	SoundPlayFilePlayer      func(uintptr, int32)
	SoundStopFilePlayer      func(uintptr)
	SoundPauseFilePlayer     func(uintptr)
	SoundIsFilePlayerPlaying func(uintptr) int32
	SoundSetFilePlayerVolume func(uintptr, float32, float32)
	SoundGetFilePlayerVolume func(uintptr, *float32, *float32)
	SoundGetFilePlayerLength func(uintptr) float32
	SoundSetFilePlayerOffset func(uintptr, float32)
	SoundGetFilePlayerOffset func(uintptr) float32
	SoundSetFilePlayerRate   func(uintptr, float32)

	// Sound - SamplePlayer
	SoundNewSamplePlayer       func() uintptr
	SoundFreeSamplePlayer      func(uintptr)
	SoundSetSamplePlayerSample func(uintptr, uintptr)
	SoundPlaySamplePlayer      func(uintptr, int32, float32)
	SoundStopSamplePlayer      func(uintptr)
	SoundIsSamplePlayerPlaying func(uintptr) int32
	SoundSetSamplePlayerVolume func(uintptr, float32, float32)

	// Sound - Sample
	SoundNewSample  func(int32) uintptr
	SoundLoadSample func(*byte) uintptr
	SoundFreeSample func(uintptr)

	// Sound - Global
	SoundGetHeadphoneState func(*int32, *int32)
	SoundSetOutputsActive  func(int32, int32)

	// Lua
	LuaGetArgCount  func() int32
	LuaGetArgType   func(int32) int32
	LuaArgIsNil     func(int32) int32
	LuaGetArgBool   func(int32) int32
	LuaGetArgInt    func(int32) int32
	LuaGetArgFloat  func(int32) float32
	LuaGetArgString func(int32) uintptr
	LuaPushNil      func()
	LuaPushBool     func(int32)
	LuaPushInt      func(int32)
	LuaPushFloat    func(float32)
	LuaPushString   func(*byte)

	// JSON
	JSONDecode func(*byte) int32

	// Scoreboards
	ScoreboardsAddScore        func(*byte, uint32, uintptr)
	ScoreboardsGetPersonalBest func(*byte, uintptr)
	ScoreboardsFreeScore       func(uintptr)
	ScoreboardsGetScoreboards  func(uintptr)
	ScoreboardsGetScores       func(*byte, uintptr)
}

// RegisterBridge registers C function implementations from main package
func RegisterBridge(b Bridge) {
	// System
	bridgeSysLog = b.SysLog
	bridgeSysError = b.SysError
	bridgeSysDrawFPS = b.SysDrawFPS
	bridgeSysGetCurrentTimeMS = b.SysGetCurrentTimeMS
	bridgeSysGetSecondsSinceEpoch = b.SysGetSecondsSinceEpoch
	bridgeSysGetButtonState = b.SysGetButtonState
	bridgeSysSetPeripheralsEnabled = b.SysSetPeripheralsEnabled
	bridgeSysGetAccelerometer = b.SysGetAccelerometer
	bridgeSysGetCrankChange = b.SysGetCrankChange
	bridgeSysGetCrankAngle = b.SysGetCrankAngle
	bridgeSysIsCrankDocked = b.SysIsCrankDocked
	bridgeSysSetCrankSoundsDisabled = b.SysSetCrankSoundsDisabled
	bridgeSysGetFlipped = b.SysGetFlipped
	bridgeSysSetAutoLockDisabled = b.SysSetAutoLockDisabled
	bridgeSysGetLanguage = b.SysGetLanguage
	bridgeSysGetBatteryPercentage = b.SysGetBatteryPercentage
	bridgeSysGetBatteryVoltage = b.SysGetBatteryVoltage

	// Graphics
	bridgeGfxClear = b.GfxClear
	bridgeGfxSetBackgroundColor = b.GfxSetBackgroundColor
	bridgeGfxSetDrawMode = b.GfxSetDrawMode
	bridgeGfxSetDrawOffset = b.GfxSetDrawOffset
	bridgeGfxSetClipRect = b.GfxSetClipRect
	bridgeGfxClearClipRect = b.GfxClearClipRect
	bridgeGfxSetLineCapStyle = b.GfxSetLineCapStyle
	bridgeGfxSetFont = b.GfxSetFont
	bridgeGfxSetTextTracking = b.GfxSetTextTracking
	bridgeGfxPushContext = b.GfxPushContext
	bridgeGfxPopContext = b.GfxPopContext
	bridgeGfxFillRect = b.GfxFillRect
	bridgeGfxDrawRect = b.GfxDrawRect
	bridgeGfxDrawLine = b.GfxDrawLine
	bridgeGfxFillTriangle = b.GfxFillTriangle
	bridgeGfxDrawEllipse = b.GfxDrawEllipse
	bridgeGfxFillEllipse = b.GfxFillEllipse
	bridgeGfxDrawText = b.GfxDrawText
	bridgeGfxGetTextWidth = b.GfxGetTextWidth
	bridgeGfxNewBitmap = b.GfxNewBitmap
	bridgeGfxFreeBitmap = b.GfxFreeBitmap
	bridgeGfxLoadBitmap = b.GfxLoadBitmap
	bridgeGfxCopyBitmap = b.GfxCopyBitmap
	bridgeGfxDrawBitmap = b.GfxDrawBitmap
	bridgeGfxTileBitmap = b.GfxTileBitmap
	bridgeGfxDrawScaledBitmap = b.GfxDrawScaledBitmap
	bridgeGfxDrawRotatedBitmap = b.GfxDrawRotatedBitmap
	bridgeGfxGetBitmapData = b.GfxGetBitmapData
	bridgeGfxClearBitmap = b.GfxClearBitmap
	bridgeGfxNewBitmapTable = b.GfxNewBitmapTable
	bridgeGfxFreeBitmapTable = b.GfxFreeBitmapTable
	bridgeGfxLoadBitmapTable = b.GfxLoadBitmapTable
	bridgeGfxGetTableBitmap = b.GfxGetTableBitmap
	bridgeGfxLoadFont = b.GfxLoadFont
	bridgeGfxGetFrame = b.GfxGetFrame
	bridgeGfxGetDisplayFrame = b.GfxGetDisplayFrame
	bridgeGfxMarkUpdatedRows = b.GfxMarkUpdatedRows
	bridgeGfxDisplay = b.GfxDisplay

	// Display
	bridgeDisplayGetWidth = b.DisplayGetWidth
	bridgeDisplayGetHeight = b.DisplayGetHeight
	bridgeDisplaySetRefreshRate = b.DisplaySetRefreshRate
	bridgeDisplayGetRefreshRate = b.DisplayGetRefreshRate
	bridgeDisplaySetInverted = b.DisplaySetInverted
	bridgeDisplaySetScale = b.DisplaySetScale
	bridgeDisplaySetMosaic = b.DisplaySetMosaic
	bridgeDisplaySetFlipped = b.DisplaySetFlipped
	bridgeDisplaySetOffset = b.DisplaySetOffset

	// Sprite
	bridgeSpriteNewSprite = b.SpriteNewSprite
	bridgeSpriteFreeSprite = b.SpriteFreeSprite
	bridgeSpriteAddSprite = b.SpriteAddSprite
	bridgeSpriteRemoveSprite = b.SpriteRemoveSprite
	bridgeSpriteRemoveAllSprites = b.SpriteRemoveAllSprites
	bridgeSpriteGetSpriteCount = b.SpriteGetSpriteCount
	bridgeSpriteSetImage = b.SpriteSetImage
	bridgeSpriteGetImage = b.SpriteGetImage
	bridgeSpriteSetBounds = b.SpriteSetBounds
	bridgeSpriteGetBounds = b.SpriteGetBounds
	bridgeSpriteMoveTo = b.SpriteMoveTo
	bridgeSpriteMoveBy = b.SpriteMoveBy
	bridgeSpriteGetPosition = b.SpriteGetPosition
	bridgeSpriteSetZIndex = b.SpriteSetZIndex
	bridgeSpriteGetZIndex = b.SpriteGetZIndex
	bridgeSpriteSetTag = b.SpriteSetTag
	bridgeSpriteGetTag = b.SpriteGetTag
	bridgeSpriteSetVisible = b.SpriteSetVisible
	bridgeSpriteIsVisible = b.SpriteIsVisible
	bridgeSpriteSetOpaque = b.SpriteSetOpaque
	bridgeSpriteSetDrawMode = b.SpriteSetDrawMode
	bridgeSpriteSetImageFlip = b.SpriteSetImageFlip
	bridgeSpriteGetImageFlip = b.SpriteGetImageFlip
	bridgeSpriteSetUpdatesEnabled = b.SpriteSetUpdatesEnabled
	bridgeSpriteMarkDirty = b.SpriteMarkDirty
	bridgeSpriteSetCollideRect = b.SpriteSetCollideRect
	bridgeSpriteGetCollideRect = b.SpriteGetCollideRect
	bridgeSpriteClearCollideRect = b.SpriteClearCollideRect
	bridgeSpriteSetCollisionsEnabled = b.SpriteSetCollisionsEnabled
	bridgeSpriteMoveWithCollisions = b.SpriteMoveWithCollisions
	bridgeSpriteCheckCollisions = b.SpriteCheckCollisions
	bridgeSpriteDrawSprites = b.SpriteDrawSprites
	bridgeSpriteUpdateAndDrawSprites = b.SpriteUpdateAndDrawSprites
	bridgeSpriteSetAlwaysRedraw = b.SpriteSetAlwaysRedraw
	bridgeSpriteResetCollisionWorld = b.SpriteResetCollisionWorld
	bridgeSpriteQuerySpritesAtPoint = b.SpriteQuerySpritesAtPoint
	bridgeSpriteQuerySpritesInRect = b.SpriteQuerySpritesInRect
	bridgeSpriteQuerySpritesAlongLine = b.SpriteQuerySpritesAlongLine
	bridgeSpriteAllOverlappingSprites = b.SpriteAllOverlappingSprites
	bridgeSpriteSetUpdateFunction = b.SpriteSetUpdateFunction
	bridgeSpriteSetDrawFunction = b.SpriteSetDrawFunction
	bridgeSpriteSetCollisionResponseFunction = b.SpriteSetCollisionResponseFunction

	// File
	bridgeFileOpen = b.FileOpen
	bridgeFileClose = b.FileClose
	bridgeFileRead = b.FileRead
	bridgeFileWrite = b.FileWrite
	bridgeFileFlush = b.FileFlush
	bridgeFileTell = b.FileTell
	bridgeFileSeek = b.FileSeek
	bridgeFileStat = b.FileStat
	bridgeFileMkdir = b.FileMkdir
	bridgeFileUnlink = b.FileUnlink
	bridgeFileRename = b.FileRename

	// Sound - FilePlayer
	bridgeSoundNewFilePlayer = b.SoundNewFilePlayer
	bridgeSoundFreeFilePlayer = b.SoundFreeFilePlayer
	bridgeSoundLoadIntoFilePlayer = b.SoundLoadIntoFilePlayer
	bridgeSoundPlayFilePlayer = b.SoundPlayFilePlayer
	bridgeSoundStopFilePlayer = b.SoundStopFilePlayer
	bridgeSoundPauseFilePlayer = b.SoundPauseFilePlayer
	bridgeSoundIsFilePlayerPlaying = b.SoundIsFilePlayerPlaying
	bridgeSoundSetFilePlayerVolume = b.SoundSetFilePlayerVolume
	bridgeSoundGetFilePlayerVolume = b.SoundGetFilePlayerVolume
	bridgeSoundGetFilePlayerLength = b.SoundGetFilePlayerLength
	bridgeSoundSetFilePlayerOffset = b.SoundSetFilePlayerOffset
	bridgeSoundGetFilePlayerOffset = b.SoundGetFilePlayerOffset
	bridgeSoundSetFilePlayerRate = b.SoundSetFilePlayerRate

	// Sound - SamplePlayer
	bridgeSoundNewSamplePlayer = b.SoundNewSamplePlayer
	bridgeSoundFreeSamplePlayer = b.SoundFreeSamplePlayer
	bridgeSoundSetSamplePlayerSample = b.SoundSetSamplePlayerSample
	bridgeSoundPlaySamplePlayer = b.SoundPlaySamplePlayer
	bridgeSoundStopSamplePlayer = b.SoundStopSamplePlayer
	bridgeSoundIsSamplePlayerPlaying = b.SoundIsSamplePlayerPlaying
	bridgeSoundSetSamplePlayerVolume = b.SoundSetSamplePlayerVolume

	// Sound - Sample
	bridgeSoundNewSample = b.SoundNewSample
	bridgeSoundLoadSample = b.SoundLoadSample
	bridgeSoundFreeSample = b.SoundFreeSample

	// Sound - Global
	bridgeSoundGetHeadphoneState = b.SoundGetHeadphoneState
	bridgeSoundSetOutputsActive = b.SoundSetOutputsActive

	// Lua
	bridgeLuaGetArgCount = b.LuaGetArgCount
	bridgeLuaGetArgType = b.LuaGetArgType
	bridgeLuaArgIsNil = b.LuaArgIsNil
	bridgeLuaGetArgBool = b.LuaGetArgBool
	bridgeLuaGetArgInt = b.LuaGetArgInt
	bridgeLuaGetArgFloat = b.LuaGetArgFloat
	bridgeLuaGetArgString = b.LuaGetArgString
	bridgeLuaPushNil = b.LuaPushNil
	bridgeLuaPushBool = b.LuaPushBool
	bridgeLuaPushInt = b.LuaPushInt
	bridgeLuaPushFloat = b.LuaPushFloat
	bridgeLuaPushString = b.LuaPushString

	// JSON
	bridgeJSONDecode = b.JSONDecode

	// Scoreboards
	bridgeScoreboardsAddScore = b.ScoreboardsAddScore
	bridgeScoreboardsGetPersonalBest = b.ScoreboardsGetPersonalBest
	bridgeScoreboardsFreeScore = b.ScoreboardsFreeScore
	bridgeScoreboardsGetScoreboards = b.ScoreboardsGetScoreboards
	bridgeScoreboardsGetScores = b.ScoreboardsGetScores
}
