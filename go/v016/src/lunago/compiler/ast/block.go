/*
* @Author: konyka
* @Date:   2019-05-03 17:56:56
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-03 17:58:47
*/

package ast


/**
 * chunk ::= block
 * block ::= {stat} [retstat]
 * retstat ::= return [explist] [‘;’]
 * explist ::= exp {‘,’ exp}
 */
type Block struct {
    LastLine int
    Stats    []Stat
    RetExps  []Exp
}









