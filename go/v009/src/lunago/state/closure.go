/*
* @Author: konyka
* @Date:   2019-04-30 17:12:09
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 21:11:15
*/


package state

import "lunago/binchunk"
import . "lunago/api"

type closure struct {
    proto *binchunk.Prototype   // lua closure
    goFunc GoFunction          // go closure
}

/**
 * 创建lua闭包
 */
func newLuaClosure(proto *binchunk.Prototype) *closure {
    return &closure{proto: proto}
}

/**
 * 创建go闭包的函数
 */
func newGoClosure(f GoFunction) *closure {
    return &closure{goFunc: f}
}



