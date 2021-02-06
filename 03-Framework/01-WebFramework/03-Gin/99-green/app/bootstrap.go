package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"green/app/router"
)

func Run() {
	app := gin.Default()

	// 路由初始化
	router.InitRouters(app)

	// start
	port := 8080
	app.Run(fmt.Sprintf(":%d", port))

	//certFile := ""
	//keyFile := ""
	//app.RunTLS(fmt.Sprintf(":%d", port), certFile, keyFile)
}