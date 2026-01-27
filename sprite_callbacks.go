// pdgo Sprite callbacks - unified CGO implementation
// Manages Go callbacks that are invoked from C code via trampolines

package pdgo

import "C"
import "unsafe"

// Callback type definitions
type spriteUpdateCallback func(*LCDSprite)
type spriteDrawCallback func(*LCDSprite, PDRect, PDRect)
type spriteCollisionCallback func(*LCDSprite, *LCDSprite) SpriteCollisionResponseType

// Callback registries - map sprite pointer to callback function
var (
	spriteUpdateCallbacks    = make(map[unsafe.Pointer]spriteUpdateCallback)
	spriteDrawCallbacks      = make(map[unsafe.Pointer]spriteDrawCallback)
	spriteCollisionCallbacks = make(map[unsafe.Pointer]spriteCollisionCallback)
)

// registerSpriteUpdateCallback registers an update callback for a sprite
func registerSpriteUpdateCallback(spritePtr unsafe.Pointer, fn spriteUpdateCallback) {
	if fn != nil {
		spriteUpdateCallbacks[spritePtr] = fn
	} else {
		delete(spriteUpdateCallbacks, spritePtr)
	}
}

// registerSpriteDrawCallback registers a draw callback for a sprite
func registerSpriteDrawCallback(spritePtr unsafe.Pointer, fn spriteDrawCallback) {
	if fn != nil {
		spriteDrawCallbacks[spritePtr] = fn
	} else {
		delete(spriteDrawCallbacks, spritePtr)
	}
}

// registerSpriteCollisionCallback registers a collision response callback
func registerSpriteCollisionCallback(spritePtr unsafe.Pointer, fn spriteCollisionCallback) {
	if fn != nil {
		spriteCollisionCallbacks[spritePtr] = fn
	} else {
		delete(spriteCollisionCallbacks, spritePtr)
	}
}

// unregisterSpriteCallbacks removes all callbacks for a sprite
func unregisterSpriteCallbacks(spritePtr unsafe.Pointer) {
	delete(spriteUpdateCallbacks, spritePtr)
	delete(spriteDrawCallbacks, spritePtr)
	delete(spriteCollisionCallbacks, spritePtr)
}

// ============== Exported trampolines called from C ==============

//export pdgo_sprite_update_trampoline
func pdgo_sprite_update_trampoline(spritePtr unsafe.Pointer) {
	if fn, ok := spriteUpdateCallbacks[spritePtr]; ok && fn != nil {
		fn(&LCDSprite{ptr: spritePtr})
	}
}

//export pdgo_sprite_draw_trampoline
func pdgo_sprite_draw_trampoline(spritePtr unsafe.Pointer, bx, by, bw, bh, dx, dy, dw, dh C.float) {
	if fn, ok := spriteDrawCallbacks[spritePtr]; ok && fn != nil {
		bounds := PDRect{X: float32(bx), Y: float32(by), Width: float32(bw), Height: float32(bh)}
		drawRect := PDRect{X: float32(dx), Y: float32(dy), Width: float32(dw), Height: float32(dh)}
		fn(&LCDSprite{ptr: spritePtr}, bounds, drawRect)
	}
}

//export pdgo_sprite_collision_trampoline
func pdgo_sprite_collision_trampoline(spritePtr, otherPtr unsafe.Pointer) C.int {
	if fn, ok := spriteCollisionCallbacks[spritePtr]; ok && fn != nil {
		return C.int(fn(&LCDSprite{ptr: spritePtr}, &LCDSprite{ptr: otherPtr}))
	}
	return C.int(CollisionTypeSlide) // default
}
