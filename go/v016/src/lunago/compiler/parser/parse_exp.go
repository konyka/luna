/*
* @Author: konyka
* @Date:   2019-05-04 08:33:46
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 08:35:15
*/

package parser

import . "lunago/compiler/ast"
import . "lunago/compiler/lexer"
import "lunago/number"

// explist ::= exp {‘,’ exp}
func parseExpList(lexer *Lexer) []Exp {
    exps := make([]Exp, 0, 4)
    exps = append(exps, parseExp(lexer))
    for lexer.LookAhead() == TOKEN_SEP_COMMA {
        lexer.NextToken()
        exps = append(exps, parseExp(lexer))
    }
    return exps
}
















