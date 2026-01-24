//go:build tinygo

// TinyGo implementation of Sprite API
// API matches the CGO version for cross-platform compatibility

package pdgo

// Sprite represents a sprite object
type Sprite struct {
	ptr uintptr
}

// LCDSprite is an alias for Sprite (for API compatibility with CGO version)
type LCDSprite = Sprite

// SpriteAPI provides access to sprite functions
type SpriteAPI struct{}

func newSprite() *SpriteAPI {
	return &SpriteAPI{}
}

// ============== Sprite Creation ==============

// NewSprite creates a new sprite
func (s *SpriteAPI) NewSprite() *LCDSprite {
	if bridgeSpriteNewSprite != nil {
		ptr := bridgeSpriteNewSprite()
		if ptr != 0 {
			return &LCDSprite{ptr: ptr}
		}
	}
	return nil
}

// FreeSprite frees a sprite
func (s *SpriteAPI) FreeSprite(sprite *LCDSprite) {
	if sprite != nil && sprite.ptr != 0 {
		// Unregister all callbacks for this sprite
		unregisterSpriteCallbacks(sprite.ptr)
		if bridgeSpriteFreeSprite != nil {
			bridgeSpriteFreeSprite(sprite.ptr)
		}
		sprite.ptr = 0
	}
}

// ============== Display List ==============

// AddSprite adds a sprite to the display list
func (s *SpriteAPI) AddSprite(sprite *LCDSprite) {
	if bridgeSpriteAddSprite != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteAddSprite(sprite.ptr)
	}
}

// RemoveSprite removes a sprite from the display list
func (s *SpriteAPI) RemoveSprite(sprite *LCDSprite) {
	if bridgeSpriteRemoveSprite != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteRemoveSprite(sprite.ptr)
	}
}

// RemoveAllSprites removes all sprites from display
func (s *SpriteAPI) RemoveAllSprites() {
	if bridgeSpriteRemoveAllSprites != nil {
		bridgeSpriteRemoveAllSprites()
	}
}

// GetSpriteCount returns number of sprites
func (s *SpriteAPI) GetSpriteCount() int {
	if bridgeSpriteGetSpriteCount != nil {
		return int(bridgeSpriteGetSpriteCount())
	}
	return 0
}

// ============== Image ==============

// SetImage sets the sprite's image
func (s *SpriteAPI) SetImage(sprite *LCDSprite, image *LCDBitmap, flip LCDBitmapFlip) {
	if bridgeSpriteSetImage != nil && sprite != nil && sprite.ptr != 0 {
		var imgPtr uintptr
		if image != nil {
			imgPtr = image.ptr
		}
		bridgeSpriteSetImage(sprite.ptr, imgPtr, int32(flip))
	}
}

// GetImage returns the sprite's image
func (s *SpriteAPI) GetImage(sprite *LCDSprite) *LCDBitmap {
	if bridgeSpriteGetImage != nil && sprite != nil && sprite.ptr != 0 {
		ptr := bridgeSpriteGetImage(sprite.ptr)
		if ptr != 0 {
			return &LCDBitmap{ptr: ptr}
		}
	}
	return nil
}

// SetImageFlip sets image flip
func (s *SpriteAPI) SetImageFlip(sprite *LCDSprite, flip LCDBitmapFlip) {
	if bridgeSpriteSetImageFlip != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetImageFlip(sprite.ptr, int32(flip))
	}
}

// GetImageFlip returns image flip
func (s *SpriteAPI) GetImageFlip(sprite *LCDSprite) LCDBitmapFlip {
	if bridgeSpriteGetImageFlip != nil && sprite != nil && sprite.ptr != 0 {
		return LCDBitmapFlip(bridgeSpriteGetImageFlip(sprite.ptr))
	}
	return BitmapUnflipped
}

// ============== Position & Bounds ==============

