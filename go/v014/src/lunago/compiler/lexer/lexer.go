/*
* @Author: konyka
* @Date:   2019-05-03 11:57:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 12:06:01
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







