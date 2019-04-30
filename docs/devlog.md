

二进制chunk文件格式

1、 lua源文件编译流程

    lua源文件-->lua编译器（luac）-->二进制字节码文件chunk（.out）-->lua解释器-->输出结果


2、luac命令
    作用：a、将lua源文件编译为二进制字节码chunk
         b、反编译二进制字节码chunk，分析其中的内容，并将信息输出到控制台。

   luac的用法为：

   usage: luac [options] [filenames]
    Available options are:
  -l       list (use -l -l for full listing)
  -o name  output to file 'name' (default is "luac.out")
  -p       parse only
  -s       strip debug information
  -v       show version information
  --       stop handling options
  -        stop handling options and process stdin

3、编译lua源文件
    luac lua源文件列表（一个或者多个）
    默认生成的.out就是二进制字节码文件chunk，可以通过制定-o 选项制定输出的二进制chunk文件名。
    其中chunk文件默认会包含调试信息--行号、变量名......，可以使用-s选项取消输出其中的调试信息。
    同时，如果仅仅是想检查语法的生正确性，而不想生成chunk文件的话，可以使用-p选项进行。

    luac lua文件 // 生成。out文件
    luac -o 制定的chunk文件名 lua源码文件 //生成参数指定名称的chunk
    luac -p lua源文件 //检查语法
    luac -s //取消输出chunk中的调试信息

    luac按照函数为单位进行编译，每个函数都会被编译为一个对应的内部数据结构----Prototype。
    其中包含：
        函数的基本信息，比如参数的数量，局部变量的数量
        字节码
        常量表
        Upvalue表
        调试信息
        子函数原型列表

    对于没有定义任何函数的lua源文件，luac会自动帮我们定义一个main函数，然后把源文件中的代码都放到
    这个main函数中，这样就有了一个函数，就可以进行编译了。如下所示其过程

    print("hello world!!")--luac---> function main(...)
                                        print("hello world!!")    
                                        return
                                     end


   这个main函数即是编译的入口，同时还是lua虚拟机执行程序的入口函数。编译成函数原型以后，编译器luac还会
   为其增加头部信息，然后和原形共同组合成二进制chunk文件。

                二进制chunk结构

    头部header        头部header
    函数原型prototy-> 基本信息basic info
            ^        字节码 bytecodes
            |        常量池 constants
            |        upvalues
            |        调试信息 debug info
            +--------子函数列表 sub functions


4、查看二进制chunk信息
    luac -l chunk文件名或者lua源码文件 //查看文件名指定的chunk信息，源码会自动编译为chunk，在查看
    比如 ：luac -l xxx.out

    出入如下信息：
    main <helloworld.lua:0,0> (4 instructions at 0x7fd05dd03140)
    0+ params, 2 slots, 1 upvalue, 0 locals, 2 constants, 0 functions
    1   [6] GETTABUP    0 0 -1  ; _ENV "print"
    2   [6] LOADK       1 -2    ; "hello world\239\188\129\239\188\129\239\188\129"
    3   [6] CALL        0 2 1
    4   [6] RETURN      0 1

    因为我们的程序只有一个打印语句，没有任何函数，因此编译出来的二进制chunk也只有一个自动添加的main函数原型，没有子函数原型。如果我们的代码有函数定义，luac反编译工具会按照顺序依次输出这些函数原型的信息。例如如下代码：

    function aaa( ... )
        -- body
        function bbb( ... )
            -- body
        end
    end

    反编译查看chunk信息


    luac -l aaa_bbb.lua 

    main <aaa_bbb.lua:0,0> (3 instructions at 0x7f9622c03610)
    0+ params, 2 slots, 1 upvalue, 0 locals, 1 constant, 1 function
    1   [11]    CLOSURE     0 0 ; 0x7f9622c03770
    2   [6] SETTABUP    0 -1 0  ; _ENV "aaa"
    3   [11]    RETURN      0 1

    function <aaa_bbb.lua:6,11> (3 instructions at 0x7f9622c03770)
    0+ params, 2 slots, 1 upvalue, 0 locals, 1 constant, 1 function
    1   [10]    CLOSURE     0 0 ; 0x7f9622c038f0
    2   [8] SETTABUP    0 -1 0  ; _ENV "bbb"
    3   [11]    RETURN      0 1

    function <aaa_bbb.lua:8,10> (1 instruction at 0x7f9622c038f0)
    0+ params, 2 slots, 0 upvalues, 0 locals, 0 constants, 0 functions
    1   [10]    RETURN      0 1


    顶部的两行是函数基本信息，下面的是指令列表
    第一行
    main <aaa_bbb.lua:0,0> (3 instructions at 0x7f9622c03610)

    格式为：function <函数所在的文件ing：函数的开始行号，函数结束行号>
    （指令数量n instructions，at 地址）

    如果是main开始的函数，说明这是编译器自动添加的入口函数开始行号和结束行号都定义为0。
    function 表示这是一个函数

    第二行
    0+ params, 2 slots, 1 upvalue, 0 locals, 1 constant, 1 function

    0+ params ：指明了参数的数量，，其中+表示这是一个可变长参数vararg
    2 slots：运行函数所必须的寄存器数量
    1 upvalue：upvalue的数量
     0 locals：局部变量的数量
     1 constant：常量的数量
     1 function：子函数的数量

     指令列表
     1   [11]    CLOSURE     0 0 ; 0x7f9622c03770
     2   [6] SETTABUP    0 -1 0  ; _ENV "aaa"
     3   [11]    RETURN      0 1

     格式为 指令编号 [对应的行号] 操作吗 操作数；注释信息


     以上的输出是精简模式，可以使用-l -l 选项输出详细信息，luac会把常量表，局部变量以及upvalue表的
     信息打印出来。

     比如
     $ luac -l -l helloworld.lua 

    main <helloworld.lua:0,0> (4 instructions at 0x7fce5a403610)
    0+ params, 2 slots, 1 upvalue, 0 locals, 2 constants, 0 functions
    1   [6] GETTABUP    0 0 -1  ; _ENV "print"
    2   [6] LOADK       1 -2    ; "hello world\239\188\129\239\188\129\239\188\129"
    3   [6] CALL        0 2 1
    4   [6] RETURN      0 1
    constants (2) for 0x7fce5a403610:
    1   "print"
    2   "hello world\239\188\129\239\188\129\239\188\129"
    locals (0) for 0x7fce5a403610:
    upvalues (1) for 0x7fce5a403610:
    0   _ENV    1   0


5、二进制chunk
    a、lua的二进制chunk格式并没有规范可循，而是完全依赖内部实现，并没有进行标准化，官方的文档也
       没有文档说明，一切以官方的源码为准。
    b、chunk的设计并没有有考虑跨平台的场景。对于多个字节表示的数据，需要考虑字节序的问题。lua的
       官方实现方法为：编译源码的时候，直接按照本地的环境生成二进制chunk，当加载chunk的时候，会
       检测文件的字节序，如果和本机不符，就拒绝加载。
    c、chunk格式的设计没有考虑不同lua版本之间的兼容性问题。lua官方的做法为：编译lua源文件的时候，
       直接按照当时的lua版本生成二进制chunk文件，当对其进行加载的时候，会检测被加载的文件的版本号，
       如果和当前lua版本不匹配，则拒绝加载。
    d、二进制chunk的文件格式并没有进行紧凑性的设计。有的时候，生成的二进制chunk文件比源文件还要大。


   二进制chunk文件中的数据类型有三种：数字 、字符串以及列表。

   1、数字
      包括字节、c语言中的整型、c语言的size_t类型、Lua整型、Lua浮点数。
      使用场景
      字节：用来存放小的整数，例如版本号、函数参数个数。
      c语言中的整型int：用来表示列表的长度。
      size_t:用来表示长字符串的长度。
      lua整型以及lua浮点数：主要用在常量表中。用来记载源码中的整数以及浮点常量。

      数字类型具体的占用字节数量纪录在头部header中。

      数据类型                 c类型               go类型        占用的字节数
      字节            lu_Byte(unsigned char)      byte            1
      c语言整型        int                         uint32          4
      c语言size_t     size_t                      uint64          8
      lua整型         lua_Integer（long long）     int64           8
      lua浮点数       lua_Number（double）         float64         8

    2、字符串

      在chunk中，字符串的表现形式为一个字节数组。字符串除了c语言中的字符串，还有一个用来表示长度的字段。
      a、对于null字符串，直接使用0x00表示。
      b、长度<=0xFD（253）的字符串，使用一个字节纪录其长度+1，然后是具体的字节数组，+1是为了保存字符串
         结束符\0。
      c、长度>254(0xfe)的字符串，第一个字节是0xFF,后跟一个size_t，记录长度+1，后跟字节数组。

      null       0x00
      n<=253     n+1 字节数组
      n>=254     0xFF (size_t n+1) 字节数组   

    3、列表
       chunk中，指令表、常量表、子函数原型表等都是按照列表的形式进行存储的。格式为：先用cint记录列表
       长度，然后是存储n个列表元素。

    package binchunk

    type  binaryChunk struct{
        header      //头部
        sizeUpvalues byte    //主函数Upvlue数量
        mainFunc    *Prototype   //主函数原型
    }


6、header头部
    chunk的头部包括签名、版本号、格式编号、不同整数类型所占用的字节数、字节序以及浮点数的格式识别信息等。
    header的go定义如下：

    type header struct{
        signature   [4]byte
        version     byte
        format      byte
        luacData    [6]byte
        cintSize    byte
        sizetSize   byte
        instructionSize byte
        luaIntegerSize  byte
        luaNumberSize   byte
        luacInt int64
        luacNum float64
    }
    1、签名signature
        lua的chunk的签名是四个字节：esc、L、u、a的ascii码，即0x 1B 4C 75 61。
        签名主要用来检测是否是合法的chunk，如果不是就拒绝加载。
    2、版本号version
        记录chunk文件生成时所使用的lua版本号。由三部分组成：Major version、Minior 
        version、Release version。5.3.5表示Major version=5，Minior version=3，
        Release version=5.
        chunk中的版本号是根据Major version、Minior version得到的。计算方法为
        Major version * 16 + Minior version。（5 * 16 + 3 = 80 = 0x53）Release version主要用于bug的修复，并不会对
        chunk的格式进行调整。虚拟机在加载chunk的时候，会检测版本号version，如果和虚拟机的版本不匹配，
        就拒绝加载这个chunk文件。
    3、格式format
        记录了chunk的格式编号，虚拟机在加载chunk的时候也会检查该格式编号，如果和虚拟机的版本不匹配，
        就拒绝加载这个chunk。Lua默认使用的format为0x00.
    4、luacData
        之后的6个字节称之为LUAC_DATA,前两个字节0x1993，就是在1993年发布的1.0版本，后续的四个字节
        依次是回车符0x0d、换行符0x0a、替换符0x1a、换行符0x0a。1993 0d0a 1a0a
        这六个字节同样起到校验的作用，匹配就加载，否则拒绝。 
    5、整数以及虚拟机指令的宽度
       后续的五个字节依次记录了cint、size_t、虚拟机的指令、lua整数以及lua浮点数类型在chunk里面占用
       的字节数。0408 0408 08
       虚拟机在加载chunk的时候会检测这几个字节，如果不匹配就会拒绝加载。
    6、luac_int
        在接下来的若干个字节存储lua的整数0x5678.因为本机lua整数占用8个字节，所以使用8个字节存储0x567878 5600 0000 0000 00。至于为什么要存储这个东西，当然也是有原因的，那就是用来检测二进制chunk
        的字节序 是大端字节序 还是小端字节序。虚拟机在记载chunk的时候，会使用这个数据来检测它的字节序
        是否和本机的字节序相匹配，一样就加载，否则拒绝。
    7、 luac_num
        最后的若干个字节存储lua浮点数370.5，在这里使用的是8个字节进行存储。（00 0000 0000 2877 40）
        同样，存储这个浮点数的目的也是为了检查二进制chunk使用的浮点数的格式。虚拟机在加载chunk文件的时
        候，会使用这个浮点数，检查它的格式是否和本机一致，一致就加载否则就拒绝。一般浮点数的格式为ieee 754浮点格式。


