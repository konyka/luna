/*
* @Author: konyka
* @Date:   2019-05-04 08:33:46
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 09:53:46
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

/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
     prefixexp | tableconstructor | exp binop exp | unop exp
*/
/*
exp   ::= exp12
exp12 ::= exp11 {or exp11}
exp11 ::= exp10 {and exp10}
exp10 ::= exp9 {(‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’) exp9}
exp9  ::= exp8 {‘|’ exp8}
exp8  ::= exp7 {‘~’ exp7}
exp7  ::= exp6 {‘&’ exp6}
exp6  ::= exp5 {(‘<<’ | ‘>>’) exp5}
exp5  ::= exp4 {‘..’ exp4}
exp4  ::= exp3 {(‘+’ | ‘-’) exp3}
exp3  ::= exp2 {(‘*’ | ‘/’ | ‘//’ | ‘%’) exp2}
exp2  ::= {(‘not’ | ‘#’ | ‘-’ | ‘~’)} exp1
exp1  ::= exp0 {‘^’ exp2}
exp0  ::= nil | false | true | Numeral | LiteralString
        | ‘...’ | functiondef | prefixexp | tableconstructor
*/
func parseExp(lexer *Lexer) Exp {
    return parseExp12(lexer)
}














