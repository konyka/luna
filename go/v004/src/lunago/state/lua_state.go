/*
* @Author: konyka
* @Date:   2019-04-28 11:24:28
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 11:30:20
*/


package state
/**
 * lusState状态机 结构体定义
 * 
 */
type luaState struct {
    stack *luaStack
}

/**
 * 用来创建luaState的实例
 */
func New() *luaState {
    return &luaState{
        stack: newLuaStack(20),
    }
}








