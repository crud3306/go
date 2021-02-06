

基本的GET请求
--------------
```golang
//基本的GET请求
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    resp, err := http.Get("http://httpbin.org/get")
    // 带参数
    // resp, err := http.Get("http://httpbin.org/get?name=zhaofan&age=23")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    fmt.Println(string(body))
    fmt.Println(resp.StatusCode)
    if resp.StatusCode == 200 {
        fmt.Println("ok")
    }
}
```


带参数的Get请求
--------------
```golang
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func main(){
    resp, err := http.Get("http://httpbin.org/get?name=zhaofan&age=23")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

但是如果我们想要把一些参数做成变量而不是直接放到url中怎么操作，代码例子如下：
```golang
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

func main(){
    params := url.Values{}
    Url, err := url.Parse("http://httpbin.org/get")
    if err != nil {
        return
    }
    params.Set("name","zhaofan")
    params.Set("age","23")
    //如果参数中有中文参数,这个方法会进行URLEncode
    Url.RawQuery = params.Encode()
    urlPath := Url.String()
    fmt.Println(urlPath) // https://httpbin.org/get?age=23&name=zhaofan

    resp,err := http.Get(urlPath)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```



解析JSON类型的返回结果
--------------
```golang
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type result struct {
    Args string `json:"args"`
    Headers map[string]string `json:"headers"`
    Origin string `json:"origin"`
    Url string `json:"url"`
}

func main() {
    resp, err := http.Get("http://httpbin.org/get")
    if err != nil {
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))

    var res result
    _ = json.Unmarshal(body,&res)
    fmt.Printf("%#v", res)
}
```


GET请求添加请求头
--------------
```golang
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    client := &http.Client{}
    req,_ := http.NewRequest("GET","http://httpbin.org/get",nil)
    req.Header.Add("name","zhaofan")
    req.Header.Add("age","3")

    resp,_ := client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf(string(body))
}


//从上述的结果可以看出我们设置的头是成功了：
{
  "args": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Age": "3", 
    "Host": "httpbin.org", 
    "Name": "zhaofan", 
    "User-Agent": "Go-http-client/1.1"
  }, 
  "origin": "211.138.20.170, 211.138.20.170", 
  "url": "https://httpbin.org/get"
}
```


golang 发起POST请求
--------------
基本的POST使用
```golang
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

func main() {
    urlValues := url.Values{}
    urlValues.Add("name","zhaofan")
    urlValues.Add("age","22")

    resp, _ := http.PostForm("http://httpbin.org/post", urlValues)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}



//结果如下：
{
  "args": {}, 
  "data": "", 
  "files": {}, 
  "form": {
    "age": "22", 
    "name": "zhaofan"
  }, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Content-Length": "19", 
    "Content-Type": "application/x-www-form-urlencoded", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/1.1"
  }, 
  "json": null, 
  "origin": "211.138.20.170, 211.138.20.170", 
  "url": "https://httpbin.org/post"
}
```


post的另外一种方式
--------------
```golang
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

func main() {
    urlValues := url.Values{
        "name":{"zhaofan"},
        "age":{"23"},
    }
    reqBody:= urlValues.Encode()
    resp, _ := http.Post("http://httpbin.org/post", "text/html", strings.NewReader(reqBody))
    body,_:= ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}

//结果如下：
{
  "args": {}, 
  "data": "age=23&name=zhaofan", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Content-Length": "19", 
    "Content-Type": "text/html", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/1.1"
  }, 
  "json": null, 
  "origin": "211.138.20.170, 211.138.20.170", 
  "url": "https://httpbin.org/post"
}
```

发送JSON数据的post请求
```golang
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    client := &http.Client{}
    data := make(map[string]interface{})
    data["name"] = "zhaofan"
    data["age"] = "23"
    bytesData, _ := json.Marshal(data)

    req, _ := http.NewRequest("POST", "http://httpbin.org/post", bytes.NewReader(bytesData))
    // set头部 
    // req.Header.Add("content-type", "application/json")
    resp, _ := client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))

}


//结果如下：
{
  "args": {}, 
  "data": "{\"age\":\"23\",\"name\":\"zhaofan\"}", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Content-Length": "29", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/1.1"
  }, 
  "json": {
    "age": "23", 
    "name": "zhaofan"
  }, 
  "origin": "211.138.20.170, 211.138.20.170", 
  "url": "https://httpbin.org/post"
}
```


不用client的post请求
--------------
```golang
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    data := make(map[string]interface{})
    data["name"] = "zhaofan"
    data["age"] = "23"
    bytesData, _ := json.Marshal(data)
    
    resp, _ := http.Post("http://httpbin.org/post","application/json", bytes.NewReader(bytesData))
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```
