/*
* @Author: konyka
* @Date:   2019-04-29 19:13:25
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 19:41:36
*/

package vm

import . "lunago/api"


/**
 * [loadNil R(A), R(A+1), ..., R(A+B) := nil]
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











