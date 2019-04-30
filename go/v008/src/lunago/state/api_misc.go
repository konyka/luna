/*
* @Author: konyka
* @Date:   2019-04-29 15:24:47
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-30 11:52:08
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
    } else if t, ok := val.(*luaTable); ok {
        self.stack.push(int64(t.len()))
    } else {
        panic("length error!")
    }
}

/**
 * [ 该方法从栈顶pop n 个值，然后对这些值进行拼接，然后把结果push 到栈顶 ]
 * @Author   konyka
 * @DateTime 2019-04-29T15:30:38+0800
 * @param    {[type]}                 self *luaState)    Concat(n int [description]
 * @return   {[type]}                      [description]
 */
func (self *luaState) Concat(n int) {
    if n == 0 {
        self.stack.push("")
    } else if n >= 2 {
        for i := 1; i < n; i++ {
            if self.IsString(-1) && self.IsString(-2) {
                s2 := self.ToString(-1)
                s1 := self.ToString(-2)
                self.stack.pop()
                self.stack.pop()
                self.stack.push(s1 + s2)
                continue
            }

            panic("concatenation error!")
        }
    }
    // n == 1, do nothing
}





