package proc

const rawMainCgoGo = `//go:build !tinygo

package main

/*
#include <stdint.h>

// Set the global pd pointer in C code
void pd_set_api(void* playdate);
*/
import "C"
import (
	"unsafe"

	"github.com/playdate-go/pdgo"
)

//export eventHandler
func eventHandler(playdateAPI unsafe.Pointer, event C.int, arg C.uint32_t) C.int {
	_ = arg
	if pdgo.PDSystemEvent(event) == pdgo.EventInit {
		// Set C global pd pointer
		C.pd_set_api(playdateAPI)
		// Initialize Go side
		pd = pdgo.Init(playdateAPI)
		initGame()
		pdgo.SetUpdateCallback(update)
	}
	return 0
}
`
