//go:build tinygo

// TinyGo implementation of Graphics API

package pdgo

import (
	"unsafe"
)

// LCDBitmap represents a bitmap image
type LCDBitmap struct {
	ptr uintptr
}

// LCDBitmapTable represents a table of bitmaps
type LCDBitmapTable struct {
	ptr uintptr
}

// LCDFont represents a font
type LCDFont struct {
	ptr uintptr
}

// Graphics provides access to Playdate graphics functions
type Graphics struct{}

func newGraphics() *Graphics {
	return &Graphics{}
}

// ============== Drawing Context ==============

// Clear clears the display with the given color
func (g *Graphics) Clear(color LCDColor) {
	if bridgeGfxClear != nil {
		bridgeGfxClear(uint32(color))
	}
}

// SetBackgroundColor sets the background color
func (g *Graphics) SetBackgroundColor(color LCDSolidColor) {
	if bridgeGfxSetBackgroundColor != nil {
		bridgeGfxSetBackgroundColor(int32(color))
	}
}

// SetDrawMode sets the drawing mode
func (g *Graphics) SetDrawMode(mode LCDBitmapDrawMode) LCDBitmapDrawMode {
	if bridgeGfxSetDrawMode != nil {
		return LCDBitmapDrawMode(bridgeGfxSetDrawMode(int32(mode)))
	}
	return mode
}

// SetDrawOffset sets the drawing offset
func (g *Graphics) SetDrawOffset(dx, dy int) {
	if bridgeGfxSetDrawOffset != nil {
		bridgeGfxSetDrawOffset(int32(dx), int32(dy))
	}
}

// SetClipRect sets the clipping rectangle
func (g *Graphics) SetClipRect(x, y, width, height int) {
	if bridgeGfxSetClipRect != nil {
		bridgeGfxSetClipRect(int32(x), int32(y), int32(width), int32(height))
	}
}

// ClearClipRect clears the clipping rectangle
func (g *Graphics) ClearClipRect() {
	if bridgeGfxClearClipRect != nil {
		bridgeGfxClearClipRect()
	}
}

// SetLineCapStyle sets line cap style
func (g *Graphics) SetLineCapStyle(style LCDLineCapStyle) {
	if bridgeGfxSetLineCapStyle != nil {
		bridgeGfxSetLineCapStyle(int32(style))
	}
}

// PushContext pushes a new drawing context
func (g *Graphics) PushContext(target *LCDBitmap) {
	if bridgeGfxPushContext != nil {
		var ptr uintptr
		if target != nil {
			ptr = target.ptr
		}
		bridgeGfxPushContext(ptr)
	}
}

// PopContext pops the current drawing context
func (g *Graphics) PopContext() {
	if bridgeGfxPopContext != nil {
		bridgeGfxPopContext()
	}
}

// ============== Drawing Primitives ==============

// FillRect fills a rectangle
func (g *Graphics) FillRect(x, y, width, height int, color LCDColor) {
	if bridgeGfxFillRect != nil {
		bridgeGfxFillRect(int32(x), int32(y), int32(width), int32(height), uint32(color))
	}
}

// DrawRect draws a rectangle outline
func (g *Graphics) DrawRect(x, y, width, height int, color LCDColor) {
	if bridgeGfxDrawRect != nil {
		bridgeGfxDrawRect(int32(x), int32(y), int32(width), int32(height), uint32(color))
	}
}

// DrawLine draws a line
func (g *Graphics) DrawLine(x1, y1, x2, y2, width int, color LCDColor) {
	if bridgeGfxDrawLine != nil {
		bridgeGfxDrawLine(int32(x1), int32(y1), int32(x2), int32(y2), int32(width), uint32(color))
	}
}

// FillTriangle fills a triangle
func (g *Graphics) FillTriangle(x1, y1, x2, y2, x3, y3 int, color LCDColor) {
	if bridgeGfxFillTriangle != nil {
		bridgeGfxFillTriangle(int32(x1), int32(y1), int32(x2), int32(y2), int32(x3), int32(y3), uint32(color))
	}
}

// DrawEllipse draws an ellipse outline
func (g *Graphics) DrawEllipse(x, y, width, height, lineWidth int, startAngle, endAngle float32, color LCDColor) {
	if bridgeGfxDrawEllipse != nil {
		bridgeGfxDrawEllipse(int32(x), int32(y), int32(width), int32(height), int32(lineWidth), startAngle, endAngle, uint32(color))
	}
}

// FillEllipse fills an ellipse
func (g *Graphics) FillEllipse(x, y, width, height int, startAngle, endAngle float32, color LCDColor) {
	if bridgeGfxFillEllipse != nil {
		bridgeGfxFillEllipse(int32(x), int32(y), int32(width), int32(height), startAngle, endAngle, uint32(color))
	}
}

