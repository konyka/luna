/*
* @Author: konyka
* @Date:   2019-05-04 11:38:40
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 13:29:17
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
    breaks    [][]int
    parent    *funcInfo
    upvalues  map[string]upvalInfo
    insts     []uint32
    //to do
}

type locVarInfo struct {
    prev     *locVarInfo
    name     string
    scopeLv  int
    slot     int
    captured bool
}

type upvalInfo struct {
    locVarSlot int
    upvalIndex int
    index      int
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

func (self *funcInfo) enterScope(breakable bool) {
    self.scopeLv++
    if breakable {
        self.breaks = append(self.breaks, []int{}) //循环块
    } else {
        self.breaks = append(self.breaks, nil) //非循环块
    }
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

/**
 * slotOfLocVar检查局部变量名是否已经和某个寄存器绑定，如果过是，则返回寄存器的索引，否则返回-1
 */
func (self *funcInfo) slotOfLocVar(name string) int {
    if locVar, found := self.locNames[name]; found {
        return locVar.slot
    }
    return -1
}

/**
 * [func exitScope，在退出作用域的时候修复调转指令  ]
 * @Author   konyka
 * @DateTime 2019-05-04T12:54:36+0800
 * @param    {[type]}                 self *funcInfo)    exitScope( [description]
 * @return   {[type]}                      [description]
 */
func (self *funcInfo) exitScope() {
    pendingBreakJmps := self.breaks[len(self.breaks)-1]
    self.breaks = self.breaks[:len(self.breaks)-1]

    a := self.getJmpArgA()
    for _, pc := range pendingBreakJmps {
        sBx := self.pc() - pc
        i := (sBx+MAXARG_sBx)<<14 | a<<6 | OP_JMP
        self.insts[pc] = uint32(i)
    }

    self.scopeLv--
    for _, locVar := range self.locNames {
        if locVar.scopeLv > self.scopeLv { // out of scope
            self.removeLocVar(locVar)
        }
    }
}

/**
 * [func 当退出作用域以后，需要删除该作用域中的局部变量（解绑局部变量名、回收寄存器）]
 * @Author   konyka
 * @DateTime 2019-05-04T12:58:53+0800
 * @param    {[type]}                 self *funcInfo)    removeLocVar(locVar *locVarInfo [description]
 * @return   {[type]}                      [description]
 */
func (self *funcInfo) removeLocVar(locVar *locVarInfo) {
    self.freeReg()
    if locVar.prev == nil {
        delete(self.locNames, locVar.name)
    } else if locVar.prev.scopeLv == locVar.scopeLv {
        self.removeLocVar(locVar.prev)
    } else {
        self.locNames[locVar.name] = locVar.prev
    }
}

/**
 * [func addBreakJmp把break语句对应的跳转指令添加到最近的循环块中。
 * 如果找不到循环块，则调用panic报错。]
 * @Author   konyka
 * @DateTime 2019-05-04T13:12:01+0800
 * @param    {[type]}                 self *funcInfo)    addBreakJmp(pc int [description]
 * @return   {[type]}                      [description]
 */
func (self *funcInfo) addBreakJmp(pc int) {
    for i := self.scopeLv; i >= 0; i-- {
        if self.breaks[i] != nil { // breakable 循环块
            self.breaks[i] = append(self.breaks[i], pc)
            return
        }
    }

    panic("<break> at line ? not inside a loop!")
}

/* upvalues */

/**
 * indexOfUpval判断名称是否已经和Upvalue绑定，
 * 如果是，返回Upvalue索引，否则尝试绑定，然后返回索引。如果绑定失败，返回-1.
 */
func (self *funcInfo) indexOfUpval(name string) int {
    if upval, ok := self.upvalues[name]; ok {
        return upval.index
    }
    if self.parent != nil {
        if locVar, found := self.parent.locNames[name]; found {
            idx := len(self.upvalues)
            self.upvalues[name] = upvalInfo{locVar.slot, -1, idx}
            locVar.captured = true
            return idx
        }
        if uvIdx := self.parent.indexOfUpval(name); uvIdx >= 0 {
            idx := len(self.upvalues)
            self.upvalues[name] = upvalInfo{-1, uvIdx, idx}
            return idx
        }
    }
    return -1
}









