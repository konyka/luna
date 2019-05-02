/*
* @Author: konyka
* @Date:   2019-05-01 19:37:07
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-02 09:46:36
*/

package vm

import . "lunago/api"


/**
 * [getUpval R(A) := UpValue[B]]
 * getupval(iABC 模式)，把当前闭包的某个Upvale值复制到目标寄存器，
 * 其中目标寄存器的索引由操作数A指定，Upvalue索引由操作数B指定，操作数C没有使用。
 * @Author   konyka
 * @DateTime 2019-05-02T09:46:14+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func getUpval(i Instruction, vm LuaVM) {
    a, b, _ := i.ABC()
    a += 1
    b += 1

    vm.Copy(LuaUpvalueIndex(b), a)
}

/**
 * [getTabUp R(A) := UpValue[B][RK(C)]]
 * @Author   konyka
 * @DateTime 2019-05-01T19:41:34+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func getTabUp(i Instruction, vm LuaVM) {
    a, _, c := i.ABC()
    a += 1

    vm.PushGlobalTable()
    vm.GetRK(c)
    vm.GetTable(-2)
    vm.Replace(a)
    vm.Pop(1)
}








