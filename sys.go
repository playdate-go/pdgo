//go:build !tinygo

package pdgo

/*
#include "pd_api.h"
#include <stdlib.h>

// System API helper functions

static void* sys_realloc(const struct playdate_sys* sys, void* ptr, size_t size) {
    return sys->realloc(ptr, size);
}

static int sys_formatString(const struct playdate_sys* sys, char** ret, const char* fmt) {
    return sys->formatString(ret, fmt);
}

static void sys_logToConsole(const struct playdate_sys* sys, const char* fmt) {
    sys->logToConsole("%s", fmt);
}

static void sys_error(const struct playdate_sys* sys, const char* fmt) {
    sys->error("%s", fmt);
}

static PDLanguage sys_getLanguage(const struct playdate_sys* sys) {
    return sys->getLanguage();
}

static unsigned int sys_getCurrentTimeMilliseconds(const struct playdate_sys* sys) {
    return sys->getCurrentTimeMilliseconds();
}

static unsigned int sys_getSecondsSinceEpoch(const struct playdate_sys* sys, unsigned int* milliseconds) {
    return sys->getSecondsSinceEpoch(milliseconds);
}

static void sys_drawFPS(const struct playdate_sys* sys, int x, int y) {
    sys->drawFPS(x, y);
}

static void sys_setUpdateCallback(const struct playdate_sys* sys, PDCallbackFunction* update, void* userdata) {
    sys->setUpdateCallback(update, userdata);
}

static void sys_getButtonState(const struct playdate_sys* sys, PDButtons* current, PDButtons* pushed, PDButtons* released) {
    sys->getButtonState(current, pushed, released);
}

static void sys_setPeripheralsEnabled(const struct playdate_sys* sys, PDPeripherals mask) {
    sys->setPeripheralsEnabled(mask);
}

static void sys_getAccelerometer(const struct playdate_sys* sys, float* outx, float* outy, float* outz) {
    sys->getAccelerometer(outx, outy, outz);
}

static float sys_getCrankChange(const struct playdate_sys* sys) {
    return sys->getCrankChange();
}

static float sys_getCrankAngle(const struct playdate_sys* sys) {
    return sys->getCrankAngle();
}

static int sys_isCrankDocked(const struct playdate_sys* sys) {
    return sys->isCrankDocked();
}

static int sys_setCrankSoundsDisabled(const struct playdate_sys* sys, int flag) {
    return sys->setCrankSoundsDisabled(flag);
}

static int sys_getFlipped(const struct playdate_sys* sys) {
    return sys->getFlipped();
}

static void sys_setAutoLockDisabled(const struct playdate_sys* sys, int disable) {
    sys->setAutoLockDisabled(disable);
}

static void sys_setMenuImage(const struct playdate_sys* sys, LCDBitmap* bitmap, int xOffset) {
    sys->setMenuImage(bitmap, xOffset);
}

static PDMenuItem* sys_addMenuItem(const struct playdate_sys* sys, const char* title, PDMenuItemCallbackFunction* callback, void* userdata) {
    return sys->addMenuItem(title, callback, userdata);
}

static PDMenuItem* sys_addCheckmarkMenuItem(const struct playdate_sys* sys, const char* title, int value, PDMenuItemCallbackFunction* callback, void* userdata) {
    return sys->addCheckmarkMenuItem(title, value, callback, userdata);
}

static PDMenuItem* sys_addOptionsMenuItem(const struct playdate_sys* sys, const char* title, const char** optionTitles, int optionsCount, PDMenuItemCallbackFunction* f, void* userdata) {
    return sys->addOptionsMenuItem(title, optionTitles, optionsCount, f, userdata);
}

static void sys_removeAllMenuItems(const struct playdate_sys* sys) {
    sys->removeAllMenuItems();
}

static void sys_removeMenuItem(const struct playdate_sys* sys, PDMenuItem* menuItem) {
    sys->removeMenuItem(menuItem);
}

static int sys_getMenuItemValue(const struct playdate_sys* sys, PDMenuItem* menuItem) {
    return sys->getMenuItemValue(menuItem);
}

static void sys_setMenuItemValue(const struct playdate_sys* sys, PDMenuItem* menuItem, int value) {
    sys->setMenuItemValue(menuItem, value);
}

static const char* sys_getMenuItemTitle(const struct playdate_sys* sys, PDMenuItem* menuItem) {
    return sys->getMenuItemTitle(menuItem);
}

static void sys_setMenuItemTitle(const struct playdate_sys* sys, PDMenuItem* menuItem, const char* title) {
    sys->setMenuItemTitle(menuItem, title);
}

static void* sys_getMenuItemUserdata(const struct playdate_sys* sys, PDMenuItem* menuItem) {
    return sys->getMenuItemUserdata(menuItem);
}

static void sys_setMenuItemUserdata(const struct playdate_sys* sys, PDMenuItem* menuItem, void* ud) {
    sys->setMenuItemUserdata(menuItem, ud);
}

static int sys_getReduceFlashing(const struct playdate_sys* sys) {
    return sys->getReduceFlashing();
}

static float sys_getElapsedTime(const struct playdate_sys* sys) {
    return sys->getElapsedTime();
}

static void sys_resetElapsedTime(const struct playdate_sys* sys) {
    sys->resetElapsedTime();
}

static float sys_getBatteryPercentage(const struct playdate_sys* sys) {
    return sys->getBatteryPercentage();
}

static float sys_getBatteryVoltage(const struct playdate_sys* sys) {
    return sys->getBatteryVoltage();
}

static int32_t sys_getTimezoneOffset(const struct playdate_sys* sys) {
    return sys->getTimezoneOffset();
}

static int sys_shouldDisplay24HourTime(const struct playdate_sys* sys) {
    return sys->shouldDisplay24HourTime();
}

static void sys_convertEpochToDateTime(const struct playdate_sys* sys, uint32_t epoch, struct PDDateTime* datetime) {
    sys->convertEpochToDateTime(epoch, datetime);
}

static uint32_t sys_convertDateTimeToEpoch(const struct playdate_sys* sys, struct PDDateTime* datetime) {
    return sys->convertDateTimeToEpoch(datetime);
}

static void sys_clearICache(const struct playdate_sys* sys) {
    sys->clearICache();
}

static void sys_setButtonCallback(const struct playdate_sys* sys, PDButtonCallbackFunction* cb, void* buttonud, int queuesize) {
    sys->setButtonCallback(cb, buttonud, queuesize);
}

static void sys_setSerialMessageCallback(const struct playdate_sys* sys, void (*callback)(const char* data)) {
    sys->setSerialMessageCallback(callback);
}

static int sys_parseString(const struct playdate_sys* sys, const char* str, const char* format) {
    return sys->parseString(str, format);
}

static void sys_delay(const struct playdate_sys* sys, uint32_t milliseconds) {
    sys->delay(milliseconds);
}

static void sys_restartGame(const struct playdate_sys* sys, const char* launchargs) {
    sys->restartGame(launchargs);
}

static const char* sys_getLaunchArgs(const struct playdate_sys* sys, const char** outpath) {
    return sys->getLaunchArgs(outpath);
}

static bool sys_sendMirrorData(const struct playdate_sys* sys, uint8_t command, void* data, int len) {
    return sys->sendMirrorData(command, data, len);
}

static const struct PDInfo* sys_getSystemInfo(const struct playdate_sys* sys) {
    return sys->getSystemInfo();
}

// Go callback wrapper
extern int goUpdateCallback(void* userdata);

static int updateCallbackWrapper(void* userdata) {
    return goUpdateCallback(userdata);
}

static void sys_setGoUpdateCallback(const struct playdate_sys* sys, void* userdata) {
    sys->setUpdateCallback(updateCallbackWrapper, userdata);
}

// Menu callback wrapper
extern void goMenuItemCallback(void* userdata);

static void menuItemCallbackWrapper(void* userdata) {
    goMenuItemCallback(userdata);
}

static PDMenuItem* sys_addMenuItemGo(const struct playdate_sys* sys, const char* title, void* userdata) {
    return sys->addMenuItem(title, menuItemCallbackWrapper, userdata);
}

static PDMenuItem* sys_addCheckmarkMenuItemGo(const struct playdate_sys* sys, const char* title, int value, void* userdata) {
    return sys->addCheckmarkMenuItem(title, value, menuItemCallbackWrapper, userdata);
}

static PDMenuItem* sys_addOptionsMenuItemGo(const struct playdate_sys* sys, const char* title, const char** optionTitles, int optionsCount, void* userdata) {
    return sys->addOptionsMenuItem(title, optionTitles, optionsCount, menuItemCallbackWrapper, userdata);
}

// Button callback wrapper
extern int goButtonCallback(PDButtons button, int down, uint32_t when, void* userdata);

static int buttonCallbackWrapper(PDButtons button, int down, uint32_t when, void* userdata) {
    return goButtonCallback(button, down, when, userdata);
}

static void sys_setGoButtonCallback(const struct playdate_sys* sys, void* userdata, int queuesize) {
    sys->setButtonCallback(buttonCallbackWrapper, userdata, queuesize);
}

// Serial message callback wrapper
extern void goSerialMessageCallback(char* data);

static void serialMessageCallbackWrapper(const char* data) {
    goSerialMessageCallback((char*)data);
}

static void sys_setGoSerialMessageCallback(const struct playdate_sys* sys) {
    sys->setSerialMessageCallback(serialMessageCallbackWrapper);
}
*/
import "C"
import (
	"sync"
	"unsafe"
)

