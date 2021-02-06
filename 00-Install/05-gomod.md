

如何使用 Modules
----------------
- 把 golang 升级到 1.11（现在1.12 已经发布了，建议使用1.12）
- 设置 GO111MODULE


go version 查看已安装的go版本
```sh
go version
#输出
#go version go1.14.5 linux/amd64
```


GO111MODULE
----------------
GO111MODULE 有三个值：off, on和auto（默认值）。

1 GO111MODULE=off  
go命令行将不会支持module功能，寻找依赖包的方式将会沿用旧版本那种通过vendor目录或者GOPATH模式来查找。

2 GO111MODULE=on  
go命令行会使用modules，而一点也不会去GOPATH目录下查找。

3 GO111MODULE=auto  
默认值，go命令行将会根据当前目录来决定是否启用module功能。这种情况下可以分为两种情形：

- 当前目录在GOPATH/src之外且该目录包含go.mod文件 
- 当前文件在包含go.mod文件的目录下面。 


当modules 功能启用时，依赖包的存放位置变更为$GOPATH/pkg，允许同一个package多个版本并存，且多个项目可以共享缓存的 module。




go mod 命令
---------------
```sh
#initialize new module in current directory（在当前目录初始化mod）
go mod init xxxx

#download modules to local cache(下载依赖包)
go mod download

#add missing and remove unused modules(拉取缺少的模块，移除不用的模块)
go mod tidy

#make vendored copy of dependencies(将依赖复制到项目下的vendor中)
go mod vendor

#edit go.mod from tools or scripts（编辑go.mod
go mod edit

#print module requirement graph (打印模块依赖图)
go mod graph

#verify dependencies have expected content (验证依赖是否正确）
go mod verify

#explain why packages or modules are needed(解释为什么需要依赖)
go mod why
```



查看go mod 的各命令的用法
---------------
```sh
go help mod xxxx命令

#例如
#go help mod init
#go help mod tidy

# go help mod init 的输出如下，主要看第一行usage：
usage: go mod init [module]
```



示例一：创建一个新项目
===========
在GOPATH 目录之外新建一个目录，并使用go mod init 初始化生成go.mod 文件
```sh
mkdir hello
cd hello

#执行go mod init xxx，初始化
go mod init hello
go: creating new go.mod: module hello

#查看当前目录，可发现多一个go.mod文件
ls
go.mod

#查看go.mod文件内容
cat go.mod
module hello

go 1.14
```


注意：  
如果当前项目中有其它文件，则import文件时，需带上mod名称做为前缀。举例：
```sh
hello/
├── api
│   └── demo.go
├── go.mod
└── main.go

#如果上这种目录结构，假设mod名为：hello， api包名：api
#在main.go引中api下的文件，需要用：import hello/api
```