// SetBounds sets the sprite's bounds (PDRect version for CGO compatibility)
func (s *SpriteAPI) SetBounds(sprite *LCDSprite, bounds PDRect) {
	if bridgeSpriteSetBounds != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetBounds(sprite.ptr, bounds.X, bounds.Y, bounds.Width, bounds.Height)
	}
}

// GetBounds returns the sprite's bounds (PDRect version for CGO compatibility)
func (s *SpriteAPI) GetBounds(sprite *LCDSprite) PDRect {
	if bridgeSpriteGetBounds != nil && sprite != nil && sprite.ptr != 0 {
		var x, y, w, h float32
		bridgeSpriteGetBounds(sprite.ptr, &x, &y, &w, &h)
		return PDRect{X: x, Y: y, Width: w, Height: h}
	}
	return PDRect{}
}

// MoveTo moves the sprite to position
func (s *SpriteAPI) MoveTo(sprite *LCDSprite, x, y float32) {
	if bridgeSpriteMoveTo != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteMoveTo(sprite.ptr, x, y)
	}
}

// MoveBy moves the sprite by delta
func (s *SpriteAPI) MoveBy(sprite *LCDSprite, dx, dy float32) {
	if bridgeSpriteMoveBy != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteMoveBy(sprite.ptr, dx, dy)
	}
}

// GetPosition returns sprite position
func (s *SpriteAPI) GetPosition(sprite *LCDSprite) (x, y float32) {
	if bridgeSpriteGetPosition != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteGetPosition(sprite.ptr, &x, &y)
	}
	return
}

// ============== Z-Index & Tags ==============

// SetZIndex sets sprite z-index
func (s *SpriteAPI) SetZIndex(sprite *LCDSprite, z int16) {
	if bridgeSpriteSetZIndex != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetZIndex(sprite.ptr, z)
	}
}

// GetZIndex returns sprite z-index
func (s *SpriteAPI) GetZIndex(sprite *LCDSprite) int16 {
	if bridgeSpriteGetZIndex != nil && sprite != nil && sprite.ptr != 0 {
		return bridgeSpriteGetZIndex(sprite.ptr)
	}
	return 0
}

// SetTag sets sprite tag
func (s *SpriteAPI) SetTag(sprite *LCDSprite, tag uint8) {
	if bridgeSpriteSetTag != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetTag(sprite.ptr, tag)
	}
}

// GetTag returns sprite tag
func (s *SpriteAPI) GetTag(sprite *LCDSprite) uint8 {
	if bridgeSpriteGetTag != nil && sprite != nil && sprite.ptr != 0 {
		return bridgeSpriteGetTag(sprite.ptr)
	}
	return 0
}

// ============== Callbacks ==============

// SetUpdateFunction sets the sprite's update callback
func (s *SpriteAPI) SetUpdateFunction(sprite *LCDSprite, callback func(*LCDSprite)) {
	if sprite == nil || sprite.ptr == 0 {
		return
	}

	// Register callback in Go-side registry
	registerSpriteUpdateCallback(sprite.ptr, callback)

	// Tell C runtime this sprite has a callback
	if bridgeSpriteSetUpdateFunction != nil {
		hasCallback := uintptr(0)
		if callback != nil {
			hasCallback = 1
		}
		bridgeSpriteSetUpdateFunction(sprite.ptr, hasCallback)
	}
}

// SetDrawFunction sets the sprite's draw callback
func (s *SpriteAPI) SetDrawFunction(sprite *LCDSprite, callback func(*LCDSprite, PDRect, PDRect)) {
	if sprite == nil || sprite.ptr == 0 {
		return
	}

	registerSpriteDrawCallback(sprite.ptr, callback)

	if bridgeSpriteSetDrawFunction != nil {
		hasCallback := uintptr(0)
		if callback != nil {
			hasCallback = 1
		}
		bridgeSpriteSetDrawFunction(sprite.ptr, hasCallback)
	}
}

