/*
* @Author: konyka
* @Date:   2019-05-03 19:24:52
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 19:26:26
*/

package ast

/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
     prefixexp | tableconstructor | exp binop exp | unop exp
prefixexp ::= var | functioncall | ‘(’ exp ‘)’
var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
*/

type Exp interface{}








