/*
* @Author: konyka
* @Date:   2019-05-04 08:41:23
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 08:49:54
*/

package parser

import . "lunago/compiler/ast"
import . "lunago/compiler/lexer"

/*
stat ::=  ‘;’
    | break
    | ‘::’ Name ‘::’
    | goto Name
    | do block end
    | while exp do block end
    | repeat block until exp
    | if exp then block {elseif exp then block} [else block] end
    | for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
    | for namelist in explist do block end
    | function funcname funcbody
    | local function Name funcbody
    | local namelist [‘=’ explist]
    | varlist ‘=’ explist
    | functioncall
*/
func parseStat(lexer *Lexer) Stat {
    switch lexer.LookAhead() {
    case TOKEN_SEP_SEMI:
        return parseEmptyStat(lexer)
    case TOKEN_KW_BREAK:
        return parseBreakStat(lexer)
    case TOKEN_SEP_LABEL:
        return parseLabelStat(lexer)
    case TOKEN_KW_GOTO:
        return parseGotoStat(lexer)
    case TOKEN_KW_DO:
        return parseDoStat(lexer)
    case TOKEN_KW_WHILE:
        return parseWhileStat(lexer)
    case TOKEN_KW_REPEAT:
        return parseRepeatStat(lexer)
    case TOKEN_KW_IF:
        return parseIfStat(lexer)
    case TOKEN_KW_FOR:
        return parseForStat(lexer)
    case TOKEN_KW_FUNCTION:
        return parseFuncDefStat(lexer)
    case TOKEN_KW_LOCAL:
        return parseLocalAssignOrFuncDefStat(lexer)
    default:
        return parseAssignOrFuncCallStat(lexer)
    }
}

/**
 * ;
 */
func parseEmptyStat(lexer *Lexer) *EmptyStat {
    lexer.NextTokenOfKind(TOKEN_SEP_SEMI)
    return _statEmpty
}


/**
 * break
 */
func parseBreakStat(lexer *Lexer) *BreakStat {
    lexer.NextTokenOfKind(TOKEN_KW_BREAK)
    return &BreakStat{lexer.Line()}
}







