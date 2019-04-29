/*
* @Author: konyka
* @Date:   2019-04-29 14:58:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 15:46:15
*/


package state

import . "lunago/api"

/**
 * 对指定索引处的两个值进行比较，返回结果。此函数不会改变栈的状态.
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


/**
 * 用于比较两个值是否相等
 */
func _eq(a, b luaValue) bool {
    switch x := a.(type) {
    case nil:
        return b == nil
    case bool:
        y, ok := b.(bool)
        return ok && x == y
    case string:
        y, ok := b.(string)
        return ok && x == y
    case int64:
        switch y := b.(type) {
        case int64:
            return x == y
        case float64:
            return float64(x) == y
        default:
            return false
        }
    case float64:
        switch y := b.(type) {
        case float64:
            return x == y
        case int64:
            return x == float64(y)
        default:
            return false
        }
    default:
        return a == b
    }
}

/**
 * 小于操作 <
 */
func _lt(a, b luaValue) bool {
    switch x := a.(type) {
    case string:
        if y, ok := b.(string); ok {
            return x < y
        }
    case int64:
        switch y := b.(type) {
        case int64:
            return x < y
        case float64:
            return float64(x) < y
        }
    case float64:
        switch y := b.(type) {
        case float64:
            return x < y
        case int64:
            return x < float64(y)
        }
    }
    panic("comparison error!")
}

/**
 * <=
 */
func _le(a, b luaValue) bool {
    switch x := a.(type) {
    case string:
        if y, ok := b.(string); ok {
            return x <= y
        }
    case int64:
        switch y := b.(type) {
        case int64:
            return x <= y
        case float64:
            return float64(x) <= y
        }
    case float64:
        switch y := b.(type) {
        case float64:
            return x <= y
        case int64:
            return x <= float64(y)
        }
    }
    panic("comparison error!")
}


