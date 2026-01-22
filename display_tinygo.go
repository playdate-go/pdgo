//go:build tinygo

// TinyGo implementation of Display API

package pdgo

// Display provides access to display settings
type Display struct{}

func newDisplay() *Display {
	return &Display{}
}

// GetWidth returns display width
func (d *Display) GetWidth() int {
	if bridgeDisplayGetWidth != nil {
		return int(bridgeDisplayGetWidth())
	}
	return LCDColumns
}

// GetHeight returns display height
func (d *Display) GetHeight() int {
	if bridgeDisplayGetHeight != nil {
		return int(bridgeDisplayGetHeight())
	}
	return LCDRows
}

// SetRefreshRate sets the display refresh rate (20-50 fps)
func (d *Display) SetRefreshRate(rate float32) {
	if bridgeDisplaySetRefreshRate != nil {
		bridgeDisplaySetRefreshRate(rate)
	}
}

// GetRefreshRate returns current refresh rate
func (d *Display) GetRefreshRate() float32 {
	if bridgeDisplayGetRefreshRate != nil {
		return bridgeDisplayGetRefreshRate()
	}
	return 30.0
}

// SetInverted sets whether display colors are inverted
func (d *Display) SetInverted(inverted bool) {
	if bridgeDisplaySetInverted != nil {
		var flag int32
		if inverted {
			flag = 1
		}
		bridgeDisplaySetInverted(flag)
	}
}

// SetScale sets display scale (1, 2, 4, or 8)
func (d *Display) SetScale(scale uint) {
	if bridgeDisplaySetScale != nil {
		bridgeDisplaySetScale(uint32(scale))
	}
}

// SetMosaic sets mosaic effect
func (d *Display) SetMosaic(x, y uint) {
	if bridgeDisplaySetMosaic != nil {
		bridgeDisplaySetMosaic(uint32(x), uint32(y))
	}
}

// SetFlipped sets whether display is flipped
func (d *Display) SetFlipped(x, y bool) {
	if bridgeDisplaySetFlipped != nil {
		var fx, fy int32
		if x {
			fx = 1
		}
		if y {
			fy = 1
		}
		bridgeDisplaySetFlipped(fx, fy)
	}
}

// SetOffset sets display offset
func (d *Display) SetOffset(x, y int) {
	if bridgeDisplaySetOffset != nil {
		bridgeDisplaySetOffset(int32(x), int32(y))
	}
}
