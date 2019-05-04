

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

    3、GetTable（）

    GetTable（） 根据key（从栈顶弹出）从表（索引由参数指定）里面取值，然后把值push到栈顶，并返回值的类型。  

    func (self *luaState) GetTable(idx int) LuaType {
        t := self.stack.get(idx)
        k := self.stack.pop()
        return self.getTable(t, k)
    }

    为了减少重复，把根据key从table里面获取值的逻辑提取为函数getTable(t, k luaValue)：

    func (self *luaState) getTable(t, k luaValue) LuaType {
        if tbl, ok := t.(*luaTable); ok {
            v := tbl.get(k)
            self.stack.push(v)
            return typeOf(v)
        }

        panic("not a table!") // todo
    }

    如果位于指定索引处的值不是table，就暂时调用panic，以后在完善。

    4、GetField（）

    GetField（）和 GetTable（） 类似，只不过key不是从栈顶pop的任意值，而是用参数传入的字符串。GetField（）用来从记录中获取字段。

    func (self *luaState) GetField(idx int, k string) LuaType {
        t := self.stack.get(idx)
        return self.getTable(t, k)
    }

    还可以使用GetTable（）方法实现GetField（）方法：

      func (self *luaState) GetField(idx int, k string) LuaType {
        self.PushString(k)
        return self.getTable(idx)
    }  

    只不过第一种方法更加高效。

    5、GetI（）
    和GetField（）方法类似，只不过由参数传入的key是数字，而非字符串，这个方法是专门给数组准备的，用来根据索引获取数组的元素，执行后，相应的数组元素被push到栈顶。

    func (self *luaState) GetI(idx int, i int64) LuaType {
        t := self.stack.get(idx)
        return self.getTable(t, i)
    }  

Set方法
    
    state/api_set.go

    1、SetTable()

    SetTable()把键值对写入表。其中key和value从栈中弹出，表则位于指定的索引处。

    func (self *luaState) SetTable(idx int) {
        t := self.stack.get(idx)
        v := self.stack.pop()
        k := self.stack.pop()
        self.setTable(t, k, v)
    }

    同样，也把写表的逻辑提取成setTable(）方法：


    func (self *luaState) setTable(t, k, v luaValue) {
        if tbl, ok := t.(*luaTable); ok {
            tbl.put(k, v)
            return
        }

        panic("not a table!")
    }

    如果位于指定索引处的值不是表，就调用panic报错。

    2、SetField（）

    SetField（）和SetTable()类似，只不过key不是从栈顶弹出的任意值，而是由参数传入的字符串。用于给记录的字段赋值。执行后，value从栈顶弹出，并被赋值给记录的相应字段。

     func (self *luaState) SetField(idx int, k string) {
        t := self.stack.get(idx)
        v := self.stack.pop()
        self.setTable(t, k, v)
    }   

    3、SetI()

    SetI() 和SetField（）类似，只不过由参数传入的key是数组，而非字符串，用于按照索引修改数组的元素。
    执行之后，值从栈顶弹出，并被写到数组中。

    func (self *luaState) SetI(idx int, i int64) {
        t := self.stack.get(idx)
        v := self.stack.pop()
        self.setTable(t, i, v)
    }


table相关的指令

    newtable：创建空表。
    gettable根据key从表里面取值。
    settable根据key往表里面写入值。
    setlist按照索引批量更新数组元素。

    把它们放到文件vm/inst_table.go中：

    1、newtable
    newtable指令(iABC模式)创建空表，并将其放到指定的寄存器。寄存器索引由操作数A指定，表的初始数组容量和哈希表容量分别由操作数B、C指定。

    R（A）：= {} (size= B,C)

    lua代码中的每条表构造器语句都会产生一条newtable指令。该指令可以通过CreateTable（）函数实现：


    func newTable(i Instruction, vm LuaVM) {
        a, b, c := i.ABC()
        a += 1

        vm.CreateTable(Fb2int(b), Fb2int(c))
        vm.Replace(a)
    }


    y因为newtable指令是iABC模式，操作数BC只有9个bit，所以当作无符号整数的话，最大不能超过512。前面说过，因为表构造器方便实用，所以lua也经常被用来描述数据，比如json数据，如果有很大的数据需要写成表构造器，但是表的初始容量又不够大，就容易导致表被频繁扩容，从而影响数据的加载效率。

    为了解决这个问题，newtable指令的BC操作数使用了中叫做浮点字节的编码（floating point byte）方式。
    这种编码方式和浮点数的编码方式类似，只是仅仅用一个字节，具体来说，如果把某个字节用2进制写成eeeexxx，那么当eeee == 0时，该字节表示的整数就是xxx，否则该字节表示的整数就是（1xxx）*2^(eeee -1).

    lua官方实现中有现成的服点字节编码、解码函数。我们把它们拿过来，转换成go的函数，保存在vm/fpb.go中

    2、gettable（）   

    gettable（iABC模式）指令根据key从表中取值，并放到目标寄存器中。其中表位于寄存器中，索引有操作数B指定；key可能位于寄存器中，也可能在常量表中，索引由操作数C指定；目标寄存器的索引则由操作数A指定。

    R(A) := R(B)[RK(C)]

    该指令对应lua代码中的表索引取值操作。

    该这令可以借助GetTable（）实现：


    func getTable(i Instruction, vm LuaVM) {
        a, b, c := i.ABC()
        a += 1
        b += 1

        vm.GetRK(c)
        vm.GetTable(b)
        vm.Replace(a)
}

    settable
    settable指令（iABC 模式）根据key向表面面赋值。其中表位于寄存器中，索引由操作数A指定；key 、value柯恩呢该位于寄存器中，也可能在常量表中，索引分别由操作数BC指定。

    R(A)[RK(B)] := RK(C)

    该指令对应lua代码里面的表索引赋值操作，该指令可以借助SetTable函数实现：

    func setTable(i Instruction, vm LuaVM) {
        a, b, c := i.ABC()
        a += 1

        vm.GetRK(b)
        vm.GetRK(c)
        vm.SetTable(a)
    }


    setlist

    settable是通用指令，每次只处理一个键值对，具体操作交给表去处理，并不关心实际写入的是表的哈希部分，还是数组部分。setlist指令（iABC 模式），则是专门给数组准备的，用于按照索引批量设置数组元素，其中数组位于寄存器中，索引由操作数A指定；需要写入数组的若干个值也在寄存器中，紧挨着数组，数量由操作数B指定；数组起始索引则由操作数C指定。

    R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B

    数组的索引到底是怎么计算的？因为C只有9个bit，所以直接使用它来表示数组的索引显然是不够的。此处的解决的办法就是，让操作数C保存批次数，然后用批次数 乘上 批次大小（对应上面的fpf),就可以计算出数组的起始索引。默认的批次大小为50，操作数c能表示的最大索引就扩大到了25600（500 * 512）

    但是，如果数组的长度大于这个数值呢？是不是后面的元素就只能用settable指令设置了？这种情况下，setlist指令后面会跟着一条extraarg指令，用它的Ax操作数来保存批次数量。如果指令的操作数C大于0，那么表示的是批次数 +1，否则，整整的批次数量保存在后续的extraarg指令里面。

    func setList(i Instruction, vm LuaVM) {
        a, b, c := i.ABC()
        a += 1

        if c > 0 {
            c = c - 1
        } else {
            c = Instruction(vm.Fetch()).Ax()
        }

        vm.CheckStack(1)
        idx := int64(c * LFIELDS_PER_FLUSH)
        for j := 1; j <= b; j++ {
            idx++
            vm.PushValue(a + j)
            vm.SetI(a, idx)
        }
    }


    先从指令中解码出操作数，然后根据批次计算出数组的起始索引，最后循环调用PushValue（）以及SetI（）
    函数，按照索引设置数组元素。这里暂时值考虑操作数B>0的情况，操作数=0的情况以后再说。
    LFIELDS_PER_FLUSH 常量表示每个批次处理的数组元素数量，定义之：


    实现了这5条指令以后，还需要把它们注册到指令表中。vm/opcodes.go，修改opcodes的初始化代码。

    var opcodes = []opcode{
    /*     T  A    B       C     mode         name       action */
    // R(A) := R(B)[RK(C)]
    opcode{0, 1, OpArgR, OpArgK, IABC , "GETTABLE", getTable}, 
    // R(A)[RK(B)] := RK(C)
    opcode{0, 0, OpArgK, OpArgK, IABC , "SETTABLE", setTable}, 
    // R(A) := {} (size = B,C)
    opcode{0, 1, OpArgU, OpArgU, IABC , "NEWTABLE", newTable}, 
    // R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
    opcode{0, 0, OpArgU, OpArgU, IABC , "SETLIST ", setList},  
    
    .............

    }

单元测试

    $ luac -l test.lua 

    main <test.lua:0,0> (14 instructions at 0x7fca4dd00a90)
    0+ params, 6 slots, 1 upvalue, 2 locals, 9 constants, 0 functions
        1   [7] NEWTABLE    0 3 0
        2   [7] LOADK       1 -1    ; "a"
        3   [7] LOADK       2 -2    ; "b"
        4   [7] LOADK       3 -3    ; "c"
        5   [7] SETLIST     0 3 1   ; 1
        6   [8] SETTABLE    0 -4 -5 ; 2 "B"
        7   [9] SETTABLE    0 -6 -7 ; "foo" "Bar"
        8   [10]    GETTABLE    1 0 -8  ; 3
        9   [10]    GETTABLE    2 0 -4  ; 2
        10  [10]    GETTABLE    3 0 -9  ; 1
        11  [10]    GETTABLE    4 0 -6  ; "foo"
        12  [10]    LEN         5 0
        13  [10]    CONCAT      1 1 5
    14  [10]    RETURN      0 1


    $ go run lunago luac.out 
    [01] NEWTABLE [table][nil][nil][nil][nil][nil]
    [02] LOADK    [table]["a"][nil][nil][nil][nil]
    [03] LOADK    [table]["a"]["b"][nil][nil][nil]
    [04] LOADK    [table]["a"]["b"]["c"][nil][nil]
    [05] SETLIST  [table]["a"]["b"]["c"][nil][nil]
    [06] SETTABLE [table]["a"]["b"]["c"][nil][nil]
    [07] SETTABLE [table]["a"]["b"]["c"][nil][nil]
    [08] GETTABLE [table]["c"]["b"]["c"][nil][nil]
    [09] GETTABLE [table]["c"]["B"]["c"][nil][nil]
    [10] GETTABLE [table]["c"]["B"]["a"][nil][nil]
    [11] GETTABLE [table]["c"]["B"]["a"]["Bar"][nil]
    [12] LEN      [table]["c"]["B"]["a"]["Bar"][3]
    [13] CONCAT   [table]["cBaBar3"]["B"]["a"]["Bar"][3]

====================

函数调用

    介绍
    1、定义lua函数的时候，可以声明固定的参数列表，但调用函数的时候，并不一定要按照声明的列表来传递参数。本质上，固定参数就是在函数调用的时候，从外部获取初始值的局部变量，如果调用函数时，传入的参数数量多于定义函数时指定的数量，多于的传入参数会被忽略（或者被收到变长参数列表中）。反之，如果调用函数时传入的参数数量少于定义函数时指定的数量，缺失的参数会获得默认值nil。

    2、lua也支持变长参数列表vararg。vararg参数使用连续的三个点...表示，必须更在固定参数列表的后面，如果有的话。如果一个函数使用了vararg参数，称这个函数为vararg函数。在这种函数内部，可以使用vararg表达式，也是三个点...，来获取实际传入的vararg参数值。能够使用vararg表达式的有 赋值语句、函数调用语句、以及表构造器。

    注意：只有在vararg函数内部才能出现vararg表达式，否则函数就无法通过编译。

    3、函数可以返回因热数量的返回值。lua会在函数调用之后，根据情况对实际返回值的数量进行适当的调整。如果函数调用世界返回的值比期望的数量多，那么多余的返回值就会被丢弃；反之，如果函数调用实际返回的值比期望的数量少，孔雀的返回值会用nil补充，此规则适用于函数调用可以出现的任何地方，称之为多退少补原则。

    对于函数调用语句，函数的返回值会完全忽略。对于赋值语句，只有当函数调用表达式出现在语句末尾，并且没有被（）括起来的时候，函数的返回值数量才会根据变量的数量进行调整，否则，函数调用的返回值会被固定调整为1.

    4、如果函数调用表达式出现在参数列表、返回语句或者表构造器等等末尾，那么lua会把函数调用个的返回值一个不漏地收纳，然后原封不动地向外传递。

    函数调用栈

    函数调用也经常借用“栈”这种数据结构来实现。为了区别于lua栈，将其称之为函数调用栈，简称调用栈call stack。lua栈里面存放的是lua值，调用栈里面粗放的是调用栈帧，简称调用帧call frame。

    当调用一个函数的时候，要先往调用栈里面push一个调用帧，然后把参数传递给调用帧，函数依托调用帧执行指令，可能会调用其他的函数，以此类推。当函数执行完以后，调用帧会留下函数需要返回的值，把调用帧从调用栈顶弹出，并把返回值返回给底部的调用帧，这样一次函数调用就结束了。

    当前函数：当前正在执行的函数，其使用的调用帧为当前帧。
    主调函数：调用其他函数的函数，其使用的带哦用帧为主调帧
    被调函数：被其他函数调用的函数，其使用的调用帧为被调帧。

    主调和被调都是相对而言的，除了调用栈底部和顶部的函数，其他函数都同时充当着主调和被调两种角色。

    调用帧的实现

    之前，已经把虚拟寄存器这个函数执行过程中之冠重要的角色交给了Lua栈，lua栈也不负众望，表现的很好。这里，会对lua栈进行升级改造，让它可以在函数调用中担当起调用帧的角色。修改结构luaStack：

    type luaStack struct {
        slots   []luaValue  //用来存放值
        top     int         //记录栈顶的索引
        /* linked list */
        prev *luaStack
        /* call info */
        closure *closure
        varargs []luaValue
        pc      int

    }

    除了pc字段，其他三个字段的意义都不是很直观。之前，把pc和函数原型直接放在了luaState结构体中。由于这两个字段属于函数执行的内部状态，所以因该放在调用帧里面会更加合适。这里把pc字段照搬拿过来，但是函数原型被转换成了闭包closure。闭包又是什么？以后再说吧，这里把它当成函数原型的实例就可以了。在不引起混淆的情况下，以后出现的lua函数和闭包指的是同一个东西。

    在里面定义closure结构体：

    package state

    import "lunago/binchunk"

    type closure struct {
        proto *binchunk.Prototype
    }

    现在目前只有一个字段，存放函数原型，以后慢慢扩展。在里面定义用于创建lua闭包的函数：

    func newLuaClosure(proto *binchunk.Prototype) *closure {
        return &closure{proto: proto}
    }

    结构closure的实例或者指针对应lua语言的函数类型，
    state/lua_value.go修改typeOf函数，添加对函数类型的支持：

    func typeOf(val luaValue) LuaType {
        switch val.(type) {
        ......
        case *closure:
            return LUA_TFUNCTION
        ......
        }
    }


   varargs 字段用于实现变长参数机制。prev字段是让调用帧成为链表的节点。
   
   调用栈的实现

    结构体luaStack保存了函数的执行状态。清除luaState里面的无用数据。

    package state

    type luaState struct {
        stack *luaStack
    }

    func New() *luaState {
        return &luaState{
            stack: newLuaStack(20),
        }
    }

    这样结构体luaState就可以从充当调用栈了。增加操作栈的方法：

    func (self *luaState) pushLuaStack(stack *luaStack) {
        stack.prev = self.stack
        self.stack = stack
    }


    使用单向链表实现函数的调用栈，头部是栈顶，尾部是栈底。向栈顶push一个调用帧相当于在链表的头部插入一个节点，并让这个节点成为新的头部。

    func (self *luaState) popLuaStack() {
        stack := self.stack
        self.stack = stack.prev
        stack.prev = nil
    }

    从栈顶弹出一个调用帧只需要从链表的头部删除一个节点即可。

    luaState结构体 华丽变身为 函数调用栈。不过对它的改动破坏了之前的LuaVM接口部分的实现方法，
    需要修复。

    修改PC（）| AddPC（）|Fetch()|GetConst()|

    func (self *luaState) PC() int {
        return self.stack.pc
    }

    func (self *luaState) AddPC(n int) {
        self.stack.pc += n
    }

    func (self *luaState) Fetch() uint32 {
        i := self.stack.closure.proto.Code[self.stack.pc]
        self.stack.pc++
        return i
    }

    func (self *luaState) GetConst(idx int) {
        c := self.stack.closure.proto.Constants[idx]
        self.stack.push(c)
    }


 函数调用api

 lua解释器执行之前，需要先把脚本组装到一个main函数中，然后把main函数编译为函数原型，最后交给lua虚拟机执行。函数原型相当于类，作用就是实例化出真正可执行的函数，即闭包。


 lua api提供了load（）函数，用来将二进制chunk加载为闭包，并放到栈顶。至于闭包的执行，
 则由Call（）负责。

 实现load（）方法
 修改文件api/lua_state.go，给接口增加上面提到的方法Load（） \ Call()

    type LuaState interface {
    ...
        Load(chunk []byte, chunkName, mode string) int
        Call(nArgs, nResults int)
    ...
    }

    Load()

    Load()加载chunk，把main函数原型实例化为闭包并push到栈顶。如果加载的是chunk，只需要杜预文件、解析主函数原型、实例化为闭包、push到栈顶；如果加载的是lua脚本，先编译，在继续。

    Load()方法接手三个参数，第一个参数是字节数组，给出了要加载的chunk的数据，第二个参数是字符串，指定了chunk的名称，供加载错误或者调试使用；第三个参数是字符串，指定加载的模式，有“b”、“t”以及“bt”。如果加载的模式是b，那么第一个参数必须是二进制chunk数据，否则会加载失败，如果是“t”，那么第一个蚕食必须是文本chunk数据，否则会加载失败，如果加载模式是“bt”，那么第一个参数可以是二进制或者文本chunk数据，Load（）方法会根据世界的数据格式进行处理。暂时先忽略后两个参数，值加载二进制chunk数据，一步一步来，
    一口吃不了一个胖子，以后在完善。

    如果Load（）无法成功加载chunk，需要在栈顶留下一条错误信息。Load（）会返回一个状态码，0表示加载成功，非0表示加载失败，暂时也忽略状态码，先返回0.

    创建文件state/api_call.go

    package state

    import "fmt"
    import "luago/binchunk"
    import "luago/vm"

    func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
        proto := binchunk.Undump(chunk) // todo
        c := newLuaClosure(proto)
        self.stack.push(c)
        return 0
    }


    Call()函数

    Call会调用Lua函数。在执行Call之前，必须先把被调用的函数push到栈顶，然后把参数一次push到栈顶，Call（）完成后，参数值和函数会被弹出栈顶，取而代之的是指定数量的返回值。Call方法接收两个参数：第一个参数指定准备传递给被调函数的参数数量，同时也隐含给出了被调函数在栈中的位置；第二个参数指定需要的返回值的数量（多退少补），如果是-1，则被调函数的返回值会全部留在栈顶。
    具体实现为：

    func (self *luaState) Call(nArgs, nResults int) {
        val := self.stack.get(-(nArgs + 1))
        if c, ok := val.(*closure); ok {
            fmt.Printf("call %s<%d,%d>\n", c.proto.Source,
                c.proto.LineDefined, c.proto.LastLineDefined)
            self.callLuaClosure(nArgs, nResults, c)
        } else {
            panic("not a function!")
        }
    }

    先按照索引找到要调用的值，然后判断是不是lua函数，如果是打印调试信息，通过callLuaClosure（）调用该函数，否则调用panic报错。


    func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
        nRegs := int(c.proto.MaxStackSize)
        nParams := int(c.proto.NumParams)
        isVararg := c.proto.IsVararg == 1

        // create new lua stack
        newStack := newLuaStack(nRegs + 20)
        newStack.closure = c

        // pass args, pop func
        funcAndArgs := self.stack.popN(nArgs + 1)
        newStack.pushN(funcAndArgs[1:], nParams)
        newStack.top = nRegs
        if nArgs > nParams && isVararg {
            newStack.varargs = funcAndArgs[nParams+1:]
        }

        // run closure
        self.pushLuaStack(newStack)
        self.runLuaClosure()
        self.popLuaStack()

        // return results
        if nResults != 0 {
            results := newStack.popN(newStack.top - nRegs)
            self.stack.check(len(results))
            self.stack.pushN(results, nResults)
        }
    }

    先从函数原型拿到编译器准备好的各种信息：执行函数需要的寄存器数量、定义函数时声明的固定参数数量，以及是不是vararg函数。然后根据寄存器的数量（适当的扩大，因为要给指令水岸函数预留少量的栈空间）创建一个新的调用帧，并把闭包和调用帧联系起来。

    新的调用帧创建好以后，调用当前帧的PopN（）函数把函数和参数值一次性从栈顶弹出，然后调用pushN（）按照固定参数数量传入参数。固定参数传递完以后，需要修改新帧的栈顶指针，让它指向最后一个寄存器。如果被调函数数vararg函数，并且传入参数的数量多于固定参数数量，还需要把vararg参数记录下来，并保存到调用帧中。

    把心的调用帧push到调用栈顶，让其成为当前帧，然后调用runLuaClosure（）执行被调函数的指令。执行完毕后，心调用帧的使命就结束了，把它从调用栈顶弹出，这样主调帧就又成了当前帧。被调函数运行完以后，返回值会留在被调帧的栈顶---寄存器之上。需要把全部的返回值从被调帧的栈顶弹出，然后根据期望的返回值数量对退少补，push到当前帧的栈顶，这样函数调用才完事。

    func (self *luaState) runLuaClosure() {
        for {
            inst := vm.Instruction(self.Fetch())
            inst.Execute(self)
            if inst.Opcode() == vm.OP_RETURN {
                break
            }
        }
    }

    state/lua_stack.go添加函数

    popN(n int)从栈顶一次性弹出多个值。

    func (self *luaStack) popN(n int) []luaValue {
        vals := make([]luaValue, n)
        for i := n - 1; i >= 0; i-- {
            vals[i] = self.pop()
        }
        return vals
    }


    pushN()：向栈顶push多个值（多退少补）

    func (self *luaStack) pushN(vals []luaValue, n int) {
        nVals := len(vals)
        if n < 0 {
            n = nVals
        }

        for i := 0; i < n; i++ {
            if i < nVals {
                self.push(vals[i])
            } else {
                self.push(nil)
            }
        }
    }


    函数调用指令的实现
    closure\call\return\vararg\tailcall\self

    vm/inst_call.go实现上面的方法：


    closure指令（iBx模式）把当前lua函数的子函数原型实例化为闭包，放到由操作数A指定的寄存器中，子函数原型来自当前函数原型的子函数原型列表，索引由操作数Bx指定

    R(A) := closure(KPROTO[Bx])

    closure 指令对应lua脚本里面的函数定义语句或者表达式。


    func closure(i Instruction, vm LuaVM) {
        a, bx := i.ABx()
        a += 1

        vm.LoadProto(bx)
        vm.Replace(a)
    }

    由于lua api值提供了加载主函数原型的load（）方法，并没有提供加载子函数原型的方法，因此需要扩展LuaVM接口，给它添加一个方法LoadProto（）。该方法吧当前函数的子函数原型（索引由参数指定）实例化为闭包，并push到栈顶。

    Call
    CALL指令（iABC模式）调用Lua函数。其中被调函数位于寄存器中，索引由操作数A指定，需要传递给被调函数的参数值也要在寄存器中，紧挨着被调函数，数量由操作数B指定，函数调用结束后，原先存放在函数和参数值的寄存器会被返回值占据，具体由多少个返回值则由操作数C指定：

    R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))

    call指令对应lua脚本连的函数调用语句或者表达式。

    func call(i Instruction, vm LuaVM) {
        a, b, c := i.ABC()
        a += 1

        // println(":::"+ vm.StackToString())
        nArgs := _pushFuncAndArgs(a, b, vm)
        vm.Call(nArgs, c-1)
        _popResults(a, c, vm)
    }


    call指令可以额借助call方法实现。先调用_pushFuncAndArgs()函数把被调函数和参数值push到栈顶，然后让call（）方法去处理函数调用的逻辑，call（）方法结束之后，函数返回值已经在栈顶，调用_popResults()函数把这个些返回值移动到适当的寄存器中就可以了。

    func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
        if b >= 1 {
            vm.CheckStack(b)
            for i := a; i < a+b; i++ {
                vm.PushValue(i)
            }
            return b - 1
        } else {
            _fixStack(a, vm)
            return vm.GetTop() - vm.RegisterCount() - 1
        }
    }


    如果操作数B>0就简单了，需要传递的参数是B -1 个，循环调用PushValue（）方法把函数和参数值push到栈顶即可。由于我们给指令预留的栈顶空间是很少的，而传入参数的数量却不确定，所以这里需要调用CheckStack（）方法确保栈顶有足够的空间可以容纳函数和参数值。当B=0的时候，以后在说。

    func _popResults(a, c int, vm LuaVM) {
        if c == 1 {
            // no results
        } else if c > 1 {
            for i := a + c - 2; i >= a; i-- {
                vm.Replace(i)
            }
        } else {
            // leave results on stack
            vm.CheckStack(1)
            vm.PushInteger(int64(a))
        }
    }


    如果操作数C>1,则返回值数量是C-1，循环调用Replace（）方法把栈顶返回值移动到相应的寄存器就可以了；
    如果操作数C=1，则返回值数量是0，不需要任何处理；如果C=0，那么需要把被调函数的返回值全部返回。对于最后这种情况，干脆就把这些返回值线留在栈顶，反正后面也是要把它们在push到栈顶的。先在栈顶push一个整数值，用来标记这些返回值原本是要移动到哪些寄存器中。

    更新_pushFuncAndArgs

    } else {//参数 B 等于 0 的情况
        _fixStack(a, vm)
        return vm.GetTop() - vm.RegisterCount() - 1
    }

    因为return指令也会有类似的情况，所以需要把相应的逻辑取到函数_fixStack()里，RegisterCount()方法返回当前函数的寄存器数量，需要在LuaVM接口中定义。

    func _fixStack(a int, vm LuaVM) {
        x := int(vm.ToInteger(-1))
        vm.Pop(1)

        vm.CheckStack(x - a)
        for i := a; i < x; i++ {
            vm.PushValue(i)
        }
        vm.Rotate(vm.RegisterCount()+1, x-a)
    }

    因为后半部分草数值已经在栈顶了，所以值需要把函数和前半部分参数值push到栈顶，然后旋转栈顶即可。

    return

    return指令（iABC 参数）把存放爱连续多个寄存器里面的值返回给主调函数。其中第一个寄存器的索引由操作数A指定，寄存器的数量由操作数B指定，操作数C没有使用。

    return R(A), ... ,R(A+B-2)


    return指令对应lua脚本里面的return语句

    func _return(i Instruction, vm LuaVM) {
        a, b, _ := i.ABC()
        a += 1

        if b == 1 {
            // no return values
        } else if b > 1 {
            // b-1 return values
            vm.CheckStack(b - 1)
            for i := a; i <= a+b-2; i++ {
                vm.PushValue(i)
            }
        } else {
            _fixStack(a, vm)
        }
    }

    需要把返回值push到栈顶，如果操作数B=1，则不需要返回任何值；
    如果操作数B>1,则需要返回B -1 个值，这些值已经在寄存器里面了，循环调用PushValue()复制到栈顶即可。
    如果操作数B=0，则一部分返回值已经在栈顶了，调用_fixStack()函数吧另一部分也push到栈顶即可。

    vararg

    vararg指令(iABC 模式)把传递给当前函数的变长参数加载到连续的多个寄存器中，其中第一个寄存器的索引由
    操作数A指定，寄存器数量由操作数B指定，操作数C没有使用。

    R(A), R(A+1), ..., R(A+B-2) = vararg

    对应lua脚本中的vararg表达式，由于编译器生成的main函数也是vararg函数，所以也可以在里面使用vararg表达式。从效果来看，vararg表达式和函数调用很像。在vararg函数内部，凡事能用函数调用表达式的地方，也能用vararg表达式 。如果把vararg表达式当作函数调用，其返回值就是变长参数，要进行多退少补。正因为如此，就可以重复使用call指定的部分代码来实现vararg指令。

    func vararg(i Instruction, vm LuaVM) {
        a, b, _ := i.ABC()
        a += 1

        if b != 1 { // b==0 or b>1
            vm.LoadVararg(b - 1)
            _popResults(a, b, vm)
        }
    }


    如果操作数B>1,表示把 b -1个vararg参数复制到寄存器；否则只能=0，表示把全部vararg参数复制到寄存器。对于这两种情况，统一调用loadVararg()把vararg参数push到栈顶，剩下的工作交给_popResults()函数即可。loadVararg()方法需要在LuaVM接口中定义 todo

    tailcall

    函数调用一般通过调用栈来实现，使用这种方法，每调用一个函数都会产生一个调用帧。如果方法调用层次很深，特别是递归调用的时候，就容易导致调用栈的溢出。有什么办法，既能发挥递归调用的便捷，又能避免调用栈的溢出呢？？？？答案就是尾递归优化，利用这种技术，被调函数就可以重用主调函数的调用帧，因此可以有效缓解调用栈溢出的问题，不过尾递归优化仅仅适用于一些特定的情况，也不是很通用。

    return会被lua编译器编译为tailcall（iABC 模式）

     return R(A)(R(A+1), ... ,R(A+B-1))


    func tailCall(i Instruction, vm LuaVM) {
        a, b, _ := i.ABC()
        a += 1

        // todo: optimize tail call!
        c := 0
        nArgs := _pushFuncAndArgs(a, b, vm)
        vm.Call(nArgs, c-1)
        _popResults(a, c, vm)
    }

    self
    lua不是面向对象的语言，不过提供了一些语法和底层的支持，利用这些，可以够再出一套面向对象的体系。self
    指令主要用来优化方法调用语法糖。

    self指令（iABC模式），把对此昂和方法复制到相邻的两个目标寄存器中。对象在寄存器中，索引由操作数B指定。方法名在常量表中，索引由操作数C指定。目标寄存器的索引由操作数A指定

     R(A+1) := R(B); R(A) := R(B)[RK(C)]


    func self(i Instruction, vm LuaVM) {
        a, b, c := i.ABC()
        a += 1
        b += 1

        vm.Copy(b, a+1)
        vm.GetRK(c)
        vm.GetTable(b)
        vm.Replace(a)
    }

    注册上面的方法分到vm/opcodes.go 的 opcodes中

    扩展LuaVM接口

    api/lua_vm.go 扩展LuaVM接口，添加方法

    type LuaVM interface {
        ......
        RegisterCount() int     //当前lua函数所操作的寄存器计数器
        LoadVararg(n int)       //把传递给当前lua函数的变长参数push到栈顶 多退少补
        LoadProto(idx int)      //把当前lua函数的子函数的原型 实例化为闭包 ，并push到栈顶
    }

    在state/api_vm.go中实现上面的方法
    
    RegisterCount() int     //当前lua函数所操作的寄存器计数器   

    func (self *luaState) RegisterCount() int {
        return int(self.stack.closure.proto.MaxStackSize)
    }


    LoadVararg(n int)       //把传递给当前lua函数的变长参数push到栈顶 多退少补

    func (self *luaState) LoadVararg(n int) {
        if n < 0 {
            n = len(self.stack.varargs)
        }

        self.stack.check(n)
        self.stack.pushN(self.stack.varargs, n)
    }

    LoadProto(idx int)      //把当前lua函数的子函数的原型 实例化为闭包 ，并push到栈顶

    func (self *luaState) LoadProto(idx int) {
        proto := self.stack.closure.proto.Protos[idx]
        closure := newLuaClosure(proto)
        self.stack.push(closure)
    }


    setlist加强

    之前的 setlist，忽略了B=0的情况，当表构造器的最后一个元素是函数调用或者 vararg 表达式的时候，lua会把它们所产生的所有制都收集起来，供 setlist 使用

    修改setList()函数，处理操作数 B = 0 的情况
    vm/inst_table.go

    func setList(i Instruction, vm LuaVM) {

        ......
        bIsZero := b == 0
        if bIsZero {
            b = int(vm.ToInteger(-1)) - a - 1
            vm.Pop(1)
        }
        .....
    }

    记录下这种情况，然后适当的调整B操作数，先按照正常的逻辑处理寄存器中的值，寄存器处理完毕以后，在处理栈顶的值

  func setList(i Instruction, vm LuaVM) {

        ......
            if bIsZero {
            for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
                idx++
                vm.PushValue(j)
                vm.SetI(a, idx)
            }

            // clear stack
            vm.SetTop(vm.RegisterCount())
        }
        .....
    }
  
    都处理好了以后，调用SetTop让栈顶恢复原始状态。


