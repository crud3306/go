
在Echo和Gin框架中使用pprof

前言
=================
这一节的重点会放在如何在Echo和Gin这两个框架中增加对pprof HTTP请求的支持，因为pprof只是提供了对net/http包的ServerMux的路由支持，这些路由想放到Echo和Gin里使用时，还是需要有点额外的集成工作。

等集成到框架里，能通过HTTP访问pprof提供的几个路由后，go tool pprof工具还是通过访问这些URL把性能数拿到本地后来分析的，后续的性能数据采集和分析的操作就跟上篇文章里介绍的完全一样了，并没有因为使用的框架不一样而有什么差别。



在Echo中使用pprof
=================

由于Echo框架使用的复用器ServerMux是自定义的，需要手动注册pprof提供的路由，网上有几个把他们封装成了包可以直接使用， 不过都不是官方提供的包。后来我看了一下pprof提供的路由Handler的源码，只需要把它转换成Echo框架的路由Handler后即可能正常处理那些pprof相关的请求，具体转换操作很简单，代码如下。

```golang
func RegisterRoutes(engine *echo.Echo) {
	router := engine.Group("")
	......
	// 下面的路由根据要采集的数据需求注册，不用全都注册
	router.GET("/debug/pprof", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/debug/pprof/allocs", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/debug/pprof/block", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/debug/pprof/goroutine", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/debug/pprof/heap", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/debug/pprof/mutex", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	router.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	router.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	router.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
}
```

注册好路由后还需要对Echo框架的写响应超时WriteTimeout做一下配置，保证发生写超时的时间设置要大于pprof做数据采集的时间，这个配置对应的是/debug/pprof路由的seconds参数，默认采集时间是30秒，比如我通常要进行60秒的数据采集，那WriteTimeout配置的时间就要超过60秒，具体配置方式如下：


如果pprof做profiling的时间超过WriteTimeout会引发一个 "profile duration exceeds server's WriteTimeout"的错误。
```golang
RegisterRoutes(engine)

err := engine.StartServer(&http.Server{
   Addr:              addr,
   ReadTimeout:       time.Second * 5,
   ReadHeaderTimeout: time.Second * 2,
   WriteTimeout:      time.Second * 90,
})
```

上面两步都设置完后就能够按照上面文件里介绍的pprof子命令进行性能分析了
```sh
➜ go tool pprof http://{server_ip}:{port}/debug/pprof/profile
Fetching profile over HTTP from http://localhost/debug/pprof/profile
Saved profile in /Users/Kev/pprof/pprof.samples.cpu.005.pb.gz
Type: cpu
Time: Nov 15, 2020 at 3:32pm (CST)
Duration: 30.01s, Total samples = 0
No samples were found with the default sample value type.
Try "sample_index" command to analyze different sample values.
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```
具体pprof常用子命令的使用方法，可以参考文章Golang程序性能分析（一）pprof和go-torch里的内容。





在Gin中使用pprof
=================

在Gin框架可以通过安装Gin项目组提供的gin-contrib/pprof包，直接引入后使用就能提供pprof相关的路由访问。
```golang
import "github.com/gin-contrib/pprof"

package main

import (
 "github.com/gin-contrib/pprof"
 "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  pprof.Register(router)
  router.Run(":8080")
}
```
这个包还支持把pprof路由划分到单独的路由组里，具体可以查阅gin-contrib/pprof的文档。




上面这些都配置完，启动服务后就能使用go tool pprof进行数据采集和分析：

内存使用信息采集  
> go tool pprof http://localhost:8080/debug/pprof/heap

CPU使用情况信息采集  
> go tool pprof http://localhost:8080/debug/pprof/profile


