//go:build tinygo

// TinyGo implementation of Lua API

package pdgo

// Lua provides access to Lua scripting
type Lua struct{}

func newLua() *Lua {
	return &Lua{}
}

// AddFunction adds a Go function callable from Lua
func (l *Lua) AddFunction(name string, callback func(args []interface{}) interface{}) error {
	// Not fully implemented for TinyGo - requires complex FFI
	return &luaError{op: "addFunction", msg: "not supported in TinyGo"}
}

// GetArgCount returns number of Lua arguments
func (l *Lua) GetArgCount() int {
	if bridgeLuaGetArgCount != nil {
		return int(bridgeLuaGetArgCount())
	}
	return 0
}

// GetArgType returns type of argument at index
func (l *Lua) GetArgType(pos int) LuaType {
	if bridgeLuaGetArgType != nil {
		return LuaType(bridgeLuaGetArgType(int32(pos)))
	}
	return LuaTypeNil
}

// ArgIsNil checks if argument is nil
func (l *Lua) ArgIsNil(pos int) bool {
	if bridgeLuaArgIsNil != nil {
		return bridgeLuaArgIsNil(int32(pos)) != 0
	}
	return true
}

// GetArgBool gets boolean argument
func (l *Lua) GetArgBool(pos int) bool {
	if bridgeLuaGetArgBool != nil {
		return bridgeLuaGetArgBool(int32(pos)) != 0
	}
	return false
}

// GetArgInt gets integer argument
func (l *Lua) GetArgInt(pos int) int {
	if bridgeLuaGetArgInt != nil {
		return int(bridgeLuaGetArgInt(int32(pos)))
	}
	return 0
}

// GetArgFloat gets float argument
func (l *Lua) GetArgFloat(pos int) float32 {
	if bridgeLuaGetArgFloat != nil {
		return bridgeLuaGetArgFloat(int32(pos))
	}
	return 0
}

// GetArgString gets string argument
func (l *Lua) GetArgString(pos int) string {
	if bridgeLuaGetArgString != nil {
		ptr := bridgeLuaGetArgString(int32(pos))
		if ptr != 0 {
			return goStringFromPtr(ptr)
		}
	}
	return ""
}

// PushNil pushes nil to Lua stack
func (l *Lua) PushNil() {
	if bridgeLuaPushNil != nil {
		bridgeLuaPushNil()
	}
}

// PushBool pushes bool to Lua stack
func (l *Lua) PushBool(val bool) {
	if bridgeLuaPushBool != nil {
		var v int32
		if val {
			v = 1
		}
		bridgeLuaPushBool(v)
	}
}

// PushInt pushes int to Lua stack
func (l *Lua) PushInt(val int) {
	if bridgeLuaPushInt != nil {
		bridgeLuaPushInt(int32(val))
	}
}

// PushFloat pushes float to Lua stack
func (l *Lua) PushFloat(val float32) {
	if bridgeLuaPushFloat != nil {
		bridgeLuaPushFloat(val)
	}
}

// PushString pushes string to Lua stack
func (l *Lua) PushString(val string) {
	if bridgeLuaPushString != nil {
		buf := make([]byte, len(val)+1)
		copy(buf, val)
		bridgeLuaPushString(&buf[0])
	}
}

// LuaType represents Lua value types
type LuaType int32

const (
	LuaTypeNil      LuaType = 0
	LuaTypeBool     LuaType = 1
	LuaTypeInt      LuaType = 2
	LuaTypeFloat    LuaType = 3
	LuaTypeString   LuaType = 4
	LuaTypeTable    LuaType = 5
	LuaTypeFunction LuaType = 6
	LuaTypeThread   LuaType = 7
	LuaTypeObject   LuaType = 8
)

type luaError struct {
	op  string
	msg string
}

func (e *luaError) Error() string {
	return "lua " + e.op + ": " + e.msg
}

// Helper to convert C string pointer to Go string
func goStringFromPtr(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	// Read bytes until null terminator
	var buf []byte
	for i := uintptr(0); ; i++ {
		b := *(*byte)(unsafePtr(ptr + i))
		if b == 0 {
			break
		}
		buf = append(buf, b)
	}
	return string(buf)
}
