/*
* @Author: konyka
* @Date:   2019-05-03 11:57:34
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 11:59:23
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
    nextToken     string
    nextTokenKind int
    nextTokenLine int
}









