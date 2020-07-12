

Go 是一个开源的编程语言，它能让构造简单、可靠且高效的软件变得容易。


Go 语言特色
-------------
- 简洁、快速、安全  
- 并行、有趣、开源  
- 内存管理、数组安全、编译迅速  


Go 语言用途
-------------
Go 语言被设计成一门应用于搭载 Web 服务器，存储集群或类似用途的巨型中央服务器的系统编程语言。

对于高性能分布式系统领域而言，Go语言无疑比大多数其它语言有着更高的开发效率。它提供了海量并行的支持，这对于游戏服务端的开发而言是再好不过了。


安装
-------------
Go的安装请见00-Install目录：
https://github.com/crud3306/go/tree/master/00-Install   


路径 $GOROOT 与 $GOPATH  
-------------
$GOROOT，$GOPATH 一般不是同一个目录  

$GOROOT  
```sh
go的安装目录  
如：/usr/local/go  
```

$GOPATH  
```sh
go项目的工作目录  
如：~/go  
```


安装依赖包  
------------
> go get xxxx/xxx  
> go get -u xxxx/xxx    
如果想要强行更新代码包，可以在执行 go get 命令时加入 -u 标记。  
   
比如：   
> go get gopkg.in/mgo.v2  
> go get github.com/go-sql-driver/mysql  
  
      

执行程序  
------------ 
go run xxx.go   
go run：go run 编译并直接运行程序，它会产生一个临时文件（但不会生成 .exe 文件），直接在命令行输出程序执行结果，方便用户调试。  

  
  
go install/build 都是用来编译包和其依赖的包
------------  

go build：用于测试编译包，主要检查是否会有编译错误，如果是一个可执行文件的源码（即是 main 包），就会直接生成一个可执行文件。

go install：的作用有两步：  
- 第一步是编译导入的包文件，所有导入的包文件编译完才会编译主程序；  
- 第二步是将编译后生成的可执行文件放到 bin 目录下（$GOPATH/bin），编译后的包文件放到 pkg 目录下（$GOPATH/pkg）。  


   

第一个 Go 程序
----------
接下来我们来编写第一个 Go 程序 hello.go（Go 语言源文件的扩展是 .go），代码如下：

hello.go 文件
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

执行以上代码输出:
```sh
$ go run hello.go 
Hello, World!
```


此外我们还可以使用 go build 命令来生成二进制文件：
```sh
$ go build hello.go 
$ ls
hello    hello.go
$ ./hello 
Hello, World!
```

   