单元测试

    之前，因为还没有办法加载以及调用lua函数，只能把main函数指令的执行过程编码在luaMain（）中。现在既然又了load和call函数，luamain就可以不使用了。

    修改main。go测试

    package main

    import "io/ioutil"
    import "os"
    import "luago/state"

    func main() {
        if len(os.Args) > 1 {
            data, err := ioutil.ReadFile(os.Args[1])
            if err != nil {
                panic(err)
            }

            ls := state.New()
            ls.Load(data, os.Args[1], "b")
            ls.Call(0, 0)
        }
    }


    4个步骤：
    1、读取chunk 2、创建luaState实例 
    3、调用Load（）把main函数加载到栈顶 
    4、调用Call（）运行nain函数，没有给它传递任何参数，也不需要什么返回值。


    go run main.go  luac.out 
    call @test.lua<0,0>
    call @test.lua<6,17>
    panic: GETTABUP

    to do fix

===================

go函数调用

仅仅使用lua函数能够做的事情有限，比如没有拌饭获取当前的时间，没有拌饭读取文件，也没有办法向控制台打印输出。如何在lua语言中调用go编写的函数？？？？

s虽然lua函数需要go函数弥补自身的不足，不过lua函数也是相当挑剔的，并不是任何go函数都能使用。在lua看来，go函数和lua函数都是function类型，没办法简单区分一个函数到底是lua函数还是go函数。go函数以后叫做go闭包。

 、添加go函数类型

 要想让lua函数调用go编写的函数，需要一种机制，能够给go函数传递参数，并接收go函数的返回值。不过lua只能操作lua栈，咋办？在执行lua函数的时候，lua栈厂荡虚拟寄存器，提供给指令操作。在调用lua函数的时候，lua栈充当栈帧，提供参数和返回值。因此，可以利用lua栈给go函数传递参数和接收返回值。

    go函数约定：接受一个LuaState接口类型的参数，返回一个整数。在go函数执行之前，lua栈里面是传入的参数值，没有别的。当go函数结束以后，吧需要返回的值留在栈顶，饭后返回一个整数表示返回值的个数。由于go函数返回了返回值的数量，因此在执行结束以后，就不用对栈进行清理了，把返回值留在栈顶就可以了。

    api/lua_state.go 添加

    type GoFunction func(LuaState) int

    state/closure.go添加

    import . "luago/api"

    然后给结构体closure添加goFunc字段

    type closure struct {
        proto *binchunk.Prototype   // lua closure
        goFunc GoFunction          // go closure
    }

    使用closure结构体表示lua以及go函数。如果proto不是nil，说明这是lua闭包。同理，goFunc。

    添加创建go闭包的函数

    func newGoClosure(f GoFunction) *closure {
        return &closure{goFunc: f}
    }

    扩展lua api

    go函数要进入lua栈，变成go闭包才能被lua所使用。lua api提供了PushGoFuncion(),作用就是把go函数转换为go闭包，并放到栈顶。另一方面，lua api提供了IsGoFunction 和 ToGoFunction函数，可以把栈里面的go闭包在转换为go函数返回给用户。

    api/lua_state.go，给LuaState添加3方法

    type LuaState interface {
        .....
        PushGoFunction(f GoFunction)
        IsGoFunction(idx int) bool
        ToGoFunction(idx int) GoFunction
    }

    1、PushGoFunction(f GoFunction)

    接受一个go函数参数，把它转换为go闭包，然后push到栈顶。
    state/api_push.go实现之

    func (self *luaState) PushGoFunction(f GoFunction) {
        self.stack.push(newGoClosure(f))
    }


    2、IsGoFunction(idx int) bool

    判断指定索引处的值是否可以转换为go函数。该方法以栈的索引为参数，返回布尔值，不会改变栈额状态。

    state/api_access.go实现之

    func (self *luaState) IsGoFunction(idx int) bool {
        val := self.stack.get(idx)
        if c, ok := val.(*closure); ok {
            return c.goFunc != nil
        }
        return false
    }

    先根据索引拿到值，然后看它是否是闭包，如果是，进一步看它是不是go闭包，只有go闭包菜可以转换为go函数。

    3、ToGoFunction(idx int) GoFunction

    把指定索引处的值转换为go函数并返回。如果值无法转换为go函数，返回nil。该方法以栈索引为参数，不改变栈状态。

    state/api_access.go

    func (self *luaState) ToGoFunction(idx int) GoFunction {
        val := self.stack.get(idx)
        if c, ok := val.(*closure); ok {
            return c.goFunc
        }
        return nil
    }


    先根据索引拿到值，看是否是闭包，如果是，直接返回goFunc字段的值（对于go
    闭包，该自动断就是期望的返回值，对于lua闭包，该字段为nil），否则返回nil。

  调用go函数
  state/api_call.go 
  修改里面的Call

  func (self *luaState) Call(nArgs, nResults int) {
        val := self.stack.get(-(nArgs + 1))
        if c, ok := val.(*closure); ok {
            if c.proto != nil {
                self.callLuaClosure(nArgs, nResults, c)
            } else {
                self.callGoClosure(nArgs, nResults, c)
            }
        } else {
            panic("not function!")
        }
    }

    主要改动在外层if中，根据proto判断要调用的是lua闭包还是go闭包，如果是lua闭包，交给allLuaClosure处理；如果是go闭包，交给callGoClosure处理。

    func (self *luaState) callGoClosure(nArgs, nResults int, c *closure) {
        // create new lua stack
        newStack := newLuaStack(nArgs+LUA_MINSTACK, self)
        newStack.closure = c

        // pass args, pop func
        if nArgs > 0 {
            args := self.stack.popN(nArgs)
            newStack.pushN(args, nArgs)
        }
        self.stack.pop()

        // run closure
        self.pushLuaStack(newStack)
        r := c.goFunc(self)
        self.popLuaStack()

        // return results
        if nResults != 0 {
            results := newStack.popN(r)
            self.stack.check(len(results))
            self.stack.pushN(results, nResults)
        }
    }

    先创建心的调用帧，然后把参数值从主调帧弹出，push到被调帧。go闭包直接从主调帧里面弹出，扔掉即可。参数传递完毕以后，把被调帧push到条用栈，让它成为当前帧，然后直接执行go函数。执行完毕以后，把被调帧从调用栈pop，这样主调帧就有成了当前帧。最后（如果有必要），还需要把返回值从被调帧里面弹出，push到主调帧（多退少补）。


    Lua注册表

    lua给用户提供了一个注册表，这个注册表实际上就是一个普通的lua表，所以用户可以在里面存放任何lua值。
    lua本身也用到了，lua的全局变量就是借助这个注册表实现的。

    添加注册表

    state/lua_state.go添加内容

    import . "lunago/api"

    type luaState struct {
        registry *luaTable//added
        stack *luaStack
    }

    由于注册表是全局的，每个lua解释器实例都有自己的注册表，因此将其放在lusState结构体里比较合适。需要在创建luaState实例的时候初始化注册表，修改New函数

    func New() *luaState {
        registry := newLuaTable(0, 0)
        registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))

        ls := &luaState{registry: registry}
        ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
        return ls
    }

    先创建注册表，然后预先在里面放置一个全局环境。全局环境也是一个普通的lua表，所有的lua全局变量都放在这里。

    创建好注册表以后，用它创建luaState结构体实例，然后往里面push一个空的lua栈（调用帧），之前，在创建lua栈的时候，吧预留的空间写成了20，这是一种很不好的做法，因此引入了一个常量LUA_MINSTACK.

    加入其他的常量
    api/consts.go

    const LUA_MINSTACK = 20
    const LUAI_MAXSTACK = 1000000
    const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000
    const LUA_RIDX_GLOBALS int64 = 2

    LUA_RIDX_GLOBALS用于定义全局环境在注册表里面的索引。

    操作注册表

    由于注册表就是一个普通的lua表，所以lua api并没有提供专门的方法来操作。任何可以操作表的api都可以用来操作注册表，不过，表操作的方法都是通过索引来访问表的，可以通过“伪索引”来访问注册表。

    一般来说，我们并不需要lua栈有非常大的容量，随意lua定义了一个常量，用来表示lua栈的最大索引，这个常量就是LUAI_MAXSTACK。负的一百万-1000就是表示注册表的伪索引，用LUA_REGISTRYINDEX表示。

    为了支持注册表的伪索引，需要从luaStack里面访问注册表。state/lua_stack.go
    添加luastack结构体字段

    type luaStack struct {
        ......
        state   *luaState
        ......

    }


    让luaStack引用luaState，这样就可以间接访问注册表了。
    修改newluaStack，给state赋值

    func newLuaStack(size int, state *luaState) *luaStack {
        return &luaStack{
            slots: make([]luaValue, size),
            top:   0,
            state: state,
        }
    }

    继续修改luaStack结构体的其他方法，让它们支持注册表伪索引。

    func (self *luaStack) absIndex(idx int) int {
        if idx >= 0 || idx <= LUA_REGISTRYINDEX {
            return idx
        }
        return idx + self.top + 1
    }


    如果索引<LUA_REGISTRYINDEX,说明是伪索引，直接返回即可。

    func (self *luaStack) isValid(idx int) bool {
        if idx == LUA_REGISTRYINDEX {
            return true
        }
        absIdx := self.absIndex(idx)
        return absIdx > 0 && absIdx <= self.top
    }


    注册表伪索引属于有效索引，所以直接返回true。

    func (self *luaStack) get(idx int) luaValue {
        if idx == LUA_REGISTRYINDEX {
            return self.state.registry
        }

        absIdx := self.absIndex(idx)
        if absIdx > 0 && absIdx <= self.top {
            return self.slots[absIdx-1]
        }
        return nil
    }

    如果索引是注册表伪索引，直接返回注册表。

    func (self *luaStack) set(idx int, val luaValue) {
        if idx == LUA_REGISTRYINDEX {
            self.state.registry = val.(*luaTable)
            return
        }

        absIdx := self.absIndex(idx)
        if absIdx > 0 && absIdx <= self.top {
            self.slots[absIdx-1] = val
            return
        }
        panic("invalid index!")
    }

    如果索引是注册表伪索引，直接修改注册表。这里并没有检查传入的值是不是真的lua表，所以如果传入的是其他类型的值可能会导致注册表变成nil。

    全局环境

    lua全局变量是保存在全局环境里面，而全局环境也只是保存在注册表中的一个普通表。扩展lua api，增加对全局环境的操作。

    使用api操作全局环境

    lua api提供了四个方法，用来操作全局环境。api/lua_state.go，给接口LuaState增加方法

    type LuaState interface {
        ......

        PushGlobalTable()
        GetGlobal(name string) LuaType
        SetGlobal(name string)
        Register(name string, f GoFunction)
    }

    1、PushGlobalTable()

    由于全局环境也只是一个普通的表。所以gettable和settable等表的方法也同样适用，不过要使用这些方法，必须吧全局环境事先放到栈中。1、PushGlobalTable就是用来做这件事的，把全局环境push到栈顶，以便以后使用。
    state/api_push.go

    func (self *luaState) PushGlobalTable() {
        global := self.registry.get(LUA_RIDX_GLOBALS)
        self.stack.push(global)
    }

    还可以痛殴注册表伪索引来访问注册表，从中提取全局环境，比如

    func (self *luaState) PushGlobalTable(){
        self.GetI(LUA_REGISTRYINDEX, LUA_RIDX_GLOBALS)
    }


    2、GetGlobal(name string) LuaType

    由于全局环境主要是用来实现Lua全局变量的，所以里面的key基本上都是字符串。全局环境主要是当作记录来使用。为了便于操作，lua pai提供了GetGlobal，可以把全局环境中的某个字段(名字由参数指定)push栈顶。
    state/api_get.go

    func (self *luaState) GetGlobal(name string) LuaType {
        t := self.registry.get(LUA_RIDX_GLOBALS)
        return self.getTable(t, name)
    }

    也可以使用PushGlobalTable()以及GetField（）实现

     func (self *luaState) GetGlobal(name string) LuaType {
        self.PushGlobalTable()
        return self.GetField(t, name)
    }

     3、SetGlobal(name string)

     往全局环境里面写入一个值，其中的字段名由参数指定，值从栈顶弹出。
    state/api_set.go

    func (self *luaState) SetGlobal(name string) {
        t := self.registry.get(LUA_RIDX_GLOBALS)
        v := self.stack.pop()
        self.setTable(t, name, v)
    }

     4、Register(name string, f GoFunction)

     用于给全局环境注册go函数值，仅仅用于操作全局环境，字段名以及go函数从参数传入，不改变lua栈的状态。

     func (self *luaState) Register(name string, f GoFunction) {
        self.PushGoFunction(f)
        self.SetGlobal(name)
    }


    在Lua里面访问全局环境
    如何在lua函数里面访问全局变量？

    临时方法
    vm/inst_upvalue.go，实现gettabup指令

    package vm

    import . "lunago/api"

    // R(A) := UpValue[B][RK(C)]
    func getTabUp(i Instruction, vm LuaVM) {
        a, _, c := i.ABC()
        a += 1

        vm.PushGlobalTable()
        vm.GetRK(c)
        vm.GetTable(-2)
        vm.Replace(a)
        vm.Pop(1)
    }

    姑且认为这个指令可以额把某个全局变量放到指定的寄存器中即可。

    编辑opcodes.go文件，修改指令表，把gettabup指令的实现方法注册进去。
    vm/opcodes.go

    添加
    var opcodes = []opcode{
        ....
    // R(A) := UpValue[B][RK(C)]
    opcode{0, 1, OpArgU, OpArgK, IABC /* */, "GETTABUP", getTabUp}, 
        ...
    }

单元测试

增加了对供函数的支持，是时候写一个go函数了。

    添加两个import

    import "fmt"
    import . "luago/api"

    定义函数print

    func print(ls LuaState) int {
        nArgs := ls.GetTop()
        for i := 1; i <= nArgs; i++ {
            if ls.IsBoolean(i) {
                fmt.Printf("%t", ls.ToBoolean(i))
            } else if ls.IsString(i) {
                fmt.Print(ls.ToString(i))
            } else {
                fmt.Print(ls.TypeName(ls.Type(i)))
            }
            if i < nArgs {
                fmt.Print("\t")
            }
        }
        fmt.Println()
        return 0
    }


    由于Lua的print（）函数可以接收任意数量的参数，因此要做的第一件事情就是调用GetTop方法，看看栈里面有多少个值，换句话说，就是看看go函数接收到了多少个参数，接着取出这些参数，然后根据类型打印到控制台，饼子啊每两个值之间打印一个tab。lua的print不会返回任何值，所以打印完信息就没事了，也不用往栈里面push任何值，直接返回0就可以了。

    将go的print注册到lua里面：

    添加一行
    func main() {
        ....
        ls.Register("print", print)
        ....

    }

    创建好luaState实例以后，调用Register方法把刚才实现的print注册到全局环境。

    $ go run main.go luac.out 
    hello world！！！

================

闭包以及Upvalue

介绍
    闭包是计算机编程语言中非常普遍的一个概念，不过Upvalue却是lua独有的。

背景知识
    1、一等函数
        如果在一门编程语言中，函数属于一等公民，就说这门语言里的函数是一等函数。其实就是说，函数用起来和其他诸如数字和字符串类型的值没什么区别，比如说可以把函数存储在数据结构中，赋值给变量、作为参数传递给其他函数或者作为返回值从其他函数里面返回等等。

        如果一个函数以其他函数为参数，或者返回其他函数，称这个函数为高阶函数；反之，称这个函数为一阶函数。

        如果可以在函数内部定义其他函数，称内部定义的函数为嵌套函数，外部的函数为外围函数。在许多支持一等函数的语言里，函数实际上都是匿名的，在这些语言中，函数名就是普通的变量名，知识变量的值恰好是函数而已。比如在lua中，函数定义语句知识函数定义表达式（函数构造器，函数字面量）以及赋值语句的语法糖。下面的两条语句等价
        function add(x, y) return x + y end
        add = function(x, y) return x + y end

    2、变量的作用域
        支持一等函数的语言有个一非常重要的问题需要处理，那就是非局部变量的作用域问题，最简单的做法就是不支持嵌套函数，比如c，这样就完全避开了这个问题。对于支持嵌套函数的语言，有2种处理方式：
        动态作用域  以及 静态作用域    

        在使用动态作用域的语言中，函数的非局部变量名具体绑定的是哪个变量，只有等到函数运行的时候才能确定。由于动态作用域会导致程序不容易理解，所以现代的编程语言大多数都采用了静态作用域。

        与动态作用域不同，静态作用域在编译时就可以确定非局部变量名称绑定的变量，因此静态作用域也叫做词法作用域。lua采用的是静态作用域

     3、闭包   

        所谓闭包，就是按照词法作用域，捕获了非局部变量的嵌套函数。因为lua函数本质上都是闭包，就算编译器生成的main函数也不例外，它从外部捕获了_ENV变量。

    Upvalue介绍

    实际上，Upvalue就是闭包内部捕获的非局部变量。

    lua编译器会把Upvalue的相关信息编译到函数原型中，存放在Upvalue表里面。函数原型的Upvalue表的每一项都有4列：第一列是序号，从0递增；第二列给出Upvalue的名称；第三列支出Upvalue捕获的是不是直接外围函数的局部变量，1是，0否；如果Upvalue捕获的是直接外围函数的局部变量，第四列给出局部变量在外围函数调用帧里面的索引。

    如果闭包捕获的是非直接外围函数的局部变量呢？那就是层层捕获，即使没有用到。

    lua这种，需要借助外围函数来捕获更外围函数局部变量的闭包，叫做扁平闭包。


    全局变量

    Upvalue是非局部变量，也就是说，就是某个外围函数中定义的局部变量。那么全局变量又是啥？

    全局变量实际上是某个特殊表的字段，而这个特殊的表正是之前实现的全局环境。lua编译器在生成main函数的时候，会在他的外围，隐式声明一个局部变量。

    类似：
    loacl _ENV
    function main(...)
    ...
    end

    然后编译器会把全局变量的读写，翻译为字段_ENV的读写，也就是说，全局变量实际上也是语法糖，去掉语法糖以后，大致类似以下形式：

    local function f（）
        local function()
            _ENV.x = _ENV.y
        end
     end   

     至于这个_ENV如何初始化，则是lua api的工作。

     lua的变量可以分为三类：
     局部变来那个在函数内部定义（本质上就是函数调用帧里面的寄存器）
     Upvalue是直接或者间接外围函数定义的局部变量
     全局变量则是全局环境表的字段（痛殴隐藏的Upvalue，也就是_ENV进行访问）

  Upvalue的底层支持   

  修改closure结构体

  闭包要捕获外围函数的局部变量，就必须有地方来存放这些变量。

  state/closure.go添加字段  upvals []*upvalue

      type closure struct {
        proto  *binchunk.Prototype // lua closure
        goFunc GoFunction          // go closure
        upvals []*upvalue
    }

    需要注意的是，对于某个Upvalue来说，对它的任何改动，都必须反映在其他该Upvalue可见的地方。另外，当嵌套函数执行的时候，外围函数的局部变量有可能已经退出作用域了。为了应对这种情况，需要增加一个间接层，使用Upvalue结构体来封装Upvalue。

    在closure。go文件中定义这个结构体

    type upvalue struct {
        val *luaValue
    }

    修改newLuaClosure函数

    func newLuaClosure(proto *binchunk.Prototype) *closure {
        c := &closure{proto: proto}
        if nUpvals := len(proto.Upvalues); nUpvals > 0 {
            c.upvals = make([]*upvalue, nUpvals)
        }
        return c
    }

    lua闭包捕获的Upvalue数量以ing由编译器计算好了，在创建Lua闭包的时候，预先分配好空间即可。初始化Upvalue则由Lua API负责。不仅lua闭包可以捕获Upvalue，go闭包也可以捕获Upvalue。与lua闭包不同的是，需要在创建go闭包的时候，明确指定Upvalue的数量。

    修改newGoClosure（）函数

    func newGoClosure(f GoFunction, nUpvals int) *closure {
        c := &closure{goFunc: f}
        if nUpvals > 0 {
            c.upvals = make([]*upvalue, nUpvals)
        }
        return c
    }

    由于给newGoClosure函数增加了一个参数，所以破坏了PushGoFunction（）的实现代码，因此修改这个方法，api_push.go

    func (self *luaState) PushGoFunction(f GoFunction) {
        self.stack.push(newGoClosure(f, 0)) //第二个参数传入0
    }


    修改api，在创建闭包的时候初始化Upvalue

    Lua 闭包的支持

    Lua函数都是闭包，就连编译器生成的main函数也是闭包，捕获了_ENV这个特殊的Upvalue，这个特殊的Upvalue的初始化是由Api Load（）负责的。具体而言，就是Load（）方法在加载闭包的时候，会看它是否需要Upvalue，如果需要，那么第一个Upvalue(对于main函数来说就是_ENV)会被初始化为全局环境，其他的Upvalue会被初始化为nil。

    修改Load方法

    state/api_call.go 

    添加Upvalue的初始化代码

    func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
        proto := binchunk.Undump(chunk) // todo
        c := newLuaClosure(proto)
        self.stack.push(c)
        if len(proto.Upvalues) > 0 {    //set _ENV
            env := self.registry.get(LUA_RIDX_GLOBALS)
            c.upvals[0] = &upvalue{&env}
        }
        return 0
    }

    因为nil的初始值已经是nil，所以我们只要把第一个Upvalue值设置为全局环境即可，上面的是main函数原型，在加载子函数原型的时候也需要初始化Upvalue。

    state/api_vm.go修改LoadProto（）方法

    func (self *luaState) LoadProto(idx int) {
        stack := self.stack
        subProto := stack.closure.proto.Protos[idx]
        closure := newLuaClosure(subProto)
        stack.push(closure)

        for i, uvInfo := range subProto.Upvalues {
            uvIdx := int(uvInfo.Idx)
            if uvInfo.Instack == 1 {
                if stack.openuvs == nil {
                    stack.openuvs = map[int]*upvalue{}
                }

                if openuv, found := stack.openuvs[uvIdx]; found {
                    closure.upvals[i] = openuv
                } else {
                    closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
                    stack.openuvs[uvIdx] = closure.upvals[i]
                }
            } else {
                closure.upvals[i] = stack.closure.upvals[uvIdx]
            }
        }
    }

    需要根据函数原型里面的Upvalue表来初始化闭包的Upvalue值，对于每个Upvalue，又有两种情况需要考虑：如果某个Upvalue捕获的是当前函数的局部变量（Instack == 1），那么只要访问当前函数的局部变量就可以了；如果某个Upvalue捕获的是更外围的函数中的局部变量（Instack == 0），该Upvalue已经被当前函数所捕获，只要把该Upvalue传递给闭包就可以了。

    对于第一种情况，如果Upvalue捕获的外围函数局部变量还在栈上，直接引用即可，称这种Upvalue处于开放open状态；反之，必须把变量的实际值保存在其他的地方，称这种Upvalue处于闭合closed状态。为了能够在合适的时机（比如局部变量退出作用域的时候）把处于开放open状态的Upvalue闭合，需要记录所有暂时还处于开放状态的Upvalue，把这些Upvalue记录在被捕获局部变量所在的栈帧中。

    state/lua_stack.go，给luaStack结构体添加openuvs字段，该字段是map类型，key是int类型，存放局部变量的寄存器索引，value是Upvalue指针

    type luaStack struct {
        。。。
    openuvs map[int]*upvalue
        。。。
    }

    go闭包支持

    不仅仅是lua函数，go函数也可以捕获Upvalue，
    api/lua_state.go，给接口LuaState添加PushGoClosure（）方法
    。。。。
    PushGoClosure(f GoFunction, n int)
    。。。

    这个函数和PushGoFunction差不多，把go函数转换成go闭包push到栈顶，区别是PushGoClosure先从栈顶弹出n个lua值，这些值会成为go闭包的Upvalue。


    state/api_push.go 实现PushGoClosure

    func (self *luaState) PushGoClosure(f GoFunction, n int) {
        closure := newGoClosure(f, n)
        for i := n; i > 0; i-- {
            val := self.stack.pop()
            closure.upvals[i-1] = &upvalue{&val}
        }
        self.stack.push(closure)
    }

    先创建go闭包，然后从栈顶弹出指定数量的值，让它们变成闭包的Upvalue，最后把闭包push到栈顶。go闭包可以携带Upvalue这没什么问题，可问题是如果访问这些Upvalue？？？lua api并没有提供专门的方法，而是像注册表那样，提供了伪索引。和lua栈索引一样，Upvalue索引也是从1开始递增的。对于任何一个Upvalue索引，用注册表伪索引减去该索引就可以得到对应的Upvalue伪索引，为了便于使用，把这个转换过程分装在函数
    LuaUpvalueIndex（）里面。

    api/lua_state.go增加这个函数LuaUpvalueIndex（）

    func LuaUpvalueIndex(i int) int {
        return LUA_REGISTRYINDEX - i
    }

    就像注册表伪索引一样，还需要修改luaStack结构体的isValid（）、get（）、set（）方法，让它们支持Upvalue伪索引。

    state/lua_stack.go，修改方法isValid

    func (self *luaStack) isValid(idx int) bool {
        if idx < LUA_REGISTRYINDEX { /* upvalues */
            uvIdx := LUA_REGISTRYINDEX - idx - 1
            c := self.closure
            return c != nil && uvIdx < len(c.upvals)
        }
        if idx == LUA_REGISTRYINDEX {
            return true
        }
        absIdx := self.absIndex(idx)
        return absIdx > 0 && absIdx <= self.top
    }


    如果索引<注册表索引，说明是Upvalue索引，把它转成真实的索引（从0开始），然后看看它是不是在有效范围内，

    修改get方法

    func (self *luaStack) get(idx int) luaValue {
        if idx < LUA_REGISTRYINDEX { /* upvalues */
            uvIdx := LUA_REGISTRYINDEX - idx - 1
            c := self.closure
            if c == nil || uvIdx >= len(c.upvals) {
                return nil
            }
            return *(c.upvals[uvIdx].val)
        }

        if idx == LUA_REGISTRYINDEX {
            return self.state.registry
        }

        absIdx := self.absIndex(idx)
        if absIdx > 0 && absIdx <= self.top {
            return self.slots[absIdx-1]
        }
        return nil
    }

    如果伪索引无效，直接返回nil，否则返回Upvalue值。修改set方法

    func (self *luaStack) set(idx int, val luaValue) {
        if idx < LUA_REGISTRYINDEX { /* upvalues */
            uvIdx := LUA_REGISTRYINDEX - idx - 1
            c := self.closure
            if c != nil && uvIdx < len(c.upvals) {
                *(c.upvals[uvIdx].val) = val
            }
            return
        }

        if idx == LUA_REGISTRYINDEX {
            self.state.registry = val.(*luaTable)
            return
        }

        absIdx := self.absIndex(idx)
        if absIdx > 0 && absIdx <= self.top {
            self.slots[absIdx-1] = val
            return
        }
        panic("invalid index!")
    }

    如果伪索引有效，就修改Upvalue值，否则直接返回。

