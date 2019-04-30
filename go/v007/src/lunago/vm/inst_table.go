/*
* @Author: konyka
* @Date:   2019-04-30 13:01:31
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 14:04:53
*/


package vm

import . "lunago/api"
import "lunago/number"


/**
 * [newTable newtable指令(iABC模式)创建空表，并将其放到指定的寄存器。
 * 寄存器索引由操作数A指定，表的初始数组容量和哈希表容量分别由操作数B、C指定
 * R(A) := {} (size = B,C) ]
 * @Author   konyka
 * @DateTime 2019-04-30T13:29:53+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func newTable(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1

    vm.CreateTable(Fb2int(b), Fb2int(c))
    vm.Replace(a)
}

// R(A) := R(B)[RK(C)]
/**
 * [getTable gettable（iABC模式）指令根据key从表中取值，并放到目标寄存器中。
 * 其中表位于寄存器中，索引有操作数B指定；
 * key可能位于寄存器中，也可能在常量表中，索引由操作数C指定；
 * 目标寄存器的索引则由操作数A指定。]
 * @Author   konyka
 * @DateTime 2019-04-30T13:49:20+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func getTable(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1
    b += 1

    vm.GetRK(c)
    vm.GetTable(b)
    vm.Replace(a)
}

/**
 * [setTable R(A)[RK(B)] := RK(C)
 * settable指令（iABC 模式）根据key向表面面赋值。其中表位于寄存器中，
 * 索引由操作数A指定；key 、value柯恩呢该位于寄存器中，也可能在常量表中，索引分别由操作数BC指定。]
 * @Author   konyka
 * @DateTime 2019-04-30T14:04:13+0800
 * @param    {[type]}                 i  Instruction [description]
 * @param    {[type]}                 vm LuaVM       [description]
 */
func setTable(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1

    vm.GetRK(b)
    vm.GetRK(c)
    vm.SetTable(a)
}

