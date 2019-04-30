/*
* @Author: konyka
* @Date:   2019-04-30 17:12:09
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 17:22:52
*/


package state

import "lunago/binchunk"

type closure struct {
    proto *binchunk.Prototype
}

/**
 * 创建闭包
 */
func newLuaClosure(proto *binchunk.Prototype) *closure {
    return &closure{proto: proto}
}





