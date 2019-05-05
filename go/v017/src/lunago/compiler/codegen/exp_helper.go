/*
* @Author: konyka
* @Date:   2019-05-05 08:17:37
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 08:18:25
*/

package codegen

import . "lunago/compiler/ast"

func isVarargOrFuncCall(exp Exp) bool {
    switch exp.(type) {
    case *VarargExp, *FuncCallExp:
        return true
    }
    return false
}

func removeTailNils(exps []Exp) []Exp {
    for n := len(exps) - 1; n >= 0; n-- {
        if _, ok := exps[n].(*NilExp); !ok {
            return exps[0 : n+1]
        }
    }
    return nil
}