// ============== Text ==============

// DrawText draws text at the given position
func (g *Graphics) DrawText(text string, x, y int) int {
	if bridgeGfxDrawText != nil {
		buf := make([]byte, len(text)+1)
		copy(buf, text)
		return int(bridgeGfxDrawText(&buf[0], int32(len(text)), int32(UTF8Encoding), int32(x), int32(y)))
	}
	return 0
}

// GetTextWidth returns the width of text
func (g *Graphics) GetTextWidth(font *LCDFont, text string, tracking int) int {
	if bridgeGfxGetTextWidth != nil {
		buf := make([]byte, len(text)+1)
		copy(buf, text)
		var fontPtr uintptr
		if font != nil {
			fontPtr = font.ptr
		}
		return int(bridgeGfxGetTextWidth(fontPtr, &buf[0], int32(len(text)), int32(UTF8Encoding), int32(tracking)))
	}
	return 0
}

// SetFont sets the current font
func (g *Graphics) SetFont(font *LCDFont) {
	if bridgeGfxSetFont != nil && font != nil {
		bridgeGfxSetFont(font.ptr)
	}
}

// SetTextTracking sets text tracking
func (g *Graphics) SetTextTracking(tracking int) {
	if bridgeGfxSetTextTracking != nil {
		bridgeGfxSetTextTracking(int32(tracking))
	}
}

// LoadFont loads a font from file
func (g *Graphics) LoadFont(path string) (*LCDFont, error) {
	if bridgeGfxLoadFont != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		ptr := bridgeGfxLoadFont(&buf[0])
		if ptr == 0 {
			return nil, &loadError{path: path}
		}
		return &LCDFont{ptr: ptr}, nil
	}
	return nil, &loadError{path: path}
}

// ============== Bitmap ==============

// NewBitmap creates a new bitmap
func (g *Graphics) NewBitmap(width, height int, bgcolor LCDColor) *LCDBitmap {
	if bridgeGfxNewBitmap != nil {
		ptr := bridgeGfxNewBitmap(int32(width), int32(height), uint32(bgcolor))
		if ptr != 0 {
			return &LCDBitmap{ptr: ptr}
		}
	}
	return nil
}

// FreeBitmap frees a bitmap
func (g *Graphics) FreeBitmap(bitmap *LCDBitmap) {
	if bridgeGfxFreeBitmap != nil && bitmap != nil && bitmap.ptr != 0 {
		bridgeGfxFreeBitmap(bitmap.ptr)
		bitmap.ptr = 0
	}
}

// LoadBitmap loads a bitmap from file
func (g *Graphics) LoadBitmap(path string) (*LCDBitmap, error) {
	if bridgeGfxLoadBitmap != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		ptr := bridgeGfxLoadBitmap(&buf[0])
		if ptr != 0 {
			return &LCDBitmap{ptr: ptr}, nil
		}
	}
	return nil, &loadError{path: path}
}

// CopyBitmap copies a bitmap
func (g *Graphics) CopyBitmap(bitmap *LCDBitmap) *LCDBitmap {
	if bridgeGfxCopyBitmap != nil && bitmap != nil && bitmap.ptr != 0 {
		ptr := bridgeGfxCopyBitmap(bitmap.ptr)
		if ptr != 0 {
			return &LCDBitmap{ptr: ptr}
		}
	}
	return nil
}

// DrawBitmap draws a bitmap
func (g *Graphics) DrawBitmap(bitmap *LCDBitmap, x, y int, flip LCDBitmapFlip) {
	if bridgeGfxDrawBitmap != nil && bitmap != nil && bitmap.ptr != 0 {
		bridgeGfxDrawBitmap(bitmap.ptr, int32(x), int32(y), int32(flip))
	}
}

// TileBitmap tiles a bitmap
func (g *Graphics) TileBitmap(bitmap *LCDBitmap, x, y, width, height int, flip LCDBitmapFlip) {
	if bridgeGfxTileBitmap != nil && bitmap != nil && bitmap.ptr != 0 {
		bridgeGfxTileBitmap(bitmap.ptr, int32(x), int32(y), int32(width), int32(height), int32(flip))
	}
}

// DrawScaledBitmap draws a scaled bitmap
func (g *Graphics) DrawScaledBitmap(bitmap *LCDBitmap, x, y int, xscale, yscale float32) {
	if bridgeGfxDrawScaledBitmap != nil && bitmap != nil && bitmap.ptr != 0 {
		bridgeGfxDrawScaledBitmap(bitmap.ptr, int32(x), int32(y), xscale, yscale)
	}
}

