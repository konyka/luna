/*
* @Author: konyka
* @Date:   2019-04-28 22:33:14
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 19:25:49
*/

package state

import "fmt"
import . "lunago/api"

/**
 * [push nil to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:35:11+0800
 * @param    {[type]}                 self *luaState)    PushNil( [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushNil() {
    self.stack.push(nil)
}

/**
 * [push boolean to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:35:43+0800
 * @param    {[type]}                 self *luaState)    PushBoolean(b bool [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushBoolean(b bool) {
    self.stack.push(b)
}

/**
 * [push integer to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:36:07+0800
 * @param    {[type]}                 self *luaState)    PushInteger(n int64 [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushInteger(n int64) {
    self.stack.push(n)
}

/**
 * [push number to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:36:17+0800
 * @param    {[type]}                 self *luaState)    PushNumber(n float64 [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushNumber(n float64) {
    self.stack.push(n)
}

/**
 * [push string to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:36:26+0800
 * @param    {[type]}                 self *luaState)    PushString(s string [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushString(s string) {
    self.stack.push(s)
}

/**
 * [func 接受一个go函数参数，把它转换为go闭包，然后push到栈顶。]
 * @Author   konyka
 * @DateTime 2019-05-01T15:36:48+0800
 * @param    {[type]}                 self *luaState)    PushGoFunction(f GoFunction [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushGoFunction(f GoFunction) {
    self.stack.push(newGoClosure(f, 0)) //第二个参数传入0
}

/**
 * [func 把全局环境push到栈顶]
 * @Author   konyka
 * @DateTime 2019-05-01T17:28:05+0800
 * @param    {[type]}                 self *luaState)    PushGlobalTable( [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushGlobalTable() {
    global := self.registry.get(LUA_RIDX_GLOBALS)
    self.stack.push(global)
}

/**
 * [func 这个函数和PushGoFunction差不多，把go函数转换成go闭包push到栈顶，
 * 区别是PushGoClosure先从栈顶弹出n个lua值，这些值会成为go闭包的Upvalue。]
 * @Author   konyka
 * @DateTime 2019-05-02T08:52:43+0800
 * @param    {[type]}                 self *luaState)    PushGoClosure(f GoFunction, n int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushGoClosure(f GoFunction, n int) {
    closure := newGoClosure(f, n)
    for i := n; i > 0; i-- {
        val := self.stack.pop()
        closure.upvals[i-1] = &upvalue{&val}
    }
    self.stack.push(closure)
}

func (self *luaState) PushFString(fmtStr string, a ...interface{}) {
    str := fmt.Sprintf(fmtStr, a...)
    self.stack.push(str)
}

/**
 * PushThread()将线程push到栈顶，返回的布尔值表示线程是不是为主线程。
 */
func (self *luaState) PushThread() bool {
    self.stack.push(self)
    return self.isMainThread()
}




