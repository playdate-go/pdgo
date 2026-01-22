//go:build !tinygo

package pdgo

/*
#include "pd_api.h"

// Display API helper functions

static int display_getWidth(const struct playdate_display* disp) {
    return disp->getWidth();
}

static int display_getHeight(const struct playdate_display* disp) {
    return disp->getHeight();
}

static void display_setRefreshRate(const struct playdate_display* disp, float rate) {
    disp->setRefreshRate(rate);
}

static void display_setInverted(const struct playdate_display* disp, int flag) {
    disp->setInverted(flag);
}

static void display_setScale(const struct playdate_display* disp, unsigned int s) {
    disp->setScale(s);
}

static void display_setMosaic(const struct playdate_display* disp, unsigned int x, unsigned int y) {
    disp->setMosaic(x, y);
}

static void display_setFlipped(const struct playdate_display* disp, int x, int y) {
    disp->setFlipped(x, y);
}

static void display_setOffset(const struct playdate_display* disp, int x, int y) {
    disp->setOffset(x, y);
}

static float display_getRefreshRate(const struct playdate_display* disp) {
    return disp->getRefreshRate();
}

static float display_getFPS(const struct playdate_display* disp) {
    return disp->getFPS();
}
*/
import "C"

// Display wraps the playdate_display API
type Display struct {
	ptr *C.struct_playdate_display
}

func newDisplay(ptr *C.struct_playdate_display) *Display {
	return &Display{ptr: ptr}
}

// GetWidth returns the display width
func (d *Display) GetWidth() int {
	return int(C.display_getWidth(d.ptr))
}

// GetHeight returns the display height
func (d *Display) GetHeight() int {
	return int(C.display_getHeight(d.ptr))
}

// SetRefreshRate sets the display refresh rate
func (d *Display) SetRefreshRate(rate float32) {
	C.display_setRefreshRate(d.ptr, C.float(rate))
}

// GetRefreshRate returns the display refresh rate
func (d *Display) GetRefreshRate() float32 {
	return float32(C.display_getRefreshRate(d.ptr))
}

// GetFPS returns the current frames per second
func (d *Display) GetFPS() float32 {
	return float32(C.display_getFPS(d.ptr))
}

// SetInverted sets whether the display is inverted
func (d *Display) SetInverted(inverted bool) {
	flag := 0
	if inverted {
		flag = 1
	}
	C.display_setInverted(d.ptr, C.int(flag))
}

// SetScale sets the display scale
func (d *Display) SetScale(scale uint) {
	C.display_setScale(d.ptr, C.uint(scale))
}

// SetMosaic sets the mosaic effect
func (d *Display) SetMosaic(x, y uint) {
	C.display_setMosaic(d.ptr, C.uint(x), C.uint(y))
}

// SetFlipped sets whether the display is flipped
func (d *Display) SetFlipped(x, y bool) {
	fx, fy := 0, 0
	if x {
		fx = 1
	}
	if y {
		fy = 1
	}
	C.display_setFlipped(d.ptr, C.int(fx), C.int(fy))
}

// SetOffset sets the display offset
func (d *Display) SetOffset(x, y int) {
	C.display_setOffset(d.ptr, C.int(x), C.int(y))
}
