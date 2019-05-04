/*
* @Author: konyka
* @Date:   2019-05-04 14:22:11
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 21:02:41
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
    fi.enterScope(false)
    cgBlock(fi, node.Block)
    fi.closeOpenUpvals()
    fi.exitScope()
}



















