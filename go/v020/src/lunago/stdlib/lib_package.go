/*
* @Author: konyka
* @Date:   2019-05-05 14:50:58
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 16:04:50
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










