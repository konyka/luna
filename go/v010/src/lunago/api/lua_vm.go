/*
* @Author: konyka
* @Date:   2019-04-29 16:42:33
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-02 10:32:49
*/

package api

type LuaVM interface {
    LuaState
    PC() int              //返回当前PC
    AddPC(n int)          //修改PC 用于实现跳转指令
    Fetch() uint32        //取出当前的指令，将PC指向下一条指令
    GetConst(idx int)     //将指定的常量push到栈顶
    GetRK(rk int)         //将指定的常量或者栈值push到栈顶
    RegisterCount() int     //当前lua函数所操作的寄存器计数器
    LoadVararg(n int)       //把传递给当前lua函数的变长参数push到栈顶 多退少补
    LoadProto(idx int)      //把当前lua函数的子函数的原型 实例化为闭包 ，并push到栈顶
    CloseUpvalues(a int)
}








