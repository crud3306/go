select
===========
golang 的 select 就是监听 IO 操作，当 IO 操作发生时，触发相应的动作。
在执行select语句的时候，运行时系统会自上而下地判断每个case中的发送或接收操作是否可以被立即执行【立即执行：意思是当前Goroutine不会因此操作而被阻塞，还需要依据通道的具体特性(缓存或非缓存)】

- 每个case语句里必须是一个IO操作
- 所有channel表达式都会被求值、所有被发送的表达式都会被求值
- 如果任意某个case可以进行，它就执行(其他被忽略)。
- 如果有多个case都可以运行，Select会随机公平地选出一个执行(其他不会执行)。
- 如果有default子句，case不满足条件时执行该语句。
- 如果没有default字句，select将阻塞，直到某个case可以运行；Go不会重新对channel或值进行求值。


select 语句用法
---------------
注意到 select 的代码形式和 switch 非常相似， 不过 select 的 case 里的操作语句只能是【IO 操作】 。
此示例里面 select 会一直等待等到某个 case 语句完成， 也就是等到成功从 ch1 或者 ch2 中读到数据,如果都不满足条件且存在default case, 那么default case会被执行。 则 select 语句结束。
示例：
```golang
package main

import (
    "fmt"
)

func main(){
    ch1 := make(chan int, 1)
    ch2 := make(chan int, 1)

    select {
        case e1 := <-ch1:
        //如果ch1通道成功读取数据，则执行该case处理语句
            fmt.Printf("1th case is selected. e1=%v",e1)
        case e2 := <-ch2:
        //如果ch2通道成功读取数据，则执行该case处理语句
            fmt.Printf("2th case is selected. e2=%v",e2)
        default:
        //如果上面case都没有成功，则进入default处理流程
            fmt.Println("default!.")
    }
}
```


select分支选择规则
-------------------
所有跟在case关键字右边的发送语句或接收语句中的通道表达式和元素表达式都会先被求值。无论它们所在的case是否有可能被选择都会这样。

求值顺序：自上而下、从左到右
示例：
```golang
package main

import (
    "fmt"
)
//定义几个变量，其中chs和numbers分别代表了包含了有限元素的通道列表和整数列表
var ch1 chan int
var ch2 chan int
var chs = []chan int{ch1, ch2}
var numbers = []int{1,2,3,4,5}

func main(){
    select {
        case getChan(0) <- getNumber(2):
            fmt.Println("1th case is selected.")
        case getChan(1) <- getNumber(3):
            fmt.Println("2th case is selected.")
        default:
            fmt.Println("default!.")
    }
}

func getNumber(i int) int {
    fmt.Printf("numbers[%d]\n", i)
    return numbers[i]
}

func getChan(i int) chan int {
    fmt.Printf("chs[%d]\n", i)
    return chs[i]
}
```

输出：
```sh
chs[0]
numbers[2]
chs[1]
numbers[3]
default!.
```
可以看出求值顺序。满足自上而下、自左而右这条规则。



随机执行case
----------------
如果同时有多个case满足条件，通过一个伪随机的算法决定哪一个case将会被执行。
示例：
```golang
package main

import (
    "fmt"
)
func main(){
    chanCap := 5
    ch7 := make(chan int, chanCap)

    for i := 0; i < chanCap; i++ {
        select {
            case ch7 <- 1:
            case ch7 <- 2:
            case ch7 <- 3:
        }
    }

    for i := 0; i < chanCap; i++ {
        fmt.Printf("%v\n", <-ch7)
    }
}
```
输出：(注：每次运行都会不一样)
```sh
3
3
2
3
1
```

一些惯用手法示例
--------------
示例一：单独启用一个Goroutine执行select,等待通道关闭后结束循环
```golang
package main

import (
    "fmt"
    "time"
)
func main(){
    //初始化通道
    ch11 := make(chan int, 1000)
    sign := make(chan int, 1)

    //给ch11通道写入数据
    for i := 0; i < 1000; i++ {
        ch11 <- i
    }
    //关闭ch11通道
    close(ch11)

    //单独起一个Goroutine执行select
    go func(){
        var e int
        ok := true

        for{
            select {
                case e,ok = <- ch11:
                if !ok {
                    fmt.Println("End.")
                    break
                }
                fmt.Printf("ch11 -> %d\n",e)
            }

            //通道关闭后退出for循环
            if !ok {
                sign <- 0
                break
            }
        }

    }()

    //惯用手法，读取sign通道数据，为了等待select的Goroutine执行。
    <- sign
}

// 输出
ch11 -> 0
ch11 -> 1
…
ch11 -> 999
End.
```


