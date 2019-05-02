/*
* @Author: konyka
* @Date:   2019-04-26 10:01:20
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-02 16:41:39
*/
package main

import "io/ioutil"
import "os"
import "lunago/state"
import "fmt"
import . "lunago/api"



func main() {
    if len(os.Args) > 1 {
        data, err := ioutil.ReadFile(os.Args[1])
        if err != nil {
            panic(err)
        }

        ls := state.New()
        ls.Register("print", print)
        ls.Register("getmetatable", getMetatable) //register
        ls.Register("setmetatable", setMetatable) //register
        ls.Load(data, os.Args[1], "b")
        ls.Call(0, 0)
    }
}


func print(ls LuaState) int {
    nArgs := ls.GetTop()
    for i := 1; i <= nArgs; i++ {
        if ls.IsBoolean(i) {
            fmt.Printf("%t", ls.ToBoolean(i))
        } else if ls.IsString(i) {
            fmt.Print(ls.ToString(i))
        } else {
            fmt.Print(ls.TypeName(ls.Type(i)))
        }
        if i < nArgs {
            fmt.Print("\t")
        }
    }
    fmt.Println()
    return 0
}


func getMetatable(ls LuaState) int {
    if !ls.GetMetatable(1) {
        ls.PushNil()
    }
    return 1
}

func setMetatable(ls LuaState) int {
    ls.SetMetatable(1)
    return 1
}