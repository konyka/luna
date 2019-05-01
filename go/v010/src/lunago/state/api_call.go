/*
* @Author: konyka
* @Date:   2019-04-30 18:39:45
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 21:17:11
*/

package state

import "lunago/binchunk"
import "lunago/vm"
import . "lunago/api"


/**
 * 加载chunk
 */
func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
    proto := binchunk.Undump(chunk) // todo
    c := newLuaClosure(proto)
    self.stack.push(c)
    return 0
}

/**
 * [func description]
 * Call会调用Lua函数。在执行Call之前，必须先把被调用的函数push到栈顶，然后把参数一次push到栈顶，
 * Call（）完成后，
 * 参数值和函数会被弹出栈顶，取而代之的是指定数量的返回值。Call方法接收两个参数：
 * 第一个参数指定准备传递给被调函数的参数数量，同时也隐含给出了被调函数在栈中的位置；
 * 第二个参数指定需要的返回值的数量（多退少补），如果是-1，则被调函数的返回值会全部留在栈顶。
 * @Author   konyka
 * @DateTime 2019-04-30T18:51:05+0800
 * @param    {[type]}                 self *luaState)    Call(nArgs, nResults int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Call(nArgs, nResults int) {
    val := self.stack.get(-(nArgs + 1))
    if c, ok := val.(*closure); ok {
        if c.proto != nil {
            self.callLuaClosure(nArgs, nResults, c)
        } else {
            self.callGoClosure(nArgs, nResults, c)
        }
    } else {
        panic("not function!")
    }
}

func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
    nRegs := int(c.proto.MaxStackSize)
    nParams := int(c.proto.NumParams)
    isVararg := c.proto.IsVararg == 1

    // create new lua stack
    newStack := newLuaStack(nRegs+LUA_MINSTACK, self)
    newStack.closure = c

    // pass args, pop func
    funcAndArgs := self.stack.popN(nArgs + 1)
    newStack.pushN(funcAndArgs[1:], nParams)
    newStack.top = nRegs
    if nArgs > nParams && isVararg {
        newStack.varargs = funcAndArgs[nParams+1:]
    }

    // run closure
    self.pushLuaStack(newStack)
    self.runLuaClosure()
    self.popLuaStack()

    // return results
    if nResults != 0 {
        results := newStack.popN(newStack.top - nRegs)
        self.stack.check(len(results))
        self.stack.pushN(results, nResults)
    }
}



func (self *luaState) runLuaClosure() {
    for {
        inst := vm.Instruction(self.Fetch())
        inst.Execute(self)
        if inst.Opcode() == vm.OP_RETURN {
            break
        }
    }
}


func (self *luaState) callGoClosure(nArgs, nResults int, c *closure) {
    // create new lua stack
    newStack := newLuaStack(nArgs+LUA_MINSTACK, self)
    newStack.closure = c

    // pass args, pop func
    if nArgs > 0 {
        args := self.stack.popN(nArgs)
        newStack.pushN(args, nArgs)
    }
    self.stack.pop()

    // run closure
    self.pushLuaStack(newStack)
    r := c.goFunc(self)
    self.popLuaStack()

    // return results
    if nResults != 0 {
        results := newStack.popN(r)
        self.stack.check(len(results))
        self.stack.pushN(results, nResults)
    }
}







