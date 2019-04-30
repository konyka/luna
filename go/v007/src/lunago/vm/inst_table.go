/*
* @Author: konyka
* @Date:   2019-04-30 13:01:31
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 15:27:23
*/


package vm

import . "lunago/api"
import "lunago/number"

/* number of list items to accumulate before a SETLIST instruction */
const LFIELDS_PER_FLUSH = 50

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

/**
 * [setList R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B]
 * settable是通用指令，每次只处理一个键值对，具体操作交给表去处理，
 *  并不关心实际写入的是表的哈希部分，
 *  还是数组部分。setlist指令（iABC 模式），
 *  则是专门给数组准备的，用于按照索引批量设置数组元素，其中数组位于寄存器中，索引由操作数A指定；
 *  需要写入数组的若干个值也在寄存器中，紧挨着数组，数量由操作数B指定；数组起始索引则由操作数C指定。
 *  数组的索引到底是怎么计算的？因为C只有9个bit，所以直接使用它来表示数组的索引显然是不够的。
 *  此处的解决的办法就是，让操作数C保存批次数，然后用批次数 乘上 批次大小（对应上面的fpf), 
 *  就可以计算出数组的起始索引。默认的批次大小为50，
 *  操作数c能表示的最大索引就扩大到了25600（500 * 512）

 *  但是，如果数组的长度大于这个数值呢？是不是后面的元素就只能用settable指令设置了？
 *  这种情况下，setlist指令后面会跟着一条extraarg指令，  
 *  用它的Ax操作数来保存批次数量。如果指令的操作数C大于0，那么表示的是批次数 +1，否则，
 *  整整的批次数量保存在后续的extraarg指令里面。
 * @Author   konyka
 * @DateTime 2019-04-30T15:20:58+0800
 * @param    {[type]}                 i  Instruction [description]
 * @param    {[type]}                 vm LuaVM       [description]
 */
func setList(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1

    if c > 0 {
        c = c - 1
    } else {
        c = Instruction(vm.Fetch()).Ax()
    }

    vm.CheckStack(1)
    idx := int64(c * LFIELDS_PER_FLUSH)
    for j := 1; j <= b; j++ {
        idx++
        vm.PushValue(a + j)
        vm.SetI(a, idx)
    }
}

