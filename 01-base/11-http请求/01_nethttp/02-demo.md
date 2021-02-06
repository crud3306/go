

简单示列
-----------------
```golang
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {

	body, err := QRequest("http://xxxx:8898/post", "POST", nil, nil)
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


// QRequest 简单请求
func QRequest(url string, method string, data interface{}, header map[string]string) ([]byte, error) {
	var err error
	var req *http.Request

	// post参数
	if data != nil {
		jsonStr, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		postData := bytes.NewBuffer(jsonStr)
		req, err = http.NewRequest(method, url, postData)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	// 设置header头
	// 默认application/x-www-form-urlencoded
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Referer", "http://localhost:8080")
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	// client := &http.Client{Timeout: 3 * time.Second}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 2)) //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 2,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 视业务决定是否需要
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status err:%d", resp.StatusCode)
	}

	// fmt.Println("response StatusCode:", resp.StatusCode)
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)

	return ioutil.ReadAll(resp.Body)
}
```