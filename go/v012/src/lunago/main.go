/*
* @Author: konyka
* @Date:   2019-04-26 10:01:20
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 09:12:30
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

func next(ls LuaState) int {
    ls.SetTop(2) /* create a 2nd argument if there isn't one */
    if ls.Next(1) {
        return 2
    } else {
        ls.PushNil()
        return 1
    }
}

func pairs(ls LuaState) int {
    ls.PushGoFunction(next) /* will return generator, */
    ls.PushValue(1)         /* state, */
    ls.PushNil()
    return 3
}


func iPairs(ls LuaState) int {
    ls.PushGoFunction(_iPairsAux) /* iteration function */
    ls.PushValue(1)               /* state */
    ls.PushInteger(0)             /* initial value */
    return 3
}

