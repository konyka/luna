/*
* @Author: konyka
* @Date:   2019-05-05 14:50:58
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 16:13:55
*/

package stdlib

import "os"
import "strings"
import . "lunago/api"

/* key, in the registry, for table of loaded modules */
const LUA_LOADED_TABLE = "_LOADED"

/* key, in the registry, for table of preloaded loaders */
const LUA_PRELOAD_TABLE = "_PRELOAD"

const (
    LUA_DIRSEP    = string(os.PathSeparator)
    LUA_PATH_SEP  = ";"
    LUA_PATH_MARK = "?"
    LUA_EXEC_DIR  = "!"
    LUA_IGMARK    = "-"
)

var llFuncs = map[string]GoFunction{
    "require": pkgRequire,
}

var pkgFuncs = map[string]GoFunction{
    "searchpath": pkgSearchPath,
    /* placeholders */
    "preload":   nil,
    "cpath":     nil,
    "path":      nil,
    "searchers": nil,
    "loaded":    nil,
}

func OpenPackageLib(ls LuaState) int {
    ls.NewLib(pkgFuncs) /* create 'package' table */
    createSearchersTable(ls)
    /* set paths */
    ls.PushString("./?.lua;./?/init.lua")
    ls.SetField(-2, "path")
    /* store config information */
    ls.PushString(LUA_DIRSEP + "\n" + LUA_PATH_SEP + "\n" +
        LUA_PATH_MARK + "\n" + LUA_EXEC_DIR + "\n" + LUA_IGMARK + "\n")
    ls.SetField(-2, "config")
    /* set field 'loaded' */
    ls.GetSubTable(LUA_REGISTRYINDEX, LUA_LOADED_TABLE)
    ls.SetField(-2, "loaded")
    /* set field 'preload' */
    ls.GetSubTable(LUA_REGISTRYINDEX, LUA_PRELOAD_TABLE)
    ls.SetField(-2, "preload")
    ls.PushGlobalTable()
    ls.PushValue(-2)        /* set 'package' as upvalue for next lib */
    ls.SetFuncs(llFuncs, 1) /* open lib into global table */
    ls.Pop(1)               /* pop global table */
    return 1                /* return 'package' table */
}

func createSearchersTable(ls LuaState) {
    searchers := []GoFunction{
        preloadSearcher,
        luaSearcher,
    }
    /* create 'searchers' table */
    ls.CreateTable(len(searchers), 0)
    /* fill it with predefined searchers */
    for idx, searcher := range searchers {
        ls.PushValue(-2) /* set 'package' as upvalue for all searchers */
        ls.PushGoClosure(searcher, 1)
        ls.RawSetI(-2, int64(idx+1))
    }
    ls.SetField(-2, "searchers") /* put it in field 'searchers' */
}

func preloadSearcher(ls LuaState) int {
    name := ls.CheckString(1)
    ls.GetField(LUA_REGISTRYINDEX, "_PRELOAD")
    if ls.GetField(-1, name) == LUA_TNIL { /* not found? */
        ls.PushString("\n\tno field package.preload['" + name + "']")
    }
    return 1
}




