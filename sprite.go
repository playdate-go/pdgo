// pdgo Sprite API - unified CGO implementation

package pdgo

/*
#include <stdint.h>

// Sprite API
void* pd_sprite_new(void);
void pd_sprite_free(void* sprite);
void pd_sprite_add(void* sprite);
void pd_sprite_remove(void* sprite);
void pd_sprite_removeAll(void);
int pd_sprite_getCount(void);
void pd_sprite_setImage(void* sprite, void* image, int flip);
void* pd_sprite_getImage(void* sprite);
void pd_sprite_setBounds(void* sprite, float x, float y, float w, float h);
void pd_sprite_getBounds(void* sprite, float* x, float* y, float* w, float* h);
void pd_sprite_moveTo(void* sprite, float x, float y);
void pd_sprite_moveBy(void* sprite, float dx, float dy);
void pd_sprite_getPosition(void* sprite, float* x, float* y);
void pd_sprite_setZIndex(void* sprite, int16_t z);
int16_t pd_sprite_getZIndex(void* sprite);
void pd_sprite_setTag(void* sprite, uint8_t tag);
uint8_t pd_sprite_getTag(void* sprite);
void pd_sprite_setVisible(void* sprite, int visible);
int pd_sprite_isVisible(void* sprite);
void pd_sprite_setOpaque(void* sprite, int opaque);
void pd_sprite_setDrawMode(void* sprite, int mode);
void pd_sprite_setImageFlip(void* sprite, int flip);
int pd_sprite_getImageFlip(void* sprite);
void pd_sprite_setUpdatesEnabled(void* sprite, int enabled);
void pd_sprite_markDirty(void* sprite);
void pd_sprite_drawSprites(void);
void pd_sprite_updateAndDrawSprites(void);
void pd_sprite_setAlwaysRedraw(int always);
void pd_sprite_setCollideRect(void* sprite, float x, float y, float w, float h);
void pd_sprite_getCollideRect(void* sprite, float* x, float* y, float* w, float* h);
void pd_sprite_clearCollideRect(void* sprite);
void pd_sprite_setCollisionsEnabled(void* sprite, int enabled);
void pd_sprite_resetCollisionWorld(void);

// Sprite callbacks
void pd_sprite_setUpdateFunction(void* sprite);
void pd_sprite_setDrawFunction(void* sprite);
void pd_sprite_setCollisionResponseFunction(void* sprite);
*/
import "C"
import "unsafe"

// Sprite represents a sprite object
type Sprite struct {
	ptr unsafe.Pointer
}

// LCDSprite is an alias for Sprite (for API compatibility)
type LCDSprite = Sprite

// SpriteAPI provides access to sprite functions
type SpriteAPI struct{}

func newSprite() *SpriteAPI {
	return &SpriteAPI{}
}

// ============== Sprite Creation ==============

// NewSprite creates a new sprite
func (s *SpriteAPI) NewSprite() *LCDSprite {
	ptr := C.pd_sprite_new()
	if ptr != nil {
		return &LCDSprite{ptr: ptr}
	}
	return nil
}

// FreeSprite frees a sprite
func (s *SpriteAPI) FreeSprite(sprite *LCDSprite) {
	if sprite != nil && sprite.ptr != nil {
		unregisterSpriteCallbacks(sprite.ptr)
		C.pd_sprite_free(sprite.ptr)
		sprite.ptr = nil
	}
}

// ============== Display List ==============

// AddSprite adds a sprite to the display list
func (s *SpriteAPI) AddSprite(sprite *LCDSprite) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_add(sprite.ptr)
	}
}

// RemoveSprite removes a sprite from the display list
func (s *SpriteAPI) RemoveSprite(sprite *LCDSprite) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_remove(sprite.ptr)
	}
}

// RemoveAllSprites removes all sprites from display
func (s *SpriteAPI) RemoveAllSprites() {
	C.pd_sprite_removeAll()
}

// GetSpriteCount returns number of sprites
func (s *SpriteAPI) GetSpriteCount() int {
	return int(C.pd_sprite_getCount())
}

// ============== Image ==============

// SetImage sets the sprite's image
func (s *SpriteAPI) SetImage(sprite *LCDSprite, image *LCDBitmap, flip LCDBitmapFlip) {
	if sprite != nil && sprite.ptr != nil {
		var imgPtr unsafe.Pointer
		if image != nil {
			imgPtr = image.ptr
		}
		C.pd_sprite_setImage(sprite.ptr, imgPtr, C.int(flip))
	}
}

// GetImage returns the sprite's image
func (s *SpriteAPI) GetImage(sprite *LCDSprite) *LCDBitmap {
	if sprite != nil && sprite.ptr != nil {
		ptr := C.pd_sprite_getImage(sprite.ptr)
		if ptr != nil {
			return &LCDBitmap{ptr: ptr}
		}
	}
	return nil
}

