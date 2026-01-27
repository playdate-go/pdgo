// pdgo Graphics API - unified CGO implementation

package pdgo

/*
#include <stdint.h>

// Graphics API
void pd_gfx_clear(uint32_t color);
void pd_gfx_setBackgroundColor(int color);
int pd_gfx_setDrawMode(int mode);
void pd_gfx_setDrawOffset(int dx, int dy);
void pd_gfx_setClipRect(int x, int y, int w, int h);
void pd_gfx_clearClipRect(void);
void pd_gfx_setLineCapStyle(int style);
void pd_gfx_setFont(void* font);
void pd_gfx_setTextTracking(int tracking);
void pd_gfx_pushContext(void* target);
void pd_gfx_popContext(void);

// Drawing primitives
void pd_gfx_fillRect(int x, int y, int w, int h, uint32_t color);
void pd_gfx_drawRect(int x, int y, int w, int h, uint32_t color);
void pd_gfx_drawLine(int x1, int y1, int x2, int y2, int width, uint32_t color);
void pd_gfx_fillTriangle(int x1, int y1, int x2, int y2, int x3, int y3, uint32_t color);
void pd_gfx_drawEllipse(int x, int y, int w, int h, int lineWidth, float startAngle, float endAngle, uint32_t color);
void pd_gfx_fillEllipse(int x, int y, int w, int h, float startAngle, float endAngle, uint32_t color);

// Text
int pd_gfx_drawText(const char* text, int len, int encoding, int x, int y);
int pd_gfx_getTextWidth(void* font, const char* text, int len, int encoding, int tracking);
void* pd_gfx_loadFont(const char* path, const char** err);

// Bitmap
void* pd_gfx_newBitmap(int w, int h, uint32_t bgcolor);
void pd_gfx_freeBitmap(void* bitmap);
void* pd_gfx_loadBitmap(const char* path, const char** err);
void* pd_gfx_copyBitmap(void* bitmap);
void pd_gfx_drawBitmap(void* bitmap, int x, int y, int flip);
void pd_gfx_tileBitmap(void* bitmap, int x, int y, int w, int h, int flip);
void pd_gfx_drawScaledBitmap(void* bitmap, int x, int y, float xscale, float yscale);
void pd_gfx_drawRotatedBitmap(void* bitmap, int x, int y, float rotation, float cx, float cy, float xscale, float yscale);
void pd_gfx_getBitmapData(void* bitmap, int* w, int* h, int* rowbytes, uint8_t** mask, uint8_t** data);
void pd_gfx_clearBitmap(void* bitmap, uint32_t bgcolor);

// BitmapTable
void* pd_gfx_newBitmapTable(int count, int w, int h);
void pd_gfx_freeBitmapTable(void* table);
void* pd_gfx_loadBitmapTable(const char* path, const char** err);
void* pd_gfx_getTableBitmap(void* table, int idx);

// Frame buffer
uint8_t* pd_gfx_getFrame(void);
uint8_t* pd_gfx_getDisplayFrame(void);
void pd_gfx_markUpdatedRows(int start, int end);
void pd_gfx_display(void);
void* pd_gfx_getDisplayBufferBitmap(void);
*/
import "C"
import "unsafe"

// LCDBitmap represents a bitmap image
type LCDBitmap struct {
	ptr unsafe.Pointer
}

// LCDBitmapTable represents a table of bitmaps
type LCDBitmapTable struct {
	ptr unsafe.Pointer
}

// LCDFont represents a font
type LCDFont struct {
	ptr unsafe.Pointer
}

// Graphics provides access to Playdate graphics functions
type Graphics struct{}

func newGraphics() *Graphics {
	return &Graphics{}
}

// ============== Drawing Context ==============

// Clear clears the display with the given color
func (g *Graphics) Clear(color LCDColor) {
	C.pd_gfx_clear(C.uint32_t(color))
}

// SetBackgroundColor sets the background color
func (g *Graphics) SetBackgroundColor(color LCDSolidColor) {
	C.pd_gfx_setBackgroundColor(C.int(color))
}

