go语言 grpc入门教程


前置知识点：
---------------
- protobuf语法
- grpc核心概念


go语言使用grpc的步骤：
---------------
- 安装go语言grpc包
- 安装protobuf编译器
- 划分目录结构
- 定义服务
- 使用protobuf编译器，编译proto协议文件，生成go代码。
- 实现服务端代码
- 实现客户端代码

提示：go版本要求：1.11版本以上，如果你的版本低于1.11版本，可以升级下版本。



1.安装go语言grpc包
================
> go get -u google.golang.org/grpc


2.安装protobuf编译器
================

用于将proto协议文件，编译成go语言代码


2.1. 安装protoc编译器
----------------
protobuf编译器就叫protoc，到下面github地址，根据自己的系统版本选择下载，解压缩安装即可。  
https://github.com/protocolbuffers/protobuf/releases

打开地址,ctrl+F查找linux-x86_64

例如3.14.0版本压缩包介绍：
```sh
protoc-3.14.0-win64.zip - windows 64版本
protoc-3.14.0-osx-x86_64.zip - mac os 64版本
protoc-3.14.0-linux-x86_64.zip - linux 64版本 (解压即可)
```

解压缩安装包之后，将 [安装目录]/bin 目录，添加到PATH环境变量。


2.2. 安装protoc编译器的go语言插件
----------------
因为目前的protoc编译器，默认没有包含go语言代码生成器，所以需要单独安装插件。

> go get -u github.com/golang/protobuf/protoc-gen-go

安装go语言插件后，需要将 $GOPATH/bin 路径加入到PATH环境变量中。



3.例子目录结构
================
本教程例子的目录结构如下：
```sh
helloworld/
├── client.go 	- 客户端代码
├── go.mod  	- go模块配置文件
├── proto     	- 协议目录
│   ├── helloworld.pb.go 	- rpc协议go版本代码
│   └── helloworld.proto 	- rpc协议文件
└── server.go  				- rpc服务端代码
```

初始化命令如下：
```sh
# 创建项目目录
mkdir helloworld

# 切换到项目目录
cd helloworld

# 创建RPC协议目录
mkdir proto

# 初始化go模块配置，用来管理第三方依赖, 本例子，项目模块名是：xxx.com/helloworld
go mod init xxx.com/helloworld
```
说明：本例子使用go module管理第三方依赖包，如果不了解可以点击学习go语言包管理



4.定义服务
================
定义服务，其实就是通过protobuf语法定义语言平台无关的接口。

文件: helloworld/proto/helloworld.proto
```golang
syntax = "proto3";
// 定义包名
package proto;

// 定义Greeter服务
service Greeter {
  // 定义SayHello方法，接受HelloRequest消息， 并返回HelloReply消息
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// 定义HelloRequest消息
message HelloRequest {
  // name字段
  string name = 1;
}

// 定义HelloReply消息
message HelloReply {
  // message字段
  string message = 1;
}
```
Greeter服务提供了一个SayHello接口，请求SayHello接口，需要传递一个包含name字段的HelloRequest消息，返回包含message字段的HelloReply消息。



5.编译proto协议文件
================

上面通过proto定义的接口，没法直接在代码中使用，因此需要通过protoc编译器，将proto协议文件，编译成go语言代码。
```sh
cd /xxx/helloworld/
# 切换到helloworld项目根目录，执行命令
protoc -I proto/ --go_out=plugins=grpc:proto proto/helloworld.proto

# protoc命令参数说明:

# -I 指定代码输出目录，忽略服务定义的包名，否则会根据包名创建目录

# --go_out 指定代码输出目录，格式：--go_out=plugins=grpc:目录名

# 命令最后面的参数是proto协议文件
```

编译成功后在proto目录生成了helloworld.pb.go文件，里面包含了，我们的服务和接口定义。



6.实现服务端代码
================
文件:helloworld/server.go
```golang
package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	// 导入grpc包
	"google.golang.org/grpc"
	// 导入刚才我们生成的代码所在的proto包。
        pb "xxx.com/helloworld/proto"
	"google.golang.org/grpc/reflection"
)


// 定义server，用来实现proto文件，里面实现的Greeter服务里面的接口
type server struct{}

// 实现SayHello接口
// 第一个参数是上下文参数，所有接口默认都要必填
// 第二个参数是我们定义的HelloRequest消息
// 返回值是我们定义的HelloReply消息，error返回值也是必须的。
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// 创建一个HelloReply消息，设置Message字段，然后直接返回。
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// 实例化grpc服务端
	s := grpc.NewServer()

    // 注册Greeter服务，注意：RegisterXxxxServer中的Xxxx即是你在.proto文件中定义的服务名
	pb.RegisterGreeterServer(s, &server{})

	// 往grpc服务端注册反射服务
	reflection.Register(s)

	// 监听127.0.0.1:50051地址
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

        // 启动grpc服务
	if err := s.Serve(lis); err != nil {
	    log.Fatalf("failed to serve: %v", err)
	}
}
```

运行服务端：
```sh
# 切换到项目根目录，运行命令
go run server.go
```


7.实现客户端代码
================
文件：helloworld/client.go
```golang
package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	// 导入grpc包
	"google.golang.org/grpc"
	// 导入刚才我们生成的代码所在的proto包。
    pb "xxx.com/helloworld/proto"
)

const (
	defaultName = "world"
)

func main() {
	// 连接grpc服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 延迟关闭连接
	defer conn.Close()


	// 初始化Greeter服务客户端，注意：NewXxxxClient中的Xxxx即是你在.proto文件中定义的服务名
	c := pb.NewGreeterClient(conn)

	// 初始化上下文，设置请求超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 延迟关闭请求会话
	defer cancel()

	// 调用SayHello接口，发送一条消息
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// 打印服务的返回的消息
	log.Printf("Greeting: %s", r.Message)
}
```


运行
```sh
# 切换到项目根目录，运行命令
go run client.go


#输出结果：
2019/09/26 00:19:27 Greeting: Hello world
```