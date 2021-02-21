
gin框架路由拆分与注册



基本的路由注册
==================

下面最基础的gin路由注册方式，适用于路由条目比较少的简单项目或者项目demo。
```golang
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}

func main() {
	r := gin.Default()
	
	r.GET("/hello", helloHandler)

	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
```


路由拆分成单独文件或包
==================
当项目的规模增大后就不太适合继续在项目的main.go文件中去实现路由注册相关逻辑了，我们会倾向于把路由部分的代码都拆分出来，形成一个单独的文件或包

此时的目录结构：
```sh
gin_demo
├── go.mod
├── go.sum
├── main.go
└── routers.go
```

在routers.go文件中定义并注册路由信息：
```golang
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/hello", helloHandler)
	return r
}
```

main.go中调用上面定义好的setupRouter函数：
```golang
func main() {
	r := setupRouter()

	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
```



把路由部分的代码单独拆分成包的话也是可以的，拆分后的目录结构如下：
```sh
gin_demo
├── go.mod
├── go.sum
├── main.go
└── routers
    └── routers.go
```

routers/routers.go需要注意此时setupRouter需要改成首字母大写：
```golang
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}

// SetupRouter 配置路由信息
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/hello", helloHandler)
	return r
}
```

main.go文件内容如下：
```golang
package main

import (
	"fmt"
	"gin_demo/routers"
)

func main() {
	r := routers.SetupRouter()
	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
```



路由拆分成多个文件
==================
当我们的业务规模继续膨胀，单独的一个routers文件或包已经满足不了我们的需求了，
```golang
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/hello", helloHandler)
	r.GET("/xx1", xxHandler1)
	...
	r.GET("/xx30", xxHandler30)
	return r
}
```

因为我们把所有的路由注册都写在一个SetupRouter函数中的话就会太复杂了。


我们可以分开定义多个路由文件，例如：
```sh
gin_demo
├── go.mod
├── go.sum
├── main.go
└── routers
    ├── blog.go
    └── shop.go
```

routers/shop.go中添加一个LoadShop的函数，将shop相关的路由注册到指定的路由器：
```golang
func LoadShop(e *gin.Engine)  {
	e.GET("/hello", helloHandler)
	e.GET("/goods", goodsHandler)
	e.GET("/checkout", checkoutHandler)
	...
}
```

routers/blog.go中添加一个LoadBlog的函数，将blog相关的路由注册到指定的路由器：
```golang
func LoadBlog(e *gin.Engine) {
	e.GET("/post", postHandler)
	e.GET("/comment", commentHandler)
	...
}
```

在main函数中实现最终的注册逻辑如下：
```golang
func main() {
	r := gin.Default()
	routers.LoadBlog(r)
	routers.LoadShop(r)
	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
```



路由拆分到不同的APP
==================
有时候项目规模实在太大，那么我们就更倾向于把业务拆分的更详细一些，例如把不同的业务代码拆分成不同的APP。

因此我们在项目目录下单独定义一个app目录，用来存放我们不同业务线的代码文件，这样就很容易进行横向扩展。大致目录结构如下：
```sh
gin_demo
├── app
│   ├── blog
│   │   ├── handler.go
│   │   └── router.go
│   └── shop
│       ├── handler.go
│       └── router.go
├── go.mod
├── go.sum
├── main.go
└── routers
    └── routers.go
```

其中app/blog/router.go用来定义post相关路由信息，具体内容如下：
```golang
func Routers(e *gin.Engine) {
	e.GET("/post", postHandler)
	e.GET("/comment", commentHandler)
}
```

app/shop/router.go用来定义shop相关路由信息，具体内容如下：
```golang
func Routers(e *gin.Engine) {
	e.GET("/goods", goodsHandler)
	e.GET("/checkout", checkoutHandler)
}
```

routers/routers.go中根据需要定义Include函数用来注册子app中定义的路由，Init函数用来进行路由的初始化操作：
```golang
type Option func(*gin.Engine)

var options = []Option{}

// 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

// 初始化
func Init() *gin.Engine {
	r := gin.Default()
	for _, opt := range options {
		opt(r)
	}
	return r
}
```

main.go中按如下方式先注册子app中的路由，然后再进行路由的初始化：
```golang
func main() {
	// 加载多个APP的路由配置
	routers.Include(shop.Routers, blog.Routers)

	// 初始化路由
	r := routers.Init()
	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
```

