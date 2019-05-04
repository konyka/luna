/*
* @Author: konyka
* @Date:   2019-05-05 07:50:00
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 07:50:17
*/

package compiler

import "lunago/binchunk"
import "lunago/compiler/codegen"
import "lunago/compiler/parser"

func Compile(chunk, chunkName string) *binchunk.Prototype {
    ast := parser.Parse(chunk, chunkName)
    return codegen.GenProto(ast)
}











