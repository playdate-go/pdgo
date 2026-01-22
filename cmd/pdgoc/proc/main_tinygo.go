package proc

const rawMainTinyGo = `//go:build tinygo

package main

import (
	"unsafe"

	"github.com/playdate-go/pdgo"
)

//export eventHandler
func eventHandler(playdateAPI unsafe.Pointer, event int32, arg uint32) int32 {
	if pdgo.PDSystemEvent(event) == pdgo.EventInit {
		pd = pdgo.Init(playdateAPI)
		initGame()
		pdgo.SetUpdateCallback(update)
	}
	return 0
}

//export updateCallback
func updateCallback(userdata unsafe.Pointer) int32 {
	return int32(pdgo.CallUpdateCallback())
}
`
