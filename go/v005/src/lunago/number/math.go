/*
* @Author: konyka
* @Date:   2019-04-29 11:26:48
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 11:43:43
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