示例二：加以改进，我们不想等到通道被关闭后再退出循环，利用一个辅助通道模拟出操作超时。
```golang
package main

import (
    "fmt"
    "time"
)

func main(){
    //初始化通道
    ch11 := make(chan int, 1000)
    sign := make(chan int, 1)

    //给ch11通道写入数据
    for i := 0; i < 1000; i++ {
        ch11 <- i
    }
    //关闭ch11通道
    close(ch11)

    //我们不想等到通道被关闭之后再推出循环，我们创建并初始化一个辅助的通道，利用它模拟出操作超时行为
    timeout := make(chan bool,1)
    go func(){
        time.Sleep(time.Millisecond) //休息1ms
        timeout <- false
    }()

    //单独起一个Goroutine执行select
    go func(){
        var e int
        ok := true

        for{
            select {
                case e,ok = <- ch11:
                    if !ok {
                        fmt.Println("End.")
                        break
                    }
                    fmt.Printf("ch11 -> %d\n",e)
                case ok = <- timeout:
                //向timeout通道发送元素false后，该case几乎马上就会被执行, ok = false
                    fmt.Println("Timeout.")
                    break
            }

            //终止for循环
            if !ok {
                sign <- 0
                break
            }
        }

    }()

    //惯用手法，读取sign通道数据，为了等待select的Goroutine执行。
    <- sign
}
ch11 -> 0
ch11 -> 1
…
ch11 -> 691
Timeout.
```



示例三：上面实现了单个操作的超时,但是那个超时触发器开始计时有点早。
```golang
package main

import (
    "fmt"
    "time"
)
func main(){
    //初始化通道
    ch11 := make(chan int, 1000)
    sign := make(chan int, 1)

    //给ch11通道写入数据
    for i := 0; i < 1000; i++ {
        ch11 <- i
    }
    //关闭ch11通道
    //close(ch11),为了看效果先注释掉

    //单独起一个Goroutine执行select
    go func(){
        var e int
        ok := true

        for{
            select {
                case e,ok = <- ch11:
                    if !ok {
                        fmt.Println("End.")
                        break
                    }
                    fmt.Printf("ch11 -> %d\n",e)
                case ok = <- func() chan bool {
                    //经过大约1ms后，该接收语句会从timeout通道接收到一个新元素并赋值给ok,从而恰当地执行了针对单个操作的超时子流程，恰当地结束当前for循环
                    timeout := make(chan bool,1)
                    go func(){
                        time.Sleep(time.Millisecond)//休息1ms
                        timeout <- false
                    }()
                    return timeout
                }():
                    fmt.Println("Timeout.")
                    break
            }
            //终止for循环
            if !ok {
                sign <- 0
                break
            }
        }

    }()

    //惯用手法，读取sign通道数据，为了等待select的Goroutine执行。
    <- sign
}
ch11 -> 0
ch11 -> 1
…
ch11 -> 999
Timeout.
```


非缓冲的Channel
---------------
我们在初始化一个通道时将其容量设置成0，或者直接忽略对容量的设置，那么就称之为非缓冲通道
```golang
ch1 := make(chan int, 1) //缓冲通道
ch2 := make(chan int, 0) //非缓冲通道
ch3 := make(chan int) //非缓冲通道
```
- 向此类通道发送元素值的操作会被阻塞，直到至少有一个针对该通道的接收操作开始进行为止。
- 从此类通道接收元素值的操作会被阻塞，直到至少有一个针对该通道的发送操作开始进行为止。
- 针对非缓冲通道的接收操作会在与之相应的发送操作完成之前完成。

对于第三条要特别注意，发送操作在向非缓冲通道发送元素值的时候，会等待能够接收该元素值的那个接收操作。并且确保该元素值被成功接收，它才会真正的完成执行。而缓冲通道中，刚好相反，由于元素值的传递是异步的，所以发送操作在成功向通道发送元素值之后就会立即结束(它不会关心是否有接收操作)。


示例一

