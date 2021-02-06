package main

import (
	"fmt"
	"net"
	_ "os"
)

func main() {
	listen_socket, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listen_socket.Close()

	fmt.Println("server is star and listen 127.0.0.1:8080")

	for {
		conn, err := listen_socket.Accept()
		CheckError(err)

		go ProcessInfo(conn)
	}
}

// 处理单个请求
func ProcessInfo(conn net.Conn) {
	buf := make([]byte, 1024)
	defer conn.Close()

	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			break
		}

		if numOfBytes != 0 {
			fmt.Printf("has received message: %s\n", string(buf))
		}
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




























