//go:build tinygo

// Sprite callback registry for TinyGo
// Manages Go callbacks that are invoked from C code via trampolines

package pdgo

// Callback type definitions (matching CGO version)
type spriteUpdateCallback func(*LCDSprite)
type spriteDrawCallback func(*LCDSprite, PDRect, PDRect)
type spriteCollisionCallback func(*LCDSprite, *LCDSprite) SpriteCollisionResponseType

// Callback registries - map sprite pointer to callback function
var (
	spriteUpdateCallbacks    = make(map[uintptr]spriteUpdateCallback)
	spriteDrawCallbacks      = make(map[uintptr]spriteDrawCallback)
	spriteCollisionCallbacks = make(map[uintptr]spriteCollisionCallback)
)

// registerSpriteUpdateCallback registers an update callback for a sprite
func registerSpriteUpdateCallback(spritePtr uintptr, fn spriteUpdateCallback) {
	if fn != nil {
		spriteUpdateCallbacks[spritePtr] = fn
	} else {
		delete(spriteUpdateCallbacks, spritePtr)
	}
}

// registerSpriteDrawCallback registers a draw callback for a sprite
func registerSpriteDrawCallback(spritePtr uintptr, fn spriteDrawCallback) {
	if fn != nil {
		spriteDrawCallbacks[spritePtr] = fn
	} else {
		delete(spriteDrawCallbacks, spritePtr)
	}
}

// registerSpriteCollisionCallback registers a collision response callback
func registerSpriteCollisionCallback(spritePtr uintptr, fn spriteCollisionCallback) {
	if fn != nil {
		spriteCollisionCallbacks[spritePtr] = fn
	} else {
		delete(spriteCollisionCallbacks, spritePtr)
	}
}

// unregisterSpriteCallbacks removes all callbacks for a sprite
func unregisterSpriteCallbacks(spritePtr uintptr) {
	delete(spriteUpdateCallbacks, spritePtr)
	delete(spriteDrawCallbacks, spritePtr)
	delete(spriteCollisionCallbacks, spritePtr)
}

// ============== Exported trampolines called from C ==============

//export pdgo_sprite_update_trampoline
func pdgo_sprite_update_trampoline(spritePtr uintptr) {
	if fn, ok := spriteUpdateCallbacks[spritePtr]; ok && fn != nil {
		fn(&LCDSprite{ptr: spritePtr})
	}
}

//export pdgo_sprite_draw_trampoline
func pdgo_sprite_draw_trampoline(spritePtr uintptr, bx, by, bw, bh, dx, dy, dw, dh float32) {
	if fn, ok := spriteDrawCallbacks[spritePtr]; ok && fn != nil {
		bounds := PDRect{X: bx, Y: by, Width: bw, Height: bh}
		drawRect := PDRect{X: dx, Y: dy, Width: dw, Height: dh}
		fn(&LCDSprite{ptr: spritePtr}, bounds, drawRect)
	}
}

//export pdgo_sprite_collision_trampoline
func pdgo_sprite_collision_trampoline(spritePtr, otherPtr uintptr) int32 {
	if fn, ok := spriteCollisionCallbacks[spritePtr]; ok && fn != nil {
		return int32(fn(&LCDSprite{ptr: spritePtr}, &LCDSprite{ptr: otherPtr}))
	}
	return int32(CollisionTypeSlide) // default
}
