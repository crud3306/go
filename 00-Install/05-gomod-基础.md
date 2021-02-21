

如何使用 Modules
----------------
- 把 golang 升级到 1.11+
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


使用 go module 管理依赖后，会在项目根目录下生成两个文件go.mod和go.sum。



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





GOPROXY
==============
Go1.11之后设置GOPROXY命令为：
```sh
export GOPROXY=https://goproxy.cn
```
Go1.13之后GOPROXY默认值为https://proxy.golang.org，在国内是无法访问的，所以十分建议大家设置GOPROXY，这里我推荐使用goproxy.cn。
```sh
go env -w GOPROXY=https://goproxy.cn,direct
```


go mod命令
==============
常用的go mod命令如下：
```sh
go mod download    下载依赖的module到本地cache（默认为$GOPATH/pkg/mod目录）
go mod edit        编辑go.mod文件
go mod graph       打印模块依赖图
go mod init        初始化当前文件夹, 创建go.mod文件
go mod tidy        增加缺少的module，删除无用的module
go mod vendor      将依赖复制到vendor下
go mod verify      校验依赖
go mod why         解释为什么需要依赖
```


go.mod
==============
go.mod文件记录了项目所有的依赖信息，其结构大致如下：
```sh
module github.com/Q1mi/studygo/blogger

go 1.12

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/satori/go.uuid v1.2.0
	google.golang.org/appengine v1.6.1 // indirect
)
```

其中
```sh
module用来定义包名
require用来定义依赖包及版本
indirect表示间接引用
```



依赖的版本
==============

go mod支持语义化版本号，比如go get foo@v1.2.3，也可以跟git的分支或tag，比如go get foo@master，当然也可以跟git提交哈希，比如go get foo@e3702bed2。关于依赖的版本支持以下几种格式：
```sh
gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
gopkg.in/vmihailenco/msgpack.v2 v2.9.1
gopkg.in/yaml.v2 <=v2.2.1
github.com/tatsushid/go-fastping v0.0.0-20160109021039-d7bb493dee3e
latest
```



replace
==============

在国内访问golang.org/x的各个包都需要翻墙，你可以在go.mod中使用replace替换成github上对应的库。

```sh
replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 => github.com/golang/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)
```


go get
==============
在项目中执行go get命令可以下载依赖包，并且还可以指定下载的版本。

```sh
运行go get -u将会升级到最新的次要版本或者修订版本(x.y.z, z是修订版本号， y是次要版本号)
运行go get -u=patch将会升级到最新的修订版本
运行go get package@version将会升级到指定的版本号version
```

如果下载所有依赖可以使用go mod download命令。



整理依赖
==============
我们在代码中删除依赖代码后，相关的依赖库并不会在go.mod文件中自动移除。这种情况下我们可以使用go mod tidy命令更新go.mod中的依赖关系。



go mod edit
==============
格式化

因为我们可以手动修改go.mod文件，所以有些时候需要格式化该文件。Go提供了一下命令：
```sh
go mod edit -fmt
```

添加依赖项
```sh
go mod edit -require=golang.org/x/text
```

移除依赖项
如果只是想修改go.mod文件中的内容，那么可以运行go mod edit -droprequire=package path，比如要在go.mod中移除golang.org/x/text包，可以使用如下命令：
```sh
go mod edit -droprequire=golang.org/x/text
```

关于go mod edit的更多用法可以通过go help mod edit查看。



在项目中使用go module
==============

既有项目
--------------
如果需要对一个已经存在的项目启用go module，可以按照以下步骤操作：

- 在项目目录下执行go mod init，生成一个go.mod文件。

- 执行go get，查找并记录当前项目的依赖，同时生成一个go.sum记录每个依赖库的版本和哈希值。


新项目
--------------
对于一个新创建的项目，我们可以在项目文件夹下按照以下步骤操作：

- 执行go mod init 项目名命令，在当前项目文件夹下创建一个go.mod文件。

- 手动编辑go.mod中的require依赖项或执行go get自动发现、维护依赖。