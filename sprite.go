//go:build !tinygo

package pdgo

/*
#cgo CFLAGS: -DTARGET_EXTENSION=1
#include "pd_api.h"
#include <stdlib.h>

// Sprite API helper functions

static void sprite_setAlwaysRedraw(const struct playdate_sprite* spr, int flag) {
    spr->setAlwaysRedraw(flag);
}

static void sprite_addDirtyRect(const struct playdate_sprite* spr, LCDRect dirtyRect) {
    spr->addDirtyRect(dirtyRect);
}

static void sprite_drawSprites(const struct playdate_sprite* spr) {
    spr->drawSprites();
}

static void sprite_updateAndDrawSprites(const struct playdate_sprite* spr) {
    spr->updateAndDrawSprites();
}

static LCDSprite* sprite_newSprite(const struct playdate_sprite* spr) {
    return spr->newSprite();
}

static void sprite_freeSprite(const struct playdate_sprite* spr, LCDSprite* sprite) {
    spr->freeSprite(sprite);
}

static LCDSprite* sprite_copy(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->copy(sprite);
}

static void sprite_addSprite(const struct playdate_sprite* spr, LCDSprite* sprite) {
    spr->addSprite(sprite);
}

static void sprite_removeSprite(const struct playdate_sprite* spr, LCDSprite* sprite) {
    spr->removeSprite(sprite);
}

static void sprite_removeSprites(const struct playdate_sprite* spr, LCDSprite** sprites, int count) {
    spr->removeSprites(sprites, count);
}

static void sprite_removeAllSprites(const struct playdate_sprite* spr) {
    spr->removeAllSprites();
}

static int sprite_getSpriteCount(const struct playdate_sprite* spr) {
    return spr->getSpriteCount();
}

static void sprite_setBounds(const struct playdate_sprite* spr, LCDSprite* sprite, PDRect bounds) {
    spr->setBounds(sprite, bounds);
}

static PDRect sprite_getBounds(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->getBounds(sprite);
}

static void sprite_moveTo(const struct playdate_sprite* spr, LCDSprite* sprite, float x, float y) {
    spr->moveTo(sprite, x, y);
}

static void sprite_moveBy(const struct playdate_sprite* spr, LCDSprite* sprite, float dx, float dy) {
    spr->moveBy(sprite, dx, dy);
}

static void sprite_setImage(const struct playdate_sprite* spr, LCDSprite* sprite, LCDBitmap* image, LCDBitmapFlip flip) {
    spr->setImage(sprite, image, flip);
}

static LCDBitmap* sprite_getImage(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->getImage(sprite);
}

static void sprite_setSize(const struct playdate_sprite* spr, LCDSprite* s, float width, float height) {
    spr->setSize(s, width, height);
}

static void sprite_setZIndex(const struct playdate_sprite* spr, LCDSprite* sprite, int16_t zIndex) {
    spr->setZIndex(sprite, zIndex);
}

static int16_t sprite_getZIndex(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->getZIndex(sprite);
}

static void sprite_setDrawMode(const struct playdate_sprite* spr, LCDSprite* sprite, LCDBitmapDrawMode mode) {
    spr->setDrawMode(sprite, mode);
}

static void sprite_setImageFlip(const struct playdate_sprite* spr, LCDSprite* sprite, LCDBitmapFlip flip) {
    spr->setImageFlip(sprite, flip);
}

static LCDBitmapFlip sprite_getImageFlip(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->getImageFlip(sprite);
}

static void sprite_setStencil(const struct playdate_sprite* spr, LCDSprite* sprite, LCDBitmap* stencil) {
    spr->setStencil(sprite, stencil);
}

static void sprite_setClipRect(const struct playdate_sprite* spr, LCDSprite* sprite, LCDRect clipRect) {
    spr->setClipRect(sprite, clipRect);
}

static void sprite_clearClipRect(const struct playdate_sprite* spr, LCDSprite* sprite) {
    spr->clearClipRect(sprite);
}

static void sprite_setClipRectsInRange(const struct playdate_sprite* spr, LCDRect clipRect, int startZ, int endZ) {
    spr->setClipRectsInRange(clipRect, startZ, endZ);
}

static void sprite_clearClipRectsInRange(const struct playdate_sprite* spr, int startZ, int endZ) {
    spr->clearClipRectsInRange(startZ, endZ);
}

static void sprite_setUpdatesEnabled(const struct playdate_sprite* spr, LCDSprite* sprite, int flag) {
    spr->setUpdatesEnabled(sprite, flag);
}

static int sprite_updatesEnabled(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->updatesEnabled(sprite);
}

static void sprite_setCollisionsEnabled(const struct playdate_sprite* spr, LCDSprite* sprite, int flag) {
    spr->setCollisionsEnabled(sprite, flag);
}

static int sprite_collisionsEnabled(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->collisionsEnabled(sprite);
}

static void sprite_setVisible(const struct playdate_sprite* spr, LCDSprite* sprite, int flag) {
    spr->setVisible(sprite, flag);
}

static int sprite_isVisible(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->isVisible(sprite);
}

static void sprite_setOpaque(const struct playdate_sprite* spr, LCDSprite* sprite, int flag) {
    spr->setOpaque(sprite, flag);
}

static void sprite_markDirty(const struct playdate_sprite* spr, LCDSprite* sprite) {
    spr->markDirty(sprite);
}

static void sprite_setTag(const struct playdate_sprite* spr, LCDSprite* sprite, uint8_t tag) {
    spr->setTag(sprite, tag);
}

static uint8_t sprite_getTag(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->getTag(sprite);
}

static void sprite_setIgnoresDrawOffset(const struct playdate_sprite* spr, LCDSprite* sprite, int flag) {
    spr->setIgnoresDrawOffset(sprite, flag);
}

static void sprite_getPosition(const struct playdate_sprite* spr, LCDSprite* sprite, float* x, float* y) {
    spr->getPosition(sprite, x, y);
}

// Collision functions
static void sprite_resetCollisionWorld(const struct playdate_sprite* spr) {
    spr->resetCollisionWorld();
}

static void sprite_setCollideRect(const struct playdate_sprite* spr, LCDSprite* sprite, PDRect collideRect) {
    spr->setCollideRect(sprite, collideRect);
}

static PDRect sprite_getCollideRect(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->getCollideRect(sprite);
}

static void sprite_clearCollideRect(const struct playdate_sprite* spr, LCDSprite* sprite) {
    spr->clearCollideRect(sprite);
}

static SpriteCollisionInfo* sprite_checkCollisions(const struct playdate_sprite* spr, LCDSprite* sprite, float goalX, float goalY, float* actualX, float* actualY, int* len) {
    return spr->checkCollisions(sprite, goalX, goalY, actualX, actualY, len);
}

static SpriteCollisionInfo* sprite_moveWithCollisions(const struct playdate_sprite* spr, LCDSprite* sprite, float goalX, float goalY, float* actualX, float* actualY, int* len) {
    return spr->moveWithCollisions(sprite, goalX, goalY, actualX, actualY, len);
}

static LCDSprite** sprite_querySpritesAtPoint(const struct playdate_sprite* spr, float x, float y, int* len) {
    return spr->querySpritesAtPoint(x, y, len);
}

static LCDSprite** sprite_querySpritesInRect(const struct playdate_sprite* spr, float x, float y, float width, float height, int* len) {
    return spr->querySpritesInRect(x, y, width, height, len);
}

static LCDSprite** sprite_querySpritesAlongLine(const struct playdate_sprite* spr, float x1, float y1, float x2, float y2, int* len) {
    return spr->querySpritesAlongLine(x1, y1, x2, y2, len);
}

static SpriteQueryInfo* sprite_querySpriteInfoAlongLine(const struct playdate_sprite* spr, float x1, float y1, float x2, float y2, int* len) {
    return spr->querySpriteInfoAlongLine(x1, y1, x2, y2, len);
}

static LCDSprite** sprite_overlappingSprites(const struct playdate_sprite* spr, LCDSprite* sprite, int* len) {
    return spr->overlappingSprites(sprite, len);
}

static LCDSprite** sprite_allOverlappingSprites(const struct playdate_sprite* spr, int* len) {
    return spr->allOverlappingSprites(len);
}

// v1.7
static void sprite_setStencilPattern(const struct playdate_sprite* spr, LCDSprite* sprite, uint8_t pattern[8]) {
    spr->setStencilPattern(sprite, pattern);
}

static void sprite_clearStencil(const struct playdate_sprite* spr, LCDSprite* sprite) {
    spr->clearStencil(sprite);
}

static void sprite_setUserdata(const struct playdate_sprite* spr, LCDSprite* sprite, void* userdata) {
    spr->setUserdata(sprite, userdata);
}

static void* sprite_getUserdata(const struct playdate_sprite* spr, LCDSprite* sprite) {
    return spr->getUserdata(sprite);
}

// v1.10
static void sprite_setStencilImage(const struct playdate_sprite* spr, LCDSprite* sprite, LCDBitmap* stencil, int tile) {
    spr->setStencilImage(sprite, stencil, tile);
}

// v2.1
static void sprite_setCenter(const struct playdate_sprite* spr, LCDSprite* s, float x, float y) {
    spr->setCenter(s, x, y);
}

static void sprite_getCenter(const struct playdate_sprite* spr, LCDSprite* s, float* x, float* y) {
    spr->getCenter(s, x, y);
}

// v2.7
static void sprite_setTilemap(const struct playdate_sprite* spr, LCDSprite* s, LCDTileMap* tilemap) {
    spr->setTilemap(s, tilemap);
}

static LCDTileMap* sprite_getTilemap(const struct playdate_sprite* spr, LCDSprite* s) {
    return spr->getTilemap(s);
}

// Callback wrappers
extern void goSpriteUpdateCallback(LCDSprite* sprite);
extern void goSpriteDrawCallback(LCDSprite* sprite, PDRect bounds, PDRect drawrect);
extern SpriteCollisionResponseType goSpriteCollisionFilterCallback(LCDSprite* sprite, LCDSprite* other);

static void spriteUpdateCallbackWrapper(LCDSprite* sprite) {
    goSpriteUpdateCallback(sprite);
}

static void spriteDrawCallbackWrapper(LCDSprite* sprite, PDRect bounds, PDRect drawrect) {
    goSpriteDrawCallback(sprite, bounds, drawrect);
}

static SpriteCollisionResponseType spriteCollisionFilterCallbackWrapper(LCDSprite* sprite, LCDSprite* other) {
    return goSpriteCollisionFilterCallback(sprite, other);
}

static void sprite_setUpdateFunction(const struct playdate_sprite* spr, LCDSprite* sprite, int hasCallback) {
    if (hasCallback) {
        spr->setUpdateFunction(sprite, spriteUpdateCallbackWrapper);
    } else {
        spr->setUpdateFunction(sprite, NULL);
    }
}

static void sprite_setDrawFunction(const struct playdate_sprite* spr, LCDSprite* sprite, int hasCallback) {
    if (hasCallback) {
        spr->setDrawFunction(sprite, spriteDrawCallbackWrapper);
    } else {
        spr->setDrawFunction(sprite, NULL);
    }
}

static void sprite_setCollisionResponseFunction(const struct playdate_sprite* spr, LCDSprite* sprite, int hasCallback) {
    if (hasCallback) {
        spr->setCollisionResponseFunction(sprite, spriteCollisionFilterCallbackWrapper);
    } else {
        spr->setCollisionResponseFunction(sprite, NULL);
    }
}
*/
import "C"
import (
	"sync"
	"unsafe"
)

