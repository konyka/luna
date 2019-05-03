/*
* @Author: konyka
* @Date:   2019-05-03 18:11:48
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 18:12:50
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






