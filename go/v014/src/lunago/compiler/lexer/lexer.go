/*
* @Author: konyka
* @Date:   2019-05-03 11:57:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 13:50:25
*/

package lexer

import "bytes"
import "fmt"
import "regexp"
import "strconv"
import "strings"

var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)
var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")
var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

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


func (self *Lexer) scanLongString() string {
    openingLongBracket := reOpeningLongBracket.FindString(self.chunk)
    if openingLongBracket == "" {
        self.error("invalid long string delimiter near '%s'",
            self.chunk[0:2])
    }

    closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
    closingLongBracketIdx := strings.Index(self.chunk, closingLongBracket)
    if closingLongBracketIdx < 0 {
        self.error("unfinished long string or comment")
    }

    str := self.chunk[len(openingLongBracket):closingLongBracketIdx]
    self.next(closingLongBracketIdx + len(closingLongBracket))

    str = reNewLine.ReplaceAllString(str, "\n")
    self.line += strings.Count(str, "\n")
    if len(str) > 0 && str[0] == '\n' {
        str = str[1:]
    }

    return str
}

func (self *Lexer) error(f string, a ...interface{}) {
    err := fmt.Sprintf(f, a...)
    err = fmt.Sprintf("%s:%d: %s", self.chunkName, self.line, err)
    panic(err)
}

func (self *Lexer) scanShortString() string {
    if str := reShortStr.FindString(self.chunk); str != "" {
        self.next(len(str))
        str = str[1 : len(str)-1]
        if strings.Index(str, `\`) >= 0 {
            self.line += len(reNewLine.FindAllString(str, -1))
            str = self.escape(str)
        }
        return str
    }
    self.error("unfinished string")
    return ""
}


func (self *Lexer) escape(str string) string {
    var buf bytes.Buffer

    for len(str) > 0 {
        if str[0] != '\\' {
            buf.WriteByte(str[0])
            str = str[1:]
            continue
        }

        if len(str) == 1 {
            self.error("unfinished string")
        }

        switch str[1] {
        case 'a':
            buf.WriteByte('\a')
            str = str[2:]
            continue
        case 'b':
            buf.WriteByte('\b')
            str = str[2:]
            continue
        case 'f':
            buf.WriteByte('\f')
            str = str[2:]
            continue
        case 'n', '\n':
            buf.WriteByte('\n')
            str = str[2:]
            continue
        case 'r':
            buf.WriteByte('\r')
            str = str[2:]
            continue
        case 't':
            buf.WriteByte('\t')
            str = str[2:]
            continue
        case 'v':
            buf.WriteByte('\v')
            str = str[2:]
            continue
        case '"':
            buf.WriteByte('"')
            str = str[2:]
            continue
        case '\'':
            buf.WriteByte('\'')
            str = str[2:]
            continue
        case '\\':
            buf.WriteByte('\\')
            str = str[2:]
            continue
        case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
            if found := reDecEscapeSeq.FindString(str); found != "" {
                d, _ := strconv.ParseInt(found[1:], 10, 32)
                if d <= 0xFF {
                    buf.WriteByte(byte(d))
                    str = str[len(found):]
                    continue
                }
                self.error("decimal escape too large near '%s'", found)
            }
        case 'x': // \xXX
            if found := reHexEscapeSeq.FindString(str); found != "" {
                d, _ := strconv.ParseInt(found[2:], 16, 32)
                buf.WriteByte(byte(d))
                str = str[len(found):]
                continue
            }
        case 'u': // \u{XXX}
            if found := reUnicodeEscapeSeq.FindString(str); found != "" {
                d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
                if err == nil && d <= 0x10FFFF {
                    buf.WriteRune(rune(d))
                    str = str[len(found):]
                    continue
                }
                self.error("UTF-8 value too large near '%s'", found)
            }
        case 'z':
            str = str[2:]
            for len(str) > 0 && isWhiteSpace(str[0]) { // todo
                str = str[1:]
            }
            continue
        }
        self.error("invalid escape sequence near '\\%c'", str[1])
    }

    return buf.String()
}




