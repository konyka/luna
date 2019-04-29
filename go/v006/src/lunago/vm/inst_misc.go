/*
* @Author: konyka
* @Date:   2019-04-29 18:41:22
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 18:42:50
*/


package vm

import . "lunago/api"

/**
 * [move R(A) := R(B)
 * 虽然说是move指令，实际上叫做copy指令可能会更贴切一些，因为源寄存器的值还原封不动的待在原地。]
 * @Author   konyka
 * @DateTime 2019-04-29T18:42:01+0800
 * @param    {[type]}                 i  Instruction   [description]
 * @param    {[type]}                 vm LuaVM         [description]
 * @return   {[type]}                    [description]
 */
func move(i Instruction, vm LuaVM) {
    a, b, _ := i.ABC()
    a += 1
    b += 1

    vm.Copy(b, a)
}







