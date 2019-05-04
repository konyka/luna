/*
* @Author: konyka
* @Date:   2019-05-05 07:31:44
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 07:41:51
*/


package codegen

import . "lunago/binchunk"

func toProto(fi *funcInfo) *Prototype {
    proto := &Prototype{
        NumParams:    byte(fi.numParams),
        MaxStackSize: byte(fi.maxRegs),
        Code:         fi.insts,
        Constants:    getConstants(fi),
        Upvalues:     getUpvalues(fi),
        Protos:       toProtos(fi.subFuncs),
        LineInfo:     []uint32{}, // debug
        LocVars:      []LocVar{}, // debug
        UpvalueNames: []string{}, // debug
    }

    if proto.MaxStackSize < 2 {
        proto.MaxStackSize = 2 // todo
    }
    if fi.isVararg {
        proto.IsVararg = 1 // todo
    }

    return proto
}

func toProtos(fis []*funcInfo) []*Prototype {
    protos := make([]*Prototype, len(fis))
    for i, fi := range fis {
        protos[i] = toProto(fi)
    }
    return protos
}

func getConstants(fi *funcInfo) []interface{} {
    consts := make([]interface{}, len(fi.constants))
    for k, idx := range fi.constants {
        consts[idx] = k
    }
    return consts
}















