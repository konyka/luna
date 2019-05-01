/*
* @Author: konyka
* @Date:   2019-04-28 11:24:28
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 16:09:14
*/


package state

import . "lunago/api"

/**
 * lusState状态机 结构体定义
 * 
 */
type luaState struct {
    registry *luaTable
    stack *luaStack

}

/**
 * 用来创建luaState的实例
 * 给New函数增加了两个参数，第一个参数用于指定Lua栈的初始容量，第二个参数传入函数原型，以初始化proto字段。
 * 由于虚拟机肯定是从第一条指令开始执行的，因此pc字段初始化为0就可以了。
 */
func New() *luaState {
    registry := newLuaTable(0, 0)
    registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))

    ls := &luaState{registry: registry}
    ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
    return ls
}

func (self *luaState) pushLuaStack(stack *luaStack) {
    stack.prev = self.stack
    self.stack = stack
}

func (self *luaState) popLuaStack() {
    stack := self.stack
    self.stack = stack.prev
    stack.prev = nil
}






