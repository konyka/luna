/*
* @Author: konyka
* @Date:   2019-05-04 22:29:18
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 22:33:32
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








