/*
* @Author: konyka
* @Date:   2019-04-28 11:43:36
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 11:47:45
*/


package state

/**
 * 返回栈顶索引
 */
func (self *luaState) GetTop() int {
    return self.stack.top
}








