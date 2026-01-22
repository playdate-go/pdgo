//go:build tinygo

// TinyGo implementation of System API

package pdgo

// System provides access to Playdate system functions
type System struct{}

func newSystem() *System {
	return &System{}
}

// LogToConsole prints a message to the Playdate console
func (s *System) LogToConsole(msg string) {
	if bridgeSysLog != nil {
		buf := make([]byte, len(msg)+1)
		copy(buf, msg)
		bridgeSysLog(&buf[0])
	}
}

// Error logs an error message
func (s *System) Error(msg string) {
	if bridgeSysError != nil {
		buf := make([]byte, len(msg)+1)
		copy(buf, msg)
		bridgeSysError(&buf[0])
	} else {
		s.LogToConsole("ERROR: " + msg)
	}
}

// DrawFPS draws the current FPS at the given position
func (s *System) DrawFPS(x, y int) {
	if bridgeSysDrawFPS != nil {
		bridgeSysDrawFPS(int32(x), int32(y))
	}
}

// GetCurrentTimeMilliseconds returns the current time in milliseconds
func (s *System) GetCurrentTimeMilliseconds() uint {
	if bridgeSysGetCurrentTimeMS != nil {
		return uint(bridgeSysGetCurrentTimeMS())
	}
	return 0
}

// GetSecondsSinceEpoch returns seconds and milliseconds since epoch
func (s *System) GetSecondsSinceEpoch() (seconds uint, milliseconds uint) {
	if bridgeSysGetSecondsSinceEpoch != nil {
		var ms uint32
		sec := bridgeSysGetSecondsSinceEpoch(&ms)
		return uint(sec), uint(ms)
	}
	return 0, 0
}

// GetButtonState returns current, pushed, and released button states
func (s *System) GetButtonState() (current, pushed, released PDButtons) {
	if bridgeSysGetButtonState != nil {
		var c, p, r uint32
		bridgeSysGetButtonState(&c, &p, &r)
		return PDButtons(c), PDButtons(p), PDButtons(r)
	}
	return 0, 0, 0
}

// SetPeripheralsEnabled enables or disables peripherals like accelerometer
func (s *System) SetPeripheralsEnabled(mask PDPeripherals) {
	if bridgeSysSetPeripheralsEnabled != nil {
		bridgeSysSetPeripheralsEnabled(uint32(mask))
	}
}

// GetAccelerometer returns accelerometer readings
func (s *System) GetAccelerometer() (x, y, z float32) {
	if bridgeSysGetAccelerometer != nil {
		bridgeSysGetAccelerometer(&x, &y, &z)
	}
	return
}

// GetCrankChange returns crank angle change since last call
func (s *System) GetCrankChange() float32 {
	if bridgeSysGetCrankChange != nil {
		return bridgeSysGetCrankChange()
	}
	return 0
}

// GetCrankAngle returns current crank angle (0-360)
func (s *System) GetCrankAngle() float32 {
	if bridgeSysGetCrankAngle != nil {
		return bridgeSysGetCrankAngle()
	}
	return 0
}

// IsCrankDocked returns true if crank is docked
func (s *System) IsCrankDocked() bool {
	if bridgeSysIsCrankDocked != nil {
		return bridgeSysIsCrankDocked() != 0
	}
	return false
}

// SetCrankSoundsDisabled disables/enables crank sounds
func (s *System) SetCrankSoundsDisabled(disabled bool) bool {
	if bridgeSysSetCrankSoundsDisabled != nil {
		var flag int32
		if disabled {
			flag = 1
		}
		return bridgeSysSetCrankSoundsDisabled(flag) != 0
	}
	return false
}

// GetFlipped returns true if the device is flipped
func (s *System) GetFlipped() bool {
	if bridgeSysGetFlipped != nil {
		return bridgeSysGetFlipped() != 0
	}
	return false
}

// SetAutoLockDisabled disables/enables auto lock
func (s *System) SetAutoLockDisabled(disabled bool) {
	if bridgeSysSetAutoLockDisabled != nil {
		var flag int32
		if disabled {
			flag = 1
		}
		bridgeSysSetAutoLockDisabled(flag)
	}
}

// GetLanguage returns the system language
func (s *System) GetLanguage() PDLanguage {
	if bridgeSysGetLanguage != nil {
		return PDLanguage(bridgeSysGetLanguage())
	}
	return LanguageEnglish
}

// GetBatteryPercentage returns battery percentage (0.0-1.0)
func (s *System) GetBatteryPercentage() float32 {
	if bridgeSysGetBatteryPercentage != nil {
		return bridgeSysGetBatteryPercentage()
	}
	return 1.0
}

// GetBatteryVoltage returns battery voltage
func (s *System) GetBatteryVoltage() float32 {
	if bridgeSysGetBatteryVoltage != nil {
		return bridgeSysGetBatteryVoltage()
	}
	return 4.2
}

// SetUpdateCallback sets the update callback function
func (s *System) SetUpdateCallback(callback func() int) {
	SetUpdateCallback(callback)
}
