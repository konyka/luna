/*
* @Author: konyka
* @Date:   2019-05-03 18:11:48
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 18:30:35
*/

package ast

/*
stat ::=  ‘;’ |
     varlist ‘=’ explist |
     functioncall |
     label |
     break |
     goto Name |
     do block end |
     while exp do block end |
     repeat block until exp |
     if exp then block {elseif exp then block} [else block] end |
     for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end |
     for namelist in explist do block end |
     function funcname funcbody |
     local function Name funcbody |
     local namelist [‘=’ explist]
*/
type Stat interface{}


type EmptyStat struct{}              // ‘;’
type BreakStat struct{ Line int }    // break
type LabelStat struct{ Name string } // ‘::’ Name ‘::’
type GotoStat struct{ Name string }  // goto Name
type DoStat struct{ Block *Block }   // do block end
type FuncCallStat = FuncCallExp      // functioncall


/**
 * while exp do block end
 */
type WhileStat struct {
    Exp   Exp
    Block *Block
}

/**
 * repeat block until exp
 */
type RepeatStat struct {
    Block *Block
    Exp   Exp
}

/**
 * if exp then block {elseif exp then block} [else block] end
 */
type IfStat struct {
    Exps   []Exp
    Blocks []*Block
}



