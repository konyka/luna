/*
* @Author: konyka
* @Date:   2019-05-04 10:36:43
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 10:49:28
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

func _finishPrefixExp(lexer *Lexer, exp Exp) Exp {
    for {
        switch lexer.LookAhead() {
        case TOKEN_SEP_LBRACK: // prefixexp ‘[’ exp ‘]’
            lexer.NextToken()                       // ‘[’
            keyExp := parseExp(lexer)               // exp
            lexer.NextTokenOfKind(TOKEN_SEP_RBRACK) // ‘]’
            exp = &TableAccessExp{lexer.Line(), exp, keyExp}
        case TOKEN_SEP_DOT: // prefixexp ‘.’ Name
            lexer.NextToken()                    // ‘.’
            line, name := lexer.NextIdentifier() // Name
            keyExp := &StringExp{line, name}
            exp = &TableAccessExp{line, exp, keyExp}
        case TOKEN_SEP_COLON, // prefixexp ‘:’ Name args
            TOKEN_SEP_LPAREN, TOKEN_SEP_LCURLY, TOKEN_STRING: // prefixexp args
            exp = _finishFuncCallExp(lexer, exp)
        default:
            return exp
        }
    }
    return exp
}

/**
 * 圆括号表达式的解析函数
 */
func parseParensExp(lexer *Lexer) Exp {
    lexer.NextTokenOfKind(TOKEN_SEP_LPAREN) // (
    exp := parseExp(lexer)                  // exp
    lexer.NextTokenOfKind(TOKEN_SEP_RPAREN) // )

    switch exp.(type) {
    case *VarargExp, *FuncCallExp, *NameExp, *TableAccessExp:
        return &ParensExp{exp}
    }

    // no need to keep parens
    return exp
}

/**
 * functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
 */
func _finishFuncCallExp(lexer *Lexer, prefixExp Exp) *FuncCallExp {
    nameExp := _parseNameExp(lexer)
    line := lexer.Line() // todo
    args := _parseArgs(lexer)
    lastLine := lexer.Line()
    return &FuncCallExp{line, lastLine, prefixExp, nameExp, args}
}





















