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
void* pd_gfx_rotatedBitmap(void* bitmap, float rotation, float xscale, float yscale);
void* pd_gfx_getBitmapMask(void* bitmap);
int pd_gfx_setBitmapMask(void* bitmap, void* mask);
void pd_gfx_setStencilImage(void* stencil, int tile);
void pd_gfx_setColorToPattern(void* color, void* bitmap, int x, int y);

// BitmapTable
void* pd_gfx_newBitmapTable(int count, int w, int h);
void pd_gfx_freeBitmapTable(void* table);
void* pd_gfx_loadBitmapTable(const char* path, const char** err);
void* pd_gfx_getTableBitmap(void* table, int idx);

// Tilemap
void* pd_gfx_tilemap_new(void);
void pd_gfx_tilemap_free(void* tilemap);
void pd_gfx_tilemap_setImageTable(void* tilemap, void* table);
void* pd_gfx_tilemap_getImageTable(void* tilemap);
void pd_gfx_tilemap_setSize(void* tilemap, int tilesWide, int tilesHigh);
void pd_gfx_tilemap_getSize(void* tilemap, int* tilesWide, int* tilesHigh);
void pd_gfx_tilemap_setTileAtPosition(void* tilemap, int x, int y, uint16_t idx);
int pd_gfx_tilemap_getTileAtPosition(void* tilemap, int x, int y);
void pd_gfx_tilemap_drawAtPoint(void* tilemap, float x, float y);

// Frame buffer
uint8_t* pd_gfx_getFrame(void);
uint8_t* pd_gfx_getDisplayFrame(void);
void pd_gfx_markUpdatedRows(int start, int end);
void pd_gfx_display(void);
void* pd_gfx_getDisplayBufferBitmap(void);
*/
import "C"
import (
	"runtime"
	"unsafe"
)

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
// Note: This function allocates memory. For hot paths, use DrawTextBytes instead.
func (g *Graphics) DrawText(text string, x, y int) int {
	cstr := make([]byte, len(text)+1)
	copy(cstr, text)
	return int(C.pd_gfx_drawText((*C.char)(unsafe.Pointer(&cstr[0])), C.int(len(text)), C.int(UTF8Encoding), C.int(x), C.int(y)))
}

// DrawTextBytes draws text from a pre-allocated byte buffer (null-terminated).
// Use this in hot paths to avoid memory allocation. The buffer must be
// null-terminated and len should not include the null terminator.
func (g *Graphics) DrawTextBytes(buf []byte, len int, x, y int) int {
	return int(C.pd_gfx_drawText((*C.char)(unsafe.Pointer(&buf[0])), C.int(len), C.int(UTF8Encoding), C.int(x), C.int(y)))
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
	font := &LCDFont{ptr: ptr}
	// Note: Playdate SDK manages font lifecycle, no explicit free needed
	// but we set a finalizer just in case for future compatibility
	runtime.SetFinalizer(font, func(f *LCDFont) {
		// Font is managed by SDK, no explicit free
	})
	return font, nil
}

// ============== Bitmap ==============

// NewBitmap creates a new bitmap
func (g *Graphics) NewBitmap(width, height int, bgcolor LCDColor) *LCDBitmap {
	ptr := C.pd_gfx_newBitmap(C.int(width), C.int(height), C.uint32_t(bgcolor))
	if ptr != nil {
		bitmap := &LCDBitmap{ptr: ptr}
		runtime.SetFinalizer(bitmap, func(b *LCDBitmap) {
			if b.ptr != nil {
				C.pd_gfx_freeBitmap(b.ptr)
			}
		})
		return bitmap
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
		bitmap := &LCDBitmap{ptr: ptr}
		runtime.SetFinalizer(bitmap, func(b *LCDBitmap) {
			if b.ptr != nil {
				C.pd_gfx_freeBitmap(b.ptr)
			}
		})
		return bitmap, nil
	}
	return nil, &loadError{path: path}
}

