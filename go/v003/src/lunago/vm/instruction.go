/*
* @Author: konyka
* @Date:   2019-04-27 16:27:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-27 16:36:16
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


func (self Instruction) ABx() (a, bx int) {
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}
























