/*
* @Author: konyka
* @Date:   2019-04-27 18:15:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-02 13:52:06
*/

package state

import . "lunago/api"
import "lunago/number"

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
	case *luaTable:
		return LUA_TTABLE
	case *closure:
		return LUA_TFUNCTION
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

/**
 * [convertToInteger 任意值转化为整数]
 * @Author   konyka
 * @DateTime 2019-04-29T13:05:52+0800
 * @param    {[type]}                 val luaValue)     (int64, bool [description]
 * @return   {[type]}                     [description]
 */
func convertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		return _stringToInteger(x)
	default:
		return 0, false
	}
}


/**
 * [ _stringToInteger  对于浮点数，可以调用之前定义的FLoatToInteger()方法将其转换为整数，
 * 对于字符串，可以先试试能够直接解析为整数，如果不能，在尝试将其解析为浮点数，然后转换为整数。]
 * @Author   konyka
 * @DateTime 2019-04-29T13:08:42+0800
 * @param    {[type]}                 s string)       (int64, bool [description]
 * @return   {[type]}                   [description]
 */
func _stringToInteger(s string) (int64, bool) {
	if i, ok := number.ParseInteger(s); ok {
		return i, true
	}
	if f, ok := number.ParseFloat(s); ok {
		return number.FloatToInteger(f)
	}
	return 0, false
}
/* metatable */

/**
 * [setMetatable setMetatable 用来给值关联元表 ]
 * @Author   konyka
 * @DateTime 2019-05-02T13:20:03+0800
 * @param    {[type]}                 val luaValue  [description]
 * @param    {[type]}                 mt  *luaTable [description]
 * @param    {[type]}                 ls  *luaState [description]
 */
func setMetatable(val luaValue, mt *luaTable, ls *luaState) {
	if t, ok := val.(*luaTable); ok {
		t.metatable = mt
		return
	}
	key := fmt.Sprintf("_MT%d", typeOf(val))
	ls.registry.put(key, mt)
}


func getMetatable(val luaValue, ls *luaState) *luaTable {
	if t, ok := val.(*luaTable); ok {
		return t.metatable
	}
	key := fmt.Sprintf("_MT%d", typeOf(val))
	if mt := ls.registry.get(key); mt != nil {
		return mt.(*luaTable)
	}
	return nil
}

/**
 * [callMetamethod callMetamethod方法负责查找并调用元方法 ]
 * @Author   konyka
 * @DateTime 2019-05-02T13:44:18+0800
 * @param    {[type]}                 a      [description]
 * @param    {[type]}                 b      luaValue      [description]
 * @param    {[type]}                 mmName string        [description]
 * @param    {[type]}                 ls     *luaState)    (luaValue,    bool [description]
 * @return   {[type]}                        [description]
 */
func callMetamethod(a, b luaValue, mmName string, ls *luaState) (luaValue, bool) {
	var mm luaValue
	if mm = getMetafield(a, mmName, ls); mm == nil {
		if mm = getMetafield(b, mmName, ls); mm == nil {
			return nil, false
		}
	}

	ls.stack.check(4)
	ls.stack.push(mm)
	ls.stack.push(a)
	ls.stack.push(b)
	ls.Call(2, 1)
	return ls.stack.pop(), true
}


func getMetafield(val luaValue, fieldName string, ls *luaState) luaValue {
	if mt := getMetatable(val, ls); mt != nil {
		return mt.get(fieldName)
	}
	return nil
}