7、函数原型
    函数原型主要包括函数的基本信息、指令表、常量表、upvalue表、子函数原型表以及调试信息；
    基本信息又包括源文件名、起止行号、固定的参数个数、是否是vararg参数、以及运行函数所必须的寄存器数量
    调试信息包括行号表、局部变量表以及upvalue名称列表


    type Prototype struct {
        Source      string
        LineDefined     uint32
        LastLineDefined uint32
        NumParams       byte
        IsVararg        byte
        MaxStackSize    byte
        Code            []uint32
        Constants       []interface{}
        Upvalues        []Upvalue
        Protos          []*Prototype
        LineInfo        []uint32
        LocVars         []LocVar
        UpvalueNames    []string
    }

    1、源文件名Source
        用来记载chunk是什么源文件编译出来的。为了避免重名，只有在main函数里，这个字段才会有真正的值。
        在其他嵌套的函数原型中，该字段为空字符串。如果使用了-s选项进行编译的话，源文件名以及调试信息
        会被从chunk中略去。

         $hexdump -C luac.out 

        00000000  1b 4c 75 61 53 00 19 93  0d 0a 1a 0a 04 08 04 08  |.LuaS...........|
        00000010  08 78 56 00 00 00 00 00  00 00 00 00 00 00 28 77  |.xV...........(w|
        00000020  40 01 10 40 68 65 6c 6c  6f 77 6f 72 6c 64 2e 6c  |@..@helloworld.l|
        00000030  75 61 00 00 00 00 00 00  00 00 00 01 02 04 00 00  |ua..............|
        00000040  00 06 00 40 00 41 40 00  00 24 40 00 01 26 00 80  |...@.A@..$@..&..|
        00000050  00 02 00 00 00 04 06 70  72 69 6e 74 04 15 68 65  |.......print..he|
        00000060  6c 6c 6f 20 77 6f 72 6c  64 ef bc 81 ef bc 81 ef  |llo world.......|
        00000070  bc 81 01 00 00 00 01 00  00 00 00 00 04 00 00 00  |................|
        00000080  06 00 00 00 06 00 00 00  06 00 00 00 06 00 00 00  |................|
        00000090  00 00 00 00 01 00 00 00  05 5f 45 4e 56           |........._ENV|


        10 40 68 65 6c 6c  6f 77 6f 72 6c 64 2e 6c 75 61

        因为字符串长度小于253，因此使用的是短字符串形式进行存储。字符串长度+1 占用一个字节（0x10），
        也就是十进制的16，在-1，字符串的长度为15。长度之后是@helloworld.lua，正好占用15个字节。
        @表示chunk是从源文件编译得到的，如果是以=开头，则表示其他的意义。比如“=stdin”，说明这个chunk
        是从标准输入编译而来，如果没有=，则表示chunk是从程序提供的字符串编译而来。去掉这个@符号后，
        才是实际的源文件名称。

    2、 起止行号LineDefined LastLineDefined
        源文件名后的两个cin整数，就是原型对应的函数在源文件中的起止行号。如果是普通的函数，起止行号应该
        都>0, 如果是main函数，则它们都为0.

    3、固定参数个数 NumParams
        起止行号之后的一个字节纪录了函数的固定参数的数量。这里所说的固定参数，是相对于变长参数Vararg而言
        的。编译器生成的main函数没有固定参数，因为该位置为0.

    4、是否是Vararg参数
        接下来的一个字节用来记录函数是否为Vararg参数，也就是说是否是变长参数，1表示是，0表示不是。
        main函数是Vararg函数，因此其值为1.

    5、寄存器数量
        下一个字节表示的是寄存器数量。Lua编译器会为每个函数生成一个指令表，也就是我们说的字节码。由于Lua
        虚拟机是基于寄存器的，因此大部分的指令也都会涉及到虚拟寄存器的操作。Lua编译器会在编译函数的时候
        ，将这个数量计算出来，并按照字节类型保存在函数的原型中。

        这个字段也称之为MaxStackSize。因为在Lua虚拟机运行函数的时候，实际使用的是一种栈结构，这种结构，除了可以进行常规的push 、pop操作之外，还可以按照索引进行访问，可以用来模拟cpu的寄存器。
    
    6、指令表
        之后就是指令表。在这个hello world中，main函数有 4 条指令，每个指令占用 4 个字节。

    7、常量表
        之后就是常量表，用于存储代码中出现的字面量，包括nil、布尔值、整数、浮点数以及字符串。
        每个常量都以1 个字节的tag 开始，用于标示后续存放的是那种类型的常量。

        tag     lua字面量类型    存储类型    
        0x00    nil             不存储
        0x01    boolean         字节（0、1）
        0x03    number          Lua浮点数
        0x13    integer         Lua整数
        0x04    string          短字符串
        0x14    string          长字符串

        定义tag常量

        const (
            TAG_NIL     = 0x00
            TAG_BOOLEAN = 0x01
            TAG_NUMBER  = 0x03
            TAG_INTEGER = 0x13
            TAG_SHORT_STR   = 0x04
            TAG_LONG_STR    = 0x14
        )

    8、Upvalue表
        之后就是Upvalue表。每个元素占用两个字节。

        定义Upvalue结构体：

        type Upvalue struct {
            Instack byte
            Idx     byte
        }

        hello world 中有一个Upvalue：01 00 00 00 01 00

    9、子函数原型列表
        之后就是子函数原型列表。因为hello world中没有定义函数，因此该列表长度为0，占用四个4字节。
    10、行号表
        之后是行号表，cin类型。行号表中的行号和指令表中的指令是一一对应的关系，分别记录每条指令在
        源码中对应的行号。hello代码只有4条指令，对应的行号都是1
    11、局部变量表
        之后就是局部变量表，用来记录局部变量名。表中的每个项都会包含变量名--使用字符串类型进行存储--
        以及起止指令的索引--使用cint类型进行存储。
        定义局部变量的结构体LocVar
        type LocVar struct {
            VarName     string
            StartPC     uint32
            EndPC       uint32
        }

        hello world中没有使用局部变量，所以局部变量表的长度为0，占用4个字节。
    12、Upvalue名称列表
        函数原型的最后一个部分就是Upvalue名称列表。该列表中的元素--使用字符串类型进行存储--）和
        之前的Upvalue表中的元素一一对应。分别记录每个Upvalue在源代码中的名称。hello world程序
        使用了一个Upvalue，名称为“_ENV”。
        00 00  05 5f 45 4e 56           |........._ENV|

        行号表，局部变量表以及Upvalue名称列表里面存储的都是调试信息，可以使用-s选项清空chunk中的
        对应字段。
    
8、Undump函数
    用来解析二进制chunk文件的信息：

    func Undump(data []byte) *Prototype {
        reader := &reader{byte}
        reader.checkHeader()    //检查头部
        reader.readByte()       //跳过Upvalue的数量
        return reader.readProto("") //读取函数原型信息
    }

9、解析chunk文件
    用来解析chunk的结构定义为：

    type reader struct {
        data []byte
    }

    结构reader只有一个data字段，存放将要被解析的二进制chunk数据。

    1、读取基本的数据类型信息
        读取基本的数据类型的方法一共有7种。其他的方法通过调用者几个方法来从chunk文件中提取信息。

        readByte方法：从字节流中读取一个字节

        func (self *reader) readByte() byte {
            b ：= self.data[0]
            self.data = self[1:]
            return b 
        }

        readUint32():使用smallend方式从字节流中读取一个c int类型的整数，在c中占用4个字节，
        对应go类型的uint32类型。

        func (self *reader) readUint32() uint32 {
            i := binary.LittleEndian.Uint32(self.data)
            self.data = self.data[4:]
            return i
        }

        readUint64():使用smallend方式从字节流中读取一个c size_t类型的整数，在c中占用8个字节，
        对应go类型的uint64类型。

        func (self *reader) readUint64() uint64 {
            i := binary.LittleEndian.Uint64(self.data)
            self.data = self.data[8:]
            return i
        }


        readLuaInteger() 利用readUint64()方法从字节流中读取一个Lua整数（占用8个字节，对应go的
        int64）.

        func (self *reader) readLuaInteger() uint64 {
            return int64(self.readUint64())
        }

        readLuaNumber()利用readUint64()方法从字节流中读取一个Lua浮点数（占用8个字节，对应go的
        float64）.

        func (self *reader) readLuaNumber() float64 {
            return math.Float64frombits(self.readUint64())
        }

        readString()从字节流中读取字符串(对应go的string类型)。

        func (self *reader) readString() string {
            size := uint(self.readByte)     //short or long string
            if 0 == size {  // null string
                return ""
            }

            if 0xFF == size {   //long string
                size = uint(self.readUint64())
            }

            bytes := self.readBytes(size - 1)
            return string(bytes)
        }

        readBytes()方法从字节流中读取n个字节：

        func (self *reader) readBytes(n uint) []byte {
            bytes := self.data[:n]
            self.data = self.data[n:]
            return bytes
        }

    2、检查header
        checkHeader()方法从字节流中读取并检测二进制chunk头的各个字段信息，如果
        发现某个字段和预期的不一致，就调用panic停止加载chunk。

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


    3、读取函数原型信息
        readProto()方法从字节流中读取函数原型信息。

        func (self *reader) readProto(parentSource string) *Prototype {
    source := self.readString()
            if "" == source {source = parentSource }

            return &Prototype {
                Source:         source,
                LineDefined:    self.readUint32(),
                LastLineDefined:self.readUint32(),
                NumParams:      self.readByte(),
                IsVararg:       self.readByte(),
                MaxStackSize:   self.readByte(),
                Code:           self.readCode(),
                Constants:      self.readConstants(),
                Upvalues:       self.readUpvalues(),
                Protos:         self.readProtos(source),
                LineInfo:       self.readLineInfo(),
                LocVars:        self.readLocVars(),
                UpvalueNames:   self.readUpvalueNames(),
            }

        }

        Lua编译器只给main函数设置了源文件名，用来减少数据的冗余，这样子函数原型就需要从父函数原型
        那里获取源文件名。

        readCode():从字节流中读取指令列表。

        func (self *reader) readCode() uint32 {
            code := make( []uint32, self.readUint32() )
            for i := range code {
                code[i] = self.readUint32()
            }
            return code
        }


        readConstants()：从字节流中读取常量表

        func (self *reader) readConstants() interface{} {
            constants := make([]interface{}, self.readUint32())
            for i := range constants {
                constants[i] = self.readConstant()
            }
            return constants
        }

        readConstant():从字节流总读取一个常量

        func (self *reader) readConstant() interface{} {
            switch self.readByte() {
                case TAG_NIL:       return nil
                case TAG_BOOLEEAN:  return self.readByte() != 0
                case TAG_INTEGET:   return self.readLuaInteger()
                case TAG_NUMBER:    return self.readLuaNumber()
                case TAG_SHORT_STR  return self.readString()
                case TAG_LONG_STR   return self.readString()
                default:            panic("corrupted")  
            }
        }

        readUpvalues():从字节流中读取Upvalue表信息。

        func (self *reader) readUpvalues() []Upvalue {
            upvalues := make([]Upvalue, self.readUint32())
            for i : = range upvalues {
                upvalues[i] = Upvalue{
                    Instack:    self.readByte(),
                    Idx:        self.readByte(),
                }
            }
            return upvalues
        }

        因为函数原型本身是递归的数据结构，因此readProto()也会递归调用，并读取子函数的原型。

        readProtos()：从字节流中读取子函数的原型列表。

        func (self *reader) readProtos(parentSource string) []*Prototype {
            protos := make([]*Prototype, self.readUint32())
            for i := range protos {
                protos[i] = self.readProto(parentSource)
            }
            return protos
        }

        readLineInfo():从字节流中读取行号表。

        func (self *reader) readLineInfo() []uint32 {
            lineInfo := make([]uint32, self.readUint32())
            for i := range lineInfo {
                lineInfo[i] = self.readUint32()
            }
            return lineInfo
        } 

        readLocVars():从字节流中读取局部变量表。

        func (self *reader) readLocVars() []LocVar {
            locVars := make([]LocVar, self.readUint32())
            for i := range locVars {
                locVars[i] = LocVar {
                    VarName:    self.readString(),
                    StartPC:    self.readUint32(),
                    EndPC:      self.readUint32(),
                }
            }
            return locVars
        }

        readUpvalueNames():从字节流中读取Upvalue名称列表。

        func (self *reader) readUpvalueNames() []string {
            names := make([]string, slef.readUint32())
            for i : = range names {
                names[i] = self.readString()
            }
            return names
        }




10、单元测试

    验证二进制chunk文件的解析效果
    lunago xxx.out（lua的二进制chunk文件）

    func main() {
        if len(os.Args) > 1 {
            data,err := ioutil.ReadFile(os.Args[1])
            if err != nil { panic("err") }
            proto := binchunk.Undump(data)
            list(proto)
        }
        println("Hello world!!!")
    }

    通过命令行参数把需要反编译的二进制chunk传递给main函数，接下来就会读取文件数据，并调用
    Undump函数解析函数原型信息，之后通过list把函数原型的信息打印出来。

    func list (f *binchunk.Prototype) {
        printHeader(f)
        printCode(f)
        printDetail(f)

        for _, p := range f.Protos {
            list(p)
        }

    }

    list首先打印函数的基本信息，之后是指令表和其他信息，然后递归调用，把子函数的信息也打印出来。

    func printHeader(f *binchunk.Prototype) {
        funcType := "main"
        if f.LineDefined > 0 { funcType = "function" }

        varargFlag := ""
        if f.IsVararg > 0 { varargFlag = "+" }

        fmt.Printf("\n%s <%s : %d, %d> (%d instructions)\n", funcType,
                    f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

        fmt.Printf("%d%s params, %d slots, %d upvalues, ", f.NumParams, varargFlag, 
                    f.MaxStackSize, len(f.Upvalues))

        fmt.Printf("%d locals, %d constants, %d functions\n", len(f.LocVars), 
                    len(f.Constants), len(f.Protos))
    }

    下面的printCode仅仅打印指令的序号 、行号 以及16进制表示。

    func printCode(f *binchunk.Prototype) {
        for pc, c := range f.Code {
            line := "-"
            if len(f.LineInfo) > 0 {
                line = fmt.Sprintf("%d", f.LineInfo[pc])
            }
            fmt.Printf("\t%d\t[%s]\t0x%08x\n", pc + 1, line, c)
        }
    }

    printDetail函数打印常量表、局部变量表以及Upvalue表。

    func printDetail(f *binchunk.Prototype) {
        fmt.Printf("constants (%d):\n", len(f.Constants))
        for i, k := range f.Constants {
            fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))      
        }

        fmt.Printf("locals (%d):\n", len(f.LocVars))
        for i,locVar := range f.LocVars {
            fmt.Print("\t%d\t%s\t%d\t%d\n", i, locVar.VarName,locVar.StartPC+1, 
                locVar.EndPC+1)
        }

        fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
        for i, upval := range f.Upvalues {
            fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalName(f, i), upval.Instack, upval.Idx)
        }

    }

    constantToString函数把常量表里面的常量转换为字符串。

    func constantToString(k interface{}) string {
        switch k.(type) {
        case nil:       return "nil"
        case bool:      return fmt.Sprintf("%t", k)
        case float64:   return fmt.Sprintf("%g", k)
        case int64:     return fmt.Sprintf("%d", k)
        case string:    return fmt.Sprintf("%q", k)
        default:    return "?"
        }

    }

    upvalName函数根据Upvalue的索引从调试信息里面查找Upvalue的名称

    func upvalName(f *binchunk.Prototype, idx int) string {
        if len(f.UpvalueNames) > 0 {
            return f.UpvalueNames[idx]
        }
        return "-"
    }

    最终的执行结果如下：

    $ ./lunago luac.out 

    main <@helloworld.lua : 0, 0> (4 instructions)
    0+ params, 2 slots, 1 upvalues, 0 locals, 2 constants, 0 functions
        1   [6] 0x00400006
        2   [6] 0x00004041
        3   [6] 0x01004024
        4   [6] 0x00800026
    constants (2):
        1   "print"
        2   "hello world！！！"
    locals (0):
    upvalues (1):
        0   _ENV    1   0
    Hello world!!!

---------------------------------

=======================================

lua虚拟机指令以及编码格式 指令集

虚拟机大致分为两类：stack based 以及register based。
a、基于stack的虚拟机使用push 、pop 入栈出栈，其他的指令则是对stack top进行操作，因此指令集相对比较大，
  不过指令的平均长度比较短；
b、基于register的虚拟机，由于可以直接对寄存器进行寻址，因此不需要push 、pop指令，指令集相对较小。不过
    由于需要把寄存器的地址编码到指令里面，所以指令的长度会比较长。

按照指令长度，指令集可以分为Fixed-width指令集 和 variable-width指令集。lua虚拟机使用的是固定长度
的指令集，每条指令占用4个字节--32bit，其中6bit用于opcode，其余26bit用于操作数operand。

