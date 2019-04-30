/*
* @Author: konyka
* @Date:   2019-04-30 10:55:53
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 11:35:09
*/


package state

import "math"
import "lunago/number"

type luaTable struct {
    arr  []luaValue
    _map map[luaValue]luaValue
}

func newLuaTable(nArr, nRec int) *luaTable {
    t := &luaTable{}
    if nArr > 0 {
        t.arr = make([]luaValue, 0, nArr)
    }
    if nRec > 0 {
        t._map = make(map[luaValue]luaValue, nRec)
    }
    return t
}

/**
 * get（）方法根据key从表里面查找值。
 */
func (self *luaTable) get(key luaValue) luaValue {
    key = _floatToInteger(key)
    if idx, ok := key.(int64); ok {
        if idx >= 1 && idx <= int64(len(self.arr)) {
            return self.arr[idx-1]
        }
    }
    return self._map[key]
}


/**
 * _floatToInteger()尝试吧服点类型的key转换为整数
 */
func _floatToInteger(key luaValue) luaValue {
    if f, ok := key.(float64); ok {
        if i, ok := number.FloatToInteger(f); ok {
            return i
        }
    }
    return key
}

/**
 * [func  put()方法向表里保存键值对。]
 * @Author   konyka
 * @DateTime 2019-04-30T11:26:14+0800
 * @param    {[type]}                 self *luaTable)    put(key, val luaValue [description]
 * @return   {[type]}                      [description]
 */
func (self *luaTable) put(key, val luaValue) {
    if key == nil {
        panic("table index is nil!")
    }
    if f, ok := key.(float64); ok && math.IsNaN(f) {
        panic("table index is NaN!")
    }

    key = _floatToInteger(key)
    if idx, ok := key.(int64); ok && idx >= 1 {
        arrLen := int64(len(self.arr))
        if idx <= arrLen {
            self.arr[idx-1] = val
            if idx == arrLen && val == nil {
                self._shrinkArray()
            }
            return
        }
        if idx == arrLen+1 {
            delete(self._map, key)
            if val != nil {
                self.arr = append(self.arr, val)
                self._expandArray()
            }
            return
        }
    }
    if val != nil {
        if self._map == nil {
            self._map = make(map[luaValue]luaValue, 8)
        }
        self._map[key] = val
    } else {
        delete(self._map, key)
    }
}


/**
 * [func 把尾部的hole全部删除]
 * @Author   konyka
 * @DateTime 2019-04-30T11:34:00+0800
 * @param    {[type]}                 self *luaTable)    _shrinkArray( [description]
 * @return   {[type]}                      [description]
 */
func (self *luaTable) _shrinkArray() {
    for i := len(self.arr) - 1; i >= 0; i-- {
        if self.arr[i] == nil {
            self.arr = self.arr[0:i]
        }
    }
}