// SetImageFlip sets image flip
func (s *SpriteAPI) SetImageFlip(sprite *LCDSprite, flip LCDBitmapFlip) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_setImageFlip(sprite.ptr, C.int(flip))
	}
}

// GetImageFlip returns image flip
func (s *SpriteAPI) GetImageFlip(sprite *LCDSprite) LCDBitmapFlip {
	if sprite != nil && sprite.ptr != nil {
		return LCDBitmapFlip(C.pd_sprite_getImageFlip(sprite.ptr))
	}
	return BitmapUnflipped
}

// ============== Position & Bounds ==============

// SetBounds sets the sprite's bounds
func (s *SpriteAPI) SetBounds(sprite *LCDSprite, bounds PDRect) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_setBounds(sprite.ptr, C.float(bounds.X), C.float(bounds.Y), C.float(bounds.Width), C.float(bounds.Height))
	}
}

// GetBounds returns the sprite's bounds
func (s *SpriteAPI) GetBounds(sprite *LCDSprite) PDRect {
	if sprite != nil && sprite.ptr != nil {
		var x, y, w, h C.float
		C.pd_sprite_getBounds(sprite.ptr, &x, &y, &w, &h)
		return PDRect{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}
	}
	return PDRect{}
}

// MoveTo moves the sprite to position
func (s *SpriteAPI) MoveTo(sprite *LCDSprite, x, y float32) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_moveTo(sprite.ptr, C.float(x), C.float(y))
	}
}

// MoveBy moves the sprite by delta
func (s *SpriteAPI) MoveBy(sprite *LCDSprite, dx, dy float32) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_moveBy(sprite.ptr, C.float(dx), C.float(dy))
	}
}

// GetPosition returns sprite position
func (s *SpriteAPI) GetPosition(sprite *LCDSprite) (x, y float32) {
	if sprite != nil && sprite.ptr != nil {
		var cx, cy C.float
		C.pd_sprite_getPosition(sprite.ptr, &cx, &cy)
		return float32(cx), float32(cy)
	}
	return 0, 0
}

// ============== Z-Index & Tags ==============

// SetZIndex sets sprite z-index
func (s *SpriteAPI) SetZIndex(sprite *LCDSprite, z int16) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_setZIndex(sprite.ptr, C.int16_t(z))
	}
}

// GetZIndex returns sprite z-index
func (s *SpriteAPI) GetZIndex(sprite *LCDSprite) int16 {
	if sprite != nil && sprite.ptr != nil {
		return int16(C.pd_sprite_getZIndex(sprite.ptr))
	}
	return 0
}

// SetTag sets sprite tag
func (s *SpriteAPI) SetTag(sprite *LCDSprite, tag uint8) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_setTag(sprite.ptr, C.uint8_t(tag))
	}
}

// GetTag returns sprite tag
func (s *SpriteAPI) GetTag(sprite *LCDSprite) uint8 {
	if sprite != nil && sprite.ptr != nil {
		return uint8(C.pd_sprite_getTag(sprite.ptr))
	}
	return 0
}

// ============== Callbacks ==============

// SetUpdateFunction sets the sprite's update callback
func (s *SpriteAPI) SetUpdateFunction(sprite *LCDSprite, callback func(*LCDSprite)) {
	if sprite == nil || sprite.ptr == nil {
		return
	}
	registerSpriteUpdateCallback(sprite.ptr, callback)
	C.pd_sprite_setUpdateFunction(sprite.ptr)
}

// SetDrawFunction sets the sprite's draw callback
func (s *SpriteAPI) SetDrawFunction(sprite *LCDSprite, callback func(*LCDSprite, PDRect, PDRect)) {
	if sprite == nil || sprite.ptr == nil {
		return
	}
	registerSpriteDrawCallback(sprite.ptr, callback)
	C.pd_sprite_setDrawFunction(sprite.ptr)
}

// SetCollisionResponseFunction sets the sprite's collision response callback
func (s *SpriteAPI) SetCollisionResponseFunction(sprite *LCDSprite, callback func(*LCDSprite, *LCDSprite) SpriteCollisionResponseType) {
	if sprite == nil || sprite.ptr == nil {
		return
	}
	registerSpriteCollisionCallback(sprite.ptr, callback)
	C.pd_sprite_setCollisionResponseFunction(sprite.ptr)
}

// ============== Visibility & Drawing ==============

// SetVisible sets sprite visibility
func (s *SpriteAPI) SetVisible(sprite *LCDSprite, visible bool) {
	if sprite != nil && sprite.ptr != nil {
		var v C.int
		if visible {
			v = 1
		}
		C.pd_sprite_setVisible(sprite.ptr, v)
	}
}

// IsVisible returns true if sprite is visible
func (s *SpriteAPI) IsVisible(sprite *LCDSprite) bool {
	if sprite != nil && sprite.ptr != nil {
		return C.pd_sprite_isVisible(sprite.ptr) != 0
	}
	return false
}

