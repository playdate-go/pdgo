//go:build tinygo

// TinyGo implementation of Sprite API

package pdgo

// Sprite represents a sprite object
type Sprite struct {
	ptr uintptr
}

// SpriteAPI provides access to sprite functions
type SpriteAPI struct{}

func newSprite() *SpriteAPI {
	return &SpriteAPI{}
}

// ============== Sprite Creation ==============

// NewSprite creates a new sprite
func (s *SpriteAPI) NewSprite() *Sprite {
	if bridgeSpriteNewSprite != nil {
		ptr := bridgeSpriteNewSprite()
		if ptr != 0 {
			return &Sprite{ptr: ptr}
		}
	}
	return nil
}

// FreeSprite frees a sprite
func (s *SpriteAPI) FreeSprite(sprite *Sprite) {
	if bridgeSpriteFreeSprite != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteFreeSprite(sprite.ptr)
		sprite.ptr = 0
	}
}

// ============== Display List ==============

// AddSprite adds a sprite to the display list
func (s *SpriteAPI) AddSprite(sprite *Sprite) {
	if bridgeSpriteAddSprite != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteAddSprite(sprite.ptr)
	}
}

// RemoveSprite removes a sprite from the display list
func (s *SpriteAPI) RemoveSprite(sprite *Sprite) {
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
func (s *SpriteAPI) SetImage(sprite *Sprite, image *LCDBitmap, flip LCDBitmapFlip) {
	if bridgeSpriteSetImage != nil && sprite != nil && sprite.ptr != 0 {
		var imgPtr uintptr
		if image != nil {
			imgPtr = image.ptr
		}
		bridgeSpriteSetImage(sprite.ptr, imgPtr, int32(flip))
	}
}

// GetImage returns the sprite's image
func (s *SpriteAPI) GetImage(sprite *Sprite) *LCDBitmap {
	if bridgeSpriteGetImage != nil && sprite != nil && sprite.ptr != 0 {
		ptr := bridgeSpriteGetImage(sprite.ptr)
		if ptr != 0 {
			return &LCDBitmap{ptr: ptr}
		}
	}
	return nil
}

// SetImageFlip sets image flip
func (s *SpriteAPI) SetImageFlip(sprite *Sprite, flip LCDBitmapFlip) {
	if bridgeSpriteSetImageFlip != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetImageFlip(sprite.ptr, int32(flip))
	}
}

// GetImageFlip returns image flip
func (s *SpriteAPI) GetImageFlip(sprite *Sprite) LCDBitmapFlip {
	if bridgeSpriteGetImageFlip != nil && sprite != nil && sprite.ptr != 0 {
		return LCDBitmapFlip(bridgeSpriteGetImageFlip(sprite.ptr))
	}
	return BitmapUnflipped
}

// ============== Position & Bounds ==============

// SetBounds sets the sprite's bounds
func (s *SpriteAPI) SetBounds(sprite *Sprite, x, y, w, h float32) {
	if bridgeSpriteSetBounds != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetBounds(sprite.ptr, x, y, w, h)
	}
}

// GetBounds returns the sprite's bounds
func (s *SpriteAPI) GetBounds(sprite *Sprite) (x, y, w, h float32) {
	if bridgeSpriteGetBounds != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteGetBounds(sprite.ptr, &x, &y, &w, &h)
	}
	return
}

// MoveTo moves the sprite to position
func (s *SpriteAPI) MoveTo(sprite *Sprite, x, y float32) {
	if bridgeSpriteMoveTo != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteMoveTo(sprite.ptr, x, y)
	}
}

// MoveBy moves the sprite by delta
func (s *SpriteAPI) MoveBy(sprite *Sprite, dx, dy float32) {
	if bridgeSpriteMoveBy != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteMoveBy(sprite.ptr, dx, dy)
	}
}

// GetPosition returns sprite position
func (s *SpriteAPI) GetPosition(sprite *Sprite) (x, y float32) {
	if bridgeSpriteGetPosition != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteGetPosition(sprite.ptr, &x, &y)
	}
	return
}

// ============== Z-Index & Tags ==============

// SetZIndex sets sprite z-index
func (s *SpriteAPI) SetZIndex(sprite *Sprite, z int16) {
	if bridgeSpriteSetZIndex != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetZIndex(sprite.ptr, z)
	}
}

// GetZIndex returns sprite z-index
func (s *SpriteAPI) GetZIndex(sprite *Sprite) int16 {
	if bridgeSpriteGetZIndex != nil && sprite != nil && sprite.ptr != 0 {
		return bridgeSpriteGetZIndex(sprite.ptr)
	}
	return 0
}

// SetTag sets sprite tag
func (s *SpriteAPI) SetTag(sprite *Sprite, tag uint8) {
	if bridgeSpriteSetTag != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetTag(sprite.ptr, tag)
	}
}

// GetTag returns sprite tag
func (s *SpriteAPI) GetTag(sprite *Sprite) uint8 {
	if bridgeSpriteGetTag != nil && sprite != nil && sprite.ptr != 0 {
		return bridgeSpriteGetTag(sprite.ptr)
	}
	return 0
}

// ============== Visibility & Drawing ==============

// SetVisible sets sprite visibility
func (s *SpriteAPI) SetVisible(sprite *Sprite, visible bool) {
	if bridgeSpriteSetVisible != nil && sprite != nil && sprite.ptr != 0 {
		var v int32
		if visible {
			v = 1
		}
		bridgeSpriteSetVisible(sprite.ptr, v)
	}
}

