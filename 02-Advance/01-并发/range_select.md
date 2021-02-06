

range channel
------------
Go提供了range关键字，将其使用在channel上时，会自动等待channel的动作一直到channel被关闭

正确的使用步骤

- a)发送器一旦停止发送数据后立即关闭channel
- b)接收器一旦停止接收内容，终止程序
- c)移除time.Sleep语句

```golang
package main

import (
    "fmt"
    "strconv"
)

func makeCakeAndSend(cs chan string, count int) {
    for i := 1; i <= count; i++ {
        cakeName := "Strawberry Cake " + strconv.Itoa(i)
        cs <- cakeName //send a strawberry cake
    }
    close(cs)
}

func receiveCakeAndPack(cs chan string) {
    for s := range cs {
        fmt.Println("Packing received cake: ", s)
    }
}

func main() {
    cs := make(chan string)
    go makeCakeAndSend(cs, 5)
    receiveCakeAndPack(cs)
}
```


select channel
------------
select关键字用于多个channel的结合，这些channel会通过类似于are-you-ready polling的机制来工作。select中会有case代码块，用于发送或接收数据——不论通过<-操作符指定的发送还是接收操作准备好时，channel也就准备好了。在select中也可以有一个default代码块，其一直是准备好的。那么，在select中，哪一个代码块被执行的算法大致如下：

- 检查每个case代码块
- 如果任意一个case代码块准备好发送或接收，执行对应内容
- 如果多余一个case代码块准备好发送或接收，随机选取一个并执行对应内容
- 如果任何一个case代码块都没有准备好，等待
- 如果有default代码块，并且没有任何case代码块准备好，执行default代码块对应内容

在下面的程序中，我们扩展蛋糕制作工厂来模拟多于一种口味的蛋糕生产的情况——现在有草莓和巧克力两种口味！但是装箱机制还是同以前一样的。由于蛋糕来自不同的channel，而装箱器不知道确切的何时会有何种蛋糕放置到某个或多个channel上，这就可以用select语句来处理所有这些情况——一旦某一个channel准备好接收蛋糕/数据，select就会完成该对应的代码块内容

注意，我们这里使用的多个返回值case cakeName, strbry_ok := <-strbry_cs，第二个返回值是一个bool类型，当其为false时说明channel被关闭了。如果是true，说明有一个值被成功传递了。我们使用这个值来判断是否应该停止等待。


```golang
package main

import (
    "fmt"
    "strconv"
)

func makeCakeAndSend(cs chan string, flavor string, count int) {
    for i := 1; i <= count; i++ {
        cakeName := flavor + " Cake " + strconv.Itoa(i)
        cs <- cakeName //send a strawberry cake
    }
    close(cs)
}

func receiveCakeAndPack(strbry_cs chan string, choco_cs chan string) {
    strbry_closed, choco_closed := false, false

    for {
        //if both channels are closed then we can stop
        if strbry_closed && choco_closed {
        	fmt.Println("no new cake ...")
            return
        }

        fmt.Println("Waiting for a new cake ...")
        select {
        case cakeName, strbry_ok := <-strbry_cs:
            if !strbry_ok {
                strbry_closed = true
                fmt.Println(" ... Strawberry channel closed!")
            } else {
                fmt.Println("Received from Strawberry channel.  Now packing", cakeName)
            }
        case cakeName, choco_ok := <-choco_cs:
            if !choco_ok {
                choco_closed = true
                fmt.Println(" ... Chocolate channel closed!")
            } else {
                fmt.Println("Received from Chocolate channel.  Now packing", cakeName)
            }
        }
    }
}

func main() {
    strbry_cs := make(chan string)
    choco_cs := make(chan string)

    //two cake makers
    go makeCakeAndSend(choco_cs, "Chocolate", 3)   //make 3 chocolate cakes and send
    go makeCakeAndSend(strbry_cs, "Strawberry", 4) //make 3 strawberry cakes and send

    //one cake receiver and packer
    receiveCakeAndPack(strbry_cs, choco_cs) //pack all cakes received on these cake channels
}
```

