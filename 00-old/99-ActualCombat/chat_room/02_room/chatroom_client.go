package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer conn.Close()


	// 发送消息到服务端
	// conn.Write([]byte("hello 123"))
	go MessageSend(conn)


	// 处理服务端返回的消息
	MessageReceive(conn)
}

// 发送消息到服务端
func MessageSend(conn net.Conn) {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}

		_, err := conn.Write([]byte(input))
		if err != nil {
			conn.Close()
			fmt.Println("client connect fail" + err.Error())
			break
		}

		fmt.Println("client send message：" + input)
	}
}

// 接收服务端返回的消息
func MessageReceive(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		CheckError(err)
		fmt.Println("receive server message : " + string(buf))		
	}
}

// 错误处理
func CheckError(err error) {
	if err != nil {
		// fmt.Println("Error: %s", err.Error())
		// os.Exit(1) // 退出并传递退出码，0是成功，非0是失败

		panic(err)
	}
}



























