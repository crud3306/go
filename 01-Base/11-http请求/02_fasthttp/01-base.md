

get
--------------
```golang
package main

import (
    "github.com/valyala/fasthttp"
)

func main() {
    url := `http://httpbin.org/get`

    status, resp, err := fasthttp.Get(nil, url)
    if err != nil {
        fmt.Println("请求失败:", err.Error())
        return
    }

    if status != fasthttp.StatusOK {
        fmt.Println("请求没有成功:", status)
        return
    }

    fmt.Println(string(resp))
}
````


post
--------------
```golang
func main() {
    url := `http://httpbin.org/post?key=123`
    
    // 填充表单，类似于net/url
    args := &fasthttp.Args{}
    args.Add("name", "test")
    args.Add("age", "18")

    status, resp, err := fasthttp.Post(nil, url, args)
    if err != nil {
        fmt.Println("请求失败:", err.Error())
        return
    }

    if status != fasthttp.StatusOK {
        fmt.Println("请求没有成功:", status)
        return
    }

    fmt.Println(string(resp))
}
```

上面两个简单的demo实现了get和post请求，这两个方法也已经实现了自动重定向，那么如果有更复杂的请求或需要手动重定向，需要怎么处理呢？比如有些API服务需要我们提供json body或者xml body也就是，Content-Type是application/json、application/xml或者其他类型。


继续看下面：
--------------
```golang
func main() {
    url := `http://httpbin.org/post?key=123`
    
    req := &fasthttp.Request{}
    req.SetRequestURI(url)
    
    requestBody := []byte(`{"request":"test"}`)
    req.SetBody(requestBody)

    // 默认是application/x-www-form-urlencoded
    req.Header.SetContentType("application/json")
    req.Header.SetMethod("POST")

    resp := &fasthttp.Response{}

    client := &fasthttp.Client{}
    if err := client.Do(req, resp);err != nil {
        fmt.Println("请求失败:", err.Error())
        return
    }

    b := resp.Body()

    fmt.Println("result:\r\n", string(b))
}
```


搞定，这样就完成了，but还有优化的空间有木有？
有文章说到它的高性能主要源自于“复用”，通过服务协程和内存变量的复用，节省了大量资源分配的成本。
好吧。。。 继续优化下。。
翻阅文档发现了他提供了几个方法：AcquireRequest()、AcquireResponse(),而且也推荐了在有性能要求的代码中，通过 AcquireRequest 和 AcquireResponse 来获取 req 和 resp。


AcquireRequest、AcquireResponse
--------------
```golang
func main() {
    url := `http://httpbin.org/post?key=123`

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
    
    // 默认是application/x-www-form-urlencoded
    req.Header.SetContentType("application/json")
    req.Header.SetMethod("POST")
    
    req.SetRequestURI(url)
    
    requestBody := []byte(`{"request":"test"}`)
    req.SetBody(requestBody)

    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

    if err := fasthttp.Do(req, resp); err != nil {
        fmt.Println("请求失败:", err.Error())
        return
    }

    b := resp.Body()

    fmt.Println("result:\r\n", string(b))
}
```
经过这样简单的改动，性能上确实是增加了一些。

