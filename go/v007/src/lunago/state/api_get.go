/*
* @Author: konyka
* @Date:   2019-04-30 12:01:30
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 12:07:47
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







