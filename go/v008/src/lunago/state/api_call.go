/*
* @Author: konyka
* @Date:   2019-04-30 18:39:45
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 18:51:23
*/

package state

import "fmt"
import "luago/binchunk"
import "luago/vm"

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
        fmt.Printf("call %s<%d,%d>\n", c.proto.Source,
            c.proto.LineDefined, c.proto.LastLineDefined)
        self.callLuaClosure(nArgs, nResults, c)
    } else {
        panic("not function!")
    }
}










