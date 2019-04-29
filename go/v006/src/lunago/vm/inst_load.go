/*
* @Author: konyka
* @Date:   2019-04-29 19:13:25
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 19:54:47
*/

package vm

import . "lunago/api"


/**
 * [loadNil R(A), R(A+1), ..., R(A+B) := nil
 * loadnil指令（iABC模式）用于给连续n个寄存器放置nil值。
 * 寄存器的起始索引由操作数A指定，寄存器数量由操作数B指定，操作数C没有使用.
 * 在lua代码里，局部变量的默认初始值就是nil。loadnil指令常用于给连续的n个局部变量设置初始值。]
 * @Author   konyka
 * @DateTime 2019-04-29T19:14:17+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func loadNil(i Instruction, vm LuaVM) {
    a, b, _ := i.ABC()
    a += 1

    vm.PushNil()
    for i := a; i <= a+b; i++ {
        vm.Copy(-1, i)
    }
    vm.Pop(1)
}


/**
 * [loadBool R(A) := (bool)B; if (C) pc++
 * loadbool指令（iABC指令）给耽搁寄存器设置布尔值。寄存器索引由操作数A指定，
 * 布尔值由操作数B指定（0表示false，非0表示true），如果寄存器C非0，则跳过下一条指令。]
 * @Author   konyka
 * @DateTime 2019-04-29T19:41:05+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func loadBool(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1

    vm.PushBoolean(b != 0)
    vm.Replace(a)

    if c != 0 {
        vm.AddPC(1)
    }
}

/**
 * [loadK ：R(A) := Kst(Bx)
 * loadk(iABx模式)将常量表里面的某个常量加载到指定的寄存器中，
 * 寄存器的索引由操作数A指定，常量表的索引由操作数Bx指定。
 * ]
 * @Author   konyka
 * @DateTime 2019-04-29T19:53:57+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func loadK(i Instruction, vm LuaVM) {
    a, bx := i.ABx()
    a += 1

    vm.GetConst(bx)
    vm.Replace(a)
}







