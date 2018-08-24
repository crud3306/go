//在 Golang 语言中，文件使用指向 os.File 类型的指针来表示的，也叫做文件句柄。
//注意，标准输入 os.Stdin 和标准输出 os.Stdout ，他们的类型都是 *os.File 哟。在任何计算机设备中，文件是都是必须的对象，而在 Web编程中,文件的操作一直是 Web程序员经常遇到的问题,文件操作在 Web应用中是必须的,非常有用的,我们经常遇到生成文件目录,文件(夹)编辑等操作。

 

// 一.文件的读取姿势
// ===========================
/*
 1 歌曲：我有一只小毛驴
 2 歌手：碧瑶
 3 专辑：《瑶谣摇》
 4 发行时间：2014-03-18
 5 词：付林
 6 曲：付林
 7 歌词：
 8 我有一只小毛驴
 9 我从来也不骑
10 有一天我心血来潮骑着去赶集
11 我手里拿着小皮鞭
12 我心里很得意
13 不知怎么哗啦啦啦啦
14 摔了一身泥
15 我有一只小毛驴
16 我从来也不骑
17 有一天我心血来潮骑着去赶集
18 我手里拿着小皮鞭
19 我心里很得意
20 不知怎么哗啦啦啦啦
21 摔了一身泥
22 我有一只小毛驴 我从来也不骑
23 有一天我心血来潮 骑着去赶集
24 我手里拿着小皮鞭 我心里很得意
25 不知怎么哗啦啦啦 摔了一身泥
26 我有一只小毛驴 我从来也不骑
27 有一天我心血来潮 骑着去赶集
28 我手里拿着小皮鞭 我心里很得意
29 不知怎么哗啦啦啦 摔了一身泥
30 我手里拿着小皮鞭 我心里很得意
31 不知怎么哗啦啦啦 摔了一身泥
我有一只小毛驴.txt（文件内容戳我）
*/

// 姿势1. 顺序按行读取文件内容

#!/usr/bin/env gorun

package main
 
import (
    "bufio"
    "fmt"
    "io"
    "os"
)

var (
    FileName string = "E:\\Code\\Golang\\Golang_Program\\文件处理\\我有一只小毛驴.txt"    ////这是我们需要打开的文件，当然你也可以把它定义到从某个配置文件来获取变量。
    InputFile  *os.File    //变量 InputFile 是 *os.File 类型的。该类型是一个结构，表示一个打开文件的描述符（文件句柄）。
    InputError error    //我们使用 os 包里的 Open 函数来打开一个文件。如果文件不存在或者程序没有足够的权限打开这个文件，Open函数会返回一个错误，InputError变量就是用来接收这个错误的。
    Count int            //这个变量是我们用来统计行号的，默认值为0.
)

func main() {
    //InputFile,InputError = os.OpenFile(FileName,os.O_CREATE|os.O_RDWR,0644) //打开FileName文件，如果不存在就创建新文件，打开的权限是可读可写，权限是644。这种打开方式相对下面的打开方式权限会更大一些。
    InputFile, InputError = os.Open(FileName) //使用 os 包里的 Open 函数来打开一个文件。该函数的参数是文件名，类型为 string 。我们以只读模式打开"FileName"文件。
    if InputError != nil {    //如果打开文件出错，那么我们可以给用户一些提示，然后在推出函数。
        fmt.Printf("An error occurred on opening the inputfile\n" +
            "Does the file exist?\n" +
            "Have you got acces to it?\n")
        return // exit the function on error
    }
    defer InputFile.Close()        //defer关键字是用在程序即将结束时执行的代码,确保在程序退出前关闭该文件。
    inputReader := bufio.NewReader(InputFile) //我们使用 bufio.NewReader()函数来获得一个读取器变量（读取器）。我们可以很方便的操作相对高层的 string 对象，而避免了去操作比较底层的字节。
    for {
        Count += 1
        inputString, readerError := inputReader.ReadString('\n')  //我们将inputReader里面的字符串按行进行读取。
        if readerError == io.EOF {
            return  //如果遇到错误就终止循环。
        }
        fmt.Printf("The %d line is: %s",Count, inputString)    //将文件的内容逐行（行结束符'\n'）读取出来。
    }
}

