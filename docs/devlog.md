



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



    4、
    5、
    6、
    7、
    8、
    9、



10、11、12、13、14、15、16、17、18、19、20、21、22、23、24、25、26、27、28、29、30、
















