package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer conn.Close()

	conn.Write([]byte("hello 123"))

	fmt.Println("has sent the message")
}

// 错误处理
func CheckError(err error) {
	if err != nil {
		// fmt.Println("Error: %s", err.Error())
		// os.Exit(1) // 退出并传递退出码，0是成功，非0是失败

		panic(err)
	}
}



























