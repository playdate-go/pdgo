// pdgo Lua API - unified CGO implementation

package pdgo

/*
// Lua API
int pd_lua_getArgCount(void);
int pd_lua_getArgType(int pos, const char** outClass);
int pd_lua_argIsNil(int pos);
int pd_lua_getArgBool(int pos);
int pd_lua_getArgInt(int pos);
float pd_lua_getArgFloat(int pos);
const char* pd_lua_getArgString(int pos);
void pd_lua_pushNil(void);
void pd_lua_pushBool(int val);
void pd_lua_pushInt(int val);
void pd_lua_pushFloat(float val);
void pd_lua_pushString(const char* str);
*/
import "C"
import "unsafe"

// Lua provides access to Lua scripting
type Lua struct{}

func newLua() *Lua {
	return &Lua{}
}

// AddFunction adds a Go function callable from Lua
func (l *Lua) AddFunction(name string, callback func(args []interface{}) interface{}) error {
	// Not implemented - requires complex callback registration
	return &luaError{op: "addFunction", msg: "not supported"}
}

// GetArgCount returns number of Lua arguments
func (l *Lua) GetArgCount() int {
	return int(C.pd_lua_getArgCount())
}

// GetArgType returns type of argument at index and class name
func (l *Lua) GetArgType(pos int) (LuaType, string) {
	var outClass *C.char
	t := C.pd_lua_getArgType(C.int(pos), &outClass)
	className := ""
	if outClass != nil {
		className = goStringFromCStr(outClass)
	}
	return LuaType(t), className
}

// ArgIsNil checks if argument is nil
func (l *Lua) ArgIsNil(pos int) bool {
	return C.pd_lua_argIsNil(C.int(pos)) != 0
}

// GetArgBool gets boolean argument
func (l *Lua) GetArgBool(pos int) bool {
	return C.pd_lua_getArgBool(C.int(pos)) != 0
}

// GetArgInt gets integer argument
func (l *Lua) GetArgInt(pos int) int {
	return int(C.pd_lua_getArgInt(C.int(pos)))
}

// GetArgFloat gets float argument
func (l *Lua) GetArgFloat(pos int) float32 {
	return float32(C.pd_lua_getArgFloat(C.int(pos)))
}

// GetArgString gets string argument
func (l *Lua) GetArgString(pos int) string {
	cstr := C.pd_lua_getArgString(C.int(pos))
	if cstr == nil {
		return ""
	}
	return goStringFromCStr(cstr)
}

// PushNil pushes nil to Lua stack
func (l *Lua) PushNil() {
	C.pd_lua_pushNil()
}

// PushBool pushes bool to Lua stack
func (l *Lua) PushBool(val bool) {
	var v C.int
	if val {
		v = 1
	}
	C.pd_lua_pushBool(v)
}

// PushInt pushes int to Lua stack
func (l *Lua) PushInt(val int) {
	C.pd_lua_pushInt(C.int(val))
}

// PushFloat pushes float to Lua stack
func (l *Lua) PushFloat(val float32) {
	C.pd_lua_pushFloat(C.float(val))
}

// PushString pushes string to Lua stack
func (l *Lua) PushString(val string) {
	cstr := make([]byte, len(val)+1)
	copy(cstr, val)
	C.pd_lua_pushString((*C.char)(unsafe.Pointer(&cstr[0])))
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

// Helper to convert C string to Go string
func goStringFromCStr(cstr *C.char) string {
	if cstr == nil {
		return ""
	}
	var buf []byte
	ptr := unsafe.Pointer(cstr)
	for i := uintptr(0); ; i++ {
		b := *(*byte)(unsafe.Pointer(uintptr(ptr) + i))
		if b == 0 {
			break
		}
		buf = append(buf, b)
	}
	return string(buf)
}