// PDButtons represents button states
type PDButtons uint32

const (
	ButtonLeft  PDButtons = C.kButtonLeft
	ButtonRight PDButtons = C.kButtonRight
	ButtonUp    PDButtons = C.kButtonUp
	ButtonDown  PDButtons = C.kButtonDown
	ButtonB     PDButtons = C.kButtonB
	ButtonA     PDButtons = C.kButtonA
)

// PDLanguage represents the system language
type PDLanguage int

const (
	LanguageEnglish  PDLanguage = C.kPDLanguageEnglish
	LanguageJapanese PDLanguage = C.kPDLanguageJapanese
	LanguageUnknown  PDLanguage = C.kPDLanguageUnknown
)

// PDPeripherals represents peripheral device flags
type PDPeripherals uint32

const (
	PeripheralNone          PDPeripherals = C.kNone
	PeripheralAccelerometer PDPeripherals = C.kAccelerometer
	PeripheralAll           PDPeripherals = C.kAllPeripherals
)

// PDDateTime represents date and time
type PDDateTime struct {
	Year    uint16
	Month   uint8 // 1-12
	Day     uint8 // 1-31
	Weekday uint8 // 1=Monday - 7=Sunday
	Hour    uint8 // 0-23
	Minute  uint8
	Second  uint8
}

