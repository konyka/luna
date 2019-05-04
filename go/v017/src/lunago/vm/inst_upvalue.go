/*
* @Author: konyka
* @Date:   2019-05-01 19:37:07
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-02 11:28:01
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
 * [setUpval UpValue[B] := R(A)]
 * setupval指令（iABC），使用寄存器中的值给当前闭包的Upvalue赋值。
 *  其中仅存起索引由操作数A指定，Upvalue索引由操作数B指定，操作数C没有用到。
 * @Author   konyka
 * @DateTime 2019-05-02T09:59:41+0800
 * @param    {[type]}                 i  Instruction [description]
 * @param    {[type]}                 vm LuaVM       [description]
 */
func setUpval(i Instruction, vm LuaVM) {
    a, b, _ := i.ABC()
    a += 1
    b += 1

    vm.Copy(a, LuaUpvalueIndex(b))
}

 
/**
 * [getTabUp R(A) := UpValue[B][RK(C)]]
 *  如果当前闭包的某个Upvalue是表，则gettabup指令（iABC模式）可以根据key从该表里面取值，
 *  然后把value放到目标寄存器中。其中目标寄存器的索引由餐做数A指定，
 *  Upvalue的索引由操作数B指定，key（可能在寄存器中，也可能在常量表中）索引由操作数C指定。
 * @Author   konyka
 * @DateTime 2019-05-02T10:07:23+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func getTabUp(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1
    b += 1

    vm.GetRK(c)
    vm.GetTable(LuaUpvalueIndex(b))
    vm.Replace(a)
}
 
/**
 * [setTabUp UpValue[A][RK(B)] := RK(C)]
 * @Author   konyka
 * @DateTime 2019-05-02T10:16:16+0800
 * @param    {[type]}                 i  Instruction [description]
 * @param    {[type]}                 vm LuaVM       [description]
 */
func setTabUp(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1

    vm.GetRK(b)
    vm.GetRK(c)
    vm.SetTable(LuaUpvalueIndex(a))
}






