/*
* @Author: konyka
* @Date:   2019-05-04 11:38:40
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 12:41:39
*/
package codegen

import . "lunago/compiler/ast"
import . "lunago/compiler/lexer"
import . "lunago/vm"

type funcInfo struct {
    constants map[interface{}]int
    usedRegs  int
    maxRegs   int
    scopeLv   int
    locVars   []*locVarInfo
    locNames  map[string]*locVarInfo
    //to do
}

type locVarInfo struct {
    prev     *locVarInfo
    name     string
    scopeLv  int
    slot     int
    captured bool
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

/**
 * allocReg分配一个寄存器，必要的时候更新最大寄存器数量，并返回寄存器的索引
 */
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

/**
 * [func freeReg回收最近分配的寄存器]
 * @Author   konyka
 * @DateTime 2019-05-04T12:12:38+0800
 * @param    {[type]}                 self *funcInfo)    freeReg( [description]
 * @return   {[type]}                      [description]
 */
func (self *funcInfo) freeReg() {
    if self.usedRegs <= 0 {
        panic("usedRegs <= 0 !")
    }
    self.usedRegs--
}

/**
 * allocRegs分配连续的n个寄存器，返回第一个寄存器的索引
 */
func (self *funcInfo) allocRegs(n int) int {
    if n <= 0 {
        panic("n <= 0 !")
    }
    for i := 0; i < n; i++ {
        self.allocReg()
    }
    return self.usedRegs - n
}

/**
 * [func freeRegs方法回收最近分配的n个寄存器 ]
 * @Author   konyka
 * @DateTime 2019-05-04T12:14:35+0800
 * @param    {[type]}                 self *funcInfo)    freeRegs(n int [description]
 * @return   {[type]}                      [description]
 */
func (self *funcInfo) freeRegs(n int) {
    if n < 0 {
        panic("n < 0 !")
    }
    for i := 0; i < n; i++ {
        self.freeReg()
    }
}


/* lexical scope */

func (self *funcInfo) enterScope() {
    self.scopeLv++
}

/**
 * addLocVar在当前作用域里面增加一个局部变量，返回其分配的寄存器索引
 */
func (self *funcInfo) addLocVar(name string) int {
    newVar := &locVarInfo{
        name:    name,
        prev:    self.locNames[name],
        scopeLv: self.scopeLv,
        slot:    self.allocReg(),
    }

    self.locVars = append(self.locVars, newVar)
    self.locNames[name] = newVar

    return newVar.slot
}