// PDInfo contains system information
type PDInfo struct {
	OSVersion uint32
	Language  PDLanguage
}

// PDMenuItem represents a menu item
type PDMenuItem struct {
	ptr *C.PDMenuItem
}

// System wraps the playdate_sys API
type System struct {
	ptr *C.struct_playdate_sys
}

// Global callback storage
var (
	updateCallback        func() int
	menuItemCallbacks     = make(map[uintptr]func())
	buttonCallback        func(button PDButtons, down bool, when uint32) int
	serialMessageCallback func(data string)
	callbackMutex         sync.RWMutex
	menuItemIDCounter     uintptr
)

//export goUpdateCallback
func goUpdateCallback(userdata unsafe.Pointer) C.int {
	callbackMutex.RLock()
	cb := updateCallback
	callbackMutex.RUnlock()

	if cb != nil {
		return C.int(cb())
	}
	return 0
}

//export goMenuItemCallback
func goMenuItemCallback(userdata unsafe.Pointer) {
	id := uintptr(userdata)
	callbackMutex.RLock()
	cb, ok := menuItemCallbacks[id]
	callbackMutex.RUnlock()

	if ok && cb != nil {
		cb()
	}
}

//export goButtonCallback
func goButtonCallback(button C.PDButtons, down C.int, when C.uint32_t, userdata unsafe.Pointer) C.int {
	callbackMutex.RLock()
	cb := buttonCallback
	callbackMutex.RUnlock()

	if cb != nil {
		return C.int(cb(PDButtons(button), down != 0, uint32(when)))
	}
	return 0
}

//export goSerialMessageCallback
func goSerialMessageCallback(data *C.char) {
	callbackMutex.RLock()
	cb := serialMessageCallback
	callbackMutex.RUnlock()

	if cb != nil {
		cb(goString(data))
	}
}

func newSystem(ptr *C.struct_playdate_sys) *System {
	return &System{ptr: ptr}
}

