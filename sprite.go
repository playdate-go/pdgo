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

// Collision and query functions
void* pd_sprite_moveWithCollisions(void* sprite, float goalX, float goalY, float* actualX, float* actualY, int* len);
void* pd_sprite_checkCollisions(void* sprite, float goalX, float goalY, float* actualX, float* actualY, int* len);
void* pd_sprite_querySpritesAtPoint(float x, float y, int* len);
void* pd_sprite_querySpritesInRect(float x, float y, float w, float h, int* len);
void* pd_sprite_querySpritesAlongLine(float x1, float y1, float x2, float y2, int* len);
void* pd_sprite_querySpriteInfoAlongLine(float x1, float y1, float x2, float y2, int* len);
void* pd_sprite_overlappingSprites(void* sprite, int* len);
void* pd_sprite_allOverlappingSprites(int* len);
void pd_sprite_freeArray(void* arr);

// Collision info accessors (portable across 32/64-bit)
void* pd_sprite_getCollisionAt(void* arr, int index);
void* pd_collision_getSprite(void* info);
void* pd_collision_getOther(void* info);
int pd_collision_getResponseType(void* info);
int pd_collision_getOverlaps(void* info);
float pd_collision_getTi(void* info);
float pd_collision_getMoveX(void* info);
float pd_collision_getMoveY(void* info);
int pd_collision_getNormalX(void* info);
int pd_collision_getNormalY(void* info);
float pd_collision_getTouchX(void* info);
float pd_collision_getTouchY(void* info);
float pd_collision_getSpriteRectX(void* info);
float pd_collision_getSpriteRectY(void* info);
float pd_collision_getSpriteRectW(void* info);
float pd_collision_getSpriteRectH(void* info);
float pd_collision_getOtherRectX(void* info);
float pd_collision_getOtherRectY(void* info);
float pd_collision_getOtherRectW(void* info);
float pd_collision_getOtherRectH(void* info);

// Query info accessors
void* pd_sprite_getQueryInfoAt(void* arr, int index);
void* pd_queryInfo_getSprite(void* info);
float pd_queryInfo_getTi1(void* info);
float pd_queryInfo_getTi2(void* info);
float pd_queryInfo_getEntryX(void* info);
float pd_queryInfo_getEntryY(void* info);
float pd_queryInfo_getExitX(void* info);
float pd_queryInfo_getExitY(void* info);

// Sprite array accessor
void* pd_spriteArray_getAt(void* arr, int index);

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

// Ptr returns the underlying C pointer for sprite comparison
func (s *Sprite) Ptr() unsafe.Pointer {
	if s == nil {
		return nil
	}
	return s.ptr
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
	if sprite == nil || sprite.ptr == nil {
		return nil, goalX, goalY
	}
	var actualX, actualY C.float
	var len C.int
	arr := C.pd_sprite_moveWithCollisions(sprite.ptr, C.float(goalX), C.float(goalY), &actualX, &actualY, &len)
	if arr == nil || len == 0 {
		return nil, float32(actualX), float32(actualY)
	}
	infos := parseCollisionInfo(arr, int(len))
	C.pd_sprite_freeArray(arr)
	return infos, float32(actualX), float32(actualY)
}

