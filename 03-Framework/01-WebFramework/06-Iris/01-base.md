


##文档地址
https://learnku.com/docs/iris-go  




###路由

####基本用法
基本用法 app.Handle
所有的 HTTP 方法都支持，开发者可以为相同路径的处理程序注册不同的方法。
第一个参数是 HTTP 方法，第二个参数是路由路径，第三个可变参数应该包含一个或者多个 iris.Handler，当获取某个资源的请求从服务器到来时，处理器按照注册顺序被调用执行。
```golang
app := iris.New()

app.Handle("GET", "/contact", func(ctx iris.Context) {
    ctx.HTML("<h1> Hello from /contact </h1>")
})
```

基本用法 app.各种Method
为了让开发者更容易开发处理程序，iris 提供了所有的 HTTP 方法。第一个参数是路由路径，第二个可变参数是一个或者多个 iris.Handler。
```golang
app := iris.New()

// 方法: "GET"
app.Get("/", handler)

// 方法: "POST"
app.Post("/", handler)

// 方法: "PUT"
app.Put("/", handler)

// 方法: "DELETE"
app.Delete("/", handler)

// 方法: "OPTIONS"
app.Options("/", handler)

// 方法: "TRACE"
app.Trace("/", handler)

// 方法: "CONNECT"
app.Connect("/", handler)

// 方法: "HEAD"
app.Head("/", handler)

// 方法: "PATCH"
app.Patch("/", handler)

// 用于所有 HTTP 方法
app.Any("/", handler)

func handler(ctx iris.Context){
    ctx.Writef("Hello from method: %s and path: %s", ctx.Method(), ctx.Path())
}
```

####路由分组 Party
一组路由可以用前缀路径分组，组之间共享相同的中间件和模板布局，组内可以嵌套组。
```golang
app := iris.New()

users := app.Party("/users", myAuthMiddlewareHandler)

// http://localhost:8080/users/42/profile
users.Get("/{id:int}/profile", userProfileHandler)
// http://localhost:8080/users/messages/1
users.Get("/messages/{id:int}", userMessageHandler)
```

```golang
app := iris.New()

app.PartyFunc("/users", func(users iris.Party) {
    users.Use(myAuthMiddlewareHandler)

    // http://localhost:8080/users/42/profile
    users.Get("/{id:int}/profile", userProfileHandler)
    // http://localhost:8080/users/messages/1
    users.Get("/messages/{id:int}", userMessageHandler)
})
```



###404，500错误处理
当程序产生了一个特定的 http 错误的时候，你可以定义你自己的错误处理代码。

错误代码是大于或者等于 400 的 http 状态码，像 404 not found 或者 500 服务器内部错误。
```golang
package main

import "github.com/kataras/iris"

func main(){
    app := iris.New()
    app.OnErrorCode(iris.StatusNotFound, notFound)
    app.OnErrorCode(iris.StatusInternalServerError, internalServerError)

    // 为所有的大于等于400的状态码注册一个处理器：
    // app.OnAnyErrorCode(handler)

    app.Get("/", index)
    app.Run(iris.Addr(":8080"))
}

func notFound(ctx iris.Context) {
   // 出现 404 的时候，就跳转到 $views_dir/errors/404.html 模板
    ctx.View("errors/404.html")
}

func internalServerError(ctx iris.Context) {
    ctx.WriteString("Oups something went wrong, try again")
}


func index(ctx context.Context) {
    ctx.View("index.html")
}
```





iris demo
```sh
config
apps
	controller
	service
	dao
library
	db
	redis
	es
	utils
		
route
main.go
```



iris demo
```sh
bin
conf
data
static
src
	service
	logic
	dao
lib
	mysql
	redis
	es
	result
		errors
		result
	utils
		
route
main.go
```



main 中调用new app， rotue，app run
route中关联controller
controller中接收数据，调用service，返回数据






sql条件
```sh

getList(cons, sets, limit, order)
getInfo(cons)


insert(sets)
update(cons, sets)
delete(cons)


cons拼装
{
	'key':'value'
	'id':1,
	'id': {'operator':'in', 'value':[]}
	'id': {'operator':'>=', 'value':1}
	'id': {'operator':'<=', 'value':1}
	'id': {'operator':'=', 'value':1}
}

cons = m


select * from tablexx 
where xxx 
order by xxx 
limit xx,xx

update tablex 
set xxxx
where xxx

insert into tablexx(xxx)
values(xxx)

delete from tablex 
where xxx
```

```

BuildWhere

从cons中提取key，Operator 拼装where语名
```








```
json
http
mysql
redis
```




modules