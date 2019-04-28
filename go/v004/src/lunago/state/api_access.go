/*
* @Author: konyka
* @Date:   2019-04-28 22:39:58
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 23:22:38
*/

package state

import "fmt"
import . "luago/api"

/**
 * TypeName()方法不需要读取任何栈数据，只是把给定的lua类型转换为对应的字符串表示
 */
func (self *luaState) TypeName(tp LuaType) string {
    switch tp {
    case LUA_TNONE:
        return "no value"
    case LUA_TNIL:
        return "nil"
    case LUA_TBOOLEAN:
        return "boolean"
    case LUA_TNUMBER:
        return "number"
    case LUA_TSTRING:
        return "string"
    case LUA_TTABLE:
        return "table"
    case LUA_TFUNCTION:
        return "function"
    case LUA_TTHREAD:
        return "thread"
    default:
        return "userdata"
    }
}

/**
 * Type()根据索引返回值的类型，如果索引无效，则返回LUA_TNONE.
 */
func (self *luaState) Type(idx int) LuaType {
    if self.stack.isValid(idx) {
        val := self.stack.get(idx)
        return typeOf(val)
    }
    return LUA_TNONE
}

/**
 * 判断是否是none
 */
func (self *luaState) IsNone(idx int) bool {
    return self.Type(idx) == LUA_TNONE
}

/**
 * 判断是否是nil
 */
func (self *luaState) IsNil(idx int) bool {
    return self.Type(idx) == LUA_TNIL
}

/**
 * 判断是否是none 或者 nil
 */
func (self *luaState) IsNoneOrNil(idx int) bool {
    return self.Type(idx) <= LUA_TNIL
}

/**
 * 判断是否是boolean 
 */
func (self *luaState) IsBoolean(idx int) bool {
    return self.Type(idx) == LUA_TBOOLEAN
}

/**
 * IsString()判断指定索引处的值是不是字符串或者数字。
 */
func (self *luaState) IsString(idx int) bool {
    t := self.Type(idx)
    return t == LUA_TSTRING || t == LUA_TNUMBER
}

/**
 * IsNumber()方法判断给定随你处的值是不是数字类型，如果可以转化为数字类型也可以。
 */
func (self *luaState) IsNumber(idx int) bool {
    _, ok := self.ToNumberX(idx)
    return ok
}




