

Go的官网：https://golang.google.cn


1 下载压缩包
------------
访问 https://golang.google.cn/dl/, 选择一个自已需要的版本下载
```sh
wget https://golang.google.cn/dl/go1.14.14.linux-amd64.tar.gz

# 如果出现SSL证书错误问题，加上--no-check-certificate选项。
# wget --no-check-certificate https://golang.google.cn/dl/go1.14.14.linux-amd64.tar.gz 
```


2 解压到指定目录下
------------
```sh
tar -C /usr/local -xzvf go1.14.14.linux-amd64.tar.gz
# 解压成功后，会在/usr/local目录下生成go目录，亦即go的安装路径是/usr/local/go
```


3 设置Go的环境变量
------------
3.1 添加/usr/local/go/bin 到PATH环境变量
```sh
vim /etc/profile

# 在文件末尾添加如下内容: 
# go安装目录下的bin目录加入PATH中，必须配置
export PATH=$PATH:/usr/local/go/bin

# 你自已的go开发目录，一般与go安装目录区分开
export GOPATH=$HOME/work/gopath

#不使用代理，比如企业自已搭建的gitlab，非必配项
GOPRIVATE=*.gitlab.com,*.gitee.com,*.corp.qihoo.net
export GOPRIVATE

#代理,非必配项
GOPROXY=https://goproxy.io
export GOPROXY
```

使profile立即生效  
> source /etc/profile

你也可以将其添加到当前用户的配置文件，即$HOME/.profile，修改后使其生效的命令: source $HOME/.profile


3.2 查看Go的版本信息
> go version

3.3 查看Go的环境变量信息
> go env






使用
=============

开发时，常用两种方式

- 基于GOPATH目录开发
	配置GOPATH目录，code需放置在该目录下
	
- go mod
	代码无需放置GOPATH目录下



设置GOPATH
-------------
新建一个目录gopath作为GOPATH 的目录，并且设置环境变量（export GOPATH=/newhome/go/gopath）。
在gopath下新建3个文件夹分别为 src、pkg、bin目录。


go语言的工作空间其实就是一个文件目录，目录中必须包含src、pkg、bin三个目录。

其中src目录用于存放go源代码，pkg目录用于package对象，bin目录用于存放可执行对象。

GOPATH目录指明了你go代码的工作空间的位置，不能与GOROOT目录相同，而且GO代码必须位于工作空间内。



添加go代码库
-------------
src的源码代码可以go get github.com/** 的方式获取，也可以从复制别的地方项目到src目录下。

go get使用
```sh
#下载项目依赖
go get ./...

#拉取最新的版本(优先择取 tag)
go get golang.org/x/text@latest

#拉取 master 分支的最新 commit
go get golang.org/x/text@master

#拉取 tag 为 v0.3.2 的 commit
go get golang.org/x/text@v0.3.2

#拉取 hash 为 342b231 的 commit，最终会被转换为 v0.3.2：
go get golang.org/x/text@342b2e

#指定版本拉取，拉取v3版本
go get github.com/smartwalle/alipay/v3

#更新
go get -u xxxx
```



部署自己项目
-------------
上传自己的项目到src目录下与github.com, golang.org等其他目录平级
```sh
ls $GOPATH/src/
#输出
github.com
golang.org
goonlinetest
```

假设goonlinetest就是我的项目 然后进入我项目执行go build main.go 会编译一个linux 可执行程序。
最后执行 ./main 就行了。
如果想让项目在后台执行：执行 nohup ./main & ，这样就可以程序在后台运行了。

