/*
* @Author: konyka
* @Date:   2019-05-04 14:05:36
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-04 14:06:33
*/


package codegen

import . "lunago/compiler/ast"

func cgBlock(fi *funcInfo, node *Block) {
    for _, stat := range node.Stats {
        cgStat(fi, stat)
    }

    if node.RetExps != nil {
        cgRetStat(fi, node.RetExps)
    }
}















