/*
* @Author: konyka
* @Date:   2019-04-27 18:09:33
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 16:14:58
*/

package api

const LUA_MINSTACK = 20
const LUAI_MAXSTACK = 1000000
const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000
const LUA_RIDX_GLOBALS int64 = 2

/* basic types */
const (
	LUA_TNONE = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

/* arithmetic functions const */
const (
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // *
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // //
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ~
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPUNM         // -
	LUA_OPBNOT        // ~
)

/* comparison functions const */
const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        // <=
)