//#以上代码执行结果如下：
//The 1 line is: 歌曲：我有一只小毛驴
//The 2 line is: 歌手：碧瑶
//...省略
//The 30 line is: 我手里拿着小皮鞭 我心里很得意
 


//姿势2. 按列读取数据
 
#!/usr/bin/env gorun
package main

import (
    "fmt"
    "os"
)

var (
    FileName = "E:\\Code\\Golang\\Golang_Program\\文件处理\\a.txt"
)

/*
20 注意：FileName的文件内容如下：
21 A B C
22 a b c
23 1 2 3
24 */

func main() {
    file, err := os.Open(FileName)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    var Column1, Column2, Column3 []string    //定义3个切片，每个切片用来保存不同列的数据。

    for {
        var FirstRowColumn, SecondRowColumn, ThirdRowColumn string
        _, err := fmt.Fscanln(file, &FirstRowColumn, &SecondRowColumn, &ThirdRowColumn)    //如果数据是按列排列并用空格分隔的，我们可以使用 fmt 包提供的以 FScan 开头的一系列函数来读取他们。
        if err != nil {
            break
        }
        Column1 = append(Column1, FirstRowColumn)  //将第一列的每一行的参数追加到空切片Column1中。以下代码类似。
        Column2 = append(Column2, SecondRowColumn)
        Column3 = append(Column3, ThirdRowColumn)
    }
    fmt.Println(Column1)
    fmt.Println(Column2)
    fmt.Println(Column3)
}

/*
#以上代码执行结果如下：
[A a 1]
[B b 2]
[C c 3]
 */



// 姿势3. 带缓冲的读取

#!/usr/bin/env gorun
 
package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
)
 
var (
    FileName = "E:\\Code\\Golang\\Golang_Program\\文件处理\\我有一只小毛驴.txt"
)

func main() {
    f,err := os.Open(FileName)
    if err != nil{
        panic(err)
    }
    defer f.Close()
    ReadSize := make([]byte,1024) //指定每次读取的大小为1024。
    ReadByte := make([]byte,4096,4096) //指定读取到的字节数。
    r := bufio.NewReader(f)
    for   {
        ActualSize,err := r.Read(ReadSize)    //回返回每次读取到的实际字节大小。
        if err != nil && err != io.EOF {
            panic(err)
        }
        if ActualSize == 0 {
            break
        }
        ReadByte = append(ReadByte,ReadSize[:ActualSize]...)    //将每次的读取到的内容都追加到我们定义的切片中。
    }
    fmt.Println(string(ReadByte))    //打印我们读取到的内容，注意，不能直接读取，因为我们的切片的类型是字节，需要转换成字符串这样我们读取起来会更方便。
}
 


// 姿势4.将整个文件的内容读到一个字节切片中

#!/usr/bin/env gorun

package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "reflect"
)

var (
    FileName = "E:\\Code\\Golang\\Golang_Program\\文件处理\\我有一只小毛驴.txt"
)

func main() {
    buf, err := ioutil.ReadFile(FileName) //将整个文件的内容读到一个字节切片中。
    fmt.Println(reflect.TypeOf(buf))
    if err != nil {
        fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
        // panic(err.Error())
    }
    fmt.Printf("%s\n", string(buf))
}
/*
33 #以上代码执行结果如下：
34 []uint8
35 歌曲：我有一只小毛驴
36 歌手：碧瑶
37 专辑：《瑶谣摇》
38 发行时间：2014-03-18
39 词：付林
40 曲：付林
41 歌词：
42 我有一只小毛驴
43 我从来也不骑
44 有一天我心血来潮骑着去赶集
45 我手里拿着小皮鞭
46 我心里很得意
47 不知怎么哗啦啦啦啦
48 摔了一身泥
49 我有一只小毛驴
50 我从来也不骑
51 有一天我心血来潮骑着去赶集
52 我手里拿着小皮鞭
53 我心里很得意
54 不知怎么哗啦啦啦啦
55 摔了一身泥
56 我有一只小毛驴 我从来也不骑
57 有一天我心血来潮 骑着去赶集
58 我手里拿着小皮鞭 我心里很得意
59 不知怎么哗啦啦啦 摔了一身泥
60 我有一只小毛驴 我从来也不骑
61 有一天我心血来潮 骑着去赶集
62 我手里拿着小皮鞭 我心里很得意
63 不知怎么哗啦啦啦 摔了一身泥
64 我手里拿着小皮鞭 我心里很得意
65 不知怎么哗啦啦啦 摔了一身泥
*/

 


