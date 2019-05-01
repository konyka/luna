/*
* @Author: konyka
* @Date:   2019-04-27 18:15:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 16:45:02
*/

package state

type luaStack struct {
	slots	[]luaValue	//用来存放值
	top		int 		//记录栈顶的索引
	/* linked list */
	prev *luaStack
	/* call info */
	state   *luaState
	closure *closure
	varargs []luaValue
	pc      int

}

/**
 * 创建指定容量的栈
 */
func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: state,
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
 * [push()方法用来将值push到栈顶，如果溢出，就先暂时调用panic（）终止程序的执行。]
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
/**
 * pop()方法从栈顶弹出一个值，如果栈是空的，则调用panic()终止程序
 */
func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		panic("stack underflow!")
	}
	self.top--
	val := self.slots[self.top]
	self.slots[self.top] = nil
	return val
}

/**
 * absIndex()方法吧索引切换成绝对索引--并没有考虑索引是否有效
 */
func (self *luaStack) absIndex(idx int) int {
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	return idx + self.top + 1
}


/**
 * isValid()判断所有是否有效
 */
func (self *luaStack) isValid(idx int) bool {
	if idx == LUA_REGISTRYINDEX {
		return true
	}
	absIdx := self.absIndex(idx)
	return absIdx > 0 && absIdx <= self.top
}

/**
 * get()根据索引从栈里面取值，如果索引无效 返回nil
 */
func (self *luaStack) get(idx int) luaValue {
	if idx == LUA_REGISTRYINDEX {
		return self.state.registry
	}

	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		return self.slots[absIdx-1]
	}
	return nil
}

/**
 * [set()根据索引向栈里面写入值，如果索引无效，调用panic（）终止]
 * @Author   konyka
 * @DateTime 2019-04-28T11:03:08+0800
 * @param    {[type]}                 self *luaStack)    set(idx int, val luaValue [description]
 * @return   {[type]}                      [description]
 */
func (self *luaStack) set(idx int, val luaValue) {
	if idx == LUA_REGISTRYINDEX {
		self.state.registry = val.(*luaTable)
		return
	}

	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		self.slots[absIdx-1] = val
		return
	}
	panic("invalid index!")
}

/**
 * [ reverse()方法就是循环交换两个位置的值。 ]
 * @Author   konyka
 * @DateTime 2019-04-28T22:17:15+0800
 * @param    {[type]}                 self *luaStack)    reverse(from, to int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

/**
 * popN(n int)从栈顶一次性弹出多个值。
 */
func (self *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = self.pop()
	}
	return vals
}

/**
 * [func pushN()：向栈顶push多个值（多退少补）]
 * @Author   konyka
 * @DateTime 2019-04-30T19:14:15+0800
 * @param    {[type]}                 self *luaStack)    pushN(vals []luaValue, n int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}

	for i := 0; i < n; i++ {
		if i < nVals {
			self.push(vals[i])
		} else {
			self.push(nil)
		}
	}
}
