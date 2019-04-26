/*
* @Author: konyka
* @Date:   2019-04-26 13:25:00
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-26 16:59:56
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





