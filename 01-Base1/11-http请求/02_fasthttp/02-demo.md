

简单示例
--------------
```golang
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

func main() {

	body, err := QRequest("http://xxx:8898/post", "POST", nil, "", 3*time.Second)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("obj", obj)
}


// QRequest ...
func QRequest(url string, method string, data interface{}, contentType string, timeOut time.Duration) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 默认是application/x-www-form-urlencoded
	if contentType != "" {
		req.Header.SetContentType("application/json")
	}
	req.Header.SetMethod(method)

	req.SetRequestURI(url)
	// fmt.Println("======= requestBody", url, method, data)

	if data != nil {
		requestBody, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		// fmt.Println("======= requestBody", string(requestBody))
		req.SetBody(requestBody)
	}

	if err := fasthttp.DoTimeout(req, resp, timeOut); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}
```