Upvalue相关的指令

    有5个：getupval、setupval、gettabup、settabup、jmp。

    1、getupval

    getupval(iABC 模式)，把当前闭包的某个Upvale值复制到目标寄存器，其中目标寄存器的索引由操作数A指定，Upvalue索引由操作数B指定，操作数C没有使用。

     R(A) := UpValue[B]

     如果在函数中访问Upvalue值，lua表一起就会在这些地方生成

     1、getupval指令
     vm/inst_upvalue.go

    func getUpval(i Instruction, vm LuaVM) {
        a, b, _ := i.ABC()
        a += 1
        b += 1

        vm.Copy(LuaUpvalueIndex(b), a)
    }


    由于可以使用伪索引，所以getupval可回忆直接使用copy实现，不过需要注意的是，在lua虚拟机指令的操作数里，Upvalue的索引也是从0开始的，但是在转换成lua栈伪索引时，Upvalue指令从1开始的。

    2、setupval

    setupval指令（iABC），使用寄存器中的值给当前闭包的Upvalue赋值。 其中仅存起索引由操作数A指定，Upvalue索引由操作数B指定，操作数C没有用到。

     UpValue[B] := R(A)
    
     如果在函数给Upvalue赋值，lua编译器就会在这些地方生成setupval指令，


    func setUpval(i Instruction, vm LuaVM) {
        a, b, _ := i.ABC()
        a += 1
        b += 1

        vm.Copy(a, LuaUpvalueIndex(b))
    }


    3、gettabup

    如果当前闭包的某个Upvalue是表，则gettabup指令（iABC模式）可以根据key从该表里面取值，然后把value放到目标寄存器中。其中目标寄存器的索引由餐做数A指定，Upvalue的索引由操作数B指定，key（可能在寄存器中，也可能在常量表中）索引由操作数C指定。gettabup相当于getupval 和 gettable两条指令的组合

     R(A) := UpValue[B][RK(C)]

     如果在函数里面按照key从Upvalue表中取值，lua编译器就会在这些地方生成3、gettabup指令

    func getTabUp(i Instruction, vm LuaVM) {
        a, b, c := i.ABC()
        a += 1
        b += 1

        vm.GetRK(c)
        vm.GetTable(LuaUpvalueIndex(b))
        vm.Replace(a)
    }

    先调用GetRK()把keypush到栈顶，然后调用GetTalbe（）从Upvalue中取值（key从栈顶弹出，value被push到栈顶），然后调用Replace（）吧值从栈顶弹出，并放到目标寄存器。


    settabup

    如果当前闭包的某个Upvalue是表，则settabup指令（iABC 模式）可以根据key向这个表里面写入值，其中
    Upvalue的索引由操作数A指定，key 、value可控在寄存器中，也肯跟在常量表中，索引分别操作数B、C指定。

     UpValue[A][RK(B)] := RK(C)
    
     如果在函数里面根据key向Upvalue里面写入值，lua编译器就会在这些地方生成 settabup 指令  

    func setTabUp(i Instruction, vm LuaVM) {
        a, b, c := i.ABC()
        a += 1

        vm.GetRK(b)
        vm.GetRK(c)
        vm.SetTable(LuaUpvalueIndex(a))
    }

    先通过GetRK（）把key value push到栈顶，然后调用SetTable（）就可以了SetTalbe会把键值对从栈顶弹出，然后根据伪索引吧键值对写道Upvalue。

    jmp

    jmp指令除了可以进行无条件跳转之外，还兼顾着闭合处于开启状态的Upvalue的责任。如果某个快内部定义的局部变量已经被嵌套函数所捕获，那么当这些局部变量退出作用域（也就是结束的）的时候，编译器就会生成一条jmp指令，指示虚拟机闭合相应的Upvalue

    如果jmp指令的sBx操作数是0，所以其实并没有起到任何跳转的作用，真正的用途就是闭合Upvalue。

    vm/inst_misc.go，修改jmp指令的代码

    func jmp(i Instruction, vm LuaVM) {
        a, sBx := i.AsBx()

        vm.AddPC(sBx)
        if a != 0 {
            vm.CloseUpvalues(a)
        }
    }


    由于lua api并没有提供闭合Upvalue的方法，所以要自己添加，api/lua_vm.go，在接口LuaVM中添加方法
    CloseUpvalues（）方法

    type LuaVM interface {
        ....
        CloseUpvalues(a int)
        ....
    }

    state/api_vm.go 实现方法CloseUpvalues

    func (self *luaState) CloseUpvalues(a int) {
        for i, openuv := range self.stack.openuvs {
            if i >= a-1 {
                val := *openuv.val
                openuv.val = &val
                delete(self.stack.openuvs, i)
            }
        }
    }


    处于开始状态的Upvalue引用了还在寄存器中的值，把这些lua值从寄存器里面复制出来，然后更新Upvalue，这样就将其更改为闭合状态了。

    var opcodes = []opcode{
        /*     T  A    B       C     mode         name       action */
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "MOVE    ", move},     // R(A) := R(B)
        opcode{0, 1, OpArgK, OpArgN, IABx /* */, "LOADK   ", loadK},    // R(A) := Kst(Bx)
        opcode{0, 1, OpArgN, OpArgN, IABx /* */, "LOADKX  ", loadKx},   // R(A) := Kst(extra arg)
        opcode{0, 1, OpArgU, OpArgU, IABC /* */, "LOADBOOL", loadBool}, // R(A) := (bool)B; if (C) pc++
        opcode{0, 1, OpArgU, OpArgN, IABC /* */, "LOADNIL ", loadNil},  // R(A), R(A+1), ..., R(A+B) := nil
        opcode{0, 1, OpArgU, OpArgN, IABC /* */, "GETUPVAL", getUpval}, // R(A) := UpValue[B]
        opcode{0, 1, OpArgU, OpArgK, IABC /* */, "GETTABUP", getTabUp}, // R(A) := UpValue[B][RK(C)]
        opcode{0, 1, OpArgR, OpArgK, IABC /* */, "GETTABLE", getTable}, // R(A) := R(B)[RK(C)]
        opcode{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABUP", setTabUp}, // UpValue[A][RK(B)] := RK(C)
        opcode{0, 0, OpArgU, OpArgN, IABC /* */, "SETUPVAL", setUpval}, // UpValue[B] := R(A)
        opcode{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABLE", setTable}, // R(A)[RK(B)] := RK(C)
        opcode{0, 1, OpArgU, OpArgU, IABC /* */, "NEWTABLE", newTable}, // R(A) := {} (size = B,C)
        opcode{0, 1, OpArgR, OpArgK, IABC /* */, "SELF    ", self},     // R(A+1) := R(B); R(A) := R(B)[RK(C)]
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "ADD     ", add},      // R(A) := RK(B) + RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SUB     ", sub},      // R(A) := RK(B) - RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "MUL     ", mul},      // R(A) := RK(B) * RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "MOD     ", mod},      // R(A) := RK(B) % RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "POW     ", pow},      // R(A) := RK(B) ^ RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "DIV     ", div},      // R(A) := RK(B) / RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "IDIV    ", idiv},     // R(A) := RK(B) // RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BAND    ", band},     // R(A) := RK(B) & RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BOR     ", bor},      // R(A) := RK(B) | RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BXOR    ", bxor},     // R(A) := RK(B) ~ RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SHL     ", shl},      // R(A) := RK(B) << RK(C)
        opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SHR     ", shr},      // R(A) := RK(B) >> RK(C)
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "UNM     ", unm},      // R(A) := -R(B)
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "BNOT    ", bnot},     // R(A) := ~R(B)
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "NOT     ", not},      // R(A) := not R(B)
        opcode{0, 1, OpArgR, OpArgN, IABC /* */, "LEN     ", length},   // R(A) := length of R(B)
        opcode{0, 1, OpArgR, OpArgR, IABC /* */, "CONCAT  ", concat},   // R(A) := R(B).. ... ..R(C)
        opcode{0, 0, OpArgR, OpArgN, IAsBx /**/, "JMP     ", jmp},      // pc+=sBx; if (A) close all upvalues >= R(A - 1)
        opcode{1, 0, OpArgK, OpArgK, IABC /* */, "EQ      ", eq},       // if ((RK(B) == RK(C)) ~= A) then pc++
        opcode{1, 0, OpArgK, OpArgK, IABC /* */, "LT      ", lt},       // if ((RK(B) <  RK(C)) ~= A) then pc++
        opcode{1, 0, OpArgK, OpArgK, IABC /* */, "LE      ", le},       // if ((RK(B) <= RK(C)) ~= A) then pc++
        opcode{1, 0, OpArgN, OpArgU, IABC /* */, "TEST    ", test},     // if not (R(A) <=> C) then pc++
        opcode{1, 1, OpArgR, OpArgU, IABC /* */, "TESTSET ", testSet},  // if (R(B) <=> C) then R(A) := R(B) else pc++
        opcode{0, 1, OpArgU, OpArgU, IABC /* */, "CALL    ", call},     // R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
        opcode{0, 1, OpArgU, OpArgU, IABC /* */, "TAILCALL", tailCall}, // return R(A)(R(A+1), ... ,R(A+B-1))
        opcode{0, 0, OpArgU, OpArgN, IABC /* */, "RETURN  ", _return},  // return R(A), ... ,R(A+B-2)
        opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORLOOP ", forLoop},  // R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }
        opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORPREP ", forPrep},  // R(A)-=R(A+2); pc+=sBx
        opcode{0, 0, OpArgN, OpArgU, IABC /* */, "TFORCALL", nil},      // R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
        opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "TFORLOOP", nil},      // if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx }
        opcode{0, 0, OpArgU, OpArgU, IABC /* */, "SETLIST ", setList},  // R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
        opcode{0, 1, OpArgU, OpArgN, IABx /* */, "CLOSURE ", closure},  // R(A) := closure(KPROTO[Bx])
        opcode{0, 1, OpArgU, OpArgN, IABC /* */, "VARARG  ", vararg},   // R(A), R(A+1), ..., R(A+B-2) = vararg
        opcode{0, 0, OpArgU, OpArgU, IAx /*  */, "EXTRAARG", nil},      // extra (larger) argument for previous opcode
    }


    单元测试

    $ go run main.go luac.out 
    1
    2
    1
    3
    2


========================

元编程

    所谓的元编程，就是能够处理程序的程序，其中的“处理”包括读取、生成、分析、转换等。而元编程，就是值编写元程序的编程技术。

    lua通过标准库debug提供了类似反射的能力，通过load（）提供了运行时执行任意代码的能力，

元表以及元方法
    元表

    在lua中，每个值都可以有一个元表，如果值的类型是表或者用户数据，则可以用欧自己特有的元表。其他类型的值则是没中类型共享一个元表。新建的表默认没有元表，nil、布尔、数字、函数类型默认也没有元表，不过string标准库给字符串类型设置了元表。lua标准库提供了getmetatable（），可以用来获取与值关联的元表。

    lua标准库也提供了setmetatable函数，不过这个函数仅仅可以给表设置元表，对于其他类型的值，则必须通过debug库里面的setmetatable设置元表

    元方法

    元表也是普通的表，没什么特别的。真正让元表变得与众不同的是元方法。比如，当对两个表进行加法运算，lua会检查这两个表是不是有元表，如果有，进一步检查是否有_add元方法，如果有，就将两个表作为参数调用这个方法，将返回结果作为运算结果。

    通过元表和元方法，lua提供了中插件的机制。由于元表是普通的表，元方法也是普通的函数，因此可以使用lua代码来编写插件，扩展lua语言。

    实际上，lua里面的每个运算符都有一个元方法与之对应，另外，表操作以及函数调用也有相应的元方法。

    支持元表

    每个表都可以拥有自己的元表，其他值则是，没中类型共享一个元表，所以要做的第一件事就是将值和元表关联起来。对于表来说，显而易见的做法就是给底层的结构体添加一个字段，用来存放元表。对于其他类型的值，可以把元表存放在注册表里。

    state/lua_table.go 给结构体luaTable添加字段metatable

    
    type luaTable struct {
        metatable *luaTable
        arr       []luaValue
        _map      map[luaValue]luaValue
    }

    state/lua_value.go 增加函数

    setMetatable 用来给值关联元表

    func setMetatable(val luaValue, mt *luaTable, ls *luaState) {
        if t, ok := val.(*luaTable); ok {
            t.metatable = mt
            return
        }
        key := fmt.Sprintf("_MT%d", typeOf(val))
        ls.registry.put(key, mt)
    }

    先判断值否不是表，如果是，直接修改其元表字段即可。否则的话，根据变量类型把元表存储在注册表里面，这样就达到了按照类型共享元表的目的。虽说注册表也是一个普通的表，不过按照约定，下划线开头，后跟大写字母的字段名是保留给lua实现使用的，所以使用了“_MT!1”这样的字段名，以免和用户（通过api）放在注册表里面的数据产生冲突。另外，如果传递给函数的元表是nil值，效果就相当于删除元表。

    getMetatable（）返回与给定值关联的元表


    func getMetatable(val luaValue, ls *luaState) *luaTable {
        if t, ok := val.(*luaTable); ok {
            return t.metatable
        }
        key := fmt.Sprintf("_MT%d", typeOf(val))
        if mt := ls.registry.get(key); mt != nil {
            return mt.(*luaTable)
        }
        return nil
    }

    同样先判断是不是表，如果是，直接返回其元表的字段即可；否则，根据值的类型从注册表中取出与该类型关联的元表并返回，如果值没有元表与之关联，返回值就是nil。这样lua值就已经可以关联元方法了。


    调用元方法

    lua中的每个运算符，以及表操作和函数调用，均有元方法与之对应。

    1、算术运算符

    与算术运算符（包括按位运算符）对应的元方法有16个。

    state/api_arith.go，修改结构体operator，添加字段，metamethod，用来存储元方法的名称

    type operator struct {
        metamethod  string
        integerFunc func(int64, int64) int64
        floatFunc   func(float64, float64) float64
    }

    修改operators变量的初始化代码，给每个运算符都关联元方法名

        var operators = []operator{
        operator{"__add", iadd, fadd},
        operator{"__sub", isub, fsub},
        operator{"__mul", imul, fmul},
        operator{"__mod", imod, fmod},
        operator{"__pow", nil, pow},
        operator{"__div", nil, div},
        operator{"__idiv", iidiv, fidiv},
        operator{"__band", band, nil},
        operator{"__bor", bor, nil},
        operator{"__bxor", bxor, nil},
        operator{"__shl", shl, nil},
        operator{"__shr", shr, nil},
        operator{"__unm", iunm, funm},
        operator{"__bnot", bnot, nil},
}

    对于算术运算符，只有当一个操作数（或者无法自动转换为）不是数字时，才会调用对应的元方法。算术运算最终是由Arith（）实现的，对其进行修改调整：


    func (self *luaState) Arith(op ArithOp) {
    >>>>>>......
        if result := _arith(a, b, operator); result != nil {
            self.stack.push(result)
            return
        }

        mm := operator.metamethod
        if result, ok := callMetamethod(a, b, mm, self); ok {
            self.stack.push(result)
            return
        }

        panic("arithmetic error!")
    }

    如果操作数都是数字，或者可以转化为数字，则执行正常的算术运算逻辑，否则尝试查找并执行算术元方法，如果找不到相应的元方法，则调用panic报错。

    callMetamethod方法负责查找并调用元方法，因为在他的地方也会用到，将其定义到文件lua_value.go中：
    state/lua_value.go

    func callMetamethod(a, b luaValue, mmName string, ls *luaState) (luaValue, bool) {
        var mm luaValue
        if mm = getMetafield(a, mmName, ls); mm == nil {
            if mm = getMetafield(b, mmName, ls); mm == nil {
                return nil, false
            }
        }

        ls.stack.check(4)
        ls.stack.push(mm)
        ls.stack.push(a)
        ls.stack.push(b)
        ls.Call(2, 1)
        return ls.stack.pop(), true
    }

    callMetamethod()接收4个参数，返回两个值，前两个参数，是算术运算的两个操作数，第三个参数给出元方法名，如果过操作数不是表，则需要通过第四个参数来访问注册表，从中查找元表。第一个返回值是元方法的执行结果，不过由于元方法可能会返回任何值，包括nil和false，因此只能使用第二个返回值来表示能否找到方法。

    先依次检查两个操作数是否有对应的元方法，如果找不到对应的元方法，则直接返回nil以及false，如果任何一个操作数都对应的元方法，则以两个操作数为参数调用元方法，将元方法调用的结果以及true返回。对于一元运算符，两个操作数可以传入同一个值，这样一元和二元元方法就可以进行统一的管理了。

    getMetafield()也放在lua_value.go中

    func getMetafield(val luaValue, fieldName string, ls *luaState) luaValue {
        if mt := getMetatable(val, ls); mt != nil {
            return mt.get(fieldName)
        }
        return nil
    }

    长度元方法

    对于长度运算（#），lua首先判断值是不是字符串，如果是，结果就是字符串的长度；否则检查值是不是有__len方法，如果有，则以值为参数调用元方法，将元方法返回值作为结果，如果还找不到对应的元方法，但值是表，结果就是表的长度。长度运算时由lua api方法Len()实现的。
    state/api_misc.go

    func (self *luaState) Len(idx int) {
        val := self.stack.get(idx)

        if s, ok := val.(string); ok {
            self.stack.push(int64(len(s)))
        } else if result, ok := callMetamethod(val, val, "__len", self); ok {
            self.stack.push(result)
        } else if t, ok := val.(*luaTable); ok {
            self.stack.push(int64(t.len()))
        } else {
            panic("length error!")
        }
    }

    拼接元方法

     对于拼接运算（...）,如果两个操作数都是字符串或者数字，则进行字符串的拼接；否则，尝试调用__concat元方法，查找和调用规则与二元算术的元方法相同。拼接运算时由api的concat（）实现的，对其进行调整

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

                b := self.stack.pop()
                a := self.stack.pop()
                if result, ok := callMetamethod(a, b, "__concat", self); ok {
                    self.stack.push(result)
                    continue
                }

                panic("concatenation error!")
            }
        }
        // n == 1, do nothing
    }


    比较元方法

    1、__eq

    对于 == 运算，当且仅当两个操作数时不同的表时，此阿辉尝试执行__eq元方法。元方法的检查和执行规则与二元算术运算符类似，但是执行结果就被转换为布尔值。


    state/api_compare.go ，修改__eq

    func _eq(a, b luaValue, ls *luaState) bool { //增加了ls参数
        switch x := a.(type) {
    。。。。。。
        case *luaTable:
            if y, ok := b.(*luaTable); ok && x != y && ls != nil {
                if result, ok := callMetamethod(x, y, "__eq", ls); ok {
                    return convertToBoolean(result)
                }
            }
            return a == b
        default:
            return a == b
        }
    }


    注意，_eq()函数增加了*luaState参数，因为查找元方法的时候需要用到这个参数，另外，如果这个参数传入的是nil，表示不希望执行_eq元方法

    2、__lt

    <运算，如果两个操作数都是数字，则进行数字比较；如果两个操作数都是字符串，则进行字符串的表，否则尝试调用__lt元方法。元方法的查找和调用规则与__eq类似。

    func _lt(a, b luaValue, ls *luaState) bool {//增加了ls参数
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
        //added
        if result, ok := callMetamethod(a, b, "__lt", ls); ok {
            return convertToBoolean(result)
        } else {
            panic("comparison error!")
        }
    }

    3、__le

    <= 规则类似于<,与之对应的元方法是_le,有所不同的是，如果lua找不到__le元方法，则会尝试调用__lt元方法（假设 a<= b <===> not (b < a)）

    func _le(a, b luaValue, ls *luaState) bool { //增加了ls参数
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
        // added
        if result, ok := callMetamethod(a, b, "__le", ls); ok {
            return convertToBoolean(result)
        } else if result, ok := callMetamethod(b, a, "__lt", ls); ok {
            return !convertToBoolean(result)
        } else {
            panic("comparison error!")
        }
    }

    因为给这几个函数增加了参数，因此还需要在state/api_compare.go中修改Compare方法，
    在调用者三个函数的时候把self作为第三个参数传入


    func (self *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
        if !self.stack.isValid(idx1) || !self.stack.isValid(idx2) {
            return false
        }

        a := self.stack.get(idx1)
        b := self.stack.get(idx2)
        switch op {
        case LUA_OPEQ:
            return _eq(a, b, self)
        case LUA_OPLT:
            return _lt(a, b, self)
        case LUA_OPLE:
            return _le(a, b, self)
        default:
            panic("invalid compare op!")
        }
    }

    索引元方法

    索引元方法有两个作用：如果一个值不是表，索引元方法可以让这个值用起来像一个表；如果一个值是表，索引元方法可以拦截表的操作。索引元方法有两个：__index 、__newindex.前者对应索引取值操作，后者对应索引赋值操作。

    1、__index

    当lua执行t[k]表达式的时候，如果t不是表，或者k在表中不存在，就会出发__index元方法。虽然名为元方法，但实际上，__index元方法既可以是函数，也可以是表。如果是函数，那么lua会以t 、k为参数调用该函数，以函数的返回值为结果；如果是表，lua会以k为key，访问该表，以value为结果（可能会继续出发__index元方法）

    对表的访问是有GetTable等api实现的，这些方法友都调用了getTable方法。对gettable方法进行修改

    state/api_get.go

        //增加了参数raw
    func (self *luaState) getTable(t, k luaValue, raw bool) LuaType {
        if tbl, ok := t.(*luaTable); ok {
            v := tbl.get(k)
            //新增加的
            if raw || v != nil || !tbl.hasMetafield("__index") {
                self.stack.push(v)
                return typeOf(v)
            }
        }

        if !raw {
            if mf := getMetafield(t, "__index", self); mf != nil {
                switch x := mf.(type) {
                case *luaTable:
                    return self.getTable(x, k, false)
                case *closure:
                    self.stack.push(mf)
                    self.stack.push(t)
                    self.stack.push(k)
                    self.Call(2, 1)
                    v := self.stack.get(-1)
                    return typeOf(v)
                }
            }
        }

        panic("index error!")
    } 

    给getTable（）增加了参数raw，如果该蛇叔值为true，表示需要忽略元方法。如果t是表，并且key已经在表里面了，或者需要忽略元方法，或者表没有__index元方法，则维持原来的逻辑不变，否则尝试调用元方法。

    如果__index是表，则把表的访问操作转发给该表，否则__index是函数，调用该函数。  

    增加hasMetafield方法

    state/lua_table.go

   func (self *luaTable) hasMetafield(fieldName string) bool {
        return self.metatable != nil &&
            self.metatable.get(fieldName) != nil
    } 


    2、__newindex

    当lua执行t[k]=v语句的时候，如果t不是表，或者k不在表中，就会触发__newindex元方法。和__index元方法一样，__newindex元方法也可以是函数或者表。如果是函数，那么lua会以t、k、v为参数调用该函数；如果是表，撸啊会以k为key，v为value给该表赋值（可能会继续触发__newindex方法）

    对表的写入是由SetTable等api实现的，这些方法又都调用了setTable（）,修改这个方法

    state/api_set.go

    func (self *luaState) setTable(t, k, v luaValue, raw bool) {
        if tbl, ok := t.(*luaTable); ok {
            if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
                tbl.put(k, v)
                return
            }
        }

        if !raw {
            if mf := getMetafield(t, "__newindex", self); mf != nil {
                switch x := mf.(type) {
                case *luaTable:
                    self.setTable(x, k, v, false)
                    return
                case *closure:
                    self.stack.push(mf)
                    self.stack.push(t)
                    self.stack.push(k)
                    self.stack.push(v)
                    self.Call(3, 0)
                    return
                }
            }
        }

        panic("index error!")
    }


    因为setTable getTable增加了参数，因此需要在api_get.go api_set.go中修改GetTable()\GetField()\GetI()\SetTable()\SetField()\SetI()
    ,它们都需要触发元方法

    func (self *luaState) GetTable(idx int) LuaType {
        t := self.stack.get(idx)
        k := self.stack.pop()
        return self.getTable(t, k, false)   //需要触发元方法
    }

    函数调用元方法

    当试图调用一个非函数类型的值时，lua会看这个值是否有__call元方法，如果有，lua会以该值作为第一个参数，后跟元方法调用的其他参数，来调用元方法，以元方法返回值为返回值。函数调用是由Api方法Call（）实现的，调整该方法

    state/api_call.go

    func (self *luaState) Call(nArgs, nResults int) {
        val := self.stack.get(-(nArgs + 1))

        c, ok := val.(*closure)
        if !ok {
            if mf := getMetafield(val, "__call", self); mf != nil {
                if c, ok = mf.(*closure); ok {
                    self.stack.push(val)
                    self.Insert(-(nArgs + 2))
                    nArgs += 1
                }
            }
        }

        if ok {
            if c.proto != nil {
                self.callLuaClosure(nArgs, nResults, c)
            } else {
                self.callGoClosure(nArgs, nResults, c)
            }
        } else {
            panic("not function!")
        }
    }

    如果被调用的值不是函数，就根据前面的规则查找并调用元方法。元表和元方法并非只有lua内部实现可以使用，标准库以及自定义的函数也可以利用元表和元方法

    扩展LuaApi

    api/lua_state.go

    为接口LuaState添加方法

    type LuaState interface {
        。。。。。。
        GetMetatable(idx int) bool
        SetMetatable(idx int)
        RawLen(idx int) uint              //state/api_access.go
        RawEqual(idx1, idx2 int) bool     //state/api_compare.go
        RawGet(idx int) LuaType           //state/api_get.go
        RawSet(idx int)                   //state/api_set.go
        RawGetI(idx int, i int64) LuaType //state/api_get.go
        RawSetI(idx int, i int64)         //state/api_set.go
    }

    GetMetatable(idx int)  SetMetatable(idx int)用于操作元表，其他的方法和不带Raw前缀的版本基本一样，指示不会去尝试查找以及调用元方法。

    GetMetatable()

    GetMetatable()，查看置顶索引处的值是不是有元表，如果有，则把元表push到栈顶并返回true；否则栈的状态不会改变，返回false。

    state/api_get.go中实现

   func (self *luaState) GetMetatable(idx int) bool {
        val := self.stack.get(idx)

        if mt := getMetatable(val, self); mt != nil {
            self.stack.push(mt)
            return true
        } else {
            return false
        }
    } 

    SetMetatable()

    SetMetatable,从栈顶弹出一个表，然后把指定索引处值的元表设置为这个表。

    func (self *luaState) SetMetatable(idx int) {
        val := self.stack.get(idx)
        mtVal := self.stack.pop()

        if mtVal == nil {
            setMetatable(val, nil, self)
        } else if mt, ok := mtVal.(*luaTable); ok {
            setMetatable(val, mt, self)
        } else {
            panic("table expected!") // todo
        }
    }

    先把栈顶弹出，如果它是nil，实际效果就是清除元表；如果它是表，就用它设置元表；否则就调用panic报错。

