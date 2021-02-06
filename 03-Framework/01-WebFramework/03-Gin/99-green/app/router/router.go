package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router interface {
	GroupName() string
	InitGroup(router *gin.RouterGroup)
}

func InitRouters(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome")
	})

	// 为路由注册各个控制器
	registerRouter(
		engine,
		new(RouterBook),
		new(RouterUser),
	);
}

func registerRouter(engine *gin.Engine, routers ...Router) {
	for _, router := range routers {
		routerGroup := engine.Group(router.GroupName())
		router.InitGroup(routerGroup)
	}
}