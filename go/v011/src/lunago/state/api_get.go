/*
* @Author: konyka
* @Date:   2019-04-30 12:01:30
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-02 15:32:28
*/

package state

import . "lunago/api"


/**
 * [func CreateTable（）创建一个空的表，将其push到栈顶。
 * 该方法提供了两个参数，用来指定数组部分和哈希表部分的初始容量。
 * 如果可以预先估计出表的使用方式和容量，
 * 那么可以使用这两个参数在创建表的时候预先分配足够的空间，用来避免后续对表进行频繁的扩容]
 * @Author   konyka
 * @DateTime 2019-04-30T12:07:13+0800
 * @param    {[type]}                 self *luaState)    CreateTable(nArr, nRec int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) CreateTable(nArr, nRec int) {
    t := newLuaTable(nArr, nRec)
    self.stack.push(t)
}

/**
 * [func 如果无法预先估计表的用法和容量，
 * 可以使用NewTable()创建表。NewTable()只是CreateTable（）的特殊情况。]
 * @Author   konyka
 * @DateTime 2019-04-30T12:10:01+0800
 * @param    {[type]}                 self *luaState)    NewTable( [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) NewTable() {
    self.CreateTable(0, 0)
}

/**
 * 根据key（从栈顶弹出）从表（索引由参数指定）里面取值，然后把值push到栈顶，并返回值的类型
 */
func (self *luaState) GetTable(idx int) LuaType {
    t := self.stack.get(idx)
    k := self.stack.pop()
    return self.getTable(t, k, false)
}


 
/**
 * 为了减少重复，把根据key从table里面获取值的逻辑提取为函数getTable(t, k luaValue)：
 * push(t[k])
 */
func (self *luaState) getTable(t, k luaValue, raw bool) LuaType {
    if tbl, ok := t.(*luaTable); ok {
        v := tbl.get(k)
        if raw || v != nil || !tbl.hasMetafield("__index") {
            self.stack.push(v)
            return typeOf(v)
        }
    }

    if !raw {
        if mf := getMetafield(t, "__index", self); mf != nil {
            switch x := mf.(type) {
            case *luaTable:
                return self.getTable(x, k, false)
            case *closure:
                self.stack.push(mf)
                self.stack.push(t)
                self.stack.push(k)
                self.Call(2, 1)
                v := self.stack.get(-1)
                return typeOf(v)
            }
        }
    }

    panic("index error!")
}

/**
 * GetField（）用来从记录中获取字段。
 */
func (self *luaState) GetField(idx int, k string) LuaType {
    t := self.stack.get(idx)
    return self.getTable(t, k, false)
}

/**
 * GetI（）这个方法是专门给数组准备的，用来根据索引获取数组的元素，执行后，相应的数组元素被push到栈顶。
 */
func (self *luaState) GetI(idx int, i int64) LuaType {
    t := self.stack.get(idx)
    return self.getTable(t, i, false)
}

func (self *luaState) GetGlobal(name string) LuaType {
    t := self.registry.get(LUA_RIDX_GLOBALS)
    return self.getTable(t, name)
}




