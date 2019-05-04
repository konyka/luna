/*
* @Author: konyka
* @Date:   2019-04-26 10:01:20
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 11:18:37
*/
package main

import "encoding/json"
import "io/ioutil"
import "os"
import "lunago/compiler/parser"

func main() {
    if len(os.Args) > 1 {
        data, err := ioutil.ReadFile(os.Args[1])
        if err != nil {
            panic(err)
        }

        testParser(string(data), os.Args[1])
    }
}

func testParser(chunk, chunkName string) {
    ast := parser.Parse(chunk, chunkName)
    b, err := json.Marshal(ast)
    if err != nil {
        panic(err)
    }
    println(string(b))
}