// SetOpaque sets whether sprite is opaque
func (s *SpriteAPI) SetOpaque(sprite *LCDSprite, opaque bool) {
	if sprite != nil && sprite.ptr != nil {
		var o C.int
		if opaque {
			o = 1
		}
		C.pd_sprite_setOpaque(sprite.ptr, o)
	}
}

// SetDrawMode sets sprite draw mode
func (s *SpriteAPI) SetDrawMode(sprite *LCDSprite, mode LCDBitmapDrawMode) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_setDrawMode(sprite.ptr, C.int(mode))
	}
}

// SetUpdatesEnabled enables/disables sprite updates
func (s *SpriteAPI) SetUpdatesEnabled(sprite *LCDSprite, enabled bool) {
	if sprite != nil && sprite.ptr != nil {
		var e C.int
		if enabled {
			e = 1
		}
		C.pd_sprite_setUpdatesEnabled(sprite.ptr, e)
	}
}

// MarkDirty marks the sprite as needing redraw
func (s *SpriteAPI) MarkDirty(sprite *LCDSprite) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_markDirty(sprite.ptr)
	}
}

// ============== Collision ==============

// SetCollideRect sets collision rectangle
func (s *SpriteAPI) SetCollideRect(sprite *LCDSprite, collideRect PDRect) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_setCollideRect(sprite.ptr, C.float(collideRect.X), C.float(collideRect.Y), C.float(collideRect.Width), C.float(collideRect.Height))
	}
}

// GetCollideRect returns collision rectangle
func (s *SpriteAPI) GetCollideRect(sprite *LCDSprite) PDRect {
	if sprite != nil && sprite.ptr != nil {
		var x, y, w, h C.float
		C.pd_sprite_getCollideRect(sprite.ptr, &x, &y, &w, &h)
		return PDRect{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}
	}
	return PDRect{}
}

// ClearCollideRect clears collision rectangle
func (s *SpriteAPI) ClearCollideRect(sprite *LCDSprite) {
	if sprite != nil && sprite.ptr != nil {
		C.pd_sprite_clearCollideRect(sprite.ptr)
	}
}

// SetCollisionsEnabled enables/disables collisions
func (s *SpriteAPI) SetCollisionsEnabled(sprite *LCDSprite, enabled bool) {
	if sprite != nil && sprite.ptr != nil {
		var e C.int
		if enabled {
			e = 1
		}
		C.pd_sprite_setCollisionsEnabled(sprite.ptr, e)
	}
}

// MoveWithCollisions moves sprite with collision detection
func (s *SpriteAPI) MoveWithCollisions(sprite *LCDSprite, goalX, goalY float32) ([]SpriteCollisionInfo, float32, float32) {
	// TODO: implement with collision info
	return nil, goalX, goalY
}

// CheckCollisions checks for collisions without moving
func (s *SpriteAPI) CheckCollisions(sprite *LCDSprite, goalX, goalY float32) ([]SpriteCollisionInfo, float32, float32) {
	// TODO: implement with collision info
	return nil, goalX, goalY
}

// ============== Global Functions ==============

// UpdateAndDrawSprites updates and draws all sprites
func (s *SpriteAPI) UpdateAndDrawSprites() {
	C.pd_sprite_updateAndDrawSprites()
}

// DrawSprites draws all sprites without updating
func (s *SpriteAPI) DrawSprites() {
	C.pd_sprite_drawSprites()
}

// SetAlwaysRedraw sets whether sprites always redraw
func (s *SpriteAPI) SetAlwaysRedraw(flag bool) {
	var f C.int
	if flag {
		f = 1
	}
	C.pd_sprite_setAlwaysRedraw(f)
}

// ResetCollisionWorld resets the collision world
func (s *SpriteAPI) ResetCollisionWorld() {
	C.pd_sprite_resetCollisionWorld()
}

// ============== Queries (stubs) ==============

// QuerySpritesAtPoint returns sprites at a point
func (s *SpriteAPI) QuerySpritesAtPoint(x, y float32) []*LCDSprite {
	// TODO: implement
	return nil
}

// QuerySpritesInRect returns sprites in a rectangle
func (s *SpriteAPI) QuerySpritesInRect(x, y, w, h float32) []*LCDSprite {
	// TODO: implement
	return nil
}

// QuerySpritesAlongLine returns sprites along a line
func (s *SpriteAPI) QuerySpritesAlongLine(x1, y1, x2, y2 float32) []*LCDSprite {
	// TODO: implement
	return nil
}

// AllOverlappingSprites returns all overlapping sprites
func (s *SpriteAPI) AllOverlappingSprites() []*LCDSprite {
	// TODO: implement
	return nil
}

// SpriteCollisionInfo for collision results
type SpriteCollisionInfo struct {
	Sprite       *LCDSprite
	Other        *LCDSprite
	ResponseType SpriteCollisionResponseType
	Overlaps     bool
	Ti           float32
	Move         CollisionPoint
	Normal       CollisionVector
	Touch        CollisionPoint
	SpriteRect   PDRect
	OtherRect    PDRect
}
