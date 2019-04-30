/*
* @Author: konyka
* @Date:   2019-04-30 19:26:44
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 00:09:59
*/


package vm

import . "lunago/api"


/**
 * [closure R(A) := closure(KPROTO[Bx])]
 *  closure指令（iBx模式）把当前lua函数的子函数原型实例化为闭包，
 *  放到由操作数A指定的寄存器中，子函数原型来自当前函数原型的子函数原型列表，索引由操作数Bx指定
 *  closure 指令对应lua脚本里面的函数定义语句或者表达式
 * @Author   konyka
 * @DateTime 2019-05-01T00:00:44+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func closure(i Instruction, vm LuaVM) {
    a, bx := i.ABx()
    a += 1

    vm.LoadProto(bx)
    vm.Replace(a)
}


/**
 * [call R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))]
 * CALL指令（iABC模式）调用Lua函数。其中被调函数位于寄存器中，索引由操作数A指定，
 * 需要传递给被调函数的参数值也要在寄存器中，紧挨着被调函数，数量由操作数B指定，
 * 函数调用结束后，原先存放在函数和参数值的寄存器会被返回值占据，具体由多少个返回值则由操作数C指定
 * @Author   konyka
 * @DateTime 2019-05-01T00:09:33+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func call(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1

    // println(":::"+ vm.StackToString())
    nArgs := _pushFuncAndArgs(a, b, vm)
    vm.Call(nArgs, c-1)
    _popResults(a, c, vm)
}