// DrawRotatedBitmap draws a rotated bitmap
func (g *Graphics) DrawRotatedBitmap(bitmap *LCDBitmap, x, y int, rotation, centerX, centerY, xscale, yscale float32) {
	if bridgeGfxDrawRotatedBitmap != nil && bitmap != nil && bitmap.ptr != 0 {
		bridgeGfxDrawRotatedBitmap(bitmap.ptr, int32(x), int32(y), rotation, centerX, centerY, xscale, yscale)
	}
}

// GetBitmapData returns bitmap data information
func (g *Graphics) GetBitmapData(bitmap *LCDBitmap) *BitmapData {
	if bridgeGfxGetBitmapData != nil && bitmap != nil && bitmap.ptr != 0 {
		var width, height, rowbytes int32
		var mask, data uintptr
		bridgeGfxGetBitmapData(bitmap.ptr, &width, &height, &rowbytes, &mask, &data)

		result := &BitmapData{
			Width:    int(width),
			Height:   int(height),
			RowBytes: int(rowbytes),
		}

		if data != 0 {
			result.Data = unsafe.Slice((*byte)(unsafe.Pointer(data)), int(height)*int(rowbytes))
		}
		if mask != 0 {
			result.Mask = unsafe.Slice((*byte)(unsafe.Pointer(mask)), int(height)*int(rowbytes))
		}
		return result
	}
	return nil
}

// ClearBitmap clears a bitmap
func (g *Graphics) ClearBitmap(bitmap *LCDBitmap, bgcolor LCDColor) {
	if bridgeGfxClearBitmap != nil && bitmap != nil && bitmap.ptr != 0 {
		bridgeGfxClearBitmap(bitmap.ptr, uint32(bgcolor))
	}
}

// ============== BitmapTable ==============

// NewBitmapTable creates a new bitmap table
func (g *Graphics) NewBitmapTable(count, width, height int) *LCDBitmapTable {
	if bridgeGfxNewBitmapTable != nil {
		ptr := bridgeGfxNewBitmapTable(int32(count), int32(width), int32(height))
		if ptr != 0 {
			return &LCDBitmapTable{ptr: ptr}
		}
	}
	return nil
}

// FreeBitmapTable frees a bitmap table
func (g *Graphics) FreeBitmapTable(table *LCDBitmapTable) {
	if bridgeGfxFreeBitmapTable != nil && table != nil && table.ptr != 0 {
		bridgeGfxFreeBitmapTable(table.ptr)
		table.ptr = 0
	}
}

// LoadBitmapTable loads a bitmap table from file
func (g *Graphics) LoadBitmapTable(path string) (*LCDBitmapTable, error) {
	if bridgeGfxLoadBitmapTable != nil {
		buf := make([]byte, len(path)+1)
		copy(buf, path)
		ptr := bridgeGfxLoadBitmapTable(&buf[0])
		if ptr != 0 {
			return &LCDBitmapTable{ptr: ptr}, nil
		}
	}
	return nil, &loadError{path: path}
}

// GetTableBitmap gets a bitmap from a table
func (g *Graphics) GetTableBitmap(table *LCDBitmapTable, idx int) *LCDBitmap {
	if bridgeGfxGetTableBitmap != nil && table != nil && table.ptr != 0 {
		ptr := bridgeGfxGetTableBitmap(table.ptr, int32(idx))
		if ptr != 0 {
			return &LCDBitmap{ptr: ptr}
		}
	}
	return nil
}

// ============== Frame Buffer ==============

// GetFrame returns the display framebuffer
func (g *Graphics) GetFrame() []byte {
	if bridgeGfxGetFrame != nil {
		ptr := bridgeGfxGetFrame()
		if ptr != 0 {
			return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), LCDRows*LCDRowSize)
		}
	}
	return nil
}

// GetDisplayFrame returns the actual display buffer
func (g *Graphics) GetDisplayFrame() []byte {
	if bridgeGfxGetDisplayFrame != nil {
		ptr := bridgeGfxGetDisplayFrame()
		if ptr != 0 {
			return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), LCDRows*LCDRowSize)
		}
	}
	return nil
}

// MarkUpdatedRows marks rows as needing update
func (g *Graphics) MarkUpdatedRows(start, end int) {
	if bridgeGfxMarkUpdatedRows != nil {
		bridgeGfxMarkUpdatedRows(int32(start), int32(end))
	}
}

// Display flushes the frame buffer to screen
func (g *Graphics) Display() {
	if bridgeGfxDisplay != nil {
		bridgeGfxDisplay()
	}
}

// loadError represents a load error
type loadError struct {
	path string
}

func (e *loadError) Error() string {
	return "failed to load: " + e.path
}
