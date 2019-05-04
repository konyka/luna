/*
* @Author: konyka
* @Date:   2019-05-04 08:12:28
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 08:25:35
*/

package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// block ::= {stat} [retstat]
func parseBlock(lexer *Lexer) *Block {
    return &Block{
        Stats:    parseStats(lexer),
        RetExps:  parseRetExps(lexer),
        LastLine: lexer.Line(),
    }
}

func parseStats(lexer *Lexer) []Stat {
    stats := make([]Stat, 0, 8)
    for !_isReturnOrBlockEnd(lexer.LookAhead()) {
        stat := parseStat(lexer)
        if _, ok := stat.(*EmptyStat); !ok {
            stats = append(stats, stat)
        }
    }
    return stats
}

func _isReturnOrBlockEnd(tokenKind int) bool {
    switch tokenKind {
    case TOKEN_KW_RETURN, TOKEN_EOF, TOKEN_KW_END,
        TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
        return true
    }
    return false
}






















