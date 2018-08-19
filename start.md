
安装
-------------
go的安装请见install目录：
https://github.com/crud3306/go-start/tree/master/install   


路径 $GOROOT 与 $GOPATH  
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
  
      
执行程序  
------------ 
go run xxx.go   
go run：go run 编译并直接运行程序，它会产生一个临时文件（但不会生成 .exe 文件），直接在命令行输出程序执行结果，方便用户调试。  

  
  
go install/build 都是用来编译包和其依赖的包
------------  

go build：用于测试编译包，主要检查是否会有编译错误，如果是一个可执行文件的源码（即是 main 包），就会直接生成一个可执行文件。

go install：的作用有两步：第一步是编译导入的包文件，所有导入的包文件编译完才会编译主程序；第二步是将编译后生成的可执行文件放到 bin 目录下（$GOPATH/bin），编译后的包文件放到 pkg 目录下（$GOPATH/pkg）。


   
    

------------



   



