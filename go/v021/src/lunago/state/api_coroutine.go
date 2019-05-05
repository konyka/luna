/*
* @Author: konyka
* @Date:   2019-05-05 19:29:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 19:35:20
*/

package state

import . "lunago/api"


func (self *luaState) NewThread() LuaState {
    t := &luaState{registry: self.registry}
    t.pushLuaStack(newLuaStack(LUA_MINSTACK, t))
    self.stack.push(t)
    return t
}