// SetCollisionResponseFunction sets the sprite's collision response callback
func (s *SpriteAPI) SetCollisionResponseFunction(sprite *LCDSprite, callback func(*LCDSprite, *LCDSprite) SpriteCollisionResponseType) {
	if sprite == nil || sprite.ptr == 0 {
		return
	}

	registerSpriteCollisionCallback(sprite.ptr, callback)

	if bridgeSpriteSetCollisionResponseFunction != nil {
		hasCallback := uintptr(0)
		if callback != nil {
			hasCallback = 1
		}
		bridgeSpriteSetCollisionResponseFunction(sprite.ptr, hasCallback)
	}
}

// ============== Visibility & Drawing ==============

// SetVisible sets sprite visibility
func (s *SpriteAPI) SetVisible(sprite *LCDSprite, visible bool) {
	if bridgeSpriteSetVisible != nil && sprite != nil && sprite.ptr != 0 {
		var v int32
		if visible {
			v = 1
		}
		bridgeSpriteSetVisible(sprite.ptr, v)
	}
}

// IsVisible returns true if sprite is visible
func (s *SpriteAPI) IsVisible(sprite *LCDSprite) bool {
	if bridgeSpriteIsVisible != nil && sprite != nil && sprite.ptr != 0 {
		return bridgeSpriteIsVisible(sprite.ptr) != 0
	}
	return false
}

// SetOpaque sets whether sprite is opaque
func (s *SpriteAPI) SetOpaque(sprite *LCDSprite, opaque bool) {
	if bridgeSpriteSetOpaque != nil && sprite != nil && sprite.ptr != 0 {
		var o int32
		if opaque {
			o = 1
		}
		bridgeSpriteSetOpaque(sprite.ptr, o)
	}
}

// SetDrawMode sets sprite draw mode
func (s *SpriteAPI) SetDrawMode(sprite *LCDSprite, mode LCDBitmapDrawMode) {
	if bridgeSpriteSetDrawMode != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetDrawMode(sprite.ptr, int32(mode))
	}
}

// SetUpdatesEnabled enables/disables sprite updates
func (s *SpriteAPI) SetUpdatesEnabled(sprite *LCDSprite, enabled bool) {
	if bridgeSpriteSetUpdatesEnabled != nil && sprite != nil && sprite.ptr != 0 {
		var e int32
		if enabled {
			e = 1
		}
		bridgeSpriteSetUpdatesEnabled(sprite.ptr, e)
	}
}

// MarkDirty marks the sprite as needing redraw
func (s *SpriteAPI) MarkDirty(sprite *LCDSprite) {
	if bridgeSpriteMarkDirty != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteMarkDirty(sprite.ptr)
	}
}

// ============== Collision ==============

// SetCollideRect sets collision rectangle (PDRect version for CGO compatibility)
func (s *SpriteAPI) SetCollideRect(sprite *LCDSprite, collideRect PDRect) {
	if bridgeSpriteSetCollideRect != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetCollideRect(sprite.ptr, collideRect.X, collideRect.Y, collideRect.Width, collideRect.Height)
	}
}

// GetCollideRect returns collision rectangle (PDRect version for CGO compatibility)
func (s *SpriteAPI) GetCollideRect(sprite *LCDSprite) PDRect {
	if bridgeSpriteGetCollideRect != nil && sprite != nil && sprite.ptr != 0 {
		var x, y, w, h float32
		bridgeSpriteGetCollideRect(sprite.ptr, &x, &y, &w, &h)
		return PDRect{X: x, Y: y, Width: w, Height: h}
	}
	return PDRect{}
}

// ClearCollideRect clears collision rectangle
func (s *SpriteAPI) ClearCollideRect(sprite *LCDSprite) {
	if bridgeSpriteClearCollideRect != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteClearCollideRect(sprite.ptr)
	}
}

// SetCollisionsEnabled enables/disables collisions
func (s *SpriteAPI) SetCollisionsEnabled(sprite *LCDSprite, enabled bool) {
	if bridgeSpriteSetCollisionsEnabled != nil && sprite != nil && sprite.ptr != 0 {
		var e int32
		if enabled {
			e = 1
		}
		bridgeSpriteSetCollisionsEnabled(sprite.ptr, e)
	}
}

