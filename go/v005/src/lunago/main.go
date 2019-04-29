/*
* @Author: konyka
* @Date:   2019-04-26 10:01:20
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 15:39:34
*/
package main

import "fmt"
import . "lunago/api"
import "lunago/state"
import _ "lunago/binchunk"


func main() {
	ls := state.New()
	
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	printStack(ls)

	ls.Arith(LUA_OPADD)
	printStack(ls)
	ls.Arith(LUA_OPBNOT)
	printStack(ls)
	ls.Len(2)
	printStack(ls)
	ls.Concat(3)
	printStack(ls)
	ls.PushBoolean(ls.Compare(1, 2, LUA_OPEQ))
	printStack(ls)

}

/**
 * [print stack info]
 * @Author   konyka
 * @DateTime 2019-04-29T10:03:13+0800
 * @param    {[type]}                 ls LuaState      [description]
 * @return   {[type]}                    [description]
 */
func printStack(ls LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}