// CheckCollisions checks for collisions without moving
func (s *SpriteAPI) CheckCollisions(sprite *LCDSprite, goalX, goalY float32) ([]SpriteCollisionInfo, float32, float32) {
	if sprite == nil || sprite.ptr == nil {
		return nil, goalX, goalY
	}
	var actualX, actualY C.float
	var len C.int
	arr := C.pd_sprite_checkCollisions(sprite.ptr, C.float(goalX), C.float(goalY), &actualX, &actualY, &len)
	if arr == nil || len == 0 {
		return nil, float32(actualX), float32(actualY)
	}
	infos := parseCollisionInfo(arr, int(len))
	C.pd_sprite_freeArray(arr)
	return infos, float32(actualX), float32(actualY)
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

// ============== Queries ==============

// QuerySpritesAtPoint returns sprites at a point
func (s *SpriteAPI) QuerySpritesAtPoint(x, y float32) []*LCDSprite {
	var len C.int
	arr := C.pd_sprite_querySpritesAtPoint(C.float(x), C.float(y), &len)
	if arr == nil || len == 0 {
		return nil
	}
	sprites := parseSpriteArray(arr, int(len))
	C.pd_sprite_freeArray(arr)
	return sprites
}

// QuerySpritesInRect returns sprites in a rectangle
func (s *SpriteAPI) QuerySpritesInRect(x, y, w, h float32) []*LCDSprite {
	var len C.int
	arr := C.pd_sprite_querySpritesInRect(C.float(x), C.float(y), C.float(w), C.float(h), &len)
	if arr == nil || len == 0 {
		return nil
	}
	sprites := parseSpriteArray(arr, int(len))
	C.pd_sprite_freeArray(arr)
	return sprites
}

// QuerySpritesAlongLine returns sprites along a line
func (s *SpriteAPI) QuerySpritesAlongLine(x1, y1, x2, y2 float32) []*LCDSprite {
	var len C.int
	arr := C.pd_sprite_querySpritesAlongLine(C.float(x1), C.float(y1), C.float(x2), C.float(y2), &len)
	if arr == nil || len == 0 {
		return nil
	}
	sprites := parseSpriteArray(arr, int(len))
	C.pd_sprite_freeArray(arr)
	return sprites
}

// QuerySpriteInfoAlongLine returns sprite query info along a line
func (s *SpriteAPI) QuerySpriteInfoAlongLine(x1, y1, x2, y2 float32) []SpriteQueryInfo {
	var len C.int
	arr := C.pd_sprite_querySpriteInfoAlongLine(C.float(x1), C.float(y1), C.float(x2), C.float(y2), &len)
	if arr == nil || len == 0 {
		return nil
	}
	infos := parseSpriteQueryInfo(arr, int(len))
	C.pd_sprite_freeArray(arr)
	return infos
}

// AllOverlappingSprites returns all overlapping sprites
func (s *SpriteAPI) AllOverlappingSprites() []*LCDSprite {
	var len C.int
	arr := C.pd_sprite_allOverlappingSprites(&len)
	if arr == nil || len == 0 {
		return nil
	}
	sprites := parseSpriteArray(arr, int(len))
	C.pd_sprite_freeArray(arr)
	return sprites
}

// SpriteQueryInfo contains information about a sprite intersection along a line
type SpriteQueryInfo struct {
	Sprite     *LCDSprite
	Ti1        float32
	Ti2        float32
	EntryPoint CollisionPoint
	ExitPoint  CollisionPoint
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

// parseCollisionInfo parses C SpriteCollisionInfo array into Go slice
// Uses C accessor functions for portability across 32/64-bit architectures
func parseCollisionInfo(arr unsafe.Pointer, len int) []SpriteCollisionInfo {
	if arr == nil || len == 0 {
		return nil
	}

	infos := make([]SpriteCollisionInfo, len)

	for i := 0; i < len; i++ {
		info := C.pd_sprite_getCollisionAt(arr, C.int(i))
		infos[i] = SpriteCollisionInfo{
			Sprite:       &LCDSprite{ptr: C.pd_collision_getSprite(info)},
			Other:        &LCDSprite{ptr: C.pd_collision_getOther(info)},
			ResponseType: SpriteCollisionResponseType(C.pd_collision_getResponseType(info)),
			Overlaps:     C.pd_collision_getOverlaps(info) != 0,
			Ti:           float32(C.pd_collision_getTi(info)),
			Move:         CollisionPoint{X: float32(C.pd_collision_getMoveX(info)), Y: float32(C.pd_collision_getMoveY(info))},
			Normal:       CollisionVector{X: int(C.pd_collision_getNormalX(info)), Y: int(C.pd_collision_getNormalY(info))},
			Touch:        CollisionPoint{X: float32(C.pd_collision_getTouchX(info)), Y: float32(C.pd_collision_getTouchY(info))},
			SpriteRect: PDRect{
				X:      float32(C.pd_collision_getSpriteRectX(info)),
				Y:      float32(C.pd_collision_getSpriteRectY(info)),
				Width:  float32(C.pd_collision_getSpriteRectW(info)),
				Height: float32(C.pd_collision_getSpriteRectH(info)),
			},
			OtherRect: PDRect{
				X:      float32(C.pd_collision_getOtherRectX(info)),
				Y:      float32(C.pd_collision_getOtherRectY(info)),
				Width:  float32(C.pd_collision_getOtherRectW(info)),
				Height: float32(C.pd_collision_getOtherRectH(info)),
			},
		}
	}
	return infos
}

// parseSpriteQueryInfo parses C SpriteQueryInfo array into Go slice
// Uses C accessor functions for portability across 32/64-bit architectures
func parseSpriteQueryInfo(arr unsafe.Pointer, len int) []SpriteQueryInfo {
	if arr == nil || len == 0 {
		return nil
	}

	infos := make([]SpriteQueryInfo, len)

	for i := 0; i < len; i++ {
		info := C.pd_sprite_getQueryInfoAt(arr, C.int(i))
		infos[i] = SpriteQueryInfo{
			Sprite: &LCDSprite{ptr: C.pd_queryInfo_getSprite(info)},
			Ti1:    float32(C.pd_queryInfo_getTi1(info)),
			Ti2:    float32(C.pd_queryInfo_getTi2(info)),
			EntryPoint: CollisionPoint{
				X: float32(C.pd_queryInfo_getEntryX(info)),
				Y: float32(C.pd_queryInfo_getEntryY(info)),
			},
			ExitPoint: CollisionPoint{
				X: float32(C.pd_queryInfo_getExitX(info)),
				Y: float32(C.pd_queryInfo_getExitY(info)),
			},
		}
	}
	return infos
}

// parseSpriteArray parses C LCDSprite* array into Go slice
// Uses C accessor function for portability across 32/64-bit architectures
func parseSpriteArray(arr unsafe.Pointer, len int) []*LCDSprite {
	if arr == nil || len == 0 {
		return nil
	}

	sprites := make([]*LCDSprite, len)

	for i := 0; i < len; i++ {
		sprites[i] = &LCDSprite{ptr: C.pd_spriteArray_getAt(arr, C.int(i))}
	}
	return sprites
}