// SpriteCollisionResponseType represents collision response types
type SpriteCollisionResponseType int

const (
	CollisionTypeSlide   SpriteCollisionResponseType = C.kCollisionTypeSlide
	CollisionTypeFreeze  SpriteCollisionResponseType = C.kCollisionTypeFreeze
	CollisionTypeOverlap SpriteCollisionResponseType = C.kCollisionTypeOverlap
	CollisionTypeBounce  SpriteCollisionResponseType = C.kCollisionTypeBounce
)

// PDRect represents a rectangle with float coordinates
type PDRect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// CollisionPoint represents a collision point
type CollisionPoint struct {
	X float32
	Y float32
}

// CollisionVector represents a collision normal vector
type CollisionVector struct {
	X int
	Y int
}

// SpriteCollisionInfo contains collision information
type SpriteCollisionInfo struct {
	Sprite       *LCDSprite                  // The sprite being moved
	Other        *LCDSprite                  // The sprite colliding with
	ResponseType SpriteCollisionResponseType // The result of collision response
	Overlaps     bool                        // True if overlapping at start
	Ti           float32                     // 0-1 position of collision
	Move         CollisionPoint              // Difference between original and actual coordinates
	Normal       CollisionVector             // Collision normal
	Touch        CollisionPoint              // Touch point
	SpriteRect   PDRect                      // Sprite rectangle at touch
	OtherRect    PDRect                      // Other sprite rectangle at touch
}

