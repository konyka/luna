/*
* @Author: konyka
* @Date:   2019-05-03 11:57:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 12:38:40
*/

package lexer

import "bytes"
import "fmt"
import "regexp"
import "strconv"
import "strings"

type Lexer struct {
    chunk         string // source code
    chunkName     string // source name
    line          int    // current line number
}

func NewLexer(chunk, chunkName string) *Lexer {
    return &Lexer{chunk, chunkName, 1}
}

func (self *Lexer) NextToken() (line, kind int, token string) {

    self.skipWhiteSpaces()
    if len(self.chunk) == 0 {
        return self.line, TOKEN_EOF, "EOF"
    }

    switch self.chunk[0] {
    case ';':
        self.next(1)
        return self.line, TOKEN_SEP_SEMI, ";"
    case ',':
        self.next(1)
        return self.line, TOKEN_SEP_COMMA, ","
    }

    return
}





