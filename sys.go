// pdgo System API - unified CGO implementation

package pdgo

/*
#include <stdint.h>

// System API
void pd_sys_log(const char* msg);
void pd_sys_error(const char* msg);
void pd_sys_drawFPS(int x, int y);
uint32_t pd_sys_getCurrentTimeMS(void);
uint32_t pd_sys_getSecondsSinceEpoch(uint32_t* ms);
void pd_sys_getButtonState(uint32_t* current, uint32_t* pushed, uint32_t* released);
void pd_sys_setPeripheralsEnabled(uint32_t mask);
void pd_sys_getAccelerometer(float* x, float* y, float* z);
float pd_sys_getCrankChange(void);
float pd_sys_getCrankAngle(void);
int pd_sys_isCrankDocked(void);
int pd_sys_setCrankSoundsDisabled(int disabled);
int pd_sys_getFlipped(void);
void pd_sys_setAutoLockDisabled(int disabled);
int pd_sys_getLanguage(void);
float pd_sys_getBatteryPercentage(void);
float pd_sys_getBatteryVoltage(void);
*/
import "C"
import "unsafe"

// System provides access to Playdate system functions
type System struct{}

func newSystem() *System {
	return &System{}
}

// LogToConsole prints a message to the Playdate console
func (s *System) LogToConsole(msg string) {
	cstr := make([]byte, len(msg)+1)
	copy(cstr, msg)
	C.pd_sys_log((*C.char)(unsafe.Pointer(&cstr[0])))
}

// Error logs an error message
func (s *System) Error(msg string) {
	cstr := make([]byte, len(msg)+1)
	copy(cstr, msg)
	C.pd_sys_error((*C.char)(unsafe.Pointer(&cstr[0])))
}

// DrawFPS draws the current FPS at the given position
func (s *System) DrawFPS(x, y int) {
	C.pd_sys_drawFPS(C.int(x), C.int(y))
}

// GetCurrentTimeMilliseconds returns the current time in milliseconds
func (s *System) GetCurrentTimeMilliseconds() uint {
	return uint(C.pd_sys_getCurrentTimeMS())
}

// GetSecondsSinceEpoch returns seconds and milliseconds since epoch
func (s *System) GetSecondsSinceEpoch() (seconds uint, milliseconds uint) {
	var ms C.uint32_t
	sec := C.pd_sys_getSecondsSinceEpoch(&ms)
	return uint(sec), uint(ms)
}

// GetButtonState returns current, pushed, and released button states
func (s *System) GetButtonState() (current, pushed, released PDButtons) {
	var c, p, r C.uint32_t
	C.pd_sys_getButtonState(&c, &p, &r)
	return PDButtons(c), PDButtons(p), PDButtons(r)
}

// SetPeripheralsEnabled enables or disables peripherals like accelerometer
func (s *System) SetPeripheralsEnabled(mask PDPeripherals) {
	C.pd_sys_setPeripheralsEnabled(C.uint32_t(mask))
}

// GetAccelerometer returns accelerometer readings
func (s *System) GetAccelerometer() (x, y, z float32) {
	var cx, cy, cz C.float
	C.pd_sys_getAccelerometer(&cx, &cy, &cz)
	return float32(cx), float32(cy), float32(cz)
}

// GetCrankChange returns crank angle change since last call
func (s *System) GetCrankChange() float32 {
	return float32(C.pd_sys_getCrankChange())
}

// GetCrankAngle returns current crank angle (0-360)
func (s *System) GetCrankAngle() float32 {
	return float32(C.pd_sys_getCrankAngle())
}

// IsCrankDocked returns true if crank is docked
func (s *System) IsCrankDocked() bool {
	return C.pd_sys_isCrankDocked() != 0
}

// SetCrankSoundsDisabled disables/enables crank sounds
func (s *System) SetCrankSoundsDisabled(disabled bool) bool {
	var flag C.int
	if disabled {
		flag = 1
	}
	return C.pd_sys_setCrankSoundsDisabled(flag) != 0
}

// GetFlipped returns true if the device is flipped
func (s *System) GetFlipped() bool {
	return C.pd_sys_getFlipped() != 0
}

// SetAutoLockDisabled disables/enables auto lock
func (s *System) SetAutoLockDisabled(disabled bool) {
	var flag C.int
	if disabled {
		flag = 1
	}
	C.pd_sys_setAutoLockDisabled(flag)
}

// GetLanguage returns the system language
func (s *System) GetLanguage() PDLanguage {
	return PDLanguage(C.pd_sys_getLanguage())
}

// GetBatteryPercentage returns battery percentage (0.0-1.0)
func (s *System) GetBatteryPercentage() float32 {
	return float32(C.pd_sys_getBatteryPercentage())
}

// GetBatteryVoltage returns battery voltage
func (s *System) GetBatteryVoltage() float32 {
	return float32(C.pd_sys_getBatteryVoltage())
}

// SetUpdateCallback sets the update callback function
func (s *System) SetUpdateCallback(callback func() int) {
	SetUpdateCallback(callback)
}