// Realloc allocates, reallocates, or frees memory
func (s *System) Realloc(ptr unsafe.Pointer, size uint) unsafe.Pointer {
	return C.sys_realloc(s.ptr, ptr, C.size_t(size))
}

// Malloc allocates memory
func (s *System) Malloc(size uint) unsafe.Pointer {
	return C.sys_realloc(s.ptr, nil, C.size_t(size))
}

// Free frees memory
func (s *System) Free(ptr unsafe.Pointer) {
	C.sys_realloc(s.ptr, ptr, 0)
}

// LogToConsole logs a message to the console
func (s *System) LogToConsole(msg string) {
	cstr := cString(msg)
	defer freeCString(cstr)
	C.sys_logToConsole(s.ptr, cstr)
}

// Error logs an error and halts execution
func (s *System) Error(msg string) {
	cstr := cString(msg)
	defer freeCString(cstr)
	C.sys_error(s.ptr, cstr)
}

// GetLanguage returns the current system language
func (s *System) GetLanguage() PDLanguage {
	return PDLanguage(C.sys_getLanguage(s.ptr))
}

// GetCurrentTimeMilliseconds returns the current time in milliseconds
func (s *System) GetCurrentTimeMilliseconds() uint {
	return uint(C.sys_getCurrentTimeMilliseconds(s.ptr))
}

// GetSecondsSinceEpoch returns seconds since Unix epoch and optionally milliseconds
func (s *System) GetSecondsSinceEpoch() (seconds uint, milliseconds uint) {
	var ms C.uint
	sec := C.sys_getSecondsSinceEpoch(s.ptr, &ms)
	return uint(sec), uint(ms)
}

// DrawFPS draws the current FPS at the given position
func (s *System) DrawFPS(x, y int) {
	C.sys_drawFPS(s.ptr, C.int(x), C.int(y))
}

// SetUpdateCallback sets the update callback function
func (s *System) SetUpdateCallback(callback func() int) {
	callbackMutex.Lock()
	updateCallback = callback
	callbackMutex.Unlock()

	C.sys_setGoUpdateCallback(s.ptr, nil)
}

// GetButtonState returns the current, pushed, and released button states
func (s *System) GetButtonState() (current, pushed, released PDButtons) {
	var c, p, r C.PDButtons
	C.sys_getButtonState(s.ptr, &c, &p, &r)
	return PDButtons(c), PDButtons(p), PDButtons(r)
}

// SetPeripheralsEnabled enables or disables peripherals
func (s *System) SetPeripheralsEnabled(mask PDPeripherals) {
	C.sys_setPeripheralsEnabled(s.ptr, C.PDPeripherals(mask))
}

// GetAccelerometer returns accelerometer values
func (s *System) GetAccelerometer() (x, y, z float32) {
	var ox, oy, oz C.float
	C.sys_getAccelerometer(s.ptr, &ox, &oy, &oz)
	return float32(ox), float32(oy), float32(oz)
}

// GetCrankChange returns the change in crank angle since last call
func (s *System) GetCrankChange() float32 {
	return float32(C.sys_getCrankChange(s.ptr))
}

// GetCrankAngle returns the current crank angle
func (s *System) GetCrankAngle() float32 {
	return float32(C.sys_getCrankAngle(s.ptr))
}

// IsCrankDocked returns true if the crank is docked
func (s *System) IsCrankDocked() bool {
	return C.sys_isCrankDocked(s.ptr) != 0
}

// SetCrankSoundsDisabled enables or disables crank sounds
func (s *System) SetCrankSoundsDisabled(disabled bool) bool {
	flag := 0
	if disabled {
		flag = 1
	}
	return C.sys_setCrankSoundsDisabled(s.ptr, C.int(flag)) != 0
}

// GetFlipped returns true if the display is flipped
func (s *System) GetFlipped() bool {
	return C.sys_getFlipped(s.ptr) != 0
}

// SetAutoLockDisabled enables or disables auto-lock
func (s *System) SetAutoLockDisabled(disabled bool) {
	flag := 0
	if disabled {
		flag = 1
	}
	C.sys_setAutoLockDisabled(s.ptr, C.int(flag))
}

