/*
* @Author: konyka
* @Date:   2019-05-04 08:41:23
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 08:55:00
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

/**
 * ‘::’ Name ‘::’
 */
func parseLabelStat(lexer *Lexer) *LabelStat {
    lexer.NextTokenOfKind(TOKEN_SEP_LABEL) // ::
    _, name := lexer.NextIdentifier()      // name
    lexer.NextTokenOfKind(TOKEN_SEP_LABEL) // ::
    return &LabelStat{name}
}

/**
 * goto Name
 */
func parseGotoStat(lexer *Lexer) *GotoStat {
    lexer.NextTokenOfKind(TOKEN_KW_GOTO) // goto
    _, name := lexer.NextIdentifier()    // name
    return &GotoStat{name}
}

/**
 * do block end
 */
func parseDoStat(lexer *Lexer) *DoStat {
    lexer.NextTokenOfKind(TOKEN_KW_DO)  // do
    block := parseBlock(lexer)          // block
    lexer.NextTokenOfKind(TOKEN_KW_END) // end
    return &DoStat{block}
}

/**
 * while exp do block end
 */
func parseWhileStat(lexer *Lexer) *WhileStat {
    lexer.NextTokenOfKind(TOKEN_KW_WHILE) // while
    exp := parseExp(lexer)                // exp
    lexer.NextTokenOfKind(TOKEN_KW_DO)    // do
    block := parseBlock(lexer)            // block
    lexer.NextTokenOfKind(TOKEN_KW_END)   // end
    return &WhileStat{exp, block}
}
