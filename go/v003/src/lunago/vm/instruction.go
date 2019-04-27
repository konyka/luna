/*
* @Author: konyka
* @Date:   2019-04-27 16:27:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-27 16:30:44
*/

package vm

type Instruction uint32

func (self Instruction) Opcode() int {
	return int(self & 0x3F)
}


















