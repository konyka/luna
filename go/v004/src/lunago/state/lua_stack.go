/*
* @Author: konyka
* @Date:   2019-04-27 18:15:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 09:20:37
*/

package state

type luaStack struct {
	slots	[]luaValue	//用来存放值
	top		int 		//记录栈顶的索引
}

/**
 * 创建指定容量的栈
 */
func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

/**
 * [check()方法检查栈的空闲空间是否还可以容纳（push）至少n个值，如果不满足这个条件，就会调用go的
    append()函数对其进行扩容。]
 * @Author   konyka
 * @DateTime 2019-04-28T09:11:03+0800
 * @param    {[type]}                 self *luaStack)    check(n int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaStack) check(n int) {
	free := len(self.slots) - self.top
	for i := free; i < n; i++ {
		self.slots = append(self.slots, nil)
	}
}

/**
 * [push()方法用来将值push到栈顶，如果溢出，就先暂时调用panic（）终止程序的执行。
]
 * @Author   konyka
 * @DateTime 2019-04-28T09:20:13+0800
 * @param    {[type]}                 self *luaStack)    push(val luaValue [description]
 * @return   {[type]}                      [description]
 */
func (self *luaStack) push(val luaValue) {
	if self.top == len(self.slots) {
		panic("stack overflow!")
	}
	self.slots[self.top] = val
	self.top++
}













