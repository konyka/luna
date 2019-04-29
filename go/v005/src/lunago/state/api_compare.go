/*
* @Author: konyka
* @Date:   2019-04-29 14:58:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 14:59:02
*/


package state

import . "luago/api"

/**
 * 
 */
func (self *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
    if !self.stack.isValid(idx1) || !self.stack.isValid(idx2) {
        return false
    }

    a := self.stack.get(idx1)
    b := self.stack.get(idx2)
    switch op {
    case LUA_OPEQ:
        return _eq(a, b)
    case LUA_OPLT:
        return _lt(a, b)
    case LUA_OPLE:
        return _le(a, b)
    default:
        panic("invalid compare op!")
    }
}









