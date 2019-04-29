/*
* @Author: konyka
* @Date:   2019-04-29 16:42:33
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 16:43:49
*/

package api

type LuaVM interface {
    LuaState
    PC() int
    AddPC(n int)
    Fetch() uint32
    GetConst(idx int)
    GetRK(rk int)
}










