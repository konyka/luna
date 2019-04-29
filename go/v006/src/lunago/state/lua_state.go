/*
* @Author: konyka
* @Date:   2019-04-28 11:24:28
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 17:39:03
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
 * 给New函数增加了两个参数，第一个参数用于指定Lua栈的初始容量，第二个参数传入函数原型，以初始化proto字段。
 * 由于虚拟机肯定是从第一条指令开始执行的，因此pc字段初始化为0就可以了。
 */
func New(stackSize int, proto *binchunk.Prototype) *luaState {
    return &luaState{
        stack: newLuaStack(stackSize),
        proto: proto,
        pc:    0,
    }
}








