
channe
-----------
无缓冲的信道在取消息和存消息的时候都会挂起当前的goroutine，除非另一端已经准备好。
如果有缓冲，当缓冲写满时也会阻塞。

主线程main 本身也是一个goroutine，当主线程结束时，其它线程哪怕没跑完，也会结束掉。

来个例子：
```go
package main 

import (
	"fmt"
	_ "time"
)

var ch chan int = make(chan int)

func foo() {
	for j:=0; j<10; j++ {
		fmt.Println("iam foo")
	    cc := <- ch // 从ch取数据，如果ch中还没放数据，那就挂起main线，直到foo函数中放数据为止
	    fmt.Println(cc)
	}
}

func main() {

	fmt.Println(111)    

    go foo()

    fmt.Println(222)    

    ch <- 0 // 向ch中加数据，如果没有其他goroutine来取走这个数据，那么挂起foo, 直到main函数把0这个数据拿走
    

	fmt.Println(333)    
}
```
上面程序执行结果：
```go
111
222
iam foo
0
iam foo
333
```

上面的执行流程是:

main主线开启
执行fmt.Println(111)，输出111
开启一个goroutine foo 但goroutine里面代码暂未执行
执行fmt.Println(222)，输出222
执行 ch<-0 ，主线main阻塞
goroutine foo 开始执行，输出 fmt.Println("iam foo")，接收到 cc := <- ch，输出fmt.Println(cc)，然后执行for的第二轮循环，输出 fmt.Println("iam foo")，再到cc := <- ch ，因为没有数据所以阻塞住
main主线继续执行，输出 fmt.Println(333)   


一定注意的是:  
  
main 中执行到 ch <- 0，会出查找是否有已开启的goroutine可以消费该ch，如果没有则直接报错。如果有主线才会阻塞，让相应的可消费chan的goroutine开始执行。  
如上面的代码改成下面这样(go foo() 放到了 ch <- 0 后面)，则程序在输出111, 222 后，直接报错。  
```go
func main() {
	fmt.Println(111)    
    fmt.Println(222)    

    ch <- 0 // 向ch中加数据，如果没有其他goroutine来取走这个数据，那么挂起foo, 
    直到main函数把0这个数据拿走

    go foo()

	fmt.Println(333)    
}
```
执行结果：
```go
111
222
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	/data/my_git/my_open/go-start/training/goroutine/02.go:26 +0xbd
exit status 2
```
报错原因：正是因为主进程找不到可消费该channe ch <- 0 的goroutine