// SetDrawMode sets the drawing mode
func (g *Graphics) SetDrawMode(mode LCDBitmapDrawMode) LCDBitmapDrawMode {
	return LCDBitmapDrawMode(C.pd_gfx_setDrawMode(C.int(mode)))
}

// SetDrawOffset sets the drawing offset
func (g *Graphics) SetDrawOffset(dx, dy int) {
	C.pd_gfx_setDrawOffset(C.int(dx), C.int(dy))
}

// SetClipRect sets the clipping rectangle
func (g *Graphics) SetClipRect(x, y, width, height int) {
	C.pd_gfx_setClipRect(C.int(x), C.int(y), C.int(width), C.int(height))
}

// ClearClipRect clears the clipping rectangle
func (g *Graphics) ClearClipRect() {
	C.pd_gfx_clearClipRect()
}

// SetLineCapStyle sets line cap style
func (g *Graphics) SetLineCapStyle(style LCDLineCapStyle) {
	C.pd_gfx_setLineCapStyle(C.int(style))
}

// PushContext pushes a new drawing context
func (g *Graphics) PushContext(target *LCDBitmap) {
	var ptr unsafe.Pointer
	if target != nil {
		ptr = target.ptr
	}
	C.pd_gfx_pushContext(ptr)
}

// PopContext pops the current drawing context
func (g *Graphics) PopContext() {
	C.pd_gfx_popContext()
}

// ============== Drawing Primitives ==============

// FillRect fills a rectangle
func (g *Graphics) FillRect(x, y, width, height int, color LCDColor) {
	C.pd_gfx_fillRect(C.int(x), C.int(y), C.int(width), C.int(height), C.uint32_t(color))
}

// DrawRect draws a rectangle outline
func (g *Graphics) DrawRect(x, y, width, height int, color LCDColor) {
	C.pd_gfx_drawRect(C.int(x), C.int(y), C.int(width), C.int(height), C.uint32_t(color))
}

// DrawLine draws a line
func (g *Graphics) DrawLine(x1, y1, x2, y2, width int, color LCDColor) {
	C.pd_gfx_drawLine(C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(width), C.uint32_t(color))
}

// FillTriangle fills a triangle
func (g *Graphics) FillTriangle(x1, y1, x2, y2, x3, y3 int, color LCDColor) {
	C.pd_gfx_fillTriangle(C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(x3), C.int(y3), C.uint32_t(color))
}

// DrawEllipse draws an ellipse outline
func (g *Graphics) DrawEllipse(x, y, width, height, lineWidth int, startAngle, endAngle float32, color LCDColor) {
	C.pd_gfx_drawEllipse(C.int(x), C.int(y), C.int(width), C.int(height), C.int(lineWidth), C.float(startAngle), C.float(endAngle), C.uint32_t(color))
}

// FillEllipse fills an ellipse
func (g *Graphics) FillEllipse(x, y, width, height int, startAngle, endAngle float32, color LCDColor) {
	C.pd_gfx_fillEllipse(C.int(x), C.int(y), C.int(width), C.int(height), C.float(startAngle), C.float(endAngle), C.uint32_t(color))
}

// ============== Text ==============

// DrawText draws text at the given position
func (g *Graphics) DrawText(text string, x, y int) int {
	cstr := make([]byte, len(text)+1)
	copy(cstr, text)
	return int(C.pd_gfx_drawText((*C.char)(unsafe.Pointer(&cstr[0])), C.int(len(text)), C.int(UTF8Encoding), C.int(x), C.int(y)))
}

// GetTextWidth returns the width of text
func (g *Graphics) GetTextWidth(font *LCDFont, text string, tracking int) int {
	cstr := make([]byte, len(text)+1)
	copy(cstr, text)
	var fontPtr unsafe.Pointer
	if font != nil {
		fontPtr = font.ptr
	}
	return int(C.pd_gfx_getTextWidth(fontPtr, (*C.char)(unsafe.Pointer(&cstr[0])), C.int(len(text)), C.int(UTF8Encoding), C.int(tracking)))
}

