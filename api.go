//go:build !tinygo

// Package pdgo provides Go bindings for the Playdate SDK C API.
// This is the CGO implementation for simulator builds. The TinyGo implementation for device builds is the same.
// CGO (simulator) and TinyGo (device) APIs are fully consistent
package pdgo

/*
#cgo CFLAGS: -DTARGET_EXTENSION=1
#include "pd_api.h"
#include <stdlib.h>
#include <string.h>

// Global API pointer - set by Init()
static PlaydateAPI* pd = NULL;

void pdgo_set_api(PlaydateAPI* api) {
    pd = api;
}

PlaydateAPI* pdgo_get_api(void) {
    return pd;
}
*/
import "C"
import (
	"unsafe"
)

// PDSystemEvent represents system events sent to the game
type PDSystemEvent int

const (
	EventInit          PDSystemEvent = C.kEventInit
	EventInitLua       PDSystemEvent = C.kEventInitLua
	EventLock          PDSystemEvent = C.kEventLock
	EventUnlock        PDSystemEvent = C.kEventUnlock
	EventPause         PDSystemEvent = C.kEventPause
	EventResume        PDSystemEvent = C.kEventResume
	EventTerminate     PDSystemEvent = C.kEventTerminate
	EventKeyPressed    PDSystemEvent = C.kEventKeyPressed
	EventKeyReleased   PDSystemEvent = C.kEventKeyReleased
	EventLowPower      PDSystemEvent = C.kEventLowPower
	EventMirrorStarted PDSystemEvent = C.kEventMirrorStarted
	EventMirrorEnded   PDSystemEvent = C.kEventMirrorEnded
)

// String returns the string representation of the event
func (e PDSystemEvent) String() string {
	switch e {
	case EventInit:
		return "Init"
	case EventInitLua:
		return "InitLua"
	case EventLock:
		return "Lock"
	case EventUnlock:
		return "Unlock"
	case EventPause:
		return "Pause"
	case EventResume:
		return "Resume"
	case EventTerminate:
		return "Terminate"
	case EventKeyPressed:
		return "KeyPressed"
	case EventKeyReleased:
		return "KeyReleased"
	case EventLowPower:
		return "LowPower"
	case EventMirrorStarted:
		return "MirrorStarted"
	case EventMirrorEnded:
		return "MirrorEnded"
	default:
		return "Unknown"
	}
}

// PlaydateAPI is the main API wrapper
type PlaydateAPI struct {
	ptr         *C.PlaydateAPI
	System      *System
	File        *File
	Graphics    *Graphics
	Sprite      *Sprite
	Display     *Display
	Sound       *Sound
	Lua         *Lua
	JSON        *JSON
	Scoreboards *Scoreboards
	Network     *Network
}

var api *PlaydateAPI

// Init initializes the Playdate API with the provided C API pointer.
// This should be called from your eventHandler function.
func Init(playdateAPI unsafe.Pointer) *PlaydateAPI {
	cAPI := (*C.PlaydateAPI)(playdateAPI)
	C.pdgo_set_api(cAPI)

	api = &PlaydateAPI{
		ptr: cAPI,
	}

	// Initialize subsystems
	api.System = newSystem(cAPI.system)
	api.File = newFile(cAPI.file)
	api.Graphics = newGraphics(cAPI.graphics)
	api.Sprite = newSprite(cAPI.sprite)
	api.Display = newDisplay(cAPI.display)
	api.Sound = newSound(cAPI.sound)
	api.Lua = newLua(cAPI.lua)
	api.JSON = newJSON(cAPI.json)
	api.Scoreboards = newScoreboards(cAPI.scoreboards)
	api.Network = newNetwork(cAPI.network)

	return api
}

// GetAPI returns the current PlaydateAPI instance
func GetAPI() *PlaydateAPI {
	return api
}

// cString converts a Go string to a C string. The caller must free the result.
func cString(s string) *C.char {
	return C.CString(s)
}

// goString converts a C string to a Go string
func goString(s *C.char) string {
	return C.GoString(s)
}

// freeCString frees a C string
func freeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
}