单元测试

    为了测试，实现简化版的
    func getMetatable(ls LuaState) int {
        if !ls.GetMetatable(1) {
            ls.PushNil()
        }
        return 1
    }


    通过api GetMetatable（）实现对应的标准库函数，由于getmetatable（）只有一个参数，所以调用GetMetatable（）的时候传入索引1就可以了，如果值有元表，方法结束后，元表已经在栈顶了，否则需要把nil值push到栈顶。最后返回1，把栈顶的值（元表或者nil）返回给Lua函数。

    func setMetatable(ls LuaState) int {
        ls.SetMetatable(1)
        return 1
    }

    将它们注册到全局环境

    func main() {
        if len(os.Args) > 1 {
            data, err := ioutil.ReadFile(os.Args[1])
            if err != nil {
                panic(err)
            }

            ls := state.New()
            ls.Register("print", print)
            ls.Register("getmetatable", getMetatable)
            ls.Register("setmetatable", setMetatable)
            ls.Load(data, os.Args[1], "b")
            ls.Call(0, 0)
        }
    }

        $ go run main.go luac.out 
        [1, 2]
        [3, 4]
        [2, 4]
        [3, 6]
        5
        false
        true
        [3, 6]

====================

迭代器

    lua语言支持两种形式的for循环：数值for 以及通用for。数值for循环用于在两个数值范围内按照一定的的step进行迭代；通用for循环常用于对表进行迭代。不过通用for循环之所以“通用”，就在于它并不仅仅适用于表，实际上，通用for循环可以对任何迭代器进行迭代。

    迭代器介绍
    为了对集合或者容器进行遍历，迭代器需要保存一些内部状态。在Lua中，是用函数来表示迭代器，内部状态可以使用闭包捕获。

    看一个对数组进行迭代的例子

    function ipairs(t)
        local i = 0
        return function()
            i = i + 1
            if t[i] == nil then
                return nil, nil
            else
                return i, t[i]
            end
        end
    end


    可以把上面的ipairs（）看作一个工厂，函数每次调用都会返回一个数组迭代器。迭代器从外部捕获了t 、i这两个变量（Upvalue），把它们作为内部状态，用于控制迭代，并痛殴第一个返回值（是不是nil）通知使用者迭代是否结束。如何创建迭代器，并利用它对数组进行遍历。

    t = {10, 20, 30}

    iter = ipairs(t)
    while true do
        local i, v = iter()
        if i == nil then
            break
        end
        print(i, v)
    end

    ----------------
    t = {10, 20, 30}

    for i, v in ipairs(t) do
        print(i, v)
    end

    上面是给数组（序列）创建迭代器，如何给关联数组创建迭代器？由于关联数组没有办法通过递增下标的方式迭代，所以lua标准库提供了next（）创建关联数组的迭代器。next（）接收两个参数---表和key，返回两个值---下一个键值对。如果传递给next（）的key是nil，表示迭代开始；如果next（）返回的kye是nil，表示迭代结束。

    function pairs(t)
        local k, v
        return function()
            k, v = next(t, k)
            return k, v
        end
    end

    t = {a=10, b=20, c=30}

    for k, v in pairs(t) do
        prinf(k,  v)
    end

    通用的for循环语句的一般形式：

    for var1,....,varn in explist do block end

    等价于下面的代码：

    do
        local _f, _s, var = explist
        while true do
            local var1,..., varn = _f(_s, _var)
            if var1 == nil then break end
            _var = var1
            block
        end
    end


    其中_f _s _var是通用for循环内部使用的，由explist求值得到（多重赋值、多退少补），_f是迭代器函数，_s是一个不变量，_var是控制变量。对比前面的pairs：_f<>next(),_s<>表,_var<>用于存放key。因此，虽然可用闭包保存迭代器内部状态，不过通用for可以帮我们保存一些状态，这样很多时候就可以免去闭包创建的烦恼。知道了这个，就可以直接使用next（）对表进行遍历

       t = {a=10, b=20, c=30} 

       for k, v in next, t, nil do
        print(k, v)
       end 

    在LUa虚拟机级别，有两个专门的指令用于实现for-in，tforcall\tforloop。

    next（）

    伪代码实现：

    function next(table, key)
        if key == nil then
            nextKey = table.firstKey()
        else
            nextKey = table.nextKey(key)
        end
        if nextKey ~= nil then
            return nextKey, table[nextKey]
        else
            return nil
    end

    Lua api提供了next（）方法，可以用来实现标准库里面的next（）函数。不过要实现api级别的Next（），必须要得到表结构体的内部支持。修改luaTable结构体，让其支持key的遍历；然后实现api层面的Next（）方法；之后实现标准库层面的next（）。

    修改luaTable结构体

    type luaTable struct {
        。。。。。。
        keys      map[luaValue]luaValue // used by next()
        lastKey   luaValue              // used by next()
        changed   bool                  // used by next()
    }

    因为go的map不保证遍历的顺序性（甚至同样内容的map，两次遍历的顺序也可能不一样），所以只能在遍历开始之前把所有的key都固定下来，保存在keys字段里面。

    添加nextKey方法

    func (self *luaTable) nextKey(key luaValue) luaValue {
        if self.keys == nil || (key == nil && self.changed) {
            self.initKeys()
            self.changed = false
        }

        nextKey := self.keys[key]
        if nextKey == nil && key != nil && key != self.lastKey {
            panic("invalid key to 'next'")
        }

        return nextKey
    }


    nextKey()根据传入的key返回表的下一个key，完全是为api方法Next（）而设计了。如果传入的key就是nil，表示遍历开始，需要把所有的key都收集到kyes中，keys的键值对记录了表的key和xiayigekey的关系。所以keys字段初始化好了以后，直接根据传入的参数取值，并返回就可以了。

    func (self *luaTable) initKeys() {
        self.keys = make(map[luaValue]luaValue)
        var key luaValue = nil
        for i, v := range self.arr {
            if v != nil {
                self.keys[key] = int64(i + 1)
                key = int64(i + 1)
            }
        }
        for k, v := range self._map {
            if v != nil {
                self.keys[key] = k
                key = k
            }
        }
        self.lastKey = key
    }

    因为标准库next()和api 方法Next（）考虑的是整个表（包括数组部分和关联数组部分），因此initKeys（）需要把数组索引和关系数组键都手机起来。

    扩展LUa api

    api/lua_state.go 

    type LuaState interface {
        .....
        Next(idx int) bool
    }

    Next()方法根据key获取表的下一个键值对，其中表的索引由参数指定，上一个键从栈顶弹出。如果从栈顶弹出的key是nil，说明刚开始遍历表，把白哦的第一个键值对push到栈顶并返回true；否则，如果遍历还没有结束，把下一个键值对push到栈顶并返回true；如果表是空的，或者遍历已经结束，不用向栈里面push任何值，直接返回false即可。

    state/api_misc.go

    func (self *luaState) Next(idx int) bool {
        val := self.stack.get(idx)
        if t, ok := val.(*luaTable); ok {
            key := self.stack.pop()
            if nextKey := t.nextKey(key); nextKey != nil {
                self.stack.push(nextKey)
                self.stack.push(t.get(nextKey))
                return true
            }
            return false
        }
        panic("table expected!")
    }

    先根据索引拿到表（如果拿到的不是表，就调用panic报错），然后从栈顶弹出上一个key，解析来使用这个key调用表的nextKey方法，如果方法返回nil，说明遍历已经结束，返回false即可；反之，把下一个键值对push到栈顶，返回true

    实现next（）函数

    修改main。go

    func next(ls LuaState) int {
        ls.SetTop(2) /* 入股参数2不存在，则设置nil */
        if ls.Next(1) {
            return 2
        } else {
            ls.PushNil()
            return 1
        }
    }

    因为next()的第二个参数（key）是可选的，所以首先调用SetTop，以便在用户没有提供这个参数的时候给它补上默认值nil。这样就能保证栈里面肯定有两个值，索引1处是表，索引2处是上一个key。然后通过索引1调用Next（），她会把key从栈顶弹出，如果Next（）返回true，说明遍历还没有结束，下一个键值对已经在栈顶了，返回2就可以。，反之，遍历已经结束，需要自己向栈顶push nil，然后返回1


    通用for循环指令

    通用for循环也是利用tforcal tforloop这两个指令实现的

    其中tforcall（iABC 模式）伪代码：

    R（A+3），... ,R(A+2+C) := R(A)(R(A+1), R(A+2))

    tforloop(iAsBx模式)

    if R（A+1）~= nil then {
        R(A)=R(A+1);pc += sBx
    }

    vm/inst_call.go实现tforcall
    
    func tForCall(i Instruction, vm LuaVM) {
        a, _, c := i.ABC()
        a += 1

        _pushFuncAndArgs(a, 3, vm)
        vm.Call(2, c)
        _popResults(a+3, c+1, vm)
    }


    vm/inst_for.go实现tforloop


    func tForLoop(i Instruction, vm LuaVM) {
        a, sBx := i.AsBx()
        a += 1

        if !vm.IsNil(a + 1) {
            vm.Copy(a+1, a)
            vm.AddPC(sBx)
        }
    }

    修改opcodes.go 修改指令表，注册上面的函数
    vm/opcodes.go


单元测试

    修改main。go添加函数


    func pairs(ls LuaState) int {
        ls.PushGoFunction(next) /* will return generator, */
        ls.PushValue(1)         /* state, */
        ls.PushNil()
        return 3
    }

    这个函数实际上就是返回了next函数（对应_f)\表（_s）\nil（_var），这三个值而已。类似lua代码

    function pairs(t)
        return next, t, nil
    end


    添加ipairs

    func iPairs(ls LuaState) int {
        ls.PushGoFunction(_iPairsAux) /* iteration function */
        ls.PushValue(1)               /* state */
        ls.PushInteger(0)             /* initial value */
        return 3
    }
    数组版的next（）
    func _iPairsAux(ls LuaState) int {
        i := ls.ToInteger(2) + 1
        ls.PushInteger(i)
        if ls.GetI(1, i) == LUA_TNIL {
            return 1
        } else {
            return 2
        }
    }

    类似lua代码

    function inext(t, i)
        local nextIdx = i + 1
        local nextVal = t[nextIdx]
        if nextVal == nil then
            retrun nil
        else
            return nextIdx, nextVal
        end
    end

    修改mian，注册pairs ipairs

    func main() {
        if len(os.Args) > 1 {
            data, err := ioutil.ReadFile(os.Args[1])
            if err != nil {
                panic(err)
            }

            ls := state.New()
            ls.Register("print", print)
            ls.Register("getmetatable", getMetatable)
            ls.Register("setmetatable", setMetatable)
            ls.Register("next", next)
            ls.Register("pairs", pairs)
            ls.Register("ipairs", iPairs)
            ls.Load(data, os.Args[1], "b")
            ls.Call(0, 0)
        }
    }

=================

错误和异常处理

    介绍

    lua并没有在语法层面上直接支持异常处理，不过在标准库中提供了一些函数，可以用来抛出异常或者捕获异常。

    function lock2seconds (lock)
        if not lock:tryLock() then
            error("Unable to acquire the lock!")
        end
        pcall(function() sleep(2000) end)
        lock:unlock()
    end

    error()抛出异常，异常抛出以后，正常的函数执行结束，然后异常逐步向外传播，直到被pcall（）捕获。通常可以使用字符串表示异常信息，实际上，error（）可以把任何lua值当作异常抛出

    error（{err="Unable to acquire the lock!"}）

    如果抛出的异常时字符串或者表，由于lua的函数调用语法糖允许载有且只有一个参数，并且该参数是字符串字面领或者表构造器的时候省略圆括号，所以，也可以让error（）看起来像个关键字

    error {err="Unable to acquire the lock!"}

    上面使用匿名函数来保护可能会抛出异常的代码块，然后交给pcall（）调用，对于这个例子来说，完全没有必要使用匿名函数，直接把被调函数和参数一起传递给pcall就可以了

    function lock2seconds(lock)
        lock:lock()
        pcall(sleep, 2000)
        lock:unlock()
    end

    pcall()会在保护模式下调用被调函数，如果一切正常，pcall（）返回true和被调函数返回的全部值；如果被调过程中有异常抛出，pcall不捕获异常，返回false以及异常。通常需要检查pcall的第一个返回值，看看调用是否正常，然后根据情况进行一步处理

    function lock2seconds (lock)
        lock:lock()
        local ok, msg = pcall(sleep, 2000)

        lock:unlock()

        if ok then
            print("ok")
        else
            print("error:" .. msg)
        end

     end   

     实际上，error（）还有一个可选参数level，pcall函数还有另一个版本xpcall（）


     错误和异常处理api

     lua api提供了两个方法 Error（）和PCall（）。对应标准库的error（）、pcall（）。

     api/lua_state.go

    type LuaState interface {
        ...
        Error() int
        PCall(nArgs, nResults, msgh int) int
    }

    定义常量
    api/consts.go


    const (
        LUA_OK = iota
        LUA_YIELD
        LUA_ERRRUN
        LUA_ERRSYNTAX
        LUA_ERRMEM
        LUA_ERRGCMM
        LUA_ERRERR
        LUA_ERRFILE
    )

    这些常量用来表示函数加载或者执行的状态，LUA_OK、LUA_ERRRUN会在PCall（）方法中用到。

    Error()

    从栈顶弹出一个值，把该值作为错误抛出，虽然方法以整数位返回值，但是实际上该方法是没有办法正常返回的，之所以有返回值，完全是为了方便用户使用。比如正在实现某个go函数，并且在某处需要调用Error（）抛出异常，就可以通过直接返回Error（）的返回值来简化代码

    func gofunc(ls LuaState) int{
        return ls.Error()
    }

    state/api_misc.go

    func (self *luaState) Error() int {
        err := self.stack.pop()
        panic(err)
    }

    从栈顶把错误对象弹出，然后调用go的panic抛出即可。

    PCall()

    与Call比较类似，区别在于PCall会捕获函数调用过程中产生的错误，而Call会错误传播出去。如果没有错误产生，它们的行为完全一样，最后妇女会LUA_OK.如果有错误产生，PCall会捕获错误，把错误对象留在栈顶，并返回相应的错误吗。PCall的第三个参数用于指定错误处理程序。

    state/api_call.go

    func (self *luaState) PCall(nArgs, nResults, msgh int) (status int) {
        caller := self.stack
        status = LUA_ERRRUN

        // catch error
        defer func() {
            if err := recover(); err != nil {
                if msgh != 0 {
                    panic(err)
                }
                for self.stack != caller {
                    self.popLuaStack()
                }
                self.stack.push(err)
            }
        }()

        self.Call(nArgs, nResults)
        status = LUA_OK
        return
    }


    既然使用go的panic抛出错误，自然就需要使用defer-recover机制来捕获异常。如果一切正常，返回LUA_OK,反之，捕获并处理错误

    defer func() {
            if err := recover(); err != nil {
                if msgh != 0 {
                    panic(err)
                }
                for self.stack != caller {
                    self.popLuaStack()
                }
                self.stack.push(err)
            }
        }()

        调用go内置的recover（）从错误中恢复，然后从调用栈顶依次弹出调用帧， 直到到达发起调用的调用帧为止，然后把所错误对象push到栈顶，返回LUA_ERRUN.


    error pcall

    在main。go实现简化版的error

    func error(ls LuaState) int {
        return ls.Error()
    }

    暂时只接收一个参数，所以错误对象已经在栈顶了，这届调用Error（）抛出错误就可以了。

    func pCall(ls LuaState) int {
        nArgs := ls.GetTop() - 1
        status := ls.PCall(nArgs, -1, 0)
        ls.PushBoolean(status == LUA_OK)
        ls.Insert(1)
        return ls.GetTop()
    }


    调用PCall并插入一个布尔类型的返回值。PCall会把被调函数和参数从栈顶弹出，然后把返回值或者错误对象留在栈顶，只需要根据PCall（）返回的状态码向栈顶push一个布尔值，然后将其移动到栈底，让它成为返回给Lua的第一个值就可以了

    简化版的error pcall就心这样了，注册到全局环境。

    func main() {
        if len(os.Args) > 1 {
            data, err := ioutil.ReadFile(os.Args[1])
            if err != nil {
                panic(err)
            }

            ls := state.New()
            ls.Register("print", print)
            ls.Register("getmetatable", getMetatable)
            ls.Register("setmetatable", setMetatable)
            ls.Register("next", next)
            ls.Register("pairs", pairs)
            ls.Register("ipairs", iPairs)
            ls.Register("error", error)
            ls.Register("pcall", pCall)
            ls.Load(data, os.Args[1], "b")
            ls.Call(0, 0)
        }
    }



单元测试

    $ go run main.go luac.out 
    true    2
    false   DIV BY ZERO !
    false   arithmetic error!

================

词法分析

    编译器介绍

    任何一种可以将编程语言（源语言）转换为另外一种编程语言（目标语言）的程序都可以称之为编译器。通常所说的编译器，一般指的是高级语言编译器。

    化繁为简，阶段化处理
    主要的编译阶段包括预处理、词法分析、语法分析、语义分析、中间代码生成、中间代码优化、目标代码生成等

    编译阶段：三个大阶段
    前段Front End、中端Middle End、后端Back End。


    Lua语言不支持宏等特性，不需要预处理。语义分析最重要的一个任务就是类型检查，由于lua是动态类型语言，在编译期间不需要进行类型检查，所以也不需要语义分析阶段。只关注三个阶段

    词法分析 语法分析 代码生成

    lua源文件---lexer---》token---parser---》ast抽象语法树--codegen---〉字节码

    自己实现词法分析器以及语法分析器，不用现成的工具

    Lua词法

    编译器在编译源码的时候，不是以字符为单位，而是以“token”为单位进行处理的，词法分析器的作用就是根据编程语言的词法规则，把源代码（字符流）分解为token流

    token按照其作用可以分为不同的类型：注释 、关键字、标识符、字面量、分隔符等。

    空白字符

    编译器除了需要使用换行符计算行号以外，会完全忽略空白自负，这样的语言为自由格式语言free-form、free-format

    lua也属于自由格式语言，使用关键字界定代码块。lua编译器会忽略换行符\n \r \t \v \f 空格符

    2、注释
    编译器可以完全忽略注释。lua支持两种形式的注释：
    -- 单行注释内容
    --[n个=[ 可以跨越多行]n个=]

    长注释其实就是--后跟一个长字符串字面量

    3、标识符
    主要用来命名变量。lua标识符以字符或者_开头，后跟数字 、字母或者_的任意组合。Lua是大小写敏感的。

    4、关键字
    由语言所保留，不能当作标识符使用。

    5、数字字面量
    3、314 当使用小数写法的时候，整数部分和小数部分都可以省略。3. 、.334。0.314E1、314e-1.

    16进制写法以0x 0X开头0xff，0x3.45fa，十六进制也可以使用科学计数法，不过指数部分要用字母p 、P表示，只能使用时禁止数子，表示的是2的多少次方，比如0xA45p-4.

    如果数字在里面量不包含小数和指数部分，也没有超过lua整数的表示范围，则会被lua解释为整数值，否则会被解释为浮点值

    6、字符串字面量
    lua字符串字面量有长短字符串两种。短字符串使用单引号或者双引号分隔，里面可以爆燃转义序列。
    lua短字符串字面量不能跨行（可以使用转移序列插入回车或者华换行符），惟一的例外就是\z，该转义序列会删除自己，以及紧随其后的空白符。

    长字符串字面量使用长方括号分隔，不支持转义序列。其中出现的换行符序列(\r\n \n\r \n \r)会被lua统一替换为\n,另外紧跟在左长方括号后面的第一个换行符会被lua删除。

    7、运算符和分隔符

