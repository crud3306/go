
安装  
-------------
下载地址：  
https://studygolang.com/dl   
或：  
https://golang.google.cn/dl/    
  

mac下两种安装方式：
=============
1 下载压缩包(xxx.tar.gz)，解压即可 (简单快速)  
2 下载安装包(xxx.pkg)，点击安装  
  
配置
------------
Download the archive and extract it into /usr/local, creating a Go tree in /usr/local/go. For example:
```
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
注：$VERSION.$OS-$ARCH 改成你自已下载的版本号
```
Add /usr/local/go/bin to the PATH environment variable. You can do this by adding this line to your /etc/profile (for a system-wide installation) or $HOME/.profile:
```
vi ~/.profile
export PATH=$PATH:/usr/local/go/bin

source ~/.profile
```

创建你自已的工作目录
> cd ~  
> mkdir go/src  
开始第一个脚本  
> vi hello.go  
```
package main

import "fmt"

func main() {
	fmt.Printf("hello, world\n")
}
```
执行  
> go run hello.go  


  
$GOROOT 与 $GOPATH  
-------------
$GOROOT  
go的安装目录  
如：/usr/local/go  
  
$GOPATH  
go项目的工作目录  
如：~/go  
  
$GOROOT，$GOPATH 一般不是同一个目录  

