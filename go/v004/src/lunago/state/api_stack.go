/*
* @Author: konyka
* @Date:   2019-04-28 11:43:36
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 22:25:01
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

/**
 * [ Insert()方法将栈顶的值弹出，然后将其值插入到指定的位置。
 * 原来idx以及之后的值则分别向上移动一个位置。]
 * @Author   konyka
 * @DateTime 2019-04-28T13:42:32+0800
 * @param    {[type]}                 self *luaState)    Insert(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Insert(idx int) {
    self.Rotate(idx, 1)
}

/**
 * [旋转操作。Rotate(idx, n int) 将[idx, top] 索引区间内的值朝着栈顶方向旋转 n 个位置。
 * 如果n是负数，那么实际的效果就是朝着栈底方向旋转。 ]
 * @Author   konyka
 * @DateTime 2019-04-28T13:49:09+0800
 * @param    {[type]}                 self *luaState)    Rotate(idx, n int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Rotate(idx, n int) {
    t := self.stack.top - 1           /* end of stack segment being rotated */
    p := self.stack.absIndex(idx) - 1 /* start of segment */
    var m int                         /* end of prefix */
    if n >= 0 {
        m = t - n
    } else {
        m = p - n - 1
    }
    self.stack.reverse(p, m)   /* reverse the prefix with length 'n' */
    self.stack.reverse(m+1, t) /* reverse the suffix */
    self.stack.reverse(p, t)   /* reverse the entire segment */
}

/**
 * [ Remove() 删除置顶索引处的值，然后将该值上面的所有值全部向下移动一个位置。]
 * @Author   konyka
 * @DateTime 2019-04-28T13:56:31+0800
 * @param    {[type]}                 self *luaState)    Remove(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Remove(idx int) {
    self.Rotate(idx, -1)
    self.Pop(1)
}

/**
 * [ SetTop()将栈顶索引设置为指定的值。如果指定的值小于当前栈顶的索引，效果则相当于弹出操作，指定值0
 * 如果指定的值 n 大于当前栈顶的索引，则效果相当于push （n - 栈顶索引） 个nil值。
 * SetTop()根据不同的情况执行push 、pop操作。]
 * @Author   konyka
 * @DateTime 2019-04-28T22:24:10+0800
 * @param    {[type]}                 self *luaState)    SetTop(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) SetTop(idx int) {
    newTop := self.stack.absIndex(idx)
    if newTop < 0 {
        panic("stack underflow!")
    }

    n := self.stack.top - newTop
    if n > 0 {
        for i := 0; i < n; i++ {
            self.stack.pop()
        }
    } else if n < 0 {
        for i := 0; i > n; i-- {
            self.stack.push(nil)
        }
    }
}


















