package main

import (
	"net"
	"syscall"
	_"encoding/json"
)

func main() {
	 // 创建一个Unix domain soceket
    socket, _ := net.Listen("unix", "/tmp/keyword_match.sock")
    // 关闭时删除绑定的文件
    defer syscall.Unlink("/tmp/keyword_match.sock") 
    // 无限循环监听和受理客户端请求
    for {
        client, _ := socket.Accept()
        
        buf := make([]byte, 1024)
        data_len, _ := client.Read(buf)
        // data := buf[0:data_len]
        // msg := string(data)
        
        // matched := trie.Match(tree, msg)

        response := []byte("[]") // 给响应一个默认值
        // if len(matched) > 0 {
        //     json_str, _ := json.Marshal(matched)
        //     response = []byte(string(json_str))
        // }

        response = []byte(string("123456 length"+string(data_len)))
        _, _ = client.Write(response)
    }
}