// CopyBitmap copies a bitmap
func (g *Graphics) CopyBitmap(bitmap *LCDBitmap) *LCDBitmap {
	if bitmap != nil && bitmap.ptr != nil {
		ptr := C.pd_gfx_copyBitmap(bitmap.ptr)
		if ptr != nil {
			copy := &LCDBitmap{ptr: ptr}
			runtime.SetFinalizer(copy, func(b *LCDBitmap) {
				if b.ptr != nil {
					C.pd_gfx_freeBitmap(b.ptr)
				}
			})
			return copy
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

// RotatedBitmap creates a rotated copy of a bitmap
func (g *Graphics) RotatedBitmap(bitmap *LCDBitmap, rotation, xscale, yscale float32) *LCDBitmap {
	if bitmap != nil && bitmap.ptr != nil {
		ptr := C.pd_gfx_rotatedBitmap(bitmap.ptr, C.float(rotation), C.float(xscale), C.float(yscale))
		if ptr != nil {
			rotated := &LCDBitmap{ptr: ptr}
			runtime.SetFinalizer(rotated, func(b *LCDBitmap) {
				if b.ptr != nil {
					C.pd_gfx_freeBitmap(b.ptr)
				}
			})
			return rotated
		}
	}
	return nil
}

// GetBitmapMask returns the mask bitmap for the given bitmap
func (g *Graphics) GetBitmapMask(bitmap *LCDBitmap) *LCDBitmap {
	if bitmap != nil && bitmap.ptr != nil {
		ptr := C.pd_gfx_getBitmapMask(bitmap.ptr)
		if ptr != nil {
			mask := &LCDBitmap{ptr: ptr}
			runtime.SetFinalizer(mask, func(b *LCDBitmap) {
				if b.ptr != nil {
					C.pd_gfx_freeBitmap(b.ptr)
				}
			})
			return mask
		}
	}
	return nil
}

// SetBitmapMask sets the mask bitmap for the given bitmap
func (g *Graphics) SetBitmapMask(bitmap, mask *LCDBitmap) bool {
	if bitmap != nil && bitmap.ptr != nil {
		var maskPtr unsafe.Pointer
		if mask != nil {
			maskPtr = mask.ptr
		}
		return C.pd_gfx_setBitmapMask(bitmap.ptr, maskPtr) != 0
	}
	return false
}

// SetStencilImage sets a stencil image for drawing
func (g *Graphics) SetStencilImage(stencil *LCDBitmap, tile bool) {
	if stencil != nil && stencil.ptr != nil {
		var tileFlag C.int
		if tile {
			tileFlag = 1
		}
		C.pd_gfx_setStencilImage(stencil.ptr, tileFlag)
	}
}

// SetColorToPattern sets up a pattern color for drawing
func (g *Graphics) SetColorToPattern(bitmap *LCDBitmap, x, y int) LCDColor {
	var color C.uint32_t
	var bitmapPtr unsafe.Pointer
	if bitmap != nil {
		bitmapPtr = bitmap.ptr
	}
	C.pd_gfx_setColorToPattern(unsafe.Pointer(&color), bitmapPtr, C.int(x), C.int(y))
	return LCDColor(color)
}

// ============== BitmapTable ==============

// NewBitmapTable creates a new bitmap table
func (g *Graphics) NewBitmapTable(count, width, height int) *LCDBitmapTable {
	ptr := C.pd_gfx_newBitmapTable(C.int(count), C.int(width), C.int(height))
	if ptr != nil {
		table := &LCDBitmapTable{ptr: ptr}
		runtime.SetFinalizer(table, func(t *LCDBitmapTable) {
			if t.ptr != nil {
				C.pd_gfx_freeBitmapTable(t.ptr)
			}
		})
		return table
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
		table := &LCDBitmapTable{ptr: ptr}
		runtime.SetFinalizer(table, func(t *LCDBitmapTable) {
			if t.ptr != nil {
				C.pd_gfx_freeBitmapTable(t.ptr)
			}
		})
		return table, nil
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

// ============== Tilemap ==============

// LCDTileMap represents a tilemap
type LCDTileMap struct {
	ptr unsafe.Pointer
}

// NewTilemap creates a new tilemap
func (g *Graphics) NewTilemap() *LCDTileMap {
	ptr := C.pd_gfx_tilemap_new()
	if ptr != nil {
		tilemap := &LCDTileMap{ptr: ptr}
		runtime.SetFinalizer(tilemap, func(t *LCDTileMap) {
			if t.ptr != nil {
				C.pd_gfx_tilemap_free(t.ptr)
			}
		})
		return tilemap
	}
	return nil
}

// Free frees the tilemap
func (t *LCDTileMap) Free() {
	if t != nil && t.ptr != nil {
		C.pd_gfx_tilemap_free(t.ptr)
		t.ptr = nil
	}
}

// SetImageTable sets the image table for the tilemap
func (t *LCDTileMap) SetImageTable(table *LCDBitmapTable) {
	if t != nil && t.ptr != nil && table != nil && table.ptr != nil {
		C.pd_gfx_tilemap_setImageTable(t.ptr, table.ptr)
	}
}

// GetImageTable gets the image table from the tilemap
func (t *LCDTileMap) GetImageTable() *LCDBitmapTable {
	if t != nil && t.ptr != nil {
		ptr := C.pd_gfx_tilemap_getImageTable(t.ptr)
		if ptr != nil {
			return &LCDBitmapTable{ptr: ptr}
		}
	}
	return nil
}

// SetSize sets the size of the tilemap in tiles
func (t *LCDTileMap) SetSize(tilesWide, tilesHigh int) {
	if t != nil && t.ptr != nil {
		C.pd_gfx_tilemap_setSize(t.ptr, C.int(tilesWide), C.int(tilesHigh))
	}
}

// GetSize returns the size of the tilemap in tiles
func (t *LCDTileMap) GetSize() (tilesWide, tilesHigh int) {
	if t != nil && t.ptr != nil {
		var wide, high C.int
		C.pd_gfx_tilemap_getSize(t.ptr, &wide, &high)
		return int(wide), int(high)
	}
	return 0, 0
}

// SetTileAtPosition sets the tile index at the given position
func (t *LCDTileMap) SetTileAtPosition(x, y int, idx uint16) {
	if t != nil && t.ptr != nil {
		C.pd_gfx_tilemap_setTileAtPosition(t.ptr, C.int(x), C.int(y), C.uint16_t(idx))
	}
}

// GetTileAtPosition gets the tile index at the given position
func (t *LCDTileMap) GetTileAtPosition(x, y int) int {
	if t != nil && t.ptr != nil {
		return int(C.pd_gfx_tilemap_getTileAtPosition(t.ptr, C.int(x), C.int(y)))
	}
	return 0
}

// DrawAtPoint draws the tilemap at the given position
func (t *LCDTileMap) DrawAtPoint(x, y float32) {
	if t != nil && t.ptr != nil {
		C.pd_gfx_tilemap_drawAtPoint(t.ptr, C.float(x), C.float(y))
	}
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
