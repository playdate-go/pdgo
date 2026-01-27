// pdgo API - unified CGO implementation for device and simulator

package pdgo

/*
void pd_sys_setUpdateCallback(void);
*/
import "C"
import (
	"unsafe"
)

// PlaydateAPI is the main API wrapper
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
var updateCallbackFn func() int

func SetUpdateCallback(callback func() int) {
	updateCallbackFn = callback
	// Register the callback with Playdate SDK via C
	C.pd_sys_setUpdateCallback()
}

// CallUpdateCallback is called from C trampoline
func CallUpdateCallback() int {
	if updateCallbackFn != nil {
		return updateCallbackFn()
	}
	return 1
}

//export pdgo_update_trampoline
func pdgo_update_trampoline() C.int {
	return C.int(CallUpdateCallback())
}

// PDSystemEvent represents system events
type PDSystemEvent int32

const (
	EventInit        PDSystemEvent = 0
	EventInitLua     PDSystemEvent = 1
	EventLock        PDSystemEvent = 2
	EventUnlock      PDSystemEvent = 3
	EventPause       PDSystemEvent = 4
	EventResume      PDSystemEvent = 5
	EventTerminate   PDSystemEvent = 6
	EventKeyPressed  PDSystemEvent = 7
	EventKeyReleased PDSystemEvent = 8
	EventLowPower    PDSystemEvent = 9
)
