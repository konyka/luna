/*
* @Author: konyka
* @Date:   2019-05-04 10:36:43
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 10:37:13
*/

package parser

import . "lunago/compiler/ast"
import . "lunago/compiler/lexer"

// prefixexp ::= var | functioncall | ‘(’ exp ‘)’
// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args

/*
prefixexp ::= Name
    | ‘(’ exp ‘)’
    | prefixexp ‘[’ exp ‘]’
    | prefixexp ‘.’ Name
    | prefixexp [‘:’ Name] args
*/
func parsePrefixExp(lexer *Lexer) Exp {
    var exp Exp
    if lexer.LookAhead() == TOKEN_IDENTIFIER {
        line, name := lexer.NextIdentifier() // Name
        exp = &NameExp{line, name}
    } else { // ‘(’ exp ‘)’
        exp = parseParensExp(lexer)
    }
    return _finishPrefixExp(lexer, exp)
}












