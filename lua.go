//go:build !tinygo

package pdgo

/*
#include "pd_api.h"
#include <stdlib.h>

// Lua API helper functions

static int lua_addFunction(const struct playdate_lua* lua, lua_CFunction f, const char* name, const char** outErr) {
    return lua->addFunction(f, name, outErr);
}

static int lua_registerClass(const struct playdate_lua* lua, const char* name, const lua_reg* reg, const lua_val* vals, int isstatic, const char** outErr) {
    return lua->registerClass(name, reg, vals, isstatic, outErr);
}

static void lua_pushFunction(const struct playdate_lua* lua, lua_CFunction f) {
    lua->pushFunction(f);
}

static int lua_indexMetatable(const struct playdate_lua* lua) {
    return lua->indexMetatable();
}

static void lua_stop(const struct playdate_lua* lua) {
    lua->stop();
}

static void lua_start(const struct playdate_lua* lua) {
    lua->start();
}

static int lua_getArgCount(const struct playdate_lua* lua) {
    return lua->getArgCount();
}

static enum LuaType lua_getArgType(const struct playdate_lua* lua, int pos, const char** outClass) {
    return lua->getArgType(pos, outClass);
}

static int lua_argIsNil(const struct playdate_lua* lua, int pos) {
    return lua->argIsNil(pos);
}

static int lua_getArgBool(const struct playdate_lua* lua, int pos) {
    return lua->getArgBool(pos);
}

static int lua_getArgInt(const struct playdate_lua* lua, int pos) {
    return lua->getArgInt(pos);
}

static float lua_getArgFloat(const struct playdate_lua* lua, int pos) {
    return lua->getArgFloat(pos);
}

static const char* lua_getArgString(const struct playdate_lua* lua, int pos) {
    return lua->getArgString(pos);
}

static const char* lua_getArgBytes(const struct playdate_lua* lua, int pos, size_t* outlen) {
    return lua->getArgBytes(pos, outlen);
}

static void* lua_getArgObject(const struct playdate_lua* lua, int pos, char* type, LuaUDObject** outud) {
    return lua->getArgObject(pos, type, outud);
}

static LCDBitmap* lua_getBitmap(const struct playdate_lua* lua, int pos) {
    return lua->getBitmap(pos);
}

static LCDSprite* lua_getSprite(const struct playdate_lua* lua, int pos) {
    return lua->getSprite(pos);
}

static void lua_pushNil(const struct playdate_lua* lua) {
    lua->pushNil();
}

static void lua_pushBool(const struct playdate_lua* lua, int val) {
    lua->pushBool(val);
}

static void lua_pushInt(const struct playdate_lua* lua, int val) {
    lua->pushInt(val);
}

static void lua_pushFloat(const struct playdate_lua* lua, float val) {
    lua->pushFloat(val);
}

static void lua_pushString(const struct playdate_lua* lua, const char* str) {
    lua->pushString(str);
}

static void lua_pushBytes(const struct playdate_lua* lua, const char* str, size_t len) {
    lua->pushBytes(str, len);
}

static void lua_pushBitmap(const struct playdate_lua* lua, LCDBitmap* bitmap) {
    lua->pushBitmap(bitmap);
}

static void lua_pushSprite(const struct playdate_lua* lua, LCDSprite* sprite) {
    lua->pushSprite(sprite);
}

static LuaUDObject* lua_pushObject(const struct playdate_lua* lua, void* obj, char* type, int nValues) {
    return lua->pushObject(obj, type, nValues);
}

static LuaUDObject* lua_retainObject(const struct playdate_lua* lua, LuaUDObject* obj) {
    return lua->retainObject(obj);
}

static void lua_releaseObject(const struct playdate_lua* lua, LuaUDObject* obj) {
    lua->releaseObject(obj);
}

static void lua_setUserValue(const struct playdate_lua* lua, LuaUDObject* obj, unsigned int slot) {
    lua->setUserValue(obj, slot);
}

static int lua_getUserValue(const struct playdate_lua* lua, LuaUDObject* obj, unsigned int slot) {
    return lua->getUserValue(obj, slot);
}

static int lua_callFunction(const struct playdate_lua* lua, const char* name, int nargs, const char** outerr) {
    return lua->callFunction(name, nargs, outerr);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

// LuaType represents Lua value types
type LuaType int

const (
	LuaTypeNil      LuaType = C.kTypeNil
	LuaTypeBool     LuaType = C.kTypeBool
	LuaTypeInt      LuaType = C.kTypeInt
	LuaTypeFloat    LuaType = C.kTypeFloat
	LuaTypeString   LuaType = C.kTypeString
	LuaTypeTable    LuaType = C.kTypeTable
	LuaTypeFunction LuaType = C.kTypeFunction
	LuaTypeThread   LuaType = C.kTypeThread
	LuaTypeObject   LuaType = C.kTypeObject
)

// LuaUDObject wraps a Lua userdata object
type LuaUDObject struct {
	ptr *C.LuaUDObject
}

// Lua wraps the playdate_lua API
type Lua struct {
	ptr *C.struct_playdate_lua
}

func newLua(ptr *C.struct_playdate_lua) *Lua {
	return &Lua{ptr: ptr}
}

// Stop stops the Lua runtime
func (l *Lua) Stop() {
	C.lua_stop(l.ptr)
}

// Start starts the Lua runtime
func (l *Lua) Start() {
	C.lua_start(l.ptr)
}

// GetArgCount returns the number of arguments
func (l *Lua) GetArgCount() int {
	return int(C.lua_getArgCount(l.ptr))
}

// GetArgType returns the type of an argument
func (l *Lua) GetArgType(pos int) (LuaType, string) {
	var outClass *C.char
	t := C.lua_getArgType(l.ptr, C.int(pos), &outClass)
	class := ""
	if outClass != nil {
		class = goString(outClass)
	}
	return LuaType(t), class
}

// ArgIsNil returns whether an argument is nil
func (l *Lua) ArgIsNil(pos int) bool {
	return C.lua_argIsNil(l.ptr, C.int(pos)) != 0
}

// GetArgBool returns a boolean argument
func (l *Lua) GetArgBool(pos int) bool {
	return C.lua_getArgBool(l.ptr, C.int(pos)) != 0
}

// GetArgInt returns an integer argument
func (l *Lua) GetArgInt(pos int) int {
	return int(C.lua_getArgInt(l.ptr, C.int(pos)))
}

// GetArgFloat returns a float argument
func (l *Lua) GetArgFloat(pos int) float32 {
	return float32(C.lua_getArgFloat(l.ptr, C.int(pos)))
}

// GetArgString returns a string argument
func (l *Lua) GetArgString(pos int) string {
	s := C.lua_getArgString(l.ptr, C.int(pos))
	if s == nil {
		return ""
	}
	return goString(s)
}

// GetArgBytes returns a bytes argument
func (l *Lua) GetArgBytes(pos int) []byte {
	var outlen C.size_t
	s := C.lua_getArgBytes(l.ptr, C.int(pos), &outlen)
	if s == nil || outlen == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(s), C.int(outlen))
}

// GetBitmap returns a bitmap argument
func (l *Lua) GetBitmap(pos int) *LCDBitmap {
	ptr := C.lua_getBitmap(l.ptr, C.int(pos))
	if ptr == nil {
		return nil
	}
	return &LCDBitmap{ptr: ptr}
}

// GetSprite returns a sprite argument
func (l *Lua) GetSprite(pos int) *LCDSprite {
	ptr := C.lua_getSprite(l.ptr, C.int(pos))
	if ptr == nil {
		return nil
	}
	return &LCDSprite{ptr: ptr}
}

// PushNil pushes nil onto the stack
func (l *Lua) PushNil() {
	C.lua_pushNil(l.ptr)
}

// PushBool pushes a boolean onto the stack
func (l *Lua) PushBool(val bool) {
	v := 0
	if val {
		v = 1
	}
	C.lua_pushBool(l.ptr, C.int(v))
}

// PushInt pushes an integer onto the stack
func (l *Lua) PushInt(val int) {
	C.lua_pushInt(l.ptr, C.int(val))
}

// PushFloat pushes a float onto the stack
func (l *Lua) PushFloat(val float32) {
	C.lua_pushFloat(l.ptr, C.float(val))
}

// PushString pushes a string onto the stack
func (l *Lua) PushString(str string) {
	cstr := cString(str)
	defer freeCString(cstr)
	C.lua_pushString(l.ptr, cstr)
}

// PushBytes pushes bytes onto the stack
func (l *Lua) PushBytes(data []byte) {
	if len(data) == 0 {
		C.lua_pushBytes(l.ptr, nil, 0)
		return
	}
	C.lua_pushBytes(l.ptr, (*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
}

// PushBitmap pushes a bitmap onto the stack
func (l *Lua) PushBitmap(bitmap *LCDBitmap) {
	var b *C.LCDBitmap
	if bitmap != nil {
		b = bitmap.ptr
	}
	C.lua_pushBitmap(l.ptr, b)
}

// PushSprite pushes a sprite onto the stack
func (l *Lua) PushSprite(sprite *LCDSprite) {
	var s *C.LCDSprite
	if sprite != nil {
		s = sprite.ptr
	}
	C.lua_pushSprite(l.ptr, s)
}

// PushObject pushes an object onto the stack
func (l *Lua) PushObject(obj unsafe.Pointer, typeName string, nValues int) *LuaUDObject {
	ctype := cString(typeName)
	defer freeCString(ctype)
	ptr := C.lua_pushObject(l.ptr, obj, ctype, C.int(nValues))
	if ptr == nil {
		return nil
	}
	return &LuaUDObject{ptr: ptr}
}

// RetainObject retains a Lua object
func (l *Lua) RetainObject(obj *LuaUDObject) *LuaUDObject {
	if obj == nil {
		return nil
	}
	ptr := C.lua_retainObject(l.ptr, obj.ptr)
	if ptr == nil {
		return nil
	}
	return &LuaUDObject{ptr: ptr}
}

// ReleaseObject releases a Lua object
func (l *Lua) ReleaseObject(obj *LuaUDObject) {
	if obj != nil && obj.ptr != nil {
		C.lua_releaseObject(l.ptr, obj.ptr)
		obj.ptr = nil
	}
}

// SetUserValue sets a user value on an object
func (l *Lua) SetUserValue(obj *LuaUDObject, slot uint) {
	if obj != nil {
		C.lua_setUserValue(l.ptr, obj.ptr, C.uint(slot))
	}
}

// GetUserValue gets a user value from an object
func (l *Lua) GetUserValue(obj *LuaUDObject, slot uint) int {
	if obj == nil {
		return 0
	}
	return int(C.lua_getUserValue(l.ptr, obj.ptr, C.uint(slot)))
}

// CallFunction calls a Lua function
func (l *Lua) CallFunction(name string, nargs int) error {
	cname := cString(name)
	defer freeCString(cname)

	var outerr *C.char
	result := C.lua_callFunction(l.ptr, cname, C.int(nargs), &outerr)
	if result == 0 {
		if outerr != nil {
			return errors.New(goString(outerr))
		}
		return errors.New("failed to call Lua function")
	}
	return nil
}

// IndexMetatable indexes the metatable
func (l *Lua) IndexMetatable() int {
	return int(C.lua_indexMetatable(l.ptr))
}