lua的指令，根据其作用，大致可以分为：常量加载指令、运算符相关指令、循环跳转指令、函数调用相关指令、
    表操作指令以及Upvalue操作指令。

1、指令编码格式
    编码格式

    每条lua指令占用4个字节，32bit，可以使用go的uint32进行表示。低6位用于操作码，高26位用于操作数。
    根据高26bit的分配以及解释方式的不同，lua虚拟机指令可以分为4种，对应4种编码模式：iABC、iABx、
    iAsBx、iAx。

    iABC格式的指令可以额携带a b c 三个操作数，分别占用8 9 9个bit
    iABx格式的指令可以额携带A 和 Bx 两个操作数，分别占用8 16个bit
    iAsBx格式的指令可以额携带As 和Bx 两个操作数，分别占用8 18个bit
    iAx格式的指令只可以额携带一个操作数，占用26个bit
    4种模式种，只有iAsBx模式下的sBx操作数会被解释为有符号整数，其他情况下，操作数都被解释为无符号整数。


    定义表示指令编码格式的常量，位于opcodes.go文件中。


    操作吗

    操作码用于识别指令。因为lua虚拟机使用6bit表示操作码，因此最多可以表示64条指令。lua 5.3定义了47
    条指令，操作码从0开始，46截止。

    操作数
    指令的参数就是操作数，每个指令可以携带1--3个操作数。其中参数A主要用来表示目标寄存器的索引，其他
    操作数按照其表示的信息，可以大致分为4种：OpArgN、OpArgU、 OpArgR、 OpArgK。

    OpArgN类型的操作数不表示任何信息，也就是说不会被使用。比如move指令只使用A、B操作数，不使用C操作数。
    move指令的格式为iABC

    OpArgU，操作数也可能表示布尔值、整数值、upvalue索引、子函数的索引等等。 

    OpArgR类型的操作数载在iABC模式下表示寄存器的索引。在iAsBx模式下表示跳转偏移量。在move指令中，
    A操作数表示目标寄存器的索引，B操作数（OpArgR类型）表示源寄存器的索引。move指令可以表示伪指令：
    RA：= RB。


    OpArgK类型的操作数表示常量表的索引或者寄存器的索引，分为两种情况：
    1、loadk指令（iABx模式，用于将常量表中的常量加载到寄存器中），该指令的Bx操作数表示常量表的索引，
        伪代码可以表示为RA ：= （Bx）。
    2、部分iABC模式的指令，它们的B或者C操作数可以表示常量表的索引，也可以表示寄存器的索引。例如add，
        RA ：= （B) + (C).对于这种既可以表示常量表的索引 ，又可以表示寄存器的索引的情况，怎么确定
        索引的类型呢？在iABC模式中，BC操作数各占9bit，如果BC操作数属于OpArgK类型，那么就只能使用
        9个bit中的低8bit，最高位的bit如果是1，表示操作数是常量表索引，如果是0表示是寄存器的索引。

    定义表示操作数类型的常量：

    /* OpArgMask */
    const (
        OpArgN = iota // argument is not used
        OpArgU        // argument is used
        OpArgR        // argument is a register or a jump offset
        OpArgK        // argument is a constant or register/constant
    )

    指令表

    lua实现把每条指令的基本信息都编码为一个字节，比如编码模式、是否设置寄存器A，操作数BC使用的类型。。
    对其模仿的时候，我们使用结构体，而不是字节，并把操作码的名称保存起来。

    定义opcode的结构体：

    type opcode struct {
        testFlag byte // operator is a test (next instruction must be a jump)
        setAFlag byte // instruction set register A
        argBMode byte // B arg mode
        argCMode byte // C arg mode
        opMode   byte // op mode
        name     string
    }

    定义指令表

    var opcodes = []opcode{
        /*     T  A    B       C     mode         name    */
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "MOVE    "}, // R(A) := R(B)
        opcode{0, 1, OpArgK, OpArgN, IABx /* */, "LOADK   "}, // R(A) := Kst(Bx)
        opcode{0, 1, OpArgN, OpArgN, IABx /* */, "LOADKX  "}, // R(A) := Kst(extra arg)
        opcode{0, 1, OpArgU, OpArgU, IABC /* */, "LOADBOOL"}, // R(A) := (bool)B; if (C) pc++
        opcode{0, 1, OpArgU, OpArgN, IABC /* */, "LOADNIL "}, // R(A), R(A+1), ..., R(A+B) := nil
        opcode{0, 1, OpArgU, OpArgN, IABC /* */, "GETUPVAL"}, // R(A) := UpValue[B]
        opcode{0, 1, OpArgU, OpArgK, IABC /* */, "GETTABUP"}, // R(A) := UpValue[B][RK(C)]
        opcode{0, 1, OpArgR, OpArgK, IABC /* */, "GETTABLE"}, // R(A) := R(B)[RK(C)]
        opcode{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABUP"}, // UpValue[A][RK(B)] := RK(C)
        opcode{0, 0, OpArgU, OpArgN, IABC /* */, "SETUPVAL"}, // UpValue[B] := R(A)
        opcode{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABLE"}, // R(A)[RK(B)] := RK(C)
        opcode{0, 1, OpArgU, OpArgU, IABC /* */, "NEWTABLE"}, // R(A) := {} (size = B,C)
        opcode{0, 1, OpArgR, OpArgK, IABC /* */, "SELF    "}, // R(A+1) := R(B); R(A) := R(B)[RK(C)]
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "ADD     "}, // R(A) := RK(B) + RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SUB     "}, // R(A) := RK(B) - RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "MUL     "}, // R(A) := RK(B) * RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "MOD     "}, // R(A) := RK(B) % RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "POW     "}, // R(A) := RK(B) ^ RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "DIV     "}, // R(A) := RK(B) / RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "IDIV    "}, // R(A) := RK(B) // RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BAND    "}, // R(A) := RK(B) & RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BOR     "}, // R(A) := RK(B) | RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BXOR    "}, // R(A) := RK(B) ~ RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SHL     "}, // R(A) := RK(B) << RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SHR     "}, // R(A) := RK(B) >> RK(C)
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "UNM     "}, // R(A) := -R(B)
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "BNOT    "}, // R(A) := ~R(B)
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "NOT     "}, // R(A) := not R(B)
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "LEN     "}, // R(A) := length of R(B)
        opcode{0, 1, OpArgR, OpArgR, IABC /* */, "CONCAT  "}, // R(A) := R(B).. ... ..R(C)
        opcode{0, 0, OpArgR, OpArgN, IAsBx /**/, "JMP     "}, // pc+=sBx; if (A) close all upvalues >= R(A - 1)
        opcode{1, 0, OpArgK, OpArgK, IABC /* */, "EQ      "}, // if ((RK(B) == RK(C)) ~= A) then pc++
        opcode{1, 0, OpArgK, OpArgK, IABC /* */, "LT      "}, // if ((RK(B) <  RK(C)) ~= A) then pc++
        opcode{1, 0, OpArgK, OpArgK, IABC /* */, "LE      "}, // if ((RK(B) <= RK(C)) ~= A) then pc++
        opcode{1, 0, OpArgN, OpArgU, IABC /* */, "TEST    "}, // if not (R(A) <=> C) then pc++
        opcode{1, 1, OpArgR, OpArgU, IABC /* */, "TESTSET "}, // if (R(B) <=> C) then R(A) := R(B) else pc++
        opcode{0, 1, OpArgU, OpArgU, IABC /* */, "CALL    "}, // R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
        opcode{0, 1, OpArgU, OpArgU, IABC /* */, "TAILCALL"}, // return R(A)(R(A+1), ... ,R(A+B-1))
        opcode{0, 0, OpArgU, OpArgN, IABC /* */, "RETURN  "}, // return R(A), ... ,R(A+B-2)
        opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORLOOP "}, // R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }
        opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORPREP "}, // R(A)-=R(A+2); pc+=sBx
        opcode{0, 0, OpArgN, OpArgU, IABC /* */, "TFORCALL"}, // R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
        opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "TFORLOOP"}, // if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx }
        opcode{0, 0, OpArgU, OpArgU, IABC /* */, "SETLIST "}, // R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
        opcode{0, 1, OpArgU, OpArgN, IABx /* */, "CLOSURE "}, // R(A) := closure(KPROTO[Bx])
        opcode{0, 1, OpArgU, OpArgN, IABC /* */, "VARARG  "}, // R(A), R(A+1), ..., R(A+B-2) = vararg
        opcode{0, 0, OpArgU, OpArgU, IAx /*  */, "EXTRAARG"}, // extra (larger) argument for previous opcode
    }

    指令的解码

    之前使用的是uint32类型来表示存在二进制chunk文件里面的指令，为了便于操作，给指令定义一个专门的
    类型。创建instruction。go，并定义Instruction类型来表示指令的类型。

    type Instruction uint32


    然后给这个类型定义一些方法，用于解码指令。


    Opcode()用于从指令中提取操作码：

    func (self Instruction) Opcode() int {
        return int(self & 0x3F)
    }

    ABC()用于从iABC模式的指令中提取参数：

    func (self Instruction) ABC() (a, b, c int) {
        a = int(self >> 6 & 0xFF)
        c = int(self >> 14 & 0x1FF)
        b = int(self >> 23 & 0x1FF)
        return
    }

    ABx()用于从iABx模式的命令中提取参数：

    func (self Instruction) ABx() (a, bx int) {
        a = int(self >> 6 & 0xFF)
        bx = int(self >> 14)
        return
    }

    AsBx()用于从iAsBx模式的指令中提取参数：

    func (self Instruction) AsBx() (a, sbx int) {
        a, bx := self.ABx()
        return a, bx - MAXARG_sBx
    }

    Ax（）用于从iAx模式的指令中提取参数：

    func (self Instruction) Ax() int {
        return int(self >> 6)
    }

    说明：sBx操作数共18bit，表示的是有符号整数。有很多方式可以把有符号整数编码为bit序列，比如2的补码。
    Lua虚拟机采用了一种称之为偏移二进制码的方法---offset binary--也称之为Excess—K。具体来说，如果
    把sBx解释成无符号整数的时候，它的数值为x，那么解释成有符号整数的时候，它的数值就是x - K。K是sBx
    所能表示的最大无符号整数值的二分之一，也就是上面中的 MAXARG_sBx。

        min                                 max
    Bx   0              131071              262143
    sBx  -131071            0               131071

    Bx 以及 sBx 的取值范围

    OpName()返回指令的操作码名称

    func (self Instruction) OpName() string {
        return opcodes[self.Opcode()].name
    }

    OpMode()返回指令的编码模式

    func (self Instruction) OpMode() byte {
        return opcodes[self.Opcode()].opMode
    }
    BMode()返回操作数B的使用模式

    func (self Instruction) BMode() byte {
        return opcodes[self.Opcode()].argBMode
    }

    CMode()返回操作数C的使用模式

    func (self Instruction) CMode() byte {
        return opcodes[self.Opcode()].argCMode
    }


    单元测试

    完善之前的反编译工具，打印指令的操作码还有操作数。
    增加
    import . "lunago/vm"
    修改如下代码

    func printCode(f *binchunk.Prototype) {
        for pc, c := range f.Code {
            line := "-"
            if len(f.LineInfo) > 0 {
                line = fmt.Sprintf("%d", f.LineInfo[pc])
            }
            //将uint32类型转换为Instruction类型
            i := Instruction(c)

            //fmt.Printf("\t%d\t[%s]\t0x%08x\n", pc + 1, line, c)
            fmt.Printf("\t%d\t[%s]\t%s \t", pc + 1, line, i.OpName())
            printOperands(i)
            fmt.Printf("\n")
        }
    }

    把指令从uint32类型转换成自定义的Instruction类型，这样就可以很方便的拿到指令的操作数了。

    printOperands（）用于打印操作数。


    func printOperands(i Instruction) {
        switch i.OpMode() {
        case IABC:
            a, b, c := i.ABC()

            fmt.Printf("%d", a)
            if i.BMode() != OpArgN {
                if b > 0xFF {
                    fmt.Printf(" %d", -1-b&0xFF)
                } else {
                    fmt.Printf(" %d", b)
                }
            }
            if i.CMode() != OpArgN {
                if c > 0xFF {
                    fmt.Printf(" %d", -1-c&0xFF)
                } else {
                    fmt.Printf(" %d", c)
                }
            }
        case IABx:
            a, bx := i.ABx()

            fmt.Printf("%d", a)
            if i.BMode() == OpArgK {
                fmt.Printf(" %d", -1-bx)
            } else if i.BMode() == OpArgU {
                fmt.Printf(" %d", bx)
            }
        case IAsBx:
            a, sbx := i.AsBx()
            fmt.Printf("%d %d", a, sbx)
        case IAx:
            ax := i.Ax()
            fmt.Printf("%d", -1-ax)
        }
    }


    对于iABC模式的指令，首先打印操作数A，而操作数BC在某些指令里面可能并没有使用，所以不一定会打印出来。
    如果操作数BC的最高位是1，那么她就表示常量表索引，使用负数来输出。

    对于iABx模式的指令，也是先打印出操作数A，然后是操作数Bx。如果操作数Bx表示常量表索引，一样使用负数
    输出。

    对于iAsBx模式的指令，先打印操作数A，在打印操作数sBx；对于iAx模式的指令，只会打印操作数Ax就足够了。

    测试结果如下：

    $ ./lunago luac.out 

    main <@helloworld.lua : 0, 0> (4 instructions)
    0+ params, 2 slots, 1 upvalues, 0 locals, 2 constants, 0 functions
        1   [6] GETTABUP    0 0 -1
        2   [6] LOADK       1 -2
        3   [6] CALL        0 2 1
        4   [6] RETURN      0 1
    constants (2):
        1   "print"
        2   "hello world！！！"
    locals (0):
    upvalues (1):
        0   _ENV    1   0
    Hello world!!!


 ========================================
 
 Lua api

 为了方便地将lua脚本嵌入到其他Host环境中，Lua是以Library的形式实现的，其他的应用程序只需要链接
 lua库就可以额使用lua提供的api，轻松获得脚本的执行能力。


 lua api是一系列以“lua_”开始的c函数，其中也包括宏定义。用书使用lua_newstate()创建lua_state,
 其他的函数则用来操作这个lua_state.

 lua栈

 lua state是lua api中非常核心的概念，所有的api函数都是围绕lua state来进行操作的。而lua state
 内部封装的最为基础的一个状态就是虚拟的栈。lua栈是宿主语言和lua语言进行沟通的桥梁，lua api函数有
 很大一部分是专门用来操作lua栈的。

 宿主语言-----> lua stack -----> lua
        <-----           <-----

 Lua的数据类型和值

 lua是动态类型语言，其中的百变量是不懈怠类型信息的，变量的值才会携带类型信息。也就是说，任何一个lua
 变量都可以被赋予任何类型的值。

 lua一共提供8种数据类型：nil、boolean、number、string、table、function、thread以及userdata。

 lua提供了type（）函数，用来获取变量的类型。

         lua数据类型 go数据类型
         nil        nil
         boolean    bool
         integer    int64
         float      float64
         string     string

 lua给没中lua数据类型都定义了一个常量值，我们需要使用宏定义转换为go语言的常量，以便以后使用。

 consts。go中定义如下常量：

 package api

    /* basic types */
    const (
        LUA_TNONE = iota - 1 // -1
        LUA_TNIL
        LUA_TBOOLEAN
        LUA_TLIGHTUSERDATA
        LUA_TNUMBER
        LUA_TSTRING
        LUA_TTABLE
        LUA_TFUNCTION
        LUA_TUSERDATA
        LUA_TTHREAD
    )       


    LUA_TNONE表示无效值。


    创建lua_state.go文件，并在里面定义用来表示lua值的luaValue类型。

    package state

    type luaState interface{}

    在这里我们仍然使用interface{}来表示各种不同类型的Lua值。

    添加typeOf（）函数，根据变量的值 返回其类型：

    func typeOf(val luaValue) LuaType {
        switch val.(type) {
        case nil:
            return LUA_TNIL
        case bool:
            return LUA_TBOOLEAN
        case int64, float64:
            return LUA_TNUMBER
        case string:
            return LUA_TSTRING
        default:
            panic("todo!")
        }
    }


    lua栈索引

    a、在大多数编程语言中，数组或者列表的索引都是从0开始的，然而在lua api里面，栈索引是从1开始的。
    b、lua中索引可以为事负值。正数索引称之为绝对索引，从栈底1开始递增，负数索引称之为相对索引。从
        栈顶-1开始递减。lua api会在内部把相对索引转换为绝对索引。
    c、如果lua栈的容量为n，栈顶索引是top（0 < top <= n），我们称位于[1，top]闭区间内的索引为有效
        Valid索引，位于位于[1,n]闭区间内的索引为可接受accettable的索引。如果要往栈里面写入值，必须
        给lua api提供 有效索引，否则可能会导致错误，甚至程序崩溃。如果仅仅是从栈里面读取值，则可以
        提供 可接受索引。对于无效的可接受索引，其行为差不多相当于该索引处存放的是nil值。


    定义luaStack 结构体

    在lua_stack。go中定义luaStack结构体。

    package state

    type luaStack struct {
        slots   []luaValue  //用来存放值
        top     int //记录栈顶的索引
    }

    定义函数newLuaStack（）,用来创建指定容量的栈。

    //创建指定容量的栈
    func newLuaStack(size int) *luaStack {
        return &luaStack{
            slots: make([]luaValue, size),
            top:   0,
        }
    }

    接下来为luaStack结构体定义一些方法。
    check()方法检查栈的空闲空间是否还可以容纳（push）至少n个值，如果不满足这个条件，就会调用go的
    append()函数对其进行扩容。

    func (self *luaStack) check(n int) {
        free := len(self.slots) - self.top
        for i := free; i < n; i++ {
            self.slots = append(self.slots, nil)
        }
    }

    push()方法用来将值push到栈顶，如果溢出，就先暂时调用panic（）终止程序的执行。

    func (self *luaStack) push(val luaValue) {
        if self.top == len(self.slots) {
            panic("stack overflow!")
        }
        self.slots[self.top] = val
        self.top++
    }


    pop()方法从栈顶弹出一个值，如果栈是空的，则调用panic()终止程序。

    func (self *luaStack) pop() luaValue {
        if self.top < 1 {
            panic("stack underflow!")
        }
        self.top--
        val := self.slots[self.top]
        self.slots[self.top] = nil
        return val
    }

    absIndex()方法吧索引切换成绝对索引--并没有考虑索引是否有效

    func (self *luaStack) absIndex(idx int) int {
        if idx >= 0 {
            return idx
        }
        return idx + self.top + 1
    }

    isValid()判断所有是否有效

    func (self *luaStack) isValid(idx int) bool {
        absIdx := self.absIndex(idx)
        return absIdx > 0 && absIdx <= self.top
    }

    get()根据索引从栈里面取值，如果索引无效 返回nil

    func (self *luaStack) get(idx int) luaValue {
        absIdx := self.absIndex(idx)
        if absIdx > 0 && absIdx <= self.top {
            return self.slots[absIdx-1]
        }
        return nil
    }

    set()根据索引向栈里面写入值，如果索引无效，调用panic（）终止

    func (self *luaStack) set(idx int, val luaValue) {
        absIdx := self.absIndex(idx)
        if absIdx > 0 && absIdx <= self.top {
            self.slots[absIdx-1] = val
            return
        }
        panic("invalid index!")
    }


    Lua State

    Lua State状态机封装了lua解释器的状态信息。暂时先将lua state中当作只有一个lua栈，简化理解，
    以后在扩展。

    定义lua state借口

    因为lua是使用c语言进行编写的，因此lua api就是很多结构体lua_State 进行操作的函数或者宏定义。
    go支持接口设计，因此可以吧这些函数整合到一个接口中。本程序仅仅实现一些基本的栈操作函数：基础栈操作
    的函数，栈访问的函数 以及 压栈的函数。

    新建lua_state.go文件，定义LuaState接口。

    package api
    /**
     * [LuaType int 类型]
     * @type {[type]}
     */
    type LuaType = int

    /**
     * LuaState 接口定义
     * 
     */
    type LuaState interface {
        /* basic stack manipulation */
        GetTop() int
        AbsIndex(idx int) int
        CheckStack(n int) bool
        Pop(n int)
        Copy(fromIdx, toIdx int)
        PushValue(idx int)
        Replace(idx int)
        Insert(idx int)
        Remove(idx int)
        Rotate(idx, n int)
        SetTop(idx int)
        /* access functions (stack -> Go) */
        TypeName(tp LuaType) string
        Type(idx int) LuaType
        IsNone(idx int) bool
        IsNil(idx int) bool
        IsNoneOrNil(idx int) bool
        IsBoolean(idx int) bool
        IsInteger(idx int) bool
        IsNumber(idx int) bool
        IsString(idx int) bool
        IsTable(idx int) bool
        IsThread(idx int) bool
        IsFunction(idx int) bool
        ToBoolean(idx int) bool
        ToInteger(idx int) int64
        ToIntegerX(idx int) (int64, bool)
        ToNumber(idx int) float64
        ToNumberX(idx int) (float64, bool)
        ToString(idx int) string
        ToStringX(idx int) (string, bool)
        /* push functions (Go -> stack) */
        PushNil()
        PushBoolean(b bool)
        PushInteger(n int64)
        PushNumber(n float64)
        PushString(s string)
    }


    定义 luaState 结构体

    有了LuaState 接口，还需要定义一个结构体来实现这个接口。

    package state

    type luaState struct {
        stack *luaStack
    }

    go并不强制要求显式实现接口，只要结构体实现了接口的全部方法，它就隐式实现了该接口。
    增加New（）函数，用来创建luaState的实例。

    func New() *luaState {
        return &luaState{
            stack: newLuaStack(20),
        }
    }

    暂时先把lua栈的容量设置为20.


    操作栈的基本方法

    GetTop() 、Pop（）、CheckStack（）。。。。。。

    api_stack.go实现文件中的代码：

    package state

    /**
     * 返回栈顶索引
     */
    func (self *luaState) GetTop() int {
        return self.stack.top
    }


    /**
     * AbsIndex(idx int)把索引转化为绝对索引。
     */
    func (self *luaState) AbsIndex(idx int) int {
        return self.stack.absIndex(idx)
    }


     /**
     * lua栈的容量不会自动增长，使用者需要检查栈的剩余空间，看看是否可以push n 个值而不会溢出。
     * 如果剩余空间足够 或者扩容成功 返回true，否则返回false.
     * n 表示需要多少个剩余空间存放数据。
     */
    func (self *luaState) CheckStack(n int) bool {
        self.stack.check(n)
        return true // ??? never fails
    }   

    扩容的逻辑已经在luaStack结构体中实现，这里只是简单的调用相关的方法。暂时忽略扩展失败的情况。


    /**
     * [ Pop(n int) 方法从栈顶弹出n 个值。]
     */
    func (self *luaState) Pop(n int) {
        for i := 0; i < n; i++ {
            self.stack.pop()
        }
    }

    Copy()方法把值从一个位置复制到另一个位置。

    func (self *luaState) Copy(fromIdx, toIdx int) {
        val := self.stack.get(fromIdx)
        self.stack.set(toIdx, val)
    }

    PushValue()方法把指定索引处的值push到栈顶。

    func (self *luaState) PushValue(idx int) {
        val := self.stack.get(idx)
        self.stack.push(val)
    }

    Replace()是PushValue()的反操作：
    将栈顶的值弹出，然后写入到指定的位置。

    func (self *luaState) Replace(idx int) {
        val := self.stack.pop()
        self.stack.set(idx, val)
    }

    Insert()方法将栈顶的值弹出，然后将其值插入到指定的位置。
    原来idx以及之后的值则分别向上移动一个位置。
    func (self *luaState) Insert(idx int) {
        self.Rotate(idx, 1)
    }

    Rotate()旋转操作。Insert就是旋转操作的一种。
    旋转操作。Rotate(idx, n int) 将[idx, top] 索引区间内的值朝着栈顶方向旋转 n 个位置。
    如果n是负数，那么实际的效果就是朝着栈底方向旋转。

    func (self *luaState) Rotate(idx, n int) {
        t := self.stack.top - 1           /* end of stack segment being rotated */
        p := self.stack.absIndex(idx) - 1 /* start of segment */
        var m int                         /* end of prefix */
        if n >= 0 {
            m = t - n
        } else {
            m = p - n - 1
        }
        self.stack.reverse(p, m)   /* reverse the prefix with length 'n' */
        self.stack.reverse(m+1, t) /* reverse the suffix */
        self.stack.reverse(p, t)   /* reverse the entire segment */
    }


    /**
     * [ Remove() 删除置顶索引处的值，然后将该值上面的所有值全部向下移动一个位置。]
    */
    func (self *luaState) Remove(idx int) {
        self.Rotate(idx, -1)
        self.Pop(1)
    }


    luaStack结构体的reverse()方法就是循环交换两个位置的值。

    在lua_stack.go中添加如下代码：

    func (self *luaStack) reverse(from, to int) {
        slots := self.slots
        for from < to {
            slots[from], slots[to] = slots[to], slots[from]
            from++
            to--
        }
    }

    SetTop()将栈顶索引设置为指定的值。如果指定的值小于当前栈顶的索引，效果则相当于弹出操作，指定值0
    相当于清空栈。
    如果指定的值 n 大于当前栈顶的索引，则效果相当于push （n - 栈顶索引） 个nil值。
    SetTop()根据不同的情况执行push 、pop操作。

    func (self *luaState) SetTop(idx int) {
        newTop := self.stack.absIndex(idx)
        if newTop < 0 {
            panic("stack underflow!")
        }

        n := self.stack.top - newTop
        if n > 0 {
            for i := 0; i < n; i++ {
                self.stack.pop()
            }
        } else if n < 0 {
            for i := 0; i > n; i-- {
                self.stack.push(nil)
            }
        }
    }

    其实前面的Pop()方法之时SetTop()方法的特殊情况，完全可以使用SetTop()方法实现。

    func (self *luaState) Pop(n int) {
        self.SetPop(-n - 1)
    }


    Push方法

    将Lua值从外部push到栈顶。先将5中基本类型的值push到栈顶，以后在添加其他的方法。

    package state
    
    func (self *luaState) PushNil() {
        self.stack.push(nil)
    }

    func (self *luaState) PushBoolean(b bool) {
        self.stack.push(b)
    }

    func (self *luaState) PushInteger(n int64) {
        self.stack.push(n)
    }

    func (self *luaState) PushNumber(n float64) {
        self.stack.push(n)
    }

    func (self *luaState) PushString(s string) {
        self.stack.push(s)
    }


    Access方法
    用于从栈中获取数据信息。差不多仅仅使用索引访问栈中存储的信息，不会改变栈的状态。
    api_access.go.

    package state

    import "fmt"
    import . "luago/api"

    TypeName()方法不需要读取任何栈数据，只是把给定的lua类型转换为对应的字符串表示。

      func (self *luaState) TypeName(tp LuaType) string {
        switch tp {
        case LUA_TNONE:
            return "no value"
        case LUA_TNIL:
            return "nil"
        case LUA_TBOOLEAN:
            return "boolean"
        case LUA_TNUMBER:
            return "number"
        case LUA_TSTRING:
            return "string"
        case LUA_TTABLE:
            return "table"
        case LUA_TFUNCTION:
            return "function"
        case LUA_TTHREAD:
            return "thread"
        default:
            return "userdata"
        }
    }  


    Type()根据索引返回值的类型，如果索引无效，则返回LUA_TNONE.

    func (self *luaState) Type(idx int) LuaType {
        if self.stack.isValid(idx) {
            val := self.stack.get(idx)
            return typeOf(val)
        }
        return LUA_TNONE
    }

    IsType()

    IsNone() IsNil() IsNoneOrNil() IsBoolean() 用来判断给定的索引处的值是否属于特定的类型，可以使用Type()来实现。

    func (self *luaState) IsNone(idx int) bool {
    return self.Type(idx) == LUA_TNONE
    }

    func (self *luaState) IsNil(idx int) bool {
        return self.Type(idx) == LUA_TNIL
    }


    func (self *luaState) IsNoneOrNil(idx int) bool {
        return self.Type(idx) <= LUA_TNIL
    }


    func (self *luaState) IsBoolean(idx int) bool {
        return self.Type(idx) == LUA_TBOOLEAN
    }

    IsString()判断指定索引处的值是不是字符串或者数字。

     func (self *luaState) IsString(idx int) bool {
        t := self.Type(idx)
        return t == LUA_TSTRING || t == LUA_TNUMBER
    }   

    IsNumber()方法判断给定随你处的值是不是数字类型，如果可以转化为数字类型也可以。

    func (self *luaState) IsNumber(idx int) bool {
        _, ok := self.ToNumberX(idx)
        return ok
    }

    IsInteger()判断指定索引处的值是不是整数类型。

    func (self *luaState) IsInteger(idx int) bool {
        val := self.stack.get(idx)
        _, ok := val.(int64)
        return ok
    }

    ToBoolean()从指定的索引处取出一个boolean值，如果值不是布尔类型，则需要进行类型转换。

    func (self *luaState) ToBoolean(idx int) bool {
        val := self.stack.get(idx)
        return convertToBoolean(val)
    }

    在Lua中，只有nil、false表示假，其他都表示真。lua_value.go定义convertToBoolean。

    func convertToBoolean(val luaValue) bool {
        switch x := val.(type) {
        case nil:
            return false
        case bool:
            return x
        default:
            return true
        }
    }

    ToNumber() ToNumberX()

    从指定的索引处取出一个数字，如果值不是数字类型，则需要进行类型转换。

    ToNumber()：如果值不是数字类型，并且也没有办法转换成数字类型，返回0.

    ToNumberX()：如果值不是数字类型，并且也没有办法转换成数字类型，则会报告转换是否成功。

    ToNumber()可以借助ToNumberX()实现。

    func (self *luaState) ToNumber(idx int) float64 {
        n, _ := self.ToNumberX(idx)
        return n
    }


    func (self *luaState) ToNumberX(idx int) (float64, bool) {
        val := self.stack.get(idx)
        switch x := val.(type) {
        case float64:
            return x, true
        case int64:
            return float64(x), true
        default:
            return 0, false
        }
    }

    在c语言中，可以通过返回值、指针类型的参数实现多个返回值的效果，go本身就支持多个返回值，所以不需要借助指针操作。

    ToInteger() ToIntegerX()

    从指定的索引处取出一个整数值，如果值不是整数类型，则需要进行类型转换。

    ToInteger()：如果值不是整数类型，并且也没有办法转换成整数类型，返回0.

    ToIntegerX()：如果值不是整数类型，并且也没有办法转换成整数类型，则会报告转换是否成功。

    func (self *luaState) ToInteger(idx int) int64 {
        i, _ := self.ToIntegerX(idx)
        return i
    }


    func (self *luaState) ToIntegerX(idx int) (int64, bool) {
        val := self.stack.get(idx)
        i, ok := val.(int64)
        return i, ok
    }

    ToString() ToStringX()

    从指定的索引处取出一个值，如果值是字符串，则返回字符串。如果值是数字，则将其转换为字符串--会修改栈，然后返回字符串。否则，返回空字符串。
    在c api中，该函数值有一个返回值，如果返回NULL，则表示指定的索引处的值不是字符串或者数字，由于go
    语言字符串类型没有对应的nil值，因此采用ToInteger() ToIntegerX()类似的做法，添加一个
    ToStringX()方法，其中返回值的第二个返回类型是布尔类型，表示转换是否成功。
    ToString()调用ToStringX(),忽略第二个返回值就可以了。

    func (self *luaState) ToString(idx int) string {
        s, _ := self.ToStringX(idx)
        return s
    }

    func (self *luaState) ToStringX(idx int) (string, bool) {
        val := self.stack.get(idx)

        switch x := val.(type) {
        case string:
            return x, true
        case int64, float64:
            s := fmt.Sprintf("%v", x) // todo 这里会修改stack
            self.stack.set(idx, s)
            return s, true
        default:
            return "", false
        }
    }

    //判断是否是table
    func (self *luaState) IsTable(idx int) bool {
        return self.Type(idx) == LUA_TTABLE
    }

    //判断是否是funciton
    func (self *luaState) IsFunction(idx int) bool {
        return self.Type(idx) == LUA_TFUNCTION
    }

    //判断是否是thread
    func (self *luaState) IsThread(idx int) bool {
        return self.Type(idx) == LUA_TTHREAD
    }
    单元测试main.go

    package main

    import "fmt"
    import . "lunago/api"
    import "lunago/state"
    import _ "lunago/binchunk"


    func main() {
        ls := state.New()
        ls.PushBoolean(true)
        printStack(ls)
        ls.PushInteger(10)
        printStack(ls)
        ls.PushNil()
        printStack(ls)
        ls.PushString("hello")
        printStack(ls)
        ls.PushValue(-4)
        printStack(ls)
        ls.Replace(3)
        printStack(ls)
        ls.SetTop(6)
        printStack(ls)
        ls.Remove(-3)
        printStack(ls)
        ls.SetTop(-5)
        printStack(ls)

    }
    func printStack(ls LuaState) {
        top := ls.GetTop()
        for i := 1; i <= top; i++ {
            t := ls.Type(i)
            switch t {
            case LUA_TBOOLEAN:
                fmt.Printf("[%t]", ls.ToBoolean(i))
            case LUA_TNUMBER:
                fmt.Printf("[%g]", ls.ToNumber(i))
            case LUA_TSTRING:
                fmt.Printf("[%q]", ls.ToString(i))
            default: // other values
                fmt.Printf("[%s]", ls.TypeName(t))
            }
        }
        fmt.Println()
    }

    输出结果如下所示：

    $ go run main.go 
    [true]
    [true][10]
    [true][10][nil]
    [true][10][nil]["hello"]
    [true][10][nil]["hello"][true]
    [true][10][true]["hello"]
    [true][10][true]["hello"][nil][nil]
    [true][10][true][nil][nil]
    [true]

