/*
* @Author: konyka
* @Date:   2019-04-30 12:34:22
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-02 15:33:30
*/

package state

import . "lunago/api"

/**
 * [func SetTable()把键值对写入表。其中key和value从栈中弹出，表则位于指定的索引处。]
 * @Author   konyka
 * @DateTime 2019-04-30T12:37:29+0800
 * @param    {[type]}                 self *luaState)    SetTable(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) SetTable(idx int) {
    t := self.stack.get(idx)
    v := self.stack.pop()
    k := self.stack.pop()
    self.setTable(t, k, v, false)
}


/**
 * [func 表的逻辑提取成setTable(）方法 t[k]=v]
 * @Author   konyka
 * @DateTime 2019-04-30T12:39:30+0800
 * @param    {[type]}                 self *luaState)    setTable(t, k, v luaValue [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) setTable(t, k, v luaValue, raw bool) {
    if tbl, ok := t.(*luaTable); ok {
        if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
            tbl.put(k, v)
            return
        }
    }

    if !raw {
        if mf := getMetafield(t, "__newindex", self); mf != nil {
            switch x := mf.(type) {
            case *luaTable:
                self.setTable(x, k, v, false)
                return
            case *closure:
                self.stack.push(mf)
                self.stack.push(t)
                self.stack.push(k)
                self.stack.push(v)
                self.Call(3, 0)
                return
            }
        }
    }

    panic("index error!")
}

/**
 * [func SetField（）和SetTable()类似，只不过key不是从栈顶弹出的任意值，而是由参数传入的字符串。
 * 用于给记录的字段赋值。执行后，value从栈顶弹出，并被赋值给记录的相应字段。]
 * @Author   konyka
 * @DateTime 2019-04-30T12:45:30+0800
 * @param    {[type]}                 self *luaState)    SetField(idx int, k string [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) SetField(idx int, k string) {
    t := self.stack.get(idx)
    v := self.stack.pop()
    self.setTable(t, k, v, false)
}

/**
 * [func SetI() 和SetField（）类似，只不过由参数传入的key是数组，而非字符串，用于按照索引修改数组的元素。
 *  执行之后，值从栈顶弹出，并被写到数组中。]
 * @Author   konyka
 * @DateTime 2019-04-30T12:55:27+0800
 * @param    {[type]}                 self *luaState)    SetI(idx int, i int64 [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) SetI(idx int, i int64) {
    t := self.stack.get(idx)
    v := self.stack.pop()
    self.setTable(t, i, v, false)
}

/**
 * [func 往全局环境里面写入一个值，其中的字段名由参数指定，值从栈顶弹出。]
 * @Author   konyka
 * @DateTime 2019-05-01T17:47:03+0800
 * @param    {[type]}                 self *luaState)    SetGlobal(name string [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) SetGlobal(name string) {
    t := self.registry.get(LUA_RIDX_GLOBALS)
    v := self.stack.pop()
    self.setTable(t, name, v)
}

/**
 * [ 用于给全局环境注册go函数值，仅仅用于操作全局环境，字段名以及go函数从参数传入，不改变lua栈的状态。 ]
 * @Author   konyka
 * @DateTime 2019-05-01T18:19:35+0800
 * @param    {[type]}                 self *luaState)    Register(name string, f GoFunction [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Register(name string, f GoFunction) {
    self.PushGoFunction(f)
    self.SetGlobal(name)
}




