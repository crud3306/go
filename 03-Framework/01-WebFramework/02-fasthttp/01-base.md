
fasthttp是一个高性能的web server框架。Golang官方的net/http性能相比fasthttp逊色很多。根据测试，fasthttp的性能可以达到net/http的10倍。所以，在一些高并发的项目中，我们经常用fasthttp来代替net/http。

地址：  
https://github.com/valyala/fasthttp


fasthttp性能之所以提高很多，是因为它使用了一个ctxPool来维护RequestCtx，每次请求都先去ctxPool中获取。如果能获取到就用池中已经存在的，如果获取不到，new出一个新的RequestCtx。这也就是fasthttp性能高的一个主要原因，复用RequestCtx可以减少创建对象所有的时间以及减少内存使用率。


注意一个问题：
如果在高并发的场景下，如果整个请求链路中有另起的goroutine，前一个RequestCtx处理完成业务逻辑以后(另起的协程还没有完成)，立刻被第二个请求使用，那就会发生前文所述的错乱的request body。 (复用RequestCtx造成的)


快速开始
----------------------
```golang
package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	// 创建路由
	router := fasthttprouter.New()

	// 不同的路由执行不同的处理函数
	router.GET("/", Index)

	router.GET("/hello", Hello)

	router.GET("/get", TestGet)

	// post方法
	router.POST("/post", TestPost)

	// 启动web服务器，监听 0.0.0.0:12345
	log.Fatal(fasthttp.ListenAndServe(":12345", router.Handler))
}

// index 页
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome")
}

// 简单路由页
func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello")
}

// 获取GET请求json数据
// 使用 ctx.QueryArgs() 方法
// Peek类似与python中dict的pop方法，取某个键对应的值
func TestGet(ctx *fasthttp.RequestCtx) {
	values := ctx.QueryArgs()
	str := fmt.Sprintf("param a:%s", string(values.Peek("a")))
	fmt.Fprint(ctx, str) // 不加string返回的byte数组

}

// 获取post的请求json数据
// 这里就有点坑是，查了很多网页说可以用 ctx.PostArgs() 取post的参数，返现不行，返回空
// 后来用 ctx.FormValue() 取表单数据就好了，难道是版本升级的问题？
// ctx.PostBody() 在上传文件的时候比较有用
func TestPost(ctx *fasthttp.RequestCtx) {
	//postValues := ctx.PostArgs()
	//fmt.Fprint(ctx, string(postValues))

	// 获取表单数据
	str := fmt.Sprintf("param a:%s", string(ctx.FormValue("a")))
	fmt.Fprint(ctx, str)

	// 这两行可以获取PostBody数据，在上传数据文件的时候有用
	postBody := ctx.PostBody()
	fmt.Fprint(ctx, string(postBody))
}
```