// SetFont sets the current font
func (g *Graphics) SetFont(font *LCDFont) {
	if font != nil {
		C.pd_gfx_setFont(font.ptr)
	}
}

// SetTextTracking sets text tracking
func (g *Graphics) SetTextTracking(tracking int) {
	C.pd_gfx_setTextTracking(C.int(tracking))
}

// LoadFont loads a font from file
func (g *Graphics) LoadFont(path string) (*LCDFont, error) {
	cpath := make([]byte, len(path)+1)
	copy(cpath, path)
	ptr := C.pd_gfx_loadFont((*C.char)(unsafe.Pointer(&cpath[0])), nil)
	if ptr == nil {
		return nil, &loadError{path: path}
	}
	return &LCDFont{ptr: ptr}, nil
}

// ============== Bitmap ==============

// NewBitmap creates a new bitmap
func (g *Graphics) NewBitmap(width, height int, bgcolor LCDColor) *LCDBitmap {
	ptr := C.pd_gfx_newBitmap(C.int(width), C.int(height), C.uint32_t(bgcolor))
	if ptr != nil {
		return &LCDBitmap{ptr: ptr}
	}
	return nil
}

// FreeBitmap frees a bitmap
func (g *Graphics) FreeBitmap(bitmap *LCDBitmap) {
	if bitmap != nil && bitmap.ptr != nil {
		C.pd_gfx_freeBitmap(bitmap.ptr)
		bitmap.ptr = nil
	}
}

// LoadBitmap loads a bitmap from file
func (g *Graphics) LoadBitmap(path string) (*LCDBitmap, error) {
	cpath := make([]byte, len(path)+1)
	copy(cpath, path)
	ptr := C.pd_gfx_loadBitmap((*C.char)(unsafe.Pointer(&cpath[0])), nil)
	if ptr != nil {
		return &LCDBitmap{ptr: ptr}, nil
	}
	return nil, &loadError{path: path}
}

// CopyBitmap copies a bitmap
func (g *Graphics) CopyBitmap(bitmap *LCDBitmap) *LCDBitmap {
	if bitmap != nil && bitmap.ptr != nil {
		ptr := C.pd_gfx_copyBitmap(bitmap.ptr)
		if ptr != nil {
			return &LCDBitmap{ptr: ptr}
		}
	}
	return nil
}

// DrawBitmap draws a bitmap
func (g *Graphics) DrawBitmap(bitmap *LCDBitmap, x, y int, flip LCDBitmapFlip) {
	if bitmap != nil && bitmap.ptr != nil {
		C.pd_gfx_drawBitmap(bitmap.ptr, C.int(x), C.int(y), C.int(flip))
	}
}

// TileBitmap tiles a bitmap
func (g *Graphics) TileBitmap(bitmap *LCDBitmap, x, y, width, height int, flip LCDBitmapFlip) {
	if bitmap != nil && bitmap.ptr != nil {
		C.pd_gfx_tileBitmap(bitmap.ptr, C.int(x), C.int(y), C.int(width), C.int(height), C.int(flip))
	}
}

// DrawScaledBitmap draws a scaled bitmap
func (g *Graphics) DrawScaledBitmap(bitmap *LCDBitmap, x, y int, xscale, yscale float32) {
	if bitmap != nil && bitmap.ptr != nil {
		C.pd_gfx_drawScaledBitmap(bitmap.ptr, C.int(x), C.int(y), C.float(xscale), C.float(yscale))
	}
}

// DrawRotatedBitmap draws a rotated bitmap
func (g *Graphics) DrawRotatedBitmap(bitmap *LCDBitmap, x, y int, rotation, centerX, centerY, xscale, yscale float32) {
	if bitmap != nil && bitmap.ptr != nil {
		C.pd_gfx_drawRotatedBitmap(bitmap.ptr, C.int(x), C.int(y), C.float(rotation), C.float(centerX), C.float(centerY), C.float(xscale), C.float(yscale))
	}
}

