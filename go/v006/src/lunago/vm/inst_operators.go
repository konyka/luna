/*
* @Author: konyka
* @Date:   2019-04-29 20:32:58
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 08:31:56
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


/**
 * [concat R(A) := R(B).. ... ..R(C)
 * cancat(iABC 模式)，将连续的n个寄存器（起止索引分别由操作数B、C指定）里的值拼接，
 * 将结果放到另一个寄存器中（索引由操作数A指定）
 * 在实现前面的指令时，最多只是往栈顶push了一两个值，所以我们可以在创建Lua栈的时候把容量设置的稍大一些，
 * 这样在push少量的值之前，就不需要检查栈的剩余空间了。
 * 但是concat指令则有所不同，因为进行拼接的值的数量不是固定的，所以在吧这些值push到栈顶之前，
 * 必须调用CheckStack（）确保还有足够的空间可以容纳这些值，否则可能会导致溢出]
 * @Author   konyka
 * @DateTime 2019-04-29T22:44:01+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func concat(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1
    b += 1
    c += 1

    n := c - b + 1
    vm.CheckStack(n)
    for i := b; i <= c; i++ {
        vm.PushValue(i)
    }
    vm.Concat(n)
    vm.Replace(a)
}


/**
 * [_compare if ((RK(B) op RK(C)) ~= A) then pc++
 * 
 * 比较指令（iABC 模式），比较寄存器或者常量表里面的两个值（索引分别由操作数B、C指定），
 * 如果比较结果和操作数A（转换为布尔值）匹配，则跳过下一条指令。比较指令不会改变寄存器的状态。
 * 
 * if（RK（B）op RK（C）～= A）then pc++]
 * @Author   konyka
 * @DateTime 2019-04-30T08:10:33+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @param    {[type]}                 op CompareOp     [description]
 * @return   {[type]}                    [description]
 */
func _compare(i Instruction, vm LuaVM, op CompareOp) {
    a, b, c := i.ABC()

    vm.GetRK(b)
    vm.GetRK(c)
    if vm.Compare(-2, -1, op) != (a != 0) {
        vm.AddPC(1)
    }
    vm.Pop(2)
}

/* compare op*/

func eq(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPEQ) } // ==
func lt(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLT) } // <
func le(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLE) } // <=

/* logical */

// R(A) := not R(B)
/**
 * [not R(A) := not R(B)]
 * @Author   konyka
 * @DateTime 2019-04-30T08:22:34+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func not(i Instruction, vm LuaVM) {
    a, b, _ := i.ABC()
    a += 1
    b += 1

    vm.PushBoolean(!vm.ToBoolean(b))
    vm.Replace(a)
}


/**
 * [testSet if (R(B) <=> C) then R(A) := R(B) else pc++]
 *testset指令（iABC 模式），破案段寄存器B（索引由操作数B指定）中的值转换为布尔值之后，
 *是否和操作数C表示的布尔值一致，
 *如果一样，则将寄存器B中的值符知道寄存器A中，索引由操作数A指定，否则跳过下一条指令。
 *
 * @Author   konyka
 * @DateTime 2019-04-30T08:31:06+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func testSet(i Instruction, vm LuaVM) {
    a, b, c := i.ABC()
    a += 1
    b += 1

    if vm.ToBoolean(b) == (c != 0) {
        vm.Copy(b, a)
    } else {
        vm.AddPC(1)
    }
}

