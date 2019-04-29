/*
* @Author: konyka
* @Date:   2019-04-29 11:26:48
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 12:32:13
*/

package number

import "math"
/**
 * 整除函数
 */
func IFloorDiv(a, b int64) int64 {
    if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
        return a / b
    } else {
        return a/b - 1
    }
}

/**
 * 整除函数
 */
func FFloorDiv(a, b float64) float64 {
    return math.Floor(a / b)
}

/**
 * a % b == a - ((a // b) * b)
 */
func IMod(a, b int64) int64 {
    return a - IFloorDiv(a, b)*b
}

/**
 * a % b == a - ((a // b) * b)
 */
func FMod(a, b float64) float64 {
    if a > 0 && math.IsInf(b, 1) || a < 0 && math.IsInf(b, -1) {
        return a
    }
    if a > 0 && math.IsInf(b, -1) || a < 0 && math.IsInf(b, 1) {
        return b
    }
    return a - math.Floor(a/b)*b
}

/**
 * << 左移
 * 因为go里面的位移运算符右边的操作时只能是无符号整数，因此在第一个分支里面对位移的数进行了类型转换。
 */
func ShiftLeft(a, n int64) int64 {
    if n >= 0 {
        return a << uint64(n)
    } else {
        return ShiftRight(a, -n)
    }
}

/**
 * >> 右移
 *  go中，如果右移运算符的左操作数是有符号整数，那么进行的就是有符号右移，空位补充1.
 *  不过我们期望的是无符号右移，空位补充0，所以在第一个分支里面需要先将左操作数转换成无符号整数
 *  在执行右移擦欧总，然后在将结果转换为有符号整数。如果要移动的位数小于0，则将右移转换为左移。
 */
func ShiftRight(a, n int64) int64 {
    if n >= 0 {
        return int64(uint64(a) >> uint64(n))
    } else {
        return ShiftLeft(a, -n)
    }
}

// todo: ？？？？？correct?
/**
 * [Float To Integer ]
 * @Author   konyka
 * @DateTime 2019-04-29T12:31:57+0800
 * @param    {[type]}                 f float64) (int64, bool [description]
 */
func FloatToInteger(f float64) (int64, bool) {
    i := int64(f)
    return i, float64(i) == f
}


