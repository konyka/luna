/*
* @Author: konyka
* @Date:   2019-04-30 10:55:53
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 11:07:50
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





