/*
* @Author: konyka
* @Date:   2019-04-29 13:40:41
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 14:08:26
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










