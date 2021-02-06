package router

import (
	"github.com/gin-gonic/gin"
	"green/app/controller"
)

type RouterUser struct {

}

func (r *RouterUser) GroupName() string {
	return "user"
}

func (r *RouterUser) InitGroup(routerGroup *gin.RouterGroup) {
	cUser := &controller.ControllerUser{}
	routerGroup.GET("/list", cUser.List)
	routerGroup.GET("/info", cUser.Info)
}