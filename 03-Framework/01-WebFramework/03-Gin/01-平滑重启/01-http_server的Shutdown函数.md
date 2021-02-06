

优雅关闭
============

golang 程序启动一个 http 服务时，若服务被意外终止或中断，会让现有请求连接突然中断，未处理完成的任务也会出现不可预知的错误，这样即会造成服务硬终止；为了解决硬终止问题，我们希望服务中断或退出时，将正在处理的请求正常返回并且等待服务停止前做一些必要的处理工作。


前言
-------------
gin底层使用的是net/http, 所以gin的优雅退出就等于
http.Server的优雅退出, Golang 1.8以后提供了Shutdown函数，可以优雅关闭http.Server
```sh
func (srv *Server) Shutdown(ctx context.Context) error
```

Shutdown(ctx) 会阻止新的连接进入并等待活跃连接处理完成后再终止程序，达到优雅退出的目的。


优雅退出的过程
-------------
1）关闭所有的监听  
2）关闭所有的空闲连接  
3）无限期等待活动的连接处理完毕转为空闲，并关闭。 如果提供了带有超时的Context，将在服务关闭前返回 Context的超时错误



```golang
package main

import (
    "net/http"
    "time"
    "os"
    "os/signal"
    "syscall"
    "fmt"
    "github.com/gin-gonic/gin"
    "context"
)

func SlowHandler(c *gin.Context) {
    fmt.Println("[start] SlowHandler")
    //time.Sleep(30 * time.Second)
    time.Sleep(30 * time.Second)
    fmt.Println("[end] SlowHandler")
    c.JSON(http.StatusOK, gin.H{

        "message": "success",
    })
}


func main() {
    r := gin.Default()
    // 1.
    r.GET("/slow", SlowHandler)


    server := &http.Server{
        Addr:           ":8080",
        Handler:        r,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    go server.ListenAndServe()

    // 设置优雅退出
    gracefulExitWeb(server)
    // gracefulExitWebWithTimeout(server)
}

func gracefulExitWeb(server *http.Server) {
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
    sig := <-ch

    fmt.Println("got a signal", sig)
    now := time.Now()
    err := server.Shutdown(context.Background())
    if err != nil{
        fmt.Println("err", err)
    }

    // 看看实际退出所耗费的时间
    fmt.Println("------exited--------", time.Since(now))
}

func gracefulExitWebWithTimeout(server *http.Server) {
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
    sig := <-ch

    fmt.Println("got a signal", sig)
    now := time.Now()

    // 超时时间到了，如果还没结束，则强行结束
    // 带超时的 Context 是在创建时就开始计时了，因此需要在接收到结束信号后再创建带超时的 Context。
    cxt, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
    err := server.Shutdown(cxt)
    if err != nil{
        fmt.Println("err", err)
    }

    // 看看实际退出所耗费的时间
    fmt.Println("------exited--------", time.Since(now))
}
```