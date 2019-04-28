/*
* @Author: konyka
* @Date:   2019-04-28 22:39:58
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 22:42:40
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









