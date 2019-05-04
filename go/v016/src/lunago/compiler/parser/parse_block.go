/*
* @Author: konyka
* @Date:   2019-05-04 08:12:28
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 08:17:00
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























