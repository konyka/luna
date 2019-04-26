



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





7、
8、
9、
10、11、12、13、14、15、16、17、18、19、20、21、22、23、24、25、26、27、28、29、30、
