============================================

    Lua 运算符

    lua有算数运算符、按位运算符、比较运算符、逻辑运算符、长度运算符、字符串拼接运算符。

    算数运算符

    +加法 -（减法、一元取反）、*（乘法）、/（除法）、//（整除）、%（取模）、^(乘方)。

    / ^ 会先把操作数转换为浮点数，然后进行计算，计算结果也是浮点数。其他6个算数运算符会先判断操作时是不是整数，如果是，则进行整数计算，结果也是整数，即使可能会溢出；如果不是整数，就把操作数转换为浮点数，然后进行计算，结果为浮点数。

    // 会将除法的结果向下取整（-负值方向）。

    在 java 和 go 语言中，整除 运算仅仅是将除法结果截断，向0 取整，这是和lua的区别。

    % ，可以使用整除 运算定义。

    a % b = a - （（a // b）* b）

    ^ 和 自负符传拼接运算符具有向右结合性，其他的二元运算符则具有向左结合性。

    + - * / 取反运算符可以直接映射为go的对应运算符，^ // % 运算符却不能直接映射。go没有^y 运算符，go的整除运算符仅仅用于整数，并且是直接阶段结果而非向下取整，%也是类似。

    新建math.go,并定义整除 和 取模 函数。

    func IFloorDiv(a, b int64) int64 {
        if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
            return a / b
        } else {
            return a/b - 1
        }
    }

    func FFloorDiv(a, b float64) float64 {
        return math.Floor(a / b)
    }

    取模函数可以使用整除函数实现

    /**
     * a % b == a - ((a // b) * b)
     */
    func IMod(a, b int64) int64 {
        return a - IFloorDiv(a, b)*b
    }

    /**
     * a % b == a - ((a // b) * b)
     */
    func FMod(a, b float64) float64 {
        if a > 0 && math.IsInf(b, 1) || a < 0 && math.IsInf(b, -1) {
            return a
        }
        if a > 0 && math.IsInf(b, -1) || a < 0 && math.IsInf(b, 1) {
            return b
        }
        return a - math.Floor(a/b)*b
    }


    按位运算符
    按位与& ，按位或| ，二元异或、一元按位取反～ ，左移<<,右移>>

    按位运算符会先把操作数转换为整数，在进行计算，结果也是整数
    右移运算符是无符号右移，空出来的bit使用0填充
    移动-n 个bit，相当于向相反的方向移动 n 个bit。

    按位与、按位或、异或、按位取反运算符可以直接映射位go的对应运算符，不过位移运算符需要处理一下。
    增加左移函数：

    func ShiftLeft(a, n int64) int64 {
        if n >= 0 {
            return a << uint64(n)
        } else {
            return ShiftRight(a, -n)
        }
    }

    因为go里面的位移运算符右边的操作时只能是无符号整数，因此在第一个分支里面对位移的数进行了类型转换。
    如果要移动的位数 < 0,则将左移操作转换为右移操作。

    右移代码如下：

    func ShiftRight(a, n int64) int64 {
        if n >= 0 {
            return int64(uint64(a) >> uint64(n))
        } else {
            return ShiftLeft(a, -n)
        }
    }

    go中，如果右移运算符的左操作数是有符号整数，那么进行的就是有符号右移，空位补充1.不过我们期望的是无符号右移，空位补充0，所以在第一个分支里面需要先将左操作数转换成无符号整数在执行右移擦欧总，然后在将结果转换为有符号整数。如果要移动的位数小于0，则将右移转换为左移。

    比较运算符

    == ～= > >= < <=
    只需要实现 == < <==，就可以了，a ～= b 可以转换为 not （a == b），其他的也可以类似转换，
    a > b ===> b < a, a >= b ====> b <= a

    逻辑运算符与或非

    and or not 。lua没有给逻辑运算符指定专用的符号，而是给她们指定了三个关键字。

    lua会对逻辑与 逻辑或 表达式进行段路求值。
    需要注意的是，逻辑与 逻辑或运算符的结果是操作数之一，并不会转换为boolean值。这个特性经常用于
    简化变量的初始化代码比如：

    function fun(t)
        t = t or {} // if not t then t = {} end
    end

    或者模拟c中的三目运算符。

    max = a > b and a or b --a > b ? a : b

    逻辑非运算符会先把操作数转换为布尔值，然后取非，所以结果也是布尔值。
    在lua中 只有false和nil的值会被转换为false，其他的值都会被转换为true。

    长度运算符

    用于提取字符串或者表的长度。

    print(#"hello") --- >5
    print(#{1,2,3}) -----> 3

    字符串拼接运算符
    用于拼接字符串和数字
    print("a" .. "b" .. "c") ----->abc
    print(1 .. 2 .. 3) --->123

    自动类型转换
    lua运算符会在适当的情况下对操作数进行自动类型转换。

    转换规则

    除了各种类型都可以转换为布尔值以外，自动类型转换主要发生在数字和字符串之间。

        算术运算符
            除法和乘方运算符
                1、如果操作数是整数，则提升为浮点数
                2、如果操作数是字符串，且可以解析为浮点数，则解析为浮点数
                3、然后进行浮点数运算，结果为浮点数。
            其他算术运算符
                1、如果操作数全部为整数，则进行整数运算，结果为整数
                2、否则将操作数转换为浮点数，规则同除法和乘方运算符
                3、然后进行浮点数运算，结果为浮点数。

        按位运算符
            1、如果操作数是整数，则无序转换
            2、如果操作数是浮点数，但实际表示的是整数值，比如1000.0，且没有超出整数的取值范围，则转换为整数。
            3、如果操作数是字符串，且可以解析为整数值，比如“1000”，则解析为整数
            4、如果操作数是字符串，但是无法解析为整数值，不过可以解析为浮点数，比如1000.0，且浮点数可以按照上看的规则转换为整数，则解析位浮点数，然后在转换为整数。
            5、然后进行整数运算，结果也是整数。               


        字符串拼接运算
            1、如果操作数是字符串，则无需转换
            2、如果操作数是数字，整数或者浮点数，则转换为字符串。
            3、然后进行字符串的拼接，结果为字符串。

    如果lua无法将操作数转换为运算符期望的类型，则会导致lua脚本运行错误。

    浮点数转换为整数

    如果浮点数的消小数部分为0，且整数部分没有超出lua整数能够表示的范围，则转换成功，否则失败。
    增加转换函数代码：

    func FloatToInteger(f float64) (int64, bool) {
        i := int64(f)
        return i, float64(i) == f
    }
        

    字符串解析为数字

    parser.go，定义将字符串解析为整数和浮点数的函数。

    package number

    import "strconv"

    func ParseInteger(str string) (int64, bool) {
        i, err := strconv.ParseInt(str, 10, 64)
        return i, err == nil
    }

    func ParseFloat(str string) (float64, bool) {
        f, err := strconv.ParseFloat(str, 64)
        return f, err == nil
    }

    这两个函数都有两个返回值，其中第一个返回值是解析后的值，第二个返回值说明解析是否成功。


    任意值转换为浮点数

    在lua_value.go中定义函数convertToFloat()：

    func convertToFloat(val luaValue) (float64, bool) {
        switch x := val.(type) {
        case int64:
            return float64(x), true
        case float64:
            return x, true
        case string:
            return number.ParseFloat(x)
        default:
            return 0, false
        }
    }

    整数可以直接转换为浮点数，字符串可以调用ParseFloat()解析为浮点数，其他类型的值不能转换为浮点数。
    使用convertToFloat(val luaValue)就可以完善之前的ToNumberX()方法了。
    api_access.go修改其代码为：

    func (self *luaState) ToNumberX(idx int) (float64, bool) {
        val := self.stack.get(idx)
        return convertToFloat(val)
    }

    任意值转化为整数
    lua_value.go 添加函数convertToInteger()方法：

    func convertToInteger(val luaValue) (int64, bool) {
        switch x := val.(type) {
        case int64:
            return x, true
        case float64:
            return number.FloatToInteger(x)
        case string:
            return _stringToInteger(x)
        default:
            return 0, false
        }
    }


    对于浮点数，可以调用之前定义的FLoatToInteger()方法将其转换为整数，对于字符串，可以先试试能够直接解析为整数，如果不能，在尝试将其解析为浮点数，然后转换为整数。_stringToInteger()代码如下：

    func _stringToInteger(s string) (int64, bool) {
        if i, ok := number.ParseInteger(s); ok {
            return i, true
        }
        if f, ok := number.ParseFloat(s); ok {
            return number.FloatToInteger(f)
        }
        return 0, false
    }

    至此，可以使用convertToInteger()函数完善api_access.go中的ToIntegerX（）函数，修改为：

    func (self *luaState) ToIntegerX(idx int) (int64, bool) {
        val := self.stack.get(idx)
        return convertToInteger(val)
    }

    扩展LuaState接口

    前面主要从lua的角度对lua运算符和核运算时可能会进行的自动类型转换进行了说明，在api层面，有几个方法专门用来支持lua运算符。在lua_state.go中，在LuaState接口中添加方法：

     package api

    type LuaType = int
    type ArithOp = int
    type CompareOp = int

    type LuaState interface {
        。。。。。。
        /* Comparison and arithmetic functions */
        Arith(op ArithOp)
        Compare(idx1, idx2 int, op CompareOp) bool
        /* miscellaneous functions */
        Len(idx int)
        Concat(n int)
    }   

     Arith(op ArithOp)：用于执行算术以及按位运算
     Compare(idx1, idx2 int, op CompareOp) bool：用于执行比较运算
     Len(idx int)：用于获取长度的运算
     Concat(n int)：用于执行字符串拼接的运算。

     Arith(op ArithOp)：为了区分具体执行的是什么运算，lua api为每个算术和按位运算符都指定了一个
     运算码，在api/consts.go中添加这些常量定义：

     /* arithmetic functions const */
    const (
        LUA_OPADD  = iota // +
        LUA_OPSUB         // -
        LUA_OPMUL         // *
        LUA_OPMOD         // %
        LUA_OPPOW         // ^
        LUA_OPDIV         // /
        LUA_OPIDIV        // //
        LUA_OPBAND        // &
        LUA_OPBOR         // |
        LUA_OPBXOR        // ~
        LUA_OPSHL         // <<
        LUA_OPSHR         // >>
        LUA_OPUNM         // -
        LUA_OPBNOT        // ~
    )

    Compare(idx1, idx2 int, op CompareOp) bool：
        lua api只给 == 、< 、<=分配了运算码。把这三个运算码常量定义在consts.go中：

        /* comparison functions */
        const (
            LUA_OPEQ = iota // ==
            LUA_OPLT        // <
            LUA_OPLE        // <=
        )

    方法的实现

     Arith(op ArithOp)：

     可以执行算术以及按位运算，局的运算由参数指定，操作数则从栈顶弹出。对于二元运算，该方法会从栈顶弹出
     两个值进行计算，然后将结果push到栈顶。

     对于一元运算，该函数从栈顶pop一个值进行计算，然后把结果push到栈顶。

     该方法是lua运算符在api层面上的映射，也需要遵守前面的自动类型转换的规则。为了减少代码的重复，我们先把算术运算和位移运算统一映射位go语言的运算符，或者之前预先定义好的函数。
     创建state/api_arith.go，并在其中定义若干函数类型的变量：

    package state

    import "math"
    import . "luago/api"
    import "luago/number"


    var (
        iadd  = func(a, b int64) int64 { return a + b }
        fadd  = func(a, b float64) float64 { return a + b }
        isub  = func(a, b int64) int64 { return a - b }
        fsub  = func(a, b float64) float64 { return a - b }
        imul  = func(a, b int64) int64 { return a * b }
        fmul  = func(a, b float64) float64 { return a * b }
        imod  = number.IMod
        fmod  = number.FMod
        pow   = math.Pow
        div   = func(a, b float64) float64 { return a / b }
        iidiv = number.IFloorDiv
        fidiv = number.FFloorDiv
        band  = func(a, b int64) int64 { return a & b }
        bor   = func(a, b int64) int64 { return a | b }
        bxor  = func(a, b int64) int64 { return a ^ b }
        shl   = number.ShiftLeft
        shr   = number.ShiftRight
        iunm  = func(a, _ int64) int64 { return -a }
        funm  = func(a, _ float64) float64 { return -a }
        bnot  = func(a, _ int64) int64 { return ^a }
    )


    使用两个参数、返回一个值的函数来统一表示lua的运算符，一元运算符只需要简单地忽略第二个参数就可以了。
    乘方 元算符可是使用go数学库中的Pow实现，整除、取模、位移运算的函数前面已经实现了，其他的运算符可以直接映射成go的运算符就可以了。
    还需要一个结构体来容纳整数以及浮点数类型的运算。在api_arith.go中定义operator结构体：

    type operator struct {
        integerFunc func(int64, int64) int64
        floatFunc   func(float64, float64) float64
    }

    然后定义一个slice，里面是各种运算，需要注意的是，要和前面定义的lua运算码常量的顺序要一致。

    var operators = []operator{
        operator{iadd, fadd},
        operator{isub, fsub},
        operator{imul, fmul},
        operator{imod, fmod},
        operator{nil, pow},
        operator{nil, div},
        operator{iidiv, fidiv},
        operator{band, nil},
        operator{bor, nil},
        operator{bxor, nil},
        operator{shl, nil},
        operator{shr, nil},
        operator{iunm, funm},
        operator{bnot, nil},
    }


    接下来就是在api_arith.go中定义Arith（）函数了

    func (self *luaState) Arith(op ArithOp) {
        var a, b luaValue // operands
        b = self.stack.pop()
        if op != LUA_OPUNM && op != LUA_OPBNOT {
            a = self.stack.pop()
        } else {
            a = b
        }

        operator := operators[op]
        if result := _arith(a, b, operator); result != nil {
            self.stack.push(result)
        } else {
            panic("arithmetic error!")
        }
    }


    该函数根据情况从lua栈中pop一个或者两个操作数，然后按照索引取出相应的operator实例，最后调用_arith() 进行计算。如果计算结果不是nil，则表示操作数（或者可以转换为）运算符规定的类型，将计算结果push到
    lua栈即可，否则调用panic()终止程序的执行。

    func _arith(a, b luaValue, op operator) luaValue {
        if op.floatFunc == nil { // bitwise
            if x, ok := convertToInteger(a); ok {
                if y, ok := convertToInteger(b); ok {
                    return op.integerFunc(x, y)
                }
            }
        } else { // arith
            if op.integerFunc != nil { // add,sub,mul,mod,idiv,unm
                if x, ok := a.(int64); ok {
                    if y, ok := b.(int64); ok {
                        return op.integerFunc(x, y)
                    }
                }
            }
            if x, ok := convertToFloat(a); ok {
                if y, ok := convertToFloat(b); ok {
                    return op.floatFunc(x, y)
                }
            }
        }
        return nil
    }


    按位运算期望操作数都是（或者可以转换为）整数，运算结果也是整数。加、减、乘、除、整除、取反运算会在两个操作数都是整数的时候进行整数运算，结果也是整数。对于其他的情况，则尝试将操作数转换为浮点数，在执行计算，运算结果也是浮点数。

    Compare()方法
    对指定索引处的两个值进行比较，返回结果。此函数不会改变栈的状态。

    在sate/api_compare.go中定义：

    package state

    import . "luago/api"


    func (self *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
        if !self.stack.isValid(idx1) || !self.stack.isValid(idx2) {
            return false
        }

        a := self.stack.get(idx1)
        b := self.stack.get(idx2)
        switch op {
        case LUA_OPEQ:
            return _eq(a, b)
        case LUA_OPLT:
            return _lt(a, b)
        case LUA_OPLE:
            return _le(a, b)
        default:
            panic("invalid compare op!")
        }
    }

    Compare()方法先按照索引取出两个操作数，然后根据操作码调用_eq()\_lt()\_le()函数进行比较。
    _eq()函数用于比较两个值是否相等。

    func _eq(a, b luaValue) bool {
        switch x := a.(type) {
        case nil:
            return b == nil
        case bool:
            y, ok := b.(bool)
            return ok && x == y
        case string:
            y, ok := b.(string)
            return ok && x == y
        case int64:
            switch y := b.(type) {
            case int64:
                return x == y
            case float64:
                return float64(x) == y
            default:
                return false
            }
        case float64:
            switch y := b.(type) {
            case float64:
                return x == y
            case int64:
                return x == float64(y)
            default:
                return false
            }
        default:
            return a == b
        }
    }


    只有当两个操作数载Lua语言层面具有相同的类型时，== 运算才有可能返回true。nil、布尔以及字符串类型的等于操作比较简单。整数和浮点数仅仅在lua实现层面有差别，在lua语言层面统一表现为数字类型，因此需要相互转换。其他类型的值暂时先按照引用进行比较，以后在完善。

    func _lt(a, b luaValue) bool {
        switch x := a.(type) {
        case string:
            if y, ok := b.(string); ok {
                return x < y
            }
        case int64:
            switch y := b.(type) {
            case int64:
                return x < y
            case float64:
                return float64(x) < y
            }
        case float64:
            switch y := b.(type) {
            case float64:
                return x < y
            case int64:
                return x < float64(y)
            }
        }
        panic("comparison error!")
    }

    < 操作仅仅对数字和字符串类型有意义，其他的情况以后再说，暂时调用panic终止程序的执行。_le（） 和
    _lt()基本一样。

    func _le(a, b luaValue) bool {
        switch x := a.(type) {
        case string:
            if y, ok := b.(string); ok {
                return x <= y
            }
        case int64:
            switch y := b.(type) {
            case int64:
                return x <= y
            case float64:
                return float64(x) <= y
            }
        case float64:
            switch y := b.(type) {
            case float64:
                return x <= y
            case int64:
                return x <= float64(y)
            }
        }
        panic("comparison error!")
    }

    Len(idx int)：用于获取长度的运算
    访问指定索引处的值取其长度，然后push到栈顶。state/api_misc.go

    package state

    func (self *luaState) Len(idx int) {
        val := self.stack.get(idx)

        if s, ok := val.(string); ok {
            self.stack.push(int64(len(s)))
        } else {
            panic("length error!")
        }
    }

     暂时值考虑字符串的长度，对其他情况先调用panic终止程序的执行，以后在完善。


     Concat(n int)：用于执行字符串拼接的运算。

     该方法从栈顶pop n 个值，然后对这些值进行拼接，然后把结果push 到栈顶。

    state/api_misc.go下添加方法：

    func (self *luaState) Concat(n int) {
        if n == 0 {
            self.stack.push("")
        } else if n >= 2 {
            for i := 1; i < n; i++ {
                if self.IsString(-1) && self.IsString(-2) {
                    s2 := self.ToString(-1)
                    s1 := self.ToString(-2)
                    self.stack.pop()
                    self.stack.pop()
                    self.stack.push(s1 + s2)
                    continue
                }

                panic("concatenation error!")
            }
        }
        // n == 1, do nothing
    }

    如果n 为0，不弹出任何值，直接往栈顶push一个空的字符串，否则，将栈顶的两个值弹出，进行拼接，然后将结果push到栈顶。这个过程一直持续，直到n个值都处理完。在弹出栈顶值之前，需要先调用ToString()将其转换为字符串，如果转换失败，暂时先调用panic终止程序的执行，以后在完善。

    单元测试

    package main

    import "fmt"
    import . "luago/api"
    import _ "luago/binchunk"
    import "luago/state"

    func main() {
        ls := state.New()
        ls.PushInteger(1)
        ls.PushString("2.0")
        ls.PushString("3.0")
        ls.PushNumber(4.0)
        printStack(ls)

        ls.Arith(LUA_OPADD)
        printStack(ls)
        ls.Arith(LUA_OPBNOT)
        printStack(ls)
        ls.Len(2)
        printStack(ls)
        ls.Concat(3)
        printStack(ls)
        ls.PushBoolean(ls.Compare(1, 2, LUA_OPEQ))
        printStack(ls)
    }

    先创建一个luaState()实例，然后push一些值，进行各种运算 ，并将整个栈打印出来。

===========================================


虚拟机小雏形--lua虚拟机指令

将lua栈当作虚拟的寄存器来使用。

    添加LuaVM接口

    我们在读书的时候，往往很难一口气读完，如果中间累了，或者需要停下来思考问题，或者有其他的事情需要处理，就暂时把当前页数记录下来，合上书，然后休息、思考或者处理问题，然后在回来继续阅读。我们也经常跳到之前已经度过的某一页，去回顾一些内容，或者翻到后面黑还没有读到的某一页，去预览一些内容，然后在回到当前页继续阅读。

    就如同人类需要使用大脑（或者书签）去记录正在阅读的书的页数一样，计算机页需要一个程序计数器(Program Counter  ,简称PC)来记录正在执行的指令。与人类不同的是，计算机不需要休息和思考，页没有烦恼的事情去处理，不过需要等待io，一旦接听了电源，就会不知疲倦地计算pc,取出当前指令、执行指令、如此往复，直到指令全部处理完或者断电。

    Lua解释器在执行一段Lua脚本之前，会先把它包装在一个main函数里面，并编译成lua虚拟机的指令序列，然后联通其他信息一起，打包成一个二进制的chunk，然后lua虚拟机就会接管这个chunk，执行里面的指令。和真实的计算机一样，lua虚拟机页需要使用程序计数器。可以使用如下的为代码来米奥术lua虚拟机的内部循环：

    while（ true ）{
        1、计算PC
        2、取出当前指令
        3、执行当前指令
    }

    由于lua虚拟机采用的是定长的指令集，每条指令占用固定的4个字节，所以上面虚幻里面的前两个步骤会比较简单。至于步骤3，则可以根据指令的操作码写成一个巨大的switch--case语句，乳以下伪代码所示：

    switch opcode{
        case OP_MOVE:执行move指令
        case OP_LOADK:加载常量
    }

    虽然lua虚拟机目前只有47条指令，不过如果真的吧所有的实现袭击都可在一个switch-case语句里面，则会导致代码太长，难以阅读和维护。基于这个原因，我们会把每条指令都单独实现成一个函数，调整后的伪代码为：

     switch opcode{
        case OP_MOVE:move（instruction）
        case OP_LOADK:laodk（instruction）
    }   


    那么对于每条指令所对应的函数来说，只有指令这么一个参数就够用了吗？显示是不够的，lua虚拟机是基于寄存器的虚拟机，因此虚拟机指令集中的大部分指令都会涉及到寄存器的操作也就没什么好奇怪的啦。也就是说，我们至少还要给指令听过寄存器，以便能够对其进行操作。


    由于lua api也允许我们使用索引来操作隐藏在Lua State内部的lua 栈，所以正好可以使用Lua State来模拟寄存器。由于大部分lua运算符在lua虚拟机指令集中也有对应的指令，因此利用Lua State也可以很容易地实现运算符相关的指令。因此，把Lua State作为第二个参数传递给指令函数是个不错的选择。

    不过要完整实现Lua虚拟机的指令集，仅仅依赖lua api还不够，比如说loadk指令就需要查看二进制chunk的常量表、从中取出某个常量，放到指定的寄存器中。但lua api并没有将常量表这种函数内部的细节暴露给用户，因此也无法通过api的方式操作常量表。因为我们不想给lua api私自添加任何方法，所以我们要引入一个新的LuaVM接口，让其扩展现有的LuaState接口，然后添加几个必要的方法以便满足指令实现函数的需要。

    定义LuaVM接口
    lua_vm.go定义这个接口：

    package api

    type LuaVM interface {
        LuaState
        PC() int    //返回当前PC
        AddPC(n int) //修改PC 用于实现跳转指令
        Fetch() uint32 //取出当前的指令，将PC指向下一条指令
        GetConst(idx int) //将指定的常量push到栈顶
        GetRK(rk int) // 将指定的常量或者栈值push到栈顶
    }

    LuaVM接口扩展了LuaState接口，并添加了5个方法，其中

    AddPC(n int) 用于修改当前的PC，这个方法是实现跳转指令所必须的。

    Fetch() uint32 用于取出当前的指令，同时递增PC，让其指向下一条指令。这个方法主要是前面提到的虚拟机循环使用，不过loadk等少数几个指令也会用到。

    GetConst(idx int) 用于从常量表中取出指定的常量并push到栈顶。loadk和ladkx这两个指令会用到这个方法。

    GetRK(rk int) 根据情况从常量表中提取常量或者从栈中提取值，然后push 到栈顶。
    
    PC() int 用于返回当前的pc，测试用。

    接下来对结构体luaState进行改造，让其实现接口LuaVM

    改造结构体luaState

    package state

    import "lunago/binchunk"

    type luaState struct {
        stack *luaStack
        proto *binchunk.Prototype //added 保存函数原型
        pc    int //added  程序计数器
    }

    这样就可以从中提取指令或者常量了。

    luaState增加字段后，New（）也需要做一下调整

    func New(stackSize int, proto *binchunk.Prototype) *luaState {
        return &luaState{
            stack: newLuaStack(stackSize),
            proto: proto,
            pc:    0,
        }
    }

    给New函数增加了两个参数，第一个参数用于指定Lua栈的初始容量，第二个参数传入函数原型，以初始化proto字段。由于虚拟机肯定是从第一条指令开始执行的，因此pc字段初始化为0就可以了。

    实现LuaVM接口

    有了改进版的luaState结构体，实现LuaVM接口就比较容易了。为了便于修改，把LuaVM接口定义的5个方法实现在api_vm.go中。

    package state

    func (self *luaState) PC() int {
        return self.pc
    }

    func (self *luaState) AddPC(n int) {
        self.pc += n
    }

    Fetch（）根据PC索引从函数原型的指令表中取出当前的指令，然后把PC+1，这样下次在调用该方法取出的就是下一条指令。

    func (self *luaState) Fetch() uint32 {
        i := self.proto.Code[self.pc]
        self.pc++
        return i
    }

    GetConst()根据索引从函数原型的常量表中取出一个常量值，然后将其push到栈顶

    func (self *luaState) GetConst(idx int) {
        c := self.proto.Constants[idx]
        self.stack.push(c)
    }  

    GetRK(）根据情况调用GetConst（）把某个常量push到栈顶，或者调用PushValue（）把某个索引处的栈值push到栈顶。

     func (self *luaState) GetRK(rk int) {
        if rk > 0xFF { // constant
            self.GetConst(rk & 0xFF)
        } else { // register
            self.PushValue(rk + 1)
        }
    }   

    传递给GetRK()的参数实际上就是iABC模式指令里面的OpArgK类型的参数。这种类型的参数占用9个bit，如果最高位是1，则参数里面存放的是常量表索引，把最高位去掉就可以得到索引值；否则如果最高位是0，参数里面存放的就是寄存器的索引值，需要注意的是，Lua虚拟机指令操作数里面携带的寄存器索引是从0开始的，而Lua api里面的栈索引是从1开始的，所以当需要把寄存器的索引当成栈的索引使用的时候，需要对寄存器的索引+1.

    实现lua虚拟机的指令    
    lua虚拟机的指令集一共有47个，其中EXTRSAARG指令实际上只能用来扩展其他指令的操作数，并不能单出执行，所以真正的指令只有46条。


    指令实现放在四个文件中，vm文件夹下，
    运算符相关指令有22 条，放在inst_operators.go;
    加载类指令放在inst_load.go;for循环指令仿造inst_for.go中，move、jmp指令放在inst_misc.go中。

    move 和 jmp指令

    MOVE
    move指令（iABC 模式）把源寄存器（索引由操作数B指定）里面值移动到目标寄存器
    （索引由操作数A指定）里。

    R（A）：= R（B）

    move指令常常用于局部变量赋值以及参数的传递。虽然说是move指令，实际上叫做copy指令可能会更贴切一些，因为源寄存器的值还原封不动的待在原地。在inst_misc.go中实现move指令：

    func move(i Instruction, vm LuaVM) {
        a, b, _ := i.ABC()
        a += 1
        b += 1

        vm.Copy(b, a)
    }

    先解码指令，得到目标寄存器和源寄存器的索引，然后把她们转换成栈索引，最后调用lua api提供的Copy（）拷贝栈值。寄存器索引+1才是对应的栈索引。

    从这里可以看到，lua代码中的局部变量实际上是存储在寄存器里面的。因为move等指令使用操作数A（占用8个bit）来表示目标寄存器的索引，因此lua函数使用的局部变量不能超过255个。实际上，lua编译器把函数的局部变量的数量限制在了200个以内，如果超过这个数量，函数就复发通过编译。

    jmp

    jmp指令（iAsBx模式）执行无条件跳转，该指令往往和TEST指令配合使用，不过也可能会单独出现，比如lua也支持标签和goto语句。

    jmp指令不会改变寄存器的状态。代码如下：

    func jmp(i Instruction, vm LuaVM) {
        a, sBx := i.AsBx()

        vm.AddPC(sBx)
        if a != 0 {
            panic("todo: jmp!")
        }
    }

    jmp指令的操作数A和Upvalue有关，以后在说，暂时先不管。

    加载指令

    用于把nil值、布尔值、或者常量表中的常量值加载到寄存器中。

    1、loadnil
        loadnil指令（iABC模式）用于给连续n个寄存器放置nil值。寄存器的起始索引由操作数A指定，寄存器数量由操作数B指定，操作数C没有使用。loadnil指令可以用如下伪代码表示：

        R（A），R（A+1），。。。，R（A+B）：= nil

        在lua代码里，局部变量的默认初始值就是nil。loadnil指令常用于给连续的n个局部变量设置初始值。

        inst_load.go中实现代码：

        func loadNil(i Instruction, vm LuaVM) {
            a, b, _ := i.ABC()
            a += 1

            vm.PushNil()
            for i := a; i <= a+b; i++ {
                vm.Copy(-1, i)
            }
            vm.Pop(1)
        }


        lua编译器在编译函数，生成指令表的时候，会把指令执行阶段所需的寄存器数量预先计算好，保存在函数的原型里面。这里嘉定虚拟机在执行第一条指令之前，已经根据这一信息调用SetTop函数保留了必要数量的占空间。有了这个假设，就可以先调用PushNil（）函数，往栈顶push一个nil值，然后连续调用Copy（）将Nil值符知道指定的寄存器中，最后调用Pop（）把一开始push到栈顶的那个nil值弹出，让栈顶指针恢复原状。

        2、loadbool
        loadbool指令（iABC指令）给单个寄存器设置布尔值。寄存器索引由操作数A指定，布尔值由操作数B指定（0表示false，非0表示true），如果寄存器C非0，则跳过下一条指令。伪代码如下：
        R（A） := (bool)B; if(C) pc++

        loadbool指令可以单独使用 也可以和比较指令结合。

         func loadBool(i Instruction, vm LuaVM) {
            a, b, c := i.ABC()
            a += 1

            vm.PushBoolean(b != 0)
            vm.Replace(a)

            if c != 0 {
                vm.AddPC(1)
            }
        }       

        3、loadk loadkx

        loadk(iABx模式)将常量表里面的某个常量加载到指定的寄存器中，寄存器的索引由操作数A指定，常量表的索引由操作数Bx指定。如果Kst（N）表示常量表中的第N个常量，那么loadk指令可以使用如下伪代码表示：

        R(A):=Kst(N)

        在lua函数里面出现的字面量（主要是字符串和数字）会被lua解释器收集起来，放到常量表里面，

        func loadK(i Instruction, vm LuaVM) {
            a, bx := i.ABx()
            a += 1

            vm.GetConst(bx)
            vm.Replace(a)
        }      

        先调用之前准备好的GetConst（）函数吧给定的常量push到栈顶，然后调用Replace（）把它移动到指定的索引处。操作数Bx占用18个bit，能表示的最大无符号整数是262143，大部分lua函数的常量表大小都不会超过这个数字，因此这个限制通常不是神没问题。不过lua也经常被当作数据描述语言使用，因此常量表的大小可能会超出这个现实也并不奇怪，为了应对这种情况，lua还提供了一条loadkx指令。

        loadkx指令（也是iABx模式）需要和EXTEAARG指令（iAx模式）配合使用。用后者的Ax操作数来指定常量的索引。Ax操作数占用26个bit，可以表达的最大无符号整数是67108864，可以满足大部分情况了。

        // R(A) := Kst(extra arg)
        func loadKx(i Instruction, vm LuaVM) {
            a, _ := i.ABx()
            a += 1
            ax := Instruction(vm.Fetch()).Ax()

            //vm.CheckStack(1)
            vm.GetConst(ax)
            vm.Replace(a)
        }

算术运算符
    1、二元算术运算指令
    二元算术运算指令（iABC 模式），对连个寄存器或者常量值（索引由操作数B、C指定）进行运算， 将结果放到另一个寄存器中（随你由操作数A指定）。如果用RK（N）表示寄存器或者常量值，那么二元算术运算指令的伪代码可以如下表示：

    R（A）:= RK(B) op RK(C)  

    vm/inst_operators.go中的实现代码：
    package vm

    import . "lunago/api"

    func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
        a, b, c := i.ABC()
        a += 1

        vm.GetRK(b)
        vm.GetRK(c)
        vm.Arith(op)
        vm.Replace(a)
    }

    _binaryArith(），以后会使用它来实现二元算术运算指令。

    先调用之前准备好的GetRK（）函数把两个操作数push 到栈顶，然后调用Arith（）进行算术运算，算术运算完毕之后，操作数已经从栈顶pop，取而代之的是运算结果，调用Replace（）把它移动到指定的寄存器即可。

    2、一元算术运算指令

    一元算术运算指令（iABC 模式），对操作数B所指定的寄存器里面的值进行运算，
    然后把结果放到操作数 A 所指定的寄存器中，操作数 C 没有使用。可以使用如下伪代码表示：

    R（A）：= op R（B）

    func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
        a, b, _ := i.ABC()
        a += 1
        b += 1

        vm.PushValue(b)
        vm.Arith(op)
        vm.Replace(a)
    }

    有了上面这两个函数，实现算术运算符就容易多了。


    func add(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPADD) }  // +
    func sub(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSUB) }  // -
    func mul(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMUL) }  // *
    func mod(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMOD) }  // %
    func pow(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPPOW) }  // ^
    func div(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPDIV) }  // /
    func idiv(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPIDIV) } // //
    func band(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBAND) } // &
    func bor(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPBOR) }  // |
    func bxor(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBXOR) } // ~
    func shl(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHL) }  // <<
    func shr(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHR) }  // >>
    func unm(i Instruction, vm LuaVM)  { _unaryArith(i, vm, LUA_OPUNM) }   // -
    func bnot(i Instruction, vm LuaVM) { _unaryArith(i, vm, LUA_OPBNOT) }  // ~   

    长度以及拼接指令

    len
    len指令（iABC模式）进行的操作和一元算术运算指令类似，伪代码表示如下：

    R（A）：= length of R(B)

    func length(i Instruction, vm LuaVM) {
        a, b, _ := i.ABC()
        a += 1
        b += 1

        vm.Len(b)
        vm.Replace(a)
    }

    2、concat
        cancat(iABC 模式)，将连续的n个寄存器（起止索引分别由操作数B、C指定）里的值拼接，将结果放到另一个寄存器中（索引由操作数A指定）。

        R（A）：= R（B）.. ... .. R(C)

        func concat(i Instruction, vm LuaVM) {
            a, b, c := i.ABC()
            a += 1
            b += 1
            c += 1

            n := c - b + 1
            vm.CheckStack(n)
            for i := b; i <= c; i++ {
                vm.PushValue(i)
            }
            vm.Concat(n)
            vm.Replace(a)
        }

    在实现前面的指令时，最多只是往栈顶push了一两个值，所以我们可以在创建Lua栈的时候把容量设置的稍大一些，这样在push少量的值之前，就不需要检查栈的剩余空间了。但是concat指令则有所不同，因为进行拼接的值的数量不是固定的，所以在吧这些值push到栈顶之前，必须调用CheckStack（）确保还有足够的空间可以容纳这些值，否则可能会导致溢出。


    比较指令
    比较指令（iABC 模式），比较寄存器或者常量表里面的两个值（索引分别由操作数B、C指定），
    如果比较结果和操作数A（转换为布尔值）匹配，则跳过下一条指令。比较指令不会改变寄存器的状态。

    if（RK（B）op RK（C）～= A）then pc++

    比较指令对应lua语言里面的比较运算符，当用于赋值的时候，需要和loadbool指令搭配使用。

    在inst_operators.go文件中定义函数：

    func _compare(i Instruction, vm LuaVM, op CompareOp) {
        a, b, c := i.ABC()

        vm.GetRK(b)
        vm.GetRK(c)
        if vm.Compare(-2, -1, op) != (a != 0) {
            vm.AddPC(1)
        }
        vm.Pop(2)
    }

    先调用GetRK（）把两个要比较的值push到栈顶，然后调用Compare（）执行比较运算，如果比较结果和操作数A一致，就把pc++。因为Compare（）方法并没有把栈顶的值弹出，因此我们需要自己调用Pop（）清理栈顶。


    func eq(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPEQ) } // ==
    func lt(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLT) } // <
    func le(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLE) } // <=

    逻辑运算指令

    对应lua中的逻辑运算符

    1、not
        not指令（iABC模式）进行的操作和一元算术运算符指令类似，对应lua中的逻辑非运算符
        R（A）：= not R（B）

        func not(i Instruction, vm LuaVM) {
            a, b, _ := i.ABC()
            a += 1
            b += 1

            vm.PushBoolean(!vm.ToBoolean(b))
            vm.Replace(a)
        }

    2、testset
        testset指令（iABC 模式），破案段寄存器B（索引由操作数B指定）中的值转换为布尔值之后，是否和操作数C表示的布尔值一致，如果一样，则将寄存器B中的值符知道寄存器A中，索引由操作数A指定，否则跳过下一条指令。

        if（R（B）布尔比较 C ） then R（A）：= R（B）else pc++

        该指令对应lua语言中的逻辑与和逻辑或运算符

        func testSet(i Instruction, vm LuaVM) {
            a, b, c := i.ABC()
            a += 1
            b += 1

            if vm.ToBoolean(b) == (c != 0) {
                vm.Copy(b, a)
            } else {
                vm.AddPC(1)
            }
        }


        3、test
        test指令（iABC模式），判断寄存器A（索引由操作数A指定）中的值转换为布尔值后是否和操作数C表示的布尔值一致，如果一致，则跳过下一条指令。test指令不使用操作数B，也不会改变寄存器的状态。

        if not（R（A）<=> C）then pc++

        test是testset的特殊形式

        func test(i Instruction, vm LuaVM) {
            a, _, c := i.ABC()
            a += 1

            if vm.ToBoolean(a) != (c != 0) {
                vm.AddPC(1)
            }
        

    for循环指令
    lua语言中的循环有两种形式：数值形式Numerical以及通用形式Generic。数值for循环用于按照一定的步长遍历某个范围内的数值，通用for循环主要用于遍历表。

    数值for循环需要借助两条指令来实现：forrep、forloop，其中forrep可以使用伪代码表示为：

    R（A）-=R（A+2）；pc+=sBx

    forloop可以表示为：

    R（A）+=R（A+2）；
    if  R（A）<?= R（A+1）then {
        pc+=sBx; R（A+3）=  R（A）
    }

    从反编译处的局部变量表可知，lua编译器为了实现for循环，使用了三个特殊的局部变量，这三个特殊的局部变量的程程都包含（），属于非法标示符，这样就避免了和程序中出现的普通变量重名的可能性。由名称可知，这三个局部变量分别存放index，limit以及step，并且在循环开始之前就已经预先初始化好了---通过三个loadk指令。

    这三个特殊的变量正好对应伪代码中的 R（A）、 R（A + 1） R（A + 2）三个寄存器，自己在for循环里面定义的变量i则对应寄存器 R（A + 3）。forrep指令执行的操作起始就是在循环开始之前预先给index - step，然后跳转到forloop指令，正式开始循环。

    forloop指令则是献给index+step，然后判断数值是否还在范围内，如果已经超出了范围，则循环结束；如果为超过范围，就把数值复制给用户定义的局部变量，然后跳转到循环体的内部开始执行具体的代码块。

    还有一点需要解释的是，“<?=”,当step是正数的时候，含义为“<=”，也就是说继续循环的条件是index<=limit;当step是负数的时候，含义是“>=”，循环继续的条件就是index>=limit。

    具体实现：vm/inst_for.go

    package vm

    import . "lunago/api"

    func forPrep(i Instruction, vm LuaVM) {
        a, sBx := i.AsBx()
        a += 1

        if vm.Type(a) == LUA_TSTRING {
            vm.PushNumber(vm.ToNumber(a))
            vm.Replace(a)
        }
        if vm.Type(a+1) == LUA_TSTRING {
            vm.PushNumber(vm.ToNumber(a + 1))
            vm.Replace(a + 1)
        }
        if vm.Type(a+2) == LUA_TSTRING {
            vm.PushNumber(vm.ToNumber(a + 2))
            vm.Replace(a + 2)
        }

        vm.PushValue(a)
        vm.PushValue(a + 2)
        vm.Arith(LUA_OPSUB)
        vm.Replace(a)
        vm.AddPC(sBx)
    }


    func forLoop(i Instruction, vm LuaVM) {
        a, sBx := i.AsBx()
        a += 1

        // R(A)+=R(A+2);
        vm.PushValue(a + 2)
        vm.PushValue(a)
        vm.Arith(LUA_OPADD)
        vm.Replace(a)

        isPositiveStep := vm.ToNumber(a+2) >= 0
        if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
            !isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {

            // pc+=sBx; R(A+3)=R(A)
            vm.AddPC(sBx)
            vm.Copy(a, a+3)
        }
    }



指令的分发
    前面提到过，可以使用一个巨大的switch-case语句进行指令的分发Dispatch。因为我们已经定义好了一个指令表，因此只要对它稍加扩展就可以通过查表的方式进行指令的分发。修改vm/opcodes.go中的结构体opcode：

    增加
    import "lunago/api"

    type opcode struct {
        ....//增加如下字段
        action   func(i Instruction, vm api.LuaVM) //存放指令的实现函数。
    }

    给opcode字典添加了字段action，用来存放指令的实现函数。继续修改vm/opcodes.go中的的代码，
    修改opcodes变量的初始化代码，给前面的指令添加上action。

    修改vm/instruction.go中的代码，给Instruction类型增加一个Execute（）函数：

    package vm

    import "lunago/api" //新增加的

    。。。。

    func (self Instruction) Execute(vm api.LuaVM) {
        action := opcodes[self.Opcode()].action
        if action != nil {
            action(self, vm)
        } else {
            panic(self.OpName())
        }
    }


    在指令表的辅助下，指令的分发变得异常简单，Execute（）先从指令里面提取操作码，然后根据操作码从指令表中查找对应的指令实现方法，最后调用指令实现方法执行指令。

单元测试
    修改main函数

    package main

    import "fmt"
    import "io/ioutil"
    import "os"
    import . "lunago/api"
    import "lunago/state"
    import "lunago/binchunk"
    import . "lunago/vm"

    func main() {
        if len(os.Args) > 1 {
            data, err := ioutil.ReadFile(os.Args[1])
            if err != nil {
                panic(err)
            }

            proto := binchunk.Undump(data)
            luaMain(proto)
        }
    }


    func luaMain(proto *binchunk.Prototype) {
        nRegs := int(proto.MaxStackSize)
        ls := state.New(nRegs+8, proto)
        ls.SetTop(nRegs)
        for {
            pc := ls.PC()
            inst := Instruction(ls.Fetch())
            if inst.Opcode() != OP_RETURN {
                inst.Execute(ls)

                fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
                printStack(ls)
            } else {
                break
            }
        }
    }

    可以从main函数原型中获取到运行该函数所需的寄存器数量，因为指令实现也需要少量的占空间，所以时间创建的Lua栈容量要比寄存器数量稍大一些。luaState结构体实例创建好了以后，调用SetTop（）方法在栈里面预留出寄存器的空间，剩余的占空间留给指令的实现函数使用，剩下的代码就是指令循环了：取出指令，递增pc，执行指令、打印指令和栈信息，知道遇到返回指令为止。

    有了lua虚拟机，准备一个测试脚本sum.lua
    反编译一下，看看都生成了什么指令

    $ luac -l sum.lua 

    main <sum.lua:0,0> (11 instructions at 0x7fba9dd00070)
    0+ params, 6 slots, 1 upvalue, 5 locals, 4 constants, 0 functions
        1   [7] LOADK       0 -1    ; 0
        2   [9] LOADK       1 -2    ; 1
        3   [9] LOADK       2 -3    ; 100
        4   [9] LOADK       3 -2    ; 1
        5   [9] FORPREP     1 4 ; to 10
        6   [10]    MOD         5 4 -4  ; - 2
        7   [10]    EQ          0 5 -1  ; - 0
        8   [10]    JMP         0 1 ; to 10
        9   [11]    ADD         0 0 4
        10  [9] FORLOOP     1 -5    ; to 6
        11  [13]    RETURN      0 1

    除了return指令，其他的都实现好了。因为遇到return就结束循环，因此暂时没什么问题。后续在完善。
    
 ================
 
 表       

lua值提供了一种数据结构---表。不仅可以直接当成数组和列表使用，也可以用它实现其他的各种数据结构。表对lua也很重要，很多东西都依赖表。

介绍
lua表本质上就是关联数组（Associative Array ，Dictionary 、Map），里面存放的是两两关联的key value对。lua提供了表构造器（table constructor）表达式，方便了表的构建。

可以使用下表报答时给key赋值，或者根据key访问value。除了nil值还有浮点数的NaN之外，任何lua值都可以当作key来使用，value可以是任意lua值，包括nil NaN。如果一个表的key全部是字符串类型，这个表为记录record。

表可以实现很多结构。按照惯例，如果某个表的key全是正整数，就是list或者array，正整数并不包括0，所以lua数组的索引从1开始，如果表的构造器省略key和=号，仅仅列出值，创建出来的的就是数组。

如果数组中存在nil值，称这些nil值为洞Hole，如果一个数组中没有nil，这个束素就是序列sequeue。对于序列，可以使用长度运算符获取它的长度。

当key在表内使用的时候，实际表示整数（且在整数范围内）的浮点数会被转换成相应的整数。


表的内部实现：

直接使用go的内置类型map就可以了。比如 type luaTable map[luaValue] luaValue

Lua 5.0之前，也是简单的使用hashtable来实现lua表的，不过由于实践中，数组的使用非常频繁，为了专门优化数组的效率，lua 5。0开始改用混合数据结构来实现表，这种结构同时包含了数组还有哈希表两部分。如果表的key全是连续的正整数，那么哈希表就是空的，值全部按照索引存储在数组中，这样lua数组无论是在空间利用率上（不需要显式存储key），还是时间效率上（可以按照索引存取值，不需要计算哈希码）都和真正的数组相差不大。如果表没有被当成数组使用，那么数据完全存储在hashtable中，数组部分就是空的，，也没什么损失。

现在参照lua官方做法，使用数组和哈希表的混合方式实现lua表。state/lua_table.go在里面定义结构体
luaTable。

    
    package state

    import "math"
    import "lunago/number"

    type luaTable struct {
        arr  []luaValue
        _map map[luaValue]luaValue
    }

    继续添加函数newLuaTable(）

    func newLuaTable(nArr, nRec int) *luaTable {
        t := &luaTable{}
        if nArr > 0 {
            t.arr = make([]luaValue, 0, nArr)
        }
        if nRec > 0 {
            t._map = make(map[luaValue]luaValue, nRec)
        }
        return t
    }

    该函数接受两个参数，用于预估表的用途和容量，如果参数nArry大于0，说明表可能是当作数组使用的，先创建数组部分；如果参数nRec大于0，说明表可能是当作记录使用的，先创建哈希表部分。

    接下来给结构体luaTable定义三个方法，用来根据key存取值以及计算数组的长度。

    func (self *luaTable) get(key luaValue) luaValue {
        key = _floatToInteger(key)
        if idx, ok := key.(int64); ok {
            if idx >= 1 && idx <= int64(len(self.arr)) {
                return self.arr[idx-1]
            }
        }
        return self._map[key]
    }

    get（）方法根据key从表里面查找值。如果key是整数(或者能够转换为整数的浮点数)，且在啊数组索引范围内，直接按照索引访问数组部分就可以了；否则从hashtable查找值。

    _floatToInteger()尝试吧服点类型的key转换为整数：

    func _floatToInteger(key luaValue) luaValue {
        if f, ok := key.(float64); ok {
            if i, ok := number.FloatToInteger(f); ok {
                return i
            }
        }
        return key
    }

    put()方法向表里保存键值对。

    func (self *luaTable) put(key, val luaValue) {
        if key == nil {
            panic("table index is nil!")
        }
        if f, ok := key.(float64); ok && math.IsNaN(f) {
            panic("table index is NaN!")
        }

        key = _floatToInteger(key)
        if idx, ok := key.(int64); ok && idx >= 1 {
            arrLen := int64(len(self.arr))
            if idx <= arrLen {
                self.arr[idx-1] = val
                if idx == arrLen && val == nil {
                    self._shrinkArray()
                }
                return
            }
            if idx == arrLen+1 {
                delete(self._map, key)
                if val != nil {
                    self.arr = append(self.arr, val)
                    self._expandArray()
                }
                return
            }
        }
        if val != nil {
            if self._map == nil {
                self._map = make(map[luaValue]luaValue, 8)
            }
            self._map[key] = val
        } else {
            delete(self._map, key)
        }
    }

    先判断key是否是nil或者NaN，如果是则调用panic（），否则上市把key转换为整数。
    如果key是整数，或者已经被转换为整数，且在数组索引范围之内，直接按照索引修改数组元素就可以了。向数组中放入nil值会产生hole，如果hole在数组末尾的话，则调用_shrinkArray()把尾部的hole全部删除。如果key是整数，并且刚刚超出数组索引的范围，并且值不是nil，就把值追加到数组的末尾，然后调用_expandArray()动态扩展数组。

    如果值不是nil，就把键值对写到哈希表，否则把key从哈希表中删除，来节约空间。因为在穿件talbe的时候并不一定创建了哈希表部分，因此第一次写入的时候，需要创建哈希表。

    func (self *luaTable) _shrinkArray() {
        for i := len(self.arr) - 1; i >= 0; i-- {
            if self.arr[i] == nil {
                self.arr = self.arr[0:i]
            }
        }
    }


    func (self *luaTable) _expandArray() {
        for idx := int64(len(self.arr)) + 1; true; idx++ {
            if val, found := self._map[idx]; found {
                delete(self._map, idx)
                self.arr = append(self.arr, val)
            } else {
                break
            }
        }
    }

    _expandArray()在数组部分动态扩展以后，吧原本保存在哈希表中的某些值也挪到数组中。

    由于set（）方法已经数组部分和哈希表部分进行了动态的调整，因此len（）方法就容易实现了：

    func (self *luaTable) len() int {
        return len(self.arr)
    }

    修改lua_value.go中的代码typeOf(）：

    func typeOf(val luaValue) LuaType {
        switch val.(type) {
        ......
        case *luaTable:
            return LUA_TTABLE
        ......
        }
    }

    修改state/api_misc.go 中的len，让它能够获取数组的长度：

    func (self *luaState) Len(idx int) {
        val := self.stack.get(idx)

        if s, ok := val.(string); ok {
            self.stack.push(int64(len(s)))
        } else if t, ok := val.(*luaTable); ok {
            self.stack.push(int64(t.len()))
        } else {
            panic("length error!")
        }
    }


表相关api

因为表的实现完全属于lua解释器的内部细节，因此lua api并没有吧表直接暴露给用户，
而是提供了一系列创建、和操作表的方法。修改api/lua_state.go，给接口LuaState添加8个方法：

    /* get functions (Lua -> stack) */
    NewTable()
    CreateTable(nArr, nRec int)
    GetTable(idx int) LuaType
    GetField(idx int, k string) LuaType
    GetI(idx int, i int64) LuaType
    /* set functions (stack -> Lua) */
    SetTable(idx int)
    SetField(idx int, k string)
    SetI(idx int, i int64)

    get方法

    把get方法放在api_get.go文件中。state/api_get.go
    1、CreateTable（）

    CreateTable（）创建一个空的表，将其push到栈顶。该方法提供了两个参数，用来指定数组部分和哈希表部分的初始容量。如果可以预先估计出表的使用方式和容量，那么可以使用这两个参数在创建表的时候预先分配足够的空间，用来避免后续对表进行频繁的扩容。

    func (self *luaState) CreateTable(nArr, nRec int) {
        t := newLuaTable(nArr, nRec)
        self.stack.push(t)
    }

    2、NewTable()

    如果无法预先估计表的用法和容量，可以使用NewTable()创建表。NewTable()只是CreateTable（）的特殊情况。

    func (self *luaState) NewTable() {
        self.CreateTable(0, 0)
    }
































































































































































































































































































































































































































































































































































