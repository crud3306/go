package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControllerUser struct {

}

func (c ControllerUser) List(ct *gin.Context) {
	ct.String(http.StatusOK, "user list")
}

func (c ControllerUser) Info(ct *gin.Context) {
	ct.String(http.StatusOK, "user info")
}
