/*
* @Author: konyka
* @Date:   2019-04-30 12:34:22
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 12:37:44
*/

package state

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
    self.setTable(t, k, v)
}















