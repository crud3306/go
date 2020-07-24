
注：UNIX/Linux/Mac OS X, 和 FreeBSD 均可采用如下安装方式

  
安装
------------
下载二进制包(xxx.tar.gz)，解压即可。(简单快速)    
如：go1.14.4.linux-amd64.tar.gz  

下载地址：  
https://studygolang.com/dl   
或：  
https://golang.google.cn/dl/    

```sh
cd /usr/local/src/
wget https://studygolang.com/dl/golang/go1.14.4.linux-amd64.tar.gz
ll
tar -xzf go1.14.4.linux-amd64.tar.gz -C /usr/local
```

设置环境变量
------------
```sh
#创建工作目录，名字随意
mkdir ~/gopath

#设置环境变量
vi ~/.bash_profile
export GOROOT=/usr/local/go   #（你解压后的目录，即安装目录）
export PATH=$GOROOT/bin:$PATH
export GOPATH=/xxx/gopath  #（你的开发地址，这个随便，你自己设置，假设~/gopath）

#刷新环境变量
source ~/.profile
```


验证是否安装成功
------------
> go version


开始第一个脚本  
------------
> vi hello.go  
```go
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
$GOROOT，$GOPATH 一般不是同一个目录  

$GOROOT：go的安装目录
```sh
  
如：/usr/local/go  
```

$GOPATH：go项目的工作目录
```sh
如：~/go  
```