实现多个Goroutine之间的同步
```golang
package main

import (
    "fmt"
    "time"
)

func main(){
    unbufChan := make(chan int)
    //unbufChan := make(chan int, 1) 有缓冲容量

    //启用一个Goroutine接收元素值操作
    go func(){
        fmt.Println("Sleep a second...")
        time.Sleep(time.Second)//休息1s
        num := <- unbufChan //接收unbufChan通道元素值
        fmt.Printf("Received a integer %d.\n", num)
    }()

    num := 1
    fmt.Printf("Send integer %d...\n", num)
    //发送元素值
    unbufChan <- num
    fmt.Println("Done.")
}
```

缓冲channel输出结果如下：
```sh
Send integer 1…
Done.
```

非缓冲channel输出结果如下：
```sh
Send integer 1…
Sleep a second…
Received a integer 1.
Done.
```
在非缓冲Channel中，从打印数据可以看出主Goroutine中的发送操作在等待一个能够与之配对的接收操作。配对成功后，元素值1才得以经由unbufChan通道被从主Goroutine传递至那个新的Goroutine.



select与非缓冲通道
------------
与操作缓冲通道的select相比，它被阻塞的概率一般会大很多。只有存在可配对的操作的时候，传递元素值的动作才能真正的开始。

示例：

发送操作间隔1s,接收操作间隔2s
分别向unbufChan通道发送小于10和大于等于10的整数，这样更容易从打印结果分辨出配对的时候哪一个case被选中了。下列案例两个case是被随机选择的。
```golang
package main

import (
    "fmt"
    "time"
)

func main(){
    unbufChan := make(chan int)
    sign := make(chan byte, 2)

    go func(){
        for i := 0; i < 10; i++ {
            select {
                case unbufChan <- i:
                case unbufChan <- i + 10:
                default:
                    fmt.Println("default!")
            }
            time.Sleep(time.Second)
        }

        close(unbufChan)
        fmt.Println("The channel is closed.")

        sign <- 0
    }()

    go func(){
        loop:
            for {
                select {
                    case e, ok := <-unbufChan:
                    if !ok {
                        fmt.Println("Closed channel.")
                        break loop
                    }
                    fmt.Printf("e: %d\n",e)
                    time.Sleep(2 * time.Second)
                }
            }
            
            sign <- 1
    }()

    <- sign
    <- sign
}

// 输出
default! //无法配对
e: 1
default!//无法配对
e: 3
default!//无法配对
e: 15
default!//无法配对
e: 17
default!//无法配对
e: 9
The channel is closed.
Closed channel.
```
default case会在收发操作无法配对的情况下被选中并执行。在这里它被选中的概率是50%。

上面的示例给予了我们这样一个启发：使用非缓冲通道能够让我们非常方便地在接收端对发送端的操作频率实施控制。
可以尝试去掉default case，看看打印结果，代码稍作修改如下：
```golang
package main

import (
    "fmt"
    "time"
)

func main(){
    unbufChan := make(chan int)
    sign := make(chan byte, 2)

    go func(){
        for i := 0; i < 10; i++ {
            select {
                case unbufChan <- i:
                case unbufChan <- i + 10:

            }
            fmt.Printf("The %d select is selected\n",i)
            time.Sleep(time.Second)
        }
        close(unbufChan)
        fmt.Println("The channel is closed.")
        sign <- 0
    }()

    go func(){
        loop:
            for {
                select {
                    case e, ok := <-unbufChan:
                    if !ok {
                        fmt.Println("Closed channel.")
                        break loop
                    }
                    fmt.Printf("e: %d\n",e)
                    time.Sleep(2 * time.Second)
                }
            }
            sign <- 1
    }()
    <- sign
    <- sign
}

// 输出
e: 0
The 0 select is selected
e: 11
The 1 select is selected
e: 12
The 2 select is selected
e: 3
The 3 select is selected
e: 14
The 4 select is selected
e: 5
The 5 select is selected
e: 16
The 6 select is selected
e: 17
The 7 select is selected
e: 8
The 8 select is selected
e: 19
The 9 select is selected
The channel is closed.
Closed channel.
```
总结：上面两个例子，第一个有default case 无法配对时执行该语句，而第二个没有default case ，无法配对case时select将阻塞，直到某个case可以运行(上述示例是直到unbufChan数据被读取操作)，不会重新对channel或值进行求值。