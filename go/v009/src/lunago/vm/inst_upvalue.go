/*
* @Author: konyka
* @Date:   2019-05-01 19:37:07
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 19:41:46
*/

package vm

import . "lunago/api"


/**
 * [getTabUp R(A) := UpValue[B][RK(C)]]
 * @Author   konyka
 * @DateTime 2019-05-01T19:41:34+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func getTabUp(i Instruction, vm LuaVM) {
    a, _, c := i.ABC()
    a += 1

    vm.PushGlobalTable()
    vm.GetRK(c)
    vm.GetTable(-2)
    vm.Replace(a)
    vm.Pop(1)
}








