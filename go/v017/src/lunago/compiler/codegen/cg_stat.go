/*
* @Author: konyka
* @Date:   2019-05-04 14:22:11
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 21:32:39
*/


package codegen

import . "lunago/compiler/ast"

func cgStat(fi *funcInfo, node Stat) {
    switch stat := node.(type) {
    case *FuncCallStat:
        cgFuncCallStat(fi, stat)
    case *BreakStat:
        cgBreakStat(fi, stat)
    case *DoStat:
        cgDoStat(fi, stat)
    case *WhileStat:
        cgWhileStat(fi, stat)
    case *RepeatStat:
        cgRepeatStat(fi, stat)
    case *IfStat:
        cgIfStat(fi, stat)
    case *ForNumStat:
        cgForNumStat(fi, stat)
    case *ForInStat:
        cgForInStat(fi, stat)
    case *AssignStat:
        cgAssignStat(fi, stat)
    case *LocalVarDeclStat:
        cgLocalVarDeclStat(fi, stat)
    case *LocalFuncDefStat:
        cgLocalFuncDefStat(fi, stat)
    case *LabelStat, *GotoStat:
        panic("label and goto statements are not supported!")
    }
}

func cgLocalFuncDefStat(fi *funcInfo, node *LocalFuncDefStat) {
    r := fi.addLocVar(node.Name)
    cgFuncDefExp(fi, node.Exp, r)
}

func cgFuncCallStat(fi *funcInfo, node *FuncCallStat) {
    r := fi.allocReg()
    cgFuncCallExp(fi, node, r, 0)
    fi.freeReg()
}

func cgBreakStat(fi *funcInfo, node *BreakStat) {
    pc := fi.emitJmp(0, 0)
    fi.addBreakJmp(pc)
}


func cgDoStat(fi *funcInfo, node *DoStat) {
    fi.enterScope(false)    //非循环块
    cgBlock(fi, node.Block)
    fi.closeOpenUpvals()
    fi.exitScope()
}

/*
           ______________
          /  false? jmp  |
         /               |
while exp do block end <-'
      ^           \
      |___________/
           jmp
*/
func cgWhileStat(fi *funcInfo, node *WhileStat) {
    pcBeforeExp := fi.pc()

    r := fi.allocReg()
    cgExp(fi, node.Exp, r, 1)
    fi.freeReg()

    fi.emitTest(r, 0)
    pcJmpToEnd := fi.emitJmp(0, 0)

    fi.enterScope(true)
    cgBlock(fi, node.Block)
    fi.closeOpenUpvals()
    fi.emitJmp(0, pcBeforeExp-fi.pc()-1)
    fi.exitScope()

    fi.fixSbx(pcJmpToEnd, fi.pc()-pcJmpToEnd)
}

/*
        ______________
       |  false? jmp  |
       V              /
repeat block until exp
*/
func cgRepeatStat(fi *funcInfo, node *RepeatStat) {
    fi.enterScope(true)

    pcBeforeBlock := fi.pc()
    cgBlock(fi, node.Block)

    r := fi.allocReg()
    cgExp(fi, node.Exp, r, 1)
    fi.freeReg()

    fi.emitTest(r, 0)
    fi.emitJmp(fi.getJmpArgA(), pcBeforeBlock-fi.pc()-1)
    fi.closeOpenUpvals()

    fi.exitScope()
}

/*
         _________________       _________________       _____________
        / false? jmp      |     / false? jmp      |     / false? jmp  |
       /                  V    /                  V    /              V
if exp1 then block1 elseif exp2 then block2 elseif true then block3 end <-.
                   \                       \                       \      |
                    \_______________________\_______________________\_____|
                    jmp                     jmp                     jmp
*/
func cgIfStat(fi *funcInfo, node *IfStat) {
    pcJmpToEnds := make([]int, len(node.Exps))
    pcJmpToNextExp := -1

    for i, exp := range node.Exps {
        if pcJmpToNextExp >= 0 {
            fi.fixSbx(pcJmpToNextExp, fi.pc()-pcJmpToNextExp)
        }

        r := fi.allocReg()
        cgExp(fi, exp, r, 1)
        fi.freeReg()

        fi.emitTest(r, 0)
        pcJmpToNextExp = fi.emitJmp(0, 0)

        fi.enterScope(false)
        cgBlock(fi, node.Blocks[i])
        fi.closeOpenUpvals()
        fi.exitScope()
        if i < len(node.Exps)-1 {
            pcJmpToEnds[i] = fi.emitJmp(0, 0)
        } else {
            pcJmpToEnds[i] = pcJmpToNextExp
        }
    }

    for _, pc := range pcJmpToEnds {
        fi.fixSbx(pc, fi.pc()-pc)
    }
}














