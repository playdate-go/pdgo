// pdgo Display API - unified CGO implementation

package pdgo

/*
// Display API
int pd_display_getWidth(void);
int pd_display_getHeight(void);
void pd_display_setRefreshRate(float rate);
float pd_display_getRefreshRate(void);
void pd_display_setInverted(int inverted);
void pd_display_setScale(unsigned int scale);
void pd_display_setMosaic(unsigned int x, unsigned int y);
void pd_display_setFlipped(int x, int y);
void pd_display_setOffset(int x, int y);
*/
import "C"

// Display provides access to display settings
type Display struct{}

func newDisplay() *Display {
	return &Display{}
}

// GetWidth returns display width
func (d *Display) GetWidth() int {
	return int(C.pd_display_getWidth())
}

// GetHeight returns display height
func (d *Display) GetHeight() int {
	return int(C.pd_display_getHeight())
}

// SetRefreshRate sets the display refresh rate (20-50 fps)
func (d *Display) SetRefreshRate(rate float32) {
	C.pd_display_setRefreshRate(C.float(rate))
}

// GetRefreshRate returns current refresh rate
func (d *Display) GetRefreshRate() float32 {
	return float32(C.pd_display_getRefreshRate())
}

// SetInverted sets whether display colors are inverted
func (d *Display) SetInverted(inverted bool) {
	var flag C.int
	if inverted {
		flag = 1
	}
	C.pd_display_setInverted(flag)
}

// SetScale sets display scale (1, 2, 4, or 8)
func (d *Display) SetScale(scale uint) {
	C.pd_display_setScale(C.uint(scale))
}

// SetMosaic sets mosaic effect
func (d *Display) SetMosaic(x, y uint) {
	C.pd_display_setMosaic(C.uint(x), C.uint(y))
}

// SetFlipped sets whether display is flipped
func (d *Display) SetFlipped(x, y bool) {
	var fx, fy C.int
	if x {
		fx = 1
	}
	if y {
		fy = 1
	}
	C.pd_display_setFlipped(fx, fy)
}

// SetOffset sets display offset
func (d *Display) SetOffset(x, y int) {
	C.pd_display_setOffset(C.int(x), C.int(y))
}
