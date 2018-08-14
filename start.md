
$GOROOT 与 $GOPATH  
-------------
$GOROOT  
go的安装目录  
如：/usr/local/go  
  
$GOPATH  
go项目的工作目录  
如：~/go  
  
$GOROOT，$GOPATH 一般不是同一个目录  



安装依赖包
------------
> go get xxxx/xxx  
> go get -u xxxx/xxx  
如果想要强行更新代码包,可以在执行 go get 命令时加入 -u 标记。
  
比如：  
> go get gopkg.in/mgo.v2
> go get github.com/go-sql-driver/mysql
  
  