// SetMenuImage sets the menu image
func (s *System) SetMenuImage(bitmap *LCDBitmap, xOffset int) {
	var bmp *C.LCDBitmap
	if bitmap != nil {
		bmp = bitmap.ptr
	}
	C.sys_setMenuImage(s.ptr, bmp, C.int(xOffset))
}

// AddMenuItem adds a menu item with the given callback
func (s *System) AddMenuItem(title string, callback func()) *PDMenuItem {
	callbackMutex.Lock()
	menuItemIDCounter++
	id := menuItemIDCounter
	menuItemCallbacks[id] = callback
	callbackMutex.Unlock()

	cstr := cString(title)
	defer freeCString(cstr)

	ptr := C.sys_addMenuItemGo(s.ptr, cstr, unsafe.Pointer(id))
	return &PDMenuItem{ptr: ptr}
}

// AddCheckmarkMenuItem adds a checkmark menu item
func (s *System) AddCheckmarkMenuItem(title string, value bool, callback func()) *PDMenuItem {
	callbackMutex.Lock()
	menuItemIDCounter++
	id := menuItemIDCounter
	menuItemCallbacks[id] = callback
	callbackMutex.Unlock()

	cstr := cString(title)
	defer freeCString(cstr)

	v := 0
	if value {
		v = 1
	}

	ptr := C.sys_addCheckmarkMenuItemGo(s.ptr, cstr, C.int(v), unsafe.Pointer(id))
	return &PDMenuItem{ptr: ptr}
}

// AddOptionsMenuItem adds an options menu item
func (s *System) AddOptionsMenuItem(title string, options []string, callback func()) *PDMenuItem {
	callbackMutex.Lock()
	menuItemIDCounter++
	id := menuItemIDCounter
	menuItemCallbacks[id] = callback
	callbackMutex.Unlock()

	cstr := cString(title)
	defer freeCString(cstr)

	// Create C string array
	cOptions := make([]*C.char, len(options))
	for i, opt := range options {
		cOptions[i] = cString(opt)
	}
	defer func() {
		for _, copt := range cOptions {
			freeCString(copt)
		}
	}()

	ptr := C.sys_addOptionsMenuItemGo(s.ptr, cstr, (**C.char)(unsafe.Pointer(&cOptions[0])), C.int(len(options)), unsafe.Pointer(id))
	return &PDMenuItem{ptr: ptr}
}

// RemoveAllMenuItems removes all menu items
func (s *System) RemoveAllMenuItems() {
	callbackMutex.Lock()
	menuItemCallbacks = make(map[uintptr]func())
	callbackMutex.Unlock()

	C.sys_removeAllMenuItems(s.ptr)
}

// RemoveMenuItem removes a specific menu item
func (s *System) RemoveMenuItem(item *PDMenuItem) {
	C.sys_removeMenuItem(s.ptr, item.ptr)
}

// GetMenuItemValue returns the value of a menu item
func (s *System) GetMenuItemValue(item *PDMenuItem) int {
	return int(C.sys_getMenuItemValue(s.ptr, item.ptr))
}

// SetMenuItemValue sets the value of a menu item
func (s *System) SetMenuItemValue(item *PDMenuItem, value int) {
	C.sys_setMenuItemValue(s.ptr, item.ptr, C.int(value))
}

// GetMenuItemTitle returns the title of a menu item
func (s *System) GetMenuItemTitle(item *PDMenuItem) string {
	return goString(C.sys_getMenuItemTitle(s.ptr, item.ptr))
}

// SetMenuItemTitle sets the title of a menu item
func (s *System) SetMenuItemTitle(item *PDMenuItem, title string) {
	cstr := cString(title)
	defer freeCString(cstr)
	C.sys_setMenuItemTitle(s.ptr, item.ptr, cstr)
}

// GetReduceFlashing returns true if reduce flashing mode is enabled
func (s *System) GetReduceFlashing() bool {
	return C.sys_getReduceFlashing(s.ptr) != 0
}

// GetElapsedTime returns elapsed time since resetElapsedTime was called
func (s *System) GetElapsedTime() float32 {
	return float32(C.sys_getElapsedTime(s.ptr))
}

// ResetElapsedTime resets the elapsed time counter
func (s *System) ResetElapsedTime() {
	C.sys_resetElapsedTime(s.ptr)
}

// GetBatteryPercentage returns the battery percentage
func (s *System) GetBatteryPercentage() float32 {
	return float32(C.sys_getBatteryPercentage(s.ptr))
}

