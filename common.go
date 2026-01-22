//go:build tinygo

// Package pdgo provides Go bindings for the Playdate SDK C API.
// This file contains common types and constants for TinyGo builds.

package pdgo

// Display constants
const (
	LCDColumns = 400
	LCDRows    = 240
	LCDRowSize = 52 // Bytes per row (400 pixels / 8 bits + padding)
)

// PDSystemEvent represents system events sent to the game
type PDSystemEvent int32

const (
	EventInit          PDSystemEvent = 0
	EventInitLua       PDSystemEvent = 1
	EventLock          PDSystemEvent = 2
	EventUnlock        PDSystemEvent = 3
	EventPause         PDSystemEvent = 4
	EventResume        PDSystemEvent = 5
	EventTerminate     PDSystemEvent = 6
	EventKeyPressed    PDSystemEvent = 7
	EventKeyReleased   PDSystemEvent = 8
	EventLowPower      PDSystemEvent = 9
	EventMirrorStarted PDSystemEvent = 10
	EventMirrorEnded   PDSystemEvent = 11
)

// String returns the string representation of the event
func (e PDSystemEvent) String() string {
	switch e {
	case EventInit:
		return "Init"
	case EventInitLua:
		return "InitLua"
	case EventLock:
		return "Lock"
	case EventUnlock:
		return "Unlock"
	case EventPause:
		return "Pause"
	case EventResume:
		return "Resume"
	case EventTerminate:
		return "Terminate"
	case EventKeyPressed:
		return "KeyPressed"
	case EventKeyReleased:
		return "KeyReleased"
	case EventLowPower:
		return "LowPower"
	case EventMirrorStarted:
		return "MirrorStarted"
	case EventMirrorEnded:
		return "MirrorEnded"
	default:
		return "Unknown"
	}
}

// PDButtons represents button state
type PDButtons uint32

const (
	ButtonLeft  PDButtons = 1 << 0
	ButtonRight PDButtons = 1 << 1
	ButtonUp    PDButtons = 1 << 2
	ButtonDown  PDButtons = 1 << 3
	ButtonB     PDButtons = 1 << 4
	ButtonA     PDButtons = 1 << 5
)

// PDPeripherals for enabling peripherals
type PDPeripherals uint32

const (
	PeripheralNone          PDPeripherals = 0
	PeripheralAccelerometer PDPeripherals = 1 << 0
)

// LCDSolidColor represents solid colors
type LCDSolidColor int32

const (
	ColorBlack LCDSolidColor = 0
	ColorWhite LCDSolidColor = 1
	ColorClear LCDSolidColor = 2
	ColorXOR   LCDSolidColor = 3
)

// LCDColor can be a solid color or a pattern pointer
type LCDColor uint32

// NewColorFromSolid creates an LCDColor from a solid color
func NewColorFromSolid(c LCDSolidColor) LCDColor {
	return LCDColor(c)
}

// LCDBitmapDrawMode for drawing modes
type LCDBitmapDrawMode int32

const (
	DrawModeCopy             LCDBitmapDrawMode = 0
	DrawModeWhiteTransparent LCDBitmapDrawMode = 1
	DrawModeBlackTransparent LCDBitmapDrawMode = 2
	DrawModeFillWhite        LCDBitmapDrawMode = 3
	DrawModeFillBlack        LCDBitmapDrawMode = 4
	DrawModeXOR              LCDBitmapDrawMode = 5
	DrawModeNXOR             LCDBitmapDrawMode = 6
	DrawModeInverted         LCDBitmapDrawMode = 7
)

// LCDBitmapFlip for bitmap flipping
type LCDBitmapFlip int32

const (
	BitmapUnflipped LCDBitmapFlip = 0
	BitmapFlippedX  LCDBitmapFlip = 1
	BitmapFlippedY  LCDBitmapFlip = 2
	BitmapFlippedXY LCDBitmapFlip = 3
)

// LCDLineCapStyle for line cap styles
type LCDLineCapStyle int32

const (
	LineCapStyleButt   LCDLineCapStyle = 0
	LineCapStyleSquare LCDLineCapStyle = 1
	LineCapStyleRound  LCDLineCapStyle = 2
)

// PDStringEncoding for text encoding
type PDStringEncoding int32

const (
	ASCIIEncoding  PDStringEncoding = 0
	UTF8Encoding   PDStringEncoding = 1
	Latin1Encoding PDStringEncoding = 2
)

// SpriteCollisionResponseType for collision responses
type SpriteCollisionResponseType int32

const (
	CollisionTypeSlide   SpriteCollisionResponseType = 0
	CollisionTypeFreeze  SpriteCollisionResponseType = 1
	CollisionTypeOverlap SpriteCollisionResponseType = 2
	CollisionTypeBounce  SpriteCollisionResponseType = 3
)

// PDLanguage represents system language
type PDLanguage int32

const (
	LanguageEnglish  PDLanguage = 0
	LanguageJapanese PDLanguage = 1
	LanguageUnknown  PDLanguage = 2
)

// BitmapData holds bitmap pixel data information
type BitmapData struct {
	Width    int
	Height   int
	RowBytes int
	Mask     []byte
	Data     []byte
}

// CollisionInfo contains collision detection results
type CollisionInfo struct {
	Sprite       *Sprite
	Other        *Sprite
	ResponseType SpriteCollisionResponseType
	Overlaps     bool
	Ti           float32
	Move         struct{ X, Y float32 }
	Normal       struct{ X, Y int32 }
	Touch        struct{ X, Y float32 }
	SpriteRect   struct{ X, Y, W, H float32 }
	OtherRect    struct{ X, Y, W, H float32 }
}

// SpriteQueryInfo for raycast results
type SpriteQueryInfo struct {
	Sprite     *Sprite
	Ti1        float32
	Ti2        float32
	EntryPoint struct{ X, Y float32 }
	ExitPoint  struct{ X, Y float32 }
}
