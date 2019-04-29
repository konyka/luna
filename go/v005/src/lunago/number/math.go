/*
* @Author: konyka
* @Date:   2019-04-29 11:26:48
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 11:56:32
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
 */
func ShiftRight(a, n int64) int64 {
    if n >= 0 {
        return int64(uint64(a) >> uint64(n))
    } else {
        return ShiftLeft(a, -n)
    }
}




