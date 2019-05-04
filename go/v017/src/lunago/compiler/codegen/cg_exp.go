/*
* @Author: konyka
* @Date:   2019-05-04 22:29:18
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 22:57:58
*/


package codegen

import . "lunago/compiler/ast"
import . "lunago/compiler/lexer"
import . "lunago/vm"

// todo: rename to evalExp()?
func cgExp(fi *funcInfo, node Exp, a, n int) {
    switch exp := node.(type) {
    case *NilExp:
        fi.emitLoadNil(a, n)
    case *FalseExp:
        fi.emitLoadBool(a, 0, 0)
    case *TrueExp:
        fi.emitLoadBool(a, 1, 0)
    case *IntegerExp:
        fi.emitLoadK(a, exp.Val)
    case *FloatExp:
        fi.emitLoadK(a, exp.Val)
    case *StringExp:
        fi.emitLoadK(a, exp.Str)
    case *ParensExp:
        cgExp(fi, exp.Exp, a, 1)
    case *VarargExp:
        cgVarargExp(fi, exp, a, n)
    case *FuncDefExp:
        cgFuncDefExp(fi, exp, a)
    case *TableConstructorExp:
        cgTableConstructorExp(fi, exp, a)
    case *UnopExp:
        cgUnopExp(fi, exp, a)
    case *BinopExp:
        cgBinopExp(fi, exp, a)
    case *ConcatExp:
        cgConcatExp(fi, exp, a)
    case *NameExp:
        cgNameExp(fi, exp, a)
    case *TableAccessExp:
        cgTableAccessExp(fi, exp, a)
    case *FuncCallExp:
        cgFuncCallExp(fi, exp, a, n)
    }
}

func cgVarargExp(fi *funcInfo, node *VarargExp, a, n int) {
    if !fi.isVararg {
        panic("cannot use '...' outside a vararg function")
    }
    fi.emitVararg(a, n)
}

// f[a] := function(args) body end
func cgFuncDefExp(fi *funcInfo, node *FuncDefExp, a int) {
    subFI := newFuncInfo(fi, node)
    fi.subFuncs = append(fi.subFuncs, subFI)

    for _, param := range node.ParList {
        subFI.addLocVar(param)
    }

    cgBlock(subFI, node.Block)
    subFI.exitScope()
    subFI.emitReturn(0, 0)

    bx := len(fi.subFuncs) - 1
    fi.emitClosure(a, bx)
}

func cgTableConstructorExp(fi *funcInfo, node *TableConstructorExp, a int) {
    nArr := 0
    for _, keyExp := range node.KeyExps {
        if keyExp == nil {
            nArr++
        }
    }
    nExps := len(node.KeyExps)
    multRet := nExps > 0 &&
        isVarargOrFuncCall(node.ValExps[nExps-1])

    fi.emitNewTable(a, nArr, nExps-nArr)

    arrIdx := 0
    for i, keyExp := range node.KeyExps {
        valExp := node.ValExps[i]

        if keyExp == nil {
            arrIdx++
            tmp := fi.allocReg()
            if i == nExps-1 && multRet {
                cgExp(fi, valExp, tmp, -1)
            } else {
                cgExp(fi, valExp, tmp, 1)
            }

            if arrIdx%50 == 0 || arrIdx == nArr { // LFIELDS_PER_FLUSH
                n := arrIdx % 50
                if n == 0 {
                    n = 50
                }
                fi.freeRegs(n)
                c := (arrIdx-1)/50 + 1 // todo: c > 0xFF
                if i == nExps-1 && multRet {
                    fi.emitSetList(a, 0, c)
                } else {
                    fi.emitSetList(a, n, c)
                }
            }

            continue
        }

        b := fi.allocReg()
        cgExp(fi, keyExp, b, 1)
        c := fi.allocReg()
        cgExp(fi, valExp, c, 1)
        fi.freeRegs(2)

        fi.emitSetTable(a, b, c)
    }
}


// r[a] := op exp
func cgUnopExp(fi *funcInfo, node *UnopExp, a int) {
    b := fi.allocReg()
    cgExp(fi, node.Exp, b, 1)
    fi.emitUnaryOp(node.Op, a, b)
    fi.freeReg()
}










