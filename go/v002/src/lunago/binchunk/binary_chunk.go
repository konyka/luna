/*
* @Author: konyka
* @Date:   2019-04-26 13:25:00
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-26 13:28:15
*/
package binchunk

type  binaryChunk struct{
	header		//头部
	sizeUpvalues byte    //主函数Upvlue数量
	mainFunc	*Prototype   //主函数原型
}