// GetBitmapData returns bitmap data information
func (g *Graphics) GetBitmapData(bitmap *LCDBitmap) *BitmapData {
	if bitmap != nil && bitmap.ptr != nil {
		var width, height, rowbytes C.int
		var mask, data *C.uint8_t
		C.pd_gfx_getBitmapData(bitmap.ptr, &width, &height, &rowbytes, &mask, &data)

		result := &BitmapData{
			Width:    int(width),
			Height:   int(height),
			RowBytes: int(rowbytes),
		}

		if data != nil {
			result.Data = unsafe.Slice((*byte)(unsafe.Pointer(data)), int(height)*int(rowbytes))
		}
		if mask != nil {
			result.Mask = unsafe.Slice((*byte)(unsafe.Pointer(mask)), int(height)*int(rowbytes))
		}
		return result
	}
	return nil
}

// ClearBitmap clears a bitmap
func (g *Graphics) ClearBitmap(bitmap *LCDBitmap, bgcolor LCDColor) {
	if bitmap != nil && bitmap.ptr != nil {
		C.pd_gfx_clearBitmap(bitmap.ptr, C.uint32_t(bgcolor))
	}
}

// ============== BitmapTable ==============

// NewBitmapTable creates a new bitmap table
func (g *Graphics) NewBitmapTable(count, width, height int) *LCDBitmapTable {
	ptr := C.pd_gfx_newBitmapTable(C.int(count), C.int(width), C.int(height))
	if ptr != nil {
		return &LCDBitmapTable{ptr: ptr}
	}
	return nil
}

// FreeBitmapTable frees a bitmap table
func (g *Graphics) FreeBitmapTable(table *LCDBitmapTable) {
	if table != nil && table.ptr != nil {
		C.pd_gfx_freeBitmapTable(table.ptr)
		table.ptr = nil
	}
}

// LoadBitmapTable loads a bitmap table from file
func (g *Graphics) LoadBitmapTable(path string) (*LCDBitmapTable, error) {
	cpath := make([]byte, len(path)+1)
	copy(cpath, path)
	ptr := C.pd_gfx_loadBitmapTable((*C.char)(unsafe.Pointer(&cpath[0])), nil)
	if ptr != nil {
		return &LCDBitmapTable{ptr: ptr}, nil
	}
	return nil, &loadError{path: path}
}

// GetTableBitmap gets a bitmap from a table
func (g *Graphics) GetTableBitmap(table *LCDBitmapTable, idx int) *LCDBitmap {
	if table != nil && table.ptr != nil {
		ptr := C.pd_gfx_getTableBitmap(table.ptr, C.int(idx))
		if ptr != nil {
			return &LCDBitmap{ptr: ptr}
		}
	}
	return nil
}

// ============== Frame Buffer ==============

// GetFrame returns the display framebuffer
func (g *Graphics) GetFrame() []byte {
	ptr := C.pd_gfx_getFrame()
	if ptr != nil {
		return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), LCDRows*LCDRowSize)
	}
	return nil
}

// GetDisplayFrame returns the actual display buffer
func (g *Graphics) GetDisplayFrame() []byte {
	ptr := C.pd_gfx_getDisplayFrame()
	if ptr != nil {
		return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), LCDRows*LCDRowSize)
	}
	return nil
}

// MarkUpdatedRows marks rows as needing update
func (g *Graphics) MarkUpdatedRows(start, end int) {
	C.pd_gfx_markUpdatedRows(C.int(start), C.int(end))
}

// Display flushes the frame buffer to screen
func (g *Graphics) Display() {
	C.pd_gfx_display()
}

// GetDisplayBufferBitmap returns the display buffer as a bitmap
func (g *Graphics) GetDisplayBufferBitmap() *LCDBitmap {
	ptr := C.pd_gfx_getDisplayBufferBitmap()
	if ptr != nil {
		return &LCDBitmap{ptr: ptr}
	}
	return nil
}

// loadError represents a load error
type loadError struct {
	path string
}

func (e *loadError) Error() string {
	return "failed to load: " + e.path
}