// IsVisible returns true if sprite is visible
func (s *SpriteAPI) IsVisible(sprite *Sprite) bool {
	if bridgeSpriteIsVisible != nil && sprite != nil && sprite.ptr != 0 {
		return bridgeSpriteIsVisible(sprite.ptr) != 0
	}
	return false
}

// SetOpaque sets whether sprite is opaque
func (s *SpriteAPI) SetOpaque(sprite *Sprite, opaque bool) {
	if bridgeSpriteSetOpaque != nil && sprite != nil && sprite.ptr != 0 {
		var o int32
		if opaque {
			o = 1
		}
		bridgeSpriteSetOpaque(sprite.ptr, o)
	}
}

// SetDrawMode sets sprite draw mode
func (s *SpriteAPI) SetDrawMode(sprite *Sprite, mode LCDBitmapDrawMode) {
	if bridgeSpriteSetDrawMode != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetDrawMode(sprite.ptr, int32(mode))
	}
}

// SetUpdatesEnabled enables/disables sprite updates
func (s *SpriteAPI) SetUpdatesEnabled(sprite *Sprite, enabled bool) {
	if bridgeSpriteSetUpdatesEnabled != nil && sprite != nil && sprite.ptr != 0 {
		var e int32
		if enabled {
			e = 1
		}
		bridgeSpriteSetUpdatesEnabled(sprite.ptr, e)
	}
}

// ============== Collision ==============

// SetCollideRect sets collision rectangle
func (s *SpriteAPI) SetCollideRect(sprite *Sprite, x, y, w, h float32) {
	if bridgeSpriteSetCollideRect != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteSetCollideRect(sprite.ptr, x, y, w, h)
	}
}

// GetCollideRect returns collision rectangle
func (s *SpriteAPI) GetCollideRect(sprite *Sprite) (x, y, w, h float32) {
	if bridgeSpriteGetCollideRect != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteGetCollideRect(sprite.ptr, &x, &y, &w, &h)
	}
	return
}

// ClearCollideRect clears collision rectangle
func (s *SpriteAPI) ClearCollideRect(sprite *Sprite) {
	if bridgeSpriteClearCollideRect != nil && sprite != nil && sprite.ptr != 0 {
		bridgeSpriteClearCollideRect(sprite.ptr)
	}
}

// SetCollisionsEnabled enables/disables collisions
func (s *SpriteAPI) SetCollisionsEnabled(sprite *Sprite, enabled bool) {
	if bridgeSpriteSetCollisionsEnabled != nil && sprite != nil && sprite.ptr != 0 {
		var e int32
		if enabled {
			e = 1
		}
		bridgeSpriteSetCollisionsEnabled(sprite.ptr, e)
	}
}

// MoveWithCollisions moves sprite with collision detection
func (s *SpriteAPI) MoveWithCollisions(sprite *Sprite, goalX, goalY float32) (actualX, actualY float32, collisions []CollisionInfo) {
	if bridgeSpriteMoveWithCollisions != nil && sprite != nil && sprite.ptr != 0 {
		var count int32
		bridgeSpriteMoveWithCollisions(sprite.ptr, goalX, goalY, &actualX, &actualY, &count)
		// Note: collision parsing would require more complex handling
		return actualX, actualY, nil
	}
	return goalX, goalY, nil
}

// CheckCollisions checks for collisions without moving
func (s *SpriteAPI) CheckCollisions(sprite *Sprite) []CollisionInfo {
	if bridgeSpriteCheckCollisions != nil && sprite != nil && sprite.ptr != 0 {
		var count int32
		bridgeSpriteCheckCollisions(sprite.ptr, &count)
		// Note: collision parsing would require more complex handling
	}
	return nil
}

// ============== Global Functions ==============

// UpdateAndDrawSprites updates and draws all sprites
func (s *SpriteAPI) UpdateAndDrawSprites() {
	if bridgeSpriteUpdateAndDrawSprites != nil {
		bridgeSpriteUpdateAndDrawSprites()
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
func (s *SpriteAPI) QuerySpritesAtPoint(x, y float32) []*Sprite {
	if bridgeSpriteQuerySpritesAtPoint != nil {
		var count int32
		bridgeSpriteQuerySpritesAtPoint(x, y, &count)
		// Note: sprite array parsing would require more complex handling
	}
	return nil
}

// QuerySpritesInRect returns sprites in a rectangle
func (s *SpriteAPI) QuerySpritesInRect(x, y, w, h float32) []*Sprite {
	if bridgeSpriteQuerySpritesInRect != nil {
		var count int32
		bridgeSpriteQuerySpritesInRect(x, y, w, h, &count)
		// Note: sprite array parsing would require more complex handling
	}
	return nil
}

// QuerySpritesAlongLine returns sprites along a line
func (s *SpriteAPI) QuerySpritesAlongLine(x1, y1, x2, y2 float32) []*Sprite {
	if bridgeSpriteQuerySpritesAlongLine != nil {
		var count int32
		bridgeSpriteQuerySpritesAlongLine(x1, y1, x2, y2, &count)
		// Note: sprite array parsing would require more complex handling
	}
	return nil
}

// AllOverlappingSprites returns all overlapping sprites
func (s *SpriteAPI) AllOverlappingSprites() []*Sprite {
	if bridgeSpriteAllOverlappingSprites != nil {
		var count int32
		bridgeSpriteAllOverlappingSprites(&count)
		// Note: sprite array parsing would require more complex handling
	}
	return nil
}
