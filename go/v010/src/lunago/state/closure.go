/*
* @Author: konyka
* @Date:   2019-04-30 17:12:09
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 22:59:14
*/


package state

import "lunago/binchunk"
import . "lunago/api"

type upvalue struct {
    val *luaValue
}

type closure struct {
    proto  *binchunk.Prototype // lua closure
    goFunc GoFunction          // go closure
    upvals []*upvalue
}
/**
 * 创建lua闭包
 */
func newLuaClosure(proto *binchunk.Prototype) *closure {
    c := &closure{proto: proto}
    if nUpvals := len(proto.Upvalues); nUpvals > 0 {
        c.upvals = make([]*upvalue, nUpvals)
    }
    return c
}

/**
 * 创建go闭包的函数
 */
func newGoClosure(f GoFunction, nUpvals int) *closure {
    c := &closure{goFunc: f}
    if nUpvals > 0 {
        c.upvals = make([]*upvalue, nUpvals)
    }
    return c
}

