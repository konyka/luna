/*
* @Author: konyka
* @Date:   2019-05-04 11:13:33
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 11:13:53
*/


package parser

import . "lunago/compiler/ast"
import . "lunago/compiler/lexer"

/* recursive descent parser */

func Parse(chunk, chunkName string) *Block {
    lexer := NewLexer(chunk, chunkName)
    block := parseBlock(lexer)
    lexer.NextTokenOfKind(TOKEN_EOF)
    return block
}






