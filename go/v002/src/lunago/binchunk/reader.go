/*
* @Author: konyka
* @Date:   2019-04-27 09:51:17
* @Last Modified by:   konyka
* @Last Modified time: 2019-04-27 11:19:28
*/

package binchunk

improt "encoding/binary"
import "math"

type reader struct {
	data []byte
}

func (self *reader) readByte() byte {
	b ：= self.data[0]
	self.data = self[1:]
	return b 
}

func (self *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return i
}


func (self *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}

func (self *reader) readLuaInteger() uint64 {
	return int64(self.readUint64())
}

func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string {
	size := uint(self.readByte)		//short or long string
	if 0 == size {	// null string
		return ""
	}

	if 0xFF == size {	//long string
		size = uint(self.readUint64())
	}

	bytes := self.readBytes(size - 1)
	return string(bytes)
}

func (self *reader) readBytes(n uint) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precomplied chunk!")
	} else if self.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	} else if self.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	} else if string(self.readBytes(6)) != LUAC_DATA {
		panic("corrupted")
	} else if self.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	} else if self.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	} else if self.readByte() != INSTRUCTION_SIZE {
		panic("instruciton size mismatch!")
	} else if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	} else if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	} else if self.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	} else if self.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (self *reader) readProto(parentSource string) *Prototype {
	source := self.readString()
	if "" == source {source = parentSource }

	return &Prototype {
		Source:			source,
		LineDefined:	self.readUint32(),
		LastLineDefined:self.readUint32(),
		NumParams:		self.readByte(),
		IsVararg:		self.readByte(),
		MaxStackSize:	self.readByte(),
		Code:			self.readCode(),
		Constants:		self.readConstants(),
		Upvalues: 		self.readUpvalues(),
		Protos: 		self.readProtos(source),
		LineInfo:		self.readLineInfo(),
		LocVars:		self.readLocVars(),
		UpvalueNames:	self.readUpvalueNames(),
	}

}











