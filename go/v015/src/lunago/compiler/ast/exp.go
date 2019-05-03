/*
* @Author: konyka
* @Date:   2019-05-03 19:24:52
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 19:44:02
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

type NilExp struct{ Line int }    // nil
type TrueExp struct{ Line int }   // true
type FalseExp struct{ Line int }  // false
type VarargExp struct{ Line int } // ...

// Numeral
type IntegerExp struct {
    Line int
    Val  int64
}
type FloatExp struct {
    Line int
    Val  float64
}

// LiteralString
type StringExp struct {
    Line int
    Str  string
}

type NameExp struct {
    Line int
    Name string
}

/**
 * unop exp
 */
type UnopExp struct {
    Line int // line of operator
    Op   int // operator
    Exp  Exp
}

/**
 * exp1 op exp2
 */
type BinopExp struct {
    Line int // line of operator
    Op   int // operator
    Exp1 Exp
    Exp2 Exp
}






















