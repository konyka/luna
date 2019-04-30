/*
* @Author: konyka
* @Date:   2019-04-28 22:33:14
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-28 22:36:31
*/

package state

/**
 * [push nil to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:35:11+0800
 * @param    {[type]}                 self *luaState)    PushNil( [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushNil() {
    self.stack.push(nil)
}

/**
 * [push boolean to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:35:43+0800
 * @param    {[type]}                 self *luaState)    PushBoolean(b bool [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushBoolean(b bool) {
    self.stack.push(b)
}

/**
 * [push integer to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:36:07+0800
 * @param    {[type]}                 self *luaState)    PushInteger(n int64 [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushInteger(n int64) {
    self.stack.push(n)
}

/**
 * [push number to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:36:17+0800
 * @param    {[type]}                 self *luaState)    PushNumber(n float64 [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushNumber(n float64) {
    self.stack.push(n)
}

/**
 * [push string to stack top]
 * @Author   konyka
 * @DateTime 2019-04-28T22:36:26+0800
 * @param    {[type]}                 self *luaState)    PushString(s string [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) PushString(s string) {
    self.stack.push(s)
}













