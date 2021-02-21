用pprof分析gRPC服务的性能

这次我们的主要内容是如何使用pprof工具对gRPC服务的程序性能进行分析。关于gRPC这个框架的文章之前已经写过不少文章了，如果你对它还不太熟悉，不知道它是用来干什么的，可以通过gRPC入门系列的文章对它先做个了解。


怎么用pprof分析gRPC的性能
=================
gRPC底层基于HTTP协议的，一个典型的gRPC服务的启动程序可能像下面这样
```golang
func main () {
  lis, err := net.Listen("tcp", 10000)
  grpcServer := grpc.NewServer()
  pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
  grpcServer.Serve(lis)
}
```
它是一个RPC框架不是Web框架，不支持浏览器用URL访问，所以也就没法向上一节给Echo和Gin框架单独注册pprof采集数据用的那些路由。但是我们可以换个角度来看这个问题，pprof做CPU分析原理是按照一定的频率采集程序CPU（包括寄存器）的使用情况，确定应用程序在主动消耗 CPU 周期时花费时间的位置。

所以我们可以在gRPC服务启动时，异步启动一个监听其他端口的HTTP服务，通过这个HTTP服务间接获取gRPC服务的分析数据。
```golang
go func() {
   http.ListenAndServe(":10001", nil)
}()
```
由于使用默认的ServerMux（服务复用器），所以只要匿名导入net/http/pprof包，这个HTTP的复用器默认就会注册pprof相关的路由。


此外建议在启动程序的最开端，调用runtime.SetBlockProfileRate(1)指示对阻塞超过1纳秒的goroutine进行数据采集。

全部代码，类似下面
```golang
func main () {
  runtime.SetBlockProfileRate(1)

  go func() {
    http.ListenAndServe(":10001", nil)
  }()
  

  lis, err := net.Listen("tcp", 10000)
  grpcServer := grpc.NewServer()
  pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
  grpcServer.Serve(lis)
}
```

服务启动后就能通过{server_ip}:10001/debug/pprof/profile采集CPU的使用情况了，具体pprof工具的使用方法的详细说明参考系列的第一篇文章。

虽然是用另外一个端口的HTTP服务拿到的分析数据，但依然能采集到监听另一个端口的gRPC服务程序的CPU使用情况。



pprof的局限
=================
pprof这些功能虽然很有用，但是想分析出程序的性能问题还是挺费事儿的，从我使用下来的感觉主要有两点。

首先，因为调用图里把所有函数调用都显示出来了，有些耗时长的还是Go底层的runtime包内函数的直接，想要在这一堆里找到慢的业务函数还是得花不少力气。

再一个现在很多服务都是分布式的，如果服务A调用了服务B，服务B里的方法执行的比较耗时的话，在A的分析数据里只能知道grpc.invoke（客户端调用gRPC方法的请求都是由invoke发出的）耗时长，这时又得去服务B上采集数据，做不到全链路服务性能的采集，这块如果谁知道好的解决方案可以在留言里说一下。
