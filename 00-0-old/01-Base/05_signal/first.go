package main

import (
	"fmt"
	_"time"
	"syscall"
	"os"
	"os/signal"
)


// go run 该文件后 
// 1.首先系统接收到ctrl+c的指令，signal接收到该指令。
// 2.signal执行 原先堵塞的 s:=<-c 这步骤，并关闭 shutdown 通道。
// 3.打印 相关消息

func main() {
	shutdown := make(chan struct{})

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		select {
		case c := <-shutdown:
			fmt.Println("shutdown", c)
			return
		}
	}()	

	s := <-c
	close(shutdown)
	fmt.Println("Got signal:", s) 
	// time.Sleep(1*time.Second)
	// time.Sleep(1*time.Microsecond)
}