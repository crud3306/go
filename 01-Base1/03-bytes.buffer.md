
bytes.buffer
============
bytes.buffer是一个缓冲byte类型的缓冲器，这个缓冲器里存放着都是byte。



创建一个缓冲器
----------
- NewBuffer
如果在使用bytes.NewBuffer的时候，参数是[]byte的slice的话，缓冲器里就是这个slice的内容，如果参数是nil的话，意思是New一个空的缓冲器里

- NewBufferString
还可以用bytes.NewBufferString("hello")来建立一个内容是hello的缓冲器


```golang
buf1:=bytes.NewBufferString("hello")
buf2:=bytes.NewBuffer([]byte("hello"))
buf3:=bytes.NewBuffer([]byte{"h","e","l","l","o"})
```
以上三者等效


```golang
buf4:=bytes.NewBufferString("")
buf5:=bytes.NewBuffer([]byte{})
```
以上两者等效



如果buffer在new的时候是空的也没关系，因为可以用Write来写入，写在尾部


写入到缓冲器（缓冲器变大）
-------------
- Write —- func (b 星号Buffer) Write(p []byte) (n int, err error)

使用Write方法，将一个byte类型的slice放到缓冲器的尾部
```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    s := []byte(" world")
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String())  //buf.String()方法是吧buf里的内容转成string，以便于打印
    buf.Write(s) //将s这个slice写到buf的尾部
    fmt.Println(buf.String())  //打印 hello world
}
```

- WriteString —- func (b 星号Buffer) WriteString(s string) (n int, err error)

使用WriteString方法，将一个字符串放到缓冲器的尾部
```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    s := " world"
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String())  //buf.String()方法是吧buf里的内容转成string，以便于打印
    buf.WriteString(s) //将s这个string写到buf的尾部
    fmt.Println(buf.String())  //打印 hello world
}
```

- WriteByte —- func (b 星号Buffer) WriteByte(c byte) error

使用WriteByte方法，将一个byte类型的数据放到缓冲器的尾部
```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    var s byte = '!'
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String())  //buf.String()方法是吧buf里的内容转成string，以便于打印
    buf.WriteByte(s) //将s这个string写到buf的尾部
    fmt.Println(buf.String())  //打印 hello!
}
```

- WriteRune —- func (b 星号Buffer) WriteRune(r rune) (n int, err error)
使用WriteRune方法，将一个rune类型的数据放到缓冲器的尾部

```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    var s rune = '好'
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String())  //buf.String()方法是吧buf里的内容转成string，以便于打印
    buf.WriteRune(s) //将s这个string写到buf的尾部
    fmt.Println(buf.String())  //打印 hello好
}
```



从缓冲器写出（缓冲器变小）
-------------
WriteTo—- func (b 星号Buffer) WriteTo(w io.Writer) (n int64, err error)
使用WriteTo方法，将一个缓冲器的数据写到w里，w是实现io.Writer的，比如os.File就是实现io.Writer

```golang
package main

import (
    "bytes"
    "fmt"
    "os"
)

func main() {
    file , _ := os.Create("text.txt")
    buf := bytes.NewBufferString("hello")
    buf.WriteTo(file) //hello写到text.txt文件中了
    fmt.Fprintf(file, buf.String()) //虽然这不在讨论范围，但这句效果同上
}
```



读出缓冲器（缓冲器变小）
---------------
- Read —- func (b 星号Buffer) Read(p []byte) (n int, err error)

给Read方法一个容器p，读完后，p就满了，缓冲器相应的减少了，返回的n为成功读的数量  

如，缓冲器是一个装满5升水的杯子，这个杯子有Read方法，给Read方法一个3升的杯子
Read完后，5升杯子里有2升水，3升的杯子满了，返回的n为3 在一次Read时，5升杯子里有0升水，3升的杯子还是满的，但其中有2升的水被新倒入的水替代了，返回的n为2
```golang
package main
import(
    "fmt"
    "bytes"
)

func main() {
    s1:=[]byte("hello")                //申明一个slice为s1
    buff:=bytes.NewBuffer(s1)     //new一个缓冲器buff，里面存着hello这5个byte
    s2:=[]byte(" world")                 //申明另一个slice为s2
    buff.Write(s2)                     //把s2写入添加到buff缓冲器内
    fmt.Println(buff.String())            //使用缓冲器的String方法转成字符串，并打印："hello world"

    s3:=make([]byte,3)               //申明一个空的slice为s3，容量为3
    buff.Read(s3)                       //把buff的内容读入到s3内，因为s3的容量为3，所以只读了3个过来
    fmt.Println(buff.String())       //buff的前3个字符被读走了，所以buff变成："lo world"
    fmt.Println(string(s3))          //空的s3被写入3个字符，所以为"hel"
    buff.Read(s3)                       //把buff的内容读入到s3内，因为s3的容量为3，所以只读了3个过来，原来s3的内容被覆盖了
    fmt.Println(buff.String())       //buff的前3个字符又被读走了，所以buff变成："world"
    fmt.Println(string(s3))          //原来的s3被从"hel"变成"lo "，因为"hel"被覆盖了
}
```

