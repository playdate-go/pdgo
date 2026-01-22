//go:build tinygo

// TinyGo implementation of pdgo API (no CGO)

package pdgo

import (
	"unsafe"
)

// PlaydateAPI is the main API wrapper for TinyGo
type PlaydateAPI struct {
	ptr         uintptr // C PlaydateAPI pointer
	System      *System
	Graphics    *Graphics
	Display     *Display
	Sprite      *SpriteAPI
	File        *File
	Sound       *Sound
	Lua         *Lua
	JSON        *JSON
	Scoreboards *Scoreboards
	Network     *Network
}

var api *PlaydateAPI
var pdPtr uintptr // Global C API pointer

// Init initializes the Playdate API with the provided C API pointer.
// This should be called from your eventHandler function.
func Init(playdateAPI unsafe.Pointer) *PlaydateAPI {
	pdPtr = uintptr(playdateAPI)

	api = &PlaydateAPI{
		ptr: pdPtr,
	}

	// Initialize all subsystems
	api.System = newSystem()
	api.Graphics = newGraphics()
	api.Display = newDisplay()
	api.Sprite = newSprite()
	api.File = newFile()
	api.Sound = newSound()
	api.Lua = newLua()
	api.JSON = newJSON()
	api.Scoreboards = newScoreboards()
	api.Network = newNetwork()

	return api
}

// GetAPI returns the current PlaydateAPI instance
func GetAPI() *PlaydateAPI {
	return api
}

// SetUpdateCallback sets the game update callback
// The callback is registered in C runtime, this just stores it for Go side
var updateCallbackFn func() int

func SetUpdateCallback(callback func() int) {
	updateCallbackFn = callback
}

// CallUpdateCallback is called from the game's updateCallback export
func CallUpdateCallback() int {
	if updateCallbackFn != nil {
		return updateCallbackFn()
	}
	return 1
}
