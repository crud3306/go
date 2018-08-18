package main

import (
	"fmt"
	"net"
	_ "os"
	_ "strings"
)

// 用于存储在线连接
var onlineConns = make(map[string]net.Conn)

var messageQueue chan string = make(chan string, 1000)
var quitChan = make(chan bool)

func main() {

	listen_socket, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listen_socket.Close()

	fmt.Println("server is start and listen 127.0.0.1:8080")

	// 消费接收的消息
	go ConsumeMessage()

	// 接收消息
	for {
		conn, err := listen_socket.Accept()
		CheckError(err)

		// 将conn存起来
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		onlineConns[addr] = conn

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
			message := string(buf[:numOfBytes])

			messageQueue <- message

			remoteAddr := conn.RemoteAddr()
			fmt.Println(remoteAddr)
			fmt.Printf("has received message: %s\n", string(buf[:numOfBytes]))

		}
	}
}

// 消费消息
func ConsumeMessage() {
	for {
		select {
			case message := <-messageQueue:
				// 对消息进行解析
				doProcessMessage(message)

			case <- quitChan:
				break
		}
	}
	
}

func doProcessMessage(message string) {
	// 群发
	for _, v := range(onlineConns) {
		_, err := v.Write([]byte(message))
		if err != nil {
			fmt.Println("server send message error")
		}
	}

	// arr := strings.Split(message, "#")
	// if len(arr) > 1 {
	// 	addr := strings.Trim(arr[0], " ")
	// 	sendMessage := arr[1]

	// 	if conn, ok := onlineConns[addr]; ok {
	// 		_, err := conn.Write([]byte(sendMessage))
	// 		if err != nil {
	// 			fmt.Println("server send message error")
	// 		}
	// 	}
	// }
}

// 错误处理
func CheckError(err error) {
	if err != nil {
		// fmt.Println("Error: %s", err.Error())
		// os.Exit(1) // 退出并传递退出码，0是成功，非0是失败

		panic(err)
	}
}




























