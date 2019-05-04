/*
* @Author: konyka
* @Date:   2019-05-04 08:12:28
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 08:14:02
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











