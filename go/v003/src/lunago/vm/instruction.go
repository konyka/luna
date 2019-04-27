/*
* @Author: konyka
* @Date:   2019-04-27 16:27:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-27 16:33:28
*/

package vm

type Instruction uint32

func (self Instruction) Opcode() int {
	return int(self & 0x3F)
}


func (self Instruction) ABC() (a, b, c int) {
	a = int(self >> 6 & 0xFF)
	c = int(self >> 14 & 0x1FF)
	b = int(self >> 23 & 0x1FF)
	return
}















