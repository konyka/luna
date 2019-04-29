/*
* @Author: konyka
* @Date:   2019-04-28 22:39:58
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 08:58:58
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

/**
 * IsInteger()判断指定索引处的值是不是整数类型。
 */
func (self *luaState) IsInteger(idx int) bool {
    val := self.stack.get(idx)
    _, ok := val.(int64)
    return ok
}

/**
 * ToBoolean()从指定的索引处取出一个boolean值，如果值不是布尔类型，则需要进行类型转换。
 */
func (self *luaState) ToBoolean(idx int) bool {
    val := self.stack.get(idx)
    return convertToBoolean(val)
}

/**
 * ToNumber()：如果值不是数字类型，并且也没有办法转换成数字类型，返回0
 */
func (self *luaState) ToNumber(idx int) float64 {
    n, _ := self.ToNumberX(idx)
    return n
}

/**
 * [ ToNumberX()：如果值不是数字类型，并且也没有办法转换成数字类型，则会报告转换是否成功。]
 * @Author   konyka
 * @DateTime 2019-04-29T08:48:02+0800
 * @param    {[type]}                 self *luaState)    ToNumberX(idx int) (float64, bool [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) ToNumberX(idx int) (float64, bool) {
    val := self.stack.get(idx)
    switch x := val.(type) {
    case float64:
        return x, true
    case int64:
        return float64(x), true
    default:
        return 0, false
    }
}


/**
 * ToInteger()：如果值不是整数类型，并且也没有办法转换成整数类型，返回0.
 */
func (self *luaState) ToInteger(idx int) int64 {
    i, _ := self.ToIntegerX(idx)
    return i
}

/**
 * [ToIntegerX()：如果值不是整数类型，并且也没有办法转换成整数类型，则会报告转换是否成功。]
 * @Author   konyka
 * @DateTime 2019-04-29T08:58:45+0800
 * @param    {[type]}                 self *luaState)    ToIntegerX(idx int) (int64, bool [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) ToIntegerX(idx int) (int64, bool) {
    val := self.stack.get(idx)
    i, ok := val.(int64)
    return i, ok
}


















