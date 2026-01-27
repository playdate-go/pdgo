// pdgo common types and constants

package pdgo

// Display constants
const (
	LCDColumns = 400
	LCDRows    = 240
	LCDRowSize = 52 // bytes per row (400 / 8 = 50, rounded up to 52)
)

// LCDColor represents a color value
type LCDColor uint32

// LCDSolidColor represents solid colors
type LCDSolidColor int32

const (
	ColorBlack LCDSolidColor = 0
	ColorWhite LCDSolidColor = 1
	ColorClear LCDSolidColor = 2
	ColorXOR   LCDSolidColor = 3
)

// Convenience color values
const (
	SolidBlack LCDColor = 0
	SolidWhite LCDColor = 1
)

// LCDBitmapDrawMode represents bitmap draw modes
type LCDBitmapDrawMode int32

const (
	DrawModeCopy        LCDBitmapDrawMode = 0
	DrawModeWhiteTransp LCDBitmapDrawMode = 1
	DrawModeBlackTransp LCDBitmapDrawMode = 2
	DrawModeFillWhite   LCDBitmapDrawMode = 3
	DrawModeFillBlack   LCDBitmapDrawMode = 4
	DrawModeXOR         LCDBitmapDrawMode = 5
	DrawModeNXOR        LCDBitmapDrawMode = 6
	DrawModeInverted    LCDBitmapDrawMode = 7
)

// LCDBitmapFlip represents bitmap flip options
type LCDBitmapFlip int32

const (
	BitmapUnflipped LCDBitmapFlip = 0
	BitmapFlippedX  LCDBitmapFlip = 1
	BitmapFlippedY  LCDBitmapFlip = 2
	BitmapFlippedXY LCDBitmapFlip = 3
)

// LCDLineCapStyle represents line cap styles
type LCDLineCapStyle int32

const (
	LineCapStyleButt   LCDLineCapStyle = 0
	LineCapStyleSquare LCDLineCapStyle = 1
	LineCapStyleRound  LCDLineCapStyle = 2
)

// Text encoding types
type PDStringEncoding int32

const (
	ASCIIEncoding  PDStringEncoding = 0
	UTF8Encoding   PDStringEncoding = 1
	Latin1Encoding PDStringEncoding = 2
)

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

// PDPeripherals represents peripheral flags
type PDPeripherals uint32

const (
	PeripheralNone          PDPeripherals = 0
	PeripheralAccelerometer PDPeripherals = 1 << 0
)

// PDLanguage represents system language
type PDLanguage int32

const (
	LanguageEnglish  PDLanguage = 0
	LanguageJapanese PDLanguage = 1
)

// PDRect represents a rectangle
type PDRect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// BitmapData contains bitmap information
type BitmapData struct {
	Width    int
	Height   int
	RowBytes int
	Mask     []byte
	Data     []byte
}

// SpriteCollisionResponseType represents collision response types
type SpriteCollisionResponseType int32

const (
	CollisionTypeSlide   SpriteCollisionResponseType = 0
	CollisionTypeFreeze  SpriteCollisionResponseType = 1
	CollisionTypeOverlap SpriteCollisionResponseType = 2
	CollisionTypeBounce  SpriteCollisionResponseType = 3
)

// CollisionPoint represents a point in collision detection
type CollisionPoint struct {
	X float32
	Y float32
}

// CollisionVector represents a vector in collision detection
type CollisionVector struct {
	X float32
	Y float32
}
