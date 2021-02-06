package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControllerBook struct {

}

func (c ControllerBook) List(ct *gin.Context) {
	ct.String(http.StatusOK, "book list")
}

func (c ControllerBook) Info(ct *gin.Context) {
	ct.String(http.StatusOK, "book info")
}