// MoveWithCollisions moves sprite with collision detection
// Returns collision info array, actual X, actual Y
func (s *SpriteAPI) MoveWithCollisions(sprite *LCDSprite, goalX, goalY float32) ([]SpriteCollisionInfo, float32, float32) {
	if bridgeSpriteMoveWithCollisions != nil && sprite != nil && sprite.ptr != 0 {
		var actualX, actualY float32
		var count int32
		bridgeSpriteMoveWithCollisions(sprite.ptr, goalX, goalY, &actualX, &actualY, &count)
		// TODO: parse collision info array from C when bridge supports it
		return nil, actualX, actualY
	}
	return nil, goalX, goalY
}

// CheckCollisions checks for collisions without moving
func (s *SpriteAPI) CheckCollisions(sprite *LCDSprite, goalX, goalY float32) ([]SpriteCollisionInfo, float32, float32) {
	if bridgeSpriteCheckCollisions != nil && sprite != nil && sprite.ptr != 0 {
		var count int32
		bridgeSpriteCheckCollisions(sprite.ptr, &count)
		// TODO: parse collision info array from C when bridge supports it
	}
	return nil, goalX, goalY
}

// ============== Global Functions ==============

// UpdateAndDrawSprites updates and draws all sprites
func (s *SpriteAPI) UpdateAndDrawSprites() {
	if bridgeSpriteUpdateAndDrawSprites != nil {
		bridgeSpriteUpdateAndDrawSprites()
	}
}

// DrawSprites draws all sprites without updating
func (s *SpriteAPI) DrawSprites() {
	if bridgeSpriteDrawSprites != nil {
		bridgeSpriteDrawSprites()
	}
}

// SetAlwaysRedraw sets whether sprites always redraw
func (s *SpriteAPI) SetAlwaysRedraw(flag bool) {
	if bridgeSpriteSetAlwaysRedraw != nil {
		var f int32
		if flag {
			f = 1
		}
		bridgeSpriteSetAlwaysRedraw(f)
	}
}

// ResetCollisionWorld resets the collision world
func (s *SpriteAPI) ResetCollisionWorld() {
	if bridgeSpriteResetCollisionWorld != nil {
		bridgeSpriteResetCollisionWorld()
	}
}

// ============== Queries ==============

// QuerySpritesAtPoint returns sprites at a point
func (s *SpriteAPI) QuerySpritesAtPoint(x, y float32) []*LCDSprite {
	if bridgeSpriteQuerySpritesAtPoint != nil {
		var count int32
		bridgeSpriteQuerySpritesAtPoint(x, y, &count)
		// TODO: parse sprite array from C when bridge supports it
	}
	return nil
}

// QuerySpritesInRect returns sprites in a rectangle
func (s *SpriteAPI) QuerySpritesInRect(x, y, w, h float32) []*LCDSprite {
	if bridgeSpriteQuerySpritesInRect != nil {
		var count int32
		bridgeSpriteQuerySpritesInRect(x, y, w, h, &count)
		// TODO: parse sprite array from C when bridge supports it
	}
	return nil
}

// QuerySpritesAlongLine returns sprites along a line
func (s *SpriteAPI) QuerySpritesAlongLine(x1, y1, x2, y2 float32) []*LCDSprite {
	if bridgeSpriteQuerySpritesAlongLine != nil {
		var count int32
		bridgeSpriteQuerySpritesAlongLine(x1, y1, x2, y2, &count)
		// TODO: parse sprite array from C when bridge supports it
	}
	return nil
}

// AllOverlappingSprites returns all overlapping sprites
func (s *SpriteAPI) AllOverlappingSprites() []*LCDSprite {
	if bridgeSpriteAllOverlappingSprites != nil {
		var count int32
		bridgeSpriteAllOverlappingSprites(&count)
		// TODO: parse sprite array from C when bridge supports it
	}
	return nil
}

// SpriteCollisionInfo for collision results (matching CGO version)
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
