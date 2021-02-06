golang 程序启动一个 http 服务时，若服务被意外终止或中断，会让现有请求连接突然中断，未处理完成的任务也会出现不可预知的错误，这样即会造成服务硬终止；为了解决硬终止问题我们希望服务中断或退出时将正在处理的请求正常返回并且等待服务停止前作的一些必要的处理工作。


我们可以看一个硬终止的例子：
-------------
```golang
mux := http.NewServeMux()
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
   time.Sleep(5 * time.Second)
   fmt.Fprintln(w, "Hello world!")
})

server := &http.Server{
   Addr:         ":8080",
   Handler:      mux,
}
server.ListenAndServe()
```

启动服务后，我们可以访问 http://127.0.0.1:8080 页面等待 5s 会输出一个 “Hello world!”， 我们可以尝试 Ctrl+C 终止程序，可以看到浏览器立刻就显示无法连接，这表示连接立刻就中断了，退出前的请求也未正常返回。

在 Golang1.8 以后 http 服务有个新特性 Shutdown 方法可以优雅的关闭一个 http 服务， 该方法需要传入一个 Context 参数，当程序终止时其中不会中断活跃的连接，会等待活跃连接闲置或 Context 终止（手动 cancle 或超时）最后才终止程序，官方文档详见：https://godoc.org/net/http#Server.Shutdown

 

在具体用应用中我们可以配合 signal.Notify 函数来监听系统退出信号来完成程序优雅退出；

特别注意：server.ListenAndServe() 方法在 Shutdown 时会立刻返回，Shutdown 方法会阻塞至所有连接闲置或 context 完成，所以 Shutdown 的方法要写在主 goroutine 中。

 

优雅退出实验1：
-----------------
```golang
func main() {
   mux := http.NewServeMux()
   mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      time.Sleep(5 * time.Second)
      fmt.Fprintln(w, "Hello world!")
   })
   server := &http.Server{
      Addr:         ":8080",
      Handler:      mux,
   }
   go server.ListenAndServe()

   listenSignal(context.Background(), server)
}

func listenSignal(ctx context.Context, httpSrv *http.Server) {
   sigs := make(chan os.Signal, 1)
   signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

   select {
   case <-sigs:
      fmt.Println("notify sigs")
      httpSrv.Shutdown(ctx)
      fmt.Println("http shutdown")
   }
}
```
我们创建了一个 listenSignal 函数来监听程序退出信号 listenSignal 函数中的 select 会一直阻塞直到收到退出信号，然后执行 Shutdown(ctx) 。

可以看到，我们是重新开启了一个 goroutine 来启动 http 服务监听，而 Shutdown(ctx) 在主 goroutine 中，这样才能等待所有连接闲置后再退出程序。

启动上述程序，我们访问  http://127.0.0.1:8080 页面等待 5s 会输出一个 “Hello world!” 在等待期间，我们可以尝试 Ctrl+C 关闭程序，可以看程序控制台会等待输出后才打印 http shutdown 同时浏览器会显示输出内容；而关闭程序之后再新开一个浏览器窗口访问 http://127.0.0.1:8080 则新开的窗口直接断开无法访问。（这些操作需要在 5s 内完成，可以适当调整处理时间方便我们观察实验结果）

通过该实验我们能看到，Shutdown(ctx) 会阻止新的连接进入并等待活跃连接处理完成后再终止程序，达到优雅退出的目的。




如果果不想一直等待，我们也可以设置超时退出。为Shutdown传递一个带超时的ctx即可。

Shutdown(ctx) 除了等待活跃连接的同时也会监听 Context 完成事件，二者有一个触发都会触发程序终止。

我们将代码稍作修改如下：
-----------------
```golang
func main() {
   mux := http.NewServeMux()
   mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      time.Sleep(10 * time.Second)
      fmt.Fprintln(w, "Hello world!")
   })
   server := &http.Server{
      Addr:         ":8080",
      Handler:      mux,
   }
   go server.ListenAndServe()

   listenSignal(context.Background(), server)
}

func listenSignal(ctx context.Context, httpSrv *http.Server) {
   sigs := make(chan os.Signal, 1)
   signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

   select {
   case <-sigs:
      timeoutCtx,_ := context.WithTimeout(ctx, 3*time.Second)
      fmt.Println("notify sigs")
      httpSrv.Shutdown(timeoutCtx)
      fmt.Println("http shutdown")
   }
}
```
我们将 http 服务处理修改成等待 10s， 监听到退出事件后 ctx 修改成 3s 超时的 Context，运行上述程序，然后 Ctrl+C 发送结束信号，我们可以直观的看到，程序在等待 3s 后就终止了，此时即使 http 服务中的处理还没完成，程序也终止了，浏览器中也直接中断连接了。

需要注意的问题：我们在 HandleFunc 中编写的处理逻辑都是在主 goroutine 中完成的和 Shotdown 方法是一个同步操作，因此 Shutdown(ctx) 会等待完成，如果我们的处理逻辑是在新的 goroutine 中或是一个像 Websock 这样的长连接，则Shutdown(ctx) 不会等待处理完成，如果需要解决这类问题还是需要利用 sync.WaitGroup 来进行同步等待。


