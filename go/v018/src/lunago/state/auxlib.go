/*
* @Author: konyka
* @Date:   2019-05-05 09:40:08
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 09:46:55
*/

package state

import "fmt"
import "io/ioutil"
import . "lunago/api"

import "lunago/stdlib"


func (self *luaState) TypeName2(idx int) string {
    return self.TypeName(self.Type(idx))
}


func (self *luaState) Len2(idx int) int64 {
    self.Len(idx)
    i, isNum := self.ToIntegerX(-1)
    if !isNum {
        self.Error2("object length is not an integer")
    }
    self.Pop(1)
    return i
}





























