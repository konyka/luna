/*
* @Author: konyka
* @Date:   2019-04-29 16:42:33
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-01 10:26:59
*/

package api

type LuaVM interface {
    LuaState
    PC() int              //返回当前PC
    AddPC(n int)          //修改PC 用于实现跳转指令
    Fetch() uint32        //取出当前的指令，将PC指向下一条指令
    GetConst(idx int)     //将指定的常量push到栈顶
    GetRK(rk int)         //将指定的常量或者栈值push到栈顶
    RegisterCount() int     //计数器
    LoadVararg(n int)       //加载vararg
    LoadProto(idx int)      //加载原型
}








