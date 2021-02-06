
Golang标准库http包提供了基础的http服务，这个服务又基于Handler接口和ServeMux结构的叫做Mutilpexer。



基本使用
=================

绑定handle两种方式：

- HandleFunc
	绑定一个func

- Handle
	绑定一个结构体


示例
-----------------
```golang
package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {

	// HandleFunc示例
	http.HandleFunc("/", HelloWorldHandler)
	http.HandleFunc("/user/login", UserLoginHandler)

	// Handle示例
	http.Handle("/count", new(countHandler))

	//监听服务
	if err := http.ListenAndServe("0.0.0.0:8880", nil); err != nil {
		fmt.Println("服务器错误", err)
	}
}


// HelloWorldHandler ...
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.Method = ", r.Method)
	fmt.Println("r.URL = ", r.URL)
	fmt.Println("r.Header = ", r.Header)
	fmt.Println("r.Body = ", r.Body)
	fmt.Fprintf(w, "HelloWorld!")
}

// UserLoginHandler ...
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler Hello")
	fmt.Fprintf(w, "Login Success")
	// io.WriteString(w, "UserLoginHandler #1!")
}


// 标准库http提供了Handler接口，用于开发者实现自己的handler。只要实现接口的ServeHTTP方法即可。
// countHandler ...
type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
}
```
当代码中不显示的创建serveMux对象，http包就默认创建一个DefaultServeMux对象用来做路由管理器mutilplexer。


实际上，go的作者设计Handler这样的接口，不仅提供了默认的ServeMux对象，开发者也可以自定义ServeMux对象。

本质上ServeMux只是一个路由管理器，而它本身也实现了Handler接口的ServeHTTP方法。 因此围绕Handler接口的方法ServeHTTP，可以轻松的写出go中的中间件。



自定义ServeMux
================
```golang 
func main() {
    mux := http.NewServeMux()
 
    mux.Handle("/", &indexHandler{})
    mux.Handle("/text", &textHandler{"TextHandler !"})
 
    http.ListenAndServe(":8000", mux)
}

// textHandler ...
type textHandler struct {
    responseText string
}
 
func (th *textHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, th.responseText)
}
 
// indexHandler ...
type indexHandler struct {}
 
func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
 
    html := `<doctype html>
        <html>
        <head>
          <title>Hello World</title>
        </head>
        <body>
        <p>
          <a href="/welcome">Welcome</a> |  <a href="/message">Message</a>
        </p>
        </body>
</html>`
    fmt.Fprintln(w, html)
}
```

上面自定义了两个handler结构，都实现了ServeHTTP方法。

我们知道，NewServeMux可以创建一个ServeMux实例，ServeMux同时也实现了ServeHTTP方法，因此代码中的mux也是一种handler。把它当作参数传给http.ListenAndServe方法，后者会把mux传给Server实例。因为指定了handler，因此整个http服务就不再是DefaultServeMux，而是mux，无论是在注册路由还是提供请求服务的时候。

有一点值得注意，这里并没有使用HandleFunc注册路由，而是直接使用了mux注册路由。当没有指定mux的时候，系统需要创建一个默认的defaultServeMux，此时我们已经有了mux，因此不需要http.HandleFunc方法了，直接使用mux的Handle方法注册即可。

此外，Handle第二个参数是一个Handler(处理器),并不是HandleFunc的一个handler函数，其原因是mux.Handle本质上就需要绑定url的pattern模式和handler(处理器)即可。既然indexHandler是handle(处理器),当然就能作为参数，一切请求的处理过程，都交给实现的接口方法ServeHTTP就行了。


 

自定义Server
==================
默认的DefaultServeMux创建的判断来自server对象，如果server对象不提供handler，才会使用默认的serveMux对象。既然ServeMux可以自定义，那么Server对象一样可以。


使用http.Server即可创建自定义的server对象
------------------
```golang
package main

import (
	"fmt"
	"net/http"
	_ "io"
)

func main(){
	http.HandleFunc("/", index)

	server := &http.Server{
	    Addr: ":8000",
	    ReadTimeout: 60 * time.Second,
	    WriteTimeout: 60 * time.Second,
	}
	server.ListenAndServe()
}

// index ...
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler index")
	fmt.Fprintf(w, "hello index")
	// io.WriteString(w, "hello index")
}
```


自定义的serverMux对象也可以传到server对象中
-----------------
```golang
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", index)
 
    server := &http.Server{
        Addr: ":8000",
        ReadTimeout: 60 * time.Second,
        WriteTimeout: 60 * time.Second,
        Handler: mux,
    }
    server.ListenAndServe()
}
```
可见go中的路由和处理函数之间关系非常密切，同时又很灵活。




中间件Middleware
==================
中间件就是连接上下级不同功能的函数或软件，通常进行一些包裹函数的行为，为被包裹函数提供一些功能或行为。

go的http中间件很简单，只要实现一个函数签名为func(http.Handler) http.Handler的函数即可。http.Handler是一个接口，接口方法我们熟悉的为serveHTTP。返回也是一个handler。因为go中的函数也可以当作变量传递或者返回，只要这个函数是一个handler即可，即实现或被handlerFunc包裹成handler处理器。
```golang
func middlewareHandler(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        // 执行handler之前的逻辑
        next.ServeHTTP(w, r)
        // 执行完毕handler后的逻辑
    })
}
```

go的实现实例：
```golang
func loggingHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
        log.Printf("Comleted %s in %v", r.URL.Path, time.Since(start))
    })
}
 
func main() {
    http.Handle("/", loggingHandler(http.HandlerFunc(index)))
    http.ListenAndServe(":8000", nil)
}
```


既然中间件是一种函数，并且签名都是一样，那么很容易就联想到函数一层包一层的中间件。再添加一个函数，然后修改main函数：
```golang
func hook(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("before hook")
        next.ServeHTTP(w, r)
        log.Println("after hook")
 
    })
}
 
func main() {
    http.Handle("/", hook(loggingHandler(http.HandlerFunc(index))))
    http.ListenAndServe(":8000", nil)
}
```
函数调用形成一条链，可以在这条链上做很多事情。


