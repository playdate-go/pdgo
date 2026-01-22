package proc

const rawMainCgoGo = `//go:build !tinygo

package main

/*
#include <stdint.h>
*/
import "C"
import (
	"unsafe"

	"github.com/playdate-go/pdgo"
)

//export eventHandler
func eventHandler(playdateAPI unsafe.Pointer, event int32, arg uint32) int32 {
	if pdgo.PDSystemEvent(event) == pdgo.EventInit {
		pd = pdgo.Init(playdateAPI)
		initGame()
		pd.System.SetUpdateCallback(update)
	}
	return 0
}
`
