/*
* @Author: konyka
* @Date:   2019-04-28 11:43:36
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 13:18:18
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




