/*
* @Author: konyka
* @Date:   2019-04-27 18:15:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 09:03:58
*/

package state

type luaStack struct {
	slots	[]luaValue	//用来存放值
	top		int 		//记录栈顶的索引
}

//创建指定容量的栈
func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}











