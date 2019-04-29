/*
* @Author: konyka
* @Date:   2019-04-29 20:32:58
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 21:15:19
*/

package vm

import . "lunago/api"


/**
 * [_binaryArith R(A) := RK(B) op RK(C)
 * 二元算术运算指令（iABC 模式），对连个寄存器或者常量值（索引由操作数B、C指定）进行运算， 
 * 将结果放到另一个寄存器中（随你由操作数A指定）。如果用RK（N）表示寄存器或者常量值，
 * 那么二元算术运算指令的伪代码可以如下表示：

    R（A）:= RK(B) op RK(C)  ]
 * @Author   konyka
 * @DateTime 2019-04-29T20:39:18+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @param    {[type]}                 op ArithOp       [description]
 * @return   {[type]}                    [description]
 */
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
    a, b, c := i.ABC()
    a += 1

    vm.GetRK(b)
    vm.GetRK(c)
    vm.Arith(op)
    vm.Replace(a)
}


/**
 * [_unaryArith R(A) := op R(B) 
 * 元算术运算指令（iABC 模式），对操作数B所指定的寄存器里面的值进行运算，
    然后把结果放到操作数 A 所指定的寄存器中，操作数 C 没有使用。]
 * @Author   konyka
 * @DateTime 2019-04-29T21:03:50+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @param    {[type]}                 op ArithOp       [description]
 * @return   {[type]}                    [description]
 */
func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
    a, b, _ := i.ABC()
    a += 1
    b += 1

    vm.PushValue(b)
    vm.Arith(op)
    vm.Replace(a)
}

/* arith */

func add(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPADD) }  // +
func sub(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSUB) }  // -
func mul(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMUL) }  // *
func mod(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMOD) }  // %
func pow(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPPOW) }  // ^
func div(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPDIV) }  // /
func idiv(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPIDIV) } // //
func band(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBAND) } // &
func bor(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPBOR) }  // |
func bxor(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBXOR) } // ~
func shl(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHL) }  // <<
func shr(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHR) }  // >>
func unm(i Instruction, vm LuaVM)  { _unaryArith(i, vm, LUA_OPUNM) }   // -
func bnot(i Instruction, vm LuaVM) { _unaryArith(i, vm, LUA_OPBNOT) }  // ~


/**
 * [length R(A) := length of R(B)]
 * @Author   konyka
 * @DateTime 2019-04-29T21:15:12+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func length(i Instruction, vm LuaVM) {
    a, b, _ := i.ABC()
    a += 1
    b += 1

    vm.Len(b)
    vm.Replace(a)
}




