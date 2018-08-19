  
官方地址：
---------
https://beego.me/quickstart （快速开始）  
https://beego.me/docs/intro/   (文档很详细，简单快速入门)  
  
实例地址：
---------
https://www.cnblogs.com/zhangboyu/p/7760693.html  
  
  
执行逻辑
---------
main监听 -> 路由 -> 参数过滤 -> controller -> model -> db -> model -> controller -> ...  
  
  
安装  
---------
你需要安装或者升级 Beego 和 Bee 的开发工具:  
> go get -u github.com/astaxie/beego  
> go get -u github.com/beego/bee  
  
注意：  
如果想用整个beego框架，则安装上面beego 和 bee包  
如果仅想使用beego的一些子包，可以按需安装，比如：  
配置文件解析  
> go get -u github.com/astaxie/beego/config  
orm映射    
> go get -u github.com/astaxie/beego/orm   

> go get -u github.com/astaxie/beego/context    
  
  
  
设置$GOPATH  
---------
vi ~/.bash_profile 添加  
> export GOPATH="$HOME/go  
> export PATH="$GOPATH/bin:$PATH  
:wq 保存  
  
注意：不一定是.bash_profile文件，以你的机器实际情况为准  
  
  
  
创建项目  
-------------  
> bee new 项目名  //创建一个web应用    
> bee api 项目名  //创建一个api应用    
  
  
运行项目   
-------------  
cd进入对应的项目目录后，bee run 项目名   
    
> cd $GOPATH/src    
> bee new hello //创建一个web应用   
> bee api hello //创建一个api应用  
   
> cd hello  
> bee run hello  
  
  
项目打包发布:  
-------------  
> bee pack  
  
  
代码生成  
-------------  
生成models：  
> bee generate model user -fields="name:string,age:int"  
  
生成controller:  
> bee generate controller user  
  
生成view:  
> bee generate view user  
  
生成文档:  
> bee generate docs  
  
  
  
应用配置  
------------  
编辑 conf/app.conf  
  
> appname = WEB  
> runmode = dev  
>   
> [dev]  
> httpport = 8080  
> [test]  
> httpport = 8081  
> [prod]  
> httpport = 8082  
  
  
在代码中要获取或修改配置  
------------
beego.BConfig.Xxx 来访问和修改系统配置   
例：  
beego.BConfig.RunMode  
   
beego.AppConfig.String("xxx") 来访问和修改应用配置  
例：  
beego.AppConfig.String("dev::mysqluser")  


入口文件 main 
------------
```go
package main

import (
    _ "quickstart/routers"
    "github.com/astaxie/beego"
)

func main() {
    beego.Run()
}
```
    
   
router
------------
```go
package routers

import (
    "quickstart/controllers"
    "github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
```
  
  
controller  
------------
```go
package controllers

import (
    "github.com/astaxie/beego"
)

type MainController struct {
    beego.Controller
}

func (this *MainController) Get() {
    this.Data["Website"] = "beego.me"
    this.Data["Email"] = "astaxie@gmail.com"
    this.TplName = "index.tpl"  // 指定模板后，系统自动调用Render函数去渲染
}
func (this *MainController) Get() {
	id = this.GetString("id")

    this.Ctx.WriteString(id)
}
```
  
接收get参数
---------
> c.GetString(key string) string  
> c.GetStrings(key string) []string  
> c.GetInt(key string) (int64, error)  
> c.GetBool(key string) (bool, error)  
> c.GetFloat(key string) (float64, error)  
  
或者  
> c.Input().Get(key string)  
  
  
  
获取/设置cookie  
---------
> c.Ctx.SetCookie(key string, value string, expire int, path string)  
> c.Ctx.GetCookie(key string)  
  
例：  
设置  
> c.Ctx.SetCookie("name", "zhangsan", 100, "/")  
获取  
> c.Ctx.GetCookie("name")  
清除  
> c.Ctx.SetCookie("name", "zhangsan", -1, "/")  
  
  
使用session  
---------
1 在main入口函数中设置如下，开启session，这样在controller中才能使用session  
> beego.BConfig.WebConfig.Session.SessionOn = true   
或者在配置文件中配置   
> sessionon = true   
  
2  使用  
> SetSession(name string, value interface{})  
> GetSession(name string) interface{}  
> DelSession(name string)   
> SessionRegenerateID()  
> DestroySession()  
  
  
  
查看地址：  
https://beego.me/docs/mvc/controller/session.md  
  
  
  
单独使用beego中的一些包  
--------------
config库（配置文件解析)   
> go get -u github.com/astaxie/beego/config  
  
```go
iniconf, err := config.NewConfig("ini", "testini.conf")
if err != nil {
	panic(err)
}
  
iniconf.String("appname")  
```
解析器对象目前支持的函数有：  
> Set(key, val string) error  
> String(key string) string  
> Int(key string)(int, error)  
> Int64(key string)(int64, error)  
> Bool(key string)(bool, error)  
> Float(key string)(float64, error)  
> DIY(key string)(interface{}, error)  
  
  
  
orm映射  
> go get -u github.com/astaxie/beego/orm    
  
  
httplib库  
主要用来模拟客户端发送http请求，类似于curl工具，支持链式操作。  
> go get -u github.com/astaxie/beego/httplib    
```go
req := httplib.Get("http://beego.me/")
str, err := req.String()
if err != nil {
	panic(err)
}
fmt.Println(str)
```
  
  
爬虫抓动态网页，可以用：  
plantomjs  
selenium  
  

  
   




  


  
  



































