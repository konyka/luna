/*
* @Author: konyka
* @Date:   2019-04-26 13:25:00
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-27 13:33:11
*/
package binchunk

type  binaryChunk struct{
	header		//头部
	sizeUpvalues byte    //主函数Upvlue数量
	mainFunc	*Prototype   //主函数原型
}


type header struct{
	signature	[4]byte
	version		byte
	format		byte
	luacData	[6]byte
	cintSize	byte
	sizetSize	byte
	instructionSize	byte
	luaIntegerSize	byte
	luaNumberSize	byte
	luacInt	int64
	luacNum	float64
}

const (
	LUA_SIGNATURE		= "\x1bLua"
	LUAC_VERSION		= 0x53
	LUAC_FORMAT			= 0x00
	LUAC_DATA			= "\x19\x93\r\n\x1a\n"
	CINT_SIZE			= 4
	CSIZET_SIZE			= 8
	INSTRUCTION_SIZE	= 4
	LUA_INTEGER_SIZE	= 8
	LUA_NUMBER_SIZE		= 8
	LUAC_INT			= 0x5678
	LUAC_NUM			=370.5
)


type Prototype struct {
	Source		string
	LineDefined 	uint32
	LastLineDefined	uint32
	NumParams		byte
	IsVararg		byte
	MaxStackSize	byte
	Code			[]uint32
	Constants		[]interface{}
	Upvalues 		[]Upvalue
	Protos 			[]*Prototype
	LineInfo		[]uint32
	LocVars			[]LocVar
	UpvalueNames	[]string
}

const (
	TAG_NIL		= 0x00
	TAG_BOOLEAN	= 0x01
	TAG_NUMBER	= 0x03
	TAG_INTEGER	= 0x13
	TAG_SHORT_STR	= 0x04
	TAG_LONG_STR	= 0x14
)

type Upvalue struct {
	Instack	byte
	Idx		byte
}

type LocVar struct {
	VarName 	string
	StartPC		uint32
	EndPC		uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()	//检查头部
	reader.readByte()		//跳过Upvalue的数量
	return reader.readProto("")	//读取函数原型信息
}

