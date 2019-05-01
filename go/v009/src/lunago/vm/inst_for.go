/*
* @Author: konyka
* @Date:   2019-04-30 09:08:09
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 09:16:18
*/

package vm

import . "lunago/api"


/**
 * [forPrep R(A)-=R(A+2); pc+=sBx]
 * @Author   konyka
 * @DateTime 2019-04-30T09:09:55+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func forPrep(i Instruction, vm LuaVM) {
    a, sBx := i.AsBx()
    a += 1

    if vm.Type(a) == LUA_TSTRING {
        vm.PushNumber(vm.ToNumber(a))
        vm.Replace(a)
    }
    if vm.Type(a+1) == LUA_TSTRING {
        vm.PushNumber(vm.ToNumber(a + 1))
        vm.Replace(a + 1)
    }
    if vm.Type(a+2) == LUA_TSTRING {
        vm.PushNumber(vm.ToNumber(a + 2))
        vm.Replace(a + 2)
    }

    vm.PushValue(a)
    vm.PushValue(a + 2)
    vm.Arith(LUA_OPSUB)
    vm.Replace(a)
    vm.AddPC(sBx)
}



/**
 * [forLoop   R(A)+=R(A+2);
 *            if R(A) <?= R(A+1) then {
 *              pc+=sBx; R(A+3)=R(A)
 *             }]
 * @Author   konyka
 * @DateTime 2019-04-30T09:15:45+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func forLoop(i Instruction, vm LuaVM) {
    a, sBx := i.AsBx()
    a += 1

    // R(A)+=R(A+2);
    vm.PushValue(a + 2)
    vm.PushValue(a)
    vm.Arith(LUA_OPADD)
    vm.Replace(a)

    isPositiveStep := vm.ToNumber(a+2) >= 0
    if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
        !isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {

        // pc+=sBx; R(A+3)=R(A)
        vm.AddPC(sBx)
        vm.Copy(a, a+3)
    }
}