// GetBatteryVoltage returns the battery voltage
func (s *System) GetBatteryVoltage() float32 {
	return float32(C.sys_getBatteryVoltage(s.ptr))
}

// GetTimezoneOffset returns the timezone offset in seconds
func (s *System) GetTimezoneOffset() int32 {
	return int32(C.sys_getTimezoneOffset(s.ptr))
}

// ShouldDisplay24HourTime returns true if 24-hour time should be displayed
func (s *System) ShouldDisplay24HourTime() bool {
	return C.sys_shouldDisplay24HourTime(s.ptr) != 0
}

// ConvertEpochToDateTime converts Unix epoch to PDDateTime
func (s *System) ConvertEpochToDateTime(epoch uint32) PDDateTime {
	var dt C.struct_PDDateTime
	C.sys_convertEpochToDateTime(s.ptr, C.uint32_t(epoch), &dt)
	return PDDateTime{
		Year:    uint16(dt.year),
		Month:   uint8(dt.month),
		Day:     uint8(dt.day),
		Weekday: uint8(dt.weekday),
		Hour:    uint8(dt.hour),
		Minute:  uint8(dt.minute),
		Second:  uint8(dt.second),
	}
}

// ConvertDateTimeToEpoch converts PDDateTime to Unix epoch
func (s *System) ConvertDateTimeToEpoch(dt PDDateTime) uint32 {
	cdt := C.struct_PDDateTime{
		year:    C.uint16_t(dt.Year),
		month:   C.uint8_t(dt.Month),
		day:     C.uint8_t(dt.Day),
		weekday: C.uint8_t(dt.Weekday),
		hour:    C.uint8_t(dt.Hour),
		minute:  C.uint8_t(dt.Minute),
		second:  C.uint8_t(dt.Second),
	}
	return uint32(C.sys_convertDateTimeToEpoch(s.ptr, &cdt))
}

// ClearICache clears the instruction cache
func (s *System) ClearICache() {
	C.sys_clearICache(s.ptr)
}

// SetButtonCallback sets the button callback function
func (s *System) SetButtonCallback(callback func(button PDButtons, down bool, when uint32) int, queueSize int) {
	callbackMutex.Lock()
	buttonCallback = callback
	callbackMutex.Unlock()

	C.sys_setGoButtonCallback(s.ptr, nil, C.int(queueSize))
}

// SetSerialMessageCallback sets the serial message callback function
func (s *System) SetSerialMessageCallback(callback func(data string)) {
	callbackMutex.Lock()
	serialMessageCallback = callback
	callbackMutex.Unlock()

	C.sys_setGoSerialMessageCallback(s.ptr)
}

// Delay pauses execution for the specified number of milliseconds
func (s *System) Delay(milliseconds uint32) {
	C.sys_delay(s.ptr, C.uint32_t(milliseconds))
}

// RestartGame restarts the game with optional launch arguments
func (s *System) RestartGame(launchArgs string) {
	var cstr *C.char
	if launchArgs != "" {
		cstr = cString(launchArgs)
		defer freeCString(cstr)
	}
	C.sys_restartGame(s.ptr, cstr)
}

// GetLaunchArgs returns the launch arguments
func (s *System) GetLaunchArgs() (args string, path string) {
	var outpath *C.char
	cargs := C.sys_getLaunchArgs(s.ptr, &outpath)
	if cargs != nil {
		args = goString(cargs)
	}
	if outpath != nil {
		path = goString(outpath)
	}
	return
}

// SendMirrorData sends data to the mirror
func (s *System) SendMirrorData(command uint8, data []byte) bool {
	if len(data) == 0 {
		return bool(C.sys_sendMirrorData(s.ptr, C.uint8_t(command), nil, 0))
	}
	return bool(C.sys_sendMirrorData(s.ptr, C.uint8_t(command), unsafe.Pointer(&data[0]), C.int(len(data))))
}

// GetSystemInfo returns system information
func (s *System) GetSystemInfo() PDInfo {
	info := C.sys_getSystemInfo(s.ptr)
	if info == nil {
		return PDInfo{}
	}
	return PDInfo{
		OSVersion: uint32(info.osversion),
		Language:  PDLanguage(info.language),
	}
}