- ReadByte —- func (b 星号Buffer) ReadByte() (c byte, err error)

返回缓冲器头部的第一个byte，缓冲器头部第一个byte被拿掉
```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String()) //buf.String()方法是吧buf里的内容转成string，>以便于打印
    b, _ := buf.ReadByte()    //读取第一个byte，赋值给b
    fmt.Println(buf.String()) //打印 ello，缓冲器头部第一个h被拿掉
    fmt.Println(string(b))    //打印 h
}
```

- ReadRune —- func (b 星号Buffer) ReadRune() (r rune, size int, err error)
ReadRune和ReadByte很像 返回缓冲器头部的第一个rune，缓冲器头部第一个rune被拿掉
```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    buf := bytes.NewBufferString("好hello")
    fmt.Println(buf.String()) //buf.String()方法是吧buf里的内容转成string，>以便于打印
    b, n, _ := buf.ReadRune() //读取第一个rune，赋值给b
    fmt.Println(buf.String()) //打印 hello
    fmt.Println(string(b))    //打印中文字： 好，缓冲器头部第一个“好”被拿掉
    fmt.Println(n)            //打印3，“好”作为utf8储存占3个byte
    b, n, _ = buf.ReadRune()  //再读取第一个rune，赋值给b
    fmt.Println(buf.String()) //打印 ello
    fmt.Println(string(b))    //打印h，缓冲器头部第一个h被拿掉
    fmt.Println(n)            //打印 1，“h”作为utf8储存占1个byte
}
```

- ReadBytes —- func (b 星号Buffer) ReadBytes(delim byte) (line []byte, err error)
ReadBytes和ReadByte根本就不是一回事
ReadBytes需要一个byte作为分隔符，读的时候从缓冲器里找第一个出现的分隔符（delim），找到后，把从缓冲器头部开始到分隔符之间的所有byte进行返回，作为byte类型的slice，返回后，缓冲器也会空掉一部分
```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    var d byte = 'e' //分隔符为e
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String()) //buf.String()方法是吧buf里的内容转成string，以便于打印
    b, _ := buf.ReadBytes(d)  //读到分隔符，并返回给b
    fmt.Println(buf.String()) //打印 llo，缓冲器被取走一些数据
    fmt.Println(string(b))    //打印 he，找到e了，将缓冲器从头开始，到e的内容都返回给b
}
```


- ReadString—- func (b 星号Buffer) ReadString(delim byte) (line string, err error)
ReadBytes和ReadString基本就是一回事  
ReadBytes需要一个byte作为分隔符，读的时候从缓冲器里找第一个出现的分隔符（delim），找到后，把从缓冲器头部开始到分隔符之间的所有byte进行返回，作为字符串，返回后，缓冲器也会空掉一部分
```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    var d byte = 'e' //分隔符为e
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String()) //buf.String()方法是吧buf里的内容转成string，以便于打印
    b, _ := buf.ReadString(d)  //读到分隔符，并返回给b
    fmt.Println(buf.String()) //打印 llo，缓冲器被取走一些数据
    fmt.Println(b)    //打印 he，找到e了，将缓冲器从头开始，到e的内容都返回给b
}
```



读入缓冲器（缓冲器变大）
-------------
ReadFrom—- func (b 星号Buffer) ReadFrom(r io.Reader) (n int64, err error)
从一个实现io.Reader接口的r，把r里的内容读到缓冲器里，n返回读的数量
```golang
package main

import (
    "bytes"
    "fmt"
    "os"
)

func main() {
    file, _ := os.Open("test.txt")  //test.txt的内容是“world”
    buf := bytes.NewBufferString("hello ")
    buf.ReadFrom(file)              //将text.txt内容追加到缓冲器的尾部
    fmt.Println(buf.String())    //打印“hello world”
}
```



从缓冲器取出（缓冲器变小）
-------------
Next —- func (b 星号Buffer) Next(n int) []byte
返回前n个byte，成为slice返回，原缓冲器变小

```golang
package main

import (
    "bytes"
    "fmt"
)

func main() {
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String())
    b := buf.Next(2)   //重头开始，取2个
    fmt.Println(buf.String())  //变小了
    fmt.Println(string(b))   //打印he
}
```