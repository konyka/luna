/*
* @Author: konyka
* @Date:   2019-05-03 11:57:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 13:25:54
*/

package lexer

import "bytes"
import "fmt"
import "regexp"
import "strconv"
import "strings"

var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)


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
    case '(':
        self.next(1)
        return self.line, TOKEN_SEP_LPAREN, "("
    case ')':
        self.next(1)
        return self.line, TOKEN_SEP_RPAREN, ")"
    case ']':
        self.next(1)
        return self.line, TOKEN_SEP_RBRACK, "]"
    case '{':
        self.next(1)
        return self.line, TOKEN_SEP_LCURLY, "{"
    case '}':
        self.next(1)
        return self.line, TOKEN_SEP_RCURLY, "}"
    case '+':
        self.next(1)
        return self.line, TOKEN_OP_ADD, "+"
    case '-':
        self.next(1)
        return self.line, TOKEN_OP_MINUS, "-"
    case '*':
        self.next(1)
        return self.line, TOKEN_OP_MUL, "*"
    case '^':
        self.next(1)
        return self.line, TOKEN_OP_POW, "^"
    case '%':
        self.next(1)
        return self.line, TOKEN_OP_MOD, "%"
    case '&':
        self.next(1)
        return self.line, TOKEN_OP_BAND, "&"
    case '|':
        self.next(1)
        return self.line, TOKEN_OP_BOR, "|"
    case '#':
        self.next(1)
        return self.line, TOKEN_OP_LEN, "#"
    case ':':
        if self.test("::") {
            self.next(2)
            return self.line, TOKEN_SEP_LABEL, "::"
        } else {
            self.next(1)
            return self.line, TOKEN_SEP_COLON, ":"
        }
    case '/':
        if self.test("//") {
            self.next(2)
            return self.line, TOKEN_OP_IDIV, "//"
        } else {
            self.next(1)
            return self.line, TOKEN_OP_DIV, "/"
        }
    case '~':
        if self.test("~=") {
            self.next(2)
            return self.line, TOKEN_OP_NE, "~="
        } else {
            self.next(1)
            return self.line, TOKEN_OP_WAVE, "~"
        }
    case '=':
        if self.test("==") {
            self.next(2)
            return self.line, TOKEN_OP_EQ, "=="
        } else {
            self.next(1)
            return self.line, TOKEN_OP_ASSIGN, "="
        }
    case '<':
        if self.test("<<") {
            self.next(2)
            return self.line, TOKEN_OP_SHL, "<<"
        } else if self.test("<=") {
            self.next(2)
            return self.line, TOKEN_OP_LE, "<="
        } else {
            self.next(1)
            return self.line, TOKEN_OP_LT, "<"
        }
    case '>':
        if self.test(">>") {
            self.next(2)
            return self.line, TOKEN_OP_SHR, ">>"
        } else if self.test(">=") {
            self.next(2)
            return self.line, TOKEN_OP_GE, ">="
        } else {
            self.next(1)
            return self.line, TOKEN_OP_GT, ">"
        }
    case '.':
        if self.test("...") {
            self.next(3)
            return self.line, TOKEN_VARARG, "..."
        } else if self.test("..") {
            self.next(2)
            return self.line, TOKEN_OP_CONCAT, ".."
        } else if len(self.chunk) == 1 || !isDigit(self.chunk[1]) {
            self.next(1)
            return self.line, TOKEN_SEP_DOT, "."
        }
    case '[':
        if self.test("[[") || self.test("[=") {
            return self.line, TOKEN_STRING, self.scanLongString()
        } else {
            self.next(1)
            return self.line, TOKEN_SEP_LBRACK, "["
        }
    case '\'', '"':
        return self.line, TOKEN_STRING, self.scanShortString()
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

/**
 * test()判断剩余的源代码是否以某种字符串开头
 */
func (self *Lexer) test(s string) bool {
    return strings.HasPrefix(self.chunk, s)
}

/**
 * [func nextf跳过n个字节]
 * @Author   konyka
 * @DateTime 2019-05-03T13:12:04+0800
 * @param    {[type]}                 self *Lexer)       next(n int [description]
 * @return   {[type]}                      [description]
 */
func (self *Lexer) next(n int) {
    self.chunk = self.chunk[n:]
}


/**
 * isWhiteSpace判断自负是不是空白符
 */
func isWhiteSpace(c byte) bool {
    switch c {
    case '\t', '\n', '\v', '\f', '\r', ' ':
        return true
    }
    return false
}


/**
 * isNewLine判断字符是不是回车或者换行
 */
func isNewLine(c byte) bool {
    return c == '\r' || c == '\n'
}

/**
 * [func 跳过注释]
 * @Author   konyka
 * @DateTime 2019-05-03T13:17:28+0800
 * @param    {[type]}                 self *Lexer)       skipComment( [description]
 * @return   {[type]}                      [description]
 */
func (self *Lexer) skipComment() {
    self.next(2) // skip --

    // long comment ?
    if self.test("[") {
        if reOpeningLongBracket.FindString(self.chunk) != "" {
            self.scanLongString()
            return
        }
    }

    // short comment
    for len(self.chunk) > 0 && !isNewLine(self.chunk[0]) {
        self.next(1)
    }
}





