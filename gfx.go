//go:build !tinygo

package pdgo

/*
#cgo CFLAGS: -DTARGET_EXTENSION=1
#include "pd_api.h"
#include <stdlib.h>

// Graphics API helper functions

static void gfx_clear(const struct playdate_graphics* gfx, LCDColor color) {
    gfx->clear(color);
}

static void gfx_setBackgroundColor(const struct playdate_graphics* gfx, LCDSolidColor color) {
    gfx->setBackgroundColor(color);
}

static void gfx_setStencil(const struct playdate_graphics* gfx, LCDBitmap* stencil) {
    gfx->setStencil(stencil);
}

static LCDBitmapDrawMode gfx_setDrawMode(const struct playdate_graphics* gfx, LCDBitmapDrawMode mode) {
    return gfx->setDrawMode(mode);
}

static void gfx_setDrawOffset(const struct playdate_graphics* gfx, int dx, int dy) {
    gfx->setDrawOffset(dx, dy);
}

static void gfx_setClipRect(const struct playdate_graphics* gfx, int x, int y, int width, int height) {
    gfx->setClipRect(x, y, width, height);
}

static void gfx_clearClipRect(const struct playdate_graphics* gfx) {
    gfx->clearClipRect();
}

static void gfx_setLineCapStyle(const struct playdate_graphics* gfx, LCDLineCapStyle endCapStyle) {
    gfx->setLineCapStyle(endCapStyle);
}

static void gfx_setFont(const struct playdate_graphics* gfx, LCDFont* font) {
    gfx->setFont(font);
}

static void gfx_setTextTracking(const struct playdate_graphics* gfx, int tracking) {
    gfx->setTextTracking(tracking);
}

static void gfx_pushContext(const struct playdate_graphics* gfx, LCDBitmap* target) {
    gfx->pushContext(target);
}

static void gfx_popContext(const struct playdate_graphics* gfx) {
    gfx->popContext();
}

static void gfx_drawBitmap(const struct playdate_graphics* gfx, LCDBitmap* bitmap, int x, int y, LCDBitmapFlip flip) {
    gfx->drawBitmap(bitmap, x, y, flip);
}

static void gfx_tileBitmap(const struct playdate_graphics* gfx, LCDBitmap* bitmap, int x, int y, int width, int height, LCDBitmapFlip flip) {
    gfx->tileBitmap(bitmap, x, y, width, height, flip);
}

static void gfx_drawLine(const struct playdate_graphics* gfx, int x1, int y1, int x2, int y2, int width, LCDColor color) {
    gfx->drawLine(x1, y1, x2, y2, width, color);
}

static void gfx_fillTriangle(const struct playdate_graphics* gfx, int x1, int y1, int x2, int y2, int x3, int y3, LCDColor color) {
    gfx->fillTriangle(x1, y1, x2, y2, x3, y3, color);
}

static void gfx_drawRect(const struct playdate_graphics* gfx, int x, int y, int width, int height, LCDColor color) {
    gfx->drawRect(x, y, width, height, color);
}

static void gfx_fillRect(const struct playdate_graphics* gfx, int x, int y, int width, int height, LCDColor color) {
    gfx->fillRect(x, y, width, height, color);
}

static void gfx_drawEllipse(const struct playdate_graphics* gfx, int x, int y, int width, int height, int lineWidth, float startAngle, float endAngle, LCDColor color) {
    gfx->drawEllipse(x, y, width, height, lineWidth, startAngle, endAngle, color);
}

static void gfx_fillEllipse(const struct playdate_graphics* gfx, int x, int y, int width, int height, float startAngle, float endAngle, LCDColor color) {
    gfx->fillEllipse(x, y, width, height, startAngle, endAngle, color);
}

static void gfx_drawScaledBitmap(const struct playdate_graphics* gfx, LCDBitmap* bitmap, int x, int y, float xscale, float yscale) {
    gfx->drawScaledBitmap(bitmap, x, y, xscale, yscale);
}

static int gfx_drawText(const struct playdate_graphics* gfx, const void* text, size_t len, PDStringEncoding encoding, int x, int y) {
    return gfx->drawText(text, len, encoding, x, y);
}

// LCDBitmap functions
static LCDBitmap* gfx_newBitmap(const struct playdate_graphics* gfx, int width, int height, LCDColor bgcolor) {
    return gfx->newBitmap(width, height, bgcolor);
}

static void gfx_freeBitmap(const struct playdate_graphics* gfx, LCDBitmap* bitmap) {
    gfx->freeBitmap(bitmap);
}

static LCDBitmap* gfx_loadBitmap(const struct playdate_graphics* gfx, const char* path, const char** outerr) {
    return gfx->loadBitmap(path, outerr);
}

static LCDBitmap* gfx_copyBitmap(const struct playdate_graphics* gfx, LCDBitmap* bitmap) {
    return gfx->copyBitmap(bitmap);
}

static void gfx_loadIntoBitmap(const struct playdate_graphics* gfx, const char* path, LCDBitmap* bitmap, const char** outerr) {
    gfx->loadIntoBitmap(path, bitmap, outerr);
}

static void gfx_getBitmapData(const struct playdate_graphics* gfx, LCDBitmap* bitmap, int* width, int* height, int* rowbytes, uint8_t** mask, uint8_t** data) {
    gfx->getBitmapData(bitmap, width, height, rowbytes, mask, data);
}

static void gfx_clearBitmap(const struct playdate_graphics* gfx, LCDBitmap* bitmap, LCDColor bgcolor) {
    gfx->clearBitmap(bitmap, bgcolor);
}

static LCDBitmap* gfx_rotatedBitmap(const struct playdate_graphics* gfx, LCDBitmap* bitmap, float rotation, float xscale, float yscale, int* allocedSize) {
    return gfx->rotatedBitmap(bitmap, rotation, xscale, yscale, allocedSize);
}

// LCDBitmapTable functions
static LCDBitmapTable* gfx_newBitmapTable(const struct playdate_graphics* gfx, int count, int width, int height) {
    return gfx->newBitmapTable(count, width, height);
}

static void gfx_freeBitmapTable(const struct playdate_graphics* gfx, LCDBitmapTable* table) {
    gfx->freeBitmapTable(table);
}

static LCDBitmapTable* gfx_loadBitmapTable(const struct playdate_graphics* gfx, const char* path, const char** outerr) {
    return gfx->loadBitmapTable(path, outerr);
}

static void gfx_loadIntoBitmapTable(const struct playdate_graphics* gfx, const char* path, LCDBitmapTable* table, const char** outerr) {
    gfx->loadIntoBitmapTable(path, table, outerr);
}

static LCDBitmap* gfx_getTableBitmap(const struct playdate_graphics* gfx, LCDBitmapTable* table, int idx) {
    return gfx->getTableBitmap(table, idx);
}

// LCDFont functions
static LCDFont* gfx_loadFont(const struct playdate_graphics* gfx, const char* path, const char** outErr) {
    return gfx->loadFont(path, outErr);
}

static LCDFontPage* gfx_getFontPage(const struct playdate_graphics* gfx, LCDFont* font, uint32_t c) {
    return gfx->getFontPage(font, c);
}

static LCDFontGlyph* gfx_getPageGlyph(const struct playdate_graphics* gfx, LCDFontPage* page, uint32_t c, LCDBitmap** bitmap, int* advance) {
    return gfx->getPageGlyph(page, c, bitmap, advance);
}

static int gfx_getGlyphKerning(const struct playdate_graphics* gfx, LCDFontGlyph* glyph, uint32_t glyphcode, uint32_t nextcode) {
    return gfx->getGlyphKerning(glyph, glyphcode, nextcode);
}

static int gfx_getTextWidth(const struct playdate_graphics* gfx, LCDFont* font, const void* text, size_t len, PDStringEncoding encoding, int tracking) {
    return gfx->getTextWidth(font, text, len, encoding, tracking);
}

// Framebuffer functions
static uint8_t* gfx_getFrame(const struct playdate_graphics* gfx) {
    return gfx->getFrame();
}

static uint8_t* gfx_getDisplayFrame(const struct playdate_graphics* gfx) {
    return gfx->getDisplayFrame();
}

static LCDBitmap* gfx_getDebugBitmap(const struct playdate_graphics* gfx) {
    if (gfx->getDebugBitmap != NULL) {
        return gfx->getDebugBitmap();
    }
    return NULL;
}

static LCDBitmap* gfx_copyFrameBufferBitmap(const struct playdate_graphics* gfx) {
    return gfx->copyFrameBufferBitmap();
}

static void gfx_markUpdatedRows(const struct playdate_graphics* gfx, int start, int end) {
    gfx->markUpdatedRows(start, end);
}

static void gfx_display(const struct playdate_graphics* gfx) {
    gfx->display();
}

// Misc util
static void gfx_setColorToPattern(const struct playdate_graphics* gfx, LCDColor* color, LCDBitmap* bitmap, int x, int y) {
    gfx->setColorToPattern(color, bitmap, x, y);
}

static int gfx_checkMaskCollision(const struct playdate_graphics* gfx, LCDBitmap* bitmap1, int x1, int y1, LCDBitmapFlip flip1, LCDBitmap* bitmap2, int x2, int y2, LCDBitmapFlip flip2, LCDRect rect) {
    return gfx->checkMaskCollision(bitmap1, x1, y1, flip1, bitmap2, x2, y2, flip2, rect);
}

// v1.1
static void gfx_setScreenClipRect(const struct playdate_graphics* gfx, int x, int y, int width, int height) {
    gfx->setScreenClipRect(x, y, width, height);
}

// v1.1.1
static void gfx_fillPolygon(const struct playdate_graphics* gfx, int nPoints, int* coords, LCDColor color, LCDPolygonFillRule fillrule) {
    gfx->fillPolygon(nPoints, coords, color, fillrule);
}

static uint8_t gfx_getFontHeight(const struct playdate_graphics* gfx, LCDFont* font) {
    return gfx->getFontHeight(font);
}

// v1.7
static LCDBitmap* gfx_getDisplayBufferBitmap(const struct playdate_graphics* gfx) {
    return gfx->getDisplayBufferBitmap();
}

static void gfx_drawRotatedBitmap(const struct playdate_graphics* gfx, LCDBitmap* bitmap, int x, int y, float rotation, float centerx, float centery, float xscale, float yscale) {
    gfx->drawRotatedBitmap(bitmap, x, y, rotation, centerx, centery, xscale, yscale);
}

static void gfx_setTextLeading(const struct playdate_graphics* gfx, int lineHeightAdjustment) {
    gfx->setTextLeading(lineHeightAdjustment);
}

// v1.8
static int gfx_setBitmapMask(const struct playdate_graphics* gfx, LCDBitmap* bitmap, LCDBitmap* mask) {
    return gfx->setBitmapMask(bitmap, mask);
}

static LCDBitmap* gfx_getBitmapMask(const struct playdate_graphics* gfx, LCDBitmap* bitmap) {
    return gfx->getBitmapMask(bitmap);
}

// v1.10
static void gfx_setStencilImage(const struct playdate_graphics* gfx, LCDBitmap* stencil, int tile) {
    gfx->setStencilImage(stencil, tile);
}

// v1.12
static LCDFont* gfx_makeFontFromData(const struct playdate_graphics* gfx, LCDFontData* data, int wide) {
    return gfx->makeFontFromData(data, wide);
}

// v2.1
static int gfx_getTextTracking(const struct playdate_graphics* gfx) {
    return gfx->getTextTracking();
}

// v2.5
static void gfx_setPixel(const struct playdate_graphics* gfx, int x, int y, LCDColor c) {
    gfx->setPixel(x, y, c);
}

static LCDSolidColor gfx_getBitmapPixel(const struct playdate_graphics* gfx, LCDBitmap* bitmap, int x, int y) {
    return gfx->getBitmapPixel(bitmap, x, y);
}

static void gfx_getBitmapTableInfo(const struct playdate_graphics* gfx, LCDBitmapTable* table, int* count, int* width) {
    gfx->getBitmapTableInfo(table, count, width);
}

// v2.6
static void gfx_drawTextInRect(const struct playdate_graphics* gfx, const void* text, size_t len, PDStringEncoding encoding, int x, int y, int width, int height, PDTextWrappingMode wrap, PDTextAlignment align) {
    gfx->drawTextInRect(text, len, encoding, x, y, width, height, wrap, align);
}

// v2.7
static int gfx_getTextHeightForMaxWidth(const struct playdate_graphics* gfx, LCDFont* font, const void* text, size_t len, int maxwidth, PDStringEncoding encoding, PDTextWrappingMode wrap, int tracking, int extraLeading) {
    return gfx->getTextHeightForMaxWidth(font, text, len, maxwidth, encoding, wrap, tracking, extraLeading);
}

static void gfx_drawRoundRect(const struct playdate_graphics* gfx, int x, int y, int width, int height, int radius, int lineWidth, LCDColor color) {
    gfx->drawRoundRect(x, y, width, height, radius, lineWidth, color);
}

static void gfx_fillRoundRect(const struct playdate_graphics* gfx, int x, int y, int width, int height, int radius, LCDColor color) {
    gfx->fillRoundRect(x, y, width, height, radius, color);
}

// Video functions
static LCDVideoPlayer* gfx_video_loadVideo(const struct playdate_graphics* gfx, const char* path) {
    return gfx->video->loadVideo(path);
}

static void gfx_video_freePlayer(const struct playdate_graphics* gfx, LCDVideoPlayer* p) {
    gfx->video->freePlayer(p);
}

static int gfx_video_setContext(const struct playdate_graphics* gfx, LCDVideoPlayer* p, LCDBitmap* context) {
    return gfx->video->setContext(p, context);
}

static void gfx_video_useScreenContext(const struct playdate_graphics* gfx, LCDVideoPlayer* p) {
    gfx->video->useScreenContext(p);
}

static int gfx_video_renderFrame(const struct playdate_graphics* gfx, LCDVideoPlayer* p, int n) {
    return gfx->video->renderFrame(p, n);
}

static const char* gfx_video_getError(const struct playdate_graphics* gfx, LCDVideoPlayer* p) {
    return gfx->video->getError(p);
}

static void gfx_video_getInfo(const struct playdate_graphics* gfx, LCDVideoPlayer* p, int* outWidth, int* outHeight, float* outFrameRate, int* outFrameCount, int* outCurrentFrame) {
    gfx->video->getInfo(p, outWidth, outHeight, outFrameRate, outFrameCount, outCurrentFrame);
}

static LCDBitmap* gfx_video_getContext(const struct playdate_graphics* gfx, LCDVideoPlayer* p) {
    return gfx->video->getContext(p);
}

// Tilemap functions (v3.0)
static LCDTileMap* gfx_tilemap_newTilemap(const struct playdate_graphics* gfx) {
    return gfx->tilemap->newTilemap();
}

static void gfx_tilemap_freeTilemap(const struct playdate_graphics* gfx, LCDTileMap* m) {
    gfx->tilemap->freeTilemap(m);
}

static void gfx_tilemap_setImageTable(const struct playdate_graphics* gfx, LCDTileMap* m, LCDBitmapTable* table) {
    gfx->tilemap->setImageTable(m, table);
}

static LCDBitmapTable* gfx_tilemap_getImageTable(const struct playdate_graphics* gfx, LCDTileMap* m) {
    return gfx->tilemap->getImageTable(m);
}

static void gfx_tilemap_setSize(const struct playdate_graphics* gfx, LCDTileMap* m, int tilesWide, int tilesHigh) {
    gfx->tilemap->setSize(m, tilesWide, tilesHigh);
}

static void gfx_tilemap_getSize(const struct playdate_graphics* gfx, LCDTileMap* m, int* tilesWide, int* tilesHigh) {
    gfx->tilemap->getSize(m, tilesWide, tilesHigh);
}

static void gfx_tilemap_getPixelSize(const struct playdate_graphics* gfx, LCDTileMap* m, uint32_t* outWidth, uint32_t* outHeight) {
    gfx->tilemap->getPixelSize(m, outWidth, outHeight);
}

static void gfx_tilemap_setTiles(const struct playdate_graphics* gfx, LCDTileMap* m, uint16_t* indexes, int count, int rowwidth) {
    gfx->tilemap->setTiles(m, indexes, count, rowwidth);
}

static void gfx_tilemap_setTileAtPosition(const struct playdate_graphics* gfx, LCDTileMap* m, int x, int y, uint16_t idx) {
    gfx->tilemap->setTileAtPosition(m, x, y, idx);
}

static int gfx_tilemap_getTileAtPosition(const struct playdate_graphics* gfx, LCDTileMap* m, int x, int y) {
    return gfx->tilemap->getTileAtPosition(m, x, y);
}

static void gfx_tilemap_drawAtPoint(const struct playdate_graphics* gfx, LCDTileMap* m, float x, float y) {
    gfx->tilemap->drawAtPoint(m, x, y);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

// LCD display constants
const (
	LCDColumns = 400
	LCDRows    = 240
	LCDRowSize = 52
)

// LCDBitmapDrawMode represents bitmap drawing modes
type LCDBitmapDrawMode int

const (
	DrawModeCopy             LCDBitmapDrawMode = C.kDrawModeCopy
	DrawModeWhiteTransparent LCDBitmapDrawMode = C.kDrawModeWhiteTransparent
	DrawModeBlackTransparent LCDBitmapDrawMode = C.kDrawModeBlackTransparent
	DrawModeFillWhite        LCDBitmapDrawMode = C.kDrawModeFillWhite
	DrawModeFillBlack        LCDBitmapDrawMode = C.kDrawModeFillBlack
	DrawModeXOR              LCDBitmapDrawMode = C.kDrawModeXOR
	DrawModeNXOR             LCDBitmapDrawMode = C.kDrawModeNXOR
	DrawModeInverted         LCDBitmapDrawMode = C.kDrawModeInverted
)

// LCDBitmapFlip represents bitmap flip modes
type LCDBitmapFlip int

const (
	BitmapUnflipped LCDBitmapFlip = C.kBitmapUnflipped
	BitmapFlippedX  LCDBitmapFlip = C.kBitmapFlippedX
	BitmapFlippedY  LCDBitmapFlip = C.kBitmapFlippedY
	BitmapFlippedXY LCDBitmapFlip = C.kBitmapFlippedXY
)

// LCDSolidColor represents solid colors
type LCDSolidColor int

const (
	ColorBlack LCDSolidColor = C.kColorBlack
	ColorWhite LCDSolidColor = C.kColorWhite
	ColorClear LCDSolidColor = C.kColorClear
	ColorXOR   LCDSolidColor = C.kColorXOR
)

// LCDLineCapStyle represents line cap styles
type LCDLineCapStyle int

const (
	LineCapStyleButt   LCDLineCapStyle = C.kLineCapStyleButt
	LineCapStyleSquare LCDLineCapStyle = C.kLineCapStyleSquare
	LineCapStyleRound  LCDLineCapStyle = C.kLineCapStyleRound
)

// PDStringEncoding represents string encodings
type PDStringEncoding int

const (
	ASCIIEncoding   PDStringEncoding = C.kASCIIEncoding
	UTF8Encoding    PDStringEncoding = C.kUTF8Encoding
	Encoding16BitLE PDStringEncoding = C.k16BitLEEncoding
)

// LCDPolygonFillRule represents polygon fill rules
type LCDPolygonFillRule int

const (
	PolygonFillNonZero LCDPolygonFillRule = C.kPolygonFillNonZero
	PolygonFillEvenOdd LCDPolygonFillRule = C.kPolygonFillEvenOdd
)

// PDTextWrappingMode represents text wrapping modes
type PDTextWrappingMode int

const (
	WrapClip      PDTextWrappingMode = C.kWrapClip
	WrapCharacter PDTextWrappingMode = C.kWrapCharacter
	WrapWord      PDTextWrappingMode = C.kWrapWord
)

// PDTextAlignment represents text alignment
type PDTextAlignment int

const (
	AlignTextLeft   PDTextAlignment = C.kAlignTextLeft
	AlignTextCenter PDTextAlignment = C.kAlignTextCenter
	AlignTextRight  PDTextAlignment = C.kAlignTextRight
)

// LCDRect represents a rectangle
type LCDRect struct {
	Left   int
	Right  int // not inclusive
	Top    int
	Bottom int // not inclusive
}

// LCDColor represents a color (solid or pattern)
type LCDColor uintptr

// NewColorFromSolid creates a color from a solid color
func NewColorFromSolid(color LCDSolidColor) LCDColor {
	return LCDColor(color)
}

// LCDPattern represents an 8x8 pattern (8 rows image data, 8 rows mask)
type LCDPattern [16]uint8

// NewColorFromPattern creates a color from a pattern
func NewColorFromPattern(pattern *LCDPattern) LCDColor {
	return LCDColor(uintptr(unsafe.Pointer(pattern)))
}

// LCDBitmap wraps a Playdate bitmap
type LCDBitmap struct {
	ptr *C.LCDBitmap
}

// LCDBitmapTable wraps a Playdate bitmap table
type LCDBitmapTable struct {
	ptr *C.LCDBitmapTable
}

// LCDFont wraps a Playdate font
type LCDFont struct {
	ptr *C.LCDFont
}

// LCDFontPage wraps a font page
type LCDFontPage struct {
	ptr *C.LCDFontPage
}

// LCDFontGlyph wraps a font glyph
type LCDFontGlyph struct {
	ptr *C.LCDFontGlyph
}

// LCDTileMap wraps a Playdate tilemap
type LCDTileMap struct {
	ptr *C.LCDTileMap
}

// LCDVideoPlayer wraps a Playdate video player
type LCDVideoPlayer struct {
	ptr *C.LCDVideoPlayer
}

// Graphics wraps the playdate_graphics API
type Graphics struct {
	ptr     *C.struct_playdate_graphics
	Video   *Video
	Tilemap *Tilemap
}

func newGraphics(ptr *C.struct_playdate_graphics) *Graphics {
	g := &Graphics{ptr: ptr}
	g.Video = &Video{gfx: ptr}
	g.Tilemap = &Tilemap{gfx: ptr}
	return g
}

// Clear clears the screen with the given color
func (g *Graphics) Clear(color LCDColor) {
	C.gfx_clear(g.ptr, C.LCDColor(color))
}

// SetBackgroundColor sets the background color
func (g *Graphics) SetBackgroundColor(color LCDSolidColor) {
	C.gfx_setBackgroundColor(g.ptr, C.LCDSolidColor(color))
}

// SetStencil sets the stencil bitmap (deprecated, use SetStencilImage)
func (g *Graphics) SetStencil(stencil *LCDBitmap) {
	var s *C.LCDBitmap
	if stencil != nil {
		s = stencil.ptr
	}
	C.gfx_setStencil(g.ptr, s)
}

// SetDrawMode sets the draw mode and returns the previous mode
func (g *Graphics) SetDrawMode(mode LCDBitmapDrawMode) LCDBitmapDrawMode {
	return LCDBitmapDrawMode(C.gfx_setDrawMode(g.ptr, C.LCDBitmapDrawMode(mode)))
}

// SetDrawOffset sets the draw offset
func (g *Graphics) SetDrawOffset(dx, dy int) {
	C.gfx_setDrawOffset(g.ptr, C.int(dx), C.int(dy))
}

// SetClipRect sets the clipping rectangle
func (g *Graphics) SetClipRect(x, y, width, height int) {
	C.gfx_setClipRect(g.ptr, C.int(x), C.int(y), C.int(width), C.int(height))
}

// ClearClipRect clears the clipping rectangle
func (g *Graphics) ClearClipRect() {
	C.gfx_clearClipRect(g.ptr)
}

// SetLineCapStyle sets the line cap style
func (g *Graphics) SetLineCapStyle(style LCDLineCapStyle) {
	C.gfx_setLineCapStyle(g.ptr, C.LCDLineCapStyle(style))
}

// SetFont sets the current font
func (g *Graphics) SetFont(font *LCDFont) {
	var f *C.LCDFont
	if font != nil {
		f = font.ptr
	}
	C.gfx_setFont(g.ptr, f)
}

// SetTextTracking sets the text tracking (letter spacing)
func (g *Graphics) SetTextTracking(tracking int) {
	C.gfx_setTextTracking(g.ptr, C.int(tracking))
}

// GetTextTracking returns the current text tracking
func (g *Graphics) GetTextTracking() int {
	return int(C.gfx_getTextTracking(g.ptr))
}

// PushContext pushes a new drawing context
func (g *Graphics) PushContext(target *LCDBitmap) {
	var t *C.LCDBitmap
	if target != nil {
		t = target.ptr
	}
	C.gfx_pushContext(g.ptr, t)
}

// PopContext pops the current drawing context
func (g *Graphics) PopContext() {
	C.gfx_popContext(g.ptr)
}

// DrawBitmap draws a bitmap at the given position
func (g *Graphics) DrawBitmap(bitmap *LCDBitmap, x, y int, flip LCDBitmapFlip) {
	if bitmap == nil {
		return
	}
	C.gfx_drawBitmap(g.ptr, bitmap.ptr, C.int(x), C.int(y), C.LCDBitmapFlip(flip))
}

// TileBitmap tiles a bitmap in the given area
func (g *Graphics) TileBitmap(bitmap *LCDBitmap, x, y, width, height int, flip LCDBitmapFlip) {
	if bitmap == nil {
		return
	}
	C.gfx_tileBitmap(g.ptr, bitmap.ptr, C.int(x), C.int(y), C.int(width), C.int(height), C.LCDBitmapFlip(flip))
}

// DrawLine draws a line
func (g *Graphics) DrawLine(x1, y1, x2, y2, width int, color LCDColor) {
	C.gfx_drawLine(g.ptr, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(width), C.LCDColor(color))
}

// FillTriangle fills a triangle
func (g *Graphics) FillTriangle(x1, y1, x2, y2, x3, y3 int, color LCDColor) {
	C.gfx_fillTriangle(g.ptr, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(x3), C.int(y3), C.LCDColor(color))
}

// DrawRect draws a rectangle
func (g *Graphics) DrawRect(x, y, width, height int, color LCDColor) {
	C.gfx_drawRect(g.ptr, C.int(x), C.int(y), C.int(width), C.int(height), C.LCDColor(color))
}

// FillRect fills a rectangle
func (g *Graphics) FillRect(x, y, width, height int, color LCDColor) {
	C.gfx_fillRect(g.ptr, C.int(x), C.int(y), C.int(width), C.int(height), C.LCDColor(color))
}

// DrawEllipse draws an ellipse
func (g *Graphics) DrawEllipse(x, y, width, height, lineWidth int, startAngle, endAngle float32, color LCDColor) {
	C.gfx_drawEllipse(g.ptr, C.int(x), C.int(y), C.int(width), C.int(height), C.int(lineWidth), C.float(startAngle), C.float(endAngle), C.LCDColor(color))
}

// FillEllipse fills an ellipse
func (g *Graphics) FillEllipse(x, y, width, height int, startAngle, endAngle float32, color LCDColor) {
	C.gfx_fillEllipse(g.ptr, C.int(x), C.int(y), C.int(width), C.int(height), C.float(startAngle), C.float(endAngle), C.LCDColor(color))
}

// DrawScaledBitmap draws a scaled bitmap
func (g *Graphics) DrawScaledBitmap(bitmap *LCDBitmap, x, y int, xscale, yscale float32) {
	if bitmap == nil {
		return
	}
	C.gfx_drawScaledBitmap(g.ptr, bitmap.ptr, C.int(x), C.int(y), C.float(xscale), C.float(yscale))
}

// DrawText draws text and returns its width
func (g *Graphics) DrawText(text string, x, y int) int {
	cstr := cString(text)
	defer freeCString(cstr)
	return int(C.gfx_drawText(g.ptr, unsafe.Pointer(cstr), C.size_t(len(text)), C.PDStringEncoding(UTF8Encoding), C.int(x), C.int(y)))
}

// DrawTextWithEncoding draws text with a specific encoding
func (g *Graphics) DrawTextWithEncoding(text string, encoding PDStringEncoding, x, y int) int {
	cstr := cString(text)
	defer freeCString(cstr)
	return int(C.gfx_drawText(g.ptr, unsafe.Pointer(cstr), C.size_t(len(text)), C.PDStringEncoding(encoding), C.int(x), C.int(y)))
}

// DrawRotatedBitmap draws a rotated bitmap
func (g *Graphics) DrawRotatedBitmap(bitmap *LCDBitmap, x, y int, rotation, centerX, centerY, xscale, yscale float32) {
	if bitmap == nil {
		return
	}
	C.gfx_drawRotatedBitmap(g.ptr, bitmap.ptr, C.int(x), C.int(y), C.float(rotation), C.float(centerX), C.float(centerY), C.float(xscale), C.float(yscale))
}

// DrawRoundRect draws a rounded rectangle
func (g *Graphics) DrawRoundRect(x, y, width, height, radius, lineWidth int, color LCDColor) {
	C.gfx_drawRoundRect(g.ptr, C.int(x), C.int(y), C.int(width), C.int(height), C.int(radius), C.int(lineWidth), C.LCDColor(color))
}

// FillRoundRect fills a rounded rectangle
func (g *Graphics) FillRoundRect(x, y, width, height, radius int, color LCDColor) {
	C.gfx_fillRoundRect(g.ptr, C.int(x), C.int(y), C.int(width), C.int(height), C.int(radius), C.LCDColor(color))
}

// DrawTextInRect draws text in a rectangle with wrapping and alignment
func (g *Graphics) DrawTextInRect(text string, x, y, width, height int, wrap PDTextWrappingMode, align PDTextAlignment) {
	cstr := cString(text)
	defer freeCString(cstr)
	C.gfx_drawTextInRect(g.ptr, unsafe.Pointer(cstr), C.size_t(len(text)), C.PDStringEncoding(UTF8Encoding), C.int(x), C.int(y), C.int(width), C.int(height), C.PDTextWrappingMode(wrap), C.PDTextAlignment(align))
}

// SetTextLeading sets the text leading (line height adjustment)
func (g *Graphics) SetTextLeading(lineHeightAdjustment int) {
	C.gfx_setTextLeading(g.ptr, C.int(lineHeightAdjustment))
}

// FillPolygon fills a polygon
func (g *Graphics) FillPolygon(points []int, color LCDColor, fillRule LCDPolygonFillRule) {
	if len(points) < 6 { // Need at least 3 points (6 coords)
		return
	}
	coords := make([]C.int, len(points))
	for i, p := range points {
		coords[i] = C.int(p)
	}
	C.gfx_fillPolygon(g.ptr, C.int(len(points)/2), (*C.int)(&coords[0]), C.LCDColor(color), C.LCDPolygonFillRule(fillRule))
}

// SetPixel sets a single pixel
func (g *Graphics) SetPixel(x, y int, color LCDColor) {
	C.gfx_setPixel(g.ptr, C.int(x), C.int(y), C.LCDColor(color))
}

// SetScreenClipRect sets the screen clipping rectangle
func (g *Graphics) SetScreenClipRect(x, y, width, height int) {
	C.gfx_setScreenClipRect(g.ptr, C.int(x), C.int(y), C.int(width), C.int(height))
}

// SetStencilImage sets the stencil image with optional tiling
func (g *Graphics) SetStencilImage(stencil *LCDBitmap, tile bool) {
	var s *C.LCDBitmap
	if stencil != nil {
		s = stencil.ptr
	}
	t := 0
	if tile {
		t = 1
	}
	C.gfx_setStencilImage(g.ptr, s, C.int(t))
}

// NewBitmap creates a new bitmap
func (g *Graphics) NewBitmap(width, height int, bgcolor LCDColor) *LCDBitmap {
	ptr := C.gfx_newBitmap(g.ptr, C.int(width), C.int(height), C.LCDColor(bgcolor))
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// FreeBitmap frees a bitmap
func (g *Graphics) FreeBitmap(bitmap *LCDBitmap) {
	if bitmap != nil && bitmap.ptr != nil {
		C.gfx_freeBitmap(g.ptr, bitmap.ptr)
		bitmap.ptr = nil
	}
}

// LoadBitmap loads a bitmap from a file
func (g *Graphics) LoadBitmap(path string) (*LCDBitmap, error) {
	cpath := cString(path)
	defer freeCString(cpath)

	var outerr *C.char
	ptr := C.gfx_loadBitmap(g.ptr, cpath, &outerr)
	if ptr == nil {
		if outerr != nil {
			return nil, errors.New(goString(outerr))
		}
		return nil, errors.New("failed to load bitmap")
	}
	return &LCDBitmap{ptr: ptr}, nil
}

// CopyBitmap creates a copy of a bitmap
func (g *Graphics) CopyBitmap(bitmap *LCDBitmap) *LCDBitmap {
	if bitmap == nil {
		return nil
	}
	ptr := C.gfx_copyBitmap(g.ptr, bitmap.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// LoadIntoBitmap loads a bitmap from a file into an existing bitmap
func (g *Graphics) LoadIntoBitmap(path string, bitmap *LCDBitmap) error {
	if bitmap == nil {
		return errors.New("bitmap is nil")
	}
	cpath := cString(path)
	defer freeCString(cpath)

	var outerr *C.char
	C.gfx_loadIntoBitmap(g.ptr, cpath, bitmap.ptr, &outerr)
	if outerr != nil {
		return errors.New(goString(outerr))
	}
	return nil
}

// BitmapData contains bitmap data information
type BitmapData struct {
	Width    int
	Height   int
	RowBytes int
	Mask     []byte
	Data     []byte
}

// GetBitmapData returns the bitmap's data
func (g *Graphics) GetBitmapData(bitmap *LCDBitmap) *BitmapData {
	if bitmap == nil {
		return nil
	}
	var width, height, rowbytes C.int
	var mask, data *C.uint8_t
	C.gfx_getBitmapData(g.ptr, bitmap.ptr, &width, &height, &rowbytes, &mask, &data)

	bd := &BitmapData{
		Width:    int(width),
		Height:   int(height),
		RowBytes: int(rowbytes),
	}

	dataSize := int(height) * int(rowbytes)
	if data != nil {
		bd.Data = C.GoBytes(unsafe.Pointer(data), C.int(dataSize))
	}
	if mask != nil {
		bd.Mask = C.GoBytes(unsafe.Pointer(mask), C.int(dataSize))
	}

	return bd
}

// ClearBitmap clears a bitmap
func (g *Graphics) ClearBitmap(bitmap *LCDBitmap, bgcolor LCDColor) {
	if bitmap != nil {
		C.gfx_clearBitmap(g.ptr, bitmap.ptr, C.LCDColor(bgcolor))
	}
}

// RotatedBitmap returns a rotated copy of a bitmap
func (g *Graphics) RotatedBitmap(bitmap *LCDBitmap, rotation, xscale, yscale float32) (*LCDBitmap, int) {
	if bitmap == nil {
		return nil, 0
	}
	var allocedSize C.int
	ptr := C.gfx_rotatedBitmap(g.ptr, bitmap.ptr, C.float(rotation), C.float(xscale), C.float(yscale), &allocedSize)
	if ptr == nil {
		return nil, 0
	}
	return &LCDBitmap{ptr: ptr}, int(allocedSize)
}

// SetBitmapMask sets the mask for a bitmap
func (g *Graphics) SetBitmapMask(bitmap, mask *LCDBitmap) int {
	if bitmap == nil {
		return 0
	}
	var m *C.LCDBitmap
	if mask != nil {
		m = mask.ptr
	}
	return int(C.gfx_setBitmapMask(g.ptr, bitmap.ptr, m))
}

// GetBitmapMask returns the mask for a bitmap
func (g *Graphics) GetBitmapMask(bitmap *LCDBitmap) *LCDBitmap {
	if bitmap == nil {
		return nil
	}
	ptr := C.gfx_getBitmapMask(g.ptr, bitmap.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// GetBitmapPixel returns the color of a pixel in a bitmap
func (g *Graphics) GetBitmapPixel(bitmap *LCDBitmap, x, y int) LCDSolidColor {
	if bitmap == nil {
		return ColorClear
	}
	return LCDSolidColor(C.gfx_getBitmapPixel(g.ptr, bitmap.ptr, C.int(x), C.int(y)))
}

// NewBitmapTable creates a new bitmap table
func (g *Graphics) NewBitmapTable(count, width, height int) *LCDBitmapTable {
	ptr := C.gfx_newBitmapTable(g.ptr, C.int(count), C.int(width), C.int(height))
	if ptr == nil {
		return nil
	}
	return &LCDBitmapTable{ptr: ptr}
}

// FreeBitmapTable frees a bitmap table
func (g *Graphics) FreeBitmapTable(table *LCDBitmapTable) {
	if table != nil && table.ptr != nil {
		C.gfx_freeBitmapTable(g.ptr, table.ptr)
		table.ptr = nil
	}
}

// LoadBitmapTable loads a bitmap table from a file
func (g *Graphics) LoadBitmapTable(path string) (*LCDBitmapTable, error) {
	cpath := cString(path)
	defer freeCString(cpath)

	var outerr *C.char
	ptr := C.gfx_loadBitmapTable(g.ptr, cpath, &outerr)
	if ptr == nil {
		if outerr != nil {
			return nil, errors.New(goString(outerr))
		}
		return nil, errors.New("failed to load bitmap table")
	}
	return &LCDBitmapTable{ptr: ptr}, nil
}

// LoadIntoBitmapTable loads a bitmap table from a file into an existing table
func (g *Graphics) LoadIntoBitmapTable(path string, table *LCDBitmapTable) error {
	if table == nil {
		return errors.New("table is nil")
	}
	cpath := cString(path)
	defer freeCString(cpath)

	var outerr *C.char
	C.gfx_loadIntoBitmapTable(g.ptr, cpath, table.ptr, &outerr)
	if outerr != nil {
		return errors.New(goString(outerr))
	}
	return nil
}

// GetTableBitmap returns a bitmap from a bitmap table
func (g *Graphics) GetTableBitmap(table *LCDBitmapTable, index int) *LCDBitmap {
	if table == nil {
		return nil
	}
	ptr := C.gfx_getTableBitmap(g.ptr, table.ptr, C.int(index))
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// GetBitmapTableInfo returns info about a bitmap table
func (g *Graphics) GetBitmapTableInfo(table *LCDBitmapTable) (count, width int) {
	if table == nil {
		return 0, 0
	}
	var c, w C.int
	C.gfx_getBitmapTableInfo(g.ptr, table.ptr, &c, &w)
	return int(c), int(w)
}

// LoadFont loads a font from a file
func (g *Graphics) LoadFont(path string) (*LCDFont, error) {
	cpath := cString(path)
	defer freeCString(cpath)

	var outerr *C.char
	ptr := C.gfx_loadFont(g.ptr, cpath, &outerr)
	if ptr == nil {
		if outerr != nil {
			return nil, errors.New(goString(outerr))
		}
		return nil, errors.New("failed to load font")
	}
	return &LCDFont{ptr: ptr}, nil
}

// GetFontPage returns a font page for a character
func (g *Graphics) GetFontPage(font *LCDFont, c uint32) *LCDFontPage {
	if font == nil {
		return nil
	}
	ptr := C.gfx_getFontPage(g.ptr, font.ptr, C.uint32_t(c))
	if ptr == nil {
		return nil
	}
	return &LCDFontPage{ptr: ptr}
}

// GetPageGlyph returns a glyph from a font page
func (g *Graphics) GetPageGlyph(page *LCDFontPage, c uint32) (*LCDFontGlyph, *LCDBitmap, int) {
	if page == nil {
		return nil, nil, 0
	}
	var bitmap *C.LCDBitmap
	var advance C.int
	ptr := C.gfx_getPageGlyph(g.ptr, page.ptr, C.uint32_t(c), &bitmap, &advance)
	if ptr == nil {
		return nil, nil, 0
	}
	var bmp *LCDBitmap
	if bitmap != nil {
		bmp = &LCDBitmap{ptr: bitmap}
	}
	return &LCDFontGlyph{ptr: ptr}, bmp, int(advance)
}

// GetGlyphKerning returns the kerning between two glyphs
func (g *Graphics) GetGlyphKerning(glyph *LCDFontGlyph, glyphCode, nextCode uint32) int {
	if glyph == nil {
		return 0
	}
	return int(C.gfx_getGlyphKerning(g.ptr, glyph.ptr, C.uint32_t(glyphCode), C.uint32_t(nextCode)))
}

// GetTextWidth returns the width of text
func (g *Graphics) GetTextWidth(font *LCDFont, text string, tracking int) int {
	if font == nil {
		return 0
	}
	cstr := cString(text)
	defer freeCString(cstr)
	return int(C.gfx_getTextWidth(g.ptr, font.ptr, unsafe.Pointer(cstr), C.size_t(len(text)), C.PDStringEncoding(UTF8Encoding), C.int(tracking)))
}

// GetFontHeight returns the height of a font
func (g *Graphics) GetFontHeight(font *LCDFont) uint8 {
	if font == nil {
		return 0
	}
	return uint8(C.gfx_getFontHeight(g.ptr, font.ptr))
}

// GetTextHeightForMaxWidth returns the height of text for a max width
func (g *Graphics) GetTextHeightForMaxWidth(font *LCDFont, text string, maxWidth int, wrap PDTextWrappingMode, tracking, extraLeading int) int {
	if font == nil {
		return 0
	}
	cstr := cString(text)
	defer freeCString(cstr)
	return int(C.gfx_getTextHeightForMaxWidth(g.ptr, font.ptr, unsafe.Pointer(cstr), C.size_t(len(text)), C.int(maxWidth), C.PDStringEncoding(UTF8Encoding), C.PDTextWrappingMode(wrap), C.int(tracking), C.int(extraLeading)))
}

// GetFrame returns the current frame buffer as a slice pointing directly to C memory.
// Modifications to this slice will affect the actual frame buffer.
func (g *Graphics) GetFrame() []byte {
	ptr := C.gfx_getFrame(g.ptr)
	if ptr == nil {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), LCDRows*LCDRowSize)
}

// GetDisplayFrame returns the display frame buffer as a slice pointing directly to C memory.
// Modifications to this slice will affect the actual display buffer.
func (g *Graphics) GetDisplayFrame() []byte {
	ptr := C.gfx_getDisplayFrame(g.ptr)
	if ptr == nil {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), LCDRows*LCDRowSize)
}

// GetDebugBitmap returns the debug bitmap (simulator only)
func (g *Graphics) GetDebugBitmap() *LCDBitmap {
	ptr := C.gfx_getDebugBitmap(g.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// CopyFrameBufferBitmap returns a copy of the frame buffer as a bitmap
func (g *Graphics) CopyFrameBufferBitmap() *LCDBitmap {
	ptr := C.gfx_copyFrameBufferBitmap(g.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// MarkUpdatedRows marks rows as updated
func (g *Graphics) MarkUpdatedRows(start, end int) {
	C.gfx_markUpdatedRows(g.ptr, C.int(start), C.int(end))
}

// Display updates the display
func (g *Graphics) Display() {
	C.gfx_display(g.ptr)
}

// GetDisplayBufferBitmap returns the display buffer as a bitmap
func (g *Graphics) GetDisplayBufferBitmap() *LCDBitmap {
	ptr := C.gfx_getDisplayBufferBitmap(g.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// SetColorToPattern sets a color to a pattern from a bitmap
func (g *Graphics) SetColorToPattern(bitmap *LCDBitmap, x, y int) LCDColor {
	if bitmap == nil {
		return 0
	}
	var color C.LCDColor
	C.gfx_setColorToPattern(g.ptr, &color, bitmap.ptr, C.int(x), C.int(y))
	return LCDColor(color)
}

// CheckMaskCollision checks for collision between two bitmaps using their masks
func (g *Graphics) CheckMaskCollision(bitmap1 *LCDBitmap, x1, y1 int, flip1 LCDBitmapFlip, bitmap2 *LCDBitmap, x2, y2 int, flip2 LCDBitmapFlip, rect LCDRect) bool {
	if bitmap1 == nil || bitmap2 == nil {
		return false
	}
	crect := C.LCDRect{
		left:   C.int(rect.Left),
		right:  C.int(rect.Right),
		top:    C.int(rect.Top),
		bottom: C.int(rect.Bottom),
	}
	return C.gfx_checkMaskCollision(g.ptr, bitmap1.ptr, C.int(x1), C.int(y1), C.LCDBitmapFlip(flip1), bitmap2.ptr, C.int(x2), C.int(y2), C.LCDBitmapFlip(flip2), crect) != 0
}

// Video wraps video functions
type Video struct {
	gfx *C.struct_playdate_graphics
}

// LoadVideo loads a video from a file
func (v *Video) LoadVideo(path string) *LCDVideoPlayer {
	cpath := cString(path)
	defer freeCString(cpath)
	ptr := C.gfx_video_loadVideo(v.gfx, cpath)
	if ptr == nil {
		return nil
	}
	return &LCDVideoPlayer{ptr: ptr}
}

// FreePlayer frees a video player
func (v *Video) FreePlayer(player *LCDVideoPlayer) {
	if player != nil && player.ptr != nil {
		C.gfx_video_freePlayer(v.gfx, player.ptr)
		player.ptr = nil
	}
}

// SetContext sets the rendering context for a video player
func (v *Video) SetContext(player *LCDVideoPlayer, context *LCDBitmap) int {
	if player == nil {
		return 0
	}
	var ctx *C.LCDBitmap
	if context != nil {
		ctx = context.ptr
	}
	return int(C.gfx_video_setContext(v.gfx, player.ptr, ctx))
}

// UseScreenContext sets the video player to use the screen context
func (v *Video) UseScreenContext(player *LCDVideoPlayer) {
	if player != nil {
		C.gfx_video_useScreenContext(v.gfx, player.ptr)
	}
}

// RenderFrame renders a frame from the video
func (v *Video) RenderFrame(player *LCDVideoPlayer, frame int) int {
	if player == nil {
		return 0
	}
	return int(C.gfx_video_renderFrame(v.gfx, player.ptr, C.int(frame)))
}

// GetError returns the last error for a video player
func (v *Video) GetError(player *LCDVideoPlayer) string {
	if player == nil {
		return ""
	}
	err := C.gfx_video_getError(v.gfx, player.ptr)
	if err == nil {
		return ""
	}
	return goString(err)
}

// VideoInfo contains video information
type VideoInfo struct {
	Width        int
	Height       int
	FrameRate    float32
	FrameCount   int
	CurrentFrame int
}

// GetInfo returns information about a video
func (v *Video) GetInfo(player *LCDVideoPlayer) *VideoInfo {
	if player == nil {
		return nil
	}
	var width, height, frameCount, currentFrame C.int
	var frameRate C.float
	C.gfx_video_getInfo(v.gfx, player.ptr, &width, &height, &frameRate, &frameCount, &currentFrame)
	return &VideoInfo{
		Width:        int(width),
		Height:       int(height),
		FrameRate:    float32(frameRate),
		FrameCount:   int(frameCount),
		CurrentFrame: int(currentFrame),
	}
}

// GetContext returns the rendering context for a video player
func (v *Video) GetContext(player *LCDVideoPlayer) *LCDBitmap {
	if player == nil {
		return nil
	}
	ptr := C.gfx_video_getContext(v.gfx, player.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// Tilemap wraps tilemap functions
type Tilemap struct {
	gfx *C.struct_playdate_graphics
}

// NewTilemap creates a new tilemap
func (t *Tilemap) NewTilemap() *LCDTileMap {
	ptr := C.gfx_tilemap_newTilemap(t.gfx)
	if ptr == nil {
		return nil
	}
	return &LCDTileMap{ptr: ptr}
}

// FreeTilemap frees a tilemap
func (t *Tilemap) FreeTilemap(m *LCDTileMap) {
	if m != nil && m.ptr != nil {
		C.gfx_tilemap_freeTilemap(t.gfx, m.ptr)
		m.ptr = nil
	}
}

// SetImageTable sets the image table for a tilemap
func (t *Tilemap) SetImageTable(m *LCDTileMap, table *LCDBitmapTable) {
	if m == nil {
		return
	}
	var tbl *C.LCDBitmapTable
	if table != nil {
		tbl = table.ptr
	}
	C.gfx_tilemap_setImageTable(t.gfx, m.ptr, tbl)
}

// GetImageTable returns the image table for a tilemap
func (t *Tilemap) GetImageTable(m *LCDTileMap) *LCDBitmapTable {
	if m == nil {
		return nil
	}
	ptr := C.gfx_tilemap_getImageTable(t.gfx, m.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDBitmapTable{ptr: ptr}
}

// SetSize sets the size of a tilemap
func (t *Tilemap) SetSize(m *LCDTileMap, tilesWide, tilesHigh int) {
	if m != nil {
		C.gfx_tilemap_setSize(t.gfx, m.ptr, C.int(tilesWide), C.int(tilesHigh))
	}
}

// GetSize returns the size of a tilemap
func (t *Tilemap) GetSize(m *LCDTileMap) (tilesWide, tilesHigh int) {
	if m == nil {
		return 0, 0
	}
	var tw, th C.int
	C.gfx_tilemap_getSize(t.gfx, m.ptr, &tw, &th)
	return int(tw), int(th)
}

// GetPixelSize returns the pixel size of a tilemap
func (t *Tilemap) GetPixelSize(m *LCDTileMap) (width, height uint32) {
	if m == nil {
		return 0, 0
	}
	var w, h C.uint32_t
	C.gfx_tilemap_getPixelSize(t.gfx, m.ptr, &w, &h)
	return uint32(w), uint32(h)
}

// SetTiles sets the tiles for a tilemap
func (t *Tilemap) SetTiles(m *LCDTileMap, indexes []uint16, rowWidth int) {
	if m == nil || len(indexes) == 0 {
		return
	}
	C.gfx_tilemap_setTiles(t.gfx, m.ptr, (*C.uint16_t)(&indexes[0]), C.int(len(indexes)), C.int(rowWidth))
}

// SetTileAtPosition sets a tile at a position
func (t *Tilemap) SetTileAtPosition(m *LCDTileMap, x, y int, idx uint16) {
	if m != nil {
		C.gfx_tilemap_setTileAtPosition(t.gfx, m.ptr, C.int(x), C.int(y), C.uint16_t(idx))
	}
}

// GetTileAtPosition returns the tile at a position
func (t *Tilemap) GetTileAtPosition(m *LCDTileMap, x, y int) int {
	if m == nil {
		return 0
	}
	return int(C.gfx_tilemap_getTileAtPosition(t.gfx, m.ptr, C.int(x), C.int(y)))
}

// DrawAtPoint draws a tilemap at a position
func (t *Tilemap) DrawAtPoint(m *LCDTileMap, x, y float32) {
	if m != nil {
		C.gfx_tilemap_drawAtPoint(t.gfx, m.ptr, C.float(x), C.float(y))
	}
}
