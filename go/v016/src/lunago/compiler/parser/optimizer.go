/*
* @Author: konyka
* @Date:   2019-05-04 10:59:33
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 11:08:29
*/

package parser

import "math"
import "lunago/number"
import . "lunago/compiler/ast"
import . "lunago/compiler/lexer"

func optimizeUnaryOp(exp *UnopExp) Exp {
    switch exp.Op {
    case TOKEN_OP_UNM:
        return optimizeUnm(exp)
    case TOKEN_OP_NOT:
        return optimizeNot(exp)
    case TOKEN_OP_BNOT:
        return optimizeBnot(exp)
    default:
        return exp
    }
}