实现词法分析器
    词法分析器一般使用有限状态机实现。

    compiler/lexer/lexer.go 定义结构体Lexer

     package lexer

    import "bytes"
    import "fmt"
    import "regexp"
    import "strconv"
    import "strings"   


    type Lexer struct {
        chunk         string // source code
        chunkName     string // source name
        line          int    // current line number
    }

    chunk用来保存将要进行词法分析的源代码
    line用来保存当前的行号
    这两个字段构成了词法分析器的内部状态。
    chunkName用来保存源文件的名称，用于在词法分析过程中，出错的是欧生成错误提示信息

    lexer.go 添加NewLexer（）方法

    func NewLexer(chunk, chunkName string) *Lexer {
        return &Lexer{chunk, chunkName, 1}
    }

    NewLexer根据源文件名以及源代码创建Lexer结构体实例，并将出事行号设置为1.

    给结构体Lexer 添加方法 NextToken

    func (self *Lexer) NextToken() (line, kind int, token string) {

        self.skipWhiteSpaces()
        if len(self.chunk) == 0 {
            return self.line, TOKEN_EOF, "EOF"
        }

        switch self.chunk[0] {
        case ';':
            self.next(1)
            return self.line, TOKEN_SEP_SEMI, ";"
        case ',':
            self.next(1)
            return self.line, TOKEN_SEP_COMMA, ","
        }

        return
    }

    NextToken()会掉过空白自负以及注释，返回下一个token（包括类型 行号），如果源代码以ing全部分析完毕，返回表示分析结果的特殊token。FSM可以使用swith-case语句来实现，词法分析器也是这样。

    定义类型Token

    lexer/token.go
    定义表示token类型的常量

    package lexer

    // token kind
    const (
        TOKEN_EOF         = iota           // end-of-file
        TOKEN_VARARG                       // ...
        TOKEN_SEP_SEMI                     // ;
        TOKEN_SEP_COMMA                    // ,
        TOKEN_SEP_DOT                      // .
        TOKEN_SEP_COLON                    // :
        TOKEN_SEP_LABEL                    // ::
        TOKEN_SEP_LPAREN                   // (
        TOKEN_SEP_RPAREN                   // )
        TOKEN_SEP_LBRACK                   // [
        TOKEN_SEP_RBRACK                   // ]
        TOKEN_SEP_LCURLY                   // {
        TOKEN_SEP_RCURLY                   // }
        TOKEN_OP_ASSIGN                    // =
        TOKEN_OP_MINUS                     // - (sub or unm)
        TOKEN_OP_WAVE                      // ~ (bnot or bxor)
        TOKEN_OP_ADD                       // +
        TOKEN_OP_MUL                       // *
        TOKEN_OP_DIV                       // /
        TOKEN_OP_IDIV                      // //
        TOKEN_OP_POW                       // ^
        TOKEN_OP_MOD                       // %
        TOKEN_OP_BAND                      // &
        TOKEN_OP_BOR                       // |
        TOKEN_OP_SHR                       // >>
        TOKEN_OP_SHL                       // <<
        TOKEN_OP_CONCAT                    // ..
        TOKEN_OP_LT                        // <
        TOKEN_OP_LE                        // <=
        TOKEN_OP_GT                        // >
        TOKEN_OP_GE                        // >=
        TOKEN_OP_EQ                        // ==
        TOKEN_OP_NE                        // ~=
        TOKEN_OP_LEN                       // #
        TOKEN_OP_AND                       // and
        TOKEN_OP_OR                        // or
        TOKEN_OP_NOT                       // not
        TOKEN_KW_BREAK                     // break
        TOKEN_KW_DO                        // do
        TOKEN_KW_ELSE                      // else
        TOKEN_KW_ELSEIF                    // elseif
        TOKEN_KW_END                       // end
        TOKEN_KW_FALSE                     // false
        TOKEN_KW_FOR                       // for
        TOKEN_KW_FUNCTION                  // function
        TOKEN_KW_GOTO                      // goto
        TOKEN_KW_IF                        // if
        TOKEN_KW_IN                        // in
        TOKEN_KW_LOCAL                     // local
        TOKEN_KW_NIL                       // nil
        TOKEN_KW_REPEAT                    // repeat
        TOKEN_KW_RETURN                    // return
        TOKEN_KW_THEN                      // then
        TOKEN_KW_TRUE                      // true
        TOKEN_KW_UNTIL                     // until
        TOKEN_KW_WHILE                     // while
        TOKEN_IDENTIFIER                   // identifier
        TOKEN_NUMBER                       // number literal
        TOKEN_STRING                       // string literal
        TOKEN_OP_UNM      = TOKEN_OP_MINUS // unary minus
        TOKEN_OP_SUB      = TOKEN_OP_MINUS
        TOKEN_OP_BNOT     = TOKEN_OP_WAVE
        TOKEN_OP_BXOR     = TOKEN_OP_WAVE
    )

    除了字面量和标识符，给其他的没中token都分配了一个单独的常量值。需要注意的是，因为在词法分析阶段，没有办法区分建好到底是二元减法运算符 好事一元取负运算符，所以将其命名为TOKEN_OP_MINUS，并定义了TOKEN_OP_UNM 和TOKEN_OP_SUB，这三个常量有相同的常量值。同样，TOKEN_OP_BNOT TOKEN_OP_BXOR TOKEN_OP_WAVE也是一样的。

    定义关联数组，将关键字和常量值一一对应

    var keywords = map[string]int{
        "and":      TOKEN_OP_AND,
        "break":    TOKEN_KW_BREAK,
        "do":       TOKEN_KW_DO,
        "else":     TOKEN_KW_ELSE,
        "elseif":   TOKEN_KW_ELSEIF,
        "end":      TOKEN_KW_END,
        "false":    TOKEN_KW_FALSE,
        "for":      TOKEN_KW_FOR,
        "function": TOKEN_KW_FUNCTION,
        "goto":     TOKEN_KW_GOTO,
        "if":       TOKEN_KW_IF,
        "in":       TOKEN_KW_IN,
        "local":    TOKEN_KW_LOCAL,
        "nil":      TOKEN_KW_NIL,
        "not":      TOKEN_OP_NOT,
        "or":       TOKEN_OP_OR,
        "repeat":   TOKEN_KW_REPEAT,
        "return":   TOKEN_KW_RETURN,
        "then":     TOKEN_KW_THEN,
        "true":     TOKEN_KW_TRUE,
        "until":    TOKEN_KW_UNTIL,
        "while":    TOKEN_KW_WHILE,
    }


    空白字符
    lexer/lexer.go

    定义方法skipWhiteSpaces

    func (self *Lexer) skipWhiteSpaces() {
        for len(self.chunk) > 0 {
            if self.test("--") {
                self.skipComment()
            } else if self.test("\r\n") || self.test("\n\r") {
                self.next(2)
                self.line += 1
            } else if isNewLine(self.chunk[0]) {
                self.next(1)
                self.line += 1
            } else if isWhiteSpace(self.chunk[0]) {
                self.next(1)
            } else {
                break
            }
        }
    }


    skipWhiteSpaces()不仅仅跳过了空白自负，更新了行号，同时还一并跳过了注释。

    test()判断剩余的源代码是否以某种字符串开头

     func (self *Lexer) test(s string) bool {
        return strings.HasPrefix(self.chunk, s)
    }   

    next跳过n个字节

    func (self *Lexer) next(n int) {
        self.chunk = self.chunk[n:]
    }

    isWhiteSpace判断自负是不是空白符

    func isWhiteSpace(c byte) bool {
        switch c {
        case '\t', '\n', '\v', '\f', '\r', ' ':
            return true
        }
        return false
    }

    判断字符是不是回车或者换行
    func isNewLine(c byte) bool {
        return c == '\r' || c == '\n'
    }

    注释

    skipComment跳过注释

     func (self *Lexer) skipComment() {
        self.next(2) // skip --

        // long comment ?
        if self.test("[") {
            if reOpeningLongBracket.FindString(self.chunk) != "" {
                self.scanLongString()
                return
            }
        }

        // short comment
        for len(self.chunk) > 0 && !isNewLine(self.chunk[0]) {
            self.next(1)
        }
    }
   
    长注释实际上是两个--，后跟一个长字符串，如果遇到长注释，只要跳过两个--，然后提取一个长字符串扔掉就可以了；如果是短字符串，则跳过两个--和后续的所有字符，直到遇到换行符为止。为了简单，借助正则表达式来检查左长方括号，定义正则表达式：

    var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)

    分隔符和运算符

     分隔符和运算符数量很多，所以NextToken的整个switch-case差不多都是用来读取这两种token的

        switch self.chunk[0] {
        case ';':
            self.next(1)
            return self.line, TOKEN_SEP_SEMI, ";"
        case ',':
            self.next(1)
            return self.line, TOKEN_SEP_COMMA, ","
        case '(':
            self.next(1)
            return self.line, TOKEN_SEP_LPAREN, "("
        case ')':
            self.next(1)
            return self.line, TOKEN_SEP_RPAREN, ")"
        case ']':
            self.next(1)
            return self.line, TOKEN_SEP_RBRACK, "]"
        case '{':
            self.next(1)
            return self.line, TOKEN_SEP_LCURLY, "{"
        case '}':
            self.next(1)
            return self.line, TOKEN_SEP_RCURLY, "}"
        case '+':
            self.next(1)
            return self.line, TOKEN_OP_ADD, "+"
        case '-':
            self.next(1)
            return self.line, TOKEN_OP_MINUS, "-"
        case '*':
            self.next(1)
            return self.line, TOKEN_OP_MUL, "*"
        case '^':
            self.next(1)
            return self.line, TOKEN_OP_POW, "^"
        case '%':
            self.next(1)
            return self.line, TOKEN_OP_MOD, "%"
        case '&':
            self.next(1)
            return self.line, TOKEN_OP_BAND, "&"
        case '|':
            self.next(1)
            return self.line, TOKEN_OP_BOR, "|"
        case '#':
            self.next(1)
            return self.line, TOKEN_OP_LEN, "#"
        case ':':
            if self.test("::") {
                self.next(2)
                return self.line, TOKEN_SEP_LABEL, "::"
            } else {
                self.next(1)
                return self.line, TOKEN_SEP_COLON, ":"
            }
        case '/':
            if self.test("//") {
                self.next(2)
                return self.line, TOKEN_OP_IDIV, "//"
            } else {
                self.next(1)
                return self.line, TOKEN_OP_DIV, "/"
            }
        case '~':
            if self.test("~=") {
                self.next(2)
                return self.line, TOKEN_OP_NE, "~="
            } else {
                self.next(1)
                return self.line, TOKEN_OP_WAVE, "~"
            }
        case '=':
            if self.test("==") {
                self.next(2)
                return self.line, TOKEN_OP_EQ, "=="
            } else {
                self.next(1)
                return self.line, TOKEN_OP_ASSIGN, "="
            }
        case '<':
            if self.test("<<") {
                self.next(2)
                return self.line, TOKEN_OP_SHL, "<<"
            } else if self.test("<=") {
                self.next(2)
                return self.line, TOKEN_OP_LE, "<="
            } else {
                self.next(1)
                return self.line, TOKEN_OP_LT, "<"
            }
        case '>':
            if self.test(">>") {
                self.next(2)
                return self.line, TOKEN_OP_SHR, ">>"
            } else if self.test(">=") {
                self.next(2)
                return self.line, TOKEN_OP_GE, ">="
            } else {
                self.next(1)
                return self.line, TOKEN_OP_GT, ">"
            }
        case '.':
            if self.test("...") {
                self.next(3)
                return self.line, TOKEN_VARARG, "..."
            } else if self.test("..") {
                self.next(2)
                return self.line, TOKEN_OP_CONCAT, ".."
            } else if len(self.chunk) == 1 || !isDigit(self.chunk[1]) {
                self.next(1)
                return self.line, TOKEN_SEP_DOT, "."
            }
        case '[':
            if self.test("[[") || self.test("[=") {
                return self.line, TOKEN_STRING, self.scanLongString()
            } else {
                self.next(1)
                return self.line, TOKEN_SEP_LBRACK, "["
            }
        case '\'', '"':
            return self.line, TOKEN_STRING, self.scanShortString()
        }


   虽然分隔符和运算符比较多，但是每个case都不难。最后两个case用来提取字符串字面量     

   长字符串字面量

      func (self *Lexer) scanLongString() string {
        openingLongBracket := reOpeningLongBracket.FindString(self.chunk)
        if openingLongBracket == "" {
            self.error("invalid long string delimiter near '%s'",
                self.chunk[0:2])
        }

        closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
        closingLongBracketIdx := strings.Index(self.chunk, closingLongBracket)
        if closingLongBracketIdx < 0 {
            self.error("unfinished long string or comment")
        }

        str := self.chunk[len(openingLongBracket):closingLongBracketIdx]
        self.next(closingLongBracketIdx + len(closingLongBracket))

        str = reNewLine.ReplaceAllString(str, "\n")
        self.line += strings.Count(str, "\n")
        if len(str) > 0 && str[0] == '\n' {
            str = str[1:]
        }

        return str
    } 

    先查找左右的长方括号，如果任何一个都找不到，则寿命源代码有语法错误，调用error回报错误并终止分析。然后是提取字符串字面量，把左右长方括号去掉，把换行符序列统一转换为换行符\n,在把开头的第一个换行符去掉，如果有的话，得到的就是最终的字符串。error使用源文件名、当前行号，以及传入的格式和参数抛出错误信息。


    func (self *Lexer) error(f string, a ...interface{}) {
        err := fmt.Sprintf(f, a...)
        err = fmt.Sprintf("%s:%d: %s", self.chunkName, self.line, err)
        panic(err)
    }

    同样使用正则表达式处理换行序列，定义正则表达式：

    var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")

    短字符串字面量

    func (self *Lexer) scanShortString() string {
        if str := reShortStr.FindString(self.chunk); str != "" {
            self.next(len(str))
            str = str[1 : len(str)-1]
            if strings.Index(str, `\`) >= 0 {
                self.line += len(reNewLine.FindAllString(str, -1))
                str = self.escape(str)
            }
            return str
        }
        self.error("unfinished string")
        return ""
    }

    先使用正则表达式提取短字符串，如果提取失败，说明源代码有语法错误，调用error报错，然后去掉字面量两边的引号，并在必要的时候调用escape对转义序列进行处理，得到最终的字符串。

    var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)

    func (self *Lexer) escape(str string) string {
        var buf bytes.Buffer

        for len(str) > 0 {
            if str[0] != '\\' {
                buf.WriteByte(str[0])
                str = str[1:]
                continue
            }

            if len(str) == 1 {
                self.error("unfinished string")
            }

            switch str[1] {
            case 'a':
                buf.WriteByte('\a')
                str = str[2:]
                continue
            case 'b':
                buf.WriteByte('\b')
                str = str[2:]
                continue
            case 'f':
                buf.WriteByte('\f')
                str = str[2:]
                continue
            case 'n', '\n':
                buf.WriteByte('\n')
                str = str[2:]
                continue
            case 'r':
                buf.WriteByte('\r')
                str = str[2:]
                continue
            case 't':
                buf.WriteByte('\t')
                str = str[2:]
                continue
            case 'v':
                buf.WriteByte('\v')
                str = str[2:]
                continue
            case '"':
                buf.WriteByte('"')
                str = str[2:]
                continue
            case '\'':
                buf.WriteByte('\'')
                str = str[2:]
                continue
            case '\\':
                buf.WriteByte('\\')
                str = str[2:]
                continue
            case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
                if found := reDecEscapeSeq.FindString(str); found != "" {
                    d, _ := strconv.ParseInt(found[1:], 10, 32)
                    if d <= 0xFF {
                        buf.WriteByte(byte(d))
                        str = str[len(found):]
                        continue
                    }
                    self.error("decimal escape too large near '%s'", found)
                }
            case 'x': // \xXX
                if found := reHexEscapeSeq.FindString(str); found != "" {
                    d, _ := strconv.ParseInt(found[2:], 16, 32)
                    buf.WriteByte(byte(d))
                    str = str[len(found):]
                    continue
                }
            case 'u': // \u{XXX}
                if found := reUnicodeEscapeSeq.FindString(str); found != "" {
                    d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
                    if err == nil && d <= 0x10FFFF {
                        buf.WriteRune(rune(d))
                        str = str[len(found):]
                        continue
                    }
                    self.error("UTF-8 value too large near '%s'", found)
                }
            case 'z':
                str = str[2:]
                for len(str) > 0 && isWhiteSpace(str[0]) { // todo
                    str = str[1:]
                }
                continue
            }
            self.error("invalid escape sequence near '\\%c'", str[1])
        }

        return buf.String()
    }    



    \ddd \xXX
    用正则表达式提取转义序列，然后将其解析为整数值，如果整数值超过0xFF,则调用error报错

    \u{XXX}

    用正则表达式提取转义序列，然后将其解析为整数值，然后调用go标准库的方法将Unicode代码转换为UTF-8编码格式的字节序列。如果代码超出范围，则调用error报错

    \z
    先跳过\z这个转义序列，跳过紧随其后的空白符。需要定义提取转义序列的正则表达式。

    var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
    var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
    var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

    数字字面量

    数字字面量以点或者数字开头，标识符和关键字以字母或者下划线开头。把提取这些token的代码写到switch-case的外面

    c := self.chunk[0]
    if c == '.' || isDigit(c) {
        token := self.scanNumber()
        return self.line, TOKEN_NUMBER, token
    }

    如果发现下一个token是数字字面量，调用scanNumber（）提取token，isDigit是用来判断字符是不是数字


    func isDigit(c byte) bool {
        return c >= '0' && c <= '9'
    }

    scanNumber简单调用scan

    func (self *Lexer) scanNumber() string {
        return self.scan(reNumber)
    }

    为了简化数字字面量和标识符的提取逻辑，scan也使用了正则表达式

    func (self *Lexer) scan(re *regexp.Regexp) string {
        if token := re.FindString(self.chunk); token != "" {
            self.next(len(token))
            return token
        }
        panic("unreachable!")
    }

    定义表示数字字面量的正则表达式

    var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)

    标识符和关键字

    if c == '_' || isLetter(c) {
        token := self.scanIdentifier()
        if kind, found := keywords[token]; found {
            return self.line, kind, token // keyword
        } else {
            return self.line, TOKEN_IDENTIFIER, token
        }
    }

    self.error("unexpected symbol near %q", c)


    如果下一个token是标识符，调用scanIdentifier提取token，取得token以后，根据它是普通的标识符还是关键字，分情况返回。isLetter判断字符是不是字母

    func isLetter(c byte) bool {
        return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
    }

    scanIdentifier也指示调用了scan

    func (self *Lexer) scanIdentifier() string {
        return self.scan(reIdentifier)
    }

    定义表示标识符的正则表达式

    var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)

    LookAhead（）

    已经定义好了结构体Lexer，实现了NextToken方法，提供了基本的词法分析功能，每次调用都会读取并返回下一个token，不过有时候并不想跳过下一个token，指示向看看下一个token是什么类型。只需要备份词法分析器的状态，然后读取下一个token，记录类型，然后恢复状态就可以了。既然读都读了，还不如缓存起来，免得做无用功。

    修改结构体

    type Lexer struct {
        chunk         string // source code
        chunkName     string // source name
        line          int    // current line number
        nextToken     string
        nextTokenKind int
        nextTokenLine int
    }

  后面的三个字段用于缓存下一个token的信息，有了它们就可以实现提前看看的功能LookAhead（）

      func (self *Lexer) LookAhead() int {
        if self.nextTokenLine > 0 {
            return self.nextTokenKind
        }
        currentLine := self.line
        line, kind, token := self.NextToken()
        self.line = currentLine
        self.nextTokenLine = line
        self.nextTokenKind = kind
        self.nextToken = token
        return kind
    }


    通过字段nextTokenLine判断缓存里面是否有下一个token信息，如果有，直接返回token的类型就可以；反之，调用NextToken提取下一个token并缓存。

    修改NextToken方法，如果发现缓存里面有下一个token信息，就直接从缓存里面读取，然后清空缓存

    func (self *Lexer) NextToken() (line, kind int, token string) {
        if self.nextTokenLine > 0 {
            line = self.nextTokenLine
            kind = self.nextTokenKind
            token = self.nextToken
            self.line = self.nextTokenLine
            self.nextTokenLine = 0
            return
        }
        ......

    }


    给结构体Lexer增加方法NextTokenOfKind，用来提取指令类型的token

    func (self *Lexer) NextTokenOfKind(kind int) (line int, token string) {
        line, _kind, token := self.NextToken()
        if kind != _kind {
            self.error("syntax error near '%s'", token)
        }
        return line, token
    }

    NextIdentifier用来提取标识符

    func (self *Lexer) NextIdentifier() (line int, token string) {
        return self.NextTokenOfKind(TOKEN_IDENTIFIER)
    }

    Line()返回当前行号

    func (self *Lexer) Line() int {
        return self.line
    }   

单元测试

    package main

    import "fmt"
    import "io/ioutil"
    import "os"
    import . "lunago/compiler/lexer"

    func main() {
        if len(os.Args) > 1 {
            data, err := ioutil.ReadFile(os.Args[1])
            if err != nil {
                panic(err)
            }

            testLexer(string(data), os.Args[1])
        }
    }

    根据命令行参数读取测试脚本，然后交给testLexer处理

    func testLexer(chunk, chunkName string) {
        lexer := NewLexer(chunk, chunkName)
        for {
            line, kind, token := lexer.NextToken()
            fmt.Printf("[%2d] [%-10s] %s\n",
                line, kindToCategory(kind), token)
            if kind == TOKEN_EOF {
                break
            }
        }
    }

    testLexer函数先根据源文件名和脚本创建Lexer结构体的实例，然后循环调用NextToken对脚本进行词法分析，打印出每个token的行号、类型和内容，直到整个甲苯都分析完成为止。kindToCategory函数把token的类型转换成更容易理解的字符串形式。

    func kindToCategory(kind int) string {
        switch {
        case kind < TOKEN_SEP_SEMI:
            return "other"
        case kind <= TOKEN_SEP_RCURLY:
            return "separator"
        case kind <= TOKEN_OP_NOT:
            return "operator"
        case kind <= TOKEN_KW_WHILE:
            return "keyword"
        case kind == TOKEN_IDENTIFIER:
            return "identifier"
        case kind == TOKEN_NUMBER:
            return "number"
        case kind == TOKEN_STRING:
            return "string"
        default:
            return "other"
        }
    }


    $ go run main.go helloworld.lua 
    [ 6] [identifier] print
    [ 6] [separator ] (
    [ 6] [string    ] hello world！！！
    [ 6] [separator ] )
    [13] [other     ] EOF


=========================================

抽象语法树

    字符可以任意组合，词法规则定义了怎样的组合可以构成合法的token，而token，也可以任意组合，语法规则定义了怎么样的组合可以构成合法的程序。

介绍
    整个源代码可以堪称是一个字符序列，词法分析阶段根据词法规则将字符序列分解为token序列，接下来的语法分析阶段根据语法规则，将token序列解析为抽象语法树ast

    与抽象语法树对应的是具体语法树（Concrete Syntax tree），也叫做解析树（Parse tree parsing tree）。cst，也是解析树，是源代码的解析结果，它完整保留了源代码里面的各种信息，不过对于编译器来说，cts包含的很多信息（比如分号，关键字等）都是多余的，这些信息在后续的编译阶段并没有太大的用处，把cst里面多余的信息去掉，仅仅留下必要的信息，化繁为简，就得到了一个ast。

    以算术表达式为例，由于与少年黁符有不同的优先级，所以在必要时，需要使用（）来改变运算符的优先级，cst会如实记录（），但是由于树状结构本身就可以隐含运算符的优先级，所以ast完全可以省略（），

    大部分编译器并不是先把源代码解析为cst，在转换为ast，而是直接生成ast。因此，cst往往只存在于概念阶段，此外，ast也不一定非要使用“树”这种数据结构。这里将直接使用结构体来拜师ast的各个节点

    计算机语言一般使用上下文无关的文法来表述，而cfg一般使用巴科斯范式或者其扩展ebnf，bnf、ebnf又可以使用语法图来表示。

    以后将结合ebnf和语法图对lua语法进行ast节点的定义。

Chunk和块

    lua中，一点完整的（可以被lua虚拟机解释执行）lua代码块成为chunk
    chunk的ebnf描述
    chunk：：=block

    在ebnf中“：：=”表示“被定义”的意思。因此，chunk就是一个代码块

    block：：= {stat} [restat]
    restat::= return [explist][';']
    exolist::= exp {',' exp}

    在ebnf中，{A}表示A可以出现任意次，0或者多次。[A]表示A可选，可以出现0或1次。因此，代码块是任意多条语句，在加上一条可选的返回语句；返回语句以关键字return后跟可选的表达式列表，以及一个可选的分号；表达式列表则是1到多个表达式，由都好分隔。


    compiler/ast/block.go

    定义结构体Block

    package ast

    type Block struct {
        LastLine int
        Stats    []Stat
        RetExps  []Exp
    }

  因为chunk实际上等同于代码块，因此只定义了Block结构体。Block结构体仅仅包含了以后要处理的必要信息，包括语句序列、返回语句里面的表达式序列等，至于关键字、分号、逗号等信息全部丢弃，这就是将其称为ast的原因。
  LastLine：记录了代码块的末尾行号，在代码生成阶段需要使用这个信息。


  语句
  在命令式编程语言中，语句statement是最基本的执行单位，表达式expression则是构成语句的元素之一。玉玦和表达式主要区别在于：语句只能执行，不能用于求值，而表达式值你能用于求值，不用单独执行

  lua中，函数调用既可以是表达式，也可以是语句。lua有15种语句，下面是lua语句的EBNF表述
  http://www.lua.org/manual/5.3/manual.html#9

  stat ::=  ‘;’ | 
         varlist ‘=’ explist | 
         functioncall | 
         label | 
         break | 
         goto Name | 
         do block end | 
         while exp do block end | 
         repeat block until exp | 
         if exp then block {elseif exp then block} [else block] end | 
         for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end | 
         for namelist in explist do block end | 
         function funcname funcbody | 
         local function Name funcbody | 
         local namelist [‘=’ explist] 

  lua语句大致可以分为控制语句、声明和赋值语句、以及其他语句三种。声明和赋值语句用于声明局部变量、给变量赋值或者向表里面写入值，包括局部变量声明语句、赋值语句、局部函数定义语句和非局部函数定义语句。控制语句用于改变执行流程，包括while、repeat、if、for、break、label、goto。其他语句包括空语句、do、函数调用。

  compiler/ast/stat.go
  定义接口Stat

    package ast


    type Stat interface{}


  1、简单语句
    空语句、break、do、函数调用、label和goto语句比较简单。定义这几种语句

    type EmptyStat struct{}              // ‘;’
    type BreakStat struct{ Line int }    // break
    type LabelStat struct{ Name string } // ‘::’ Name ‘::’
    type GotoStat struct{ Name string }  // goto Name
    type DoStat struct{ Block *Block }   // do block end
    type FuncCallStat = FuncCallExp      // functioncall

    空语句没有任何语义，仅仅起到分隔的作用，所以结构体EmptyStat也没有任何字段。
    break语句会在代码生成阶段产生一条跳转指令，所以需要记录其行号
    label和goto语句搭配使用，用来实现任意跳转，所以需要记录标签名
    do语句仅仅是为了给语句快引入新的作用域，所以结构体DoStat也只有一个字段Block
    函数调用既可以是语句，也可以是表达式，所以仅仅是起了一个别名。

  2、while repeat
    用于实现条件循环,ebnf表示

    while exp do block end
    repeat block until exp

    定义语句

    type WhileStat struct {
        Exp   Exp
        Block *Block
    }

    type RepeatStat struct {
        Block *Block
        Exp   Exp
    }

 3、if
    Ebnf:
    if exp then block {elseif exp then block} [else block] end

    为了简化st和后面的代码生成，对if语句进行改造，把最后可选的else块改为elseif块：

    if exp then block {elseif exp then block} [elseif true then block] end

    如果把elseif合并。ebnf可以简化为

    if exp then block {elseif exp then block} end

    定义if语句


    type IfStat struct {
        Exps   []Exp
        Blocks []*Block
    }


    把表达式收集到Exps总，把语句快收集到Blocks 字段中。表达式和语句块都是按照索引意义对应，索引0处是if-then表达式和块，其余索引处是elseif-then表达式和块。

4、数值for

    lua有阆中形式的for 循环：数值for和通用for循环
    数值for的ebnf：
    for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end

    数值for循环以关键字for开始，然后是标识符和等号，然后是逗号分隔的初始值、限制以及可选的步长表达式，后跟一条do语句

    定义数值for循环语句


    type ForNumStat struct {
        LineOfFor int
        LineOfDo  int
        VarName   string
        InitExp   Exp
        LimitExp  Exp
        StepExp   Exp
        Block     *Block
    }

    需要把关键字for和do所在行号记录下来，供代码生成阶段使用。


5、通用for

    for namelist in explist do block end
    namelist ::= Name {‘,’ Name}
    explist ::= exp {‘,’ exp}

    以关键字for开始，然后是逗号分隔的标识符列表，然后是关键字in和逗号分隔的表达式列表，最后是do语句

    定义通用for语句

    type ForInStat struct {
        LineOfDo int
        NameList []string
        ExpList  []Exp
        Block    *Block
    }

    需要把关键字do所在行号记录下来，供代码生成阶段使用。和if语句类似，把关键字in左边的标识符列表记录在NameList中，右侧的表达式列表记录在ExpList中。不过和if语句不通，标识符和表达式不是一一对应的。

6、局部变量声明    
    由于lua支持多重赋值，所以声明和赋值语句比较复杂，其中局部变量声明语句用于声明（并初始化）新的局部变量

    local namelist [‘=’ explist]
    namelist ::= Name {‘,’ Name}
    explist ::= exp {‘,’ exp}

    局部变量声明语句以关键字local开始，后跟逗号分隔的标识符列表，然后是可选的等号以及逗号分隔的表达式列表。

    定义局部变量的声明语句

    type LocalVarDeclStat struct {
        LastLine int
        NameList []string
        ExpList  []Exp
    }

    需要把末尾的行号记录下来，供代码生成阶段使用。和通用for循环类似，把等号左侧的标识符列表记录在NameList中，右侧的表达式列表记录在ExpList中，标识符和表达式也不是一一对应的，且表达式列表可以为空。

7、赋值语句

    赋值语句用于给变量（包括已经声明的局部变量、Upvalue或者全局变量）赋值、根据key给表赋值、或者修改记录的字段

    varlist ‘=’ explist
    varlist ::= var {‘,’ var}
    var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
    explist ::= exp {‘,’ exp}

    赋值语句被等号分成了两部分，左边是逗号分隔的var表达式列表，右边是逗啊后分隔的任意表达式列表。

    定义赋值语句

    type AssignStat struct {
        LastLine int
        VarList  []Exp
        ExpList  []Exp
    }

    和局部变量声明语句一样，需要把末尾的行号也记录下来，供代码生成阶段使用。等到左侧的var表达式列表记录在VarList中，右侧的任意表达式列表记录在ExpList中

 8、非局部函数定义语句   
    函数定义语句实际上是赋值语句的语法糖。函数定义语句又分为局部函数定义语句和非局部函数定义语句两种。非局部函数的定义：

    function ::= funcname funcbody
    funcname ::= Name {‘.’ Name} [‘:’ Name]
    funcbody ::= ‘(’ [parlist] ‘)’ block end

    非局部函数定义语句 以关键值function开始，然后是函数名（并非简单的标识符），然后是被圆括号扩起来的可选参数列表，最后是语句块以及关键字end。

    参数列表的ednf：

    parlist ::= namelist [‘,’ ‘...’] | ‘...’
    namelist ::= Name {‘,’ Name}

    参数列表可以是逗号分隔的标识符列表，后跟可选的逗号以及vararg符号，或者就单一个vararg符号。

    函数名中的“:”写法，实际上也是lua添加的语法糖，用来模拟面向对象的方法定义。比如以下三天语句在语义上是完全等价的。

    function t.a.b.c:f (params) body end          //方法定义
    function t.a.b.c.f (self, params) body end     //函数定义
    t.a.b.c.f = function (self, params) body end   //赋值

    在语法分析阶段，会把被局部函数定义语句的冒号语法糖去掉，并会把它转换为赋值语句，所以不需要给它定义专用的结构体

9、局部函数定义

    和非局部函数定义语句 类似，局部函数定义语句实际上是局部变量声明语句的语法糖

    local function Name funcbody

    局部函数定义语句以关键字local开始，后跟关键字function，然后是标识符。从圆括号开始，剩下的部分和非局部函数定义语句完全一样。

    注意，和非局部函数定义语句略有不同，为了方便递归函数的编写，局部函数定义语句会被转换为局部变量声明和赋值两条语句，比如
    local function f (params) body end
    实际上会被转化为下面两条
    local f; f = function(params) body end

    为了简化代码生成，我们部分地保留了这个语法糖。定义局部函数定义语句


    type LocalFuncDefStat struct {
        Name string
        Exp  *FuncDefExp
    }

    Name 对应函数名，Exp对应函数定义表达式

表达式
    lua有11种表达式，分为5类：字面量表达式、构造器表达式、混算符表达式、vararg表达式以及前缀表达式。字面量表达式包括nil、布尔、数字和字符串表达式。构造器表达式包括表构造器和函数构造器表达式。运算符表达式包括一元和二元运算符表达式。

    exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef | 
         prefixexp | tableconstructor | exp binop exp | unop exp   

    ast/exp.go中定义接口Exp

    package ast

    type Exp interface{}

    简单表达式

    字面量、vararg、名称表达式的定义ast/exp.go

    type NilExp struct{ Line int }    // nil
    type TrueExp struct{ Line int }   // true
    type FalseExp struct{ Line int }  // false
    type VarargExp struct{ Line int } // ...

    // Numeral
    type IntegerExp struct {
        Line int
        Val  int64
    }
    type FloatExp struct {
        Line int
        Val  float64
    }

    type StringExp struct {
        Line int
        Str  string
    }

    type NameExp struct {
        Line int
        Name string
    }

    布尔、nil、vararg表达式只需要记录行号就可以了，为了简化代码生成器，把布尔表达式进一步分成了真假两种。数字字面量表达式除了记录行号，还需要解析并记录数值。同样为了简化代码生成器，把数字直面量表达式又进一步细分为整数和赋电视两种，字符串字面量表达式需要记录行号和字符串自身，名称表达式需要记录行号和名字，也就是标识符。

  运算符表达式  

  分为一元和二元运算符表达式。定义一元运算符表达式

    type UnopExp struct {
        Line int // line of operator
        Op   int // operator
        Exp  Exp
    }

    记录了表达式以及运算符和运算符所在的行号。

    定义二元运算符表达式


    type BinopExp struct {
        Line int // line of operator
        Op   int // operator
        Exp1 Exp
        Exp2 Exp
    }

    由于拼接运算符的特殊性，需要单独定义一个结构体  

    type ConcatExp struct {
        Line int // line of last ..
        Exps []Exp
    }

    在语法分析阶段，会把连续的多个拼接操作整合在一起，这样就可以很方便的在代码生成阶段使用一条指令CONCAT优化拼接操作。

    表构造表达式

     tableconstructor ::= ‘{’ [fieldlist] ‘}’
     fieldlist ::= field {fieldsep field} [fieldsep]
     field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp
     fieldsep ::= ‘,’ | ‘;’


     表构造器由花括号、其中是可选的字段列表组成；字段列表由逗号或者分号分隔，并且末尾可以有一个可选的逗号或者分号。字段可以类似索引赋值或者变量赋值，也可以是一个简单的表达式。


    之所以这么复杂，是因为lua在李米昂放了几个字段语法糖。在{}中，k = v 完全等价于["k"] = v,单独的表达式exp则基本上等价于[n] = exp(n是整数，从1开始递增)。

        type TableConstructorExp struct {
        Line     int // line of `{` ?
        LastLine int // line of `}`
        KeyExps  []Exp
        ValExps  []Exp
    }

    需要记录{}所在的行号，供以后代码生成阶段使用。在语法分析阶段，会把字段的语法糖去掉，这样就可以统一把key value表达式整理到KeyExps和ValExps中

函数定义表达式
    函数定义表达式也叫做函数构造器，其求值的结果是一个匿名函数

    functiondef ::= function funcbody
    funcbody ::= ‘(’ [parlist] ‘)’ block end
    parlist ::= namelist [‘,’ ‘...’] | ‘...’
    namelist ::= Name {‘,’ Name}

    函数定义表达式和函数定义语句很像，区别在于前者省略了函数名，定义函数定义表达式

    type FuncDefExp struct {
        Line     int
        LastLine int // line of `end`
        ParList  []string
        IsVararg bool
        Block    *Block
    }

    需要记录{}所在的行号，供代码生成阶段使用。

前缀表达式
    prefixexp ::= var | functioncall | ‘(’ exp ‘)’
    var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name 
    functioncall ::=  prefixexp args | prefixexp ‘:’ Name args 

    前缀表达式包括var表达式、函数调用表达式、和圆括号表达式三种。var表达式是能够出现在赋值语句等号左侧的表达式，又包括名称表达式、表访问表达式和记录访问表达式三种。

    为了更直观的理解前缀表达式，对上面的写法进行改造

  
    prefixexp ::= Name |
                  ‘(’ exp ‘)’ |
                  prefixexp ‘[’ exp ‘]’ |
                  prefixexp ‘.’ Name |
                  prefixexp ‘:’ Name args |
                  prefixexp args
    
    可见，之所以叫做“前缀”表达式，就是因为可以作为表访问表达式、记录表达式和函数调用表达式的前缀。
    
    名称表达式茜米昂已经定义了。记录访问表达式其实是表访问表达式的语法糖，在语义上，t.k完全等价于t["k"]。在语法分析阶段会把记录访问表达式转换成表访问表达式，所以没有必要专门定义记录访问表达式。              
圆括号表达式

    圆括号表达式主要有两个用途：改变运算符的优先级或者结合性，或者在多重赋值的时候将vararg和函数调用的返回值固定为1.运算符的优先级可以隐含在ast中，所以可以在语法分析阶段扔掉这部分的圆括号。但是对于其他情况，为了简化语法分析和代码生成阶段，还是需要保留圆括号。

    定义圆括号表达式

    type ParensExp struct {
        Exp Exp
    }

表访问表达式
    定义表访问表达式

    type TableAccessExp struct {
        LastLine  int // line of `]` ?
        PrefixExp Exp
        KeyExp    Exp
    }

    把有方括号所在的行号记录在字段LastLine中，供词法分析阶段使用。

函数调用表达式

    functioncall ::= prefixexp [‘:’ Name] args
    args ::= ‘(’ [explist] ‘)’ | tableconstructor | LiteralString 

    函数调用表达式以前缀表达式开始，后跟可选的冒号以及标识符，然后是表构造器、字符串字面量或者由圆括号包围起来的可选参数列表。

    lua在函数调用表达式里面加了两个语法糖：第一个允许在有且仅有一个参数，且参数是字符串字面量或者表构造器的时候省略圆括号；第二个是为了模拟面向对象语言里面的方法调用而加入的。在语义上，v:name(args)完全等价于v.name(v, args).

    定义函数调用表达式

    type FuncCallExp struct {
        Line      int // line of `(` ?
        LastLine  int // line of ')'
        PrefixExp Exp
        NameExp   *StringExp
        Args      []Exp
    }


    圆括号所在的行号记录在Line和LastLine中，着两个行号在代码生成阶段会用到。虽然方法调用时语法糖，不过lua编译器会针对方法调用生成self指令。所以，需要把这个语法糖保留起来，记录在NameExp中。省略圆括号的语法糖则可以完全去掉，这样就可以把参数表统一起来。记录在Args中

  ==============
  
  语法分析  语法分析器

  说明
  计算机语言一般使用上下文无关文法cfg描述，成折射语言为上下文无关语言。cfg一般使用巴科斯范式bnf或者其扩展形式ebnf书写，语法分析器的作用就是按照某种语言的bnf描述，将这种语言的源代码转换成抽象语法树ast，供后续阶段使用。

  歧义
  假定某种语言L事上下文无关语言，那么可以把用L语言编写（也符合L语言的语法规则）的源代码转换成解析树cst。对于任何一段L源代码，如果只能被转换成唯一的一棵cst，那么就称L语言无歧义，反之，成L语言有歧义

  关于歧义，其中最著名的例子就是c中的悬挂else问题，如下时c语言的if-else语句的ebnf的相关部分

    stat ::= ... 其他部分
            | if '(' exp ')' stat [else stat]

    比如下面这条语句：        

    if (a) if(b) s1;else s2;

    如果不佳其他的限制，就会产生歧义，可以转换成两颗cst。

    语法规则中存在的歧义必须痛殴其他规则去除，否则语法分析器遇到有歧义的代码就会不知所措。C语言规定else和离它最近的那个if关联，因此前面的那条例子语句会被转换成一个cst，就解决了歧义和悬挂else的问题。

    Lua的if语句语法规则更加严格，所以没有歧义。不过lua语法中仍然存在有歧义的地方，就是一元和二元运算符表达式。以二元运算符表达式为例

    exp binop exp

    a + b * c 应该解释为 （a + b）* c,还是a + （b * c）?可以通过运算符优先级和结合律解决歧义。

前瞻 以及 回顾

    如果手动把一段lua代码转换为语法树，会怎么做？

    首先会扮演词法分析器的角色，从源代码中提取出一个token，然后切换到语法分析器的角色，根据这个token，看看下一步该干嘛。比如这个token是if，就会尝试解析一个if语句，然后继续，如果是for，就会尝试解析for循环语句，然后继续。

    可是如果拿到的token是local该怎么处理？
    是尝试解析局部变量声明语句  还是局部函数定义语句？
    没有办法了，只好从源代码中在多提取一个token，如果它是关键字function，就尝试解析局部函数定义语句，否则就尝试解析局部变量声明语句。

    可是有的时候，即使提取再多的token也没用，比如赋值语句还有函数调用语下面是着两种语句相关语法的ebnf

    stat ::= ......
            | varlist '=' explist
            | functioncall
    varlist ::= varlist ::= var {‘,’ var}

    var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name 
    prefixexp ::= var | functioncall | ‘(’ exp ‘)’
    functioncall ::=  prefixexp args | prefixexp ‘:’ Name args 

    如果已经通过下一个token排除了其他的可能，值徐亚尝试解析赋值语句或者函数调用语句，可问题是，这两种语句都以前缀表达式开始，而前缀表达式有可以是任意长度，所以没有办法知道到底需要提取多少token才能最终决定要解析哪种语句，那该咋办？

    也不是没有办法，可以先记录源代码已经分析到了哪里，然后尝试解析额赋值语句，如果运气好，真的就解析出了一条赋值语句，就继续后面的工作即可；如果运气不好，那当然就解析失败了，回到刚才的位置，重新尝试解析额函数调用语句就可以了。

    像这种痛殴预先读取后面的几个token来决定下一步解析策略的方法叫做 前瞻 lookahead

    前瞻失败后，记录张泰惊醒尝试，并可能会退的做法叫做 回溯 backtracking

    如果上下无关语言L不需要借助回溯就可以完成解析，那么就成L为确定性（deterministic）语言。确定性语言一定没有歧义。

    回溯 最大的问题就是会导致解析起无法在线性时间内完成工作，因此要尽可能的避免。

    下面会使用一些简单的技巧来解决lua语法规则里面的不确定性，从而避免回溯


解析方式
    显示世界中的树是向上生长的，树根在地底下，不不过计算机科学里面的“树”正好相反，在最上面。前面说了，如何把牟总语言的源代码转换成语法树，使用这种方法，需要先构造根节点，然后是非叶子结点，最后是叶子结点，如此往复，直到构造出一棵完整的语法树。由于根在上面，所以这种方式叫做 自顶向下法Top-down.

    实际上，还有另一种方法，可以先构造叶子结点，然后是非叶子结点，最后是根节点，这样也可以构造出语法树，这种方式叫做 自底向上法 bottom-up

    自底向上的解析器包括LR解析器和CYK解析器等，自顶向下的解析器包括LL解析器和递归下降解析器（Recursive Descent Parser）等。

    编写Lua递归下降解析器

解析块    
    递归下降解析器曹勇自顶向下的方式进行解析，所以直接从ast的根节点入手。由于lua脚本实际上就是一个代码块，所以解析结果应该是一个Block结构体的实例。

    compiler/parser/parse_block.go

    定义函数 parseBlock（）

    package parser

    import . "lunago/compiler/ast"
    import . "lunago/compiler/lexer"

    // block ::= {stat} [retstat]
    func parseBlock(lexer *Lexer) *Block {
        return &Block{
            Stats:    parseStats(lexer),
            RetExps:  parseRetExps(lexer),
            LastLine: lexer.Line(),
        }
    }

    创建结构体Block的实例，调用parseBlock（）解析语句序列，调用parseRetExps（）解析额可选的返回语句，并记录末尾的行号，

    parseStats（）

    func parseStats(lexer *Lexer) []Stat {
        stats := make([]Stat, 0, 8)
        for !_isReturnOrBlockEnd(lexer.LookAhead()) {
            stat := parseStat(lexer)
            if _, ok := stat.(*EmptyStat); !ok {
                stats = append(stats, stat)
            }
        }
        return stats
    }


    循环调用parseStat（）解析语句，知道痛殴前瞻看到关键字return或者发现块已经结束为止。
    那么如何知道块在什么地方结束呢？
    只要把所有和块相关的语法规则都列出来，就可以找到规律

    block ::= {stat} [retstat]
    stat ::=  ‘;’ | 
         varlist ‘=’ explist | 
         functioncall | 
         label | 
         break | 
         goto Name | 


         do block end | 
         while exp do block end | 
         repeat block until exp | 
         if exp then block {elseif exp then block} [else block] end | 
         for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end | 
         for namelist in explist do block end | 
         function funcname funcbody | 
         local function Name funcbody | 



         local namelist [‘=’ explist] 
 
   可见，如果出现在其他语句里面，快后面的值你鞥是关键字end、else、esleif或者until。如果作为单独的chunk，并且没有返回语句，那么块的后面不能再有任何非空白token。
   把这些判断放到函数_isReturnOrBlockEnd里：

   func _isReturnOrBlockEnd(tokenKind int) bool {
        switch tokenKind {
        case TOKEN_KW_RETURN, TOKEN_EOF, TOKEN_KW_END,
            TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
            return true
        }
        return false
    }

    parseRetExps函数

    // retstat ::= return [explist] [‘;’]
    // explist ::= exp {‘,’ exp}
    func parseRetExps(lexer *Lexer) []Exp {
        if lexer.LookAhead() != TOKEN_KW_RETURN {
            return nil
        }

        lexer.NextToken()
        switch lexer.LookAhead() {
        case TOKEN_EOF, TOKEN_KW_END,
            TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
            return []Exp{}
        case TOKEN_SEP_SEMI:
            lexer.NextToken()
            return []Exp{}
        default:
            exps := parseExpList(lexer)
            if lexer.LookAhead() == TOKEN_SEP_SEMI {
                lexer.NextToken()
            }
            return exps
        }
    }   

    通过词法分析器前瞻下一个token，如果不是return，说明没有返回语句，直接返回nil；否则掉过关键值return，前瞻下一个token，如果返现块已经结束或者分号，那么返回教育没有任何表达式，跳过分号，如果有的话，返回空的Exp列表就可以了，否则调用parseExpList()函数解析表达式序列，并跳过可选的分号

    compiler/parser/parse_exp.go


    parseExpList函数

    package parser

    import . "lunago/compiler/ast"
    import . "lunago/compiler/lexer"
    import "lunago/number"

    // explist ::= exp {‘,’ exp}
    func parseExpList(lexer *Lexer) []Exp {
        exps := make([]Exp, 0, 4)
        exps = append(exps, parseExp(lexer))
        for lexer.LookAhead() == TOKEN_SEP_COMMA {
            lexer.NextToken()
            exps = append(exps, parseExp(lexer))
        }
        return exps
    }

解析语句

    lua有15中语句

    stat ::=  ‘;’ | 
         varlist ‘=’ explist | 
         functioncall | 
         label | 
         break | 
         goto Name | 
         do block end | 
         while exp do block end | 
         repeat block until exp | 
         if exp then block {elseif exp then block} [else block] end | 
         for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end | 
         for namelist in explist do block end | 
         function funcname funcbody | 
         local function Name funcbody | 
         local namelist [‘=’ explist] 


  通过前瞻一个token，可以锁定其中的13种语句，对于局部变量声明和局部函数定义语句，需要前瞻2个token才能确定要解析哪一种；对于for循环语句，则需要前瞻3个token才能确定解析数值for循环 还是通用for循环。
  剩下的函数调用和赋值两种，只能想别的办法了。

  compiler/parser/parse_stat.go

  定义函数

    package parser

    import . "lunago/compiler/ast"
    import . "lunago/compiler/lexer"

    /*
    stat ::=  ‘;’
        | break
        | ‘::’ Name ‘::’
        | goto Name
        | do block end
        | while exp do block end
        | repeat block until exp
        | if exp then block {elseif exp then block} [else block] end
        | for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
        | for namelist in explist do block end
        | function funcname funcbody
        | local function Name funcbody
        | local namelist [‘=’ explist]
        | varlist ‘=’ explist
        | functioncall
    */
    func parseStat(lexer *Lexer) Stat {
        switch lexer.LookAhead() {
        case TOKEN_SEP_SEMI:
            return parseEmptyStat(lexer)
        case TOKEN_KW_BREAK:
            return parseBreakStat(lexer)
        case TOKEN_SEP_LABEL:
            return parseLabelStat(lexer)
        case TOKEN_KW_GOTO:
            return parseGotoStat(lexer)
        case TOKEN_KW_DO:
            return parseDoStat(lexer)
        case TOKEN_KW_WHILE:
            return parseWhileStat(lexer)
        case TOKEN_KW_REPEAT:
            return parseRepeatStat(lexer)
        case TOKEN_KW_IF:
            return parseIfStat(lexer)
        case TOKEN_KW_FOR:
            return parseForStat(lexer)
        case TOKEN_KW_FUNCTION:
            return parseFuncDefStat(lexer)
        case TOKEN_KW_LOCAL:
            return parseLocalAssignOrFuncDefStat(lexer)
        default:
            return parseAssignOrFuncCallStat(lexer)
        }
    }


    前瞻一个token，然后根据类型调用对应的函数解析语句

简单语句
    空语句、break、label、goto、do、while、repeat，这几种语句的解析函数很简单，对于空语句，跳过分号就可以可

    // ;
    func parseEmptyStat(lexer *Lexer) *EmptyStat {
        lexer.NextTokenOfKind(TOKEN_SEP_SEMI)
        return _statEmpty
    }

    对于break语句，则掉过关键字并记录行号就可以了

     // break
    func parseBreakStat(lexer *Lexer) *BreakStat {
        lexer.NextTokenOfKind(TOKEN_KW_BREAK)
        return &BreakStat{lexer.Line()}
    }   

    对于label语句，跳过分隔符并记录标签名就可以了 

    // ‘::’ Name ‘::’
    func parseLabelStat(lexer *Lexer) *LabelStat {
        lexer.NextTokenOfKind(TOKEN_SEP_LABEL) // ::
        _, name := lexer.NextIdentifier()      // name
        lexer.NextTokenOfKind(TOKEN_SEP_LABEL) // ::
        return &LabelStat{name}
    }

    对于goto语句，跳过关键字并记录标签名就可以

    // goto Name
    func parseGotoStat(lexer *Lexer) *GotoStat {
        lexer.NextTokenOfKind(TOKEN_KW_GOTO) // goto
        _, name := lexer.NextIdentifier()    // name
        return &GotoStat{name}
    }

    对于do语句，先跳过关键字do，然后调用函数parseBlock（）解析额块，最后跳过关键字end

    // do block end
    func parseDoStat(lexer *Lexer) *DoStat {
        lexer.NextTokenOfKind(TOKEN_KW_DO)  // do
        block := parseBlock(lexer)          // block
        lexer.NextTokenOfKind(TOKEN_KW_END) // end
        return &DoStat{block}
    }

    while语句解析

    // while exp do block end
    func parseWhileStat(lexer *Lexer) *WhileStat {
        lexer.NextTokenOfKind(TOKEN_KW_WHILE) // while
        exp := parseExp(lexer)                // exp
        lexer.NextTokenOfKind(TOKEN_KW_DO)    // do
        block := parseBlock(lexer)            // block
        lexer.NextTokenOfKind(TOKEN_KW_END)   // end
        return &WhileStat{exp, block}
    }

    repeat语句

    // repeat block until exp
    func parseRepeatStat(lexer *Lexer) *RepeatStat {
        lexer.NextTokenOfKind(TOKEN_KW_REPEAT) // repeat
        block := parseBlock(lexer)             // block
        lexer.NextTokenOfKind(TOKEN_KW_UNTIL)  // until
        exp := parseExp(lexer)                 // exp
        return &RepeatStat{block, exp}
    }

    需要注意的是，parseBlock函数最后可能会调用到parseDotStat函数，而parseDotStat函数等又会调用到
parseBlock函数，所以块和语句解析额函数之间存在递归调用的关系，递归下降解析器因此得名。

if语句    

    定义if语句的解析函数

    // if exp then block {elseif exp then block} [else block] end
    func parseIfStat(lexer *Lexer) *IfStat {
        exps := make([]Exp, 0, 4)
        blocks := make([]*Block, 0, 4)

        lexer.NextTokenOfKind(TOKEN_KW_IF)         // if
        exps = append(exps, parseExp(lexer))       // exp
        lexer.NextTokenOfKind(TOKEN_KW_THEN)       // then
        blocks = append(blocks, parseBlock(lexer)) // block

        for lexer.LookAhead() == TOKEN_KW_ELSEIF {
            lexer.NextToken()                          // elseif
            exps = append(exps, parseExp(lexer))       // exp
            lexer.NextTokenOfKind(TOKEN_KW_THEN)       // then
            blocks = append(blocks, parseBlock(lexer)) // block
        }

        // else block => elseif true then block
        if lexer.LookAhead() == TOKEN_KW_ELSE {
            lexer.NextToken()                           // else
            exps = append(exps, &TrueExp{lexer.Line()}) //
            blocks = append(blocks, parseBlock(lexer))  // block
        }

        lexer.NextTokenOfKind(TOKEN_KW_END) // end
        return &IfStat{exps, blocks}
    }

    按照步骤解析额每个语法元素就可以了，需要说明的是，为了简化代码生成器，把最后可选的else块换成了elseif块

for循环语句    

    lua有两种for循环：数值for和通用for循环。

    定义for循环语句的解析函数

    // for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
    // for namelist in explist do block end
    func parseForStat(lexer *Lexer) Stat {
        lineOfFor, _ := lexer.NextTokenOfKind(TOKEN_KW_FOR)
        _, name := lexer.NextIdentifier()
        if lexer.LookAhead() == TOKEN_OP_ASSIGN {
            return _finishForNumStat(lexer, lineOfFor, name)
        } else {
            return _finishForInStat(lexer, name)
        }
    }

    因为两种for循环都以关键值for开始，后跟一个标识符，所以跳过关键字for以后，还需要前瞻两个token，才能判断到底是哪种for循环，不过为了简化词法分析器，尽量通过只前瞻一个token来完成解析。这里采用了不同的做法，先跳过关键字for，然后提取标识符，然后前瞻一个token，如果是等号，就按照数值for循环解析，否则按照通用for循环解析。

    数值for循环的解析函数：

    // for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
    func _finishForNumStat(lexer *Lexer, lineOfFor int, varName string) *ForNumStat {
        lexer.NextTokenOfKind(TOKEN_OP_ASSIGN) // for name =
        initExp := parseExp(lexer)             // exp
        lexer.NextTokenOfKind(TOKEN_SEP_COMMA) // ,
        limitExp := parseExp(lexer)            // exp

        var stepExp Exp
        if lexer.LookAhead() == TOKEN_SEP_COMMA {
            lexer.NextToken()         // ,
            stepExp = parseExp(lexer) // exp
        } else {
            stepExp = &IntegerExp{lexer.Line(), 1}
        }

        lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO) // do
        block := parseBlock(lexer)                        // block
        lexer.NextTokenOfKind(TOKEN_KW_END)               // end

        return &ForNumStat{lineOfFor, lineOfDo,
            varName, initExp, limitExp, stepExp, block}
    }

    关键字for和标识符已经读取，直接从等号开始解析就可以了，需要说明的是，为了简化代码生成器，给不上不上了默认值1.

    通用for循环解析函数

    // for namelist in explist do block end
    // namelist ::= Name {‘,’ Name}
    // explist ::= exp {‘,’ exp}
    func _finishForInStat(lexer *Lexer, name0 string) *ForInStat {
        nameList := _finishNameList(lexer, name0)         // for namelist
        lexer.NextTokenOfKind(TOKEN_KW_IN)                // in
        expList := parseExpList(lexer)                    // explist
        lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO) // do
        block := parseBlock(lexer)                        // block
        lexer.NextTokenOfKind(TOKEN_KW_END)               // end
        return &ForInStat{lineOfDo, nameList, expList, block}
    }

    关键字for和第一个标识符已经读取，继续吧便是福列表解析完成，然后按照步骤解析其他语法元素就可以了。

    // namelist ::= Name {‘,’ Name}
    func _finishNameList(lexer *Lexer, name0 string) []string {
        names := []string{name0}
        for lexer.LookAhead() == TOKEN_SEP_COMMA {
            lexer.NextToken()                 // ,
            _, name := lexer.NextIdentifier() // Name
            names = append(names, name)
        }
        return names
    }

局部变量声明以及函数定义语句 

  局部变量声明以及函数定义语句都以关键字local开始，下面是着两句语的解析代码  

    // local function Name funcbody
    // local namelist [‘=’ explist]
    func parseLocalAssignOrFuncDefStat(lexer *Lexer) Stat {
        lexer.NextTokenOfKind(TOKEN_KW_LOCAL)
        if lexer.LookAhead() == TOKEN_KW_FUNCTION {
            return _finishLocalFuncDefStat(lexer)
        } else {
            return _finishLocalVarDeclStat(lexer)
        }
    }

    跳过local关键值，然后前瞻一个token，如果是是关键字funciton，就解析局部函数定义语句，否则解析局部变量声明语句。
    局部函数定义语句的解析函数：

    func _finishLocalFuncDefStat(lexer *Lexer) *LocalFuncDefStat {
        lexer.NextTokenOfKind(TOKEN_KW_FUNCTION) // local function
        _, name := lexer.NextIdentifier()        // name
        fdExp := parseFuncDefExp(lexer)          // funcbody
        return &LocalFuncDefStat{name, fdExp}
    }

    跳过关键字function，读取标识符，剩下的工作交给parseFuncDefExp完成，等会在说它。

    局部变量声明语句的解析函数

    // local namelist [‘=’ explist]
    func _finishLocalVarDeclStat(lexer *Lexer) *LocalVarDeclStat {
        _, name0 := lexer.NextIdentifier()        // local Name
        nameList := _finishNameList(lexer, name0) // { , Name }
        var expList []Exp = nil
        if lexer.LookAhead() == TOKEN_OP_ASSIGN {
            lexer.NextToken()             // ==
            expList = parseExpList(lexer) // explist
        }
        lastLine := lexer.Line()
        return &LocalVarDeclStat{lastLine, nameList, expList}
    }

赋值和函数调用语句
    赋值和函数调用语句都以前缀表达式开始，而前缀表达式又是任意长度，所以需要有栈栈无限个token的能力才能区分着两种语句，或者借助回溯来解析。不过分析这两种语句的语法规则后，不难看到，函数调用既可以是语句，也可以是前缀表达式，单一定不是var表达式，据此，可以先解析一个前缀表达式，然后看看它是不是函数调用，如果是，那么解析出来的实际上就是一条函数调用语句；反之，解析出来的必须是一个var表达式，继续解析剩余的复制语句就可以了。

    parseAssignOrFuncCallStat（）函数的定义如下

    // varlist ‘=’ explist
    // functioncall
    func parseAssignOrFuncCallStat(lexer *Lexer) Stat {
        prefixExp := parsePrefixExp(lexer)
        if fc, ok := prefixExp.(*FuncCallExp); ok {
            return fc
        } else {
            return parseAssignStat(lexer, prefixExp)
        }
    }


    赋值语句解析函数

    // varlist ‘=’ explist |
    func parseAssignStat(lexer *Lexer, var0 Exp) *AssignStat {
        varList := _finishVarList(lexer, var0) // varlist
        lexer.NextTokenOfKind(TOKEN_OP_ASSIGN) // =
        expList := parseExpList(lexer)         // explist
        lastLine := lexer.Line()
        return &AssignStat{lastLine, varList, expList}
    }

    _finishVarList代码如下：

    // varlist ::= var {‘,’ var}
    func _finishVarList(lexer *Lexer, var0 Exp) []Exp {
        vars := []Exp{_checkVar(lexer, var0)}      // var
        for lexer.LookAhead() == TOKEN_SEP_COMMA { // {
            lexer.NextToken()                          // ,
            exp := parsePrefixExp(lexer)               // var
            vars = append(vars, _checkVar(lexer, exp)) //
        } // }
        return vars
    }


    通过_checkVar确保解析出来的都是var表达式，否则借助词法分析器报错。

    // var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
    func _checkVar(lexer *Lexer, exp Exp) Exp {
        switch exp.(type) {
        case *NameExp, *TableAccessExp:
            return exp
        }
        lexer.NextTokenOfKind(-1) // trigger error
        panic("unreachable!")
    }

非局部函数定义语句

    非局部函数定义语句的解析代码

    // function funcname funcbody
    // funcname ::= Name {‘.’ Name} [‘:’ Name]
    // funcbody ::= ‘(’ [parlist] ‘)’ block end
    // parlist ::= namelist [‘,’ ‘...’] | ‘...’
    // namelist ::= Name {‘,’ Name}
    func parseFuncDefStat(lexer *Lexer) *AssignStat {
        lexer.NextTokenOfKind(TOKEN_KW_FUNCTION) // function
        fnExp, hasColon := _parseFuncName(lexer) // funcname
        fdExp := parseFuncDefExp(lexer)          // funcbody
        if hasColon {                            // insert self
            fdExp.ParList = append(fdExp.ParList, "")
            copy(fdExp.ParList[1:], fdExp.ParList)
            fdExp.ParList[0] = "self"
        }

        return &AssignStat{
            LastLine: fdExp.Line,
            VarList:  []Exp{fnExp},
            ExpList:  []Exp{fdExp},
        }
    }

    需要说明的是：
    1、去掉了冒号语法糖，在参数列表中插入了self参数
    2、会把非局部函数定义语句转换成了赋值语句

    _parseFuncName的代码如下：

    // funcname ::= Name {‘.’ Name} [‘:’ Name]
    func _parseFuncName(lexer *Lexer) (exp Exp, hasColon bool) {
        line, name := lexer.NextIdentifier()
        exp = &NameExp{line, name}

        for lexer.LookAhead() == TOKEN_SEP_DOT {
            lexer.NextToken()
            line, name := lexer.NextIdentifier()
            idx := &StringExp{line, name}
            exp = &TableAccessExp{line, exp, idx}
        }
        if lexer.LookAhead() == TOKEN_SEP_COLON {
            lexer.NextToken()
            line, name := lexer.NextIdentifier()
            idx := &StringExp{line, name}
            exp = &TableAccessExp{line, exp, idx}
            hasColon = true
        }

        return
    }


    去掉冒号语法糖以后，函数名实际上就是一蹿记录访问表达式。记录访问也是语法糖，去掉以后就是一段表索引访问表达式。

解析表达式    
    
    之前说了怎么使用前瞻和一些技巧来解析语句，不过由于运算符表达式语法存在歧义，所以这个方法并不能直接用来解析表达式，为了套用这个方法，必须借助语义规则，把表达式语法中存在的歧义消除掉，这个语义规则就是lua运算符的优先级和结合性。

    根据运算符和优先级把表达式的语法规则重写

    exp   ::= exp12
    exp12 ::= exp11 {or exp11}
    exp11 ::= exp10 {and exp10}
    exp10 ::= exp9 {(‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’) exp9}
    exp9  ::= exp8 {‘|’ exp8}
    exp8  ::= exp7 {‘~’ exp7}
    exp7  ::= exp6 {‘&’ exp6}
    exp6  ::= exp5 {(‘<<’ | ‘>>’) exp5}
    exp5  ::= exp4 {‘..’ exp4}
    exp4  ::= exp3 {(‘+’ | ‘-’) exp3}
    exp3  ::= exp2 {(‘*’ | ‘/’ | ‘//’ | ‘%’) exp2}
    exp2  ::= {(‘not’ | ‘#’ | ‘-’ | ‘~’)} exp1
    exp1  ::= exp0 {‘^’ exp2}
    exp0  ::= nil | false | true | Numeral | LiteralString
            | ‘...’ | functiondef | prefixexp | tableconstructor

    有了重写的语法规则，表达式解析函数也就水到渠成了。

    compiler/parser/parse_exp.go定义parseExp（）函数

    func parseExp(lexer *Lexer) Exp {
        return parseExp12(lexer)
    }

    因为运算符分12级，所以需要编写parseExp12到parseExp1 这12个函数。


运算符表达式

    逻辑或表达式的解析函数

    // x or y
    func parseExp12(lexer *Lexer) Exp {
        exp := parseExp11(lexer)
        for lexer.LookAhead() == TOKEN_OP_OR {
            line, op, _ := lexer.NextToken()
            lor := &BinopExp{line, op, exp, parseExp11(lexer)}
            exp = optimizeLogicalOr(lor)
        }
        return exp
    }


    在二元运算符中，只有拼接“..” 和乘方“^”具有右结合性，其他均具有左结合性。由于逻辑或运算符具有左结合性，所以在循环里面调用了函数
    parseExp11（）解析更好优先级的运算符表达式。比如a or b or c,解析后的ast如下

            or
         or          
      a     b     c

    对比一下，看看成方混算符表达式的解析函数

    /**
     * x ^ y exp0 {'^' exp2}
     */
    func parseExp1(lexer *Lexer) Exp { // pow is right associative
        exp := parseExp0(lexer)
        if lexer.LookAhead() == TOKEN_OP_POW {
            line, op, _ := lexer.NextToken()
            exp = &BinopExp{line, op, exp, parseExp2(lexer)}
        }
        return optimizePow(exp)
    }

    因为乘方运算符具有右结合性，所以函数parseExp1（）递归调用自己，解析后面的乘方运算符表达式。
    以a ^ b ^ c为例，解析后的ast如下

                ^
                  ^
             a   b  c     


    一元运算符的解析函数

    /**
     * unary {(‘not’ | ‘#’ | ‘-’ | ‘~’)} exp1
     */
    func parseExp2(lexer *Lexer) Exp {
        switch lexer.LookAhead() {
        case TOKEN_OP_UNM, TOKEN_OP_BNOT, TOKEN_OP_LEN, TOKEN_OP_NOT:
            line, op, _ := lexer.NextToken()
            exp := &UnopExp{line, op, parseExp2(lexer)}
            return optimizeUnaryOp(exp)
        }
        return parseExp1(lexer)
    }


    可以说一元运算符也具有右结合性，所以parseExp2需要调用自身解析后面的一段怒算符表达式。

    拼接运算符表达式的解析函数

    /**
     * a .. b exp4 {‘..’ exp4}
     */
    func parseExp5(lexer *Lexer) Exp {
        exp := parseExp4(lexer)
        if lexer.LookAhead() != TOKEN_OP_CONCAT {
            return exp
        }

        line := 0
        exps := []Exp{exp}
        for lexer.LookAhead() == TOKEN_OP_CONCAT {
            line, _, _ = lexer.NextToken()
            exps = append(exps, parseExp4(lexer))
        }
        return &ConcatExp{line, exps}
    }

    虽然凭借运算符也具有右结合性，不过由于其对应的lua虚拟机指令CONCAT比较特别，所以对它进行了特殊的处理。对于拼接运算符表达式，解析生成的并不是二叉树，而是多叉树，比如a.. b .. c,解析后ast如下

                ..

            a    b    c
            
     非运算符表达式           

     运算符表达式之外额其他表达式由parseExp0解析

     func parseExp0(lexer *Lexer) Exp {
        switch lexer.LookAhead() {
        case TOKEN_VARARG: // ...
            line, _, _ := lexer.NextToken()
            return &VarargExp{line}
        case TOKEN_KW_NIL: // nil
            line, _, _ := lexer.NextToken()
            return &NilExp{line}
        case TOKEN_KW_TRUE: // true
            line, _, _ := lexer.NextToken()
            return &TrueExp{line}
        case TOKEN_KW_FALSE: // false
            line, _, _ := lexer.NextToken()
            return &FalseExp{line}
        case TOKEN_STRING: // LiteralString
            line, _, token := lexer.NextToken()
            return &StringExp{line, token}
        case TOKEN_NUMBER: // Numeral
            return parseNumberExp(lexer)
        case TOKEN_SEP_LCURLY: // tableconstructor
            return parseTableConstructorExp(lexer)
        case TOKEN_KW_FUNCTION: // functiondef
            lexer.NextToken()
            return parseFuncDefExp(lexer)
        default: // prefixexp
            return parsePrefixExp(lexer)
        }
    }


    和语句类似，前瞻一个token来决定具体要解析哪种表达式，由于vararg和非数字字面量表达式较为简单，所以直接写在了case中。

    数字字面量表达式的解析函数：

    func parseNumberExp(lexer *Lexer) Exp {
        line, _, token := lexer.NextToken()
        if i, ok := number.ParseInteger(token); ok {
            return &IntegerExp{line, i}
        } else if f, ok := number.ParseFloat(token); ok {
            return &FloatExp{line, f}
        } else { // todo
            panic("not a number: " + token)
        }
    }    

函数定义表达式    
    
    函数定义表达式的解析函数

    // functiondef ::= function funcbody
    // funcbody ::= ‘(’ [parlist] ‘)’ block end
    func parseFuncDefExp(lexer *Lexer) *FuncDefExp {
        line := lexer.Line()                               // function
        lexer.NextTokenOfKind(TOKEN_SEP_LPAREN)            // (
        parList, isVararg := _parseParList(lexer)          // [parlist]
        lexer.NextTokenOfKind(TOKEN_SEP_RPAREN)            // )
        block := parseBlock(lexer)                         // block
        lastLine, _ := lexer.NextTokenOfKind(TOKEN_KW_END) // end
        return &FuncDefExp{line, lastLine, parList, isVararg, block}
    }


    为了在其他地方可以重用这个方法，跳过了关键字function，只解析函数定义表达式的其他部分。可选的参数列表由函数 _parseParList 解析：

    // [parlist]
    // parlist ::= namelist [‘,’ ‘...’] | ‘...’
    func _parseParList(lexer *Lexer) (names []string, isVararg bool) {
        switch lexer.LookAhead() {
        case TOKEN_SEP_RPAREN:
            return nil, false
        case TOKEN_VARARG:
            lexer.NextToken()
            return nil, true
        }

        _, name := lexer.NextIdentifier()
        names = append(names, name)
        for lexer.LookAhead() == TOKEN_SEP_COMMA {
            lexer.NextToken()
            if lexer.LookAhead() == TOKEN_IDENTIFIER {
                _, name := lexer.NextIdentifier()
                names = append(names, name)
            } else {
                lexer.NextTokenOfKind(TOKEN_VARARG)
                isVararg = true
                break
            }
        }
        return
    }

表构造表达式   

    // tableconstructor ::= ‘{’ [fieldlist] ‘}’
    func parseTableConstructorExp(lexer *Lexer) *TableConstructorExp {
        line := lexer.Line()
        lexer.NextTokenOfKind(TOKEN_SEP_LCURLY)    // {
        keyExps, valExps := _parseFieldList(lexer) // [fieldlist]
        lexer.NextTokenOfKind(TOKEN_SEP_RCURLY)    // }
        lastLine := lexer.Line()
        return &TableConstructorExp{line, lastLine, keyExps, valExps}
    }


    可选的字段列表由函数 _parseFieldList 解析

    // fieldlist ::= field {fieldsep field} [fieldsep]
    func _parseFieldList(lexer *Lexer) (ks, vs []Exp) {
        if lexer.LookAhead() != TOKEN_SEP_RCURLY {
            k, v := _parseField(lexer)
            ks = append(ks, k)
            vs = append(vs, v)

            for _isFieldSep(lexer.LookAhead()) {
                lexer.NextToken()
                if lexer.LookAhead() != TOKEN_SEP_RCURLY {
                    k, v := _parseField(lexer)
                    ks = append(ks, k)
                    vs = append(vs, v)
                } else {
                    break
                }
            }
        }
        return
    }

    因为字段列表的末尾允许有可以可选的分隔符，所以有点麻烦。
    字段分隔符可以是逗号或者分号

    // fieldsep ::= ‘,’ | ‘;’
    func _isFieldSep(tokenKind int) bool {
        return tokenKind == TOKEN_SEP_COMMA || tokenKind == TOKEN_SEP_SEMI
    }  


    字段由函数 _parseField 解析

    // field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp
    func _parseField(lexer *Lexer) (k, v Exp) {
        if lexer.LookAhead() == TOKEN_SEP_LBRACK {
            lexer.NextToken()                       // [
            k = parseExp(lexer)                     // exp
            lexer.NextTokenOfKind(TOKEN_SEP_RBRACK) // ]
            lexer.NextTokenOfKind(TOKEN_OP_ASSIGN)  // =
            v = parseExp(lexer)                     // exp
            return
        }

        exp := parseExp(lexer)
        if nameExp, ok := exp.(*NameExp); ok {
            if lexer.LookAhead() == TOKEN_OP_ASSIGN {
                // Name ‘=’ exp => ‘[’ LiteralString ‘]’ = exp
                lexer.NextToken()
                k = &StringExp{nameExp.Line, nameExp.Name}
                v = parseExp(lexer)
                return
            }
        }

        return nil, exp
    }

    字段分三种情况解析，需要说明的是，去掉了k = v中的语法糖，将其变回["k"] = v。

前缀表达式

    为了避免parse_exp.go文件过长，把前缀表达式相关的解析函数放在单独的文件中。

    compiler/parser/parse_prefix_exp.go，定义解析函数parsePrefixExp

    package parser

    import . "lunago/compiler/ast"
    import . "lunago/compiler/lexer"

    func parsePrefixExp(lexer *Lexer) Exp {
        var exp Exp
        if lexer.LookAhead() == TOKEN_IDENTIFIER {
            line, name := lexer.NextIdentifier() // Name
            exp = &NameExp{line, name}
        } else { // ‘(’ exp ‘)’
            exp = parseParensExp(lexer)
        }
        return _finishPrefixExp(lexer, exp)
    }

    前缀表达式只能以标识符或者左圆括号开始，所以先前瞻一个token，根据情况解析出标识符或者圆括号表达式，然后调用函数 _finishPrefixExp 完成后续的解析工作。

    func _finishPrefixExp(lexer *Lexer, exp Exp) Exp {
        for {
            switch lexer.LookAhead() {
            case TOKEN_SEP_LBRACK: // prefixexp ‘[’ exp ‘]’
                lexer.NextToken()                       // ‘[’
                keyExp := parseExp(lexer)               // exp
                lexer.NextTokenOfKind(TOKEN_SEP_RBRACK) // ‘]’
                exp = &TableAccessExp{lexer.Line(), exp, keyExp}
            case TOKEN_SEP_DOT: // prefixexp ‘.’ Name
                lexer.NextToken()                    // ‘.’
                line, name := lexer.NextIdentifier() // Name
                keyExp := &StringExp{line, name}
                exp = &TableAccessExp{line, exp, keyExp}
            case TOKEN_SEP_COLON, // prefixexp ‘:’ Name args
                TOKEN_SEP_LPAREN, TOKEN_SEP_LCURLY, TOKEN_STRING: // prefixexp args
                exp = _finishFuncCallExp(lexer, exp)
            default:
                return exp
            }
        }
        return exp
    }

    表和字段访问表达式比较简单，直接在case中解析，顺便把字段访问表达式转换成了表访问表达式。函数调用表达式有点复杂，放到函数 _finishFuncCallExp 里面解析


圆括号表达式
    圆括号表达式的解析函数
    
    func parseParensExp(lexer *Lexer) Exp {
        lexer.NextTokenOfKind(TOKEN_SEP_LPAREN) // (
        exp := parseExp(lexer)                  // exp
        lexer.NextTokenOfKind(TOKEN_SEP_RPAREN) // )

        switch exp.(type) {
        case *VarargExp, *FuncCallExp, *NameExp, *TableAccessExp:
            return &ParensExp{exp}
        }

        // no need to keep parens
        return exp
    }

    由于圆括号会改变vararg和函数调用表达式的语义，所以需要保留着两种语句的圆括号，队医var表达式，也需要保留圆括号，否则之前的_checkVar()函数就会出现问题。奇遇的表达式两侧的圆括号，则完全没有必要留在ast中。

函数调用表达式

       函数调用表达式的解析函数 

       // functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
    func _finishFuncCallExp(lexer *Lexer, prefixExp Exp) *FuncCallExp {
        nameExp := _parseNameExp(lexer)
        line := lexer.Line() // todo
        args := _parseArgs(lexer)
        lastLine := lexer.Line()
        return &FuncCallExp{line, lastLine, prefixExp, nameExp, args}
    }


    可选的方法名由函数 _parseNameExp 解析

    func _parseNameExp(lexer *Lexer) *StringExp {
        if lexer.LookAhead() == TOKEN_SEP_COLON {
            lexer.NextToken()
            line, name := lexer.NextIdentifier()
            return &StringExp{line, name}
        }
        return nil
    }

    参数列表由函数 _parseArgs 解析

    // args ::=  ‘(’ [explist] ‘)’ | tableconstructor | LiteralString
    func _parseArgs(lexer *Lexer) (args []Exp) {
        switch lexer.LookAhead() {
        case TOKEN_SEP_LPAREN: // ‘(’ [explist] ‘)’
            lexer.NextToken() // TOKEN_SEP_LPAREN
            if lexer.LookAhead() != TOKEN_SEP_RPAREN {
                args = parseExpList(lexer)
            }
            lexer.NextTokenOfKind(TOKEN_SEP_RPAREN)
        case TOKEN_SEP_LCURLY: // ‘{’ [fieldlist] ‘}’
            args = []Exp{parseTableConstructorExp(lexer)}
        default: // LiteralString
            line, str := lexer.NextTokenOfKind(TOKEN_STRING)
            args = []Exp{&StringExp{line, str}}
        }
        return
    }


表达式的优化

    为了降低开发难度并提高可重用性、可以吧编译器分为前端、中端、后端三部分，优化，一般在中端和后端进行。不过这里仅仅设计了前段的词法分析和语法分析，以及后端的代码生成阶段，为了简化代码生成器，会在语法分析阶段进行一点优化。

    对全部由字面量参与的算术、按位、逻辑运算符表达式进行优化

    如果运算符表达式的值能够在编译期间计算出来，lua编译器会完全将其优化掉。可以在解析额运算符表达式的时候，顺便做一下这个优化，比如一元运算符，增加优化逻辑以后的解析函数代码如下：


    func parseExp2(lexer *Lexer) Exp {
        switch lexer.LookAhead() {
        case TOKEN_OP_UNM, TOKEN_OP_BNOT, TOKEN_OP_LEN, TOKEN_OP_NOT:
            line, op, _ := lexer.NextToken()
            exp := &UnopExp{line, op, parseExp2(lexer)}
            return optimizeUnaryOp(exp)
        }
        return parseExp1(lexer)
    }

    compiler/parser/optimizer.go ，优化相关的代码放在这里，定义函数
    optimizeUnaryOp（）

    package parser

    import "math"
    import "lunago/number"
    import . "lunago/compiler/ast"
    import . "lunago/compiler/lexer"

    func optimizeUnaryOp(exp *UnopExp) Exp {
        switch exp.Op {
        case TOKEN_OP_UNM:
            return optimizeUnm(exp)
        case TOKEN_OP_NOT:
            return optimizeNot(exp)
        case TOKEN_OP_BNOT:
            return optimizeBnot(exp)
        default:
            return exp
        }
    }

    一元取负运算符表达式的优化

    func optimizeUnm(exp *UnopExp) Exp {
        switch x := exp.Exp.(type) { // number?
        case *IntegerExp:
            x.Val = -x.Val
            return x
        case *FloatExp:
            if x.Val != 0 {
                x.Val = -x.Val
                return x
            }
        }
        return exp
    }


    其他的优化函数也类似。。。。。具体看代码

    这样，解析器基本上就差不多了，不过还需要一个入口函数，创建文件compiler/parser/parser.go，并定义函数Parse（）：


    package parser

    import . "lunago/compiler/ast"
    import . "lunago/compiler/lexer"

    /* recursive descent parser */

    func Parse(chunk, chunkName string) *Block {
        lexer := NewLexer(chunk, chunkName)
        block := parseBlock(lexer)
        lexer.NextTokenOfKind(TOKEN_EOF)
        return block
    } 

    先创建词法分析器，然后解析一个语法块，最后确保所有的源代码都已经解析完了。

单元测试

    
    package main

    import "encoding/json"
    import "io/ioutil"
    import "os"
    import "lunago/compiler/parser"

    func main() {
        if len(os.Args) > 1 {
            data, err := ioutil.ReadFile(os.Args[1])
            if err != nil {
                panic(err)
            }

            testParser(string(data), os.Args[1])
        }
    }

    testParser对语法分析器进行测试

    func testParser(chunk, chunkName string) {
        ast := parser.Parse(chunk, chunkName)
        b, err := json.Marshal(ast)
        if err != nil {
            panic(err)
        }
        println(string(b))
    }

    先对源代码进行解析，得到ast，然后把ast以json的格式打印。

    $ go run main.go helloworld.lua 
{"LastLine":6,"Stats":[{"Line":6,"LastLine":6,"PrefixExp":{"Line":6,"Name":"print"},"NameExp":null,"Args":[{"Line":6,"Str":"hello world！！！"}]}],"RetExps":null}

============


代码生成

有了语法分析器，就可以把lua源代码解析成一棵抽象语法树ast,进一步处理，利用它就可以生成lua字节码和函数原型，并最终输出二进制chunk文件。


定义funcInfo结构体

    每个lua函数都会被编译为函数原型存放在二进制chunk中，另外lua编译器还会生成一个main函数。以后的任务就是编写代码生成器（“代码”，在这里指的是存放在函数原型里面的lua虚拟机的字节码），把语法分析器输出的ast转换为函数原型。为了简单起见，把代码生成分成两个阶段：

    1、对ast进行处理，生成自定义的内部结构
    2、把内部的结构转换为函数原型

    compiler/codegen/func_info.go

    定义用来比奥市函数便捷结果的数据结构：

    package codegen

    import . "lunago/compiler/ast"
    import . "lunago/compiler/lexer"
    import . "lunago/vm"

    type funcInfo struct {
        //to do
    }

常量表    
    
    每个函数原型都有自己的常量表，里面存放函数体中出现的nil、布尔、数字或者字符串字面量。
    添加字段，用来表示常量表

    type funcInfo struct {
        constants map[interface{}]int
        //to do
    }

    为了便于查找，使用了map来存储常量，其中key是常量值，value是常量在表中的索引。给结构体定义一个方法：

    /* constants */

    func (self *funcInfo) indexOfConstant(k interface{}) int {
        if idx, found := self.constants[k]; found {
            return idx
        }

        idx := len(self.constants)
        self.constants[k] = idx
        return idx
    }

    该方法返回常量在表总的索引，如果常量不在表中，先把常量放到表中，在返回索引，索引从0开始递增。

寄存器的分配    

    lua虚拟机是基于寄存器的，所以在生成指令的时候需要进行寄存器的分配。简单的说，需要给每个局部变量和临时变量都分配一个寄存器，在局部变量退出作用域或者临时变量使用完毕以后，回收寄存器。给结构体funcInfo
    增加字段

        type funcInfo struct {
            constants map[interface{}]int
            usedRegs  int
            maxRegs   int
            //to do
        }

    主记录已经分配的寄存器数量和需要的足底啊寄存器数量就可以了。allocReg分配一个寄存器，必要的时候更新最大寄存器数量，并返回寄存器的索引。

    /* registers */

    func (self *funcInfo) allocReg() int {
        self.usedRegs++
        if self.usedRegs >= 255 {
            panic("function or expression needs too many registers")
        }
        if self.usedRegs > self.maxRegs {
            self.maxRegs = self.usedRegs
        }
        return self.usedRegs - 1
    }

    寄存器的索引是从0开始的，并且不能超过255.freeReg回收最近分配的寄存器

    func (self *funcInfo) freeReg() {
        if self.usedRegs <= 0 {
            panic("usedRegs <= 0 !")
        }
        self.usedRegs--
    }

    allocRegs分配连续的n个寄存器，返回第一个寄存器的索引：

    func (self *funcInfo) allocRegs(n int) int {
        if n <= 0 {
            panic("n <= 0 !")
        }
        for i := 0; i < n; i++ {
            self.allocReg()
        }
        return self.usedRegs - n
    }

     freeRegs方法回收最近分配的n个寄存器   

    func (self *funcInfo) freeRegs(n int) {
        if n < 0 {
            panic("n < 0 !")
        }
        for i := 0; i < n; i++ {
            self.freeReg()
        }
    }


局部变量表

    lua采用词法作用域。在函数内部，某个局部变量的作用域是包围该变量的最内层语句块。简单的说，每个块都会制造一个新的作用域，在块的内部可以使用局部变量声明语句声明（并初始化）局部变量，每个局部变量都会占用一个寄存器索引，当块结束以后，作用域也随之消失。局部变量不复存在，占用的寄存器也会被回收。对于repeat和for语句中的块，情况有一些不同。

    由于同一个局部变量名可以先后绑定不同的寄存器，为了简化代码，使用单项链表来串联同名的局部变量。

    定义结构体locVarInfo

    type locVarInfo struct {
        prev     *locVarInfo
        name     string
        scopeLv  int
        slot     int
        captured bool
    }

    其中字段prev使结构体locVarInfo成为单项链表的节点。name记录局部变量名，scopeLv记录记不变量所在的作用域层次，slot记录与局部变量绑定的寄存器列表，captured表示局部变量是否被闭包所捕获。

    修改结构体funcInfo，增加三个字段

    type funcInfo struct {
         constants map[interface{}]int
         usedRegs  int
         maxRegs   int
         scopeLv   int
         locVars   []*locVarInfo
         locNames  map[string]*locVarInfo
        //to do
    }


    其中scopeLv记录当前作用域的层次，locVars按照顺序记录函数内部声明的全部局部变量，locNames记录当前生效的局部变量。作用域层次从0开始，没今日一个作用域就加1。

    enterScope 进入新的作用域

    /* lexical scope */

    func (self *funcInfo) enterScope() {
        self.scopeLv++
    }

    addLocVar在当前作用域里面增加一个局部变量，返回其分配的寄存器索引

    func (self *funcInfo) addLocVar(name string) int {
        newVar := &locVarInfo{
            name:    name,
            prev:    self.locNames[name],
            scopeLv: self.scopeLv,
            slot:    self.allocReg(),
        }

        self.locVars = append(self.locVars, newVar)
        self.locNames[name] = newVar

        return newVar.slot
    }   


  slotOfLocVar检查局部变量名是否已经和某个寄存器绑定，如果过是，则返回寄存器的索引，否则返回-1

    func (self *funcInfo) slotOfLocVar(name string) int {
        if locVar, found := self.locNames[name]; found {
            return locVar.slot
        }
        return -1
    }

    exitScope退出作用域

    func (self *funcInfo) exitScope() {
        self.scopeLv--
        for _, locVar := range self.locNames {
            if locVar.scopeLv > self.scopeLv { // out of scope
                self.removeLocVar(locVar)
            }
        }
    }

    当退出作用域以后，需要删除该作用域中的局部变量（解绑局部变量名、回收寄存器），把这个逻辑封装在
removeLocVar中

    func (self *funcInfo) removeLocVar(locVar *locVarInfo) {
        self.freeReg()
        if locVar.prev == nil {
            delete(self.locNames, locVar.name)
        } else if locVar.prev.scopeLv == locVar.scopeLv {
            self.removeLocVar(locVar.prev)
        } else {
            self.locNames[locVar.name] = locVar.prev
        }
    }   

    先回收寄存器，然后检查是否有其他同名的局部变量，如果没有，则直接解绑局部变量名即可。如果有，并且在同一个作用域中，则递归调用removeLocVar进行处理。否则，同名的局部变零就在更外层的作用域中，需要把局部变量名与该局部变量重新绑定。

 break表   

    在for、repeat、while语句块的内部，可以使用break语句打断循环，。break语句处理起来有两个难点：
    1、break 语句可能在更深层次的块中，所以需要穿透块找到距离break语句最接近的那个or、repeat、while块
    2、break语句泗洪跳转指令实现，但是在处理break语句的时候，块可能还没有结束，所以跳转的目的地址还不确定。

    为了解决这两个问题，需要把跳转指定的地址记录在对应的or、repeat、while块中，等到块结束的时候在修复跳转的目标地址。

    给结构体funcInfo增加breaks字段

    type funcInfo struct {
         constants map[interface{}]int
         usedRegs  int
         maxRegs   int
         scopeLv   int
         locVars   []*locVarInfo
         locNames  map[string]*locVarInfo
         breaks    [][]int
        //to do
    }

    以后将for、repeat、while循环语句块称之为 循环快。使用数字来保存循环块中待处理的跳转指令，数组长度和块的深度对应，通过判断数组元素（也是数组）是否为nil，可以知道对应的块是不是循环块。

    修改enterScope，对进入的作用域是不是属于循环块进行标记

    func (self *funcInfo) enterScope(breakable bool) {
        self.scopeLv++
        if breakable {
            self.breaks = append(self.breaks, []int{}) //循环块
        } else {
            self.breaks = append(self.breaks, nil) //非循环块
        }
    }

    addBreakJmp把break语句对应的跳转指令添加到最近的循环块中。如果找不到循环块，则调用panic报错。

    func (self *funcInfo) addBreakJmp(pc int) {
        for i := self.scopeLv; i >= 0; i-- {
            if self.breaks[i] != nil { // breakable 循环块
                self.breaks[i] = append(self.breaks[i], pc)
                return
            }
        }

        panic("<break> at line ? not inside a loop!")
    }

    exitScope，在退出作用域的时候修复调转指令

    func (self *funcInfo) exitScope() {
        pendingBreakJmps := self.breaks[len(self.breaks)-1]
        self.breaks = self.breaks[:len(self.breaks)-1]

        a := self.getJmpArgA()
        for _, pc := range pendingBreakJmps {
            sBx := self.pc() - pc
            i := (sBx+MAXARG_sBx)<<14 | a<<6 | OP_JMP
            self.insts[pc] = uint32(i)
        }

        self.scopeLv--
        for _, locVar := range self.locNames {
            if locVar.scopeLv > self.scopeLv { // out of scope
                self.removeLocVar(locVar)
            }
        }
    }    


Upvalue表

    Upvalue实际上就是闭包，按照词法作用域捕获的外围函数中的局部变量。和局部变量类似，也需要把Upvalue名称和外围函数的局部变量绑定。不过由于Upvalue名称仅仅能绑定唯一的Upvalue，所以不需要使用链表结构。

    定义结构体upvalInfo

    type upvalInfo struct {
        locVarSlot int
        upvalIndex int
        index      int
    }

    如果Upvalue捕获的是直接外围函数的局部变量，则locVarSlot 保存该局部变量所占用的寄存器索引；否则Upvalue已经被直接的外围函数所捕获，upvalIndex记录该Upvalue在直接外围函数表中的索引。index记录Upvalue在函数中出现的顺序。修改结构体funcInfo，增加字段parent、Upvalues

    type funcInfo struct {
         constants map[interface{}]int
         usedRegs  int
         maxRegs   int
         scopeLv   int
         locVars   []*locVarInfo
         locNames  map[string]*locVarInfo
         breaks    [][]int
         parent    *funcInfo
         upvalues  map[string]upvalInfo
        //to do
    }    

   其中upvalues保存Upvalue表，parent则能够定位到外围函数的局部变量表和Upvalue表。 

   indexOfUpval判断名称是否已经和Upvalue绑定，如果是，返回Upvalue索引，否则尝试绑定，然后返回索引。如果绑定失败，返回-1.


    func (self *funcInfo) indexOfUpval(name string) int {
        if upval, ok := self.upvalues[name]; ok {
            return upval.index
        }
        if self.parent != nil {
            if locVar, found := self.parent.locNames[name]; found {
                idx := len(self.upvalues)
                self.upvalues[name] = upvalInfo{locVar.slot, -1, idx}
                locVar.captured = true
                return idx
            }
            if uvIdx := self.parent.indexOfUpval(name); uvIdx >= 0 {
                idx := len(self.upvalues)
                self.upvalues[name] = upvalInfo{-1, uvIdx, idx}
                return idx
            }
        }
        return -1
    }


字节码

    字节码，也就是lua虚拟机指令才是函数原型的主角，其他信息都是配角。修改funcInfo，为主角添加字段

    type funcInfo struct {
         constants map[interface{}]int
         usedRegs  int
         maxRegs   int
         scopeLv   int
         locVars   []*locVarInfo
         locNames  map[string]*locVarInfo
         breaks    [][]int
         parent    *funcInfo
         upvalues  map[string]upvalInfo
         insts     []uint32
        //to do
    }     

    直接存储编码后的指令，然后给结构体funcInfo添加一些方法，用来生成四种编码模式的指令。

    func (self *funcInfo) emitABC(opcode, a, b, c int) {
        i := b<<23 | c<<14 | a<<6 | opcode
        self.insts = append(self.insts, uint32(i))
    }

    func (self *funcInfo) emitABx(opcode, a, bx int) {
        i := bx<<14 | a<<6 | opcode
        self.insts = append(self.insts, uint32(i))
    }

    func (self *funcInfo) emitAsBx(opcode, a, b int) {
        i := (b+MAXARG_sBx)<<14 | a<<6 | opcode
        self.insts = append(self.insts, uint32(i))
    }

    func (self *funcInfo) emitAx(opcode, ax int) {
        i := ax<<6 | opcode
        self.insts = append(self.insts, uint32(i))
    }

    为了提高可读性，还需要定义emitMove、emitLoadNil、emitVararg等针对单个指令的方法。

    为了方便跳转指令的生成，定义个pc方法，这个方法返回已经生成的最后一条指令的成功女婿计数器（program  counter），也就是该指令的索引。

     func (self *funcInfo) pc() int {
        return len(self.insts) - 1
    }   

    fixSbx修复Sbx


    func (self *funcInfo) fixSbx(pc, sBx int) {
        i := self.insts[pc]
        i = i << 18 >> 18                  // clear sBx
        i = i | uint32(sBx+MAXARG_sBx)<<14 // reset sBx
        self.insts[pc] = i
    }

其他的一些信息    
    函数原型是递归结构，因为其内部还可以有子函数原型，与此对应，funcInfo结构体也需要是递归结构。对其进行修改，增加字段

    type funcInfo struct {
         constants map[interface{}]int
         usedRegs  int
         maxRegs   int
         scopeLv   int
         locVars   []*locVarInfo
         locNames  map[string]*locVarInfo
         breaks    [][]int
         parent    *funcInfo
         upvalues  map[string]upvalInfo
         insts     []uint32
         subFuncs  []*funcInfo
         numParams int
         isVararg  bool
        //to do
    }      

    subFuncs用于存放子函数的信息。其他两个字段是生成函数原型的时候所必须的，也需要记录下来。
    在定义一个newFuncInfo函数，用来创建结构体实例

     func newFuncInfo(parent *funcInfo, fd *FuncDefExp) *funcInfo {
        return &funcInfo{
            parent:    parent,
            subFuncs:  []*funcInfo{},
            locVars:   make([]*locVarInfo, 0, 8),
            locNames:  map[string]*locVarInfo{},
            upvalues:  map[string]upvalInfo{},
            constants: map[interface{}]int{},
            breaks:    make([][]int, 1),
            insts:     make([]uint32, 0, 8),
            numParams: len(fd.ParList),
            isVararg:  fd.IsVararg,
        }
    }   

 ast===>funcInfo实例   

编译块

    由于函数的主题实际上就是语句块，所以从块入手。
    compiler/codegen/cg_block.go

    定义函数cgBlock（）

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

    函数体由任意语句和一条可选的返回语句组成，所以先循环调用cgStat（）处理每条语句，如果有返回语句，就调用cgRetStat进行处理。

    cgRetStat的代码：

    func cgRetStat(fi *funcInfo, exps []Exp) {
        nExps := len(exps)
        if nExps == 0 {
            fi.emitReturn(0, 0)
            return
        }

        if nExps == 1 {
            if nameExp, ok := exps[0].(*NameExp); ok {
                if r := fi.slotOfLocVar(nameExp.Name); r >= 0 {
                    fi.emitReturn(r, 1)
                    return
                }
            }
            if fcExp, ok := exps[0].(*FuncCallExp); ok {
                r := fi.allocReg()
                cgTailCallExp(fi, fcExp, r)
                fi.freeReg()
                fi.emitReturn(r, -1)
                return
            }
        }

        multRet := isVarargOrFuncCall(exps[nExps-1])
        for i, exp := range exps {
            r := fi.allocReg()
            if i == nExps-1 && multRet {
                cgExp(fi, exp, r, -1)
            } else {
                cgExp(fi, exp, r, 1)
            }
        }
        fi.freeRegs(nExps)

        a := fi.usedRegs // correct?
        if multRet {
            fi.emitReturn(a, -1)
        } else {
            fi.emitReturn(a, nExps)
        }
    }    


    如果返回语句后面没有任何表达式，那么只要生成一条RETURN指令就可以了。

            nExps := len(exps)
        if nExps == 0 {
            fi.emitReturn(0, 0)
            return
        }

    如果返回语句后面还有表达式，要先对表达式进行处理，然后在生成RETUEN指令。

            multRet := isVarargOrFuncCall(exps[nExps-1])
        for i, exp := range exps {
            r := fi.allocReg()
            if i == nExps-1 && multRet {
                cgExp(fi, exp, r, -1)
            } else {
                cgExp(fi, exp, r, 1)
            }
        }
        fi.freeRegs(nExps)

         a := fi.usedRegs // correct?
        if multRet {
            fi.emitReturn(a, -1)
        } else {
            fi.emitReturn(a, nExps)
        }   


    有点繁琐，因为需要对最后一个表达式为vararg或者函数调用的情况进行特殊处理。另外，这里没有处理尾递归调用的情况。

    isVarargOrFuncCall位于compiler/parser/optimizer.go

    func isVarargOrFuncCall(exp Exp) bool {
        switch exp.(type) {
        case *VarargExp, *FuncCallExp:
            return true
        }
        return false
    }    

编译语句

    lua一共有15种语句，其中空语句不需要处理，非局部函数定义语句已经被转换为赋值语句。

    compiler/codegen/cg_stat.go 定义cgStat，对11种语句进行处理，不包括标签和goto

    package codegen

    import . "lunago/compiler/ast"

    func cgStat(fi *funcInfo, node Stat) {
        switch stat := node.(type) {
        case *FuncCallStat:
            cgFuncCallStat(fi, stat)
        case *BreakStat:
            cgBreakStat(fi, stat)
        case *DoStat:
            cgDoStat(fi, stat)
        case *WhileStat:
            cgWhileStat(fi, stat)
        case *RepeatStat:
            cgRepeatStat(fi, stat)
        case *IfStat:
            cgIfStat(fi, stat)
        case *ForNumStat:
            cgForNumStat(fi, stat)
        case *ForInStat:
            cgForInStat(fi, stat)
        case *AssignStat:
            cgAssignStat(fi, stat)
        case *LocalVarDeclStat:
            cgLocalVarDeclStat(fi, stat)
        case *LocalFuncDefStat:
            cgLocalFuncDefStat(fi, stat)
        case *LabelStat, *GotoStat:
            panic("label and goto statements are not supported!")
        }
    }

简单语句 

    local function f() end <===>local f; f = function () end, 
    所以局部函数定义语句处理起来也比较简单

    func cgLocalFuncDefStat(fi *funcInfo, node *LocalFuncDefStat) {
        r := fi.addLocVar(node.Name)
        cgFuncDefExp(fi, node.Exp, r)
    }

    对于函数调用语句，可以认为是对函数调用表达式进行求值，但是不需要任何返回值，所以处理起来一样简单

     func cgFuncCallStat(fi *funcInfo, node *FuncCallStat) {
        r := fi.allocReg()
        cgFuncCallExp(fi, node, r, 0)
        fi.freeReg()
    }   

    对于break语句，生成一条jmp指令，并吧地址保存在break表中就可以了，等到块退出以后，在修补跳转的偏移量。break语句的处理代码：

    func cgBreakStat(fi *funcInfo, node *BreakStat) {
        pc := fi.emitJmp(0, 0)
        fi.addBreakJmp(pc)
    }    


    由于do语句本质上就是一个块，作用就是引入新的作用域，所以也不难

     func cgDoStat(fi *funcInfo, node *DoStat) {
        fi.enterScope(false)    //非循环块
        cgBlock(fi, node.Block)
        fi.closeOpenUpvals()
        fi.exitScope()
    }   

    当居于变量退出作用域的时候，调用 closeOpenUpvals 把处于开启状态的Upvalue 闭合。如果有需要处理的局部变量，这个方法就会产生一个jmp指令，其操作数A指出了需要处理的第一个局部变量的寄存器索引。代码：
    compiler/codegen/func_info.go

    func (self *funcInfo) closeOpenUpvals() {
        a := self.getJmpArgA()
        if a > 0 {
            self.emitJmp(a, 0)
        }
    }

    通过getJmpArgA获取jmp指令的操作数A：

    func (self *funcInfo) getJmpArgA() int {
        hasCapturedLocVars := false
        minSlotOfLocVars := self.maxRegs
        for _, locVar := range self.locNames {
            if locVar.scopeLv == self.scopeLv {
                for v := locVar; v != nil && v.scopeLv == self.scopeLv; v = v.prev {
                    if v.captured {
                        hasCapturedLocVars = true
                    }
                    if v.slot < minSlotOfLocVars && v.name[0] != '(' {
                        minSlotOfLocVars = v.slot
                    }
                }
            }
        }
        if hasCapturedLocVars {
            return minSlotOfLocVars + 1
        } else {
            return 0
        }
    }    


while repeat语句

    while语句循环对表达式求值，如果结果为真，则执行块，继续循环，否则就跳过块结束循环。

    /*
               ______________
              /  false? jmp  |
             /               |
    while exp do block end <-'
          ^           \
          |___________/
               jmp
    */
    func cgWhileStat(fi *funcInfo, node *WhileStat) {
        //1
        pcBeforeExp := fi.pc()
        //2
        r := fi.allocReg()
        cgExp(fi, node.Exp, r, 1)
        fi.freeReg()
        //3
        fi.emitTest(r, 0)
        pcJmpToEnd := fi.emitJmp(0, 0)
        //4
        fi.enterScope(true)
        cgBlock(fi, node.Block)
        fi.closeOpenUpvals()
        fi.emitJmp(0, pcBeforeExp-fi.pc()-1)
        fi.exitScope()
        //5
        fi.fixSbx(pcJmpToEnd, fi.pc()-pcJmpToEnd)
    }

   1、先保存当前的Pc，因为后面计算跳转偏移量的时候会用到； 
   2、分配一个临时变量，对表达式求值，然后释放临时变量
   3、生成test和jmp指令，实现条件跳转，由于此时还没有对块进行处理，所以跳转的偏移量还没有办法给出
   4、对块进行处理，生成一条jmp指令，跳转到最开始
   5、修复第一条jmp指令的偏移量

   repeat语句先执行语句块，然后再对表达式求值，如果结果为真，则循环结束，否则继续循环。

    /*
            ______________
           |  false? jmp  |
           V              /
    repeat block until exp
    */
    func cgRepeatStat(fi *funcInfo, node *RepeatStat) {
        fi.enterScope(true)

        pcBeforeBlock := fi.pc()
        cgBlock(fi, node.Block)

        r := fi.allocReg()
        cgExp(fi, node.Exp, r, 1)
        fi.freeReg()

        fi.emitTest(r, 0)
        fi.emitJmp(fi.getJmpArgA(), pcBeforeBlock-fi.pc()-1)
        fi.closeOpenUpvals()

        fi.exitScope()
    }

     需要说明的是，repeat语句块的作用域是把后面的表达式也覆盖在内了，所以在表达式里，可以访问到块里声明的局部变量   

if语句

    lua的if语句包括必须出现的if表达式和块、任意多个elseif表达式和块，以及可选的else表达式和块。不过在语法分析阶段，已经对if语句进行了简化，把可选的else部分合并到了elseif部分里面

    /*
             _________________       _________________       _____________
            / false? jmp      |     / false? jmp      |     / false? jmp  |
           /                  V    /                  V    /              V
    if exp1 then block1 elseif exp2 then block2 elseif true then block3 end <-.
                       \                       \                       \      |
                        \_______________________\_______________________\_____|
                        jmp                     jmp                     jmp
    */
    func cgIfStat(fi *funcInfo, node *IfStat) {
        pcJmpToEnds := make([]int, len(node.Exps))
        pcJmpToNextExp := -1

        for i, exp := range node.Exps {
            if pcJmpToNextExp >= 0 {
                fi.fixSbx(pcJmpToNextExp, fi.pc()-pcJmpToNextExp)
            }

            r := fi.allocReg()
            cgExp(fi, exp, r, 1)
            fi.freeReg()

            fi.emitTest(r, 0)
            pcJmpToNextExp = fi.emitJmp(0, 0)

            fi.enterScope(false)
            cgBlock(fi, node.Blocks[i])
            fi.closeOpenUpvals()
            fi.exitScope()
            if i < len(node.Exps)-1 {
                pcJmpToEnds[i] = fi.emitJmp(0, 0)
            } else {
                pcJmpToEnds[i] = pcJmpToNextExp
            }
        }

        for _, pc := range pcJmpToEnds {
            fi.fixSbx(pc, fi.pc()-pc)
        }
    }    


for循环语句

    lua有两种for循环，其中数值for循环需要借助forprep和forloop指令实现。通用for循环需要借助tforcall、tforloop指令实现

    数值for处理：

    func cgForNumStat(fi *funcInfo, node *ForNumStat) {
        fi.enterScope(true)
        //1
        cgLocalVarDeclStat(fi, &LocalVarDeclStat{
            NameList: []string{"(for index)", "(for limit)", "(for step)"},
            ExpList:  []Exp{node.InitExp, node.LimitExp, node.StepExp},
        })
        fi.addLocVar(node.VarName)
        //2
        a := fi.usedRegs - 4
        pcForPrep := fi.emitForPrep(a, 0)
        cgBlock(fi, node.Block)
        fi.closeOpenUpvals()
        pcForLoop := fi.emitForLoop(a, 0)
        //3
        fi.fixSbx(pcForPrep, pcForLoop-pcForPrep-1)
        fi.fixSbx(pcForLoop, pcForPrep-pcForLoop)

        fi.exitScope()
    }    


    数值for使用了三个特殊的局部变量，分别保存索引、限制还有步长。
    1、声明三个局部变量，并且使用初值、限制以及步长表达式的值初始化这三个变量，另外，跟在关键字for后面的名称也声明了一个局部变量
    2、生成forprep指令，处理块，然后程程forloop指令
    3、吧指令里面的调整偏移量修复就可以了

    通用for循环

    func cgForInStat(fi *funcInfo, node *ForInStat) {
        fi.enterScope(true)

        cgLocalVarDeclStat(fi, &LocalVarDeclStat{
            NameList: []string{"(for generator)", "(for state)", "(for control)"},
            ExpList:  node.ExpList,
        })
        for _, name := range node.NameList {
            fi.addLocVar(name)
        }

        pcJmpToTFC := fi.emitJmp(0, 0)
        cgBlock(fi, node.Block)
        fi.closeOpenUpvals()
        fi.fixSbx(pcJmpToTFC, fi.pc()-pcJmpToTFC)

        rGenerator := fi.slotOfLocVar("(for generator)")
        fi.emitTForCall(rGenerator, len(node.NameList))
        fi.emitTForLoop(rGenerator+2, pcJmpToTFC-fi.pc()-1)

        fi.exitScope()
    }    

    通用for循环也使用了三个特殊的局部变量，另外跟在关键字for后面的名称也全部被声明为局部变量。

局部变量声明语句

    因为局部变量声明语句可以一次声明多个局部变量、并对变量进行初始化，所以有点麻烦

    func cgLocalVarDeclStat(fi *funcInfo, node *LocalVarDeclStat) {
        exps := removeTailNils(node.ExpList)
        nExps := len(exps)
        nNames := len(node.NameList)

        oldRegs := fi.usedRegs
        if nExps == nNames {
            for _, exp := range exps {
                a := fi.allocReg()
                cgExp(fi, exp, a, 1)
            }
        } else if nExps > nNames {
            for i, exp := range exps {
                a := fi.allocReg()
                if i == nExps-1 && isVarargOrFuncCall(exp) {
                    cgExp(fi, exp, a, 0)
                } else {
                    cgExp(fi, exp, a, 1)
                }
            }
        } else { // nNames > nExps
            multRet := false
            for i, exp := range exps {
                a := fi.allocReg()
                if i == nExps-1 && isVarargOrFuncCall(exp) {
                    multRet = true
                    n := nNames - nExps + 1
                    cgExp(fi, exp, a, n)
                    fi.allocRegs(n - 1)
                } else {
                    cgExp(fi, exp, a, 1)
                }
            }
            if !multRet {
                n := nNames - nExps
                a := fi.allocRegs(n)
                fi.emitLoadNil(a, n)
            }
        }

        fi.usedRegs = oldRegs
        for _, name := range node.NameList {
            fi.addLocVar(name)
        }
    }

    如果等号左侧声明的局部变量和等号右侧提供的表达式数量一样，那么处理起来就容易了。需要注意的是，新声明的局部变量，只有在声明语句结束以后才会生效，所以先分配临时的变量，对表达式求值，然后在讲局部变量名和寄存器绑定。如果表达式比局部变量多，处理起来也不麻烦

    多余的表达式也一样要求值，另外如果最后一个表达式是vararg或者函数调用，还需要特别处理一下，如果表达式比局部变量少，处理起来最麻烦

    这里有两种情况，如果最后一个表达式是vararg或者函数调用，则需要使用多重赋值初始化其余的局部变量，否则必须生成loadnil指令来初始化剩余的局部变量。最后释放临时变量，声明局部变量就可以了。

赋值语句   

    赋值语句不仅仅支持多重赋值，而且还可以同时给局部变量、Upvalue、全局变量赋值或者根据索引修改表，所以处理起来很麻烦。

    func cgAssignStat(fi *funcInfo, node *AssignStat) {
        exps := removeTailNils(node.ExpList)
        nExps := len(exps)
        nVars := len(node.VarList)

        tRegs := make([]int, nVars)
        kRegs := make([]int, nVars)
        vRegs := make([]int, nVars)
        oldRegs := fi.usedRegs

        for i, exp := range node.VarList {
            if taExp, ok := exp.(*TableAccessExp); ok {
                tRegs[i] = fi.allocReg()
                cgExp(fi, taExp.PrefixExp, tRegs[i], 1)
                kRegs[i] = fi.allocReg()
                cgExp(fi, taExp.KeyExp, kRegs[i], 1)
            } else {
                name := exp.(*NameExp).Name
                if fi.slotOfLocVar(name) < 0 && fi.indexOfUpval(name) < 0 {
                    // global var
                    kRegs[i] = -1
                    if fi.indexOfConstant(name) > 0xFF {
                        kRegs[i] = fi.allocReg()
                    }
                }
            }
        }
        for i := 0; i < nVars; i++ {
            vRegs[i] = fi.usedRegs + i
        }

        if nExps >= nVars {
            for i, exp := range exps {
                a := fi.allocReg()
                if i >= nVars && i == nExps-1 && isVarargOrFuncCall(exp) {
                    cgExp(fi, exp, a, 0)
                } else {
                    cgExp(fi, exp, a, 1)
                }
            }
        } else { // nVars > nExps
            multRet := false
            for i, exp := range exps {
                a := fi.allocReg()
                if i == nExps-1 && isVarargOrFuncCall(exp) {
                    multRet = true
                    n := nVars - nExps + 1
                    cgExp(fi, exp, a, n)
                    fi.allocRegs(n - 1)
                } else {
                    cgExp(fi, exp, a, 1)
                }
            }
            if !multRet {
                n := nVars - nExps
                a := fi.allocRegs(n)
                fi.emitLoadNil(a, n)
            }
        }

        for i, exp := range node.VarList {
            if nameExp, ok := exp.(*NameExp); ok {
                varName := nameExp.Name
                if a := fi.slotOfLocVar(varName); a >= 0 {
                    fi.emitMove(a, vRegs[i])
                } else if b := fi.indexOfUpval(varName); b >= 0 {
                    fi.emitSetUpval(vRegs[i], b)
                } else if a := fi.slotOfLocVar("_ENV"); a >= 0 {
                    if kRegs[i] < 0 {
                        b := 0x100 + fi.indexOfConstant(varName)
                        fi.emitSetTable(a, b, vRegs[i])
                    } else {
                        fi.emitSetTable(a, kRegs[i], vRegs[i])
                    }
                } else { // global var
                    a := fi.indexOfUpval("_ENV")
                    if kRegs[i] < 0 {
                        b := 0x100 + fi.indexOfConstant(varName)
                        fi.emitSetTabUp(a, b, vRegs[i])
                    } else {
                        fi.emitSetTabUp(a, kRegs[i], vRegs[i])
                    }
                }
            } else {
                fi.emitSetTable(tRegs[i], kRegs[i], vRegs[i])
            }
        }

        // todo
        fi.usedRegs = oldRegs
    }


    由于赋值语句等号左边可以出现t[k]这样的表达式，灯油右边可以出现任意表达式，所以需要先分配临时变量，对这些表达式进行求值，然后在统一生成赋值指令。tRegs\kRegs\vRegs这三个数组分别记录为表、键、值分配的临时变量。

    先处理等号左侧的索引表达式，分配临时变量，并对表和键求值。然后统一为等号右侧的表达式计算寄存器的索引。

    这里的做法和局部变量声明语句类似，需要考虑多重赋值，也需要在必要的时候补上loadnil指令。

    在循环中对赋值进行处理：如果给局部变量赋值，生成move指令；如果给Upvalue赋值，生成Setupval指令；如果给全局变量赋值，生成settabup指令；如果值按照索引给表赋值，则生成settable指令。循环结束以后，需要释放掉所有的临时变量。


编译表达式

    lua表达式大值可以分为：字面量表达式、构造器表达式、运算符表达式、前缀表达式、vararg表达式5种。

    compiler/codegen/cg_exp.go 定义函数cgExp（）函数：

    package codegen

    import . "lunago/compiler/ast"
    import . "lunago/compiler/lexer"
    import . "lunago/vm"

    func cgExp(fi *funcInfo, node Exp, a, n int) {
        switch exp := node.(type) {
        case *NilExp:
            fi.emitLoadNil(a, n)
        case *FalseExp:
            fi.emitLoadBool(a, 0, 0)
        case *TrueExp:
            fi.emitLoadBool(a, 1, 0)
        case *IntegerExp:
            fi.emitLoadK(a, exp.Val)
        case *FloatExp:
            fi.emitLoadK(a, exp.Val)
        case *StringExp:
            fi.emitLoadK(a, exp.Str)
        case *ParensExp:
            cgExp(fi, exp.Exp, a, 1)
        case *VarargExp:
            cgVarargExp(fi, exp, a, n)
        case *FuncDefExp:
            cgFuncDefExp(fi, exp, a)
        case *TableConstructorExp:
            cgTableConstructorExp(fi, exp, a)
        case *UnopExp:
            cgUnopExp(fi, exp, a)
        case *BinopExp:
            cgBinopExp(fi, exp, a)
        case *ConcatExp:
            cgConcatExp(fi, exp, a)
        case *NameExp:
            cgNameExp(fi, exp, a)
        case *TableAccessExp:
            cgTableAccessExp(fi, exp, a)
        case *FuncCallExp:
            cgFuncCallExp(fi, exp, a, n)
        }
    }

    字面量表达式，只要生成相应的load指令就可以了，所以直接在case语句里面处理。圆括号表达式也比较好处理，也直接写到了case种，其他的表达式需要使用专门的函数。vararg表达式的处理函数如下：

    func cgVarargExp(fi *funcInfo, node *VarargExp, a, n int) {
        if !fi.isVararg {
            panic("cannot use '...' outside a vararg function")
        }
        fi.emitVararg(a, n)
    }   





























































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































































