package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)
	router.POST("/user/:user_name", Login)

	return router
}

func main() {
	r := RegisterHandlers()
	// main阻塞在这里
	// listen -> RegisterHandlers -> handlers -> validation -> business -> response，每一个handlers都是在一个goroutin里执行的
	http.ListenAndServe(":8000", r)
}































