/*
* @Author: konyka
* @Date:   2019-05-05 10:51:37
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 11:07:29
*/

package stdlib

import "fmt"
import "strconv"
import "strings"
import . "lunago/api"


var baseFuncs = map[string]GoFunction{
    "print":        basePrint,
    "assert":       baseAssert,
    "error":        baseError,
    "select":       baseSelect,
    "ipairs":       baseIPairs,
    "pairs":        basePairs,
    "next":         baseNext,
    "load":         baseLoad,
    "loadfile":     baseLoadFile,
    "dofile":       baseDoFile,
    "pcall":        basePCall,
    "xpcall":       baseXPCall,
    "getmetatable": baseGetMetatable,
    "setmetatable": baseSetMetatable,
    "rawequal":     baseRawEqual,
    "rawlen":       baseRawLen,
    "rawget":       baseRawGet,
    "rawset":       baseRawSet,
    "type":         baseType,
    "tostring":     baseToString,
    "tonumber":     baseToNumber,
    /* placeholders */
    "_G":       nil,
    "_VERSION": nil,
}


func basePrint(ls LuaState) int {
    n := ls.GetTop() /* number of arguments */
    ls.GetGlobal("tostring")
    for i := 1; i <= n; i++ {
        ls.PushValue(-1) /* function to be called */
        ls.PushValue(i)  /* value to print */
        ls.Call(1, 1)
        s, ok := ls.ToStringX(-1) /* get result */
        if !ok {
            return ls.Error2("'tostring' must return a string to 'print'")
        }
        if i > 1 {
            fmt.Print("\t")
        }
        fmt.Print(s)
        ls.Pop(1) /* pop result */
    }
    fmt.Println()
    return 0
}

func OpenBaseLib(ls LuaState) int {
    /* open lib into global table */
    ls.PushGlobalTable()
    ls.SetFuncs(baseFuncs, 0)
    /* set global _G */
    ls.PushValue(-1)
    ls.SetField(-2, "_G")
    /* set global _VERSION */
    ls.PushString("Lua 5.3") // todo
    ls.SetField(-2, "_VERSION")
    return 1
}

/**
 * select (index, ···)
 */
func baseSelect(ls LuaState) int {
    n := int64(ls.GetTop())
    if ls.Type(1) == LUA_TSTRING && ls.CheckString(1) == "#" {
        ls.PushInteger(n - 1)
        return 1
    } else {
        i := ls.CheckInteger(1)
        if i < 0 {
            i = n + i
        } else if i > n {
            i = n
        }
        ls.ArgCheck(1 <= i, 1, "index out of range")
        return int(n - i)
    }
}









