package router

import (
	"github.com/gin-gonic/gin"
	"green/app/controller"
)

type RouterBook struct {

}

func (r *RouterBook) GroupName() string {
	return "book"
}

func (r *RouterBook) InitGroup(routerGroup *gin.RouterGroup) {
	cBook := &controller.ControllerBook{}
	routerGroup.GET("/list", cBook.List)
	routerGroup.GET("/info", cBook.Info)
}