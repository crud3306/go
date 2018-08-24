package main

import "fmt"
import "time"

var quit chan int

func foo(id int) {
    fmt.Println(id)
    // time.Sleep(time.Second) // 停顿一秒
    quit <- 0 // 发消息：我执行完啦！
}


func main() {
    count := 2
    // quit = make(chan int) // 无缓冲
    quit = make(chan int, 1000) // 缓冲1000个数据

    for i := 0; i < count; i++ { //开1000个goroutine
        go foo(i)
    }

    for i :=0 ; i < count; i++ { // 等待所有完成消息发送完毕。
        <- quit
    }
}