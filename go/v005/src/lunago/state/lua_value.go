/*
* @Author: konyka
* @Date:   2019-04-27 18:15:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 10:39:59
*/

package state

import . "lunago/api"

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64, float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	default:
		panic("todo!")
	}
}

/**
 * 在Lua中，只有nil、false表示假，其他都表示真。lua_value.go定义convertToBoolean。
 */
func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}








