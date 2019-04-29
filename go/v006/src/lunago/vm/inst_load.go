/*
* @Author: konyka
* @Date:   2019-04-29 19:13:25
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 20:18:26
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


/**
 * [loadKx R(A) := Kst(extra arg) 
 * 先调用之前准备好的GetConst（）函数吧给定的常量push到栈顶，然后调用Replace（）把它移动到指定的索引处。
 * 操作数Bx占用18个bit，能表示的最大无符号整数是262143，大部分lua函数的常量表大小都不会超过这个数字，
 * 因此这个限制通常不是神没问题。不过lua也经常被当作数据描述语言使用，
 * 因此常量表的大小可能会超出这个现实也并不奇怪，为了应对这种情况，lua还提供了一条loadkx指令。

 * loadkx指令（也是iABx模式）需要和EXTEAARG指令（iAx模式）配合使用。
 * 用后者的Ax操作数来指定常量的索引。Ax操作数占用26个bit，可以表达的最大无符号整数是67108864，
 * 可以满足大部分情况了。]
 * @Author   konyka
 * @DateTime 2019-04-29T20:17:26+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func loadKx(i Instruction, vm LuaVM) {
    a, _ := i.ABx()
    a += 1
    ax := Instruction(vm.Fetch()).Ax()

    //vm.CheckStack(1)
    vm.GetConst(ax)
    vm.Replace(a)
}