// 二.文件的写入姿势
// ===========================

// 姿势1：打开一个文件，如果没有就创建，如果有这个文件就清空文件内容（相当于python中的"w"）

#!/usr/bin/env gorun
 
package main

import (
    "os"
    "log"
)

func main()  {
    f,err := os.Create("a.txt") //姿势一：打开一个文件，如果没有就创建，如果有这个文件就清空文件内容,需要用两个变量接受相应的参数
    if err != nil {
        log.Fatal(err)
    }
    f.WriteString("yinzhengjie\n") //往文件写入相应的字符串。
    f.Close()
}


//姿势2：以追加的方式打开一个文件（相当于python中的"a"）
/*
　　OpenFile 函数有三个参数：文件名、一个或多个标志（使用逻辑运算符“|”连接），使用的文件权限。我们通常会用到以下标志：

　　　　1>.os.O_RDONLY ：只读

　　　　2>.os.WRONLY ：只写

　　　　3>.os.O_CREATE ：创建：如果指定文件不存在，就创建该文件。

　　　　4>.os.O_TRUNC ：截断：如果指定文件已存在，就将该文件的长度截为0。

　　在读文件的时候，文件的权限是被忽略的，所以在使用 OpenFile 时传入的第三个参数可以用0。
*/
 
#!/usr/bin/env gorun


package main

import (
    "os"
    "log"
)

