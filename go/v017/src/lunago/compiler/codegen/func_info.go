/*
* @Author: konyka
* @Date:   2019-05-04 11:38:40
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 12:10:20
*/
package codegen

import . "lunago/compiler/ast"
import . "lunago/compiler/lexer"
import . "lunago/vm"

type funcInfo struct {
    constants map[interface{}]int
    usedRegs  int
    maxRegs   int
    //to do
}

/* constants */

func (self *funcInfo) indexOfConstant(k interface{}) int {
    if idx, found := self.constants[k]; found {
        return idx
    }

    idx := len(self.constants)
    self.constants[k] = idx
    return idx
}

/* registers */

func (self *funcInfo) allocReg() int {
    self.usedRegs++
    if self.usedRegs >= 255 {
        panic("function or expression needs too many registers")
    }
    if self.usedRegs > self.maxRegs {
        self.maxRegs = self.usedRegs
    }
    return self.usedRegs - 1
}


func (self *funcInfo) freeReg() {
    if self.usedRegs <= 0 {
        panic("usedRegs <= 0 !")
    }
    self.usedRegs--
}













