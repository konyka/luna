/*
* @Author: konyka
* @Date:   2019-04-28 11:43:36
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 13:41:34
*/


package state

/**
 * 返回栈顶索引
 */
func (self *luaState) GetTop() int {
    return self.stack.top
}

/**
 * AbsIndex(idx int)把索引转化为绝对索引。
 */
func (self *luaState) AbsIndex(idx int) int {
    return self.stack.absIndex(idx)
}

/**
 * CheckStack(n int) 检查栈中是否有 n 个剩余空间可用
 * lua栈的容量不会自动增长，使用者需要检查栈的剩余空间，看看是否可以push n 个值而不会溢出。
 * 如果剩余空间足够 或者扩容成功 返回true，否则返回false.
 * n 表示需要多少个剩余空间存放数据。
 * 
 */
func (self *luaState) CheckStack(n int) bool {
    self.stack.check(n)
    return true // ??? never fails
}

/**
 * [ Pop(n int) 方法从栈顶弹出n 个值。]
 * @Author   konyka
 * @DateTime 2019-04-28T13:27:39+0800
 * @param    {[type]}                 self *luaState)    Pop(n int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Pop(n int) {
    for i := 0; i < n; i++ {
        self.stack.pop()
    }
}
/**
 * [ Copy()方法把值从一个位置复制到另一个位置。 ]
 * @Author   konyka
 * @DateTime 2019-04-28T13:30:02+0800
 * @param    {[type]}                 self *luaState)    Copy(fromIdx, toIdx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Copy(fromIdx, toIdx int) {
    val := self.stack.get(fromIdx)
    self.stack.set(toIdx, val)
}

/**
 * [ PushValue()方法把指定索引处的值push到栈顶。 ]
 * @Author   konyka
 * @DateTime 2019-04-28T13:32:52+0800
 * @param    {[type]}                 self *luaState)    PushValue(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushValue(idx int) {
    val := self.stack.get(idx)
    self.stack.push(val)
}

/**
 * [ 将栈顶的值弹出，然后写入到指定的位置。]
 * @Author   konyka
 * @DateTime 2019-04-28T13:41:20+0800
 * @param    {[type]}                 self *luaState)    Replace(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Replace(idx int) {
    val := self.stack.pop()
    self.stack.set(idx, val)
}






