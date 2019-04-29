/*
* @Author: konyka
* @Date:   2019-04-27 18:15:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 12:49:01
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

/**
 * [convertToFloat description]
 * @Author   konyka
 * @DateTime 2019-04-29T12:49:00+0800
 * @param    {[type]}                 val luaValue)     (float64, bool [description]
 * @return   {[type]}                     [description]
 */
func convertToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case int64:
		return float64(x), true
	case float64:
		return x, true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}






