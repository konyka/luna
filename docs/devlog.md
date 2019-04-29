

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































