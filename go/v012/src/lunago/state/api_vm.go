/*
* @Author: konyka
* @Date:   2019-04-29 17:46:01
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-02 10:39:07
*/


package state

func (self *luaState) PC() int {
    return self.stack.pc
}

func (self *luaState) AddPC(n int) {
    self.stack.pc += n
}

/**
 * Fetch（）根据PC索引从函数原型的指令表中取出当前的指令，
 * 然后把PC+1，这样下次在调用该方法取出的就是下一条指令
 */
func (self *luaState) Fetch() uint32 {
    i := self.stack.closure.proto.Code[self.stack.pc]
    self.stack.pc++
    return i
}

/**
 * [func GetConst()根据索引从函数原型的常量表中取出一个常量值，然后将其push到栈顶]
 * @Author   konyka
 * @DateTime 2019-04-29T17:56:37+0800
 * @param    {[type]}                 self *luaState)    GetConst(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) GetConst(idx int) {
    c := self.stack.closure.proto.Constants[idx]
    self.stack.push(c)
}

/**
 * [func GetRK(）根据情况调用GetConst（）把某个常量push到栈顶，
 * 或者调用PushValue（）把某个索引处的栈值push到栈顶。]
 * @Author   konyka
 * @DateTime 2019-04-29T18:03:05+0800
 * @param    {[type]}                 self *luaState)    GetRK(rk int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) GetRK(rk int) {
    if rk > 0xFF { // constant
        self.GetConst(rk & 0xFF)
    } else { // register
        self.PushValue(rk + 1)
    }
}

/**
 *  RegisterCount() 当前lua函数所操作的寄存器计数器
 */
func (self *luaState) RegisterCount() int {
    return int(self.stack.closure.proto.MaxStackSize)
}

/**
 * [func 把传递给当前lua函数的变长参数push到栈顶 多退少补]
 * @Author   konyka
 * @DateTime 2019-05-01T10:36:48+0800
 * @param    {[type]}                 self *luaState)    LoadVararg(n int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) LoadVararg(n int) {
    if n < 0 {
        n = len(self.stack.varargs)
    }

    self.stack.check(n)
    self.stack.pushN(self.stack.varargs, n)
}

/**
 * [func LoadProto(idx int)  把当前lua函数的子函数的原型 实例化为闭包 ，并push到栈顶]
 * @Author   konyka
 * @DateTime 2019-05-01T10:38:43+0800
 * @param    {[type]}                 self *luaState)    LoadProto(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) LoadProto(idx int) {
    stack := self.stack
    subProto := stack.closure.proto.Protos[idx]
    closure := newLuaClosure(subProto)
    stack.push(closure)

    for i, uvInfo := range subProto.Upvalues {
        uvIdx := int(uvInfo.Idx)
        if uvInfo.Instack == 1 {
            if stack.openuvs == nil {
                stack.openuvs = map[int]*upvalue{}
            }

            if openuv, found := stack.openuvs[uvIdx]; found {
                closure.upvals[i] = openuv
            } else {
                closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
                stack.openuvs[uvIdx] = closure.upvals[i]
            }
        } else {
            closure.upvals[i] = stack.closure.upvals[uvIdx]
        }
    }
}

func (self *luaState) CloseUpvalues(a int) {
    for i, openuv := range self.stack.openuvs {
        if i >= a-1 {
            val := *openuv.val
            openuv.val = &val
            delete(self.stack.openuvs, i)
        }
    }
}