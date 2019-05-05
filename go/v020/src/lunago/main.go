/*
* @Author: konyka
* @Date:   2019-04-26 10:01:20
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 11:12:34
*/
package main

import "os"
import "lunago/state"

func main() {
    if len(os.Args) > 1 {
        ls := state.New()
        ls.OpenLibs()
        ls.LoadFile(os.Args[1])
        ls.Call(0, -1)
    }
}