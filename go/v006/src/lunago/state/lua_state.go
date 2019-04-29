/*
* @Author: konyka
* @Date:   2019-04-28 11:24:28
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 17:23:50
*/


package state

import "lunago/binchunk"
/**
 * lusState状态机 结构体定义
 * 
 */
type luaState struct {
    stack *luaStack
    proto *binchunk.Prototype //保存函数原型
    pc    int //程序计数器
}

/**
 * 用来创建luaState的实例
 */
func New() *luaState {
    return &luaState{
        stack: newLuaStack(20),
    }
}








