/*
* @Author: konyka
* @Date:   2019-04-26 13:25:00
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-26 14:21:14
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






