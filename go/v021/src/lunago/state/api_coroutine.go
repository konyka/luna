/*
* @Author: konyka
* @Date:   2019-05-05 19:29:13
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 19:49:48
*/

package state

import . "lunago/api"


func (self *luaState) NewThread() LuaState {
    t := &luaState{registry: self.registry}
    t.pushLuaStack(newLuaStack(LUA_MINSTACK, t))
    self.stack.push(t)
    return t
}

func (self *luaState) Resume(from LuaState, nArgs int) int {
    lsFrom := from.(*luaState)
    if lsFrom.coChan == nil {
        lsFrom.coChan = make(chan int)
    }

    if self.coChan == nil {
        // start coroutine
        self.coChan = make(chan int)
        self.coCaller = lsFrom
        go func() {
            self.coStatus = self.PCall(nArgs, -1, 0)
            lsFrom.coChan <- 1
        }()
    } else {
        // resume coroutine
        if self.coStatus != LUA_YIELD { // todo
            self.stack.push("cannot resume non-suspended coroutine")
            return LUA_ERRRUN
        }
        self.coStatus = LUA_OK
        self.coChan <- 1
    }

    <-lsFrom.coChan // wait coroutine to finish or yield
    return self.coStatus
}


func (self *luaState) Yield(nResults int) int {
    if self.coCaller == nil { // todo
        panic("attempt to yield from outside a coroutine")
    }
    self.coStatus = LUA_YIELD
    self.coCaller.coChan <- 1
    <-self.coChan
    return self.GetTop()
}


func (self *luaState) Status() int {
    return self.coStatus
}

/**
 * debug
 */
func (self *luaState) GetStack() bool {
    return self.stack.prev != nil
}

func (self *luaState) IsYieldable() bool {
    if self.isMainThread() {
        return false
    }
    return self.coStatus != LUA_YIELD // todo
}



