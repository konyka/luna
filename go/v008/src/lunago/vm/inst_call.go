/*
* @Author: konyka
* @Date:   2019-04-30 19:26:44
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 00:01:28
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