// SpriteQueryInfo contains sprite query information
type SpriteQueryInfo struct {
	Sprite     *LCDSprite     // The sprite being intersected
	Ti1        float32        // Entry point (0-1)
	Ti2        float32        // Exit point (0-1)
	EntryPoint CollisionPoint // First intersection coordinates
	ExitPoint  CollisionPoint // Second intersection coordinates
}

// LCDSprite wraps a Playdate sprite
type LCDSprite struct {
	ptr *C.LCDSprite
}

// Sprite wraps the playdate_sprite API
type Sprite struct {
	ptr *C.struct_playdate_sprite
}

// Sprite callback storage
var (
	spriteUpdateCallbacks    = make(map[*C.LCDSprite]func(*LCDSprite))
	spriteDrawCallbacks      = make(map[*C.LCDSprite]func(*LCDSprite, PDRect, PDRect))
	spriteCollisionCallbacks = make(map[*C.LCDSprite]func(*LCDSprite, *LCDSprite) SpriteCollisionResponseType)
	spriteMutex              sync.RWMutex
)

//export goSpriteUpdateCallback
func goSpriteUpdateCallback(sprite *C.LCDSprite) {
	spriteMutex.RLock()
	cb, ok := spriteUpdateCallbacks[sprite]
	spriteMutex.RUnlock()

	if ok && cb != nil {
		cb(&LCDSprite{ptr: sprite})
	}
}

//export goSpriteDrawCallback
func goSpriteDrawCallback(sprite *C.LCDSprite, bounds C.PDRect, drawrect C.PDRect) {
	spriteMutex.RLock()
	cb, ok := spriteDrawCallbacks[sprite]
	spriteMutex.RUnlock()

	if ok && cb != nil {
		cb(&LCDSprite{ptr: sprite},
			PDRect{X: float32(bounds.x), Y: float32(bounds.y), Width: float32(bounds.width), Height: float32(bounds.height)},
			PDRect{X: float32(drawrect.x), Y: float32(drawrect.y), Width: float32(drawrect.width), Height: float32(drawrect.height)})
	}
}

//export goSpriteCollisionFilterCallback
func goSpriteCollisionFilterCallback(sprite *C.LCDSprite, other *C.LCDSprite) C.SpriteCollisionResponseType {
	spriteMutex.RLock()
	cb, ok := spriteCollisionCallbacks[sprite]
	spriteMutex.RUnlock()

	if ok && cb != nil {
		return C.SpriteCollisionResponseType(cb(&LCDSprite{ptr: sprite}, &LCDSprite{ptr: other}))
	}
	return C.SpriteCollisionResponseType(CollisionTypeSlide)
}

func newSprite(ptr *C.struct_playdate_sprite) *Sprite {
	return &Sprite{ptr: ptr}
}

