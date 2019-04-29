/*
* @Author: konyka
* @Date:   2019-04-29 15:24:47
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-29 15:25:20
*/

package state

/**
 * [func description]
 * @Author   konyka
 * @DateTime 2019-04-29T15:25:16+0800
 * @param    {[type]}                 self *luaState)    Len(idx int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Len(idx int) {
    val := self.stack.get(idx)

    if s, ok := val.(string); ok {
        self.stack.push(int64(len(s)))
    } else {
        panic("length error!")
    }
}







