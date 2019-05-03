/*
* @Author: konyka
* @Date:   2019-04-26 10:01:20
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 14:41:40
*/
package main

import "fmt"
import "io/ioutil"
import "os"
import . "lunago/compiler/lexer"

func main() {
    if len(os.Args) > 1 {
        data, err := ioutil.ReadFile(os.Args[1])
        if err != nil {
            panic(err)
        }

        testLexer(string(data), os.Args[1])
    }
}

func testLexer(chunk, chunkName string) {
    lexer := NewLexer(chunk, chunkName)
    for {
        line, kind, token := lexer.NextToken()
        fmt.Printf("[%2d] [%-10s] %s\n",
            line, kindToCategory(kind), token)
        if kind == TOKEN_EOF {
            break
        }
    }
}


func kindToCategory(kind int) string {
    switch {
    case kind < TOKEN_SEP_SEMI:
        return "other"
    case kind <= TOKEN_SEP_RCURLY:
        return "separator"
    case kind <= TOKEN_OP_NOT:
        return "operator"
    case kind <= TOKEN_KW_WHILE:
        return "keyword"
    case kind == TOKEN_IDENTIFIER:
        return "identifier"
    case kind == TOKEN_NUMBER:
        return "number"
    case kind == TOKEN_STRING:
        return "string"
    default:
        return "other"
    }
}