func main()  {
    f,err := os.OpenFile("a.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR,0644) //表示最佳的方式打开文件，如果不存在就创建，打开的模式是可读可写，权限是644
    if err    != nil {
        log.Fatal(err)
    }
    f.WriteString("yinzhengjie\n")
    f.Close()
}


//姿势3：修改文件内容-随机写入（自定义插入的位置，相当python重的seek方法）
 
#!/usr/bin/env gorun
 
package main

import (
    "os"
    "log"
)


func main()  {
	f,err := os.OpenFile("a.txt",os.O_CREATE|os.O_RDWR,0644)
    if err != nil {
        log.Fatal(err)
    }

    f.WriteString("yinzhengjie\n")
    f.Seek(1,os.SEEK_SET) //表示文件的其实位置，从第二个字符往后写入。
    f.WriteString("$$$")
    f.Close()
}



// 姿势4.ioutil方法创建文件

#!/usr/bin/env gorun

package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

var (
    FileName = "E:\\Code\\Golang\\Golang_Program\\文件处理\\我有一只小毛驴.txt"
    OutputFile = "E:\\Code\\Golang\\Golang_Program\\文件处理\\复制的小毛驴.txt"
)
 
func main() {
    buf, err := ioutil.ReadFile(FileName) //将整个文件的内容读到一个切片中。
    if err != nil {
        fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
        // panic(err.Error())
    }
    //fmt.Printf("%s\n", string(buf))
    err = ioutil.WriteFile(OutputFile, buf, 0x644)    //我们将读取到的内容又重新写入到另外一个OutputFile文件中去。
    if err != nil {
        panic(err. Error())
    }
    
    /*注意，在执行该代码之后，就会生成一个OutputFile文件，其内容和FileName的内容是一致的哟！*/
}
 



// 三.文件的删除操作
// ===========================
#!/usr/bin/env gorun
  
package main

import "os"

var (
    FileName = "E:\\Code\\Golang\\Golang_Program\\文件处理\\复制的小毛驴.txt"
)

func main() {
    os.Remove(FileName)    //删除文件比较方便，直接用OS包就搞定啦的Remove方法就搞定案啦！
}




// 四.目录的操作姿势
// ===========================
// 1.目录的创建操作

#!/usr/bin/env gorun
 
package main

import (
    "os"
)

var (
    OneLevelDirectory = "yinzhengjie"
    MultilevelDirectory = "yinzhengjie/golang/code"
)
func main() {
    os.Mkdir(OneLevelDirectory, 0777)    //创建名称为OneLevelDirectory的目录，设置权限为0777。相当于Linux系统中的“mkdir yinzhengjie”
    os.MkdirAll(MultilevelDirectory, 0777)    //创建MultilevelDirectory多级子目录，设置权限为0777。相当于Linux中的 “mkdir -p yinzhengjie/golang/code”
}
 

2.目录的删除操作

#!/usr/bin/env gorun 

package main
 
import (
    "fmt"
    "os"
)

var (
    OneLevelDirectory = "yinzhengjie"
    MultilevelDirectory = "yinzhengjie/golang/code"
)

func main() {
    err := os.Remove(MultilevelDirectory) //删除名称为OneLevelDirectory的目录，当目录下有文件或者其他目录是会出错。
    if err != nil {
        fmt.Println(err)
    }
    os.RemoveAll(OneLevelDirectory) //根据path删除多级子目录，如果 path是单个名称，那么该目录不删除。
}




// 五.文件处理进阶知识
// ===========================
 1 [root@yinzhengjie code]# ll
 2 总用量 16
 3 -rw-r--r--+ 1 root root 891 11月  8 13:55 littleDonkey
 4 -rw-r--r--+ 1 root root 734 11月  8 13:57 readFile.go
 5 [root@yinzhengjie code]# more littleDonkey 
 6 歌曲：我有一只小毛驴
 7 歌手：碧瑶
 8 专辑：《瑶谣摇》
 9 发行时间：2014-03-18
10 词：付林
11 曲：付林
12 歌词：
13 我有一只小毛驴
14 我从来也不骑
15 有一天我心血来潮骑着去赶集
16 我手里拿着小皮鞭
17 我心里很得意
18 不知怎么哗啦啦啦啦
19 摔了一身泥
20 我有一只小毛驴
21 我从来也不骑
22 有一天我心血来潮骑着去赶集
23 我手里拿着小皮鞭
24 我心里很得意
25 不知怎么哗啦啦啦啦
26 摔了一身泥
27 我有一只小毛驴 我从来也不骑
28 有一天我心血来潮 骑着去赶集
29 我手里拿着小皮鞭 我心里很得意
30 不知怎么哗啦啦啦 摔了一身泥
31 我有一只小毛驴 我从来也不骑
32 有一天我心血来潮 骑着去赶集
33 我手里拿着小皮鞭 我心里很得意
34 不知怎么哗啦啦啦 摔了一身泥
35 我手里拿着小皮鞭 我心里很得意
36 不知怎么哗啦啦啦 摔了一身泥
37 [root@yinzhengjie code]# tar zcf littleDonkey.tar.gz littleDonkey     #创建压缩文件
38 [root@yinzhengjie code]# go run readFile.go littleDonkey.tar.gz 
39 littleDonkey0000644000000000000000000000157313200516146012152 0ustar  rootroot歌曲：我有一只小毛驴
40 歌手：碧瑶
41 专辑：《瑶谣摇》
42 发行时间：2014-03-18
43 词：付林
44 曲：付林
45 歌词：
46 我有一只小毛驴
47 我从来也不骑
48 有一天我心血来潮骑着去赶集
49 我手里拿着小皮鞭
50 我心里很得意
51 不知怎么哗啦啦啦啦
52 摔了一身泥
53 我有一只小毛驴
54 我从来也不骑
55 有一天我心血来潮骑着去赶集
56 我手里拿着小皮鞭
57 我心里很得意
58 不知怎么哗啦啦啦啦
59 摔了一身泥
60 我有一只小毛驴 我从来也不骑
61 有一天我心血来潮 骑着去赶集
62 我手里拿着小皮鞭 我心里很得意
63 不知怎么哗啦啦啦 摔了一身泥
64 我有一只小毛驴 我从来也不骑
65 有一天我心血来潮 骑着去赶集
66 我手里拿着小皮鞭 我心里很得意
67 不知怎么哗啦啦啦 摔了一身泥
68 我手里拿着小皮鞭 我心里很得意
69 不知怎么哗啦啦啦 摔了一身泥
70 Done reading file
71 [root@yinzhengjie code]# 