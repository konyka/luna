/*
* @Author: konyka
* @Date:   2019-04-30 18:39:45
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 18:40:42
*/

package state

import "fmt"
import "luago/binchunk"
import "luago/vm"

/**
 * 加载chunk
 */
func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
    proto := binchunk.Undump(chunk) // todo
    c := newLuaClosure(proto)
    self.stack.push(c)
    return 0
}












