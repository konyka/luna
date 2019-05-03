/*
* @Author: konyka
* @Date:   2019-05-03 11:57:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 13:10:50
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

func (self *Lexer) skipWhiteSpaces() {
    for len(self.chunk) > 0 {
        if self.test("--") {
            self.skipComment()
        } else if self.test("\r\n") || self.test("\n\r") {
            self.next(2)
            self.line += 1
        } else if isNewLine(self.chunk[0]) {
            self.next(1)
            self.line += 1
        } else if isWhiteSpace(self.chunk[0]) {
            self.next(1)
        } else {
            break
        }
    }
}

func (self *Lexer) test(s string) bool {
    return strings.HasPrefix(self.chunk, s)
}