// SetAlwaysRedraw sets whether sprites should always be redrawn
func (s *Sprite) SetAlwaysRedraw(flag bool) {
	f := 0
	if flag {
		f = 1
	}
	C.sprite_setAlwaysRedraw(s.ptr, C.int(f))
}

// AddDirtyRect adds a dirty rectangle
func (s *Sprite) AddDirtyRect(rect LCDRect) {
	crect := C.LCDRect{
		left:   C.int(rect.Left),
		right:  C.int(rect.Right),
		top:    C.int(rect.Top),
		bottom: C.int(rect.Bottom),
	}
	C.sprite_addDirtyRect(s.ptr, crect)
}

// DrawSprites draws all sprites
func (s *Sprite) DrawSprites() {
	C.sprite_drawSprites(s.ptr)
}

// UpdateAndDrawSprites updates and draws all sprites
func (s *Sprite) UpdateAndDrawSprites() {
	C.sprite_updateAndDrawSprites(s.ptr)
}

// NewSprite creates a new sprite
func (s *Sprite) NewSprite() *LCDSprite {
	ptr := C.sprite_newSprite(s.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDSprite{ptr: ptr}
}

// FreeSprite frees a sprite
func (s *Sprite) FreeSprite(sprite *LCDSprite) {
	if sprite != nil && sprite.ptr != nil {
		spriteMutex.Lock()
		delete(spriteUpdateCallbacks, sprite.ptr)
		delete(spriteDrawCallbacks, sprite.ptr)
		delete(spriteCollisionCallbacks, sprite.ptr)
		spriteMutex.Unlock()

		C.sprite_freeSprite(s.ptr, sprite.ptr)
		sprite.ptr = nil
	}
}

// Copy creates a copy of a sprite
func (s *Sprite) Copy(sprite *LCDSprite) *LCDSprite {
	if sprite == nil {
		return nil
	}
	ptr := C.sprite_copy(s.ptr, sprite.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDSprite{ptr: ptr}
}

// AddSprite adds a sprite to the display list
func (s *Sprite) AddSprite(sprite *LCDSprite) {
	if sprite != nil {
		C.sprite_addSprite(s.ptr, sprite.ptr)
	}
}

// RemoveSprite removes a sprite from the display list
func (s *Sprite) RemoveSprite(sprite *LCDSprite) {
	if sprite != nil {
		C.sprite_removeSprite(s.ptr, sprite.ptr)
	}
}

// RemoveSprites removes multiple sprites
func (s *Sprite) RemoveSprites(sprites []*LCDSprite) {
	if len(sprites) == 0 {
		return
	}
	ptrs := make([]*C.LCDSprite, len(sprites))
	for i, spr := range sprites {
		if spr != nil {
			ptrs[i] = spr.ptr
		}
	}
	C.sprite_removeSprites(s.ptr, &ptrs[0], C.int(len(sprites)))
}

// RemoveAllSprites removes all sprites
func (s *Sprite) RemoveAllSprites() {
	C.sprite_removeAllSprites(s.ptr)
}

// GetSpriteCount returns the number of sprites
func (s *Sprite) GetSpriteCount() int {
	return int(C.sprite_getSpriteCount(s.ptr))
}

// SetBounds sets the bounds of a sprite
func (s *Sprite) SetBounds(sprite *LCDSprite, bounds PDRect) {
	if sprite == nil {
		return
	}
	cbounds := C.PDRect{
		x:      C.float(bounds.X),
		y:      C.float(bounds.Y),
		width:  C.float(bounds.Width),
		height: C.float(bounds.Height),
	}
	C.sprite_setBounds(s.ptr, sprite.ptr, cbounds)
}

// GetBounds returns the bounds of a sprite
func (s *Sprite) GetBounds(sprite *LCDSprite) PDRect {
	if sprite == nil {
		return PDRect{}
	}
	bounds := C.sprite_getBounds(s.ptr, sprite.ptr)
	return PDRect{
		X:      float32(bounds.x),
		Y:      float32(bounds.y),
		Width:  float32(bounds.width),
		Height: float32(bounds.height),
	}
}

// MoveTo moves a sprite to a position
func (s *Sprite) MoveTo(sprite *LCDSprite, x, y float32) {
	if sprite != nil {
		C.sprite_moveTo(s.ptr, sprite.ptr, C.float(x), C.float(y))
	}
}

// MoveBy moves a sprite by an offset
func (s *Sprite) MoveBy(sprite *LCDSprite, dx, dy float32) {
	if sprite != nil {
		C.sprite_moveBy(s.ptr, sprite.ptr, C.float(dx), C.float(dy))
	}
}

// SetImage sets the image for a sprite
func (s *Sprite) SetImage(sprite *LCDSprite, image *LCDBitmap, flip LCDBitmapFlip) {
	if sprite == nil {
		return
	}
	var img *C.LCDBitmap
	if image != nil {
		img = image.ptr
	}
	C.sprite_setImage(s.ptr, sprite.ptr, img, C.LCDBitmapFlip(flip))
}

// GetImage returns the image of a sprite
func (s *Sprite) GetImage(sprite *LCDSprite) *LCDBitmap {
	if sprite == nil {
		return nil
	}
	ptr := C.sprite_getImage(s.ptr, sprite.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// SetSize sets the size of a sprite
func (s *Sprite) SetSize(sprite *LCDSprite, width, height float32) {
	if sprite != nil {
		C.sprite_setSize(s.ptr, sprite.ptr, C.float(width), C.float(height))
	}
}

// SetZIndex sets the Z index of a sprite
func (s *Sprite) SetZIndex(sprite *LCDSprite, zIndex int16) {
	if sprite != nil {
		C.sprite_setZIndex(s.ptr, sprite.ptr, C.int16_t(zIndex))
	}
}

// GetZIndex returns the Z index of a sprite
func (s *Sprite) GetZIndex(sprite *LCDSprite) int16 {
	if sprite == nil {
		return 0
	}
	return int16(C.sprite_getZIndex(s.ptr, sprite.ptr))
}

// SetDrawMode sets the draw mode of a sprite
func (s *Sprite) SetDrawMode(sprite *LCDSprite, mode LCDBitmapDrawMode) {
	if sprite != nil {
		C.sprite_setDrawMode(s.ptr, sprite.ptr, C.LCDBitmapDrawMode(mode))
	}
}

// SetImageFlip sets the image flip of a sprite
func (s *Sprite) SetImageFlip(sprite *LCDSprite, flip LCDBitmapFlip) {
	if sprite != nil {
		C.sprite_setImageFlip(s.ptr, sprite.ptr, C.LCDBitmapFlip(flip))
	}
}

// GetImageFlip returns the image flip of a sprite
func (s *Sprite) GetImageFlip(sprite *LCDSprite) LCDBitmapFlip {
	if sprite == nil {
		return BitmapUnflipped
	}
	return LCDBitmapFlip(C.sprite_getImageFlip(s.ptr, sprite.ptr))
}

// SetStencil sets the stencil of a sprite (deprecated)
func (s *Sprite) SetStencil(sprite *LCDSprite, stencil *LCDBitmap) {
	if sprite == nil {
		return
	}
	var st *C.LCDBitmap
	if stencil != nil {
		st = stencil.ptr
	}
	C.sprite_setStencil(s.ptr, sprite.ptr, st)
}

// SetClipRect sets the clip rect of a sprite
func (s *Sprite) SetClipRect(sprite *LCDSprite, clipRect LCDRect) {
	if sprite == nil {
		return
	}
	crect := C.LCDRect{
		left:   C.int(clipRect.Left),
		right:  C.int(clipRect.Right),
		top:    C.int(clipRect.Top),
		bottom: C.int(clipRect.Bottom),
	}
	C.sprite_setClipRect(s.ptr, sprite.ptr, crect)
}

// ClearClipRect clears the clip rect of a sprite
func (s *Sprite) ClearClipRect(sprite *LCDSprite) {
	if sprite != nil {
		C.sprite_clearClipRect(s.ptr, sprite.ptr)
	}
}

// SetClipRectsInRange sets clip rects for sprites in a Z range
func (s *Sprite) SetClipRectsInRange(clipRect LCDRect, startZ, endZ int) {
	crect := C.LCDRect{
		left:   C.int(clipRect.Left),
		right:  C.int(clipRect.Right),
		top:    C.int(clipRect.Top),
		bottom: C.int(clipRect.Bottom),
	}
	C.sprite_setClipRectsInRange(s.ptr, crect, C.int(startZ), C.int(endZ))
}

// ClearClipRectsInRange clears clip rects for sprites in a Z range
func (s *Sprite) ClearClipRectsInRange(startZ, endZ int) {
	C.sprite_clearClipRectsInRange(s.ptr, C.int(startZ), C.int(endZ))
}

// SetUpdatesEnabled enables or disables updates for a sprite
func (s *Sprite) SetUpdatesEnabled(sprite *LCDSprite, enabled bool) {
	if sprite == nil {
		return
	}
	f := 0
	if enabled {
		f = 1
	}
	C.sprite_setUpdatesEnabled(s.ptr, sprite.ptr, C.int(f))
}

// UpdatesEnabled returns whether updates are enabled for a sprite
func (s *Sprite) UpdatesEnabled(sprite *LCDSprite) bool {
	if sprite == nil {
		return false
	}
	return C.sprite_updatesEnabled(s.ptr, sprite.ptr) != 0
}

// SetCollisionsEnabled enables or disables collisions for a sprite
func (s *Sprite) SetCollisionsEnabled(sprite *LCDSprite, enabled bool) {
	if sprite == nil {
		return
	}
	f := 0
	if enabled {
		f = 1
	}
	C.sprite_setCollisionsEnabled(s.ptr, sprite.ptr, C.int(f))
}

// CollisionsEnabled returns whether collisions are enabled for a sprite
func (s *Sprite) CollisionsEnabled(sprite *LCDSprite) bool {
	if sprite == nil {
		return false
	}
	return C.sprite_collisionsEnabled(s.ptr, sprite.ptr) != 0
}

// SetVisible sets the visibility of a sprite
func (s *Sprite) SetVisible(sprite *LCDSprite, visible bool) {
	if sprite == nil {
		return
	}
	f := 0
	if visible {
		f = 1
	}
	C.sprite_setVisible(s.ptr, sprite.ptr, C.int(f))
}

// IsVisible returns whether a sprite is visible
func (s *Sprite) IsVisible(sprite *LCDSprite) bool {
	if sprite == nil {
		return false
	}
	return C.sprite_isVisible(s.ptr, sprite.ptr) != 0
}

// SetOpaque sets whether a sprite is opaque
func (s *Sprite) SetOpaque(sprite *LCDSprite, opaque bool) {
	if sprite == nil {
		return
	}
	f := 0
	if opaque {
		f = 1
	}
	C.sprite_setOpaque(s.ptr, sprite.ptr, C.int(f))
}

// MarkDirty marks a sprite as dirty
func (s *Sprite) MarkDirty(sprite *LCDSprite) {
	if sprite != nil {
		C.sprite_markDirty(s.ptr, sprite.ptr)
	}
}

// SetTag sets the tag of a sprite
func (s *Sprite) SetTag(sprite *LCDSprite, tag uint8) {
	if sprite != nil {
		C.sprite_setTag(s.ptr, sprite.ptr, C.uint8_t(tag))
	}
}

// GetTag returns the tag of a sprite
func (s *Sprite) GetTag(sprite *LCDSprite) uint8 {
	if sprite == nil {
		return 0
	}
	return uint8(C.sprite_getTag(s.ptr, sprite.ptr))
}

// SetIgnoresDrawOffset sets whether a sprite ignores the draw offset
func (s *Sprite) SetIgnoresDrawOffset(sprite *LCDSprite, ignores bool) {
	if sprite == nil {
		return
	}
	f := 0
	if ignores {
		f = 1
	}
	C.sprite_setIgnoresDrawOffset(s.ptr, sprite.ptr, C.int(f))
}

// SetUpdateFunction sets the update function for a sprite
func (s *Sprite) SetUpdateFunction(sprite *LCDSprite, callback func(*LCDSprite)) {
	if sprite == nil {
		return
	}
	spriteMutex.Lock()
	if callback != nil {
		spriteUpdateCallbacks[sprite.ptr] = callback
	} else {
		delete(spriteUpdateCallbacks, sprite.ptr)
	}
	spriteMutex.Unlock()

	hasCallback := 0
	if callback != nil {
		hasCallback = 1
	}
	C.sprite_setUpdateFunction(s.ptr, sprite.ptr, C.int(hasCallback))
}

// SetDrawFunction sets the draw function for a sprite
func (s *Sprite) SetDrawFunction(sprite *LCDSprite, callback func(*LCDSprite, PDRect, PDRect)) {
	if sprite == nil {
		return
	}
	spriteMutex.Lock()
	if callback != nil {
		spriteDrawCallbacks[sprite.ptr] = callback
	} else {
		delete(spriteDrawCallbacks, sprite.ptr)
	}
	spriteMutex.Unlock()

	hasCallback := 0
	if callback != nil {
		hasCallback = 1
	}
	C.sprite_setDrawFunction(s.ptr, sprite.ptr, C.int(hasCallback))
}

// GetPosition returns the position of a sprite
func (s *Sprite) GetPosition(sprite *LCDSprite) (x, y float32) {
	if sprite == nil {
		return 0, 0
	}
	var cx, cy C.float
	C.sprite_getPosition(s.ptr, sprite.ptr, &cx, &cy)
	return float32(cx), float32(cy)
}

// ResetCollisionWorld resets the collision world
func (s *Sprite) ResetCollisionWorld() {
	C.sprite_resetCollisionWorld(s.ptr)
}

// SetCollideRect sets the collision rect of a sprite
func (s *Sprite) SetCollideRect(sprite *LCDSprite, collideRect PDRect) {
	if sprite == nil {
		return
	}
	crect := C.PDRect{
		x:      C.float(collideRect.X),
		y:      C.float(collideRect.Y),
		width:  C.float(collideRect.Width),
		height: C.float(collideRect.Height),
	}
	C.sprite_setCollideRect(s.ptr, sprite.ptr, crect)
}

// GetCollideRect returns the collision rect of a sprite
func (s *Sprite) GetCollideRect(sprite *LCDSprite) PDRect {
	if sprite == nil {
		return PDRect{}
	}
	rect := C.sprite_getCollideRect(s.ptr, sprite.ptr)
	return PDRect{
		X:      float32(rect.x),
		Y:      float32(rect.y),
		Width:  float32(rect.width),
		Height: float32(rect.height),
	}
}

// ClearCollideRect clears the collision rect of a sprite
func (s *Sprite) ClearCollideRect(sprite *LCDSprite) {
	if sprite != nil {
		C.sprite_clearCollideRect(s.ptr, sprite.ptr)
	}
}

// SetCollisionResponseFunction sets the collision response function
func (s *Sprite) SetCollisionResponseFunction(sprite *LCDSprite, callback func(*LCDSprite, *LCDSprite) SpriteCollisionResponseType) {
	if sprite == nil {
		return
	}
	spriteMutex.Lock()
	if callback != nil {
		spriteCollisionCallbacks[sprite.ptr] = callback
	} else {
		delete(spriteCollisionCallbacks, sprite.ptr)
	}
	spriteMutex.Unlock()

	hasCallback := 0
	if callback != nil {
		hasCallback = 1
	}
	C.sprite_setCollisionResponseFunction(s.ptr, sprite.ptr, C.int(hasCallback))
}

// CheckCollisions checks for collisions when moving a sprite
func (s *Sprite) CheckCollisions(sprite *LCDSprite, goalX, goalY float32) ([]SpriteCollisionInfo, float32, float32) {
	if sprite == nil {
		return nil, goalX, goalY
	}
	var actualX, actualY C.float
	var length C.int
	result := C.sprite_checkCollisions(s.ptr, sprite.ptr, C.float(goalX), C.float(goalY), &actualX, &actualY, &length)
	if result == nil || length == 0 {
		return nil, float32(actualX), float32(actualY)
	}
	defer C.free(unsafe.Pointer(result))

	collisions := make([]SpriteCollisionInfo, int(length))
	cinfo := (*[1 << 20]C.SpriteCollisionInfo)(unsafe.Pointer(result))[:length:length]
	for i := 0; i < int(length); i++ {
		collisions[i] = convertCollisionInfo(&cinfo[i])
	}
	return collisions, float32(actualX), float32(actualY)
}

// MoveWithCollisions moves a sprite with collision detection
func (s *Sprite) MoveWithCollisions(sprite *LCDSprite, goalX, goalY float32) ([]SpriteCollisionInfo, float32, float32) {
	if sprite == nil {
		return nil, goalX, goalY
	}
	var actualX, actualY C.float
	var length C.int
	result := C.sprite_moveWithCollisions(s.ptr, sprite.ptr, C.float(goalX), C.float(goalY), &actualX, &actualY, &length)
	if result == nil || length == 0 {
		return nil, float32(actualX), float32(actualY)
	}
	defer C.free(unsafe.Pointer(result))

	collisions := make([]SpriteCollisionInfo, int(length))
	cinfo := (*[1 << 20]C.SpriteCollisionInfo)(unsafe.Pointer(result))[:length:length]
	for i := 0; i < int(length); i++ {
		collisions[i] = convertCollisionInfo(&cinfo[i])
	}
	return collisions, float32(actualX), float32(actualY)
}

func convertCollisionInfo(info *C.SpriteCollisionInfo) SpriteCollisionInfo {
	return SpriteCollisionInfo{
		Sprite:       &LCDSprite{ptr: info.sprite},
		Other:        &LCDSprite{ptr: info.other},
		ResponseType: SpriteCollisionResponseType(info.responseType),
		Overlaps:     info.overlaps != 0,
		Ti:           float32(info.ti),
		Move:         CollisionPoint{X: float32(info.move.x), Y: float32(info.move.y)},
		Normal:       CollisionVector{X: int(info.normal.x), Y: int(info.normal.y)},
		Touch:        CollisionPoint{X: float32(info.touch.x), Y: float32(info.touch.y)},
		SpriteRect: PDRect{
			X:      float32(info.spriteRect.x),
			Y:      float32(info.spriteRect.y),
			Width:  float32(info.spriteRect.width),
			Height: float32(info.spriteRect.height),
		},
		OtherRect: PDRect{
			X:      float32(info.otherRect.x),
			Y:      float32(info.otherRect.y),
			Width:  float32(info.otherRect.width),
			Height: float32(info.otherRect.height),
		},
	}
}

// QuerySpritesAtPoint queries sprites at a point
func (s *Sprite) QuerySpritesAtPoint(x, y float32) []*LCDSprite {
	var length C.int
	result := C.sprite_querySpritesAtPoint(s.ptr, C.float(x), C.float(y), &length)
	if result == nil || length == 0 {
		return nil
	}
	defer C.free(unsafe.Pointer(result))

	sprites := make([]*LCDSprite, int(length))
	csprites := (*[1 << 20]*C.LCDSprite)(unsafe.Pointer(result))[:length:length]
	for i := 0; i < int(length); i++ {
		sprites[i] = &LCDSprite{ptr: csprites[i]}
	}
	return sprites
}

// QuerySpritesInRect queries sprites in a rectangle
func (s *Sprite) QuerySpritesInRect(x, y, width, height float32) []*LCDSprite {
	var length C.int
	result := C.sprite_querySpritesInRect(s.ptr, C.float(x), C.float(y), C.float(width), C.float(height), &length)
	if result == nil || length == 0 {
		return nil
	}
	defer C.free(unsafe.Pointer(result))

	sprites := make([]*LCDSprite, int(length))
	csprites := (*[1 << 20]*C.LCDSprite)(unsafe.Pointer(result))[:length:length]
	for i := 0; i < int(length); i++ {
		sprites[i] = &LCDSprite{ptr: csprites[i]}
	}
	return sprites
}

// QuerySpritesAlongLine queries sprites along a line
func (s *Sprite) QuerySpritesAlongLine(x1, y1, x2, y2 float32) []*LCDSprite {
	var length C.int
	result := C.sprite_querySpritesAlongLine(s.ptr, C.float(x1), C.float(y1), C.float(x2), C.float(y2), &length)
	if result == nil || length == 0 {
		return nil
	}
	defer C.free(unsafe.Pointer(result))

	sprites := make([]*LCDSprite, int(length))
	csprites := (*[1 << 20]*C.LCDSprite)(unsafe.Pointer(result))[:length:length]
	for i := 0; i < int(length); i++ {
		sprites[i] = &LCDSprite{ptr: csprites[i]}
	}
	return sprites
}

// QuerySpriteInfoAlongLine queries sprite info along a line
func (s *Sprite) QuerySpriteInfoAlongLine(x1, y1, x2, y2 float32) []SpriteQueryInfo {
	var length C.int
	result := C.sprite_querySpriteInfoAlongLine(s.ptr, C.float(x1), C.float(y1), C.float(x2), C.float(y2), &length)
	if result == nil || length == 0 {
		return nil
	}
	defer C.free(unsafe.Pointer(result))

	infos := make([]SpriteQueryInfo, int(length))
	cinfos := (*[1 << 20]C.SpriteQueryInfo)(unsafe.Pointer(result))[:length:length]
	for i := 0; i < int(length); i++ {
		infos[i] = SpriteQueryInfo{
			Sprite:     &LCDSprite{ptr: cinfos[i].sprite},
			Ti1:        float32(cinfos[i].ti1),
			Ti2:        float32(cinfos[i].ti2),
			EntryPoint: CollisionPoint{X: float32(cinfos[i].entryPoint.x), Y: float32(cinfos[i].entryPoint.y)},
			ExitPoint:  CollisionPoint{X: float32(cinfos[i].exitPoint.x), Y: float32(cinfos[i].exitPoint.y)},
		}
	}
	return infos
}

// OverlappingSprites returns sprites overlapping with a sprite
func (s *Sprite) OverlappingSprites(sprite *LCDSprite) []*LCDSprite {
	if sprite == nil {
		return nil
	}
	var length C.int
	result := C.sprite_overlappingSprites(s.ptr, sprite.ptr, &length)
	if result == nil || length == 0 {
		return nil
	}
	defer C.free(unsafe.Pointer(result))

	sprites := make([]*LCDSprite, int(length))
	csprites := (*[1 << 20]*C.LCDSprite)(unsafe.Pointer(result))[:length:length]
	for i := 0; i < int(length); i++ {
		sprites[i] = &LCDSprite{ptr: csprites[i]}
	}
	return sprites
}

// AllOverlappingSprites returns all overlapping sprites
func (s *Sprite) AllOverlappingSprites() []*LCDSprite {
	var length C.int
	result := C.sprite_allOverlappingSprites(s.ptr, &length)
	if result == nil || length == 0 {
		return nil
	}
	defer C.free(unsafe.Pointer(result))

	sprites := make([]*LCDSprite, int(length))
	csprites := (*[1 << 20]*C.LCDSprite)(unsafe.Pointer(result))[:length:length]
	for i := 0; i < int(length); i++ {
		sprites[i] = &LCDSprite{ptr: csprites[i]}
	}
	return sprites
}

// SetStencilPattern sets the stencil pattern for a sprite
func (s *Sprite) SetStencilPattern(sprite *LCDSprite, pattern [8]uint8) {
	if sprite != nil {
		C.sprite_setStencilPattern(s.ptr, sprite.ptr, (*C.uint8_t)(&pattern[0]))
	}
}

// ClearStencil clears the stencil for a sprite
func (s *Sprite) ClearStencil(sprite *LCDSprite) {
	if sprite != nil {
		C.sprite_clearStencil(s.ptr, sprite.ptr)
	}
}

// SetUserdata sets user data for a sprite
func (s *Sprite) SetUserdata(sprite *LCDSprite, userdata unsafe.Pointer) {
	if sprite != nil {
		C.sprite_setUserdata(s.ptr, sprite.ptr, userdata)
	}
}

// GetUserdata returns user data for a sprite
func (s *Sprite) GetUserdata(sprite *LCDSprite) unsafe.Pointer {
	if sprite == nil {
		return nil
	}
	return C.sprite_getUserdata(s.ptr, sprite.ptr)
}

// SetStencilImage sets the stencil image for a sprite
func (s *Sprite) SetStencilImage(sprite *LCDSprite, stencil *LCDBitmap, tile bool) {
	if sprite == nil {
		return
	}
	var st *C.LCDBitmap
	if stencil != nil {
		st = stencil.ptr
	}
	t := 0
	if tile {
		t = 1
	}
	C.sprite_setStencilImage(s.ptr, sprite.ptr, st, C.int(t))
}

// SetCenter sets the center of a sprite
func (s *Sprite) SetCenter(sprite *LCDSprite, x, y float32) {
	if sprite != nil {
		C.sprite_setCenter(s.ptr, sprite.ptr, C.float(x), C.float(y))
	}
}

// GetCenter returns the center of a sprite
func (s *Sprite) GetCenter(sprite *LCDSprite) (x, y float32) {
	if sprite == nil {
		return 0.5, 0.5
	}
	var cx, cy C.float
	C.sprite_getCenter(s.ptr, sprite.ptr, &cx, &cy)
	return float32(cx), float32(cy)
}

// SetTilemap sets the tilemap for a sprite
func (s *Sprite) SetTilemap(sprite *LCDSprite, tilemap *LCDTileMap) {
	if sprite == nil {
		return
	}
	var tm *C.LCDTileMap
	if tilemap != nil {
		tm = tilemap.ptr
	}
	C.sprite_setTilemap(s.ptr, sprite.ptr, tm)
}

// GetTilemap returns the tilemap for a sprite
func (s *Sprite) GetTilemap(sprite *LCDSprite) *LCDTileMap {
	if sprite == nil {
		return nil
	}
	ptr := C.sprite_getTilemap(s.ptr, sprite.ptr)
	if ptr == nil {
		return nil
	}
	return &LCDTileMap{ptr: ptr}
}
