/*
* @Author: konyka
* @Date:   2019-05-05 07:44:38
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 07:45:02
*/

package codegen

import . "lunago/binchunk"
import . "lunago/compiler/ast"

func GenProto(chunk *Block) *Prototype {
    fd := &FuncDefExp{
        IsVararg: true,
        Block:    chunk,
    }

    fi := newFuncInfo(nil, fd)
    fi.addLocVar("_ENV")
    cgFuncDefExp(fi, fd, 0)
    return toProto(fi.subFuncs[0])
}












