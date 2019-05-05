/*
* @Author: konyka
* @Date:   2019-05-05 09:40:08
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 10:22:29
*/

package state

import "fmt"
import "io/ioutil"
import . "lunago/api"

import "lunago/stdlib"


func (self *luaState) TypeName2(idx int) string {
    return self.TypeName(self.Type(idx))
}


func (self *luaState) Len2(idx int) int64 {
    self.Len(idx)
    i, isNum := self.ToIntegerX(-1)
    if !isNum {
        self.Error2("object length is not an integer")
    }
    self.Pop(1)
    return i
}

func (self *luaState) CheckStack2(sz int, msg string) {
    if !self.CheckStack(sz) {
        if msg != "" {
            self.Error2("stack overflow (%s)", msg)
        } else {
            self.Error2("stack overflow")
        }
    }
}

func (self *luaState) Error2(fmt string, a ...interface{}) int {
    self.PushFString(fmt, a...) // todo
    return self.Error()
}

func (self *luaState) ToString2(idx int) string {
    if self.CallMeta(idx, "__tostring") { /* metafield? */
        if !self.IsString(-1) {
            self.Error2("'__tostring' must return a string")
        }
    } else {
        switch self.Type(idx) {
        case LUA_TNUMBER:
            if self.IsInteger(idx) {
                self.PushString(fmt.Sprintf("%d", self.ToInteger(idx))) // todo
            } else {
                self.PushString(fmt.Sprintf("%g", self.ToNumber(idx))) // todo
            }
        case LUA_TSTRING:
            self.PushValue(idx)
        case LUA_TBOOLEAN:
            if self.ToBoolean(idx) {
                self.PushString("true")
            } else {
                self.PushString("false")
            }
        case LUA_TNIL:
            self.PushString("nil")
        default:
            tt := self.GetMetafield(idx, "__name") /* try name */
            var kind string
            if tt == LUA_TSTRING {
                kind = self.CheckString(-1)
            } else {
                kind = self.TypeName2(idx)
            }

            self.PushString(fmt.Sprintf("%s: %p", kind, self.ToPointer(idx)))
            if tt != LUA_TNIL {
                self.Remove(-2) /* remove '__name' */
            }
        }
    }
    return self.CheckString(-1)
}

func (self *luaState) LoadString(s string) int {
    return self.Load([]byte(s), s, "bt")
}

func (self *luaState) LoadFileX(filename, mode string) int {
    if data, err := ioutil.ReadFile(filename); err == nil {
        return self.Load(data, "@" + filename, mode)
    }
    return LUA_ERRFILE
}

func (self *luaState) LoadFile(filename string) int {
    return self.LoadFileX(filename, "bt")
}

func (self *luaState) DoString(str string) bool {
    return self.LoadString(str) == LUA_OK &&
        self.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

func (self *luaState) DoFile(filename string) bool {
    return self.LoadFile(filename) == LUA_OK &&
        self.PCall(0, LUA_MULTRET, 0) == LUA_OK
}


func (self *luaState) ArgError(arg int, extraMsg string) int {
    // bad argument #arg to 'funcname' (extramsg)
    return self.Error2("bad argument #%d (%s)", arg, extraMsg) // todo
}

func (self *luaState) ArgCheck(cond bool, arg int, extraMsg string) {
    if !cond {
        self.ArgError(arg, extraMsg)
    }
}

func (self *luaState) CheckAny(arg int) {
    if self.Type(arg) == LUA_TNONE {
        self.ArgError(arg, "value expected")
    }
}


func (self *luaState) CheckType(arg int, t LuaType) {
    if self.Type(arg) != t {
        self.tagError(arg, t)
    }
}

func (self *luaState) CheckInteger(arg int) int64 {
    i, ok := self.ToIntegerX(arg)
    if !ok {
        self.intError(arg)
    }
    return i
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checknumber
func (self *luaState) CheckNumber(arg int) float64 {
    f, ok := self.ToNumberX(arg)
    if !ok {
        self.tagError(arg, LUA_TNUMBER)
    }
    return f
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkstring
// http://www.lua.org/manual/5.3/manual.html#luaL_checklstring
func (self *luaState) CheckString(arg int) string {
    s, ok := self.ToStringX(arg)
    if !ok {
        self.tagError(arg, LUA_TSTRING)
    }
    return s
}














