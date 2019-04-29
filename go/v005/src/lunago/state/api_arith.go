/*
* @Author: konyka
* @Date:   2019-04-29 13:40:41
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 14:29:21
*/


package state

import "math"
import . "luago/api"
import "luago/number"


var (
    iadd  = func(a, b int64) int64 { return a + b }
    fadd  = func(a, b float64) float64 { return a + b }
    isub  = func(a, b int64) int64 { return a - b }
    fsub  = func(a, b float64) float64 { return a - b }
    imul  = func(a, b int64) int64 { return a * b }
    fmul  = func(a, b float64) float64 { return a * b }
    imod  = number.IMod
    fmod  = number.FMod
    pow   = math.Pow
    div   = func(a, b float64) float64 { return a / b }
    iidiv = number.IFloorDiv
    fidiv = number.FFloorDiv
    band  = func(a, b int64) int64 { return a & b }
    bor   = func(a, b int64) int64 { return a | b }
    bxor  = func(a, b int64) int64 { return a ^ b }
    shl   = number.ShiftLeft
    shr   = number.ShiftRight
    iunm  = func(a, _ int64) int64 { return -a }
    funm  = func(a, _ float64) float64 { return -a }
    bnot  = func(a, _ int64) int64 { return ^a }
)


type operator struct {
    integerFunc func(int64, int64) int64
    floatFunc   func(float64, float64) float64
}
/**
 * [operators 定义一个slice，里面是各种运算，需要注意的是，要和前面定义的lua运算码常量的顺序要一致]
 * @type {Array}
 */
var operators = []operator{
    operator{iadd, fadd},
    operator{isub, fsub},
    operator{imul, fmul},
    operator{imod, fmod},
    operator{nil, pow},
    operator{nil, div},
    operator{iidiv, fidiv},
    operator{band, nil},
    operator{bor, nil},
    operator{bxor, nil},
    operator{shl, nil},
    operator{shr, nil},
    operator{iunm, funm},
    operator{bnot, nil},
}

/**
 * [func [-(2|1), +1, e]]
 * @Author   konyka
 * @DateTime 2019-04-29T14:28:40+0800
 * @param    {[type]}                 self *luaState)    Arith(op ArithOp [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Arith(op ArithOp) {
    var a, b luaValue // operands
    b = self.stack.pop()
    if op != LUA_OPUNM && op != LUA_OPBNOT {
        a = self.stack.pop()
    } else {
        a = b
    }

    operator := operators[op]
    if result := _arith(a, b, operator); result != nil {
        self.stack.push(result)
    } else {
        panic("arithmetic error!")
    }
}





